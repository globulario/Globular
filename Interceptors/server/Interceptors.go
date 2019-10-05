package Interceptors

import (
	"fmt"
	"log"

	"strings"

	"io/ioutil"
	"os"

	"github.com/davecourtois/Globular/Interceptors/Authenticate"
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

func validateToken(token string) error {

	// Initialize a new instance of `Claims`
	claims := &Interceptors.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	log.Println("---> receive token: ", token)
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Get the key from the local temp file.
		jwtKey, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "globular_key")
		return jwtKey, err
	})

	if err != nil {
		log.Println("Invalid token error line 40", err)
		return err
	}

	if !tkn.Valid {
		return fmt.Errorf("invalid token!")
	}

	return nil
}

// authenticateAgent check the client credentials
func authenticateClient(ctx context.Context) (string, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		token := strings.Join(md["token"], "")
		// In that case no token was given...
		if len(token) == 0 {
			return "", nil
		}

		err := validateToken(token)
		if err != nil {
			return "", err
		}

		if len(token) > 0 {
			return "42", nil
		}
	}

	return "", fmt.Errorf("missing credentials")
}

// unaryInterceptor calls authenticateClient with current context
func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	clientID, err := authenticateClient(ctx)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, clientIDKey, clientID)

	return handler(ctx, req)
}
