package ressource

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// admin Client Service
////////////////////////////////////////////////////////////////////////////////

type Ressource_Client struct {
	cc *grpc.ClientConn
	c  RessourceServiceClient

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
func NewRessource_Client(address string, name string) (*Ressource_Client, error) {
	client := new(Ressource_Client)
	err := api.InitClient(client, address, name)
	if err != nil {
		return nil, err
	}

	client.cc, err = api.GetClientConnection(client)
	if err != nil {
		return nil, err
	}

	client.c = NewRessourceServiceClient(client.cc)

	return client, nil
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
	rqst := &RegisterAccountRqst{
		Account: &Account{
			Name:     name,
			Email:    email,
			Password: "",
		},
		Password:        password,
		ConfirmPassword: confirmation_password,
	}

	_, err := self.c.RegisterAccount(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

// Delete an account.
func (self *Ressource_Client) DeleteAccount(id string) error {
	rqst := &DeleteAccountRqst{
		Id: id,
	}

	_, err := self.c.DeleteAccount(api.GetClientContext(self), rqst)
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

	rqst := &AuthenticateRqst{
		Name:     name,
		Password: password,
	}

	rsp, err := self.c.Authenticate(api.GetClientContext(self), rqst)
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
 * Create a new role with given action list.
 */
func (self *Ressource_Client) CreateRole(name string, actions []string) error {
	rqst := new(CreateRoleRqst)
	role := new(Role)
	role.Name = name
	role.Actions = actions
	rqst.Role = role
	_, err := self.c.CreateRole(api.GetClientContext(self), rqst)

	return err
}

func (self *Ressource_Client) DeleteRole(name string) error {
	rqst := new(DeleteRoleRqst)
	rqst.RoleId = name

	_, err := self.c.DeleteRole(api.GetClientContext(self), rqst)

	return err
}

/**
 * Add a action to a given role.
 */
func (self *Ressource_Client) AddRoleAction(roleId string, action string) error {
	rqst := &AddRoleActionRqst{
		RoleId: roleId,
		Action: action,
	}
	_, err := self.c.AddRoleAction(api.GetClientContext(self), rqst)

	return err
}

/**
 * Remove action from a given role.
 */
func (self *Ressource_Client) RemoveRoleAction(roleId string, action string) error {
	rqst := &RemoveRoleActionRqst{
		RoleId: roleId,
		Action: action,
	}
	_, err := self.c.RemoveRoleAction(api.GetClientContext(self), rqst)

	return err
}

/**
 * Add a action to a given application.
 */
func (self *Ressource_Client) AddApplicationAction(applicationId string, action string) error {
	rqst := &AddApplicationActionRqst{
		ApplicationId: applicationId,
		Action:        action,
	}
	_, err := self.c.AddApplicationAction(api.GetClientContext(self), rqst)

	return err
}

/**
 * Remove action from a given application.
 */
func (self *Ressource_Client) RemoveApplicationAction(applicationId string, action string) error {
	rqst := &RemoveApplicationActionRqst{
		ApplicationId: applicationId,
		Action:        action,
	}
	_, err := self.c.RemoveApplicationAction(api.GetClientContext(self), rqst)

	return err
}

/**
 * Set role to a account
 */
func (self *Ressource_Client) AddAccountRole(accountId string, roleId string) error {
	rqst := &AddAccountRoleRqst{
		AccountId: accountId,
		RoleId:    roleId,
	}
	_, err := self.c.AddAccountRole(api.GetClientContext(self), rqst)

	return err
}

/**
 * Remove role from an account
 */
func (self *Ressource_Client) RemoveAccountRole(accountId string, roleId string) error {
	rqst := &RemoveAccountRoleRqst{
		AccountId: accountId,
		RoleId:    roleId,
	}
	_, err := self.c.RemoveAccountRole(api.GetClientContext(self), rqst)

	return err
}

/**
 * Return the list of all available actions on the server.
 */
func (self *Ressource_Client) GetAllActions() ([]string, error) {
	rqst := &GetAllActionsRqst{}
	rsp, err := self.c.GetAllActions(api.GetClientContext(self), rqst)
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
	rqst := &SetPermissionRqst{
		Permission: &RessourcePermission{
			Owner: &RessourcePermission_User{
				User: userId,
			},
			Path:   path,
			Number: permission,
		},
	}

	_, err := self.c.SetPermission(api.GetClientContext(self), rqst)
	return err
}

/**
 * Set file permission for a given role.
 */
func (self *Ressource_Client) SetRessourcePermissionByRole(roleId string, path string, permission int32) error {
	rqst := &SetPermissionRqst{
		Permission: &RessourcePermission{
			Owner: &RessourcePermission_Role{
				Role: roleId,
			},
			Path:   path,
			Number: permission,
		},
	}

	_, err := self.c.SetPermission(api.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) GetRessourcePermissions(path string) (string, error) {
	rqst := &GetPermissionsRqst{
		Path: path,
	}

	rsp, err := self.c.GetPermissions(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}
	return rsp.GetPermissions(), nil
}

func (self *Ressource_Client) DeleteRessourcePermissions(path string, owner string) error {
	rqst := &DeletePermissionsRqst{
		Path:  path,
		Owner: owner,
	}

	_, err := self.c.DeletePermissions(api.GetClientContext(self), rqst)

	return err
}

func (self *Ressource_Client) GetAllFilesInfo() (string, error) {
	rqst := &GetAllFilesInfoRqst{}

	rsp, err := self.c.GetAllFilesInfo(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidateUserRessourceAccess(token string, path string, method string, permission int32) (bool, error) {
	rqst := &ValidateUserRessourceAccessRqst{}
	rqst.Token = token
	rqst.Path = path
	rqst.Method = method
	rqst.Permission = permission

	rsp, err := self.c.ValidateUserRessourceAccess(api.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidateApplicationRessourceAccess(application string, path string, method string, permission int32) (bool, error) {
	rqst := &ValidateApplicationRessourceAccessRqst{}
	rqst.Name = application
	rqst.Path = path
	rqst.Method = method
	rqst.Permission = permission

	rsp, err := self.c.ValidateApplicationRessourceAccess(api.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidatePeerRessourceAccess(name string, path string, method string, permission int32) (bool, error) {
	rqst := &ValidatePeerRessourceAccessRqst{}
	rqst.Name = name
	rqst.Path = path
	rqst.Method = method
	rqst.Permission = permission

	rsp, err := self.c.ValidatePeerRessourceAccess(api.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidateUserAccess(token string, method string) (bool, error) {
	rqst := &ValidateUserAccessRqst{}
	rqst.Token = token
	rqst.Method = method

	rsp, err := self.c.ValidateUserAccess(api.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidateApplicationAccess(application string, method string) (bool, error) {
	rqst := &ValidateApplicationAccessRqst{}
	rqst.Name = application
	rqst.Method = method
	rsp, err := self.c.ValidateApplicationAccess(api.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidatePeerAccess(name string, method string) (bool, error) {
	rqst := &ValidatePeerAccessRqst{}
	rqst.Name = name
	rqst.Method = method
	rsp, err := self.c.ValidatePeerAccess(api.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) CreateDirPermissions(token string, path string, name string) error {
	rqst := &CreateDirPermissionsRqst{
		Token: token,
		Path:  path,
		Name:  name,
	}
	_, err := self.c.CreateDirPermissions(api.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) RenameFilePermission(path string, oldName string, newName string) error {
	rqst := &RenameFilePermissionRqst{
		Path:    path,
		OldName: oldName,
		NewName: newName,
	}

	_, err := self.c.RenameFilePermission(api.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) DeleteDirPermissions(path string) error {
	rqst := &DeleteDirPermissionsRqst{
		Path: path,
	}
	_, err := self.c.DeleteDirPermissions(api.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) DeleteFilePermissions(path string) error {
	rqst := &DeleteFilePermissionsRqst{
		Path: path,
	}
	_, err := self.c.DeleteFilePermissions(api.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) DeleteRolePermissions(id string) error {
	rqst := &DeleteRolePermissionsRqst{
		Id: id,
	}
	_, err := self.c.DeleteRolePermissions(api.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) DeleteAccountPermissions(id string) error {
	rqst := &DeleteAccountPermissionsRqst{
		Id: id,
	}
	_, err := self.c.DeleteAccountPermissions(api.GetClientContext(self), rqst)
	return err
}

func (self *Ressource_Client) GetActionPermission(action string) (int32, error) {
	rqst := &GetActionPermissionRqst{
		Action: action,
	}

	rsp, err := self.c.GetActionPermission(api.GetClientContext(self), rqst)
	if err != nil {
		return -1, err
	}

	return rsp.Permission, nil
}

func (self *Ressource_Client) SetRessource(name string, path string, modified int64, size int64, token string) error {
	ressource := &Ressource{
		Name:     name,
		Path:     path,
		Modified: modified,
		Size:     size,
	}

	rqst := &SetRessourceRqst{
		Ressource: ressource,
	}
	var err error
	if len(token) > 0 {
		md := metadata.New(map[string]string{"token": string(token), "domain": self.GetDomain(), "mac": Utility.MyMacAddr(), "ip": Utility.MyIP()})
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		_, err = self.c.SetRessource(ctx, rqst)
	} else {
		_, err = self.c.SetRessource(api.GetClientContext(self), rqst)
	}

	return err
}

func (self *Ressource_Client) SetRessourceOwner(owner string, path string, token string) error {
	rqst := &SetRessourceOwnerRqst{
		Owner: owner,
		Path:  path,
	}
	var err error
	if len(token) > 0 {
		md := metadata.New(map[string]string{"token": string(token), "domain": self.GetDomain(), "mac": Utility.MyMacAddr(), "ip": Utility.MyIP()})
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		_, err = self.c.SetRessourceOwner(ctx, rqst)
	} else {
		_, err = self.c.SetRessourceOwner(api.GetClientContext(self), rqst)
	}

	return err
}

// Set action permission
func (self *Ressource_Client) SetActionPermission(action string, permission int32, token string) error {
	var err error
	// Set action permission.
	rqst := &SetActionPermissionRqst{
		Action:     action,
		Permission: permission,
	}
	if len(token) > 0 {
		md := metadata.New(map[string]string{"token": string(token), "domain": self.GetDomain(), "mac": Utility.MyMacAddr(), "ip": Utility.MyIP()})
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		// Set action permission.
		_, err = self.c.SetActionPermission(ctx, rqst)
	} else {
		_, err = self.c.SetActionPermission(api.GetClientContext(self), rqst)
	}

	return err
}

/////////////////////// Log ////////////////////////

// Append a new log information.
func (self *Ressource_Client) Log(application string, user string, method string, err_ error) error {

	// Here I set a log information.
	rqst := new(LogRqst)
	info := new(LogInfo)
	info.Application = application
	info.UserName = user
	info.Method = method
	info.Date = time.Now().Unix()
	if err_ != nil {
		info.Message = err_.Error()
		info.Type = LogType_ERROR
	} else {
		info.Type = LogType_INFO
	}
	rqst.Info = info

	_, err := self.c.Log(api.GetClientContext(self), rqst)

	return err
}

/////////////////////////////////// Peer's  ///////////////////////////////////

// Register a peer with a given name and mac address.
func (self *Ressource_Client) RegisterPeer(name string, mac string) error {
	rqst := &RegisterPeerRqst{
		Peer: &Peer{
			MacAddress: mac,
			Name:       name,
		},
	}

	_, err := self.c.RegisterPeer(api.GetClientContext(self), rqst)
	return err

}
