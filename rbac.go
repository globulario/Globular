package main

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"

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

//* Set the action resources permissions *
func (self *Globule) SetActionResourcesPermissions(ctx context.Context, rqst *resourcepb.SetActionResourcesPermissionsRqst) (*resourcepb.SetActionResourcesPermissionsRsp, error) {
	// TODO implement it.
	return nil, errors.New("method SetActionResourcesPermissions not yet implemented!")
}

//* Get the action ressouces permission *
func (self *Globule) GetActionResourcesPermissions(ctx context.Context, rqst *resourcepb.GetActionResourcesPermissionsRqst) (*resourcepb.GetActionResourcesPermissionsRsp, error) {
	// TODO implement it.
	return nil, errors.New("method GetActionResourcesPermissions not yet implemented!")
}

//* Set resource permissions this method will replace existing permission at once *
func (self *Globule) SetResourcePermissions(ctx context.Context, rqst *resourcepb.SetResourcePermissionsRqst) (*resourcepb.SetResourcePermissionsRqst, error) {
	// TODO implement it.
	return nil, errors.New("method SetResourcePermissions not yet implemented!")
}

//* Delete a resource permissions (when a resource is deleted) *
func (self *Globule) DeleteResourcePermissions(ctx context.Context, rqst *resourcepb.DeleteResourcePermissionsRqst) (*resourcepb.DeleteResourcePermissionsRqst, error) {
	// TODO implement it.
	return nil, errors.New("method DeleteResourcePermissions not yet implemented!")
}

//* Delete a specific resource permission *
func (self *Globule) DeleteResourcePermission(ctx context.Context, rqst *resourcepb.DeleteResourcePermissionRqst) (*resourcepb.DeleteResourcePermissionRqst, error) {
	// TODO implement it.
	return nil, errors.New("method DeleteResourcePermission not yet implemented!")
}

//* Set specific resource permission  ex. read permission... *
func (self *Globule) SetResourcePermission(ctx context.Context, rqst *resourcepb.SetResourcePermissionRqst) (*resourcepb.SetResourcePermissionRsp, error) {
	// TODO implement it.
	return nil, errors.New("method SetResourcePermission not yet implemented!")
}

//* Get a specific resource access *
func (self *Globule) GetResourcePermission(ctx context.Context, rqst *resourcepb.GetResourcePermissionRqst) (*resourcepb.GetResourcePermissionRsp, error) {
	// TODO implement it.
	return nil, errors.New("method GetResourcePermission not yet implemented!")
}

//* Get resource permissions *
func (self *Globule) GetResourcePermissions(ctx context.Context, rqst *resourcepb.GetResourcePermissionsRqst) (*resourcepb.GetResourcePermissionsRsp, error) {
	// TODO implement it.
	return nil, errors.New("method GetResourcePermissions not yet implemented!")
}

//* Add resource owner do nothing if it already exist
func (self *Globule) AddResourceOwner(ctx context.Context, rqst *resourcepb.AddResourceOwnerRqst) (*resourcepb.AddResourceOwnerRsp, error) {
	// TODO implement it.
	return nil, nil
}

//* Remove resource owner
func (self *Globule) RemoveResourceOwner(ctx context.Context, rqst *resourcepb.AddResourceOwnerRqst) (*resourcepb.AddResourceOwnerRsp, error) {
	// TODO implement it.
	return nil, nil
}

//* That function must be call when a subject is removed to clean up permissions.
func (self *Globule) DeleteAllAccess(ctx context.Context, rqst *resourcepb.DeleteAllAccessRqst) (*resourcepb.DeleteAllAccessRsp, error) {
	// TODO implement it.
	return nil, nil
}

//* Validate if a user can get access to a given ressource for a given operation (read, write...) *
func (self *Globule) ValidateAccess(ctx context.Context, rqst *resourcepb.ValidateAccessRqst) (*resourcepb.ValidateAccessRsp, error) {
	// TODO implement it.
	return nil, nil
}

//* Return the list of access for a given subject
func (self *Globule) GetAccesses(ctx context.Context, rqst *resourcepb.GetAccessesRqst) (*resourcepb.GetAccessesRsp, error) {
	// TODO implement it.
	return nil, nil
}
