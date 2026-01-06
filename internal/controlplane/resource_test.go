package controlplane

import (
	"testing"

	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
)

func TestMakeDownstreamTLSWithoutCA(t *testing.T) {
	ts := makeDownstreamTLS("cert.crt", "key.key", "")
	ctx := &tls_v3.DownstreamTlsContext{}
	if err := ts.GetTypedConfig().UnmarshalTo(ctx); err != nil {
		t.Fatalf("unmarshal TLS context: %v", err)
	}
	if ctx.GetCommonTlsContext().GetValidationContextType() != nil {
		t.Fatalf("expected validation context to be nil when CA is missing")
	}
}

func TestMakeDownstreamTLSWithCA(t *testing.T) {
	ts := makeDownstreamTLS("cert.crt", "key.key", "/etc/globular/certs/ca.pem")
	ctx := &tls_v3.DownstreamTlsContext{}
	if err := ts.GetTypedConfig().UnmarshalTo(ctx); err != nil {
		t.Fatalf("unmarshal TLS context: %v", err)
	}
	if ctx.GetCommonTlsContext().GetValidationContextType() == nil {
		t.Fatalf("expected validation context to be set when CA is provided")
	}
}
