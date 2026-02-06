package watchers

import (
	"testing"

	"github.com/globulario/Globular/internal/xds/builder"
)

// TestConfigureModeResources_IngressMode tests ingress mode configuration
func TestConfigureModeResources_IngressMode(t *testing.T) {
	w := &Watcher{}

	ingressSpec := &IngressSpec{
		Listener: builder.Listener{
			Port:     443,
			Host:     "", // Should default to 0.0.0.0
			CertFile: "/path/to/cert.pem",
			KeyFile:  "/path/to/key.pem",
		},
		Clusters: []builder.Cluster{
			{Name: "test_cluster"},
		},
		Routes: []builder.Route{
			{Prefix: "/", Cluster: "test_cluster"},
		},
		HTTPPort:           80,
		EnableHTTPRedirect: true,
		GatewayPort:        8080,
	}

	result, err := w.configureModeResources(ingressSpec, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify listener defaults applied
	if result.listener.Host != "0.0.0.0" {
		t.Errorf("expected default Host 0.0.0.0, got: %s", result.listener.Host)
	}
	if result.listener.RouteName != defaultRouteName {
		t.Errorf("expected default RouteName %s, got: %s", defaultRouteName, result.listener.RouteName)
	}
	if result.listener.Name == "" {
		t.Error("expected listener Name to be set")
	}

	// Verify resources passed through
	if len(result.clusters) != 1 {
		t.Errorf("expected 1 cluster, got: %d", len(result.clusters))
	}
	if len(result.routes) != 1 {
		t.Errorf("expected 1 route, got: %d", len(result.routes))
	}

	// Verify HTTP settings
	if result.ingressHTTPPort != 80 {
		t.Errorf("expected ingressHTTPPort 80, got: %d", result.ingressHTTPPort)
	}
	if !result.enableHTTPRedirect {
		t.Error("expected enableHTTPRedirect true")
	}

	t.Log("✓ Ingress mode configuration applied correctly")
}

// TestConfigureModeResources_IngressModeWithDefaults tests that all listener defaults are applied
func TestConfigureModeResources_IngressModeWithDefaults(t *testing.T) {
	w := &Watcher{}

	ingressSpec := &IngressSpec{
		Listener: builder.Listener{
			Port: 443,
			// All optional fields empty - should get defaults
		},
	}

	result, err := w.configureModeResources(ingressSpec, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all defaults applied
	if result.listener.Host != "0.0.0.0" {
		t.Errorf("expected default Host 0.0.0.0, got: %s", result.listener.Host)
	}
	if result.listener.RouteName != defaultRouteName {
		t.Errorf("expected default RouteName %s, got: %s", defaultRouteName, result.listener.RouteName)
	}
	if result.listener.Name == "" {
		t.Error("expected listener Name to be generated")
	}
	// Verify Name contains port
	expectedName := "ingress_listener_443"
	if result.listener.Name != expectedName {
		t.Errorf("expected listener Name %s, got: %s", expectedName, result.listener.Name)
	}

	t.Log("✓ Ingress mode defaults applied correctly")
}

// TestConfigureModeResources_IngressModePreservesValues tests that explicit values are preserved
func TestConfigureModeResources_IngressModePreservesValues(t *testing.T) {
	w := &Watcher{}

	ingressSpec := &IngressSpec{
		Listener: builder.Listener{
			Port:      8443,
			Host:      "192.168.1.1",     // Explicit host - should be preserved
			RouteName: "custom_routes",   // Explicit route name - should be preserved
			Name:      "custom_listener", // Explicit name - should be preserved
		},
	}

	result, err := w.configureModeResources(ingressSpec, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify explicit values preserved
	if result.listener.Host != "192.168.1.1" {
		t.Errorf("expected Host 192.168.1.1, got: %s", result.listener.Host)
	}
	if result.listener.RouteName != "custom_routes" {
		t.Errorf("expected RouteName custom_routes, got: %s", result.listener.RouteName)
	}
	if result.listener.Name != "custom_listener" {
		t.Errorf("expected Name custom_listener, got: %s", result.listener.Name)
	}

	t.Log("✓ Ingress mode explicit values preserved")
}

// TestConfigureModeResources_LegacyMode tests legacy mode fallback
func TestConfigureModeResources_LegacyMode(t *testing.T) {
	// This test requires mocking buildLegacyGatewayResources
	// For now, we'll skip it and test via integration tests
	t.Skip("Requires mocking buildLegacyGatewayResources - covered by integration tests")
}

// TestConfigureModeResources_FinalNormalization tests that RouteName is always set
func TestConfigureModeResources_FinalNormalization(t *testing.T) {
	w := &Watcher{}

	// Test with ingress spec that has no RouteName
	ingressSpec := &IngressSpec{
		Listener: builder.Listener{
			Port:      443,
			Host:      "0.0.0.0",
			RouteName: "", // Empty - should get default
		},
	}

	result, err := w.configureModeResources(ingressSpec, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Final normalization should always set RouteName
	if result.listener.RouteName != defaultRouteName {
		t.Errorf("expected RouteName %s after final normalization, got: %s",
			defaultRouteName, result.listener.RouteName)
	}

	t.Log("✓ Final normalization ensures RouteName always set")
}

// TestConfigureModeResources_ListenerNameGeneration tests listener name generation logic
func TestConfigureModeResources_ListenerNameGeneration(t *testing.T) {
	w := &Watcher{}

	testCases := []struct {
		name         string
		port         uint32
		host         string
		expectedName string
	}{
		{
			name:         "standard_https_port",
			port:         443,
			host:         "0.0.0.0",
			expectedName: "ingress_listener_443",
		},
		{
			name:         "custom_port",
			port:         8443,
			host:         "0.0.0.0",
			expectedName: "ingress_listener_8443",
		},
		{
			name:         "port_80",
			port:         80,
			host:         "0.0.0.0",
			expectedName: "ingress_listener_80",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ingressSpec := &IngressSpec{
				Listener: builder.Listener{
					Port: tc.port,
					Host: tc.host,
					Name: "", // Empty - should be generated
				},
			}

			result, err := w.configureModeResources(ingressSpec, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.listener.Name != tc.expectedName {
				t.Errorf("expected Name %s, got: %s", tc.expectedName, result.listener.Name)
			} else {
				t.Logf("✓ Generated listener name: %s", result.listener.Name)
			}
		})
	}
}

// TestConfigureModeResources_HTTPRedirectLogic tests HTTP redirect configuration
func TestConfigureModeResources_HTTPRedirectLogic(t *testing.T) {
	w := &Watcher{}

	testCases := []struct {
		name             string
		ingressHTTPPort  uint32
		enableRedirect   bool
		expectedHTTPPort uint32
		expectedRedirect bool
	}{
		{
			name:             "redirect_enabled",
			ingressHTTPPort:  80,
			enableRedirect:   true,
			expectedHTTPPort: 80,
			expectedRedirect: true,
		},
		{
			name:             "redirect_disabled",
			ingressHTTPPort:  0,
			enableRedirect:   false,
			expectedHTTPPort: 0,
			expectedRedirect: false,
		},
		{
			name:             "redirect_custom_port",
			ingressHTTPPort:  8080,
			enableRedirect:   true,
			expectedHTTPPort: 8080,
			expectedRedirect: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ingressSpec := &IngressSpec{
				Listener: builder.Listener{
					Port: 443,
				},
				HTTPPort:           tc.ingressHTTPPort,
				EnableHTTPRedirect: tc.enableRedirect,
			}

			result, err := w.configureModeResources(ingressSpec, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.ingressHTTPPort != tc.expectedHTTPPort {
				t.Errorf("expected ingressHTTPPort %d, got: %d",
					tc.expectedHTTPPort, result.ingressHTTPPort)
			}
			if result.enableHTTPRedirect != tc.expectedRedirect {
				t.Errorf("expected enableHTTPRedirect %v, got: %v",
					tc.expectedRedirect, result.enableHTTPRedirect)
			}

			t.Logf("✓ HTTP redirect configured: port=%d, enabled=%v",
				result.ingressHTTPPort, result.enableHTTPRedirect)
		})
	}
}
