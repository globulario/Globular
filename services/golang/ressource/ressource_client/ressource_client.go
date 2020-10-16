package ressource_client

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	globular "github.com/davecourtois/Globular/services/golang/globular_client"
	"github.com/davecourtois/Globular/services/golang/ressource/ressourcepb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

////////////////////////////////////////////////////////////////////////////////
// admin Client Service
////////////////////////////////////////////////////////////////////////////////

type Ressource_Client struct {
	cc *grpc.ClientConn
	c  ressourcepb.RessourceServiceClient

	// The id of the service
	id string

	// The name of the service
	name string

	// The client domain
	domain string

	// The port number
	port int

	// is the connection is secure?
	hasTLS bool

	// Link to client key file
	keyFile string

	// Link to client certificate file.
	certFile string

	// certificate authority file
	caFile string
}

// Create a connection to the service.
func NewRessourceService_Client(address string, id string) (*Ressource_Client, error) {
	client := new(Ressource_Client)
	err := globular.InitClient(client, address, id)
	if err != nil {
		return nil, err
	}

	client.cc, err = globular.GetClientConnection(client)
	if err != nil {
		return nil, err
	}

	client.c = ressourcepb.NewRessourceServiceClient(client.cc)

	return client, nil
}

func (self *Ressource_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the ipv4 address
// Return the address
func (self *Ressource_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the domain
func (self *Ressource_Client) GetDomain() string {
	return self.domain
}

// Return the id of the service instance
func (self *Ressource_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *Ressource_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Ressource_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *Ressource_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *Ressource_Client) SetId(id string) {
	self.id = id
}

// Set the client name.
func (self *Ressource_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Ressource_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Ressource_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Ressource_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Ressource_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Ressource_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Ressource_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Ressource_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Ressource_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Ressource_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

////////////// API ////////////////
// Register a new Account.
func (self *Ressource_Client) RegisterAccount(name string, email string, password string, confirmation_password string) error {
	rqst := &ressourcepb.RegisterAccountRqst{
		Account: &ressourcepb.Account{
			Name:     name,
			Email:    email,
			Password: "",
		},
		Password:        password,
		ConfirmPassword: confirmation_password,
	}

	_, err := self.c.RegisterAccount(globular.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

// Delete an account.
func (self *Ressource_Client) DeleteAccount(id string) error {
	rqst := &ressourcepb.DeleteAccountRqst{
		Id: id,
	}

	_, err := self.c.DeleteAccount(globular.GetClientContext(self), rqst)
	return err
}

// Authenticate a user.
func (self *Ressource_Client) Authenticate(name string, password string) (string, error) {
	// In case of other domain than localhost I will rip off the token file
	// before each authentication.
	path := os.TempDir() + string(os.PathSeparator) + self.GetDomain() + "_token"
	if !Utility.IsLocal(self.GetDomain()) {
		// remove the file if it already exist.
		os.Remove(path)
	}

	rqst := &ressourcepb.AuthenticateRqst{
		Name:     name,
		Password: password,
	}

	rsp, err := self.c.Authenticate(globular.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	// Here I will save the token into the temporary directory the token will be valid for a given time (default is 15 minutes)
	// it's the responsability of the client to keep it refresh... see Refresh token from the server...
	if !Utility.IsLocal(self.GetDomain()) {
		err = ioutil.WriteFile(path, []byte(rsp.Token), 0644)
		if err != nil {
			return "", err
		}
	}

	return rsp.Token, nil
}

/**
 *  Generate a new token from expired one.
 */
func (self *Ressource_Client) RefreshToken(token string) (string, error) {
	rqst := new(ressourcepb.RefreshTokenRqst)
	rqst.Token = token

	rsp, err := self.c.RefreshToken(globular.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	return rsp.Token, nil
}

/**
 * Create a new role with given action list.
 */
func (self *Ressource_Client) CreateRole(name string, actions []string) error {
	rqst := new(ressourcepb.CreateRoleRqst)
	role := new(ressourcepb.Role)
	role.Name = name
	role.Actions = actions
	rqst.Role = role
	_, err := self.c.CreateRole(globular.GetClientContext(self), rqst)

	return err
}

func (self *Ressource_Client) DeleteRole(name string) error {
	rqst := new(ressourcepb.DeleteRoleRqst)
	rqst.RoleId = name

	_, err := self.c.DeleteRole(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Add a action to a given role.
 */
func (self *Ressource_Client) AddRoleAction(roleId string, action string) error {
	rqst := &ressourcepb.AddRoleActionRqst{
		RoleId: roleId,
		Action: action,
	}
	_, err := self.c.AddRoleAction(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Remove action from a given role.
 */
func (self *Ressource_Client) RemoveRoleAction(roleId string, action string) error {
	rqst := &ressourcepb.RemoveRoleActionRqst{
		RoleId: roleId,
		Action: action,
	}
	_, err := self.c.RemoveRoleAction(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Add a action to a given application.
 */
func (self *Ressource_Client) AddApplicationAction(applicationId string, action string) error {
	rqst := &ressourcepb.AddApplicationActionRqst{
		ApplicationId: applicationId,
		Action:        action,
	}
	_, err := self.c.AddApplicationAction(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Remove action from a given application.
 */
func (self *Ressource_Client) RemoveApplicationAction(applicationId string, action string) error {
	rqst := &ressourcepb.RemoveApplicationActionRqst{
		ApplicationId: applicationId,
		Action:        action,
	}
	_, err := self.c.RemoveApplicationAction(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Set role to a account
 */
func (self *Ressource_Client) AddAccountRole(accountId string, roleId string) error {
	rqst := &ressourcepb.AddAccountRoleRqst{
		AccountId: accountId,
		RoleId:    roleId,
	}
	_, err := self.c.AddAccountRole(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Remove role from an account
 */
func (self *Ressource_Client) RemoveAccountRole(accountId string, roleId string) error {
	rqst := &ressourcepb.RemoveAccountRoleRqst{
		AccountId: accountId,
		RoleId:    roleId,
	}
	_, err := self.c.RemoveAccountRole(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Return the list of all available actions on the server.
 */
func (self *Ressource_Client) GetAllActions() ([]string, error) {
	rqst := &ressourcepb.GetAllActionsRqst{}
	rsp, err := self.c.GetAllActions(globular.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}
	return rsp.Actions, err
}

/////////////////////////////// Ressouce permissions ///////////////////////////////

/**
 * Set file permission for a given user.
 */
func (self *Ressource_Client) SetRessourcePermissionByUser(userId string, path string, permission int32) error {
	rqst := &ressourcepb.SetPermissionRqst{
		Permission: &ressourcepb.RessourcePermission{
			Owner: &ressourcepb.RessourcePermission_User{
				User: userId,
			},
			Path:   path,
			Number: permission,
		},
	}

	_, err := self.c.SetPermission(globular.GetClientContext(self), rqst)
	return err
}

/**
 * Set file permission for a given role.
 */
func (self *Ressource_Client) SetRessourcePermissionByRole(roleId string, path string, permission int32) error {
	rqst := &ressourcepb.SetPermissionRqst{
		Permission: &ressourcepb.RessourcePermission{
			Owner: &ressourcepb.RessourcePermission_Role{
				Role: roleId,
			},
			Path:   path,
			Number: permission,
		},
	}

	_, err := self.c.SetPermission(globular.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) GetRessourcePermissions(path string) (string, error) {
	rqst := &ressourcepb.GetPermissionsRqst{
		Path: path,
	}

	rsp, err := self.c.GetPermissions(globular.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}
	return rsp.GetPermissions(), nil
}

func (self *Ressource_Client) DeleteRessourcePermissions(path string, owner string) error {
	rqst := &ressourcepb.DeletePermissionsRqst{
		Path:  path,
		Owner: owner,
	}

	_, err := self.c.DeletePermissions(globular.GetClientContext(self), rqst)

	return err
}

func (self *Ressource_Client) GetAllFilesInfo() (string, error) {
	rqst := &ressourcepb.GetAllFilesInfoRqst{}

	rsp, err := self.c.GetAllFilesInfo(globular.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidateUserRessourceAccess(token string, path string, method string, permission int32) (bool, error) {
	rqst := &ressourcepb.ValidateUserRessourceAccessRqst{}
	rqst.Token = token
	rqst.Path = path
	rqst.Method = method
	rqst.Permission = permission

	rsp, err := self.c.ValidateUserRessourceAccess(globular.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidateApplicationRessourceAccess(application string, path string, method string, permission int32) (bool, error) {
	rqst := &ressourcepb.ValidateApplicationRessourceAccessRqst{}
	rqst.Name = application
	rqst.Path = path
	rqst.Method = method
	rqst.Permission = permission

	rsp, err := self.c.ValidateApplicationRessourceAccess(globular.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidateUserAccess(token string, method string) (bool, error) {
	rqst := &ressourcepb.ValidateUserAccessRqst{}
	rqst.Token = token
	rqst.Method = method

	rsp, err := self.c.ValidateUserAccess(globular.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidateApplicationAccess(application string, method string) (bool, error) {
	rqst := &ressourcepb.ValidateApplicationAccessRqst{}
	rqst.Name = application
	rqst.Method = method
	rsp, err := self.c.ValidateApplicationAccess(globular.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) DeleteRolePermissions(id string) error {
	rqst := &ressourcepb.DeleteRolePermissionsRqst{
		Id: id,
	}
	_, err := self.c.DeleteRolePermissions(globular.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) DeleteAccountPermissions(id string) error {
	rqst := &ressourcepb.DeleteAccountPermissionsRqst{
		Id: id,
	}
	_, err := self.c.DeleteAccountPermissions(globular.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) GetActionPermission(action string) ([]*ressourcepb.ActionParameterRessourcePermission, error) {
	rqst := &ressourcepb.GetActionPermissionRqst{
		Action: action,
	}

	rsp, err := self.c.GetActionPermission(globular.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	return rsp.ActionParameterRessourcePermissions, nil
}

func (self *Ressource_Client) SetRessource(name string, path string, modified int64, size int64, token string) error {
	ressource := &ressourcepb.Ressource{
		Name:     name,
		Path:     path,
		Modified: modified,
		Size:     size,
	}

	rqst := &ressourcepb.SetRessourceRqst{
		Ressource: ressource,
	}
	var err error
	if len(token) > 0 {
		md := metadata.New(map[string]string{"token": string(token), "domain": self.GetDomain(), "mac": Utility.MyMacAddr(), "ip": Utility.MyIP()})
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		_, err = self.c.SetRessource(ctx, rqst)
	} else {
		_, err = self.c.SetRessource(globular.GetClientContext(self), rqst)
	}

	return err
}

func (self *Ressource_Client) SetRessourceOwner(owner string, path string, token string) error {
	rqst := &ressourcepb.SetRessourceOwnerRqst{
		Owner: owner,
		Path:  path,
	}
	var err error
	if len(token) > 0 {
		md := metadata.New(map[string]string{"token": string(token), "domain": self.GetDomain(), "mac": Utility.MyMacAddr(), "ip": Utility.MyIP()})
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		_, err = self.c.SetRessourceOwner(ctx, rqst)
	} else {
		_, err = self.c.SetRessourceOwner(globular.GetClientContext(self), rqst)
	}

	return err
}

// Set action permission
func (self *Ressource_Client) SetActionPermission(action string, actionParameterRessourcePermissions []*ressourcepb.ActionParameterRessourcePermission, token string) error {
	var err error

	// Set action permission.
	rqst := &ressourcepb.SetActionPermissionRqst{
		Action:                              action,
		ActionParameterRessourcePermissions: actionParameterRessourcePermissions,
	}

	if len(token) > 0 {
		md := metadata.New(map[string]string{"token": string(token), "domain": self.GetDomain(), "mac": Utility.MyMacAddr(), "ip": Utility.MyIP()})
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		// Set action permission.
		_, err = self.c.SetActionPermission(ctx, rqst)
	} else {
		_, err = self.c.SetActionPermission(globular.GetClientContext(self), rqst)
	}

	return err
}

/////////////////////// Log ////////////////////////

// Append a new log information.
func (self *Ressource_Client) Log(application string, user string, method string, err_ error) error {

	// Here I set a log information.
	rqst := new(ressourcepb.LogRqst)
	info := new(ressourcepb.LogInfo)
	info.Application = application
	info.UserName = user
	info.Method = method
	info.Date = time.Now().Unix()
	if err_ != nil {
		info.Message = err_.Error()
		info.Type = ressourcepb.LogType_ERROR_MESSAGE
	} else {
		info.Type = ressourcepb.LogType_INFO_MESSAGE
	}
	rqst.Info = info

	_, err := self.c.Log(globular.GetClientContext(self), rqst)

	return err
}
