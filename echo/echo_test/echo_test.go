package Globular

import (
	"context"
	"fmt"
	"log"

	"github.com/davecourtois/Globular/echo/echopb"
	"google.golang.org/grpc"

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
	var err error
	var cc *grpc.ClientConn
	if cc == nil {
		cc, err = grpc.Dial(addresse, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
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
