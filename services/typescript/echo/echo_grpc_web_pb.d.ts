import * as grpcWeb from 'grpc-web';

import * as echo_pb from './echo_pb';


export class EchoServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: echo_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: echo_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<echo_pb.StopResponse>;

  echo(
    request: echo_pb.EchoRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: echo_pb.EchoResponse) => void
  ): grpcWeb.ClientReadableStream<echo_pb.EchoResponse>;

}

export class EchoServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: echo_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<echo_pb.StopResponse>;

  echo(
    request: echo_pb.EchoRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<echo_pb.EchoResponse>;

}

