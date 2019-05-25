package main

import (
	"context"
	"log"

	"io"
	"os"

	"github.com/davecourtois/Globular/echo/echopb"
	"github.com/davecourtois/Globular/file/filepb"

	/*"github.com/davecourtois/Globular/ldap/ldappb"*/
	"strings"

	"encoding/json"

	"github.com/davecourtois/Globular/persistence/persistencepb"
	"github.com/davecourtois/Globular/smtp/smtppb"
	"github.com/davecourtois/Globular/spc/spcpb"
	"github.com/davecourtois/Globular/sql/sqlpb"
	"github.com/davecourtois/Globular/storage/storagepb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

// The client service interface.
type Client interface {

	// Close the client.
	Close()
}

/**
 * Get the client connection.
 */
func getClientConnection(addresse string) *grpc.ClientConn {
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

////////////////////////////////////////////////////////////////////////////////
// File Client Service
////////////////////////////////////////////////////////////////////////////////

type File_Client struct {
	cc *grpc.ClientConn
	c  filepb.FileServiceClient
}

// Create a connection to the service.
func NewFile_Client(addresse string) *File_Client {
	client := new(File_Client)
	client.cc = getClientConnection(addresse)
	client.c = filepb.NewFileServiceClient(client.cc)

	return client
}

// must be close when no more needed.
func (self *File_Client) Close() {
	self.cc.Close()
}

// Read the content of a dir and return it info.
func (self *File_Client) ReadDir(path interface{}, recursive interface{}, thumbnailHeight interface{}, thumbnailWidth interface{}) (string, error) {

	// Create a new client service...
	rqst := &filepb.ReadDirRequest{
		Path:           Utility.ToString(path),
		Recursive:      Utility.ToBool(recursive),
		ThumnailHeight: int32(Utility.ToInt(thumbnailHeight)),
		ThumnailWidth:  int32(Utility.ToInt(thumbnailWidth)),
	}

	stream, err := self.c.ReadDir(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	// Here I will create the final array
	data := make([]byte, 0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}

		data = append(data, msg.Data...)
		if err != nil {
			return "", err
		}
	}

	return string(data), nil
}

/**
 * Create a new directory on the server.
 */
func (self *File_Client) CreateDir(path interface{}, name interface{}) error {

	rqst := &filepb.CreateDirRequest{
		Path: Utility.ToString(path),
		Name: Utility.ToString(name),
	}

	_, err := self.c.CreateDir(context.Background(), rqst)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Read file data
 */
func (self *File_Client) ReadFile(path interface{}) ([]byte, error) {

	rqst := &filepb.ReadFileRequest{
		Path: Utility.ToString(path),
	}

	stream, err := self.c.ReadFile(context.Background(), rqst)
	if err != nil {
		return nil, err
	}

	// Here I will create the final array
	data := make([]byte, 0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}

		data = append(data, msg.Data...)
		if err != nil {
			return nil, err
		}
	}

	return data, err
}

/**
 * Rename a directory.
 */
func (self *File_Client) RenameDir(path interface{}, oldname interface{}, newname interface{}) error {

	rqst := &filepb.RenameRequest{
		Path:    Utility.ToString(path),
		OldName: Utility.ToString(oldname),
		NewName: Utility.ToString(newname),
	}

	_, err := self.c.Rename(context.Background(), rqst)

	return err
}

/**
 * Delete a directory
 */
func (self *File_Client) DeleteDir(path string) error {
	rqst := &filepb.DeleteDirRequest{
		Path: Utility.ToString(path),
	}

	_, err := self.c.DeleteDir(context.Background(), rqst)
	return err
}

/**
 * Get a single file info.
 */
func (self *File_Client) GetFileInfo(path interface{}, recursive interface{}, thumbnailHeight interface{}, thumbnailWidth interface{}) (string, error) {

	rqst := &filepb.GetFileInfoRequest{
		Path:           Utility.ToString(path),
		ThumnailHeight: int32(Utility.ToInt(thumbnailHeight)),
		ThumnailWidth:  int32(Utility.ToInt(thumbnailWidth)),
	}

	rsp, err := self.c.GetFileInfo(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return string(rsp.Data), nil
}

/**
 * That function move a file from a directory to another... (mv) in unix.
 */
func (self *File_Client) MoveFile(path interface{}, dest interface{}) error {

	// Open the stream...
	stream, err := self.c.SaveFile(context.Background())
	if err != nil {
		return err
	}

	err = stream.Send(&filepb.SaveFileRequest{
		File: &filepb.SaveFileRequest_Path{
			Path: Utility.ToString(dest), // Where the file will be save...
		},
	})

	if err != nil {
		return err
	}

	// Where the file is read from.
	file, err := os.Open(Utility.ToString(path))
	if err != nil {
		return err
	}

	// close the file when done.
	defer file.Close()

	const BufferSize = 1024 * 5 // the chunck size.
	buffer := make([]byte, BufferSize)
	for {
		bytesread, err := file.Read(buffer)
		if bytesread > 0 {
			rqst := &filepb.SaveFileRequest{
				File: &filepb.SaveFileRequest_Data{
					Data: buffer[:bytesread],
				},
			}
			err = stream.Send(rqst)
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

/**
 * Delete a file whit a given path.
 */
func (self *File_Client) DeleteFile(path string) error {

	rqst := &filepb.DeleteFileRequest{
		Path: Utility.ToString(path),
	}

	_, err := self.c.DeleteFile(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while testing get file info: %v", err)
	}

	return err
}

// Read the content of a dir and return all images as thumbnails.
func (self *File_Client) GetThumbnails(path interface{}, recursive interface{}, thumbnailHeight interface{}, thumbnailWidth interface{}) (string, error) {

	// Create a new client service...
	rqst := &filepb.GetThumbnailsRequest{
		Path:           Utility.ToString(path),
		Recursive:      Utility.ToBool(recursive),
		ThumnailHeight: int32(Utility.ToInt(thumbnailHeight)),
		ThumnailWidth:  int32(Utility.ToInt(thumbnailWidth)),
	}

	stream, err := self.c.GetThumbnails(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	// Here I will create the final array
	data := make([]byte, 0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}

		data = append(data, msg.Data...)
		if err != nil {
			return "", err
		}
	}

	return string(data), nil
}

////////////////////////////////////////////////////////////////////////////////
// SQL Client Service
////////////////////////////////////////////////////////////////////////////////
type SQL_Client struct {
	cc *grpc.ClientConn
	c  sqlpb.SqlServiceClient
}

// Create a connection to the service.
func NewSql_Client(addresse string) *SQL_Client {
	client := new(SQL_Client)
	client.cc = getClientConnection(addresse)
	client.c = sqlpb.NewSqlServiceClient(client.cc)
	return client
}

// must be close when no more needed.
func (self *SQL_Client) Close() {
	self.cc.Close()
}

// Test if a connection is found
func (self *SQL_Client) Ping(connectionId interface{}) (string, error) {

	// Here I will try to ping a non-existing connection.
	rqst := &sqlpb.PingConnectionRqst{
		Id: Utility.ToString(connectionId),
	}

	rsp, err := self.c.Ping(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.Result, err
}

// That function return the json string with all element in it.
func (self *SQL_Client) QueryContext(connectionId interface{}, query interface{}, parameters interface{}) (string, error) {

	parameters_ := strings.Split(parameters.(string), ",")
	parametersStr, _ := Utility.ToJson(parameters_)

	// The query and all it parameters.
	rqst := &sqlpb.QueryContextRqst{
		Query: &sqlpb.Query{
			ConnectionId: Utility.ToString(connectionId),
			Query:        Utility.ToString(query),
			Parameters:   parametersStr,
		},
	}

	// Because number of values can be high I will use a stream.
	stream, err := self.c.QueryContext(context.Background(), rqst)
	if err != nil {
		return "", err
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
			return "", err
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

	// Create object result and put header and data in it.
	result := make(map[string]interface{}, 0)
	result["header"] = header
	result["data"] = data
	resultStr, _ := json.Marshal(result)
	return string(resultStr), nil
}

func (self *SQL_Client) ExecContext(connectionId interface{}, query interface{}, parameters interface{}, tx interface{}) (string, error) {

	parameters_ := strings.Split(parameters.(string), ",")
	parametersStr, _ := Utility.ToJson(parameters_)

	rqst := &sqlpb.ExecContextRqst{
		Query: &sqlpb.Query{
			ConnectionId: Utility.ToString(connectionId),
			Query:        Utility.ToString(query),
			Parameters:   parametersStr,
		},
		Tx: Utility.ToBool(tx),
	}

	rsp, err := self.c.ExecContext(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	result := make(map[string]interface{}, 0)
	result["affectRows"] = rsp.AffectedRows
	result["lastId"] = rsp.LastId
	resultStr, _ := json.Marshal(result)

	return string(resultStr), nil
}

//** Create connection and Delete connection are left on the server side only...

////////////////////////////////////////////////////////////////////////////////
// LDAP Client Service
////////////////////////////////////////////////////////////////////////////////

type LDAP_Client struct {
	cc *grpc.ClientConn
	//c  ldappb
}

// Create a connection to the service.
func NewLdap_Client(addresse string) *LDAP_Client {
	client := new(LDAP_Client)
	client.cc = getClientConnection(addresse)
	//client.c = ldappb.(client.cc)

	return client
}

// must be close when no more needed.
func (self *LDAP_Client) Close() {
	self.cc.Close()
}

////////////////////////////////////////////////////////////////////////////////
// SMTP Client Service
////////////////////////////////////////////////////////////////////////////////
type SMTP_Client struct {
	cc *grpc.ClientConn
	c  smtppb.SmtpServiceClient
}

// Create a connection to the service.
func NewSmtp_Client(addresse string) *SMTP_Client {
	client := new(SMTP_Client)
	client.cc = getClientConnection(addresse)
	client.c = smtppb.NewSmtpServiceClient(client.cc)

	return client
}

// must be close when no more needed.
func (self *SMTP_Client) Close() {
	self.cc.Close()
}

////////////////////////////////////////////////////////////////////////////////
// Persitence Client Service
////////////////////////////////////////////////////////////////////////////////
type Persistence_Client struct {
	cc *grpc.ClientConn
	c  persistencepb.PersistenceServiceClient
}

// Create a connection to the service.
func NewPersistence_Client(addresse string) *Persistence_Client {
	client := new(Persistence_Client)
	client.cc = getClientConnection(addresse)
	client.c = persistencepb.NewPersistenceServiceClient(client.cc)

	return client
}

// must be close when no more needed.
func (self *Persistence_Client) Close() {
	self.cc.Close()
}

////////////////////////////////////////////////////////////////////////////////
// storage Client Service
////////////////////////////////////////////////////////////////////////////////

type Storage_Client struct {
	cc *grpc.ClientConn
	c  storagepb.StorageServiceClient
}

// Create a connection to the service.
func NewStorage_Client(addresse string) *Storage_Client {
	client := new(Storage_Client)
	client.cc = getClientConnection(addresse)
	client.c = storagepb.NewStorageServiceClient(client.cc)
	return client
}

// must be close when no more needed.
func (self *Storage_Client) Close() {
	self.cc.Close()
}

////////////////////////////////////////////////////////////////////////////////
// SPC Client Service
////////////////////////////////////////////////////////////////////////////////
type SPC_Client struct {
	cc *grpc.ClientConn
	c  spcpb.SpcServiceClient
}

// Create a connection to the service.
func NewSpc_Client(addresse string) *SPC_Client {
	client := new(SPC_Client)
	client.cc = getClientConnection(addresse)
	client.c = spcpb.NewSpcServiceClient(client.cc)

	return client
}

// must be close when no more needed.
func (self *SPC_Client) Close() {
	self.cc.Close()
}

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type Echo_Client struct {
	cc *grpc.ClientConn
	c  echopb.EchoServiceClient
}

// Create a connection to the service.
func NewEcho_Client(addresse string) *Echo_Client {
	client := new(Echo_Client)
	client.cc = getClientConnection(addresse)
	client.c = echopb.NewEchoServiceClient(client.cc)
	return client
}

// must be close when no more needed.
func (self *Echo_Client) Close() {
	self.cc.Close()
}

func (self *Echo_Client) Echo(msg interface{}) (string, error) {
	rqst := &echopb.EchoRequest{
		Message: Utility.ToString(msg),
	}

	rsp, err := self.c.Echo(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.Message, nil
}
