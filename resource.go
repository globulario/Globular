package main

import (
	"context"
	"os"

	"log"

	"errors"
	"strings"

	"time"

	"encoding/json"
	"reflect"

	"fmt"
	"path/filepath"

	"net"

	"strconv"

	"github.com/davecourtois/Utility"
	"github.com/emicklei/proto"
	"github.com/globulario/Globular/Interceptors"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/golang/protobuf/jsonpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (self *Globule) startResourceService() error {
	id := string(resourcepb.File_proto_resource_proto.Services().Get(0).FullName())
	resource_server, err := self.startInternalService(id, resourcepb.File_proto_resource_proto.Path(), self.ResourcePort, self.ResourceProxy, self.Protocol == "https", self.unaryResourceInterceptor, self.streamResourceInterceptor)
	if err == nil {
		self.inernalServices = append(self.inernalServices, resource_server)

		// Create the channel to listen on resource port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.ResourcePort))
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
			self.saveConfig()
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

func (self *Globule) setServiceMethods(name string, path string) {
	log.Println("-------> ", path)
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
			"/resource.ResourceService/AccountExist",
			"/resource.ResourceService/Authenticate",
			"/resource.ResourceService/RefreshToken",
			"/resource.ResourceService/GetPermissions",
			"/resource.ResourceService/GetAllFilesInfo",
			"/resource.ResourceService/GetAllApplicationsInfo",
			"/resource.ResourceService/GetResourceOwners",
			"/resource.ResourceService/ValidateUserAccess",
			"/resource.ResourceService/ValidateUserResourceAccess",
			"/resource.ResourceService/ValidateApplicationAccess",
			"/resource.ResourceService/ValidateApplicationResourceAccess",
			"/event.EventService/Subscribe",
			"/event.EventService/UnSubscribe", "/event.EventService/OnEvent",
			"/event.EventService/Quit",
			"/event.EventService/Publish",
			"/services.ServiceDiscovery/FindServices",
			"/services.ServiceDiscovery/GetServiceDescriptor",
			"/services.ServiceDiscovery/GetServicesDescriptor",
			"/services.ServiceRepository/downloadBundle",
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

	// Create connection application.
	self.createApplicationConnection()

	// Now I will save the action permission.
	for i := 0; i < len(self.actionPermissions); i++ {
		permission := self.actionPermissions[i].(map[string]interface{})
		self.setActionPermission(permission["action"].(string), permission["actionParameterResourcePermissions"].([]interface{}))
	}

	return nil
}

func (self *Globule) createApplicationConnection() error {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	store, err := self.getPersistenceStore()
	applications, _ := store.Find(context.Background(), "local_resource", "local_resource", "Applications", "{}", "")
	if err != nil {
		return err
	}
	address, port := self.getBackendAddress()
	for i := 0; i < len(applications); i++ {
		application := applications[i].(map[string]interface{})
		// Open the user database connection.

		err = p.CreateConnection(application["_id"].(string)+"_db", application["_id"].(string)+"_db", address, float64(port), 0, application["_id"].(string), application["password"].(string), 5000, "", false)
		if err != nil {
			return err
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

	if rqst.SyncInfo.UserSyncInfos == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No LDAP sync users infos was given!")))
	}

	syncInfo, err := Utility.ToMap(rqst.SyncInfo)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)] == nil {
		self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)] = make([]interface{}, 0)
		self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)] = append(self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)].([]interface{}), syncInfo)
	} else {
		syncInfos := self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)].([]interface{})
		exist := false
		for i := 0; i < len(syncInfos); i++ {
			if syncInfos[i].(map[string]interface{})["ldapSeriveId"] == syncInfo["ldapSeriveId"] {
				if syncInfos[i].(map[string]interface{})["connectionId"] == syncInfo["connectionId"] {
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
			self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)] = append(self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)].([]interface{}), syncInfo)
			// save the config.
			self.saveConfig()

		}
	}

	// Cast the the correct type.

	// Searh for roles.
	ldap_, err := self.getLdapClient()
	if err != nil {
		return nil, err
	}
	rolesInfo, err := ldap_.Search(rqst.SyncInfo.ConnectionId, rqst.SyncInfo.GroupSyncInfos.Base, rqst.SyncInfo.GroupSyncInfos.Query, []string{rqst.SyncInfo.GroupSyncInfos.Id, "distinguishedName"})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Print role info.
	for i := 0; i < len(rolesInfo); i++ {
		name := rolesInfo[i][0].([]interface{})[0].(string)
		id := Utility.GenerateUUID(rolesInfo[i][1].([]interface{})[0].(string))
		self.createRole(id, name, []string{})
	}

	// Synchronize account and user info...
	accountsInfo, err := ldap_.Search(rqst.SyncInfo.ConnectionId, rqst.SyncInfo.UserSyncInfos.Base, rqst.SyncInfo.UserSyncInfos.Query, []string{rqst.SyncInfo.UserSyncInfos.Id, rqst.SyncInfo.UserSyncInfos.Email, "distinguishedName", "memberOf"})

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(accountsInfo); i++ {
		// Print the list of account...
		// I will not set the password...
		name := strings.ToLower(accountsInfo[i][0].([]interface{})[0].(string))

		if len(accountsInfo[i][1].([]interface{})) > 0 {
			email := strings.ToLower(accountsInfo[i][1].([]interface{})[0].(string))

			if len(email) > 0 {

				id := Utility.GenerateUUID(strings.ToLower(accountsInfo[i][2].([]interface{})[0].(string)))
				if len(id) > 0 {
					roles := make([]interface{}, 0)
					roles = append(roles, "guest")
					// Here I will set the roles of the user.
					if len(accountsInfo[i][3].([]interface{})) > 0 {
						for j := 0; j < len(accountsInfo[i][3].([]interface{})); j++ {
							roles = append(roles, Utility.GenerateUUID(accountsInfo[i][3].([]interface{})[j].(string)))
						}
					}

					// Try to create account...
					err := self.registerAccount(id, name, email, id, roles)
					if err != nil {
						rolesStr := `[{"$ref":"Roles","$id":"guest","$db":"local_resource"}`
						for j := 0; j < len(accountsInfo[i][3].([]interface{})); j++ {
							rolesStr += `,{"$ref":"Roles","$id":"` + Utility.GenerateUUID(accountsInfo[i][3].([]interface{})[j].(string)) + `","$db":"local_resource"}`
						}

						rolesStr += `]`
						err := p.UpdateOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+id+`"}`, `{ "$set":{"roles":`+rolesStr+`}}`, "")
						if err != nil {
							return nil, status.Errorf(
								codes.Internal,
								Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
						}

					}
				}
			} else {

				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("account "+strings.ToLower(accountsInfo[i][2].([]interface{})[0].(string))+" has no email configured! ")))
			}
		}
	}

	if rqst.SyncInfo.GroupSyncInfos == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No LDAP sync groups infos was given!")))
	}

	return &resourcepb.SynchronizeLdapRsp{
		Result: true,
	}, nil
}

