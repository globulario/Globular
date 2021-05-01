package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"golang.org/x/net/html"

	// "golang.org/x/sys/windows/registry"
	"path/filepath"

	"regexp"
	"strings"
	"sync"

	//	"time"
	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/interceptors"
	"github.com/globulario/services/golang/rbac/rbacpb"

	//"github.com/globulario/services/golang/resource/resourcepb"
	"encoding/json"
	"os/exec"
	"reflect"

	"github.com/globulario/services/golang/packages/packagespb"
	"github.com/globulario/services/golang/security"
	"github.com/golang/protobuf/jsonpb"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"strconv"

	"net"

	"github.com/davecourtois/Utility"
	"github.com/globulario/services/golang/admin/adminpb"
	globular "github.com/globulario/services/golang/globular_service"

	//"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/packages/packages_client"
	"google.golang.org/protobuf/types/known/structpb"

	"google.golang.org/grpc/codes"

	// "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (globule *Globule) startAdminService() error {
	id := string(adminpb.File_proto_admin_proto.Services().Get(0).FullName())
	admin_server, port, err := globule.startInternalService(id, adminpb.File_proto_admin_proto.Path(), globule.Protocol == "https", interceptors.ServerUnaryInterceptor, interceptors.ServerStreamInterceptor) // must be accessible to all clients...
	if err == nil {
		globule.inernalServices = append(globule.inernalServices, admin_server)
		// First of all I will creat a listener.
		// Create the channel to listen on admin port.

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
		if err != nil {
			return err
		}

		adminpb.RegisterAdminServiceServer(admin_server, globule)

		// Here I will make a signal hook to interrupt to exit cleanly.

		go func() {
			// no web-rpc server.
			if err := admin_server.Serve(lis); err != nil {
				log.Println(err)
			}
			// Close it proxy process
			s := globule.getService(id)
			pid := getIntVal(s, "ProxyProcess")
			Utility.TerminateProcess(pid, 0)
			s.Store("ProxyProcess", -1)
		}()

	}
	return err

}

func (globule *Globule) getConfig() map[string]interface{} {

	config := make(map[string]interface{})
	config["Name"] = globule.Name
	config["PortHttp"] = globule.PortHttp
	config["PortHttps"] = globule.PortHttps
	config["AdminEmail"] = globule.AdminEmail
	config["AlternateDomains"] = globule.AlternateDomains
	config["SessionTimeout"] = globule.SessionTimeout
	config["Discoveries"] = globule.Discoveries
	config["PortsRange"] = globule.PortsRange
	config["Version"] = globule.Version
	config["Build"] = globule.Build
	config["Platform"] = globule.Platform
	config["DNS"] = globule.DNS
	config["Protocol"] = globule.Protocol
	config["Domain"] = globule.getDomain()
	config["CertExpirationDelay"] = globule.CertExpirationDelay
	config["ExternalApplications"] = globule.ExternalApplications
	config["CertURL"] = globule.CertURL
	config["CertStableURL"] = globule.CertStableURL
	config["CertExpirationDelay"] = globule.CertExpirationDelay
	config["CertPassword"] = globule.CertPassword
	config["Country"] = globule.Country
	config["State"] = globule.State
	config["City"] = globule.City
	config["Organization"] = globule.Organization
	config["IndexApplication"] = globule.IndexApplication
	config["Path"] = globule.path
	config["Files"] = globule.data + "/files"
	config["Data"] = globule.data
	config["Config"] = globule.config

	// return the full service configuration.
	// Here I will give only the basic services informations and keep
	// all other infromation secret.
	config["Services"] = make(map[string]interface{})

	for _, service_config := range globule.getServices() {
		s := make(map[string]interface{})
		s["Domain"] = getStringVal(service_config, "Domain")
		s["Port"] = getIntVal(service_config, "Port")
		s["Proxy"] = getIntVal(service_config, "Proxy")
		s["TLS"] = getBoolVal(service_config, "TLS")
		s["Version"] = getStringVal(service_config, "Version")
		s["PublisherId"] = getStringVal(service_config, "PublisherId")
		s["KeepUpToDate"] = getBoolVal(service_config, "KeepUpToDate")
		s["KeepAlive"] = getBoolVal(service_config, "KeepAlive")
		s["Description"] = getStringVal(service_config, "Description")
		s["Keywords"] = getVal(service_config, "Keywords")
		s["Repositories"] = getVal(service_config, "Repositories")
		s["Discoveries"] = getVal(service_config, "Discoveries")
		s["State"] = getStringVal(service_config, "State")
		s["Id"] = getStringVal(service_config, "Id")
		s["Name"] = getStringVal(service_config, "Name")
		s["CertFile"] = getStringVal(service_config, "CertFile")
		s["KeyFile"] = getStringVal(service_config, "KeyFile")
		s["CertAuthorityTrust"] = getStringVal(service_config, "CertAuthorityTrust")

		config["Services"].(map[string]interface{})[s["Id"].(string)] = s
	}

	return config

}

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)

	}

	return nil
}

func (globule *Globule) watchConfigFile() {

	doneChan := make(chan bool)

	go func(doneChan chan bool) {
		defer func() {
			doneChan <- true
		}()

		err := watchFile(globule.config + "/config.json")
		if err != nil {
			fmt.Println(err)
		}
		// Run only if the server is running
		if !globule.exit_ {
			// Here I will test if the configuration has change.
			log.Println("configuration was changed and save from external actions.")

			// Here I will read the file.
			data, _ := ioutil.ReadFile(globule.config + "/config.json")
			config := make(map[string]interface{})
			json.Unmarshal(data, &config)
			globule.setConfig(config)

			globule.watchConfigFile() // watch again...
		}
	}(doneChan)

	<-doneChan

}

// Save the configuration file.
func (globule *Globule) saveConfig() {
	if globule.exit_ {
		return // not writting the config file when the server is closing.
	}
	// Here I will save the server attribute
	str, err := Utility.ToJson(globule.toMap())
	if err == nil {
		ioutil.WriteFile(globule.config+"/"+"config.json", []byte(str), 0644)
	} else {
		log.Panicln(err)
	}
}

/**
 * Test if a process with a given name is Running on the server.
 * By default that function is accessible by sa only.
 */
func (globule *Globule) HasRunningProcess(ctx context.Context, rqst *adminpb.HasRunningProcessRequest) (*adminpb.HasRunningProcessResponse, error) {
	ids, err := Utility.GetProcessIdsByName(rqst.Name)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if len(ids) == 0 {
		return &adminpb.HasRunningProcessResponse{
			Result: false,
		}, nil
	}

	return &adminpb.HasRunningProcessResponse{
		Result: true,
	}, nil
}

/**
 * Get the server configuration with all it install services configuration.
 * That function is protected, by default only sa role has access to this information.
 */
func (globule *Globule) GetFullConfig(ctx context.Context, rqst *adminpb.GetConfigRequest) (*adminpb.GetConfigResponse, error) {

	obj, err := structpb.NewStruct(globule.toMap())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.GetConfigResponse{
		Result: obj,
	}, nil

}

// That function return minimal configuration information.
// that function must be accessible by all role, by default guess role has access to this information.
//
func (globule *Globule) GetConfig(ctx context.Context, rqst *adminpb.GetConfigRequest) (*adminpb.GetConfigResponse, error) {
	map_, _ := Utility.ToMap(globule.getConfig())
	obj, err := structpb.NewStruct(map_)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.GetConfigResponse{
		Result: obj,
	}, nil
}

