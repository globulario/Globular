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

//import "github.com/davecourtois/Utility"

var (
	ressource_client *ressource.Ressource_Client
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
	//key := Utility.GenerateUUID(token + method + path + Utility.ToString(permission))
	//log.Println("---> key", key)
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
	// key := Utility.GenerateUUID(applicationName + method + path + Utility.ToString(permission))
	// log.Println("---> key", key)
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
	// key := Utility.GenerateUUID(token + method)
	// log.Println("---> key", key)
	clientId, _, _ := ValidateToken(token)
	fmt.Println("--------> validate user access: ", domain, clientId, method)
	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return false, err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidateUserAccess(token, method)
	return hasAccess, err
}

func ValidateApplicationAccess(domain string, application string, method string) (bool, error) {
	// key := Utility.GenerateUUID(application + method)
	// log.Println("---> key", key)

	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return false, err
	}

	// get access from remote source.
	hasAccess, err := ressource_client.ValidateApplicationAccess(application, method)
	return hasAccess, err
}

// That interceptor is use by all services except the ressource service who has
// it own interceptor.
func ServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// The token and the application id.
	var token string
	var application string
	var path string
	var domain string // the domain of the ressource manager (where globular run)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
		path = strings.Join(md["path"], "")
		domain = strings.Join(md["domain"], "")
	}

	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return nil, err
	}

	method := info.FullMethod

	// If the call come from a local client it has hasAccess
	hasAccess := false // ip == Utility.MyLocalIP()

	// needed to get access to the system.
	if method == "/admin.AdminService/GetConfig" {
		hasAccess = true
	}

	clientId, _, err := ValidateToken(token)
	if err == nil {
		if clientId == "sa" {
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

	if !hasAccess {
		err := errors.New("Permission denied to execute method " + method)
		ressource_client, err := getRessourceClient(domain)
		if err != nil {
			return nil, err
		}

		ressource_client.Log(application, clientId, method, err)
		log.Println(err)
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
		method != "/persistence.PersistenceService/Count" {
		// Here I will retreive the permission from the database if there is some...
		// the path will be found in the parameter of the method.
		permission, err := ressource_client.GetActionPermission(method)
		if err == nil && permission != -1 {
			// So here I will try to validate each parameter that begin with a '/'
			if strings.HasPrefix(path, "/") {

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
	}

	// Execute the action.
	result, err := handler(ctx, req)

	// Send log event...
	ressource_client.Log(application, clientId, method, err)

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
	//var path string

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
		//path = strings.Join(md["path"], "")
		domain = strings.Join(md["domain"], "")
	}

	ressource_client, err := getRessourceClient(domain)
	if err != nil {
		return err
	}

	method := info.FullMethod
	hasAccess := false

	if method == "/persistence.PersistenceService/Find" {
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

	if !hasAccess {
		return errors.New("Permission denied to execute method " + method)
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		hasAccess, err = ValidateApplicationAccess(domain, application, method)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	err = handler(srv, stream)

	// TODO find when the stream is closing and log only one time.
	//if err == io.EOF {
	// Send log event...
	clientId, _, _ := ValidateToken(token)
	ressource_client.Log(application, clientId, method, err)
	//}

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
