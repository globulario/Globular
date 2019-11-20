package ressource

import (
	"strconv"

	"io/ioutil"
	"os"

	"github.com/davecourtois/Globular/api"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// admin Client Service
////////////////////////////////////////////////////////////////////////////////

type Ressource_Client struct {
	cc *grpc.ClientConn
	c  RessourceServiceClient

	// The name of the service
	name string

	// The client domain
	domain string

	// The port number
	port int

	// is the connection is secure?
	hasTLS bool

	// Link to client key file
	keyFile string

	// Link to client certificate file.
	certFile string

	// certificate authority file
	caFile string
}

// Create a connection to the service.
func NewRessource_Client(domain string, port int, hasTLS bool, keyFile string, certFile string, caFile string) *Ressource_Client {

	client := new(Ressource_Client)
	client.port = port
	client.domain = domain
	client.name = "ressource"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client)
	client.c = NewRessourceServiceClient(client.cc)

	return client
}

// Return the ipv4 address
// Return the address
func (self *Ressource_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the domain
func (self *Ressource_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *Ressource_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Ressource_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Ressource_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Ressource_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Ressource_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Ressource_Client) GetCaFile() string {
	return self.caFile
}

// Register a new Account.
func (self *Ressource_Client) RegisterAccount(name string, email string, password string, confirmation_password string) error {
	rqst := &RegisterAccountRqst{
		Account: &Account{
			Name:     name,
			Email:    email,
			Password: "",
		},
		Password:        password,
		ConfirmPassword: confirmation_password,
	}

	_, err := self.c.RegisterAccount(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

// Delete an account.
func (self *Ressource_Client) DeleteAccount(name string) error {
	rqst := &DeleteAccountRqst{
		Name: name,
	}

	_, err := self.c.DeleteAccount(api.GetClientContext(self), rqst)
	return err
}

// Authenticate a user.
func (self *Ressource_Client) Authenticate(name string, password string) (string, error) {
	// remove the file if it already exist.
	os.Remove(os.TempDir() + string(os.PathSeparator) + self.GetDomain() + "_token")

	rqst := &AuthenticateRqst{
		Name:     name,
		Password: password,
	}

	rsp, err := self.c.Authenticate(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	// Here I will save the token into the temporary directory the token will be valid for a given time (default is 15 minutes)
	// it's the responsability of the client to keep it refresh... see Refresh token from the server...
	err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+self.GetDomain()+"_token", []byte(rsp.Token), 0644)
	if err != nil {
		return "", err
	}

	return rsp.Token, nil
}
