package services

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/davecourtois/Globular/api"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// Service Discovery Client
////////////////////////////////////////////////////////////////////////////////
type ServicesDiscovery_Client struct {
	cc *grpc.ClientConn
	c  ServiceDiscoveryClient

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
func NewServicesDiscovery_Client(address string, name string) *ServicesDiscovery_Client {
	client := new(ServicesDiscovery_Client)
	api.InitClient(client, address, name)
	client.cc = api.GetClientConnection(client)
	client.c = NewServiceDiscoveryClient(client.cc)

	return client
}

// Return the ipv4 address
func (self *ServicesDiscovery_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the domain
func (self *ServicesDiscovery_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *ServicesDiscovery_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *ServicesDiscovery_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *ServicesDiscovery_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *ServicesDiscovery_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *ServicesDiscovery_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *ServicesDiscovery_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *ServicesDiscovery_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *ServicesDiscovery_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *ServicesDiscovery_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *ServicesDiscovery_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *ServicesDiscovery_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *ServicesDiscovery_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *ServicesDiscovery_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

///////////////////////// API /////////////////////////

/**
 * Find a services by keywords.
 */
func (self *ServicesDiscovery_Client) FindServices(keywords []string) ([]*ServiceDescriptor, error) {
	rqst := new(FindServicesDescriptorRequest)
	rqst.Keywords = keywords

	rsp, err := self.c.FindServices(api.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	return rsp.GetResults(), nil
}

/**
 * Get list of service descriptor for one service with  various version.
 */
func (self *ServicesDiscovery_Client) GetServicesDescriptor(service_id string) ([]*ServiceDescriptor, error) {
	rqst := new(GetServicesDescriptorRequest)
	rsp, err := self.c.GetServicesDescriptor(api.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	return rsp.GetResults(), nil
}

/**
 * Get a list of all services descriptor for a given server.
 */
func (self *ServicesDiscovery_Client) GetServiceDescriptor(service_id string, publisher_id string) ([]*ServiceDescriptor, error) {
	rqst := new(GetServiceDescriptorRequest)
	rqst.ServiceId = service_id
	rqst.PublisherId = publisher_id

	rsp, err := self.c.GetServiceDescriptor(api.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	return rsp.Results, nil
}

/** Publish a service to service discovery **/
func (self *ServicesDiscovery_Client) PublishServiceDescriptor(descriptor *ServiceDescriptor) error {
	rqst := new(PublishServiceDescriptorRequest)
	rqst.Descriptor_ = descriptor

	// publish a service descriptor on the network.
	_, err := self.c.PublishServiceDescriptor(api.GetClientContext(self), rqst)

	return err
}

////////////////////////////////////////////////////////////////////////////////
// Service Repository Client
////////////////////////////////////////////////////////////////////////////////
type ServicesRepository_Client struct {
	cc *grpc.ClientConn
	c  ServiceRepositoryClient

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
func NewServicesRepository_Client(address string, name string) *ServicesRepository_Client {
	client := new(ServicesRepository_Client)
	api.InitClient(client, address, name)
	client.cc = api.GetClientConnection(client)
	client.c = NewServiceRepositoryClient(client.cc)
	return client
}

// Return the address
func (self *ServicesRepository_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the domain
func (self *ServicesRepository_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *ServicesRepository_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *ServicesRepository_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *ServicesRepository_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *ServicesRepository_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *ServicesRepository_Client) SetDomain(domain string) {
	self.domain = domain
}

///////////////////////// TLS /////////////////////////

// Get if the client is secure.
func (self *ServicesRepository_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *ServicesRepository_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *ServicesRepository_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *ServicesRepository_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *ServicesRepository_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *ServicesRepository_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *ServicesRepository_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *ServicesRepository_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

///////////////////////// API /////////////////////////

/**
 * Download bundle from a repository and return it as an object in memory.
 */
func (self *ServicesRepository_Client) DownloadBundle(descriptor *ServiceDescriptor, platform Platform) (*ServiceBundle, error) {

	rqst := &DownloadBundleRequest{
		Descriptor_: descriptor,
		Plaform:     platform,
	}

	stream, err := self.c.DownloadBundle(api.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	// Here I will create the final array
	var buffer bytes.Buffer
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}
		if err != nil {
			return nil, err
		}

		_, err = buffer.Write(msg.Data)
		if err != nil {
			return nil, err
		}
	}

	// The buffer that contain the
	dec := gob.NewDecoder(&buffer)
	bundle := new(ServiceBundle)
	err = dec.Decode(bundle)
	if err != nil {
		return nil, err
	}

	return bundle, err
}

/**
 * Upload a service bundle.
 */
func (self *ServicesRepository_Client) UploadBundle(discoveryId string, serviceId string, publisherId string, platform int32, packagePath string) error {

	// The service bundle...
	bundle := new(ServiceBundle)
	bundle.Plaform = Platform(platform)

	// Here I will find the service descriptor from the given information.
	discoveryService := NewServicesDiscovery_Client(discoveryId, "services_discovery")
	descriptors, err := discoveryService.GetServiceDescriptor(serviceId, publisherId)
	if err != nil {
		return err
	}

	bundle.Descriptor_ = descriptors[0]

	/*bundle.Binairies*/
	data, err := ioutil.ReadFile(packagePath)
	if err == nil {
		bundle.Binairies = data
	}

	return self.uploadBundle(bundle)
}

/**
 * Upload a bundle into the service repository.
 */
func (self *ServicesRepository_Client) uploadBundle(bundle *ServiceBundle) error {

	// Open the stream...
	stream, err := self.c.UploadBundle(api.GetClientContext(self))
	if err != nil {
		log.Fatalf("error while TestSendEmailWithAttachements: %v", err)
	}

	const BufferSize = 1024 * 5 // the chunck size.
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer) // Will write to network.
	err = enc.Encode(bundle)
	if err != nil {
		return err
	}

	for {
		var data [BufferSize]byte
		bytesread, err := buffer.Read(data[0:BufferSize])
		if bytesread > 0 {
			rqst := &UploadBundleRequest{
				Data: data[0:bytesread],
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
		return err
	}

	return nil

}