// return true if the configuation has change.
func (globule *Globule) saveServiceConfig(config *sync.Map) bool {
	root, _ := ioutil.ReadFile(os.TempDir() + "/GLOBULAR_ROOT")
	root_ := string(root)[0:strings.Index(string(root), ":")]

	if !Utility.IsLocal(getStringVal(config, "Domain")) && root_ != globule.path {
		return false
	}

	// Here I will
	configPath := globule.getServiceConfigPath(config)
	if len(configPath) == 0 {
		return false
	}

	// set the domain of the service.
	config.Store("Domain", globule.getDomain())

	// format the path's
	config.Store("Path", strings.ReplaceAll(getStringVal(config, "Path"), "\\", "/"))
	config.Store("Proto", strings.ReplaceAll(getStringVal(config, "Proto"), "\\", "/"))

	// so here I will get the previous information...
	f, err := os.Open(configPath)

	if err == nil {
		b, err := ioutil.ReadAll(f)
		if err == nil {
			// get previous configuration...
			config_ := make(map[string]interface{})
			json.Unmarshal(b, &config_)
			config__ := make(map[string]interface{})
			config.Range(func(k, v interface{}) bool {
				config__[k.(string)] = v
				return true
			})

			// Test if there change from the original value's
			if reflect.DeepEqual(config_, config__) {
				log.Println("the values has not change since last read ", configPath)
				f.Close()
				// set back the path's info.
				return false
			}

			// sync the data/config file with the service file.
			jsonStr, _ := Utility.ToJson(config__)
			// here I will write the file
			err = ioutil.WriteFile(configPath, []byte(jsonStr), 0644)
			if err != nil {
				return false
			}

			hub, err := globule.getEventHub()
			if err == nil {
				// Publish event on the network to update the configuration...
				err := hub.Publish("update_globular_service_configuration_evt", []byte(jsonStr))
				if err != nil {
					log.Println("fail to publish event with error: ", err)
				}
			}

		}
	}
	f.Close()

	// Load the services permissions.
	// Here I will get the list of service permission and set it...
	permissions, hasPermissions := config.Load("Permissions")
	if hasPermissions {
		if permissions != nil {
			for i := 0; i < len(permissions.([]interface{})); i++ {
				permission := permissions.([]interface{})[i].(map[string]interface{})
				globule.setActionResourcesPermissions(permission)
			}
		}
	}

	return true
}

func (globule *Globule) setConfig(config map[string]interface{}) {

	// if the configuration is one of services...
	if config["Id"] != nil {
		srv := globule.getService(config["Id"].(string))
		if srv != nil {
			setValues(srv, config)
			globule.initService(srv)
		}

	} else if config["Services"] != nil {

		// Here I will save the configuration
		globule.AdminEmail = config["AdminEmail"].(string)

		globule.Country = config["Country"].(string)
		globule.State = config["State"].(string)
		globule.City = config["City"].(string)
		globule.Organization = config["Organization"].(string)
		globule.CertExpirationDelay = Utility.ToInt(config["CertExpirationDelay"].(float64))
		globule.Name = config["Name"].(string)

		if config["DnsUpdateIpInfos"] != nil {
			globule.DnsUpdateIpInfos = config["DnsUpdateIpInfos"].([]interface{})
		}

		if config["AlternateDomains"] != nil {
			globule.AlternateDomains = config["AlternateDomains"].([]interface{})
		}

		restartServices := false

		httpPort := Utility.ToInt(config["PortHttp"].(float64))
		if httpPort != globule.PortHttp {
			globule.PortHttp = httpPort
			restartServices = true
		}

		httpsPort := Utility.ToInt(config["PortHttps"].(float64))
		if httpsPort != globule.PortHttps {
			globule.PortHttps = httpsPort
			restartServices = true
		}

		protocol := config["Protocol"].(string)
		if globule.Protocol != protocol {
			globule.Protocol = protocol
			restartServices = true
		}

		// The port range
		portsRange := config["PortsRange"].(string)
		if portsRange != globule.PortsRange {
			globule.PortsRange = config["PortsRange"].(string)
			restartServices = true
		}

		domain := config["Domain"].(string)

		if globule.Domain != domain {
			globule.Domain = domain
			restartServices = true
		}

		if config["LdapSyncInfos"] != nil {
			for _, ldapSyncInfos := range config["LdapSyncInfos"].(map[string]interface{}) {
				// update each ldap infos...
				for i := 0; i < len(ldapSyncInfos.([]interface{})); i++ {
					globule.synchronizeLdap(ldapSyncInfos.([]interface{})[i].(map[string]interface{}))
				}
			}
		}

		if restartServices {
			// This will restart the service.
			defer globule.restartServices()
		}

		// Save Discoveries.
		globule.Discoveries = make([]string, 0)
		for i := 0; i < len(config["Discoveries"].([]interface{})); i++ {
			globule.Discoveries = append(globule.Discoveries, config["Discoveries"].([]interface{})[i].(string))
		}

		// Save DNS
		globule.DNS = make([]interface{}, 0)
		for i := 0; i < len(config["DNS"].([]interface{})); i++ {
			globule.DNS = append(globule.DNS, config["DNS"].([]interface{})[i].(string))
		}

		// That will save the services if they have changed.
		for id, s := range config["Services"].(map[string]interface{}) {
			// Attach the actual process and proxy process to the configuration object.
			s_ := globule.getService(id)
			if s_ == nil {
				s_ = new(sync.Map)
			}
			setValues(s_, s.(map[string]interface{}))
			s_.Store("Domain", domain)
			globule.initService(s_)
			globule.setService(s_)
		}
	}
}

// Save a server/service configuration.
// That function must be accessible by Root only.
func (globule *Globule) SaveConfig(ctx context.Context, rqst *adminpb.SaveConfigRequest) (*adminpb.SaveConfigResponse, error) {
	// Save service...
	config := make(map[string]interface{})
	err := json.Unmarshal([]byte(rqst.Config), &config)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will save the server attribute
	str, err := Utility.ToJson(config)
	if err == nil {
		err := ioutil.WriteFile(globule.config+"/"+"config.json", []byte(str), 0644)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	} else {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// return the new configuration file...
	result, _ := Utility.ToJson(config)
	return &adminpb.SaveConfigResponse{
		Result: result,
	}, nil
}

// Uninstall application...
func (globule *Globule) UninstallApplication(ctx context.Context, rqst *adminpb.UninstallApplicationRequest) (*adminpb.UninstallApplicationResponse, error) {

	// Here I will also remove the application permissions...
	if rqst.DeletePermissions {
		log.Println("remove applicaiton permissions...")
	}

	log.Println("remove applicaiton ", rqst.ApplicationId)

	// Same as delete applicaitons.
	err := globule.deleteApplication(rqst.ApplicationId)
	if err != nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.UninstallApplicationResponse{
		Result: true,
	}, nil
}

// Install web Application
func (globule *Globule) InstallApplication(ctx context.Context, rqst *adminpb.InstallApplicationRequest) (*adminpb.InstallApplicationResponse, error) {
	// Get the package bundle from the repository and install it on the server.
	log.Println("Try to install application " + rqst.ApplicationId)

	// Connect to the dicovery services
	package_discovery, err := packages_client.NewPackagesDiscoveryService_Client(rqst.DicorveryId, "packages.PackageDiscovery")

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.DicorveryId)))
	}

	descriptors, err := package_discovery.GetPackageDescriptor(rqst.ApplicationId, rqst.PublisherId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	log.Println("step 1: get application descriptor")
	// The first element in the array is the most recent descriptor
	// so if no version is given the most recent will be taken.
	descriptor := descriptors[0]
	for i := 0; i < len(descriptors); i++ {
		if descriptors[i].Version == rqst.Version {
			descriptor = descriptors[i]
			break
		}
	}

	log.Println("step 2: try to dowload application bundle")
	if len(descriptor.Repositories) == 0 {

		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No service repository was found for application "+descriptor.Id)))
		}

	}

	for i := 0; i < len(descriptor.Repositories); i++ {

		package_repository, err := packages_client.NewServicesRepositoryService_Client(descriptor.Repositories[i], "packages.PackageRepository")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		log.Println("--> try to download application bundle from ", descriptor.Repositories[i])
		bundle, err := package_repository.DownloadBundle(descriptor, "webapp")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		// Create the file.
		r := bytes.NewReader(bundle.Binairies)

		// Now I will install the applicaiton.
		err = globule.installApplication(rqst.Domain, descriptor.Id, descriptor.PublisherId, descriptor.Version, descriptor.Description, descriptor.Icon, descriptor.Alias, r, descriptor.Actions, descriptor.Keywords, descriptor.Roles, descriptor.Groups)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

	}

	log.Println("application was install!")
	return &adminpb.InstallApplicationResponse{
		Result: true,
	}, nil

}

