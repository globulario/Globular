import * as grpcWeb from 'grpc-web';

import {
  AddAccountRoleRqst,
  AddAccountRoleRsp,
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
  DeleteRoleRqst,
  DeleteRoleRsp,
  GetAllActionsRqst,
  GetAllActionsRsp,
  GetAllFilesInfoRqst,
  GetAllFilesInfoRsp,
  GetPermissionsRqst,
  GetPermissionsRsp,
  RefreshTokenRqst,
  RefreshTokenRsp,
  RegisterAccountRqst,
  RegisterAccountRsp,
  RemoveAccountRoleRqst,
  RemoveAccountRoleRsp,
  RemoveRoleActionRqst,
  RemoveRoleActionRsp,
  SetPermissionRqst,
  SetPermissionRsp} from './ressource_pb';

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

  getAllFilesInfo(
    request: GetAllFilesInfoRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetAllFilesInfoRsp) => void
  ): grpcWeb.ClientReadableStream<GetAllFilesInfoRsp>;

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

  getAllFilesInfo(
    request: GetAllFilesInfoRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<GetAllFilesInfoRsp>;

}

