import * as grpcWeb from 'grpc-web';

import {
  AddAccountRoleRqst,
  AddAccountRoleRsp,
  AddApplicationActionRqst,
  AddApplicationActionRsp,
  AddRoleActionRqst,
  AddRoleActionRsp,
  AuthenticateRqst,
  AuthenticateRsp,
  ClearAllLogRqst,
  ClearAllLogRsp,
  CreateDirPermissionsRqst,
  CreateDirPermissionsRsp,
  CreateRoleRqst,
  CreateRoleRsp,
  DeleteAccountPermissionsRqst,
  DeleteAccountPermissionsRsp,
  DeleteAccountRqst,
  DeleteAccountRsp,
  DeleteApplicationRqst,
  DeleteApplicationRsp,
  DeleteDirPermissionsRqst,
  DeleteDirPermissionsRsp,
  DeleteFilePermissionsRqst,
  DeleteFilePermissionsRsp,
  DeleteLogRqst,
  DeleteLogRsp,
  DeletePermissionsRqst,
  DeletePermissionsRsp,
  DeleteRessourceOwnerRqst,
  DeleteRessourceOwnerRsp,
  DeleteRessourceOwnersRqst,
  DeleteRessourceOwnersRsp,
  DeleteRolePermissionsRqst,
  DeleteRolePermissionsRsp,
  DeleteRoleRqst,
  DeleteRoleRsp,
  GetAllActionsRqst,
  GetAllActionsRsp,
  GetAllApplicationsInfoRqst,
  GetAllApplicationsInfoRsp,
  GetAllFilesInfoRqst,
  GetAllFilesInfoRsp,
  GetLogRqst,
  GetLogRsp,
  GetPermissionsRqst,
  GetPermissionsRsp,
  GetRessourceOwnersRqst,
  GetRessourceOwnersRsp,
  LogRqst,
  LogRsp,
  RefreshTokenRqst,
  RefreshTokenRsp,
  RegisterAccountRqst,
  RegisterAccountRsp,
  RemoveAccountRoleRqst,
  RemoveAccountRoleRsp,
  RemoveApplicationActionRqst,
  RemoveApplicationActionRsp,
  RemoveRoleActionRqst,
  RemoveRoleActionRsp,
  RenameFilePermissionRqst,
  RenameFilePermissionRsp,
  ResetLogRqst,
  ResetLogRsp,
  SetLogRqst,
  SetLogRsp,
  SetPermissionRqst,
  SetPermissionRsp,
  SetRessourceOwnerRqst,
  SetRessourceOwnerRsp,
  SynchronizeLdapRqst,
  SynchronizeLdapRsp,
  ValidateApplicationAccessRqst,
  ValidateApplicationAccessRsp,
  ValidateApplicationFileAccessRqst,
  ValidateApplicationFileAccessRsp,
  ValidateUserAccessRqst,
  ValidateUserAccessRsp,
  ValidateUserFileAccessRqst,
  ValidateUserFileAccessRsp} from './ressource_pb';

