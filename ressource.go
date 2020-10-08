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

	"io/ioutil"

	"path/filepath"

	"net"

	"strconv"

	"os/signal"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/services/golang/ressource/ressourcepb"
	"github.com/davecourtois/Utility"
	"github.com/emicklei/proto"
	"github.com/golang/protobuf/jsonpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (self *Globule) startRessourceService() error {
	ressource_server, err := self.startInternalService(string(ressourcepb.File_services_proto_ressource_proto.Services().Get(0).FullName()), ressourcepb.File_services_proto_ressource_proto.Path(), self.RessourcePort, self.RessourceProxy, self.Protocol == "https", self.unaryRessourceInterceptor, self.streamRessourceInterceptor)
	if err == nil {

		// Create the channel to listen on ressource port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.RessourcePort))
		if err != nil {
			log.Fatalf("could not start ressource service %s: %s", self.getDomain(), err)
		}

		ressourcepb.RegisterRessourceServiceServer(ressource_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			go func() {
				// no web-rpc server.
				if err = ressource_server.Serve(lis); err != nil {
					log.Println(err)
				}
			}()

			// Wait for signal to stop.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch

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

	s := self.Services[name]
	if s == nil {
		s = self.Services[name+"_server"]
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
	err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"sa"}`, jsonStr, `[{"upsert":true}]`)
	if err != nil {
		return err
	}

	log.Println("role sa with was updated!")

	// I will also create the guest role, the basic one
	count, err := p.Count(context.Background(), "local_ressource", "local_ressource", "Roles", `{ "_id":"guest"}`, "")
	guest := make(map[string]interface{})
	if err != nil {
		return err
	} else if count == 0 {
		log.Println("need to create roles guest...")
		guest["_id"] = "guest"
		guest["name"] = "guest"
		guest["actions"] = []string{
			"/admin.AdminService/GetConfig",
			"/ressource.RessourceService/RegisterAccount",
			"/ressource.RessourceService/AccountExist",
			"/ressource.RessourceService/Authenticate",
			"/ressource.RessourceService/RefreshToken",
			"/ressource.RessourceService/GetPermissions",
			"/ressource.RessourceService/GetAllFilesInfo",
			"/ressource.RessourceService/GetAllApplicationsInfo",
			"/ressource.RessourceService/GetRessourceOwners",
			"/ressource.RessourceService/ValidateUserAccess",
			"/ressource.RessourceService/ValidateUserRessourceAccess",
			"/ressource.RessourceService/ValidateApplicationAccess",
			"/ressource.RessourceService/ValidateApplicationRessourceAccess",
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
			"/ressource.RessourceService/GetAllActions"}

		_, err := p.InsertOne(context.Background(), "local_ressource", "local_ressource", "Roles", guest, "")
		if err != nil {
			return err
		}
		log.Println("role guest was created!")
	}

	// Create connection application.
	self.createApplicationConnection()

	// Here I will also set permssion for local services.
	for i := 0; i < len(self.actionPermissions); i++ {
		permission := self.actionPermissions[i].(map[string]interface{})
		self.setActionPermission(permission["action"].(string), int32(Utility.ToInt(permission["permission"])))
	}

	return nil
}

func (self *Globule) createApplicationConnection() error {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	store, err := self.getPersistenceStore()
	applications, _ := store.Find(context.Background(), "local_ressource", "local_ressource", "Applications", "{}", "")
	if err != nil {
		return err
	}

	for i := 0; i < len(applications); i++ {
		application := applications[i].(map[string]interface{})
		// Open the user database connection.
		err = p.CreateConnection(application["_id"].(string)+"_db", application["_id"].(string)+"_db", "0.0.0.0", 27017, 0, application["_id"].(string), application["password"].(string), 5000, "", false)
		if err != nil {
			return err
		}
	}

	return nil
}

/** Append new LDAP synchronization informations. **/
func (self *Globule) SynchronizeLdap(ctx context.Context, rqst *ressourcepb.SynchronizeLdapRqst) (*ressourcepb.SynchronizeLdapRsp, error) {

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
						rolesStr := `[{"$ref":"Roles","$id":"guest","$db":"local_ressource"}`
						for j := 0; j < len(accountsInfo[i][3].([]interface{})); j++ {
							rolesStr += `,{"$ref":"Roles","$id":"` + Utility.GenerateUUID(accountsInfo[i][3].([]interface{})[j].(string)) + `","$db":"local_ressource"}`
						}

						rolesStr += `]`
						err := p.UpdateOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+id+`"}`, `{ "$set":{"roles":`+rolesStr+`}}`, "")
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

	return &ressourcepb.SynchronizeLdapRsp{
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
	count, err := p.Count(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"name":"`+name+`"}`, "")
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
		role["$db"] = "local_ressource"
		account["roles"] = append(account["roles"].([]map[string]interface{}), role)
	}

	// Here I will insert the account in the database.
	_, err = p.InsertOne(context.Background(), "local_ressource", "local_ressource", "Accounts", account, "")

	// replace @ and . by _
	name = strings.ReplaceAll(strings.ReplaceAll(name, "@", "_"), ".", "_")

	// Each account will have their own database and a use that can read and write
	// into it.
	// Here I will wrote the script for mongoDB...
	createUserScript := fmt.Sprintf(
		"db=db.getSiblingDB('%s_db');db.createCollection('user_data');db=db.getSiblingDB('admin');db.createUser({user: '%s', pwd: '%s',roles: [{ role: 'dbOwner', db: '%s_db' }]});",
		name, name, password, name)

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd(context.Background(), "local_ressource", "sa", self.RootPassword, createUserScript)
	if err != nil {
		return err
	}

	p_, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	err = p_.CreateConnection(name+"_db", name+"_db", "0.0.0.0", 27017, 0, name, password, 5000, "", false)
	if err != nil {
		return errors.New("No persistence service are available to store ressource information.")
	}

	return nil

}

/* Register a new Account */
func (self *Globule) RegisterAccount(ctx context.Context, rqst *ressourcepb.RegisterAccountRqst) (*ressourcepb.RegisterAccountRsp, error) {
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
	_, err = p.InsertOne(context.Background(), "local_ressource", "local_ressource", "Tokens", map[string]interface{}{"_id": name, "expireAt": Utility.ToString(expireAt)}, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Now I will
	return &ressourcepb.RegisterAccountRsp{
		Result: tokenString, // Return the token string.
	}, nil
}

func (self *Globule) AccountExist(ctx context.Context, rqst *ressourcepb.AccountExistRqst) (*ressourcepb.AccountExistRsp, error) {
	var exist bool
	// Get the persistence connection
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}
	// Test with the _id
	count, _ := p.Count(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.Id+`"}`, "")
	if count > 0 {
		exist = true
	}

	// Test with the name
	if !exist {
		count, _ := p.Count(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"name":"`+rqst.Id+`"}`, "")
		if count > 0 {
			exist = true
		}
	}

	// Test with the email.
	if !exist {
		count, _ := p.Count(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"email":"`+rqst.Id+`"}`, "")
		if count > 0 {
			exist = true
		}
	}
	if exist {
		return &ressourcepb.AccountExistRsp{
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
// to validate permission over the requested ressourcepb.
func (self *Globule) Authenticate(ctx context.Context, rqst *ressourcepb.AuthenticateRqst) (*ressourcepb.AuthenticateRsp, error) {
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

		/** Return the token only **/
		return &ressourcepb.AuthenticateRsp{
			Token: tokenString,
		}, nil
	}

	values, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"name":"`+rqst.Name+`"}`, "")
	if len(values) == 0 {
		values, err = p.Find(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"email":"`+rqst.Name+`"}`, "")
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
	err = p_.CreateConnection(name_+"_db", name_+"_db", "0.0.0.0", 27017, 0, name_, rqst.Password, 5000, "", false)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No persistence service are available to store ressource information.")))
	}

	// save the newly create token into the database.
	name, _, expireAt, _ := Interceptors.ValidateToken(tokenString)
	err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`"}`, `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// So here I will start a timer

	// Here I got the token I will now put it in the cache.
	return &ressourcepb.AuthenticateRsp{
		Token: tokenString,
	}, nil
}

/**
 * Refresh token get a new token.
 */
func (self *Globule) RefreshToken(ctx context.Context, rqst *ressourcepb.RefreshTokenRqst) (*ressourcepb.RefreshTokenRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// first of all I will validate the current token.
	name, email, expireAt, _ := Interceptors.ValidateToken(rqst.Token)

	// If the token is older than seven day without being refresh then I retrun an error.
	if time.Unix(expireAt, 0).Before(time.Now().AddDate(0, 0, -7)) {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The token cannot be refresh after 7 day")))
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`"}`, `[{"Projection":{"expireAt":1}}]`)
	if err == nil {
		lastTokenInfo := values.(map[string]interface{})
		savedTokenExpireAt := time.Unix(int64(lastTokenInfo["expireAt"].(int32)), 0)

		// That mean a newer token was already refresh.
		if savedTokenExpireAt.Before(time.Unix(expireAt, 0)) {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("That token cannot not be refresh because a newer one already exist. You need to re-authenticate in order to get a new token.")))
		}
	}

	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, name, email)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// get back the new expireAt
	name, _, expireAt, _ = Interceptors.ValidateToken(tokenString)

	err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`"}`, `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, `[{"upsert":true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// return the token string.
	return &ressourcepb.RefreshTokenRsp{
		Token: tokenString,
	}, nil
}

//* Delete an account *
func (self *Globule) DeleteAccount(ctx context.Context, rqst *ressourcepb.DeleteAccountRqst) (*ressourcepb.DeleteAccountRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.Id+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	account := values.(map[string]interface{})

	// Try to delete the account...
	err = p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.Id+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Delete permissions
	err = p.Delete(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"owner":"`+rqst.Id+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Delete the token.
	err = p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "Tokens", `{"_id":"`+rqst.Id+`"}`, "")
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
	err = p.RunAdminCmd(context.Background(), "local_ressource", "sa", self.RootPassword, dropUserScript)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Remove the user database.
	err = p.DeleteDatabase(context.Background(), "local_ressource", name+"_db")
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

	return &ressourcepb.DeleteAccountRsp{
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
	_, err = p.FindOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"`+id+`"}`, ``)
	if err == nil {
		return errors.New("Role named " + name + "already exist!")
	}

	// Here will create the new role.
	role := make(map[string]interface{}, 0)
	role["_id"] = id
	role["name"] = name
	role["actions"] = actions

	_, err = p.InsertOne(context.Background(), "local_ressource", "local_ressource", "Roles", role, "")
	if err != nil {
		return err
	}

	return nil
}

//* Create a role with given action list *
func (self *Globule) CreateRole(ctx context.Context, rqst *ressourcepb.CreateRoleRqst) (*ressourcepb.CreateRoleRsp, error) {
	// That service made user of persistence service.
	err := self.createRole(rqst.Role.Id, rqst.Role.Name, rqst.Role.Actions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.CreateRoleRsp{Result: true}, nil
}

//* Delete a role with a given id *
func (self *Globule) DeleteRole(ctx context.Context, rqst *ressourcepb.DeleteRoleRqst) (*ressourcepb.DeleteRoleRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	accounts, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Accounts", `{}`, ``)
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

					err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"name":"`+account["name"].(string)+`"}`, jsonStr, ``)
					if err != nil {
						return nil, status.Errorf(
							codes.Internal,
							Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
					}
				}
			}
		}
	}

	err = p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Delete permissions
	err = p.Delete(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"owner":"`+rqst.RoleId+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeleteRoleRsp{Result: true}, nil
}

