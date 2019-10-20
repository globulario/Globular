package ressource

import (
	//"encoding/json"
	"log"
	"testing"
)

var (
	// Connect to the plc client.
	crt    = "E:/Project/src/github.com/davecourtois/Globular/creds/client.crt"
	key    = "E:/Project/src/github.com/davecourtois/Globular/creds/client.pem"
	ca     = "E:/Project//src/github.com/davecourtois/Globular/creds/ca.crt"
	client = NewRessource_Client("localhost", "127.0.0.1:10003", true, key, crt, ca)
	//client = NewRessource_Client("globular_server", "135.19.213.242:10003", false, key, crt, ca)
)

// Remove an account.
func TestDeleteAccount(t *testing.T) {

	log.Println("---> test remove existing account.")
	err := client.DeleteAccount("davecourtois")
	if err != nil {

		log.Println("---> ", err)
	}
}

// Test various function here.
func TestRegisterAccount(t *testing.T) {

	log.Println("---> test register a new account.")
	err := client.RegisterAccount("davecourtois", "dave.courtois60@gmail.com", "1234", "1234")
	if err != nil {
		log.Println("---> ", err)
	}
}

func TestAuthenticate(t *testing.T) {

	log.Println("---> test authenticate account.")
	//token, err := client.Authenticate("dave.courtois60@gmail.com", "1234")
	token, err := client.Authenticate("davecourtois", "1234")
	if err != nil {
		log.Println("---> ", err)
	} else {
		log.Println("---> ", token)
	}
}