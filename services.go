package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/services/golang/event/event_client"
	"github.com/davecourtois/Globular/services/golang/event/eventpb"
	"github.com/davecourtois/Globular/services/golang/services/service_client"
	"github.com/davecourtois/Globular/services/golang/services/servicespb"
	"github.com/davecourtois/Utility"
	"github.com/golang/protobuf/jsonpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
 * Subscribe to Discoverie's and repositories to keep services up to date.
 */
func (self *Globule) keepServicesUpToDate() map[string]map[string][]string {

	// append itself to service discoveries...
	subscribers := make(map[string]map[string][]string, 0)

	// Connect to service update events...
	for i := 0; i < len(self.Discoveries); i++ {
		log.Println("Connect to discovery event hub ", self.Discoveries[i])
		address := self.Discoveries[i]
		if address == self.Domain {
			address += ":" + Utility.ToString(self.PortHttp)
		}

		eventHub, err := event_client.NewEventService_Client(address, "event.EventService")
		if err == nil {
			log.Println("Connected with event service at ", self.Discoveries[i])
			if subscribers[self.Discoveries[i]] == nil {
				subscribers[self.Discoveries[i]] = make(map[string][]string)
			}

			for _, s := range self.getServices() {
				_, hasPublisherId := s.Load("PublisherId")
				if hasPublisherId {
					id := getStringVal(s, "PublisherId") + ":" + getStringVal(s, "Name") + ":SERVICE_PUBLISH_EVENT"

					if subscribers[self.Discoveries[i]][id] == nil {
						subscribers[self.Discoveries[i]][id] = make([]string, 0)
					}

					// each channel has it event...
					uuid := Utility.RandomUUID()
					fct := func(evt *eventpb.Event) {
						descriptor := new(servicespb.ServiceDescriptor)
						jsonpb.UnmarshalString(string(evt.GetData()), descriptor)

						// here I will update the service if it's version is lower
						for _, s := range self.getServices() {
							_, hasPublisherId := s.Load("PublisherId")
							if hasPublisherId {
								if getStringVal(s, "Name") == descriptor.GetId() && getStringVal(s, "PublisherId") == descriptor.GetPublisherId() {

									if getBoolVal(s, "KeepUpToDate") {
										// Test if update is needed...
										version := getStringVal(s, "Version")
										if Utility.ToInt(strings.Split(version, ".")[0]) <= Utility.ToInt(strings.Split(descriptor.Version, ".")[0]) {
											if Utility.ToInt(strings.Split(version, ".")[1]) <= Utility.ToInt(strings.Split(descriptor.Version, ".")[1]) {
												if Utility.ToInt(strings.Split(version, ".")[2]) < Utility.ToInt(strings.Split(descriptor.Version, ".")[2]) {
													self.stopService(getStringVal(s, "Id"))
													self.deleteService(getStringVal(s, "Id"))
													err := self.installService(descriptor)
													if err != nil {
														fmt.Println("fail to install service ", descriptor.GetPublisherId(), descriptor.GetId(), descriptor.GetVersion(), err)
													} else {
														s.Store("KeepUpToDate", true)
														self.saveConfig()
														fmt.Println("service was update!", descriptor.GetPublisherId(), descriptor.GetId(), descriptor.GetVersion())
													}
												}
											}
										}
									}

								}
							}
						}
					}

					// So here I will subscribe to service update event.
					// try 5 time wait 5 second before given up.
					registered := false
					for nbTry := 5; registered == false && nbTry > 0; nbTry-- {
						err := eventHub.Subscribe(id, uuid, fct)
						if err == nil {
							subscribers[self.Discoveries[i]][id] = append(subscribers[self.Discoveries[i]][id], uuid)
							log.Println("subscription to ", id, " succeed!")
							registered = true
						} else {
							log.Println("fail to subscribe to ", id)
							nbTry--
							time.Sleep(1 * time.Second)
						}
					}
				}
			}
			// keep on memorie...
			self.discorveriesEventHub[self.Discoveries[i]] = eventHub
		}

	}
	return subscribers

}

