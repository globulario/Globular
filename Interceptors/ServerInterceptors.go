package Interceptors

// TODO for the validation, use a map to store valid method/token/ressource/access
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
	"github.com/globulario/services/golang/ressource/ressource_client"
	"github.com/globulario/services/golang/storage/storage_store"
	"github.com/shirou/gopsutil/load"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	// The ressource client
	ressource_client_ *ressource_client.Ressource_Client

	// The load balancer client.
	lb_client *load_balancing_client.Lb_Client

	// The map will contain connection with other server of same kind to load
	// balance the server charge.
	clients map[string]globular.Client

	// That will contain the permission in memory to limit the number
	// of ressource request...
	cache *storage_store.BigCache_store
)

/**
 * Get a the local ressource client.
 */
func getLoadBalancingClient(domain string, serverId string, serviceName string, serverDomain string, serverPort int32) (*load_balancing_client.Lb_Client, error) {

	var err error
	if lb_client == nil {
		lb_client, err = load_balancing_client.NewLbService_Client(domain, "lb.LoadBalancingService")
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

					lb_client.ReportLoadInfo(load_info)
				}
			}
		}()
	}

	return lb_client, nil
}

/**
 * Get a the local ressource client.
 */
func GetRessourceClient(domain string) (*ressource_client.Ressource_Client, error) {
	var err error
	if ressource_client_ == nil {
		ressource_client_, err = ressource_client.NewRessourceService_Client(domain, "ressource.RessourceService")
		if err != nil {
			return nil, err
		}
	}

	return ressource_client_, nil
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

/**
 * Validate user file permission.
 */
func ValidateUserRessourceAccess(domain string, token string, method string, path string, permission int32) error {

	// keep the values in the map for the lifetime of the token and validate it
	// from local map.
	ressource_client, err := GetRessourceClient(domain)
	if err != nil {
		return err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidateUserRessourceAccess(token, path, method, permission)
	if err != nil {
		return err
	}

	if !hasAccess {
		user, _, _, _ := ValidateToken(token)
		return errors.New("Permission denied for user " + user + " to execute methode " + method + " on ressource " + path)
	}

	return nil
}

/**
 * Validate application ressource permission.
 */
func ValidateApplicationRessourceAccess(domain string, applicationName string, method string, path string, permission int32) error {

	// keep the values in the map for the lifetime of the token and validate it
	// from local map.
	ressource_client, err := GetRessourceClient(domain)
	if err != nil {
		return err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidateApplicationRessourceAccess(applicationName, path, method, permission)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("Permission denied for application " + applicationName + " to execute method " + method + " on ressource " + path)
	}

	return nil
}

func ValidateUserAccess(domain string, token string, method string) (bool, error) {
	clientId, _, expire, err := ValidateToken(token)

	key := Utility.GenerateUUID(clientId + method)
	if err != nil || expire < time.Now().Unix() {
		getCache().RemoveItem(key)
	}

	_, err = getCache().GetItem(key)
	if err == nil {
		// Here a value exist in the store...
		return true, nil
	}

	ressource_client, err := GetRessourceClient(domain)
	if err != nil {
		return false, err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidateUserAccess(token, method)

	if hasAccess {
		getCache().SetItem(key, []byte(""))
	}

	return hasAccess, err
}

/**
 * Validate peer ressouce permission.
 */
func ValidatePeerRessourceAccess(domain string, name string, method string, path string, permission int32) error {

	// keep the values in the map for the lifetime of the token and validate it
	// from local map.
	ressource_client, err := GetRessourceClient(domain)
	if err != nil {
		return err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidatePeerRessourceAccess(name, path, method, permission)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("Permission denied! for ressource " + path)
	}

	return nil
}

func ValidatePeerAccess(domain string, peer string, method string) (bool, error) {

	key := Utility.GenerateUUID(peer + method)

	_, err := getCache().GetItem(key)
	if err == nil {
		// Here a value exist in the store...
		return true, nil
	}

	ressource_client, err := GetRessourceClient(domain)
	if err != nil {
		return false, err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidatePeerAccess(peer, method)
	if hasAccess {
		getCache().SetItem(key, []byte(""))

		// Here I will set a timeout for the permission.
		timeout := time.NewTimer(15 * time.Minute)
		go func() {
			<-timeout.C
			getCache().RemoveItem(key)
		}()
	}
	return hasAccess, err
}

/**
 * Validate Application method access.
 */
func ValidateApplicationAccess(domain string, application string, method string) (bool, error) {
	key := Utility.GenerateUUID(application + method)

	_, err := getCache().GetItem(key)
	if err == nil {
		// Here a value exist in the store...
		return true, nil
	}

	ressource_client, err := GetRessourceClient(domain)
	if err != nil {
		return false, err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidateApplicationAccess(application, method)
	if hasAccess {
		getCache().SetItem(key, []byte(""))

		// Here I will set a timeout for the permission.
		timeout := time.NewTimer(15 * time.Minute)
		go func() {
			<-timeout.C
			getCache().RemoveItem(key)
		}()
	}
	return hasAccess, err
}

// Refresh a token.
func refreshToken(domain string, token string) (string, error) {
	ressource_client, err := GetRessourceClient(domain)
	if err != nil {
		return "", err
	}

	return ressource_client.RefreshToken(token)
}

// That interceptor is use by all services except the ressource service who has
// it own interceptor.
func ServerUnaryInterceptor(ctx context.Context, rqst interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// The token and the application id.
	var token string
	var application string
	var path string
	var domain string // This is the target domain, the one use in TLS certificate.
	var load_balanced bool

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
		// in case of ressource path.
		path = strings.Join(md["path"], "")
		domain = strings.Join(md["domain"], "")
		if strings.Index(domain, ":") == 0 {
			port := strings.Join(md["port"], "")
			if len(port) > 0 {
				// set the address in that particular case.
				domain += ":" + port
			}
		}

		//TODO secure it
		load_balanced_ := strings.Join(md["load_balanced"], "")
		ctx = metadata.AppendToOutgoingContext(ctx, "load_balanced", "") // Set back the value to nothing.
		load_balanced = load_balanced_ == "true"

	}

	p, _ := peer.FromContext(ctx)
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
		method == "/ressource.RessourceService/Log" {
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
			log.Println("-----> 361 has access " + method + " user:" + clientId + " domain:" + domain + " application:" + application)
			hasAccess = true
		}
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		hasAccess, _ = ValidateApplicationAccess(domain, application, method)
	}

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		hasAccess, _ = ValidateUserAccess(domain, token, method)
	}

	// Test if peer has access
	if !hasAccess {
		// hasAccess, _ = ValidatePeerAccess(domain, "globular.io", method)
	}

	// Connect to the ressource services for the given domain.
	ressource_client, err := GetRessourceClient(domain)
	if err != nil {
		log.Println("fail to get ressource validator client ", err)
		return nil, err
	}

	log.Println("validate call from ", p.Addr.String(), "application", application, "domain", domain)
	log.Println("Validate access method with result: ", method, hasAccess)

	if !hasAccess {
		err := errors.New("Permission denied to execute method " + method + " user:" + clientId + " domain:" + domain + " application:" + application)
		fmt.Println(err)
		log.Println("validation fail ", err)
		ressource_client.Log(application, clientId, method, err)
		return nil, err
	}

	// Now I will test file permission.
	if clientId != "sa" {
		// Here I will retreive the permission from the database if there is some...
		// the path will be found in the parameter of the method.
		actionParameterRessourcesPermissions, err := ressource_client.GetActionPermission(method)
		if err == nil {
			val, _ := Utility.CallMethod(rqst, "ProtoReflect", []interface{}{})
			rqst_ := val.(protoreflect.Message)
			if rqst_.Descriptor().Fields().Len() > 0 {
				for i := 0; i < len(actionParameterRessourcesPermissions); i++ {
					permission := actionParameterRessourcesPermissions[i].Permission

					// Here I will get the paremeter that represent the path of a ressource.
					param := rqst_.Descriptor().Fields().Get(int(actionParameterRessourcesPermissions[i].Index))
					path, _ := Utility.CallMethod(rqst, "Get"+strings.ToUpper(string(param.Name())[0:1])+string(param.Name())[1:], []interface{}{})

					fmt.Println("validate ressource ", path, permission)

					// I will test if the user has file permission.
					err := ValidateUserRessourceAccess(domain, token, method, Utility.ToString(path), permission)
					if err != nil {
						if len(application) == 0 {
							return nil, err
						}
						err = ValidateApplicationRessourceAccess(domain, application, method, Utility.ToString(path), permission)
						if err != nil {
							err = ValidatePeerRessourceAccess(domain, "globular.io", method, Utility.ToString(path), permission)
							if err != nil {
								fmt.Println(err)
								return nil, err
							}
						}
					}

				}
			}
		}
	}

	// Here I will exclude local service from the load balancing.
	var candidates []*lbpb.ServerInfo
	// I will try to get the list of candidates for load balancing
	if Utility.GetProperty(info.Server, "Port") != nil {

		// Here I will refresh the load balance of the server to keep track of real load.
		lb_client, err := getLoadBalancingClient(domain, Utility.GetProperty(info.Server, "Id").(string), Utility.GetProperty(info.Server, "Name").(string), Utility.GetProperty(info.Server, "Domain").(string), int32(Utility.GetProperty(info.Server, "Port").(int)))
		if err != nil {
			return nil, err
		}

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

		// if load balanced is false I will get list of candidate.
		if load_balanced == false {
			candidates, _ = lb_client.GetCandidates(Utility.GetProperty(info.Server, "Name").(string))
		}

	}

	var result interface{}

	// Execute the action.
	if candidates != nil {
		serverId := Utility.GetProperty(info.Server, "Id").(string)
		// Here there is some candidate in the list.
		for i := 0; i < len(candidates); i++ {
			candidate := candidates[i]

			if candidate.GetId() == serverId {
				// In that case the handler is the actual server.
				result, err = handler(ctx, rqst)

				break // stop the loop...
			} else {
				// Here the canditade is the actual server so I will dispatch the request to the candidate.
				if clients[candidate.GetId()] == nil && len(method) > 1 {
					// Here I will create an instance of the client.
					newClientFct := method[1:strings.LastIndex(method, "/")]
					newClientFct = method[strings.Index(newClientFct, ".")+1:]
					newClientFct = "New" + newClientFct + "_Client"

					// Here I will create a connection with the other server in order to be able to dispatch the request.
					results, err := Utility.CallFunction(newClientFct, candidate.GetDomain(), candidate.GetId())
					if err != nil {
						fmt.Println(err)
						continue // skip to the next client.
					}
					// So here I will keep the client inside the map.
					clients[candidate.GetId()] = results[0].Interface().(globular.Client)
				}

				// Here I will invoke the request on the server whit the same context, so permission and token etc will be kept the save.
				result, err = clients[candidate.GetId()].Invoke(method, rqst, metadata.AppendToOutgoingContext(ctx, "load_balanced", "true", "domain", Utility.GetProperty(info.Server, "Domain").(string), "path", path, "application", application, "token", token))
				if err != nil {
					fmt.Println(err)
					continue // skip to the next client.
				} else {
					break
				}
			}
		}

	} else {
		result, err = handler(ctx, rqst)
	}

	// Send log event...
	if (len(application) > 0 && len(clientId) > 0 && clientId != "sa") || err != nil {
		ressource_client.Log(application, clientId, method, err)
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
 * rqst, so I can validate it ressources.
 */
func (l ServerStreamInterceptorStream) RecvMsg(rqst interface{}) error {

	var err error
	ressource_client, err := GetRessourceClient(l.domain)
	if err != nil {
		return err
	}

	if l.clientId != "sa" {
		val, _ := Utility.CallMethod(rqst, "ProtoReflect", []interface{}{})
		rqst_ := val.(protoreflect.Message)
		actionParameterRessourcesPermissions, err := ressource_client.GetActionPermission(l.method)
		if err == nil {
			// The token and the application id.
			if rqst_.Descriptor().Fields().Len() > 0 {
				for i := 0; i < len(actionParameterRessourcesPermissions); i++ {
					permission := actionParameterRessourcesPermissions[i].Permission
					// Here I will get the paremeter that represent the path of a ressource.
					param := rqst_.Descriptor().Fields().Get(int(actionParameterRessourcesPermissions[i].Index))
					path, _ := Utility.CallMethod(rqst, "Get"+strings.ToUpper(string(param.Name())[0:1])+string(param.Name())[1:], []interface{}{})

					// I will test if the user has file permission.
					err := ValidateUserRessourceAccess(l.domain, l.token, l.method, Utility.ToString(path), permission)
					if err != nil {
						if len(l.application) == 0 {
							ressource_client.Log("", l.clientId, l.method, err)
							return err
						}
						err = ValidateApplicationRessourceAccess(l.domain, l.application, l.method, Utility.ToString(path), permission)
						if err != nil {
							ressource_client.Log(l.application, l.clientId, l.method, err)
							err = ValidatePeerRessourceAccess(l.domain, l.peer, l.method, Utility.ToString(path), permission)
							if err != nil {
								ressource_client.Log("", l.clientId, l.method, err)
								return err
							}
							return err
						}

					}
				}
			}
		}
	}
	log.Println("---> 606 ", l.domain)
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
			log.Println("-----> 631 has access")
			hasAccess = true
		}
	}

	// If the call come from a local client it has hasAccess
	hasAccess := false //strings.HasPrefix(p.Addr.String(), "127.0.0.1") || strings.HasPrefix(p.Addr.String(), Utility.MyIP())

	// needed by the admin.
	if application == "admin" ||
		method == "/services.ServiceDiscovery/GetServicesDescriptor" ||
		method == "/services.ServiceRepository/DownloadBundle" {
		hasAccess = true
	}

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		hasAccess, _ = ValidateUserAccess(domain, token, method)
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		hasAccess, _ = ValidateApplicationAccess(domain, application, method)
	}

	if !hasAccess {
		// hasAccess, _ = ValidatePeerAccess(domain, "globular.io", method)
	}
	log.Println("validate call from ", p.Addr.String(), "application", application, "domain", domain)
	log.Println("Validate access method with result: ", method, hasAccess)
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
