package handlers

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/globulario/Globular/internal/controllerclient"
	adminHandlers "github.com/globulario/Globular/internal/gateway/handlers/admin"
	clusterHandlers "github.com/globulario/Globular/internal/gateway/handlers/cluster"
	cfgHandlers "github.com/globulario/Globular/internal/gateway/handlers/config"
	domainHandlers "github.com/globulario/Globular/internal/gateway/handlers/domains"
	filesHandlers "github.com/globulario/Globular/internal/gateway/handlers/files"
	mediaHandlers "github.com/globulario/Globular/internal/gateway/handlers/media"
	statsHandlers "github.com/globulario/Globular/internal/gateway/handlers/stats"
	httplib "github.com/globulario/Globular/internal/gateway/http"
	middleware "github.com/globulario/Globular/internal/gateway/http/middleware"
	globpkg "github.com/globulario/Globular/internal/globule"
	Utility "github.com/globulario/utility"
)

var serviceConfigCache = cfgHandlers.NewServiceConfigCache()

// HandlerConfig holds knobs consumed by the HTTP surface.
type HandlerConfig struct {
	MaxUpload      int64
	RateRPS        int
	RateBurst      int
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
	mux.Handle("/serve", wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/serve/", http.StatusMovedPermanently)
	})))
	mux.Handle("/globular/webroot/", wrap(http.StripPrefix("/globular/webroot", serve)))

	h.wireConfig(mux, wrap)
	h.wireFiles(mux, wrap)
	h.wireObjectStoreHealth(mux, wrap)
	h.wireMedia(mux, wrap)
	h.wireCluster(mux, wrap)
	h.wireStats(mux, wrap)
	h.wireAdmin(mux, wrap)
	h.wireDomains(mux, wrap)

	return mux
}

func (h *GatewayHandlers) Middleware() func(http.Handler) http.Handler {
	return middleware.WithRedirectAndPreflight(h, h.setHeaders)
}

func (h *GatewayHandlers) RedirectTo(host string) (bool, *middleware.Target) {
	return false, nil
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

	u, err := url.Parse(scheme + "://" + addr)
	if err != nil {
		httplib.WriteJSONError(w, http.StatusInternalServerError, "invalid redirect target URL")
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)

	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	proxy.ServeHTTP(w, r)
}

