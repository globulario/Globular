package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/davecourtois/Utility"
	"github.com/globulario/services/golang/resource/resourcepb"
)

func (self *Globule) startRbacService() error {
	id := string(resourcepb.File_proto_resource_proto.Services().Get(1).FullName())
	rbac_server, err := self.startInternalService(id, resourcepb.File_proto_resource_proto.Path(), self.RbacPort, self.RbacProxy, self.Protocol == "https", self.unaryResourceInterceptor, self.streamResourceInterceptor)
	if err == nil {
		self.inernalServices = append(self.inernalServices, rbac_server)

		// Create the channel to listen on resource port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.RbacPort))
		if err != nil {
			log.Fatalf("could not start resource service %s: %s", self.getDomain(), err)
		}

		resourcepb.RegisterRbacServiceServer(rbac_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			// no web-rpc server.
			if err = rbac_server.Serve(lis); err != nil {
				log.Println(err)
			}
			s := self.getService(id)
			pid := getIntVal(s, "ProxyProcess")
			Utility.TerminateProcess(pid, 0)
			s.Store("ProxyProcess", -1)
			self.saveConfig()
		}()
	}

	return err
}

func (self *Globule) setEntityResourcePermissions(entity string, path string) error {
	// Here I will retreive the actual list of paths use by this user.
	data, err := self.permissions.GetItem(entity)

	paths := make([]string, 0)
	if err == nil {
		err := json.Unmarshal(data, &paths)
		if err != nil {
			return err
		}
	}

	if !Utility.Contains(paths, path) {
		paths = append(paths, path)
	} else {
		return nil // nothing todo here...
	}

	// simply marshal the permission and put it into the store.
	data, err = json.Marshal(paths)
	if err != nil {
		return err
	}

	return self.permissions.SetItem(entity, data)
}

func (self *Globule) setResourcePermissions(path string, permissions *resourcepb.Permissions) error {

	// First of all I need to remove the existing permission.
	self.deleteResourcePermissions(path, permissions)

	// Allowed resources
	allowed := permissions.Allowed
	for i := 0; i < len(allowed); i++ {

		// Accounts
		for j := 0; j < len(allowed[i].Accounts); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Accounts[j], path)
			if err != nil {
				return err
			}
		}

		// Groups
		for j := 0; j < len(allowed[i].Groups); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Groups[j], path)
			if err != nil {
				return err
			}
		}

		// Organizations
		for j := 0; j < len(allowed[i].Organizations); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Organizations[j], path)
			if err != nil {
				return err
			}
		}

		// Applications
		for j := 0; j < len(allowed[i].Applications); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Applications[j], path)
			if err != nil {
				return err
			}
		}

		// Peers
		for j := 0; j < len(allowed[i].Peers); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Peers[j], path)
			if err != nil {
				return err
			}
		}
	}

	// Denied resources
	denied := permissions.Denied
	for i := 0; i < len(denied); i++ {
		// Acccounts
		for j := 0; j < len(denied[i].Accounts); j++ {
			err := self.setEntityResourcePermissions(denied[i].Accounts[j], path)
			if err != nil {
				return err
			}
		}
		// Applications
		for j := 0; j < len(denied[i].Applications); j++ {
			err := self.setEntityResourcePermissions(denied[i].Applications[j], path)
			if err != nil {
				return err
			}
		}

		// Peers
		for j := 0; j < len(denied[i].Peers); j++ {
			err := self.setEntityResourcePermissions(denied[i].Peers[j], path)
			if err != nil {
				return err
			}
		}

		// Groups
		for j := 0; j < len(denied[i].Groups); j++ {
			err := self.setEntityResourcePermissions(denied[i].Groups[j], path)
			if err != nil {
				return err
			}
		}

		// Organizations
		for j := 0; j < len(denied[i].Organizations); j++ {
			err := self.setEntityResourcePermissions(denied[i].Organizations[j], path)
			if err != nil {
				return err
			}
		}
	}

	// Owned resources
	owners := permissions.Owners
	// Acccounts
	for j := 0; j < len(owners.Accounts); j++ {
		err := self.setEntityResourcePermissions(owners.Accounts[j], path)
		if err != nil {
			return err
		}
	}

	// Applications
	for j := 0; j < len(owners.Applications); j++ {
		err := self.setEntityResourcePermissions(owners.Applications[j], path)
		if err != nil {
			return err
		}
	}

	// Peers
	for j := 0; j < len(owners.Peers); j++ {
		err := self.setEntityResourcePermissions(owners.Peers[j], path)
		if err != nil {
			return err
		}
	}

	// Groups
	for j := 0; j < len(owners.Groups); j++ {
		err := self.setEntityResourcePermissions(owners.Groups[j], path)
		if err != nil {
			return err
		}
	}

	// Organizations
	for j := 0; j < len(owners.Organizations); j++ {
		err := self.setEntityResourcePermissions(owners.Organizations[j], path)
		if err != nil {
			return err
		}
	}

	// simply marshal the permission and put it into the store.
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	err = self.permissions.SetItem(path, data)
	if err != nil {
		return err
	}

	return nil
}

