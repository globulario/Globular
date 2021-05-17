package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/davecourtois/Utility"
	"github.com/globulario/services/golang/interceptors"
	"github.com/globulario/services/golang/packages/packages_client"
	"github.com/globulario/services/golang/packages/packagespb"
	"github.com/golang/protobuf/jsonpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
 * Subscribe to Discoverie's and repositories to keep services up to date.
 */
func (globule *Globule) keepServicesToDate() {

	ticker := time.NewTicker(time.Duration(globule.WatchUpdateDelay) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				// Connect to service update events...
				for _, s := range globule.getServices() {
					if getIntVal(s, "Process") != -1 {
						// TODO implement the api where to validate the the actual service configuration and
						// get the new one as necessary.
						log.Println("Keep service up to date", getStringVal(s, "Name"), getStringVal(s, "Id"), getStringVal(s, "Version"))
					}
				}
			}
		}
	}()
}

/**
 *
 */
func (globule *Globule) startPackagesDiscoveryService() error {
	// The service discovery.
	id := string(packagespb.File_proto_packages_proto.Services().Get(0).FullName())
	services_discovery_server, port, err := globule.startInternalService(id, packagespb.File_proto_packages_proto.Path(), globule.Protocol == "https", interceptors.ServerUnaryInterceptor, interceptors.ServerStreamInterceptor)
	if err == nil {
		globule.inernalServices = append(globule.inernalServices, services_discovery_server)
		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
		if err != nil {
			log.Fatalf("could not start services discovery service %s: %s", globule.getDomain(), err)
		}

		packagespb.RegisterPackageDiscoveryServer(services_discovery_server, globule)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			// no web-rpc server.
			if err := services_discovery_server.Serve(lis); err != nil {
				log.Println(err)
			}
			s := globule.getService(id)
			pid := getIntVal(s, "ProxyProcess")
			Utility.TerminateProcess(pid, 0)
			s.Store("ProxyProcess", -1)
			globule.setService(s)
		}()
	}
	return err

}

