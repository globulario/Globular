import * as grpcWeb from 'grpc-web';

import {
  GetConfigRequest,
  GetConfigResponse} from './admin_pb';

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

}

