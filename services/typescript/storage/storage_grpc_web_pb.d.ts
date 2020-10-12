import * as grpcWeb from 'grpc-web';

import * as storage_pb from './storage_pb';


export class StorageServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: storage_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<storage_pb.StopResponse>;

  open(
    request: storage_pb.OpenRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.OpenRsp) => void
  ): grpcWeb.ClientReadableStream<storage_pb.OpenRsp>;

  close(
    request: storage_pb.CloseRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.CloseRsp) => void
  ): grpcWeb.ClientReadableStream<storage_pb.CloseRsp>;

  createConnection(
    request: storage_pb.CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<storage_pb.CreateConnectionRsp>;

  deleteConnection(
    request: storage_pb.DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<storage_pb.DeleteConnectionRsp>;

  setItem(
    request: storage_pb.SetItemRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.SetItemResponse) => void
  ): grpcWeb.ClientReadableStream<storage_pb.SetItemResponse>;

  getItem(
    request: storage_pb.GetItemRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.GetItemResponse) => void
  ): grpcWeb.ClientReadableStream<storage_pb.GetItemResponse>;

  removeItem(
    request: storage_pb.RemoveItemRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.RemoveItemResponse) => void
  ): grpcWeb.ClientReadableStream<storage_pb.RemoveItemResponse>;

  clear(
    request: storage_pb.ClearRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.ClearResponse) => void
  ): grpcWeb.ClientReadableStream<storage_pb.ClearResponse>;

  drop(
    request: storage_pb.DropRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: storage_pb.DropResponse) => void
  ): grpcWeb.ClientReadableStream<storage_pb.DropResponse>;

}

export class StorageServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: storage_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.StopResponse>;

  open(
    request: storage_pb.OpenRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.OpenRsp>;

  close(
    request: storage_pb.CloseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.CloseRsp>;

  createConnection(
    request: storage_pb.CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.CreateConnectionRsp>;

  deleteConnection(
    request: storage_pb.DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.DeleteConnectionRsp>;

  setItem(
    request: storage_pb.SetItemRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.SetItemResponse>;

  getItem(
    request: storage_pb.GetItemRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.GetItemResponse>;

  removeItem(
    request: storage_pb.RemoveItemRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.RemoveItemResponse>;

  clear(
    request: storage_pb.ClearRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.ClearResponse>;

  drop(
    request: storage_pb.DropRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<storage_pb.DropResponse>;

}

