package persistence_client

import (
	//"fmt"
	//"io/ioutil"
	"log"

	"testing"

	"github.com/globulario/Globular/persistence/persistence_client"
	//"github.com/davecourtois/Utility"
)

// start: mongod --dbpath E:\Project\src\github.com\davecourtois\Globular\data\mongodb-data
// Set the correct addresse here as needed.
var (

	// Connect to the plc client.
	client = persistence_client.NewPersistence_Client("localhost", "persistence.PersistenceService")
)

// First test create a fresh new connection...
/*
func TestCreateConnection(t *testing.T) {
	//fmt.Println("Connection creation test.")
	user := "sa"
	pwd := "adminadmin"
	err := client.CreateConnection("mongo_db_test_connection", "TestMongoDB", "localhost", 27017, 0, user, pwd, 500, "", true)
	if err != nil {
		log.Println("fail to create connection! ", err)
	}

}
*/
func TestAggregate(t *testing.T) {
	//fmt.Println("Aggregate")
	/*user := "sa"
	pwd := "adminadmin"
	err := client.CreateConnection("mongo_db_test_connection", "local_ressource", "localhost", 27017, 0, user, pwd, 500, "", true)
	if err != nil {
		log.Println("fail to create connection! ", err)
	}

	Id := "mongo_db_test_connection"
	Database := "local_ressource"
	Collection := "Employees"

	results, err := client.Aggregate(Id, Database, Collection, `[{"$count":"toto"}]`, "")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("---> ", results)
	*/
}

/*
func TestPingConnection(t *testing.T) {
	log.Println("Test ping connection")

	err := client.Ping("mongo_db_test_connection")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Ping mongo_db_test_connection successed!")
}


// First test create a fresh new connection...
/*func TestPersistMany(t *testing.T) {
	fmt.Println("Persist many test.")

	b, err := ioutil.ReadFile("/tmp/Employees.json")
	if err != nil {
		log.Fatal(err)
	}

	Id := "mongo_db_test_connection"
	Database := "TestCreateAndDelete_DB"
	Collection := "Employees"

	ids, err := client.InsertMany(Id, Database, Collection, string(b), "")
	if err != nil {
		log.Fatalf("TestPersistMany fail %v", err)
	}

	log.Println(ids)
}*/

/*func TestPersistOne(t *testing.T) {

	Id := "mongo_db_test_connection"
	Database := "TestMongoDB"
	Collection := "Employees"
	JsonStr := `{"hire_date":"2007-07-01", "last_name":"Courtois", "first_name":"Dave", "birth_date":"1976-01-28", "emp_no":200000, "gender":"M"}`
	id, err := client.InsertOne(Id, Database, Collection, JsonStr, "")

	if err != nil {
		log.Fatalf("TestPersistOne fail %v", err)
	}

	log.Println("one entity persist with id ", id)
}*/
/*
func TestCreateConnection(t *testing.T) {
	fmt.Println("Connection creation test.")
	user := "chitchat"
	pwd := "8b6021d3-5e4a-38a0-9057-0e127754938b"
	err := client.CreateConnection("chitchat_db", "chitchat_db", "localhost", 27017, 0, user, pwd, 500, "", true)
	if err != nil {
		log.Println("fail to create connection! ", err)
	}
}*/

/*func TestPersistOne(t *testing.T) {

	Id := "chitchat_db"
	Database := "chitchat_db"
	Collection := "Room"
	JsonStr := `{"_id":"test"}`
	id, err := client.InsertOne(Id, Database, Collection, JsonStr, "")

	if err != nil {
		log.Fatalf("TestPersistOne fail %v", err)
	}

	log.Println("one entity persist with id ", id)
}*/

/** Test find one **/
/*func TestUpdate(t *testing.T) {
	fmt.Println("Update test.")

	Id := "mongo_db_test_connection"
	Database := "TestMongoDB"
	Collection := "Employees"
	Query := `{"emp_no": 200000}`
	Value := `{"$set":{"gender":"F"}}`

	err := client.Update(Id, Database, Collection, Query, Value, "")
	if err != nil {
		log.Fatalf("TestUpdate fail %v", err)
	}
	log.Println("---> update success!")
}*/

/** Test find many **/
/*func TestFind(t *testing.T) {
	fmt.Println("Find many test.")

	Id := "mongo_db_test_connection"
	Database := "TestMongoDB"
	Collection := "Employees"
	Query := `{"first_name": "Dave"}`

	values, err := client.Find(Id, Database, Collection, Query, `[{"Projection":{"first_name":1}}]`)
	if err != nil {
		log.Fatalf("TestFind fail %v", err)
	}

	log.Println(values)
	log.Println("--> end of find!")

}*/

/** Test find one **/
/*func TestFindOne(t *testing.T) {
	fmt.Println("Find one test.")

	Id := "mongo_db_test_connection"
	Database := "TestMongoDB"
	Collection := "Employees"
	Query := `{"first_name": "Dave"}`

	values, err := client.FindOne(Id, Database, Collection, Query, "")
	if err != nil {
		log.Fatalf("TestFind fail %v", err)
	}

	log.Println(values)
}*/

/** Test remove **/
/*func TestRemove(t *testing.T) {
	fmt.Println("Test Remove")

	Id := "mongo_db_test_connection"
	Database := "TestMongoDB"
	Collection := "Employees"
	Query := `{"emp_no": 200000}`

	err := client.DeleteOne(Id, Database, Collection, Query, "")
	if err != nil {
		log.Fatalf("DeleteOne fail %v", err)
	}

	log.Println("---> Delete success!")
}*/

/*func TestRemoveMany(t *testing.T) {
	fmt.Println("Test Remove")

	Id := "mongo_db_test_connection"
	Database := "TestMongoDB"
	Collection := "Employees"
	Query := `{"emp_no": 200000}`

	err := client.Delete(Id, Database, Collection, Query, "")
	if err != nil {
		log.Fatalf("DeleteOne fail %v", err)
	}

	log.Println("---> Delete success!")
}*/

// Test create a db, create a collection and remove it after...
/*func TestCreateAndDelete(t *testing.T) {
	fmt.Println("Test Create And Delete")

	Id := "mongo_db_test_connection"
	Database := "TestCreateAndDelete_DB"
	Collection := "Employees"
	JsonStr := `{"hire_date":"2007-07-01", "last_name":"Courtois", "first_name":"Dave", "birth_data":"1976-01-28", "emp_no":200000, "gender":"M"}`

	id, err := client.InsertOne(Id, Database, Collection, JsonStr, "")

	var c int
	c, err = client.Count(Id, Database, Collection, "{}", "")

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("---> count is ", c)

	// Test drop collection.
	err = client.DeleteCollection(Id, Database, Collection)
	if err != nil {
		log.Panicln(err)
	}

	err = client.DeleteDatabase(Id, Database)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(id)

}*/

/*func TestDeleteConnection(t *testing.T) {
	fmt.Println("Connection creation test.")
	err := client.DeleteConnection("mongo_db_test_connection")
	if err != nil {
		log.Println("fail to delete connection! ", err)
	}
}*/
