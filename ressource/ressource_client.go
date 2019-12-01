package ressource

import (
	"strconv"

	"io/ioutil"
	"os"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Utility"
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
func NewRessource_Client(address string, name string) *Ressource_Client {

	client := new(Ressource_Client)
	api.InitClient(client, address, name)
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

// Set grpc_service port.
func (self *Ressource_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *Ressource_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Ressource_Client) SetDomain(domain string) {
	self.domain = domain
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

// Set the client is a secure client.
func (self *Ressource_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Ressource_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Ressource_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Ressource_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

////////////// API ////////////////
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
	// In case of other domain than localhost I will rip off the token file
	// before each authentication.
	path := os.TempDir() + string(os.PathSeparator) + self.GetDomain() + "_token"
	if !Utility.IsLocal(self.GetDomain()) {
		// remove the file if it already exist.
		os.Remove(path)
	}

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
	if !Utility.IsLocal(self.GetDomain()) {
		err = ioutil.WriteFile(path, []byte(rsp.Token), 0644)
		if err != nil {
			return "", err
		}
	}

	return rsp.Token, nil
}
