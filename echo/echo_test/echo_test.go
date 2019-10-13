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
	crt = "E:/Project/src/github.com/davecourtois/Globular/creds/client.crt"
	key = "E:/Project/src/github.com/davecourtois/Globular/creds/client.pem"
	ca  = "E:/Project//src/github.com/davecourtois/Globular/creds/ca.crt"

	// The token is print in the sever console. It can be taken from the local temp file if the client service run on the same machine.
	// It will be regenerate each time the server is started and valid until the token expire.
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Ijg4ZmE5ZDIyLWJkNmMtNDMxYi04MDU1LTU0MjRkNmJhMDUxNCIsImV4cCI6NTE3MDA0NDg4OX0.GKMUPWsPTYSde4pe1nNRbJPu2HeZQIidxxP-SxRFbEs"
)

// Test various function here.
func TestEcho(t *testing.T) {

	// Connect to the plc client.
	client := echo_client.NewEcho_Client("localhost", "127.0.0.1:10029", true, key, crt, ca, token)

	val, err := client.Echo("Ceci est un test")
	if err != nil {
		log.Println("---> ", err)
	} else {
		log.Println("---> ", val)
	}
}
