package Globular

import (
	"context"
	"fmt"

	"io"
	"log"

	"encoding/json"

	"io/ioutil"
	"testing"

	"github.com/davecourtois/Globular/sql/sqlpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Set the correct addresse here as needed.
var (
	addresse = "localhost:10009"
)

/**
Before using this you must have MySql install at the local host.
- https://support.rackspace.com/how-to/installing-mysql-server-on-ubuntu/
You also have to install the employee test database,
- https://github.com/datacharmer/test_db
*/

func getClientConnection() *grpc.ClientConn {
	// So here I will read the server configuration to see if the connection
	// is secure...
	config := make(map[string]interface{})
	data, err := ioutil.ReadFile("../sql_server/config.json")
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
	c := sqlpb.NewSqlServiceClient(cc)

	rqst := &sqlpb.CreateConnectionRqst{
		Connection: &sqlpb.Connection{
			Id:       "employees_db",
			Name:     "employees",
			User:     "davecourtois",
			Password: "400zm89a",
			Port:     3306,
			Driver:   "mysql",
			Host:     "mysqlserver.cdlyvemhjtsp.us-west-2.rds.amazonaws.com",
			Charset:  "utf8",
		},
	}

	// Test with sql server.
	/*rqst := &sqlpb.CreateConnectionRqst{
		Connection: &sqlpb.Connection{
			Id:       "bris_outil",
			Name:     "BrisOutil",
			User:     "dbprog",
			Password: "dbprog",
			Port:     1433,
			Driver:   "odbc",
			Host:     "mon-sql-v01",
			Charset:  "utf8",
		},
	}*/

	rsp, err := c.CreateConnection(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}

// Ping a connection,
// ** there is 1 second delay before the ping give up..
func TestPingConnection(t *testing.T) {
	fmt.Println("Ping connectio test.")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	// Here I will try a success case...

	rqst := &sqlpb.PingConnectionRqst{
		Id: "employees_db",
	}

	/*
		rqst := &sqlpb.PingConnectionRqst{
			Id: "bris_outil",
		}
	*/
	rsp, err := c.Ping(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)

	// Here I will try to ping a non-existing connection.
	rqst = &sqlpb.PingConnectionRqst{
		Id: "non_existing_connection_id",
	}

	rsp, err = c.Ping(api.GetClientContext(self), rqst)
	if err != nil {
		t.Log("Expected error: ", err.Error())
	}
}

// Test some sql queries here...

// Test a simple query that return first_name and last_name.
func TestQueryContext(t *testing.T) {

	fmt.Println("Test running a sql query")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	// The query and all it parameters.
	query := "SELECT * FROM employees.employees WHERE gender=?"
	parameters, _ := Utility.ToJson([]string{"F"})

	log.Println(parameters)
	rqst := &sqlpb.QueryContextRqst{
		Query: &sqlpb.Query{
			ConnectionId: "employees_db",
			Query:        query,
			Parameters:   parameters,
		},
	}

	/*
		query := "SELECT * FROM [BrisOutil].[dbo].[Bris] WHERE product_id LIKE ?"
		parameters, _ := Utility.ToJson([]string{"50-%"})

		rqst := &sqlpb.QueryContextRqst{
			Query: &sqlpb.Query{
				ConnectionId: "bris_outil",
				Query:        query,
				Parameters:   parameters,
			},
		}
	*/
	// Because number of values can be high I will use a stream.
	stream, err := c.QueryContext(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("Query error %v", err)
	}

	// Here I will create the final array
	data := make([]interface{}, 0)
	header := make([]map[string]interface{}, 0)

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}
		if err != nil {
			log.Fatalf("error while CreateConnection: %v", err)
		}

		// Get the result...
		switch v := msg.Result.(type) {
		case *sqlpb.QueryContextRsp_Header:
			// Here I receive the header information.
			json.Unmarshal([]byte(v.Header), &header)
		case *sqlpb.QueryContextRsp_Rows:
			rows := make([]interface{}, 0)
			json.Unmarshal([]byte(v.Rows), &rows)
			data = append(data, rows...)
		}

	}

	log.Println("---> all data was here ", len(data), header)
	exportSqlToJson("Employee", header, data)
}

// Here I will export json value to a json file and use it in persistence test.
func exportSqlToJson(tableName string, header []map[string]interface{}, rows []interface{}) {
	if Utility.Exists("/tmp/" + tableName + ".json") {
		return
	}
	var objects = make([]map[string]interface{}, 0)
	for i := 0; i < len(rows); i++ {
		obj := make(map[string]interface{})
		// The tow value are needed by the store.
		for j := 0; j < len(header); j++ {
			obj[header[j]["name"].(string)] = rows[i].([]interface{})[j]
		}
		objects = append(objects, obj)
	}

	jsonStr, _ := Utility.ToJson(objects)
	_ = ioutil.WriteFile("/tmp/"+tableName+"s.json", []byte(jsonStr), 0644)

}

// Test a simple query that return first_name and last_name.
func TestInsertValue(t *testing.T) {
	// Test create query...
	query := "INSERT INTO EmailLst (account_id, product_name) VALUES (?, ?);"

	fmt.Println("Test running a sql query")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	parameters, _ := Utility.ToJson([]string{"test", "1212"})

	rqst := &sqlpb.ExecContextRqst{
		Query: &sqlpb.Query{
			ConnectionId: "bris_outil",
			Query:        query,
			Parameters:   parameters,
		},
		Tx: false,
	}

	rsp, err := c.ExecContext(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalln("Fail ExecContext ", err)
	}

	log.Println("Value insert number of rows affect: ", rsp.AffectedRows, " lastId ", rsp.LastId)

}

// Test upatade value
func TestUpdateValue(t *testing.T) {
	// Test create query...
	query := "UPDATE EmailLst SET product_name=? WHERE account_id = ?;"

	fmt.Println("Test running a sql query")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	parameters, _ := Utility.ToJson([]string{"3434", "test"})

	rqst := &sqlpb.ExecContextRqst{
		Query: &sqlpb.Query{
			ConnectionId: "bris_outil",
			Query:        query,
			Parameters:   parameters,
		},
		Tx: false,
	}

	rsp, err := c.ExecContext(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalln("Fail ExecContext ", err)
	}

	log.Println("Value insert number of rows affect: ", rsp.AffectedRows, " lastId ", rsp.LastId)
}

// Test delete value
func TestDeleteValue(t *testing.T) {
	// Test create query...
	query := "DELETE FROM EmailLst WHERE account_id = ?;"

	fmt.Println("Test running a sql query")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	parameters, _ := Utility.ToJson([]string{"test"})

	rqst := &sqlpb.ExecContextRqst{
		Query: &sqlpb.Query{
			ConnectionId: "bris_outil",
			Query:        query,
			Parameters:   parameters,
		},
		Tx: false,
	}

	rsp, err := c.ExecContext(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalln("Fail ExecContext ", err)
	}

	log.Println("Value insert number of rows affect: ", rsp.AffectedRows, " lastId ", rsp.LastId)
}

/*
// Remove the test connection from the service.
func TestDeleteConnection(t *testing.T) {

	fmt.Println("Connection delete test.")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := sqlpb.NewSqlServiceClient(cc)

	rqst := &sqlpb.DeleteConnectionRqst{
		Id: "employees_db",
	}

	rqst := &sqlpb.DeleteConnectionRqst{
		Id: "bris_outil",
	}

	rsp, err := c.DeleteConnection(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("error while Delete connection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}
*/
