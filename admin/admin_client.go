package admin

import (
	"context"
	// "log"

	"github.com/davecourtois/Globular/api"
	//	"github.com/davecourtois/Utility"
	"strconv"

	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// admin Client Service
////////////////////////////////////////////////////////////////////////////////

type Admin_Client struct {
	cc *grpc.ClientConn
	c  AdminServiceClient

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
func NewAdmin_Client(domain string, port int, hasTLS bool, keyFile string, certFile string, caFile string, token string) *Admin_Client {
	client := new(Admin_Client)

	client.port = port
	client.domain = domain
	client.name = "admin"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client, token)
	client.c = NewAdminServiceClient(client.cc)

	return client
}

// Return the address
func (self *Admin_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the domain
func (self *Admin_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *Admin_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Admin_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Admin_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Admin_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Admin_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Admin_Client) GetCaFile() string {
	return self.caFile
}

// Get server configuration.
func (self *Admin_Client) GetConfig() (string, error) {
	rqst := new(GetConfigRequest)

	rsp, err := self.c.GetConfig(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

// Get the server configuration with all detail must be secured.
func (self *Admin_Client) GetFullConfig() (string, error) {
	rqst := new(GetConfigRequest)

	rsp, err := self.c.GetFullConfig(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

func (self *Admin_Client) SaveConfig(config string) error {
	rqst := &SaveConfigRequest{
		Config: config,
	}

	_, err := self.c.SaveConfig(context.Background(), rqst)
	if err != nil {
		return err
	}

	return nil
}

func (self *Admin_Client) StartService(id string) (int, int, error) {
	rqst := new(StartServiceRequest)
	rqst.ServiceId = id
	rsp, err := self.c.StartService(context.Background(), rqst)
	if err != nil {
		return -1, -1, err
	}

	return int(rsp.ServicePid), int(rsp.ProxyPid), nil
}

func (self *Admin_Client) StopService(id string) error {
	rqst := new(StopServiceRequest)
	rqst.ServiceId = id
	_, err := self.c.StopService(context.Background(), rqst)
	if err != nil {
		return err
	}

	return nil
}

// Register and start an application.
func (self *Admin_Client) RegisterExternalApplication(id string, path string, args []string) (int, error) {
	rqst := &RegisterExternalApplicationRequest{
		ServiceId: id,
		Path:      path,
		Args:      args,
	}

	rsp, err := self.c.RegisterExternalApplication(context.Background(), rqst)

	if err != nil {
		return -1, err
	}

	return int(rsp.ServicePid), nil
}

/////////////////////////// Services management functions ////////////////////////

/**
 * Publish a service from a runing globular server.
 */
func (self *Admin_Client) PublishService(serviceId string, discoveryAddress string, repositoryAddress string, description string, keywords []string) error {

	rqst := new(PublishServiceRequest)
	rqst.Description = description
	rqst.DicorveryId = discoveryAddress
	rqst.RepositoryId = repositoryAddress
	rqst.Keywords = keywords
	rqst.ServiceId = serviceId

	_, err := self.c.PublishService(context.Background(), rqst)

	return err
}

/**
 * Intall a new service or update an existing one.
 */
func (self *Admin_Client) InstallService(discoveryId string, publisherId string, serviceId string) error {

	rqst := new(InstallServiceRequest)
	rqst.DicorveryId = discoveryId
	rqst.PublisherId = publisherId
	rqst.ServiceId = serviceId

	_, err := self.c.InstallService(context.Background(), rqst)

	return err
}
