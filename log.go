package main

import (
	"time"

	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/davecourtois/Utility"
	"github.com/globulario/Globular/Interceptors"
	globular "github.com/globulario/services/golang/globular_client"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Log_Client struct {
	cc *grpc.ClientConn
	c  resourcepb.LogServiceClient

	// The id of the service
	id string

	// The name of the service
	name string

	// The client domain
	domain string

	// The port number
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
func NewResourceService_Client(address string, id string) (*Log_Client, error) {

	client := new(Log_Client)
	err := globular.InitClient(client, address, id)
	if err != nil {

		return nil, err
	}

	client.cc, err = globular.GetClientConnection(client)
	if err != nil {

		return nil, err
	}

	client.c = resourcepb.NewLogServiceClient(client.cc)

	return client, nil
}

func (self *Log_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the ipv4 address
// Return the address
func (self *Log_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the domain
func (self *Log_Client) GetDomain() string {
	return self.domain
}

// Return the id of the service instance
func (self *Log_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *Log_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Log_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *Log_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *Log_Client) SetId(id string) {
	self.id = id
}

// Set the client name.
func (self *Log_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Log_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Log_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Log_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Log_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Log_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Log_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Log_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Log_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Log_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

////////////////////////////////////////////////////////////////////////////////
// Api
////////////////////////////////////////////////////////////////////////////////

func (self *Globule) logServiceInfo(service string, message string) error {

	// Here I will use event to publish log information...
	info := new(resourcepb.LogInfo)
	info.Application = ""
	info.UserId = "globular"
	info.UserName = "globular"
	info.Method = service
	info.Date = time.Now().Unix()
	info.Message = message
	info.Type = resourcepb.LogType_ERROR_MESSAGE // not necessarely errors..
	self.log(info)

	return nil
}

// Log err and info...
func (self *Globule) logInfo(application string, method string, token string, err_ error) error {

	// Remove cyclic calls
	if method == "/resource.ResourceService/Log" {
		return errors.New("Method " + method + " cannot not be log because it will cause a circular call to itself!")
	}

	// Here I will use event to publish log information...
	info := new(resourcepb.LogInfo)
	info.Application = application
	info.UserId = token
	info.UserName = token
	info.Method = method
	info.Date = time.Now().Unix()
	if err_ != nil {
		info.Message = err_.Error()
		info.Type = resourcepb.LogType_ERROR_MESSAGE
		logger.Error(info.Message)
	} else {
		info.Type = resourcepb.LogType_INFO_MESSAGE
		logger.Info(info.Message)
	}

	self.log(info)

	return nil
}

func (self *Globule) getLogInfoKeyValue(info *resourcepb.LogInfo) (string, string, error) {
	marshaler := new(jsonpb.Marshaler)
	jsonStr, err := marshaler.MarshalToString(info)
	if err != nil {
		return "", "", err
	}

	key := ""
	if info.GetType() == resourcepb.LogType_INFO_MESSAGE {

		// Append the log in leveldb
		key += "/infos/" + info.Method + Utility.ToString(info.Date)

		// Set the application in the path
		if len(info.Application) > 0 {
			key += "/" + info.Application
		}
		// Set the User Name if available.
		if len(info.UserName) > 0 {
			key += "/" + info.UserName
		}

		key += "/" + Utility.GenerateUUID(jsonStr)

	} else {

		key += "/errors/" + info.Method + Utility.ToString(info.Date)

		// Set the application in the path
		if len(info.Application) > 0 {
			key += "/" + info.Application
		}
		// Set the User Name if available.
		if len(info.UserName) > 0 {
			key += "/" + info.UserName
		}

		key += "/" + Utility.GenerateUUID(jsonStr)

	}
	return key, jsonStr, nil
}

func (self *Globule) log(info *resourcepb.LogInfo) error {

	// The userId can be a single string or a JWT token.
	if len(info.UserName) > 0 {
		name, _, _, err := Interceptors.ValidateToken(info.UserName)
		if err == nil {
			info.UserName = name
		}
		info.UserId = info.UserName // keep only the user name
		if info.UserName == "sa" {
			return nil // not log sa activities...
		}
	} else {
		return nil
	}

	key, jsonStr, err := self.getLogInfoKeyValue(info)
	if err != nil {
		return err
	}

	// Append the error in leveldb
	self.logs.SetItem(key, []byte(jsonStr))
	eventHub, err := self.getEventHub()
	if err != nil {
		return err
	}

	eventHub.Publish(info.Method, []byte(jsonStr))

	return nil
}

// Log error or information into the data base *
func (self *Globule) Log(ctx context.Context, rqst *resourcepb.LogRqst) (*resourcepb.LogRsp, error) {
	// Publish event...
	self.log(rqst.Info)

	return &resourcepb.LogRsp{
		Result: true,
	}, nil
}

// Log error or information into the data base *
// Retreive log infos (the query must be something like /infos/'date'/'applicationName'/'userName'
func (self *Globule) GetLog(rqst *resourcepb.GetLogRqst, stream resourcepb.LogService_GetLogServer) error {

	query := rqst.Query
	if len(query) == 0 {
		query = "/*"
	}

	data, err := self.logs.GetItem(query)

	jsonDecoder := json.NewDecoder(strings.NewReader(string(data)))

	// read open bracket
	_, err = jsonDecoder.Token()
	if err != nil {
		return err
	}

	infos := make([]*resourcepb.LogInfo, 0)
	i := 0
	max := 100
	for jsonDecoder.More() {
		info := resourcepb.LogInfo{}
		err := jsonpb.UnmarshalNext(jsonDecoder, &info)
		if err != nil {
			return err
		}
		// append the info inside the stream.
		infos = append(infos, &info)
		if i == max-1 {
			// I will send the stream at each 100 logs...
			rsp := &resourcepb.GetLogRsp{
				Info: infos,
			}
			// Send the infos
			err = stream.Send(rsp)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			infos = make([]*resourcepb.LogInfo, 0)
			i = 0
		}
		i++
	}

	// Send the last infos...
	if len(infos) > 0 {
		rsp := &resourcepb.GetLogRsp{
			Info: infos,
		}
		err = stream.Send(rsp)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return nil
}

func (self *Globule) deleteLog(query string) error {

	// First of all I will retreive the log info with a given date.
	data, err := self.logs.GetItem(query)

	jsonDecoder := json.NewDecoder(strings.NewReader(string(data)))

	// read open bracket
	_, err = jsonDecoder.Token()
	if err != nil {
		return err
	}

	for jsonDecoder.More() {
		info := resourcepb.LogInfo{}

		err := jsonpb.UnmarshalNext(jsonDecoder, &info)
		if err != nil {
			return err
		}

		key, _, err := self.getLogInfoKeyValue(&info)
		if err != nil {
			return err
		}
		self.logs.RemoveItem(key)

	}

	return nil
}

//* Delete a log info *
func (self *Globule) DeleteLog(ctx context.Context, rqst *resourcepb.DeleteLogRqst) (*resourcepb.DeleteLogRsp, error) {

	key, _, _ := self.getLogInfoKeyValue(rqst.Log)
	err := self.logs.RemoveItem(key)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteLogRsp{
		Result: true,
	}, nil
}

//* Clear logs. info or errors *
func (self *Globule) ClearAllLog(ctx context.Context, rqst *resourcepb.ClearAllLogRqst) (*resourcepb.ClearAllLogRsp, error) {
	var err error

	if rqst.Type == resourcepb.LogType_ERROR_MESSAGE {
		err = self.deleteLog("/errors/*")
	} else {
		err = self.deleteLog("/infos/*")
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.ClearAllLogRsp{
		Result: true,
	}, nil
}
