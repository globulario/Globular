package globule

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/event/eventpb"
	"github.com/globulario/services/golang/globular_client"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/resource/resourcepb"
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

// -------------------- nodes public API --------------------

// getNodeIdentities returns the list of stored node identities from ResourceService.
func (g *Globule) getNodeIdentities() ([]*resourcepb.NodeIdentity, error) {
	addr, _ := config.GetAddress()
	rc, err := getResourceClient(addr)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var nodes []*resourcepb.NodeIdentity
	nodes, err = rc.ListNodeIdentities("", "")
	if err == nil {
		return nodes, nil
	}

	return nil, fmt.Errorf("getNodeIdentities: %w", err)
}

// initNodeIdentity updates /etc/hosts (or Windows hosts), fetches public key if missing,
// and refreshes this node's info on the remote ResourceService node.
func (g *Globule) initNodeIdentity(p *resourcepb.NodeIdentity) error {
	// Build host for hosts mapping.
	host := g.nodeHostname(p)

	// If we're on the same external IP, add local IP→hostname mapping (optional).
	if g.MutateHostsFile && Utility.MyIP() == p.GetExternalIpAddress() {
		if err := g.setHost(p.GetLocalIpAddress(), host); err != nil {
			g.log.Warn("initNodeIdentity: setHost failed", "err", err)
		}
	}

	apiAddr := g.nodeAPIAddress(p)

	// Ensure we have the node's public key locally
	keyPath := filepath.Join(g.configDir, "keys", strings.ReplaceAll(p.Mac, ":", "_")+"_public")
	if !Utility.Exists(keyPath) {
		url := p.Protocol + "://" + apiAddr + "/public_key"
		resp, err := http.Get(url) // #nosec G107 — controlled URL
		if err == nil && resp != nil {
			defer resp.Body.Close()
			body, rErr := io.ReadAll(resp.Body)
			if rErr == nil {
				if wErr := osWriteFile0600(keyPath, body); wErr != nil {
					fmt.Println("initNodeIdentity: save public key:", wErr)
				}
			}
		}
	}

	// Refresh our own node identity on that remote node (best-effort)
	rc, err := getResourceClient(apiAddr)
	if err != nil {
		return err
	}
	defer rc.Close()

	nodes, err := rc.ListNodeIdentities(`{"mac":"`+g.Mac+`"}`, "")
	if err != nil || len(nodes) == 0 {
		return errors.New("initNodeIdentity: failed to read local node on remote")
	}

	me := nodes[0]
	me.Protocol = g.Protocol
	me.LocalIpAddress = config.GetLocalIP()
	me.ExternalIpAddress = Utility.MyIP()
	me.PortHttp = int32(g.PortHTTP)
	me.PortHttps = int32(g.PortHTTPS)
	me.Domain = g.Domain

	if g.EnablePeerUpserts {
		if _, err := rc.UpsertNodeIdentity(me); err != nil {
			return err
		}
	} else {
		g.log.Debug("node upserts disabled; skipping remote identity patch")
	}
	return nil
}

// saveNodeIdentities persists the current nodes into g.Nodes and saves config.
func (g *Globule) saveNodeIdentities() error {
	// Build serializable view
	nodes := make([]interface{}, 0)
	g.nodes.Range(func(_, v any) bool {
		p := v.(*resourcepb.NodeIdentity)
		port := p.PortHttp
		if p.Protocol == "https" {
			port = p.PortHttps
		}
		nodes = append(nodes, map[string]interface{}{
			"Hostname": p.Hostname,
			"Domain":   p.Domain,
			"Mac":      p.Mac,
			"Port":     int(port),
		})
		return true
	})
	g.Nodes = nodes
	return g.SaveConfig()
}

// -------------------- event handlers --------------------

