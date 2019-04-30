package Globular

import (
	"context"
	"fmt"
	"log"

	"github.com/davecourtois/Globular/ldap/ldappb"
	"google.golang.org/grpc"

	"encoding/json"
	"testing"
)

/**
TODO Create TLS connection and test it. Also use OpenLDAP as ldap server.
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
	c := ldappb.NewLdapServiceClient(cc)

	rqst := &ldappb.CreateConnectionRqst{
		Connection: &ldappb.Connection{
			Id:       "test_ldap",
			User:     "mrmfct037@UD6.UF6",
			Password: "Dowty123",
			Port:     389,
			Host:     "mon-dc-p01.UD6.UF6",
		},
	}

	rsp, err := c.CreateConnection(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}

// Test a ldap query.
func TestSearch(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := ldappb.NewLdapServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &ldappb.SearchRqst{
		Search: &ldappb.Search{
			Id:         "test_ldap",
			BaseDN:     "OU=Users,OU=MON,OU=CA,DC=UD6,DC=UF6",
			Filter:     "(objectClass=user)",
			Attributes: []string{"sAMAccountName", "givenName", "mail", "telephoneNumber", "userPrincipalName", "distinguishedName"},
		},
	}

	rsp, err := c.Search(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	values := make([][]string, 0)
	json.Unmarshal([]byte(rsp.Result), &values)
	for i := 0; i < len(values); i++ {
		log.Println(values[i])
	}

	// Now I will close the connection...
	rqst_ := &ldappb.CloseRqst{
		Id: "test_ldap",
	}

	_, err = c.Close(context.Background(), rqst_)
	if err != nil {
		log.Fatalf("error while closing the connection: %v", err)
	}

	log.Println("Connection was close successfully!")
}

// Test a ldap query.
func TestDeleteConnection(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := ldappb.NewLdapServiceClient(cc)

	// I will execute a simple ldap search here...
	rqst := &ldappb.DeleteConnectionRqst{
		Id: "test_ldap",
	}

	_, err := c.DeleteConnection(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while deleting the connection: %v", err)
	}

	log.Println("Delete connection success!")
}
