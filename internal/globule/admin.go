package globule

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/globulario/services/golang/authentication/authentication_client"
	serviceManagerClient "github.com/globulario/services/golang/services_manager/services_manager_client"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/rbac/rbacpb"
	Utility "github.com/globulario/utility"
)

// addResourceOwner adds the given subject as an owner of the resource at path.
func (globule *Globule) addResourceOwner(path, resourceType, subject string, subjectType rbacpb.SubjectType) error {

	rbacClient, err := getRbacClient()
	if err != nil {
		return err
	}
	return rbacClient.AddResourceOwner(path, resourceType, subject, subjectType)
}

// getAuthenticationClient returns an authenticated client to the authentication service.
func (g *Globule) getAuthenticationClient(address string) (*authentication_client.Authentication_Client, error) {
	// Authenticate the user in order to get the token
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		g.log.Error("fail to create authentication client:", "err", err)
		return nil, err
	}
	return authenticationClient, nil
}

// getResourceClient returns a client to the resource service.
func (g *Globule) getServiceManagerClient(domain string) (*serviceManagerClient.Services_Manager_Client, error) {
	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	servicesManagerCli, err := serviceManagerClient.NewServicesManagerService_Client(domain, "services_manager.ServicesManagerService")
	if err != nil {
		g.log.Error("fail to connect to services manager at ", "domain", domain, "err", err)
		return nil, err
	}
	return servicesManagerCli, nil
}

// BootstrapAdmin authenticates "sa" from config, ensures admin role exists with all actions,
// ensures the "sa" account exists in Resource, and assigns the admin role to it.
func (g *Globule) BootstrapAdmin() error {
	address, err := config.GetAddress()
	if err != nil {
		return fmt.Errorf("get address: %w", err)
	}

	// 1) Clients
	resourceClient, err := getResourceClient(address)
	if err != nil {
		return fmt.Errorf("resource client: %w", err)
	}
	defer resourceClient.Close()

	authClient, err := g.getAuthenticationClient(address)
	if err != nil {
		return fmt.Errorf("auth client: %w", err)
	}
	defer authClient.Close()

	// 2) Read config for domain/email/pwd
	cfgBytes, err := os.ReadFile(config.GetConfigDir() + "/config.json")
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	cfg := map[string]any{}
	if err := json.Unmarshal(cfgBytes, &cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	domain, _ := cfg["Domain"].(string)
	adminEmail, _ := cfg["AdminEmail"].(string)
	rootPwd, _ := cfg["RootPassword"].(string)
	if domain == "" || adminEmail == "" || rootPwd == "" {
		return fmt.Errorf("missing Domain/AdminEmail/RootPassword in config")
	}

	// 3) Authenticate "sa" against Authentication service (config-backed path)
	token, err := authClient.Authenticate("sa", rootPwd)
	if err != nil {
		return fmt.Errorf("authenticate sa: %w", err)
	}

	// 4) Ensure "admin" role exists with all actions
	servicesManager, err := g.getServiceManagerClient(address)
	if err != nil {
		return fmt.Errorf("services manager client: %w", err)
	}
	actions, err := servicesManager.GetAllActions()
	if err != nil {
		return fmt.Errorf("get all actions: %w", err)
	}
	if err := resourceClient.CreateRole(token, "admin", "admin", actions); err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			return fmt.Errorf("create admin role: %w", err)
		}
	}

	// 5) Ensure "sa" account exists in Resource; if not, create it
	accts, err := resourceClient.GetAccounts(`{"_id":"sa"}`)
	if err != nil {
		return fmt.Errorf("get accounts: %w", err)
	}
	if len(accts) == 0 {
		// Create "sa" account in Resource. The Authentication service already validates root via config,
		// but we also want an account row for directories/ownership/UI consistency.
		if err := resourceClient.RegisterAccount(domain, "sa", "sa", adminEmail, rootPwd, rootPwd); err != nil {
			return fmt.Errorf("register sa: %w", err)
		}
	}

	// 6) Assign "admin" role to "sa"
	if err := resourceClient.AddAccountRole("sa", "admin"); err != nil &&
		!strings.Contains(strings.ToLower(err.Error()), "already") {
		return err
	}

	// 7) Ensure home folder + ownership
	userHome := config.GetDataDir() + "/files/users/sa@" + domain
	if !Utility.Exists(userHome) {
		if err := Utility.CreateDirIfNotExist(userHome); err != nil {
			return fmt.Errorf("create sa home: %w", err)
		}
	}
	if err := g.addResourceOwner("/users/sa@"+domain, "file", "sa@"+domain, rbacpb.SubjectType_ACCOUNT); err != nil {
		// non-fatal, but useful
		g.log.Warn("addResourceOwner", "err", err)
	}

	g.log.Info("BootstrapAdmin: ok", "domain", domain)
	return nil
}
