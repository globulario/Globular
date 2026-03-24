package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/StalkR/imdb"
	mediaHandlers "github.com/globulario/Globular/internal/gateway/handlers/media"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/title/titlepb"
	Utility "github.com/globulario/utility"
)

// tmdbAPIKey returns the TMDb API key, trying env var first then title service config.
var (
	tmdbKeyMu    sync.Mutex
	tmdbKeyValue string
	tmdbKeyDone  bool
)

func tmdbAPIKey() string {
	tmdbKeyMu.Lock()
	defer tmdbKeyMu.Unlock()

	// If already resolved successfully, return cached value.
	if tmdbKeyDone {
		return tmdbKeyValue
	}

	// 1. Environment variable
	if k := strings.TrimSpace(os.Getenv("TMDB_API_KEY")); k != "" {
		tmdbKeyValue = k
		tmdbKeyDone = true
		return tmdbKeyValue
	}

	// 2. Title or media service config in etcd (may not be available yet at startup).
	for _, svcName := range []string{"title.TitleService", "media.MediaService"} {
		cfg, err := config.GetServiceConfigurationById(svcName)
		if err != nil {
			slog.Debug("tmdbAPIKey: GetServiceConfigurationById failed", "svc", svcName, "err", err)
			cfgs, err2 := config.GetServicesConfigurationsByName(svcName)
			if err2 == nil && len(cfgs) > 0 {
				cfg = cfgs[0]
			}
		}
		if cfg != nil {
			if v, ok := cfg["TmdbApiKey"].(string); ok && v != "" {
				tmdbKeyValue = strings.TrimSpace(v)
				tmdbKeyDone = true
				slog.Info("TMDb API key loaded from service config", "svc", svcName)
				return tmdbKeyValue
			}
		}
	}

	slog.Warn("TMDb API key not found — set TMDB_API_KEY env var or TmdbApiKey in title service config")
	// Not found yet — will retry on next call.
	return ""
}

// ----- IMDb helpers (shared by media adapters) ----

// -- Season/Episode scraping (best-effort) --

type imdbSeasonEpisode struct{ Client *http.Client }

type imdbTrailer struct{}

type imdbPoster struct{}

type imdbTitles struct{}

type imdbHeaderTransport struct {
	base http.RoundTripper
}

func (s imdbSeasonEpisode) ResolveSeasonEpisode(titleID string) (int, int, string, error) {
	if s.Client == nil {
		s.Client = &http.Client{Timeout: 10 * time.Second}
	}
	req, _ := http.NewRequest(http.MethodGet, "https://www.imdb.com/title/"+titleID+"/", nil)
	req.Header.Set("User-Agent", defaultUA())
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	resp, err := s.Client.Do(req)
	if err != nil {
		return -1, -1, "", err
	}
	defer resp.Body.Close()

	page, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, "", err
	}

	reSE := regexp.MustCompile(`>S(\d+)<!-- -->\.<!-- -->E(\d+)<`)
	season, episode := 0, 0
	if m := reSE.FindSubmatch(page); len(m) == 3 {
		if v, e := strconv.Atoi(string(m[1])); e == nil {
			season = v
		}
		if v, e := strconv.Atoi(string(m[2])); e == nil {
			episode = v
		}
	}

	reSeries := regexp.MustCompile(`(?s)data-testid="hero-title-block__series-link".*?href="/title/(tt\d{7,8})/`)
	seriesID := ""
	if m := reSeries.FindSubmatch(page); len(m) == 2 {
		seriesID = string(m[1])
	} else {
		// Alternate pattern — IMDb sometimes reorders attributes.
		re2 := regexp.MustCompile(`href="/title/(tt\d{7,8})/[^"]*"[^>]*data-testid="hero-title-block__series-link"`)
		if m2 := re2.FindSubmatch(page); len(m2) == 2 {
			seriesID = string(m2[1])
		}
	}
	return season, episode, seriesID, nil
}

func (imdbTrailer) FetchIMDBTrailer(id string) (string, string, string, error) {
	return fetchIMDBTrailer(id)
}

func (imdbPoster) FetchIMDBPoster(id, size string) ([]byte, string, string, error) {
	u, err := fetchIMDBPosterURL(id)
	if err != nil {
		return nil, "", "", err
	}
	if size != "" {
		u = rewriteIMDBImageSize(u, size)
	}
	return nil, "", u, nil
}

func GetIMDBPoster(imdbID string) (string, error) { return fetchIMDBPosterURL(imdbID) }

func GetIMDBPosterSized(imdbID, size string) (string, error) {
	u, err := fetchIMDBPosterURL(imdbID)
	if err != nil {
		return "", err
	}
	return rewriteIMDBImageSize(u, size), nil
}

func fetchIMDBTrailer(imdbID string) (string, string, string, error) {
	page, err := fetchIMDBHTML("https://www.imdb.com/title/" + imdbID + "/")
	if err != nil {
		return "", "", "", err
	}

	if u, img := extractTrailerFromTitle(page); u != "" {
		if strings.Contains(u, "/video/") {
			videoSrc, err := fetchVideoSource(u)
			if err != nil {
				return u, img, "", nil
			}
			return u, img, videoSrc, nil
		}
		if videoURL, videoImg, err := findTrailerInGallery(u); err == nil && videoURL != "" {
			if img == "" {
				img = videoImg
			}
			videoSrc, err := fetchVideoSource(videoURL)
			if err != nil {
				return videoURL, img, "", nil
			}
			return videoURL, img, videoSrc, nil
		}
		return u, img, "", nil
	}

	return "", "", "", nil
}

