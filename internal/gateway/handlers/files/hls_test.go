package files

import "testing"

// TestIsHLSFile tests the HLS file detection function
func TestIsHLSFile(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
		desc     string
	}{
		// HLS segment files
		{"/videos/segment.ts", true, "MPEG-TS segment"},
		{"/path/to/video.ts", true, "TS file"},

		// HLS playlist files (various resolutions)
		{"/videos/240p.m3u8", true, "240p playlist"},
		{"/videos/360p.m3u8", true, "360p playlist"},
		{"/videos/480p.m3u8", true, "480p playlist"},
		{"/videos/720p.m3u8", true, "720p playlist"},
		{"/videos/1080p.m3u8", true, "1080p playlist"},
		{"/videos/2160p.m3u8", true, "4K playlist"},

		// Non-HLS files
		{"/videos/video.mp4", false, "MP4 video"},
		{"/videos/video.mkv", false, "MKV video"},
		{"/documents/file.txt", false, "Text file"},
		{"/documents/file.pdf", false, "PDF file"},
		{"/images/photo.jpg", false, "JPEG image"},
		{"/scripts/test.js", false, "JavaScript file"},

		// Edge cases
		{"", false, "Empty path"},
		{"/", false, "Root path"},
		{"/videos/test", false, "No extension"},
		{"/videos/test.m3u8", false, "Generic m3u8 (not resolution-specific)"},
		{"/videos/master.m3u8", false, "Master playlist"},
		{"/videos/test.tsv", false, "TSV file (not .ts)"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := isHLSFile(tc.path)
			if result != tc.expected {
				t.Errorf("isHLSFile(%q) = %v, expected %v", tc.path, result, tc.expected)
			}
		})
	}
}

// TestIsHLSFile_CaseSensitive tests that HLS detection is case-sensitive
func TestIsHLSFile_CaseSensitive(t *testing.T) {
	// These should NOT match (case-sensitive)
	testCases := []string{
		"/videos/segment.TS", // Uppercase .TS
		"/videos/240P.m3u8",  // Uppercase P
		"/videos/720p.M3U8",  // Uppercase M3U8
	}

	for _, path := range testCases {
		if isHLSFile(path) {
			t.Errorf("isHLSFile(%q) should be false (case-sensitive)", path)
		}
	}
}

// TestIsHLSFile_Comprehensive tests real-world HLS streaming patterns
func TestIsHLSFile_Comprehensive(t *testing.T) {
	// Real-world HLS paths that should be allowed
	hlsPaths := []string{
		"/users/alice/.hidden/stream/240p.m3u8",
		"/users/bob/videos/movie/segment-001.ts",
		"/webroot/content/live/720p.m3u8",
		"/applications/video-player/assets/1080p.m3u8",
	}

	for _, path := range hlsPaths {
		if !isHLSFile(path) {
			t.Errorf("isHLSFile(%q) should be true for valid HLS path", path)
		}
	}

	// Paths that should NOT be HLS files
	nonHLSPaths := []string{
		"/users/alice/document.pdf",
		"/users/bob/videos/movie.mp4",
		"/webroot/styles/main.css",
		"/applications/app/index.html",
	}

	for _, path := range nonHLSPaths {
		if isHLSFile(path) {
			t.Errorf("isHLSFile(%q) should be false for non-HLS path", path)
		}
	}
}

// BenchmarkIsHLSFile benchmarks the HLS detection performance
func BenchmarkIsHLSFile(b *testing.B) {
	testPaths := []string{
		"/videos/segment.ts",
		"/videos/720p.m3u8",
		"/videos/video.mp4",
		"/documents/file.txt",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			_ = isHLSFile(path)
		}
	}
}
