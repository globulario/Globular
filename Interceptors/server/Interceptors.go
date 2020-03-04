package Interceptors

import (
	"fmt"
	"log"

	"strings"

	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/davecourtois/Globular/Interceptors/Authenticate"
	//"github.com/davecourtois/Globular/admin"
	"github.com/davecourtois/Globular/file/filepb"
	"github.com/davecourtois/Globular/persistence/persistence_client"
	"github.com/davecourtois/Globular/ressource"
	"github.com/davecourtois/Utility"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// private type for Context keys
type contextKey int

const (
	clientIDKey contextKey = iota
)

var (
	client   *persistence_client.Persistence_Client
	rootPath string
)

/**
 * Get the persistence servcice connection.
 * That connection is used to retreive roles and application information.
 */
func getPersistenceClient() (*persistence_client.Persistence_Client, error) {
	// Here I will need the persistence client to read user permission.
	// Here I will read the server token, the service must run on the
	// same computer as globular.
	if client == nil {
		// The root password to be able to perform query over persistence service.
		infoStr, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "globular_sa")
		if err != nil {
			return nil, err
		}

		infos := make(map[string]interface{}, 0)
		err = json.Unmarshal(infoStr, &infos)
		if err != nil {
			return nil, err
		}

		// Local to the server so the information will be taken from
		// information in the file.
		root := infos["pwd"].(string)

		// Get the root path from the tmp file.
		rootPath = strings.ReplaceAll(infos["rootPath"].(string), "\\", "/")

		// close the
		if client != nil {
			client.Close()
		}

		// Use the client sa connection.
		client = persistence_client.NewPersistence_Client(infos["address"].(string), infos["name"].(string))
		err = client.CreateConnection("local_ressource", "local_ressource", "localhost", 27017, 0, "sa", root, 5000, "", false)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func ValidateToken(token string) (string, int64, error) {

	// Initialize a new instance of `Claims`
	claims := &Interceptors.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Get the key from the local temp file.
		jwtKey, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "globular_key")
		return jwtKey, err
	})

	if err != nil {
		return "", 0, err
	}

	if !tkn.Valid {
		return "", 0, fmt.Errorf("invalid token!")
	}

	return claims.Username, claims.ExpiresAt, nil
}

// authenticateAgent check the client credentials
func authenticateClient(ctx context.Context) (string, string, int64, error) {
	var userId string
	var applicationId string
	var expired int64
	var err error

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		applicationId = strings.Join(md["application"], "")
		token := strings.Join(md["token"], "")
		// In that case no token was given...
		if len(token) > 0 {
			userId, expired, err = ValidateToken(token)
		}
		return applicationId, userId, expired, err
	}

	return "", "", 0, fmt.Errorf("missing credentials")
}

// Test if a role can use action.
func canRunAction(roleName string, method string) error {

	Id := "local_ressource"
	Database := "local_ressource"
	Collection := "Roles"
	Query := `{"_id":"` + roleName + `"}`
	client, err := getPersistenceClient()
	if err != nil {
		return err
	}

	values, err := client.FindOne(Id, Database, Collection, Query, `[{"Projection":{"actions":1}}]`)
	if err != nil {
		return err
	}

	role := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &role)
	if err != nil {
		return err
	}

	// append all action into the actions
	for i := 0; i < len(role["actions"].([]interface{})); i++ {
		if role["actions"].([]interface{})[i].(string) == method {
			return nil
		}
	}

	// Here I will test if the user has write to execute the methode.
	return errors.New("Permission denied!")
}

/**
 * Return the file permission (unix number) necessary for a given method.
 */
