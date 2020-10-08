import * as grpcWeb from 'grpc-web';

import {
  CreateAnalyseRqst,
  CreateAnalyseRsp,
  StopRequest,
  StopResponse} from './spc_pb';

export class SpcServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StopResponse) => void
  ): grpcWeb.ClientReadableStream<StopResponse>;

  createAnalyse(
    request: CreateAnalyseRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateAnalyseRsp) => void
  ): grpcWeb.ClientReadableStream<CreateAnalyseRsp>;

}

export class SpcServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StopResponse>;

  createAnalyse(
    request: CreateAnalyseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateAnalyseRsp>;

}

