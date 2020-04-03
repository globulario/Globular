package Interceptors

// TODO for the validation, use a map to store valid method/token/ressource/access
// the validation will be renew only if the token expire. And when a token expire
// the value in the map will be discard. That way it will put less charge on the server
// side.

import "context"
import "fmt"
import "log"
import "google.golang.org/grpc"
import "github.com/davecourtois/Globular/ressource"
import "google.golang.org/grpc/metadata"
import "strings"
import "errors"
import "github.com/davecourtois/Globular/file/filepb"

//import "google.golang.org/grpc/peer"
import "github.com/davecourtois/Utility"
import "github.com/davecourtois/Globular/storage/storage_store"
import "time"

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
			log.Println(err)
		}
	}
	return cache
}

/**
 * Return the file permission necessary for a given method.
 */
func getFilePermissionForMethod(method string, req interface{}) (string, int32) {
	var path string
	var permission int32

	if method == "/file.FileService/ReadDir" {
		rqst := req.(*filepb.ReadDirRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 4
	} else if method == "/file.FileService/CreateDir" {
		rqst := req.(*filepb.CreateDirRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 2
	} else if method == "/file.FileService/DeleteDir" {
		rqst := req.(*filepb.DeleteDirRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 1
	} else if method == "/file.FileService/Rename" {
		rqst := req.(*filepb.RenameRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath() + "/" + rqst.GetOldName()
		} else {
			path = rqst.GetOldName()
		}
		permission = 2
	} else if method == "/file.FileService/GetFileInfo" {
		rqst := req.(*filepb.ReadFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 4
	} else if method == "/file.FileService/ReadFile" {
		rqst := req.(*filepb.SaveFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 4
	} else if method == "/file.FileService/SaveFile" {
		rqst := req.(*filepb.SaveFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 2
	} else if method == "/file.FileService/DeleteFile" {
		rqst := req.(*filepb.DeleteFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 1
	} else if method == "/file.FileService/GetThumbnails" {
		rqst := req.(*filepb.GetThumbnailsRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 4
	} else if method == "/file.FileService/WriteExcelFile" {
		rqst := req.(*filepb.WriteExcelFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 2
	} else if method == "/file.FileService/CreateAchive" {
		rqst := req.(*filepb.CreateArchiveRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 4
	} else if method == "/ressource.RessourceService/SetPermission" {
		rqst := req.(*ressource.SetPermissionRqst)
		if len(rqst.Permission.GetPath()) > 1 {
			path = rqst.Permission.GetPath()
		}
		permission = 2
	} else if method == "/ressource.RessourceService/DeletePermissions" {
		rqst := req.(*ressource.DeletePermissionsRqst)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 1
	} else if method == "/ressource.RessourceService/SetRessourceOwner" {
		rqst := req.(*ressource.SetRessourceOwnerRqst)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 2
	} else if method == "/ressource.RessourceService/DeleteRessourceOwner" {
		rqst := req.(*ressource.DeleteRessourceOwnerRqst)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = 1
	}

	path = strings.ReplaceAll(path, "\\", "/")
	return path, permission

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
 * Validate application file permission.
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
		return errors.New("Permission denied! for file " + path)
	}

	return nil
}

func ValidateUserAccess(domain string, token string, method string) (bool, error) {
	clientId, expire, err := ValidateToken(token)
	key := Utility.GenerateUUID(clientId + method)
	if err != nil || expire < time.Now().Unix() {
		getCache().RemoveItem(key)
	}

	_, err = getCache().GetItem(key)
	if err == nil {
		fmt.Println("----> permission found ", domain, clientId, method)
		// Here a value exist in the store...
		return true, nil
	}

	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		fmt.Println("----> no ressource client found!", domain, err)
		return false, err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidateUserAccess(token, method)
	if hasAccess {
		getCache().SetItem(key, []byte(""))
		fmt.Println("----> permission saved ", domain, clientId, method)
	}

	return hasAccess, err
}

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

// That interceptor is use by all services except the ressource service who has
// it own interceptor.
func ServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// The token and the application id.
	var token string
	var application string
	var path string
	var domain string // This is the target domain, the one use in TLS certificate.
	//var ip string
	//var mac string

	// Get the caller ip address.
	//p, _ := peer.FromContext(ctx)
	//ip = p.Addr.String()

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
		path = strings.Join(md["path"], "")
		if !strings.HasPrefix("/", path) {
			path = "/" + path
		}
		domain = strings.Join(md["domain"], "")
		//mac = strings.Join(md["mac"], "")
	}

	method := info.FullMethod

	// If the call come from a local client it has hasAccess
	hasAccess := false //strings.HasPrefix(p.Addr.String(), "127.0.0.1") || strings.HasPrefix(ip, Utility.MyIP())

	// needed to get access to the system.
	if method == "/admin.AdminService/GetConfig" || method == "/dns.DnsService/GetA" || method == "/dns.DnsService/GetAAAA" || method == "/ressource.RessourceService/Log" {
		hasAccess = true
	}

	clientId, timeout, err := ValidateToken(token)

	if err == nil && timeout > time.Now().Unix() {
		if clientId == "sa" {
			// Execute the action.
			return handler(ctx, req)
		}
	} else if len(clientId) > 0 {
		// Here the error is not nil and the token is not a valid
		return nil, err
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		hasAccess, _ = ValidateApplicationAccess(domain, application, method)
	}

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		hasAccess, _ = ValidateUserAccess(domain, token, method)
	}

	// Connect to the ressource services for the given domain.
	log.Println("Validate access for ", clientId, application, domain, method)
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
	if strings.HasPrefix(method, "/file.FileService/") {
		path, permission := getFilePermissionForMethod(method, req)

		// I will test if the user has file permission.
		err = ValidateUserRessourceAccess(domain, token, path, method, permission)
		if err != nil {
			err = ValidateApplicationRessourceAccess(domain, application, path, method, permission)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}

	} else if method != "/persistence.PersistenceService/CreateConnection" && // Those method will make the server run in infinite loop
		method != "/persistence.PersistenceService/FindOne" &&
		method != "/persistence.PersistenceService/Count" && len(path) > 0 {
		// Here I will retreive the permission from the database if there is some...
		// the path will be found in the parameter of the method.
		log.Println("---> validate ressource permission ", path, " for method ", method)
		permission, err := ressource_client.GetActionPermission(method)
		if err == nil && permission != -1 {
			// I will test if the user has file permission.
			err = ValidateUserRessourceAccess(domain, token, path, method, permission)
			if err != nil {
				err = ValidateApplicationRessourceAccess(domain, application, path, method, permission)
				if err != nil {
					log.Println(err)
					return nil, err
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
		if method == "/file.FileService/CreateDir" {
			rqst := req.(*filepb.CreateDirRequest)
			err := ressource_client.CreateDirPermissions(token, rqst.GetPath(), rqst.GetName())
			if err != nil {
				log.Println(err)
				return nil, err
			}
		} else if method == "/file.FileService/Rename" {
			rqst := req.(*filepb.RenameRequest)
			err := ressource_client.RenameFilePermission(rqst.GetPath(), rqst.GetOldName(), rqst.GetNewName())
			if err != nil {
				log.Println(err)
				return nil, err
			}
		} else if method == "/file.FileService/DeleteFile" {
			rqst := req.(*filepb.DeleteFileRequest)
			err := ressource_client.DeleteFilePermissions(rqst.GetPath())
			if err != nil {
				log.Println(err)
				return nil, err
			}
		} else if method == "/file.FileService/DeleteDir" {
			rqst := req.(*filepb.DeleteDirRequest)
			err := ressource_client.DeleteDirPermissions(rqst.GetPath())
			if err != nil {
				log.Println(err)
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
		if !strings.HasPrefix("/", path) {
			path = "/" + path
		}
		domain = strings.Join(md["domain"], "")
	}

	method := info.FullMethod
	clientId, timeout, err := ValidateToken(token)

	if err == nil && timeout > time.Now().Unix() {
		if clientId == "sa" {
			// Execute the action.
			return handler(srv, stream)
		}
	} else if len(clientId) > 0 {
		// Here the error is not nil and the token is not a valid
		return err
	}

	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return err
	}

	// If the call come from a local client it has hasAccess
	hasAccess := false // strings.HasPrefix(p.Addr.String(), "127.0.0.1") || strings.HasPrefix(ip, Utility.MyIP())
	// needed by the admin.
	if application == "admin" && method == "/persistence.PersistenceService/Find" {
		hasAccess = true
	}

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		hasAccess, err = ValidateUserAccess(domain, token, method)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		hasAccess, err = ValidateApplicationAccess(domain, application, method)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Return here if access is denied.
	if !hasAccess {
		return errors.New("Permission denied to execute method " + method)
	}

	// Now the permissions
	if len(path) > 0 {
		permission, err := ressource_client.GetActionPermission(method)
		if err == nil && permission != -1 {
			// I will test if the user has file permission.
			err = ValidateUserRessourceAccess(domain, token, path, method, permission)
			if err != nil {
				err = ValidateApplicationRessourceAccess(domain, application, path, method, permission)
				if err != nil {
					log.Println(err)
					return err
				}
			}

		}
	}

	err = handler(srv, stream)

	// TODO find when the stream is closing and log only one time.
	//if err == io.EOF {
	// Send log event...
	ressource_client.Log(application, clientId, method, err)
	//}

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
