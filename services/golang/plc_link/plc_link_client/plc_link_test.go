package plc_link_client

import (
	"log"
	"testing"
)

var (
// Connect to the plc client.
//client = NewPlcLink_Client.("localhost", 10027, false, "", "", "")
//client = plc_link_client.NewPlcLink_Client("localhost", 10027, false, "", "", "")
)

// Test various function here.
func TestCreateConnection(t *testing.T) {

	// Create a connection.
	err := client.CreateConnection("ab_connection", "stevePc", "135.19.213.242:10023")
	if err != nil {
		log.Println("---> error: ", err)
	}

	err = client.CreateConnection("simens_connection", "stevePc", "135.19.213.242:10025")
	if err != nil {
		log.Println("---> error: ", err)
	}
}

func TestLinkTag(t *testing.T) {
	var id string = "lnk_1"
	var frequency int32 = 500
	var src_domain string = "stevePc"
	var src_address string = "135.19.213.242:10023"
	var src_connectionId string = "test"
	var src_tag_name string = "Data[0]"
	var src_tag_label string = "data_array_index0"
	var src_tag_typeName string = "REAL"
	var src_offset int32 = 0
	var trg_domain string = "stevePc"
	var trg_address string = "135.19.213.242:10025"
	var trg_connectionId string = "test_0"
	var trg_tag_name string = "DB200"
	var trg_tag_label string = "DB200"
	var trg_tag_typeName string = "REAL"
	var trg_offset int32 = 8

	err := client.Link(id, frequency, src_domain, src_address, src_connectionId, src_tag_name, src_tag_label, src_tag_typeName, src_offset, trg_domain, trg_address, trg_connectionId, trg_tag_name, trg_tag_label, trg_tag_typeName, trg_offset)

	if err != nil {
		log.Println("error found ", err)
	}
}
