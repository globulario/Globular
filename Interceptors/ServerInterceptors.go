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
	"time"

	"github.com/davecourtois/Utility"
	globular "github.com/globulario/services/golang/globular_client"
	"github.com/globulario/services/golang/lb/lbpb"
	"github.com/globulario/services/golang/lb/load_balancing_client"
	"github.com/globulario/services/golang/rbac/rbac_client"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/storage/storage_store"
	"github.com/shirou/gopsutil/load"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	// The resource client
	resource_client_ *resource_client.Resource_Client

	// The rbac client
	rbac_client_ *rbac_client.Rbac_Client

	// The load balancer client.
	lb_client_ *load_balancing_client.Lb_Client

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
 * Get the rbac client.
 */
func GetRbacClient(domain string) (*rbac_client.Rbac_Client, error) {
	var err error
	if rbac_client_ == nil {
		rbac_client_, err = rbac_client.NewRbacService_Client(domain, "rbac.RbacService")
		if err != nil {
			return nil, err
		}
	}

	return rbac_client_, nil
}

/**
 * Get a the local ressource client.
 */
func getLoadBalancingClient(domain string, serverId string, serviceName string, serverDomain string, serverPort int32) (*load_balancing_client.Lb_Client, error) {

	var err error
	if lb_client_ == nil {
		lb_client_, err = load_balancing_client.NewLbService_Client(domain, "lb.LoadBalancingService")
		if err != nil {
			return nil, err
		}

		// Here I will create the client map.
		clients = make(map[string]globular.Client)

		// Now I will start reporting load at each minutes.
		ticker := time.NewTicker(1 * time.Minute)
		go func() {
			for {
				select {
				case <-ticker.C:
					stats, err := load.Avg()
					if err != nil {
						break
					}
					load_info := &lbpb.LoadInfo{
						ServerInfo: &lbpb.ServerInfo{
							Id:     serverId,
							Name:   serviceName,
							Domain: serverDomain,
							Port:   serverPort,
						},
						Load1:  stats.Load1,
						Load5:  stats.Load5,
						Load15: stats.Load15,
					}

					lb_client_.ReportLoadInfo(load_info)
				}
			}
		}()
	}

	return lb_client_, nil
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

func validateActionRequest(rqst interface{}, method string, subject string, subjectType rbacpb.SubjectType, domain string) (bool, error) {

	rbac_client_, err := GetRbacClient(domain)
	if err != nil {
		return false, err
	}
	hasAccess := false
	infos, err := rbac_client_.GetActionResourceInfos(method)
	if err != nil {
		infos = make([]*rbacpb.ResourceInfos, 0)
	} else {
		// Here I will get the params...
		val, _ := Utility.CallMethod(rqst, "ProtoReflect", []interface{}{})
		rqst_ := val.(protoreflect.Message)
		if rqst_.Descriptor().Fields().Len() > 0 {
			for i := 0; i < len(infos); i++ {
				// Get the path value from retreive infos.
				param := rqst_.Descriptor().Fields().Get(Utility.ToInt(infos[i].Index))
				val := rqst_.Get(param)
				infos[i].Path = val.String()
			}
		}
	}

	// TODO keep to value in cache for keep speed.

	hasAccess, err = rbac_client_.ValidateAction(method, subject, subjectType, infos)
	if err != nil {
		return false, err
	}

	// Here I will store the permission for further use...
	return hasAccess, nil
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

	// Get the peer information.
	p, _ := peer.FromContext(ctx)

	// Here I will test if the
	method := info.FullMethod
	if Utility.IsLocal(p.Addr.String()) {
		domain = "localhost"
	}

	if len(domain) == 0 {
		if strings.Index(p.Addr.String(), ":") != 0 {
			domain = p.Addr.String()[0:strings.Index(p.Addr.String(), ":")]
		} else {
			domain = p.Addr.String()
		}
	}

	// If the call come from a local client it has hasAccess
	hasAccess := false

	var pwd string
	if Utility.GetProperty(info.Server, "RootPassword") != nil {
		pwd = Utility.GetProperty(info.Server, "RootPassword").(string)
	}

	// needed to get access to the system.
	if method == "/admin.AdminService/GetConfig" ||
		method == "/admin.AdminService/HasRunningProcess" ||
		method == "/admin.AdminService/InstallCertificates" ||
		method == "/services.PackageDiscovery/FindServices" ||
		method == "/services.PackageDiscovery/GetServiceDescriptor" ||
		method == "/services.PackageDiscovery/GetServicesDescriptor" ||
		method == "/rbac.RbacService/ValidateAction" ||
		method == "/rbac.RbacService/ValidateAccess" ||
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
	if !hasAccess && len(clientId) > 0 {
		hasAccess, _ = validateActionRequest(rqst, method, clientId, rbacpb.SubjectType_ACCOUNT, domain)
	}

	if !hasAccess && len(application) > 0 {
		hasAccess, _ = validateActionRequest(rqst, method, application, rbacpb.SubjectType_APPLICATION, domain)
	}

	if !hasAccess {
		hasAccess, _ = validateActionRequest(rqst, method, p.Addr.String(), rbacpb.SubjectType_PEER, domain)
	}

	if !hasAccess {
		err := errors.New("Permission denied to execute method " + method + " user:" + clientId + " domain:" + domain + " application:" + application)
		fmt.Println(err)
		log.Println("validation fail ", err)
		return nil, err
	}

	// So here the user has access to the ressource...

	// I will try to get the list of candidates for load balancing
	if Utility.GetProperty(info.Server, "Port") != nil {

		// Here I will refresh the load balance of the server to keep track of real load.
		lb_client, err := getLoadBalancingClient(domain, Utility.GetProperty(info.Server, "Id").(string), Utility.GetProperty(info.Server, "Name").(string), Utility.GetProperty(info.Server, "Domain").(string), int32(Utility.GetProperty(info.Server, "Port").(int)))
		if err != nil {
			return nil, err
		}

		// At each call I will report the load of the server.
		stats, _ := load.Avg()
		load_info := &lbpb.LoadInfo{
			ServerInfo: &lbpb.ServerInfo{
				Id:     Utility.GetProperty(info.Server, "Id").(string),
				Name:   Utility.GetProperty(info.Server, "Name").(string),
				Domain: Utility.GetProperty(info.Server, "Domain").(string),
				Port:   int32(Utility.GetProperty(info.Server, "Port").(int)),
			},
			Load1:  stats.Load1,
			Load5:  stats.Load5,
			Load15: stats.Load15,
		}
		lb_client.ReportLoadInfo(load_info)
	}

	var result interface{}

	if Utility.GetProperty(info.Server, "Id") != nil {
		fmt.Println("421 execute request ", Utility.GetProperty(info.Server, "Id").(string))
	}
	result, err = handler(ctx, rqst)
	// Send log message.
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

func (l ServerStreamInterceptorStream) SendMsg(rqst interface{}) error {
	return l.inner.SendMsg(rqst)
}

/**
 * Here I will wrap the original stream into this one to get access to the original
 * rqst, so I can validate it resources.
 */
func (l ServerStreamInterceptorStream) RecvMsg(rqst interface{}) error {
	// First of all i will get the message.
	l.inner.RecvMsg(rqst)

	hasAccess := l.clientId == "sa" ||
		l.method == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo" ||
		l.method == "/admin.AdminService/DeployApplication"

	// Test if peer has access
	if !hasAccess && len(l.clientId) > 0 {
		hasAccess, _ = validateActionRequest(rqst, l.method, l.clientId, rbacpb.SubjectType_ACCOUNT, l.domain)
	}

	if !hasAccess && len(l.application) > 0 {
		hasAccess, _ = validateActionRequest(rqst, l.method, l.application, rbacpb.SubjectType_APPLICATION, l.domain)
	}

	if !hasAccess && len(l.peer) > 0 {
		hasAccess, _ = validateActionRequest(rqst, l.method, l.peer, rbacpb.SubjectType_PEER, l.domain)
	}

	if !hasAccess {
		err := errors.New("Permission denied to execute method " + l.method + " user:" + l.clientId + " domain:" + l.domain + " application:" + l.application)
		fmt.Println(err)
		log.Println("validation fail ", err)
		return err
	}

	return nil
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

	method := info.FullMethod

	var clientId string
	var err error
	// If the call come from a local client it has hasAccess
	hasAccess := false

	peer_, _ := peer.FromContext(stream.Context())
	address := peer_.Addr.String()
	address = address[0:strings.Index(address, ":")]
	if Utility.IsLocal(address) {
		hasAccess = true
	}

	if len(token) > 0 {
		clientId, _, _, err = ValidateToken(token)
		if err != nil {
			return err
		}
		if clientId == "sa" {
			hasAccess = true
		}
	}

	if Utility.GetProperty(srv, "Id") != nil {
		serverId := Utility.GetProperty(srv, "Id").(string)
		serverName := Utility.GetProperty(srv, "Name").(string)
		serverDomain := Utility.GetProperty(srv, "Domain").(string)
		serverPort := int32(Utility.GetProperty(srv, "Port").(int))
		// Set load balancing informations.
		lb_client, err := getLoadBalancingClient(domain, serverId, serverName, serverDomain, serverPort)
		if err != nil {
			return err
		}

		// At each call I will report the load of the server.
		stats, _ := load.Avg()
		load_info := &lbpb.LoadInfo{
			ServerInfo: &lbpb.ServerInfo{
				Id:     serverId,
				Name:   serverName,
				Domain: serverDomain,
				Port:   serverPort,
			},
			Load1:  stats.Load1,
			Load5:  stats.Load5,
			Load15: stats.Load15,
		}

		lb_client.ReportLoadInfo(load_info)
	}

	// needed by the admin.
	if application == "admin" ||
		method == "/services.PackageDiscovery/GetServicesDescriptor" ||
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
