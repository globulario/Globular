import * as grpcWeb from 'grpc-web';

import * as sql_pb from './sql_pb';


export class SqlServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: sql_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: sql_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<sql_pb.StopResponse>;

  createConnection(
    request: sql_pb.CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: sql_pb.CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<sql_pb.CreateConnectionRsp>;

  deleteConnection(
    request: sql_pb.DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: sql_pb.DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<sql_pb.DeleteConnectionRsp>;

  ping(
    request: sql_pb.PingConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: sql_pb.PingConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<sql_pb.PingConnectionRsp>;

  queryContext(
    request: sql_pb.QueryContextRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<sql_pb.QueryContextRsp>;

  execContext(
    request: sql_pb.ExecContextRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: sql_pb.ExecContextRsp) => void
  ): grpcWeb.ClientReadableStream<sql_pb.ExecContextRsp>;

}

export class SqlServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: sql_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<sql_pb.StopResponse>;

  createConnection(
    request: sql_pb.CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<sql_pb.CreateConnectionRsp>;

  deleteConnection(
    request: sql_pb.DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<sql_pb.DeleteConnectionRsp>;

  ping(
    request: sql_pb.PingConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<sql_pb.PingConnectionRsp>;

  queryContext(
    request: sql_pb.QueryContextRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<sql_pb.QueryContextRsp>;

  execContext(
    request: sql_pb.ExecContextRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<sql_pb.ExecContextRsp>;

}

