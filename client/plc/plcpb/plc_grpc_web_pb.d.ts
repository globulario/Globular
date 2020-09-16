import * as grpcWeb from 'grpc-web';

import {
  CloseConnectionRqst,
  CloseConnectionRsp,
  CreateConnectionRqst,
  CreateConnectionRsp,
  DeleteConnectionRqst,
  DeleteConnectionRsp,
  GetConnectionRqst,
  GetConnectionRsp,
  ReadTagRqst,
  ReadTagRsp,
  WriteTagRqst,
  WriteTagRsp} from './plc_pb';

export class PlcServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  createConnection(
    request: CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<CreateConnectionRsp>;

  getConnection(
    request: GetConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<GetConnectionRsp>;

  closeConnection(
    request: CloseConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CloseConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<CloseConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteConnectionRsp>;

  readTag(
    request: ReadTagRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ReadTagRsp) => void
  ): grpcWeb.ClientReadableStream<ReadTagRsp>;

  writeTag(
    request: WriteTagRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: WriteTagRsp) => void
  ): grpcWeb.ClientReadableStream<WriteTagRsp>;

}

export class PlcServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  createConnection(
    request: CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateConnectionRsp>;

  getConnection(
    request: GetConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<GetConnectionRsp>;

  closeConnection(
    request: CloseConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CloseConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteConnectionRsp>;

  readTag(
    request: ReadTagRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ReadTagRsp>;

  writeTag(
    request: WriteTagRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<WriteTagRsp>;

}

