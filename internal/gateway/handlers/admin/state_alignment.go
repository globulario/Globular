package admin

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/installed_state"
	"github.com/globulario/services/golang/repository/repository_client"
	"github.com/globulario/services/golang/versionutil"
)

// stateAlignmentStatus describes the state alignment of a single package
// across the artifact and installed-observed layers.
type stateAlignmentStatus struct {
	Name              string `json:"name"`
	Kind              string `json:"kind"`
	Status            string `json:"status"`
	InstalledVersion  string `json:"installed_version,omitempty"`
	InstalledBuildNum int64  `json:"installed_build_number,omitempty"`
	RepoVersion       string `json:"repo_version,omitempty"`
	RepoBuildNum      int64  `json:"repo_build_number,omitempty"`
	Message           string `json:"message,omitempty"`
}

// stateAlignmentReport is the result of a read-only alignment check.
type stateAlignmentReport struct {
	Packages       []*stateAlignmentStatus `json:"packages"`
	Aligned        int                     `json:"aligned"`
	Drifted        int                     `json:"drifted"`
	Unmanaged      int                     `json:"unmanaged"`
	MissingInRepo  int                     `json:"missing_in_repo"`
	RepositoryAddr string                  `json:"repository_addr,omitempty"`
}

// NewStateAlignmentHandler returns an HTTP handler that serves
// GET /admin/state-alignment — a read-only artifact+installed alignment report.
//
// Reads from:
//   - installed-state registry (etcd) — canonical observed state
//   - repository service (gRPC, auto-discovered) — artifact state
func NewStateAlignmentHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Auto-discover repository address.
		repoAddr := config.ResolveServiceAddr("repository.PackageRepository", "")

		report := &stateAlignmentReport{
			RepositoryAddr: repoAddr,
		}

		// Step 1: Collect all installed packages from the registry.
		type pkgInfo struct {
			name        string
			version     string
			buildNumber int64
			kind        string
		}
		installed := make(map[string]pkgInfo)

		for _, kind := range []string{"SERVICE", "APPLICATION", "INFRASTRUCTURE", "COMMAND"} {
			pkgs, err := installed_state.ListAllNodes(r.Context(), kind, "")
			if err != nil {
				continue
			}
			for _, pkg := range pkgs {
				name := strings.TrimSpace(pkg.GetName())
				if name == "" {
					continue
				}
				key := kind + "/" + name
				if _, exists := installed[key]; !exists {
					ver := pkg.GetVersion()
					if cv, err := versionutil.Canonical(ver); err == nil {
						ver = cv
					}
					installed[key] = pkgInfo{name: name, version: ver, buildNumber: pkg.GetBuildNumber(), kind: kind}
				}
			}
		}

		// Step 2: Query repository for available artifact versions.
		type repoInfo struct {
			version     string
			buildNumber int64
		}
		repoVersions := make(map[string]repoInfo)
		if repoAddr != "" {
			rc, err := repository_client.NewRepositoryService_Client(repoAddr, "repository.PackageRepository")
			if err == nil {
				if arts, err := rc.ListArtifacts(); err == nil {
					for _, m := range arts {
						if m.GetRef() == nil {
							continue
						}
						name := m.GetRef().GetName()
						ver := m.GetRef().GetVersion()
						if cv, err := versionutil.Canonical(ver); err == nil {
							ver = cv
						}
						key := strings.ToLower(name)
						existing, ok := repoVersions[key]
						if !ok {
							repoVersions[key] = repoInfo{version: ver, buildNumber: m.GetBuildNumber()}
						} else if cmp, cerr := versionutil.CompareFull(existing.version, existing.buildNumber, ver, m.GetBuildNumber()); cerr == nil && cmp < 0 {
							repoVersions[key] = repoInfo{version: ver, buildNumber: m.GetBuildNumber()}
						}
					}
				}
				rc.Close()
			}
		}

		// Step 3: Cross-reference installed vs. repo.
		for _, pkg := range installed {
			repo := repoVersions[strings.ToLower(pkg.name)]

			entry := &stateAlignmentStatus{
				Name:              pkg.name,
				Kind:              pkg.kind,
				InstalledVersion:  pkg.version,
				InstalledBuildNum: pkg.buildNumber,
				RepoVersion:       repo.version,
				RepoBuildNum:      repo.buildNumber,
			}

			switch {
			case repo.version == "":
				entry.Status = "missing_in_repo"
				entry.Message = "installed but artifact not found in repository"
				report.MissingInRepo++
			case !versionutil.EqualFull(pkg.version, pkg.buildNumber, repo.version, repo.buildNumber):
				entry.Status = "drifted"
				entry.Message = "installed version differs from repository latest"
				report.Drifted++
			default:
				entry.Status = "aligned"
				report.Aligned++
			}

			report.Packages = append(report.Packages, entry)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	})
}
