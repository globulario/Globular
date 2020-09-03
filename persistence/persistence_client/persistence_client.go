package persistence_client

import (
	"encoding/json"
	"io"

	"strconv"

	"context"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/persistence/persistencepb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// Persitence Client Service
////////////////////////////////////////////////////////////////////////////////
type Persistence_Client struct {
	cc *grpc.ClientConn
	c  persistencepb.PersistenceServiceClient

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
func NewPersistence_Client(address string, id string) (*Persistence_Client, error) {
	client := new(Persistence_Client)
	err := api.InitClient(client, address, id)
	if err != nil {
		return nil, err
	}
	client.cc, err = api.GetClientConnection(client)
	if err != nil {
		return nil, err
	}
	client.c = persistencepb.NewPersistenceServiceClient(client.cc)

	return client, nil
}

func (self *Persistence_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = api.GetClientContext(self)
	}
	return api.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the domain
func (self *Persistence_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *Persistence_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the id of the service instance
func (self *Persistence_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *Persistence_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Persistence_Client) Close() {
	if self.cc != nil {
		self.cc.Close()
	}
}

// Set grpc_service port.
func (self *Persistence_Client) SetPort(port int) {
	self.port = port
}

// Set the client instance id.
func (self *Persistence_Client) SetId(id string) {
	self.id = id
}

// Set the client name.
func (self *Persistence_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Persistence_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Persistence_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Persistence_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Persistence_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Persistence_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Persistence_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Persistence_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Persistence_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Persistence_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

///////////////////////// API /////////////////////

// Create a new datastore connection.
func (self *Persistence_Client) CreateConnection(connectionId string, name string, host string, port float64, storeType float64, user string, pwd string, timeout float64, options string, save bool) error {
	rqst := &persistencepb.CreateConnectionRqst{
		Connection: &persistencepb.Connection{
			Id:       connectionId,
			Name:     name,
			Host:     host,
			Port:     int32(Utility.ToInt(port)),
			Store:    persistencepb.StoreType(storeType),
			User:     user,
			Password: pwd,
			Timeout:  int32(Utility.ToInt(timeout)),
			Options:  options,
		},
		Save: save,
	}

	_, err := self.c.CreateConnection(api.GetClientContext(self), rqst)
	return err
}

func (self *Persistence_Client) DeleteConnection(connectionId string) error {
	rqst := &persistencepb.DeleteConnectionRqst{
		Id: connectionId,
	}

	_, err := self.c.DeleteConnection(api.GetClientContext(self), rqst)
	return err
}

func (self *Persistence_Client) Connect(id string, password string) error {
	rqst := &persistencepb.ConnectRqst{
		ConnectionId: id,
		Password:     password,
	}

	_, err := self.c.Connect(api.GetClientContext(self), rqst)
	return err
}

func (self *Persistence_Client) Disconnect(connectionId string) error {

	rqst := &persistencepb.DisconnectRqst{
		ConnectionId: connectionId,
	}

	_, err := self.c.Disconnect(api.GetClientContext(self), rqst)

	return err
}

func (self *Persistence_Client) Ping(connectionId string) error {

	rqst := &persistencepb.PingConnectionRqst{
		Id: connectionId,
	}

	_, err := self.c.Ping(api.GetClientContext(self), rqst)

	return err
}

func (self *Persistence_Client) FindOne(connectionId string, database string, collection string, jsonStr string, options string) (string, error) {

	// Retreive a single value...
	rqst := &persistencepb.FindOneRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Query:      jsonStr,
		Options:    options,
	}

	rsp, err := self.c.FindOne(api.GetClientContext(self), rqst)

	if err != nil {
		return "", err
	}

	return rsp.GetJsonStr(), err
}

func (self *Persistence_Client) Find(connectionId string, database string, collection string, query string, options string) (string, error) {

	// Retreive a single value...
	rqst := &persistencepb.FindRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Query:      query,
		Options:    options,
	}

	stream, err := self.c.Find(api.GetClientContext(self), rqst)

	if err != nil {
		return "", err
	}

	valuesStr := "["
	for {
		results, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// end of stream...
				break
			} else {
				return "", err
			}
		} else {
			if len(valuesStr) > 1 {
				valuesStr += ","
			}
			valuesStr += results.JsonStr[1 : len(results.JsonStr)-1]
		}

	}

	valuesStr += "]"

	return valuesStr, nil
}

/**
 * Usefull function to query and transform document.
 */
func (self *Persistence_Client) Aggregate(connectionId, database string, collection string, pipeline string, options string) (string, error) {
	// Retreive a single value...
	rqst := &persistencepb.AggregateRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Pipeline:   pipeline,
		Options:    options,
	}

	stream, err := self.c.Aggregate(api.GetClientContext(self), rqst)

	if err != nil {
		return "", err
	}

	valuesStr := "["
	for {
		results, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// end of stream...
				break
			} else {
				return "", err
			}
		} else {
			if len(valuesStr) > 1 {
				valuesStr += ","
			}
			valuesStr += results.JsonStr[1 : len(results.JsonStr)-1]
		}

	}

	valuesStr += "]"

	return string(valuesStr), nil
}

