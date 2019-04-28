package Globular

import (
	"context"
	"fmt"

	"io"
	"log"

	"encoding/json"

	"github.com/davecourtois/Globular/sql/sqlpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"

	"testing"
)

/**
Before using this you must have MySql install at the local host.
- https://support.rackspace.com/how-to/installing-mysql-server-on-ubuntu/
You also have to install the employee test database,
- https://github.com/datacharmer/test_db
*/

func getClientConnection() *grpc.ClientConn {
	var err error
	var cc *grpc.ClientConn
	if cc == nil {
		cc, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}

	}
	return cc
}

// First test create a fresh new connection...
func TestCreateConnection(t *testing.T) {

	fmt.Println("Connection creation test.")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	rqst := &sqlpb.CreateConnectionRqst{
		Connection: &sqlpb.Connection{
			Id:       "employees_db",
			Name:     "employees",
			User:     "test",
			Password: "password",
			Port:     3306,
			Driver:   "mysql",
			Host:     "localhost",
			Charset:  "utf8",
		},
	}

	rsp, err := c.CreateConnection(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}

// Ping a connection,
// ** there is 1 second delay before the ping give up..
func TestPingConnection(t *testing.T) {
	fmt.Println("Ping connectio test.")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	// Here I will try a success case...
	rqst := &sqlpb.PingConnectionRqst{
		Id: "employees_db",
	}

	rsp, err := c.Ping(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)

	// Here I will try to ping a non-existing connection.
	rqst = &sqlpb.PingConnectionRqst{
		Id: "non_existing_connection_id",
	}

	rsp, err = c.Ping(context.Background(), rqst)
	if err != nil {
		t.Log("Expected error: ", err.Error())
	}
}

// Test some sql queries here...

// Test a simple query that return first_name and last_name.
func TestQueryContext(t *testing.T) {

	fmt.Println("Test running a sql query")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	// The query and all it parameters.
	query := "SELECT first_name, last_name FROM employees.employees WHERE gender=?"
	parameters, _ := Utility.ToJson([]string{"F"})

	log.Println(parameters)
	rqst := &sqlpb.QueryContextRqst{
		Query: &sqlpb.Query{
			ConnectionId: "employees_db",
			Query:        query,
			Parameters:   parameters,
		},
	}

	// Because number of values can be high I will use a stream.
	stream, err := c.QueryContext(context.Background(), rqst)
	if err != nil {
		log.Fatalf("Query error %v", err)
	}

	// Here I will create the final array
	data := make([]interface{}, 0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}
		if err != nil {
			log.Fatalf("error while CreateConnection: %v", err)
		}

		// Get the result...
		if len(msg.GetHeader()) > 0 {
			// Here I receive the header information.
			log.Println(msg.GetHeader())
		} else {
			result := msg.GetRows()
			rows := make([]interface{}, 0)
			json.Unmarshal([]byte(result), &rows)
			data = append(data, rows...)
		}

	}
	log.Println("---> all data was here ", len(data))
}

// Remove the test connection from the service.
func TestDeleteConnection(t *testing.T) {

	fmt.Println("Connection delete test.")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	rqst := &sqlpb.DeleteConnectionRqst{
		Id: "employees_db",
	}

	rsp, err := c.DeleteConnection(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while Delete connection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}
