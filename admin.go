package main

import (
	"context"

	"io/ioutil"
	"log"
	"os"

	"bytes"

	"time"

	"io"

	"strings"

	"errors"
	"fmt"
	"regexp"

	"github.com/davecourtois/Globular/services"
	"github.com/davecourtois/Globular/services/servicespb"
	"github.com/golang/protobuf/jsonpb"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"encoding/json"
	"os/exec"
	"reflect"
	"runtime"

	"github.com/davecourtois/Globular/lb/lbpb"

	"strconv"

	"net"

	"os/signal"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/admin/adminpb"
	"github.com/davecourtois/Globular/ressource/ressourcepb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (self *Globule) startAdminService() error {
	admin_server, err := self.startInternalService(string(adminpb.File_admin_admin_proto.Services().Get(0).FullName()), adminpb.File_admin_admin_proto.Path(), self.AdminPort, self.AdminProxy, self.Protocol == "https", Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor) // must be accessible to all clients...
	if err == nil {
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
			go func() {
				// no web-rpc server.
				if err := admin_server.Serve(lis); err != nil {
					log.Println(err)
				}
			}()
			// Wait for signal to stop.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			Utility.KillProcessByName("mongod")
			Utility.KillProcessByName("prometheus")
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
	config["RessourcePort"] = self.RessourcePort
	config["RessourceProxy"] = self.RessourceProxy
	config["ServicesDiscoveryPort"] = self.ServicesDiscoveryPort
	config["ServicesDiscoveryProxy"] = self.ServicesDiscoveryProxy
	config["ServicesRepositoryPort"] = self.ServicesRepositoryPort
	config["ServicesRepositoryProxy"] = self.ServicesRepositoryProxy
	config["LoadBalancingServiceProxy"] = self.LoadBalancingServiceProxy
	config["SessionTimeout"] = self.SessionTimeout
	config["Discoveries"] = self.Discoveries
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

	// return the full service configuration.
	// Here I will give only the basic services informations and keep
	// all other infromation secret.
	config["Services"] = make(map[string]interface{}) //self.Services

	for name, service_config := range self.Services {
		s := make(map[string]interface{})
		s["Domain"] = service_config.(map[string]interface{})["Domain"]
		s["Port"] = service_config.(map[string]interface{})["Port"]
		s["Proxy"] = service_config.(map[string]interface{})["Proxy"]
		s["TLS"] = service_config.(map[string]interface{})["TLS"]
		s["Version"] = service_config.(map[string]interface{})["Version"]
		s["PublisherId"] = service_config.(map[string]interface{})["PublisherId"]
		s["KeepUpToDate"] = service_config.(map[string]interface{})["KeepUpToDate"]
		s["KeepAlive"] = service_config.(map[string]interface{})["KeepAlive"]
		s["State"] = service_config.(map[string]interface{})["State"]
		s["Id"] = name
		s["Name"] = service_config.(map[string]interface{})["Name"]
		s["CertFile"] = service_config.(map[string]interface{})["CertFile"]
		s["KeyFile"] = service_config.(map[string]interface{})["KeyFile"]
		s["CertAuthorityTrust"] = service_config.(map[string]interface{})["CertAuthorityTrust"]

		config["Services"].(map[string]interface{})[name] = s
	}

	return config

}

func (self *Globule) saveConfig() {
	// Here I will save the server attribute
	config, err := Utility.ToMap(self)
	if err == nil {
		services := make(map[string]interface{}, 0)
		if config["Services"] != nil {
			services = config["Services"].(map[string]interface{})
		}

		for _, service := range services {
			// remove running information...
			delete(service.(map[string]interface{}), "Process")
			delete(service.(map[string]interface{}), "ProxyProcess")
		}
		str, err := Utility.ToJson(config)
		if err == nil {
			ioutil.WriteFile(self.config+string(os.PathSeparator)+"config.json", []byte(str), 0644)
		} else {
			log.Panicln(err)
		}
	} else {
		log.Panicln(err)
	}
}

/**
 * Return globular configuration.
 */
func (self *Globule) GetFullConfig(ctx context.Context, rqst *adminpb.GetConfigRequest) (*adminpb.GetConfigResponse, error) {

	config, err := Utility.ToMap(self)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	services := config["Services"].(map[string]interface{})
	for _, service := range services {
		// remove running information...
		delete(service.(map[string]interface{}), "Process")
		delete(service.(map[string]interface{}), "ProxyProcess")

	}

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
func (self *Globule) saveServiceConfig(config map[string]interface{}) bool {
	root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")
	root_ := string(root)[0:strings.Index(string(root), ":")]

	if !Utility.IsLocal(config["Domain"].(string)) && root_ != self.path {
		return false
	}

	if config["configPath"] == nil {
		return false
	}

	// Here I will

	// set the domain of the service.
	config["Domain"] = self.getDomain()

	// get the config path.
	var process interface{}
	var proxyProcess interface{}

	process = config["Process"]
	proxyProcess = config["ProxyProcess"]

	// remove unused information...
	delete(config, "Process")
	delete(config, "ProxyProcess")

	// so here I will get the previous information...
	f, err := os.Open(config["configPath"].(string))
	if err == nil {
		b, err := ioutil.ReadAll(f)
		if err == nil {
			config_ := make(map[string]interface{})
			json.Unmarshal(b, &config_)
			if reflect.DeepEqual(config_, config) {
				f.Close()
				// set back the path's info.
				config["Process"] = process
				config["ProxyProcess"] = proxyProcess
				return false
			}
		}
	}
	f.Close()

	// set back internal infos...
	config["Process"] = process
	config["ProxyProcess"] = proxyProcess

	// sync the data/config file with the service file.
	jsonStr, _ := Utility.ToJson(config)

	// here I will write the file
	err = ioutil.WriteFile(config["configPath"].(string), []byte(jsonStr), 0644)
	if err != nil {
		return false
	}

	// Here I will get the list of service permission and set it...
	if config["Permissions"] != nil {
		for i := 0; i < len(config["Permissions"].([]interface{})); i++ {
			permission := config["Permissions"].([]interface{})[i].(map[string]interface{})
			self.actionPermissions = append(self.actionPermissions, permission)
		}
	}

	return true
}

// Save a service configuration
func (self *Globule) SaveConfig(ctx context.Context, rqst *adminpb.SaveConfigRequest) (*adminpb.SaveConfigResponse, error) {

	config := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(rqst.Config), &config)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// if the configuration is one of servicespb...
	if config["Id"] != nil {
		srv := self.Services[config["Id"].(string)]
		if srv != nil {
			// Attach the actual process and proxy process to the configuration object.
			config["Process"] = srv.(map[string]interface{})["Process"]
			config["ProxyProcess"] = srv.(map[string]interface{})["ProxyProcess"]
			self.initService(config)
		}
	} else if config["Services"] != nil {
		// Here I will save the configuration
		self.Name = config["Name"].(string)
		self.PortHttp = Utility.ToInt(config["PortHttp"].(float64))
		self.PortHttps = Utility.ToInt(config["PortHttps"].(float64))
		self.AdminEmail = config["AdminEmail"].(string)
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

		self.Protocol = config["Protocol"].(string)
		self.Domain = config["Domain"].(string)
		self.CertExpirationDelay = Utility.ToInt(config["CertExpirationDelay"].(float64))

		// That will save the services if they have changed.
		for n, s := range config["Services"].(map[string]interface{}) {
			// Attach the actual process and proxy process to the configuration object.
			s.(map[string]interface{})["Process"] = self.Services[n].(map[string]interface{})["Process"]
			s.(map[string]interface{})["ProxyProcess"] = self.Services[n].(map[string]interface{})["ProxyProcess"]
			self.initService(s.(map[string]interface{}))
		}

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
	abosolutePath := self.webRoot + string(os.PathSeparator) + name

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
	application["path"] = "/" + name                 // The path must be the same as the application name.
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
	err := mongod.Start()
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
	services_discovery, err := services.NewServicesDiscovery_Client(rqst.DicorveryId, "services_discovery")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.DicorveryId)))
	}

	// Connect to the repository servicespb.
	services_repository, err := services.NewServicesRepository_Client(rqst.RepositoryId, "services_repository")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.RepositoryId)))
	}

	// Now I will upload the service to the repository...
	serviceDescriptor := &servicespb.ServiceDescriptor{
		Id:           rqst.ServiceId,
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
	err = services_repository.UploadBundle(rqst.DicorveryId, serviceDescriptor.Id, serviceDescriptor.PublisherId, int32(rqst.Platform), rqst.Path)
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
	// eventHub, err := event_client.NewEvent_Client(rqst.DicorveryId, "event.EventService")
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
	if len(descriptor.Repositories) == 0 {
		return errors.New("No service repository was found for service " + descriptor.Id)
	}

	var platform servicespb.Platform
	// The first step will be to create the archive.
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			platform = servicespb.Platform_WIN64
		} else if runtime.GOARCH == "386" {
			platform = servicespb.Platform_WIN32
		}
	} else if runtime.GOOS == "linux" { // also can be specified to FreeBSD
		if runtime.GOARCH == "amd64" {
			platform = servicespb.Platform_LINUX64
		} else if runtime.GOARCH == "386" {
			platform = servicespb.Platform_LINUX32
		}
	} else if runtime.GOOS == "darwin" {
		/** TODO Deploy services on other platforme here... **/
	}

	for i := 0; i < len(descriptor.Repositories); i++ {
		services_repository, err := services.NewServicesRepository_Client(descriptor.Repositories[i], "services_repository")
		if err != nil {
			return err
		}

		bundle, err := services_repository.DownloadBundle(descriptor, platform)
		if err == nil {
			id := descriptor.PublisherId + "%" + descriptor.Id + "%" + descriptor.Version
			if platform == servicespb.Platform_LINUX32 {
				id += "%LINUX32"
			} else if platform == servicespb.Platform_LINUX64 {
				id += "%LINUX64"
			} else if platform == servicespb.Platform_WIN32 {
				id += "%WIN32"
			} else if platform == servicespb.Platform_WIN64 {
				id += "%WIN64"
			}

			// Create the file.
			r := bytes.NewReader(bundle.Binairies)
			Utility.ExtractTarGz(r)

			// I will save the binairy in file...
			dest := "services" + string(os.PathSeparator) + strings.ReplaceAll(id, "%", string(os.PathSeparator))
			Utility.CreateDirIfNotExist(dest)
			Utility.CopyDirContent(self.path+string(os.PathSeparator)+id, self.path+string(os.PathSeparator)+dest)

			// remove the file...
			os.RemoveAll(self.path + string(os.PathSeparator) + id)

			// I will repalce the service configuration with the new one...
			jsonStr, err := ioutil.ReadFile(self.path + string(os.PathSeparator) + dest + string(os.PathSeparator) + "config.json")
			if err != nil {
				return err
			}

			config := make(map[string]interface{})
			json.Unmarshal(jsonStr, &config)

			config["configPath"] = strings.ReplaceAll(string(os.PathSeparator)+dest+string(os.PathSeparator)+"config.json", string(os.PathSeparator), "/")

			// save the new paths...

			// Take the existing servicePath from the configuration.
			servicePath := string(os.PathSeparator) + dest
			servicePath += string(os.PathSeparator) + config["Path"].(string)[strings.LastIndex(config["Path"].(string), "/"):]
			config["Path"] = strings.ReplaceAll(servicePath, string(os.PathSeparator), "/")

			// Here I will try to find .proto file.
			files, err := ioutil.ReadDir(self.path + string(os.PathSeparator) + dest)
			if err != nil {
				return err
			}

			// proto file dosen't have the save name as the service itself.
			protoPath := string(os.PathSeparator) + dest
			for i := 0; i < len(files); i++ {
				f := files[i]
				if strings.HasSuffix(f.Name(), ".proto") {
					protoPath += string(os.PathSeparator) + f.Name()
					break
				}
			}

			config["Proto"] = strings.ReplaceAll(protoPath, string(os.PathSeparator), "/")

			// Set execute permission
			err = os.Chmod(config["Path"].(string), 0755)
			if err != nil {
				return err
			}

			// Set service tls true if the protocol is https
			config["TLS"] = self.Protocol == "https"

			// Set the id with the descriptor id.
			config["Id"] = strings.Replace(descriptor.Id, ".exe", "", -1)

			// initialyse the new service.
			err = self.initService(config)
			if err != nil {
				return err
			}

			self.saveConfig() // save the configuration with the newly install service...

			// Here I will set the service method...
			self.setServiceMethods(config["Name"].(string), self.path+string(os.PathSeparator)+config["Proto"].(string))
			self.registerMethods()

			break
		}
	}

	return nil

}

