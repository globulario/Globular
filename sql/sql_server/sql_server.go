package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"os"
	"strconv"

	//"strings"
	"time"

	//	"net/http"
	"reflect"
	"runtime"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/sql/sql_client"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/sql/sqlpb"
	"github.com/davecourtois/Utility"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	//"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	// The list of available drivers...
	// feel free to append the driver you need.
	// dont forgot the set correction string if you do so.

	_ "github.com/alexbrainman/odbc"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	defaultPort  = 10039
	defaultProxy = 10040

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// The default domain
	domain string = "localhost"
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

	} else if c.Driver == "postgres" {
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

	} else if c.Driver == "sqlite3" {
		connectionString += c.Host + string(os.PathSeparator) + c.Name // The directory...
	}

	return connectionString
}

type server struct {

	// The global attribute of the services.
	Id                 string
	Name               string
	Proto              string
	Port               int
	Path               string
	Proxy              int
	Protocol           string
	AllowAllOrigins    bool
	AllowedOrigins     string // comma separated string.
	Domain             string
	CertAuthorityTrust string
	CertFile           string
	KeyFile            string
	TLS                bool
	Version            string
	PublisherId        string
	KeepUpToDate       bool
	KeepAlive          bool
	Permissions        []interface{} // contains the action permission for the services.

	// The grpc server.
	grpcServer *grpc.Server

	// The map of connection...
	Connections map[string]connection
}

// Globular services implementation...
// The id of a particular service instance.
func (self *server) GetId() string {
	return self.Id
}
func (self *server) SetId(id string) {
	self.Id = id
}

// The name of a service, must be the gRpc Service name.
func (self *server) GetName() string {
	return self.Name
}
func (self *server) SetName(name string) {
	self.Name = name
}

// The path of the executable.
func (self *server) GetPath() string {
	return self.Path
}
func (self *server) SetPath(path string) {
	self.Path = path
}

// The path of the .proto file.
func (self *server) GetProto() string {
	return self.Proto
}
func (self *server) SetProto(proto string) {
	self.Proto = proto
}

// The gRpc port.
func (self *server) GetPort() int {
	return self.Port
}
func (self *server) SetPort(port int) {
	self.Port = port
}

// The reverse proxy port (use by gRpc Web)
func (self *server) GetProxy() int {
	return self.Proxy
}
func (self *server) SetProxy(proxy int) {
	self.Proxy = proxy
}

// Can be one of http/https/tls
func (self *server) GetProtocol() string {
	return self.Protocol
}
func (self *server) SetProtocol(protocol string) {
	self.Protocol = protocol
}

// Return true if all Origins are allowed to access the mircoservice.
func (self *server) GetAllowAllOrigins() bool {
	return self.AllowAllOrigins
}
func (self *server) SetAllowAllOrigins(allowAllOrigins bool) {
	self.AllowAllOrigins = allowAllOrigins
}

// If AllowAllOrigins is false then AllowedOrigins will contain the
// list of address that can reach the services.
func (self *server) GetAllowedOrigins() string {
	return self.AllowedOrigins
}

func (self *server) SetAllowedOrigins(allowedOrigins string) {
	self.AllowedOrigins = allowedOrigins
}

// Can be a ip address or domain name.
func (self *server) GetDomain() string {
	return self.Domain
}
func (self *server) SetDomain(domain string) {
	self.Domain = domain
}

// TLS section

// If true the service run with TLS. The
func (self *server) GetTls() bool {
	return self.TLS
}
func (self *server) SetTls(hasTls bool) {
	self.TLS = hasTls
}

// The certificate authority file
func (self *server) GetCertAuthorityTrust() string {
	return self.CertAuthorityTrust
}
func (self *server) SetCertAuthorityTrust(ca string) {
	self.CertAuthorityTrust = ca
}

// The certificate file.
func (self *server) GetCertFile() string {
	return self.CertFile
}
func (self *server) SetCertFile(certFile string) {
	self.CertFile = certFile
}

// The key file.
func (self *server) GetKeyFile() string {
	return self.KeyFile
}
func (self *server) SetKeyFile(keyFile string) {
	self.KeyFile = keyFile
}

// The service version
func (self *server) GetVersion() string {
	return self.Version
}
func (self *server) SetVersion(version string) {
	self.Version = version
}

// The publisher id.
func (self *server) GetPublisherId() string {
	return self.PublisherId
}
func (self *server) SetPublisherId(publisherId string) {
	self.PublisherId = publisherId
}

