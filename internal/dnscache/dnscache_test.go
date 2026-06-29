package dnscache

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestLookupReturnsStaleEntryOnRefreshFailure(t *testing.T) {
	cache := New(50*time.Millisecond, "127.0.0.1:1")
	stale := net.ParseIP("10.0.0.8")
	cache.entries["svc.cluster.local"] = &cacheEntry{
		ips:       []net.IP{stale},
		expiresAt: time.Now().Add(-time.Second),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	got, err := cache.Lookup(ctx, "svc.cluster.local")
	if err != nil {
		t.Fatalf("Lookup returned error instead of stale answer: %v", err)
	}
	if len(got) != 1 || !got[0].Equal(stale) {
		t.Fatalf("Lookup=%v want stale %v", got, stale)
	}
}

func TestInvalidateAllClearsAAndSRVCaches(t *testing.T) {
	cache := New(time.Minute)
	cache.entries["svc.cluster.local"] = &cacheEntry{
		ips:       []net.IP{net.ParseIP("10.0.0.9")},
		expiresAt: time.Now().Add(time.Minute),
	}
	cache.srvEntries["_svc._tcp.cluster.local"] = &srvCacheEntry{
		records:   []*SRVRecord{{Priority: 1, Weight: 10, Port: 8080, Target: "svc.cluster.local."}},
		expiresAt: time.Now().Add(time.Minute),
	}

	if size := cache.Size(); size != 2 {
		t.Fatalf("Size before InvalidateAll=%d want 2", size)
	}

	cache.InvalidateAll()

	if size := cache.Size(); size != 0 {
		t.Fatalf("Size after InvalidateAll=%d want 0", size)
	}
}

func TestLookupSRVReturnsStaleEntryOnRefreshFailure(t *testing.T) {
	cache := New(50*time.Millisecond, "127.0.0.1:1")
	stale := []*SRVRecord{{Priority: 1, Weight: 10, Port: 8443, Target: "svc.cluster.local."}}
	cache.srvEntries["_svc._tcp.cluster.local"] = &srvCacheEntry{
		records:   stale,
		expiresAt: time.Now().Add(-time.Second),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	got, err := cache.LookupSRV(ctx, "svc", "tcp", "cluster.local")
	if err != nil {
		t.Fatalf("LookupSRV returned error instead of stale answer: %v", err)
	}
	if len(got) != 1 || got[0].Target != stale[0].Target || got[0].Port != stale[0].Port {
		t.Fatalf("LookupSRV=%v want stale %v", got, stale)
	}
}
