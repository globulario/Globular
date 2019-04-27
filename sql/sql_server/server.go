package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"io/ioutil"
	"os"
	"strconv"

	"runtime"

	"github.com/davecourtois/Globular/sql/sqlpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"

	// The list of available drivers...
	// feel free to append the driver you need.
	// dont forgot the set correction string if you do so.

	_ "github.com/alexbrainman/odbc"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

// Keep connection information here.
type connection struct {
	Id       string // The connection id
	Name     string // The database name
	Host     string // can also be ipv4 addresse.
	Charset  string
	Driver   string // The name of the driver.
	User     string
	Password string
	Port     int32
}

func (c *connection) getConnectionString() string {
	var connectionString string

	if c.Driver == "mssql" {
		/** Connect to Microsoft Sql server here... **/
		// So I will create the connection string from info...
		connectionString += "server=" + c.Host + ";"
		connectionString += "user=" + c.User + ";"
		connectionString += "password=" + c.Password + ";"
		connectionString += "port=" + strconv.Itoa(int(c.Port)) + ";"
		connectionString += "database=" + c.Name + ";"
		connectionString += "driver=mssql"
		connectionString += "charset=" + c.Charset + ";"
	} else if c.Driver == "mysql" {
		/** Connect to oracle MySql server here... **/
		connectionString += c.User + ":"
		connectionString += c.Password + "@tcp("
		connectionString += c.Host + ":" + strconv.Itoa(int(c.Port)) + ")"
		connectionString += "/" + c.Name
		//connectionString += "encrypt=false;"
		connectionString += "?"
		connectionString += "charset=" + c.Charset + ";"

	} else if c.Driver == "odbc" {
		/** Connect with ODBC here... **/
		if runtime.GOOS == "windows" {
			connectionString += "driver=sql server;"
		} else {
			connectionString += "driver=freetds;"
		}
		connectionString += "server=" + c.Host + ";"
		connectionString += "database=" + c.Name + ";"
		connectionString += "uid=" + c.User + ";"
		connectionString += "pwd=" + c.Password + ";"
		connectionString += "port=" + strconv.Itoa(int(c.Port)) + ";"
		connectionString += "charset=" + c.Charset + ";"

	}

	return connectionString
}

type server struct {
	// The map of connection...
	connections map[string]connection
}

// Create a new SQL connection and store it for futur use. If the connection already
// exist it will be replace by the new one.
func (self *server) CreateConnection(ctx context.Context, rsqt *sqlpb.CreateConnectionRqst) (*sqlpb.CreateConnectionRsp, error) {
	// sqlpb
	fmt.Println("Try to create a new connection")
	var c connection

	// Set the connection info from the request.
	c.Id = rsqt.Connection.Id
	c.Name = rsqt.Connection.Name
	c.Host = rsqt.Connection.Host
	c.Port = rsqt.Connection.Port
	c.User = rsqt.Connection.User
	c.Password = rsqt.Connection.Password
	c.Driver = rsqt.Connection.Driver
	c.Charset = rsqt.Connection.Charset

	db, err := sql.Open(c.Driver, c.getConnectionString())

	if err != nil {
		log.Fatal(Utility.FileLine(), Utility.FunctionName(), c.getConnectionString(), err)
	}

	log.Println("------------> vagin!")
	rows, err := db.QueryContext(ctx, "SELECT first_name, last_name FROM employees.employees WHERE gender=?", "F")

	if err != nil {
		log.Fatal(Utility.FileLine(), Utility.FunctionName(), err)
	}

	defer rows.Close()

	// test only...
	for rows.Next() {
		var first_name string
		var last_name string
		if err := rows.Scan(&first_name, &last_name); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(Utility.FileLine(), Utility.FunctionName(), err)
		}
		log.Println("---> ", first_name, last_name)
	}

	// close the connection when done.
	defer db.Close()

	if err != nil {
		log.Println(Utility.FileLine(), Utility.FunctionName(), "err: ", err)
		return nil, err
	}

	// set or update the connection and save it in json file.
	self.connections[c.Id] = c

	// In that case I will save it in file.
	str, err := Utility.ToJson(self.connections)
	if err == nil {
		err := ioutil.WriteFile(os.Args[0]+".json", []byte(str), 0644)
		if err != nil {
			log.Println(Utility.FileLine(), Utility.FunctionName(), "fail to create file: ", os.Args[0]+".json")
			return nil, err
		}
		log.Println(Utility.FileLine(), Utility.FunctionName(), "file: ", os.Args[0]+".json contain the new connection.")
	} else {
		log.Println(Utility.FileLine(), Utility.FunctionName(), "fail to serialze json: ", self.connections, " error ", err.Error())
	}

	// Print the success message here.
	log.Println("Connection " + c.Id + " was created with success!")

	return &sqlpb.CreateConnectionRsp{
		Result: true,
	}, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {
	fmt.Println("Start Sql Service")

	// The first argument must be the port number to listen to.
	port := "50051"
	if len(os.Args) > 1 {
		port = os.Args[1] // The second argument must be the port number
	}

	// First of all I will creat a listener.
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// The actual server implementation.
	s_impl := new(server)
	s_impl.connections = make(map[string]connection)

	// Here I will retreive the list of connections from file if there are some...
	file, err := ioutil.ReadFile(os.Args[0] + ".json")

	if err == nil {
		json.Unmarshal([]byte(file), &s_impl.connections)
	}

	sqlpb.RegisterSqlServiceServer(s, s_impl)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