//* Set resource permissions this method will replace existing permission at once *
func (self *Globule) SetResourcePermissions(ctx context.Context, rqst *resourcepb.SetResourcePermissionsRqst) (*resourcepb.SetResourcePermissionsRqst, error) {
	err := self.setResourcePermissions(rqst.Path, rqst.Permissions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &resourcepb.SetResourcePermissionsRqst{}, nil
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

/**
 * Remove a resource path for an entity.
 */
func (self *Globule) deleteEntityResourcePermissions(entity string, path string) error {
	data, err := self.permissions.GetItem(entity)
	if err != nil {
		return err
	}

	paths := make([]string, 0)
	err = json.Unmarshal(data, &paths)
	if err != nil {
		return err
	}

	// Here I will
	if Utility.Contains(paths, path) {
		paths = remove(paths, path)
		data, err = json.Marshal(paths)
		if err != nil {
			return err
		}

		return self.permissions.SetItem(entity, data)
	}

	return nil
}

// Remouve a ressource permission
func (self *Globule) deleteResourcePermissions(path string, permissions *resourcepb.Permissions) error {

	// Allowed resources
	allowed := permissions.Allowed
	for i := 0; i < len(allowed); i++ {

		// Accounts
		for j := 0; j < len(allowed[i].Accounts); j++ {
			err := self.deleteEntityResourcePermissions(allowed[i].Accounts[j], path)
			if err != nil {
				return err
			}
		}

		// Groups
		for j := 0; j < len(allowed[i].Groups); j++ {
			err := self.deleteEntityResourcePermissions(allowed[i].Groups[j], path)
			if err != nil {
				return err
			}
		}

		// Organizations
		for j := 0; j < len(allowed[i].Organizations); j++ {
			err := self.deleteEntityResourcePermissions(allowed[i].Organizations[j], path)
			if err != nil {
				return err
			}
		}

		// Applications
		for j := 0; j < len(allowed[i].Applications); j++ {
			err := self.deleteEntityResourcePermissions(allowed[i].Applications[j], path)
			if err != nil {
				return err
			}
		}

		// Peers
		for j := 0; j < len(allowed[i].Peers); j++ {
			err := self.deleteEntityResourcePermissions(allowed[i].Peers[j], path)
			if err != nil {
				return err
			}
		}
	}

	// Denied resources
	denied := permissions.Denied
	for i := 0; i < len(denied); i++ {
		// Acccounts
		for j := 0; j < len(denied[i].Accounts); j++ {
			err := self.deleteEntityResourcePermissions(denied[i].Accounts[j], path)
			if err != nil {
				return err
			}
		}
		// Applications
		for j := 0; j < len(denied[i].Applications); j++ {
			err := self.deleteEntityResourcePermissions(denied[i].Applications[j], path)
			if err != nil {
				return err
			}
		}

		// Peers
		for j := 0; j < len(denied[i].Peers); j++ {
			err := self.deleteEntityResourcePermissions(denied[i].Peers[j], path)
			if err != nil {
				return err
			}
		}

		// Groups
		for j := 0; j < len(denied[i].Groups); j++ {
			err := self.deleteEntityResourcePermissions(denied[i].Groups[j], path)
			if err != nil {
				return err
			}
		}

		// Organizations
		for j := 0; j < len(denied[i].Organizations); j++ {
			err := self.deleteEntityResourcePermissions(denied[i].Organizations[j], path)
			if err != nil {
				return err
			}
		}
	}

	// Owned resources
	owners := permissions.Owners
	// Acccounts
	for j := 0; j < len(owners.Accounts); j++ {
		err := self.deleteEntityResourcePermissions(owners.Accounts[j], path)
		if err != nil {
			return err
		}
	}

	// Applications
	for j := 0; j < len(owners.Applications); j++ {
		err := self.deleteEntityResourcePermissions(owners.Applications[j], path)
		if err != nil {
			return err
		}
	}

	// Peers
	for j := 0; j < len(owners.Peers); j++ {
		err := self.deleteEntityResourcePermissions(owners.Peers[j], path)
		if err != nil {
			return err
		}
	}

	// Groups
	for j := 0; j < len(owners.Groups); j++ {
		err := self.deleteEntityResourcePermissions(owners.Groups[j], path)
		if err != nil {
			return err
		}
	}

	// Organizations
	for j := 0; j < len(owners.Organizations); j++ {
		err := self.deleteEntityResourcePermissions(owners.Organizations[j], path)
		if err != nil {
			return err
		}
	}

	return self.permissions.RemoveItem(path)

}

func (self *Globule) getResourcePermissions(path string) (*resourcepb.Permissions, error) {
	data, err := self.permissions.GetItem(path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	permissions := new(resourcepb.Permissions)
	err = json.Unmarshal(data, &permissions)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

//* Delete a resource permissions (when a resource is deleted) *
func (self *Globule) DeleteResourcePermissions(ctx context.Context, rqst *resourcepb.DeleteResourcePermissionsRqst) (*resourcepb.DeleteResourcePermissionsRqst, error) {

	permissions, err := self.getResourcePermissions(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = self.deleteResourcePermissions(rqst.Path, permissions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteResourcePermissionsRqst{}, nil
}

//* Delete a specific resource permission *
func (self *Globule) DeleteResourcePermission(ctx context.Context, rqst *resourcepb.DeleteResourcePermissionRqst) (*resourcepb.DeleteResourcePermissionRqst, error) {

	permissions, err := self.getResourcePermissions(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	if rqst.Type == resourcepb.PermissionType_ALLOWED {
		// Remove the permission from the allowed permission
		allowed := make([]*resourcepb.Permission, 0)
		for i := 0; i < len(permissions.Allowed); i++ {
			if permissions.Allowed[i].Name != rqst.Name {
				allowed = append(allowed, permissions.Allowed[i])
			}
		}
		permissions.Allowed = allowed
	} else if rqst.Type == resourcepb.PermissionType_DENIED {
		// Remove the permission from the allowed permission.
		denied := make([]*resourcepb.Permission, 0)
		for i := 0; i < len(permissions.Denied); i++ {
			if permissions.Denied[i].Name != rqst.Name {
				denied = append(denied, permissions.Denied[i])
			}
		}
		permissions.Denied = denied
	}
	err = self.setResourcePermissions(rqst.Path, permissions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteResourcePermissionRqst{}, nil
}

//* Get the ressource Permission.
func (self *Globule) GetResourcePermission(ctx context.Context, rqst *resourcepb.GetResourcePermissionRqst) (*resourcepb.GetResourcePermissionRsp, error) {
	permissions, err := self.getResourcePermissions(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Search on allowed permission
	if rqst.Type == resourcepb.PermissionType_ALLOWED {
		for i := 0; i < len(permissions.Allowed); i++ {
			if permissions.Allowed[i].Name == rqst.Name {
				return &resourcepb.GetResourcePermissionRsp{Permission: permissions.Allowed[i]}, nil
			}
		}
	} else if rqst.Type == resourcepb.PermissionType_DENIED { // search in denied permissions.

		for i := 0; i < len(permissions.Denied); i++ {
			if permissions.Denied[i].Name == rqst.Name {
				return &resourcepb.GetResourcePermissionRsp{Permission: permissions.Allowed[i]}, nil
			}
		}
	}

	return nil, status.Errorf(
		codes.Internal,
		Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No permission found with name "+rqst.Name)))
}

//* Set specific resource permission  ex. read permission... *
func (self *Globule) SetResourcePermission(ctx context.Context, rqst *resourcepb.SetResourcePermissionRqst) (*resourcepb.SetResourcePermissionRsp, error) {
	permissions, err := self.getResourcePermissions(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Remove the permission from the allowed permission
	if rqst.Type == resourcepb.PermissionType_ALLOWED {
		allowed := make([]*resourcepb.Permission, 0)
		for i := 0; i < len(permissions.Allowed); i++ {
			if permissions.Allowed[i].Name == rqst.Permission.Name {
				allowed = append(allowed, permissions.Allowed[i])
			} else {
				allowed = append(allowed, rqst.Permission)
			}
		}
		permissions.Allowed = allowed
	} else if rqst.Type == resourcepb.PermissionType_DENIED {

		// Remove the permission from the allowed permission.
		denied := make([]*resourcepb.Permission, 0)
		for i := 0; i < len(permissions.Denied); i++ {
			if permissions.Denied[i].Name == rqst.Permission.Name {
				denied = append(denied, permissions.Denied[i])
			} else {
				denied = append(denied, rqst.Permission)
			}
		}
		permissions.Denied = denied
	}
	err = self.setResourcePermissions(rqst.Path, permissions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.SetResourcePermissionRsp{}, nil
}

//* Get resource permissions *
func (self *Globule) GetResourcePermissions(ctx context.Context, rqst *resourcepb.GetResourcePermissionsRqst) (*resourcepb.GetResourcePermissionsRsp, error) {
	permissions, err := self.getResourcePermissions(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.GetResourcePermissionsRsp{Permissions: permissions}, nil
}

//* Add resource owner do nothing if it already exist
func (self *Globule) AddResourceOwner(ctx context.Context, rqst *resourcepb.AddResourceOwnerRqst) (*resourcepb.AddResourceOwnerRsp, error) {
	permissions, err := self.getResourcePermissions(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Owned resources
	owners := permissions.Owners
	if rqst.Type == resourcepb.SubjectType_ACCOUNT {
		if !Utility.Contains(owners.Accounts, rqst.Subject) {
			owners.Accounts = append(owners.Accounts, rqst.Subject)
		}
	} else if rqst.Type == resourcepb.SubjectType_APPLICATION {
		if !Utility.Contains(owners.Applications, rqst.Subject) {
			owners.Applications = append(owners.Applications, rqst.Subject)
		}
	} else if rqst.Type == resourcepb.SubjectType_GROUP {
		if !Utility.Contains(owners.Groups, rqst.Subject) {
			owners.Groups = append(owners.Groups, rqst.Subject)
		}
	} else if rqst.Type == resourcepb.SubjectType_ORGANIZATION {
		if !Utility.Contains(owners.Organizations, rqst.Subject) {
			owners.Organizations = append(owners.Organizations, rqst.Subject)
		}
	} else if rqst.Type == resourcepb.SubjectType_PEER {
		if !Utility.Contains(owners.Peers, rqst.Subject) {
			owners.Peers = append(owners.Peers, rqst.Subject)
		}
	}
	permissions.Owners = owners
	err = self.setResourcePermissions(rqst.Path, permissions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.AddResourceOwnerRsp{}, nil
}

func (self *Globule) removeResourceOwner(owner string, subjectType resourcepb.SubjectType, path string) error {
	permissions, err := self.getResourcePermissions(path)
	if err != nil {
		return err
	}

	// Owned resources
	owners := permissions.Owners
	if subjectType == resourcepb.SubjectType_ACCOUNT {
		if Utility.Contains(owners.Accounts, owner) {
			owners.Accounts = remove(owners.Accounts, owner)
		}
	} else if subjectType == resourcepb.SubjectType_APPLICATION {
		if Utility.Contains(owners.Applications, owner) {
			owners.Applications = remove(owners.Applications, owner)
		}
	} else if subjectType == resourcepb.SubjectType_GROUP {
		if Utility.Contains(owners.Groups, owner) {
			owners.Groups = remove(owners.Groups, owner)
		}
	} else if subjectType == resourcepb.SubjectType_ORGANIZATION {
		if Utility.Contains(owners.Organizations, owner) {
			owners.Organizations = remove(owners.Organizations, owner)
		}
	} else if subjectType == resourcepb.SubjectType_PEER {
		if Utility.Contains(owners.Peers, owner) {
			owners.Peers = remove(owners.Peers, owner)
		}
	}

	permissions.Owners = owners
	err = self.setResourcePermissions(path, permissions)
	if err != nil {
		return err
	}

	return nil
}

// Remove a Subject from denied list and allowed list.
func (self *Globule) removeResourceSubject(subject string, subjectType resourcepb.SubjectType, path string) error {
	permissions, err := self.getResourcePermissions(path)
	if err != nil {
		return err
	}

	// Allowed resources
	allowed := permissions.Allowed
	for i := 0; i < len(allowed); i++ {
		// Accounts
		if subjectType == resourcepb.SubjectType_ACCOUNT {
			accounts := make([]string, 0)
			for j := 0; j < len(allowed[i].Accounts); j++ {
				if subject != allowed[i].Accounts[j] {
					accounts = append(accounts, allowed[i].Accounts[j])
				}

			}
			allowed[i].Accounts = accounts
		}

		// Groups
		if subjectType == resourcepb.SubjectType_GROUP {
			groups := make([]string, 0)
			for j := 0; j < len(allowed[i].Groups); j++ {
				if subject != allowed[i].Groups[j] {
					groups = append(groups, allowed[i].Groups[j])
				}
			}
			allowed[i].Groups = groups
		}

		// Organizations
		if subjectType == resourcepb.SubjectType_ORGANIZATION {
			organizations := make([]string, 0)
			for j := 0; j < len(allowed[i].Organizations); j++ {
				if subject != allowed[i].Organizations[j] {
					organizations = append(organizations, allowed[i].Organizations[j])
				}
			}
			allowed[i].Organizations = organizations
		}

		// Applications
		if subjectType == resourcepb.SubjectType_APPLICATION {
			applications := make([]string, 0)
			for j := 0; j < len(allowed[i].Applications); j++ {
				if subject != allowed[i].Applications[j] {
					applications = append(applications, allowed[i].Applications[j])
				}
			}
			allowed[i].Applications = applications
		}

		// Peers
		if subjectType == resourcepb.SubjectType_PEER {
			peers := make([]string, 0)
			for j := 0; j < len(allowed[i].Peers); j++ {
				if subject != allowed[i].Peers[j] {
					peers = append(peers, allowed[i].Peers[j])
				}
			}
			allowed[i].Peers = peers
		}
	}

	// Denied resources
	denied := permissions.Denied
	for i := 0; i < len(denied); i++ {
		// Accounts
		if subjectType == resourcepb.SubjectType_ACCOUNT {
			accounts := make([]string, 0)
			for j := 0; j < len(denied[i].Accounts); j++ {
				if subject != denied[i].Accounts[j] {
					accounts = append(accounts, denied[i].Accounts[j])
				}

			}
			denied[i].Accounts = accounts
		}

		// Groups
		if subjectType == resourcepb.SubjectType_GROUP {
			groups := make([]string, 0)
			for j := 0; j < len(denied[i].Groups); j++ {
				if subject != denied[i].Groups[j] {
					groups = append(groups, denied[i].Groups[j])
				}
			}
			denied[i].Groups = groups
		}

		// Organizations
		if subjectType == resourcepb.SubjectType_ORGANIZATION {
			organizations := make([]string, 0)
			for j := 0; j < len(denied[i].Organizations); j++ {
				if subject != denied[i].Organizations[j] {
					organizations = append(organizations, denied[i].Organizations[j])
				}
			}
			denied[i].Organizations = organizations
		}

		// Applications
		if subjectType == resourcepb.SubjectType_APPLICATION {
			applications := make([]string, 0)
			for j := 0; j < len(denied[i].Applications); j++ {
				if subject != denied[i].Applications[j] {
					applications = append(applications, denied[i].Applications[j])
				}
			}
			denied[i].Applications = applications
		}

		// Peers
		if subjectType == resourcepb.SubjectType_PEER {
			peers := make([]string, 0)
			for j := 0; j < len(denied[i].Peers); j++ {
				if subject != denied[i].Peers[j] {
					peers = append(peers, denied[i].Peers[j])
				}
			}
			denied[i].Peers = peers
		}
	}

	err = self.setResourcePermissions(path, permissions)
	if err != nil {
		return err
	}

	return nil
}

//* Remove resource owner
func (self *Globule) RemoveResourceOwner(ctx context.Context, rqst *resourcepb.RemoveResourceOwnerRqst) (*resourcepb.RemoveResourceOwnerRsp, error) {
	err := self.removeResourceOwner(rqst.Subject, rqst.Type, rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.RemoveResourceOwnerRsp{}, nil
}

//* That function must be call when a subject is removed to clean up permissions.
func (self *Globule) DeleteAllAccess(ctx context.Context, rqst *resourcepb.DeleteAllAccessRqst) (*resourcepb.DeleteAllAccessRsp, error) {

	// Here I must remove the subject from all permissions.
	data, err := self.permissions.GetItem(rqst.Subject)
	paths := make([]string, 0)
	err = json.Unmarshal(data, &paths)
	for i := 0; i < len(paths); i++ {

		// Remove from owner
		self.removeResourceOwner(rqst.Subject, rqst.Type, paths[i])

		// Remove from subject.
		self.removeResourceSubject(rqst.Subject, rqst.Type, paths[i])

	}

	err = self.permissions.RemoveItem(rqst.Subject)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteAllAccessRsp{}, nil
}

// Return  accessAllowed, accessDenied, error
func (self *Globule) validateAccess(subject string, subjectType resourcepb.SubjectType, name string, path string) (bool, bool, error) {
	permissions, err := self.getResourcePermissions(path)
	if err != nil {
		return false, false, err
	}
	// Test if the Subject is owner of the ressource in that case I will git him access.
	owners := permissions.Owners
	isOwner := false
	subjectStr := ""
	if subjectType == resourcepb.SubjectType_ACCOUNT {
		subjectStr = "Account"
		if Utility.Contains(owners.Accounts, subject) {
			isOwner = true
		}
	} else if subjectType == resourcepb.SubjectType_APPLICATION {
		subjectStr = "Application"
		if Utility.Contains(owners.Applications, subject) {
			isOwner = true
		}
	} else if subjectType == resourcepb.SubjectType_GROUP {
		subjectStr = "Group"
		if Utility.Contains(owners.Groups, subject) {
			isOwner = true
		}
	} else if subjectType == resourcepb.SubjectType_ORGANIZATION {
		subjectStr = "Organization"
		if Utility.Contains(owners.Organizations, subject) {
			isOwner = true
		}
	} else if subjectType == resourcepb.SubjectType_PEER {
		subjectStr = "Peer"
		if Utility.Contains(owners.Peers, subject) {
			isOwner = true
		}
	}

	// If the user is the owner no other validation are required.
	if isOwner {
		log.Println("----> is owner ", subject, path)
		return true, false, nil
	}

	// First I will validate that the permission is not denied...
	var denied *resourcepb.Permission
	for i := 0; i < len(permissions.Denied); i++ {
		if permissions.Denied[i].Name == name {
			denied = permissions.Denied[i]
			break
		}
	}

	////////////////////// Test if the access is denied first. /////////////////////

	accessDenied := false

	// Here the Subject is not the owner...
	if denied != nil {
		if subjectType == resourcepb.SubjectType_ACCOUNT {

			// Here the subject is an account.
			if denied.Accounts != nil {
				accessDenied = Utility.Contains(denied.Accounts, subject)

				// The access is not denied for the account itself, I will validate
				// that the account is not part of denied group.
				if !accessDenied {
					// I will test if one of the group account if part of hare access denied.
					p, err := self.getPersistenceStore()
					if err != nil {
						return false, false, err
					}

					// Here I will test if a newer token exist for that user if it's the case
					// I will not refresh that token.
					values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+subject+`"}`, ``)
					if err != nil {
						return false, false, errors.New("No account named " + subject + " exist!")
					}

					// from the account I will get the list of group.
					account := values.(map[string]interface{})
					if account["groups"] != nil {
						groups := []interface{}(account["groups"].(primitive.A))
						if groups != nil {
							for i := 0; i < len(groups); i++ {
								groupId := groups[i].(map[string]interface{})["$id"].(string)
								_, accessDenied_, _ := self.validateAccess(groupId, resourcepb.SubjectType_GROUP, name, path)
								if accessDenied_ {
									return false, true, nil
								}
							}
						}
					}

					// from the account I will get the list of group.
					if account["organizations"] != nil {
						organizations := []interface{}(account["organizations"].(primitive.A))
						if organizations != nil {
							for i := 0; i < len(organizations); i++ {
								organizationId := organizations[i].(map[string]interface{})["$id"].(string)
								_, accessDenied_, _ := self.validateAccess(organizationId, resourcepb.SubjectType_ORGANIZATION, name, path)
								if accessDenied_ {
									return false, true, nil
								}
							}
						}
					}
				}
			}

		} else if subjectType == resourcepb.SubjectType_APPLICATION {
			// Here the Subject is an application.
			accessDenied = Utility.Contains(denied.Applications, subject)
		} else if subjectType == resourcepb.SubjectType_GROUP {
			// Here the Subject is a group
			if denied.Groups != nil {
				accessDenied = Utility.Contains(denied.Groups, subject)

				// The access is not denied for the account itself, I will validate
				// that the account is not part of denied group.
				if !accessDenied {
					// I will test if one of the group account if part of hare access denied.
					p, err := self.getPersistenceStore()
					if err != nil {
						return false, false, err
					}

					// Here I will test if a newer token exist for that user if it's the case
					// I will not refresh that token.
					values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Groups", `{"_id":"`+subject+`"}`, ``)
					if err != nil {
						return false, false, errors.New("No account named " + subject + " exist!")
					}

					// from the account I will get the list of group.
					group := values.(map[string]interface{})
					if group["organizations"] != nil {
						organizations := []interface{}(group["organizations"].(primitive.A))
						if organizations != nil {
							for i := 0; i < len(organizations); i++ {
								organizationId := organizations[i].(map[string]interface{})["$id"].(string)
								_, accessDenied_, _ := self.validateAccess(organizationId, resourcepb.SubjectType_ORGANIZATION, name, path)
								if accessDenied_ {
									return false, true, errors.New("Access denied for " + subjectStr + " " + organizationId + "!")
								}
							}
						}
					}
				}
			}
		} else if subjectType == resourcepb.SubjectType_ORGANIZATION {
			// Here the Subject is an Organisations.
			accessDenied = Utility.Contains(denied.Organizations, subject)
		} else if subjectType == resourcepb.SubjectType_PEER {
			// Here the Subject is a Peer.
			accessDenied = Utility.Contains(denied.Peers, subject)
		}
	}

	if accessDenied {
		err := errors.New("Access denied for " + subjectStr + " " + subject + "!")
		return false, true, err
	}

	var allowed *resourcepb.Permission
	for i := 0; i < len(permissions.Allowed); i++ {
		if permissions.Allowed[i].Name == name {
			allowed = permissions.Allowed[i]
			break
		}
	}

	hasAccess := false

	// Test if the access is allowed
	if subjectType == resourcepb.SubjectType_ACCOUNT {
		hasAccess = Utility.Contains(allowed.Accounts, subject)
		if !hasAccess {
			// I will test if one of the group account if part of hare access denied.
			p, err := self.getPersistenceStore()
			if err != nil {
				return false, false, err
			}

			// Here I will test if a newer token exist for that user if it's the case
			// I will not refresh that token.
			values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Accounts", `{"_id":"`+subject+`"}`, ``)
			if err == nil {
				// from the account I will get the list of group.
				account := values.(map[string]interface{})
				if account["groups"] != nil {
					groups := []interface{}(account["groups"].(primitive.A))
					if groups != nil {
						for i := 0; i < len(groups); i++ {
							groupId := groups[i].(map[string]interface{})["$id"].(string)
							hasAccess_, _, _ := self.validateAccess(groupId, resourcepb.SubjectType_GROUP, name, path)
							if hasAccess_ {
								return true, false, nil
							}
						}
					}
				}

				// from the account I will get the list of group.
				if account["organizations"] != nil {
					organizations := []interface{}(account["organizations"].(primitive.A))
					if organizations != nil {
						for i := 0; i < len(organizations); i++ {
							organizationId := organizations[i].(map[string]interface{})["$id"].(string)
							hasAccess_, _, _ := self.validateAccess(organizationId, resourcepb.SubjectType_ORGANIZATION, name, path)
							if hasAccess_ {
								return true, false, nil
							}
						}
					}
				}
			}
		}

	} else if subjectType == resourcepb.SubjectType_GROUP {
		// validate the group access
		hasAccess = Utility.Contains(allowed.Groups, subject)
		if !hasAccess {
			// I will test if one of the group account if part of hare access denied.
			p, err := self.getPersistenceStore()
			if err != nil {
				return false, false, err
			}

			// Here I will test if a newer token exist for that user if it's the case
			// I will not refresh that token.
			values, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Groups", `{"_id":"`+subject+`"}`, ``)
			if err == nil {
				// Validate group organization
				group := values.(map[string]interface{})
				if group["organizations"] != nil {
					organizations := []interface{}(group["organizations"].(primitive.A))
					if organizations != nil {
						for i := 0; i < len(organizations); i++ {
							organizationId := organizations[i].(map[string]interface{})["$id"].(string)
							hasAccess_, _, _ := self.validateAccess(organizationId, resourcepb.SubjectType_ORGANIZATION, name, path)
							if hasAccess_ {
								return true, false, nil
							}
						}
					}
				}
			}
		}
	} else if subjectType == resourcepb.SubjectType_ORGANIZATION {
		hasAccess = Utility.Contains(allowed.Organizations, subject)
	} else if subjectType == resourcepb.SubjectType_PEER {
		// Here the Subject is an application.
		hasAccess = Utility.Contains(allowed.Peers, subject)
	} else if subjectType == resourcepb.SubjectType_APPLICATION {
		// Here the Subject is an application.
		hasAccess = Utility.Contains(allowed.Applications, subject)
	}

	if !hasAccess {
		err := errors.New("Access denied for " + subjectStr + " " + subject + "!")
		return false, false, err
	}

	// The permission is set.
	return true, false, nil
}

//* Validate if a account can get access to a given ressource for a given operation (read, write...) That function is recursive. *
func (self *Globule) ValidateAccess(ctx context.Context, rqst *resourcepb.ValidateAccessRqst) (*resourcepb.ValidateAccessRsp, error) {
	hasAccess, accessDenied, err := self.validateAccess(rqst.Subject, rqst.Type, rqst.Permission, rqst.Path)

	if err != nil || !hasAccess || accessDenied {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The permission is set.
	return &resourcepb.ValidateAccessRsp{Result: true}, nil
}
