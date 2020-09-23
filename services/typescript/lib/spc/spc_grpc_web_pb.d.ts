import * as grpcWeb from 'grpc-web';

import {
  CreateAnalyseRqst,
  CreateAnalyseRsp} from './spc_pb';

export class SpcServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

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

  createAnalyse(
    request: CreateAnalyseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateAnalyseRsp>;

}

