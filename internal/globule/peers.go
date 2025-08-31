package globule

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/event/eventpb"
	"github.com/globulario/services/golang/globular_client"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
)

// -------------------- resource client helper --------------------

func getResourceClient(address string) (*resource_client.Resource_Client, error) {
	Utility.RegisterFunction("NewResourceService_Client", resource_client.NewResourceService_Client)
	c, err := globular_client.GetClient(address, "resource.ResourceService", "NewResourceService_Client")
	if err != nil {
		return nil, err
	}
	return c.(*resource_client.Resource_Client), nil
}

// -------------------- peers public API --------------------

// getPeers returns the list of registered peers from ResourceService.
func (g *Globule) getPeers() ([]*resourcepb.Peer, error) {
	addr, _ := config.GetAddress()
	rc, err := getResourceClient(addr)
	if err != nil {
		return nil, err
	}

	var peers []*resourcepb.Peer
	peers, err = rc.GetPeers("")
	if err == nil {
		return peers, nil
	}

	return nil, fmt.Errorf("getPeers: %w", err)
}

// initPeer updates /etc/hosts (or Windows hosts), fetches public key if missing,
// and refreshes this node's info on the remote ResourceService peer.
func (g *Globule) initPeer(p *resourcepb.Peer) error {
	// Build address for hosts mapping
	address := p.Hostname
	if p.Domain != "localhost" {
		address += "." + p.Domain
	} else if g.Domain != "localhost" && p.Protocol == "https" {
		address += "." + g.Domain
	}

	// If we're on the same external IP, add local IP→hostname mapping
	if Utility.MyIP() == p.GetExternalIpAddress() {
		if err := g.setHost(p.GetLocalIpAddress(), address); err != nil {
			fmt.Println("initPeer: setHost:", err)
		}
	}

	// Determine port and full HTTP(S) address for fetching public key
	if p.Protocol == "https" {
		address += ":" + Utility.ToString(int(p.GetPortHttps()))
	} else {
		address += ":" + Utility.ToString(int(p.GetPortHttp()))
	}

	// Ensure we have the peer's public key locally
	keyPath := filepath.Join(g.configDir, "keys", strings.ReplaceAll(p.Mac, ":", "_")+"_public")
	if !Utility.Exists(keyPath) {
		url := p.Protocol + "://" + address + "/public_key"
		resp, err := http.Get(url) // #nosec G107 — controlled URL
		if err == nil && resp != nil {
			defer resp.Body.Close()
			body, rErr := io.ReadAll(resp.Body)
			if rErr == nil {
				if wErr := osWriteFile0600(keyPath, body); wErr != nil {
					fmt.Println("initPeer: save public key:", wErr)
				}
			}
		}
	}

	// Refresh our own peer info on that remote node (best-effort)
	token, err := security.GenerateToken(g.SessionTimeout, p.GetMac(), "sa", "", g.AdminEmail, g.Domain)
	if err != nil {
		return err
	}

	rc, err := getResourceClient(address)
	if err != nil {
		return err
	}

	peers, err := rc.GetPeers(`{"mac":"` + g.Mac + `"}`)
	if err != nil || len(peers) == 0 {
		return errors.New("initPeer: failed to read local peer on remote")
	}

	me := peers[0]
	me.Protocol = g.Protocol
	me.LocalIpAddress = config.GetLocalIP()
	me.ExternalIpAddress = Utility.MyIP()
	me.PortHttp = int32(g.PortHTTP)
	me.PortHttps = int32(g.PortHTTPS)
	me.Domain = g.Domain

	if err := rc.UpdatePeer(token, me); err != nil {
		return err
	}
	return nil
}

// savePeers persists the current peers into g.Peers and saves config.
func (g *Globule) savePeers() error {
	// Build serializable view
	peers := make([]interface{}, 0)
	g.peers.Range(func(_, v any) bool {
		p := v.(*resourcepb.Peer)
		port := p.PortHttp
		if p.Protocol == "https" {
			port = p.PortHttps
		}
		peers = append(peers, map[string]interface{}{
			"Hostname": p.Hostname,
			"Domain":   p.Domain,
			"Mac":      p.Mac,
			"Port":     int(port),
		})
		return true
	})
	g.Peers = peers
	return g.SaveConfig()
}

