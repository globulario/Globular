import * as grpcWeb from 'grpc-web';

import {
  GetCanditatesRequest,
  GetCanditatesResponse,
  ReportLoadInfoRequest,
  ReportLoadInfoResponse} from './lb_pb';

export class LoadBalancingServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  getCanditates(
    request: GetCanditatesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetCanditatesResponse) => void
  ): grpcWeb.ClientReadableStream<GetCanditatesResponse>;

}

export class LoadBalancingServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  getCanditates(
    request: GetCanditatesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetCanditatesResponse>;

}