func getFilePermissionForMethod(method string, req interface{}) (string, string) {
	var path string
	var permission string

	if method == "/file.FileService/ReadDir" {
		rqst := req.(*filepb.ReadDirRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "read"
	} else if method == "/file.FileService/CreateDir" {
		rqst := req.(*filepb.CreateDirRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "write"
	} else if method == "/file.FileService/DeleteDir" {
		rqst := req.(*filepb.DeleteDirRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "write"
	} else if method == "/file.FileService/Rename" {
		rqst := req.(*filepb.RenameRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath() + "/" + rqst.GetOldName()
		} else {
			path = rqst.GetOldName()
		}

		permission = "write"
	} else if method == "/file.FileService/GetFileInfo" {
		rqst := req.(*filepb.ReadFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "read"
	} else if method == "/file.FileService/ReadFile" {
		rqst := req.(*filepb.SaveFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "read"
	} else if method == "/file.FileService/SaveFile" {
		rqst := req.(*filepb.SaveFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "write"
	} else if method == "/file.FileService/DeleteFile" {
		rqst := req.(*filepb.DeleteFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}

		permission = "write"
	} else if method == "/file.FileService/GetThumbnails" {
		rqst := req.(*filepb.GetThumbnailsRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "read"
	} else if method == "/file.FileService/WriteExcelFile" {
		rqst := req.(*filepb.WriteExcelFileRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "write"
	} else if method == "/file.FileService/CreateAchive" {
		rqst := req.(*filepb.CreateArchiveRequest)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "read"
	} else if method == "/ressource.RessourceService/SetPermission" {
		rqst := req.(*ressource.SetPermissionRqst)
		if len(rqst.Permission.GetPath()) > 1 {
			path = rqst.Permission.GetPath()
		}
		permission = "write"
	} else if method == "/ressource.RessourceService/DeletePermissions" {
		rqst := req.(*ressource.DeletePermissionsRqst)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "write"
	} else if method == "/ressource.RessourceService/SetRessourceOwner" {
		rqst := req.(*ressource.SetRessourceOwnerRqst)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "write"
	} else if method == "/ressource.RessourceService/DeleteRessourceOwner" {
		rqst := req.(*ressource.DeleteRessourceOwnerRqst)
		if len(rqst.GetPath()) > 1 {
			path = rqst.GetPath()
		}
		permission = "write"
	}

	// make sure the path does not contain // anywhere...
	if len(path) == 1 {
		path = rootPath
	} else if strings.HasPrefix(path, "/") {
		path = rootPath + path
	} else {
		path = rootPath + "/" + path
	}

	path = strings.ReplaceAll(path, "\\", "/")
	return path, permission

}

func isOwner(name string, path string) bool {
	// get the client...
	client, _ := getPersistenceClient()

	// Now I will get the user roles and validate if the user can execute the
	// method.
	Id := "local_ressource"
	Database := "local_ressource"
	Collection := "Accounts"
	Query := `{"name":"` + name + `"}`

	values, err := client.FindOne(Id, Database, Collection, Query, `[{"Projection":{"roles":1}}]`)
	if err != nil {
		return false
	}

	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &account)
	if err != nil {
		return false
	}

	path = strings.ReplaceAll(path, "\\", "/")
	// If the user is the owner of the ressource it has the permission
	log.Println("---> ", `{"path":"`+path+`","owner":"`+name+`"}`)
	count, err := client.Count("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`","owner":"`+account["_id"].(string)+`"}`, ``)
	if err == nil {
		if count > 0 {
			log.Println("--> ", name, " is owner of ", path)
			return true
		}
	} else {
		log.Println(err)
	}
	return false
}

func hasPermission(name string, path string, permission string) (bool, int) {
	// Set the path with / instead of \\ in case of windows...
	path = strings.ReplaceAll(path, "\\", "/")
	log.Println("--> validate " + name + " has " + permission + " permission on file " + path)

	// If the user is the owner of the ressource it has all permission
	if isOwner(name, path) {
		return true, 0
	}

	count, err := client.Count("local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, ``)
	if err != nil {
		return false, 0
	}

	permissionsStr, err := client.Find("local_ressource", "local_ressource", "Permissions", `{"owner":"`+name+`", "path":"`+path+`"}`, ``)
	if err == nil {
		permissions := make([]map[string]interface{}, 0)
		json.Unmarshal([]byte(permissionsStr), &permissions)
		if len(permissions) == 0 {
			return false, count
		}
		for i := 0; i < len(permissions); i++ {
			permission_ := int32(Utility.ToInt(permissions[i]["permission"]))
			if permission == "read" {
				if permission_ > 3 {
					return true, count
				}
			} else if permission == "write" {
				if permission_ == 2 || permission_ == 3 || permission_ == 6 || permission_ == 7 {
					return true, count
				}
			} else if permission == "execute" {
				if permission_ == 1 || permission_ == 3 || permission_ == 7 {
					return true, count
				}
			}
		}

		return false, count
	}

	return false, count
}

func ValidateUserFileAccess(token string, method string, path string, permission string) error {
	// first of all I will validate the token.
	clientId, expiredAt, err := ValidateToken(token)
	if err != nil {
		return err
	}

	if expiredAt < time.Now().Unix() {
		return errors.New("The token is expired!")
	}

	return validateUserFileAccess(clientId, method, path, permission)
}

/**
 * Validate application file permission.
 */
func ValidateApplicationFileAccess(applicationName string, method string, path string, permission string) error {

	hasApplicationPermission, count := hasPermission(applicationName, path, permission)
	if hasApplicationPermission {
		return nil
	}

	// If theres permission definied and we are here it's means the user dosent
	// have write to execute the action on the ressource.
	if count > 0 {
		return errors.New("Permission Denied for " + applicationName)
	}

	return nil
}

/**
 * Validate if a user, a role or an application has write to do operation on a file or a directorty.
 */
func validateUserFileAccess(userName string, method string, path string, permission string) error {
	path = strings.ReplaceAll(path, "\\", "/")
	if !strings.HasPrefix(method, "/file.FileService") {
		return nil
	}
	log.Println("--> validate file access for ", userName, " with method ", method, " on file ", path, " and acess ", permission)
	if len(userName) == 0 {
		return errors.New("No user  name was given to validate method access " + method)
	}

	// if the user is the super admin no validation is required.
	if userName == "sa" {
		return nil
	}

	Id := "local_ressource"
	Database := "local_ressource"
	Collection := "Accounts"
	Query := `{"name":"` + userName + `"}`

	client, err := getPersistenceClient()
	if err != nil {
		return err
	}

	// Find the user role.
	values, err := client.FindOne(Id, Database, Collection, Query, `[{"Projection":{"roles":1}}]`)
	if err != nil {
		return err
	}

	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &account)
	if err != nil {
		return err
	}

	count := 0

	hasUserPermission, hasUserPermissionCount := hasPermission(userName, path, permission)
	if hasUserPermission {
		return nil
	}

	count += hasUserPermissionCount
	roles := account["roles"].([]interface{})
	for i := 0; i < len(roles); i++ {
		role := roles[i].(map[string]interface{})
		hasRolePermission, hasRolePermissionCount := hasPermission(role["$id"].(string), path, permission)
		count += hasRolePermissionCount
		if hasRolePermission {
			return nil
		}
	}

	// If theres permission definied and we are here it's means the user dosent
	// have write to execute the action on the ressource.
	if count > 0 {
		return errors.New("Permission Denied for " + userName)
	}

	return nil
}

/**
 * Validate user access by role
 */
func validateUserAccess(userName string, method string) error {

	if len(userName) == 0 {
		return errors.New("No user  name was given to validate method access " + method)
	}

	// if the user is the super admin no validation is required.
	if userName == "sa" {
		return nil
	}

	// if guest can run the action...
	if canRunAction("guest", method) == nil {
		// everybody can run the action in that case.
		return nil
	}

	client, err := getPersistenceClient()
	if err != nil {
		return err
	}

	// Now I will get the user roles and validate if the user can execute the
	// method.
	Id := "local_ressource"
	Database := "local_ressource"
	Collection := "Accounts"
	Query := `{"name":"` + userName + `"}`

	values, err := client.FindOne(Id, Database, Collection, Query, `[{"Projection":{"roles":1}}]`)
	if err != nil {
		log.Println(err)
		return err
	}

	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &account)
	if err != nil {
		log.Println(err)
		return err
	}

	roles := account["roles"].([]interface{})
	for i := 0; i < len(roles); i++ {
		role := roles[i].(map[string]interface{})
		if canRunAction(role["$id"].(string), method) == nil {
			return nil
		}
	}

	err = errors.New("permission denied! account " + userName + " cannot execute methode '" + method + "'")

	return err
}

/**
 * Log error in database.
 */
func logError(ctx context.Context, method string, err error) {
	// The name of the applicaition.
	application := "undefined"
	userId := "undefined"

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token := strings.Join(md["token"], "")
		userId, _, _ = ValidateToken(token)
	}
	saveInfo(application, userId, method, err)
}

func saveInfo(application string, userId string, method string, err_ error) error {

	p, err := getPersistenceClient()
	if err != nil {
		return err
	}

	info := make(map[string]interface{})
	info["application"] = application
	info["userId"] = userId
	info["method"] = method
	info["date"] = time.Now().Unix() // save it as unix time.if

	db := "Logs"
	if err_ != nil {
		db = "Errors"
		info["error"] = err_.Error()
	}

	data, _ := json.Marshal(&info)
	_, err = p.InsertOne("local_ressource", "local_ressource", db, string(data), "")

	if err != nil {
		return err
	}

	return nil
}

/**
 * Here I will log application actions. That's can be usefull to monitor
 * application utilisation.
 */
func logAction(ctx context.Context, method string, result interface{}) {

	var application string
	var userId string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token := strings.Join(md["token"], "")
		userId, _, _ = ValidateToken(token)
	}

	// Here I will save only some methode informations...
	if method == "/ressource.RessourceService/RegisterAccount" {
		token := result.(*ressource.RegisterAccountRsp).Result
		userId, _, _ = ValidateToken(token)
		saveInfo(application, userId, method, nil)
	} else if method == "/ressource.RessourceService/Authenticate" {
		token := result.(*ressource.AuthenticateRsp).Token
		userId, _, _ = ValidateToken(token)
		saveInfo(application, userId, method, nil)
	} else if method == "/admin.AdminService/DeployApplication" || method == "/admin.AdminService/PublishService" || method == "/admin.AdminService/RegisterExternalApplication" || method == "/admin.AdminService/UninstallService" || method == "/admin.AdminService/InstallService" {
		saveInfo(application, userId, method, nil)
	}
}

/*
type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {

	switch v := m.(type) {
	case *admin.DeployApplicationRequest:
		log.Println("----> 601 name ", v.GetName(), " data ", len(v.GetData()))
	}
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	return nil
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}
*/

// Stream interceptor.
func StreamAuthInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// Call 'handler' to invoke the stream handler before this function returns
	method := info.FullMethod
	applicationID, clientID, _, err := authenticateClient(stream.Context())
	log.Println("---> validate method", method, applicationID, clientID)

	if err != nil {
		return err
	}

	// Here I will test if the user can run that function or not...
	err = validateUserAccess(clientID, method)
	if err != nil {
		return err
	}

	err = handler(srv, stream)
	if err != nil {
		return err
	}

	if method == "/admin.AdminService/DeployApplication" {
		log.Println("------------> deploy application: ", applicationID)
		// Now I here I will validate that the ClientID has write access
		// or is the owner of the application.
		isOwner_ := isOwner(clientID, applicationID)
		if isOwner_ {
			return nil
		}

		// Now I will test if the user has application write permission.

		if err == io.EOF {
			return nil
		} else {
			return err
		}
	}

	return nil
}

// unaryInterceptor calls authenticateClient with current context
func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	applicationID, clientID, _, err := authenticateClient(ctx)
	if err != nil {
		log.Println("---> ", err)
		return nil, err
	}

	method := info.FullMethod

	// Validate the user access.
	log.Println("---> validate method", method, " for ", clientID)

	if len(applicationID) > 0 {
		// TODO validate application action here.
		// log.Println("---> validate application permission: ", applicationID)
	}

	if len(clientID) > 0 {
		path, permission := getFilePermissionForMethod(method, req)
		// Test if the user is owner...
		isOwner_ := false
		if strings.HasPrefix(method, "/file.FileService/") || method == "/ressource.RessourceService/SetPermission" ||
			method == "/ressource.RessourceService/DeletePermissions" || method == "/ressource.RessourceService/SetRessourceOwner" ||
			method == "/ressource.RessourceService/DeleteRessourceOwner" {
			isOwner_ = isOwner(clientID, path)
			log.Println("---> is owner?", clientID, path, isOwner_)
		}

		if !isOwner_ {
			// Validate the user access
			err = validateUserAccess(clientID, method)
			if err != nil {
				return nil, err
			}
			// Validate file access
			err = validateUserFileAccess(clientID, method, path, permission)
			if err != nil {
				return nil, err
			}
		}
	}

	// Execute the action.
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	result, err := handler(ctx, req)

	if err != nil {
		logError(ctx, method, err)
		return nil, err
	}

	// Log the action as needed for info.
	logAction(ctx, method, result)

	// Here I will set permission depending of actions...
	client, err := getPersistenceClient()
	if err != nil {
		return nil, err
	}

	var path string // must be set call after calling getPersistClient

	if method == "/file.FileService/CreateDir" {

		// A new directory will take the parent permissions by default...
		rqst := req.(*filepb.CreateDirRequest)
		path += rqst.GetPath()
		if len(path) > 1 {
			if strings.HasPrefix(path, "/") {
				path = rootPath + path
			} else {
				path = rootPath + "/" + path
			}
		} else {
			path = rootPath
		}

		permissionsStr, err := client.Find("local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, "")
		if err != nil {
			return nil, err
		}

		// Create permission object.
		permissions := make([]interface{}, 0)
		err = json.Unmarshal([]byte(permissionsStr), &permissions)
		if err != nil {
			return nil, err
		}

		// Now I will create the new permission of the created directory.
		for i := 0; i < len(permissions); i++ {
			// Copye the permission.
			permission := permissions[i].(map[string]interface{})
			permission_ := make(map[string]interface{}, 0)
			permission_["owner"] = permission["owner"]
			permission_["path"] = path + "/" + rqst.GetName()
			permission_["permission"] = permission["permission"]
			permissionStr, _ := Utility.ToJson(permission_)
			client.InsertOne("local_ressource", "local_ressource", "Permissions", permissionStr, "")
		}

		ressourceOwnersStr, err := client.Find("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, "")
		if err != nil {
			return nil, err
		}

		// Create permission object.
		ressourceOwners := make([]interface{}, 0)
		err = json.Unmarshal([]byte(ressourceOwnersStr), &ressourceOwners)
		if err != nil {
			return nil, err
		}

		// Now I will create the new permission of the created directory.
		for i := 0; i < len(ressourceOwners); i++ {
			// Copye the permission.
			ressourceOwner := ressourceOwners[i].(map[string]interface{})
			ressourceOwner_ := make(map[string]interface{}, 0)
			ressourceOwner_["owner"] = ressourceOwner["owner"]
			ressourceOwner_["path"] = path + "/" + rqst.GetName()
			ressourceOwnerStr, _ := Utility.ToJson(ressourceOwner_)
			client.InsertOne("local_ressource", "local_ressource", "RessourceOwners", ressourceOwnerStr, "")
		}

		// The user who create a directory will be the owner of the
		// directory.
		if clientID != "sa" && clientID != "guest" {
			ressourceOwner := make(map[string]interface{}, 0)
			ressourceOwner["owner"] = clientID
			ressourceOwner["path"] = path + "/" + rqst.GetName()
			ressourceOwnerStr, _ := Utility.ToJson(ressourceOwner)
			client.InsertOne("local_ressource", "local_ressource", "RessourceOwners", ressourceOwnerStr, `[{"upsert":true}]`)
		}

	} else if method == "/file.FileService/Rename" {
		rqst := req.(*filepb.RenameRequest)

		path := rqst.GetPath()
		path = strings.ReplaceAll(path, "\\", "/")
		oldPath := rqst.OldName
		newPath := rqst.NewName

		if strings.HasPrefix(path, string(os.PathSeparator)) {
			if len(path) > 1 {
				oldPath = path + "/" + rqst.OldName
				newPath = path + "/" + rqst.NewName
			} else {
				oldPath = rqst.OldName
				newPath = rqst.NewName
			}
		}

		err := client.Update("local_ressource", "local_ressource", "Permissions", `{"path":"`+rootPath+"/"+oldPath+`"}`, `{"$set":{"path":"`+rootPath+"/"+newPath+`"}}`, "")
		if err != nil {
			log.Println(err)
		}

		err = client.Update("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+rootPath+"/"+oldPath+`"}`, `{"$set":{"path":"`+rootPath+"/"+newPath+`"}}`, "")
		if err != nil {
			log.Println(err)
		}

		// Replace all files of subdirectories.
		permissionsStr, err := client.Find("local_ressource", "local_ressource", "Permissions", `{"path":{"$regex":"/.*`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"}}`, "")
		if err == nil {
			permissions := make([]interface{}, 0)
			json.Unmarshal([]byte(permissionsStr), &permissions)
			for i := 0; i < len(permissions); i++ { // stringnify and save it...
				permission := permissions[i].(map[string]interface{})
				err := client.Update("local_ressource", "local_ressource", "Permissions", `{"path":"`+permission["path"].(string)+`"}`, `{"$set":{"path":"`+strings.ReplaceAll(permission["path"].(string), oldPath, newPath)+`"}}`, "")
				if err != nil {
					log.Println(err)
				}

				err = client.Update("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+permission["path"].(string)+`"}`, `{"$set":{"path":"`+strings.ReplaceAll(permission["path"].(string), oldPath, newPath)+`"}}`, "")
				if err != nil {
					log.Println(err)
				}
			}
		}

	} else if method == "/file.FileService/DeleteFile" {
		rqst := req.(*filepb.DeleteFileRequest)
		path := rqst.GetPath()
		path = strings.ReplaceAll(path, "\\", "/")
		if len(path) > 1 {
			if strings.HasPrefix(path, "/") {
				path = rootPath + path
			} else {
				path = rootPath + "/" + path
			}
		} else {
			path = rootPath
		}

		err = client.Delete("local_ressource", "local_ressource", "Permissions", `{"path":"`+strings.ReplaceAll(path, "\\", "/")+`"}`, "")
		if err != nil {
			log.Println(err)
		}

		err = client.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+strings.ReplaceAll(path, "\\", "/")+`"}`, "")
		if err != nil {
			log.Println(err)
		}

	} else if method == "/file.FileService/DeleteDir" {
		rqst := req.(*filepb.DeleteDirRequest)
		path += rqst.GetPath()

		path = strings.ReplaceAll(path, "\\", "/")
		path = strings.ReplaceAll(path, "/", "\\/") // replace . by \. and / by \/
		path = strings.ReplaceAll(path, ".", "\\.") // TODO fix the nasty bug  with regex.

		// Delete Permissions
		err = client.Delete("local_ressource", "local_ressource", "Permissions", `{"path":{"$regex":"/.*`+path+`.*/"}}`, "")
		if err != nil {
			log.Println(err)
		}

		err = client.Delete("local_ressource", "local_ressource", "Permissions", `{"path":"`+rootPath+path+`"}`, "")
		if err != nil {
			log.Println(err)
		}

		// Delete Owners
		err = client.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":{"$regex":"/.*`+path+`.*/"}}`, "")
		if err != nil {
			log.Println(err)
		}

		err = client.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+rootPath+path+`"}`, "")
		if err != nil {
			log.Println(err)
		}
	} else if method == "/ressource.RessourceService/RemoveApplicationAction" {

		rqst := req.(*ressource.RemoveApplicationActionRqst)
		applicationStr, _ := client.FindOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, "")
		application := make(map[string]interface{})
		json.Unmarshal([]byte(applicationStr), &application)
		path := application["path"].(string)
		path = strings.ReplaceAll(path, "\\", "/")
		path = strings.ReplaceAll(path, "/", "\\/") // replace . by \. and / by \/
		path = strings.ReplaceAll(path, ".", "\\.") // TODO fix the nasty bug  with regex.

		// Delete file permissions
		// Delete Permissions
		err = client.Delete("local_ressource", "local_ressource", "Permissions", `{"path":{"$regex":"/.*`+path+`.*/"}}`, "")
		if err != nil {
			log.Println(err)
		}

		err = client.Delete("local_ressource", "local_ressource", "Permissions", `{"path":"`+rootPath+path+`"}`, "")
		if err != nil {
			log.Println(err)
		}

		// Delete Owners
		err = client.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":{"$regex":"/.*`+path+`.*/"}}`, "")
		if err != nil {
			log.Println(err)
		}

		err = client.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+rootPath+path+`"}`, "")
		if err != nil {
			log.Println(err)
		}
	}

	return result, nil
}