func fetchIMDBHTML(url string) (string, error) {
	client := &http.Client{Timeout: 12 * time.Second}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", defaultUA())
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func extractTrailerFromTitle(page string) (string, string) {
	reOG := regexp.MustCompile(`(?i)<meta\s+(?:property|name)=["']og:video(?::secure_url)?["']\s+content=["']([^"']+)["']`)
	if m := reOG.FindStringSubmatch(page); len(m) == 2 && m[1] != "" {
		img := ""
		if m2 := regexp.MustCompile(`(?i)<meta\s+(?:property|name)=["']og:image(?::secure_url)?["']\s+content=["']([^"']+)["']`).FindStringSubmatch(page); len(m2) == 2 {
			img = m2[1]
		}
		return m[1], img
	}

	reMedia := regexp.MustCompile(`(?s)data-testid="[^\"]*hero-media__slate[^\"]*".*?<a[^>]+href="([^\"]+)"[^>]*>.*?<img[^>]+src="([^\"]+)"[^>]*`)
	if m := reMedia.FindStringSubmatch(page); len(m) == 3 {
		return absoluteIMDBURL(m[1]), m[2]
	}

	cands := collectVideoCandidates(page)
	return chooseCandidate(cands)
}

func findTrailerInGallery(galleryURL string) (string, string, error) {
	page, err := fetchIMDBHTML(absoluteIMDBURL(galleryURL))
	if err != nil {
		return "", "", err
	}
	if url, img := chooseCandidate(collectVideoCandidates(page)); url != "" {
		return url, img, nil
	}
	return "", "", nil
}

type videoCandidate struct {
	url   string
	image string
	label string
}

func collectVideoCandidates(page string) []videoCandidate {
	var out []videoCandidate
	reBlock := regexp.MustCompile(`(?s)<div[^>]+class="[^\"]*ipc-media__img[^\"]*".*?<img[^>]+src="([^\"]+)"[^>]*>.*?</div>\s*<a[^>]+href="(/video/[^"]+)"[^>]*?(?:aria-label|title)="([^\"]+)"`)
	for _, m := range reBlock.FindAllStringSubmatch(page, -1) {
		out = append(out, videoCandidate{url: absoluteIMDBURL(m[2]), image: m[1], label: m[3]})
	}
	reGeneric := regexp.MustCompile(`(?s)<a[^>]+href="(/video/[^"]+)"[^>]*?(?:aria-label|title)="([^\"]+)"`)
	for _, m := range reGeneric.FindAllStringSubmatch(page, -1) {
		out = append(out, videoCandidate{url: absoluteIMDBURL(m[1]), label: m[2]})
	}
	return out
}

func chooseCandidate(cands []videoCandidate) (string, string) {
	for _, c := range cands {
		if strings.Contains(strings.ToLower(c.label), "trailer") {
			return c.url, c.image
		}
	}
	if len(cands) > 0 {
		return cands[0].url, cands[0].image
	}
	return "", ""
}

func fetchVideoSource(videoPage string) (string, error) {
	page, err := fetchIMDBHTML(absoluteIMDBURL(videoPage))
	if err != nil {
		return "", err
	}
	reNext := regexp.MustCompile(`(?s)<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)
	if m := reNext.FindStringSubmatch(page); len(m) == 2 {
		var data map[string]any
		if json.Unmarshal([]byte(m[1]), &data) == nil {
			if url := pickVideoURL(lookupPath(data, "props", "pageProps", "videoPlaybackData", "videoLegacyEncodings")); url != "" {
				return url, nil
			}
			if url := pickVideoURL(lookupPath(data, "props", "pageProps", "videoPlaybackData", "playbackURLs")); url != "" {
				return url, nil
			}
		}
	}
	reURL := regexp.MustCompile(`"videoUrl":"([^\"]+)"`)
	if m := reURL.FindStringSubmatch(page); len(m) == 2 {
		return strings.ReplaceAll(m[1], "\\/", "/"), nil
	}
	return "", fmt.Errorf("video source not found")
}

func pickVideoURL(node any) string {
	arr, ok := node.([]any)
	if !ok {
		return ""
	}
	var fallback string
	for _, it := range arr {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		url := Utility.ToString(m["url"])
		mime := strings.ToLower(Utility.ToString(m["mimeType"]))
		if strings.Contains(mime, "mp4") && url != "" {
			return strings.ReplaceAll(url, "\\/", "/")
		}
		if fallback == "" && url != "" {
			fallback = strings.ReplaceAll(url, "\\/", "/")
		}
	}
	return fallback
}

// -- Titles (search) using IMDb suggestion API --

const imdbIDPattern = `^tt\d+$`

var imdbIDRE = regexp.MustCompile(imdbIDPattern)

func (t imdbHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	if r.Header.Get("User-Agent") == "" {
		r.Header.Set("User-Agent", defaultUA())
	}
	if r.Header.Get("Accept-Language") == "" {
		r.Header.Set("Accept-Language", "en-US,en;q=0.9")
	}
	if r.Header.Get("Accept") == "" {
		r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	}
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(r)
}

func newIMDBClient(timeout time.Duration) *http.Client {
	base := http.DefaultTransport
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		base = t.Clone()
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: imdbHeaderTransport{base: base},
	}
}

func (imdbTitles) SearchIMDBTitles(q mediaHandlers.TitlesQuery) ([]map[string]any, error) {
	query := strings.TrimSpace(q.Q)
	if query == "" {
		return nil, fmt.Errorf("empty query")
	}

	apiKey := tmdbAPIKey()

	// Direct IMDb ID lookup — use TMDb find API (reliable, no WAF).
	if imdbIDRE.MatchString(query) {
		if apiKey != "" {
			if t := tmdbFindByIMDbID(apiKey, query); t != nil {
				return []map[string]any{titleProtoToMap(t)}, nil
			}
		}
		// TMDb miss or no API key — try IMDb suggestion API as last resort.
		client := newIMDBClient(10 * time.Second)
		resolver := imdbSeasonEpisode{Client: client}
		if t, err := fetchTitleFromSuggestionAPI(client, query, resolver); err == nil {
			return []map[string]any{titleProtoToMap(t)}, nil
		}
		return nil, fmt.Errorf("title %s not found", query)
	}

	// Text search — use TMDb multi-search for rich results.
	if apiKey != "" {
		results := tmdbSearchMulti(apiKey, query, q.Year, q.Limit, q.Offset)
		if len(results) > 0 {
			out := make([]map[string]any, 0, len(results))
			for _, t := range results {
				out = append(out, titleProtoToMap(t))
			}
			return out, nil
		}
	}

	// Fallback: IMDb search (may fail with WAF).
	client := newIMDBClient(10 * time.Second)
	resolver := imdbSeasonEpisode{Client: client}
	results, err := imdb.SearchTitle(client, query)
	if err != nil {
		return nil, err
	}

	filtered := make([]imdb.Title, 0, len(results))
	for _, it := range results {
		if it.ID == "" || it.Name == "" {
			continue
		}
		if q.Year > 0 && it.Year != q.Year {
			continue
		}
		if qt := strings.TrimSpace(q.Type); qt != "" && !strings.Contains(strings.ToLower(it.Type), strings.ToLower(qt)) {
			continue
		}
		filtered = append(filtered, it)
	}
	if len(filtered) == 0 {
		filtered = results
	}

	start := q.Offset
	if start < 0 {
		start = 0
	}
	if start > len(filtered) {
		return []map[string]any{}, nil
	}
	end := start + q.Limit
	if end > len(filtered) {
		end = len(filtered)
	}
	selected := filtered[start:end]

	out := make([]map[string]any, 0, len(selected))
	for _, it := range selected {
		out = append(out, titleProtoToMap(buildTitleProto(it, resolver)))
	}
	return out, nil
}

func titleProtoToMap(t *titlepb.Title) map[string]any {
	if t == nil {
		return nil
	}
	out := map[string]any{
		"ID":          t.GetID(),
		"URL":         t.GetURL(),
		"Name":        t.GetName(),
		"Type":        t.GetType(),
		"Year":        t.GetYear(),
		"Rating":      t.GetRating(),
		"RatingCount": t.GetRatingCount(),
		"Description": t.GetDescription(),
		"Duration":    t.GetDuration(),
	}
	if v := t.GetSerie(); v != "" {
		out["Serie"] = v
	}
	if v := t.GetSeason(); v > 0 {
		out["Season"] = v
	}
	if v := t.GetEpisode(); v > 0 {
		out["Episode"] = v
	}
	if v := t.GetUUID(); v != "" {
		out["UUID"] = v
	}
	if langs := t.GetLanguage(); len(langs) > 0 {
		out["Languages"] = append([]string(nil), langs...)
	}
	if genres := t.GetGenres(); len(genres) > 0 {
		out["Genres"] = append([]string(nil), genres...)
	}
	if nat := t.GetNationalities(); len(nat) > 0 {
		out["Nationalities"] = append([]string(nil), nat...)
	}
	if aka := t.GetAKA(); len(aka) > 0 {
		out["AKA"] = append([]string(nil), aka...)
	}
	if dirs := t.GetDirectors(); len(dirs) > 0 {
		out["Directors"] = personsToMaps(dirs)
	}
	if writers := t.GetWriters(); len(writers) > 0 {
		out["Writers"] = personsToMaps(writers)
	}
	if actors := t.GetActors(); len(actors) > 0 {
		out["Actors"] = personsToMaps(actors)
	}
	if poster := t.GetPoster(); poster != nil {
		pm := map[string]any{}
		if v := poster.GetID(); v != "" {
			pm["ID"] = v
		}
		if v := poster.GetTitleId(); v != "" {
			pm["TitleID"] = v
		}
		if v := poster.GetURL(); v != "" {
			pm["URL"] = v
		}
		if v := poster.GetContentUrl(); v != "" {
			pm["ContentURL"] = v
		}
		if len(pm) > 0 {
			out["Poster"] = pm
		}
	}
	return out
}

func buildTitleProto(it imdb.Title, resolver imdbSeasonEpisode) *titlepb.Title {
	var rating float32
	if it.Rating != "" {
		if v, err := strconv.ParseFloat(it.Rating, 32); err == nil {
			rating = float32(v)
		}
	}

	title := &titlepb.Title{
		ID:            it.ID,
		URL:           it.URL,
		Name:          it.Name,
		Type:          it.Type,
		Year:          int32(it.Year),
		Rating:        rating,
		RatingCount:   int32(it.RatingCount),
		Description:   it.Description,
		Genres:        append([]string(nil), it.Genres...),
		Language:      append([]string(nil), it.Languages...),
		Nationalities: append([]string(nil), it.Nationalities...),
		Duration:      it.Duration,
		Directors:     namesToPersons(it.Directors),
		Writers:       namesToPersons(it.Writers),
		Actors:        namesToPersons(it.Actors),
	}

	if strings.EqualFold(it.Type, "TVEpisode") {
		if season, episode, serie, err := resolver.ResolveSeasonEpisode(it.ID); err == nil && (season > 0 || episode > 0 || serie != "") {
			title.Season = int32(season)
			title.Episode = int32(episode)
			title.Serie = serie
		} else {
			// IMDb HTML scraping failed (WAF) — fall back to TMDb API.
			if s, e, sid, ok := tmdbResolveEpisodeInfo(it.ID); ok {
				title.Season = int32(s)
				title.Episode = int32(e)
				title.Serie = sid
			}
		}
	}

	if poster, err := GetIMDBPoster(it.ID); err == nil && poster != "" {
		title.Poster = &titlepb.Poster{URL: poster}
	}

	return title
}

func namesToPersons(names []imdb.Name) []*titlepb.Person {
	out := make([]*titlepb.Person, 0, len(names))
	for _, n := range names {
		out = append(out, &titlepb.Person{
			ID:       n.ID,
			URL:      n.URL,
			FullName: n.FullName,
		})
	}
	return out
}

func personsToMaps(list []*titlepb.Person) []map[string]any {
	out := make([]map[string]any, 0, len(list))
	for _, p := range list {
		if p == nil {
			continue
		}
		m := map[string]any{}
		if v := p.GetID(); v != "" {
			m["ID"] = v
		}
		if v := p.GetURL(); v != "" {
			m["URL"] = v
		}
		if v := p.GetFullName(); v != "" {
			m["FullName"] = v
		}
		if aliases := p.GetAliases(); len(aliases) > 0 {
			m["Aliases"] = append([]string(nil), aliases...)
		}
		if len(m) > 0 {
			out = append(out, m)
		}
	}
	return out
}

func fetchIMDBPosterURL(imdbID string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Try the suggestion API first — it's reliable and not blocked by WAF
	if u := fetchPosterFromSuggestionAPI(client, imdbID); u != "" {
		return u, nil
	}

	// Fall back to HTML scraping (may fail with WAF challenge)
	req, _ := http.NewRequest(http.MethodGet, "https://www.imdb.com/title/"+imdbID+"/", nil)
	req.Header.Set("User-Agent", defaultUA())
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// If WAF blocks us (202/403), return what we have
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("imdb returned %d (WAF challenge)", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	page := string(b)

	reOG := regexp.MustCompile(`(?i)<meta\s+(?:property|name)=["']og:image["']\s+content=["']([^"']+)["']`)
	if m := reOG.FindStringSubmatch(page); len(m) == 2 && m[1] != "" {
		return m[1], nil
	}

	reLD := regexp.MustCompile(`(?s)<script[^>]+type=["']application/ld\+json["'][^>]*>(.*?)</script>`)
	for _, m := range reLD.FindAllStringSubmatch(page, -1) {
		var v any
		if json.Unmarshal([]byte(m[1]), &v) == nil {
			if u := extractImageFromLDJSON(v); u != "" {
				return u, nil
			}
		}
	}

	rePosterLink := regexp.MustCompile(`(?s)data-testid="hero-media__poster".*?href="([^\"]+mediaviewer[^\"]+)"`)
	viewerPath := ""
	if m := rePosterLink.FindStringSubmatch(page); len(m) == 2 {
		viewerPath = m[1]
	}
	if viewerPath == "" {
		reViewer := regexp.MustCompile(`href="(/title/tt\d+/mediaviewer/[^"]+)"`)
		if m := reViewer.FindStringSubmatch(page); len(m) == 2 {
			viewerPath = m[1]
		}
	}
	if viewerPath == "" {
		return "", fmt.Errorf("poster not found")
	}

	vReq, _ := http.NewRequest(http.MethodGet, "https://www.imdb.com"+viewerPath, nil)
	vReq.Header.Set("User-Agent", defaultUA())
	vReq.Header.Set("Accept-Language", "en-US,en;q=0.9")
	vResp, err := client.Do(vReq)
	if err != nil {
		return "", err
	}
	defer vResp.Body.Close()

	vb, err := io.ReadAll(vResp.Body)
	if err != nil {
		return "", err
	}
	view := string(vb)

	reSrcset := regexp.MustCompile(`(?s)<img[^>]+srcset="([^\"]+)"[^>]*>`)
	if m := reSrcset.FindStringSubmatch(view); len(m) == 2 {
		if u := pickMaxWidthFromSrcset(m[1]); u != "" {
			return u, nil
		}
	}

	reSrc := regexp.MustCompile(`(?s)<img[^>]+src="([^\"]+)"[^>]*>`)
	if m := reSrc.FindStringSubmatch(view); len(m) == 2 && m[1] != "" {
		return m[1], nil
	}

	return "", fmt.Errorf("poster not found")
}

// fetchPosterFromSuggestionAPI gets the poster URL from the IMDb suggestion API.
func fetchPosterFromSuggestionAPI(client *http.Client, imdbID string) string {
	if len(imdbID) < 3 {
		return ""
	}
	url := fmt.Sprintf("https://v2.sg.media-imdb.com/suggestion/%s/%s.json",
		string(imdbID[0]), imdbID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", defaultUA())

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var data struct {
		D []struct {
			ID string `json:"id"`
			I  struct {
				ImageURL string `json:"imageUrl"`
			} `json:"i"`
		} `json:"d"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ""
	}
	for _, d := range data.D {
		if d.ID == imdbID && d.I.ImageURL != "" {
			return d.I.ImageURL
		}
	}
	return ""
}

func pickMaxWidthFromSrcset(srcset string) string {
	var maxURL string
	maxW := -1
	for _, part := range strings.Split(srcset, ",") {
		part = strings.TrimSpace(part)
		items := strings.Fields(part)
		if len(items) != 2 {
			continue
		}
		u := items[0]
		w := strings.TrimSuffix(items[1], "w")
		if n, err := strconv.Atoi(w); err == nil && n > maxW {
			maxW, maxURL = n, u
		}
	}
	return maxURL
}

func rewriteIMDBImageSize(u, size string) string {
	size = strings.ToLower(size)
	if size == "" || size == "orig" {
		return u
	}
	target := map[string]string{
		"small":  "UX400",
		"medium": "UX800",
		"large":  "UX1200",
	}[size]
	if target == "" {
		return u
	}
	re := regexp.MustCompile(`\._V1_[^_.]*(_[A-Z]\w+)*(\.[a-zA-Z0-9]+)$`)
	return re.ReplaceAllString(u, "._V1_"+target+"$2")
}

func extractImageFromLDJSON(v any) string {
	switch x := v.(type) {
	case map[string]any:
		if img := x["image"]; img != nil {
			if u := imageURLFromAny(img); u != "" {
				return u
			}
		}
		if g, ok := x["@graph"]; ok {
			if u := extractImageFromLDJSON(g); u != "" {
				return u
			}
		}
	case []any:
		for _, it := range x {
			if u := extractImageFromLDJSON(it); u != "" {
				return u
			}
		}
	}
	return ""
}

func imageURLFromAny(a any) string {
	switch y := a.(type) {
	case string:
		return y
	case map[string]any:
		if u, _ := y["url"].(string); u != "" {
			return u
		}
	case []any:
		for _, it := range y {
			if u := imageURLFromAny(it); u != "" {
				return u
			}
		}
	}
	return ""
}

func absoluteIMDBURL(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	if strings.HasPrefix(path, "//") {
		return "https:" + path
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "https://www.imdb.com" + path
}

func lookupPath(v any, keys ...string) any {
	cur := v
	for _, k := range keys {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[k]
	}
	return cur
}

func defaultUA() string {
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36"
}

// fetchTitleFromSuggestionAPI uses the IMDb suggestion API (JSON, no WAF) as a
// fallback when HTML scraping is blocked. Returns basic title info.
func fetchTitleFromSuggestionAPI(client *http.Client, imdbID string, resolver imdbSeasonEpisode) (*titlepb.Title, error) {
	// The suggestion API indexes by the first letter after "tt"
	if len(imdbID) < 3 {
		return nil, fmt.Errorf("invalid imdb id: %s", imdbID)
	}
	url := fmt.Sprintf("https://v2.sg.media-imdb.com/suggestion/%s/%s.json",
		string(imdbID[0]), imdbID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaultUA())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("suggestion API returned %d", resp.StatusCode)
	}

	var data struct {
		D []struct {
			ID  string `json:"id"`
			L   string `json:"l"`   // title name
			Q   string `json:"q"`   // type label (e.g. "TV episode")
			QID string `json:"qid"` // type key (e.g. "tvEpisode")
			S   string `json:"s"`   // stars
			Y   int    `json:"y"`   // year
			I   struct {
				ImageURL string `json:"imageUrl"`
				Height   int    `json:"height"`
				Width    int    `json:"width"`
			} `json:"i"`
		} `json:"d"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Find our specific title in the results
	for _, d := range data.D {
		if d.ID != imdbID {
			continue
		}

		titleType := mapSuggestionType(d.QID)
		title := &titlepb.Title{
			ID:   d.ID,
			URL:  fmt.Sprintf("https://www.imdb.com/title/%s/", d.ID),
			Name: d.L,
			Type: titleType,
			Year: int32(d.Y),
		}

		if d.I.ImageURL != "" {
			title.Poster = &titlepb.Poster{URL: d.I.ImageURL}
		}

		// Parse stars into actors
		if d.S != "" {
			for _, name := range strings.Split(d.S, ", ") {
				name = strings.TrimSpace(name)
				if name != "" {
					title.Actors = append(title.Actors, &titlepb.Person{FullName: name})
				}
			}
		}

		// Resolve season/episode for TV episodes
		if strings.EqualFold(titleType, "TVEpisode") || d.QID == "tvEpisode" {
			if season, episode, serie, err := resolver.ResolveSeasonEpisode(d.ID); err == nil && (season > 0 || episode > 0 || serie != "") {
				title.Season = int32(season)
				title.Episode = int32(episode)
				title.Serie = serie
			} else {
				// IMDb HTML scraping failed (WAF) — fall back to TMDb API.
				if s, e, sid, ok := tmdbResolveEpisodeInfo(d.ID); ok {
					title.Season = int32(s)
					title.Episode = int32(e)
					title.Serie = sid
				}
			}
		}

		return title, nil
	}

	return nil, fmt.Errorf("title %s not found in suggestion API", imdbID)
}

// mapSuggestionType converts IMDb suggestion API qid values to display types.
func mapSuggestionType(qid string) string {
	switch qid {
	case "movie":
		return "Movie"
	case "tvSeries":
		return "TV Series"
	case "tvEpisode":
		return "TVEpisode"
	case "tvMiniSeries":
		return "TV Mini Series"
	case "tvMovie":
		return "TV Movie"
	case "tvSpecial":
		return "TV Special"
	case "short":
		return "Short"
	case "videoGame":
		return "Video Game"
	default:
		return qid
	}
}

// ---- TMDb-first title resolution ----

const tmdbImageBase = "https://image.tmdb.org/t/p/w500"

// tmdbFindByIMDbID looks up a title by IMDb ID via TMDb's find endpoint,
// then fetches full details (cast, genres, description, rating, etc.).
func tmdbFindByIMDbID(apiKey, imdbID string) *titlepb.Title {
	findURL := fmt.Sprintf("https://api.themoviedb.org/3/find/%s?api_key=%s&external_source=imdb_id",
		url.PathEscape(imdbID), url.QueryEscape(apiKey))

	var findResp struct {
		MovieResults []struct {
			ID               int     `json:"id"`
			Title            string  `json:"title"`
			Overview         string  `json:"overview"`
			ReleaseDate      string  `json:"release_date"`
			VoteAverage      float32 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
			PosterPath       string  `json:"poster_path"`
			GenreIDs         []int   `json:"genre_ids"`
			OriginalLanguage string  `json:"original_language"`
		} `json:"movie_results"`
		TvResults []struct {
			ID               int     `json:"id"`
			Name             string  `json:"name"`
			Overview         string  `json:"overview"`
			FirstAirDate     string  `json:"first_air_date"`
			VoteAverage      float32 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
			PosterPath       string  `json:"poster_path"`
			GenreIDs         []int   `json:"genre_ids"`
			OriginalLanguage string  `json:"original_language"`
		} `json:"tv_results"`
		TvEpisodeResults []struct {
			ID            int     `json:"id"`
			Name          string  `json:"name"`
			Overview      string  `json:"overview"`
			AirDate       string  `json:"air_date"`
			VoteAverage   float32 `json:"vote_average"`
			VoteCount     int     `json:"vote_count"`
			StillPath     string  `json:"still_path"`
			ShowID        int     `json:"show_id"`
			SeasonNumber  int     `json:"season_number"`
			EpisodeNumber int     `json:"episode_number"`
		} `json:"tv_episode_results"`
	}
	if err := tmdbGET(findURL, &findResp); err != nil {
		slog.Warn("tmdb find failed", "imdbID", imdbID, "err", err)
		return nil
	}

	if len(findResp.MovieResults) > 0 {
		r := findResp.MovieResults[0]
		title := &titlepb.Title{
			ID:          imdbID,
			URL:         fmt.Sprintf("https://www.imdb.com/title/%s/", imdbID),
			Name:        r.Title,
			Type:        "Movie",
			Description: r.Overview,
			Rating:      r.VoteAverage,
			RatingCount: int32(r.VoteCount),
		}
		if y := extractYear(r.ReleaseDate); y > 0 {
			title.Year = int32(y)
		}
		if r.PosterPath != "" {
			title.Poster = &titlepb.Poster{URL: tmdbImageBase + r.PosterPath}
		}
		title.Genres = tmdbGenreIDsToNames(r.GenreIDs)
		tmdbEnrichMovieDetails(apiKey, r.ID, title)
		return title
	}

	if len(findResp.TvResults) > 0 {
		r := findResp.TvResults[0]
		title := &titlepb.Title{
			ID:          imdbID,
			URL:         fmt.Sprintf("https://www.imdb.com/title/%s/", imdbID),
			Name:        r.Name,
			Type:        "TVSeries",
			Description: r.Overview,
			Rating:      r.VoteAverage,
			RatingCount: int32(r.VoteCount),
		}
		if y := extractYear(r.FirstAirDate); y > 0 {
			title.Year = int32(y)
		}
		if r.PosterPath != "" {
			title.Poster = &titlepb.Poster{URL: tmdbImageBase + r.PosterPath}
		}
		title.Genres = tmdbGenreIDsToNames(r.GenreIDs)
		tmdbEnrichTVDetails(apiKey, r.ID, title)
		return title
	}

	if len(findResp.TvEpisodeResults) > 0 {
		r := findResp.TvEpisodeResults[0]
		title := &titlepb.Title{
			ID:          imdbID,
			URL:         fmt.Sprintf("https://www.imdb.com/title/%s/", imdbID),
			Name:        r.Name,
			Type:        "TVEpisode",
			Description: r.Overview,
			Rating:      r.VoteAverage,
			RatingCount: int32(r.VoteCount),
			Season:      int32(r.SeasonNumber),
			Episode:     int32(r.EpisodeNumber),
		}
		if y := extractYear(r.AirDate); y > 0 {
			title.Year = int32(y)
		}
		if r.StillPath != "" {
			title.Poster = &titlepb.Poster{URL: tmdbImageBase + r.StillPath}
		}
		// Resolve parent series IMDb ID.
		if r.ShowID != 0 {
			if sid := tmdbGetSeriesIMDbID(apiKey, r.ShowID); sid != "" {
				title.Serie = sid
			}
			// Enrich episode with series genres and credits.
			tmdbEnrichEpisodeDetails(apiKey, r.ShowID, r.SeasonNumber, r.EpisodeNumber, title)
		}
		return title
	}

	return nil
}

// tmdbSearchMulti performs a TMDb multi-search and returns titlepb.Title results.
func tmdbSearchMulti(apiKey, query string, year, limit, offset int) []*titlepb.Title {
	page := (offset / max(limit, 1)) + 1
	searchURL := fmt.Sprintf("https://api.themoviedb.org/3/search/multi?api_key=%s&query=%s&page=%d",
		url.QueryEscape(apiKey), url.QueryEscape(query), page)
	if year > 0 {
		searchURL += fmt.Sprintf("&year=%d", year)
	}

	var searchResp struct {
		Results []struct {
			ID               int     `json:"id"`
			MediaType        string  `json:"media_type"` // movie, tv, person
			Title            string  `json:"title"`      // movie
			Name             string  `json:"name"`       // tv/person
			Overview         string  `json:"overview"`
			ReleaseDate      string  `json:"release_date"`
			FirstAirDate     string  `json:"first_air_date"`
			VoteAverage      float32 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
			PosterPath       string  `json:"poster_path"`
			ProfilePath      string  `json:"profile_path"`
			GenreIDs         []int   `json:"genre_ids"`
			OriginalLanguage string  `json:"original_language"`
		} `json:"results"`
	}
	if err := tmdbGET(searchURL, &searchResp); err != nil {
		slog.Warn("tmdb multi-search failed", "query", query, "err", err)
		return nil
	}

	var out []*titlepb.Title
	for _, r := range searchResp.Results {
		if r.MediaType == "person" {
			continue // skip person results
		}

		// Get the IMDb ID for this TMDb entry.
		imdbID := ""
		switch r.MediaType {
		case "movie":
			imdbID = tmdbGetMovieIMDbID(apiKey, r.ID)
		case "tv":
			imdbID = tmdbGetSeriesIMDbID(apiKey, r.ID)
		}
		if imdbID == "" {
			continue
		}

		title := &titlepb.Title{
			ID:          imdbID,
			URL:         fmt.Sprintf("https://www.imdb.com/title/%s/", imdbID),
			Description: r.Overview,
			Rating:      r.VoteAverage,
			RatingCount: int32(r.VoteCount),
			Genres:      tmdbGenreIDsToNames(r.GenreIDs),
		}

		switch r.MediaType {
		case "movie":
			title.Name = r.Title
			title.Type = "Movie"
			if y := extractYear(r.ReleaseDate); y > 0 {
				title.Year = int32(y)
			}
			if r.PosterPath != "" {
				title.Poster = &titlepb.Poster{URL: tmdbImageBase + r.PosterPath}
			}
			tmdbEnrichMovieDetails(apiKey, r.ID, title)
		case "tv":
			title.Name = r.Name
			title.Type = "TVSeries"
			if y := extractYear(r.FirstAirDate); y > 0 {
				title.Year = int32(y)
			}
			if r.PosterPath != "" {
				title.Poster = &titlepb.Poster{URL: tmdbImageBase + r.PosterPath}
			}
			tmdbEnrichTVDetails(apiKey, r.ID, title)
		}

		out = append(out, title)
		if len(out) >= limit {
			break
		}
	}
	return out
}

// tmdbGetSeriesIMDbID resolves a TMDb TV show ID to its IMDb ID.
func tmdbGetSeriesIMDbID(apiKey string, tmdbShowID int) string {
	extURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/external_ids?api_key=%s",
		tmdbShowID, url.QueryEscape(apiKey))
	var resp struct {
		IMDbID string `json:"imdb_id"`
	}
	if err := tmdbGET(extURL, &resp); err == nil {
		return resp.IMDbID
	}
	return ""
}

// tmdbGetMovieIMDbID resolves a TMDb movie ID to its IMDb ID.
func tmdbGetMovieIMDbID(apiKey string, tmdbMovieID int) string {
	extURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/external_ids?api_key=%s",
		tmdbMovieID, url.QueryEscape(apiKey))
	var resp struct {
		IMDbID string `json:"imdb_id"`
	}
	if err := tmdbGET(extURL, &resp); err == nil {
		return resp.IMDbID
	}
	return ""
}

// tmdbEnrichMovieDetails adds cast/crew/duration/language from TMDb movie details+credits.
func tmdbEnrichMovieDetails(apiKey string, tmdbID int, title *titlepb.Title) {
	// Details (runtime, languages, genres)
	detailURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s",
		tmdbID, url.QueryEscape(apiKey))
	var details struct {
		Runtime         int `json:"runtime"`
		SpokenLanguages []struct {
			Name string `json:"english_name"`
		} `json:"spoken_languages"`
		ProductionCountries []struct {
			Name string `json:"name"`
		} `json:"production_countries"`
		Genres []struct {
			Name string `json:"name"`
		} `json:"genres"`
	}
	if err := tmdbGET(detailURL, &details); err == nil {
		if title.Duration == "" && details.Runtime > 0 {
			title.Duration = fmt.Sprintf("%dm", details.Runtime)
		}
		if len(title.Language) == 0 {
			for _, l := range details.SpokenLanguages {
				if l.Name != "" {
					title.Language = append(title.Language, l.Name)
				}
			}
		}
		if len(title.Nationalities) == 0 {
			for _, c := range details.ProductionCountries {
				if c.Name != "" {
					title.Nationalities = append(title.Nationalities, c.Name)
				}
			}
		}
		if len(title.Genres) == 0 && len(details.Genres) > 0 {
			for _, g := range details.Genres {
				if g.Name != "" {
					title.Genres = append(title.Genres, g.Name)
				}
			}
		}
	}

	// Credits (cast + crew)
	tmdbEnrichCredits(apiKey, fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/credits?api_key=%s",
		tmdbID, url.QueryEscape(apiKey)), title)
}

// tmdbEnrichTVDetails adds cast/crew/duration/language from TMDb TV details+credits.
func tmdbEnrichTVDetails(apiKey string, tmdbID int, title *titlepb.Title) {
	detailURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d?api_key=%s",
		tmdbID, url.QueryEscape(apiKey))
	var details struct {
		EpisodeRunTime  []int `json:"episode_run_time"`
		SpokenLanguages []struct {
			Name string `json:"english_name"`
		} `json:"spoken_languages"`
		ProductionCountries []struct {
			Name string `json:"name"`
		} `json:"production_countries"`
		Genres []struct {
			Name string `json:"name"`
		} `json:"genres"`
		CreatedBy []struct {
			Name string `json:"name"`
		} `json:"created_by"`
	}
	if err := tmdbGET(detailURL, &details); err == nil {
		if title.Duration == "" && len(details.EpisodeRunTime) > 0 && details.EpisodeRunTime[0] > 0 {
			title.Duration = fmt.Sprintf("%dm", details.EpisodeRunTime[0])
		}
		if len(title.Language) == 0 {
			for _, l := range details.SpokenLanguages {
				if l.Name != "" {
					title.Language = append(title.Language, l.Name)
				}
			}
		}
		if len(title.Nationalities) == 0 {
			for _, c := range details.ProductionCountries {
				if c.Name != "" {
					title.Nationalities = append(title.Nationalities, c.Name)
				}
			}
		}
		if len(title.Genres) == 0 && len(details.Genres) > 0 {
			for _, g := range details.Genres {
				if g.Name != "" {
					title.Genres = append(title.Genres, g.Name)
				}
			}
		}
	}

	tmdbEnrichCredits(apiKey, fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/credits?api_key=%s",
		tmdbID, url.QueryEscape(apiKey)), title)
}

