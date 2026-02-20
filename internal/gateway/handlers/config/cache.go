package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	config_ "github.com/globulario/services/golang/config"
	Utility "github.com/globulario/utility"
)

type ServiceConfigCache struct {
	mu      sync.RWMutex
	entries map[string]map[string]any
}

func NewServiceConfigCache() *ServiceConfigCache {
	c := &ServiceConfigCache{entries: make(map[string]map[string]any)}
	c.reload()
	return c
}

func (c *ServiceConfigCache) reload() {
	dir := config_.GetServicesConfigDir()
	if dir == "" {
		return
	}

	entries := make(map[string]map[string]any)
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var parsed map[string]any
		if err := json.Unmarshal(raw, &parsed); err != nil {
			continue
		}
		id := Utility.ToString(parsed["Id"])
		if id == "" {
			id = strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		}
		for _, key := range []string{id, strings.ToLower(id)} {
			if key == "" {
				continue
			}
			if _, ok := entries[key]; ok {
				continue
			}
			entries[key] = parsed
		}
		if name, _ := parsed["Name"].(string); name != "" {
			entries[strings.TrimSpace(strings.ToLower(name))] = parsed
		}
	}

	c.mu.Lock()
	c.entries = entries
	c.mu.Unlock()
}
func (c *ServiceConfigCache) Get(idOrName string) (map[string]any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[strings.ToLower(strings.TrimSpace(idOrName))]
	if !ok || entry == nil {
		return nil, false
	}
	return cloneMap(entry), true
}

func cloneMap(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