// Start service discovery
func (self *Globule) startDiscoveryService() error {
	// The service discovery.
	id := string(servicespb.File_services_proto_services_proto.Services().Get(0).FullName())
	services_discovery_server, err := self.startInternalService(id, servicespb.File_services_proto_services_proto.Path(), self.ServicesDiscoveryPort, self.ServicesDiscoveryProxy, self.Protocol == "https", Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor)
	if err == nil {
		self.inernalServices = append(self.inernalServices, services_discovery_server)
		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.ServicesDiscoveryPort))
		if err != nil {
			log.Fatalf("could not start services discovery service %s: %s", self.getDomain(), err)
		}

		servicespb.RegisterServiceDiscoveryServer(services_discovery_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.

		go func() {
			// no web-rpc server.
			if err := services_discovery_server.Serve(lis); err != nil {
				log.Println(err)
			}
			s := self.getService(id)
			pid := getIntVal(s, "ProxyProcess")
			Utility.TerminateProcess(pid, 0)
			s.Store("ProxyProcess", -1)
			self.saveConfig()
			return
		}()
	}
	return err

}

// Start service repository
func (self *Globule) startRepositoryService() error {
	services_repository_server, err := self.startInternalService(string(servicespb.File_services_proto_services_proto.Services().Get(1).FullName()), servicespb.File_services_proto_services_proto.Path(), self.ServicesRepositoryPort, self.ServicesRepositoryProxy,
		self.Protocol == "https",
		Interceptors.ServerUnaryInterceptor,
		Interceptors.ServerStreamInterceptor)

	if err == nil {
		self.inernalServices = append(self.inernalServices, services_repository_server)

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.ServicesRepositoryPort))
		if err != nil {
			log.Fatalf("could not start services repository service %s: %s", self.getDomain(), err)
		}

		servicespb.RegisterServiceRepositoryServer(services_repository_server, self)

		go func() {
			// no web-rpc server.
			if err := services_repository_server.Serve(lis); err != nil {
				log.Println(err)

			}
		}()
	}
	return err

}

//////////////////////////////// Services management  //////////////////////////

// TODO synchronize it!
/**
 * Return the list of service configuaration with a given name.
 **/
func (self *Globule) getServiceConfigByName(name string) []map[string]interface{} {
	configs := make([]map[string]interface{}, 0)

	for _, config := range self.getConfig()["Services"].(map[string]interface{}) {
		if config.(map[string]interface{})["Name"].(string) == name {
			configs = append(configs, config.(map[string]interface{}))
		}
	}

	return configs
}

// Discovery
func (self *Globule) FindServices(ctx context.Context, rqst *servicespb.FindServicesDescriptorRequest) (*servicespb.FindServicesDescriptorResponse, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	kewordsStr, err := Utility.ToJson(rqst.Keywords)
	if err != nil {
		return nil, err
	}

	// Test...
	query := `{"keywords": { "$all" : ` + kewordsStr + `}}`

	data, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Services", query, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	descriptors := make([]*servicespb.ServiceDescriptor, len(data))
	for i := 0; i < len(data); i++ {
		descriptor := data[i].(map[string]interface{})
		descriptors[i] = new(servicespb.ServiceDescriptor)
		descriptors[i].Id = descriptor["id"].(string)
		descriptors[i].Name = descriptor["name"].(string)
		descriptors[i].Description = descriptor["description"].(string)
		descriptors[i].PublisherId = descriptor["publisherid"].(string)
		descriptors[i].Version = descriptor["version"].(string)
		if descriptor["keywords"] != nil {
			descriptor["keywords"] = []interface{}(descriptor["keywords"].(primitive.A))
			descriptors[i].Keywords = make([]string, len(descriptor["keywords"].([]interface{})))
			for j := 0; j < len(descriptor["keywords"].([]interface{})); j++ {
				descriptors[i].Keywords[j] = descriptor["keywords"].([]interface{})[j].(string)
			}
		}
		if descriptor["discoveries"] != nil {
			descriptor["discoveries"] = []interface{}(descriptor["discoveries"].(primitive.A))
			descriptors[i].Discoveries = make([]string, len(descriptor["discoveries"].([]interface{})))
			for j := 0; j < len(descriptor["discoveries"].([]interface{})); j++ {
				descriptors[i].Discoveries[j] = descriptor["discoveries"].([]interface{})[j].(string)
			}
		}

		if descriptor["repositories"] != nil {
			descriptor["repositories"] = []interface{}(descriptor["repositories"].(primitive.A))
			descriptors[i].Repositories = make([]string, len(descriptor["repositories"].([]interface{})))
			for j := 0; j < len(descriptor["repositories"].([]interface{})); j++ {
				descriptors[i].Repositories[j] = descriptor["repositories"].([]interface{})[j].(string)
			}
		}
	}

	// Return the list of Service Descriptor.
	return &servicespb.FindServicesDescriptorResponse{
		Results: descriptors,
	}, nil
}

