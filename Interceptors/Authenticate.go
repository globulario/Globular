package Interceptors

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
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
	Email    string `json:"email"`

	jwt.StandardClaims
}

// Generate a token for a ginven user.
func GenerateToken(jwtKey []byte, timeout time.Duration, userName string, email string) (string, error) {

	// Declare the expiration time of the token
	expirationTime := time.Now().Add(timeout * time.Millisecond)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: userName,
		Email:    email,
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

	return tokenString, nil
}

/** Validate a Token **/
func ValidateToken(token string) (string, string, int64, error) {

	// Initialize a new instance of `Claims`
	claims := &Claims{}

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
		return claims.Username, claims.Email, claims.ExpiresAt, err
	}

	if !tkn.Valid {
		return claims.Username, claims.Email, claims.ExpiresAt, fmt.Errorf("invalid token!")
	}

	return claims.Username, claims.Email, claims.ExpiresAt, nil
}