func (h *GatewayHandlers) setHeaders(w http.ResponseWriter, r *http.Request) {
	policy := h.globule.GetCorsPolicy()

	// If CORS is disabled, skip all headers.
	if !policy.Enabled {
		return
	}

	origin := r.Header.Get("Origin")
	allowedOrigin := ""

	if policy.AllowAllOrigins {
		if policy.AllowCredentials {
			// Spec: wildcard + credentials is invalid; reflect the origin instead.
			allowedOrigin = origin
		} else {
			allowedOrigin = "*"
		}
	} else if origin != "" {
		for _, allowed := range policy.AllowedOrigins {
			if allowed == origin {
				allowedOrigin = origin
				break
			}
		}
	}
	// Fallback: if no origin matched, use protocol://domain for same-site.
	if allowedOrigin == "" && origin == "" {
		allowedOrigin = h.globule.Protocol + "://" + h.globule.Domain
	}

	// Also check legacy fields as fallback (backward compat).
	if allowedOrigin == "" && origin != "" {
		for _, allowed := range h.globule.AllowedOrigins {
			if allowed == "*" || allowed == origin {
				allowedOrigin = origin
				break
			}
		}
	}

	header := w.Header()
	if allowedOrigin != "" {
		header.Set("Access-Control-Allow-Origin", allowedOrigin)
		if policy.AllowCredentials && allowedOrigin != "*" {
			header.Set("Access-Control-Allow-Credentials", "true")
		}
		header.Add("Vary", "Origin")
	}

	// Methods
	if len(policy.AllowedMethods) > 0 {
		header.Set("Access-Control-Allow-Methods", strings.Join(policy.AllowedMethods, ","))
	} else {
		header.Set("Access-Control-Allow-Methods", strings.Join(h.globule.AllowedMethods, ","))
	}

	// Headers
	if len(policy.AllowedHeaders) > 0 {
		header.Set("Access-Control-Allow-Headers", strings.Join(policy.AllowedHeaders, ","))
	} else {
		header.Set("Access-Control-Allow-Headers", strings.Join(h.globule.AllowedHeaders, ","))
	}

	// Exposed headers
	if len(policy.ExposedHeaders) > 0 {
		header.Set("Access-Control-Expose-Headers", strings.Join(policy.ExposedHeaders, ","))
	}

	// Private network
	if policy.AllowPrivateNetwork {
		header.Set("Access-Control-Allow-Private-Network", "true")
	}

	if r.Method == http.MethodOptions {
		if policy.MaxAgeSeconds > 0 {
			header.Set("Access-Control-Max-Age", strconv.Itoa(policy.MaxAgeSeconds))
		} else {
			header.Set("Access-Control-Max-Age", "3600")
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func (h *GatewayHandlers) wireConfig(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	getConfig := cfgHandlers.NewGetConfig(cfgProvider{globule: h.globule, cache: serviceConfigCache})
	getServiceConfig := cfgHandlers.NewGetServiceConfig(cfgProvider{globule: h.globule})
	saveConfig := cfgHandlers.NewSaveConfig(cfgSaver{globule: h.globule}, tokenValidator{})
	saveServiceConfig := cfgHandlers.NewSaveServiceConfig(cfgSaver{globule: h.globule}, tokenValidator{})
	getSvcPerms := cfgHandlers.NewGetServicePermissions(svcPermsProvider{globule: h.globule})
	describeService := cfgHandlers.NewDescribeService(describeProvider{})

	ca := cfgHandlers.NewCAProvider()
	getServicesCors := cfgHandlers.NewGetServicesCors(cfgProvider{globule: h.globule, cache: serviceConfigCache})
	setServiceCors := cfgHandlers.NewSetServiceCors(cfgSaver{globule: h.globule})

	// Structured CORS policy handlers (PR1)
	corsProvider := cfgProvider{globule: h.globule, cache: serviceConfigCache}
	corsSaver := cfgSaver{globule: h.globule}

	cfgHandlers.Mount(mux, cfgHandlers.Deps{
		GetConfig:             wrap(getConfig),
		GetServiceConfig:      wrap(getServiceConfig),
		SaveConfig:            wrap(saveConfig),
		SaveServiceConfig:     wrap(saveServiceConfig),
		GetServicePermissions: wrap(getSvcPerms),
		DescribeService:       wrap(describeService),
		GetCACertificate:      wrap(cfgHandlers.NewGetCACertificate(ca)),
		SignCACertificate:     wrap(cfgHandlers.NewSignCACertificate(ca)),
		GetSANConf:            wrap(cfgHandlers.NewGetSANConf(ca)),
		GetServicesCors:       wrap(getServicesCors),
		SetServiceCors:        wrap(setServiceCors),

		// Structured CORS policy (PR1)
		GetGatewayCorsPolicy:     wrap(cfgHandlers.NewGetGatewayCorsPolicy(corsProvider)),
		SetGatewayCorsPolicy:     wrap(cfgHandlers.NewSetGatewayCorsPolicy(corsSaver)),
		GetServiceCorsPolicy:     wrap(cfgHandlers.NewGetServiceCorsPolicy(corsProvider)),
		SetServiceCorsPolicy:     wrap(cfgHandlers.NewSetServiceCorsPolicy(corsSaver)),
		GetAllServicesCorsPolicy: wrap(cfgHandlers.NewGetAllServicesCorsPolicy(corsProvider, corsProvider)),
		CorsDiagnostics:          wrap(cfgHandlers.NewCorsDiagnostics(corsProvider)),
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

func (h *GatewayHandlers) wireObjectStoreHealth(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strict := r.URL.Query().Get("strict") == "1"
		var (
			cfg *filesHandlers.MinioProxyConfig
			err error
		)
		if strict {
			cfg, err = h.newServeProvider().FileServiceMinioConfigStrict(r.Context())
		} else {
			cfg, err = h.newServeProvider().FileServiceMinioConfig()
		}
		if err != nil || cfg == nil {
			httplib.WriteJSONError(w, http.StatusServiceUnavailable, "object store unavailable")
			return
		}
		httplib.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("/health/objectstore", wrap(handler))
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
		Controller: controllerclient.New(h.cfg.ControllerAddr),
	}
	clusterHandlers.Mount(mux, clusterHandlers.Deps{
		JoinToken:   wrap(clusterHandlers.NewJoinTokenHandler(deps)),
		Nodes:       wrap(clusterHandlers.NewNodesHandler(deps)),
		NodeActions: wrap(clusterHandlers.NewNodeActionsHandler(deps)),
	})
}

func (h *GatewayHandlers) wireStats(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	statsHandlers.Mount(mux, statsHandlers.Deps{
		Stats: wrap(statsHandlers.NewStatsHandler(h.globule)),
	})
}

func (h *GatewayHandlers) wireAdmin(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	prov := adminProvider{globule: h.globule}
	certProv := certProvider{globule: h.globule}
	controller := controllerclient.New(h.cfg.ControllerAddr)
	adminHandlers.Mount(mux, adminHandlers.Deps{
		MetricsServices:      wrap(adminHandlers.NewServicesHandler(prov)),
		MetricsStorage:       wrap(adminHandlers.NewStorageHandler(prov)),
		MetricsEnvoy:         wrap(adminHandlers.NewEnvoyHandler()),
		ServiceLogs:          wrap(adminHandlers.NewLogsHandler(journalAdapter{})),
		CertificatesOverview: wrap(adminHandlers.NewCertificatesHandler(certProv)),
		CertificatesCluster: wrap(adminHandlers.NewClusterCertificatesHandler(adminHandlers.ClusterCertDeps{
			Controller:  controller,
			GatewayPort: h.globule.PortHTTP,
			LocalProv:   certProv,
		})),
		RenewPublic:        wrap(adminHandlers.NewRenewPublicHandler(certProv)),
		RegenerateInternal: wrap(adminHandlers.NewRegenerateInternalHandler(certProv)),
	})
}

func (h *GatewayHandlers) wireDomains(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	prov := domainStoreProvider{}
	domainHandlers.Mount(mux, domainHandlers.Deps{
		ListProviders:  wrap(domainHandlers.NewListProviders(prov)),
		GetProvider:    wrap(domainHandlers.NewGetProvider(prov)),
		PutProvider:    wrap(domainHandlers.NewPutProvider(prov)),
		DeleteProvider: wrap(domainHandlers.NewDeleteProvider(prov)),
		ListDomains:    wrap(domainHandlers.NewListDomains(prov)),
		GetDomain:      wrap(domainHandlers.NewGetDomain(prov)),
		PutDomain:      wrap(domainHandlers.NewPutDomain(prov)),
		DeleteDomain:   wrap(domainHandlers.NewDeleteDomain(prov)),
	})
}
