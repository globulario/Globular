package main

import (
	"context"

	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/globulario/Globular/services/golang/event/event_client"
	"github.com/globulario/Globular/services/golang/services/servicespb"
	"github.com/golang/protobuf/jsonpb"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"encoding/json"
	"os/exec"
	"reflect"

	"github.com/globulario/Globular/services/golang/lb/lbpb"

	"strconv"

	"net"

	"github.com/globulario/Globular/Interceptors"
	"github.com/globulario/Globular/services/golang/admin/adminpb"
	globular "github.com/globulario/Globular/services/golang/globular_service"
	"github.com/globulario/Globular/services/golang/ressource/ressourcepb"
	"github.com/globulario/Globular/services/golang/services/service_client"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (self *Globule) startAdminService() error {
	id := string(adminpb.File_services_proto_admin_proto.Services().Get(0).FullName())
	admin_server, err := self.startInternalService(id, adminpb.File_services_proto_admin_proto.Path(), self.AdminPort, self.AdminProxy, self.Protocol == "https", Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor) // must be accessible to all clients...
	if err == nil {
		self.inernalServices = append(self.inernalServices, admin_server)
		// First of all I will creat a listener.
		// Create the channel to listen on admin port.

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.AdminPort))
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
			self.saveConfig()
		}()

	}
	return err

}

