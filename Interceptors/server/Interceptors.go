package Interceptors

import (
	"fmt"
	"log"

	"strings"

	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/davecourtois/Globular/Interceptors/Authenticate"
	"github.com/davecourtois/Globular/persistence/persistence_client"
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
	token_ string // the last know sa token.
)

func getPersistenceClient() (*persistence_client.Persistence_Client, error) {
	// Here I will need the persistence client to read user permission.
	// Here I will read the server token, the service must run on the
	// same computer as globular.
	token, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "globular_token")
	if err != nil {
		return nil, err
	}

	if client == nil || token_ != string(token) {
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

		root := infos["pwd"].(string)
		addresse := infos["address"].(string)
		crt := infos["certFile"].(string)
		key := infos["keyFile"].(string)
		ca := infos["certAuthorityTrust"].(string)

		// close the
		if client != nil {
			client.Close()
		}

		client = persistence_client.NewPersistence_Client("localhost", addresse, true, key, crt, ca, string(token))
		client.Connect("local_ressource", root)

		// keep the token for futher use
		token_ = string(token)
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
			log.Println("no token was given.")
			return "", 0, nil
		}

		log.Println("token from incoming: ", token)

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

	return errors.New("permission denied! account " + userName + " cannot execute methode '" + method + "'")
}

// unaryInterceptor calls authenticateClient with current context
func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	clientID, _, err := authenticateClient(ctx)
	if err != nil {
		return nil, err
	}

	err = validateUserAccess(clientID, info.FullMethod)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, clientIDKey, clientID)

	return handler(ctx, req)
}
