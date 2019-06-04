package api

import (
	"io/ioutil"
	"log"

	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// The client service interface.
type Client interface {
	// Return the ipv4 address
	GetAddress() string

	// Return the name of the service
	GetName() string

	// Close the client.
	Close()

	////////////////// TLS ///////////////////

	//if the client is secure.
	HasTLS() bool

	// Get the TLS certificate file path
	GetCertFile() string

	// Get the TLS key file path
	GetKeyFile() string

	// Get the TLS key file path
	GetCaFile() string
}

/**
 * Get the client connection.
 */
func GetClientConnection(client Client) *grpc.ClientConn {

	var cc *grpc.ClientConn
	var err error
	if cc == nil {
		if client.HasTLS() {
			if len(client.GetKeyFile()) == 0 {
				log.Println("no key file is available for client ")
			}

			if len(client.GetCertFile()) == 0 {
				log.Println("no certificate file is available for client ")
			}

			certificate, err := tls.LoadX509KeyPair(client.GetCertFile(), client.GetKeyFile())
			if err != nil {
				log.Fatalf("could not load client key pair: %s", err)
			}

			// Create a certificate pool from the certificate authority
			certPool := x509.NewCertPool()
			ca, err := ioutil.ReadFile(client.GetCaFile())
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
			cc, err = grpc.Dial(client.GetAddress(), grpc.WithTransportCredentials(creds))
			if err != nil {
				log.Fatalf("could not dial %s: %s", client.GetAddress(), err)
			}
		} else {
			cc, err = grpc.Dial(client.GetAddress(), grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
		}

	}
	return cc
}
