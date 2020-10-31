package load_balancing_client

import (
	"strconv"

	"context"

	globular "github.com/globulario/Globular/services/golang/globular_client"
	"github.com/globulario/Globular/services/golang/lb/lbpb"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// Ca Client Service
////////////////////////////////////////////////////////////////////////////////

type Lb_Client struct {
	cc *grpc.ClientConn
	c  lbpb.LoadBalancingServiceClient

	// The id of the service
	id string

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

	// when the client need to close.
	close_channel chan bool

	// the channel to report load info.
	report_load_info_channel chan *lbpb.LoadInfo
}

// Create a connection to the service.
func NewLbService_Client(address string, id string) (*Lb_Client, error) {
	client := new(Lb_Client)

	err := globular.InitClient(client, address, id)
	if err != nil {
		return nil, err
	}
	client.cc, err = globular.GetClientConnection(client)
	if err != nil {
		return nil, err
	}
	client.c = lbpb.NewLoadBalancingServiceClient(client.cc)

	// Start processing load info.
	go func() {
		client.startReportLoadInfo()
	}()

	return client, nil
}

func (self *Lb_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the address
func (self *Lb_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the domain
func (self *Lb_Client) GetDomain() string {
	return self.domain
}

// Return the id of the service instance
func (self *Lb_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *Lb_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Lb_Client) Close() {
	// Close the load report loop.
	self.close_channel <- true
	self.cc.Close()
}

// Set grpc_service port.
func (self *Lb_Client) SetPort(port int) {
	self.port = port
}

// Set the client instance id.
func (self *Lb_Client) SetId(id string) {
	self.id = id
}

// Set the client name.
func (self *Lb_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Lb_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Lb_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Lb_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Lb_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Lb_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Lb_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Lb_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Lb_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Lb_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

////////////////////////////////////////////////////////////////////////////////
// Load balancing functions.
////////////////////////////////////////////////////////////////////////////////

// Start reporting client load infos.
func (self *Lb_Client) startReportLoadInfo() error {
	self.close_channel = make(chan bool)
	self.report_load_info_channel = make(chan *lbpb.LoadInfo)

	// Open the stream...
	stream, err := self.c.ReportLoadInfo(globular.GetClientContext(self))
	if err != nil {
		return err
	}

	for {
		select {
		case <-self.close_channel:
			// exit.
			break

		case load_info := <-self.report_load_info_channel:
			rqst := &lbpb.ReportLoadInfoRequest{
				Info: load_info,
			}
			stream.Send(rqst)
		}

	}

	// Close the stream.
	_, err = stream.CloseAndRecv()

	return err

}

// Simply report the load info to the load balancer service.
func (self *Lb_Client) ReportLoadInfo(load_info *lbpb.LoadInfo) {
	if self.report_load_info_channel == nil {
		return // the service is not ready to get info.
	}
	self.report_load_info_channel <- load_info
}

// Get the list of candidate for a given services.
func (self *Lb_Client) GetCandidates(serviceName string) ([]*lbpb.ServerInfo, error) {
	rqst := &lbpb.GetCanditatesRequest{
		ServiceName: serviceName,
	}

	resp, err := self.c.GetCanditates(globular.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	return resp.GetServers(), nil
}
