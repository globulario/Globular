package Globular

import (
	"context"
	"fmt"

	"io/ioutil"
	"log"

	"github.com/davecourtois/Globular/persistence/persistencepb"
	"google.golang.org/grpc"

	"encoding/json"
	"testing"
	//"github.com/davecourtois/Utility"
)

// Set the correct addresse here as needed.
var (
	addresse = "localhost:10005"
)

/**
 * Get the client connection.
 */
func getClientConnection() *grpc.ClientConn {
	var err error
	var cc *grpc.ClientConn
	if cc == nil {
		cc, err = grpc.Dial(addresse, grpc.WithInsecure())
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
	c := persistencepb.NewPersistenceServiceClient(cc)

	rqst := &persistencepb.CreateConnectionRqst{
		Connection: &persistencepb.Connection{
			Id:       "mongo_db_test_connection",
			Name:     "TestMongoDB",
			User:     "",
			Password: "",
			Port:     27017,
			Host:     "localhost",
			Store:    persistencepb.StoreType_MONGO,
			Timeout:  10,
		},
	}

	rsp, err := c.CreateConnection(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}

// First test create a fresh new connection...
/*func TestPersist(t *testing.T) {
	fmt.Println("Connection creation test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := persistencepb.NewPersistenceServiceClient(cc)

	stream, err := c.InsertMany(context.Background())
	if err != nil {
		log.Fatalf("error while TestSendEmailWithAttachements: %v", err)
	}

	// here you must run the sql service test before runing this test in order
	// to generate the file Employees.json
	employes := make([]map[string]interface{}, 0)

	b, err := ioutil.ReadFile("/tmp/Employees.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &employes)
	// Persist 500 rows at time to save marshaling unmarshaling cycle time.
	for i := 0; i < len(employes); i++ {
		employes_ := make([]interface{}, 0)
		for j := 0; j < 500 && i < len(employes); j++ {
			employes_ = append(employes_, employes[i])
			i++
		}
		var jsonStr []byte
		jsonStr, err = json.Marshal(employes_)
		if err != nil {
			log.Fatalln(err)
		}

		rqst := &persistencepb.InsertManyRqst{
			Id:         "mongo_db_test_connection",
			Database:   "TestMongoDB",
			Collection: "Employees",
			JsonStr:    string(jsonStr),
		}

		err = stream.Send(rqst)
		if err != nil {
			log.Fatalf("error while TestSendEmailWithAttachements: %v", err)
		}

		log.Println(i, "/", len(employes))
	}

	rsp, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Persist entities fail %v", err)
	}

	log.Println("Persist entities succed: ", rsp.Ids)
}*/
