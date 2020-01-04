package Globular

import (
	//"encoding/json"
	"log"
	"testing"

	"github.com/davecourtois/Globular/dns/dns_client"
	"github.com/davecourtois/Utility"
)

var (

	// The token is print in the sever console. It can be taken from the local temp file if the client service run on the same machine.
	// It will be regenerate each time the server is started and valid until the token expire.
	// client = dns_client.NewDns_Client("localhost", 10033, false, key, crt, ca, token)

	// Try to connect to a nameserver.
	client = dns_client.NewDns_Client("ns2.globular.app", "dns_server")
	//client = dns_client.NewDns_Client("www.omniscient.app", "35.183.163.145:10033", false, key, crt, ca, token)
	//client = dns_client.NewDns_Client("www.omniscient.app", "44.225.184.139:10033", false, key, crt, ca, token)

	token = ""
)

// Test various function here.
func TestSetA(t *testing.T) {
	// Set ip address
	domain, err := client.SetA("globular4", Utility.MyIP(), 1000)
	if err == nil {
		log.Println("--> your domain is ", domain)
	} else {
		log.Panicln(err)
	}
}

func TestResolve(t *testing.T) {

	// Connect to the plc client.
	log.Println("---> test resolve A")
	ipv4, err := client.GetA("globular4.omniscient.app")
	if err == nil {
		log.Println("--> your ip is ", ipv4)
	} else {
		log.Panicln(err)
	}
}

/*func TestRemoveA(t *testing.T) {

	// Connect to the plc client.
	log.Println("---> test resolve A")
	err := client.RemoveA("toto")
	if err == nil {
		log.Println("--> your A is remove!")
	} else {
		log.Panicln(err)
	}
}*/

/*func TestTextValue(t *testing.T) {
	// Connect to the plc client.
	log.Println("---> test set text")
	err := client.SetText("toto", []string{"toto", "titi", "tata"})
	if err != nil {
		log.Panicln(err)
	}

	log.Println("---> test get text")
	values, err := client.GetText("toto")
	if err != nil {
		log.Panicln(err)
	}

	log.Println("--> values retreive: ", values)

	log.Println("---> test remove text")
	err = client.RemoveText("toto")
	if err != nil {
		log.Panicln(err)
	}
}*/
