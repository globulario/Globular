package admin

import (
	"bytes"
	"io"
	"os"
	"strconv"

	"github.com/davecourtois/Globular/api"

	"github.com/davecourtois/Utility"
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
func NewAdmin_Client(domain string, port int, hasTLS bool, keyFile string, certFile string, caFile string) *Admin_Client {
	client := new(Admin_Client)

	client.port = port
	client.domain = domain
	client.name = "admin"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client)
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

	rsp, err := self.c.GetConfig(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

// Get the server configuration with all detail must be secured.
func (self *Admin_Client) GetFullConfig() (string, error) {
	rqst := new(GetConfigRequest)

	rsp, err := self.c.GetFullConfig(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

func (self *Admin_Client) SaveConfig(config string) error {
	rqst := &SaveConfigRequest{
		Config: config,
	}

	_, err := self.c.SaveConfig(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

func (self *Admin_Client) StartService(id string) (int, int, error) {
	rqst := new(StartServiceRequest)
	rqst.ServiceId = id
	rsp, err := self.c.StartService(api.GetClientContext(self), rqst)
	if err != nil {
		return -1, -1, err
	}

	return int(rsp.ServicePid), int(rsp.ProxyPid), nil
}

func (self *Admin_Client) StopService(id string) error {
	rqst := new(StopServiceRequest)
	rqst.ServiceId = id
	_, err := self.c.StopService(api.GetClientContext(self), rqst)
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

	rsp, err := self.c.RegisterExternalApplication(api.GetClientContext(self), rqst)

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

	_, err := self.c.PublishService(api.GetClientContext(self), rqst)

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

	_, err := self.c.InstallService(api.GetClientContext(self), rqst)

	return err
}

/**
 * Intall a new service or update an existing one.
 */
func (self *Admin_Client) UninstallService(publisherId string, serviceId string, version string) error {

	rqst := new(UninstallServiceRequest)
	rqst.PublisherId = publisherId
	rqst.ServiceId = serviceId
	rqst.Version = version

	_, err := self.c.UninstallService(api.GetClientContext(self), rqst)

	return err
}

/**
 * Deploy the content of an application with a given name to the server.
 */
func (self *Admin_Client) DeployApplication(name string, path string) error {

	rqst := new(DeployApplicationRequest)
	rqst.Name = name

	Utility.CreateDirIfNotExist(Utility.GenerateUUID(name))
	Utility.CopyDirContent(path, Utility.GenerateUUID(name))

	// Now I will open the data and create a archive from it.
	var buffer bytes.Buffer
	err := Utility.CompressDir(Utility.GenerateUUID(name), &buffer)
	if err != nil {
		return err
	}

	os.RemoveAll(Utility.GenerateUUID(name))

	// Open the stream...
	stream, err := self.c.DeployApplication(api.GetClientContext(self))
	if err != nil {
		return err
	}

	const BufferSize = 1024 * 5 // the chunck size.
	for {
		var data [BufferSize]byte
		bytesread, err := buffer.Read(data[0:BufferSize])
		if bytesread > 0 {
			rqst := &DeployApplicationRequest{
				Data: data[0:bytesread],
				Name: name,
			}
			// send the data to the server.
			err = stream.Send(rqst)
		}

		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		os.RemoveAll(Utility.GenerateUUID(name))
		return err
	}

	return nil

}