func (self *Globule) registerAccount(id string, name string, email string, password string, roles []interface{}) error {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// first of all the Persistence service must be active.
	count, err := p.Count(context.Background(), "local_resource", "local_resource", "Accounts", `{"name":"`+name+`"}`, "")
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

	account["roles"] = make([]map[string]interface{}, 0)
	for i := 0; i < len(roles); i++ {
		role := make(map[string]interface{}, 0)
		role["$id"] = roles[i]
		role["$ref"] = "Roles"
		role["$db"] = "local_resource"
		account["roles"] = append(account["roles"].([]map[string]interface{}), role)
	}

	// Here I will insert the account in the database.
	_, err = p.InsertOne(context.Background(), "local_resource", "local_resource", "Accounts", account, "")

	// replace @ and . by _
	name = strings.ReplaceAll(strings.ReplaceAll(name, "@", "_"), ".", "_")

	// Each account will have their own database and a use that can read and write
	// into it.
	// Here I will wrote the script for mongoDB...
	createUserScript := fmt.Sprintf(
		"db=db.getSiblingDB('%s_db');db.createCollection('user_data');db=db.getSiblingDB('admin');db.createUser({user: '%s', pwd: '%s',roles: [{ role: 'dbOwner', db: '%s_db' }]});",
		name, name, password, name)

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

	return nil

}

/* Register a new Account */
func (self *Globule) RegisterAccount(ctx context.Context, rqst *resourcepb.RegisterAccountRqst) (*resourcepb.RegisterAccountRsp, error) {
	if rqst.ConfirmPassword != rqst.Password {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to confirm your password!")))

	}

	if rqst.Account == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No account information was given!")))

	}

	err := self.registerAccount(rqst.Account.Name, rqst.Account.Name, rqst.Account.Email, rqst.Password, []interface{}{"guest"})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Generate a token to identify the user.
	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, rqst.Account.Name, rqst.Account.Email)
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

	name, _, expireAt, _ := Interceptors.ValidateToken(tokenString)
	_, err = p.InsertOne(context.Background(), "local_resource", "local_resource", "Tokens", map[string]interface{}{"_id": name, "expireAt": Utility.ToString(expireAt)}, "")
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

func (self *Globule) AccountExist(ctx context.Context, rqst *resourcepb.AccountExistRqst) (*resourcepb.AccountExistRsp, error) {
	var exist bool
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}
	// Test with the _id
	count, _ := p.Count(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+rqst.Id+`"}`, "")
	if count > 0 {
		exist = true
	}

	// Test with the name
	if !exist {
		count, _ := p.Count(context.Background(), "local_resource", "local_resource", "Accounts", `{"name":"`+rqst.Id+`"}`, "")
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
		tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, "sa", self.AdminEmail)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		name, _, expireAt, _ := Interceptors.ValidateToken(tokenString)
		err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+name+`"}`, `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

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

	// Generate a token to identify the user.
	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, values[0].(map[string]interface{})["name"].(string), values[0].(map[string]interface{})["email"].(string))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

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
	name, _, expireAt, _ := Interceptors.ValidateToken(tokenString)
	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+name+`"}`, `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// So here I will start a timer

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
	name, email, expireAt, err := Interceptors.ValidateToken(rqst.Token)

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
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+name+`"}`, `[{"Projection":{"expireAt":1}}]`)
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

	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, name, email)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// get back the new expireAt
	name, _, expireAt, _ = Interceptors.ValidateToken(tokenString)

	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Tokens", `{"_id":"`+name+`"}`, `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, `[{"upsert":true}]`)
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

//* Delete an account *
func (self *Globule) DeleteAccount(ctx context.Context, rqst *resourcepb.DeleteAccountRqst) (*resourcepb.DeleteAccountRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+rqst.Id+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	account := values.(map[string]interface{})

	// Try to delete the account...
	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+rqst.Id+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Delete permissions
	err = self.deleteAccountPermissions("/" + rqst.Id)
	if err != nil {
		return nil, err
	}

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

	return &resourcepb.CreateRoleRsp{Result: true}, nil
}

//* Delete a role with a given id *
func (self *Globule) DeleteRole(ctx context.Context, rqst *resourcepb.DeleteRoleRqst) (*resourcepb.DeleteRoleRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	accounts, err := p.Find(context.Background(), "local_resource", "local_resource", "Accounts", `{}`, ``)
	if err == nil {
		for i := 0; i < len(accounts); i++ {
			account := accounts[i].(map[string]interface{})
			if account["roles"] != nil {
				roles := []interface{}(account["roles"].(primitive.A))
				roles_ := make([]interface{}, 0)
				needSave := false
				for j := 0; j < len(roles); j++ {
					// TODO remove the role with name rqst.roleId from the account.
					role := roles[j].(map[string]interface{})
					if role["$id"] == rqst.RoleId {
						needSave = true
					} else {
						roles_ = append(roles_, role)
					}
				}

				// Here I will save the role.
				if needSave {
					account["roles"] = roles_
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
						return nil, status.Errorf(
							codes.Internal,
							Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
					}
				}
			}
		}
	}

	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = self.deleteRolePermissions("/" + rqst.RoleId)
	if err != nil {
		return nil, err
	}

	return &resourcepb.DeleteRoleRsp{Result: true}, nil
}

//* Append an action to existing role. *
func (self *Globule) AddRoleAction(ctx context.Context, rqst *resourcepb.AddRoleActionRqst) (*resourcepb.AddRoleActionRsp, error) {
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
		actions := []interface{}(role["actions"].(primitive.A))
		for i := 0; i < len(actions); i++ {
			if actions[i].(string) == rqst.Action {
				exist = true
				break
			}
		}
		if !exist {
			actions = append(actions, rqst.Action)
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Role named "+rqst.RoleId+"already contain actions named "+rqst.Action+"!")))
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

	return &resourcepb.AddRoleActionRsp{Result: true}, nil
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
		jsonStr, _ := json.Marshal(role)
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
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No account named "+rqst.AccountId+" exist!")))
	}

	account := values.(map[string]interface{})

	// Now I will test if the account already contain the role.
	if account["roles"] != nil {
		roles := []interface{}(account["roles"].(primitive.A))
		for j := 0; j < len(roles); j++ {
			if roles[j].(map[string]interface{})["$id"] == rqst.RoleId {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Role named "+rqst.RoleId+" aleready exist in account "+rqst.AccountId+"!")))
			}
		}

		// append the newly created role.
		account["roles"] = append(roles, map[string]interface{}{"$ref": "Roles", "$id": rqst.RoleId, "$db": "local_resource"})

		// Here I will save the role.
		jsonStr := "{"
		jsonStr += `"name":"` + account["name"].(string) + `",`
		jsonStr += `"email":"` + account["email"].(string) + `",`
		jsonStr += `"password":"` + account["password"].(string) + `",`
		jsonStr += `"roles":[`

		if reflect.TypeOf(account["roles"]).String() == "primitive.A" {
			account["roles"] = []interface{}(account["roles"].(primitive.A))
		}

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

		err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, jsonStr, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

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
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No account named "+rqst.AccountId+" exist!")))
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

		// Here I will save the role.
		jsonStr := "{"
		jsonStr += `"name":"` + account["name"].(string) + `",`
		jsonStr += `"email":"` + account["email"].(string) + `",`
		jsonStr += `"password":"` + account["password"].(string) + `",`
		jsonStr += `"roles":[`
		if reflect.TypeOf(account["roles"]).String() == "primitive.A" {
			account["roles"] = []interface{}(account["roles"].(primitive.A))
		}
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

		err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, jsonStr, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

	}
	return &resourcepb.RemoveAccountRoleRsp{Result: true}, nil
}

