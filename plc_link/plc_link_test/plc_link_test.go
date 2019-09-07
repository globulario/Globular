package Globular

import (
	"log"
	"testing"

	"github.com/davecourtois/Globular/plc_link/plc_link_client"
)

var (
	// Connect to the plc client.

	// client = NewPlcLink_Client.("stevePc", "135.19.213.242:10027", false, "", "", "")
	client = plc_link_client.NewPlcLink_Client("stevePc", "127.0.0.1:10027", false, "", "", "")
)

// Test various function here.
func TestCreateConnection(t *testing.T) {

	// Create a connection.
	err := client.CreateConnection("local_server_1", "localhost", "127.0.0.1:10023")
	if err != nil {
		log.Println("---> error: ", err)
	}
}
