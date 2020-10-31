package plc_client

import (
	"strconv"

	"context"
	"encoding/json"

	globular "github.com/globulario/Globular/services/golang/globular_client"
	"github.com/globulario/Globular/services/golang/plc/plcpb"

	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type Plc_Client struct {
	cc *grpc.ClientConn
	c  plcpb.PlcServiceClient

	// The id of the service
	id string

	// The name of the service
	name string

	// The client domain
	domain string

	// The service port
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

// Create a connection to the service. The port must be
// the http(s) port where configuration can be reach.
// After the configuration is set the port will be change to
// the actual service port.
func NewPlcService_Client(address string, id string) (*Plc_Client, error) {
	client := new(Plc_Client)
	err := globular.InitClient(client, address, id)
	if err != nil {
		return nil, err
	}
	client.cc, err = globular.GetClientConnection(client)
	if err != nil {
		return nil, err
	}
	client.c = plcpb.NewPlcServiceClient(client.cc)

	return client, nil
}

func (self *Plc_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the domain
func (self *Plc_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *Plc_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the id of the service instance
func (self *Plc_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *Plc_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Plc_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *Plc_Client) SetPort(port int) {
	self.port = port
}

// Set the client id.
func (self *Plc_Client) SetId(id string) {
	self.id = id
}

// Set the client name.
func (self *Plc_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Plc_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Plc_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Plc_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Plc_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Plc_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Plc_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Plc_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Plc_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Plc_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

/////////////////// Services functions /////////////////////////

// Stop the service.
func (self *Plc_Client) StopService() {
	self.c.Stop(globular.GetClientContext(self), &plcpb.StopRequest{})
}

// Create a new datastore connection.
func (self *Plc_Client) CreateConnection(connectionId string, ip string, cpuType float64, protocolType float64, portType float64, slot float64, timeout float64, save bool) error {
	// Create a new connection
	rqst := &plcpb.CreateConnectionRqst{
		Connection: &plcpb.Connection{
			Id:       connectionId,
			Ip:       ip,
			Protocol: plcpb.ProtocolType(protocolType),
			Cpu:      plcpb.CpuType(cpuType),
			PortType: plcpb.PortType(portType),
			Slot:     int32(Utility.ToInt(slot)),
			Timeout:  int64(Utility.ToInt(timeout)),
		},
		Save: save,
	}

	_, err := self.c.CreateConnection(globular.GetClientContext(self), rqst)
	return err
}

// Delete a connection
func (self *Plc_Client) DeleteConnection(connectionId string) error {
	// Create a new connection
	rqst := &plcpb.DeleteConnectionRqst{
		Id: connectionId,
	}

	_, err := self.c.DeleteConnection(globular.GetClientContext(self), rqst)
	return err
}

// Now the plc action function.

// Write tag value.
func (self *Plc_Client) WriteTag(connectionId string, name string, tagType float64, values []interface{}, offset int32, length int32, unsigned bool) error {
	jsonStr, _ := Utility.ToJson(values)
	rqst := &plcpb.WriteTagRqst{
		ConnectionId: connectionId,
		Name:         name,
		Type:         plcpb.TagType(Utility.ToInt(tagType)),
		Values:       jsonStr,
		Offset:       offset,
		Length:       length,
		Unsigned:     unsigned,
	}

	_, err := self.c.WriteTag(globular.GetClientContext(self), rqst)
	return err
}

// Read tag value.
func (self *Plc_Client) ReadTag(connectionId string, name string, tagType float64, offset int32, length int32, unsigned bool) ([]interface{}, error) {
	rqst := &plcpb.ReadTagRqst{
		ConnectionId: connectionId,
		Name:         name,
		Type:         plcpb.TagType(Utility.ToInt(tagType)),
		Offset:       offset,
		Length:       length,
		Unsigned:     unsigned,
	}

	rsp, err := self.c.ReadTag(globular.GetClientContext(self), rqst)
	values := make([]interface{}, 0)
	json.Unmarshal([]byte(rsp.GetValues()), &values)

	if err == nil {
		return values, nil
	}
	return nil, err
}

// Get connection
func (self *Plc_Client) GetConnection(connectionId string) (*plcpb.Connection, error) {
	rqst := &plcpb.GetConnectionRqst{
		Id: connectionId,
	}

	rsp, err := self.c.GetConnection(globular.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	return rsp.Connection, err
}

// Delete All tags.
func (self *Plc_Client) CloseConnection(connectionId string) error {
	rqst := &plcpb.CloseConnectionRqst{
		ConnectionId: connectionId,
	}

	_, err := self.c.CloseConnection(globular.GetClientContext(self), rqst)

	return err
}
