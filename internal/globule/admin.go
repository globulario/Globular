package globule

import (
	"fmt"
	"strings"

	"github.com/globulario/services/golang/authentication/authentication_client"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/rbac/rbacpb"
	serviceManagerClient "github.com/globulario/services/golang/services_manager/services_manager_client"
	Utility "github.com/globulario/utility"
)

// getAuthenticationClient returns an authenticated client to the authentication service.
func (g *Globule) getAuthenticationClient(address string) (*authentication_client.Authentication_Client, error) {
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		g.log.Error("fail to create authentication client:", "err", err)
		return nil, err
	}
	return authenticationClient, nil
}

// getServiceManagerClient returns a services manager client.
func (g *Globule) getServiceManagerClient(domain string) (*serviceManagerClient.Services_Manager_Client, error) {
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

	domain := strings.TrimSpace(g.Domain)
	adminEmail := strings.TrimSpace(g.AdminEmail)
	rootPwd := strings.TrimSpace(g.RootPassword)
	if domain == "" || adminEmail == "" || rootPwd == "" {
		return fmt.Errorf("missing Domain/AdminEmail/RootPassword in config")
	}

	token, err := authClient.Authenticate("sa", rootPwd)
	if err != nil {
		return fmt.Errorf("authenticate sa: %w", err)
	}

	servicesManager, err := g.getServiceManagerClient(address)
	if err != nil {
		return fmt.Errorf("services manager client: %w", err)
	}
	defer servicesManager.Close()

	actions, err := servicesManager.GetAllActions()
	if err != nil {
		return fmt.Errorf("get all actions: %w", err)
	}
	if err := resourceClient.CreateRole(token, "admin", "admin", actions); err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			return fmt.Errorf("create admin role: %w", err)
		}
	}

	accts, err := resourceClient.GetAccounts(`{"_id":"sa"}`)
	if err != nil {
		return fmt.Errorf("get accounts: %w", err)
	}
	if len(accts) == 0 {
		if err := resourceClient.RegisterAccount(domain, "sa", "sa", adminEmail, rootPwd, rootPwd); err != nil {
			return fmt.Errorf("register sa: %w", err)
		}
	}

	if err := resourceClient.AddAccountRole(token, "sa", "admin"); err != nil &&
		!strings.Contains(strings.ToLower(err.Error()), "already") {
		return err
	}

	userHome := config.GetDataDir() + "/files/users/sa@" + domain
	if !Utility.Exists(userHome) {
		if err := Utility.CreateDirIfNotExist(userHome); err != nil {
			return fmt.Errorf("create sa home: %w", err)
		}
	}

	if err := g.AddResourceOwner("/users/sa@"+domain, "file", "sa@"+domain, rbacpb.SubjectType_ACCOUNT); err != nil {
		g.log.Warn("addResourceOwner", "err", err)
	}

	g.log.Info("BootstrapAdmin: ok", "domain", domain)
	return nil
}
