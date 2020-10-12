import * as grpcWeb from 'grpc-web';

import * as lb_pb from './lb_pb';


export class LoadBalancingServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getCanditates(
    request: lb_pb.GetCanditatesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: lb_pb.GetCanditatesResponse) => void
  ): grpcWeb.ClientReadableStream<lb_pb.GetCanditatesResponse>;

}

export class LoadBalancingServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getCanditates(
    request: lb_pb.GetCanditatesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<lb_pb.GetCanditatesResponse>;

}

