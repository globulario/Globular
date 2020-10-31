package sql_client

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/globulario/Globular/sql/sql_client"
)

// Set the correct addresse here as needed.
var (
	client = sql_client.NewSql_Client("localhost", "sql_server") // connect with the local service.
)

// First test create a fresh new connection...
func TestCreateConnection(t *testing.T) {
	// The test will be made with the sqlite for simplicity.
	fmt.Println("Connection creation test.")
	err := client.CreateConnection("employees_db", "employees_db", "sqlite3", "", "", "/tmp", 0, "")
	if err != nil {
		log.Println("Fail to run CreateConnection ", err)
		return
	}

	log.Println("Succed to create sql connection")
}

// Ping a connection,
// ** there is 1 second delay before the ping give up..
func TestPingConnection(t *testing.T) {
	fmt.Println("Ping connectio test.")
	pong, err := client.Ping("employees_db")
	if err != nil {
		log.Println("Fail to run Ping ", err)
	}
	log.Println("Ping success ", pong)
}

// Test some sql queries here...
func TestCreateTable(t *testing.T) {
	query := "CREATE TABLE IF NOT EXISTS employees (id INTERGER PRIMARY KEY, firstname TEXT, lastname TEXT, gender TEXT)"
	_, err := client.ExecContext("employees_db", query, "[]", nil)
	if err != nil {
		log.Println("Fail to run TestCreateTable ", err)
	}
	log.Println("TestCreateTable success ")
}

// Test a simple query that return first_name and last_name.
func TestInsertValue(t *testing.T) {
	// Test create query...
	query := "INSERT INTO employees (id, firstname, lastname, gender) VALUES (?, ?, ?, ?);"

	data, err := client.ExecContext("employees_db", query, `[1, "Dave", "Courtois", "M"]`, nil)
	if err != nil {
		log.Println("------> fail to insert a new employe ", err)
		return
	}

	log.Println("Value insert number: ", data)

}

// Test a simple query that return first_name and last_name.
func TestQueryContext(t *testing.T) {

	fmt.Println("Test running a sql query")

	// The query and all it parameters.
	query := `SELECT * FROM employees WHERE firstname LIKE ?`

	// The employee db.
	data, err := client.QueryContext("employees_db", query, `["D%"]`)
	if err != nil {
		log.Println("Fail to read employee db", err)
	}

	results := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(data), &results)
	log.Println(results["data"])

}

// Test upatade value
func TestUpdateValue(t *testing.T) {
	// Test create query...
	query := "UPDATE employees SET firstname=? WHERE id = ?;"

	data, err := client.ExecContext("employees_db", query, `["David", 1]`, nil)
	if err != nil {
		log.Println("Fail to update employee db", err)
	}

	log.Println(data)
}

// Test delete value
func TestDeleteValue(t *testing.T) {
	// Test create query...
	query := "DELETE FROM employees WHERE id = ?;"
	data, err := client.ExecContext("employees_db", query, `[1]`, nil)
	if err != nil {
		log.Println("Fail to update employee db", err)
	}

	log.Println(data)
}

func TestDeleteTable(t *testing.T) {
	// Test create query...
	query := "DROP TABLE IF EXISTS employees"
	data, err := client.ExecContext("employees_db", query, `[]`, nil)
	if err != nil {
		log.Println("Fail to drop employee table", err)
	}

	log.Println(data)
}

// Remove the test connection from the service.
func TestDeleteConnection(t *testing.T) {

	fmt.Println("Connection delete test.")
	err := client.DeleteConnection("employees_db")
	if err != nil {
		log.Println("----> fail to delete connection.", err)
		return
	}
	log.Println("TestDeleteConnection succed!")
}
