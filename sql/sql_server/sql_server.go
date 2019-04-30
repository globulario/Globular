package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"net"
	"path/filepath"

	//	"strings"
	"time"

	"io/ioutil"
	"os"
	"strconv"

	//	"net/http"
	"os/signal"
	"reflect"
	"runtime"

	"github.com/davecourtois/Globular/sql/sqlpb"
	"github.com/davecourtois/Utility"

	//"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

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

	// The global attribute of the services.
	Name     string
	Path     string
	Port     int
	Protocol string

	// The map of connection...
	Connections map[string]connection
}

func (self *server) init() {
	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := ioutil.ReadFile(dir + "/config.json")
	if err == nil {
		json.Unmarshal([]byte(file), self)
	} else {
		self.save()
	}
}

func (self *server) save() error {
	// Create the file...
	str, err := Utility.ToJson(self)
	if err != nil {
		return err
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	ioutil.WriteFile(dir+"/config.json", []byte(str), 0644)
	return nil
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
		// codes.
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// close the connection when done.
	defer db.Close()

	// set or update the connection and save it in json file.
	self.Connections[c.Id] = c

	// In that case I will save it in file.
	err = self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// test if the connection is reacheable.
	_, err = self.ping(ctx, c.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Print the success message here.
	log.Println("Connection " + c.Id + " was created with success!")

	return &sqlpb.CreateConnectionRsp{
		Result: true,
	}, nil
}

// Remove a connection from the map and the file.
func (self *server) DeleteConnection(ctx context.Context, rqst *sqlpb.DeleteConnectionRqst) (*sqlpb.DeleteConnectionRsp, error) {
	id := rqst.GetId()
	if _, ok := self.Connections[id]; !ok {
		return &sqlpb.DeleteConnectionRsp{
			Result: true,
		}, nil
	}

	delete(self.Connections, id)

	// In that case I will save it in file.
	err := self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// return success.
	return &sqlpb.DeleteConnectionRsp{
		Result: true,
	}, nil
}

// local implementation.
func (self *server) ping(ctx context.Context, id string) (string, error) {
	if _, ok := self.Connections[id]; !ok {
		return "", errors.New("connection with id " + id + " dosent exist.")
	}

	c := self.Connections[id]

	// First of all I will try to
	db, err := sql.Open(c.Driver, c.getConnectionString())
	if err != nil {
		return "", status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// If there is no answer from the database after one second
	if err := db.PingContext(ctx); err != nil {
		return "", err
	}

	return "pong", nil
}

// Ping a sql connection.
func (self *server) Ping(ctx context.Context, rsqt *sqlpb.PingConnectionRqst) (*sqlpb.PingConnectionRsp, error) {
	pong, err := self.ping(ctx, rsqt.GetId())

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &sqlpb.PingConnectionRsp{
		Result: pong,
	}, nil
}

// The maximum results size before send it over the network.
// if the number is to big network fragmentation will slow down the transfer
// if is to low the serialisation cost will be very hight...
var maxSize = uint(16000) // Value in bytes...

// Now the execute query.
func (self *server) QueryContext(rqst *sqlpb.QueryContextRqst, stream sqlpb.SqlService_QueryContextServer) error {

	// Be sure the connection is there.
	if _, ok := self.Connections[rqst.Query.ConnectionId]; !ok {
		return errors.New("connection with id " + rqst.Query.ConnectionId + " dosent exist.")
	}

	// Now I will open the connection.
	c := self.Connections[rqst.Query.ConnectionId]

	// First of all I will try to
	db, err := sql.Open(c.Driver, c.getConnectionString())
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The query
	query := rqst.Query.Query

	// The list of parameters
	parameters := make([]interface{}, 0)
	json.Unmarshal([]byte(rqst.Query.Parameters), &parameters)

	log.Println("Execute query: ", query, " whit parameters: ", parameters)

	// Here I the sql works.
	rows, err := db.QueryContext(stream.Context(), query, parameters...)

	if err != nil {
		log.Println(Utility.FileLine(), Utility.FunctionName(), err)
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	defer rows.Close()

	// First of all I will get the information about columns
	columns, err := rows.Columns()
	if err != nil {
		log.Println(Utility.FileLine(), Utility.FunctionName(), err)
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The columns type.
	columnsType, err := rows.ColumnTypes()
	if err != nil {
		log.Println(Utility.FileLine(), Utility.FunctionName(), err)
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// In header is not guaranty to contain a column type.
	header := make([]interface{}, len(columns))

	for i := 0; i < len(columnsType); i++ {
		column := columns[i]

		// So here I will extract type information.
		typeInfo := make(map[string]interface{})
		typeInfo["DatabaseTypeName"] = columnsType[i].DatabaseTypeName()
		typeInfo["Name"] = columnsType[i].DatabaseTypeName()

		// If the type is decimal.
		precision, scale, isDecimal := columnsType[i].DecimalSize()
		if isDecimal {
			typeInfo["Scale"] = scale
			typeInfo["Precision"] = precision
		}

		length, hasLength := columnsType[i].Length()
		if hasLength {
			typeInfo["Precision"] = length
		}

		isNull, isNullable := columnsType[i].Nullable()
		typeInfo["IsNullable"] = isNullable
		if isNullable {
			typeInfo["IsNull"] = isNull
		}

		header[i] = map[string]interface{}{"name": column, "typeInfo": typeInfo}
	}

	// serialyse the header in json and send it as first message.
	headerStr, _ := Utility.ToJson(header)

	// So the first message I will send will alway be the header...
	stream.Send(&sqlpb.QueryContextRsp{
		Header: headerStr,
	})

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	rows_ := make([]interface{}, 0)

	for rows.Next() {
		row := make([]interface{}, count)
		err := rows.Scan(scanArgs...)
		if err != nil {
			return err
		}

		for i, v := range values {
			// So here I will convert the values to Number, Boolean or String
			if v == nil {
				row[i] = nil // NULL value.
			} else {
				if Utility.IsNumeric(v) {
					row[i] = Utility.ToNumeric(v)
				} else if Utility.IsBool(v) {
					row[i] = Utility.ToBool(v)
				} else {
					// here I will simply return the sting value.
					row[i] = Utility.ToString(v)
				}
			}
		}

		rows_ = append(rows_, row)
		size := uint(uintptr(len(rows_)) * reflect.TypeOf(rows_).Elem().Size())

		if size > maxSize {
			rowStr, _ := json.Marshal(rows_)
			stream.Send(&sqlpb.QueryContextRsp{
				Rows: string(rowStr),
			})
			rows_ = make([]interface{}, 0)
		}
	}

	if len(rows_) > 0 {
		rowStr, _ := json.Marshal(rows_)
		stream.Send(&sqlpb.QueryContextRsp{
			Rows: string(rowStr),
		})
	}

	return nil
}

// Exec Query SQL CREATE and INSERT. Return the affected rows.
// Now the execute query.
func (self *server) ExecContext(ctx context.Context, rqst *sqlpb.ExecContextRqst) (*sqlpb.ExecContextRsp, error) {

	// Be sure the connection is there.
	if _, ok := self.Connections[rqst.Query.ConnectionId]; !ok {
		return nil, errors.New("connection with id " + rqst.Query.ConnectionId + " dosent exist.")
	}

	// Now I will open the connection.
	c := self.Connections[rqst.Query.ConnectionId]

	// First of all I will try to
	db, err := sql.Open(c.Driver, c.getConnectionString())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The query
	query := rqst.Query.Query

	// The list of parameters
	parameters := make([]interface{}, 0)
	json.Unmarshal([]byte(rqst.Query.Parameters), &parameters)

	log.Println("Execute query: ", query, " whit parameters: ", parameters)

	// Execute the query here.
	var lastId, affectedRows int64
	var result sql.Result

	if rqst.Tx {
		// with transaction
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.Println(Utility.FileLine(), Utility.FunctionName(), err)
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		var execErr error
		result, execErr = tx.ExecContext(ctx, query, parameters...)
		if execErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.New(fmt.Sprint("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr))
				log.Println(Utility.FileLine(), Utility.FunctionName(), err)
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}

			err = errors.New(fmt.Sprint("update failed: %v", execErr))
			log.Println(Utility.FileLine(), Utility.FunctionName(), err)
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		if err := tx.Commit(); err != nil {
			log.Println(Utility.FileLine(), Utility.FunctionName(), err)
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	} else {
		// without transaction
		result, err = db.ExecContext(ctx, query, parameters...)
	}

	if err != nil {
		log.Println(Utility.FileLine(), Utility.FunctionName(), err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// So here I will stream affected row if there one.
	affectedRows, err = result.RowsAffected()
	if err != nil {
		log.Println(Utility.FileLine(), Utility.FunctionName(), err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if affectedRows != 1 {
		err := errors.New(fmt.Sprint("expected to affect 1 row, affected %d", affectedRows))
		log.Println(Utility.FileLine(), Utility.FunctionName(), err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// I will send back the last id and the number of affected rows to the caller.
	lastId, _ = result.LastInsertId()
	return &sqlpb.ExecContextRsp{
		LastId:       lastId,
		AffectedRows: affectedRows,
	}, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {
	log.Println("Sql grpc service is starting")
	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "sql_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := 50051
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// First of all I will creat a listener.
	//
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Connections = make(map[string]connection)
	s_impl.Name = "sql_server"
	s_impl.Port = port
	s_impl.Protocol = "grpc"
	s_impl.Path = os.Args[0] // keep the execution path here...

	// Here I will retreive the list of connections from file if there are some...
	s_impl.init()

	sqlpb.RegisterSqlServiceServer(grpcServer, s_impl)

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		// no web-rpc server.
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		/*wrappedServer := grpcweb.WrapServer(grpcServer)
		handler := func(resp http.ResponseWriter, req *http.Request) {
			// Redirect gRPC and gRPC-Web requests to the gRPC-Web Websocket Proxy server
			if req.ProtoMajor == 2 && strings.Contains(req.Header.Get("Content-Type"), "application/grpc") {
				wrappedServer.ServeHTTP(resp, req)
				return
			}
		}

		httpServer := http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: http.HandlerFunc(handler),
		}

		if err := httpServer.ListenAndServe(); err != nil {
			grpclog.Fatalf("failed starting http server: %v", err)
		}*/

	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	log.Println("Sql grpc service is closed")
}