// Install/Update a service on globular instance.
func (self *Globule) InstallService(ctx context.Context, rqst *adminpb.InstallServiceRequest) (*adminpb.InstallServiceResponse, error) {

	// Connect to the dicovery services
	var services_discovery *services.ServicesDiscovery_Client
	var err error
	services_discovery, err = services.NewServicesDiscovery_Client(rqst.DicorveryId, "services_discovery")

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
	for id, service := range self.Services {
		// Stop the instance of the service.
		s := service.(map[string]interface{})

		if s["Id"] != nil {
			if s["PublisherId"].(string) == rqst.PublisherId && s["Id"].(string) == rqst.ServiceId && s["Version"].(string) == rqst.Version {
				self.stopService(s["Id"].(string))
				delete(self.Services, id)

				// Get the list of method to remove from the list of actions.
				toDelete := self.getServiceMethods(s["Name"].(string), self.path+string(os.PathListSeparator)+s["Proto"].(string))
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
	if (self.Services[serviceId]) == nil {
		return errors.New("no service with id " + serviceId + " is define on the server!")
	}
	s := self.Services[serviceId].(map[string]interface{})
	if s == nil {
		return errors.New("No service found with id " + serviceId)
	}

	// Set keep alive to false...
	s["KeepAlive"] = false

	if s["Process"] == nil {
		return errors.New("No process running")
	}

	if s["Process"].(*exec.Cmd).Process == nil {
		return errors.New("No process running")
	}

	err := s["Process"].(*exec.Cmd).Process.Kill()
	// time.Sleep(time.Second * 1)

	if err != nil {
		return err
	}

	if s["ProxyProcess"] != nil {
		err := s["ProxyProcess"].(*exec.Cmd).Process.Kill()
		// time.Sleep(time.Second * 1)
		if err != nil {
			return err
		}
	}

	s["State"] = "stopped"

	self.logServiceInfo(s["Name"].(string), time.Now().String()+"Service "+s["Name"].(string)+" was stopped!")

	// I will remove the service from the load balancer.
	self.lb_remove_candidate_info_channel <- &lbpb.ServerInfo{
		Id:     s["Id"].(string),
		Name:   s["Name"].(string),
		Domain: s["Domain"].(string),
		Port:   int32(s["Port"].(float64)),
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

	s := self.Services[rqst.ServiceId]
	if s == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No service found with id "+rqst.ServiceId)))
	}

	service_pid, proxy_pid, err := self.startService(s.(map[string]interface{}))
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
	return self.ExternalApplications[serviceId].srv.Process.Kill()
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