// tmdbEnrichEpisodeDetails adds credits and genres from the parent show.
func tmdbEnrichEpisodeDetails(apiKey string, showID, season, episode int, title *titlepb.Title) {
	// Get episode credits
	creditsURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d/episode/%d/credits?api_key=%s",
		showID, season, episode, url.QueryEscape(apiKey))
	tmdbEnrichCredits(apiKey, creditsURL, title)

	// Get episode details for runtime
	epURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d/episode/%d?api_key=%s",
		showID, season, episode, url.QueryEscape(apiKey))
	var epDetails struct {
		Runtime int `json:"runtime"`
	}
	if err := tmdbGET(epURL, &epDetails); err == nil && epDetails.Runtime > 0 && title.Duration == "" {
		title.Duration = fmt.Sprintf("%dm", epDetails.Runtime)
	}

	// Get genres from parent series
	if len(title.Genres) == 0 {
		showURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d?api_key=%s",
			showID, url.QueryEscape(apiKey))
		var showDetails struct {
			Genres []struct {
				Name string `json:"name"`
			} `json:"genres"`
		}
		if err := tmdbGET(showURL, &showDetails); err == nil {
			for _, g := range showDetails.Genres {
				if g.Name != "" {
					title.Genres = append(title.Genres, g.Name)
				}
			}
		}
	}
}

