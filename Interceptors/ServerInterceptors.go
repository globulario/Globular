package Interceptors

// TODO for the validation, use a map to store valid method/token/resource/access
// the validation will be renew only if the token expire. And when a token expire
// the value in the map will be discard. That way it will put less charge on the server
// side.

import (
	"context"
	"errors"
	"fmt"

	"log"
	"strings"

	"github.com/davecourtois/Utility"
	globular "github.com/globulario/services/golang/globular_client"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/storage/storage_store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	// The resource client
	resource_client_ *resource_client.Resource_Client

	// The map will contain connection with other server of same kind to load
	// balance the server charge.
	clients map[string]globular.Client

	// That will contain the permission in memory to limit the number
	// of resource request...
	cache *storage_store.BigCache_store
)

/**
 * Get a the local resource client.
 */
func GetResourceClient(domain string) (*resource_client.Resource_Client, error) {
	var err error
	if resource_client_ == nil {
		resource_client_, err = resource_client.NewResourceService_Client(domain, "resource.ResourceService")
		if err != nil {
			return nil, err
		}
	}

	return resource_client_, nil
}

/**
 * A singleton use to access the cache.
 */
func getCache() *storage_store.BigCache_store {
	if cache == nil {
		cache = storage_store.NewBigCache_store()
		err := cache.Open("")
		if err != nil {
			fmt.Println(err)
		}
	}
	return cache
}

// Refresh a token.
func refreshToken(domain string, token string) (string, error) {
	resource_client, err := GetResourceClient(domain)
	if err != nil {
		return "", err
	}

	return resource_client.RefreshToken(token)
}

// That interceptor is use by all services except the resource service who has
// it own interceptor.
func ServerUnaryInterceptor(ctx context.Context, rqst interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// The token and the application id.
	var token string
	var application string
	var domain string // This is the target domain, the one use in TLS certificate.

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")

		// in case of resource path.
		domain = strings.Join(md["domain"], "")
		if strings.Index(domain, ":") == 0 {
			port := strings.Join(md["port"], "")
			if len(port) > 0 {
				// set the address in that particular case.
				domain += ":" + port
			}
		}
	}

	// p, _ := peer.FromContext(ctx)
	// Here I will test if the
	method := info.FullMethod

	if len(domain) == 0 {
		return nil, errors.New("No domain was given for method call '" + method + "'")
	}

	// If the call come from a local client it has hasAccess
	hasAccess := false //strings.HasPrefix(p.Addr.String(), "127.0.0.1") || strings.HasPrefix(p.Addr.String(), Utility.MyIP())

	var pwd string
	if Utility.GetProperty(info.Server, "RootPassword") != nil {
		pwd = Utility.GetProperty(info.Server, "RootPassword").(string)
	}

	// needed to get access to the system.
	if method == "/admin.AdminService/GetConfig" ||
		method == "/admin.AdminService/HasRunningProcess" ||
		method == "/services.ServiceDiscovery/FindServices" ||
		method == "/services.ServiceDiscovery/GetServiceDescriptor" ||
		method == "/services.ServiceDiscovery/GetServicesDescriptor" ||
		method == "/dns.DnsService/GetA" ||
		method == "/dns.DnsService/GetAAAA" ||
		method == "/resource.ResourceService/Log" {
		hasAccess = true
	} else if (method == "/admin.AdminService/SetRootEmail" || method == "/admin.AdminService/SetRootPassword") && ((domain == "127.0.0.1" || domain == "localhost") || pwd == "adminadmin") {
		hasAccess = true
	}

	var clientId string
	var err error

	if len(token) > 0 {
		clientId, _, _, err = ValidateToken(token)
		if err != nil {
			log.Println("token validation fail with error: ", err)
			return nil, err
		}
		if clientId == "sa" {
			hasAccess = true
		}
	}

	// Test if peer has access
	if !hasAccess {
		// TODO validate gRPC method access here.
		hasAccess = true // for convenience...
		resource_client_, err := GetResourceClient(domain)
		if err == nil {
			infos, err := resource_client_.GetActionResourceInfos(method)
			log.Println("-----------> info", infos, err)
		}
	}

	if !hasAccess {
		err := errors.New("Permission denied to execute method " + method + " user:" + clientId + " domain:" + domain + " application:" + application)
		fmt.Println(err)
		log.Println("validation fail ", err)
		//resource_client.Log(application, clientId, method, err)
		return nil, err
	}

	// Now I will test file permission.
	if clientId != "sa" {
		// TODO validate gRPC ressource access here.
	}

	var result interface{}

	result, err = handler(ctx, rqst)

	// Send log event...
	if (len(application) > 0 && len(clientId) > 0 && clientId != "sa") || err != nil {
		// resource_client.Log(application, clientId, method, err)
	}

	return result, err

}

// A wrapper for the real grpc.ServerStream
type ServerStreamInterceptorStream struct {
	inner       grpc.ServerStream
	method      string
	domain      string
	peer        string
	token       string
	application string
	clientId    string
}

func (l ServerStreamInterceptorStream) SetHeader(m metadata.MD) error {
	return l.inner.SetHeader(m)
}

func (l ServerStreamInterceptorStream) SendHeader(m metadata.MD) error {
	return l.inner.SendHeader(m)
}

func (l ServerStreamInterceptorStream) SetTrailer(m metadata.MD) {
	l.inner.SetTrailer(m)
}

func (l ServerStreamInterceptorStream) Context() context.Context {
	return l.inner.Context()
}

func (l ServerStreamInterceptorStream) SendMsg(m interface{}) error {
	return l.inner.SendMsg(m)
}

/**
 * Here I will wrap the original stream into this one to get access to the original
 * rqst, so I can validate it resources.
 */
func (l ServerStreamInterceptorStream) RecvMsg(rqst interface{}) error {

	if l.clientId != "sa" {
		// TODO validate gRpc method and ressource access. here.
	}
	return l.inner.RecvMsg(rqst)
}

// Stream interceptor.
func ServerStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	// The token and the application id.
	var token string
	var application string

	// The peer domain.
	var domain string

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
		domain = strings.Join(md["domain"], "")

		if strings.Index(domain, ":") == 0 {
			port := strings.Join(md["port"], "")
			if len(port) > 0 {
				// set the address in that particular case.
				domain += ":" + port
			}
		}
	}

	p, _ := peer.FromContext(stream.Context())

	method := info.FullMethod

	var clientId string
	var err error
	// If the call come from a local client it has hasAccess
	hasAccess := strings.HasPrefix(p.Addr.String(), "127.0.0.1") || strings.HasPrefix(p.Addr.String(), Utility.MyIP())

	if len(token) > 0 {
		clientId, _, _, err = ValidateToken(token)
		if err != nil {
			return err
		}
		if clientId == "sa" {
			hasAccess = true
		}
	}

	// needed by the admin.
	if application == "admin" ||
		method == "/services.ServiceDiscovery/GetServicesDescriptor" ||
		method == "/services.ServiceRepository/DownloadBundle" {
		hasAccess = true
	}

	if !hasAccess {
		// TODO validate gRpc method access
		hasAccess = true // for convenience...
	}

	// Return here if access is denied.
	if !hasAccess {
		return errors.New("Permission denied to execute method " + method)
	}

	// Now the permissions
	err = handler(srv, ServerStreamInterceptorStream{inner: stream, method: method, domain: domain, token: token, application: application, clientId: clientId, peer: domain})

	if err != nil {
		return err
	}

	return nil
}