/**
 * Count the number of document that match the query.
 */
func (self *Persistence_Client) Count(connectionId string, database string, collection string, query string, options string) (int, error) {

	rqst := &persistencepb.CountRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Query:      query,
		Options:    options,
	}

	rsp, err := self.c.Count(api.GetClientContext(self), rqst)

	if err != nil {
		return 0, err
	}

	return int(rsp.Result), err
}

/**
 * Insert one value in the database.
 */
func (self *Persistence_Client) InsertOne(connectionId string, database string, collection string, jsonStr string, options string) (string, error) {

	rqst := &persistencepb.InsertOneRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		JsonStr:    jsonStr,
		Options:    options,
	}

	rsp, err := self.c.InsertOne(api.GetClientContext(self), rqst)

	if err != nil {
		return "", err
	}

	return rsp.GetId(), err
}

func (self *Persistence_Client) InsertMany(connectionId string, database string, collection string, jsonStr string, options string) (string, error) {

	stream, err := self.c.InsertMany(api.GetClientContext(self))
	if err != nil {
		return "", err
	}

	// here you must run the sql service test before runing this test in order
	// to generate the file Employees.json
	data := make([]map[string]interface{}, 0)

	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	// Persist 500 rows at time to save marshaling unmarshaling cycle time.
	chunkSize := 500

	for i := 0; i < len(data); i += chunkSize {

		var jsonStr []byte
		if i+chunkSize < len(data) {
			jsonStr, err = json.Marshal(data[i : i+chunkSize])
			if err != nil {
				return "", err
			}
		} else {
			jsonStr, err = json.Marshal(data[i:])
			if err != nil {
				return "", err
			}
		}

		rqst := &persistencepb.InsertManyRqst{
			Id:         connectionId,
			Database:   database,
			Collection: collection,
			JsonStr:    "[" + string(jsonStr) + "]",
		}

		err = stream.Send(rqst)
		if err != nil {
			return "", err
		}
	}

	rsp, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return rsp.Ids, nil
}

/**
 * Insert one value in the database.
 */
func (self *Persistence_Client) ReplaceOne(connectionId string, database string, collection string, query string, value string, options string) error {

	rqst := &persistencepb.ReplaceOneRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Query:      query,
		Value:      value,
		Options:    options,
	}

	_, err := self.c.ReplaceOne(api.GetClientContext(self), rqst)

	return err
}

func (self *Persistence_Client) UpdateOne(connectionId string, database string, collection string, query string, value string, options string) error {

	rqst := &persistencepb.UpdateOneRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Query:      query,
		Value:      value,
		Options:    options,
	}

	_, err := self.c.UpdateOne(api.GetClientContext(self), rqst)

	return err
}

/**
 * Update one or more document.
 */
func (self *Persistence_Client) Update(connectionId string, database string, collection string, query string, value string, options string) error {

	rqst := &persistencepb.UpdateRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Query:      query,
		Value:      value,
		Options:    options,
	}

	_, err := self.c.Update(api.GetClientContext(self), rqst)

	return err
}

/**
 * Delete one document from the db
 */
func (self *Persistence_Client) DeleteOne(connectionId string, database string, collection string, query string, options string) error {

	rqst := &persistencepb.DeleteOneRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Query:      query,
		Options:    options,
	}

	_, err := self.c.DeleteOne(api.GetClientContext(self), rqst)

	if err != nil {
		return err
	}

	return err
}

/**
 * Delete many document from the db.
 */
func (self *Persistence_Client) Delete(connectionId string, database string, collection string, query string, options string) error {

	rqst := &persistencepb.DeleteRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Query:      query,
		Options:    options,
	}

	_, err := self.c.Delete(api.GetClientContext(self), rqst)

	if err != nil {
		return err
	}

	return err
}

/**
 * Drop a collection.
 */
func (self *Persistence_Client) DeleteCollection(connectionId string, database string, collection string) error {
	// Test drop collection.
	rqst_drop_collection := &persistencepb.DeleteCollectionRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
	}
	_, err := self.c.DeleteCollection(api.GetClientContext(self), rqst_drop_collection)

	return err
}

/**
 * Drop a database.
 */
func (self *Persistence_Client) DeleteDatabase(connectionId string, database string) error {
	// Test drop collection.
	rqst_drop_db := &persistencepb.DeleteDatabaseRqst{
		Id:       connectionId,
		Database: database,
	}

	_, err := self.c.DeleteDatabase(api.GetClientContext(self), rqst_drop_db)

	return err
}

/**
 * Admin function, that must be protected.
 */
func (self *Persistence_Client) RunAdminCmd(connectionId string, user string, pwd string, script string) error {
	// Test drop collection.
	rqst_drop_db := &persistencepb.RunAdminCmdRqst{
		ConnectionId: connectionId,
		Script:       script,
		User:         user,
		Password:     pwd,
	}

	_, err := self.c.RunAdminCmd(api.GetClientContext(self), rqst_drop_db)

	return err
}
