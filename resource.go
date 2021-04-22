package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	//"reflect"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/davecourtois/Utility"
	"github.com/emicklei/proto"
	"github.com/globulario/Globular/Interceptors"
	"github.com/globulario/services/golang/persistence/persistence_store"
	"github.com/globulario/services/golang/resource/resourcepb"

	"github.com/globulario/services/golang/rbac/rbacpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func (self *Globule) startResourceService() error {
	id := string(resourcepb.File_proto_resource_proto.Services().Get(0).FullName())
	resource_server, port, err := self.startInternalService(id, resourcepb.File_proto_resource_proto.Path(), self.Protocol == "https", self.unaryResourceInterceptor, self.streamResourceInterceptor)
	if err == nil {
		self.inernalServices = append(self.inernalServices, resource_server)

		// Create the channel to listen on resource port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
		if err != nil {
			log.Fatalf("could not start resource service %s: %s", self.getDomain(), err)
		}

		resourcepb.RegisterResourceServiceServer(resource_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			// no web-rpc server.
			if err = resource_server.Serve(lis); err != nil {
				log.Println(err)
			}
			s := self.getService(id)
			pid := getIntVal(s, "ProxyProcess")
			Utility.TerminateProcess(pid, 0)
			s.Store("ProxyProcess", -1)
			self.setService(s)
		}()

		// In order to be able to give permission to a server
		// I must register it to the globule associated
		// with the base domain.

		// set the ip into the DNS servers.
		ticker_ := time.NewTicker(5 * time.Second)
		go func() {
			ip := Utility.MyIP()
			self.registerIpToDns()
			for {
				select {
				case <-ticker_.C:
					if ip != Utility.MyIP() {
						self.registerIpToDns()
					}
					self.removeExpiredSessions()
					// Get the list of video
					convertVideo()
				case <-self.exit:
					return
				}
			}
		}()
	}
	return err
}

/**
 * Return the list of method's for a given service, the path is the path of the
 * proto file.
 */
func (self *Globule) getServiceMethods(name string, path string) []string {
	methods := make([]string, 0)

	configs := self.getServiceConfigByName(name)

	if len(configs) == 0 {
		// Test for name with pattern _server
		configs = self.getServiceConfigByName(name + "_server")
		if len(configs) == 0 {
			return nil
		}
	}

	// here I will parse the service defintion file to extract the
	// service difinition.
	reader, _ := os.Open(path)
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	// Stack values from walking tree
	stack := make([]interface{}, 0)

	handlePackage := func(stack *[]interface{}) func(*proto.Package) {
		return func(p *proto.Package) {
			*stack = append(*stack, p)
		}
	}(&stack)

	handleService := func(stack *[]interface{}) func(*proto.Service) {
		return func(s *proto.Service) {
			*stack = append(*stack, s)
		}
	}(&stack)

	handleRpc := func(stack *[]interface{}) func(*proto.RPC) {
		return func(r *proto.RPC) {
			*stack = append(*stack, r)
		}
	}(&stack)

	// Walk this way
	proto.Walk(definition,
		proto.WithPackage(handlePackage),
		proto.WithService(handleService),
		proto.WithRPC(handleRpc))

	var packageName string
	var serviceName string
	var methodName string

	for len(stack) > 0 {
		var x interface{}
		x, stack = stack[0], stack[1:]
		switch v := x.(type) {
		case *proto.Package:
			packageName = v.Name
		case *proto.Service:
			serviceName = v.Name
		case *proto.RPC:
			methodName = v.Name
			path := "/" + packageName + "." + serviceName + "/" + methodName
			// So here I will register the method into the backend.
			methods = append(methods, path)
		}
	}

	return methods
}

// Remove expired session...
func (self *Globule) removeExpiredSession(accountId string, expiredAt int64) error {
	log.Println("----------> session need to be remove: ", accountId)
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	session := make(map[string]interface{}, 0)
	session["_id"] = accountId
	session["state"] = 1
	session["lastStateTime"] = time.Unix(expiredAt, 0).UTC().Format("2006-01-02T15:04:05-0700")
	jsonStr, err := Utility.ToJson(session)
	if err != nil {
		return err
	}

	dbName := strings.ReplaceAll(accountId, ".", "_")
	dbName = strings.ReplaceAll(dbName, "@", "_")

	err = p.ReplaceOne(context.Background(), "local_resource", dbName+"_db", "Sessions", `{"_id":"`+accountId+`"}`, jsonStr, "")
	if err != nil {
		return err
	}

	evtHub, err := self.getEventHub()
	if err != nil {
		return err
	}

	// Publish the event on the network...
	evtHub.Publish("session_state_"+accountId+"_change_event", []byte(jsonStr))

	// Now I will remove the token...
	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+accountId+`"}`, "")
	if err != nil {
		return err
	}

	return nil

}

func (self *Globule) removeExpiredSessions() error {
	// I will iterate over the list token and close expired session...

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	tokens, err := p.Find(context.Background(), "local_resource", "local_resource", "Tokens", `{}`, ``)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i].(map[string]interface{})
		accountId := token["_id"].(string)
		expiredAt := int64(Utility.ToInt(token["expireAt"]))
		if expiredAt < time.Now().Unix() {
			self.removeExpiredSession(accountId, expiredAt)
		}
	}

	return nil

}

func (self *Globule) setServiceMethods(name string, path string) {

	// here I will parse the service defintion file to extract the
	// service difinition.
	reader, _ := os.Open(path)
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	// Stack values from walking tree
	stack := make([]interface{}, 0)

	handlePackage := func(stack *[]interface{}) func(*proto.Package) {
		return func(p *proto.Package) {
			*stack = append(*stack, p)
		}
	}(&stack)

	handleService := func(stack *[]interface{}) func(*proto.Service) {
		return func(s *proto.Service) {
			*stack = append(*stack, s)
		}
	}(&stack)

	handleRpc := func(stack *[]interface{}) func(*proto.RPC) {
		return func(r *proto.RPC) {
			*stack = append(*stack, r)
		}
	}(&stack)

	// Walk this way
	proto.Walk(definition,
		proto.WithPackage(handlePackage),
		proto.WithService(handleService),
		proto.WithRPC(handleRpc))

	var packageName string
	var serviceName string
	var methodName string

	for len(stack) > 0 {
		var x interface{}
		x, stack = stack[0], stack[1:]
		switch v := x.(type) {
		case *proto.Package:
			packageName = v.Name
		case *proto.Service:
			serviceName = v.Name
		case *proto.RPC:
			methodName = v.Name
			path := "/" + packageName + "." + serviceName + "/" + methodName
			//log.Println(path)
			// So here I will register the method into the backend.
			self.methods = append(self.methods, path)
		}
	}
}

