package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"

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

func (self *Globule) setEntityResourcePermissions(entity string, permission_type string, name string, path string) error {
	// Here I will retreive the actual list of paths use by this user.
	data, err := self.permissions.GetItem(entity + ":" + permission_type + ":" + name)
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

	return self.permissions.SetItem(entity+":"+permission_type+":"+name, data)
}

func (self *Globule) setResourcePermissions(path string, permissions *resourcepb.Permissions) error {

	// First of all I need to remove the existing permission.
	self.deleteResourcePermissions(path, permissions)

	// Allowed resources
	allowed := permissions.Allowed
	for i := 0; i < len(allowed); i++ {

		// Accounts
		for j := 0; j < len(allowed[i].Accounts); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Accounts[j], "allowed", allowed[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Groups
		for j := 0; j < len(allowed[i].Groups); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Groups[j], "allowed", allowed[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Organizations
		for j := 0; j < len(allowed[i].Organizations); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Organizations[j], "allowed", allowed[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Applications
		for j := 0; j < len(allowed[i].Applications); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Applications[j], "allowed", allowed[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Peers
		for j := 0; j < len(allowed[i].Peers); j++ {
			err := self.setEntityResourcePermissions(allowed[i].Peers[j], "allowed", allowed[i].Name, path)
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
			err := self.setEntityResourcePermissions(denied[i].Accounts[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}
		// Applications
		for j := 0; j < len(denied[i].Applications); j++ {
			err := self.setEntityResourcePermissions(denied[i].Applications[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Peers
		for j := 0; j < len(denied[i].Peers); j++ {
			err := self.setEntityResourcePermissions(denied[i].Peers[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Groups
		for j := 0; j < len(denied[i].Groups); j++ {
			err := self.setEntityResourcePermissions(denied[i].Groups[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Organizations
		for j := 0; j < len(denied[i].Organizations); j++ {
			err := self.setEntityResourcePermissions(denied[i].Organizations[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}
	}

	// Owned resources
	owners := permissions.Owners
	// Acccounts
	for j := 0; j < len(owners.Accounts); j++ {
		err := self.setEntityResourcePermissions(owners.Accounts[j], "owned", owners.Name, path)
		if err != nil {
			return err
		}
	}

	// Applications
	for j := 0; j < len(owners.Applications); j++ {
		err := self.setEntityResourcePermissions(owners.Applications[j], "owned", owners.Name, path)
		if err != nil {
			return err
		}
	}

	// Peers
	for j := 0; j < len(owners.Peers); j++ {
		err := self.setEntityResourcePermissions(owners.Peers[j], "owned", owners.Name, path)
		if err != nil {
			return err
		}
	}

	// Groups
	for j := 0; j < len(owners.Groups); j++ {
		err := self.setEntityResourcePermissions(owners.Groups[j], "owned", owners.Name, path)
		if err != nil {
			return err
		}
	}

	// Organizations
	for j := 0; j < len(owners.Organizations); j++ {
		err := self.setEntityResourcePermissions(owners.Organizations[j], "owned", owners.Name, path)
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
func (self *Globule) deleteEntityResourcePermissions(entity string, permission_type string, name string, path string) error {
	data, err := self.permissions.GetItem(entity + ":" + permission_type + ":" + name)
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

		return self.permissions.SetItem(entity+":"+permission_type+":"+name, data)
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
			err := self.deleteEntityResourcePermissions(allowed[i].Accounts[j], "allowed", allowed[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Groups
		for j := 0; j < len(allowed[i].Groups); j++ {
			err := self.deleteEntityResourcePermissions(allowed[i].Groups[j], "allowed", allowed[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Organizations
		for j := 0; j < len(allowed[i].Organizations); j++ {
			err := self.deleteEntityResourcePermissions(allowed[i].Organizations[j], "allowed", allowed[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Applications
		for j := 0; j < len(allowed[i].Applications); j++ {
			err := self.deleteEntityResourcePermissions(allowed[i].Applications[j], "allowed", allowed[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Peers
		for j := 0; j < len(allowed[i].Peers); j++ {
			err := self.deleteEntityResourcePermissions(allowed[i].Peers[j], "allowed", allowed[i].Name, path)
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
			err := self.deleteEntityResourcePermissions(denied[i].Accounts[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}
		// Applications
		for j := 0; j < len(denied[i].Applications); j++ {
			err := self.deleteEntityResourcePermissions(denied[i].Applications[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Peers
		for j := 0; j < len(denied[i].Peers); j++ {
			err := self.deleteEntityResourcePermissions(denied[i].Peers[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Groups
		for j := 0; j < len(denied[i].Groups); j++ {
			err := self.deleteEntityResourcePermissions(denied[i].Groups[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}

		// Organizations
		for j := 0; j < len(denied[i].Organizations); j++ {
			err := self.deleteEntityResourcePermissions(denied[i].Organizations[j], "denied", denied[i].Name, path)
			if err != nil {
				return err
			}
		}
	}

	// Owned resources
	owners := permissions.Owners
	// Acccounts
	for j := 0; j < len(owners.Accounts); j++ {
		err := self.deleteEntityResourcePermissions(owners.Accounts[j], "owned", owners.Name, path)
		if err != nil {
			return err
		}
	}

	// Applications
	for j := 0; j < len(owners.Applications); j++ {
		err := self.deleteEntityResourcePermissions(owners.Applications[j], "owned", owners.Name, path)
		if err != nil {
			return err
		}
	}

	// Peers
	for j := 0; j < len(owners.Peers); j++ {
		err := self.deleteEntityResourcePermissions(owners.Peers[j], "owned", owners.Name, path)
		if err != nil {
			return err
		}
	}

	// Groups
	for j := 0; j < len(owners.Groups); j++ {
		err := self.deleteEntityResourcePermissions(owners.Groups[j], "owned", owners.Name, path)
		if err != nil {
			return err
		}
	}

	// Organizations
	for j := 0; j < len(owners.Organizations); j++ {
		err := self.deleteEntityResourcePermissions(owners.Organizations[j], "owned", owners.Name, path)
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

	// Remove the permission from the allowed permission
	allowed := make([]*resourcepb.Permission, 0)
	for i := 0; i < len(permissions.Allowed); i++ {
		if permissions.Allowed[i].Name != rqst.Name {
			allowed = append(allowed, permissions.Allowed[i])
		}
	}
	permissions.Allowed = allowed

	// Remove the permission from the allowed permission.
	denied := make([]*resourcepb.Permission, 0)
	for i := 0; i < len(permissions.Denied); i++ {
		if permissions.Denied[i].Name != rqst.Name {
			allowed = append(denied, permissions.Denied[i])
		}
	}
	permissions.Denied = denied

	err = self.setResourcePermissions(rqst.Path, permissions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteResourcePermissionRqst{}, nil
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
	allowed := make([]*resourcepb.Permission, 0)
	for i := 0; i < len(permissions.Allowed); i++ {
		if permissions.Allowed[i].Name == rqst.Permission.Name {
			allowed = append(allowed, permissions.Allowed[i])
		} else {
			allowed = append(allowed, rqst.Permission)
		}
	}
	permissions.Allowed = allowed

	// Remove the permission from the allowed permission.
	denied := make([]*resourcepb.Permission, 0)
	for i := 0; i < len(permissions.Denied); i++ {
		if permissions.Denied[i].Name == rqst.Permission.Name {
			allowed = append(allowed, permissions.Allowed[i])
		} else {
			allowed = append(allowed, rqst.Permission)
		}
	}
	permissions.Denied = denied

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
	if rqst.Subject == resourcepb.SubjectType_ACCOUNT {
		if !Utility.Contains(owners.Accounts, rqst.Owner) {
			owners.Accounts = append(owners.Accounts, rqst.Owner)
		}
	} else if rqst.Subject == resourcepb.SubjectType_APPLICATION {
		if !Utility.Contains(owners.Applications, rqst.Owner) {
			owners.Applications = append(owners.Applications, rqst.Owner)
		}
	} else if rqst.Subject == resourcepb.SubjectType_GROUP {
		if !Utility.Contains(owners.Groups, rqst.Owner) {
			owners.Groups = append(owners.Groups, rqst.Owner)
		}
	} else if rqst.Subject == resourcepb.SubjectType_ORGANIZATION {
		if !Utility.Contains(owners.Organizations, rqst.Owner) {
			owners.Organizations = append(owners.Organizations, rqst.Owner)
		}
	} else if rqst.Subject == resourcepb.SubjectType_PEER {
		if !Utility.Contains(owners.Peers, rqst.Owner) {
			owners.Peers = append(owners.Peers, rqst.Owner)
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

//* Remove resource owner
func (self *Globule) RemoveResourceOwner(ctx context.Context, rqst *resourcepb.AddResourceOwnerRqst) (*resourcepb.AddResourceOwnerRsp, error) {
	permissions, err := self.getResourcePermissions(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Owned resources
	owners := permissions.Owners
	if rqst.Subject == resourcepb.SubjectType_ACCOUNT {
		if Utility.Contains(owners.Accounts, rqst.Owner) {
			owners.Accounts = remove(owners.Accounts, rqst.Owner)
		}
	} else if rqst.Subject == resourcepb.SubjectType_APPLICATION {
		if Utility.Contains(owners.Applications, rqst.Owner) {
			owners.Applications = remove(owners.Applications, rqst.Owner)
		}
	} else if rqst.Subject == resourcepb.SubjectType_GROUP {
		if Utility.Contains(owners.Groups, rqst.Owner) {
			owners.Groups = remove(owners.Groups, rqst.Owner)
		}
	} else if rqst.Subject == resourcepb.SubjectType_ORGANIZATION {
		if Utility.Contains(owners.Organizations, rqst.Owner) {
			owners.Organizations = remove(owners.Organizations, rqst.Owner)
		}
	} else if rqst.Subject == resourcepb.SubjectType_PEER {
		if Utility.Contains(owners.Peers, rqst.Owner) {
			owners.Peers = remove(owners.Peers, rqst.Owner)
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

//* That function must be call when a subject is removed to clean up permissions.
func (self *Globule) DeleteAllAccess(ctx context.Context, rqst *resourcepb.DeleteAllAccessRqst) (*resourcepb.DeleteAllAccessRsp, error) {
	err := self.permissions.RemoveItem(rqst.Subject + "*")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &resourcepb.DeleteAllAccessRsp{}, nil
}

//* Validate if a user can get access to a given ressource for a given operation (read, write...) *
func (self *Globule) ValidateAccess(ctx context.Context, rqst *resourcepb.ValidateAccessRqst) (*resourcepb.ValidateAccessRsp, error) {
	permissions, err := self.getResourcePermissions(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	owners := permissions.Owners
	hasAccess := false
	subjectType := ""
	if rqst.Type == resourcepb.SubjectType_ACCOUNT {
		subjectType = "Account"
		if Utility.Contains(owners.Accounts, rqst.Subject) {
			hasAccess = true
		}
	} else if rqst.Type == resourcepb.SubjectType_APPLICATION {
		subjectType = "Application"
		if Utility.Contains(owners.Applications, rqst.Subject) {
			hasAccess = true
		}
	} else if rqst.Type == resourcepb.SubjectType_GROUP {
		subjectType = "Group"
		if Utility.Contains(owners.Groups, rqst.Subject) {
			hasAccess = true
		}
	} else if rqst.Type == resourcepb.SubjectType_ORGANIZATION {
		subjectType = "Organization"
		if Utility.Contains(owners.Organizations, rqst.Subject) {
			hasAccess = true
		}
	} else if rqst.Type == resourcepb.SubjectType_PEER {
		subjectType = "Peer"
		if Utility.Contains(owners.Peers, rqst.Subject) {
			hasAccess = true
		}
	}

	if !hasAccess {
		err := errors.New("Permission denied for " + subjectType + " " + rqst.Subject + "!")
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))

	}

	// The permission is set.
	return &resourcepb.ValidateAccessRsp{}, nil
}
