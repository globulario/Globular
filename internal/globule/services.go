package globule

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/security"
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

// DescribeService exposes describeServiceByName to external callers (e.g., HTTP handlers).
func (g *Globule) DescribeService(name string, timeout time.Duration) (config.ServiceDesc, string, error) {
	desc, bin, err := g.describeServiceByName(name, timeout)
	return config.ServiceDesc(desc), bin, err
}

func (g *Globule) describeServiceByName(name string, timeout time.Duration) (serviceDesc, string, error) {
	root := config.GetServicesRoot()
	if root == "" {
		return serviceDesc{}, "", errors.New("describeServiceByName: ServicesRoot is empty")
	}

	bins, err := config.DiscoverExecutables(root)
	if err != nil {
		return serviceDesc{}, "", fmt.Errorf("describeServiceByName: discover executables: %w", err)
	}
	if len(bins) == 0 {
		return serviceDesc{}, "", fmt.Errorf("describeServiceByName: no service executables found under %s", root)
	}

	// Extract short key ("dns" from "dns.DnsService") just to cheaply skip unrelated paths.
	short := name
	if i := strings.IndexByte(name, '.'); i > 0 {
		short = name[:i]
	}
	short = strings.ToLower(short)

	// Per-binary time budget: at least 5s.
	per := timeout
	if per < 5*time.Second {
		per = 5 * time.Second
	}

	// Same env you already use in startServicesEtcd.
	env := map[string]string{
		"GLOBULAR_DOMAIN":  g.Domain,
		"GLOBULAR_ADDRESS": g.localIPAddress,
	}

	var (
		bestDesc serviceDesc
		bestPath string
		haveBest bool
		errs     []string
	)

	for _, bin := range bins {
		lpath := strings.ToLower(strings.ReplaceAll(bin, "\\", "/"))
		// quick skip if binary path doesn't even mention the short name
		if !strings.Contains(lpath, "/"+short+"/") && !strings.Contains(lpath, "/"+short+"_server/") &&
			!strings.Contains(lpath, "/"+short+"_server/") && !strings.Contains(filepath.Base(lpath), short) {
			continue
		}

		desc, derr := config.RunDescribe(bin, per, env)
		if derr != nil {
			if len(errs) < 6 {
				errs = append(errs, fmt.Sprintf("- %s: %v", bin, derr))
			}
			continue
		}

		// Normalize and verify target Name
		g.normalizeDescriptor(&desc)
		if desc.Name != name {
			// Not our target, skip.
			continue
		}

		// First match wins until we see a higher version.
		if !haveBest || compareSemver(desc.Version, bestDesc.Version) > 0 {
			bestDesc, bestPath, haveBest = desc, bin, true
			// keep going in case another binary has an even higher version
		}
	}

	if !haveBest {
		tail := ""
		if len(errs) > 0 {
			tail = "\n" + strings.Join(errs, "\n")
		}
		return serviceDesc{}, "", fmt.Errorf("describeServiceByName: no %s found under %s%s", name, root, tail)
	}

	return bestDesc, bestPath, nil
}

// compareSemver compares two semantic version strings "major.minor.patch[-pre]"
// Returns 1 if a>b, -1 if a<b, 0 if equal. Missing/invalid pieces are treated as 0.
// Pre-releases are considered lower than the same version without pre-release tag.
func compareSemver(a, b string) int {
	parse := func(s string) (maj, min, pat int, pre string) {
		// Trim and split prerelease
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, 0, 0, ""
		}
		main := s
		if i := strings.IndexAny(s, "-+"); i >= 0 {
			main, pre = s[:i], s[i+1:]
		}
		parts := strings.Split(main, ".")
		toInt := func(i int) int {
			if i >= len(parts) {
				return 0
			}
			n, _ := strconv.Atoi(parts[i])
			return n
		}
		return toInt(0), toInt(1), toInt(2), pre
	}

	amaj, amin, apat, apre := parse(a)
	bmaj, bmin, bpat, bpre := parse(b)

	switch {
	case amaj != bmaj:
		if amaj > bmaj {
			return 1
		}
		return -1
	case amin != bmin:
		if amin > bmin {
			return 1
		}
		return -1
	case apat != bpat:
		if apat > bpat {
			return 1
		}
		return -1
	}

	// same numeric version; handle pre-release (empty means stable → higher)
	switch {
	case apre == "" && bpre != "":
		return 1
	case apre != "" && bpre == "":
		return -1
	case apre > bpre:
		return 1
	case apre < bpre:
		return -1
	default:
		return 0
	}
}