// Method must be register in order to be assign to role.
func (self *Globule) registerMethods() error {
	// Here I will create the sa role if it dosen't exist.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// Here I will persit the sa role if it dosent already exist.
	admin := make(map[string]interface{})
	admin["_id"] = "sa"
	admin["name"] = "sa"
	admin["actions"] = self.methods
	jsonStr, _ := Utility.ToJson(admin)

	// I will set the role actions...
	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"sa"}`, jsonStr, `[{"upsert":true}]`)
	if err != nil {
		return err
	}

	log.Println("role sa with was updated!")

	// I will also create the guest role, the basic one
	count, err := p.Count(context.Background(), "local_resource", "local_resource", "Roles", `{ "_id":"guest"}`, "")
	guest := make(map[string]interface{})
	if err != nil {
		return err
	} else if count == 0 {
		log.Println("need to create roles guest...")
		guest["_id"] = "guest"
		guest["name"] = "guest"
		guest["actions"] = []string{
			"/admin.AdminService/GetConfig",
			"/resource.ResourceService/RegisterAccount",
			"/resource.ResourceService/GetAccounts",
			"/resource.ResourceService/AccountExist",
			"/resource.ResourceService/Authenticate",
			"/resource.ResourceService/RefreshToken",
			"/resource.ResourceService/GetPermissions",
			"/resource.ResourceService/GetAllFilesInfo",
			"/resource.ResourceService/GetAllApplicationsInfo",
			"/resource.ResourceService/GetResourceOwners",
			"/rbac.RbacService/ValidateAction",
			"/rbac.RbacService/ValidateAccess",
			"/event.EventService/Subscribe",
			"/event.EventService/UnSubscribe", "/event.EventService/OnEvent",
			"/event.EventService/Quit",
			"/event.EventService/Publish",
			"/packages.PackageDiscovery/FindServices",
			"/packages.PackageDiscovery/GetServiceDescriptor",
			"/packages.PackageDiscovery/GetServicesDescriptor",
			"/packages.PackageRepository/downloadBundle",
			"/persistence.PersistenceService/Find",
			"/persistence.PersistenceService/FindOne",
			"/persistence.PersistenceService/Count",
			"/resource.ResourceService/GetAllActions"}

		_, err := p.InsertOne(context.Background(), "local_resource", "local_resource", "Roles", guest, "")
		if err != nil {
			return err
		}
		log.Println("role guest was created!")
	}

	return nil
}

/** Set action permissions.
When gRPC service methode are called they must validate the ressource pass in parameters.
So each service is reponsible to give access permissions requirement.
*/
func (self *Globule) setActionResourcesPermissions(permissions map[string]interface{}) error {
	// So here I will keep values in local storage.cap()
	data, err := json.Marshal(permissions["resources"])
	if err != nil {
		return err
	}
	self.permissions.SetItem(permissions["action"].(string), data)
	return nil
}

// Retreive the ressource infos from the database.
func (self *Globule) getActionResourcesPermissions(action string) ([]*rbacpb.ResourceInfos, error) {
	data, err := self.permissions.GetItem(action)
	infos_ := make([]*rbacpb.ResourceInfos, 0)
	if err != nil {
		if strings.Index(err.Error(), "leveldb: not found") == -1 {
			return nil, err
		} else {
			// no infos_ found...
			return infos_, nil
		}
	}
	infos := make([]interface{}, 0)
	err = json.Unmarshal(data, &infos)

	for i := 0; i < len(infos); i++ {
		info := infos[i].(map[string]interface{})
		infos_ = append(infos_, &rbacpb.ResourceInfos{Index: int32(Utility.ToInt(info["index"])), Permission: info["permission"].(string)})
	}

	return infos_, err
}

func (self *Globule) createApplicationConnection() error {
	log.Println("create applications connections")
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	store, err := self.getPersistenceStore()
	applications, _ := store.Find(context.Background(), "local_resource", "local_resource", "Applications", "{}", "")
	if err != nil {
		log.Println("fail to get applications informations ", err)
		return err
	}
	address, port := self.getBackendAddress()
	for i := 0; i < len(applications); i++ {
		application := applications[i].(map[string]interface{})

		err = p.CreateConnection(application["_id"].(string)+"_db", application["_id"].(string)+"_db", address, float64(port), 0, application["_id"].(string), application["password"].(string), 5000, "", false)
		if err != nil {
			log.Println("** fail to create connection for application ", application["_id"].(string), err)

		}

		// Here I will se ownership to files in the directories.
		err = self.addResourceOwner("/applications/"+application["_id"].(string), application["_id"].(string), rbacpb.SubjectType_APPLICATION)
		if err != nil {
			return err
		}
		log.Println("-> connection created for application: ", application["_id"].(string))
	}

	return nil
}

func (self *Globule) deleteReference(p persistence_store.Store, refId, targetId, targetField, targetCollection string) error {

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", targetCollection, `{"_id":"`+targetId+`"}`, ``)
	if err != nil {
		return err
	}

	target := values.(map[string]interface{})

	if target[targetField] == nil {
		return errors.New("No field named " + targetField + " was found in object with id " + targetId + "!")
	}

	references := []interface{}(target[targetField].(primitive.A))
	references_ := make([]interface{}, 0)
	for j := 0; j < len(references); j++ {
		if references[j].(map[string]interface{})["$id"] != refId {
			references_ = append(references_, references[j])
		}
	}

	target[targetField] = references_
	jsonStr := serialyseObject(target)

	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", targetCollection, `{"_id":"`+targetId+`"}`, jsonStr, ``)
	if err != nil {
		return err
	}

	return nil
}

func (self *Globule) createReference(p persistence_store.Store, id, sourceCollection, field, targetId, targetCollection string) error {
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", sourceCollection, `{"_id":"`+id+`"}`, ``)
	if err != nil {
		return err
	}

	log.Println("create reference ", id, " source ", sourceCollection, " field ", field, " field ", targetId, " target ", targetCollection)
	source := values.(map[string]interface{})
	references := make([]interface{}, 0)
	if source[field] != nil {
		references = []interface{}(source[field].(primitive.A))
	}

	for j := 0; j < len(references); j++ {
		if references[j].(map[string]interface{})["$id"] == targetId {
			return errors.New(" named " + targetId + " aleready exist in  " + field + "!")
		}
	}

	// append the account.
	source[field] = append(references, map[string]interface{}{"$ref": targetCollection, "$id": targetId, "$db": "local_resource"})
	jsonStr := serialyseObject(source)

	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", sourceCollection, `{"_id":"`+id+`"}`, jsonStr, ``)
	if err != nil {
		return err
	}

	return nil
}

func (self *Globule) createCrossReferences(sourceId, sourceCollection, sourceField, targetId, targetCollection, targetField string) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	err = self.createReference(p, targetId, targetCollection, targetField, sourceId, sourceCollection)
	if err != nil {
		//return err
		log.Println(err)
	}

	err = self.createReference(p, sourceId, sourceCollection, sourceField, targetId, targetCollection)

	return err

}

// TODO update the account and group if they already exist...
func (self *Globule) synchronizeLdap(syncInfo map[string]interface{}) error {

	if self.LdapSyncInfos[syncInfo["LdapServiceId"].(string)] == nil {
		self.LdapSyncInfos[syncInfo["LdapServiceId"].(string)] = make([]interface{}, 0)
		self.LdapSyncInfos[syncInfo["LdapServiceId"].(string)] = append(self.LdapSyncInfos[syncInfo["LdapServiceId"].(string)].([]interface{}), syncInfo)
	} else {
		syncInfos := self.LdapSyncInfos[syncInfo["LdapServiceId"].(string)].([]interface{})
		exist := false
		for i := 0; i < len(syncInfos); i++ {
			if syncInfos[i].(map[string]interface{})["LdapServiceId"] == syncInfo["LdapServiceId"] {
				if syncInfos[i].(map[string]interface{})["ConnectionId"] == syncInfo["ConnectionId"] {
					// set the connection info.
					syncInfos[i] = syncInfo
					exist = true
					// save the config.
					self.saveConfig()
					break
				}
			}
		}

		if !exist {
			self.LdapSyncInfos[syncInfo["LdapServiceId"].(string)] = append(self.LdapSyncInfos[syncInfo["LdapServiceId"].(string)].([]interface{}), syncInfo)
			// save the config.
			self.saveConfig()
		}
	}

	ldap_, err := self.getLdapClient()
	if err != nil {
		return err
	}

	connectionId := syncInfo["ConnectionId"].(string)
	groupSyncInfos := syncInfo["GroupSyncInfos"].(map[string]interface{})

	groupsInfo, err := ldap_.Search(connectionId, groupSyncInfos["Base"].(string), groupSyncInfos["Query"].(string), []string{groupSyncInfos["Id"].(string), "distinguishedName"})
	if err != nil {
		log.Println("-----> ", len(groupsInfo), "found!")
		return err
	}

	// Print group info.
	for i := 0; i < len(groupsInfo); i++ {
		name := groupsInfo[i][0].([]interface{})[0].(string)
		id := Utility.GenerateUUID(groupsInfo[i][1].([]interface{})[0].(string))
		self.createGroup(id, name, []string{}) // The group member will be set latter in that function.
		log.Println("synchronize group ", name)
	}

	// Synchronize account and user info...
	userSyncInfos := syncInfo["UserSyncInfos"].(map[string]interface{})
	accountsInfo, err := ldap_.Search(connectionId, userSyncInfos["Base"].(string), userSyncInfos["Query"].(string), []string{userSyncInfos["Id"].(string), userSyncInfos["Email"].(string), "distinguishedName", "memberOf"})
	if err != nil {
		return err
	}

	for i := 0; i < len(accountsInfo); i++ {
		// Print the list of account...
		// I will not set the password...
		name := strings.ToLower(accountsInfo[i][0].([]interface{})[0].(string))
		log.Println("synchronize account ", name)

		if len(accountsInfo[i][1].([]interface{})) > 0 {
			email := strings.ToLower(accountsInfo[i][1].([]interface{})[0].(string))

			if len(email) > 0 {
				id := Utility.GenerateUUID(strings.ToLower(accountsInfo[i][2].([]interface{})[0].(string)))

				if len(id) > 0 {
					groups := make([]string, 0)
					// Here I will set the user groups.
					if len(accountsInfo[i][3].([]interface{})) > 0 {
						for j := 0; j < len(accountsInfo[i][3].([]interface{})); j++ {
							groups = append(groups, Utility.GenerateUUID(accountsInfo[i][3].([]interface{})[j].(string)))
						}
					}

					// Try to create account...
					err := self.registerAccount(id, name, email, id, []string{}, []string{}, []string{"guest"}, groups)
					if err != nil {
						log.Println("---> error ", err)
						// TODO if the account already exist simply udate it values.
					}
				}
			} else {
				return errors.New("account " + strings.ToLower(accountsInfo[i][2].([]interface{})[0].(string)) + " has no email configured! ")
			}
		}
	}
	return nil
}

/** Append new LDAP synchronization informations. **/
func (self *Globule) SynchronizeLdap(ctx context.Context, rqst *resourcepb.SynchronizeLdapRqst) (*resourcepb.SynchronizeLdapRsp, error) {

	if rqst.SyncInfo == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No LDAP sync infos was given!")))
	}

	/** Ldap user informations **/
	if rqst.SyncInfo.UserSyncInfos == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No LDAP sync users infos was given!")))
	}

	/** Ldap group informations **/
	if rqst.SyncInfo.GroupSyncInfos == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No LDAP sync groups infos was given!")))
	}

	syncInfo, err := Utility.ToMap(rqst.SyncInfo)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = self.synchronizeLdap(syncInfo)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.SynchronizeLdapRsp{
		Result: true,
	}, nil
}

func (self *Globule) registerAccount(id string, name string, email string, password string, organizations []string, contacts []string, roles []string, groups []string) error {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// first of all the Persistence service must be active.
	count, err := p.Count(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+id+`"},{"name":"`+id+`"} ]}`, "")
	if err != nil {
		return err
	}

	// one account already exist for the name.
	if count == 1 {
		return errors.New("account with name " + name + " already exist!")
	}

	// set the account object and set it basic roles.
	account := make(map[string]interface{})
	account["_id"] = id
	account["name"] = name
	account["email"] = email
	account["password"] = Utility.GenerateUUID(password) // hide the password...
	// List of aggregation.
	account["roles"] = make([]interface{}, 0)
	account["groups"] = make([]interface{}, 0)
	account["organizations"] = make([]interface{}, 0)

	// append guest role if not already exist.
	if !Utility.Contains(roles, "guest") {
		roles = append(roles, "guest")
	}

	// Here I will insert the account in the database.
	_, err = p.InsertOne(context.Background(), "local_resource", "local_resource", "Accounts", account, "")

	// replace @ and . by _
	name = strings.ReplaceAll(strings.ReplaceAll(name, "@", "_"), ".", "_")

	// Each account will have their own database and a use that can read and write
	// into it.
	// Here I will wrote the script for mongoDB...
	createUserScript := fmt.Sprintf("db=db.getSiblingDB('%s_db');db.createCollection('user_data');db=db.getSiblingDB('admin');db.createUser({user: '%s', pwd: '%s',roles: [{ role: 'dbOwner', db: '%s_db' }]});", name, name, password, name)

	p_, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// I will execute the sript with the admin function.
	address, port := self.getBackendAddress()
	if address == "0.0.0.0" {
		err = p.RunAdminCmd(context.Background(), "local_resource", "sa", self.RootPassword, createUserScript)
		if err != nil {
			return err
		}
	} else {
		p_, err := self.getPersistenceSaConnection()
		if err != nil {
			return err
		}
		log.Println("----> try to create account ", account)
		log.Println(createUserScript)
		err = p_.RunAdminCmd("local_resource", "sa", self.RootPassword, createUserScript)
		if err != nil {
			return err
		}
	}

	err = p_.CreateConnection(name+"_db", name+"_db", address, float64(port), 0, name, password, 5000, "", false)
	if err != nil {
		return errors.New("No persistence service are available to store resource information.")
	}

	// Now I will set the reference for
	// Contact...
	for i := 0; i < len(contacts); i++ {
		self.addAccountContact(id, contacts[i])
	}

	// Organizations
	for i := 0; i < len(organizations); i++ {
		self.createCrossReferences(organizations[i], "Organizations", "accounts", id, "Accounts", "organizations")
	}

	// Roles
	for i := 0; i < len(roles); i++ {
		self.createCrossReferences(roles[i], "Roles", "members", id, "Accounts", "roles")
	}

	// Groups
	for i := 0; i < len(groups); i++ {
		self.createCrossReferences(groups[i], "Groups", "members", id, "Accounts", "groups")
	}

	return nil

}

