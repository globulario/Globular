package globule

import (
	"fmt"

	config "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/globular_client"
	rbac_client "github.com/globulario/services/golang/rbac/rbac_client"
	"github.com/globulario/services/golang/rbac/rbacpb"
	Utility "github.com/globulario/utility"
)

// getRbacClient returns an RBAC client bound to the current node's address.
func getRbacClient() (*rbac_client.Rbac_Client, error) {
	addr, _ := config.GetAddress()
	Utility.RegisterFunction("NewRbacService_Client", rbac_client.NewRbacService_Client)
	c, err := globular_client.GetClient(addr, "rbac.RbacService", "NewRbacService_Client")
	if err != nil {
		return nil, fmt.Errorf("getRbacClient: %w", err)
	}
	return c.(*rbac_client.Rbac_Client), nil
}

// validateAccess bridges main.go’s accessControl to RBAC.ValidateAccess.
//
// subject      -> "user@domain" or app id
// subjectType  -> rbacpb.SubjectType_ACCOUNT or _APPLICATION
// action       -> arbitrary action string (e.g., "files.read")
// path         -> resource path being accessed
//
// returns: hasAccess, hasAccessDenied, err
func (g *Globule) ValidateAccess(subject string, subjectType rbacpb.SubjectType, action, path string) (bool, bool, error) {
	rbac, err := getRbacClient()
	if err != nil {
		return false, false, err
	}
	return rbac.ValidateAccess(subject, subjectType, action, path)
}