//* Append an action to existing role. *
func (self *Globule) AddRoleAction(ctx context.Context, rqst *ressourcepb.AddRoleActionRqst) (*ressourcepb.AddRoleActionRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
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
		err := p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressourcepb.AddRoleActionRsp{Result: true}, nil
}

//* Remove an action to existing role. *
func (self *Globule) RemoveRoleAction(ctx context.Context, rqst *ressourcepb.RemoveRoleActionRqst) (*ressourcepb.RemoveRoleActionRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
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
		err := p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressourcepb.RemoveRoleActionRsp{Result: true}, nil
}

//* Add role to a given account *
func (self *Globule) AddAccountRole(ctx context.Context, rqst *ressourcepb.AddAccountRoleRqst) (*ressourcepb.AddAccountRoleRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, ``)
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
		account["roles"] = append(roles, map[string]interface{}{"$ref": "Roles", "$id": rqst.RoleId, "$db": "local_ressource"})

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

		err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, jsonStr, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

	}
	return &ressourcepb.AddAccountRoleRsp{Result: true}, nil
}

//* Remove a role from a given account *
func (self *Globule) RemoveAccountRole(ctx context.Context, rqst *ressourcepb.RemoveAccountRoleRqst) (*ressourcepb.RemoveAccountRoleRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, ``)
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

		err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, jsonStr, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

	}
	return &ressourcepb.RemoveAccountRoleRsp{Result: true}, nil
}

