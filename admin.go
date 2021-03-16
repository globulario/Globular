package main

import (
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

	"regexp"
	"strings"
	"sync"

	//	"time"

	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/rbac/rbacpb"

	//"github.com/globulario/services/golang/resource/resourcepb"
	"encoding/json"
	"os/exec"
	"reflect"

	"github.com/globulario/Globular/security"
	"github.com/globulario/services/golang/packages/packagespb"
	"github.com/golang/protobuf/jsonpb"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"strconv"

	"net"

	"github.com/davecourtois/Utility"
	"github.com/globulario/Globular/Interceptors"
	"github.com/globulario/services/golang/admin/adminpb"
	globular "github.com/globulario/services/golang/globular_service"

	//"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/packages/packages_client"
	"google.golang.org/protobuf/types/known/structpb"

	"google.golang.org/grpc/codes"

	// "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (self *Globule) startAdminService() error {
	id := string(adminpb.File_proto_admin_proto.Services().Get(0).FullName())
	admin_server, port, err := self.startInternalService(id, adminpb.File_proto_admin_proto.Path(), self.Protocol == "https", Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor) // must be accessible to all clients...
	if err == nil {
		self.inernalServices = append(self.inernalServices, admin_server)
		// First of all I will creat a listener.
		// Create the channel to listen on admin port.

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
		if err != nil {
			return err
		}

		adminpb.RegisterAdminServiceServer(admin_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.

		go func() {
			// no web-rpc server.
			if err := admin_server.Serve(lis); err != nil {
				log.Println(err)
			}
			// Close it proxy process
			s := self.getService(id)
			pid := getIntVal(s, "ProxyProcess")
			Utility.TerminateProcess(pid, 0)
			s.Store("ProxyProcess", -1)
		}()

	}
	return err

}

func (self *Globule) getConfig() map[string]interface{} {

	config := make(map[string]interface{}, 0)
	config["Name"] = self.Name
	config["PortHttp"] = self.PortHttp
	config["PortHttps"] = self.PortHttps
	config["AdminEmail"] = self.AdminEmail
	config["AlternateDomains"] = self.AlternateDomains
	config["SessionTimeout"] = self.SessionTimeout
	config["Discoveries"] = self.Discoveries
	config["PortsRange"] = self.PortsRange
	config["Version"] = self.Version
	config["Platform"] = self.Platform
	config["DNS"] = self.DNS
	config["Protocol"] = self.Protocol
	config["Domain"] = self.getDomain()
	config["CertExpirationDelay"] = self.CertExpirationDelay
	config["ExternalApplications"] = self.ExternalApplications
	config["CertURL"] = self.CertURL
	config["CertStableURL"] = self.CertStableURL
	config["CertExpirationDelay"] = self.CertExpirationDelay
	config["CertPassword"] = self.CertPassword
	config["Country"] = self.Country
	config["State"] = self.State
	config["City"] = self.City
	config["Organization"] = self.Organization

	// return the full service configuration.
	// Here I will give only the basic services informations and keep
	// all other infromation secret.
	config["Services"] = make(map[string]interface{})

	for _, service_config := range self.getServices() {
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

func (self *Globule) watchConfigFile() {

	doneChan := make(chan bool)

	go func(doneChan chan bool) {
		defer func() {
			doneChan <- true
		}()

		err := watchFile(self.config + string(os.PathSeparator) + "config.json")
		if err != nil {
			fmt.Println(err)
		}
		// Run only if the server is running
		if !self.exit_ {
			// Here I will test if the configuration has change.
			log.Println("configuration was changed and save from external actions.")

			// Here I will read the file.
			data, _ := ioutil.ReadFile(self.config + string(os.PathSeparator) + "config.json")
			config := make(map[string]interface{}, 0)
			json.Unmarshal(data, &config)
			self.setConfig(config)

			self.watchConfigFile() // watch again...
		}
	}(doneChan)

	<-doneChan

}

// Save the configuration file.
func (self *Globule) saveConfig() {
	if self.exit_ {
		return // not writting the config file when the server is closing.
	}
	// Here I will save the server attribute
	str, err := Utility.ToJson(self.toMap())
	if err == nil {
		ioutil.WriteFile(self.config+string(os.PathSeparator)+"config.json", []byte(str), 0644)
	} else {
		log.Panicln(err)
	}
}

/**
 * Test if a process with a given name is Running on the server.
 * By default that function is accessible by sa only.
 */
func (self *Globule) HasRunningProcess(ctx context.Context, rqst *adminpb.HasRunningProcessRequest) (*adminpb.HasRunningProcessResponse, error) {
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
func (self *Globule) GetFullConfig(ctx context.Context, rqst *adminpb.GetConfigRequest) (*adminpb.GetConfigResponse, error) {

	obj, err := structpb.NewStruct(self.toMap())
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
func (self *Globule) GetConfig(ctx context.Context, rqst *adminpb.GetConfigRequest) (*adminpb.GetConfigResponse, error) {
	map_, _ := Utility.ToMap(self.getConfig())
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
func (self *Globule) saveServiceConfig(config *sync.Map) bool {
	root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")
	root_ := string(root)[0:strings.Index(string(root), ":")]

	if !Utility.IsLocal(getStringVal(config, "Domain")) && root_ != self.path {
		return false
	}

	// Here I will
	configPath := self.getServiceConfigPath(config)
	if len(configPath) == 0 {
		return false
	}

	// set the domain of the service.
	config.Store("Domain", self.getDomain())

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

			hub, err := self.getEventHub()
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
				self.setActionResourcesPermissions(permission)
			}
		}
	}

	return true
}

func (self *Globule) setConfig(config map[string]interface{}) {

	// if the configuration is one of services...
	if config["Id"] != nil {
		srv := self.getService(config["Id"].(string))
		if srv != nil {
			setValues(srv, config)
			self.initService(srv)
		}

	} else if config["Services"] != nil {

		// Here I will save the configuration
		self.AdminEmail = config["AdminEmail"].(string)

		self.Country = config["Country"].(string)
		self.State = config["State"].(string)
		self.City = config["City"].(string)
		self.Organization = config["Organization"].(string)
		self.CertExpirationDelay = Utility.ToInt(config["CertExpirationDelay"].(float64))
		self.Name = config["Name"].(string)

		if config["DnsUpdateIpInfos"] != nil {
			self.DnsUpdateIpInfos = config["DnsUpdateIpInfos"].([]interface{})
		}

		if config["AlternateDomains"] != nil {
			self.AlternateDomains = config["AlternateDomains"].([]interface{})
		}

		restartServices := false

		httpPort := Utility.ToInt(config["PortHttp"].(float64))
		if httpPort != self.PortHttp {
			self.PortHttp = httpPort
			restartServices = true
		}

		httpsPort := Utility.ToInt(config["PortHttps"].(float64))
		if httpsPort != self.PortHttps {
			self.PortHttps = httpsPort
			restartServices = true
		}

		protocol := config["Protocol"].(string)
		if self.Protocol != protocol {
			self.Protocol = protocol
			restartServices = true
		}

		// The port range
		portsRange := config["PortsRange"].(string)
		if portsRange != self.PortsRange {
			self.PortsRange = config["PortsRange"].(string)
			restartServices = true
		}

		domain := config["Domain"].(string)

		if self.Domain != domain {
			self.Domain = domain
			restartServices = true
		}

		if config["LdapSyncInfos"] != nil {
			for _, ldapSyncInfos := range config["LdapSyncInfos"].(map[string]interface{}) {
				// update each ldap infos...
				for i := 0; i < len(ldapSyncInfos.([]interface{})); i++ {
					self.synchronizeLdap(ldapSyncInfos.([]interface{})[i].(map[string]interface{}))
				}
			}
		}

		if restartServices {
			// This will restart the service.
			defer self.restartServices()
		}

		// Save Discoveries.
		self.Discoveries = make([]string, 0)
		for i := 0; i < len(config["Discoveries"].([]interface{})); i++ {
			self.Discoveries = append(self.Discoveries, config["Discoveries"].([]interface{})[i].(string))
		}

		// Save DNS
		self.DNS = make([]interface{}, 0)
		for i := 0; i < len(config["DNS"].([]interface{})); i++ {
			self.DNS = append(self.DNS, config["DNS"].([]interface{})[i].(string))
		}

		// That will save the services if they have changed.
		for id, s := range config["Services"].(map[string]interface{}) {
			// Attach the actual process and proxy process to the configuration object.
			s_ := self.getService(id)
			if s_ == nil {
				s_ = new(sync.Map)
			}
			setValues(s_, s.(map[string]interface{}))
			s_.Store("Domain", domain)
			self.initService(s_)
			self.setService(s_)
		}
	}
}

// Save a server/service configuration.
// That function must be accessible by Root only.
func (self *Globule) SaveConfig(ctx context.Context, rqst *adminpb.SaveConfigRequest) (*adminpb.SaveConfigResponse, error) {
	// Save service...
	config := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(rqst.Config), &config)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will save the server attribute
	str, err := Utility.ToJson(config)
	if err == nil {
		err := ioutil.WriteFile(self.config+string(os.PathSeparator)+"config.json", []byte(str), 0644)
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
func (self *Globule) UninstallApplication(ctx context.Context, rqst *adminpb.UninstallApplicationRequest) (*adminpb.UninstallApplicationResponse, error) {
	// TODO Implement it!
	return nil, status.Errorf(
		codes.Internal,
		Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Not implemented!")))
}

// Install web Application
func (self *Globule) InstallApplication(ctx context.Context, rqst *adminpb.InstallApplicationRequest) (*adminpb.InstallApplicationResponse, error) {
	// Get the package bundle from the repository and install it on the server.
	// TODO Implement it!
	return nil, status.Errorf(
		codes.Internal,
		Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Not implemented!")))
}

// Intall
func (self *Globule) installApplication(domain, name, organization, version, description string, r io.Reader, actions []string, keywords []string) error {
	// Here I will extract the file.
	Utility.ExtractTarGz(r)

	// Copy the files to it final destination
	abosolutePath := self.webRoot

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

	// remove temporary files.
	defer os.RemoveAll(Utility.GenerateUUID(name))

	// Recreate the dir and move file in it.
	Utility.CreateDirIfNotExist(abosolutePath)
	Utility.CopyDirContent(Utility.GenerateUUID(name), abosolutePath)

	// Now I will create the application database in the persistence store,
	// and the Application entry in the database.
	// That service made user of persistence service.
	store, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	count, err := store.Count(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+name+`"}`, "")
	application := make(map[string]interface{})
	application["_id"] = name
	application["password"] = Utility.GenerateUUID(name)
	application["path"] = "/" + name // The path must be the same as the application name.
	application["organization"] = organization
	application["version"] = version
	application["description"] = description
	application["actions"] = actions
	application["keywords"] = keywords

	if len(domain) > 0 {
		if Utility.Exists(self.webRoot + "/" + domain) {
			application["path"] = "/" + domain + "/" + application["path"].(string)
		}
	}

	application["last_deployed"] = time.Now().Unix() // save it as unix time.

	// Here I will set the resource to manage the applicaiton access permission.
	if err != nil || count == 0 {
		address, port := self.getBackendAddress()
		// create the application database.
		createApplicationUserDbScript := fmt.Sprintf(
			"db=db.getSiblingDB('%s_db');db.createCollection('application_data');db=db.getSiblingDB('admin');db.createUser({user: '%s', pwd: '%s',roles: [{ role: 'dbOwner', db: '%s_db' }]});",
			name, name, application["password"].(string), name)

		if address == "0.0.0.0" {
			err = store.RunAdminCmd(context.Background(), "local_resource", "sa", self.RootPassword, createApplicationUserDbScript)
			if err != nil {
				return err
			}
		} else {
			// in the case of remote data store.
			p_, err := self.getPersistenceSaConnection()
			if err != nil {
				return err
			}
			err = p_.RunAdminCmd("local_resource", "sa", self.RootPassword, createApplicationUserDbScript)
			if err != nil {
				return err
			}
		}

		application["creation_date"] = time.Now().Unix() // save it as unix time.
		_, err := store.InsertOne(context.Background(), "local_resource", "local_resource", "Applications", application, "")
		if err != nil {
			return err
		}
		p, err := self.getPersistenceSaConnection()

		err = p.CreateConnection(name+"_db", name+"_db", address, float64(port), 0, name, application["password"].(string), 5000, "", false)
		if err != nil {
			return err
		}

	} else {
		actions_, _ := Utility.ToJson(actions)
		keywords_, _ := Utility.ToJson(keywords)

		err := store.UpdateOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+name+`"}`, `{ "$set":{ "last_deployed":`+Utility.ToString(time.Now().Unix())+` }, "$set":{"keywords":`+keywords_+`}, "$set":{"actions":`+actions_+`},"$set":{"organization":"`+organization+`"},"$set":{"description":"`+description+`"}, "$set":{"version":"`+version+`"}}`, "")
		if err != nil {
			return err
		}
	}

	// here is a little workaround to be sure the bundle.js file will not be cached in the brower...
	indexHtml, err := ioutil.ReadFile(abosolutePath + "/index.html")

	// Parse the index html file to be sure the file is valid.
	_, err = html.Parse(strings.NewReader(string(indexHtml)))
	if err != nil {
		return err
	}

	if err == nil {
		var re = regexp.MustCompile(`\/bundle\.js(\?updated=\d*)?`)
		indexHtml_ := re.ReplaceAllString(string(indexHtml), "/bundle.js?updated="+Utility.ToString(time.Now().Unix()))
		// save it back.
		ioutil.WriteFile(abosolutePath+"/index.html", []byte(indexHtml_), 0644)
	}

	return err
}

func (self *Globule) publishApplication(user, organization, path, name, domain, version, description, repositoryId, discoveryId string, keywords []string) error {
	publisherId := user
	if len(organization) > 0 {
		publisherId = organization
		if !self.isOrganizationMemeber(user, organization) {
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
		Type:         packagespb.PackageType_APPLICATION_TYPE,
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

	err := self.publishPackage(user, organization, discoveryId, repositoryId, "webapp", path, descriptor)

	// Create the permission...
	permissions := &rbacpb.Permissions{
		Allowed: []*rbacpb.Permission{
			//  Exemple of possible permission values.
			&rbacpb.Permission{
				Name:          "read", // member of the organization can publish the service.
				Applications:  []string{name},
				Accounts:      []string{},
				Groups:        []string{},
				Peers:         []string{},
				Organizations: []string{organization},
			},
		},
		Denied: []*rbacpb.Permission{},
		Owners: &rbacpb.Permission{
			Name:     "owner",
			Accounts: []string{user},
		},
	}

	// Set the permissions.
	err = self.setResourcePermissions("/"+name, permissions)

	if err != nil {
		return err
	}

	return nil
}

// Deloyed a web application to a globular node. Mostly use a develeopment time.
func (self *Globule) DeployApplication(stream adminpb.AdminService_DeployApplicationServer) error {

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
	for {
		msg, err := stream.Recv()
		if msg == nil {
			return errors.New("fail to run action adminpb.AdminService.DeployApplication!")
		}
		if len(msg.Name) > 0 {
			name = msg.Name
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

		if len(msg.Actions) > 0 {
			actions = msg.Actions
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

	// Now I will save the bundle into a file in the temp directory.
	path := os.TempDir() + "/" + Utility.RandomUUID()
	defer os.RemoveAll(path)

	err := ioutil.WriteFile(path, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	err = self.publishApplication(user, organization, path, name, domain, version, description, repositoryId, discoveryId, keywords)
	if err != nil {
		return err
	}

	// Read bytes and extract it in the current directory.
	r := bytes.NewReader(buffer.Bytes())
	err = self.installApplication(domain, name, organization, version, description, r, actions, keywords)
	if err != nil {
		return err
	}

	// Now I will send the update application event.
	eventClient, err := self.getEventHub()
	if err == nil {
		log.Println("-------> sent event ", "update_"+domain+"_"+name+"_evt", []byte(version))
		eventClient.Publish("update_"+strings.Split(domain, ":")[0]+"_"+name+"_evt", []byte(version))
	}

	return err
}

/** Create the super administrator in the db. **/
func (self *Globule) registerSa() error {

	configs := self.getServiceConfigByName("persistence.PersistenceService")
	if len(configs) == 0 {
		logger.Info("No persistence service was configure on that globule!")
		return errors.New("No persistence service was configure on that globule!")
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
	dataPath := self.data + string(os.PathSeparator) + "mongodb-data"

	if !Utility.Exists(dataPath) {
		// Kill mongo db server if the process already run...
		self.stopMongod()

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

		self.waitForMongo(60, false)

		// Now I will create a new user name sa and give it all admin write.
		createSaScript := fmt.Sprintf(
			`db=db.getSiblingDB('admin');db.createUser({ user: '%s', pwd: '%s', roles: ['userAdminAnyDatabase','userAdmin','readWrite','dbAdmin','clusterAdmin','readWriteAnyDatabase','dbAdminAnyDatabase']});`, "sa", self.RootPassword) // must be change...

		createSaCmd := exec.Command("mongo", "--eval", createSaScript)
		err = createSaCmd.Run()
		if err != nil {
			// remove the mongodb-data
			os.RemoveAll(dataPath)
			return err
		}
		self.stopMongod()
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

	self.waitForMongo(60, true)

	// Get the list of all services method.
	return self.registerMethods()
}

// Set the root password
func (self *Globule) SetRootPassword(ctx context.Context, rqst *adminpb.SetRootPasswordRequest) (*adminpb.SetRootPasswordResponse, error) {
	// Here I will set the root password.
	if self.RootPassword != rqst.OldPassword {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Wrong password given!")))
	}

	// Now I will update de sa password.
	self.RootPassword = rqst.NewPassword

	// Now update the sa password in mongo db.
	changeRootPasswordScript := fmt.Sprintf(
		"db=db.getSiblingDB('admin');db.changeUserPassword('%s','%s');", "sa", rqst.NewPassword)

	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// I will execute the script with the admin function.
	address, _ := self.getBackendAddress()
	if address == "0.0.0.0" {
		err = p.RunAdminCmd(context.Background(), "local_resource", "sa", rqst.OldPassword, changeRootPasswordScript)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	} else {
		p_, err := self.getPersistenceSaConnection()
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		err = p_.RunAdminCmd("local_resource", "sa", self.RootPassword, changeRootPasswordScript)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	token, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")
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
func (self *Globule) setPassword(accountId string, oldPassword string, newPassword string) error {

	// First of all I will get the user information from the database.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+accountId+`"},{"name":"`+accountId+`"} ]}`, ``)
	if err != nil {
		return err
	}

	account := values.(map[string]interface{})

	if len(oldPassword) == 0 {
		return errors.New("You must give your old password!")
	}

	// Test the old password.
	if oldPassword != account["password"] {
		if Utility.GenerateUUID(oldPassword) != account["password"] {
			return errors.New("Wrong password given!")
		}
	}

	// Now update the sa password in mongo db.
	name := account["name"].(string)
	name = strings.ReplaceAll(strings.ReplaceAll(name, ".", "_"), "@", "_")

	changePasswordScript := fmt.Sprintf(
		"db=db.getSiblingDB('admin');db.changeUserPassword('%s','%s');", name, newPassword)

	// I will execute the sript with the admin function.
	address, _ := self.getBackendAddress()
	if address == "0.0.0.0" {
		err = p.RunAdminCmd(context.Background(), "local_resource", "sa", self.RootPassword, changePasswordScript)
		if err != nil {
			return err
		}
	} else {
		p_, err := self.getPersistenceSaConnection()
		if err != nil {
			return err
		}
		err = p_.RunAdminCmd("local_resource", "sa", self.RootPassword, changePasswordScript)
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
func (self *Globule) SetPassword(ctx context.Context, rqst *adminpb.SetPasswordRequest) (*adminpb.SetPasswordResponse, error) {

	// First of all I will get the user information from the database.
	err := self.setPassword(rqst.AccountId, rqst.OldPassword, rqst.NewPassword)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	token, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")
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
func (self *Globule) SetEmail(ctx context.Context, rqst *adminpb.SetEmailRequest) (*adminpb.SetEmailResponse, error) {

	// Here I will set the root password.
	// First of all I will get the user information from the database.
	p, err := self.getPersistenceStore()
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
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Wrong email given!")))
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
	token, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")
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
func (self *Globule) SetRootEmail(ctx context.Context, rqst *adminpb.SetRootEmailRequest) (*adminpb.SetRootEmailResponse, error) {
	// Here I will set the root password.
	if self.AdminEmail != rqst.OldEmail {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Wrong email given, old email is "+rqst.OldEmail+" not "+self.AdminEmail+"!")))
	}

	// Now I will update de sa password.
	self.AdminEmail = rqst.NewEmail

	// save the configuration.
	self.saveConfig()

	// read the local token.
	token, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")
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
func (self *Globule) UploadServicePackage(stream adminpb.AdminService_UploadServicePackageServer) error {
	// The bundle will cantain the necessary information to install the service.
	path := os.TempDir() + string(os.PathSeparator) + Utility.RandomUUID()

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
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to upload service package!")))

		}

		if err == nil {
			if len(msg.Organization) > 0 {
				if !self.isOrganizationMemeber(msg.User, msg.Organization) {
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
func (self *Globule) publishPackage(user string, organization string, discovery string, repository string, platform string, path string, descriptor *packagespb.PackageDescriptor) error {

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
	permissions, err = self.getResourcePermissions(path_)
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
				Name:     "owner",
				Accounts: []string{user},
			},
		}

		// Set the permissions.
		err = self.setResourcePermissions(path_, permissions)
		if err != nil {
			return err
		}
	}

	// Test the permission before actualy publish the package.
	hasAccess, isDenied, err := self.validateAccess(user, rbacpb.SubjectType_ACCOUNT, "publish", path_)
	if !hasAccess || isDenied || err != nil {
		log.Println(err)
		return err
	}

	// Append the user into the list of owner if is not already part of it.
	if !Utility.Contains(permissions.Owners.Accounts, user) {
		permissions.Owners.Accounts = append(permissions.Owners.Accounts, user)
	}

	// Save the permissions.
	err = self.setResourcePermissions(path_, permissions)
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
	if self.discorveriesEventHub[discovery] == nil {
		client, err := event_client.NewEventService_Client(discovery, "event.EventService")
		if err != nil {
			log.Println("-->", err)
			return err
		}
		self.discorveriesEventHub[discovery] = client
	}

	eventId := descriptor.PublisherId + ":" + descriptor.Id
	if descriptor.Type == packagespb.PackageType_SERVICE_TYPE {
		eventId += ":SERVICE_PUBLISH_EVENT"
	} else if descriptor.Type == packagespb.PackageType_APPLICATION_TYPE {
		eventId += ":APPLICATION_PUBLISH_EVENT"
	}

	return self.discorveriesEventHub[discovery].Publish(eventId, []byte(data))
}

// Publish a service. The service must be install localy on the server.
func (self *Globule) PublishService(ctx context.Context, rqst *adminpb.PublishServiceRequest) (*adminpb.PublishServiceResponse, error) {

	// Make sure the user is part of the organization if one is given.
	publisherId := rqst.User
	if len(rqst.Organization) > 0 {
		publisherId = rqst.Organization
		if !self.isOrganizationMemeber(rqst.User, rqst.Organization) {
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

	err := self.publishPackage(rqst.User, rqst.Organization, rqst.DicorveryId, rqst.RepositoryId, rqst.Platform, rqst.Path, descriptor)

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
func (self *Globule) installService(descriptor *packagespb.PackageDescriptor) error {
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
			Utility.ExtractTarGz(r)

			// This is the directory path inside the archive.
			id := descriptor.PublisherId + "%" + descriptor.Name + "%" + descriptor.Version + "%" + descriptor.Id + "%" + globular.GetPlatform()

			// I will save the binairy in file...
			Utility.CreateDirIfNotExist(self.path + "/services")
			Utility.CopyDirContent(self.path+"/"+id, self.path+"/services")

			defer os.RemoveAll(self.path + "/" + id)
			path := self.path + "/services/" + descriptor.PublisherId + "/" + descriptor.Name + "/" + descriptor.Version + "/" + descriptor.Id
			configs, _ := Utility.FindFileByName(path, "config.json")

			if len(configs) == 0 {
				log.Println("No configuration file was found at at path ", path)
				return errors.New("No configuration file was found")
			}

			s := make(map[string]interface{}, 0)
			data, err := ioutil.ReadFile(configs[0])
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &s)
			if err != nil {
				return err
			}

			protos, _ := Utility.FindFileByName(self.path+"/services/"+descriptor.PublisherId+"/"+descriptor.Name+"/"+descriptor.Version, ".proto")
			if len(protos) == 0 {
				log.Println("No prototype file was found at at path ", self.path+"/services/"+descriptor.PublisherId+"/"+descriptor.Name+"/"+descriptor.Version)
				return errors.New("No configuration file was found")
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
			err = self.initService(s_)
			if err != nil {
				return err
			}

			// Here I will set the service method...
			self.setServiceMethods(s["Name"].(string), s["Proto"].(string))

			self.registerMethods()

			break
		} else {
			log.Println("fail to download error with error ", err)
			return err
		}
	}

	return nil

}

// Install/Update a service on globular instance.
func (self *Globule) InstallService(ctx context.Context, rqst *adminpb.InstallServiceRequest) (*adminpb.InstallServiceResponse, error) {
	log.Println("Try to install new service...")
	// Connect to the dicovery services
	var services_discovery *packages_client.PackagesDiscovery_Client
	var err error
	services_discovery, err = packages_client.NewPackagesDiscoveryService_Client(rqst.DicorveryId, "packages.PackageDiscovery")

	if services_discovery == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.DicorveryId)))
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

	err = self.installService(descriptor)
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
func (self *Globule) UninstallService(ctx context.Context, rqst *adminpb.UninstallServiceRequest) (*adminpb.UninstallServiceResponse, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// First of all I will stop the running service(s) instance.
	for _, s := range self.getServices() {
		// Stop the instance of the service.
		id, ok := s.Load("Id")
		if ok {
			if getStringVal(s, "PublisherId") == rqst.PublisherId && id == rqst.ServiceId && getStringVal(s, "Version") == rqst.Version {

				self.stopService(s)
				self.deleteService(id.(string))

				// Get the list of method to remove from the list of actions.
				toDelete := self.getServiceMethods(getStringVal(s, "Name"), getStringVal(s, "Proto"))
				methods := make([]string, 0)
				for i := 0; i < len(self.methods); i++ {
					if !Utility.Contains(toDelete, self.methods[i]) {
						methods = append(methods, self.methods[i])
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

				self.methods = methods
				self.registerMethods() // refresh methods.
			}
		}
	}

	// Now I will remove the service.
	// Service are located into the packagespb...
	path := self.path + string(os.PathSeparator) + "services" + string(os.PathSeparator) + rqst.PublisherId + string(os.PathSeparator) + rqst.ServiceId + string(os.PathSeparator) + rqst.Version

	// remove directory and sub-directory.
	err = os.RemoveAll(path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// save the config.
	self.saveConfig()

	return &adminpb.UninstallServiceResponse{
		Result: true,
	}, nil
}

/**
 * Retunr the path of config.json for a given services.
 */
func (self *Globule) getServiceConfigPath(s *sync.Map) string {

	path := getStringVal(s, "Path")
	index := strings.LastIndex(path, "/")
	if index == -1 {
		return ""
	}

	path = path[0:index] + "/config.json"
	return path
}

func (self *Globule) stopService(s *sync.Map) error {

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
	configPath := self.getServiceConfigPath(s)
	if len(configPath) > 0 {
		err := ioutil.WriteFile(configPath, []byte(jsonStr), 0644)
		if err != nil {
			return err
		}
	}

	// self.logServiceInfo(getStringVal(s, "Name"), time.Now().String()+"Service "+getStringVal(s, "Name")+" was stopped!")
	self.saveConfig()
	return nil
}

// Stop a service
func (self *Globule) StopService(ctx context.Context, rqst *adminpb.StopServiceRequest) (*adminpb.StopServiceResponse, error) {

	s := self.getService(rqst.ServiceId)
	if s != nil {
		err := self.stopService(s)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	} else {
		// Close all services with a given name.
		services := self.getServiceConfigByName(rqst.ServiceId)
		for i := 0; i < len(services); i++ {
			serviceId := services[i]["Id"].(string)
			s := self.getService(serviceId)
			if s == nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No service found with id "+serviceId)))
			}
			err := self.stopService(s)
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
func (self *Globule) StartService(ctx context.Context, rqst *adminpb.StartServiceRequest) (*adminpb.StartServiceResponse, error) {

	s := self.getService(rqst.ServiceId)
	proxy_pid := int64(-1)
	service_pid := int64(-1)

	if s == nil {
		services := self.getServiceConfigByName(rqst.ServiceId)
		for i := 0; i < len(services); i++ {
			id := services[i]["Id"].(string)
			s := self.getService(id)
			service_pid_, proxy_pid_, err := self.startService(s)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			proxy_pid = int64(proxy_pid_)
			service_pid = int64(service_pid_)
		}
	} else {
		service_pid_, proxy_pid_, err := self.startService(s)
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
func (self *Globule) RestartServices(ctx context.Context, rqst *adminpb.RestartServicesRequest) (*adminpb.RestartServicesResponse, error) {
	log.Println("restart service... ")
	self.restartServices()

	return &adminpb.RestartServicesResponse{}, nil
}

// That command is use to restart a new instance of the globular.
func rerunDetached() error {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	cmd := exec.Command(os.Args[0], []string{}...)
	cmd.Dir = cwd
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return err
	}
	cmd.Process.Release()
	log.Println(os.Args[0])
	return nil
}

func (self *Globule) restartServices() {
	if self.exit_ {
		return // already restarting I will ingnore the call.
	}

	// Stop all internal services
	self.stopInternalServices()

	// Stop all external services.
	self.stopServices()

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
func (self *Globule) startExternalApplication(serviceId string) (int, error) {

	if service, ok := self.ExternalApplications[serviceId]; !ok {
		return -1, errors.New("No external service found with name " + serviceId)
	} else {

		service.srv = exec.Command(service.Path, service.Args...)

		err := service.srv.Start()
		if err != nil {
			return -1, err
		}

		// save back the service in the map.
		self.ExternalApplications[serviceId] = service

		return service.srv.Process.Pid, nil
	}

}

// Stop external service.
func (self *Globule) stopExternalApplication(serviceId string) error {
	if _, ok := self.ExternalApplications[serviceId]; !ok {
		return errors.New("No external service found with name " + serviceId)
	}

	// if no command was created
	if self.ExternalApplications[serviceId].srv == nil {
		return nil
	}

	// if no process running
	if self.ExternalApplications[serviceId].srv.Process == nil {
		return nil
	}

	// kill the process.
	return self.ExternalApplications[serviceId].srv.Process.Signal(os.Interrupt)
}

// Kill process by id
func (self *Globule) KillProcess(ctx context.Context, rqst *adminpb.KillProcessRequest) (*adminpb.KillProcessResponse, error) {
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
func (self *Globule) KillProcesses(ctx context.Context, rqst *adminpb.KillProcessesRequest) (*adminpb.KillProcessesResponse, error) {
	err := Utility.KillProcessByName(rqst.Name)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.KillProcessesResponse{}, nil
}

// Return the list of process id with a given name.
func (self *Globule) GetPids(ctx context.Context, rqst *adminpb.GetPidsRequest) (*adminpb.GetPidsResponse, error) {
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
func (self *Globule) RegisterExternalApplication(ctx context.Context, rqst *adminpb.RegisterExternalApplicationRequest) (*adminpb.RegisterExternalApplicationResponse, error) {

	// Here I will get the command path.
	externalCmd := ExternalApplication{
		Id:   rqst.ServiceId,
		Path: rqst.Path,
		Args: rqst.Args,
	}

	self.ExternalApplications[externalCmd.Id] = externalCmd

	// save the config.
	self.saveConfig()

	// start the external service.
	pid, err := self.startExternalApplication(externalCmd.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.RegisterExternalApplicationResponse{
		ServicePid: int64(pid),
	}, nil
}

// Run an external command must be use with care.
func (self *Globule) RunCmd(ctx context.Context, rqst *adminpb.RunCmdRequest) (*adminpb.RunCmdResponse, error) {
	baseCmd := rqst.Cmd
	cmdArgs := rqst.Args
	isBlocking := rqst.Blocking
	pid := -1
	cmd := exec.Command(baseCmd, cmdArgs...)
	output := ""

	if isBlocking {
		out, err := cmd.Output()
		if cmd.Process != nil {
			pid = cmd.Process.Pid
		}
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		output = string(out)
	} else {
		var err error
		err = cmd.Start()
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		if cmd.Process != nil {
			pid = cmd.Process.Pid
		}
	}

	return &adminpb.RunCmdResponse{
		Result: string(output),
		Pid:    int32(pid),
	}, nil
}

// Set environement variable.
func (self *Globule) SetEnvironmentVariable(ctx context.Context, rqst *adminpb.SetEnvironmentVariableRequest) (*adminpb.SetEnvironmentVariableResponse, error) {
	err := Utility.SetEnvironmentVariable(rqst.Name, rqst.Value)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &adminpb.SetEnvironmentVariableResponse{}, nil
}

// Get environement variable.
func (self *Globule) GetEnvironmentVariable(ctx context.Context, rqst *adminpb.GetEnvironmentVariableRequest) (*adminpb.GetEnvironmentVariableResponse, error) {
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
func (self *Globule) UnsetEnvironmentVariable(ctx context.Context, rqst *adminpb.UnsetEnvironmentVariableRequest) (*adminpb.UnsetEnvironmentVariableResponse, error) {

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
func (self *Globule) InstallCertificates(ctx context.Context, rqst *adminpb.InstallCertificatesRequest) (*adminpb.InstallCertificatesResponse, error) {
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
