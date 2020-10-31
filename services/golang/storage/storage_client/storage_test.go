package storage_client

import (
	"fmt"
	"log"

	"testing"

	"github.com/globulario/Globular/storage/storage_client"
)

// Set the correct addresse here as needed.
var (
	// Test with a secure connection.
	// Those files must be accessible to the client to be able to call
	// gRpc function on the server.
	crt = "E:/Project/src/github.com/globulario/Globular/config/tls/client.crt"
	key = "E:/Project/src/github.com/globulario/Globular/config/tls/client.pem"
	ca  = "E:/Project/src/github.com/globulario/Globular/config/tls/ca.crt"

	// The token is print in the sever console. It can be taken from the local temp file if the client service run on the same machine.
	// It will be regenerate each time the server is started and valid until the token expire.
	token  = ""
	client = storage_client.NewStorage_Client("localhost", "127.0.0.1:10013", false, key, crt, ca, token)
)

// First test create a fresh new connection...
func TestCreateConnection(t *testing.T) {
	fmt.Println("Connection creation test.")
	err := client.CreateConnection("test_storage", "storage_test", 0.0)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}
}

func TestOpenConnection(t *testing.T) {
	err := client.OpenConnection("test_storage", `{"path":"C:/temp", "name":"storage_test"}`)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}
	log.Println("Open connection success!")
}

// Test set item.
func TestSetItem(t *testing.T) {
	err := client.SetItem("test_storage", "1", []byte(`{"prop_1":"This is a test!", "prop_2":1212}`))
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}
	log.Println("Set item success!")
}

func TestGetItem(t *testing.T) {
	values, err := client.GetItem("test_storage", "1")
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}
	log.Println("Get item success with value", string(values))
}

func TestRemoveItem(t *testing.T) {
	err := client.RemoveItem("test_storage", "1")
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}
	log.Println("Remove item success!")
}

func TestClear(t *testing.T) {
	err := client.Clear("test_storage")
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}
	log.Println("Clear all items success!")
}

func TestDrop(t *testing.T) {

	err := client.Drop("test_storage")
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}
	log.Println("Drop store success!")
}

func TestCloseConnection(t *testing.T) {

	err := client.CloseConnection("test_storage")
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}
	log.Println("close connection success!")
}

// Test a ldap query.
func TestDeleteConnection(t *testing.T) {
	err := client.DeleteConnection("test_storage")
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}
	log.Println("Delete connection success!")
}
