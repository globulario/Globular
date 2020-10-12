import * as grpcWeb from 'grpc-web';

import * as services_pb from './services_pb';


export class ServiceDiscoveryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  findServices(
    request: services_pb.FindServicesDescriptorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: services_pb.FindServicesDescriptorResponse) => void
  ): grpcWeb.ClientReadableStream<services_pb.FindServicesDescriptorResponse>;

  getServiceDescriptor(
    request: services_pb.GetServiceDescriptorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: services_pb.GetServiceDescriptorResponse) => void
  ): grpcWeb.ClientReadableStream<services_pb.GetServiceDescriptorResponse>;

  getServicesDescriptor(
    request: services_pb.GetServicesDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.GetServicesDescriptorResponse>;

  setServiceDescriptor(
    request: services_pb.SetServiceDescriptorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: services_pb.SetServiceDescriptorResponse) => void
  ): grpcWeb.ClientReadableStream<services_pb.SetServiceDescriptorResponse>;

  publishServiceDescriptor(
    request: services_pb.PublishServiceDescriptorRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: services_pb.PublishServiceDescriptorResponse) => void
  ): grpcWeb.ClientReadableStream<services_pb.PublishServiceDescriptorResponse>;

}

export class ServiceRepositoryClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  downloadBundle(
    request: services_pb.DownloadBundleRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.DownloadBundleResponse>;

}

export class ServiceDiscoveryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  findServices(
    request: services_pb.FindServicesDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.FindServicesDescriptorResponse>;

  getServiceDescriptor(
    request: services_pb.GetServiceDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.GetServiceDescriptorResponse>;

  getServicesDescriptor(
    request: services_pb.GetServicesDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.GetServicesDescriptorResponse>;

  setServiceDescriptor(
    request: services_pb.SetServiceDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.SetServiceDescriptorResponse>;

  publishServiceDescriptor(
    request: services_pb.PublishServiceDescriptorRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<services_pb.PublishServiceDescriptorResponse>;

}

export class ServiceRepositoryPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  downloadBundle(
    request: services_pb.DownloadBundleRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<services_pb.DownloadBundleResponse>;

}