func (self *Globule) getConfig() map[string]interface{} {

	config := make(map[string]interface{}, 0)
	config["Name"] = self.Name
	config["PortHttp"] = self.PortHttp
	config["PortHttps"] = self.PortHttps
	config["AdminPort"] = self.AdminPort
	config["AdminProxy"] = self.AdminProxy
	config["AdminEmail"] = self.AdminEmail
	config["AlternateDomains"] = self.AlternateDomains
	config["RessourcePort"] = self.RessourcePort
	config["RessourceProxy"] = self.RessourceProxy
	config["ServicesDiscoveryPort"] = self.ServicesDiscoveryPort
	config["ServicesDiscoveryProxy"] = self.ServicesDiscoveryProxy
	config["ServicesRepositoryPort"] = self.ServicesRepositoryPort
	config["ServicesRepositoryProxy"] = self.ServicesRepositoryProxy
	config["LoadBalancingServiceProxy"] = self.LoadBalancingServiceProxy
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
	config["CertificateAuthorityPort"] = self.CertificateAuthorityPort
	config["CertificateAuthorityProxy"] = self.CertificateAuthorityProxy
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

func (self *Globule) saveConfig() {
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
 */
func (self *Globule) HasRunningProcess(ctx context.Context, rqst *adminpb.HasRunningProcessRequest) (*adminpb.HasRunningProcessResponse, error) {
	ids, err := Utility.GetProcessIdsByName(rqst.Name)
	if err != nil {
		return nil, err
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
 * Return globular configuration.
 */
func (self *Globule) GetFullConfig(ctx context.Context, rqst *adminpb.GetConfigRequest) (*adminpb.GetConfigResponse, error) {

	str, err := Utility.ToJson(self.toMap())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.GetConfigResponse{
		Result: str,
	}, nil

}

// Return the configuration.
func (self *Globule) GetConfig(ctx context.Context, rqst *adminpb.GetConfigRequest) (*adminpb.GetConfigResponse, error) {

	config := self.getConfig()

	str, err := Utility.ToJson(config)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.GetConfigResponse{
		Result: str,
	}, nil
}

// return true if the configuation has change.
func (self *Globule) saveServiceConfig(config *sync.Map) bool {
	root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")
	root_ := string(root)[0:strings.Index(string(root), ":")]

	if !Utility.IsLocal(getStringVal(config, "Domain")) && root_ != self.path {
		return false
	}

	_, hasConfigPath := config.Load("configPath")
	if !hasConfigPath {
		config_ := self.getService(getStringVal(config, "Id"))
		if config_ != nil {
			config.Range(func(k, v interface{}) bool {
				config_.Store(k, v)
				return true
			})
			// save the globule configuration.
			self.saveConfig()
		}
		return false
	}

	// Here I will

	// set the domain of the service.
	config.Store("Domain", self.getDomain())

	// format the path's
	config.Store("Path", strings.ReplaceAll(getStringVal(config, "Path"), "\\", "/"))
	config.Store("Proto", strings.ReplaceAll(getStringVal(config, "Proto"), "\\", "/"))

	_, hasConfigPath = config.Load("hasConfigPath")
	if hasConfigPath {
		config.Store("configPath", strings.ReplaceAll(getStringVal(config, "configPath"), "\\", "/"))
	}

	// so here I will get the previous information...
	f, err := os.Open(getStringVal(config, "configPath"))

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
			if reflect.DeepEqual(config_, config__) {
				f.Close()
				// set back the path's info.
				return false
			}

			// sync the data/config file with the service file.
			jsonStr, _ := Utility.ToJson(config__)
			// here I will write the file
			err = ioutil.WriteFile(getStringVal(config, "configPath"), []byte(jsonStr), 0644)
			if err != nil {
				return false
			}
		}
	}
	f.Close()

	// Here I will get the list of service permission and set it...
	permissions, hasPermissions := config.Load("permissions")

	if hasPermissions {
		for i := 0; i < len(permissions.([]interface{})); i++ {
			permission := permissions.([]interface{})[i].(map[string]interface{})
			self.actionPermissions = append(self.actionPermissions, permission)
		}
	}

	return true
}

// Save a service configuration
func (self *Globule) SaveConfig(ctx context.Context, rqst *adminpb.SaveConfigRequest) (*adminpb.SaveConfigResponse, error) {
	// Save service...
	config := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(rqst.Config), &config)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// if the configuration is one of services...
	if config["Id"] != nil {
		srv := self.getService(config["Id"].(string))
		if srv != nil {
			setValues(srv, config)
			self.initService(srv)
		}

	} else if config["Services"] != nil {
		// Here I will save the configuration
		self.Name = config["Name"].(string)
		self.PortHttp = Utility.ToInt(config["PortHttp"].(float64))
		self.PortHttps = Utility.ToInt(config["PortHttps"].(float64))
		self.PortsRange = config["PortsRange"].(string)
		self.AdminEmail = config["AdminEmail"].(string)
		self.Country = config["Country"].(string)
		self.State = config["State"].(string)
		self.City = config["City"].(string)
		self.Organization = config["Organization"].(string)
		self.AdminPort = Utility.ToInt(config["AdminPort"].(float64))
		self.AdminProxy = Utility.ToInt(config["AdminProxy"].(float64))
		self.RessourcePort = Utility.ToInt(config["RessourcePort"].(float64))
		self.RessourceProxy = Utility.ToInt(config["RessourceProxy"].(float64))
		self.ServicesDiscoveryPort = Utility.ToInt(config["ServicesDiscoveryPort"].(float64))
		self.ServicesDiscoveryProxy = Utility.ToInt(config["ServicesDiscoveryProxy"].(float64))
		self.ServicesRepositoryPort = Utility.ToInt(config["ServicesRepositoryPort"].(float64))
		self.ServicesRepositoryProxy = Utility.ToInt(config["ServicesRepositoryProxy"].(float64))
		self.CertificateAuthorityPort = Utility.ToInt(config["CertificateAuthorityPort"].(float64))
		self.CertificateAuthorityProxy = Utility.ToInt(config["CertificateAuthorityProxy"].(float64))

		if config["DnsUpdateIpInfos"] != nil {
			self.DnsUpdateIpInfos = config["DnsUpdateIpInfos"].([]interface{})
		}

		if config["AlternateDomains"] != nil {
			self.AlternateDomains = config["AlternateDomains"].([]interface{})
		}

		self.Protocol = config["Protocol"].(string)
		self.Domain = config["Domain"].(string)
		self.CertExpirationDelay = Utility.ToInt(config["CertExpirationDelay"].(float64))

		// Save Discoveries.
		self.Discoveries = make([]string, 0)
		for i := 0; i < len(config["Discoveries"].([]interface{})); i++ {
			self.Discoveries = append(self.Discoveries, config["Discoveries"].([]interface{})[i].(string))
		}

		// Save DNS
		self.DNS = make([]string, 0)
		for i := 0; i < len(config["DNS"].([]interface{})); i++ {
			self.DNS = append(self.DNS, config["DNS"].([]interface{})[i].(string))
		}

		// That will save the services if they have changed.
		for id, s := range config["Services"].(map[string]interface{}) {
			// Attach the actual process and proxy process to the configuration object.
			s_ := self.getService(id)
			setValues(s_, s.(map[string]interface{}))
			self.initService(s_)
		}

		// save the application server.
		self.saveConfig()
	}

	// return the new configuration file...
	result, _ := Utility.ToJson(config)
	return &adminpb.SaveConfigResponse{
		Result: result,
	}, nil
}

// Deloyed a web application to a globular node.
func (self *Globule) DeployApplication(stream adminpb.AdminService_DeployApplicationServer) error {

	// The bundle will cantain the necessary information to install the service.
	var buffer bytes.Buffer

	var name string
	var domain string
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			stream.SendAndClose(&adminpb.DeployApplicationResponse{
				Result: true,
			})
			break
		} else if err != nil {
			return err
		} else {
			name = msg.Name
			domain = msg.Domain
			buffer.Write(msg.Data)
		}
	}

	// Before extract I will keep the archive.
	backupPath := self.webRoot + string(os.PathSeparator) + "_old" + string(os.PathSeparator) + name

	Utility.CreateIfNotExists(backupPath, 0644)
	backupPath += string(os.PathSeparator) + Utility.ToString(time.Now().Unix()) + ".tar.gz"

	err := ioutil.WriteFile(backupPath, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	// Read bytes and extract it in the current directory.
	r := bytes.NewReader(buffer.Bytes())

	// Before I will
	Utility.ExtractTarGz(r)

	// Copy the files to it final destination
	abosolutePath := self.webRoot
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
	Utility.CopyDirContent(Utility.GenerateUUID(name), abosolutePath)

	// remove temporary files.
	os.RemoveAll(Utility.GenerateUUID(name))

	// Now I will create the application database in the persistence store,
	// and the Application entry in the database.
	// That service made user of persistence service.
	store, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	count, err := store.Count(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+name+`"}`, "")
	application := make(map[string]interface{})
	application["_id"] = name
	application["password"] = Utility.GenerateUUID(name)
	application["path"] = "/" + name // The path must be the same as the application name.
	if len(domain) > 0 {
		if Utility.Exists(self.webRoot + "/" + domain) {
			application["path"] = "/" + domain + "/" + application["path"].(string)
		}
	}
	application["last_deployed"] = time.Now().Unix() // save it as unix time.

	// Here I will set the ressource to manage the applicaiton access permission.
	ctx := stream.Context()
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		user := strings.Join(md["user"], "")

		// Set user owner of the application directory.
		self.setRessourceOwner(user, "/"+name)

		path := strings.Join(md["path"], "")
		res := &ressourcepb.Ressource{
			Path:     path,
			Modified: time.Now().Unix(),
			Size:     int64(buffer.Len()),
			Name:     name,
		}
		self.setRessource(res)
		self.setRessourceOwner(user, path+"/"+name)

	}

	if err != nil || count == 0 {

		// create the application database.
		createApplicationUserDbScript := fmt.Sprintf(
			"db=db.getSiblingDB('%s_db');db.createCollection('application_data');db=db.getSiblingDB('admin');db.createUser({user: '%s', pwd: '%s',roles: [{ role: 'dbOwner', db: '%s_db' }]});",
			name, name, application["password"].(string), name)

		err = store.RunAdminCmd(context.Background(), "local_ressource", "sa", self.RootPassword, createApplicationUserDbScript)
		if err != nil {
			return err
		}

		application["creation_date"] = time.Now().Unix() // save it as unix time.

		_, err := store.InsertOne(context.Background(), "local_ressource", "local_ressource", "Applications", application, "")
		if err != nil {
			return err
		}

		p, err := self.getPersistenceSaConnection()

		err = p.CreateConnection(name+"_db", name+"_db", "0.0.0.0", 27017, 0, name, application["password"].(string), 5000, "", false)
		if err != nil {
			return err
		}

	} else {

		err := store.UpdateOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+name+`"}`, `{ "$set":{ "last_deployed":`+Utility.ToString(time.Now().Unix())+` }}`, "")
		if err != nil {
			return err
		}
	}

	// here is a little workaround to be sure the bundle.js file will not be cached in the brower...
	indexHtml, err := ioutil.ReadFile(abosolutePath + "/index.html")
	if err == nil {
		var re = regexp.MustCompile(`\/bundle\.js(\?updated=\d*)?`)
		indexHtml_ := re.ReplaceAllString(string(indexHtml), "/bundle.js?updated="+Utility.ToString(time.Now().Unix()))
		// save it back.
		ioutil.WriteFile(abosolutePath+"/index.html", []byte(indexHtml_), 0644)

	}

	return nil
}

