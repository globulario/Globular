package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// The client service interface.
type Client interface {
	// Return the ipv4 address
	GetAddress() string

	// Get Domain return the client domain.
	GetDomain() string

	// Return the name of the service
	GetName() string

	// Close the client.
	Close()

	// At firt the port contain the http(s) address of the globular server.
	// The configuration will be get from that address and the port will
	// be set back to the correct address.
	SetPort(int)

	// Set the name of the client
	SetName(string)

	// Set the domain of the client
	SetDomain(string)

	////////////////// TLS ///////////////////

	//if the client is secure.
	HasTLS() bool

	// Get the TLS certificate file path
	GetCertFile() string

	// Get the TLS key file path
	GetKeyFile() string

	// Get the TLS key file path
	GetCaFile() string

	// Set the client is a secure client.
	SetTLS(bool)

	// Set TLS certificate file path
	SetCertFile(string)

	// Set TLS key file path
	SetKeyFile(string)

	// Set TLS authority trust certificate file path
	SetCaFile(string)
}

/**
 * A simple function to get the client configuration from http.
 */
func getClientConfig(address string, name string) (map[string]interface{}, error) {

	address = strings.Split(address, ":")[0]
	var err error

	// Here I will get the configuration information from http...
	var resp *http.Response

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err = client.Get("http://localhost:10000/client_config?address=" + address + "&name=" + name)

	if err != nil {

		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	err = json.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return config, nil
}

/**
 * Initialyse the client security and set it port to
 */
func InitClient(client Client, address string, name string) {
	// Set the domain and the name from the incomming...
	client.SetDomain(address)
	client.SetName(name)

	// Here I will initialyse the client
	config, err := getClientConfig(address, name)
	if err != nil {
		log.Print(err)
		return
	}

	// Set client attributes.
	client.SetPort(int(config["Port"].(float64)))

	// Set security values.
	client.SetKeyFile(config["KeyFile"].(string))
	client.SetCertFile(config["CertFile"].(string))
	client.SetCaFile(config["CertAuthorityTrust"].(string))
	client.SetTLS(config["TLS"].(bool))
}

/**
 * Get the client connection. The token is require to control access to ressource
 */
func GetClientConnection(client Client) *grpc.ClientConn {
	// initialyse the client values.
	var cc *grpc.ClientConn
	var err error
	if cc == nil {
		address := client.GetAddress()
		if client.HasTLS() {

			log.Println("Secure client ", client.GetDomain(), client.GetName())
			// Setup the login/pass simple test...

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
				ServerName:   client.GetDomain(), // NOTE: this is required!
				Certificates: []tls.Certificate{certificate},
				RootCAs:      certPool,
			})

			// Create a connection with the TLS credentials
			cc, err = grpc.Dial(address, grpc.WithTransportCredentials(creds))

			if err != nil {
				log.Fatalf("could not dial %s: %s", client.GetAddress(), err)
			}
		} else {
			cc, err = grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
		}
	}

	return cc
}

/**
 * That function is use to get the client context. If a token is found in the
 * tmp directory for the client domain it's set in the metadata.
 */
func GetClientContext(client Client) context.Context {
	// Token's are kept in temporary directorys
	domain := client.GetDomain()
	path := os.TempDir() + string(os.PathSeparator) + domain + "_token"
	token, err := ioutil.ReadFile(path)
	if err == nil {
		md := metadata.New(map[string]string{"token": string(token)})
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		return ctx
	}
	return context.Background()
}
