package echo_client

import (
	//"encoding/json"
	"log"
	"testing"

	"github.com/globulario/Globular/echo/echo_client"
)

// Test various function here.
func TestEcho(t *testing.T) {

	// Connect to the plc client.
	client, err := echo_client.NewEcho_Client("localhost" /*"echo.EchoService"*/, "002cd8d6-f933-4442-aac8-ccde5e6a50cf")
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
