package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/imdb"
	mediaHandlers "github.com/globulario/Globular/internal/gateway/handlers/media"
	"github.com/globulario/services/golang/title/titlepb"
	Utility "github.com/globulario/utility"
)

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

	reSE := regexp.MustCompile(`>S(\d{1,2})<!-- -->\.<!-- -->E(\d{1,2})<`)
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

	client := newIMDBClient(10 * time.Second)
	resolver := imdbSeasonEpisode{Client: client}

	if imdbIDRE.MatchString(query) {
		it, err := imdb.NewTitle(client, query)
		if err != nil {
			return nil, err
		}
		return []map[string]any{titleProtoToMap(buildTitleProto(*it, resolver))}, nil
	}

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
		if season, episode, serie, err := resolver.ResolveSeasonEpisode(it.ID); err == nil {
			title.Season = int32(season)
			title.Episode = int32(episode)
			title.Serie = serie
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

	req, _ := http.NewRequest(http.MethodGet, "https://www.imdb.com/title/"+imdbID+"/", nil)
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