/* Register a new Account */
func (self *Globule) RegisterAccount(ctx context.Context, rqst *resourcepb.RegisterAccountRqst) (*resourcepb.RegisterAccountRsp, error) {
	if rqst.ConfirmPassword != rqst.Account.Password {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to confirm your password!")))

	}

	if rqst.Account == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No account information was given!")))

	}

	err := self.registerAccount(rqst.Account.Name, rqst.Account.Name, rqst.Account.Email, rqst.Account.Password, rqst.Account.Organizations, rqst.Account.Contacts, rqst.Account.Roles, rqst.Account.Groups)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Generate a token to identify the user.
	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, rqst.Account.Id, rqst.Account.Name, rqst.Account.Email)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	id, _, _, expireAt, _ := Interceptors.ValidateToken(tokenString)
	_, err = p.InsertOne(context.Background(), "local_resource", "local_resource", "Tokens", map[string]interface{}{"_id": id, "expireAt": Utility.ToString(expireAt)}, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Now I will
	return &resourcepb.RegisterAccountRsp{
		Result: tokenString, // Return the token string.
	}, nil
}

//* Return the list accounts *
func (self *Globule) GetAccounts(rqst *resourcepb.GetAccountsRqst, stream resourcepb.ResourceService_GetAccountsServer) error {
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	query := rqst.Query
	if len(query) == 0 {
		query = "{}"
	}

	accounts, err := p.Find(context.Background(), "local_resource", "local_resource", "Accounts", query, ``)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// No I will stream the result over the networks.
	maxSize := 100
	values := make([]*resourcepb.Account, 0)

	for i := 0; i < len(accounts); i++ {
		account := accounts[i].(map[string]interface{})
		a := &resourcepb.Account{Id: account["_id"].(string), Name: account["name"].(string), Email: account["email"].(string)}

		if account["groups"] != nil {
			groups := []interface{}(account["groups"].(primitive.A))
			if groups != nil {
				for i := 0; i < len(groups); i++ {
					groupId := groups[i].(map[string]interface{})["$id"].(string)
					a.Groups = append(a.Groups, groupId)
				}
			}
		}

		if account["roles"] != nil {
			roles := []interface{}(account["roles"].(primitive.A))
			if roles != nil {
				for i := 0; i < len(roles); i++ {
					roleId := roles[i].(map[string]interface{})["$id"].(string)
					a.Roles = append(a.Roles, roleId)
				}
			}
		}

		if account["organizations"] != nil {
			organizations := []interface{}(account["organizations"].(primitive.A))
			if organizations != nil {
				for i := 0; i < len(organizations); i++ {
					organizationId := organizations[i].(map[string]interface{})["$id"].(string)
					a.Organizations = append(a.Organizations, organizationId)
				}
			}
		}

		if account["contacts"] != nil {
			contacts := []interface{}(account["contacts"].(primitive.A))
			if contacts != nil {
				for i := 0; i < len(contacts); i++ {
					contactId := contacts[i].(map[string]interface{})["$id"].(string)
					a.Contacts = append(a.Contacts, contactId)
				}
			}
		}

		values = append(values, a)

		if len(values) >= maxSize {
			err := stream.Send(
				&resourcepb.GetAccountsRsp{
					Accounts: values,
				},
			)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			values = make([]*resourcepb.Account, 0)
		}
	}

	// Send reminding values.
	err = stream.Send(
		&resourcepb.GetAccountsRsp{
			Accounts: values,
		},
	)

	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return nil
}

func (self *Globule) addAccountContact(accountId string, contactId string) error {

	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	err = self.createReference(p, accountId, "Accounts", "contacts", contactId, "Accounts")
	if err != nil {
		return err
	}

	return nil
}

