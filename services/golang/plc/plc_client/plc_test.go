package plc_client

import (
	"log"
	"testing"

	"github.com/globulario/Globular/plc/plc_client"
)

var (
	// Connect to the plc client.
	ab_client = plc_client.NewPlc_Client("localhost", "plc_server_ab")
	//simens_client = plc_client.NewPlc_Client("localhost", 10025, false, "", "", "", "")
	// client = plc_client.NewPlc_Client("localhost", 10023, false, "", "", "")
)

/*
func TestCreateConnection(t *testing.T) {

	err := ab_client.CreateConnection("test", "192.168.0.154", 1.0, 0.0, 0.0, 0.0, 3001)
	if err != nil {
		log.Println("---> error: ", err)
	}
}*/

/*func TestWriteTag(t *testing.T) {
	// Write tag value.
	err := ab_client.WriteTag("test", "Data[20]", 4.0, "0.1", 0)
	if err != nil {
		log.Println(err)
	}
}*/

/*func TestSimensReadTag(t *testing.T) {

	val, err := simens_client.ReadTag("PLC-Siemens", "DB200", 4.0, 0)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("---> value: ", val)
}*/

func TestAbReadTag(t *testing.T) {

	val, err := ab_client.ReadTag("plc_ab_connection_1", "VALUE_REAL_OUT[25]", 4.0, 0)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("---> value: ", val)
}

func TestDeleteConnection(t *testing.T) {

	//client.DeleteConnection("test")
}
