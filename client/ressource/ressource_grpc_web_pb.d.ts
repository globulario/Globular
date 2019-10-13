import * as grpcWeb from 'grpc-web';

import {
  AuthenticateRqst,
  AuthenticateRsp,
  DeleteAccountRqst,
  DeleteAccountRsp,
  RefreshTokenRqst,
  RefreshTokenRsp,
  RegisterAccountRqst,
  RegisterAccountRsp} from './ressource_pb';

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

}

