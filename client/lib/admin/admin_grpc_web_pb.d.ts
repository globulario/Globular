import * as grpcWeb from 'grpc-web';

import {
  GetConfigRequest,
  GetConfigResponse,
  RegisterExternalApplicationRequest,
  RegisterExternalApplicationResponse,
  SaveConfigRequest,
  SaveConfigResponse,
  SetRootPasswordRqst,
  SetRootPasswordRsp,
  StartServiceRequest,
  StartServiceResponse,
  StopServiceRequest,
  StopServiceResponse} from './admin_pb';

export class AdminServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  setRootPassword(
    request: SetRootPasswordRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SetRootPasswordRsp) => void
  ): grpcWeb.ClientReadableStream<SetRootPasswordRsp>;

  getConfig(
    request: GetConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetConfigResponse) => void
  ): grpcWeb.ClientReadableStream<GetConfigResponse>;

  getFullConfig(
    request: GetConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetConfigResponse) => void
  ): grpcWeb.ClientReadableStream<GetConfigResponse>;

  saveConfig(
    request: SaveConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveConfigResponse) => void
  ): grpcWeb.ClientReadableStream<SaveConfigResponse>;

  stopService(
    request: StopServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StopServiceResponse) => void
  ): grpcWeb.ClientReadableStream<StopServiceResponse>;

  startService(
    request: StartServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StartServiceResponse) => void
  ): grpcWeb.ClientReadableStream<StartServiceResponse>;

  registerExternalApplication(
    request: RegisterExternalApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RegisterExternalApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<RegisterExternalApplicationResponse>;

}

export class AdminServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  setRootPassword(
    request: SetRootPasswordRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<SetRootPasswordRsp>;

  getConfig(
    request: GetConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetConfigResponse>;

  getFullConfig(
    request: GetConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetConfigResponse>;

  saveConfig(
    request: SaveConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveConfigResponse>;

  stopService(
    request: StopServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StopServiceResponse>;

  startService(
    request: StartServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StartServiceResponse>;

  registerExternalApplication(
    request: RegisterExternalApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<RegisterExternalApplicationResponse>;

}

