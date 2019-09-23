package ressource

import (
	//"encoding/json"
	"log"
	"testing"
)

var (
	// Connect to the plc client.
	crt = "E:/Project/src/github.com/davecourtois/Globular/creds/client.crt"
	key = "E:/Project/src/github.com/davecourtois/Globular/creds/client.pem"
	ca  = "E:/Project//src/github.com/davecourtois/Globular/creds/ca.crt"

	client = NewRessource_Client("localhost", "127.0.0.1:10003", true, key, crt, ca)
)

// Test various function here.
func TestRegisterAccount(t *testing.T) {

	log.Println("---> test register a new account.")
	err := client.RegisterAccount("davecourtois", "dave.courtois60@gmail.com", "123", "123")
	if err != nil {

		log.Println("---> ", err)
	}
}

func TestAuthenticate(t *testing.T) {

	log.Println("---> test authenticate account.")
	token, err := client.Authenticate("davecourtois", "123")

	if err != nil {
		log.Println("---> ", err)
	} else {
		log.Println("---> ", token)
	}
}

// Remove an account.
func TestDeleteAccount(t *testing.T) {

	log.Println("---> test register a new account.")
	err := client.DeleteAccount("davecourtois")
	if err != nil {

		log.Println("---> ", err)
	}
}
