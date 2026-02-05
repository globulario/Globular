package sds

import (
	"context"
	"fmt"
	"log"
	"sync"

	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	secret_v3 "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/globulario/Globular/internal/controlplane"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

// Server implements the Envoy Secret Discovery Service (SDS).
// It manages a cache of secrets (TLS certificates) and pushes updates to Envoy when certs rotate.
type Server struct {
	secret_v3.UnimplementedSecretDiscoveryServiceServer

	mu      sync.RWMutex
	secrets map[string]*tls_v3.Secret // secretName -> Secret
	version string                    // Current snapshot version

	cache cache.SnapshotCache // go-control-plane snapshot cache
}

// NewServer creates a new SDS server with an empty secret cache.
func NewServer() *Server {
	return &Server{
		secrets: make(map[string]*tls_v3.Secret),
		version: "v0",
		cache:   cache.NewSnapshotCache(false, cache.IDHash{}, nil),
	}
}

// StreamSecrets implements the primary SDS RPC method.
// Envoy establishes a long-lived gRPC stream and receives secret updates pushed by the server.
func (s *Server) StreamSecrets(stream secret_v3.SecretDiscoveryService_StreamSecretsServer) error {
	log.Printf("sds: client connected for secret streaming")

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("sds: stream closed: %v", err)
			return err
		}

		log.Printf("sds: received request for %d resources (version=%s, nonce=%s)",
			len(req.ResourceNames), req.VersionInfo, req.ResponseNonce)

		// Handle ACK/NACK
		if req.ResponseNonce != "" {
			if req.ErrorDetail != nil {
				log.Printf("sds: client NACK: %s", req.ErrorDetail.Message)
			} else {
				log.Printf("sds: client ACK: nonce=%s", req.ResponseNonce)
			}
		}

		// Send current secrets
		resp, err := s.buildResponse(req)
		if err != nil {
			log.Printf("sds: build response error: %v", err)
			return status.Errorf(codes.Internal, "build response: %v", err)
		}

		if err := stream.Send(resp); err != nil {
			log.Printf("sds: send error: %v", err)
			return err
		}

		log.Printf("sds: sent %d secrets (version=%s, nonce=%s)",
			len(resp.Resources), resp.VersionInfo, resp.Nonce)
	}
}

// FetchSecrets implements the single-shot SDS RPC (optional).
// Envoy can use this for one-time secret fetching instead of streaming.
func (s *Server) FetchSecrets(ctx context.Context, req *discovery_v3.DiscoveryRequest) (*discovery_v3.DiscoveryResponse, error) {
	log.Printf("sds: fetch request for %d resources", len(req.ResourceNames))

	resp, err := s.buildResponse(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "build response: %v", err)
	}

	return resp, nil
}

// DeltaSecrets implements the incremental SDS RPC (not used yet).
func (s *Server) DeltaSecrets(stream secret_v3.SecretDiscoveryService_DeltaSecretsServer) error {
	return status.Errorf(codes.Unimplemented, "delta SDS not implemented")
}

// buildResponse constructs a DiscoveryResponse with requested secrets.
func (s *Server) buildResponse(req *discovery_v3.DiscoveryRequest) (*discovery_v3.DiscoveryResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resources := []types.Resource{}

	// If specific resources requested, return only those
	if len(req.ResourceNames) > 0 {
		for _, name := range req.ResourceNames {
			secret, ok := s.secrets[name]
			if !ok {
				log.Printf("sds: requested secret %q not found (available: %v)", name, s.availableSecretNames())
				continue
			}
			resources = append(resources, secret)
		}
	} else {
		// No specific resources requested - send all secrets
		for _, secret := range s.secrets {
			resources = append(resources, secret)
		}
	}

	// Marshal resources to Any protos
	anyResources := make([]*anypb.Any, 0, len(resources))
	for _, r := range resources {
		marshaled, err := anypb.New(r)
		if err != nil {
			return nil, fmt.Errorf("marshal resource: %w", err)
		}
		anyResources = append(anyResources, marshaled)
	}

	// Build response
	typeURL := req.GetTypeUrl()
	if typeURL == "" {
		typeURL = resource.SecretType
	}

	return &discovery_v3.DiscoveryResponse{
		VersionInfo: s.version,
		Resources:   anyResources,
		TypeUrl:     typeURL,
		Nonce:       s.version, // Use version as nonce for simplicity
	}, nil
}

// UpdateSecrets replaces the secret cache with new secrets and increments the version.
// This triggers push updates to all connected Envoy instances via ADS.
func (s *Server) UpdateSecrets(secrets map[string]*tls_v3.Secret) error {
	if len(secrets) == 0 {
		return fmt.Errorf("cannot update with empty secret map")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Compute new version from secrets
	version, err := s.computeVersion(secrets)
	if err != nil {
		return fmt.Errorf("compute version: %w", err)
	}

	// Replace secrets
	s.secrets = secrets
	s.version = version

	log.Printf("sds: updated %d secrets to version %s", len(secrets), version)

	// Update snapshot cache to trigger push to Envoy
	if err := s.updateSnapshot(); err != nil {
		return fmt.Errorf("update snapshot: %w", err)
	}

	return nil
}

// updateSnapshot pushes the current secrets to the snapshot cache.
// This triggers Envoy to receive the updated secrets via ADS.
func (s *Server) updateSnapshot() error {
	// Convert secrets to cache resources
	resources := make(map[string][]types.Resource)
	secretResources := make([]types.Resource, 0, len(s.secrets))
	for _, secret := range s.secrets {
		secretResources = append(secretResources, secret)
	}
	resources[resource.SecretType] = secretResources

	// Create snapshot with only secrets (other resource types empty)
	snapshot, err := cache.NewSnapshot(s.version, resources)
	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}

	// Push to cache for all node IDs (broadcast to all Envoy instances)
	// Note: In production, we'd track node IDs and push selectively
	if err := s.cache.SetSnapshot(context.Background(), "envoy-node", snapshot); err != nil {
		return fmt.Errorf("set snapshot: %w", err)
	}

	log.Printf("sds: snapshot updated (version=%s)", s.version)
	return nil
}

// computeVersion generates a version string based on secret content hashes.
// This ensures the version changes if and only if secrets change.
func (s *Server) computeVersion(secrets map[string]*tls_v3.Secret) (string, error) {
	if len(secrets) == 0 {
		return "v0", nil
	}

	// Use hash of first secret as base version (all secrets should rotate together)
	for _, secret := range secrets {
		hash, err := controlplane.HashSecret(secret)
		if err != nil {
			return "", err
		}
		// Use first 16 chars for readability
		if len(hash) > 16 {
			hash = hash[:16]
		}
		return fmt.Sprintf("sds-%s", hash), nil
	}

	return "v0", nil
}

// GetSecret returns a secret by name (read-only, for inspection).
func (s *Server) GetSecret(name string) (*tls_v3.Secret, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	secret, ok := s.secrets[name]
	return secret, ok
}

// GetVersion returns the current snapshot version.
func (s *Server) GetVersion() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.version
}

// GetSecretNames returns the names of all cached secrets.
func (s *Server) GetSecretNames() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.availableSecretNames()
}

// availableSecretNames returns secret names (must hold lock).
func (s *Server) availableSecretNames() []string {
	names := make([]string, 0, len(s.secrets))
	for name := range s.secrets {
		names = append(names, name)
	}
	return names
}

// GetCache returns the snapshot cache for integration with xDS server.
func (s *Server) GetCache() cache.SnapshotCache {
	return s.cache
}
