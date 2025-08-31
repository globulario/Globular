package globule

import (
	"context"
	"crypto/tls"
	"fmt"
	"path/filepath"
	"time"

	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
)

// BootstrapHTTPOnly prepares folders/config so the HTTP server can run.
// It does NOT start servers – your main() already does that via Supervisor.
func (g *Globule) BootstrapHTTPOnly() error {
	// Prepare directories, load or create config.json, set webroot/data/creds, etc.
	if err := g.InitFS(); err != nil {
		return fmt.Errorf("bootstrap: initDirectories: %w", err)
	}
	// Persist whatever we have so far (safe no-op if unchanged).
	if err := g.SaveConfig(); err != nil {
		return fmt.Errorf("bootstrap: saveConfig: %w", err)
	}
	return nil
}

// EnsureTLS validates/refreshes certs if you decide to run HTTPS.
// Safe to call even if you only run HTTP (it’ll just return nil quickly).
func (g *Globule) EnsureTLS(ctx context.Context) error {
	// Only do work when we actually want HTTPS
	if g.Protocol != "https" {
		return nil
	}

	// Validate existing cert (if present). If expired/missing, regenerate.
	if Utility.Exists(g.creds + "/" + g.Certificate) {
		if err := g.validateAndMaybeCleanupCert(); err != nil {
			return fmt.Errorf("ensureTLS: %w", err)
		}
		return nil
	}

	// No cert available -> generate local service certs and try ACME flow
	if err := g.issueCertificates(ctx); err != nil {
		return fmt.Errorf("ensureTLS: issueCertificates: %w", err)
	}
	return nil
}

func (g *Globule) generateLocalServiceCerts() error {
	// Fresh internal service certs (CA/client/server) — required before ACME CSR
	if err := security.GenerateServicesCertificates(
		g.CertPassword, g.CertExpirationDelay, g.localDomain(), g.creds,
		g.Country, g.State, g.City, g.Organization, g.AlternateDomains,
	); err != nil {
		return fmt.Errorf("generate service certs: %w", err)
	}
	return nil
}

// validateAndMaybeCleanupCert checks if the current TLS cert is still valid.
// If expired or invalid, it deletes all cert files in creds dir and clears g.Certificate.
func (g *Globule) validateAndMaybeCleanupCert() error {
	if g.Certificate == "" {
		return nil
	}
	certFile := filepath.Join(g.creds, g.Certificate)
	keyFile := filepath.Join(g.creds, "server.pem")

	if Utility.Exists(certFile) {
		if err := security.ValidateCertificateExpiration(certFile, keyFile); err == nil {
			return nil // cert is valid
		}
		// expired or invalid -> wipe creds dir
		if err := Utility.RemoveDirContents(g.creds); err != nil {
			return fmt.Errorf("cleanup creds dir: %w", err)
		}
		g.Certificate = ""
	}
	return nil
}

// issueCertificates performs your existing “generate local, then obtain ACME” flow.
// It stores paths on g and saves config. Kept small here by delegating to your existing helpers.
func (g *Globule) issueCertificates(ctx context.Context) error {
	// Make local keypair/CSR for services (already implemented in your code)
	if err := g.generateLocalServiceCerts(); err != nil {
		return err
	}

	// If you have a DNS service configured, register/update host records
	if err := g.maybeStartDNSAndRegister(ctx); err != nil {
		// Non-fatal if you want to keep HTTP running, but we’ll bubble it up
		return err
	}

	// Obtain ACME certificate for the CSR (your lego flow)
	cctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	if err := g.obtainCertificateForCSR(cctx); err != nil {
		return err
	}

	return g.SaveConfig()
}

// TLSConfig is optional; expose it if you want Supervisor to use a tls.Config.
func (g *Globule) TLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: g.localDomain(),
	}
}
