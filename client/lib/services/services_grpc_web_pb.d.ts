import * as grpcWeb from 'grpc-web';

import {
  DownloadBundleRequest,
  DownloadBundleResponse,
  GetServiceDescriptorRequest,
  GetServiceDescriptorResponse,
  GetServicesDescriptorRequest,
  GetServicesDescriptorResponse,
  PublishServiceRequest,
  PublishServiceResponse,
  UploadBundleRequest,
  UploadBundleResponse} from './services_pb';

export class ServiceDiscoveryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  getServiceDescriptor(
    request: GetServiceDescriptorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetServiceDescriptorResponse) => void
  ): grpcWeb.ClientReadableStream<GetServiceDescriptorResponse>;

  getServicesDescriptor(
    request: GetServicesDescriptorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetServicesDescriptorResponse) => void
  ): grpcWeb.ClientReadableStream<GetServicesDescriptorResponse>;

  publishService(
    request: PublishServiceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PublishServiceResponse) => void
  ): grpcWeb.ClientReadableStream<PublishServiceResponse>;

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

  getServiceDescriptor(
    request: GetServiceDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetServiceDescriptorResponse>;

  getServicesDescriptor(
    request: GetServicesDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetServicesDescriptorResponse>;

  publishService(
    request: PublishServiceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<PublishServiceResponse>;

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

