package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	filesHandlers "github.com/globulario/Globular/internal/gateway/handlers/files"
)

func TestMinioConfigCache_StrictProbeRunsOnce(t *testing.T) {
	cache := &minioConfigCache{
		ttl:         minioConfigCacheTTL,
		strictProbe: nil,
	}
	probeCalls := 0
	wantCfg := &filesHandlers.MinioProxyConfig{Bucket: "b"}
	cache.strictProbe = func(context.Context) (*filesHandlers.MinioProxyConfig, error) {
		probeCalls++
		return wantCfg, nil
	}

	cfg1, err1 := cache.get()
	if err1 != nil || cfg1 != wantCfg {
		t.Fatalf("first get: cfg=%v err=%v", cfg1, err1)
	}
	cfg2, err2 := cache.get()
	if err2 != nil || cfg2 != wantCfg {
		t.Fatalf("second get: cfg=%v err=%v", cfg2, err2)
	}
	if probeCalls != 1 {
		t.Fatalf("expected strict probe once, got %d", probeCalls)
	}
	if !cache.strictOnce {
		t.Fatalf("strictOnce not set")
	}
}

func TestMinioConfigCache_StrictProbeBackoff(t *testing.T) {
	cache := &minioConfigCache{
		ttl:         minioConfigCacheTTL,
		strictProbe: nil,
	}
	probeCalls := 0
	probeErr := errors.New("boom")
	cache.strictProbe = func(context.Context) (*filesHandlers.MinioProxyConfig, error) {
		probeCalls++
		return nil, probeErr
	}

	cfg1, err1 := cache.get()
	if err1 == nil || !errors.Is(err1, probeErr) || cfg1 != nil {
		t.Fatalf("first get expected error probe: cfg=%v err=%v", cfg1, err1)
	}
	cfg2, err2 := cache.get()
	if err2 == nil || !errors.Is(err2, probeErr) || cfg2 != nil {
		t.Fatalf("second get expected cached error: cfg=%v err=%v", cfg2, err2)
	}
	if probeCalls != 1 {
		t.Fatalf("expected strict probe once during backoff, got %d", probeCalls)
	}
	if cache.strictUntil.IsZero() || cache.strictUntil.Before(time.Now()) {
		t.Fatalf("strictUntil not set for backoff")
	}
}
