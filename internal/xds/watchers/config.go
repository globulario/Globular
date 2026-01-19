package watchers

const (
	defaultIngressHTTPPort    = uint32(80)
	defaultIngressHTTPSPort   = uint32(443)
	defaultGatewayListenPort  = uint32(8080)
	defaultIngressCertPath    = "/var/lib/globular/tls/fullchain.pem"
	defaultIngressKeyPath     = "/var/lib/globular/tls/privkey.pem"
	defaultIngressProvisionCA = "/var/lib/globular/tls/ca.pem"
)

// XDSConfig describes the static configuration needed by globular-xds.
type XDSConfig struct {
	EtcdEndpoints       []string        `json:"etcd_endpoints"`
	SyncIntervalSeconds int             `json:"sync_interval_seconds"`
	Ingress             IngressConfig   `json:"ingress"`
	Gateway             GatewayConfig   `json:"gateway"`
	Fallback            *FallbackConfig `json:"fallback"`
}

// IngressConfig defines HTTP/HTTPS ports, TLS paths, and redirect behavior.
type IngressConfig struct {
	HTTPPort           uint32      `json:"http_port"`
	HTTPSPort          uint32      `json:"https_port"`
	EnableHTTPRedirect *bool       `json:"enable_http_redirect"`
	TLS                IngressTLS  `json:"tls"`
	MTLS               IngressMTLS `json:"mtls"`
}

// IngressTLS describes TLS assets for ingress listeners.
type IngressTLS struct {
	Enabled        *bool  `json:"enabled"`
	CertChainPath  string `json:"cert_chain_path"`
	PrivateKeyPath string `json:"private_key_path"`
}

// TLSPaths points to PEM files on disk.
type TLSPaths struct {
	CertFile   string `json:"cert_file"`
	KeyFile    string `json:"key_file"`
	IssuerFile string `json:"issuer_file"`
}

// IngressMTLS configures mutual TLS validation.
type IngressMTLS struct {
	Enabled bool   `json:"enabled"`
	CAPath  string `json:"ca_path"`
}

// GatewayConfig exposes the gateway listen port for port collision checks.
type GatewayConfig struct {
	ListenPort uint32 `json:"listen_port"`
}

// FallbackConfig provides a minimal ingress definition when etcd is unavailable.
type FallbackConfig struct {
	Enabled  bool              `json:"enabled"`
	Ingress  *FallbackIngress  `json:"ingress"`
	Clusters []FallbackCluster `json:"clusters"`
}

// FallbackIngress describes listener/route metadata for fallback mode.
type FallbackIngress struct {
	ListenerHost       string          `json:"listener_host"`
	HTTPSPort          uint32          `json:"https_port"`
	HTTPPort           uint32          `json:"http_port"`
	EnableHTTPRedirect *bool           `json:"enable_http_redirect"`
	TLS                TLSPaths        `json:"tls"`
	Routes             []FallbackRoute `json:"routes"`
}

// FallbackRoute mirrors builder.Route fields for ingress fallback.
type FallbackRoute struct {
	Prefix      string `json:"prefix"`
	Cluster     string `json:"cluster"`
	HostRewrite string `json:"host_rewrite,omitempty"`
	Authority   string `json:"authority,omitempty"`
	Domains     string `json:"domains,omitempty"`
}

// FallbackCluster describes an upstream cluster and its endpoints in fallback mode.
type FallbackCluster struct {
	Name      string             `json:"name"`
	Endpoints []FallbackEndpoint `json:"endpoints"`
	TLS       TLSPaths           `json:"tls,omitempty"`
}

// FallbackEndpoint describes a single cluster endpoint.
type FallbackEndpoint struct {
	Host     string `json:"host"`
	Port     uint32 `json:"port"`
	Priority uint32 `json:"priority,omitempty"`
}

func boolPtr(v bool) *bool {
	return &v
}

func boolValue(ptr *bool, defaults bool) bool {
	if ptr == nil {
		return defaults
	}
	return *ptr
}

func (cfg *XDSConfig) ingressRedirectEnabled() bool {
	return boolValue(cfg.Ingress.EnableHTTPRedirect, true)
}

func (cfg *XDSConfig) gatewayPort() uint32 {
	if cfg == nil || cfg.Gateway.ListenPort == 0 {
		return defaultGatewayListenPort
	}
	return cfg.Gateway.ListenPort
}

func (cfg *XDSConfig) normalize() {
	cfg.Ingress.normalize()
	if cfg.Gateway.ListenPort == 0 {
		cfg.Gateway.ListenPort = defaultGatewayListenPort
	}
	if cfg.Fallback != nil {
		cfg.Fallback.normalize()
	}
}

func (ic *IngressConfig) normalize() {
	if ic.HTTPPort == 0 {
		ic.HTTPPort = defaultIngressHTTPPort
	}
	if ic.HTTPSPort == 0 {
		ic.HTTPSPort = defaultIngressHTTPSPort
	}
	if ic.EnableHTTPRedirect == nil {
		ic.EnableHTTPRedirect = boolPtr(true)
	}
	if ic.TLS.CertChainPath == "" {
		ic.TLS.CertChainPath = defaultIngressCertPath
	}
	if ic.TLS.PrivateKeyPath == "" {
		ic.TLS.PrivateKeyPath = defaultIngressKeyPath
	}
	if ic.TLS.Enabled == nil {
		ic.TLS.Enabled = boolPtr(true)
	}
	if ic.MTLS.CAPath == "" {
		ic.MTLS.CAPath = defaultIngressProvisionCA
	}
}

func (ic *IngressConfig) tlsEnabled() bool {
	if ic.TLS.Enabled == nil {
		return true
	}
	return *ic.TLS.Enabled
}

func (fb *FallbackConfig) normalize() {
	if fb == nil || fb.Ingress == nil {
		return
	}
	if fb.Ingress.HTTPPort == 0 {
		fb.Ingress.HTTPPort = defaultIngressHTTPPort
	}
	if fb.Ingress.HTTPSPort == 0 {
		fb.Ingress.HTTPSPort = defaultIngressHTTPSPort
	}
	if fb.Ingress.EnableHTTPRedirect == nil {
		fb.Ingress.EnableHTTPRedirect = boolPtr(true)
	}
	if fb.Ingress.ListenerHost == "" {
		fb.Ingress.ListenerHost = "0.0.0.0"
	}
}
