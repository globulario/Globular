package ldap_client

import (
	"fmt"
	"log"
	"testing"

	"github.com/globulario/Globular/ldap/ldap_client"
)

var (
	// Connect to the plc client.
	client = ldap_client.NewLdap_Client("localhost", "ldap_server")
)

// First test create a fresh new connection...
func TestCreateConnection(t *testing.T) {
	fmt.Println("Connection creation test.")

	err := client.CreateConnection("test_ldap", "mrmfct037@UD6.UF6", "Dowty123", "mon-dc-p01.UD6.UF6", 389)
	if err != nil {
		log.Println(err)
	}
	log.Println("Connection created!")
}

// Test a ldap query.
func TestSearch(t *testing.T) {

	// I will execute a simple ldap search here...
	results, err := client.Search("test_ldap", "OU=Users,OU=MON,OU=CA,DC=UD6,DC=UF6", "(&(!(givenName=Machine*))(objectClass=user))", []string{"sAMAccountName", "mail", "memberOf"})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("results found: ", len(results))
	for i := 0; i < len(results); i++ {
		log.Println(results[i])
	}
}

// Test a ldap query.
func TestDeleteConnection(t *testing.T) {
	err := client.DeleteConnection("test_ldap")
	if err != nil {
		log.Println(err)
	}
	log.Println("Connection deleted!")
}
