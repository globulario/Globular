package dnscache

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// cacheEntry holds a cached DNS lookup result with expiry time
type cacheEntry struct {
	ips       []net.IP
	expiresAt time.Time
}

// srvCacheEntry holds cached SRV lookup results with expiry (PR4.1)
type srvCacheEntry struct {
	records   []*SRVRecord
	expiresAt time.Time
}

// SRVRecord represents a DNS SRV record result (PR4.1)
type SRVRecord struct {
	Priority uint16
	Weight   uint16
	Port     uint16
	Target   string
}

// Cache provides thread-safe DNS caching with TTL support
type Cache struct {
	mu         sync.RWMutex
	entries    map[string]*cacheEntry    // A/AAAA cache
	srvEntries map[string]*srvCacheEntry // SRV cache (PR4.1)
	ttl        time.Duration

	// Metrics (PR5)
	aHit    uint64 // A/AAAA cache hits
	aMiss   uint64 // A/AAAA cache misses
	srvHit  uint64 // SRV cache hits
	srvMiss uint64 // SRV cache misses
}

// New creates a new DNS cache with the specified TTL.
// If ttl is 0 or negative, defaults to 30 seconds.
func New(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &Cache{
		entries:    make(map[string]*cacheEntry),
		srvEntries: make(map[string]*srvCacheEntry),
		ttl:        ttl,
	}
}

// Lookup performs a DNS A/AAAA lookup for the given FQDN.
// Returns cached result if available and not expired.
// On cache miss or expiry, performs fresh lookup and caches the result.
// Returns error if DNS lookup fails.
func (c *Cache) Lookup(ctx context.Context, fqdn string) ([]net.IP, error) {
	// Check cache first
	c.mu.RLock()
	entry, found := c.entries[fqdn]
	c.mu.RUnlock()

	if found && time.Now().Before(entry.expiresAt) {
		atomic.AddUint64(&c.aHit, 1) // PR5: Track cache hit
		return entry.ips, nil
	}

	// Cache miss or expired - perform fresh lookup
	atomic.AddUint64(&c.aMiss, 1) // PR5: Track cache miss
	ips, err := c.lookupFresh(ctx, fqdn)
	if err != nil {
		// If lookup fails but we have stale cache entry, use it as fallback
		if found && entry.ips != nil {
			return entry.ips, nil
		}
		return nil, err
	}

	// Update cache
	c.mu.Lock()
	c.entries[fqdn] = &cacheEntry{
		ips:       ips,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()

	return ips, nil
}

// lookupFresh performs a fresh DNS lookup without consulting the cache
func (c *Cache) lookupFresh(ctx context.Context, fqdn string) ([]net.IP, error) {
	resolver := &net.Resolver{}
	ips, err := resolver.LookupIP(ctx, "ip", fqdn)
	if err != nil {
		return nil, fmt.Errorf("dns lookup %s: %w", fqdn, err)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("dns lookup %s: no records found", fqdn)
	}
	return ips, nil
}

// Invalidate removes a specific FQDN from the cache
func (c *Cache) Invalidate(fqdn string) {
	c.mu.Lock()
	delete(c.entries, fqdn)
	c.mu.Unlock()
}

// InvalidateAll clears the entire cache
func (c *Cache) InvalidateAll() {
	c.mu.Lock()
	c.entries = make(map[string]*cacheEntry)
	c.srvEntries = make(map[string]*srvCacheEntry)
	c.mu.Unlock()
}

// Size returns the number of entries in the cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries) + len(c.srvEntries)
}

// LookupSRV performs a DNS SRV lookup with caching (PR4.1)
// Service and proto should NOT include leading underscores (e.g., "echo-echoservice", "tcp")
// Domain is the base domain (e.g., "cluster.local")
// Returns SRV records sorted by priority (ascending) then weight (descending)
func (c *Cache) LookupSRV(ctx context.Context, service, proto, domain string) ([]*SRVRecord, error) {
	// Format: _service._proto.domain
	name := fmt.Sprintf("_%s._%s.%s", service, proto, domain)

	// Check cache first
	c.mu.RLock()
	entry, found := c.srvEntries[name]
	c.mu.RUnlock()

	if found && time.Now().Before(entry.expiresAt) {
		atomic.AddUint64(&c.srvHit, 1) // PR5: Track cache hit
		return entry.records, nil
	}

	// Cache miss or expired - perform fresh lookup
	atomic.AddUint64(&c.srvMiss, 1) // PR5: Track cache miss
	records, err := c.lookupSRVFresh(ctx, service, proto, domain)
	if err != nil {
		// Fall back to stale cache if available
		if found && entry.records != nil {
			return entry.records, nil
		}
		return nil, err
	}

	// Update cache
	c.mu.Lock()
	c.srvEntries[name] = &srvCacheEntry{
		records:   records,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()

	return records, nil
}

// lookupSRVFresh performs a fresh DNS SRV lookup without consulting the cache
func (c *Cache) lookupSRVFresh(ctx context.Context, service, proto, domain string) ([]*SRVRecord, error) {
	resolver := &net.Resolver{}
	_, srvs, err := resolver.LookupSRV(ctx, service, proto, domain)
	if err != nil {
		return nil, fmt.Errorf("srv lookup %s.%s.%s: %w", service, proto, domain, err)
	}
	if len(srvs) == 0 {
		return nil, fmt.Errorf("srv lookup %s.%s.%s: no records found", service, proto, domain)
	}

	// Convert net.SRV to our SRVRecord
	records := make([]*SRVRecord, len(srvs))
	for i, srv := range srvs {
		records[i] = &SRVRecord{
			Priority: srv.Priority,
			Weight:   srv.Weight,
			Port:     srv.Port,
			Target:   srv.Target,
		}
	}

	return records, nil
}

// InvalidateSRV removes a specific SRV record from the cache (PR4.1)
func (c *Cache) InvalidateSRV(service, proto, domain string) {
	name := fmt.Sprintf("_%s._%s.%s", service, proto, domain)
	c.mu.Lock()
	delete(c.srvEntries, name)
	c.mu.Unlock()
}

// CacheStats holds DNS cache statistics (PR5)
type CacheStats struct {
	AHit    uint64 // A/AAAA cache hits
	AMiss   uint64 // A/AAAA cache misses
	SRVHit  uint64 // SRV cache hits
	SRVMiss uint64 // SRV cache misses
}

// Stats returns cache statistics (PR5)
func (c *Cache) Stats() CacheStats {
	return CacheStats{
		AHit:    atomic.LoadUint64(&c.aHit),
		AMiss:   atomic.LoadUint64(&c.aMiss),
		SRVHit:  atomic.LoadUint64(&c.srvHit),
		SRVMiss: atomic.LoadUint64(&c.srvMiss),
	}
}
