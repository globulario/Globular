package Interceptors

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Authentication holds the login/password
type Authentication struct {
	Token string
}

// GetRequestMetadata gets the current request metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": a.Token,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Generate a token for a ginven user.
func GenerateToken(jwtKey []byte, timeout time.Duration, userName string) (string, error) {

	// Declare the expiration time of the token
	expirationTime := time.Now().Add(timeout * time.Millisecond)
	log.Println("token expire in", timeout, "millisecond")

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	if string(jwtKey) == userName {
		// This is globular server itself.
		log.Println("-----> generate token for: GLOBULAR", tokenString)
	} else {
		log.Println("-----> generate token for: ", userName, tokenString)
	}

	return tokenString, nil
}
