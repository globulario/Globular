package Globular

import (
	"context"
	"fmt"
	"log"

	"encoding/json"
	"io/ioutil"

	"github.com/davecourtois/Globular/storage/storagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"testing"
)

// Set the correct addresse here as needed.
var (
	addresse = "localhost:10013"
)

/**
 * Get the client connection.
 */
func getClientConnection() *grpc.ClientConn {
	// So here I will read the server configuration to see if the connection
	// is secure...
	config := make(map[string]interface{})
	data, err := ioutil.ReadFile("../storage_server/config.json")
	if err != nil {
		log.Fatal("fail to read configuration")
	}

	// Read the config file.
	json.Unmarshal(data, &config)

	var cc *grpc.ClientConn
	if cc == nil {
		if config["TLS"].(bool) {
			creds, sslErr := credentials.NewClientTLSFromFile(config["CertAuthorityTrust"].(string), "")
			if err != nil {
				log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			}
			opts := grpc.WithTransportCredentials(creds)
			cc, err = grpc.Dial(addresse, opts)
		} else {
			cc, err = grpc.Dial(addresse, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
		}

	}
	return cc
}

// First test create a fresh new connection...
func TestCreateConnection(t *testing.T) {
	fmt.Println("Connection creation test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := storagepb.NewStorageServiceClient(cc)

	rqst := &storagepb.CreateConnectionRqst{
		Connection: &storagepb.Connection{
			Id:   "test_storage",
			Name: "storage_test",
			//Type: storagepb.StoreType_BIG_CACHE, // Memory store (volatile)
			Type: storagepb.StoreType_LEVEL_DB, // Disk store (persistent)
		},
	}

	rsp, err := c.CreateConnection(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}

func TestOpenConnection(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := storagepb.NewStorageServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &storagepb.OpenRqst{
		Id:      "test_storage",
		Options: `{"path":"C:\\temp", "name":"storage_test"}`,
	}

	_, err := c.Open(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}

	log.Println("Open connection success!")
}

// Test set item.
func TestSetItem(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := storagepb.NewStorageServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &storagepb.SetItemRequest{
		Id:    "test_storage",
		Key:   "1",
		Value: []byte(`{"prop_1":"This is a test!", "prop_2":1212}`),
	}

	_, err := c.SetItem(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}

	log.Println("Set item success!")
}

func TestGetItem(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := storagepb.NewStorageServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &storagepb.GetItemRequest{
		Id:  "test_storage",
		Key: "1",
	}

	rsp, err := c.GetItem(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}

	log.Println("Get item success with value", string(rsp.GetResult()))
}

func TestRemoveItem(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := storagepb.NewStorageServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &storagepb.RemoveItemRequest{
		Id:  "test_storage",
		Key: "1",
	}

	_, err := c.RemoveItem(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}

	log.Println("Remove item success!")
}

func TestClear(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := storagepb.NewStorageServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &storagepb.ClearRequest{
		Id: "test_storage",
	}

	_, err := c.Clear(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}

	log.Println("Clear all items success!")
}

func TestDrop(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := storagepb.NewStorageServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &storagepb.DropRequest{
		Id: "test_storage",
	}

	_, err := c.Drop(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}

	log.Println("Drop store success!")
}

func TestCloseConnection(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := storagepb.NewStorageServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &storagepb.CloseRqst{
		Id: "test_storage",
	}

	_, err := c.Close(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}

	log.Println("close connection success!")
}

// Test a ldap query.
func TestDeleteConnection(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := storagepb.NewStorageServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &storagepb.DeleteConnectionRqst{
		Id: "test_storage",
	}

	_, err := c.DeleteConnection(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}

	log.Println("Delete connection success!")
}
