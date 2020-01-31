import * as grpcWeb from 'grpc-web';

import {
  CreateConnectionRqst,
  CreateConnectionRsp,
  DeleteConnectionRqst,
  DeleteConnectionRsp,
  SaveItemDefinitionRequest,
  SaveItemDefinitionResponse,
  SaveItemInstanceRequest,
  SaveItemInstanceResponse,
  SavePropertyDefinitionRequest,
  SavePropertyDefinitionResponse,
  SaveUnitOfMesureRequest,
  SaveUnitOfMesureResponse} from './catalog_pb';

export class CatalogServiceClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

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

  saveUnitOfMesure(
    request: SaveUnitOfMesureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveUnitOfMesureResponse) => void
  ): grpcWeb.ClientReadableStream<SaveUnitOfMesureResponse>;

  savePropertyDefinition(
    request: SavePropertyDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SavePropertyDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<SavePropertyDefinitionResponse>;

  saveItemDefinition(
    request: SaveItemDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveItemDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<SaveItemDefinitionResponse>;

  saveItemInstance(
    request: SaveItemInstanceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveItemInstanceResponse) => void
  ): grpcWeb.ClientReadableStream<SaveItemInstanceResponse>;

}

export class CatalogServicePromiseClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  createConnection(
    request: CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteConnectionRsp>;

  saveUnitOfMesure(
    request: SaveUnitOfMesureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveUnitOfMesureResponse>;

  savePropertyDefinition(
    request: SavePropertyDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SavePropertyDefinitionResponse>;

  saveItemDefinition(
    request: SaveItemDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveItemDefinitionResponse>;

  saveItemInstance(
    request: SaveItemInstanceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveItemInstanceResponse>;

}

