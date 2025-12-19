package handlers

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/globulario/Globular/internal/controllerclient"
	globpkg "github.com/globulario/Globular/internal/globule"
	clusterHandlers "github.com/globulario/Globular/internal/handlers/cluster"
	cfgHandlers "github.com/globulario/Globular/internal/handlers/config"
	filesHandlers "github.com/globulario/Globular/internal/handlers/files"
	mediaHandlers "github.com/globulario/Globular/internal/handlers/media"
	httplib "github.com/globulario/Globular/internal/http"
	middleware "github.com/globulario/Globular/internal/http/middleware"
	Utility "github.com/globulario/utility"
)

// HandlerConfig holds knobs consumed by the HTTP surface.
type HandlerConfig struct {
	MaxUpload      int64
	RateRPS        int
	RateBurst      int
	NodeAgentAddr  string
	ControllerAddr string
	EnvoyHTTPAddr  string
	Mode           string
}

// GatewayHandlers owns the HTTP middleware, providers, and wiring logic.
type GatewayHandlers struct {
	globule *globpkg.Globule
	cfg     HandlerConfig
}

// New builds a handler set tied to the given Globule instance.
func New(globule *globpkg.Globule, cfg HandlerConfig) *GatewayHandlers {
	return &GatewayHandlers{globule: globule, cfg: cfg}
}

// Router creates the mux, mounts static routes, and wires all handlers.
func (h *GatewayHandlers) Router(logger *slog.Logger) *http.ServeMux {
	wrap := h.Middleware()
	serve := filesHandlers.NewServeFile(h.newServeProvider())

	mux := httplib.NewRouter(logger, httplib.Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		RateRPS:        h.cfg.RateRPS,
		RateBurst:      h.cfg.RateBurst,
	}, wrap(serve))

	mux.Handle("/serve/", wrap(http.StripPrefix("/serve", serve)))

	h.wireConfig(mux, wrap)
	h.wireFiles(mux, wrap)
	h.wireMedia(mux, wrap)
	h.wireCluster(mux, wrap)

	return mux
}

func (h *GatewayHandlers) Middleware() func(http.Handler) http.Handler {
	return middleware.WithRedirectAndPreflight(h, h.setHeaders)
}

func (h *GatewayHandlers) RedirectTo(host string) (bool, *middleware.Target) {
	ok, p := h.globule.RedirectTo(host)
	if !ok || p == nil {
		return false, nil
	}
	return true, &middleware.Target{
		Hostname:   p.Hostname,
		Domain:     p.Domain,
		Protocol:   p.Protocol,
		PortHTTP:   int(p.PortHttp),
		PortHTTPS:  int(p.PortHttps),
		LocalIP:    p.LocalIpAddress,
		ExternalIP: p.ExternalIpAddress,
		Raw:        p,
	}
}

func (h *GatewayHandlers) HandleRedirect(to *middleware.Target, w http.ResponseWriter, r *http.Request) {
	addr := to.Domain
	scheme := "http"
	if to.Protocol == "https" {
		addr += ":" + Utility.ToString(to.PortHTTPS)
	} else {
		addr += ":" + Utility.ToString(to.PortHTTP)
	}
	addr = strings.ReplaceAll(addr, ".localhost", "")

	u, _ := url.Parse(scheme + "://" + addr)
	proxy := httputil.NewSingleHostReverseProxy(u)

	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	proxy.ServeHTTP(w, r)
}

