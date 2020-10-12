import * as grpcWeb from 'grpc-web';

import * as plc_pb from './plc_pb';


export class PlcServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: plc_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<plc_pb.StopResponse>;

  createConnection(
    request: plc_pb.CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_pb.CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<plc_pb.CreateConnectionRsp>;

  getConnection(
    request: plc_pb.GetConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_pb.GetConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<plc_pb.GetConnectionRsp>;

  closeConnection(
    request: plc_pb.CloseConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_pb.CloseConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<plc_pb.CloseConnectionRsp>;

  deleteConnection(
    request: plc_pb.DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_pb.DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<plc_pb.DeleteConnectionRsp>;

  readTag(
    request: plc_pb.ReadTagRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_pb.ReadTagRsp) => void
  ): grpcWeb.ClientReadableStream<plc_pb.ReadTagRsp>;

  writeTag(
    request: plc_pb.WriteTagRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_pb.WriteTagRsp) => void
  ): grpcWeb.ClientReadableStream<plc_pb.WriteTagRsp>;

}

export class PlcServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: plc_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_pb.StopResponse>;

  createConnection(
    request: plc_pb.CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_pb.CreateConnectionRsp>;

  getConnection(
    request: plc_pb.GetConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_pb.GetConnectionRsp>;

  closeConnection(
    request: plc_pb.CloseConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_pb.CloseConnectionRsp>;

  deleteConnection(
    request: plc_pb.DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_pb.DeleteConnectionRsp>;

  readTag(
    request: plc_pb.ReadTagRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_pb.ReadTagRsp>;

  writeTag(
    request: plc_pb.WriteTagRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_pb.WriteTagRsp>;

}

