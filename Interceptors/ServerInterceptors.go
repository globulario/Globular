package Interceptors

// TODO for the validation, use a map to store valid method/token/ressource/access
// the validation will be renew only if the token expire. And when a token expire
// the value in the map will be discard. That way it will put less charge on the server
// side.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"net"

	"github.com/davecourtois/Globular/file/filepb"
	"github.com/davecourtois/Globular/ressource"
	"github.com/davecourtois/Globular/storage/storage_store"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	//"google.golang.org/grpc/peer"
)

var (
	ressource_client *ressource.Ressource_Client

	// That will contain the permission in memory to limit the number
	// of ressource request...
	cache *storage_store.BigCache_store
)

/**
 * Get a the local ressource client.
 */
func getRessourceClient(domain string) (*ressource.Ressource_Client, error) {
	var err error
	if ressource_client == nil {
		ressource_client, err = ressource.NewRessource_Client(domain, "ressource")
		if err != nil {
			return nil, err
		}
	}

	return ressource_client, nil
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
	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidateUserRessourceAccess(token, path, method, permission)
	if err != nil {
		return err
	}

	if !hasAccess {
		return errors.New("Permission denied! for file " + path)
	}

	return nil
}

/**
 * Validate application ressource permission.
 */
func ValidateApplicationRessourceAccess(domain string, applicationName string, method string, path string, permission int32) error {

	// keep the values in the map for the lifetime of the token and validate it
	// from local map.
	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidateApplicationRessourceAccess(applicationName, path, method, permission)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("Permission denied! for ressource " + path)
	}

	return nil
}

/**
 * Validate peer ressouce permission.
 */
func ValidatePeerRessourceAccess(domain string, name string, method string, path string, permission int32) error {

	// keep the values in the map for the lifetime of the token and validate it
	// from local map.
	ressource_client, err := getRessourceClient(domain)
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

	ressource_client, err := getRessourceClient(domain)
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
 * Validate Application method access.
 */
func ValidateApplicationAccess(domain string, application string, method string) (bool, error) {
	key := Utility.GenerateUUID(application + method)

	_, err := getCache().GetItem(key)
	if err == nil {
		// Here a value exist in the store...
		return true, nil
	}

	ressource_client, err := getRessourceClient(domain)
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

func ValidatePeerAccess(domain string, name string, method string) (bool, error) {
	key := Utility.GenerateUUID(name + method)

	_, err := getCache().GetItem(key)
	if err == nil {
		// Here a value exist in the store...
		return true, nil
	}

	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return false, err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidatePeerAccess(name, method)
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
	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return "", err
	}

	return ressource_client.RefreshToken(token)
}

// That interceptor is use by all services except the ressource service who has
// it own interceptor.
func ServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// The token and the application id.
	var token string
	var application string
	var path string
	var domain string  // This is the target domain, the one use in TLS certificate.
	var peer_id string // The name of the peer

	//var ip string
	//var mac string

	// Get the caller ip address.
	//p, _ := peer.FromContext(ctx)
	//ip = p.Addr.String()

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
		// in case of ressource path.
		path = strings.Join(md["path"], "")
		domain = strings.Join(md["domain"], "")
		//mac = strings.Join(md["mac"], "")
	}

	// Here I will test if the
	method := info.FullMethod

	if len(domain) == 0 {
		return nil, errors.New("No domain was given for method call '" + method + "'")
	}

	// now I will test if the given domain fit the caller address.
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	// display the ips.
	fmt.Println(ips)
	//fmt.Println(ip)

	// If the call come from a local client it has hasAccess
	hasAccess := false // strings.HasPrefix(ip, "127.0.0.1") || strings.HasPrefix(ip, Utility.MyIP())

	// needed to get access to the system.
	if method == "/admin.AdminService/GetConfig" ||
		method == "/services.ServiceDiscovery/FindServices" ||
		method == "/services.ServiceDiscovery/FindServices/GetServiceDescriptor" ||
		method == "/services.ServiceDiscovery/FindServices/GetServicesDescriptor" ||
		method == "/services.ServiceDiscovery/FindServices/GetServicesDescriptor" ||
		method == "/dns.DnsService/GetA" || method == "/dns.DnsService/GetAAAA" ||
		method == "/ressource.RessourceService/Log" {
		hasAccess = true
	}

	var clientId string
	//var err error

	if len(token) > 0 {
		clientId, _, _, err = ValidateToken(token)
		if err != nil {
			return nil, err
		}
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		hasAccess, _ = ValidateApplicationAccess(domain, application, method)
	}

	// Test if peer has access
	if len(peer_id) > 0 && !hasAccess {
		hasAccess, _ = ValidatePeerAccess(domain, peer_id, method)
	}

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		hasAccess, _ = ValidateUserAccess(domain, token, method)
	}

	// Connect to the ressource services for the given domain.

	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		err := errors.New("Permission denied to execute method " + method + " user:" + clientId + " domain:" + domain + " application:" + application)
		fmt.Println(err)
		ressource_client.Log(application, clientId, method, err)
		return nil, err
	}

	// Now I will test file permission.
	if clientId != "sa" {

		// Here I will retreive the permission from the database if there is some...
		// the path will be found in the parameter of the method.
		permission, err := ressource_client.GetActionPermission(method)
		if err == nil && permission != -1 {
			// I will test if the user has file permission.
			err = ValidateUserRessourceAccess(domain, token, method, path, permission)
			if err != nil {
				if len(application) == 0 {
					return nil, err
				}
				err = ValidateApplicationRessourceAccess(domain, application, method, path, permission)
				if err != nil {
					if len(peer_id) == 0 {
						return nil, err
					}
					err = ValidatePeerRessourceAccess(domain, peer_id, method, path, permission)
					if err != nil {
						return nil, err
					}
				}
			}

		}

	}

	// Execute the action.
	result, err := handler(ctx, req)

	// Send log event...
	if (len(application) > 0 && len(clientId) > 0 && clientId != "sa") || err != nil {
		ressource_client.Log(application, clientId, method, err)
	}

	// Here depending of the request I will execute more actions.
	if err == nil {
		if method == "/file.FileService/CreateDir" && clientId != "sa" {
			rqst := req.(*filepb.CreateDirRequest)
			err := ressource_client.CreateDirPermissions(token, rqst.GetPath(), rqst.GetName())
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			// Here I will set the ressource owner for the directory.
			if strings.HasSuffix(rqst.GetPath(), "/") {
				ressource_client.SetRessourceOwner(clientId, rqst.GetPath()+rqst.GetName(), "")
			} else {
				ressource_client.SetRessourceOwner(clientId, rqst.GetPath()+"/"+rqst.GetName(), "")
			}

		} else if method == "/file.FileService/Rename" {
			rqst := req.(*filepb.RenameRequest)
			err := ressource_client.RenameFilePermission(rqst.GetPath(), rqst.GetOldName(), rqst.GetNewName())
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		} else if method == "/file.FileService/DeleteFile" {
			rqst := req.(*filepb.DeleteFileRequest)
			err := ressource_client.DeleteFilePermissions(rqst.GetPath())
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		} else if method == "/file.FileService/DeleteDir" {
			rqst := req.(*filepb.DeleteDirRequest)
			err := ressource_client.DeleteDirPermissions(rqst.GetPath())
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		}
	}

	return result, err

}