// updateNodesEvent handles "update_peers_evt" payloads.
func (g *Globule) updateNodesEvent(evt *eventpb.Event) {
	node := &resourcepb.NodeIdentity{}
	var m map[string]interface{}
	if err := json.Unmarshal(evt.Data, &m); err != nil {
		fmt.Println("updateNodesEvent: bad payload:", err)
		return
	}

	node.Domain = asString_(m["domain"])
	node.Hostname = asString_(m["hostname"])
	node.Mac = asString_(m["mac"])
	node.Protocol = asString_(m["protocol"])
	node.LocalIpAddress = firstNonEmpty(asString_(m["local_ip_address"]), asString_(m["localIPAddress"]))
	node.ExternalIpAddress = firstNonEmpty(asString_(m["external_ip_address"]), asString_(m["ExternalIPAddress"]))
	node.PortHttp = int32(Utility.ToInt(m["PortHTTP"]))
	node.PortHttps = int32(Utility.ToInt(m["PortHTTPS"]))
	node.NodeId = firstNonEmpty(asString_(m["node_id"]), asString_(m["_id"]))
	node.Fingerprint = asString_(m["fingerprint"])
	node.Status = asString_(m["status"])
	node.Enabled = Utility.ToBool(m["enabled"])

	g.nodes.Store(node.Mac, node)
	if err := g.saveNodeIdentities(); err != nil {
		fmt.Println("updateNodesEvent: saveNodeIdentities:", err)
	}

	// If same external IP, maintain local hosts mapping
	if g.MutateHostsFile && Utility.MyIP() == node.ExternalIpAddress {
		if err := g.setHost(node.LocalIpAddress, g.nodeHostname(node)); err != nil {
			g.log.Warn("updateNodesEvent: setHost failed", "err", err)
		}
	}
}

// deleteNodesEvent handles "delete_peer_evt" payloads.
func (g *Globule) deleteNodesEvent(evt *eventpb.Event) {
	key := string(evt.Data) // Mac is the typical key
	g.nodes.Delete(key)
	if err := g.saveNodeIdentities(); err != nil {
		fmt.Println("deleteNodesEvent: saveNodeIdentities:", err)
	}
}

// Nodes & events (trimmed but behaviorally equivalent)
func (g *Globule) initNodes() error {
	nodesList, err := g.getNodeIdentities()
	if err != nil {
		return err
	}
	for _, p := range nodesList {
		pp := p // capture
		g.nodes.Store(pp.Mac, pp)
		go func() {
			if err := g.initNodeIdentity(pp); err != nil {
				g.nodes.Delete(pp.Mac)
				_ = g.saveNodeIdentities()
			}
		}()
	}
	if err := g.subscribe("update_peers_evt", g.updateNodesEvent); err != nil {
		g.log.Warn("subscribe update_peers_evt failed", "err", err)
	}
	if err := g.subscribe("delete_peer_evt", g.deleteNodesEvent); err != nil {
		g.log.Warn("subscribe delete_peer_evt failed", "err", err)
	}
	return g.saveNodeIdentities()
}

// -------------------- helpers --------------------

func asString_(v interface{}) string {
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

func (g *Globule) nodeHostname(p *resourcepb.NodeIdentity) string {
	host := p.Hostname
	if p.Domain != "localhost" {
		host += "." + p.Domain
	} else if g.Domain != "localhost" && p.Protocol == "https" {
		host += "." + g.Domain
	}
	return host
}

func (g *Globule) nodeAPIAddress(p *resourcepb.NodeIdentity) string {
	host := g.nodeHostname(p)
	var port string
	if p.Protocol == "https" {
		port = Utility.ToString(int(p.GetPortHttps()))
	} else {
		port = Utility.ToString(int(p.GetPortHttp()))
	}
	if port == "" {
		return host
	}
	return net.JoinHostPort(host, port)
}

func osWriteFile0600(path string, data []byte) error {
	// isolated small helper to keep imports tidy
	return os.WriteFile(path, data, 0o600)
}
