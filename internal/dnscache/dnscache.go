package dnscache

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// cacheEntry holds a cached DNS lookup result with expiry time
type cacheEntry struct {
	ips       []net.IP
	expiresAt time.Time
}

// Cache provides thread-safe DNS caching with TTL support
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*cacheEntry
	ttl     time.Duration
}

// New creates a new DNS cache with the specified TTL.
// If ttl is 0 or negative, defaults to 30 seconds.
func New(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &Cache{
		entries: make(map[string]*cacheEntry),
		ttl:     ttl,
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
		return entry.ips, nil
	}

	// Cache miss or expired - perform fresh lookup
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
	c.mu.Unlock()
}

// Size returns the number of entries in the cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
