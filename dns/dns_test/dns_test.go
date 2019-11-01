package Globular

import (
	//"encoding/json"
	"log"
	"testing"

	"github.com/davecourtois/Globular/dns/dns_client"
	"github.com/davecourtois/Utility"
)

var (

	// Test with a secure connection.
	// Those files must be accessible to the client to be able to call
	// gRpc function on the server.
	crt = "E:/Project/src/github.com/davecourtois/Globular/config/grpc_tls/client.crt"
	key = "E:/Project/src/github.com/davecourtois/Globular/config/grpc_tls/client.pem"
	ca  = "E:/Project//src/github.com/davecourtois/Globular/config/grpc_tls/ca.crt"

	// The token is print in the sever console. It can be taken from the local temp file if the client service run on the same machine.
	// It will be regenerate each time the server is started and valid until the token expire.
	// client = dns_client.NewDns_Client("localhost", "127.0.0.1:10033", false, key, crt, ca, token)

	// Try to connect to a nameserver.
	client = dns_client.NewDns_Client("localhost", "127.0.0.1:10033", false, key, crt, ca, token)
	//client = dns_client.NewDns_Client("www.omniscient.app", "35.183.163.145:10033", false, key, crt, ca, token)
	//client = dns_client.NewDns_Client("www.omniscient.app", "44.225.184.139:10033", false, key, crt, ca, token)

	token = ""
)

// Test various function here.
func TestSetEntry(t *testing.T) {
	// Connect to the plc client.
	log.Println("---> test set entry")
	domain, err := client.SetEntry("toto", Utility.MyIP())
	if err == nil {
		log.Println("--> your domain is ", domain)
	} else {
		log.Panicln(err)
	}
}

func TestResolve(t *testing.T) {

	// Connect to the plc client.
	log.Println("---> test resolve entry")
	ipv4, err := client.Resolve("toto.example.com.")
	if err == nil {
		log.Println("--> your ip is ", ipv4)
	} else {
		log.Panicln(err)
	}
}

func TestRemoveEntry(t *testing.T) {

	// Connect to the plc client.
	log.Println("---> test resolve entry")
	err := client.RemoveEntry("toto")
	if err == nil {
		log.Println("--> your entry is remove!")
	} else {
		log.Panicln(err)
	}
}

func TestTextValue(t *testing.T) {
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
}