// -------------------- event handlers --------------------

// updatePeersEvent handles "update_peers_evt" payloads.
func (g *Globule) updatePeersEvent(evt *eventpb.Event) {
	peer := &resourcepb.Peer{}
	var m map[string]interface{}
	if err := json.Unmarshal(evt.Data, &m); err != nil {
		fmt.Println("updatePeersEvent: bad payload:", err)
		return
	}

	peer.Domain = asString(m["domain"])
	peer.Hostname = asString(m["hostname"])
	peer.Mac = asString(m["mac"])
	peer.Protocol = asString(m["protocol"])
	peer.LocalIpAddress = firstNonEmpty(asString(m["local_ip_address"]), asString(m["localIPAddress"]))
	peer.ExternalIpAddress = firstNonEmpty(asString(m["external_ip_address"]), asString(m["ExternalIPAddress"]))
	peer.PortHttp = int32(Utility.ToInt(m["PortHTTP"]))
	peer.PortHttps = int32(Utility.ToInt(m["PortHTTPS"]))

	// actions (optional)
	if a, ok := m["actions"].([]interface{}); ok {
		peer.Actions = make([]string, 0, len(a))
		for _, x := range a {
			if s, ok := x.(string); ok {
				peer.Actions = append(peer.Actions, s)
			}
		}
	}

	g.peers.Store(peer.Mac, peer)
	if err := g.savePeers(); err != nil {
		fmt.Println("updatePeersEvent: savePeers:", err)
	}

	// If same external IP, maintain local hosts mapping
	if Utility.MyIP() == peer.ExternalIpAddress {
		if err := g.setHost(peer.LocalIpAddress, peer.Hostname+"."+peer.Domain); err != nil {
			fmt.Println("updatePeersEvent: setHost:", err)
		}
	}
}

// deletePeersEvent handles "delete_peer_evt" payloads.
func (g *Globule) deletePeersEvent(evt *eventpb.Event) {
	key := string(evt.Data) // Mac is the typical key
	g.peers.Delete(key)
	if err := g.savePeers(); err != nil {
		fmt.Println("deletePeersEvent: savePeers:", err)
	}
}

// Peers & events (trimmed but behaviorally equivalent)
func (g *Globule) initPeers() error {
	peers, err := g.getPeers()
	if err != nil {
		return err
	}
	for _, p := range peers {
		pp := p // capture
		g.peers.Store(pp.Mac, pp)
		go func() {
			if err := g.initPeer(pp); err != nil {
				g.peers.Delete(pp.Mac)
				_ = g.savePeers()
			}
		}()
	}
	if err := g.subscribe("update_peers_evt", g.updatePeersEvent); err != nil {
		g.log.Warn("subscribe update_peers_evt failed", "err", err)
	}
	if err := g.subscribe("delete_peer_evt", g.deletePeersEvent); err != nil {
		g.log.Warn("subscribe delete_peer_evt failed", "err", err)
	}
	return g.savePeers()
}

// -------------------- helpers --------------------

func asString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func osWriteFile0600(path string, data []byte) error {
	// isolated small helper to keep imports tidy
	return os.WriteFile(path, data, 0o600)
}

// redirectTo looks up the peer for a given host (domain).
// Returns (ok, *Peer) if found.
func (g *Globule) RedirectTo(host string) (bool, *resourcepb.Peer) {

	if host == "" {
		return false, nil
	}

	if host == "localhost" || strings.HasPrefix(host, "localhost:") {
		return false, nil
	}

	peers, err := g.getPeers()
	if err != nil {
		return false, nil
	}

	// Strip port if present
	h := host
	if i := strings.Index(host, ":"); i != -1 {
		h = host[:i]
	}

	var found *resourcepb.Peer
	for _, p := range peers {
		// Match domain or hostname+domain
		if p.Domain == h || (p.Hostname+"."+p.Domain) == h {
			found = p
			break
		}
	}

	if found != nil {
		return true, found
	}

	return false, nil
}