// Intall
func (globule *Globule) installApplication(domain, name, publisherId, version, description string, icon string, alias string, r io.Reader, actions []string, keywords []string, roles []*packagespb.Role, groups []*packagespb.Group) error {

	// Here I will extract the file.
	__extracted_path__, err := Utility.ExtractTarGz(r)
	if err != nil {
		return err
	}

	// remove temporary files.
	defer os.RemoveAll(__extracted_path__)

	// Here I will test that the index.html file is not corrupted...
	__indexHtml__, err := ioutil.ReadFile(__extracted_path__ + "/index.html")
	if err != nil {
		return err
	}

	// The file must contain a linq to a bundle.js file.
	if !strings.Contains(string(__indexHtml__), "./bundle.js") {
		return errors.New("539 something wrong append the index.html file does not contain the bundle.js file... " + string(__indexHtml__))
	}

	// Copy the files to it final destination
	abosolutePath := globule.webRoot

	// If a domain is given.
	if len(domain) > 0 {
		if Utility.Exists(abosolutePath + "/" + domain) {
			abosolutePath += "/" + domain
		}
	}

	// set the absolute application domain.
	abosolutePath += "/" + name

	// Remove the existing files.
	if Utility.Exists(abosolutePath) {
		os.RemoveAll(abosolutePath)
	}

	// Recreate the dir and move file in it.
	Utility.CreateDirIfNotExist(abosolutePath)
	Utility.CopyDir(__extracted_path__+"/.", abosolutePath)

	// Now I will create the application database in the persistence store,
	// and the Application entry in the database.
	// That service made user of persistence service.
	store, err := globule.getPersistenceStore()
	if err != nil {
		return err
	}

	count, err := store.Count(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+name+`"}`, "")
	application := make(map[string]interface{})
	application["_id"] = name
	application["password"] = Utility.GenerateUUID(name)
	application["path"] = "/" + name // The path must be the same as the application name.
	application["publisherid"] = publisherId
	application["version"] = version
	application["description"] = description
	application["actions"] = actions
	application["keywords"] = keywords
	application["icon"] = icon
	application["alias"] = alias

	if len(domain) > 0 {
		if Utility.Exists(globule.webRoot + "/" + domain) {
			application["path"] = "/" + domain + "/" + application["path"].(string)
		}
	}

	application["last_deployed"] = time.Now().Unix() // save it as unix time.

	// Here I will set the resource to manage the applicaiton access permission.
	if err != nil || count == 0 {
		address, port := globule.getBackendAddress()
		// create the application database.
		createApplicationUserDbScript := fmt.Sprintf(
			"db=db.getSiblingDB('%s_db');db.createCollection('application_data');db=db.getSiblingDB('admin');db.createUser({user: '%s', pwd: '%s',roles: [{ role: 'dbOwner', db: '%s_db' }]});",
			name, name, application["password"].(string), name)

		if address == "0.0.0.0" {
			err = store.RunAdminCmd(context.Background(), "local_resource", "sa", globule.RootPassword, createApplicationUserDbScript)
			if err != nil {
				return err
			}
		} else {
			// in the case of remote data store.
			p_, err := globule.getPersistenceSaConnection()
			if err != nil {
				return err
			}
			err = p_.RunAdminCmd("local_resource", "sa", globule.RootPassword, createApplicationUserDbScript)
			if err != nil {
				return err
			}
		}

		application["creation_date"] = time.Now().Unix() // save it as unix time.
		_, err := store.InsertOne(context.Background(), "local_resource", "local_resource", "Applications", application, "")
		if err != nil {
			return err
		}

		p, err := globule.getPersistenceSaConnection()
		if err != nil {
			return err
		}

		err = p.CreateConnection(name+"_db", name+"_db", address, float64(port), 0, name, application["password"].(string), 5000, "", false)
		if err != nil {
			return err
		}

	} else {
		actions_, _ := Utility.ToJson(actions)
		keywords_, _ := Utility.ToJson(keywords)

		err := store.UpdateOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+name+`"}`, `{ "$set":{ "last_deployed":`+Utility.ToString(time.Now().Unix())+` }, "$set":{"keywords":`+keywords_+`}, "$set":{"actions":`+actions_+`},"$set":{"publisherid":"`+publisherId+`"},"$set":{"description":"`+description+`"},"$set":{"alias":"`+alias+`"},"$set":{"icon":"`+icon+`"}, "$set":{"version":"`+version+`"}}`, "")

		if err != nil {
			return err
		}
	}

	// Now I will create/update roles define in the application descriptor...
	for i := 0; i < len(roles); i++ {
		role := roles[i]
		count, err := store.Count(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+role.Id+`"}`, "")
		if err != nil || count == 0 {
			r := make(map[string]interface{})
			r["_id"] = role.Id
			r["name"] = role.Name
			r["actions"] = role.Actions
			r["members"] = []string{}
			_, err := store.InsertOne(context.Background(), "local_resource", "local_resource", "Roles", r, "")
			if err != nil {
				return err
			}
		} else {
			actions_, err := Utility.ToJson(role.Actions)
			if err != nil {
				return err
			}
			err = store.UpdateOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+role.Id+`"}`, `{ "$set":{"name":"`+role.Name+`"}}, { "$set":{"actions":`+actions_+`}}`, "")
			if err != nil {
				return err
			}
		}
	}

	for i := 0; i < len(groups); i++ {
		group := groups[i]
		count, err := store.Count(context.Background(), "local_resource", "local_resource", "Groups", `{"_id":"`+group.Id+`"}`, "")
		if err != nil || count == 0 {
			g := make(map[string]interface{})
			g["_id"] = group.Id
			g["name"] = group.Name
			g["members"] = []string{}
			_, err := store.InsertOne(context.Background(), "local_resource", "local_resource", "Groups", g, "")
			if err != nil {
				return err
			}
		} else {

			err = store.UpdateOne(context.Background(), "local_resource", "local_resource", "Groups", `{"_id":"`+group.Id+`"}`, `{ "$set":{"name":"`+group.Name+`"}}`, "")
			if err != nil {
				return err
			}
		}

	}

	// here is a little workaround to be sure the bundle.js file will not be cached in the brower...
	indexHtml, err := ioutil.ReadFile(abosolutePath + "/index.html")
	if err != nil {
		return err
	}

	// Parse the index html file to be sure the file is valid.
	_, err = html.Parse(strings.NewReader(string(indexHtml)))
	if err != nil {
		return err
	}

	if err == nil {
		var re = regexp.MustCompile(`\/bundle\.js(\?updated=\d*)?`)
		indexHtml_ := re.ReplaceAllString(string(indexHtml), "/bundle.js?updated="+Utility.ToString(time.Now().Unix()))
		if !strings.Contains(indexHtml_, "/bundle.js?updated=") {
			return errors.New("651 something wrong append the index.html file does not contain the bundle.js file... " + indexHtml_)
		}
		// save it back.
		ioutil.WriteFile(abosolutePath+"/index.html", []byte(indexHtml_), 0644)
	}

	return err
}