/** Create the super administrator in the db. **/
func (self *Globule) registerSa() error {

	configs := self.getServiceConfigByName("persistence.PersistenceService")
	if len(configs) == 0 {
		return errors.New("No persistence service was configure on that globule!")
	}

	// Here I will test if mongo db exist on the server.
	existMongo := exec.Command("mongod", "--version")
	err := existMongo.Run()
	if err != nil {
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
	err = mongod.Start()
	if err != nil {
		return err
	}

	// wait 15 seconds that the server restart.
	self.waitForMongo(60, true)

	// Get the list of all services method.
	return self.registerMethods()
}

//Set the root password
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

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd(context.Background(), "local_ressource", "sa", rqst.OldPassword, changeRootPasswordScript)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.saveConfig()

	token, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")

	return &adminpb.SetRootPasswordResponse{
		Token: string(token),
	}, nil

}

func (self *Globule) setPassword(accountId string, oldPassword string, newPassword string) error {

	// First of all I will get the user information from the database.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+accountId+`"}`, ``)
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
	err = p.RunAdminCmd(context.Background(), "local_ressource", "sa", self.RootPassword, changePasswordScript)
	if err != nil {
		return err
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

	err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"name":"`+account["name"].(string)+`"}`, jsonStr, ``)
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

	token, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")
	return &adminpb.SetPasswordResponse{
		Token: string(token),
	}, nil

}

//Set the root password
func (self *Globule) SetEmail(ctx context.Context, rqst *adminpb.SetEmailRequest) (*adminpb.SetEmailResponse, error) {

	// Here I will set the root password.
	// First of all I will get the user information from the database.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, ``)
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

	err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"name":"`+account["name"].(string)+`"}`, jsonStr, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// read the token.
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
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Wrong email given!")))
	}

	// Now I will update de sa password.
	self.AdminEmail = rqst.NewEmail

	// save the configuration.
	self.saveConfig()

	// read the token.
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
		if err == io.EOF {
			// end of stream...
			stream.SendAndClose(&adminpb.UploadServicePackageResponse{
				Path: path,
			})
			break
		} else if err != nil {
			return err
		} else {
			fo.Write(msg.Data)
		}
	}

	return nil
}

