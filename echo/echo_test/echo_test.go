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
	client, err := echo_client.NewEcho_Client("globular.io", "echo.EchoService")
	if err != nil {
		log.Println("17 ---> ", err)
		return
	}

	val, err := client.Echo("Ceci est un test")
	if err != nil {

		log.Println("23 ---> ", err)
	} else {
		log.Println("25 ---> ", val)
	}
}
