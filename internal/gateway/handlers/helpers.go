package handlers

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func findHashedFile(p string) (string, error) {
	dir := filepath.Dir(p)
	base := filepath.Base(p)

	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fname := e.Name()
		if strings.HasPrefix(fname, name+".") && strings.HasSuffix(fname, ext) {
			return filepath.Join(dir, fname), nil
		}
	}
	return "", errors.New("hashed file not found for " + base)
}

func resolveImportPath(base, line string) (string, error) {
	re := regexp.MustCompile(`from\s+['"]([^'"]+)['"]`)
	m := re.FindStringSubmatch(line)
	if len(m) != 2 {
		return line, nil
	}
	importPath := m[1]

	if strings.HasPrefix(importPath, ".") {
		target := filepath.Join(base, importPath)
		if hashed, err := findHashedFile(target); err == nil {
			hashedRel := strings.TrimPrefix(hashed, base+string(filepath.Separator))
			return strings.Replace(line, importPath, "./"+hashedRel, 1), nil
		}
	}
	return line, nil
}

func streamHandlerMaybe(name string, w http.ResponseWriter, r *http.Request) bool {
	clean := filepath.Clean(name)
	if strings.Contains(clean, "..") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return true
	}

	f, err := os.Open(clean)
	if err != nil {
		return false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "stat "+clean+": "+err.Error(), http.StatusInternalServerError)
		return true
	}
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
	return true
}