func (globule *Globule) publishApplication(user, organization, path, name, domain, version, description, icon, alias, repositoryId, discoveryId string, actions, keywords []string, roles []*adminpb.Role) error {

	publisherId := user
	if len(organization) > 0 {
		publisherId = organization
		if !globule.isOrganizationMemeber(user, organization) {
			return errors.New(user + " is not a member of " + organization + "!")
		}
	}

	descriptor := &packagespb.PackageDescriptor{
		Id:           name,
		Name:         name,
		PublisherId:  publisherId,
		Version:      version,
		Description:  description,
		Repositories: []string{},
		Icon:         icon,
		Alias:        alias,
		Actions:      []string{},
		Type:         packagespb.PackageType_APPLICATION_TYPE,
		Roles:        []*packagespb.Role{},
	}

	if len(repositoryId) > 0 {
		descriptor.Repositories = append(descriptor.Repositories, repositoryId)
	}

	if len(discoveryId) > 0 {
		descriptor.Discoveries = append(descriptor.Discoveries, discoveryId)
	}

	if len(keywords) > 0 {
		descriptor.Keywords = keywords
	}

	if len(actions) > 0 {
		descriptor.Actions = actions
	}

	err := globule.publishPackage(user, organization, discoveryId, repositoryId, "webapp", path, descriptor)

	// Set the path of the directory where the application can store date.
	Utility.CreateDirIfNotExist(globule.applications + "/" + name)
	if err != nil {
		return err
	}

	err = globule.addResourceOwner("/applications/"+name, name, rbacpb.SubjectType_APPLICATION)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Update Globular itglobule with a new version.
 */
func (globule *Globule) Update(stream adminpb.AdminService_UpdateServer) error {
	// The buffer that will receive the service executable.
	var buffer bytes.Buffer
	var platform string
	for {
		msg, err := stream.Recv()
		if msg == nil {
			return errors.New("fail to run action adminpb.AdminService_UpdateServer")
		}

		if len(msg.Platform) > 0 {
			platform = msg.Platform
		}

		if err == io.EOF {
			// end of stream...
			stream.SendAndClose(&adminpb.UpdateResponse{})
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

	if len(platform) == 0 {
		return errors.New("no platform was given")
	}

	platform_ := runtime.GOOS + ":" + runtime.GOARCH
	if platform != platform_ {
		return errors.New("Wrong executable platform to update from! wants " + platform_ + " not " + platform)
	}

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	path := filepath.Dir(ex)

	path += "/Globular"
	if runtime.GOOS == "windows" {
		path += ".exe"
	}

	// Move the actual file to other file...
	err = os.Rename(path, path+"_"+Utility.ToString(globule.Build))
	if err != nil {
		return err
	}

	/** So here I will change the current server path and save the new executable **/
	err = ioutil.WriteFile(path, buffer.Bytes(), 0755)
	if err != nil {
		return err
	}

	// set the build time to now...
	globule.Build = time.Now().Unix()
	globule.saveConfig()

	// This will restart the service.
	defer globule.restartServices()

	return nil
}

// Download the actual globular exec file.
func (globule *Globule) DownloadGlobular(rqst *adminpb.DownloadGlobularRequest, stream adminpb.AdminService_DownloadGlobularServer) error {
	platform := rqst.Platform

	if len(platform) == 0 {
		return errors.New("no platform was given")
	}

	platform_ := runtime.GOOS + ":" + runtime.GOARCH
	if platform != platform_ {
		return errors.New("Wrong executable platform to update from! wants " + platform_ + " not " + platform)
	}

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	path := filepath.Dir(ex)
	path += "/Globular"
	if runtime.GOOS == "windows" {
		path += ".exe"
	}

	// No I will stream the result over the networks.
	data, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	reader := bufio.NewReader(data)

	const BufferSize = 1024 * 5 // the chunck size.

	for {
		var data [BufferSize]byte
		bytesread, err := reader.Read(data[0:BufferSize])
		if bytesread > 0 {
			rqst := &adminpb.DownloadGlobularResponse{
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

// Deloyed a web application to a globular node. Mostly use a develeopment time.
func (globule *Globule) DeployApplication(stream adminpb.AdminService_DeployApplicationServer) error {

	// - Get the information from the package.json (npm package, the version, the keywords and set the package descriptor with it.

	// The bundle will cantain the necessary information to install the service.
	var buffer bytes.Buffer

	// Here is necessary information to publish an application.
	var name string
	var domain string
	var user string
	var organization string
	var version string
	var description string
	var repositoryId string
	var discoveryId string
	var keywords []string
	var actions []string
	var icon string
	var alias string
	var roles []*adminpb.Role
	var groups []*adminpb.Group

	for {
		msg, err := stream.Recv()
		if msg == nil {
			return errors.New("fail to run action adminpb.AdminService.DeployApplication")
		}
		if len(msg.Name) > 0 {
			name = msg.Name
		}

		if len(msg.Alias) > 0 {
			alias = msg.Alias
		}

		if len(msg.Domain) > 0 {
			domain = msg.Domain
		}
		if len(msg.Organization) > 0 {
			organization = msg.Organization
		}
		if len(msg.User) > 0 {
			user = msg.User
		}
		if len(msg.Version) > 0 {
			version = msg.Version
		}
		if len(msg.Description) > 0 {
			description = msg.Description
		}
		if msg.Keywords != nil {
			keywords = msg.Keywords
		}
		if len(msg.Repository) > 0 {
			repositoryId = msg.Repository
		}
		if len(msg.Discovery) > 0 {
			discoveryId = msg.Discovery
		}

		if len(msg.Roles) > 0 {
			roles = msg.Roles
		}

		if len(msg.Groups) > 0 {
			groups = msg.Groups
		}

		if len(msg.Actions) > 0 {
			actions = msg.Actions
		}

		if len(msg.Icon) > 0 {
			icon = msg.Icon
		}

		if err == io.EOF {
			// end of stream...
			stream.SendAndClose(&adminpb.DeployApplicationResponse{
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

	if len(repositoryId) == 0 {
		repositoryId = domain
	}

	if len(discoveryId) == 0 {
		discoveryId = domain
	}

	p, err := globule.getPersistenceStore()
	if err != nil {
		return err
	}

	var previousVersion string
	previous, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+name+`"}`, `[{"Projection":{"version":1}}]`)
	if err == nil {
		if previous != nil {
			if previous.(map[string]interface{})["version"] != nil {
				previousVersion = previous.(map[string]interface{})["version"].(string)
			}
		}
	}

	// Now I will save the bundle into a file in the temp directory.
	path := os.TempDir() + "/" + Utility.RandomUUID()
	defer os.RemoveAll(path)

	err = ioutil.WriteFile(path, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	err = globule.publishApplication(user, organization, path, name, domain, version, description, icon, alias, repositoryId, discoveryId, actions, keywords, roles)

	if err != nil {
		return err
	}

	// convert struct...
	roles_ := make([]*packagespb.Role, len(roles))
	for i := 0; i < len(roles); i++ {
		roles_[i] = new(packagespb.Role)
		roles_[i].Id = roles[i].Id
		roles_[i].Name = roles[i].Name
		roles_[i].Actions = roles[i].Actions
	}

	groups_ := make([]*packagespb.Group, len(groups))
	for i := 0; i < len(groups); i++ {
		groups_[i] = new(packagespb.Group)
		groups_[i].Id = groups[i].Id
		groups_[i].Name = groups[i].Name
	}

	// Read bytes and extract it in the current directory.
	r := bytes.NewReader(buffer.Bytes())
	err = globule.installApplication(domain, name, organization, version, description, icon, alias, r, actions, keywords, roles_, groups_)
	if err != nil {
		return err
	}

	// If the version has change I will notify current users and undate the applications.
	if previousVersion != version {
		eventClient, err := globule.getEventHub()
		if err == nil {
			eventClient.Publish("update_"+strings.Split(domain, ":")[0]+"_"+name+"_evt", []byte(version))
		}

		/** The message to send from the notification */
		message := `<div style="display: flex; flex-direction: column">
              <div>A new version of <span style="font-weight: 500;">` + name + `</span> (v.` + version + `) is available.
              </div>
              <div>
                Press <span style="font-weight: 500;">f5</span> to refresh the page.
              </div>
            </div>
            `

		return globule.sendApplicationNotification(name, message)
	}
	return nil
}

/** Create the super administrator in the db. **/
func (globule *Globule) registerSa() error {

	configs := globule.getServiceConfigByName("persistence.PersistenceService")
	if len(configs) == 0 {
		logger.Info("No persistence service was configure on that globule!")
		return errors.New("no persistence service was configure on that globule")
	}

	// Here I will test if mongo db exist on the server.
	existMongo := exec.Command("mongod", "--version")
	err := existMongo.Run()
	if err != nil {
		logger.Info("Failt to run command mongod --version ", err.Error())
		log.Println("fail to start mongo db!", err)
		return err
	}

	// Here I will create super admin if it not already exist.
	dataPath := globule.data + "/mongodb-data"

	if !Utility.Exists(dataPath) {
		// Kill mongo db server if the process already run...
		globule.stopMongod()

		// Here I will create the directory
		err := os.MkdirAll(dataPath, os.ModeDir)
		if err != nil {
			return err
		}

		// Now I will start the command
		mongod := exec.Command("mongod", "--port", "27017", "--dbpath", dataPath)
		err = mongod.Start()
		if err != nil {
			return err
		}

		globule.waitForMongo(60, false)

		// Now I will create a new user name sa and give it all admin write.
		createSaScript := fmt.Sprintf(
			`db=db.getSiblingDB('admin');db.createUser({ user: '%s', pwd: '%s', roles: ['userAdminAnyDatabase','userAdmin','readWrite','dbAdmin','clusterAdmin','readWriteAnyDatabase','dbAdminAnyDatabase']});`, "sa", globule.RootPassword) // must be change...

		createSaCmd := exec.Command("mongo", "--eval", createSaScript)
		err = createSaCmd.Run()
		if err != nil {
			// remove the mongodb-data
			os.RemoveAll(dataPath)
			return err
		}
		globule.stopMongod()
	}

	// Now I will start mongod with auth available.
	mongod := exec.Command("mongod", "--auth", "--port", "27017", "--bind_ip", "0.0.0.0", "--dbpath", dataPath)
	logger.Info("try to start mongod with path ", dataPath)
	err = mongod.Start()
	if err != nil {

		logger.Info("Fail to strart mongo with error ", err)
		return err
	}

	// wait 15 seconds that the server restart.
	globule.waitForMongo(60, true)

	// Get the list of all services method.
	return globule.registerMethods()
}

// Set the root password
func (globule *Globule) SetRootPassword(ctx context.Context, rqst *adminpb.SetRootPasswordRequest) (*adminpb.SetRootPasswordResponse, error) {
	// Here I will set the root password.
	if globule.RootPassword != rqst.OldPassword {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("wrong password given")))
	}

	// Now I will update de sa password.
	globule.RootPassword = rqst.NewPassword

	// Now update the sa password in mongo db.
	changeRootPasswordScript := fmt.Sprintf(
		"db=db.getSiblingDB('admin');db.changeUserPassword('%s','%s');", "sa", rqst.NewPassword)

	p, err := globule.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// I will execute the script with the admin function.
	address, _ := globule.getBackendAddress()
	if address == "0.0.0.0" {
		err = p.RunAdminCmd(context.Background(), "local_resource", "sa", rqst.OldPassword, changeRootPasswordScript)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	} else {
		p_, err := globule.getPersistenceSaConnection()
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		err = p_.RunAdminCmd("local_resource", "sa", globule.RootPassword, changeRootPasswordScript)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	token, err := ioutil.ReadFile(os.TempDir() + "/" + globule.getDomain() + "_token")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &adminpb.SetRootPasswordResponse{
		Token: string(token),
	}, nil

}

/**
 * Set account passowrd.
 *
 * new_password: new password /TODO validate
 * old_password: old password to authenticate the user.
 */
func (globule *Globule) setPassword(accountId string, oldPassword string, newPassword string) error {

	// First of all I will get the user information from the database.
	p, err := globule.getPersistenceStore()
	if err != nil {
		return err
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+accountId+`"},{"name":"`+accountId+`"} ]}`, ``)
	if err != nil {
		return err
	}

	account := values.(map[string]interface{})

	if len(oldPassword) == 0 {
		return errors.New("you must give your old password")
	}

	// Test the old password.
	if oldPassword != account["password"] {
		if Utility.GenerateUUID(oldPassword) != account["password"] {
			return errors.New("wrong password given")
		}
	}

	// Now update the sa password in mongo db.
	name := account["name"].(string)
	name = strings.ReplaceAll(strings.ReplaceAll(name, ".", "_"), "@", "_")

	changePasswordScript := fmt.Sprintf(
		"db=db.getSiblingDB('admin');db.changeUserPassword('%s','%s');", name, newPassword)

	// I will execute the sript with the admin function.
	address, _ := globule.getBackendAddress()
	if address == "0.0.0.0" {
		err = p.RunAdminCmd(context.Background(), "local_resource", "sa", globule.RootPassword, changePasswordScript)
		if err != nil {
			return err
		}
	} else {
		p_, err := globule.getPersistenceSaConnection()
		if err != nil {
			return err
		}
		err = p_.RunAdminCmd("local_resource", "sa", globule.RootPassword, changePasswordScript)
		if err != nil {
			return err
		}
	}

	// Here I will update the user information.
	account["password"] = Utility.GenerateUUID(newPassword)

	// Here I will save the role.
	jsonStr := "{"
	jsonStr += `"name":"` + account["name"].(string) + `",`
	jsonStr += `"email":"` + account["email"].(string) + `",`
	jsonStr += `"password":"` + account["password"].(string) + `",`
	jsonStr += `"roles":[`

	account["roles"] = []interface{}(account["roles"].(primitive.A))
	for j := 0; j < len(account["roles"].([]interface{})); j++ {
		db := account["roles"].([]interface{})[j].(map[string]interface{})["$db"].(string)
		db = strings.ReplaceAll(db, "@", "_")
		db = strings.ReplaceAll(db, ".", "_")

		jsonStr += `{`
		jsonStr += `"$ref":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$ref"].(string) + `",`
		jsonStr += `"$id":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$id"].(string) + `",`
		jsonStr += `"$db":"` + db + `"`
		jsonStr += `}`
		if j < len(account["roles"].([]interface{}))-1 {
			jsonStr += `,`
		}
	}
	jsonStr += `]`
	jsonStr += "}"

	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"name":"`+account["name"].(string)+`"}`, jsonStr, ``)
	if err != nil {
		return err
	}

	return nil

}

//Set the root password
func (globule *Globule) SetPassword(ctx context.Context, rqst *adminpb.SetPasswordRequest) (*adminpb.SetPasswordResponse, error) {

	// First of all I will get the user information from the database.
	err := globule.setPassword(rqst.AccountId, rqst.OldPassword, rqst.NewPassword)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	token, err := ioutil.ReadFile(os.TempDir() + "/" + globule.getDomain() + "_token")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.SetPasswordResponse{
		Token: string(token),
	}, nil

}

// Set the root password
//
func (globule *Globule) SetEmail(ctx context.Context, rqst *adminpb.SetEmailRequest) (*adminpb.SetEmailResponse, error) {

	// Here I will set the root password.
	// First of all I will get the user information from the database.
	p, err := globule.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	accountId := rqst.AccountId
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+accountId+`"},{"name":"`+accountId+`"} ]}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	account := values.(map[string]interface{})

	if account["email"].(string) != rqst.OldEmail {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("wrong email given")))
	}

	account["email"] = rqst.NewEmail

	// Here I will save the role.
	jsonStr := "{"
	jsonStr += `"name":"` + account["name"].(string) + `",`
	jsonStr += `"email":"` + account["email"].(string) + `",`
	jsonStr += `"password":"` + account["password"].(string) + `",`
	jsonStr += `"roles":[`
	account["roles"] = []interface{}(account["roles"].(primitive.A))
	for j := 0; j < len(account["roles"].([]interface{})); j++ {
		db := account["roles"].([]interface{})[j].(map[string]interface{})["$db"].(string)
		db = strings.ReplaceAll(db, "@", "_")
		db = strings.ReplaceAll(db, ".", "_")
		jsonStr += `{`
		jsonStr += `"$ref":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$ref"].(string) + `",`
		jsonStr += `"$id":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$id"].(string) + `",`
		jsonStr += `"$db":"` + db + `"`
		jsonStr += `}`
		if j < len(account["roles"].([]interface{}))-1 {
			jsonStr += `,`
		}
	}
	jsonStr += `]`
	jsonStr += "}"

	// set the new email.
	account["email"] = rqst.NewEmail

	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"name":"`+account["name"].(string)+`"}`, jsonStr, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// read the local token.
	token, err := ioutil.ReadFile(os.TempDir() + "/" + globule.getDomain() + "_token")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Return the token.
	return &adminpb.SetEmailResponse{
		Token: string(token),
	}, nil
}

//Set the root email
func (globule *Globule) SetRootEmail(ctx context.Context, rqst *adminpb.SetRootEmailRequest) (*adminpb.SetRootEmailResponse, error) {
	// Here I will set the root password.
	if globule.AdminEmail != rqst.OldEmail {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Wrong email given, old email is "+rqst.OldEmail+" not "+globule.AdminEmail+"!")))
	}

	// Now I will update de sa password.
	globule.AdminEmail = rqst.NewEmail

	// save the configuration.
	globule.saveConfig()

	// read the local token.
	token, err := ioutil.ReadFile(os.TempDir() + "/" + globule.getDomain() + "_token")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Return the token.
	return &adminpb.SetRootEmailResponse{
		Token: string(token),
	}, nil
}

// Upload a service package.
func (globule *Globule) UploadServicePackage(stream adminpb.AdminService_UploadServicePackageServer) error {
	// The bundle will cantain the necessary information to install the service.
	path := os.TempDir() + "/" + Utility.RandomUUID()

	fo, err := os.Create(path)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	defer fo.Close()

	for {
		msg, err := stream.Recv()
		if msg == nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("fail to upload service package")))

		}

		if err == nil {
			if len(msg.Organization) > 0 {
				if !globule.isOrganizationMemeber(msg.User, msg.Organization) {
					return status.Errorf(
						codes.Internal,
						Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New(msg.User+" is not a member of "+msg.Organization)))
				}
			}
		}

		if err == io.EOF || len(msg.Data) == 0 {
			// end of stream...
			stream.SendAndClose(&adminpb.UploadServicePackageResponse{
				Path: path,
			})
			err = nil
			break
		} else if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		} else {
			fo.Write(msg.Data)
		}
	}
	return nil
}

// Publish a package, the package can contain an application or a services.
func (globule *Globule) publishPackage(user string, organization string, discovery string, repository string, platform string, path string, descriptor *packagespb.PackageDescriptor) error {

	// Connect to the dicovery services
	services_discovery, err := packages_client.NewPackagesDiscoveryService_Client(discovery, "packages.PackageDiscovery")
	if err != nil {
		return errors.New("Fail to connect to package discovery at " + discovery)
	}

	// Connect to the repository packagespb.
	services_repository, err := packages_client.NewServicesRepositoryService_Client(repository, "packages.PackageRepository")
	if err != nil {
		return errors.New("Fail to connect to package repository at " + repository)
	}

	// Ladies and Gentlemans After one year after tow years services as resource!
	path_ := descriptor.PublisherId + "/" + descriptor.Name + "/" + descriptor.Id + "/" + descriptor.Version

	// So here I will set the permissions
	var permissions *rbacpb.Permissions
	permissions, err = globule.getResourcePermissions(path_)
	if err != nil {
		// Create the permission...
		permissions = &rbacpb.Permissions{
			Allowed: []*rbacpb.Permission{
				//  Exemple of possible permission values.
				&rbacpb.Permission{
					Name:          "publish", // member of the organization can publish the service.
					Applications:  []string{},
					Accounts:      []string{},
					Groups:        []string{},
					Peers:         []string{},
					Organizations: []string{organization},
				},
			},
			Denied: []*rbacpb.Permission{},
			Owners: &rbacpb.Permission{
				Name:          "owner",
				Accounts:      []string{user},
				Applications:  []string{},
				Groups:        []string{},
				Peers:         []string{},
				Organizations: []string{},
			},
		}

		// Set the permissions.
		err = globule.setResourcePermissions(path_, permissions)
		if err != nil {
			return err
		}
	}

	// Test the permission before actualy publish the package.
	hasAccess, isDenied, err := globule.validateAccess(user, rbacpb.SubjectType_ACCOUNT, "publish", path_)
	if !hasAccess || isDenied || err != nil {
		log.Println(err)
		return err
	}

	// Append the user into the list of owner if is not already part of it.
	if !Utility.Contains(permissions.Owners.Accounts, user) {
		permissions.Owners.Accounts = append(permissions.Owners.Accounts, user)
	}

	// Save the permissions.
	err = globule.setResourcePermissions(path_, permissions)
	if err != nil {
		log.Println(err)
		return err
	}

	// Fist of all publish the package descriptor.
	err = services_discovery.PublishPackageDescriptor(descriptor)
	if err != nil {
		log.Println(err)
		return err
	}

	// Upload the service to the repository.
	err = services_repository.UploadBundle(discovery, descriptor.Id, descriptor.PublisherId, platform, path)
	if err != nil {
		log.Println(err)
		return err
	}

	// Here I will send a event to be sure all server will update...
	var marshaler jsonpb.Marshaler
	data, err := marshaler.MarshalToString(descriptor)
	if err != nil {
		log.Println(err)
		return err
	}

	// Here I will send an event that the service has a new version...
	if globule.discorveriesEventHub[discovery] == nil {
		client, err := event_client.NewEventService_Client(discovery, "event.EventService")
		if err != nil {
			log.Println("-->", err)
			return err
		}
		globule.discorveriesEventHub[discovery] = client
	}

	eventId := descriptor.PublisherId + ":" + descriptor.Id
	if descriptor.Type == packagespb.PackageType_SERVICE_TYPE {
		eventId += ":SERVICE_PUBLISH_EVENT"
	} else if descriptor.Type == packagespb.PackageType_APPLICATION_TYPE {
		eventId += ":APPLICATION_PUBLISH_EVENT"
	}

	return globule.discorveriesEventHub[discovery].Publish(eventId, []byte(data))
}

// Publish a service. The service must be install localy on the server.
func (globule *Globule) PublishService(ctx context.Context, rqst *adminpb.PublishServiceRequest) (*adminpb.PublishServiceResponse, error) {

	// Make sure the user is part of the organization if one is given.
	publisherId := rqst.User
	if len(rqst.Organization) > 0 {
		publisherId = rqst.Organization
		if !globule.isOrganizationMemeber(rqst.User, rqst.Organization) {
			return nil, errors.New(rqst.User + " is not member of " + rqst.Organization)
		}
	}

	// Now I will upload the service to the repository...
	descriptor := &packagespb.PackageDescriptor{
		Id:           rqst.ServiceId,
		Name:         rqst.ServiceName,
		PublisherId:  publisherId,
		Version:      rqst.Version,
		Description:  rqst.Description,
		Keywords:     rqst.Keywords,
		Repositories: []string{rqst.RepositoryId},
		Type:         packagespb.PackageType_SERVICE_TYPE,
	}

	err := globule.publishPackage(rqst.User, rqst.Organization, rqst.DicorveryId, rqst.RepositoryId, rqst.Platform, rqst.Path, descriptor)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.PublishServiceResponse{
		Result: true,
	}, nil

}

// Install/Update a service on globular instance.
func (globule *Globule) installService(descriptor *packagespb.PackageDescriptor) error {
	// repository must exist...
	log.Println("step 2: try to dowload service bundle")
	if len(descriptor.Repositories) == 0 {
		return errors.New("No service repository was found for service " + descriptor.Id)
	}

	for i := 0; i < len(descriptor.Repositories); i++ {

		services_repository, err := packages_client.NewServicesRepositoryService_Client(descriptor.Repositories[i], "packages.PackageRepository")
		if err != nil {
			return err
		}

		log.Println("--> try to download service from ", descriptor.Repositories[i])

		bundle, err := services_repository.DownloadBundle(descriptor, globular.GetPlatform())

		if err == nil {
			// Create the file.
			r := bytes.NewReader(bundle.Binairies)
			_extracted_path_, err := Utility.ExtractTarGz(r)
			defer os.RemoveAll(_extracted_path_)
			if err != nil {
				if err.Error() != "EOF" {
					// report the error and try to continue...
					log.Println(err)
				}
			}

			// I will save the binairy in file...
			Utility.CreateDirIfNotExist(globule.path + "/services/")
			err = Utility.CopyDir(_extracted_path_+"/"+descriptor.PublisherId, globule.path+"/services/")
			if err != nil {
				return err
			}

			path := globule.path + "/services/" + descriptor.PublisherId + "/" + descriptor.Name + "/" + descriptor.Version + "/" + descriptor.Id
			configs, _ := Utility.FindFileByName(path, "config.json")

			if len(configs) == 0 {
				log.Println("No configuration file was found at at path ", path)
				return errors.New("no configuration file was found")
			}

			s := make(map[string]interface{})
			data, err := ioutil.ReadFile(configs[0])
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &s)
			if err != nil {
				return err
			}

			protos, _ := Utility.FindFileByName(globule.path+"/services/"+descriptor.PublisherId+"/"+descriptor.Name+"/"+descriptor.Version, ".proto")
			if len(protos) == 0 {
				log.Println("No prototype file was found at at path ", globule.path+"/services/"+descriptor.PublisherId+"/"+descriptor.Name+"/"+descriptor.Version)
				return errors.New("no configuration file was found")
			}

			// I will replace the path inside the config...
			execName := s["Path"].(string)[strings.LastIndex(s["Path"].(string), "/")+1:]
			s["Path"] = path + "/" + execName
			s["Proto"] = protos[0]

			err = os.Chmod(s["Path"].(string), 0755)
			if err != nil {
				log.Println(err)
			}

			jsonStr, _ := Utility.ToJson(s)
			ioutil.WriteFile(configs[0], []byte(jsonStr), 0644)

			// set the service in the map.
			s_ := new(sync.Map)
			setValues(s_, s)
			log.Println("Service is install successfully!")
			// initialyse the new service.
			err = globule.initService(s_)
			if err != nil {
				return err
			}

			// Here I will set the service method...
			globule.setServiceMethods(s["Name"].(string), s["Proto"].(string))

			globule.registerMethods()

			break
		} else {
			log.Println("fail to download error with error ", err)
			return err
		}
	}

	return nil

}

// Install/Update a service on globular instance.
func (globule *Globule) InstallService(ctx context.Context, rqst *adminpb.InstallServiceRequest) (*adminpb.InstallServiceResponse, error) {
	log.Println("Try to install new service...")

	// Connect to the dicovery services
	services_discovery, err := packages_client.NewPackagesDiscoveryService_Client(rqst.DicorveryId, "packages.PackageDiscovery")

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("fail to connect to "+rqst.DicorveryId)))
	}

	descriptors, err := services_discovery.GetPackageDescriptor(rqst.ServiceId, rqst.PublisherId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	log.Println("step 1: get service dscriptor")
	// The first element in the array is the most recent descriptor
	// so if no version is given the most recent will be taken.
	descriptor := descriptors[0]
	for i := 0; i < len(descriptors); i++ {
		if descriptors[i].Version == rqst.Version {
			descriptor = descriptors[i]
			break
		}
	}

	err = globule.installService(descriptor)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	log.Println("Service was install!")
	return &adminpb.InstallServiceResponse{
		Result: true,
	}, nil

}

// Uninstall a service...
func (globule *Globule) UninstallService(ctx context.Context, rqst *adminpb.UninstallServiceRequest) (*adminpb.UninstallServiceResponse, error) {
	p, err := globule.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// First of all I will stop the running service(s) instance.
	for _, s := range globule.getServices() {
		// Stop the instance of the service.
		id, ok := s.Load("Id")
		if ok {
			if getStringVal(s, "PublisherId") == rqst.PublisherId && id == rqst.ServiceId && getStringVal(s, "Version") == rqst.Version {

				globule.stopService(s)
				globule.deleteService(id.(string))

				// Get the list of method to remove from the list of actions.
				toDelete := globule.getServiceMethods(getStringVal(s, "Name"), getStringVal(s, "Proto"))
				methods := make([]string, 0)
				for i := 0; i < len(globule.methods); i++ {
					if !Utility.Contains(toDelete, globule.methods[i]) {
						methods = append(methods, globule.methods[i])
					}
				}

				// Keep permissions use when we update a service.
				if rqst.DeletePermissions {
					// Now I will remove action permissions
					for i := 0; i < len(toDelete); i++ {
						p.Delete(context.Background(), "local_resource", "local_resource", "ActionPermission", `{"action":"`+toDelete[i]+`"}`, "")

						// Delete it from Role.
						p.Update(context.Background(), "local_resource", "local_resource", "Roles", `{}`, `{"$pull":{"actions":"`+toDelete[i]+`"}}`, "")

						// Delete it from Application.
						p.Update(context.Background(), "local_resource", "local_resource", "Applications", `{}`, `{"$pull":{"actions":"`+toDelete[i]+`"}}`, "")

						// Delete it from Peer.
						p.Update(context.Background(), "local_resource", "local_resource", "Peers", `{}`, `{"$pull":{"actions":"`+toDelete[i]+`"}}`, "")

					}
				}

				globule.methods = methods
				globule.registerMethods() // refresh methods.
			}
		}
	}

	// Now I will remove the service.
	// Service are located into the packagespb...
	path := globule.path + "/services/" + rqst.PublisherId + "/" + rqst.ServiceId + "/" + rqst.Version

	// remove directory and sub-directory.
	err = os.RemoveAll(path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// save the config.
	globule.saveConfig()

	return &adminpb.UninstallServiceResponse{
		Result: true,
	}, nil
}

/**
 * Retunr the path of config.json for a given services.
 */
func (globule *Globule) getServiceConfigPath(s *sync.Map) string {

	path := getStringVal(s, "Path")
	index := strings.LastIndex(path, "/")
	if index == -1 {
		return ""
	}

	path = path[0:index] + "/config.json"
	return path
}

func (globule *Globule) stopService(s *sync.Map) error {

	// Set keep alive to false...
	s.Store("KeepAlive", false)
	_, hasProcessPid := s.Load("Process")
	if !hasProcessPid {
		s.Store("Process", -1)
	}

	pid := getIntVal(s, "Process")
	if pid != -1 {

		if runtime.GOOS == "windows" {
			// Program written with dotnet on window need this command to stop...
			kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(pid))
			kill.Stderr = os.Stderr
			kill.Stdout = os.Stdout
			kill.Run()
		} else {
			err := Utility.TerminateProcess(pid, 0)
			if err != nil {
				log.Println("fail to teminate process ", pid)
			}
		}
	}

	_, hasProxyProcessPid := s.Load("ProxyProcess")
	if !hasProxyProcessPid {
		s.Store("ProxyProcess", -1)
	}
	pid = getIntVal(s, "ProxyProcess")
	if pid != -1 {
		err := Utility.TerminateProcess(pid, 0)
		if err != nil {
			log.Println("fail to teminate proxy process ", pid)
		}
	}

	s.Store("Process", -1)
	s.Store("ProxyProcess", -1)
	s.Store("State", "stopped")

	config := make(map[string]interface{})
	s.Range(func(k, v interface{}) bool {
		config[k.(string)] = v
		return true
	})

	// sync the data/config file with the service file.
	jsonStr, _ := Utility.ToJson(config)

	// here I will write the file
	configPath := globule.getServiceConfigPath(s)
	if len(configPath) > 0 {
		err := ioutil.WriteFile(configPath, []byte(jsonStr), 0644)
		if err != nil {
			return err
		}
	}

	// globule.logServiceInfo(getStringVal(s, "Name"), time.Now().String()+"Service "+getStringVal(s, "Name")+" was stopped!")
	globule.saveConfig()
	return nil
}

// Stop a service
func (globule *Globule) StopService(ctx context.Context, rqst *adminpb.StopServiceRequest) (*adminpb.StopServiceResponse, error) {

	s := globule.getService(rqst.ServiceId)
	if s != nil {
		err := globule.stopService(s)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	} else {
		// Close all services with a given name.
		services := globule.getServiceConfigByName(rqst.ServiceId)
		for i := 0; i < len(services); i++ {
			serviceId := services[i]["Id"].(string)
			s := globule.getService(serviceId)
			if s == nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No service found with id "+serviceId)))
			}
			err := globule.stopService(s)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		}
	}

	return &adminpb.StopServiceResponse{
		Result: true,
	}, nil

}

// Start a service
func (globule *Globule) StartService(ctx context.Context, rqst *adminpb.StartServiceRequest) (*adminpb.StartServiceResponse, error) {

	s := globule.getService(rqst.ServiceId)
	proxy_pid := int64(-1)
	service_pid := int64(-1)

	if s == nil {
		services := globule.getServiceConfigByName(rqst.ServiceId)
		for i := 0; i < len(services); i++ {
			id := services[i]["Id"].(string)
			s := globule.getService(id)
			service_pid_, proxy_pid_, err := globule.startService(s)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			proxy_pid = int64(proxy_pid_)
			service_pid = int64(service_pid_)
		}
	} else {
		service_pid_, proxy_pid_, err := globule.startService(s)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		proxy_pid = int64(proxy_pid_)
		service_pid = int64(service_pid_)
	}

	return &adminpb.StartServiceResponse{
		ProxyPid:   proxy_pid,
		ServicePid: service_pid,
	}, nil
}

// Restart all Services also the http(s)
func (globule *Globule) RestartServices(ctx context.Context, rqst *adminpb.RestartServicesRequest) (*adminpb.RestartServicesResponse, error) {
	log.Println("restart service... ")
	globule.restartServices()

	return &adminpb.RestartServicesResponse{}, nil
}

// That command is use to restart a new instance of the globular.
func rerunDetached() error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := os.Args[0]
	cmd := exec.Command(path, []string{}...)
	cmd.Dir = cwd
	err = cmd.Start()
	if err != nil {
		return err
	}
	cmd.Process.Release()
	return nil
}

func (globule *Globule) restartServices() {
	if globule.exit_ {
		return // already restarting I will ingnore the call.
	}

	// Stop all internal services
	globule.stopInternalServices()

	// Stop all external services.
	globule.stopServices()

	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	delay := 5
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:

				fmt.Println("restart in ", delay, "second")
				delay--
			}
		}
	}()

	time.Sleep(time.Duration(delay) * time.Second)
	ticker.Stop()
	done <- true

	log.Println("Globular process will now restart...")

	// Restart a new process.
	rerunDetached()

	os.Exit(0) // exit the process to release all ressources.

}

// Start an external service here.
func (globule *Globule) startExternalApplication(serviceId string) (int, error) {

	if service, ok := globule.ExternalApplications[serviceId]; !ok {
		return -1, errors.New("No external service found with name " + serviceId)
	} else {

		service.srv = exec.Command(service.Path, service.Args...)

		err := service.srv.Start()
		if err != nil {
			return -1, err
		}

		// save back the service in the map.
		globule.ExternalApplications[serviceId] = service

		return service.srv.Process.Pid, nil
	}

}

// Stop external service.
func (globule *Globule) stopExternalApplication(serviceId string) error {
	if _, ok := globule.ExternalApplications[serviceId]; !ok {
		return errors.New("No external service found with name " + serviceId)
	}

	// if no command was created
	if globule.ExternalApplications[serviceId].srv == nil {
		return nil
	}

	// if no process running
	if globule.ExternalApplications[serviceId].srv.Process == nil {
		return nil
	}

	// kill the process.
	return globule.ExternalApplications[serviceId].srv.Process.Signal(os.Interrupt)
}

// Kill process by id
func (globule *Globule) KillProcess(ctx context.Context, rqst *adminpb.KillProcessRequest) (*adminpb.KillProcessResponse, error) {
	pid := int(rqst.Pid)
	err := Utility.TerminateProcess(pid, 0)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.KillProcessResponse{}, nil
}

// Kill process by name
func (globule *Globule) KillProcesses(ctx context.Context, rqst *adminpb.KillProcessesRequest) (*adminpb.KillProcessesResponse, error) {
	err := Utility.KillProcessByName(rqst.Name)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.KillProcessesResponse{}, nil
}

// Return the list of process id with a given name.
func (globule *Globule) GetPids(ctx context.Context, rqst *adminpb.GetPidsRequest) (*adminpb.GetPidsResponse, error) {
	pids_, err := Utility.GetProcessIdsByName(rqst.Name)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	pids := make([]int32, len(pids_))
	for i := 0; i < len(pids_); i++ {
		pids[i] = int32(pids_[i])
	}

	return &adminpb.GetPidsResponse{
		Pids: pids,
	}, nil
}

// Register external service to be start by Globular in order to run
func (globule *Globule) RegisterExternalApplication(ctx context.Context, rqst *adminpb.RegisterExternalApplicationRequest) (*adminpb.RegisterExternalApplicationResponse, error) {

	// Here I will get the command path.
	externalCmd := ExternalApplication{
		Id:   rqst.ServiceId,
		Path: rqst.Path,
		Args: rqst.Args,
	}

	globule.ExternalApplications[externalCmd.Id] = externalCmd

	// save the config.
	globule.saveConfig()

	// start the external service.
	pid, err := globule.startExternalApplication(externalCmd.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.RegisterExternalApplicationResponse{
		ServicePid: int64(pid),
	}, nil
}

/**
 * Read output and send it to a channel.
 */
func ReadOutput(output chan string, rc io.ReadCloser) {

	cutset := "\r\n"
	for {
		buf := make([]byte, 3000)
		n, err := rc.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			if n == 0 {
				break
			}
		}
		text := strings.TrimSpace(string(buf[:n]))
		for {
			// Take the index of any of the given cutset
			n := strings.IndexAny(text, cutset)
			if n == -1 {
				// If not found, but still have data, send it
				if len(text) > 0 {
					output <- text
				}
				break
			}
			// Send data up to the found cutset
			output <- text[:n]
			// If cutset is last element, stop there.
			if n == len(text) {
				break
			}
			// Shift the text and start again.
			text = text[n+1:]
		}
	}

}

// Run an external command must be use with care.
func (globule *Globule) RunCmd(rqst *adminpb.RunCmdRequest, stream adminpb.AdminService_RunCmdServer) error {

	baseCmd := rqst.Cmd
	cmdArgs := rqst.Args
	isBlocking := rqst.Blocking
	pid := -1
	cmd := exec.Command(baseCmd, cmdArgs...)
	if isBlocking {

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		output := make(chan string)
		done := make(chan bool)

		// Process message util the command is done.
		go func() {
			for {
				select {
				case <-done:
					break

				case result := <-output:
					if cmd.Process != nil {
						pid = cmd.Process.Pid
					}

					stream.Send(
						&adminpb.RunCmdResponse{
							Pid:    int32(pid),
							Result: result,
						},
					)
				}
			}

		}()

		// Start reading the output
		go ReadOutput(output, stdout)

		cmd.Run()

		cmd.Wait()

		// Close the output.
		stdout.Close()
		done <- true

	} else {
		err := cmd.Start()
		if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		if cmd.Process != nil {
			pid = cmd.Process.Pid
		}

		stream.Send(
			&adminpb.RunCmdResponse{
				Pid:    int32(pid),
				Result: "",
			},
		)

	}

	return nil
}

// Set environement variable.
func (globule *Globule) SetEnvironmentVariable(ctx context.Context, rqst *adminpb.SetEnvironmentVariableRequest) (*adminpb.SetEnvironmentVariableResponse, error) {
	err := Utility.SetEnvironmentVariable(rqst.Name, rqst.Value)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &adminpb.SetEnvironmentVariableResponse{}, nil
}

// Get environement variable.
func (globule *Globule) GetEnvironmentVariable(ctx context.Context, rqst *adminpb.GetEnvironmentVariableRequest) (*adminpb.GetEnvironmentVariableResponse, error) {
	value, err := Utility.GetEnvironmentVariable(rqst.Name)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &adminpb.GetEnvironmentVariableResponse{
		Value: value,
	}, nil
}

// Delete environement variable.
func (globule *Globule) UnsetEnvironmentVariable(ctx context.Context, rqst *adminpb.UnsetEnvironmentVariableRequest) (*adminpb.UnsetEnvironmentVariableResponse, error) {

	err := Utility.UnsetEnvironmentVariable(rqst.Name)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &adminpb.UnsetEnvironmentVariableResponse{}, nil
}

///////////////////////////////////// API //////////////////////////////////////

// Get certificates from the server and copy them into the the a given directory.
// path: The path where to copy the certificates
// port: The server configuration port the default is 80.
//
// ex. Here is an exemple of the command run from the shell,
//
// Globular certificates -domain=globular.cloud -path=/tmp -port=80
//
// The command can return
func (globule *Globule) InstallCertificates(ctx context.Context, rqst *adminpb.InstallCertificatesRequest) (*adminpb.InstallCertificatesResponse, error) {
	path := rqst.Path
	if len(path) == 0 {
		path = os.TempDir()
	}

	port := 80
	if rqst.Port != 0 {
		port = int(rqst.Port)
	}

	key, cert, ca, err := security.InstallCertificates(rqst.Domain, port, path)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.InstallCertificatesResponse{
		Certkey: key,
		Cert:    cert,
		Cacert:  ca,
	}, nil
}
