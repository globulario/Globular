package globule

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/globulario/services/golang/dns/dns_client"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
	"github.com/go-acme/lego/v4/challenge/dns01"
)

type registrationResource = registration.Resource // alias for decoupling

// DNS provider used by lego (DNS-01) backed by your dns.DnsService.
type dnsProvider struct {
	apiToken string
	addr     string
	globule  *Globule
}

// Present creates a DNS TXT record to fulfill the ACME DNS-01 challenge for the specified domain.
// It uses the provided key authorization to generate the record value, connects to the DNS service,
// and sets the TXT record with a 30-second TTL. Returns an error if the DNS address is not configured,
// if the DNS client or token generation fails, or if setting the TXT record fails.
//
// Parameters:
//
//	domain  - The domain for which to present the DNS challenge.
//	_       - Unused parameter (typically for compatibility with interface).
//	keyAuth - The ACME key authorization string.
//
// Returns:
//
//	error - Non-nil if any step fails.
//
// Present creates the TXT record for DNS-01.
func (p *dnsProvider) Present(domain, _ string, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	if p.addr == "" {
		return errors.New("dns address not configured")
	}
	c, err := dns_client.NewDnsService_Client(p.addr, "dns.DnsService")
	if err != nil {
		return fmt.Errorf("dns client: %w", err)
	}
	defer c.Close()

	tk, err := security.GenerateToken(p.globule.SessionTimeout, c.GetMac(), "sa", "", p.globule.AdminEmail, p.globule.Domain)
	if err != nil {
		return fmt.Errorf("token: %w", err)
	}

	// Use a 60s TTL like the rest of your records.
	if err := c.SetText(tk, fqdn, []string{value}, 60); err != nil {
		return fmt.Errorf("set TXT %q: %w", fqdn, err)
	}
	return nil
}

// CleanUp removes the DNS-01 challenge TXT record for the specified domain.
// It connects to the DNS service using the configured address, generates an authentication token,
// and requests the removal of the TXT record associated with the provided key authorization.
// Returns an error if any step fails.
func (p *dnsProvider) CleanUp(domain, _ string, keyAuth string) error {
	fqdn, _ := dns01.GetRecord(domain, keyAuth)
	if p.addr == "" {
		return nil
	}
	c, err := dns_client.NewDnsService_Client(p.addr, "dns.DnsService")
	if err != nil {
		return fmt.Errorf("dns client: %w", err)
	}
	defer c.Close()

	tk, err := security.GenerateToken(p.globule.SessionTimeout, c.GetMac(), "sa", "", p.globule.AdminEmail, p.globule.Domain)
	if err != nil {
		return fmt.Errorf("token: %w", err)
	}

	if err := c.RemoveText(tk, fqdn); err != nil {
		return fmt.Errorf("remove TXT %q: %w", fqdn, err)
	}
	return nil
}