// tmdbEnrichCredits fetches credits from a TMDb credits URL and adds actors/directors/writers.
func tmdbEnrichCredits(apiKey, creditsURL string, title *titlepb.Title) {
	var credits struct {
		Cast []struct {
			Name string `json:"name"`
		} `json:"cast"`
		GuestStars []struct {
			Name string `json:"name"`
		} `json:"guest_stars"`
		Crew []struct {
			Name       string `json:"name"`
			Department string `json:"department"`
			Job        string `json:"job"`
		} `json:"crew"`
	}
	if err := tmdbGET(creditsURL, &credits); err != nil {
		return
	}

	if len(title.Actors) == 0 {
		seen := make(map[string]bool)
		for _, c := range credits.Cast {
			if c.Name != "" && !seen[c.Name] {
				title.Actors = append(title.Actors, &titlepb.Person{FullName: c.Name})
				seen[c.Name] = true
			}
		}
		for _, c := range credits.GuestStars {
			if c.Name != "" && !seen[c.Name] {
				title.Actors = append(title.Actors, &titlepb.Person{FullName: c.Name})
				seen[c.Name] = true
			}
		}
	}
	if len(title.Directors) == 0 {
		for _, c := range credits.Crew {
			if c.Job == "Director" && c.Name != "" {
				title.Directors = append(title.Directors, &titlepb.Person{FullName: c.Name})
			}
		}
	}
	if len(title.Writers) == 0 {
		for _, c := range credits.Crew {
			if (c.Department == "Writing" || c.Job == "Writer" || c.Job == "Screenplay" || c.Job == "Story") && c.Name != "" {
				title.Writers = append(title.Writers, &titlepb.Person{FullName: c.Name})
			}
		}
	}
}

