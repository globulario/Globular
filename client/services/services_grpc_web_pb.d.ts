import * as grpcWeb from 'grpc-web';

import {
  DownloadBundleRequest,
  DownloadBundleResponse,
  FindServicesDescriptorRequest,
  FindServicesDescriptorResponse,
  GetServiceDescriptorRequest,
  GetServiceDescriptorResponse,
  GetServicesDescriptorRequest,
  GetServicesDescriptorResponse,
  PublishServiceDescriptorRequest,
  PublishServiceDescriptorResponse,
  UploadBundleRequest,
  UploadBundleResponse} from './services_pb';

export class ServiceDiscoveryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  findServices(
    request: FindServicesDescriptorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: FindServicesDescriptorResponse) => void
  ): grpcWeb.ClientReadableStream<FindServicesDescriptorResponse>;

  getServiceDescriptor(
    request: GetServiceDescriptorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetServiceDescriptorResponse) => void
  ): grpcWeb.ClientReadableStream<GetServiceDescriptorResponse>;

  getServicesDescriptor(
    request: GetServicesDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<GetServicesDescriptorResponse>;

  publishServiceDescriptor(
    request: PublishServiceDescriptorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PublishServiceDescriptorResponse) => void
  ): grpcWeb.ClientReadableStream<PublishServiceDescriptorResponse>;

}

export class ServiceRepositoryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  downloadBundle(
    request: DownloadBundleRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<DownloadBundleResponse>;

}

export class ServiceDiscoveryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  findServices(
    request: FindServicesDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<FindServicesDescriptorResponse>;

  getServiceDescriptor(
    request: GetServiceDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetServiceDescriptorResponse>;

  getServicesDescriptor(
    request: GetServicesDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<GetServicesDescriptorResponse>;

  publishServiceDescriptor(
    request: PublishServiceDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<PublishServiceDescriptorResponse>;

}

export class ServiceRepositoryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  downloadBundle(
    request: DownloadBundleRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<DownloadBundleResponse>;

}

