package globule

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/globulario/Globular/internal/logsink"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/dns/dns_client"
	"github.com/globulario/services/golang/log/logpb"
	"github.com/globulario/services/golang/process"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
	"github.com/txn2/txeh"
)

// maybeStartDNSAndRegister starts dns.DnsService if configured locally, then registers A/AAAA/MX, etc.
// starts it (if not running), waits for readiness, then updates DNS records.
func (g *Globule) maybeStartDNSAndRegister(ctx context.Context) error {
	const svcName = "dns.DnsService"

	// 1) Ensure desired exists for DNS (describe/merge/allocate like the main boot)
	desc, bin, err := g.describeServiceByName(svcName, 5*time.Second)
	if err != nil {
		g.log.Warn("dns describe failed (service binary not found yet?)", "err", err)
		// best-effort: still attempt register using external DNS if your registerIPToDNS() supports it
		return g.registerIPToDNS()
	}

	alloc, err := config.NewDefaultPortAllocator()
	if err != nil {
		return err
	}

	// Merge or create desired config for DNS service
	desired, err := g.mergeOrCreateDesired(desc, alloc)
	if err != nil {
		return fmt.Errorf("dns desired creation failed: %w", err)
	}

	desired["Path"] = bin
	desired["Domain"] = g.Domain
	desired["Address"] = g.localIPAddress
	desired["Mac"] = g.Mac

	if err := config.SaveServiceConfiguration(desired); err != nil {

		return fmt.Errorf("dns desired save failed: %w", err)
	}

	// 2) If already running, skip start
	id := Utility.ToString(desired["Id"])
	if id == "" {
		return fmt.Errorf("dns desired missing Id")
	}
	state := ""
	if cfg, _ := config.GetServiceConfigurationById(id); cfg != nil {
		state = Utility.ToString(cfg["State"])
	}
	if state == "running" {
		// already up; proceed to register
		return g.registerIPToDNS()
	}

	// 3) Start DNS with logs and wait readiness
	address, _ := config.GetAddress()
	name := svcName
	port := Utility.ToInt(desired["Port"])
	proxy := Utility.ToInt(desired["Proxy"])
	outW := logsink.NewServiceLogWriter(address, name, "sa", "/"+name+"/stdout", logpb.LogLevel_INFO_MESSAGE, os.Stdout)
	errW := logsink.NewServiceLogWriter(address, name, "sa", "/"+name+"/stderr", logpb.LogLevel_ERROR_MESSAGE, os.Stderr)

	g.log.Info("starting dns", "port", port, "proxy", proxy, "path", desired["Path"])
	pid, err := process.StartServiceProcessWithWriters(desired, port, outW, errW)
	if err != nil {
		_ = config.PutRuntime(id, map[string]any{"Process": -1, "State": "failed", "LastError": err.Error()})
		g.log.Warn("dns start failed", "err", err)
		// Still attempt external registration so HTTPS bootstrap can proceed if possible.
		return g.registerIPToDNS()
	}
	_ = config.PutRuntime(id, map[string]any{"Process": pid, "State": "starting", "LastError": ""})

	// set pid in desired for proxy use
	desired["Process"] = pid

	// optional: start grpc-web proxy
	if !g.UseEnvoy {
		if _, err := process.StartServiceProxyProcess(desired, config.GetLocalCertificateAuthorityBundle(), config.GetLocalCertificate()); err != nil {
			g.log.Warn("dns proxy start failed", "err", err)
		}
	}

	// bounded readiness wait
	addr, _ := config.GetHostname()
	addr += ":" + Utility.ToString(port)

	ok := g.waitServiceReady(name, addr, 15*time.Second)
	if !ok {
		_ = config.PutRuntime(id, map[string]any{"State": "failed", "LastError": "dns startup timeout"})
		g.log.Warn("dns failed to become ready")
		// Attempt external register anyway.
		return g.registerIPToDNS()
	}
	_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})

	// 4) Now we can safely register IP → DNS
	return g.registerIPToDNS()
}