//* Append an action to existing application. *
func (self *Globule) AddApplicationAction(ctx context.Context, rqst *ressourcepb.AddApplicationActionRqst) (*ressourcepb.AddApplicationActionRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, ``)
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
		err := p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressourcepb.AddApplicationActionRsp{Result: true}, nil
}

//* Remove an action to existing application. *
func (self *Globule) RemoveApplicationAction(ctx context.Context, rqst *ressourcepb.RemoveApplicationActionRqst) (*ressourcepb.RemoveApplicationActionRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, ``)
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
		err := p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressourcepb.RemoveApplicationActionRsp{Result: true}, nil
}

//* Delete an application from the server. *
func (self *Globule) DeleteApplication(ctx context.Context, rqst *ressourcepb.DeleteApplicationRqst) (*ressourcepb.DeleteApplicationRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, ``)
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
	err = p.DeleteDatabase(context.Background(), "local_ressource", rqst.ApplicationId+"_db")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Finaly I will remove the entry in  the table.
	err = p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Delete permissions
	err = p.Delete(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"owner":"`+rqst.ApplicationId+`"}`, "")
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
	err = p.RunAdminCmd(context.Background(), "local_ressource", "sa", self.RootPassword, dropUserScript)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeleteApplicationRsp{
		Result: true,
	}, nil
}

func (self *Globule) deleteAccountPermissions(name string) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	err = p.Delete(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"name":"`+name+`"}`, "")
	if err != nil {
		return err
	}

	return nil
}

//* Delete all permission for a given account *
func (self *Globule) DeleteAccountPermissions(ctx context.Context, rqst *ressourcepb.DeleteAccountPermissionsRqst) (*ressourcepb.DeleteAccountPermissionsRsp, error) {

	err := self.deleteAccountPermissions(rqst.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeleteAccountPermissionsRsp{
		Result: true,
	}, nil
}

func (self *Globule) deleteRolePermissions(name string) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	err = p.Delete(context.Background(), "local_ressource", "local_ressource", "Roles", `{"name":"`+name+`"}`, "")
	if err != nil {
		return err
	}

	return nil
}

