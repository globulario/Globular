import * as grpcWeb from 'grpc-web';

import * as admin_pb from './admin_pb';


export class AdminServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  setRootPassword(
    request: admin_pb.SetRootPasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.SetRootPasswordResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.SetRootPasswordResponse>;

  setRootEmail(
    request: admin_pb.SetRootEmailRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.SetRootEmailResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.SetRootEmailResponse>;

  setPassword(
    request: admin_pb.SetPasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.SetPasswordResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.SetPasswordResponse>;

  setEmail(
    request: admin_pb.SetEmailRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.SetEmailResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.SetEmailResponse>;

  getConfig(
    request: admin_pb.GetConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.GetConfigResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.GetConfigResponse>;

  getFullConfig(
    request: admin_pb.GetConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.GetConfigResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.GetConfigResponse>;

  saveConfig(
    request: admin_pb.SaveConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.SaveConfigResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.SaveConfigResponse>;

  stopService(
    request: admin_pb.StopServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.StopServiceResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.StopServiceResponse>;

  startService(
    request: admin_pb.StartServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.StartServiceResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.StartServiceResponse>;

  publishService(
    request: admin_pb.PublishServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.PublishServiceResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.PublishServiceResponse>;

  installService(
    request: admin_pb.InstallServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.InstallServiceResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.InstallServiceResponse>;

  uninstallService(
    request: admin_pb.UninstallServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.UninstallServiceResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.UninstallServiceResponse>;

  registerExternalApplication(
    request: admin_pb.RegisterExternalApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.RegisterExternalApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.RegisterExternalApplicationResponse>;

  hasRuningProcess(
    request: admin_pb.HasRuningProcessRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: admin_pb.HasRuningProcessResponse) => void
  ): grpcWeb.ClientReadableStream<admin_pb.HasRuningProcessResponse>;

}

export class AdminServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  setRootPassword(
    request: admin_pb.SetRootPasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.SetRootPasswordResponse>;

  setRootEmail(
    request: admin_pb.SetRootEmailRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.SetRootEmailResponse>;

  setPassword(
    request: admin_pb.SetPasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.SetPasswordResponse>;

  setEmail(
    request: admin_pb.SetEmailRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.SetEmailResponse>;

  getConfig(
    request: admin_pb.GetConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.GetConfigResponse>;

  getFullConfig(
    request: admin_pb.GetConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.GetConfigResponse>;

  saveConfig(
    request: admin_pb.SaveConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.SaveConfigResponse>;

  stopService(
    request: admin_pb.StopServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.StopServiceResponse>;

  startService(
    request: admin_pb.StartServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.StartServiceResponse>;

  publishService(
    request: admin_pb.PublishServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.PublishServiceResponse>;

  installService(
    request: admin_pb.InstallServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.InstallServiceResponse>;

  uninstallService(
    request: admin_pb.UninstallServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.UninstallServiceResponse>;

  registerExternalApplication(
    request: admin_pb.RegisterExternalApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.RegisterExternalApplicationResponse>;

  hasRuningProcess(
    request: admin_pb.HasRuningProcessRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<admin_pb.HasRuningProcessResponse>;

}