// Start service repository
func (globule *Globule) startPackagesRepositoryService() error {
	id := string(packagespb.File_proto_packages_proto.Services().Get(1).FullName())

	services_repository_server, port, err := globule.startInternalService(id, packagespb.File_proto_packages_proto.Path(),
		globule.Protocol == "https",
		interceptors.ServerUnaryInterceptor,
		interceptors.ServerStreamInterceptor)

	if err == nil {
		globule.inernalServices = append(globule.inernalServices, services_repository_server)

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
		if err != nil {
			log.Fatalf("could not start services repository service %s: %s", globule.getDomain(), err)
		}

		packagespb.RegisterPackageRepositoryServer(services_repository_server, globule)

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
func (globule *Globule) getServiceConfigByName(name string) []map[string]interface{} {
	configs := make([]map[string]interface{}, 0)

	for _, config := range globule.getConfig()["Services"].(map[string]interface{}) {
		if config.(map[string]interface{})["Name"].(string) == name {
			configs = append(configs, config.(map[string]interface{}))
		}
	}

	return configs
}

// Discovery
func (globule *Globule) FindPackages(ctx context.Context, rqst *packagespb.FindPackagesDescriptorRequest) (*packagespb.FindPackagesDescriptorResponse, error) {
	// That service made user of persistence service.
	p, err := globule.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	kewordsStr, err := Utility.ToJson(rqst.Keywords)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Test...
	query := `{"keywords": { "$all" : ` + kewordsStr + `}}`

	data, err := p.Find(context.Background(), "local_resource", "local_resource", "Packages", query, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	descriptors := make([]*packagespb.PackageDescriptor, len(data))
	for i := 0; i < len(data); i++ {
		descriptor := data[i].(map[string]interface{})
		descriptors[i] = new(packagespb.PackageDescriptor)
		descriptors[i].Id = descriptor["id"].(string)
		descriptors[i].Name = descriptor["name"].(string)
		descriptors[i].Description = descriptor["description"].(string)
		descriptors[i].PublisherId = descriptor["publisherid"].(string)
		descriptors[i].Version = descriptor["version"].(string)
		descriptors[i].Icon = descriptor["icon"].(string)
		descriptors[i].Alias = descriptor["alias"].(string)
		if descriptor["keywords"] != nil {
			descriptor["keywords"] = []interface{}(descriptor["keywords"].(primitive.A))
			descriptors[i].Keywords = make([]string, len(descriptor["keywords"].([]interface{})))
			for j := 0; j < len(descriptor["keywords"].([]interface{})); j++ {
				descriptors[i].Keywords[j] = descriptor["keywords"].([]interface{})[j].(string)
			}
		}
		if descriptor["actions"] != nil {
			descriptor["actions"] = []interface{}(descriptor["actions"].(primitive.A))
			descriptors[i].Actions = make([]string, len(descriptor["actions"].([]interface{})))
			for j := 0; j < len(descriptor["actions"].([]interface{})); j++ {
				descriptors[i].Actions[j] = descriptor["actions"].([]interface{})[j].(string)
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
	return &packagespb.FindPackagesDescriptorResponse{
		Results: descriptors,
	}, nil
}

//* Return the list of all services *
func (globule *Globule) GetPackageDescriptor(ctx context.Context, rqst *packagespb.GetPackageDescriptorRequest) (*packagespb.GetPackageDescriptorResponse, error) {
	p, err := globule.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	query := `{"id":"` + rqst.ServiceId + `", "publisherid":"` + rqst.PublisherId + `"}`

	values, err := p.Find(context.Background(), "local_resource", "local_resource", "Packages", query, "")
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

	descriptors := make([]*packagespb.PackageDescriptor, len(values))
	for i := 0; i < len(values); i++ {

		descriptor := values[i].(map[string]interface{})
		descriptors[i] = new(packagespb.PackageDescriptor)
		descriptors[i].Id = descriptor["id"].(string)
		descriptors[i].Name = descriptor["name"].(string)
		if descriptor["alias"] != nil {
			descriptors[i].Alias = descriptor["alias"].(string)
		} else {
			descriptors[i].Alias = descriptors[i].Name
		}
		if descriptor["icon"] != nil {
			descriptors[i].Icon = descriptor["icon"].(string)
		}
		if descriptor["description"] != nil {
			descriptors[i].Description = descriptor["description"].(string)
		}
		if descriptor["publisherid"] != nil {
			descriptors[i].PublisherId = descriptor["publisherid"].(string)
		}
		if descriptor["publisherid"] != nil {
			descriptors[i].Version = descriptor["version"].(string)
		}
		descriptors[i].Type = packagespb.PackageType(Utility.ToInt(descriptor["type"]))

		if descriptor["keywords"] != nil {
			descriptor["keywords"] = []interface{}(descriptor["keywords"].(primitive.A))
			descriptors[i].Keywords = make([]string, len(descriptor["keywords"].([]interface{})))
			for j := 0; j < len(descriptor["keywords"].([]interface{})); j++ {
				descriptors[i].Keywords[j] = descriptor["keywords"].([]interface{})[j].(string)
			}
		}

		if descriptor["actions"] != nil {
			descriptor["actions"] = []interface{}(descriptor["actions"].(primitive.A))
			descriptors[i].Actions = make([]string, len(descriptor["actions"].([]interface{})))
			for j := 0; j < len(descriptor["actions"].([]interface{})); j++ {
				descriptors[i].Actions[j] = descriptor["actions"].([]interface{})[j].(string)
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
	return &packagespb.GetPackageDescriptorResponse{
		Results: descriptors,
	}, nil
}

//* Return the list of all services *
func (globule *Globule) GetPackagesDescriptor(rqst *packagespb.GetPackagesDescriptorRequest, stream packagespb.PackageDiscovery_GetPackagesDescriptorServer) error {
	p, err := globule.getPersistenceStore()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	data, err := p.Find(context.Background(), "local_resource", "local_resource", "Services", `{}`, "")
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	descriptors := make([]*packagespb.PackageDescriptor, 0)
	for i := 0; i < len(data); i++ {
		descriptor := new(packagespb.PackageDescriptor)

		descriptor.Id = data[i].(map[string]interface{})["id"].(string)
		descriptor.Name = data[i].(map[string]interface{})["name"].(string)
		descriptor.Description = data[i].(map[string]interface{})["description"].(string)
		descriptor.PublisherId = data[i].(map[string]interface{})["publisherid"].(string)
		descriptor.Version = data[i].(map[string]interface{})["version"].(string)
		descriptor.Icon = data[i].(map[string]interface{})["icon"].(string)
		descriptor.Alias = data[i].(map[string]interface{})["alias"].(string)
		descriptor.Type = packagespb.PackageType(Utility.ToInt(data[i].(map[string]interface{})["type"]))

		if data[i].(map[string]interface{})["keywords"] != nil {
			data[i].(map[string]interface{})["keywords"] = []interface{}(data[i].(map[string]interface{})["keywords"].(primitive.A))
			descriptor.Keywords = make([]string, len(data[i].(map[string]interface{})["keywords"].([]interface{})))
			for j := 0; j < len(data[i].(map[string]interface{})["keywords"].([]interface{})); j++ {
				descriptor.Keywords[j] = data[i].(map[string]interface{})["keywords"].([]interface{})[j].(string)
			}
		}

		if data[i].(map[string]interface{})["actions"] != nil {
			data[i].(map[string]interface{})["actions"] = []interface{}(data[i].(map[string]interface{})["actions"].(primitive.A))
			descriptor.Actions = make([]string, len(data[i].(map[string]interface{})["actions"].([]interface{})))
			for j := 0; j < len(data[i].(map[string]interface{})["actions"].([]interface{})); j++ {
				descriptor.Actions[j] = data[i].(map[string]interface{})["actions"].([]interface{})[j].(string)
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
			stream.Send(&packagespb.GetPackagesDescriptorResponse{
				Results: descriptors,
			})
			descriptors = make([]*packagespb.PackageDescriptor, 0)
		}
	}

	if len(descriptors) > 0 {
		stream.Send(&packagespb.GetPackagesDescriptorResponse{
			Results: descriptors,
		})
	}

	// Return the list of Service Descriptor.
	return nil
}

/**
 */
func (globule *Globule) SetPackageDescriptor(ctx context.Context, rqst *packagespb.SetPackageDescriptorRequest) (*packagespb.SetPackageDescriptorResponse, error) {
	p, err := globule.getPersistenceStore()
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
	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Services", `{"id":"`+rqst.Descriptor_.Id+`", "publisherid":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &packagespb.SetPackageDescriptorResponse{
		Result: true,
	}, nil
}

//* Publish a service to service discovery *
func (globule *Globule) PublishPackageDescriptor(ctx context.Context, rqst *packagespb.PublishPackageDescriptorRequest) (*packagespb.PublishPackageDescriptorResponse, error) {

	// Here I will save the descriptor inside the storage...
	p, err := globule.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Append the globule domain to the list of discoveries where the services can be found.
	if !Utility.Contains(rqst.Descriptor_.Discoveries, globule.getDomain()) {
		rqst.Descriptor_.Discoveries = append(rqst.Descriptor_.Discoveries, globule.getDomain())
	}

	// Here I will test if the services already exist...
	_, err = p.FindOne(context.Background(), "local_resource", "local_resource", "Packages", `{"id":"`+rqst.Descriptor_.Id+`", "publisherid":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, "")
	if err == nil {
		// Update existing descriptor.

		// The list of discoveries...
		discoveries, err := Utility.ToJson(rqst.Descriptor_.Discoveries)
		if err == nil {
			values := `{"$set":{"discoveries":` + discoveries + `}}`
			err = p.Update(context.Background(), "local_resource", "local_resource", "Packages", `{"id":"`+rqst.Descriptor_.Id+`", "publisherid":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, values, "")
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
			err = p.Update(context.Background(), "local_resource", "local_resource", "Packages", `{"id":"`+rqst.Descriptor_.Id+`", "publisherid":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, values, "")
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		}

	}

	// The key will be the descriptor string itself.
	jsonStr, err := Utility.ToJson(rqst.Descriptor_)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	jsonStr = strings.ReplaceAll(jsonStr, "publisherId", "publisherid")

	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Packages", `{"id":"`+rqst.Descriptor_.Id+`", "publisherid":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &packagespb.PublishPackageDescriptorResponse{
		Result: true,
	}, nil
}

// Repository
/** Download a service from a service directory **/
func (globule *Globule) DownloadBundle(rqst *packagespb.DownloadBundleRequest, stream packagespb.PackageRepository_DownloadBundleServer) error {
	bundle := new(packagespb.PackageBundle)
	bundle.Plaform = rqst.Plaform
	bundle.Descriptor_ = rqst.Descriptor_

	// Generate the bundle id....
	id := bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Name + "%" + bundle.Descriptor_.Version + "%" + bundle.Descriptor_.Id + "%" + rqst.Plaform

	path := globule.data + "/packages-repository"

	var err error
	// the file must be a zipped archive that contain a .proto, .config and executable.
	bundle.Binairies, err = ioutil.ReadFile(path + "/" + id + ".tar.gz")
	if err != nil {
		return err
	}

	p, err := globule.getPersistenceStore()
	if err != nil {
		return err
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "PackageBundle", `{"_id":"`+id+`"}`, "")
	if err != nil {
		return err
	}

	// init the map with json values.
	checksum := values.(map[string]interface{})

	// Test if the values change over time.
	if Utility.CreateDataChecksum(bundle.Binairies) != checksum["checksum"].(string) {
		return errors.New("the bundle data cheksum is not valid")
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
			rqst := &packagespb.DownloadBundleResponse{
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
func (globule *Globule) UploadBundle(stream packagespb.PackageRepository_UploadBundleServer) error {

	// The bundle will cantain the necessary information to install the service.
	var buffer bytes.Buffer
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			stream.SendAndClose(&packagespb.UploadBundleResponse{
				Result: true,
			})
			err = nil
			break
		} else if err != nil {
			return err
		} else if len(msg.Data) == 0 {
			break
		} else {
			buffer.Write(msg.Data)
		}
	}

	// The buffer that contain the
	dec := gob.NewDecoder(&buffer)
	bundle := new(packagespb.PackageBundle)
	err := dec.Decode(bundle)
	if err != nil {
		return err
	}

	// Generate the bundle id....
	id := bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Name + "%" + bundle.Descriptor_.Version + "%" + bundle.Descriptor_.Id + "%" + bundle.Plaform
	log.Println(id)

	repositoryId := globule.Domain
	if len(repositoryId) > 0 {
		// Now I will append the address of the repository into the service descriptor.
		if !Utility.Contains(bundle.Descriptor_.Repositories, repositoryId) {
			bundle.Descriptor_.Repositories = append(bundle.Descriptor_.Repositories, repositoryId)
			// Publish change into discoveries...
			for i := 0; i < len(bundle.Descriptor_.Discoveries); i++ {
				discoveryId := bundle.Descriptor_.Discoveries[i]
				discoveryService, err := packages_client.NewPackagesDiscoveryService_Client(discoveryId, "packages.PackageDiscovery")
				if err != nil {
					return err
				}
				discoveryService.PublishPackageDescriptor(bundle.Descriptor_)
			}
		}
	}

	path := globule.data + "/packages-repository"
	Utility.CreateDirIfNotExist(path)

	// the file must be a zipped archive that contain a .proto, .config and executable.
	err = ioutil.WriteFile(path+"/"+id+".tar.gz", bundle.Binairies, 0644)
	if err != nil {
		return err
	}

	checksum := Utility.CreateDataChecksum(bundle.Binairies)
	p, err := globule.getPersistenceStore()
	if err != nil {
		return err
	}

	jsonStr, err := Utility.ToJson(map[string]interface{}{"_id": id, "checksum": checksum, "platform": bundle.Plaform, "publisherid": bundle.Descriptor_.PublisherId, "servicename": bundle.Descriptor_.Name, "serviceid": bundle.Descriptor_.Id, "modified": time.Now().Unix(), "size": len(bundle.Binairies)})
	if err != nil {
		return err
	}

	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "PackageBundle", `{"_id":"`+id+`"}`, jsonStr, `[{"upsert": true}]`)

	return err
}
