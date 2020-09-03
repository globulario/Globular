import * as grpcWeb from 'grpc-web';

import {
  ClearRequest,
  ClearResponse,
  CloseRqst,
  CloseRsp,
  CreateConnectionRqst,
  CreateConnectionRsp,
  DeleteConnectionRqst,
  DeleteConnectionRsp,
  DropRequest,
  DropResponse,
  GetItemRequest,
  GetItemResponse,
  OpenRqst,
  OpenRsp,
  RemoveItemRequest,
  RemoveItemResponse,
  SetItemRequest,
  SetItemResponse} from './storage_pb';

export class StorageServiceClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  open(
    request: OpenRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OpenRsp) => void
  ): grpcWeb.ClientReadableStream<OpenRsp>;

  close(
    request: CloseRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CloseRsp) => void
  ): grpcWeb.ClientReadableStream<CloseRsp>;

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

  setItem(
    request: SetItemRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SetItemResponse) => void
  ): grpcWeb.ClientReadableStream<SetItemResponse>;

  getItem(
    request: GetItemRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetItemResponse) => void
  ): grpcWeb.ClientReadableStream<GetItemResponse>;

  removeItem(
    request: RemoveItemRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RemoveItemResponse) => void
  ): grpcWeb.ClientReadableStream<RemoveItemResponse>;

  clear(
    request: ClearRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ClearResponse) => void
  ): grpcWeb.ClientReadableStream<ClearResponse>;

  drop(
    request: DropRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DropResponse) => void
  ): grpcWeb.ClientReadableStream<DropResponse>;

}

export class StorageServicePromiseClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  open(
    request: OpenRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<OpenRsp>;

  close(
    request: CloseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CloseRsp>;

  createConnection(
    request: CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteConnectionRsp>;

  setItem(
    request: SetItemRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SetItemResponse>;

  getItem(
    request: GetItemRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetItemResponse>;

  removeItem(
    request: RemoveItemRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<RemoveItemResponse>;

  clear(
    request: ClearRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ClearResponse>;

  drop(
    request: DropRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DropResponse>;

}

