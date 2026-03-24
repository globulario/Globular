package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	repopb "github.com/globulario/services/golang/repository/repositorypb"
	"google.golang.org/protobuf/encoding/protojson"
)

// RepositoryProvider is the interface for querying the artifact repository.
type RepositoryProvider interface {
	SearchArtifacts(query, kind, publisher, platform, pageToken string, pageSize int32) ([]*repopb.ArtifactManifest, string, int32, error)
	GetArtifactManifest(ref *repopb.ArtifactRef, buildNumber int64) (*repopb.ArtifactManifest, error)
	GetArtifactVersions(publisherID, name, platform string) ([]*repopb.ArtifactManifest, error)
	DeleteArtifact(ref *repopb.ArtifactRef) error
}

// NewRepositorySearchHandler returns a handler for GET /admin/repository/search.
//
// Query params: q, kind, publisher, platform, page_token, page_size
func NewRepositorySearchHandler(prov RepositoryProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()
		query := strings.TrimSpace(q.Get("q"))
		kind := strings.TrimSpace(q.Get("kind"))
		publisher := strings.TrimSpace(q.Get("publisher"))
		platform := strings.TrimSpace(q.Get("platform"))
		pageToken := strings.TrimSpace(q.Get("page_token"))
		var pageSize int32 = 50
		if ps := strings.TrimSpace(q.Get("page_size")); ps != "" {
			if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 200 {
				pageSize = int32(v)
			}
		}

		artifacts, nextToken, total, err := prov.SearchArtifacts(query, kind, publisher, platform, pageToken, pageSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		items := make([]json.RawMessage, 0, len(artifacts))
		for _, a := range artifacts {
			data, err := protojson.Marshal(a)
			if err != nil {
				continue
			}
			items = append(items, json.RawMessage(data))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"artifacts":       items,
			"count":           len(items),
			"total":           total,
			"next_page_token": nextToken,
		})
	})
}

// NewRepositoryManifestHandler returns a handler for GET /admin/repository/manifest.
//
// Query params: publisher, name, version, platform
func NewRepositoryManifestHandler(prov RepositoryProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()
		ref := &repopb.ArtifactRef{
			PublisherId: strings.TrimSpace(q.Get("publisher")),
			Name:        strings.TrimSpace(q.Get("name")),
			Version:     strings.TrimSpace(q.Get("version")),
			Platform:    strings.TrimSpace(q.Get("platform")),
		}
		if ref.Name == "" {
			http.Error(w, "name parameter required", http.StatusBadRequest)
			return
		}
		if ref.Version == "" {
			http.Error(w, "version parameter required", http.StatusBadRequest)
			return
		}

		manifest, err := prov.GetArtifactManifest(ref, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := protojson.Marshal(manifest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
}

// NewRepositoryVersionsHandler returns a handler for GET /admin/repository/versions.
//
// Query params: publisher, name, platform
func NewRepositoryVersionsHandler(prov RepositoryProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()
		publisher := strings.TrimSpace(q.Get("publisher"))
		name := strings.TrimSpace(q.Get("name"))
		platform := strings.TrimSpace(q.Get("platform"))

		if name == "" {
			http.Error(w, "name parameter required", http.StatusBadRequest)
			return
		}

		versions, err := prov.GetArtifactVersions(publisher, name, platform)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		items := make([]json.RawMessage, 0, len(versions))
		for _, v := range versions {
			data, err := protojson.Marshal(v)
			if err != nil {
				continue
			}
			items = append(items, json.RawMessage(data))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"versions": items,
			"count":    len(items),
		})
	})
}

// NewRepositoryDeleteHandler returns a handler for DELETE /admin/repository/artifact.
//
// Query params: publisher, name, version, platform
func NewRepositoryDeleteHandler(prov RepositoryProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()
		ref := &repopb.ArtifactRef{
			PublisherId: strings.TrimSpace(q.Get("publisher")),
			Name:        strings.TrimSpace(q.Get("name")),
			Version:     strings.TrimSpace(q.Get("version")),
			Platform:    strings.TrimSpace(q.Get("platform")),
		}
		if ref.Name == "" || ref.Version == "" {
			http.Error(w, "name and version parameters required", http.StatusBadRequest)
			return
		}

		if err := prov.DeleteArtifact(ref); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"deleted": true,
			"ref": map[string]string{
				"publisher": ref.PublisherId,
				"name":      ref.Name,
				"version":   ref.Version,
				"platform":  ref.Platform,
			},
		})
	})
}
