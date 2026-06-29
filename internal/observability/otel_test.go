package observability

import (
	"context"
	"testing"
)

func TestInitReturnsCallableShutdown(t *testing.T) {
	shutdown := Init(context.Background())
	if shutdown == nil {
		t.Fatal("Init returned nil shutdown")
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown returned error: %v", err)
	}
}