//* Append an action to existing application. *
func (self *Globule) AddApplicationAction(ctx context.Context, rqst *resourcepb.AddApplicationActionRqst) (*resourcepb.AddApplicationActionRsp, error) {
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
		application["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		application["actions"] = []interface{}(application["actions"].(primitive.A))
		for i := 0; i < len(application["actions"].([]interface{})); i++ {
			if application["actions"].([]interface{})[i].(string) == rqst.Action {
				exist = true
				break
			}
		}
		if !exist {
			application["actions"] = append(application["actions"].([]interface{}), rqst.Action)
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Application named "+rqst.ApplicationId+" already contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr, _ := json.Marshal(application)
		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.AddApplicationActionRsp{Result: true}, nil
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
		jsonStr, _ := json.Marshal(application)
		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.RemoveApplicationActionRsp{Result: true}, nil
}

//* Delete an application from the server. *
func (self *Globule) DeleteApplication(ctx context.Context, rqst *resourcepb.DeleteApplicationRqst) (*resourcepb.DeleteApplicationRsp, error) {

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

	// First of all I will remove the directory.
	err = os.RemoveAll(application["path"].(string))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Now I will remove the database create for the application.
	err = p.DeleteDatabase(context.Background(), "local_resource", rqst.ApplicationId+"_db")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Finaly I will remove the entry in  the table.
	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Delete permissions
	err = p.Delete(context.Background(), "local_resource", "local_resource", "Permissions", `{"owner":"`+rqst.ApplicationId+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Drop the application user.
	// Here I will drop the db user.
	dropUserScript := fmt.Sprintf(
		`db=db.getSiblingDB('admin');db.dropUser('%s', {w: 'majority', wtimeout: 4000})`,
		rqst.ApplicationId)

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
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	// TODO delete dir permission associate with the application.

	return &resourcepb.DeleteApplicationRsp{
		Result: true,
	}, nil
}

func (self *Globule) deleteAccountPermissions(name string) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	err = p.Delete(context.Background(), "local_resource", "local_resource", "Accounts", `{"name":"`+name+`"}`, "")
	if err != nil {
		return err
	}

	return nil
}

//* Delete all permission for a given account *
func (self *Globule) DeleteAccountPermissions(ctx context.Context, rqst *resourcepb.DeleteAccountPermissionsRqst) (*resourcepb.DeleteAccountPermissionsRsp, error) {

	err := self.deleteAccountPermissions(rqst.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteAccountPermissionsRsp{
		Result: true,
	}, nil
}

func (self *Globule) deleteRolePermissions(name string) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	err = p.Delete(context.Background(), "local_resource", "local_resource", "Roles", `{"name":"`+name+`"}`, "")
	if err != nil {
		return err
	}

	return nil
}

//* Delete all permission for a given role *
func (self *Globule) DeleteRolePermissions(ctx context.Context, rqst *resourcepb.DeleteRolePermissionsRqst) (*resourcepb.DeleteRolePermissionsRsp, error) {

	err := self.deleteRolePermissions(rqst.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteRolePermissionsRsp{
		Result: true,
	}, nil
}

/**
 * Return the list of all actions avalaible on the server.
 */
func (self *Globule) GetAllActions(ctx context.Context, rqst *resourcepb.GetAllActionsRqst) (*resourcepb.GetAllActionsRsp, error) {
	return &resourcepb.GetAllActionsRsp{Actions: self.methods}, nil
}

/////////////////////// Resource permission owner /////////////////////////////
func (self *Globule) setResourceOwner(owner string, path string) error {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	if strings.TrimSpace(path) == "/" {
		return errors.New("Root has no owner!")
	}

	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}

	// Here I will set resources whit that path, be sure to have different
	// path than application and webroot path if you dont want permission follow each other.
	resources, err := self.getResources(path)
	if err == nil {
		for i := 0; i < len(resources); i++ {
			if resources[i].GetPath() != path {
				path_ := resources[i].GetPath()[len(path)+1:]
				paths := strings.Split(path_, "/")
				path_ = path
				// set sub-path...
				for j := 0; j < len(paths); j++ {
					path_ += "/" + paths[j]
					self.setResourceOwner(owner, path_)
				}
			}

			if strings.HasSuffix(resources[i].GetPath(), "/") {
				self.setResourceOwner(owner, resources[i].GetPath()+resources[i].GetName())
			} else {
				self.setResourceOwner(owner, resources[i].GetPath()+"/"+resources[i].GetName())
			}

		}
	}

	// Here I will set the resource owner.
	resourceOwner := make(map[string]interface{})
	resourceOwner["owner"] = owner
	resourceOwner["path"] = path
	resourceOwner["_id"] = Utility.GenerateUUID(owner + path)

	// Here if the
	jsonStr, err := Utility.ToJson(&resourceOwner)
	if err != nil {
		return err
	}

	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "ResourceOwners", jsonStr, jsonStr, `[{"upsert":true}]`)
	if err != nil {
		return err
	}

	return nil
}

//* Set Resource owner *
func (self *Globule) SetResourceOwner(ctx context.Context, rqst *resourcepb.SetResourceOwnerRqst) (*resourcepb.SetResourceOwnerRsp, error) {
	// That service made user of persistence service.
	path := rqst.GetPath()

	err := self.setResourceOwner(rqst.GetOwner(), path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.SetResourceOwnerRsp{
		Result: true,
	}, nil
}

//* Get the resource owner *
func (self *Globule) GetResourceOwners(ctx context.Context, rqst *resourcepb.GetResourceOwnersRqst) (*resourcepb.GetResourceOwnersRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Get the absolute path
	path := rqst.GetPath()

	// find the resource with it id
	resourceOwners, err := p.Find(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	owners := make([]string, 0)
	for i := 0; i < len(resourceOwners); i++ {
		owners = append(owners, resourceOwners[i].(map[string]interface{})["owner"].(string))
	}

	return &resourcepb.GetResourceOwnersRsp{
		Owners: owners,
	}, nil
}

//* Get the resource owner *
func (self *Globule) DeleteResourceOwner(ctx context.Context, rqst *resourcepb.DeleteResourceOwnerRqst) (*resourcepb.DeleteResourceOwnerRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	path := rqst.GetPath()

	// Delete the resource owner for a given path.
	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+path+`","owner":"`+rqst.GetOwner()+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteResourceOwnerRsp{
		Result: true,
	}, nil
}

func (self *Globule) DeleteResourceOwners(ctx context.Context, rqst *resourcepb.DeleteResourceOwnersRqst) (*resourcepb.DeleteResourceOwnersRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	path := rqst.GetPath()

	// delete the resource owners with it path
	err = p.Delete(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteResourceOwnersRsp{
		Result: true,
	}, nil
}

//////////////////////////// Loggin info ///////////////////////////////////////
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

// unaryInterceptor calls authenticateClient with current context
func (self *Globule) unaryResourceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	method := info.FullMethod

	// The token and the application id.
	var token string
	var application string
	var clientId string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
	}

	hasAccess := false
	var err error

	// Here some method are accessible by default.
	if method == "/resource.ResourceService/GetAllActions" ||
		method == "/resource.ResourceService/RegisterAccount" ||
		method == "/resource.ResourceService/RegisterPeer" ||
		method == "/resource.ResourceService/GetPeers" ||
		method == "/resource.ResourceService/AccountExist" ||
		method == "/resource.ResourceService/Authenticate" ||
		method == "/resource.ResourceService/RefreshToken" ||
		method == "/resource.ResourceService/GetPermissions" ||
		method == "/resource.ResourceService/GetAllFilesInfo" ||
		method == "/resource.ResourceService/GetAllApplicationsInfo" ||
		method == "/resource.ResourceService/GetResourceOwners" ||
		method == "/resource.ResourceService/ValidateToken" ||
		method == "/resource.ResourceService/ValidateUserAccess" ||
		method == "/resource.ResourceService/ValidateApplicationAccess" ||
		method == "/resource.ResourceService/ValidateApplicationResourceAccess" ||
		method == "/resource.ResourceService/ValidateUserResourceAccess" ||
		method == "/resource.ResourceService/GetActionPermission" ||
		method == "/resource.ResourceService/Log" ||
		method == "/resource.ResourceService/GetLog" {
		hasAccess = true
	}

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		var expiredAt int64
		clientId, _, expiredAt, err = Interceptors.ValidateToken(token)
		if err != nil {
			return nil, err
		}

		if expiredAt < time.Now().Unix() {
			return nil, errors.New("The token is expired!")
		}

		if clientId == "sa" {
			hasAccess = true
		} else {
			// special case that need ownership of the resource or be sa
			err = self.validateUserAccess(clientId, method)
			if err == nil {
				hasAccess = true
			}

			// If a resource is set set it resource owner.
			if method == "/resource.ResourceService/SetResource" {
				if clientId != "sa" {
					rqst := req.(*resourcepb.SetResourceRqst)

					// In that case i will set the userId as the owner of the resourcepb.
					err := self.setResourceOwner(clientId, rqst.Resource.Path+"/"+rqst.Resource.Name)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		err = self.validateApplicationAccess(application, method)
		if err == nil {
			hasAccess = true
		}
	}

	if !hasAccess {
		err := errors.New("Permission denied to execute method " + method)
		self.logInfo(application, method, token, err)
		return nil, err
	}

	// Execute the action.
	result, err := handler(ctx, req)

	if err == nil {

	}

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

//* Log error or information into the data base *
func (self *Globule) Log(ctx context.Context, rqst *resourcepb.LogRqst) (*resourcepb.LogRsp, error) {
	// Publish event...
	self.log(rqst.Info)

	return &resourcepb.LogRsp{
		Result: true,
	}, nil
}

//* Log error or information into the data base *
// Retreive log infos (the query must be something like /infos/'date'/'applicationName'/'userName'
func (self *Globule) GetLog(rqst *resourcepb.GetLogRqst, stream resourcepb.ResourceService_GetLogServer) error {

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

///////////////////////  resource management. /////////////////
func (self *Globule) setResource(r *resourcepb.Resource) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	var marshaler jsonpb.Marshaler

	jsonStr, err := marshaler.MarshalToString(r)
	if err != nil {
		return err
	}

	// Always create a new if not already exist.
	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Resources", `{ "path" : "`+r.Path+`", "name":"`+r.Name+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return err
	}

	return nil
}

//* Set a resource from a client (custom service) to globular
func (self *Globule) SetResource(ctx context.Context, rqst *resourcepb.SetResourceRqst) (*resourcepb.SetResourceRsp, error) {
	if rqst.Resource == nil {

		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no resource was given!")))

	}
	err := self.setResource(rqst.Resource)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.SetResourceRsp{
		Result: true,
	}, nil
}

func (self *Globule) getResources(path string) ([]*resourcepb.Resource, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	data, err := p.Find(context.Background(), "local_resource", "local_resource", "Resources", `{}`, `[{"Projection":{"_id":0}}]`)
	if err != nil {
		return nil, err
	}

	resources := make([]*resourcepb.Resource, 0)

	for i := 0; i < len(data); i++ {
		res := data[i].(map[string]interface{})
		var size int64
		if res["size"] != nil {
			size = int64(Utility.ToInt(res["size"].(string)))
		}

		var modified int64
		if res["modified"] != nil {
			modified = int64(Utility.ToInt(res["modified"].(string)))
		}

		if res["path"] == nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No path was given for the resource!")))
		}

		if res["name"] == nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No name was given for the resource!")))
		}

		// append the info inside the stream.
		if strings.HasPrefix(res["path"].(string), path) {
			resources = append(resources, &resourcepb.Resource{Path: res["path"].(string), Name: res["name"].(string), Modified: modified, Size: size})
		}
	}
	return resources, nil
}

//* Get all resources
func (self *Globule) GetResources(rqst *resourcepb.GetResourcesRqst, stream resourcepb.ResourceService_GetResourcesServer) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	query := make(map[string]string, 0)
	if len(rqst.Name) > 0 {
		query["name"] = rqst.Name
	}

	if len(rqst.Path) > 0 {
		query["path"] = rqst.Path
	}

	queryStr, _ := Utility.ToJson(query)

	data, err := p.Find(context.Background(), "local_resource", "local_resource", "Resources", queryStr, ``)
	if err != nil {
		return err
	}

	resources := make([]*resourcepb.Resource, 0)
	max := 100
	for i := 0; i < len(data); i++ {
		res := data[i].(map[string]interface{})
		var size int64
		if res["size"] != nil {
			size = int64(Utility.ToInt(res["size"].(string)))
		}

		var modified int64
		if res["modified"] != nil {
			modified = int64(Utility.ToInt(res["modified"].(string)))
		}

		if res["path"] == nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No path was given for the resource!")))
		}

		if res["name"] == nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No name was given for the resource!")))
		}

		// append the info inside the stream.
		resources = append(resources, &resourcepb.Resource{Path: res["path"].(string), Name: res["name"].(string), Modified: modified, Size: size})
		if len(resources) == max-1 {
			// I will send the stream at each 100 logs...
			rsp := &resourcepb.GetResourcesRsp{
				Resources: resources,
			}
			// Send the infos
			err = stream.Send(rsp)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			resources = make([]*resourcepb.Resource, 0)
		}
	}

	// Send the last infos...
	if len(resources) > 0 {
		rsp := &resourcepb.GetResourcesRsp{
			Resources: resources,
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

//* Remove a resource from a client (custom service) to globular
func (self *Globule) RemoveResource(ctx context.Context, rqst *resourcepb.RemoveResourceRqst) (*resourcepb.RemoveResourceRsp, error) {

	// Because regex dosent work properly I retreive all the resources.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// get all resource with that path.
	resources, err := self.getResources(rqst.Resource.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	toDelete := make([]*resourcepb.Resource, 0)
	// Remove resource that match...
	for i := 0; i < len(resources); i++ {
		res := resources[i]

		// In case the resource is a sub-resource I will remove it...
		if len(rqst.Resource.Name) > 0 {
			if rqst.Resource.Name == res.GetName() {
				toDelete = append(toDelete, res) // mark to be delete.
			}
		} else {
			toDelete = append(toDelete, res) // mark to be delete
		}

	}

	// Now I will delete the resourcepb.
	for i := 0; i < len(toDelete); i++ {
		err := p.DeleteOne(context.Background(), "local_resource", "local_resource", "Resources", `{"path":"`+toDelete[i].Path+`", "name":"`+toDelete[i].Name+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		// Delete the permissions ascosiated permission
		self.deletePermissions(toDelete[i].Path+"/"+toDelete[i].Name, "")
		err = p.Delete(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+toDelete[i].Path+"/"+toDelete[i].Name+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	// Remove the resource owner's
	if len(rqst.Resource.Name) == 0 {
		self.deletePermissions(rqst.Resource.Path, "")
		err = p.Delete(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+rqst.Resource.Path+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.RemoveResourceRsp{
		Result: true,
	}, nil
}

func (self *Globule) setActionPermission(action string, actionParameterResourcePermissions []interface{}) error {

	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	actionPermission := make(map[string]interface{}, 0)
	actionPermission["action"] = action

	// Set value in the actionPermission
	actionPermission["actionParameterResourcePermissions"] = actionParameterResourcePermissions

	actionPermission["_id"] = Utility.GenerateUUID(action)

	// Serialyse it
	actionPermissionStr, _ := Utility.ToJson(actionPermission)
	err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "ActionPermission", `{"_id":"`+actionPermission["_id"].(string)+`"}`, actionPermissionStr, `[{"upsert":true}]`)
	if err != nil {
		return err
	}

	return nil
}

//* Set a resource from a client (custom service) to globular
func (self *Globule) SetActionPermission(ctx context.Context, rqst *resourcepb.SetActionPermissionRqst) (*resourcepb.SetActionPermissionRsp, error) {
	actionParameterResourcePermissions := make([]interface{}, len(rqst.ActionParameterResourcePermissions))
	for i := 0; i < len(rqst.ActionParameterResourcePermissions); i++ {
		actionParameterResourcePermissions[i] = map[string]interface{}{"Index": rqst.ActionParameterResourcePermissions[i].Index, "Permission": rqst.ActionParameterResourcePermissions[i].Permission}
	}

	err := self.setActionPermission(rqst.Action, actionParameterResourcePermissions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.SetActionPermissionRsp{
		Result: true,
	}, nil
}

//* Remove a resource from a client (custom service) to globular
func (self *Globule) RemoveActionPermission(ctx context.Context, rqst *resourcepb.RemoveActionPermissionRqst) (*resourcepb.RemoveActionPermissionRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Try to delete the account...
	err = p.DeleteOne(context.Background(), "local_resource", "local_resource", "ActionPermission", `{"_id":"`+Utility.GenerateUUID(rqst.Action)+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.RemoveActionPermissionRsp{
		Result: true,
	}, nil
}

func (self *Globule) getActionPermission(action string) (map[string]interface{}, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "ActionPermission", `{"_id":"`+Utility.GenerateUUID(action)+`"}`, "")
	if err != nil {
		return nil, err
	}

	return values.(map[string]interface{}), nil
}

//* Remove a resource from a client (custom service) to globular
func (self *Globule) GetActionPermission(ctx context.Context, rqst *resourcepb.GetActionPermissionRqst) (*resourcepb.GetActionPermissionRsp, error) {

	actionPermission, err := self.getActionPermission(rqst.Action)
	if err != nil {
		return nil, err
	}

	actionParameterResourcePermissions := make([]*resourcepb.ActionParameterResourcePermission, len(actionPermission["actionParameterResourcePermissions"].(primitive.A)))
	for i := 0; i < len(actionPermission["actionParameterResourcePermissions"].(primitive.A)); i++ {
		a := actionPermission["actionParameterResourcePermissions"].(primitive.A)[i].(map[string]interface{})
		actionParameterResourcePermissions[i] = &resourcepb.ActionParameterResourcePermission{
			Index:      int32(Utility.ToInt(a["Index"])),
			Permission: int32(Utility.ToInt(a["Permission"])),
		}
	}

	return &resourcepb.GetActionPermissionRsp{
		ActionParameterResourcePermissions: actionParameterResourcePermissions,
	}, nil
}

func (self *Globule) savePermission(owner string, path string, permission int32) error {

	// dir cannot be executable....
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// Here I will insert one or replcace one depending if permission already exist or not.
	query := `{"owner":"` + owner + `","path":"` + path + `"}`
	jsonStr := `{"owner":"` + owner + `","path":"` + path + `","permission":` + Utility.ToString(permission) + `}`

	count, _ := p.Count(context.Background(), "local_resource", "local_resource", "Permissions", query, "")

	if count == 0 {
		_, err = p.InsertOne(context.Background(), "local_resource", "local_resource", "Permissions", map[string]interface{}{"owner": owner, "path": path, "permission": permission}, "")

	} else {
		err = p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Permissions", query, jsonStr, "")
	}

	return err
}

// Set resource permission.
func (self *Globule) setResourcePermission(owner string, path string, permission int32) error {
	return self.savePermission(owner, path, permission)
}

//* Set a resource permission, create new one if not already exist. *
func (self *Globule) SetPermission(ctx context.Context, rqst *resourcepb.SetPermissionRqst) (*resourcepb.SetPermissionRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The first thing I will do is test if the file exist.
	path := rqst.GetPermission().GetPath()
	path = strings.ReplaceAll(path, "\\", "/")
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-2]
	}

	// Now I will test if the user or the role exist.
	var owner map[string]interface{}

	switch v := rqst.Permission.GetOwner().(type) {
	case *resourcepb.ResourcePermission_User:
		// In that case I will try to find a user with that id
		values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+v.User+`"}`, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		owner = values.(map[string]interface{})
	case *resourcepb.ResourcePermission_Role:
		// In that case I will try to find a role with that id
		values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+v.Role+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		owner = values.(map[string]interface{})
	case *resourcepb.ResourcePermission_Application:
		// In that case I will try to find a role with that id
		values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+v.Application+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		owner = values.(map[string]interface{})
	}

	resources, err := self.getResources(path)
	if err == nil {
		for i := 0; i < len(resources); i++ {
			if resources[i].GetPath() != path {
				path_ := resources[i].GetPath()[len(path)+1:]
				paths := strings.Split(path_, "/")
				path_ = path
				// set sub-path...
				for j := 0; j < len(paths); j++ {
					path_ += "/" + paths[j]
					err := self.setResourcePermission(owner["_id"].(string), path_, rqst.Permission.Number)
					if err != nil {
						return nil, status.Errorf(
							codes.Internal,
							Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
					}
				}
			}

			// create resource permission
			err := self.setResourcePermission(owner["_id"].(string), resources[i].GetPath()+"/"+resources[i].GetName(), rqst.Permission.Number)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}

		}
		// save resource path.
		err = self.setResourcePermission(owner["_id"].(string), path, rqst.Permission.Number)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.SetPermissionRsp{
		Result: true,
	}, nil
}

func (self *Globule) setPermissionOwner(owner string, permission *resourcepb.ResourcePermission) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// Here I will try to find the owner in the user table
	_, err = p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+owner+`"}`, ``)
	if err == nil {
		permission.Owner = &resourcepb.ResourcePermission_User{
			User: owner,
		}
		return nil
	}

	// In the role.
	// In that case I will try to find a role with that id
	_, err = p.FindOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+owner+`"}`, "")
	if err == nil {
		permission.Owner = &resourcepb.ResourcePermission_Role{
			Role: owner,
		}
		return nil
	}

	_, err = p.FindOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+owner+`"}`, "")
	if err == nil {
		permission.Owner = &resourcepb.ResourcePermission_Application{
			Application: owner,
		}
		return nil
	}

	return errors.New("No Role or User found with id " + owner)
}

func (self *Globule) getResourcePermissions(path string) ([]*resourcepb.ResourcePermission, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// So here I will get the list of retreived permission.
	permissions_, err := p.Find(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, err
	}

	permissions := make([]*resourcepb.ResourcePermission, 0)

	for i := 0; i < len(permissions_); i++ {
		permission_ := permissions_[i].(map[string]interface{})
		permission := &resourcepb.ResourcePermission{Path: permission_["path"].(string), Owner: nil, Number: permission_["permission"].(int32)}
		err = self.setPermissionOwner(permission_["owner"].(string), permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (self *Globule) getPermissions(path string) ([]*resourcepb.ResourcePermission, error) {

	// do file stuff
	permissions, err := self.getResourcePermissions(path)
	if err != nil {
		return nil, err
	}

	return permissions, nil

}

//* Get All permissions for a given file/dir *
func (self *Globule) GetPermissions(ctx context.Context, rqst *resourcepb.GetPermissionsRqst) (*resourcepb.GetPermissionsRsp, error) {

	path := rqst.GetPath()
	permissions, err := self.getPermissions(path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	permissions_ := make([]map[string]interface{}, len(permissions))
	for i := 0; i < len(permissions); i++ {
		permission := make(map[string]interface{}, 0)
		// Set the values.
		permission["path"] = permissions[i].GetPath()
		permission["number"] = permissions[i].GetNumber()
		permission["user"] = permissions[i].GetUser()
		permission["role"] = permissions[i].GetRole()
		permission["application"] = permissions[i].GetApplication()
		permissions_[i] = permission
	}

	jsonStr, err := json.Marshal(&permissions_)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.GetPermissionsRsp{
		Permissions: string(jsonStr),
	}, nil
}

func (self *Globule) deletePermissions(path string, owner string) error {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// First of all I will retreive the permissions for the given path...
	permissions, err := self.getPermissions(path)
	if err != nil {
		return err
	}

	// Get list of all permission with a given path.
	for i := 0; i < len(permissions); i++ {
		permission := permissions[i]
		if len(owner) > 0 {
			switch v := permission.GetOwner().(type) {
			case *resourcepb.ResourcePermission_User:
				if v.User == owner {
					err := p.DeleteOne(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+permission.GetPath()+`","owner":"`+v.User+`"}`, "")
					if err != nil {
						return err
					}
				}
			case *resourcepb.ResourcePermission_Role:
				if v.Role == owner {
					err := p.DeleteOne(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+permission.GetPath()+`","owner":"`+v.Role+`"}`, "")
					if err != nil {
						return err
					}
				}

			case *resourcepb.ResourcePermission_Application:
				if v.Application == owner {
					err := p.DeleteOne(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+permission.GetPath()+`","owner":"`+v.Application+`"}`, "")
					if err != nil {
						return err
					}
				}
			}
		} else {
			err := p.DeleteOne(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+permission.GetPath()+`"}`, "")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//* Delete a file permission *
func (self *Globule) DeletePermissions(ctx context.Context, rqst *resourcepb.DeletePermissionsRqst) (*resourcepb.DeletePermissionsRsp, error) {

	// That service made user of persistence service.
	err := self.deletePermissions(rqst.GetPath(), rqst.GetOwner())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeletePermissionsRsp{
		Result: true,
	}, nil
}

//* Validate a token *
func (self *Globule) ValidateToken(ctx context.Context, rqst *resourcepb.ValidateTokenRqst) (*resourcepb.ValidateTokenRsp, error) {
	clientId, _, expireAt, err := Interceptors.ValidateToken(rqst.Token)
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

/**
 * Validate application access by role
 */
func (self *Globule) validateApplicationAccess(name string, method string) error {
	if len(name) == 0 {
		return errors.New("No application was given to validate method access " + method)
	}

	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Applications", `{"path":"/`+name+`"}`, ``)
	if err != nil {
		return err
	}

	application := values.(map[string]interface{})

	err = errors.New("permission denied! application " + name + " cannot execute methode '" + method + "'")
	if application["actions"] == nil {
		return err
	}

	actions := []interface{}(application["actions"].(primitive.A))
	if actions == nil {
		return err
	}

	for i := 0; i < len(actions); i++ {
		if actions[i].(string) == method {
			return nil
		}
	}
	return err
}

/**
 * Validate user access by role
 */
func (self *Globule) validateUserAccess(userName string, method string) error {
	if len(userName) == 0 {
		return errors.New("No user  name was given to validate method access " + method)
	}

	// if the user is the super admin no validation is required.
	if userName == "sa" {
		return nil
	}

	// if guest can run the action...
	if self.canRunAction("guest", method) == nil {
		// everybody can run the action in that case.
		return nil
	}

	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// Now I will get the user roles and validate if the user can execute the
	// method.
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"name":"`+userName+`"}`, `[{"Projection":{"roles":1}}]`)
	if err != nil {
		return err
	}

	account := values.(map[string]interface{})
	roles := []interface{}(account["roles"].(primitive.A))
	for i := 0; i < len(roles); i++ {
		role := roles[i].(map[string]interface{})
		if self.canRunAction(role["$id"].(string), method) == nil {
			return nil
		}
	}

	return errors.New("permission denied! account " + userName + " cannot execute methode '" + method + "'")
}

// Test if a role can use action.
func (self *Globule) canRunAction(roleName string, method string) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Roles", `{"_id":"`+roleName+`"}`, `[{"Projection":{"actions":1}}]`)
	if err != nil {
		return err
	}

	role := values.(map[string]interface{})

	role["actions"] = []interface{}(role["actions"].(primitive.A))

	// append all action into the actions
	for i := 0; i < len(role["actions"].([]interface{})); i++ {
		if strings.ToLower(role["actions"].([]interface{})[i].(string)) == strings.ToLower(method) {
			return nil
		}
	}

	// Here I will test if the user has write to execute the methode.
	return errors.New("Permission denied!")
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
			userId, _, expired, err = Interceptors.ValidateToken(token)
		}
		return applicationId, userId, expired, err
	}
	return "", "", 0, fmt.Errorf("missing credentials")
}

func (self *Globule) isOwner(name string, path string) bool {

	// get the client...
	client, err := self.getPersistenceStore()
	if err != nil {
		return false
	}

	// If the user is the owner of the resource it has the permission
	count, err := client.Count(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+path+`","owner":"`+name+`"}`, ``)
	if err == nil {
		if count > 0 {
			return true
		}
	}

	return false
}

func (self *Globule) hasPermission(name string, path string, permission int32) (bool, int64) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return false, 0
	}

	// If the user is the owner of the resource it has all permission
	if self.isOwner(name, path) {
		return true, 0
	}

	count, err := p.Count(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+path+`"}`, ``)
	if err != nil {
		return false, 0
	}

	if count == 0 {
		count, err = p.Count(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+path+`"}`, ``)
		if err != nil {
			return false, 0
		}
	}

	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Permissions", `{"owner":"`+name+`", "path":"`+path+`"}`, ``)
	if err == nil {
		permissions := values.(map[string]interface{})

		p := permissions["permission"].(int32)

		// Here the owner have all permissions.
		if p == 7 {
			return true, count
		}

		if permission == p {
			return true, count
		}

		// Delete permission
		if permission == 1 {
			if p == 1 || p == 3 || p == 5 {
				return true, count
			}
		}

		// Write permission
		if permission == 2 {
			if p == 2 || p == 3 || p == 6 {
				return true, count
			}
		}

		// Read permission
		if permission == 4 {
			if p == 4 || p == 5 || p == 6 {
				return true, count
			}
		}

		return false, count
	}

	return false, count
}

/**
 * Validate if a user, a role or an application has write to do operation on a file or a directorty.
 */
func (self *Globule) validateUserResourceAccess(userName string, method string, path string, permission int32) error {

	if len(userName) == 0 {
		return errors.New("No user name was given to validate method access " + method)
	}

	// if the user is the super admin no validation is required.
	if userName == "sa" {
		return nil
	}

	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// Find the user role.
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"name":"`+userName+`"}`, `[{"Projection":{"roles":1}}]`)
	if err != nil {
		return err
	}

	account := values.(map[string]interface{})

	var count int64
	hasUserPermission, hasUserPermissionCount := self.hasPermission(userName, path, permission)
	if hasUserPermission {
		return nil
	}

	count += hasUserPermissionCount
	roles := []interface{}(account["roles"].(primitive.A))
	for i := 0; i < len(roles); i++ {
		role := roles[i].(map[string]interface{})
		hasRolePermission, hasRolePermissionCount := self.hasPermission(role["$id"].(string), path, permission)
		count += hasRolePermissionCount
		if hasRolePermission {
			return nil
		}
	}

	// If theres permission definied and we are here it's means the user dosent
	// have write to execute the action on the resourcepb.
	if count > 0 {
		return errors.New("Permission Denied for " + userName)
	}

	count, err = p.Count(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+path+`"}`, ``)
	if err != nil {
		if count > 0 {
			return errors.New("Permission Denied for " + userName)
		}
	}

	return nil
}

//* Validate if user can access a given file. *
func (self *Globule) ValidateUserResourceAccess(ctx context.Context, rqst *resourcepb.ValidateUserResourceAccessRqst) (*resourcepb.ValidateUserResourceAccessRsp, error) {

	path := rqst.GetPath() // The path of the resourcepb.

	// first of all I will validate the token.
	clientId, _, expiredAt, err := Interceptors.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if expiredAt < time.Now().Unix() {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The token is expired!")))
	}

	err = self.validateUserResourceAccess(clientId, rqst.Method, path, rqst.Permission)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.ValidateUserResourceAccessRsp{
		Result: true,
	}, nil
}

//* Validate if application can access a given file. *
func (self *Globule) ValidateApplicationResourceAccess(ctx context.Context, rqst *resourcepb.ValidateApplicationResourceAccessRqst) (*resourcepb.ValidateApplicationResourceAccessRsp, error) {

	path := rqst.GetPath()

	hasApplicationPermission, count := self.hasPermission(rqst.Name, path, rqst.Permission)
	if hasApplicationPermission {
		return &resourcepb.ValidateApplicationResourceAccessRsp{
			Result: true,
		}, nil
	}

	// If theres permission definied and we are here it's means the user dosent
	// have write to execute the action on the resourcepb.
	if count > 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Permission Denied for "+rqst.Name)))

	}

	return &resourcepb.ValidateApplicationResourceAccessRsp{
		Result: true,
	}, nil
}

//* Validate if user can access a given method. *
func (self *Globule) ValidateUserAccess(ctx context.Context, rqst *resourcepb.ValidateUserAccessRqst) (*resourcepb.ValidateUserAccessRsp, error) {

	// first of all I will validate the token.
	clientID, _, expiredAt, err := Interceptors.ValidateToken(rqst.Token)
	log.Println("----> clientID ", clientID, " expire at ", time.Unix(expiredAt, 0).String(), err)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if expiredAt < time.Now().Unix() {
		log.Println("----> token is expired!", expiredAt, time.Now().Unix())
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The token is expired!")))
	}

	// Here I will test if the user can run that function or not...
	err = self.validateUserAccess(clientID, rqst.Method)

	if err != nil {
		log.Println("----> ", rqst.Method, err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.ValidateUserAccessRsp{
		Result: true,
	}, nil
}

//* Validate if application can access a given method. *
func (self *Globule) ValidateApplicationAccess(ctx context.Context, rqst *resourcepb.ValidateApplicationAccessRqst) (*resourcepb.ValidateApplicationAccessRsp, error) {

	err := self.validateApplicationAccess(rqst.GetName(), rqst.Method)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.ValidateApplicationAccessRsp{
		Result: true,
	}, nil
}

//* Retrun a json string with all file info *
func (self *Globule) GetAllFilesInfo(ctx context.Context, rqst *resourcepb.GetAllFilesInfoRqst) (*resourcepb.GetAllFilesInfoRsp, error) {
	// That map will contain the list of all directories.
	dirs := make(map[string]map[string]interface{})

	err := filepath.Walk(self.webRoot,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				dir := make(map[string]interface{})
				dir["name"] = info.Name()
				dir["size"] = info.Size()
				dir["last_modified"] = info.ModTime().Unix()
				dir["path"] = strings.ReplaceAll(strings.Replace(path, self.path, "", -1), "\\", "/")
				dir["files"] = make([]interface{}, 0)
				dirs[dir["path"].(string)] = dir
				parent := dirs[dir["path"].(string)[0:strings.LastIndex(dir["path"].(string), "/")]]
				if parent != nil {
					parent["files"] = append(parent["files"].([]interface{}), dir)
				}
			} else {
				file := make(map[string]interface{})
				file["name"] = info.Name()
				file["size"] = info.Size()
				file["last_modified"] = info.ModTime().Unix()
				file["path"] = strings.ReplaceAll(strings.Replace(path, self.path, "", -1), "\\", "/")
				dir := dirs[file["path"].(string)[0:strings.LastIndex(file["path"].(string), "/")]]
				dir["files"] = append(dir["files"].([]interface{}), file)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	jsonStr, err := json.Marshal(dirs[strings.ReplaceAll(strings.Replace(self.webRoot, self.path, "", -1), "\\", "/")])
	if err != nil {
		return nil, err
	}
	return &resourcepb.GetAllFilesInfoRsp{Result: string(jsonStr)}, nil
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

	jsonStr, _ := Utility.ToJson(values)

	return &resourcepb.GetAllApplicationsInfoRsp{
		Result: jsonStr,
	}, nil

}

//* Create Permission for a dir (recursive) *
func (self *Globule) CreateDirPermissions(ctx context.Context, rqst *resourcepb.CreateDirPermissionsRqst) (*resourcepb.CreateDirPermissionsRsp, error) {

	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	clientId, _, expiredAt, err := Interceptors.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if expiredAt < time.Now().Unix() {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The token is expired!")))
	}

	// A new directory will take the parent permissions by default...
	path := rqst.GetPath()

	permissions, err := p.Find(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Now I will create the new permission of the created directory.
	for i := 0; i < len(permissions); i++ {
		// Copy the permission.
		permission := permissions[i].(map[string]interface{})
		permission_ := make(map[string]interface{}, 0)
		permission_["owner"] = permission["owner"]
		permission_["path"] = path + "/" + rqst.GetName()
		permission_["permission"] = permission["permission"].(int32)

		p.InsertOne(context.Background(), "local_resource", "local_resource", "Permissions", permission_, "")
	}

	resourceOwners, err := p.Find(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Now I will create the new permission of the created directory.
	for i := 0; i < len(resourceOwners); i++ {
		// Copye the permission.
		resourceOwner := resourceOwners[i].(map[string]interface{})
		resourceOwner_ := make(map[string]interface{}, 0)
		resourceOwner_["owner"] = resourceOwner["owner"]
		resourceOwner_["path"] = path + "/" + rqst.GetName()

		p.InsertOne(context.Background(), "local_resource", "local_resource", "ResourceOwners", resourceOwner_, "")
	}

	// The user who create a directory will be the owner of the
	// directory.
	if clientId != "sa" && clientId != "guest" && len(rqst.GetName()) > 0 {
		resourceOwner := make(map[string]interface{}, 0)
		resourceOwner["owner"] = clientId
		resourceOwner["path"] = path + "/" + rqst.GetName()
		resourceOwnerStr, _ := Utility.ToJson(resourceOwner)
		p.ReplaceOne(context.Background(), "local_resource", "local_resource", "ResourceOwners", resourceOwnerStr, resourceOwnerStr, `[{"upsert":true}]`)
	}

	return &resourcepb.CreateDirPermissionsRsp{
		Result: true,
	}, nil
}

//* Rename file/dir permission *
func (self *Globule) RenameFilePermission(ctx context.Context, rqst *resourcepb.RenameFilePermissionRqst) (*resourcepb.RenameFilePermissionRsp, error) {

	path := rqst.GetPath()
	path = strings.ReplaceAll(path, "\\", "/")
	if len(path) > 1 {
		if strings.HasSuffix(path, "/") {
			path = path[0 : len(path)-2]
		}
	}

	oldPath := rqst.OldName
	newPath := rqst.NewName

	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if strings.HasPrefix(path, "/") {
		if len(path) > 1 {
			oldPath = path + "/" + rqst.OldName
			newPath = path + "/" + rqst.NewName

		} else {
			oldPath = "/" + rqst.OldName
			newPath = "/" + rqst.NewName
		}
	}

	// Replace permission path... regex "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"} not work.
	permissions, err := p.Find(context.Background(), "local_resource", "local_resource", "Permissions", `{}`, "")
	if err == nil {
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), oldPath) {
				path := newPath + permission["path"].(string)[len(oldPath):]
				err := p.Update(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+permission["path"].(string)+`"}`, `{"$set":{"path":"`+path+`"}}`, "")
				if err != nil {
					return nil, status.Errorf(
						codes.Internal,
						Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
				}
			}
		}
	}

	// Replace file owner path... regex not work... "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"}
	permissions, err = p.Find(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{}`, "")
	if err == nil {
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), oldPath) {
				path := newPath + permission["path"].(string)[len(oldPath):]
				err = p.Update(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+permission["path"].(string)+`"}`, `{"$set":{"path":"`+path+`"}}`, "")
				if err != nil {
					return nil, status.Errorf(
						codes.Internal,
						Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
				}
			}
		}
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &resourcepb.RenameFilePermissionRsp{
		Result: true,
	}, nil
}

func (self *Globule) deleteDirPermissions(path string) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	path = strings.ReplaceAll(path, "\\", "/")
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-2]
	}

	// Replace permission path... regex "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"} not work.
	permissions, err := p.Find(context.Background(), "local_resource", "local_resource", "Permissions", `{}`, "")
	if err == nil {
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), path) {
				err := p.Delete(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+permission["path"].(string)+`"}`, "")
				if err != nil {
					return err
				}
			}
		}
	}

	// Replace file owner path... regex not work... "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"}
	permissions, err = p.Find(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{}`, "")
	if err == nil {
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), path) {
				err = p.Delete(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+permission["path"].(string)+`"}`, "")
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//* Delete Permission for a dir (recursive) *
func (self *Globule) DeleteDirPermissions(ctx context.Context, rqst *resourcepb.DeleteDirPermissionsRqst) (*resourcepb.DeleteDirPermissionsRsp, error) {
	err := self.deleteDirPermissions(rqst.GetPath())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteDirPermissionsRsp{
		Result: true,
	}, nil
}

//* Delete a single file permission *
func (self *Globule) DeleteFilePermissions(ctx context.Context, rqst *resourcepb.DeleteFilePermissionsRqst) (*resourcepb.DeleteFilePermissionsRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	path := rqst.GetPath()

	err = p.Delete(context.Background(), "local_resource", "local_resource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = p.Delete(context.Background(), "local_resource", "local_resource", "ResourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteFilePermissionsRsp{
		Result: true,
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
func (self *Globule) AddPeerAction(ctx context.Context, rqst *resourcepb.AddPeerActionRqst) (*resourcepb.AddPeerActionRsp, error) {
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
		peer["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		for i := 0; i < len(peer["actions"].(primitive.A)); i++ {

			if peer["actions"].(primitive.A)[i].(string) == rqst.Action {
				exist = true
				break
			}
		}
		if !exist {
			peer["actions"] = append(peer["actions"].(primitive.A), rqst.Action)
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Peer named "+rqst.Domain+" already contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr, _ := json.Marshal(peer)
		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Peers", `{"_id":"`+_id+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.AddPeerActionRsp{Result: true}, nil

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
		jsonStr, _ := json.Marshal(peer)
		err := p.ReplaceOne(context.Background(), "local_resource", "local_resource", "Peers", `{"_id":"`+_id+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &resourcepb.RemovePeerActionRsp{Result: true}, nil
}

/**
 * Validate application access by role
 */
func (self *Globule) validatePeerAccess(domain string, method string) error {

	if len(domain) == 0 {
		return errors.New("No domain was given to validate peer method access " + method)
	}

	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	_id := Utility.GenerateUUID(domain)
	values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Peers", `{"_id":"`+_id+`"}`, ``)
	if err != nil {
		return err
	}

	peer := values.(map[string]interface{})

	err = errors.New("permission denied! peer with domain " + domain + " cannot execute methode '" + method + "'")
	if peer["actions"] == nil {
		return err
	}

	actions := []interface{}(peer["actions"].(primitive.A))
	if actions == nil {
		return err
	}

	for i := 0; i < len(actions); i++ {
		if actions[i].(string) == method {
			return nil
		}
	}

	return err
}

//* Validate if a peer can access a given resource. *
func (self *Globule) ValidatePeerResourceAccess(ctx context.Context, rqst *resourcepb.ValidatePeerResourceAccessRqst) (*resourcepb.ValidatePeerResourceAccessRsp, error) {
	path := rqst.GetPath()

	hasPermission, count := self.hasPermission(rqst.Domain, path, rqst.Permission)
	if hasPermission {
		return &resourcepb.ValidatePeerResourceAccessRsp{
			Result: true,
		}, nil
	}

	// If theres permission definied and we are here it's means the user dosent
	// have write to execute the action on the resource.
	if count > 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Permission Denied for "+rqst.Domain)))

	}

	return &resourcepb.ValidatePeerResourceAccessRsp{
		Result: true,
	}, nil
}

//* Validate if a peer can access a given method. *
func (self *Globule) ValidatePeerAccess(ctx context.Context, rqst *resourcepb.ValidatePeerAccessRqst) (*resourcepb.ValidatePeerAccessRsp, error) {
	err := self.validatePeerAccess(rqst.Domain, rqst.Method)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.ValidatePeerAccessRsp{
		Result: true,
	}, nil
}
