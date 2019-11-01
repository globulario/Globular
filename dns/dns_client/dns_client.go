package dns_client

import (
	"context"
	// "log"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/dns/dnspb"

	// "github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type DNS_Client struct {
	cc *grpc.ClientConn
	c  dnspb.DnsServiceClient

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
func NewDns_Client(domain string, addresse string, hasTLS bool, keyFile string, certFile string, caFile string, token string) *DNS_Client {
	client := new(DNS_Client)
	client.addresse = addresse
	client.domain = domain
	client.name = "dns"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client, token)
	client.c = dnspb.NewDnsServiceClient(client.cc)
	return client
}

// Return the ipv4 address
func (self *DNS_Client) GetAddress() string {
	return self.addresse
}

// Return the domain
func (self *DNS_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *DNS_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *DNS_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *DNS_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *DNS_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *DNS_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *DNS_Client) GetCaFile() string {
	return self.caFile
}

func (self *DNS_Client) Resolve(domain string) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.ResolveRequest{
		Domain: domain,
	}

	rsp, err := self.c.Resolve(context.Background(), rqst)
	if err != nil {
		return "", err
	}
	return rsp.Ipv4, nil
}

func (self *DNS_Client) SetEntry(name string, ipv4 string) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetEntryRequest{
		Name: name,
		Ipv4: ipv4,
	}

	rsp, err := self.c.SetEntry(context.Background(), rqst)
	if err != nil {
		return "", err
	}
	return rsp.Message, nil
}

func (self *DNS_Client) RemoveEntry(name string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveEntryRequest{
		Name: name,
	}

	_, err := self.c.RemoveEntry(context.Background(), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetText(id string) ([]string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetTextRequest{
		Id: id,
	}

	rsp, err := self.c.GetText(context.Background(), rqst)
	if err != nil {
		return nil, err
	}

	return rsp.GetValues(), nil
}

func (self *DNS_Client) SetText(id string, values []string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetTextRequest{
		Id:     id,
		Values: values,
	}

	_, err := self.c.SetText(context.Background(), rqst)
	return err
}

func (self *DNS_Client) RemoveText(id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveTextRequest{
		Id: id,
	}

	_, err := self.c.RemoveText(context.Background(), rqst)
	if err != nil {
		return err
	}
	return nil
}