// Publish a service. The service must be install localy on the server.
func (self *Globule) PublishService(ctx context.Context, rqst *adminpb.PublishServiceRequest) (*adminpb.PublishServiceResponse, error) {

	// Connect to the dicovery services
	services_discovery, err := service_client.NewServicesDiscoveryService_Client(rqst.DicorveryId, "services.ServiceDiscovery")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.DicorveryId)))
	}

	// Connect to the repository servicespb.
	services_repository, err := service_client.NewServicesRepositoryService_Client(rqst.RepositoryId, "services.ServiceRepository")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.RepositoryId)))
	}

	// Now I will upload the service to the repository...
	serviceDescriptor := &servicespb.ServiceDescriptor{
		Id:           rqst.ServiceId,
		Name:         rqst.ServiceName,
		PublisherId:  rqst.PublisherId,
		Version:      rqst.Version,
		Description:  rqst.Description,
		Keywords:     rqst.Keywords,
		Repositories: []string{rqst.RepositoryId},
	}

	err = services_discovery.PublishServiceDescriptor(serviceDescriptor)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Upload the service to the repository.
	err = services_repository.UploadBundle(rqst.DicorveryId, serviceDescriptor.Id, serviceDescriptor.PublisherId, rqst.Platform, rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	fi, err := os.Stat(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// So here I will send an plublish event...
	err = os.Remove(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will send a event to be sure all server will update...
	var marshaler jsonpb.Marshaler
	data, err := marshaler.MarshalToString(serviceDescriptor)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will send an event that the service has a new version...
	if self.discorveriesEventHub[rqst.DicorveryId] == nil {
		client, err := event_client.NewEventService_Client(rqst.DicorveryId, "event.EventService")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		self.discorveriesEventHub[rqst.DicorveryId] = client
	}

	self.discorveriesEventHub[rqst.DicorveryId].Publish(serviceDescriptor.PublisherId+":"+serviceDescriptor.Id+":SERVICE_PUBLISH_EVENT", []byte(data))

	// Here I will set the ressource to manage the applicaiton access permission.
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		user := strings.Join(md["user"], "")
		path := strings.Join(md["path"], "")
		res := &ressourcepb.Ressource{
			Path:     path,
			Modified: time.Now().Unix(),
			Size:     fi.Size(),
			Name:     rqst.ServiceId,
		}
		self.setRessource(res)
		self.setRessourceOwner(user, "/services/"+rqst.PublisherId)
	}

	return &adminpb.PublishServiceResponse{
		Result: true,
	}, nil
}

// Install/Update a service on globular instance.
func (self *Globule) installService(descriptor *servicespb.ServiceDescriptor) error {
	// repository must exist...
	log.Println("step 2: try to dowload service bundle")
	if len(descriptor.Repositories) == 0 {
		return errors.New("No service repository was found for service " + descriptor.Id)
	}

	for i := 0; i < len(descriptor.Repositories); i++ {

		services_repository, err := service_client.NewServicesRepositoryService_Client(descriptor.Repositories[i], "services.ServiceRepository")
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
			log.Println("--> service downloaded successfully")
			log.Println("--> try to install service to ", path)

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
			s["configPath"] = configs[0]
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
	var services_discovery *service_client.ServicesDiscovery_Client
	var err error
	services_discovery, err = service_client.NewServicesDiscoveryService_Client(rqst.DicorveryId, "services.ServiceDiscovery")

	if services_discovery == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.DicorveryId)))
	}

	descriptors, err := services_discovery.GetServiceDescriptor(rqst.ServiceId, rqst.PublisherId)
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
	log.Println("--> service descriptor receiced whit id ", descriptor.GetId())
	err = self.installService(descriptor)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

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

				self.stopService(id.(string))
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
						p.Delete(context.Background(), "local_ressource", "local_ressource", "ActionPermission", `{"action":"`+toDelete[i]+`"}`, "")

						// Delete it from Role.
						p.Update(context.Background(), "local_ressource", "local_ressource", "Roles", `{}`, `{"$pull":{"actions":"`+toDelete[i]+`"}}`, "")

						// Delete it from Application.
						p.Update(context.Background(), "local_ressource", "local_ressource", "Applications", `{}`, `{"$pull":{"actions":"`+toDelete[i]+`"}}`, "")

					}
				}

				self.methods = methods
				self.registerMethods() // refresh methods.
			}
		}
	}

	// Now I will remove the service.
	// Service are located into the servicespb...
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

