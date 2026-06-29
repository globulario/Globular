package faststart

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOptimizeMissingFFmpegLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "video.mp4")
	original := []byte("original media payload")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatalf("write original: %v", err)
	}

	t.Setenv("PATH", dir)
	Optimize(path)

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read optimized file: %v", err)
	}
	if string(got) != string(original) {
		t.Fatalf("file changed without ffmpeg: got %q want %q", got, original)
	}
}

func TestOptimizeSuccessfulRunReplacesOriginal(t *testing.T) {
	dir := t.TempDir()
	ffmpegPath := filepath.Join(dir, "ffmpeg")
	script := `#!/bin/sh
in=""
out=""
prev=""
for arg in "$@"; do
  if [ "$prev" = "-i" ]; then in="$arg"; fi
  prev="$arg"
  out="$arg"
done
/bin/cp "$in" "$out"
/usr/bin/python3 - <<'EOF' >> "$out"
print('x' * 5000, end='')
EOF
`
	if err := os.WriteFile(ffmpegPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake ffmpeg: %v", err)
	}

	path := filepath.Join(dir, "video.mp4")
	original := []byte("original media payload")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatalf("write original: %v", err)
	}

	t.Setenv("PATH", dir)
	Optimize(path)

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read optimized file: %v", err)
	}
	if len(got) <= len(original) {
		t.Fatalf("optimized file size=%d want > %d", len(got), len(original))
	}
	if string(got[:len(original)]) != string(original) {
		t.Fatalf("optimized file prefix=%q want %q", got[:len(original)], original)
	}
	if _, err := os.Stat(path + ".faststart.tmp"); !os.IsNotExist(err) {
		t.Fatalf("temp file still present: %v", err)
	}
}

func TestOptimizeSmallOutputLeavesOriginalUntouched(t *testing.T) {
	dir := t.TempDir()
	ffmpegPath := filepath.Join(dir, "ffmpeg")
	script := `#!/bin/sh
out=""
for arg in "$@"; do
  out="$arg"
done
printf 'tiny' > "$out"
`
	if err := os.WriteFile(ffmpegPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake ffmpeg: %v", err)
	}

	path := filepath.Join(dir, "video.mp4")
	original := []byte("original media payload")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatalf("write original: %v", err)
	}

	t.Setenv("PATH", dir)
	Optimize(path)

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file after small output: %v", err)
	}
	if string(got) != string(original) {
		t.Fatalf("file changed after small output: got %q want %q", got, original)
	}
	if _, err := os.Stat(path + ".faststart.tmp"); !os.IsNotExist(err) {
		t.Fatalf("temp file still present after discard: %v", err)
	}
}