//* Add contact to a given account *
func (self *Globule) AddAccountContact(ctx context.Context, rqst *resourcepb.AddAccountContactRqst) (*resourcepb.AddAccountContactRsp, error) {

	err := self.addAccountContact(rqst.AccountId, rqst.ContactId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.AddAccountContactRsp{Result: true}, nil
}

//* Remove a contact from a given account *
func (self *Globule) RemoveAccountContact(ctx context.Context, rqst *resourcepb.RemoveAccountContactRqst) (*resourcepb.RemoveAccountContactRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// That service made user of persistence service.
	err = self.deleteReference(p, rqst.AccountId, rqst.ContactId, "contacts", "Accounts")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.RemoveAccountContactRsp{Result: true}, nil
}

func (self *Globule) AccountExist(ctx context.Context, rqst *resourcepb.AccountExistRqst) (*resourcepb.AccountExistRsp, error) {
	var exist bool
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}
	// Test with the _id
	accountId := rqst.Id
	count, _ := p.Count(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+accountId+`"},{"name":"`+accountId+`"} ]}`, "")
	if count > 0 {
		exist = true
	}

	// Test with the name
	if !exist {
		count, _ := p.Count(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+accountId+`"},{"name":"`+accountId+`"} ]}`, "")
		if count > 0 {
			exist = true
		}
	}

	// Test with the email.
	if !exist {
		count, _ := p.Count(context.Background(), "local_resource", "local_resource", "Accounts", `{"email":"`+rqst.Id+`"}`, "")
		if count > 0 {
			exist = true
		}
	}
	if exist {
		return &resourcepb.AccountExistRsp{
			Result: true,
		}, nil
	}

	return nil, status.Errorf(
		codes.Internal,
		Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Account with id name or email '"+rqst.Id+"' dosent exist!")))

}

//* Authenticate a account by it name or email.
// That function test if the password is the correct one for a given user
// if it is a token is generate and that token will be use by other service
// to validate permission over the requested resourcepb.
func (self *Globule) Authenticate(ctx context.Context, rqst *resourcepb.AuthenticateRqst) (*resourcepb.AuthenticateRsp, error) {

	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// in case of sa user.(admin)
	if (rqst.Password == self.RootPassword && rqst.Name == "sa") || (rqst.Password == self.RootPassword && rqst.Name == self.AdminEmail) {
		// Generate a token to identify the user.
		tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, "sa", "sa", self.AdminEmail)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		id, _, _, expireAt, _ := Interceptors.ValidateToken(tokenString)
		err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+id+`"}`, `{"_id":"`+id+`","expireAt":`+Utility.ToString(expireAt)+`}`, `[{"upsert":true}]`)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		// self.setTokenByAddress(ctx, tokenString)

		/** Return the token only **/
		return &resourcepb.AuthenticateRsp{
			Token: tokenString,
		}, nil
	}

	values, err := p.Find(context.Background(), "local_resource", "local_resource", "Accounts", `{"name":"`+rqst.Name+`"}`, "")
	if len(values) == 0 {
		values, err = p.Find(context.Background(), "local_resource", "local_resource", "Accounts", `{"email":"`+rqst.Name+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	if len(values) == 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No account found for entry "+rqst.Name)))
	}

	if values[0].(map[string]interface{})["password"].(string) != Utility.GenerateUUID(rqst.Password) {

		ldap_, err := self.getLdapClient()
		if err != nil {
			return nil, err
		}

		// Here I will try to made use of ldap if there is a service configure.ldap
		err = ldap_.Authenticate("", values[0].(map[string]interface{})["name"].(string), rqst.Password)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		// Set the password whit
		err = self.setPassword(values[0].(map[string]interface{})["_id"].(string), values[0].(map[string]interface{})["password"].(string), rqst.Password)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	user := values[0].(map[string]interface{})["_id"].(string)
	// Generate a token to identify the user.
	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, user, values[0].(map[string]interface{})["name"].(string), values[0].(map[string]interface{})["email"].(string))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Create the user file directory.
	path := "/users/" + values[0].(map[string]interface{})["name"].(string)
	Utility.CreateDirIfNotExist(self.data + "/files" + path)

	self.addResourceOwner(path, user, rbacpb.SubjectType_ACCOUNT)

	name_ := values[0].(map[string]interface{})["name"].(string)
	name_ = strings.ReplaceAll(strings.ReplaceAll(name_, ".", "_"), "@", "_")

	p_, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Open the user database connection.
	address, port := self.getBackendAddress()
	err = p_.CreateConnection(name_+"_db", name_+"_db", address, float64(port), 0, name_, rqst.Password, 5000, "", false)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No persistence service are available to store resource information.")))
	}

	// save the newly create token into the database.
	id, _, _, expireAt, _ := Interceptors.ValidateToken(tokenString)
	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+id+`"}`, `{"_id":"`+id+`","expireAt":`+Utility.ToString(expireAt)+`}`, `[{"upsert":true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// save the token for the default session time.
	//self.setTokenByAddress(ctx, tokenString)

	// Here I got the token I will now put it in the cache.
	return &resourcepb.AuthenticateRsp{
		Token: tokenString,
	}, nil
}

/**
 * Refresh token get a new token.
 */
func (self *Globule) RefreshToken(ctx context.Context, rqst *resourcepb.RefreshTokenRqst) (*resourcepb.RefreshTokenRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// first of all I will validate the current token.
	id, name, email, expireAt, err := Interceptors.ValidateToken(rqst.Token)

	if err != nil {
		if !strings.HasPrefix(err.Error(), "token is expired") {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	// If the token is older than seven day without being refresh then I retrun an error.
	if time.Unix(expireAt, 0).Before(time.Now().AddDate(0, 0, -7)) {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The token cannot be refresh after 7 day")))
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+id+`"}`, `[{"Projection":{"expireAt":1}}]`)
	if err == nil && values != nil {
		lastTokenInfo := values.(map[string]interface{})
		savedTokenExpireAt := time.Unix(int64(lastTokenInfo["expireAt"].(int32)), 0)
		log.Println("already existing token expire at ", savedTokenExpireAt.String())
		log.Println("newly created token expire at ", time.Unix(expireAt, 0).String())
		// That mean a newer token was already refresh.
		if time.Unix(expireAt, 0).Before(savedTokenExpireAt) {
			err := errors.New("That token cannot not be refresh because a newer one already exist. You need to re-authenticate in order to get a new token.")
			log.Println(err)
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, id, name, email)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// get back the new expireAt
	id, name, _, expireAt, _ = Interceptors.ValidateToken(tokenString)

	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+id+`"}`, `{"_id":"`+id+`","expireAt":`+Utility.ToString(expireAt)+`}`, `[{"upsert":true}]`)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	log.Println("token was refresh with success for entity named ", name, "!!")

	// return the token string.
	return &resourcepb.RefreshTokenRsp{
		Token: tokenString,
	}, nil
}

// Test if account is a member of organisation.
func (self *Globule) isOrganizationMemeber(account string, organization string) bool {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return false
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+account+`"},{"name":"`+account+`"} ]}`, ``)
	if err != nil {
		return false
	}

	account_ := values.(map[string]interface{})
	if account_["organizations"] != nil {
		organizations := []interface{}(account_["organizations"].(primitive.A))
		for i := 0; i < len(organizations); i++ {
			organizationId := organizations[i].(map[string]interface{})["$id"].(string)
			if organization == organizationId {
				return true
			}
		}
	}

	return false

}

//* Delete an account *
func (self *Globule) DeleteAccount(ctx context.Context, rqst *resourcepb.DeleteAccountRqst) (*resourcepb.DeleteAccountRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}
	accountId := rqst.Id
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+accountId+`"},{"name":"`+accountId+`"} ]}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	account := values.(map[string]interface{})

	// Remove references.
	if account["organizations"] != nil {
		organizations := []interface{}(account["organizations"].(primitive.A))
		for i := 0; i < len(organizations); i++ {
			organizationId := organizations[i].(map[string]interface{})["$id"].(string)
			self.deleteReference(p, rqst.Id, organizationId, "accounts", "Accounts")
		}
	}

	if account["groups"] != nil {
		groups := []interface{}(account["groups"].(primitive.A))
		for i := 0; i < len(groups); i++ {
			groupId := groups[i].(map[string]interface{})["$id"].(string)
			self.deleteReference(p, rqst.Id, groupId, "members", "Accounts")
		}
	}

	if account["roles"] != nil {
		roles := []interface{}(account["roles"].(primitive.A))
		for i := 0; i < len(roles); i++ {
			roleId := roles[i].(map[string]interface{})["$id"].(string)
			self.deleteReference(p, rqst.Id, roleId, "members", "Accounts")
		}
	}

	// Try to delete the account...
	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+accountId+`"},{"name":"`+accountId+`"} ]}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Delete permissions
	// TODO delete account permissions

	// Delete the token.
	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+rqst.Id+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	name := account["name"].(string)
	name = strings.ReplaceAll(strings.ReplaceAll(name, ".", "_"), "@", "_")

	// Here I will drop the db user.
	dropUserScript := fmt.Sprintf(
		`db=db.getSiblingDB('admin');db.dropUser('%s', {w: 'majority', wtimeout: 4000})`,
		name)

	// I will execute the sript with the admin function.
	address, _ := self.getBackendAddress()
	if address == "0.0.0.0" {
		err = p.RunAdminCmd(context.Background(), "local_resource", "sa", self.RootPassword, dropUserScript)
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
		err = p_.RunAdminCmd("local_resource", "sa", self.RootPassword, dropUserScript)
		if err != nil {
			log.Println(err)
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	// Remove the user database.
	err = p.DeleteDatabase(context.Background(), "local_resource", name+"_db")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	p_, _ := self.getPersistenceSaConnection()
	err = p_.DeleteConnection(name + "_db")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteAccountRsp{
		Result: rqst.Id,
	}, nil
}

func (self *Globule) createRole(id string, name string, actions []string) error {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	_, err = p.FindOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+id+`"}`, ``)
	if err == nil {
		return errors.New("Role named " + name + "already exist!")
	}

	// Here will create the new role.
	role := make(map[string]interface{}, 0)
	role["_id"] = id
	role["name"] = name
	role["actions"] = actions

	_, err = p.InsertOne(context.Background(), "local_resource", "local_resource", "Roles", role, "")
	if err != nil {
		return err
	}

	return nil
}

