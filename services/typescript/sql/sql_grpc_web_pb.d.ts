import * as grpcWeb from 'grpc-web';

import {
  CreateConnectionRqst,
  CreateConnectionRsp,
  DeleteConnectionRqst,
  DeleteConnectionRsp,
  ExecContextRqst,
  ExecContextRsp,
  PingConnectionRqst,
  PingConnectionRsp,
  QueryContextRqst,
  QueryContextRsp,
  StopRequest,
  StopResponse} from './sql_pb';

export class SqlServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StopResponse) => void
  ): grpcWeb.ClientReadableStream<StopResponse>;

  createConnection(
    request: CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteConnectionRsp>;

  ping(
    request: PingConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PingConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<PingConnectionRsp>;

  queryContext(
    request: QueryContextRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<QueryContextRsp>;

  execContext(
    request: ExecContextRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ExecContextRsp) => void
  ): grpcWeb.ClientReadableStream<ExecContextRsp>;

}

export class SqlServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StopResponse>;

  createConnection(
    request: CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteConnectionRsp>;

  ping(
    request: PingConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<PingConnectionRsp>;

  queryContext(
    request: QueryContextRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<QueryContextRsp>;

  execContext(
    request: ExecContextRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ExecContextRsp>;

}

