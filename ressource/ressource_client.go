package ressource

import (
	"strconv"

	"io/ioutil"
	"log"
	"os"
	"time"

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
func NewRessource_Client(address string, name string) *Ressource_Client {

	client := new(Ressource_Client)
	api.InitClient(client, address, name)
	client.cc = api.GetClientConnection(client)
	client.c = NewRessourceServiceClient(client.cc)

	return client
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
	log.Println("---> path ", path)
	log.Println("---> is Local ", self.GetDomain(), Utility.IsLocal(self.GetDomain()))
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

/////////////////////////////// File permissions ///////////////////////////////

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

func (self *Ressource_Client) ValidateUserFileAccess(token string, path string, method string, permission string) (bool, error) {
	rqst := &ValidateUserFileAccessRqst{}
	rqst.Token = token
	rqst.Path = path
	rqst.Method = method
	rqst.Permission = permission

	rsp, err := self.c.ValidateUserFileAccess(api.GetClientContext(self), rqst)
	if err != nil {
		return false, err
	}

	return rsp.GetResult(), nil
}

func (self *Ressource_Client) ValidateApplicationFileAccess(application string, path string, method string, permission string) (bool, error) {
	rqst := &ValidateApplicationFileAccessRqst{}
	rqst.Name = application
	rqst.Path = path
	rqst.Method = method
	rqst.Permission = permission

	rsp, err := self.c.ValidateApplicationFileAccess(api.GetClientContext(self), rqst)
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

/////////////////////// Log ////////////////////////

// Append a new log information.
func (self *Ressource_Client) Log(application string, userId string, method string, err_ error) error {

	// Here I set a log information.
	rqst := new(LogRqst)
	info := new(LogInfo)
	info.Application = application
	info.UserId = userId
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

// Return all log method.
func (self *Ressource_Client) GetLogMethods() ([]string, error) {
	rsp, err := self.c.GetLogMethods(api.GetClientContext(self), &GetLogMethodsRqst{})
	if err != nil {
		return nil, err
	}

	return rsp.Methods, nil
}
