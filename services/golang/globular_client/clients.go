package globular_client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"reflect"

	"github.com/globulario/Globular/security"
	"github.com/davecourtois/Utility"
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

	// Return the id of the service.
	GetId() string

	// Return the name of the service.
	GetName() string

	// Close the client.
	Close()

	// At firt the port contain the http(s) address of the globular server.
	// The configuration will be get from that address and the port will
	// be set back to the correct address.
	SetPort(int)

	// Set the id of the client
	SetId(string)

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

	// Invoque a request on the client and return it grpc reply.
	Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error)
}

/**
 * Initialyse the client security and set it port to
 */
func InitClient(client Client, address string, id string) error {

	// Set the domain and the name from the incomming...
	client.SetDomain(address)
	address_ := strings.Split(address, ":")

	port := 80 // the default http port...
	if len(address_) == 2 {
		address = address_[0]
		port = Utility.ToInt(address_[1])
	}

	// Here I will initialyse the client
	config, err := security.GetClientConfig(address, id, port)
	if err != nil {
		return err
	}

	// Set client attributes.
	if config["Id"] != nil {
		client.SetId(config["Id"].(string))
	} else {
		client.SetId(config["Name"].(string))
	}

	client.SetName(config["Name"].(string))
	client.SetPort(int(config["Port"].(float64)))
	client.SetDomain(config["Domain"].(string))

	// Set security values.
	client.SetKeyFile(config["KeyFile"].(string))
	client.SetCertFile(config["CertFile"].(string))
	client.SetCaFile(config["CertAuthorityTrust"].(string))
	client.SetTLS(config["TLS"].(bool))

	return nil
}

/**
 * Get the client connection. The token is require to control access to ressource
 */
func GetClientConnection(client Client) (*grpc.ClientConn, error) {
	// initialyse the client values.
	var cc *grpc.ClientConn
	var err error
	if cc == nil {
		address := client.GetAddress()
		if client.HasTLS() {

			// Setup the login/pass simple test...
			if len(client.GetKeyFile()) == 0 {
				log.Println("no key file is available for client ")
			}

			if len(client.GetCertFile()) == 0 {

				return nil, errors.New("no certificate file is available for client")
			}

			certificate, err := tls.LoadX509KeyPair(client.GetCertFile(), client.GetKeyFile())
			if err != nil {
				return nil, err
			}

			// Create a certificate pool from the certificate authority
			certPool := x509.NewCertPool()

			ca, err := ioutil.ReadFile(client.GetCaFile())
			if err != nil {
				return nil, err
			}

			// Append the certificates from the CA
			if ok := certPool.AppendCertsFromPEM(ca); !ok {
				return nil, errors.New("failed to append ca certs")
			}

			creds := credentials.NewTLS(&tls.Config{
				ServerName:   client.GetDomain(), // NOTE: this is required!
				Certificates: []tls.Certificate{certificate},
				RootCAs:      certPool,
			})

			// Create a connection with the TLS credentials
			cc, err = grpc.Dial(address, grpc.WithTransportCredentials(creds))

			if err != nil {
				return nil, err
			}
		} else {
			cc, err = grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return nil, err
			}
		}
	}

	return cc, nil
}

/**
 * That function is use to get the client context. If a token is found in the
 * tmp directory for the client domain it's set in the metadata.
 */
func GetClientContext(client Client) context.Context {

	// if the address is local.
	path := os.TempDir() + string(os.PathSeparator) + client.GetDomain() + "_token"
	token, err := ioutil.ReadFile(path)
	address := client.GetDomain()
	if Utility.Exists(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT") {
		root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")
		port := Utility.ToInt(string(root)[strings.LastIndex(string(root), ":")+1:])
		address += ":" + Utility.ToString(port)
	}

	if err == nil {
		md := metadata.New(map[string]string{"token": string(token), "domain": address, "mac": Utility.MyMacAddr(), "ip": Utility.MyIP()})
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		return ctx
	}

	md := metadata.New(map[string]string{"token": "", "domain": address, "mac": Utility.MyMacAddr(), "ip": Utility.MyIP()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx

}

/**
 * Invoke a method on a client. The client is
 * ctx is the client request context.
 * method is the rpc method to run.
 * rqst is the request to run.
 */
func InvokeClientRequest(client interface{}, ctx context.Context, method string, rqst interface{}) (interface{}, error) {
	methodName := method[strings.LastIndex(method, "/")+1:]
	var err error
	reply, err_ := Utility.CallMethod(client, methodName, []interface{}{ctx, rqst})
	if err_ != nil {
		if reflect.TypeOf(err_).Kind() == reflect.String {
			err = errors.New(err_.(string))
		} else {
			err = err_.(error)
		}
	}

	return reply, err
}