//* Delete all permission for a given role *
func (self *Globule) DeleteRolePermissions(ctx context.Context, rqst *ressourcepb.DeleteRolePermissionsRqst) (*ressourcepb.DeleteRolePermissionsRsp, error) {

	err := self.deleteRolePermissions(rqst.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeleteRolePermissionsRsp{
		Result: true,
	}, nil
}

/**
 * Return the list of all actions avalaible on the server.
 */
func (self *Globule) GetAllActions(ctx context.Context, rqst *ressourcepb.GetAllActionsRqst) (*ressourcepb.GetAllActionsRsp, error) {
	return &ressourcepb.GetAllActionsRsp{Actions: self.methods}, nil
}

/////////////////////// Ressource permission owner /////////////////////////////
func (self *Globule) setRessourceOwner(owner string, path string) error {
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

	// here I if the ressource is a directory I will set the permission on
	// subdirectory and files...
	fileInfo, err := os.Stat(self.GetAbsolutePath(path))
	if err == nil {
		if fileInfo.IsDir() {
			files, err := ioutil.ReadDir(self.GetAbsolutePath(path))
			if err == nil {
				for i := 0; i < len(files); i++ {
					file := files[i]
					if strings.HasSuffix(path, "/") {
						self.setRessourceOwner(owner, path+file.Name())
					} else {
						self.setRessourceOwner(owner, path+"/"+file.Name())
					}
				}
			} else {
				return err
			}
		}
	}

	// Here I will set ressources whit that path, be sure to have different
	// path than application and webroot path if you dont want permission follow each other.
	ressources, err := self.getRessources(path)
	if err == nil {
		for i := 0; i < len(ressources); i++ {
			if ressources[i].GetPath() != path {
				path_ := ressources[i].GetPath()[len(path)+1:]
				paths := strings.Split(path_, "/")
				path_ = path
				// set sub-path...
				for j := 0; j < len(paths); j++ {
					path_ += "/" + paths[j]
					ressourceOwner := make(map[string]interface{})
					ressourceOwner["owner"] = owner
					ressourceOwner["path"] = path_
					// force the id to be the same for ressource with the same owner and path.
					ressourceOwner["_id"] = Utility.GenerateUUID(owner + path_)

					// Here if the
					jsonStr, err := Utility.ToJson(&ressourceOwner)
					if err != nil {
						return err
					}

					err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", jsonStr, jsonStr, `[{"upsert":true}]`)
					if err != nil {
						return err
					}
				}
			}

			if strings.HasSuffix(ressources[i].GetPath(), "/") {
				self.setRessourceOwner(owner, ressources[i].GetPath()+ressources[i].GetName())
			} else {
				self.setRessourceOwner(owner, ressources[i].GetPath()+"/"+ressources[i].GetName())
			}

		}
	}

	// Here I will set the ressource owner.
	ressourceOwner := make(map[string]interface{})
	ressourceOwner["owner"] = owner
	ressourceOwner["path"] = path
	ressourceOwner["_id"] = Utility.GenerateUUID(owner + path)

	// Here if the
	jsonStr, err := Utility.ToJson(&ressourceOwner)
	if err != nil {
		return err
	}

	err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", jsonStr, jsonStr, `[{"upsert":true}]`)
	if err != nil {
		return err
	}

	return nil
}

//* Set Ressource owner *
func (self *Globule) SetRessourceOwner(ctx context.Context, rqst *ressourcepb.SetRessourceOwnerRqst) (*ressourcepb.SetRessourceOwnerRsp, error) {
	// That service made user of persistence service.
	path := rqst.GetPath()

	err := self.setRessourceOwner(rqst.GetOwner(), path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.SetRessourceOwnerRsp{
		Result: true,
	}, nil
}

//* Get the ressource owner *
func (self *Globule) GetRessourceOwners(ctx context.Context, rqst *ressourcepb.GetRessourceOwnersRqst) (*ressourcepb.GetRessourceOwnersRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Get the absolute path
	path := rqst.GetPath()

	// find the ressource with it id
	ressourceOwners, err := p.Find(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	owners := make([]string, 0)
	for i := 0; i < len(ressourceOwners); i++ {
		owners = append(owners, ressourceOwners[i].(map[string]interface{})["owner"].(string))
	}

	return &ressourcepb.GetRessourceOwnersRsp{
		Owners: owners,
	}, nil
}

//* Get the ressource owner *
func (self *Globule) DeleteRessourceOwner(ctx context.Context, rqst *ressourcepb.DeleteRessourceOwnerRqst) (*ressourcepb.DeleteRessourceOwnerRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	path := rqst.GetPath()

	// Delete the ressource owner for a given path.
	err = p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`","owner":"`+rqst.GetOwner()+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeleteRessourceOwnerRsp{
		Result: true,
	}, nil
}

func (self *Globule) DeleteRessourceOwners(ctx context.Context, rqst *ressourcepb.DeleteRessourceOwnersRqst) (*ressourcepb.DeleteRessourceOwnersRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	path := rqst.GetPath()

	// delete the ressource owners with it path
	err = p.Delete(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeleteRessourceOwnersRsp{
		Result: true,
	}, nil
}

//////////////////////////// Loggin info ///////////////////////////////////////
func (self *Globule) logServiceInfo(service string, message string) error {

	// Here I will use event to publish log information...
	info := new(ressourcepb.LogInfo)
	info.Application = ""
	info.UserId = "globular"
	info.UserName = "globular"
	info.Method = service
	info.Date = time.Now().Unix()
	info.Message = message
	info.Type = ressourcepb.LogType_ERROR_MESSAGE // not necessarely errors..
	self.log(info)

	return nil
}

// Log err and info...
func (self *Globule) logInfo(application string, method string, token string, err_ error) error {

	// Remove cyclic calls
	if method == "/ressource.RessourceService/Log" {
		return errors.New("Method " + method + " cannot not be log because it will cause a circular call to itself!")
	}

	// Here I will use event to publish log information...
	info := new(ressourcepb.LogInfo)
	info.Application = application
	info.UserId = token
	info.UserName = token
	info.Method = method
	info.Date = time.Now().Unix()
	if err_ != nil {
		info.Message = err_.Error()
		info.Type = ressourcepb.LogType_ERROR_MESSAGE
		logger.Error(info.Message)
	} else {
		info.Type = ressourcepb.LogType_INFO_MESSAGE
		logger.Info(info.Message)
	}

	self.log(info)

	return nil
}

// unaryInterceptor calls authenticateClient with current context
func (self *Globule) unaryRessourceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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
	if method == "/ressource.RessourceService/GetAllActions" ||
		method == "/ressource.RessourceService/RegisterAccount" ||
		method == "/ressource.RessourceService/AccountExist" ||
		method == "/ressource.RessourceService/Authenticate" ||
		method == "/ressource.RessourceService/RefreshToken" ||
		method == "/ressource.RessourceService/GetPermissions" ||
		method == "/ressource.RessourceService/GetAllFilesInfo" ||
		method == "/ressource.RessourceService/GetAllApplicationsInfo" ||
		method == "/ressource.RessourceService/GetRessourceOwners" ||
		method == "/ressource.RessourceService/ValidateToken" ||
		method == "/ressource.RessourceService/ValidateUserAccess" ||
		method == "/ressource.RessourceService/ValidateApplicationAccess" ||
		method == "/ressource.RessourceService/ValidateApplicationRessourceAccess" ||
		method == "/ressource.RessourceService/ValidateUserRessourceAccess" ||
		method == "/ressource.RessourceService/GetActionPermission" ||
		method == "/ressource.RessourceService/Log" ||
		method == "/ressource.RessourceService/GetLog" {
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
			// special case that need ownership of the ressource or be sa
			if method == "/ressource.RessourceService/SetPermission" || method == "/ressource.RessourceService/DeletePermissions" ||
				method == "/ressource.RessourceService/SetRessourceOwner" || method == "/ressource.RessourceService/DeleteRessourceOwner" ||
				method == "/ressource.RessourceService/CreateDirPermissions" {
				var path string
				if method == "/ressource.RessourceService/SetPermission" {
					rqst := req.(*ressourcepb.SetPermissionRqst)
					// Here I will validate that the user is the owner.
					path = rqst.Permission.GetPath()
				} else if method == "/ressource.RessourceService/DeletePermissions" {
					rqst := req.(*ressourcepb.DeletePermissionsRqst)
					// Here I will validate that the user is the owner.
					path = rqst.Path
				} else if method == "/ressource.RessourceService/SetRessourceOwner" {
					rqst := req.(*ressourcepb.SetRessourceOwnerRqst)
					// Here I will validate that the user is the owner.
					path = rqst.Path
				} else if method == "/ressource.RessourceService/DeleteRessourceOwner" {
					rqst := req.(*ressourcepb.DeleteRessourceOwnerRqst)
					// Here I will validate that the user is the owner.
					path = rqst.Path
				} else if method == "/ressource.RessourceService/CreateDirPermissions" {
					rqst := req.(*ressourcepb.CreateDirPermissionsRqst)
					// Here I will validate that the user is the owner.
					path = rqst.Path
				}

				// If the use is the ressource owner he can run the method
				if self.isOwner(clientId, path) {
					hasAccess = true
				}

			} else {
				err = self.validateUserAccess(clientId, method)
				if err == nil {
					hasAccess = true
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
		// Set permissions in case one of those methode is called.
		if method == "/ressource.RessourceService/DeleteApplication" {
			rqst := req.(*ressourcepb.DeleteApplicationRqst)
			err := self.deleteDirPermissions("/" + rqst.ApplicationId)
			if err != nil {
				return nil, err
			}
		} else if method == "/ressource.RessourceService/DeleteRole" {
			rqst := req.(*ressourcepb.DeleteRoleRqst)
			err := self.deleteRolePermissions("/" + rqst.RoleId)
			if err != nil {
				return nil, err
			}
		} else if method == "/ressource.RessourceService/DeleteAccount" {
			rqst := req.(*ressourcepb.DeleteAccountRqst)
			err := self.deleteAccountPermissions("/" + rqst.Id)
			if err != nil {
				return nil, err
			}
		} else if method == "/ressource.RessourceService/SetRessource" {
			if clientId != "sa" {
				rqst := req.(*ressourcepb.SetRessourceRqst)

				// In that case i will set the userId as the owner of the ressourcepb.
				err := self.setRessourceOwner(clientId, rqst.Ressource.Path+"/"+rqst.Ressource.Name)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return result, err

}

// Stream interceptor.
func (self *Globule) streamRessourceInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	err := handler(srv, stream)
	if err != nil {
		return err
	}

	return nil
}

func (self *Globule) getLogInfoKeyValue(info *ressourcepb.LogInfo) (string, string, error) {
	marshaler := new(jsonpb.Marshaler)
	jsonStr, err := marshaler.MarshalToString(info)
	if err != nil {
		return "", "", err
	}

	key := ""
	if info.GetType() == ressourcepb.LogType_INFO_MESSAGE {
		// Increnment prometheus counter,
		self.methodsCounterLog.WithLabelValues("INFO", info.Method).Inc()

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
		// Increnment prometheus counter,
		self.methodsCounterLog.WithLabelValues("ERROR", info.Method).Inc()
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

func (self *Globule) log(info *ressourcepb.LogInfo) error {

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
func (self *Globule) Log(ctx context.Context, rqst *ressourcepb.LogRqst) (*ressourcepb.LogRsp, error) {
	// Publish event...
	self.log(rqst.Info)

	return &ressourcepb.LogRsp{
		Result: true,
	}, nil
}

//* Log error or information into the data base *
// Retreive log infos (the query must be something like /infos/'date'/'applicationName'/'userName'
func (self *Globule) GetLog(rqst *ressourcepb.GetLogRqst, stream ressourcepb.RessourceService_GetLogServer) error {

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

	infos := make([]*ressourcepb.LogInfo, 0)
	i := 0
	max := 100
	for jsonDecoder.More() {
		info := ressourcepb.LogInfo{}
		err := jsonpb.UnmarshalNext(jsonDecoder, &info)
		if err != nil {
			return err
		}
		// append the info inside the stream.
		infos = append(infos, &info)
		if i == max-1 {
			// I will send the stream at each 100 logs...
			rsp := &ressourcepb.GetLogRsp{
				Info: infos,
			}
			// Send the infos
			err = stream.Send(rsp)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			infos = make([]*ressourcepb.LogInfo, 0)
			i = 0
		}
		i++
	}

	// Send the last infos...
	if len(infos) > 0 {
		rsp := &ressourcepb.GetLogRsp{
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
		info := ressourcepb.LogInfo{}

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
func (self *Globule) DeleteLog(ctx context.Context, rqst *ressourcepb.DeleteLogRqst) (*ressourcepb.DeleteLogRsp, error) {

	key, _, _ := self.getLogInfoKeyValue(rqst.Log)
	err := self.logs.RemoveItem(key)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeleteLogRsp{
		Result: true,
	}, nil
}

//* Clear logs. info or errors *
func (self *Globule) ClearAllLog(ctx context.Context, rqst *ressourcepb.ClearAllLogRqst) (*ressourcepb.ClearAllLogRsp, error) {
	var err error

	if rqst.Type == ressourcepb.LogType_ERROR_MESSAGE {
		err = self.deleteLog("/errors/*")
	} else {
		err = self.deleteLog("/infos/*")
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.ClearAllLogRsp{
		Result: true,
	}, nil
}

///////////////////////  ressource management. /////////////////
func (self *Globule) setRessource(r *ressourcepb.Ressource) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	var marshaler jsonpb.Marshaler

	jsonStr, err := marshaler.MarshalToString(r)
	if err != nil {
		return err
	}

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(r.Path + r.Name)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// Always create a new if not already exist.
	err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Ressources", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return err
	}

	return nil
}

//* Set a ressource from a client (custom service) to globular
func (self *Globule) SetRessource(ctx context.Context, rqst *ressourcepb.SetRessourceRqst) (*ressourcepb.SetRessourceRsp, error) {
	if rqst.Ressource == nil {

		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no ressource was given!")))

	}
	err := self.setRessource(rqst.Ressource)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.SetRessourceRsp{
		Result: true,
	}, nil
}

func (self *Globule) getRessources(path string) ([]*ressourcepb.Ressource, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	data, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Ressources", `{}`, `[{"Projection":{"_id":0}}]`)
	if err != nil {
		return nil, err
	}

	ressources := make([]*ressourcepb.Ressource, 0)

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
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No path was given for the ressource!")))
		}

		if res["name"] == nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No name was given for the ressource!")))
		}

		// append the info inside the stream.
		if strings.HasPrefix(res["path"].(string), path) {
			ressources = append(ressources, &ressourcepb.Ressource{Path: res["path"].(string), Name: res["name"].(string), Modified: modified, Size: size})
		}
	}
	return ressources, nil
}

//* Get all ressources
func (self *Globule) GetRessources(rqst *ressourcepb.GetRessourcesRqst, stream ressourcepb.RessourceService_GetRessourcesServer) error {
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

	data, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Ressources", queryStr, ``)
	if err != nil {
		return err
	}

	ressources := make([]*ressourcepb.Ressource, 0)
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
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No path was given for the ressource!")))
		}

		if res["name"] == nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No name was given for the ressource!")))
		}

		// append the info inside the stream.
		ressources = append(ressources, &ressourcepb.Ressource{Path: res["path"].(string), Name: res["name"].(string), Modified: modified, Size: size})
		if len(ressources) == max-1 {
			// I will send the stream at each 100 logs...
			rsp := &ressourcepb.GetRessourcesRsp{
				Ressources: ressources,
			}
			// Send the infos
			err = stream.Send(rsp)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			ressources = make([]*ressourcepb.Ressource, 0)
		}
	}

	// Send the last infos...
	if len(ressources) > 0 {

		rsp := &ressourcepb.GetRessourcesRsp{
			Ressources: ressources,
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

//* Remove a ressource from a client (custom service) to globular
func (self *Globule) RemoveRessource(ctx context.Context, rqst *ressourcepb.RemoveRessourceRqst) (*ressourcepb.RemoveRessourceRsp, error) {

	// Because regex dosent work properly I retreive all the ressources.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// get all ressource with that path.
	ressources, err := self.getRessources(rqst.Ressource.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	toDelete := make([]*ressourcepb.Ressource, 0)
	// Remove ressource that match...
	for i := 0; i < len(ressources); i++ {
		res := ressources[i]

		// In case the ressource is a sub-ressource I will remove it...
		if len(rqst.Ressource.Name) > 0 {
			if rqst.Ressource.Name == res.GetName() {
				toDelete = append(toDelete, res) // mark to be delete.
			}
		} else {
			toDelete = append(toDelete, res) // mark to be delete
		}

	}

	// Now I will delete the ressourcepb.
	for i := 0; i < len(toDelete); i++ {
		id := Utility.GenerateUUID(toDelete[i].Path + toDelete[i].Name)
		err := p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "Ressources", `{"_id":"`+id+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		// Delete the permissions ascosiated permission
		self.deletePermissions(toDelete[i].Path+"/"+toDelete[i].Name, "")
		err = p.Delete(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+toDelete[i].Path+"/"+toDelete[i].Name+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	// In that case the
	if len(rqst.Ressource.Name) == 0 {
		self.deletePermissions(rqst.Ressource.Path, "")
		err = p.Delete(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+rqst.Ressource.Path+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressourcepb.RemoveRessourceRsp{
		Result: true,
	}, nil
}

func (self *Globule) setActionPermission(action string, permission int32) error {

	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	actionPermission := make(map[string]interface{}, 0)
	actionPermission["action"] = action
	actionPermission["permission"] = permission
	actionPermission["_id"] = Utility.GenerateUUID(action)

	actionPermissionStr, _ := Utility.ToJson(actionPermission)
	err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "ActionPermission", `{"_id":"`+actionPermission["_id"].(string)+`"}`, actionPermissionStr, `[{"upsert":true}]`)
	if err != nil {
		return err
	}

	return nil
}

//* Set a ressource from a client (custom service) to globular
func (self *Globule) SetActionPermission(ctx context.Context, rqst *ressourcepb.SetActionPermissionRqst) (*ressourcepb.SetActionPermissionRsp, error) {

	err := self.setActionPermission(rqst.Action, rqst.Permission)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.SetActionPermissionRsp{
		Result: true,
	}, nil
}

//* Remove a ressource from a client (custom service) to globular
func (self *Globule) RemoveActionPermission(ctx context.Context, rqst *ressourcepb.RemoveActionPermissionRqst) (*ressourcepb.RemoveActionPermissionRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Try to delete the account...
	err = p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "ActionPermission", `{"_id":"`+Utility.GenerateUUID(rqst.Action)+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.RemoveActionPermissionRsp{
		Result: true,
	}, nil
}

//* Remove a ressource from a client (custom service) to globular
func (self *Globule) GetActionPermission(ctx context.Context, rqst *ressourcepb.GetActionPermissionRqst) (*ressourcepb.GetActionPermissionRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Try to delete the account...
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "ActionPermission", `{"_id":"`+Utility.GenerateUUID(rqst.Action)+`"}`, "")
	if err != nil {
		return &ressourcepb.GetActionPermissionRsp{
			Permission: int32(-1),
		}, nil
	}

	actionPermission := values.(map[string]interface{})

	return &ressourcepb.GetActionPermissionRsp{
		Permission: actionPermission["permission"].(int32),
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

	count, _ := p.Count(context.Background(), "local_ressource", "local_ressource", "Permissions", query, "")

	if count == 0 {
		_, err = p.InsertOne(context.Background(), "local_ressource", "local_ressource", "Permissions", map[string]interface{}{"owner": owner, "path": path, "permission": permission}, "")

	} else {
		err = p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "Permissions", query, jsonStr, "")
	}

	return err
}

// Set directory permission
func (self *Globule) setDirPermission(owner string, path string, permission int32) error {

	err := self.savePermission(owner, path, permission)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(self.GetAbsolutePath(path))
	if err != nil {
		return err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.IsDir() {
			err := self.setDirPermission(owner, path+"/"+file.Name(), permission)
			if err != nil {
				return err
			}
		} else {
			err := self.setRessourcePermission(owner, path+"/"+file.Name(), permission)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// Set file permission.
func (self *Globule) setRessourcePermission(owner string, path string, permission int32) error {
	return self.savePermission(owner, path, permission)
}

//* Set a file permission, create new one if not already exist. *
func (self *Globule) SetPermission(ctx context.Context, rqst *ressourcepb.SetPermissionRqst) (*ressourcepb.SetPermissionRsp, error) {
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

	// Now if the permission exist I will read the file info.
	fileInfo, _ := os.Stat(self.GetAbsolutePath(path))

	// Now I will test if the user or the role exist.
	var owner map[string]interface{}

	switch v := rqst.Permission.GetOwner().(type) {
	case *ressourcepb.RessourcePermission_User:
		// In that case I will try to find a user with that id
		values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+v.User+`"}`, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		owner = values.(map[string]interface{})
	case *ressourcepb.RessourcePermission_Role:
		// In that case I will try to find a role with that id
		values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"`+v.Role+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		owner = values.(map[string]interface{})
	case *ressourcepb.RessourcePermission_Application:
		// In that case I will try to find a role with that id
		values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+v.Application+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		owner = values.(map[string]interface{})
	}

	if fileInfo != nil {

		switch mode := fileInfo.Mode(); {
		case mode.IsDir():
			// do directory stuff
			err := self.setDirPermission(owner["_id"].(string), path, rqst.Permission.Number)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		case mode.IsRegular():
			// do file stuff
			err := self.setRessourcePermission(owner["_id"].(string), path, rqst.Permission.Number)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		}
	}

	ressources, err := self.getRessources(path)
	if err == nil {
		for i := 0; i < len(ressources); i++ {
			if ressources[i].GetPath() != path {
				path_ := ressources[i].GetPath()[len(path)+1:]
				paths := strings.Split(path_, "/")
				path_ = path
				// set sub-path...
				for j := 0; j < len(paths); j++ {
					path_ += "/" + paths[j]
					err := self.setRessourcePermission(owner["_id"].(string), path_, rqst.Permission.Number)
					if err != nil {
						return nil, status.Errorf(
							codes.Internal,
							Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
					}
				}
			}

			// create ressource permission
			err := self.setRessourcePermission(owner["_id"].(string), ressources[i].GetPath()+"/"+ressources[i].GetName(), rqst.Permission.Number)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}

		}
		// save ressource path.
		err = self.setRessourcePermission(owner["_id"].(string), path, rqst.Permission.Number)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressourcepb.SetPermissionRsp{
		Result: true,
	}, nil
}

func (self *Globule) setPermissionOwner(owner string, permission *ressourcepb.RessourcePermission) error {
	p, err := self.getPersistenceStore()
	if err != nil {
		return err
	}

	// Here I will try to find the owner in the user table
	_, err = p.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"_id":"`+owner+`"}`, ``)
	if err == nil {
		permission.Owner = &ressourcepb.RessourcePermission_User{
			User: owner,
		}
		return nil
	}

	// In the role.
	// In that case I will try to find a role with that id
	_, err = p.FindOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"`+owner+`"}`, "")
	if err == nil {
		permission.Owner = &ressourcepb.RessourcePermission_Role{
			Role: owner,
		}
		return nil
	}

	_, err = p.FindOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"_id":"`+owner+`"}`, "")
	if err == nil {
		permission.Owner = &ressourcepb.RessourcePermission_Application{
			Application: owner,
		}
		return nil
	}

	return errors.New("No Role or User found with id " + owner)
}

/**
 *
 */
func (self *Globule) getDirPermissions(path string) ([]*ressourcepb.RessourcePermission, error) {
	if !Utility.Exists(self.GetAbsolutePath(path)) {
		return nil, errors.New("No directory found with path " + self.GetAbsolutePath(path))
	}

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// So here I will get the list of retreived permission.
	permissions_, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Now I will read the permission values.
	permissions := make([]*ressourcepb.RessourcePermission, 0)
	for i := 0; i < len(permissions_); i++ {
		permission_ := permissions_[i].(map[string]interface{})
		permission := &ressourcepb.RessourcePermission{Path: permission_["path"].(string), Owner: nil, Number: permission_["permission"].(int32)}
		err = self.setPermissionOwner(permission_["owner"].(string), permission)
		if err != nil {
			return nil, err
		}

		// append into the permissions.
		permissions = append(permissions, permission)
	}

	// No the recursion.
	files, err := ioutil.ReadDir(self.GetAbsolutePath(path))
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.IsDir() {
			permissions_, err := self.getDirPermissions(path + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			// append to the permissions
			permissions = append(permissions, permissions_...)
		} else {
			permissions_, err := self.getRessourcePermissions(path + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			// append to the permissions
			permissions = append(permissions, permissions_...)
		}
	}

	return permissions, nil
}

func (self *Globule) getRessourcePermissions(path string) ([]*ressourcepb.RessourcePermission, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// So here I will get the list of retreived permission.
	permissions_, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, err
	}

	permissions := make([]*ressourcepb.RessourcePermission, 0)

	for i := 0; i < len(permissions_); i++ {
		permission_ := permissions_[i].(map[string]interface{})
		permission := &ressourcepb.RessourcePermission{Path: permission_["path"].(string), Owner: nil, Number: permission_["permission"].(int32)}
		err = self.setPermissionOwner(permission_["owner"].(string), permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (self *Globule) getPermissions(path string) ([]*ressourcepb.RessourcePermission, error) {

	fileInfo, err := os.Stat(self.GetAbsolutePath(path))
	if err == nil {

		switch mode := fileInfo.Mode(); {
		case mode.IsDir():
			// do directory stuff
			permissions, err := self.getDirPermissions(path)
			if err != nil {
				return nil, err
			}

			return permissions, nil

		case mode.IsRegular():
			// do file stuff
			permissions, err := self.getRessourcePermissions(path)
			if err != nil {
				return nil, err
			}

			return permissions, nil
		}
	} else {
		// do file stuff
		permissions, err := self.getRessourcePermissions(path)
		if err != nil {
			return nil, err
		}

		return permissions, nil
	}

	return nil, nil
}

//* Get All permissions for a given file/dir *
func (self *Globule) GetPermissions(ctx context.Context, rqst *ressourcepb.GetPermissionsRqst) (*ressourcepb.GetPermissionsRsp, error) {

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

	return &ressourcepb.GetPermissionsRsp{
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
			case *ressourcepb.RessourcePermission_User:
				if v.User == owner {
					err := p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+permission.GetPath()+`","owner":"`+v.User+`"}`, "")
					if err != nil {
						return err
					}
				}
			case *ressourcepb.RessourcePermission_Role:
				if v.Role == owner {
					err := p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+permission.GetPath()+`","owner":"`+v.Role+`"}`, "")
					if err != nil {
						return err
					}
				}

			case *ressourcepb.RessourcePermission_Application:
				if v.Application == owner {
					err := p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+permission.GetPath()+`","owner":"`+v.Application+`"}`, "")
					if err != nil {
						return err
					}
				}
			}
		} else {
			err := p.DeleteOne(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+permission.GetPath()+`"}`, "")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//* Delete a file permission *
func (self *Globule) DeletePermissions(ctx context.Context, rqst *ressourcepb.DeletePermissionsRqst) (*ressourcepb.DeletePermissionsRsp, error) {

	// That service made user of persistence service.
	err := self.deletePermissions(rqst.GetPath(), rqst.GetOwner())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeletePermissionsRsp{
		Result: true,
	}, nil
}

//* Create Permission for a dir (recursive) *
func (self *Globule) CreateDirPermissions(ctx context.Context, rqst *ressourcepb.CreateDirPermissionsRqst) (*ressourcepb.CreateDirPermissionsRsp, error) {

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

	permissions, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, "")
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

		p.InsertOne(context.Background(), "local_ressource", "local_ressource", "Permissions", permission_, "")
	}

	ressourceOwners, err := p.Find(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Now I will create the new permission of the created directory.
	for i := 0; i < len(ressourceOwners); i++ {
		// Copye the permission.
		ressourceOwner := ressourceOwners[i].(map[string]interface{})
		ressourceOwner_ := make(map[string]interface{}, 0)
		ressourceOwner_["owner"] = ressourceOwner["owner"]
		ressourceOwner_["path"] = path + "/" + rqst.GetName()

		p.InsertOne(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", ressourceOwner_, "")
	}

	// The user who create a directory will be the owner of the
	// directory.
	if clientId != "sa" && clientId != "guest" && len(rqst.GetName()) > 0 {
		ressourceOwner := make(map[string]interface{}, 0)
		ressourceOwner["owner"] = clientId
		ressourceOwner["path"] = path + "/" + rqst.GetName()
		ressourceOwnerStr, _ := Utility.ToJson(ressourceOwner)
		p.ReplaceOne(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", ressourceOwnerStr, ressourceOwnerStr, `[{"upsert":true}]`)
	}

	return &ressourcepb.CreateDirPermissionsRsp{
		Result: true,
	}, nil
}

//* Rename file/dir permission *
func (self *Globule) RenameFilePermission(ctx context.Context, rqst *ressourcepb.RenameFilePermissionRqst) (*ressourcepb.RenameFilePermissionRsp, error) {

	path := rqst.GetPath()
	path = strings.ReplaceAll(path, "\\", "/")
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-2]
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
	permissions, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Permissions", `{}`, "")
	if err == nil {
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), oldPath) {
				path := newPath + permission["path"].(string)[len(oldPath):]
				err := p.Update(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+permission["path"].(string)+`"}`, `{"$set":{"path":"`+path+`"}}`, "")
				if err != nil {
					return nil, status.Errorf(
						codes.Internal,
						Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
				}
			}
		}
	}

	// Replace file owner path... regex not work... "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"}
	permissions, err = p.Find(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{}`, "")
	if err == nil {
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), oldPath) {
				path := newPath + permission["path"].(string)[len(oldPath):]
				err = p.Update(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+permission["path"].(string)+`"}`, `{"$set":{"path":"`+path+`"}}`, "")
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
	return &ressourcepb.RenameFilePermissionRsp{
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
	permissions, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Permissions", `{}`, "")
	if err == nil {
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), path) {
				err := p.Delete(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+permission["path"].(string)+`"}`, "")
				if err != nil {
					return err
				}
			}
		}
	}

	// Replace file owner path... regex not work... "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"}
	permissions, err = p.Find(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{}`, "")
	if err == nil {
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), path) {
				err = p.Delete(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+permission["path"].(string)+`"}`, "")
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//* Delete Permission for a dir (recursive) *
func (self *Globule) DeleteDirPermissions(ctx context.Context, rqst *ressourcepb.DeleteDirPermissionsRqst) (*ressourcepb.DeleteDirPermissionsRsp, error) {
	err := self.deleteDirPermissions(rqst.GetPath())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeleteDirPermissionsRsp{
		Result: true,
	}, nil
}

//* Delete a single file permission *
func (self *Globule) DeleteFilePermissions(ctx context.Context, rqst *ressourcepb.DeleteFilePermissionsRqst) (*ressourcepb.DeleteFilePermissionsRsp, error) {
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	path := rqst.GetPath()

	err = p.Delete(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = p.Delete(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.DeleteFilePermissionsRsp{
		Result: true,
	}, nil
}

//* Validate a token *
func (self *Globule) ValidateToken(ctx context.Context, rqst *ressourcepb.ValidateTokenRqst) (*ressourcepb.ValidateTokenRsp, error) {
	clientId, _, expireAt, err := Interceptors.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &ressourcepb.ValidateTokenRsp{
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

	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Applications", `{"path":"/`+name+`"}`, ``)
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
			log.Println("Application", name, "can run "+method)
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
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"name":"`+userName+`"}`, `[{"Projection":{"roles":1}}]`)
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

	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Roles", `{"_id":"`+roleName+`"}`, `[{"Projection":{"actions":1}}]`)
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

	// If the user is the owner of the ressource it has the permission
	count, err := client.Count(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`","owner":"`+name+`"}`, ``)
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

	// If the user is the owner of the ressource it has all permission
	if self.isOwner(name, path) {
		return true, 0
	}

	count, err := p.Count(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, ``)
	if err != nil {
		return false, 0
	}

	if count == 0 {
		count, err = p.Count(context.Background(), "local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, ``)
		if err != nil {
			return false, 0
		}
	}

	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"owner":"`+name+`", "path":"`+path+`"}`, ``)
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
func (self *Globule) validateUserRessourceAccess(userName string, method string, path string, permission int32) error {

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
	values, err := p.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", `{"name":"`+userName+`"}`, `[{"Projection":{"roles":1}}]`)
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
	// have write to execute the action on the ressourcepb.
	if count > 0 {
		return errors.New("Permission Denied for " + userName)
	}

	count, err = p.Count(context.Background(), "local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, ``)
	if err != nil {
		if count > 0 {
			return errors.New("Permission Denied for " + userName)
		}
	}

	return nil
}

//* Validate if user can access a given file. *
func (self *Globule) ValidateUserRessourceAccess(ctx context.Context, rqst *ressourcepb.ValidateUserRessourceAccessRqst) (*ressourcepb.ValidateUserRessourceAccessRsp, error) {

	path := rqst.GetPath() // The path of the ressourcepb.

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

	err = self.validateUserRessourceAccess(clientId, rqst.Method, path, rqst.Permission)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.ValidateUserRessourceAccessRsp{
		Result: true,
	}, nil
}

//* Validate if application can access a given file. *
func (self *Globule) ValidateApplicationRessourceAccess(ctx context.Context, rqst *ressourcepb.ValidateApplicationRessourceAccessRqst) (*ressourcepb.ValidateApplicationRessourceAccessRsp, error) {

	path := rqst.GetPath()

	hasApplicationPermission, count := self.hasPermission(rqst.Name, path, rqst.Permission)
	if hasApplicationPermission {
		return &ressourcepb.ValidateApplicationRessourceAccessRsp{
			Result: true,
		}, nil
	}

	// If theres permission definied and we are here it's means the user dosent
	// have write to execute the action on the ressourcepb.
	if count > 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Permission Denied for "+rqst.Name)))

	}

	return &ressourcepb.ValidateApplicationRessourceAccessRsp{
		Result: true,
	}, nil
}

//* Validate if user can access a given method. *
func (self *Globule) ValidateUserAccess(ctx context.Context, rqst *ressourcepb.ValidateUserAccessRqst) (*ressourcepb.ValidateUserAccessRsp, error) {

	// first of all I will validate the token.
	clientID, _, expiredAt, err := Interceptors.ValidateToken(rqst.Token)
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

	// Here I will test if the user can run that function or not...
	err = self.validateUserAccess(clientID, rqst.Method)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.ValidateUserAccessRsp{
		Result: true,
	}, nil
}

//* Validate if application can access a given method. *
func (self *Globule) ValidateApplicationAccess(ctx context.Context, rqst *ressourcepb.ValidateApplicationAccessRqst) (*ressourcepb.ValidateApplicationAccessRsp, error) {
	log.Println("-------> validate application access ", rqst.GetName(), rqst.Method)
	err := self.validateApplicationAccess(rqst.GetName(), rqst.Method)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressourcepb.ValidateApplicationAccessRsp{
		Result: true,
	}, nil
}

//* Retrun a json string with all file info *
func (self *Globule) GetAllFilesInfo(ctx context.Context, rqst *ressourcepb.GetAllFilesInfoRqst) (*ressourcepb.GetAllFilesInfoRsp, error) {
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
	return &ressourcepb.GetAllFilesInfoRsp{Result: string(jsonStr)}, nil
}

func (self *Globule) GetAllApplicationsInfo(ctx context.Context, rqst *ressourcepb.GetAllApplicationsInfoRqst) (*ressourcepb.GetAllApplicationsInfoRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceStore()
	if err != nil {
		return nil, err
	}

	// So here I will get the list of retreived permission.
	values, err := p.Find(context.Background(), "local_ressource", "local_ressource", "Applications", `{}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	jsonStr, _ := Utility.ToJson(values)

	return &ressourcepb.GetAllApplicationsInfoRsp{
		Result: jsonStr,
	}, nil

}