func (self *Globule) stopService(serviceId string) error {

	s := self.getService(serviceId)
	if s == nil {
		return errors.New("No service found with id " + serviceId)
	}

	// Set keep alive to false...
	s.Store("KeepAlive", false)
	_, hasProcessPid := s.Load("Process")
	if !hasProcessPid {
		s.Store("Process", -1)
	}

	pid := getIntVal(s, "Process")
	if pid != -1 {
		err := Utility.TerminateProcess(pid, 0)
		if err != nil {
			log.Println("fail to teminate process ", pid)
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
	configPath := getStringVal(s, "configPath")
	if len(configPath) > 0 {
		// here I will write the file
		err := ioutil.WriteFile(configPath, []byte(jsonStr), 0644)
		if err != nil {
			return err
		}
	}

	self.logServiceInfo(getStringVal(s, "Name"), time.Now().String()+"Service "+getStringVal(s, "Name")+" was stopped!")

	// I will remove the service from the load balancer.
	self.lb_remove_candidate_info_channel <- &lbpb.ServerInfo{
		Id:     getStringVal(s, "Id"),
		Name:   getStringVal(s, "Name"),
		Domain: getStringVal(s, "Domain"),
		Port:   int32(getIntVal(s, "Port")),
	}

	return nil
}

// Stop a service
func (self *Globule) StopService(ctx context.Context, rqst *adminpb.StopServiceRequest) (*adminpb.StopServiceResponse, error) {
	err := self.stopService(rqst.ServiceId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &adminpb.StopServiceResponse{
		Result: true,
	}, nil
}

// Start a service
func (self *Globule) StartService(ctx context.Context, rqst *adminpb.StartServiceRequest) (*adminpb.StartServiceResponse, error) {

	s := self.getService(rqst.ServiceId)
	if s == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No service found with id "+rqst.ServiceId)))
	}

	service_pid, proxy_pid, err := self.startService(s)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &adminpb.StartServiceResponse{
		ProxyPid:   int64(proxy_pid),
		ServicePid: int64(service_pid),
	}, nil
}

// Restart all Services also the http(s)
func (self *Globule) RestartServices(ctx context.Context, rqst *adminpb.RestartServicesRequest) (*adminpb.RestartServicesResponse, error) {

	defer self.restartServices()

	return &adminpb.RestartServicesResponse{}, nil
}

func (self *Globule) restartServices() {
	self.exit <- struct{}{}

	// Stop all internal services
	self.stopInternalServices()

	// Stop all external services.
	self.stopServices()

	// Start services
	self.Serve()

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