func extractYear(dateStr string) int {
	if len(dateStr) >= 4 {
		if y, err := strconv.Atoi(dateStr[:4]); err == nil {
			return y
		}
	}
	return 0
}

// tmdbGenreIDsToNames converts TMDb genre IDs to names using a static map.
func tmdbGenreIDsToNames(ids []int) []string {
	if len(ids) == 0 {
		return nil
	}
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		if name, ok := tmdbGenreMap[id]; ok {
			names = append(names, name)
		}
	}
	return names
}

// TMDb genre ID → name (combined movie + TV).
var tmdbGenreMap = map[int]string{
	28: "Action", 12: "Adventure", 16: "Animation", 35: "Comedy",
	80: "Crime", 99: "Documentary", 18: "Drama", 10751: "Family",
	14: "Fantasy", 36: "History", 27: "Horror", 10402: "Music",
	9648: "Mystery", 10749: "Romance", 878: "Science Fiction",
	10770: "TV Movie", 53: "Thriller", 10752: "War", 37: "Western",
	10759: "Action & Adventure", 10762: "Kids", 10763: "News",
	10764: "Reality", 10765: "Sci-Fi & Fantasy", 10766: "Soap",
	10767: "Talk", 10768: "War & Politics",
}

// tmdbResolveEpisodeInfo uses TMDb to resolve episode season/episode/serie
// when IMDb HTML scraping is blocked. Used by buildTitleProto and fetchTitleFromSuggestionAPI.
func tmdbResolveEpisodeInfo(imdbID string) (season, episode int, serieIMDbID string, ok bool) {
	apiKey := tmdbAPIKey()
	if apiKey == "" {
		return 0, 0, "", false
	}
	t := tmdbFindByIMDbID(apiKey, imdbID)
	if t == nil || t.Type != "TVEpisode" {
		return 0, 0, "", false
	}
	return int(t.Season), int(t.Episode), t.Serie, t.Season > 0 || t.Episode > 0 || t.Serie != ""
}

func tmdbGET(u string, out any) error {
	c := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("tmdb returned %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
