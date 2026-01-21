package globular

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var (
	forbiddenPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)"tls/"\s*\+\s*domain`),
		regexp.MustCompile(`(?i)"pki/"\s*\+\s*domain`),
		regexp.MustCompile(`(?i)GetConfigDir\(\)\s*\+\s*"/tls/"\s*\+\s*domain`),
		regexp.MustCompile(`(?i)GetConfigDir\(\)\s*\+\s*"/pki/"\s*\+\s*domain`),
		regexp.MustCompile(`(?i)tls/\{\{\s*\.\s*Domain\s*\}\}`),
		regexp.MustCompile(`(?i)pki/\{\{\s*\.\s*Domain\s*\}\}`),
		regexp.MustCompile(`(?i)filepath\.Join\([^)]*tlsDir[^)]*domain[^)]*\)`),
		regexp.MustCompile(`(?i)filepath\.Join\([^)]*GetRuntimeTLSDir\(\)[^)]*domain[^)]*\)`),
		regexp.MustCompile(`(?i)filepath\.Join\([^)]*,\s*"(tls|pki)"\s*,[^)]*(domain|host)[^)]*\)`),
	}
	selfPath        string
	defaultSkipDirs = map[string]struct{}{
		".git":         {},
		"vendor":       {},
		"third_party":  {},
		"generated":    {},
		"node_modules": {},
		"dist":         {},
		"build":        {},
		"out":          {},
	}
)

func init() {
	if _, callerFile, _, ok := runtime.Caller(0); ok {
		selfPath = callerFile
	}
}

func scanPaths(paths []string) error {
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %s: %w", path, err)
		}
		for _, re := range forbiddenPatterns {
			loc := re.FindIndex(data)
			if loc == nil {
				continue
			}
			line := 1 + bytes.Count(data[:loc[0]], []byte("\n"))
			return fmt.Errorf("forbidden pattern matched in %s:%d (%s)", path, line, re.String())
		}
	}
	return nil
}

func scanRepo(root string) error {
	var files []string
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if _, skip := defaultSkipDirs[d.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		if path == selfPath {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return err
	}
	return scanPaths(files)
}

func findRepoRoot() string {
	_, callerFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(callerFile)
	for {
		for _, marker := range []string{"go.work", ".git", "go.mod"} {
			if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

func TestTLSScanRejectsCredsDomain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.go")
	joinExpr := "filepath.Join(%s, %s)"
	content := fmt.Sprintf(`package x
import "path/filepath"
func f(tlsDir string, credsDomain string) string { return `+joinExpr+` }`, "tlsDir", "credsDomain")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	err := scanPaths([]string{path})
	if err == nil {
		t.Fatalf("expected forbidden pattern to be detected in %s", path)
	}
	if !strings.Contains(err.Error(), path) {
		t.Fatalf("error does not mention offending file: %v", err)
	}
}

func TestTLSScanRejectsConcatDomain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "concat.go")
	content := `package x
func f(domain string) string { return "tls/" + domain }`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if err := scanPaths([]string{path}); err == nil {
		t.Fatalf("expected concat pattern to be detected")
	}
}

func TestTLSScanRepo(t *testing.T) {
	root := findRepoRoot()
	if err := scanRepo(root); err != nil {
		t.Fatalf("tls scan failed: %v", err)
	}
}
