package Globular

import (
	"context"
	"fmt"
	"io"

	"io/ioutil"
	"log"

	"encoding/json"
	"testing"

	"github.com/davecourtois/Globular/persistence/persistencepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	// So here I will read the server configuration to see if the connection
	// is secure...
	config := make(map[string]interface{})
	data, err := ioutil.ReadFile("../persistence_server/config.json")
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
	fmt.Println("Persist many test.")

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

// Test create a db, create a collection and remove it after...
func TestCreateAndDelete(t *testing.T) {
	fmt.Println("Persist one test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := persistencepb.NewPersistenceServiceClient(cc)

	rqst := &persistencepb.InsertOneRqst{
		Id:         "mongo_db_test_connection",
		Database:   "TestCreateAndDelete_DB",
		Collection: "Employees",
		JsonStr:    `{"hire_date":"2007-07-01", "last_name":"Courtois", "first_name":"Dave", "birth_data":"1976-01-28", "emp_no":200000, "gender":"M"}`,
	}

	rsp, err := c.InsertOne(context.Background(), rqst)

	if err != nil {
		log.Fatalf("TestPersistOne fail %v", err)
	}

	rqst_count := &persistencepb.CountRqst{
		Id:         "mongo_db_test_connection",
		Database:   "TestCreateAndDelete_DB",
		Collection: "Employees",
		Query:      "{}", // count all!
	}

	countRsp, err := c.Count(context.Background(), rqst_count)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("---> count is ", countRsp.Result)

	// Test drop collection.
	rqst_drop_collection := &persistencepb.DeleteCollectionRqst{
		Id:         "mongo_db_test_connection",
		Database:   "TestCreateAndDelete_DB",
		Collection: "Employees",
	}

	_, err = c.DeleteCollection(context.Background(), rqst_drop_collection)
	if err != nil {
		log.Panicln(err)
	}

	rqst_drop_db := &persistencepb.DeleteDatabaseRqst{
		Id:       "mongo_db_test_connection",
		Database: "TestCreateAndDelete_DB",
	}

	_, err = c.DeleteDatabase(context.Background(), rqst_drop_db)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(rsp.GetId())
}

func TestPersistOne(t *testing.T) {
	fmt.Println("Persist one test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := persistencepb.NewPersistenceServiceClient(cc)

	rqst := &persistencepb.InsertOneRqst{
		Id:         "mongo_db_test_connection",
		Database:   "TestMongoDB",
		Collection: "Employees",
		JsonStr:    `{"hire_date":"2007-07-01", "last_name":"Courtois", "first_name":"Dave", "birth_data":"1976-01-28", "emp_no":200000, "gender":"M"}`,
	}

	rsp, err := c.InsertOne(context.Background(), rqst)

	if err != nil {
		log.Fatalf("TestPersistOne fail %v", err)
	}

	log.Println(rsp.GetId())
}

/** Test find one **/
func TestFindOne(t *testing.T) {
	fmt.Println("Find one test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := persistencepb.NewPersistenceServiceClient(cc)

	// Retreive a single value...
	rqst := &persistencepb.FindOneRqst{
		Id:         "mongo_db_test_connection",
		Database:   "TestMongoDB",
		Collection: "Employees",
		Query:      `{"first_name": "Anneke", "last_name": "Viele"}, {"birth_date":1}`,
	}

	rsp, err := c.FindOne(context.Background(), rqst)

	if err != nil {
		log.Fatalf("FindOne fail %v", err)
	}

	log.Println(rsp.GetJsonStr())
}

/** Test find one **/
func TestFind(t *testing.T) {
	fmt.Println("Find many test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := persistencepb.NewPersistenceServiceClient(cc)

	// Retreive a single value...
	rqst := &persistencepb.FindRqst{
		Id:         "mongo_db_test_connection",
		Database:   "TestMongoDB",
		Collection: "Employees",
		Query:      `{"first_name": "Anneke"},{"birth_date":1}`,
	}

	stream, err := c.Find(context.Background(), rqst)

	if err != nil {
		log.Fatalf("TestFind fail %v", err)
	}

	for {
		results, err := stream.Recv()

		if err == io.EOF {
			// end of stream...
			break
		}
		if err != nil {
			log.Fatalf("error while CreateConnection: %v", err)
		}

		// Get the result...
		log.Println(results.JsonStr)
	}

	log.Println("--> end of find!")
}

/** Test find one **/
func TestUpdate(t *testing.T) {
	fmt.Println("Update test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := persistencepb.NewPersistenceServiceClient(cc)

	// Retreive a single value...
	rqst := &persistencepb.UpdateRqst{
		Id:         "mongo_db_test_connection",
		Database:   "TestMongoDB",
		Collection: "Employees",
		Query:      `{"emp_no": 59717}`,
		Value:      `{"$set":{"gender":"F"}}`,
	}

	rsp, err := c.Update(context.Background(), rqst)
	if err != nil {
		log.Println("---> ", err)
	}

	log.Println("---> update success!", rsp.Result)
}

/** Test remove **/
func TestRemove(t *testing.T) {
	fmt.Println("Connection creation test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := persistencepb.NewPersistenceServiceClient(cc)

	// Retreive a single value...
	rqst := &persistencepb.DeleteRqst{
		Id:         "mongo_db_test_connection",
		Database:   "TestMongoDB",
		Collection: "Employees",
		Query:      `{"emp_no": 200000}`,
	}

	rsp, err := c.Delete(context.Background(), rqst)
	if err != nil {
		log.Println("---> ", err)
	}

	log.Println("---> Delete success!", rsp.Result)
}
