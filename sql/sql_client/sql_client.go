// +build wasm js
package main

import (
	"context"
	"log"
	"syscall/js"

	// grpc stuff.
	"github.com/davecourtois/Globular/sql/sqlpb"
	// "github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

var (
	c  chan bool
	cc *grpc.ClientConn
)

// Connect to the grpc server.
func getClientConnection(address string) *grpc.ClientConn {
	var err error
	var cc *grpc.ClientConn
	if cc == nil {
		cc, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}

	}
	log.Println("Connection whit ", address, "succed!")
	return cc
}

func init() {

	// Here the address must exist in the server configuration file.
	cc = getClientConnection("localhost:50051")
	c = make(chan bool)
}

// Print message test
func printMessage(this js.Value, inputs []js.Value) interface{} {
	callback := inputs[len(inputs)-1:][0]
	message := inputs[0].String()

	callback.Invoke(js.Null(), "Did you say "+message)

	return nil
}

func main() {
	// Register the function.
	js.Global().Set("printMessage", js.FuncOf(printMessage))

	// simple test...

	// Create a new client service...
	cs := sqlpb.NewSqlServiceClient(cc)

	// Here I will try a success case...
	rqst := &sqlpb.PingConnectionRqst{
		Id: "employees_db",
	}

	rsp, err := cs.Ping(context.Background(), rqst)
	if err != nil {
		log.Println("68 error while CreateConnection: %v", err)
	} else {
		log.Println("---> 71 ", rsp.Result)
	}

	<-c
}