// setHost maps an IPv4 to a hostname in the system hosts file.
// - If the current mapping is local (127.* / 10.* / 192.168.* / etc.), we refuse to overwrite it with a non-local IP.
// - Handles "localhost" specially (forces 127.0.0.1).
func (g *Globule) setHost(ipv4, address string) error {
	if strings.HasSuffix(address, ".localhost") {
		return nil
	}
	if address == "localhost" {
		ipv4 = "127.0.0.1"
	}
	if ipv4 == "" {
		return errors.New("setHost: empty ipv4")
	}
	if address == "" {
		return errors.New("setHost: empty address")
	}

	h, err := txeh.NewHostsDefault()
	if err != nil {
		return err
	}

	exists, prev, _ := h.HostAddressLookup(address, txeh.IPFamilyV4)
	if exists && Utility.IsLocal(prev) && !Utility.IsLocal(ipv4) {
		// Keep existing local mapping; admin can override manually if desired.
		return nil
	}

	h.AddHost(ipv4, address)
	return h.Save()
}

// registerIPToDNS reproduces your previous behavior, but keeps errors informative and short.
func (g *Globule) registerIPToDNS() error {

	if g.DNS == "" {
		return nil
	}

	domain := g.Domain
	localFQDN := g.Name + "." + g.Domain

	// Keep resolv.conf pointed at public recursors (matches original behavior).
	const resolvHdr = "# Generated by Globular at startup. Original saved as resolv.conf_\n"
	resolvConf := resolvHdr + "nameserver 8.8.8.8\nnameserver 1.1.1.1\n"
	if Utility.Exists("/etc/resolv.conf") {
		_ = Utility.MoveFile("/etc/resolv.conf", "/etc/resolv.conf_")
		_ = Utility.WriteStringToFile("/etc/resolv.conf", resolvConf)
	}

	c, err := dns_client.NewDnsService_Client(g.DNS, "dns.DnsService")
	if err != nil {
		return fmt.Errorf("dns client: %w", err)
	}
	defer c.Close()

	// If our DNS authority is local, wait until it reports "running" (original code did this).
	if g.DNS == localFQDN {
		const maxTry = 20
		for try := 0; try < maxTry; try++ {
			cfg, err := config.GetServiceConfigurationById(c.GetId())
			if err == nil && cfg["State"] == "running" {
				break
			}

			time.Sleep(time.Second)
			if try == maxTry-1 {
				return fmt.Errorf("dns server not running")
			}
		}
	}

	tk, err := security.GenerateToken(g.SessionTimeout, c.GetMac(), "sa", "", g.AdminEmail, g.Domain)
	if err != nil {
		return fmt.Errorf("token: %w", err)
	}

	ipv4 := Utility.MyIP()
	if ipv4 == "" {
		return errors.New("no public IPv4")
	}
	ipv6, _ := Utility.MyIPv6()

	// first of all I will set the domains for the dns server...
	domains := []string{}

	for _, v := range g.AlternateDomains {
		alt := strings.TrimPrefix(fmt.Sprint(v), "*.")
		domains = append(domains, alt)
	}

	if !Utility.Contains(domains, domain) {
		domains = append(domains, domain)
	}

	err = c.SetDomains(tk, domains)
	if err != nil {
		return fmt.Errorf("set domains: %w", err)
	}

	// --- Primary host records (A/AAAA) ---
	_ = c.RemoveA(tk, localFQDN)
	if _, err := c.SetA(tk, localFQDN, ipv4, 60); err != nil {

		return fmt.Errorf("set A %s: %w", localFQDN, err)
	}
	if ipv6 != "" {
		if _, err := c.SetAAAA(tk, localFQDN, ipv6, 60); err != nil {
			return fmt.Errorf("set AAAA %s: %w", localFQDN, err)
		}
	}

	// --- Alternates (A/AAAA), NS and SOA so public resolvers can find the zone ---
	for _, v := range g.AlternateDomains {
		alt := strings.TrimPrefix(fmt.Sprint(v), "*.")
		if _, err := c.SetA(tk, alt, ipv4, 60); err != nil {
			return fmt.Errorf("set A %s: %w", alt, err)
		}
		if ipv6 != "" {
			if _, err := c.SetAAAA(tk, alt, ipv6, 60); err != nil {
				return fmt.Errorf("set AAAA %s: %w", alt, err)
			}
		}

		// NS glue within the zone
		for _, rawNS := range g.NS {
			ns := fmt.Sprint(rawNS)
			if _, err := c.SetA(tk, ns, ipv4, 60); err != nil {
				return fmt.Errorf("set A NS %s: %w", ns, err)
			}
			if ipv6 != "" {
				if _, err := c.SetAAAA(tk, ns, ipv6, 60); err != nil {
					return fmt.Errorf("set AAAA NS %s: %w", ns, err)
				}
			}
		}

		// Delegate NS at the zone apex (alt must be FQDN with trailing dot for NS/SOA APIs)
		altDot := alt
		if !strings.HasSuffix(altDot, ".") {
			altDot += "."
		}
		for _, rawNS := range g.NS {
			ns := fmt.Sprint(rawNS)
			if !strings.HasSuffix(ns, ".") {
				ns += "."
			}
			if err := c.SetNs(tk, altDot, ns, 60); err != nil {
				return fmt.Errorf("set NS %s -> %s: %w", altDot, ns, err)
			}
		}

		// SOA (primary NS + RNAME)
		primaryNS := ""
		if len(g.NS) > 0 {
			primaryNS = fmt.Sprint(g.NS[0])
			if !strings.HasSuffix(primaryNS, ".") {
				primaryNS += "."
			}
		} else {
			primaryNS = "ns1." + domain + "."
		}

		email := g.AdminEmail
		if email == "" || !strings.Contains(email, "@") {
			email = "admin@" + domain
		}
		if !strings.HasSuffix(email, ".") {
			email += "."
		}
		const (
			serial  = uint32(1)
			refresh = uint32(86400)
			retry   = uint32(7200)
			expire  = uint32(4000000)
			ttl     = uint32(11200)
		)
		if err := c.SetSoa(tk, altDot, primaryNS, email, serial, refresh, retry, expire, ttl, ttl); err != nil {
			return fmt.Errorf("set SOA %s: %w", altDot, err)
		}

		// CAA allows Let's Encrypt
		if err := c.SetCaa(tk, alt+".", 0, "issue", "letsencrypt.org", 60); err != nil {
			return fmt.Errorf("set CAA %s: %w", alt, err)
		}
	}

	// --- Mail, SPF, DMARC, MTA-STS ---
	mailHost := "mail." + domain
	if _, err := c.SetA(tk, mailHost, ipv4, 60); err != nil {
		return fmt.Errorf("set A %s: %w", mailHost, err)
	}
	if ipv6 != "" {
		if _, err := c.SetAAAA(tk, mailHost, ipv6, 60); err != nil {
			return fmt.Errorf("set AAAA %s: %w", mailHost, err)
		}
	}

	if err := c.SetMx(tk, domain, 10, mailHost, 60); err != nil {
		return fmt.Errorf("set MX %s: %w", domain, err)
	}

	_ = c.RemoveText(tk, domain+".")
	spf := fmt.Sprintf(`v=spf1 mx ip4:%s include:_spf.google.com ~all`, ipv4)
	if err := c.SetText(tk, domain+".", []string{spf}, 60); err != nil {
		return fmt.Errorf("set SPF TXT: %w", err)
	}

	_ = c.RemoveText(tk, "_dmarc."+domain+".")
	dmarc := fmt.Sprintf(`v=DMARC1;p=quarantine;rua=mailto:%s;ruf=mailto:%s;adkim=r;aspf=r;pct=100`, g.AdminEmail, g.AdminEmail)
	if err := c.SetText(tk, "_dmarc."+domain+".", []string{dmarc}, 60); err != nil {
		return fmt.Errorf("set DMARC TXT: %w", err)
	}

	// MTA-STS file and pointers
	mtaPath := filepath.Join(g.configDir, "tls", localFQDN, "mta-sts.txt")
	if !Utility.Exists(mtaPath) {
		policy := fmt.Sprintf("version: STSv1\nmode: enforce\nmx: %s\nttl: 86400\n", domain)
		_ = os.WriteFile(mtaPath, []byte(policy), 0600)
	}
	if _, err := c.SetA(tk, "mta-sts."+domain, ipv4, 60); err != nil {
		return fmt.Errorf("set A mta-sts: %w", err)
	}
	_ = c.RemoveText(tk, "_mta-sts."+domain+".")
	if err := c.SetText(tk, "_mta-sts."+domain+".", []string{"v=STSv1; id=globular;"}, 60); err != nil {
		return fmt.Errorf("set _mta-sts TXT: %w", err)
	}

	// Optional: vendor API updates (GoDaddy, etc.)
	for _, raw := range g.DNSUpdateIPInfos {
		m := raw.(map[string]interface{})
		setA, key, secret := m["SetA"].(string), m["Key"].(string), m["Secret"].(string)
		body := `[{"data":"` + ipv4 + `"}]`
		req, _ := http.NewRequest(http.MethodPut, setA, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "sso-key "+key+":"+secret)
		_, _ = (&http.Client{Timeout: 10 * time.Second}).Do(req)
	}

	// Keep /etc/hosts in sync (don’t override local -> public).
	return g.setHost(config.GetLocalIP(), localFQDN)
}
