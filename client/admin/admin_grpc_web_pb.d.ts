import * as grpcWeb from 'grpc-web';

import {
  DeployApplicationRequest,
  DeployApplicationResponse,
  GetConfigRequest,
  GetConfigResponse,
  InstallServiceRequest,
  InstallServiceResponse,
  PublishServiceRequest,
  PublishServiceResponse,
  RegisterExternalApplicationRequest,
  RegisterExternalApplicationResponse,
  SaveConfigRequest,
  SaveConfigResponse,
  SetEmailRequest,
  SetEmailResponse,
  SetPasswordRequest,
  SetPasswordResponse,
  SetRootEmailRequest,
  SetRootEmailResponse,
  SetRootPasswordRequest,
  SetRootPasswordResponse,
  StartServiceRequest,
  StartServiceResponse,
  StopServiceRequest,
  StopServiceResponse,
  UninstallServiceRequest,
  UninstallServiceResponse,
  UploadServicePackageRequest,
  UploadServicePackageResponse} from './admin_pb';

export class AdminServiceClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  setRootPassword(
    request: SetRootPasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SetRootPasswordResponse) => void
  ): grpcWeb.ClientReadableStream<SetRootPasswordResponse>;

  setRootEmail(
    request: SetRootEmailRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SetRootEmailResponse) => void
  ): grpcWeb.ClientReadableStream<SetRootEmailResponse>;

  setPassword(
    request: SetPasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SetPasswordResponse) => void
  ): grpcWeb.ClientReadableStream<SetPasswordResponse>;

  setEmail(
    request: SetEmailRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SetEmailResponse) => void
  ): grpcWeb.ClientReadableStream<SetEmailResponse>;

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

  publishService(
    request: PublishServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PublishServiceResponse) => void
  ): grpcWeb.ClientReadableStream<PublishServiceResponse>;

  installService(
    request: InstallServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: InstallServiceResponse) => void
  ): grpcWeb.ClientReadableStream<InstallServiceResponse>;

  uninstallService(
    request: UninstallServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UninstallServiceResponse) => void
  ): grpcWeb.ClientReadableStream<UninstallServiceResponse>;

  registerExternalApplication(
    request: RegisterExternalApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RegisterExternalApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<RegisterExternalApplicationResponse>;

}

export class AdminServicePromiseClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  setRootPassword(
    request: SetRootPasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SetRootPasswordResponse>;

  setRootEmail(
    request: SetRootEmailRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SetRootEmailResponse>;

  setPassword(
    request: SetPasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SetPasswordResponse>;

  setEmail(
    request: SetEmailRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SetEmailResponse>;

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

  publishService(
    request: PublishServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<PublishServiceResponse>;

  installService(
    request: InstallServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<InstallServiceResponse>;

  uninstallService(
    request: UninstallServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UninstallServiceResponse>;

  registerExternalApplication(
    request: RegisterExternalApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<RegisterExternalApplicationResponse>;

}

