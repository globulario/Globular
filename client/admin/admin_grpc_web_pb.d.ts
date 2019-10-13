import * as grpcWeb from 'grpc-web';

import {
  GetConfigRequest,
  GetConfigResponse,
  RegisterExternalServiceRequest,
  RegisterExternalServiceResponse,
  SaveConfigRequest,
  SaveConfigResponse,
  StartServiceRequest,
  StartServiceResponse,
  StopServiceRequest,
  StopServiceResponse} from './admin_pb';

export class AdminServiceClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

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

  registerExternalService(
    request: RegisterExternalServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RegisterExternalServiceResponse) => void
  ): grpcWeb.ClientReadableStream<RegisterExternalServiceResponse>;

}

export class AdminServicePromiseClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

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

  registerExternalService(
    request: RegisterExternalServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<RegisterExternalServiceResponse>;

}

