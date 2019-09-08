package Globular

import (
	"context"
	"fmt"
	"log"

	"crypto/tls"
	"crypto/x509"

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

			// Load the client certificates from disk

			crt := "/media/dave/60B6E593B6E569CC/Project/src/github.com/davecourtois/Globular/creds/client.crt"
			key := "/media/dave/60B6E593B6E569CC/Project/src/github.com/davecourtois/Globular/creds/client.pem"
			certificate, err := tls.LoadX509KeyPair(crt, key)
			if err != nil {
				log.Fatalf("could not load client key pair: %s", err)
			}

			// Create a certificate pool from the certificate authority
			certPool := x509.NewCertPool()
			ca, err := ioutil.ReadFile(config["CertAuthorityTrust"].(string))
			if err != nil {
				log.Fatalf("could not read ca certificate: %s", err)
			}

			// Append the certificates from the CA
			if ok := certPool.AppendCertsFromPEM(ca); !ok {
				log.Fatalf("failed to append ca certs")
			}

			creds := credentials.NewTLS(&tls.Config{
				ServerName:   "localhost", // NOTE: this is required!
				Certificates: []tls.Certificate{certificate},
				RootCAs:      certPool,
			})

			// Create a connection with the TLS credentials
			cc, err = grpc.Dial(addresse, grpc.WithTransportCredentials(creds))
			if err != nil {
				log.Fatalf("could not dial %s: %s", addresse, err)
			}
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
