package Interceptors

// TODO for the validation, use a map to store valid method/token/ressource/access
// the validation will be renew only if the token expire. And when a token expire
// the value in the map will be discard. That way it will put less charge on the server
// side.

import "context"

//import "fmt"
import "log"
import "google.golang.org/grpc"
import "github.com/davecourtois/Globular/ressource"
import "google.golang.org/grpc/metadata"
import "strings"
import "errors"
import "github.com/davecourtois/Globular/file/filepb"

import "github.com/davecourtois/Utility"

var (
	ressource_client *ressource.Ressource_Client
)

/**
 * Get a the local ressource client.
 */
func getRessourceClient(domain string) *ressource.Ressource_Client {

	if ressource_client == nil {
		ressource_client = ressource.NewRessource_Client(domain, "ressource")
	}

	return ressource_client
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
	key := Utility.GenerateUUID(token + method + path + Utility.ToString(permission))
	log.Println("---> key", key)

	// get access from remote source.
	hasAccess, err := getRessourceClient(domain).ValidateUserRessourceAccess(token, path, method, permission)
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
	key := Utility.GenerateUUID(applicationName + method + path + Utility.ToString(permission))
	log.Println("---> key", key)

	hasAccess, err := getRessourceClient(domain).ValidateApplicationRessourceAccess(applicationName, path, method, permission)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("Permission denied! for file " + path)
	}
	return nil
}

func ValidateUserAccess(domain string, token string, method string) (bool, error) {
	key := Utility.GenerateUUID(token + method)
	log.Println("---> key", key)

	hasAccess, err := getRessourceClient(domain).ValidateUserAccess(token, method)
	return hasAccess, err
}

func ValidateApplicationAccess(domain string, application string, method string) (bool, error) {
	key := Utility.GenerateUUID(application + method)
	log.Println("---> key", key)

	hasAccess, err := getRessourceClient(domain).ValidateApplicationAccess(application, method)
	return hasAccess, err
}

// That interceptor is use by all services except the ressource service who has
// it own interceptor.
func ServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// The token and the application id.
	var token string
	var application string
	var path string
	var domain string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println("-----> metadata found!")
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
		path = strings.Join(md["path"], "")
		domain = strings.Join(md["domain"], "")
	}
	log.Println("-----------> application:", application)
	log.Println("-----------> token:", token)
	log.Println("-----------> path:", path)
	log.Println("-----------> domain:", domain)

	method := info.FullMethod
	hasAccess := false

	var err error
	// Here is the list of method accessible by default.
	if method == "/admin.AdminService/GetConfig" ||
		method == "/event.EventService/Subscribe" ||
		method == "/event.EventService/UnSubscribe" ||
		method == "/event.EventService/OnEvent" ||
		method == "/event.EventService/Quit" ||
		method == "/event.EventService/Publish" ||
		method == "/dns.DnsService/SetA" ||
		method == "/ca.CertificateAuthority/GetCaCertificate" ||
		method == "/ca.CertificateAuthority/SignCertificate" ||
		method == "/services.ServiceDiscovery/FindServices" ||
		method == "/services.ServiceDiscovery/GetServiceDescriptor" ||
		method == "/services.ServiceDiscovery/GetServicesDescriptor" ||
		method == "/services.ServiceRepository/downloadBundle" ||
		method == "/persistence.PersistenceService/CreateConnection" ||
		method == "/persistence.PersistenceService/Find" ||
		method == "/persistence.PersistenceService/FindOne" ||
		method == "/persistence.PersistenceService/Count" {
		hasAccess = true
	}
	clientId, _, _ := ValidateToken(token)

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		hasAccess, _ = ValidateApplicationAccess(domain, application, method)
	}

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		hasAccess, _ = ValidateUserAccess(domain, token, method)
	}

	if !hasAccess {
		err := errors.New("Permission denied to execute method " + method)
		getRessourceClient(domain).Log(application, clientId, method, err)
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
				return nil, err
			}
		}

	} else if method != "/persistence.PersistenceService/CreateConnection" &&
		method != "/persistence.PersistenceService/FindOne" &&
		method != "/persistence.PersistenceService/Count" {
		// Here I will retreive the permission from the database if there is some...
		// the path will be found in the parameter of the method.
		permission, err := getRessourceClient(domain).GetActionPermission(method)
		if err == nil && permission != -1 {
			// So here I will try to validate each parameter that begin with a '/'
			if strings.HasPrefix(path, "/") {

				// I will test if the user has file permission.
				err = ValidateUserRessourceAccess(domain, token, path, method, permission)
				if err != nil {

					err = ValidateApplicationRessourceAccess(domain, application, path, method, permission)
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
	getRessourceClient(domain).Log(application, clientId, method, err)

	// Here depending of the request I will execute more actions.
	if err == nil {
		if method == "/file.FileService/CreateDir" {
			rqst := req.(*filepb.CreateDirRequest)
			err := getRessourceClient(domain).CreateDirPermissions(token, rqst.GetPath(), rqst.GetName())
			if err != nil {
				log.Println(err)
			}
		} else if method == "/file.FileService/Rename" {
			rqst := req.(*filepb.RenameRequest)
			err := getRessourceClient(domain).RenameFilePermission(rqst.GetPath(), rqst.GetOldName(), rqst.GetNewName())
			if err != nil {
				log.Println(err)
			}
		} else if method == "/file.FileService/DeleteFile" {
			rqst := req.(*filepb.DeleteFileRequest)
			err := getRessourceClient(domain).DeleteFilePermissions(rqst.GetPath())
			if err != nil {
				log.Println(err)
			}
		} else if method == "/file.FileService/DeleteDir" {
			rqst := req.(*filepb.DeleteDirRequest)
			err := getRessourceClient(domain).DeleteDirPermissions(rqst.GetPath())
			if err != nil {
				log.Println(err)
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

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Println("-----> metadata found!")
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
		path = strings.Join(md["path"], "")
		domain = strings.Join(md["domain"], "")
	}

	method := info.FullMethod
	hasAccess := false

	log.Println("-----------> method:", method)
	log.Println("------------------> application:", application)
	log.Println("------------------> token:", token)
	log.Println("------------------> path:", path)
	log.Println("------------------> domain:", domain)

	var err error
	if method == "/persistence.PersistenceService/Find" {
		hasAccess = true
	}

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		hasAccess, err = ValidateUserAccess(domain, token, method)
		if err != nil {
			return err
		}
	}

	if !hasAccess {
		return errors.New("Permission denied to execute method " + method)
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		hasAccess, err = ValidateApplicationAccess(domain, application, method)
		if err != nil {
			return err
		}
	}

	err = handler(srv, stream)

	// TODO find when the stream is closing and log only one time.
	//if err == io.EOF {
	// Send log event...
	clientId, _, _ := ValidateToken(token)
	getRessourceClient(domain).Log(application, clientId, method, err)
	//}

	if err != nil {
		return err
	}

	return nil
}
