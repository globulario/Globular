package Globular

import (
	//"encoding/json"
	"log"
	"testing"

	"github.com/davecourtois/Globular/echo/echo_client"
)

var (

	// Test with a secure connection.
	// Those files must be accessible to the client to be able to call
	// gRpc function on the server.
	crt = "E:/Project/src/github.com/davecourtois/Globular/config/grpc_tls/client.crt"
	key = "E:/Project/src/github.com/davecourtois/Globular/config/grpc_tls/client.pem"
	ca  = "E:/Project/src/github.com/davecourtois/Globular/config/grpc_tls/ca.crt"
)

// Test various function here.
func TestEcho(t *testing.T) {

	// Connect to the plc client.
	client := echo_client.NewEcho_Client("localhost", 10029, true, key, crt, ca)

	val, err := client.Echo("Ceci est un test")
	if err != nil {
		log.Println("---> ", err)
	} else {
		log.Println("---> ", val)
	}
}