export class RessourceServiceClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  registerAccount(
    request: RegisterAccountRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RegisterAccountRsp) => void
  ): grpcWeb.ClientReadableStream<RegisterAccountRsp>;

  deleteAccount(
    request: DeleteAccountRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteAccountRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteAccountRsp>;

  authenticate(
    request: AuthenticateRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AuthenticateRsp) => void
  ): grpcWeb.ClientReadableStream<AuthenticateRsp>;

  synchronizeLdap(
    request: SynchronizeLdapRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SynchronizeLdapRsp) => void
  ): grpcWeb.ClientReadableStream<SynchronizeLdapRsp>;

  refreshToken(
    request: RefreshTokenRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RefreshTokenRsp) => void
  ): grpcWeb.ClientReadableStream<RefreshTokenRsp>;

  addAccountRole(
    request: AddAccountRoleRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AddAccountRoleRsp) => void
  ): grpcWeb.ClientReadableStream<AddAccountRoleRsp>;

  removeAccountRole(
    request: RemoveAccountRoleRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RemoveAccountRoleRsp) => void
  ): grpcWeb.ClientReadableStream<RemoveAccountRoleRsp>;

  createRole(
    request: CreateRoleRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateRoleRsp) => void
  ): grpcWeb.ClientReadableStream<CreateRoleRsp>;

  deleteRole(
    request: DeleteRoleRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteRoleRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteRoleRsp>;

  addRoleAction(
    request: AddRoleActionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AddRoleActionRsp) => void
  ): grpcWeb.ClientReadableStream<AddRoleActionRsp>;

  removeRoleAction(
    request: RemoveRoleActionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RemoveRoleActionRsp) => void
  ): grpcWeb.ClientReadableStream<RemoveRoleActionRsp>;

  addApplicationAction(
    request: AddApplicationActionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AddApplicationActionRsp) => void
  ): grpcWeb.ClientReadableStream<AddApplicationActionRsp>;

  removeApplicationAction(
    request: RemoveApplicationActionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RemoveApplicationActionRsp) => void
  ): grpcWeb.ClientReadableStream<RemoveApplicationActionRsp>;

  getAllActions(
    request: GetAllActionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetAllActionsRsp) => void
  ): grpcWeb.ClientReadableStream<GetAllActionsRsp>;

  getPermissions(
    request: GetPermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetPermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<GetPermissionsRsp>;

  setPermission(
    request: SetPermissionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SetPermissionRsp) => void
  ): grpcWeb.ClientReadableStream<SetPermissionRsp>;

  deletePermissions(
    request: DeletePermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeletePermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<DeletePermissionsRsp>;

  setRessourceOwner(
    request: SetRessourceOwnerRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SetRessourceOwnerRsp) => void
  ): grpcWeb.ClientReadableStream<SetRessourceOwnerRsp>;

  getRessourceOwners(
    request: GetRessourceOwnersRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetRessourceOwnersRsp) => void
  ): grpcWeb.ClientReadableStream<GetRessourceOwnersRsp>;

  deleteRessourceOwner(
    request: DeleteRessourceOwnerRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteRessourceOwnerRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteRessourceOwnerRsp>;

  deleteRessourceOwners(
    request: DeleteRessourceOwnersRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteRessourceOwnersRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteRessourceOwnersRsp>;

  getAllFilesInfo(
    request: GetAllFilesInfoRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetAllFilesInfoRsp) => void
  ): grpcWeb.ClientReadableStream<GetAllFilesInfoRsp>;

  validateUserFileAccess(
    request: ValidateUserFileAccessRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ValidateUserFileAccessRsp) => void
  ): grpcWeb.ClientReadableStream<ValidateUserFileAccessRsp>;

  validateApplicationFileAccess(
    request: ValidateApplicationFileAccessRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ValidateApplicationFileAccessRsp) => void
  ): grpcWeb.ClientReadableStream<ValidateApplicationFileAccessRsp>;

  validateUserAccess(
    request: ValidateUserAccessRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ValidateUserAccessRsp) => void
  ): grpcWeb.ClientReadableStream<ValidateUserAccessRsp>;

  validateApplicationAccess(
    request: ValidateApplicationAccessRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ValidateApplicationAccessRsp) => void
  ): grpcWeb.ClientReadableStream<ValidateApplicationAccessRsp>;

  createDirPermissions(
    request: CreateDirPermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateDirPermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<CreateDirPermissionsRsp>;

  renameFilePermission(
    request: RenameFilePermissionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RenameFilePermissionRsp) => void
  ): grpcWeb.ClientReadableStream<RenameFilePermissionRsp>;

  deleteDirPermissions(
    request: DeleteDirPermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteDirPermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteDirPermissionsRsp>;

  deleteFilePermissions(
    request: DeleteFilePermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteFilePermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteFilePermissionsRsp>;

  deleteAccountPermissions(
    request: DeleteAccountPermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteAccountPermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteAccountPermissionsRsp>;

  deleteRolePermissions(
    request: DeleteRolePermissionsRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteRolePermissionsRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteRolePermissionsRsp>;

  getAllApplicationsInfo(
    request: GetAllApplicationsInfoRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetAllApplicationsInfoRsp) => void
  ): grpcWeb.ClientReadableStream<GetAllApplicationsInfoRsp>;

  deleteApplication(
    request: DeleteApplicationRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteApplicationRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteApplicationRsp>;

  log(
    request: LogRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: LogRsp) => void
  ): grpcWeb.ClientReadableStream<LogRsp>;

  setLog(
    request: SetLogRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SetLogRsp) => void
  ): grpcWeb.ClientReadableStream<SetLogRsp>;

  resetLog(
    request: ResetLogRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ResetLogRsp) => void
  ): grpcWeb.ClientReadableStream<ResetLogRsp>;

  getLog(
    request: GetLogRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<GetLogRsp>;

  deleteLog(
    request: DeleteLogRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteLogRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteLogRsp>;

  clearAllLog(
    request: ClearAllLogRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ClearAllLogRsp) => void
  ): grpcWeb.ClientReadableStream<ClearAllLogRsp>;

}

export class RessourceServicePromiseClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  registerAccount(
    request: RegisterAccountRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<RegisterAccountRsp>;

  deleteAccount(
    request: DeleteAccountRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteAccountRsp>;

  authenticate(
    request: AuthenticateRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<AuthenticateRsp>;

  synchronizeLdap(
    request: SynchronizeLdapRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<SynchronizeLdapRsp>;

  refreshToken(
    request: RefreshTokenRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<RefreshTokenRsp>;

  addAccountRole(
    request: AddAccountRoleRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<AddAccountRoleRsp>;

  removeAccountRole(
    request: RemoveAccountRoleRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<RemoveAccountRoleRsp>;

  createRole(
    request: CreateRoleRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateRoleRsp>;

  deleteRole(
    request: DeleteRoleRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteRoleRsp>;

  addRoleAction(
    request: AddRoleActionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<AddRoleActionRsp>;

  removeRoleAction(
    request: RemoveRoleActionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<RemoveRoleActionRsp>;

  addApplicationAction(
    request: AddApplicationActionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<AddApplicationActionRsp>;

  removeApplicationAction(
    request: RemoveApplicationActionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<RemoveApplicationActionRsp>;

  getAllActions(
    request: GetAllActionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<GetAllActionsRsp>;

  getPermissions(
    request: GetPermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<GetPermissionsRsp>;

  setPermission(
    request: SetPermissionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<SetPermissionRsp>;

  deletePermissions(
    request: DeletePermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeletePermissionsRsp>;

  setRessourceOwner(
    request: SetRessourceOwnerRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<SetRessourceOwnerRsp>;

  getRessourceOwners(
    request: GetRessourceOwnersRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<GetRessourceOwnersRsp>;

  deleteRessourceOwner(
    request: DeleteRessourceOwnerRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteRessourceOwnerRsp>;

  deleteRessourceOwners(
    request: DeleteRessourceOwnersRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteRessourceOwnersRsp>;

  getAllFilesInfo(
    request: GetAllFilesInfoRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<GetAllFilesInfoRsp>;

  validateUserFileAccess(
    request: ValidateUserFileAccessRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ValidateUserFileAccessRsp>;

  validateApplicationFileAccess(
    request: ValidateApplicationFileAccessRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ValidateApplicationFileAccessRsp>;

  validateUserAccess(
    request: ValidateUserAccessRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ValidateUserAccessRsp>;

  validateApplicationAccess(
    request: ValidateApplicationAccessRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ValidateApplicationAccessRsp>;

  createDirPermissions(
    request: CreateDirPermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateDirPermissionsRsp>;

  renameFilePermission(
    request: RenameFilePermissionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<RenameFilePermissionRsp>;

  deleteDirPermissions(
    request: DeleteDirPermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteDirPermissionsRsp>;

  deleteFilePermissions(
    request: DeleteFilePermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteFilePermissionsRsp>;

  deleteAccountPermissions(
    request: DeleteAccountPermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteAccountPermissionsRsp>;

  deleteRolePermissions(
    request: DeleteRolePermissionsRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteRolePermissionsRsp>;

  getAllApplicationsInfo(
    request: GetAllApplicationsInfoRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<GetAllApplicationsInfoRsp>;

  deleteApplication(
    request: DeleteApplicationRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteApplicationRsp>;

  log(
    request: LogRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<LogRsp>;

  setLog(
    request: SetLogRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<SetLogRsp>;

  resetLog(
    request: ResetLogRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ResetLogRsp>;

  getLog(
    request: GetLogRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<GetLogRsp>;

  deleteLog(
    request: DeleteLogRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteLogRsp>;

  clearAllLog(
    request: ClearAllLogRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ClearAllLogRsp>;

}