// Stream interceptor.
func ServerStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	// The token and the application id.
	var token string
	var application string
	var domain string
	var path string
	var peer_id string
	//var ip string

	//var mac string
	//Get the caller ip address.
	//p, _ := peer.FromContext(stream.Context())
	//ip = p.Addr.String()

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
		path = strings.Join(md["path"], "")
		//mac = strings.Join(md["mac"], "")
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		domain = strings.Join(md["domain"], "")
	}

	method := info.FullMethod

	var clientId string
	var err error

	if len(token) > 0 {
		clientId, _, _, err = ValidateToken(token)
		if err != nil {
			return err
		}
	}

	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return err
	}

	// If the call come from a local client it has hasAccess
	hasAccess := false // strings.HasPrefix(p.Addr.String(), "127.0.0.1") || strings.HasPrefix(ip, Utility.MyIP())
	// needed by the admin.
	if application == "admin" ||
		method == "/services.ServiceDiscovery/FindServices/GetServicesDescriptor" {
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

	if len(peer_id) > 0 && !hasAccess {
		hasAccess, _ = ValidatePeerAccess(domain, peer_id, method)
	}

	// Return here if access is denied.
	if !hasAccess {
		return errors.New("Permission denied to execute method " + method)
	}

	// Now the permissions
	if len(path) > 0 && clientId != "sa" {
		permission, err := ressource_client.GetActionPermission(method)
		if err == nil && permission != -1 {
			// I will test if the user has file permission.
			err = ValidateUserRessourceAccess(domain, token, method, path, permission)
			if err != nil {
				if len(application) == 0 {
					return err
				}
				err = ValidateApplicationRessourceAccess(domain, application, method, path, permission)
				if err != nil {
					if len(peer_id) == 0 {
						return err
					}
					err = ValidatePeerRessourceAccess(domain, peer_id, method, path, permission)
					if err != nil {
						return err
					}
				}
			}

		}
	} else if clientId != "sa" {
		permission, err := ressource_client.GetActionPermission(method)
		if err == nil && permission != -1 {
			return errors.New("Permission denied to execute method " + method)
		}
	}

	err = handler(srv, stream)

	// TODO find when the stream is closing and log only one time.
	//if err == io.EOF {
	// Send log event...
	ressource_client.Log(application, clientId, method, err)
	//}

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
