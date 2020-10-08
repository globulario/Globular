package search_client

import (
	"strconv"

	"context"

	globular "github.com/davecourtois/Globular/services/golang/globular_client"
	"github.com/davecourtois/Globular/services/golang/search/searchpb"

	//"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type Search_Client struct {
	cc *grpc.ClientConn
	c  searchpb.SearchServiceClient

	// The id of the service
	id string

	// The name of the service
	name string

	// The client domain
	domain string

	// The port
	port int

	// is the connection is secure?
	hasTLS bool

	// Link to client key file
	keyFile string

	// Link to client certificate file.
	certFile string

	// certificate authority file
	caFile string
}

// Create a connection to the service.
func NewSearchService_Client(address string, id string) (*Search_Client, error) {
	client := new(Search_Client)
	err := globular.InitClient(client, address, id)
	if err != nil {
		return nil, err
	}
	client.cc, err = globular.GetClientConnection(client)
	if err != nil {
		return nil, err
	}
	client.c = searchpb.NewSearchServiceClient(client.cc)

	return client, nil
}

func (self *Search_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the domain
func (self *Search_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *Search_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the id of the service instance
func (self *Search_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *Search_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Search_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *Search_Client) SetPort(port int) {
	self.port = port
}

// Set the client service id.
func (self *Search_Client) SetId(id string) {
	self.id = id
}

// Set the client name.
func (self *Search_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Search_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Search_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Search_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Search_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Search_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Search_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Search_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Search_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Search_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

////////////////// Api //////////////////////
// Stop the service.
func (self *Search_Client) StopService() {
	self.c.Stop(globular.GetClientContext(self), &searchpb.StopRequest{})
}

/**
 * Return the version of the underlying search engine.
 */
func (self *Search_Client) GetVersion() (string, error) {

	rqst := &searchpb.GetVersionRequest{}
	ctx := globular.GetClientContext(self)
	rsp, err := self.c.GetVersion(ctx, rqst)
	if err != nil {
		return "", err
	}
	return rsp.Message, nil
}

/**
 * Index a JSON object / array
 */
func (self *Search_Client) IndexJsonObject(path string, jsonStr string, language string, id string, indexs []string, data string) error {
	rqst := &searchpb.IndexJsonObjectRequest{
		JsonStr:  jsonStr,
		Language: language,
		Id:       id,
		Indexs:   indexs,
		Data:     data,
		Path:     path,
	}
	ctx := globular.GetClientContext(self)
	_, err := self.c.IndexJsonObject(ctx, rqst)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Index a text file.
 * -dbPath the database path.
 * -filePath the file path must be reachable from the server.
 */
func (self *Search_Client) IndexFile(dbPath string, filePath string, language string) error {
	rqst := &searchpb.IndexFileRequest{
		DbPath:   dbPath,
		FilePath: filePath,
		Language: language,
	}

	ctx := globular.GetClientContext(self)
	_, err := self.c.IndexFile(ctx, rqst)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Index a text file.
 * -dbPath the database path.
 * -dirPath the file path must be reachable from the server.
 */
func (self *Search_Client) IndexDir(dbPath string, dirPath string, language string) error {
	rqst := &searchpb.IndexDirRequest{
		DbPath:   dbPath,
		DirPath:  dirPath,
		Language: language,
	}

	ctx := globular.GetClientContext(self)
	_, err := self.c.IndexDir(ctx, rqst)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 	Execute a search over the db.
 *  -path The path of the db
 *  -query The query string
 *  -language The language of the db
 *  -fields The list of fields
 *  -offset The results offset
 *  -pageSize The number of result to be return.
 *  -snippetLength The length of the snippet.
 */
func (self *Search_Client) SearchDocuments(paths []string, query string, language string, fields []string, offset int32, pageSize int32, snippetLength int32) ([]*searchpb.SearchResult, error) {
	rqst := &searchpb.SearchDocumentsRequest{
		Paths:         paths,
		Query:         query,
		Language:      language,
		Fields:        fields,
		Offset:        offset,
		PageSize:      pageSize,
		SnippetLength: snippetLength,
	}

	ctx := globular.GetClientContext(self)
	rsp, err := self.c.SearchDocuments(ctx, rqst)
	if err != nil {
		return nil, err
	}

	return rsp.GetResults(), nil

}

/**
 * Count the number of document in a given database.
 */
func (self *Search_Client) Count(path string) (int32, error) {
	rqst := &searchpb.CountRequest{
		Path: path,
	}

	ctx := globular.GetClientContext(self)
	rsp, err := self.c.Count(ctx, rqst)

	if err != nil {
		return -1, err
	}

	return rsp.Result, nil
}

/**
 * Delete a docuement from the database.
 */
func (self *Search_Client) DeleteDocument(path string, id string) error {
	rqst := &searchpb.DeleteDocumentRequest{
		Path: path,
		Id:   id,
	}

	ctx := globular.GetClientContext(self)
	_, err := self.c.DeleteDocument(ctx, rqst)

	if err != nil {
		return err
	}

	return nil
}
