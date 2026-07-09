package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeSignBackend struct {
	planJSON []byte
	planErr  error
	ips      []string
	ipsErr   error
}

func (f fakeSignBackend) GetJoinRequestStatusPlanJSON(_ context.Context, _ string) ([]byte, error) {
	return f.planJSON, f.planErr
}

func (f fakeSignBackend) ListNodeIPs(_ context.Context) ([]string, error) {
	return f.ips, f.ipsErr
}

func newTestAuthorizer(backend signControllerClient, localIPs []string) signAuthorizer {
	return signAuthorizer{
		controller: backend,
		localAddrs: func() []string { return localIPs },
	}
}

func signReq(remoteAddr, joinID string, headers map[string]string) *http.Request {
	url := "/sign_ca_certificate?csr=abc"
	if joinID != "" {
		url += "&join_id=" + joinID
	}
	r := httptest.NewRequest(http.MethodGet, url, nil)
	r.RemoteAddr = remoteAddr
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	return r
}

// A valid controller-issued join_id (with a stored JoinPlan) authorizes
// signing regardless of the requester's network position.
func TestSignAuthorizer_ValidJoinIDAllows(t *testing.T) {
	a := newTestAuthorizer(fakeSignBackend{planJSON: []byte(`{"join_id":"j1"}`)}, nil)
	if err := a.AuthorizeSign(signReq("203.0.113.7:4242", "j1", nil)); err != nil {
		t.Fatalf("valid join_id must authorize signing, got: %v", err)
	}
}

// A join_id the controller does not recognize is denied — the join_id is the
// credential, so an unverifiable one must fail closed.
func TestSignAuthorizer_UnknownJoinIDDenies(t *testing.T) {
	a := newTestAuthorizer(fakeSignBackend{planErr: errors.New("not found")}, nil)
	if err := a.AuthorizeSign(signReq("203.0.113.7:4242", "bogus", nil)); err == nil {
		t.Fatal("unknown join_id must be denied")
	}
}

// A join record without a stored JoinPlan was never authorized (denied at
// preflight) — it must not grant certificate issuance.
func TestSignAuthorizer_JoinIDWithoutPlanDenies(t *testing.T) {
	a := newTestAuthorizer(fakeSignBackend{planJSON: nil}, nil)
	if err := a.AuthorizeSign(signReq("203.0.113.7:4242", "j1", nil)); err == nil {
		t.Fatal("join_id without an authorized plan must be denied")
	}
}

// Direct loopback requests are the local node itself (Day-0 bootstrap, local
// CLI) and stay allowed even with no controller available.
func TestSignAuthorizer_LoopbackAllows(t *testing.T) {
	a := newTestAuthorizer(fakeSignBackend{ipsErr: errors.New("controller down")}, nil)
	if err := a.AuthorizeSign(signReq("127.0.0.1:5555", "", nil)); err != nil {
		t.Fatalf("loopback requester must be allowed, got: %v", err)
	}
}

// Direct requests from the host's own interface addresses are the local node
// (a process dialing its own gateway by node IP).
func TestSignAuthorizer_LocalInterfaceIPAllows(t *testing.T) {
	a := newTestAuthorizer(fakeSignBackend{ipsErr: errors.New("controller down")}, []string{"10.0.0.63"})
	if err := a.AuthorizeSign(signReq("10.0.0.63:5555", "", nil)); err != nil {
		t.Fatalf("local interface requester must be allowed, got: %v", err)
	}
}

// Direct requests from registered cluster member IPs are cert renewals.
func TestSignAuthorizer_MemberIPAllows(t *testing.T) {
	a := newTestAuthorizer(fakeSignBackend{ips: []string{"10.0.0.8", "10.0.0.20"}}, nil)
	if err := a.AuthorizeSign(signReq("10.0.0.20:39218", "", nil)); err != nil {
		t.Fatalf("cluster member requester must be allowed, got: %v", err)
	}
}

// Anyone else — reachable but neither local, member, nor join_id-bearing —
// must be denied. This is the hole the gate exists to close: a CA-signed cert
// is an etcd mTLS credential.
func TestSignAuthorizer_StrangerDenied(t *testing.T) {
	a := newTestAuthorizer(fakeSignBackend{ips: []string{"10.0.0.8"}}, []string{"10.0.0.63"})
	if err := a.AuthorizeSign(signReq("203.0.113.7:4242", "", nil)); err == nil {
		t.Fatal("non-member without join_id must be denied")
	}
}

// Proxied requests (X-Forwarded-For present) must not qualify for IP-based
// trust: the TCP peer is the proxy (possibly a cluster member forwarding
// external traffic), not the requester.
func TestSignAuthorizer_ProxiedRequestRequiresJoinID(t *testing.T) {
	a := newTestAuthorizer(fakeSignBackend{ips: []string{"10.0.0.8"}}, nil)
	req := signReq("10.0.0.8:39218", "", map[string]string{"X-Forwarded-For": "203.0.113.7"})
	if err := a.AuthorizeSign(req); err == nil {
		t.Fatal("proxied request without join_id must be denied even from a member peer IP")
	}
	// ...but a proxied request WITH a valid join_id is fine.
	a2 := newTestAuthorizer(fakeSignBackend{planJSON: []byte(`{}`)}, nil)
	req2 := signReq("10.0.0.8:39218", "j1", map[string]string{"X-Forwarded-For": "203.0.113.7"})
	if err := a2.AuthorizeSign(req2); err != nil {
		t.Fatalf("proxied request with valid join_id must be allowed, got: %v", err)
	}
}

// Controller unreachable → membership cannot be proven → fail closed for
// non-local requesters.
func TestSignAuthorizer_ControllerDownFailsClosed(t *testing.T) {
	a := newTestAuthorizer(fakeSignBackend{ipsErr: errors.New("controller down")}, nil)
	if err := a.AuthorizeSign(signReq("10.0.0.20:39218", "", nil)); err == nil {
		t.Fatal("membership check failure must deny (fail closed)")
	}
}
