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
  CreateRoleRqst,
  CreateRoleRsp,
  DeleteAccountRqst,
  DeleteAccountRsp,
  DeletePermissionsRqst,
  DeletePermissionsRsp,
  DeleteRessourceOwnerRqst,
  DeleteRessourceOwnerRsp,
  DeleteRessourceOwnersRqst,
  DeleteRessourceOwnersRsp,
  DeleteRoleRqst,
  DeleteRoleRsp,
  GetAllActionsRqst,
  GetAllActionsRsp,
  GetAllApplicationsInfoRqst,
  GetAllApplicationsInfoRsp,
  GetAllFilesInfoRqst,
  GetAllFilesInfoRsp,
  GetPermissionsRqst,
  GetPermissionsRsp,
  GetRessourceOwnersRqst,
  GetRessourceOwnersRsp,
  RefreshTokenRqst,
  RefreshTokenRsp,
  RegisterAccountRqst,
  RegisterAccountRsp,
  RemoveAccountRoleRqst,
  RemoveAccountRoleRsp,
  RemoveApplicationActionRqst,
  RemoveApplicationActionRsp,
  RemoveApplicationRqst,
  RemoveApplicationRsp,
  RemoveRoleActionRqst,
  RemoveRoleActionRsp,
  SetPermissionRqst,
  SetPermissionRsp,
  SetRessourceOwnerRqst,
  SetRessourceOwnerRsp,
  SynchronizeLdapRqst,
  SynchronizeLdapRsp} from './ressource_pb';

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

  getAllApplicationsInfo(
    request: GetAllApplicationsInfoRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetAllApplicationsInfoRsp) => void
  ): grpcWeb.ClientReadableStream<GetAllApplicationsInfoRsp>;

  removeApplication(
    request: RemoveApplicationRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RemoveApplicationRsp) => void
  ): grpcWeb.ClientReadableStream<RemoveApplicationRsp>;

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

  getAllApplicationsInfo(
    request: GetAllApplicationsInfoRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<GetAllApplicationsInfoRsp>;

  removeApplication(
    request: RemoveApplicationRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<RemoveApplicationRsp>;

}