//* Return the list of all services *
func (self *Globule) GetServiceDescriptor(ctx context.Context, rqst *servicespb.GetServiceDescriptorRequest) (*servicespb.GetServiceDescriptorResponse, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	query := `{"id":"` + rqst.ServiceId + `", "publisherid":"` + rqst.PublisherId + `"}`

	values, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Services", query, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if len(values) == 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No service descriptor with id "+rqst.ServiceId+" was found for publisher id "+rqst.PublisherId)))
	}

	descriptors := make([]*servicespb.ServiceDescriptor, len(values))
	for i := 0; i < len(values); i++ {

		descriptor := values[i].(map[string]interface{})
		descriptors[i] = new(servicespb.ServiceDescriptor)
		descriptors[i].Id = descriptor["id"].(string)
		descriptors[i].Name = descriptor["name"].(string)
		descriptors[i].Description = descriptor["description"].(string)
		descriptors[i].PublisherId = descriptor["publisherid"].(string)
		descriptors[i].Version = descriptor["version"].(string)
		if descriptor["keywords"] != nil {
			descriptor["keywords"] = []interface{}(descriptor["keywords"].(primitive.A))
			descriptors[i].Keywords = make([]string, len(descriptor["keywords"].([]interface{})))
			for j := 0; j < len(descriptor["keywords"].([]interface{})); j++ {
				descriptors[i].Keywords[j] = descriptor["keywords"].([]interface{})[j].(string)
			}
		}
		if descriptor["discoveries"] != nil {
			descriptor["discoveries"] = []interface{}(descriptor["discoveries"].(primitive.A))
			descriptors[i].Discoveries = make([]string, len(descriptor["discoveries"].([]interface{})))
			for j := 0; j < len(descriptor["discoveries"].([]interface{})); j++ {
				descriptors[i].Discoveries[j] = descriptor["discoveries"].([]interface{})[j].(string)
			}
		}

		if descriptor["repositories"] != nil {
			descriptor["repositories"] = []interface{}(descriptor["repositories"].(primitive.A))
			descriptors[i].Repositories = make([]string, len(descriptor["repositories"].([]interface{})))
			for j := 0; j < len(descriptor["repositories"].([]interface{})); j++ {
				descriptors[i].Repositories[j] = descriptor["repositories"].([]interface{})[j].(string)
			}
		}

	}

	sort.Slice(descriptors[:], func(i, j int) bool {
		return descriptors[i].Version > descriptors[j].Version
	})

	// Return the list of Service Descriptor.
	return &servicespb.GetServiceDescriptorResponse{
		Results: descriptors,
	}, nil
}

