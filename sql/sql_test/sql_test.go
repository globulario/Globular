package Globular

import (
	"context"
	"fmt"

	"log"

	"github.com/davecourtois/Globular/sql/sqlpb"
	"google.golang.org/grpc"

	"testing"
)

/**
Before using this you must have MySql install at the local host.
- https://support.rackspace.com/how-to/installing-mysql-server-on-ubuntu/
You also have to install the employee test database,
- https://github.com/datacharmer/test_db
*/

// First test create a fresh new connection...
func TestCreateConnection(t *testing.T) {

	fmt.Println("Connection creation test.")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

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