func (g *Globule) obtainCertificateForCSR(ctx context.Context) error {
	cfg := lego.NewConfig(g)
	cfg.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("lego client: %w", err)
	}

	if g.DNS != "" {
		dc, err := dns_client.NewDnsService_Client(g.DNS, "dns.DnsService")
		if err != nil {
			return fmt.Errorf("dns client: %w", err)
		}
		defer dc.Close()

		tk, err := security.GenerateToken(g.SessionTimeout, dc.GetMac(), "sa", "", g.AdminEmail, g.Domain)
		if err != nil {
			return fmt.Errorf("token: %w", err)
		}

		prov := &dnsProvider{apiToken: tk, addr: g.DNS, globule: g}
		if err := client.Challenge.SetDNS01Provider(prov); err != nil {
			return fmt.Errorf("lego set dns01: %w", err)
		}
	} else {
		prov := http01.NewProviderServer("", strconv.Itoa(g.PortHTTP))
		if err := client.Challenge.SetHTTP01Provider(prov); err != nil {
			return fmt.Errorf("lego set http01: %w", err)
		}
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return fmt.Errorf("lego register: %w", err)
	}
	g.registration = reg

	csrPEM, err := os.ReadFile(filepath.Join(g.creds, "server.csr"))
	if err != nil {
		return fmt.Errorf("read CSR: %w", err)
	}
	block, _ := pem.Decode(csrPEM)
	if block == nil || (block.Type != "CERTIFICATE REQUEST" && block.Type != "NEW CERTIFICATE REQUEST") {
		return errors.New("server.csr: invalid PEM")
	}
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse CSR: %w", err)
	}

	res, err := client.Certificate.ObtainForCSR(*csr, true)
	if err != nil {
		return fmt.Errorf("lego obtain: %w", err)
	}

	g.CertURL = res.CertURL
	g.CertStableURL = res.CertStableURL
	g.Certificate = g.Domain + ".crt"
	g.CertificateAuthorityBundle = g.Domain + ".issuer.crt"

	leafPath := filepath.Join(g.creds, g.Certificate)
	issuerPath := filepath.Join(g.creds, g.CertificateAuthorityBundle)
	fullPath := filepath.Join(g.creds, "fullchain.pem")
	chainPath := filepath.Join(g.creds, "chain.pem") // optional (issuer only)

	// Write leaf & issuer
	if err := os.WriteFile(leafPath, res.Certificate, 0o400); err != nil {
		return fmt.Errorf("write cert: %w", err)
	}
	if len(res.IssuerCertificate) > 0 {
		if err := os.WriteFile(issuerPath, res.IssuerCertificate, 0o400); err != nil {
			return fmt.Errorf("write issuer: %w", err)
		}
		// Write chain.pem (issuer only) for stacks that expect it
		_ = os.WriteFile(chainPath, res.IssuerCertificate, 0o444)
	}

	// Always try to write fullchain.pem (leaf + issuer)
	full := append([]byte{}, res.Certificate...)
	if len(res.IssuerCertificate) > 0 {
		full = append(full, res.IssuerCertificate...)
	}
	if err := os.WriteFile(fullPath, full, 0o444); err != nil {
		return fmt.Errorf("write fullchain.pem: %w", err)
	}

	return nil
}

// GetEmail / GetRegistration / GetPrivateKey — lego user hooks
func (g *Globule) GetEmail() string                        { return g.AdminEmail }
func (g *Globule) GetRegistration() *registration.Resource { return g.registration }
func (g *Globule) GetPrivateKey() crypto.PrivateKey {
	b, err := os.ReadFile(filepath.Join(g.creds, "client.pem"))
	if err != nil {
		return nil
	}
	p, _ := pem.Decode(b)
	k, err := x509.ParsePKCS8PrivateKey(p.Bytes)
	if err != nil {
		return nil
	}
	return k
}

// ensureAccountKeyAndCSR creates a lego account key (client.pem) if missing,
// and a server private key + CSR (server.key, server.csr) for g.Domain and SANs.
func (g *Globule) ensureAccountKeyAndCSR() error {
	// Ensure creds dir exists
	if err := Utility.CreateDirIfNotExist(g.creds); err != nil {
		return err
	}

	// 1) Account key for ACME user: client.pem (PKCS#8)
	acct := filepath.Join(g.creds, "client.pem")
	if !Utility.Exists(acct) {
		// Generate client.key then convert to client.pem
		if err := security.GenerateClientPrivateKey(g.creds, g.CertPassword); err != nil {
			return fmt.Errorf("generate client.key: %w", err)
		}
		if err := security.KeyToPem("client", g.creds, g.CertPassword); err != nil {
			return fmt.Errorf("write client.pem: %w", err)
		}
	}

	// 2) Server key + CSR
	sk := filepath.Join(g.creds, "server.key")
	csr := filepath.Join(g.creds, "server.csr")
	if !Utility.Exists(sk) || !Utility.Exists(csr) {
		// Generate server.key
		if err := security.GenerateSeverPrivateKey(g.creds, g.CertPassword); err != nil {
			return fmt.Errorf("generate server.key: %w", err)
		}

		// Build SANs: primary + alternates (GenerateSanConfig reads these)
		var sans []string
		if d := strings.TrimSpace(g.Domain); d != "" {
			sans = append(sans, d)
		}
		for _, v := range g.AlternateDomains {
			if ad := strings.TrimSpace(fmt.Sprint(v)); ad != "" {
				sans = append(sans, ad)
			}
		}

		// Write san.conf then create server.csr
		if err := security.GenerateSanConfig(g.Domain, g.creds, g.Country, g.State, g.City, g.Organization, sans); err != nil {
			return fmt.Errorf("generate san.conf: %w", err)
		}
		if err := security.GenerateServerCertificateSigningRequest(g.creds, g.CertPassword, g.Domain); err != nil {
			return fmt.Errorf("generate server.csr: %w", err)
		}
	}

	return nil
}

func exists(p string) bool { _, err := os.Stat(p); return err == nil }
