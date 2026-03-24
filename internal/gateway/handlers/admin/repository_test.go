package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	repopb "github.com/globulario/services/golang/repository/repositorypb"
)

// mockRepoProvider implements RepositoryProvider for testing.
type mockRepoProvider struct {
	manifests []*repopb.ArtifactManifest
	deleteErr error
}

func (m *mockRepoProvider) SearchArtifacts(query, kind, publisher, platform, pageToken string, pageSize int32) ([]*repopb.ArtifactManifest, string, int32, error) {
	var result []*repopb.ArtifactManifest
	for _, a := range m.manifests {
		if query != "" {
			if a.Ref == nil || a.Ref.Name != query {
				continue
			}
		}
		if publisher != "" && a.Ref != nil && a.Ref.PublisherId != publisher {
			continue
		}
		result = append(result, a)
	}
	return result, "", int32(len(result)), nil
}

func (m *mockRepoProvider) GetArtifactManifest(ref *repopb.ArtifactRef, buildNumber int64) (*repopb.ArtifactManifest, error) {
	for _, a := range m.manifests {
		if a.Ref != nil && a.Ref.Name == ref.Name && a.Ref.Version == ref.Version {
			return a, nil
		}
	}
	return &repopb.ArtifactManifest{}, nil
}

func (m *mockRepoProvider) GetArtifactVersions(publisherID, name, platform string) ([]*repopb.ArtifactManifest, error) {
	var result []*repopb.ArtifactManifest
	for _, a := range m.manifests {
		if a.Ref != nil && a.Ref.Name == name {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockRepoProvider) DeleteArtifact(ref *repopb.ArtifactRef) error {
	return m.deleteErr
}

func testManifests() []*repopb.ArtifactManifest {
	return []*repopb.ArtifactManifest{
		{
			Ref:      &repopb.ArtifactRef{PublisherId: "core@globular.io", Name: "gateway", Version: "1.0.0", Platform: "linux_amd64", Kind: repopb.ArtifactKind_SERVICE},
			Checksum: "aaa",
		},
		{
			Ref:      &repopb.ArtifactRef{PublisherId: "core@globular.io", Name: "gateway", Version: "2.0.0", Platform: "linux_amd64", Kind: repopb.ArtifactKind_SERVICE},
			Checksum: "bbb",
		},
		{
			Ref:      &repopb.ArtifactRef{PublisherId: "core@globular.io", Name: "webadmin", Version: "1.0.0", Platform: "linux_amd64", Kind: repopb.ArtifactKind_APPLICATION},
			Checksum: "ccc",
		},
	}
}

func TestRepoSearch_AllArtifacts(t *testing.T) {
	prov := &mockRepoProvider{manifests: testManifests()}
	h := NewRepositorySearchHandler(prov)

	req := httptest.NewRequest(http.MethodGet, "/admin/repository/search", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	if int(body["count"].(float64)) != 3 {
		t.Errorf("count = %v, want 3", body["count"])
	}
}

func TestRepoSearch_ByQuery(t *testing.T) {
	prov := &mockRepoProvider{manifests: testManifests()}
	h := NewRepositorySearchHandler(prov)

	req := httptest.NewRequest(http.MethodGet, "/admin/repository/search?q=webadmin", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	if int(body["count"].(float64)) != 1 {
		t.Errorf("count = %v, want 1", body["count"])
	}
}

func TestRepoSearch_ByPublisher(t *testing.T) {
	prov := &mockRepoProvider{manifests: testManifests()}
	h := NewRepositorySearchHandler(prov)

	req := httptest.NewRequest(http.MethodGet, "/admin/repository/search?publisher=nonexistent", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	if int(body["count"].(float64)) != 0 {
		t.Errorf("count = %v, want 0", body["count"])
	}
}

func TestRepoSearch_MethodNotAllowed(t *testing.T) {
	prov := &mockRepoProvider{}
	h := NewRepositorySearchHandler(prov)

	req := httptest.NewRequest(http.MethodPost, "/admin/repository/search", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", w.Code)
	}
}

func TestRepoManifest_Basic(t *testing.T) {
	prov := &mockRepoProvider{manifests: testManifests()}
	h := NewRepositoryManifestHandler(prov)

	req := httptest.NewRequest(http.MethodGet, "/admin/repository/manifest?name=gateway&version=1.0.0", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestRepoManifest_MissingName(t *testing.T) {
	prov := &mockRepoProvider{}
	h := NewRepositoryManifestHandler(prov)

	req := httptest.NewRequest(http.MethodGet, "/admin/repository/manifest?version=1.0.0", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestRepoVersions_Basic(t *testing.T) {
	prov := &mockRepoProvider{manifests: testManifests()}
	h := NewRepositoryVersionsHandler(prov)

	req := httptest.NewRequest(http.MethodGet, "/admin/repository/versions?name=gateway", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	if int(body["count"].(float64)) != 2 {
		t.Errorf("count = %v, want 2 (gateway has 2 versions)", body["count"])
	}
}

func TestRepoVersions_MissingName(t *testing.T) {
	prov := &mockRepoProvider{}
	h := NewRepositoryVersionsHandler(prov)

	req := httptest.NewRequest(http.MethodGet, "/admin/repository/versions", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestRepoDelete_Basic(t *testing.T) {
	prov := &mockRepoProvider{manifests: testManifests()}
	h := NewRepositoryDeleteHandler(prov)

	req := httptest.NewRequest(http.MethodDelete, "/admin/repository/artifact?name=gateway&version=1.0.0", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	if body["deleted"] != true {
		t.Error("expected deleted: true")
	}
}

func TestRepoDelete_MissingParams(t *testing.T) {
	prov := &mockRepoProvider{}
	h := NewRepositoryDeleteHandler(prov)

	req := httptest.NewRequest(http.MethodDelete, "/admin/repository/artifact?name=gateway", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestRepoDelete_WrongMethod(t *testing.T) {
	prov := &mockRepoProvider{}
	h := NewRepositoryDeleteHandler(prov)

	req := httptest.NewRequest(http.MethodGet, "/admin/repository/artifact?name=gateway&version=1.0.0", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", w.Code)
	}
}