func (self *server) GetKeepUpToDate() bool {
	return self.KeepUpToDate
}
func (self *server) SetKeepUptoDate(val bool) {
	self.KeepUpToDate = val
}

func (self *server) GetKeepAlive() bool {
	return self.KeepAlive
}
func (self *server) SetKeepAlive(val bool) {
	self.KeepAlive = val
}

func (self *server) GetPermissions() []interface{} {
	return self.Permissions
}
func (self *server) SetPermissions(permissions []interface{}) {
	self.Permissions = permissions
}

// Create the configuration file if is not already exist.
func (self *server) Init() error {

	// That function is use to get access to other server.
	Utility.RegisterFunction("NewSql_Client", sql_client.NewSql_Client)

	// Get the configuration path.
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	err := api.InitService(dir+"/config.json", self)
	if err != nil {
		return err
	}

	// Initialyse GRPC server.
	self.grpcServer, err = api.InitGrpcServer(self, Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor)
	if err != nil {
		return err
	}

	return nil

}

// Save the configuration values.
func (self *server) Save() error {
	// Create the file...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return api.SaveService(dir+"/config.json", self)
}

func (self *server) Start() error {
	return api.StartService(self, self.grpcServer)
}

func (self *server) Stop() error {
	return api.StopService(self)
}

//////////////////////// SQL Specific services /////////////////////////////////

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
	err = self.Save()
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
	err := self.Save()
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

	defer db.Close()

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

	defer db.Close()

	// The query
	query := rqst.Query.Query

	// The list of parameters
	parameters := make([]interface{}, 0)
	if len(rqst.Query.Parameters) > 0 {
		err = json.Unmarshal([]byte(rqst.Query.Parameters), &parameters)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))

		}
	}

	// Here I the sql works.
	rows, err := db.QueryContext(stream.Context(), query, parameters...)

	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	defer rows.Close()

	// First of all I will get the information about columns
	columns, err := rows.Columns()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The columns type.
	columnsType, err := rows.ColumnTypes()
	if err != nil {
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
		Result: &sqlpb.QueryContextRsp_Header{
			Header: headerStr,
		},
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
				Result: &sqlpb.QueryContextRsp_Rows{
					Rows: string(rowStr),
				},
			})
			rows_ = make([]interface{}, 0)
		}
	}

	if len(rows_) > 0 {
		rowStr, _ := json.Marshal(rows_)
		stream.Send(&sqlpb.QueryContextRsp{
			Result: &sqlpb.QueryContextRsp_Rows{
				Rows: string(rowStr),
			},
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

	defer db.Close()

	// The query
	query := rqst.Query.Query

	// The list of parameters
	parameters := make([]interface{}, 0)
	json.Unmarshal([]byte(rqst.Query.Parameters), &parameters)

	// Execute the query here.
	var lastId, affectedRows int64
	var result sql.Result

	if rqst.Tx {
		// with transaction
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		var execErr error
		result, execErr = tx.ExecContext(ctx, query, parameters...)
		if execErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.New(fmt.Sprint("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr))
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}

			err = errors.New(fmt.Sprint("update failed: %v", execErr))
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		if err := tx.Commit(); err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	} else {
		// without transaction
		result, err = db.ExecContext(ctx, query, parameters...)
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// So here I will stream affected row if there one.
	affectedRows, err = result.RowsAffected()
	if err != nil {
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

	// set the logger.
	//grpclog.SetLogger(log.New(os.Stdout, "sql_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := defaultPort
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Connections = make(map[string]connection)
	s_impl.Name = string(sqlpb.File_sql_sqlpb_sql_proto.Services().Get(0).FullName())
	s_impl.Proto = sqlpb.File_sql_sqlpb_sql_proto.Path()
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	// TODO set it from the program arguments...
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins
	s_impl.PublisherId = domain
	s_impl.Permissions = make([]interface{}, 0)

	// Here I will retreive the list of connections from file if there are some...
	err := s_impl.Init()
	if err != nil {
		log.Fatalf("Fail to initialyse service %s: %s", s_impl.Name, s_impl.Id, err)
	}

	// Register the echo services
	sqlpb.RegisterSqlServiceServer(s_impl.grpcServer, s_impl)
	reflection.Register(s_impl.grpcServer)

	// Start the service.
	s_impl.Start()
}
