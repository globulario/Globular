package globule

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/globulario/services/golang/config"
	Utility "github.com/globulario/utility"
)

// bootstrapServiceConfigsFromDisk syncs service definition files into etcd
// without clobbering existing entries; only missing services are restored.
func (g *Globule) bootstrapServiceConfigsFromDisk() error {
	dir := config.GetServicesConfigDir()
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("bootstrap services: stat %s: %w", dir, err)
	}
	if !fi.IsDir() {
		return nil
	}

	existing := map[string]struct{}{}
	if svcs, err := config.GetServicesConfigurations(); err == nil {
		for _, svc := range svcs {
			id := strings.TrimSpace(Utility.ToString(svc["Id"]))
			if id != "" {
				existing[id] = struct{}{}
			}
		}
	} else {
		g.log.Warn("bootstrap services: list existing configs failed", "err", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("bootstrap services: readdir %s: %w", dir, err)
	}

	var loaded int
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".json") {
			continue
		}

		path := filepath.Join(dir, name)
		raw, err := os.ReadFile(path)
		if err != nil {
			g.log.Warn("bootstrap services: read config failed", "file", path, "err", err)
			continue
		}

		var desired map[string]interface{}
		if err := json.Unmarshal(raw, &desired); err != nil {
			g.log.Warn("bootstrap services: parse config failed", "file", path, "err", err)
			continue
		}

		id := strings.TrimSpace(Utility.ToString(desired["Id"]))
		if id == "" {
			base := strings.TrimSuffix(name, filepath.Ext(name))
			id = base
			desired["Id"] = id
		}

		if id == "" {
			g.log.Warn("bootstrap services: missing service id", "file", path)
			continue
		}

		if _, ok := existing[id]; ok {
			continue
		}

		if err := config.SaveServiceConfiguration(desired); err != nil {
			g.log.Warn("bootstrap services: save config failed", "id", id, "file", path, "err", err)
			continue
		}

		existing[id] = struct{}{}
		loaded++
	}

	if loaded > 0 {
		g.log.Info("bootstrap services: registered missing service configs", "count", loaded)
	}
	return nil
}

func (g *Globule) stopServicesEtcd(ctx context.Context) error {
	svcs, err := config.GetServicesConfigurations()
	if err != nil {
		return err
	}

	if ctx != nil {
		if cerr := ctx.Err(); cerr != nil {
			return cerr
		}
	}

	updated := 0
	for _, svc := range svcs {
		id := Utility.ToString(svc["Id"])
		if id == "" {
			continue
		}
		cur, err := config.GetServiceConfigurationById(id)
		if err != nil || cur == nil {
			continue
		}
		if Utility.ToString(cur["State"]) == "closing" {
			continue
		}
		cur["State"] = "closing"
		cur["LastError"] = ""
		if err := config.SaveServiceConfiguration(cur); err != nil {
			g.log.Warn("stopServicesEtcd: save desired closing failed", "id", id, "err", err)
			continue
		}
		updated++
	}

	g.log.Info("stopServicesEtcd: requested closing state", "services", updated)
	return nil
}
