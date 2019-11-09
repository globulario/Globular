package sql_client

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/sql/sqlpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// SQL Client Service
////////////////////////////////////////////////////////////////////////////////
type SQL_Client struct {
	cc *grpc.ClientConn
	c  sqlpb.SqlServiceClient

	// The name of the service
	name string

	// The client domain
	domain string

	// The port
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
func NewSql_Client(domain string, addresse string, hasTLS bool, keyFile string, certFile string, caFile string, token string) *SQL_Client {

	client := new(SQL_Client)

	client.domain = domain
	client.name = "sql"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client, token)
	client.c = sqlpb.NewSqlServiceClient(client.cc)

	return client
}

// Return the domain
func (self *SQL_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *SQL_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the name of the service
func (self *SQL_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *SQL_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *SQL_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *SQL_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *SQL_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *SQL_Client) GetCaFile() string {
	return self.caFile
}

// Test if a connection is found
func (self *SQL_Client) Ping(connectionId interface{}) (string, error) {

	// Here I will try to ping a non-existing connection.
	rqst := &sqlpb.PingConnectionRqst{
		Id: Utility.ToString(connectionId),
	}

	rsp, err := self.c.Ping(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.Result, err
}

// That function return the json string with all element in it.
func (self *SQL_Client) QueryContext(connectionId interface{}, query interface{}, parameters interface{}) (string, error) {

	parameters_ := strings.Split(parameters.(string), ",")
	parametersStr, _ := Utility.ToJson(parameters_)

	// The query and all it parameters.
	rqst := &sqlpb.QueryContextRqst{
		Query: &sqlpb.Query{
			ConnectionId: Utility.ToString(connectionId),
			Query:        Utility.ToString(query),
			Parameters:   parametersStr,
		},
	}

	// Because number of values can be high I will use a stream.
	stream, err := self.c.QueryContext(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	// Here I will create the final array
	data := make([]interface{}, 0)
	header := make([]map[string]interface{}, 0)

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}
		if err != nil {
			return "", err
		}

		// Get the result...
		switch v := msg.Result.(type) {
		case *sqlpb.QueryContextRsp_Header:
			// Here I receive the header information.
			json.Unmarshal([]byte(v.Header), &header)
		case *sqlpb.QueryContextRsp_Rows:
			rows := make([]interface{}, 0)
			json.Unmarshal([]byte(v.Rows), &rows)
			data = append(data, rows...)
		}
	}

	// Create object result and put header and data in it.
	result := make(map[string]interface{}, 0)
	result["header"] = header
	result["data"] = data
	resultStr, _ := json.Marshal(result)
	return string(resultStr), nil
}

func (self *SQL_Client) ExecContext(connectionId interface{}, query interface{}, parameters interface{}, tx interface{}) (string, error) {

	parameters_ := strings.Split(parameters.(string), ",")
	parametersStr, _ := Utility.ToJson(parameters_)

	rqst := &sqlpb.ExecContextRqst{
		Query: &sqlpb.Query{
			ConnectionId: Utility.ToString(connectionId),
			Query:        Utility.ToString(query),
			Parameters:   parametersStr,
		},
		Tx: Utility.ToBool(tx),
	}

	rsp, err := self.c.ExecContext(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	result := make(map[string]interface{}, 0)
	result["affectRows"] = rsp.AffectedRows
	result["lastId"] = rsp.LastId
	resultStr, _ := json.Marshal(result)

	return string(resultStr), nil
}
