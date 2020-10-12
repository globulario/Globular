import * as grpcWeb from 'grpc-web';

import * as spc_pb from './spc_pb';


export class SpcServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: spc_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: spc_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<spc_pb.StopResponse>;

  createAnalyse(
    request: spc_pb.CreateAnalyseRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: spc_pb.CreateAnalyseRsp) => void
  ): grpcWeb.ClientReadableStream<spc_pb.CreateAnalyseRsp>;

}

export class SpcServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: spc_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<spc_pb.StopResponse>;

  createAnalyse(
    request: spc_pb.CreateAnalyseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<spc_pb.CreateAnalyseRsp>;

}

