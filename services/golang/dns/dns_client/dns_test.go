package dns_client

import (
	//"encoding/json"
	"log"
	"testing"

	"github.com/globulario/Globular/dns/dns_client"
	"github.com/davecourtois/Utility"
)

var (
	// Try to connect to a nameserver.
	client *dns_client.DNS_Client
)

// Test various function here.
func TestSetA(t *testing.T) {
	var err error
	client, err = dns_client.NewDns_Client("ns2.globular.app", "dns_server")
	if err != nil {
		log.Println("fail to get domain ", err)
		return
	}

	// Set ip address
	domain, err := client.SetA("globular.io", Utility.MyIP(), 1000)
	if err == nil {
		log.Println(err)
	}
	log.Println("----> domain registered "+domain, Utility.MyIP())
}

func TestResolve(t *testing.T) {

	// Connect to the plc client.
	log.Println("---> test resolve A")
	ipv4, err := client.GetA("globular.io")
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
