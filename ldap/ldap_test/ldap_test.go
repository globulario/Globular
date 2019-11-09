package Globular

import (
	"context"
	"fmt"
	"log"

	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/davecourtois/Globular/ldap/ldappb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	addresse = "localhost:10033"
)

/**
TODO Create TLS connection and test it. Also use OpenLDAP as ldap server.
*/

func getClientConnection() *grpc.ClientConn {
	// So here I will read the server configuration to see if the connection
	// is secure...
	config := make(map[string]interface{})
	data, err := ioutil.ReadFile("../ldap_server/config.json")
	if err != nil {
		log.Fatal("fail to read configuration")
	}

	// Read the config file.
	json.Unmarshal(data, &config)

	var cc *grpc.ClientConn
	if cc == nil {
		if config["TLS"].(bool) {
			creds, sslErr := credentials.NewClientTLSFromFile(config["CertAuthorityTrust"].(string), "")
			if err != nil {
				log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			}
			opts := grpc.WithTransportCredentials(creds)
			cc, err = grpc.Dial(addresse, opts)
		} else {
			cc, err = grpc.Dial(addresse, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
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

	// Create a new connection
	rqst := &ldappb.CreateConnectionRqst{
		Connection: &ldappb.Connection{
			Id:       "test_ldap",
			User:     "mrmfct037@UD6.UF6",
			Password: "Dowty123",
			Port:     389,
			Host:     "www.globular.app", //"mon-dc-p01.UD6.UF6",
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
			BaseDN:     "OU=Shared,OU=Users,OU=MON,OU=CA,DC=UD6,DC=UF6",
			Filter:     "(&(givenName=Machine*)(objectClass=user))",
			Attributes: []string{"sAMAccountName", "givenName", "mail"},
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
