package Globular

import (
	//"encoding/json"
	"log"
	"testing"

	"github.com/davecourtois/Globular/echo/echo_client"
	"github.com/davecourtois/Globular/ressource"
)

var (

	// Test with a secure connection.
	crt = "E:/Project/src/github.com/davecourtois/Globular/creds/client.crt"
	key = "E:/Project/src/github.com/davecourtois/Globular/creds/client.pem"
	ca  = "E:/Project//src/github.com/davecourtois/Globular/creds/ca.crt"

	// Create a new connection to globular ressource manager.
	globular = ressource.NewRessource_Client("MOND599", "10.67.44.131:10029", true, key, crt, ca)
)

// Test various function here.
func TestEcho(t *testing.T) {
	log.Println("---> test Echo")
	token, err := globular.Authenticate("admin", "adminadmin")
	if err != nil {
		log.Println("---> error ", err)
		return
	}

	// Connect to the plc client.
	client := echo_client.NewEcho_Client("MOND599", "10.67.44.131:10029", true, key, crt, ca, token)

	val, err := client.Echo("Ceci est un test")
	if err != nil {
		log.Println("---> ", err)
	} else {
		log.Println("---> ", val)
	}
}
