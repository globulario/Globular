package spc_client

import (
	"context"
	"fmt"
	"log"

	"github.com/globulario/Globular/spc/spcpb"
	"google.golang.org/grpc"

	"testing"
)

// Set the correct addresse here as needed.
var (
	addresse = "localhost:50051"
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
func TestSpc(t *testing.T) {

	data := `[["CDJ5796C",0.017,"2019-06-26 09:57:51.0320000",true],["CDJ5753C",0.062,"2019-06-26 07:44:08.2640000",true],["CDJ5749C",0.055,"2019-06-26 05:44:10.6910000",true],["CDJ5755C",0.051,"2019-06-25 15:46:31.6100000",true],["CDJ5794C",0.022,"2019-06-25 09:41:38.2100000",true],["CDJ5751C",0.098,"2019-06-23 21:44:54.9700000",true],["CDJ1245C",0.045,"2019-06-23 10:15:48.2550000",true],["CDJ1212C",0.061,"2019-06-23 09:20:34.7800000",true],["CDJ5792C",0.014,"2019-06-22 17:25:22.0210000",true],["CDJ5741C",0.057,"2019-06-22 11:10:15.8360000",true],["CDJ5786C",0.026,"2019-06-22 10:23:21.1310000",true],["CDJ5788C",0.051,"2019-06-22 09:12:28.1980000",true],["CDJ1247C",0.032,"2019-06-21 21:38:29.3520000",true],["CDJ1204C",0.048,"2019-06-21 12:19:22.1110000",true],["CDJ5747C",0.033,"2019-06-21 08:22:05.0790000",true],["CDJ5784C",0.04,"2019-06-20 09:05:25.3130000",true],["CDJ1210C",0.058,"2019-06-20 08:59:52.7690000",true],["CDJ1237C",0.034,"2019-06-20 08:33:02.4060000",true],["CDJ5782C",0.052,"2019-06-18 13:25:05.2320000",true],["CDJ5778C",0.062,"2019-06-18 11:26:25.8300000",true],["CDJ5743C",0.065,"2019-06-18 11:03:02.0290000",true],["CDJ5780C",0.036,"2019-06-18 07:46:28.6510000",true],["CDJ5776C",0.014,"2019-06-16 16:12:46.8560000",true],["CDJ5745C",0.025,"2019-06-16 15:19:39.1420000",true],["CDJ5739C",0.074,"2019-06-15 12:38:46.1920000",true],["CDJ5737C",0.034,"2019-06-14 12:55:34.0130000",true],["CDJ1229C",0.031,"2019-06-14 12:24:54.0470000",true],["CDJ5772C",0.025,"2019-06-14 07:39:28.9290000",true],["CDJ5711C",0.041,"2019-06-13 14:11:07.2420000",true],["CDJ5723C",0.052,"2019-06-13 13:46:38.2780000",true]]`
	tests := `[{"state":true,"value":0},{"state":true,"value":9},{"state":true,"value":6},{"state":true,"value":14},{"state":true,"value":2},{"state":true,"value":4},{"state":true,"value":15},{"state":true,"value":8}]`
	tolzon := 0.0
	lotol := 0.0
	uptol := 0.5
	toltype := 0.0
	ispopulation := false

	fmt.Println("Connection creation test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := spcpb.NewSpcServiceClient(cc)

	rqst := &spcpb.CreateAnalyseRqst{
		Data:         data,
		Tests:        tests,
		Tolzon:       tolzon,
		Lotol:        lotol,
		Uptol:        uptol,
		Toltype:      toltype,
		Ispopulation: ispopulation,
	}

	rsp, err := c.CreateAnalyse(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}
