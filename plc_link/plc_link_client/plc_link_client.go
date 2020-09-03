package plc_link_client

import (
	"context"
	//	"log"
	"strconv"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/plc_link/plc_linkpb"

	//	"github.com/davecourtois/Utility"

	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type PlcLink_Client struct {
	cc *grpc.ClientConn
	c  plc_link_pb.PlcLinkServiceClient

	// The id of the service
	id string

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
func NewPlcLink_Client(domain string, port int, hasTLS bool, keyFile string, certFile string, caFile string, token string) (*PlcLink_Client, error) {
	client := new(PlcLink_Client)

	client.domain = domain
	client.name = "plc_link"
	client.port = port
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	var err error
	client.cc, err = api.GetClientConnection(client)
	client.c = plc_link_pb.NewPlcLinkServiceClient(client.cc)

	return client, err
}

func (self *PlcLink_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = api.GetClientContext(self)
	}
	return api.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the domain
func (self *PlcLink_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *PlcLink_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the id of the service instance
func (self *PlcLink_Client) GetId() string {
	return self.id
}

// Set the client id.
func (self *PlcLink_Client) SetId(id string) {
	self.id = id
}

// Return the name of the service
func (self *PlcLink_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *PlcLink_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *PlcLink_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *PlcLink_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *PlcLink_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *PlcLink_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *PlcLink_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *PlcLink_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *PlcLink_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *PlcLink_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

// Set grpc_service port.
func (self *PlcLink_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *PlcLink_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *PlcLink_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// API ///////////////////

/**
 * Link tow tag together.
 */
func (self *PlcLink_Client) Link(id string, frequency int32, src_domain string, src_address string, src_connectionId string, src_tag_name string, src_tag_label string, src_tag_typeName string, src_offset int32, trg_domain string, trg_address string, trg_connectionId string, trg_tag_name string, trg_tag_label string, trg_tag_typeName string, trg_offset int32) error {

	rqst := &plc_link_pb.LinkRqst{
		Lnk: &plc_link_pb.Link{
			Id:        id,
			Frequency: frequency,
			Source: &plc_link_pb.Tag{
				Domain:       src_domain,
				ConnectionId: src_connectionId,
				Name:         src_tag_name,
				TypeName:     src_tag_typeName,
				Offset:       src_offset,
			},
			Target: &plc_link_pb.Tag{
				Domain:       trg_domain,
				ConnectionId: trg_connectionId,
				Name:         trg_tag_name,
				TypeName:     trg_tag_typeName,
				Offset:       trg_offset,
			},
		},
	}

	_, err := self.c.Link(context.Background(), rqst)
	return err
}

/**
 * UnLink a tag.
 */
func (self *PlcLink_Client) UnLink(id string) error {
	rqst := &plc_link_pb.UnLinkRqst{
		Id: id,
	}

	_, err := self.c.UnLink(context.Background(), rqst)
	return err
}

/**
 * Suspend Tag synchronization.
 */
func (self *PlcLink_Client) Suspend(id string) error {
	rqst := &plc_link_pb.SuspendRqst{
		Id: id,
	}

	_, err := self.c.Suspend(context.Background(), rqst)
	return err
}

/**
 * Resume Tag synchronization.
 */
func (self *PlcLink_Client) Resume(id string) error {
	rqst := &plc_link_pb.ResumeRqst{
		Id: id,
	}

	_, err := self.c.Resume(context.Background(), rqst)
	return err
}
