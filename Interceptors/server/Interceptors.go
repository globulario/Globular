package Interceptors

import (
	"fmt"
	"log"

	"strings"

	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/davecourtois/Globular/Interceptors/Authenticate"
	"github.com/davecourtois/Globular/file/filepb"
	"github.com/davecourtois/Globular/persistence/persistence_client"
	"github.com/davecourtois/Globular/ressource"
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
	client *persistence_client.Persistence_Client
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

		// close the
		if client != nil {
			client.Close()
		}

		// Use the client sa connection.
		client = persistence_client.NewPersistence_Client(infos["address"].(string), infos["name"].(string))

		err = client.CreateConnection("local_ressource", "local_ressource", "localhost", 27017, 0, "sa", root, 5000, "", false)
		if err != nil {
			log.Println(`--> Fail to create  the connection "local_ressource"`)
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
func authenticateClient(ctx context.Context) (string, int64, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		token := strings.Join(md["token"], "")
		// In that case no token was given...
		if len(token) == 0 {
			return "", 0, nil
		}

		return ValidateToken(token)
	}

	return "", 0, fmt.Errorf("missing credentials")
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
 * Return the file permission necessary for a given method.
 */
func getFilePermissionForMethod(method string, req interface{}) (string, int32) {
	var path string
	var permission int32
	if method == "/file.FileService/ReadDir" {
		rqst := req.(filepb.ReadDirRequest)
		path = rqst.GetPath()
		permission = 4
	} else if method == "/file.FileService/CreateDir" {
		rqst := req.(filepb.CreateDirRequest)
		path = rqst.GetPath()
		permission = 2
	} else if method == "/file.FileService/DeleteDir" {
		rqst := req.(filepb.DeleteDirRequest)
		path = rqst.GetPath()
		permission = 2
	} else if method == "/file.FileService/Rename" {
		rqst := req.(filepb.RenameRequest)
		path = rqst.GetPath() + "/" + rqst.GetOldName()
		permission = 2
	} else if method == "/file.FileService/GetFileInfo" {
		rqst := req.(filepb.ReadFileRequest)
		path = rqst.GetPath()
		permission = 4
	} else if method == "/file.FileService/ReadFile" {
		rqst := req.(filepb.SaveFileRequest)
		path = rqst.GetPath()
		permission = 4
	} else if method == "/file.FileService/SaveFile" {
		rqst := req.(filepb.SaveFileRequest)
		path = rqst.GetPath()
		permission = 2
	} else if method == "/file.FileService/DeleteFile" {
		rqst := req.(filepb.DeleteFileRequest)
		path = rqst.GetPath()
		permission = 2
	} else if method == "/file.FileService/GetThumbnails" {
		rqst := req.(filepb.GetThumbnailsRequest)
		path = rqst.GetPath()
		permission = 4
	} else if method == "/file.FileService/WriteExcelFile" {
		rqst := req.(filepb.WriteExcelFileRequest)
		path = rqst.GetPath()
		permission = 2
	}

	return path, permission
}

/**
 * Validate if a user or a role has write to do operation on a file or a directorty.
 */
func validateFileAccess(userName string, method string, req interface{}) error {
	if !strings.HasPrefix(method, "/file.FileService") {
		return nil
	}

	// if the user is the super admin no validation is required.
	if userName == "sa" {
		return nil
	}

	// Now I will get the user roles and validate if the user can execute the
	// method.
	Id := "local_ressource"
	Database := "local_ressource"
	Collection := "Accounts"
	Query := `{"name":"` + userName + `"}`

	client, err := getPersistenceClient()
	if err != nil {
		return err
	}

	path, permission := getFilePermissionForMethod(method, req)
	log.Println(path, permission)
	// Find file permissions.

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

	roles := account["roles"].([]interface{})
	for i := 0; i < len(roles); i++ {
		//role := roles[i].(map[string]interface{})
	}
	return nil
}

/**
 * Validate user access by role
 */
func validateUserAccess(userName string, method string) error {

	// if the user is the super admin no validation is required.
	if userName == "sa" {
		return nil
	}

	// if guest can run the action...
	if canRunAction("guest", method) == nil {
		// everybody can run the action in that case.
		return nil
	}

	// Now I will get the user roles and validate if the user can execute the
	// method.
	Id := "local_ressource"
	Database := "local_ressource"
	Collection := "Accounts"
	Query := `{"name":"` + userName + `"}`

	client, err := getPersistenceClient()
	if err != nil {
		return err
	}

	values, err := client.FindOne(Id, Database, Collection, Query, `[{"Projection":{"roles":1}}]`)
	if err != nil {
		return err
	}

	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &account)
	if err != nil {
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
	log.Println(err)
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

	log.Println(application, userId, method, err_)

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

// unaryInterceptor calls authenticateClient with current context
func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	clientID, _, err := authenticateClient(ctx)
	if err != nil {
		return nil, err
	}

	// Validate the user access.
	err = validateUserAccess(clientID, info.FullMethod)
	if err != nil {
		return nil, err
	}

	// Validate file access
	err = validateFileAccess(clientID, info.FullMethod, req)
	if err != nil {
		return nil, err
	}

	// Execute the action.
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	result, err := handler(ctx, req)

	if err != nil {
		logError(ctx, info.FullMethod, err)
		return nil, err
	}

	// Log the action as needed for info.
	logAction(ctx, info.FullMethod, result)

	return result, nil
}
