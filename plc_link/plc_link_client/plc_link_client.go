package plc_link_client

import (
	"context"
	//	"log"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/plc_link/plc_link_pb"

	//	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type PlcLink_Client struct {
	cc *grpc.ClientConn
	c  plc_link_pb.PlcLinkServiceClient

	// The name of the service
	name string

	// The ipv4 address
	addresse string

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
func NewPlcLink_Client(domain string, addresse string, hasTLS bool, keyFile string, certFile string, caFile string) *PlcLink_Client {
	client := new(PlcLink_Client)

	client.addresse = addresse
	client.domain = domain
	client.name = "plc_link"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client)
	client.c = plc_link_pb.NewPlcLinkServiceClient(client.cc)

	return client
}

// Return the ipv4 address
func (self *PlcLink_Client) GetAddress() string {
	return self.addresse
}

// Return the domain
func (self *PlcLink_Client) GetDomain() string {
	return self.domain
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

////////////////// API ///////////////////

// Create a new connection with a plc service.
func (self *PlcLink_Client) CreateConnection(id string, domain string, address string) error {
	// Create a new connection
	rqst := &plc_link_pb.CreateConnectionRqst{
		Connection: &plc_link_pb.Connection{
			Id:      id,
			Domain:  domain,
			Address: address,
		},
	}

	_, err := self.c.CreateConnection(context.Background(), rqst)
	return err
}

// Delete a connection
func (self *PlcLink_Client) DeleteConnection(connectionId string) error {
	// Create a new connection
	rqst := &plc_link_pb.DeleteConnectionRqst{
		Id: connectionId,
	}

	_, err := self.c.DeleteConnection(context.Background(), rqst)
	return err
}

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
				Address:      src_address,
				ConnectionId: src_connectionId,
				Name:         src_tag_name,
				Label:        src_tag_label,
				TypeName:     src_tag_typeName,
				Offset:       src_offset,
			},
			Target: &plc_link_pb.Tag{
				Domain:       trg_domain,
				Address:      trg_address,
				ConnectionId: trg_connectionId,
				Name:         trg_tag_name,
				Label:        trg_tag_label,
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
