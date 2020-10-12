import * as grpcWeb from 'grpc-web';

import * as ressource_pb from './ressource_pb';


export class RessourceServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  registerAccount(
    request: ressource_pb.RegisterAccountRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.RegisterAccountRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.RegisterAccountRsp>;

  deleteAccount(
    request: ressource_pb.DeleteAccountRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteAccountRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteAccountRsp>;

  authenticate(
    request: ressource_pb.AuthenticateRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.AuthenticateRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.AuthenticateRsp>;

  synchronizeLdap(
    request: ressource_pb.SynchronizeLdapRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.SynchronizeLdapRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.SynchronizeLdapRsp>;

  refreshToken(
    request: ressource_pb.RefreshTokenRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.RefreshTokenRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.RefreshTokenRsp>;

  addAccountRole(
    request: ressource_pb.AddAccountRoleRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.AddAccountRoleRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.AddAccountRoleRsp>;

  removeAccountRole(
    request: ressource_pb.RemoveAccountRoleRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.RemoveAccountRoleRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.RemoveAccountRoleRsp>;

  createRole(
    request: ressource_pb.CreateRoleRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.CreateRoleRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.CreateRoleRsp>;

  deleteRole(
    request: ressource_pb.DeleteRoleRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteRoleRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteRoleRsp>;

  addRoleAction(
    request: ressource_pb.AddRoleActionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.AddRoleActionRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.AddRoleActionRsp>;

  removeRoleAction(
    request: ressource_pb.RemoveRoleActionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.RemoveRoleActionRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.RemoveRoleActionRsp>;

  addApplicationAction(
    request: ressource_pb.AddApplicationActionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.AddApplicationActionRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.AddApplicationActionRsp>;

  removeApplicationAction(
    request: ressource_pb.RemoveApplicationActionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.RemoveApplicationActionRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.RemoveApplicationActionRsp>;

  getAllActions(
    request: ressource_pb.GetAllActionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.GetAllActionsRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.GetAllActionsRsp>;

  getPermissions(
    request: ressource_pb.GetPermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.GetPermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.GetPermissionsRsp>;

  setPermission(
    request: ressource_pb.SetPermissionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.SetPermissionRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.SetPermissionRsp>;

  deletePermissions(
    request: ressource_pb.DeletePermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeletePermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeletePermissionsRsp>;

  setRessourceOwner(
    request: ressource_pb.SetRessourceOwnerRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.SetRessourceOwnerRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.SetRessourceOwnerRsp>;

  getRessourceOwners(
    request: ressource_pb.GetRessourceOwnersRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.GetRessourceOwnersRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.GetRessourceOwnersRsp>;

  deleteRessourceOwner(
    request: ressource_pb.DeleteRessourceOwnerRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteRessourceOwnerRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteRessourceOwnerRsp>;

  deleteRessourceOwners(
    request: ressource_pb.DeleteRessourceOwnersRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteRessourceOwnersRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteRessourceOwnersRsp>;

  getAllFilesInfo(
    request: ressource_pb.GetAllFilesInfoRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.GetAllFilesInfoRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.GetAllFilesInfoRsp>;

  validateToken(
    request: ressource_pb.ValidateTokenRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.ValidateTokenRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.ValidateTokenRsp>;

  validateUserRessourceAccess(
    request: ressource_pb.ValidateUserRessourceAccessRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.ValidateUserRessourceAccessRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.ValidateUserRessourceAccessRsp>;

  validateApplicationRessourceAccess(
    request: ressource_pb.ValidateApplicationRessourceAccessRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.ValidateApplicationRessourceAccessRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.ValidateApplicationRessourceAccessRsp>;

  validateUserAccess(
    request: ressource_pb.ValidateUserAccessRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.ValidateUserAccessRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.ValidateUserAccessRsp>;

  validateApplicationAccess(
    request: ressource_pb.ValidateApplicationAccessRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.ValidateApplicationAccessRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.ValidateApplicationAccessRsp>;

  createDirPermissions(
    request: ressource_pb.CreateDirPermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.CreateDirPermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.CreateDirPermissionsRsp>;

  renameFilePermission(
    request: ressource_pb.RenameFilePermissionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.RenameFilePermissionRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.RenameFilePermissionRsp>;

  deleteDirPermissions(
    request: ressource_pb.DeleteDirPermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteDirPermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteDirPermissionsRsp>;

  deleteFilePermissions(
    request: ressource_pb.DeleteFilePermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteFilePermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteFilePermissionsRsp>;

  deleteAccountPermissions(
    request: ressource_pb.DeleteAccountPermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteAccountPermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteAccountPermissionsRsp>;

  deleteRolePermissions(
    request: ressource_pb.DeleteRolePermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteRolePermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteRolePermissionsRsp>;

  getAllApplicationsInfo(
    request: ressource_pb.GetAllApplicationsInfoRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.GetAllApplicationsInfoRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.GetAllApplicationsInfoRsp>;

  deleteApplication(
    request: ressource_pb.DeleteApplicationRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteApplicationRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteApplicationRsp>;

  log(
    request: ressource_pb.LogRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.LogRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.LogRsp>;

  getLog(
    request: ressource_pb.GetLogRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<ressource_pb.GetLogRsp>;

  deleteLog(
    request: ressource_pb.DeleteLogRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.DeleteLogRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.DeleteLogRsp>;

  clearAllLog(
    request: ressource_pb.ClearAllLogRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.ClearAllLogRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.ClearAllLogRsp>;

  getRessources(
    request: ressource_pb.GetRessourcesRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<ressource_pb.GetRessourcesRsp>;

  setRessource(
    request: ressource_pb.SetRessourceRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.SetRessourceRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.SetRessourceRsp>;

  removeRessource(
    request: ressource_pb.RemoveRessourceRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.RemoveRessourceRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.RemoveRessourceRsp>;

  setActionPermission(
    request: ressource_pb.SetActionPermissionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.SetActionPermissionRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.SetActionPermissionRsp>;

  removeActionPermission(
    request: ressource_pb.RemoveActionPermissionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.RemoveActionPermissionRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.RemoveActionPermissionRsp>;

  getActionPermission(
    request: ressource_pb.GetActionPermissionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ressource_pb.GetActionPermissionRsp) => void
  ): grpcWeb.ClientReadableStream<ressource_pb.GetActionPermissionRsp>;

}

export class RessourceServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  registerAccount(
    request: ressource_pb.RegisterAccountRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.RegisterAccountRsp>;

  deleteAccount(
    request: ressource_pb.DeleteAccountRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteAccountRsp>;

  authenticate(
    request: ressource_pb.AuthenticateRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.AuthenticateRsp>;

  synchronizeLdap(
    request: ressource_pb.SynchronizeLdapRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.SynchronizeLdapRsp>;

  refreshToken(
    request: ressource_pb.RefreshTokenRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.RefreshTokenRsp>;

  addAccountRole(
    request: ressource_pb.AddAccountRoleRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.AddAccountRoleRsp>;

  removeAccountRole(
    request: ressource_pb.RemoveAccountRoleRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.RemoveAccountRoleRsp>;

  createRole(
    request: ressource_pb.CreateRoleRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.CreateRoleRsp>;

  deleteRole(
    request: ressource_pb.DeleteRoleRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteRoleRsp>;

  addRoleAction(
    request: ressource_pb.AddRoleActionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.AddRoleActionRsp>;

  removeRoleAction(
    request: ressource_pb.RemoveRoleActionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.RemoveRoleActionRsp>;

  addApplicationAction(
    request: ressource_pb.AddApplicationActionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.AddApplicationActionRsp>;

  removeApplicationAction(
    request: ressource_pb.RemoveApplicationActionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.RemoveApplicationActionRsp>;

  getAllActions(
    request: ressource_pb.GetAllActionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.GetAllActionsRsp>;

  getPermissions(
    request: ressource_pb.GetPermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.GetPermissionsRsp>;

  setPermission(
    request: ressource_pb.SetPermissionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.SetPermissionRsp>;

  deletePermissions(
    request: ressource_pb.DeletePermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeletePermissionsRsp>;

  setRessourceOwner(
    request: ressource_pb.SetRessourceOwnerRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.SetRessourceOwnerRsp>;

  getRessourceOwners(
    request: ressource_pb.GetRessourceOwnersRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.GetRessourceOwnersRsp>;

  deleteRessourceOwner(
    request: ressource_pb.DeleteRessourceOwnerRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteRessourceOwnerRsp>;

  deleteRessourceOwners(
    request: ressource_pb.DeleteRessourceOwnersRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteRessourceOwnersRsp>;

  getAllFilesInfo(
    request: ressource_pb.GetAllFilesInfoRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.GetAllFilesInfoRsp>;

  validateToken(
    request: ressource_pb.ValidateTokenRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.ValidateTokenRsp>;

  validateUserRessourceAccess(
    request: ressource_pb.ValidateUserRessourceAccessRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.ValidateUserRessourceAccessRsp>;

  validateApplicationRessourceAccess(
    request: ressource_pb.ValidateApplicationRessourceAccessRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.ValidateApplicationRessourceAccessRsp>;

  validateUserAccess(
    request: ressource_pb.ValidateUserAccessRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.ValidateUserAccessRsp>;

  validateApplicationAccess(
    request: ressource_pb.ValidateApplicationAccessRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.ValidateApplicationAccessRsp>;

  createDirPermissions(
    request: ressource_pb.CreateDirPermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.CreateDirPermissionsRsp>;

  renameFilePermission(
    request: ressource_pb.RenameFilePermissionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.RenameFilePermissionRsp>;

  deleteDirPermissions(
    request: ressource_pb.DeleteDirPermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteDirPermissionsRsp>;

  deleteFilePermissions(
    request: ressource_pb.DeleteFilePermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteFilePermissionsRsp>;

  deleteAccountPermissions(
    request: ressource_pb.DeleteAccountPermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteAccountPermissionsRsp>;

  deleteRolePermissions(
    request: ressource_pb.DeleteRolePermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteRolePermissionsRsp>;

  getAllApplicationsInfo(
    request: ressource_pb.GetAllApplicationsInfoRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.GetAllApplicationsInfoRsp>;

  deleteApplication(
    request: ressource_pb.DeleteApplicationRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteApplicationRsp>;

  log(
    request: ressource_pb.LogRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.LogRsp>;

  getLog(
    request: ressource_pb.GetLogRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<ressource_pb.GetLogRsp>;

  deleteLog(
    request: ressource_pb.DeleteLogRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.DeleteLogRsp>;

  clearAllLog(
    request: ressource_pb.ClearAllLogRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.ClearAllLogRsp>;

  getRessources(
    request: ressource_pb.GetRessourcesRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<ressource_pb.GetRessourcesRsp>;

  setRessource(
    request: ressource_pb.SetRessourceRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.SetRessourceRsp>;

  removeRessource(
    request: ressource_pb.RemoveRessourceRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.RemoveRessourceRsp>;

  setActionPermission(
    request: ressource_pb.SetActionPermissionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.SetActionPermissionRsp>;

  removeActionPermission(
    request: ressource_pb.RemoveActionPermissionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.RemoveActionPermissionRsp>;

  getActionPermission(
    request: ressource_pb.GetActionPermissionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ressource_pb.GetActionPermissionRsp>;

}