//* Create a role with given action list *
func (self *Globule) CreateRole(ctx context.Context, rqst *resourcepb.CreateRoleRqst) (*resourcepb.CreateRoleRsp, error) {
	// That service made user of persistence service.
	err := self.createRole(rqst.Role.Id, rqst.Role.Name, rqst.Role.Actions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Now I will set the reference for

	// members...
	for i := 0; i < len(rqst.Role.Members); i++ {
		self.createCrossReferences(rqst.Role.Members[i], "Accounts", "roles", rqst.Role.GetId(), "Roles", "members")
	}

	// Organizations
	for i := 0; i < len(rqst.Role.Organizations); i++ {
		self.createCrossReferences(rqst.Role.Organizations[i], "Organizations", "roles", rqst.Role.GetId(), "Roles", "organizations")
	}

	return &resourcepb.CreateRoleRsp{Result: true}, nil
}

func (self *Globule) GetRoles(rqst *resourcepb.GetRolesRqst, stream resourcepb.ResourceService_GetRolesServer) error {
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	query := rqst.Query
	if len(query) == 0 {
		query = "{}"
	}

	roles, err := p.Find(context.Background(), "local_resource", "local_resource", "Roles", query, ``)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// No I will stream the result over the networks.
	maxSize := 100
	values := make([]*resourcepb.Role, 0)

	for i := 0; i < len(roles); i++ {
		role := roles[i].(map[string]interface{})
		r := &resourcepb.Role{Id: role["_id"].(string), Name: role["name"].(string), Actions: make([]string, 0)}

		if role["actions"] != nil {
			actions := []interface{}(role["actions"].(primitive.A))
			if actions != nil {
				for i := 0; i < len(actions); i++ {
					r.Actions = append(r.Actions, actions[i].(string))
				}
			}
		}

		if role["organizations"] != nil {
			organizations := []interface{}(role["organizations"].(primitive.A))
			if organizations != nil {
				for i := 0; i < len(organizations); i++ {
					organizationId := organizations[i].(map[string]interface{})["$id"].(string)
					r.Organizations = append(r.Organizations, organizationId)
				}
			}
		}

		if role["members"] != nil {
			members := []interface{}(role["members"].(primitive.A))
			if members != nil {
				for i := 0; i < len(members); i++ {
					memberId := members[i].(map[string]interface{})["$id"].(string)
					r.Members = append(r.Members, memberId)
				}
			}
		}

		values = append(values, r)

		if len(values) >= maxSize {
			err := stream.Send(
				&resourcepb.GetRolesRsp{
					Roles: values,
				},
			)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			values = make([]*resourcepb.Role, 0)
		}
	}

	// Send reminding values.
	err = stream.Send(
		&resourcepb.GetRolesRsp{
			Roles: values,
		},
	)

	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return nil
}

// That function is necessary to serialyse reference and kept field orders
func serialyseObject(obj map[string]interface{}) string {
	// Here I will save the role.
	jsonStr, _ := Utility.ToJson(obj)
	jsonStr = strings.ReplaceAll(jsonStr, `"$ref"`, `"__a__"`)
	jsonStr = strings.ReplaceAll(jsonStr, `"$id"`, `"__b__"`)
	jsonStr = strings.ReplaceAll(jsonStr, `"$db"`, `"__c__"`)

	obj_ := make(map[string]interface{}, 0)

	json.Unmarshal([]byte(jsonStr), &obj_)
	jsonStr, _ = Utility.ToJson(obj_)
	jsonStr = strings.ReplaceAll(jsonStr, `"__a__"`, `"$ref"`)
	jsonStr = strings.ReplaceAll(jsonStr, `"__b__"`, `"$id"`)
	jsonStr = strings.ReplaceAll(jsonStr, `"__c__"`, `"$db"`)

	return jsonStr
}

//* Delete a role with a given id *
func (self *Globule) DeleteRole(ctx context.Context, rqst *resourcepb.DeleteRoleRqst) (*resourcepb.DeleteRoleRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Remove references
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
	role := values.(map[string]interface{})
	roleId := role["_id"].(string)

	// Remove it from the accounts
	if role["members"] != nil {
		accounts := []interface{}(role["members"].(primitive.A))
		for i := 0; i < len(accounts); i++ {
			accountId := accounts[i].(map[string]interface{})["$id"].(string)
			self.deleteReference(p, accountId, roleId, "roles", "Accounts")
		}
	}

	// I will remove it from organizations...
	if role["organizations"] != nil {
		organizations := []interface{}(role["organizations"].(primitive.A))
		for i := 0; i < len(organizations); i++ {
			organizationId := organizations[i].(map[string]interface{})["$id"].(string)
			self.deleteReference(p, rqst.RoleId, organizationId, "roles", "Roles")
		}
	}

	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// TODO delete role permissions

	return &resourcepb.DeleteRoleRsp{Result: true}, nil
}

//* Append an action to existing role. *
func (self *Globule) AddRoleActions(ctx context.Context, rqst *resourcepb.AddRoleActionsRqst) (*resourcepb.AddRoleActionsRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	role := values.(map[string]interface{})

	needSave := false
	if role["actions"] == nil {
		role["actions"] = rqst.Actions
		needSave = true
	} else {
		actions := []interface{}(role["actions"].(primitive.A))
		for j := 0; j < len(rqst.Actions); j++ {
			exist := false
			for i := 0; i < len(actions); i++ {
				if actions[i].(string) == rqst.Actions[j] {
					exist = true
					break
				}
			}
			if !exist {
				// append only if not already there.
				actions = append(actions, rqst.Actions[j])
				needSave = true
			}
		}
		role["actions"] = actions
	}

	if needSave {

		jsonStr, _ := json.Marshal(role)
		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.AddRoleActionsRsp{Result: true}, nil
}

//* Remove an action to existing role. *
func (self *Globule) RemoveRoleAction(ctx context.Context, rqst *resourcepb.RemoveRoleActionRqst) (*resourcepb.RemoveRoleActionRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	role := values.(map[string]interface{})

	needSave := false
	if role["actions"] == nil {
		role["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		actions := make([]interface{}, 0)
		actions_ := []interface{}(role["actions"].(primitive.A))
		for i := 0; i < len(actions_); i++ {
			if actions_[i].(string) == rqst.Action {
				exist = true
			} else {
				actions = append(actions, actions_[i])
			}
		}
		if exist {
			role["actions"] = actions
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Role named "+rqst.RoleId+"not contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		// jsonStr, _ := json.Marshal(role)
		jsonStr := serialyseObject(role)

		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.RemoveRoleActionRsp{Result: true}, nil
}

//* Add role to a given account *
func (self *Globule) AddAccountRole(ctx context.Context, rqst *resourcepb.AddAccountRoleRqst) (*resourcepb.AddAccountRoleRsp, error) {
	// That service made user of persistence service.
	err := self.createCrossReferences(rqst.RoleId, "Roles", "members", rqst.AccountId, "Accounts", "roles")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.AddAccountRoleRsp{Result: true}, nil
}

//* Remove a role from a given account *
func (self *Globule) RemoveAccountRole(ctx context.Context, rqst *resourcepb.RemoveAccountRoleRqst) (*resourcepb.RemoveAccountRoleRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	accountId := rqst.AccountId
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+accountId+`"},{"name":"`+accountId+`"} ]}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No account named "+accountId+" exist!")))
	}

	account := values.(map[string]interface{})

	// Now I will test if the account already contain the role.
	if account["roles"] != nil {
		roles := make([]interface{}, 0)
		roles_ := []interface{}(account["roles"].(primitive.A))
		needSave := false
		for j := 0; j < len(roles_); j++ {
			if roles_[j].(map[string]interface{})["$id"] == rqst.RoleId {
				needSave = true
			} else {
				roles = append(roles, roles_[j])
			}
		}

		if needSave == false {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Account named "+rqst.AccountId+" does not contain role "+rqst.RoleId+"!")))
		}

		// append the newly created role.
		account["roles"] = roles
		jsonStr := serialyseObject(account)

		err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"$or":[{"_id":"`+rqst.AccountId+`"},{"name":"`+rqst.AccountId+`"} ]}`, jsonStr, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

	}
	return &resourcepb.RemoveAccountRoleRsp{Result: true}, nil
}

func (self *Globule) deleteApplication(applicationId string) error {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+applicationId+`"}`, ``)
	if err != nil {
		return err
	}

	application := values.(map[string]interface{})

	// I will remove it from organization...
	if application["organizations"] != nil {
		organizations := []interface{}(application["organizations"].(primitive.A))

		for i := 0; i < len(organizations); i++ {
			organizationId := organizations[i].(map[string]interface{})["$id"].(string)
			self.deleteReference(p, applicationId, organizationId, "applications", "Applications")
		}
	}

	// I will remove the directory.
	err = os.RemoveAll(application["path"].(string))
	if err != nil {
		return err
	}

	// Now I will remove the database create for the application.
	err = p.DeleteDatabase(context.Background(), "local_resource", applicationId+"_db")
	if err != nil {
		return err
	}

	// Finaly I will remove the entry in  the table.
	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+applicationId+`"}`, "")
	if err != nil {
		return err
	}

	// Delete permissions
	err = p.Delete(context.Background(), "local_resource", "local_resource", "Permissions", `{"owner":"`+applicationId+`"}`, "")
	if err != nil {
		return err
	}

	// Drop the application user.
	// Here I will drop the db user.
	dropUserScript := fmt.Sprintf(
		`db=db.getSiblingDB('admin');db.dropUser('%s', {w: 'majority', wtimeout: 4000})`,
		applicationId)

	// I will execute the sript with the admin function.
	address, _ := self.getBackendAddress()
	if address == "0.0.0.0" {
		err = p.RunAdminCmd(context.Background(), "local_resource", "sa", self.RootPassword, dropUserScript)
		if err != nil {
			return err
		}
	} else {
		p_, err := self.getPersistenceSaConnection()
		if err != nil {
			return err
		}
		err = p_.RunAdminCmd("local_resource", "sa", self.RootPassword, dropUserScript)
		if err != nil {
			return err
		}
	}

	// Now I will remove the application directory from the webroot...
	log.Println("remove directory ", self.webRoot+"/"+applicationId)
	os.RemoveAll(self.webRoot + "/" + applicationId)

	return nil

}

//* Delete an application from the server. *
func (self *Globule) DeleteApplication(ctx context.Context, rqst *resourcepb.DeleteApplicationRqst) (*resourcepb.DeleteApplicationRsp, error) {

	// That service made user of persistence service.
	err := self.deleteApplication(rqst.ApplicationId)
	return nil, status.Errorf(
		codes.Internal,
		Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))

	// TODO delete dir permission associate with the application.

	return &resourcepb.DeleteApplicationRsp{
		Result: true,
	}, nil
}

//* Append an action to existing application. *
func (self *Globule) AddApplicationActions(ctx context.Context, rqst *resourcepb.AddApplicationActionsRqst) (*resourcepb.AddApplicationActionsRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, ``)
	if err != nil {

		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	application := values.(map[string]interface{})
	needSave := false
	if application["actions"] == nil {
		application["actions"] = rqst.Actions
		needSave = true
	} else {
		application["actions"] = []interface{}(application["actions"].(primitive.A))
		for j := 0; j < len(rqst.Actions); j++ {
			exist := false
			for i := 0; i < len(application["actions"].([]interface{})); i++ {
				if application["actions"].([]interface{})[i].(string) == rqst.Actions[j] {
					exist = true
					break
				}
				if !exist {
					application["actions"] = append(application["actions"].([]interface{}), rqst.Actions[j])
					needSave = true
				}
			}
		}

	}

	if needSave {
		jsonStr := serialyseObject(application)
		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.AddApplicationActionsRsp{Result: true}, nil
}

//* Remove an action to existing application. *
func (self *Globule) RemoveApplicationAction(ctx context.Context, rqst *resourcepb.RemoveApplicationActionRqst) (*resourcepb.RemoveApplicationActionRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	application := values.(map[string]interface{})

	needSave := false
	if application["actions"] == nil {
		application["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		actions := make([]interface{}, 0)
		application["actions"] = []interface{}(application["actions"].(primitive.A))
		for i := 0; i < len(application["actions"].([]interface{})); i++ {
			if application["actions"].([]interface{})[i].(string) == rqst.Action {
				exist = true
			} else {
				actions = append(actions, application["actions"].([]interface{})[i])
			}
		}
		if exist {
			application["actions"] = actions
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Application named "+rqst.ApplicationId+" not contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr := serialyseObject(application)
		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.RemoveApplicationActionRsp{Result: true}, nil
}

/**
 * Return the list of all actions avalaible on the server.
 */
func (self *Globule) GetAllActions(ctx context.Context, rqst *resourcepb.GetAllActionsRqst) (*resourcepb.GetAllActionsRsp, error) {
	return &resourcepb.GetAllActionsRsp{Actions: self.methods}, nil
}

//////////////////////////// Loggin info ///////////////////////////////////////
func (self *Globule) validateActionRequest(rqst interface{}, method string, subject string, subjectType rbacpb.SubjectType) (bool, error) {
	hasAccess := false
	infos, err := self.getActionResourcesPermissions(method)
	if err != nil {
		infos = make([]*rbacpb.ResourceInfos, 0)
	} else {
		// Here I will get the params...
		val, _ := Utility.CallMethod(rqst, "ProtoReflect", []interface{}{})
		rqst_ := val.(protoreflect.Message)
		if rqst_.Descriptor().Fields().Len() > 0 {
			for i := 0; i < len(infos); i++ {
				// Get the path value from retreive infos.
				param := rqst_.Descriptor().Fields().Get(Utility.ToInt(infos[i].Index))
				path, _ := Utility.CallMethod(rqst, "Get"+strings.ToUpper(string(param.Name())[0:1])+string(param.Name())[1:], []interface{}{})
				infos[i].Path = path.(string)
			}
		}
	}

	hasAccess, err = self.validateAction(method, subject, rbacpb.SubjectType_ACCOUNT, infos)
	if err != nil {
		return false, err
	}

	return hasAccess, nil
}

// unaryInterceptor calls authenticateClient with current context
func (self *Globule) unaryResourceInterceptor(ctx context.Context, rqst interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	method := info.FullMethod

	// The token and the application id.
	var token string
	var application string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
	}

	hasAccess := false

	// If the address is local I will give the permission.
	/* Test only
	peer_, _ := peer.FromContext(ctx)
	address := peer_.Addr.String()
	address = address[0:strings.Index(address, ":")]
	if Utility.IsLocal(address) {
		hasAccess = true
	}*/

	var err error
	// Here some method are accessible by default.
	if method == "/resource.ResourceService/GetAllActions" ||
		method == "/resource.ResourceService/RegisterAccount" ||
		method == "/resource.ResourceService/GetAccounts" ||
		method == "/resource.ResourceService/RegisterPeer" ||
		method == "/resource.ResourceService/GetPeers" ||
		method == "/resource.ResourceService/AccountExist" ||
		method == "/resource.ResourceService/Authenticate" ||
		method == "/resource.ResourceService/RefreshToken" ||
		method == "/resource.ResourceService/GetAllFilesInfo" ||
		method == "/resource.ResourceService/GetAllApplicationsInfo" ||
		method == "/resource.ResourceService/ValidateToken" ||
		method == "/rbac.RbacService/ValidateAction" ||
		method == "/rbac.RbacService/ValidateAccess" ||
		method == "/rbac.RbacService/GetResourcePermissions" ||
		method == "/rbac.RbacService/GetResourcePermission" ||
		method == "/log.LogService/Log" ||
		method == "/log.LogService/DeleteLog" ||
		method == "/log.LogService/GetLog" ||
		method == "/log.LogService/ClearAllLog" {
		hasAccess = true
	}
	var clientId string

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		var expiredAt int64
		var err error

		/*clientId*/
		clientId, _, _, expiredAt, err = Interceptors.ValidateToken(token)
		if err != nil {
			return nil, err
		}

		if expiredAt < time.Now().Unix() {
			return nil, errors.New("The token is expired!")
		}

		hasAccess = clientId == "sa"
		if !hasAccess {
			hasAccess, _ = self.validateActionRequest(rqst, method, clientId, rbacpb.SubjectType_ACCOUNT)
		}
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		// TODO validate rpc method access
		hasAccess, _ = self.validateActionRequest(rqst, method, application, rbacpb.SubjectType_APPLICATION)
	}

	if !hasAccess {

		err := errors.New("Permission denied to execute method " + method + " user:" + clientId + " application:" + application)
		self.logInfo(application, method, token, err)

		return nil, err
	}

	// Execute the action.
	result, err := handler(ctx, rqst)

	return result, err

}

// Stream interceptor.
func (self *Globule) streamResourceInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	err := handler(srv, stream)
	if err != nil {
		return err
	}

	return nil
}

///////////////////////  resource management. /////////////////

//* Validate a token *
func (self *Globule) ValidateToken(ctx context.Context, rqst *resourcepb.ValidateTokenRqst) (*resourcepb.ValidateTokenRsp, error) {
	clientId, _, _, expireAt, err := Interceptors.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &resourcepb.ValidateTokenRsp{
		ClientId: clientId,
		Expired:  expireAt,
	}, nil
}

// authenticateAgent check the client credentials
func (self *Globule) authenticateClient(ctx context.Context) (string, string, int64, error) {
	var userId string
	var applicationId string
	var expired int64
	var err error

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		applicationId = strings.Join(md["application"], "")
		token := strings.Join(md["token"], "")
		// In that case no token was given...
		if len(token) > 0 {
			userId, _, _, expired, err = Interceptors.ValidateToken(token)
		}
		return applicationId, userId, expired, err
	}
	return "", "", 0, fmt.Errorf("missing credentials")
}

func (self *Globule) GetAllApplicationsInfo(ctx context.Context, rqst *resourcepb.GetAllApplicationsInfoRqst) (*resourcepb.GetAllApplicationsInfoRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// So here I will get the list of retreived permission.
	values, err := p.Find(context.Background(), "local_resource", "local_resource", "Applications", `{}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Convert to struct.
	infos := make([]*structpb.Struct, 0)
	for i := 0; i < len(values); i++ {
		values_ := values[i].(map[string]interface{})

		if values_["icon"] == nil {
			values_["icon"] = ""
		}

		if values_["alias"] == nil {
			values_["alias"] = ""
		}

		info, err := structpb.NewStruct(map[string]interface{}{"_id": values_["_id"], "name": values_["_id"], "path": values_["path"], "creation_date": values_["creation_date"], "last_deployed": values_["last_deployed"], "alias": values_["alias"], "icon": values_["icon"], "description": values_["description"]})
		if err == nil {
			infos = append(infos, info)
		} else {
			log.Println(err)
		}
	}

	return &resourcepb.GetAllApplicationsInfoRsp{
		Applications: infos,
	}, nil

}

////////////////////////////////////////////////////////////////////////////////
// Peer's Authorization and Authentication code.
////////////////////////////////////////////////////////////////////////////////

//* Register a new Peer on the network *
func (self *Globule) RegisterPeer(ctx context.Context, rqst *resourcepb.RegisterPeerRqst) (*resourcepb.RegisterPeerRsp, error) {
	// A peer want to be part of the network.

	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will first look if a peer with a same name already exist on the
	// resources...
	_id := Utility.GenerateUUID(rqst.Peer.Domain)

	count, _ := p.Count(context.Background(), "local_resource", "local_resource", "Peers", `{"_id":"`+_id+`"}`, "")
	if count > 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Peer with name '"+rqst.Peer.Domain+"' already exist!")))
	}

	// No authorization exist for that peer I will insert it.
	// Here will create the new peer.
	peer := make(map[string]interface{}, 0)
	peer["_id"] = _id
	peer["domain"] = rqst.Peer.Domain
	peer["actions"] = make([]interface{}, 0)

	_, err = p.InsertOne(context.Background(), "local_resource", "local_resource", "Peers", peer, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.RegisterPeerRsp{
		Result: true,
	}, nil
}

//* Return the list of authorized peers *
func (self *Globule) GetPeers(rqst *resourcepb.GetPeersRqst, stream resourcepb.ResourceService_GetPeersServer) error {
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	query := rqst.Query
	if len(query) == 0 {
		query = "{}"
	}

	peers, err := p.Find(context.Background(), "local_resource", "local_resource", "Peers", query, ``)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// No I will stream the result over the networks.
	maxSize := 100
	values := make([]*resourcepb.Peer, 0)

	for i := 0; i < len(peers); i++ {
		p := &resourcepb.Peer{Domain: peers[i].(map[string]interface{})["domain"].(string), Actions: make([]string, 0)}
		peers[i].(map[string]interface{})["actions"] = []interface{}(peers[i].(map[string]interface{})["actions"].(primitive.A))
		for j := 0; j < len(peers[i].(map[string]interface{})["actions"].([]interface{})); j++ {
			p.Actions = append(p.Actions, peers[i].(map[string]interface{})["actions"].([]interface{})[j].(string))
		}

		values = append(values, p)

		if len(values) >= maxSize {
			err := stream.Send(
				&resourcepb.GetPeersRsp{
					Peers: values,
				},
			)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			values = make([]*resourcepb.Peer, 0)
		}
	}

	// Send reminding values.
	err = stream.Send(
		&resourcepb.GetPeersRsp{
			Peers: values,
		},
	)

	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return nil
}

//* Remove a peer from the network *
func (self *Globule) DeletePeer(ctx context.Context, rqst *resourcepb.DeletePeerRqst) (*resourcepb.DeletePeerRsp, error) {
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// No authorization exist for that peer I will insert it.
	// Here will create the new peer.
	_id := Utility.GenerateUUID(rqst.Peer.Domain)

	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Peers", `{"_id":"`+_id+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Delete permissions
	err = p.Delete(context.Background(), "local_resource", "local_resource", "Permissions", `{"owner":"`+_id+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeletePeerRsp{
		Result: true,
	}, nil
}

//* Add peer action permission *
func (self *Globule) AddPeerActions(ctx context.Context, rqst *resourcepb.AddPeerActionsRqst) (*resourcepb.AddPeerActionsRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}
	_id := Utility.GenerateUUID(rqst.Domain)

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Peers", `{"_id":"`+_id+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	peer := values.(map[string]interface{})

	needSave := false
	if peer["actions"] == nil {
		peer["actions"] = rqst.Actions
		needSave = true
	} else {
		actions := []interface{}(peer["actions"].(primitive.A))
		for j := 0; j < len(rqst.Actions); j++ {
			exist := false
			for i := 0; i < len(peer["actions"].(primitive.A)); i++ {
				if peer["actions"].(primitive.A)[i].(string) == rqst.Actions[j] {
					exist = true
					break
				}
			}
			if !exist {
				actions = append(actions, rqst.Actions[j])
				needSave = true
			}
		}
		peer["actions"] = actions
	}

	if needSave {
		jsonStr := serialyseObject(peer)
		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Peers", `{"_id":"`+_id+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.AddPeerActionsRsp{Result: true}, nil

}

//* Remove peer action permission *
func (self *Globule) RemovePeerAction(ctx context.Context, rqst *resourcepb.RemovePeerActionRqst) (*resourcepb.RemovePeerActionRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}
	_id := Utility.GenerateUUID(rqst.Domain)
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Peers", `{"_id":"`+_id+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	peer := values.(map[string]interface{})

	needSave := false
	if peer["actions"] == nil {
		peer["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		actions := make([]interface{}, 0)
		for i := 0; i < len(peer["actions"].(primitive.A)); i++ {
			if peer["actions"].(primitive.A)[i].(string) == rqst.Action {
				exist = true
			} else {
				actions = append(actions, peer["actions"].(primitive.A)[i])
			}
		}
		if exist {
			peer["actions"] = actions
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Peer named "+rqst.Domain+" not contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr := serialyseObject(peer)
		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Peers", `{"_id":"`+_id+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.RemovePeerActionRsp{Result: true}, nil
}

//* Register a new organization
func (self *Globule) CreateOrganization(ctx context.Context, rqst *resourcepb.CreateOrganizationRqst) (*resourcepb.CreateOrganizationRsp, error) {

	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will first look if a peer with a same name already exist on the
	// resources...
	count, _ := p.Count(context.Background(), "local_resource", "local_resource", "Organizations", `{"_id":"`+rqst.Organization.Id+`"}`, "")
	if count > 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Organization with name '"+rqst.Organization.Id+"' already exist!")))
	}

	// No authorization exist for that peer I will insert it.
	// Here will create the new peer.
	g := make(map[string]interface{}, 0)
	g["_id"] = rqst.Organization.Id
	g["name"] = rqst.Organization.Name

	// Those are the list of entity linked to the organisation
	g["accounts"] = make([]interface{}, 0)
	g["groups"] = make([]interface{}, 0)
	g["roles"] = make([]interface{}, 0)
	g["applications"] = make([]interface{}, 0)

	_, err = p.InsertOne(context.Background(), "local_resource", "local_resource", "Organizations", g, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// accounts...
	for i := 0; i < len(rqst.Organization.Accounts); i++ {
		self.createCrossReferences(rqst.Organization.Accounts[i], "Accounts", "organizations", rqst.Organization.GetId(), "Organizations", "accounts")
	}

	// groups...
	for i := 0; i < len(rqst.Organization.Groups); i++ {
		self.createCrossReferences(rqst.Organization.Groups[i], "Groups", "organizations", rqst.Organization.GetId(), "Organizations", "groups")
	}

	// roles...
	for i := 0; i < len(rqst.Organization.Roles); i++ {
		self.createCrossReferences(rqst.Organization.Roles[i], "Roles", "organizations", rqst.Organization.GetId(), "Organizations", "roles")
	}

	// applications...
	for i := 0; i < len(rqst.Organization.Applications); i++ {
		self.createCrossReferences(rqst.Organization.Roles[i], "Applications", "organizations", rqst.Organization.GetId(), "Organizations", "applications")
	}

	return &resourcepb.CreateOrganizationRsp{
		Result: true,
	}, nil
}

//* Return the list of organizations
func (self *Globule) GetOrganizations(rqst *resourcepb.GetOrganizationsRqst, stream resourcepb.ResourceService_GetOrganizationsServer) error {
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	query := rqst.Query
	if len(query) == 0 {
		query = "{}"
	}

	organizations, err := p.Find(context.Background(), "local_resource", "local_resource", "Organizations", query, ``)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// No I will stream the result over the networks.
	maxSize := 50
	values := make([]*resourcepb.Organization, 0)
	for i := 0; i < len(organizations); i++ {
		o := organizations[i].(map[string]interface{})

		organization := new(resourcepb.Organization)
		organization.Id = o["_id"].(string)
		organization.Name = o["name"].(string)

		// Here I will set the aggregation.

		// Groups
		if o["groups"] != nil {
			groups := []interface{}(o["groups"].(primitive.A))
			if groups != nil {
				for i := 0; i < len(groups); i++ {
					groupId := groups[i].(map[string]interface{})["$id"].(string)
					organization.Groups = append(organization.Groups, groupId)
				}
			}
		}

		// Roles
		if o["roles"] != nil {
			roles := []interface{}(o["roles"].(primitive.A))
			if roles != nil {
				for i := 0; i < len(roles); i++ {
					roleId := roles[i].(map[string]interface{})["$id"].(string)
					organization.Roles = append(organization.Roles, roleId)
				}
			}
		}

		// Accounts
		if o["accounts"] != nil {
			accounts := []interface{}(o["accounts"].(primitive.A))
			if accounts != nil {
				for i := 0; i < len(accounts); i++ {
					accountId := accounts[i].(map[string]interface{})["$id"].(string)
					organization.Accounts = append(organization.Accounts, accountId)
				}
			}
		}

		// Applications
		if o["applications"] != nil {
			applications := []interface{}(o["applications"].(primitive.A))
			if applications != nil {
				for i := 0; i < len(applications); i++ {
					applicationId := applications[i].(map[string]interface{})["$id"].(string)
					organization.Applications = append(organization.Applications, applicationId)
				}
			}
		}

		values = append(values, organization)
		if len(values) >= maxSize {
			err := stream.Send(
				&resourcepb.GetOrganizationsRsp{
					Organizations: values,
				},
			)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			values = make([]*resourcepb.Organization, 0)
		}
	}

	// Send reminding values.
	err = stream.Send(
		&resourcepb.GetOrganizationsRsp{
			Organizations: values,
		},
	)

	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return nil
}

//* Add Account *
func (self *Globule) AddOrganizationAccount(ctx context.Context, rqst *resourcepb.AddOrganizationAccountRqst) (*resourcepb.AddOrganizationAccountRsp, error) {
	err := self.createCrossReferences(rqst.OrganizationId, "Organizations", "accounts", rqst.AccountId, "Accounts", "organizations")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.AddOrganizationAccountRsp{Result: true}, nil
}

//* Add Group *
func (self *Globule) AddOrganizationGroup(ctx context.Context, rqst *resourcepb.AddOrganizationGroupRqst) (*resourcepb.AddOrganizationGroupRsp, error) {
	err := self.createCrossReferences(rqst.OrganizationId, "Organizations", "groups", rqst.GroupId, "Groups", "organizations")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.AddOrganizationGroupRsp{Result: true}, nil
}

//* Add Role *
func (self *Globule) AddOrganizationRole(ctx context.Context, rqst *resourcepb.AddOrganizationRoleRqst) (*resourcepb.AddOrganizationRoleRsp, error) {
	err := self.createCrossReferences(rqst.OrganizationId, "Organizations", "roles", rqst.RoleId, "Roles", "organizations")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.AddOrganizationRoleRsp{Result: true}, nil
}

//* Add Application *
func (self *Globule) AddOrganizationApplication(ctx context.Context, rqst *resourcepb.AddOrganizationApplicationRqst) (*resourcepb.AddOrganizationApplicationRsp, error) {
	err := self.createCrossReferences(rqst.OrganizationId, "Organizations", "applications", rqst.ApplicationId, "Applications", "organizations")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.AddOrganizationApplicationRsp{Result: true}, nil
}

//* Remove Account *
func (self *Globule) RemoveOrganizationAccount(ctx context.Context, rqst *resourcepb.RemoveOrganizationAccountRqst) (*resourcepb.RemoveOrganizationAccountRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	err = self.deleteReference(p, rqst.AccountId, rqst.OrganizationId, "accounts", "Organizations")
	if err != nil {
		return nil, err
	}

	err = self.deleteReference(p, rqst.OrganizationId, rqst.AccountId, "organizations", "Accounts")
	if err != nil {
		return nil, err
	}

	return &resourcepb.RemoveOrganizationAccountRsp{Result: true}, nil
}

//* Remove Group *
func (self *Globule) RemoveOrganizationGroup(ctx context.Context, rqst *resourcepb.RemoveOrganizationGroupRqst) (*resourcepb.RemoveOrganizationGroupRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	err = self.deleteReference(p, rqst.GroupId, rqst.OrganizationId, "groups", "Organizations")
	if err != nil {
		return nil, err
	}

	err = self.deleteReference(p, rqst.OrganizationId, rqst.GroupId, "organizations", "Groups")
	if err != nil {
		return nil, err
	}

	return &resourcepb.RemoveOrganizationGroupRsp{Result: true}, nil
}

//* Remove Role *
func (self *Globule) RemoveOrganizationRole(ctx context.Context, rqst *resourcepb.RemoveOrganizationRoleRqst) (*resourcepb.RemoveOrganizationRoleRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	err = self.deleteReference(p, rqst.RoleId, rqst.OrganizationId, "roles", "Organizations")
	if err != nil {
		return nil, err
	}

	err = self.deleteReference(p, rqst.OrganizationId, rqst.RoleId, "organizations", "Roles")
	if err != nil {
		return nil, err
	}

	return &resourcepb.RemoveOrganizationRoleRsp{Result: true}, nil
}

//* Remove Application *
func (self *Globule) RemoveOrganizationApplication(ctx context.Context, rqst *resourcepb.RemoveOrganizationApplicationRqst) (*resourcepb.RemoveOrganizationApplicationRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	err = self.deleteReference(p, rqst.ApplicationId, rqst.OrganizationId, "applications", "Organizations")
	if err != nil {
		return nil, err
	}

	err = self.deleteReference(p, rqst.OrganizationId, rqst.ApplicationId, "organizations", "Applications")
	if err != nil {
		return nil, err
	}

	return &resourcepb.RemoveOrganizationApplicationRsp{Result: true}, nil
}

//* Delete organization
func (self *Globule) DeleteOrganization(ctx context.Context, rqst *resourcepb.DeleteOrganizationRqst) (*resourcepb.DeleteOrganizationRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Organizations", `{"_id":"`+rqst.Organization+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	organization := values.(map[string]interface{})
	if organization["groups"] != nil {
		groups := []interface{}(organization["groups"].(primitive.A))
		if groups != nil {
			for i := 0; i < len(groups); i++ {
				groupId := groups[i].(map[string]interface{})["$id"].(string)
				self.deleteReference(p, rqst.Organization, groupId, "organizations", "Organizations")
			}
		}
	}

	if organization["roles"].(primitive.A) != nil {
		roles := []interface{}(organization["roles"].(primitive.A))
		if roles != nil {
			for i := 0; i < len(roles); i++ {
				roleId := roles[i].(map[string]interface{})["$id"].(string)
				self.deleteReference(p, rqst.Organization, roleId, "organizations", "Organizations")
			}
		}
	}

	if organization["applications"].(primitive.A) != nil {
		applications := []interface{}(organization["applications"].(primitive.A))
		if applications != nil {
			for i := 0; i < len(applications); i++ {
				applicationId := applications[i].(map[string]interface{})["$id"].(string)
				self.deleteReference(p, rqst.Organization, applicationId, "organizations", "Organizations")
			}
		}
	}

	if organization["accounts"].(primitive.A) != nil {
		accounts := []interface{}(organization["accounts"].(primitive.A))
		if accounts != nil {
			for i := 0; i < len(accounts); i++ {
				accountsId := accounts[i].(map[string]interface{})["$id"].(string)
				self.deleteReference(p, rqst.Organization, accountsId, "organizations", "Organizations")
			}
		}
	}

	// Try to delete the account...
	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Organizations", `{"_id":"`+rqst.Organization+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteOrganizationRsp{Result: true}, nil
}

func (self *Globule) createGroup(id, name string, members []string) error {
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// Here I will first look if a peer with a same name already exist on the
	// resources...
	count, _ := p.Count(context.Background(), "local_resource", "local_resource", "Groups", `{"_id":"`+id+`"}`, "")
	if count > 0 {
		return errors.New("Group with name '" + id + "' already exist!")
	}

	// No authorization exist for that peer I will insert it.
	// Here will create the new peer.
	g := make(map[string]interface{}, 0)
	g["_id"] = id
	g["name"] = name

	_, err = p.InsertOne(context.Background(), "local_resource", "local_resource", "Groups", g, "")
	if err != nil {
		return err
	}

	// Create references.
	for i := 0; i < len(members); i++ {
		err := self.createCrossReferences(id, "Groups", "members", members[i], "Accounts", "groups")
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

//* Register a new group
func (self *Globule) CreateGroup(ctx context.Context, rqst *resourcepb.CreateGroupRqst) (*resourcepb.CreateGroupRsp, error) {
	// Get the persistence connection
	err := self.createGroup(rqst.Group.Id, rqst.Group.Name, rqst.Group.Members)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.CreateGroupRsp{
		Result: true,
	}, nil
}

//* Return the list of organizations
func (self *Globule) GetGroups(rqst *resourcepb.GetGroupsRqst, stream resourcepb.ResourceService_GetGroupsServer) error {
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	query := rqst.Query
	if len(query) == 0 {
		query = "{}"
	}

	groups, err := p.Find(context.Background(), "local_resource", "local_resource", "Groups", query, ``)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// No I will stream the result over the networks.
	maxSize := 50
	values := make([]*resourcepb.Group, 0)
	for i := 0; i < len(groups); i++ {

		g := &resourcepb.Group{Name: groups[i].(map[string]interface{})["name"].(string), Id: groups[i].(map[string]interface{})["_id"].(string), Members: make([]string, 0)}

		if groups[i].(map[string]interface{})["members"] != nil {
			members := []interface{}(groups[i].(map[string]interface{})["members"].(primitive.A))
			g.Members = make([]string, 0)
			for j := 0; j < len(members); j++ {
				g.Members = append(g.Members, members[j].(map[string]interface{})["$id"].(string))
			}

			values = append(values, g)
			if len(values) >= maxSize {
				err := stream.Send(
					&resourcepb.GetGroupsRsp{
						Groups: values,
					},
				)
				if err != nil {
					return status.Errorf(
						codes.Internal,
						Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
				}
				values = make([]*resourcepb.Group, 0)
			}
		}
	}

	// Send reminding values.
	err = stream.Send(
		&resourcepb.GetGroupsRsp{
			Groups: values,
		},
	)

	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return nil
}

//* Delete organization
func (self *Globule) DeleteGroup(ctx context.Context, rqst *resourcepb.DeleteGroupRqst) (*resourcepb.DeleteGroupRsp, error) {
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Groups", `{"_id":"`+rqst.Group+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	group := values.(map[string]interface{})

	// I will remove it from accounts...

	if group["members"] != nil {
		members := []interface{}(group["members"].(primitive.A))
		for j := 0; j < len(members); j++ {
			self.deleteReference(p, rqst.Group, members[j].(map[string]interface{})["$id"].(string), "groups", members[j].(map[string]interface{})["$ref"].(string))
		}
	}

	// I will remove it from organizations...
	if group["organizations"] != nil {
		organizations := []interface{}(group["organizations"].(primitive.A))
		if organizations != nil {
			for i := 0; i < len(organizations); i++ {
				organizationId := organizations[i].(map[string]interface{})["$id"].(string)
				self.deleteReference(p, rqst.Group, organizationId, "groups", "Groups")
			}
		}
	}

	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Groups", `{"_id":"`+rqst.Group+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteGroupRsp{
		Result: true,
	}, nil

}

//* Add a member account to the group *
func (self *Globule) AddGroupMemberAccount(ctx context.Context, rqst *resourcepb.AddGroupMemberAccountRqst) (*resourcepb.AddGroupMemberAccountRsp, error) {

	err := self.createCrossReferences(rqst.GroupId, "Groups", "members", rqst.AccountId, "Accounts", "groups")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.AddGroupMemberAccountRsp{Result: true}, nil
}

//* Remove member account from the group *
func (self *Globule) RemoveGroupMemberAccount(ctx context.Context, rqst *resourcepb.RemoveGroupMemberAccountRqst) (*resourcepb.RemoveGroupMemberAccountRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// That service made user of persistence service.
	err = self.deleteReference(p, rqst.AccountId, rqst.GroupId, "members", "Groups")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = self.deleteReference(p, rqst.GroupId, rqst.AccountId, "groups", "Accounts")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.RemoveGroupMemberAccountRsp{Result: true}, nil
}
