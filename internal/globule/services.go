package globule

import (
	"context"
	"errors"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/process"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
)

// refreshTokenPeriodically keeps the local SA token fresh for microservices.
func refreshTokenPeriodically(ctx context.Context, g *Globule) {
	interval := time.Duration(g.SessionTimeout)*time.Minute - 10*time.Second
	if interval <= 0 {
		interval = time.Minute
	}
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			_ = security.SetLocalToken(g.Mac, g.Domain, "sa", "sa", g.AdminEmail, g.SessionTimeout)
		}
	}
}

// Public entry points kept from original behavior (no HTTP server logic here).
func (g *Globule) startServices(ctx context.Context) error {

	// peer keys & local token
	if err := security.GeneratePeerKeys(g.Mac); err != nil {
		return err
	}
	if err := security.SetLocalToken(g.Mac, g.Domain, "sa", "sa", g.AdminEmail, g.SessionTimeout); err != nil {
		return err
	}

	services, err := config.GetOrderedServicesConfigurations()
	if err != nil {
		return err
	}

	go refreshTokenPeriodically(ctx, g)

	start := Utility.ToInt(strings.Split(g.PortsRange, "-")[0])
	end := Utility.ToInt(strings.Split(g.PortsRange, "-")[1])

	for i := range services {
		if start >= end {
			return errors.New("no more available ports")
		}
		s := services[i]
		s["State"] = "starting"
		s["ProxyProcess"] = -1

		port := start + (i * 2)
		name := s["Name"].(string)

		g.log.Info("start service", "name", name, "port", port, "proxy", port+1)
		pid, err := process.StartServiceProcess(s, port)
		if err != nil {
			g.log.Warn("service start failed", "name", name, "err", err)
			continue
		}
		s["Process"] = pid
		s["ProxyProcess"] = -1
		if _, err := process.StartServiceProxyProcess(s, config.GetLocalCertificateAuthorityBundle(), config.GetLocalCertificate()); err != nil {
			g.log.Warn("proxy start failed", "name", name, "err", err)
		}
	}

	// Wait until services become running (bounded)
	all := false
	for tries := 20; !all && tries > 0; tries-- {
		all = true
		lst, err := config.GetServicesConfigurations()
		if err != nil {
			return err
		}
		for _, s := range lst {
			if s["State"].(string) != "running" {
				all = false
				time.Sleep(1 * time.Second)
				break
			}
		}
	}
	if !all {
		return errors.New("some services failed to start")
	}

	// Register DNS
	for range 20 {
		if err := g.registerIPToDNS(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (g *Globule) stopServices() error {
	svcs, err := config.GetServicesConfigurations()
	if err != nil {
		return err
	}
	for _, s := range svcs {
		pid := Utility.ToInt(s["Process"])
		proxy := Utility.ToInt(s["ProxyProcess"])
		s["State"] = "killed"
		s["ProxyProcess"] = -1
		if err := config.SaveServiceConfiguration(s); err == nil {
			if pid > 0 {
				if p, err := os.FindProcess(pid); err == nil {
					_ = p.Signal(syscall.SIGTERM)
				}
			}
			if proxy > 0 {
				if p, err := os.FindProcess(proxy); err == nil {
					_ = p.Signal(syscall.SIGTERM)
				}
			}
		}
	}
	return nil
}
