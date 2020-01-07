package Globular

import (
	//"encoding/json"
	"log"
	"testing"

	"github.com/davecourtois/Globular/echo/echo_client"
)

// Test various function here.
func TestEcho(t *testing.T) {

	// Connect to the plc client.
	client := echo_client.NewEcho_Client("globular4.omniscient.app", "echo_server")

	val, err := client.Echo("Ceci est un test")
	if err != nil {
		log.Println("---> ", err)
	} else {
		log.Println("---> ", val)
	}
}