//* Return the list of all services *
func (self *Globule) GetServicesDescriptor(rqst *servicespb.GetServicesDescriptorRequest, stream servicespb.ServiceDiscovery_GetServicesDescriptorServer) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	data, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Services", `{}`, "")
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	descriptors := make([]*servicespb.ServiceDescriptor, 0)
	for i := 0; i < len(data); i++ {
		descriptor := new(servicespb.ServiceDescriptor)

		descriptor.Id = data[i].(map[string]interface{})["id"].(string)
		descriptor.Name = data[i].(map[string]interface{})["name"].(string)
		descriptor.Description = data[i].(map[string]interface{})["description"].(string)
		descriptor.PublisherId = data[i].(map[string]interface{})["publisherid"].(string)
		descriptor.Version = data[i].(map[string]interface{})["version"].(string)
		if data[i].(map[string]interface{})["keywords"] != nil {
			data[i].(map[string]interface{})["keywords"] = []interface{}(data[i].(map[string]interface{})["keywords"].(primitive.A))
			descriptor.Keywords = make([]string, len(data[i].(map[string]interface{})["keywords"].([]interface{})))
			for j := 0; j < len(data[i].(map[string]interface{})["keywords"].([]interface{})); j++ {
				descriptor.Keywords[j] = data[i].(map[string]interface{})["keywords"].([]interface{})[j].(string)
			}
		}
		if data[i].(map[string]interface{})["discoveries"] != nil {
			data[i].(map[string]interface{})["discoveries"] = []interface{}(data[i].(map[string]interface{})["discoveries"].(primitive.A))
			descriptor.Discoveries = make([]string, len(data[i].(map[string]interface{})["discoveries"].([]interface{})))
			for j := 0; j < len(data[i].(map[string]interface{})["discoveries"].([]interface{})); j++ {
				descriptor.Discoveries[j] = data[i].(map[string]interface{})["discoveries"].([]interface{})[j].(string)
			}
		}

		if data[i].(map[string]interface{})["repositories"] != nil {
			data[i].(map[string]interface{})["repositories"] = []interface{}(data[i].(map[string]interface{})["repositories"].(primitive.A))
			descriptor.Repositories = make([]string, len(data[i].(map[string]interface{})["repositories"].([]interface{})))
			for j := 0; j < len(data[i].(map[string]interface{})["repositories"].([]interface{})); j++ {
				descriptor.Repositories[j] = data[i].(map[string]interface{})["repositories"].([]interface{})[j].(string)
			}
		}

		descriptors = append(descriptors, descriptor)
		// send at each 20
		if i%20 == 0 {
			stream.Send(&servicespb.GetServicesDescriptorResponse{
				Results: descriptors,
			})
			descriptors = make([]*servicespb.ServiceDescriptor, 0)
		}
	}

	if len(descriptors) > 0 {
		stream.Send(&servicespb.GetServicesDescriptorResponse{
			Results: descriptors,
		})
	}

	// Return the list of Service Descriptor.
	return nil
}

/**
 */
func (self *Globule) SetServiceDescriptor(ctx context.Context, rqst *servicespb.SetServiceDescriptorRequest) (*servicespb.SetServiceDescriptorResponse, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	var marshaler jsonpb.Marshaler

	jsonStr, err := marshaler.MarshalToString(rqst.Descriptor_)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// little fix...
	jsonStr = strings.ReplaceAll(jsonStr, "publisherId", "publisherid")

	// Always create a new if not already exist.
	err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Services", `{"id":"`+rqst.Descriptor_.Id+`", "publisherid":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &servicespb.SetServiceDescriptorResponse{
		Result: true,
	}, nil
}