func (h *GatewayHandlers) setHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	allowedOrigin := h.globule.Protocol + "://" + h.globule.Domain
	if origin != "" {
		for _, allowed := range h.globule.AllowedOrigins {
			if allowed == "*" || allowed == origin {
				allowedOrigin = origin
				break
			}
		}
	}

	allowedMethods := strings.Join(h.globule.AllowedMethods, ",")
	allowedHeaders := strings.Join(h.globule.AllowedHeaders, ",")

	header := w.Header()
	if allowedOrigin != "" {
		header.Set("Access-Control-Allow-Origin", allowedOrigin)
		if allowedOrigin != "*" {
			header.Set("Access-Control-Allow-Credentials", "true")
		}
		header.Add("Vary", "Origin")
	}
	header.Set("Access-Control-Allow-Methods", allowedMethods)
	header.Set("Access-Control-Allow-Headers", allowedHeaders)
	header.Set("Access-Control-Allow-Private-Network", "true")

	if r.Method == http.MethodOptions {
		header.Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func (h *GatewayHandlers) wireConfig(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	getConfig := cfgHandlers.NewGetConfig(cfgProvider{globule: h.globule})
	getServiceConfig := cfgHandlers.NewGetServiceConfig(cfgProvider{globule: h.globule})
	saveConfig := cfgHandlers.NewSaveConfig(cfgSaver{globule: h.globule}, tokenValidator{})
	getSvcPerms := cfgHandlers.NewGetServicePermissions(svcPermsProvider{globule: h.globule})
	describeService := cfgHandlers.NewDescribeService(describeProvider{globule: h.globule})

	ca := cfgHandlers.NewCAProvider()

	cfgHandlers.Mount(mux, cfgHandlers.Deps{
		GetConfig:             wrap(getConfig),
		GetServiceConfig:      wrap(getServiceConfig),
		SaveConfig:            wrap(saveConfig),
		GetServicePermissions: wrap(getSvcPerms),
		DescribeService:       wrap(describeService),
		GetCACertificate:      wrap(cfgHandlers.NewGetCACertificate(ca)),
		SignCACertificate:     wrap(cfgHandlers.NewSignCACertificate(ca)),
		GetSANConf:            wrap(cfgHandlers.NewGetSANConf(ca)),
	})
}

func (h *GatewayHandlers) wireFiles(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	getImages := filesHandlers.NewGetImages(imgLister{})
	upload := filesHandlers.NewUploadFileWithOptions(h.newUploadProvider(), filesHandlers.UploadOptions{
		MaxBytes:    h.cfg.MaxUpload,
		AllowedExts: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".txt", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".mp4", ".webm", ".mov", ".avi", ".mkv", ".mp3", ".wav", ".zip", ".rar", ".7z", ".tar", ".gz", ".csv", ".json", ".xml", ".md", ".html", ".css", ".js", ".svg", ".ttf", ".otf", ".woff", ".woff2", ".eot", ".tgz"},
	})
	filesHandlers.Mount(mux, filesHandlers.Deps{
		GetImages: wrap(getImages),
		Upload:    wrap(upload),
	})
}

func (h *GatewayHandlers) wireMedia(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	titles := mediaHandlers.NewGetIMDBTitles(imdbTitles{})
	poster := mediaHandlers.NewGetIMDBPoster(imdbPoster{})
	seasonEpisode := mediaHandlers.NewGetIMDBSeasonEpisode(imdbSeasonEpisode{})
	trailer := mediaHandlers.NewGetIMDBTrailer(imdbTrailer{})

	mediaHandlers.Mount(mux, mediaHandlers.Deps{
		GetIMDBTitles:        wrap(titles),
		GetIMDBPoster:        wrap(poster),
		GetIMDBSeasonEpisode: wrap(seasonEpisode),
		GetIMDBTrailer:       wrap(trailer),
	})
}

func (h *GatewayHandlers) wireCluster(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	deps := clusterHandlers.HandlerDeps{
		Controller:    controllerclient.New(h.cfg.ControllerAddr),
		NodeAgentAddr: h.cfg.NodeAgentAddr,
	}
	clusterHandlers.Mount(mux, clusterHandlers.Deps{
		JoinToken:   wrap(clusterHandlers.NewJoinTokenHandler(deps)),
		Nodes:       wrap(clusterHandlers.NewNodesHandler(deps)),
		NodeActions: wrap(clusterHandlers.NewNodeActionsHandler(deps)),
	})
}
