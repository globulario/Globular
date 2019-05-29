package Globular

import (
	"context"
	"fmt"
	"log"

	"github.com/davecourtois/Globular/echo/echopb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"encoding/json"
	"io/ioutil"
	"testing"
)

// Set the correct addresse here as needed.
var (
	addresse = "localhost:10001"
)

/**
 * Get the client connection.
 */
func getClientConnection() *grpc.ClientConn {
	// So here I will read the server configuration to see if the connection
	// is secure...
	config := make(map[string]interface{})
	data, err := ioutil.ReadFile("../echo_server/config.json")
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
func TestEcho(t *testing.T) {
	fmt.Println("Connection creation test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := echopb.NewEchoServiceClient(cc)

	rqst := &echopb.EchoRequest{
		Message: "Hello Globular",
	}

	rsp, err := c.Echo(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Message)
}
