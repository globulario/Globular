package ca

import (
	"strconv"

	"github.com/davecourtois/Globular/api"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// Ca Client Service
////////////////////////////////////////////////////////////////////////////////

type Ca_Client struct {
	cc *grpc.ClientConn
	c  CertificateAuthorityClient

	// The name of the service
	name string

	// The port
	port int

	// The client domain
	domain string

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
func NewCa_Client(config map[string]interface{}) *Ca_Client {
	client := new(Ca_Client)

	api.InitClient(client, config)
	client.cc = api.GetClientConnection(client)
	client.c = NewCertificateAuthorityClient(client.cc)

	return client
}

// Return the address
func (self *Ca_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the domain
func (self *Ca_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *Ca_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Ca_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *Ca_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *Ca_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Ca_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Ca_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Ca_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Ca_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Ca_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Ca_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Ca_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Ca_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Ca_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

////////////////////////////////////////////////////////////////////////////////
// CA functions.
////////////////////////////////////////////////////////////////////////////////

/**
 * Take signing request and made it sign by the server. If succed a signed
 * certificate is return.
 */
func (self *Ca_Client) SignCertificate(csr string) (string, error) {
	// The certificate request.
	rqst := new(SignCertificateRequest)
	rqst.Csr = csr

	rsp, err := self.c.SignCertificate(api.GetClientContext(self), rqst)
	if err == nil {
		return rsp.Crt, nil
	}

	return "", err
}

/**
 * Get the ca.crt file content.
 */
func (self *Ca_Client) GetCaCertificate() (string, error) {
	rqst := new(GetCaCertificateRequest)

	rsp, err := self.c.GetCaCertificate(api.GetClientContext(self), rqst)
	if err == nil {
		return rsp.Ca, nil
	}
	return "", err
}