//* Publish a service to service discovery *
func (self *Globule) PublishServiceDescriptor(ctx context.Context, rqst *servicespb.PublishServiceDescriptorRequest) (*servicespb.PublishServiceDescriptorResponse, error) {

	// Here I will save the descriptor inside the storage...
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Append the self domain to the list of discoveries where the services can be found.
	if !Utility.Contains(rqst.Descriptor_.Discoveries, self.getDomain()) {
		rqst.Descriptor_.Discoveries = append(rqst.Descriptor_.Discoveries, self.getDomain())
	}

	// Here I will test if the services already exist...
	_, err = p.FindOne(context.Background(), "local_ressource", "local_ressource", "Services", `{"id":"`+rqst.Descriptor_.Id+`", "publisherid":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, "")
	if err == nil {
		// Update existing descriptor.

		// The list of discoveries...
		discoveries, err := Utility.ToJson(rqst.Descriptor_.Discoveries)
		if err == nil {
			values := `{"$set":{"discoveries":` + discoveries + `}}`
			err = p.Update(context.Background(), "local_ressource", "local_ressource", "Services", `{"id":"`+rqst.Descriptor_.Id+`", "publisherid":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, values, "")
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		}

		// The list of repositories
		repositories, err := Utility.ToJson(rqst.Descriptor_.Repositories)
		if err == nil {
			values := `{"$set":{"repositories":` + repositories + `}}`
			err = p.Update(context.Background(), "local_ressource", "local_ressource", "Services", `{"id":"`+rqst.Descriptor_.Id+`", "publisherid":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, values, "")
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		}

	} else {

		// The key will be the descriptor string itself.
		_, err = p.InsertOne(context.Background(), "local_ressource", "local_ressource", "Services", rqst.Descriptor_, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

	}

	return &servicespb.PublishServiceDescriptorResponse{
		Result: true,
	}, nil
}

// Repository
/** Download a service from a service directory **/
func (self *Globule) DownloadBundle(rqst *servicespb.DownloadBundleRequest, stream servicespb.ServiceRepository_DownloadBundleServer) error {
	bundle := new(servicespb.ServiceBundle)
	bundle.Plaform = rqst.Plaform
	bundle.Descriptor_ = rqst.Descriptor_

	// Generate the bundle id....
	var id string
	id = bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Name + "%" + bundle.Descriptor_.Version + "%" + bundle.Descriptor_.Id + "%" + rqst.Plaform

	path := self.data + string(os.PathSeparator) + "service-repository"

	var err error
	// the file must be a zipped archive that contain a .proto, .config and executable.
	bundle.Binairies, err = ioutil.ReadFile(path + string(os.PathSeparator) + id + ".tar.gz")
	if err != nil {
		return err
	}

	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "ServiceBundle", `{"_id":"`+id+`"}`, "")
	if err != nil {
		return err
	}

	// init the map with json values.
	checksum := values.(map[string]interface{})

	// Test if the values change over time.
	if Utility.CreateDataChecksum(bundle.Binairies) != checksum["checksum"].(string) {
		return errors.New("The bundle data cheksum is not valid!")
	}

	const BufferSize = 1024 * 5 // the chunck size.
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer) // Will write to network.
	err = enc.Encode(bundle)
	if err != nil {
		return err
	}

	for {
		var data [BufferSize]byte
		bytesread, err := buffer.Read(data[0:BufferSize])
		if bytesread > 0 {
			rqst := &servicespb.DownloadBundleResponse{
				Data: data[0:bytesread],
			}
			// send the data to the server.
			err = stream.Send(rqst)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}

/** Upload a service to a service directory **/
func (self *Globule) UploadBundle(stream servicespb.ServiceRepository_UploadBundleServer) error {

	// The bundle will cantain the necessary information to install the service.
	var buffer bytes.Buffer
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			stream.SendAndClose(&servicespb.UploadBundleResponse{
				Result: true,
			})
			break
		} else if err != nil {
			return err
		} else {
			buffer.Write(msg.Data)
		}
	}

	// The buffer that contain the
	dec := gob.NewDecoder(&buffer)
	bundle := new(servicespb.ServiceBundle)
	err := dec.Decode(bundle)
	if err != nil {
		return err
	}

	// Generate the bundle id....
	id := bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Name + "%" + bundle.Descriptor_.Version + "%" + bundle.Descriptor_.Id + "%" + bundle.Plaform
	log.Println(id)

	repositoryId := self.getDomain()
	// Now I will append the address of the repository into the service descriptor.
	if !Utility.Contains(bundle.Descriptor_.Repositories, repositoryId) {
		bundle.Descriptor_.Repositories = append(bundle.Descriptor_.Repositories, repositoryId)
		// Publish change into discoveries...
		for i := 0; i < len(bundle.Descriptor_.Discoveries); i++ {
			discoveryId := bundle.Descriptor_.Discoveries[i]
			discoveryService, err := service_client.NewServicesDiscoveryService_Client(discoveryId, "services.ServiceDiscovery")
			if err != nil {
				return err
			}
			discoveryService.PublishServiceDescriptor(bundle.Descriptor_)
		}
	}

	path := self.data + string(os.PathSeparator) + "service-repository"
	Utility.CreateDirIfNotExist(path)

	// the file must be a zipped archive that contain a .proto, .config and executable.
	err = ioutil.WriteFile(path+"/"+id+".tar.gz", bundle.Binairies, 777)
	if err != nil {
		return err
	}

	checksum := Utility.CreateDataChecksum(bundle.Binairies)
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	_, err = p.InsertOne(context.Background(), "local_ressource", "local_ressource", "ServiceBundle", map[string]interface{}{"_id": id, "checksum": checksum, "platform": bundle.Plaform, "publisherid": bundle.Descriptor_.PublisherId, "servicename": bundle.Descriptor_.Name, "serviceid": bundle.Descriptor_.Id, "modified": time.Now().Unix(), "size": len(bundle.Binairies)}, "")

	return err
}
