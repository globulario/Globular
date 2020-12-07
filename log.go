package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/davecourtois/Utility"
	"github.com/globulario/Globular/Interceptors"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Start internal logging services.
func (self *Globule) startLogService() error {
	id := string(resourcepb.File_proto_resource_proto.Services().Get(2).FullName())
	log_server, err := self.startInternalService(id, resourcepb.File_proto_resource_proto.Path(), self.LogPort, self.LogProxy, self.Protocol == "https", self.unaryResourceInterceptor, self.streamResourceInterceptor)
	if err == nil {
		self.inernalServices = append(self.inernalServices, log_server)

		// Create the channel to listen on resource port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.LogPort))
		if err != nil {
			log.Fatalf("could not start resource service %s: %s", self.getDomain(), err)
		}

		resourcepb.RegisterLogServiceServer(log_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			// no web-rpc server.
			if err = log_server.Serve(lis); err != nil {
				log.Println(err)
			}
			s := self.getService(id)
			pid := getIntVal(s, "ProxyProcess")
			Utility.TerminateProcess(pid, 0)
			s.Store("ProxyProcess", -1)
			self.saveConfig()
		}()
	}

	return err
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
	if method == "/resource.LogService/Log" {
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
