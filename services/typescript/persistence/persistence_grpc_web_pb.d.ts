import * as grpcWeb from 'grpc-web';

import {
  AggregateResp,
  AggregateRqst,
  ConnectRqst,
  ConnectRsp,
  CountRqst,
  CountRsp,
  CreateCollectionRqst,
  CreateCollectionRsp,
  CreateConnectionRqst,
  CreateConnectionRsp,
  CreateDatabaseRqst,
  CreateDatabaseRsp,
  DeleteCollectionRqst,
  DeleteCollectionRsp,
  DeleteConnectionRqst,
  DeleteConnectionRsp,
  DeleteDatabaseRqst,
  DeleteDatabaseRsp,
  DeleteOneRqst,
  DeleteOneRsp,
  DeleteRqst,
  DeleteRsp,
  DisconnectRqst,
  DisconnectRsp,
  FindOneResp,
  FindOneRqst,
  FindResp,
  FindRqst,
  InsertManyRqst,
  InsertManyRsp,
  InsertOneRqst,
  InsertOneRsp,
  PingConnectionRqst,
  PingConnectionRsp,
  ReplaceOneRqst,
  ReplaceOneRsp,
  RunAdminCmdRqst,
  RunAdminCmdRsp,
  StopRequest,
  StopResponse,
  UpdateOneRqst,
  UpdateOneRsp,
  UpdateRqst,
  UpdateRsp} from './persistence_pb';

export class PersistenceServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StopResponse) => void
  ): grpcWeb.ClientReadableStream<StopResponse>;

  createDatabase(
    request: CreateDatabaseRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateDatabaseRsp) => void
  ): grpcWeb.ClientReadableStream<CreateDatabaseRsp>;

  connect(
    request: ConnectRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ConnectRsp) => void
  ): grpcWeb.ClientReadableStream<ConnectRsp>;

  disconnect(
    request: DisconnectRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DisconnectRsp) => void
  ): grpcWeb.ClientReadableStream<DisconnectRsp>;

  deleteDatabase(
    request: DeleteDatabaseRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteDatabaseRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteDatabaseRsp>;

  createCollection(
    request: CreateCollectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateCollectionRsp) => void
  ): grpcWeb.ClientReadableStream<CreateCollectionRsp>;

  deleteCollection(
    request: DeleteCollectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteCollectionRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteCollectionRsp>;

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

  count(
    request: CountRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CountRsp) => void
  ): grpcWeb.ClientReadableStream<CountRsp>;

  insertOne(
    request: InsertOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: InsertOneRsp) => void
  ): grpcWeb.ClientReadableStream<InsertOneRsp>;

  find(
    request: FindRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<FindResp>;

  findOne(
    request: FindOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: FindOneResp) => void
  ): grpcWeb.ClientReadableStream<FindOneResp>;

  aggregate(
    request: AggregateRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<AggregateResp>;

  update(
    request: UpdateRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UpdateRsp) => void
  ): grpcWeb.ClientReadableStream<UpdateRsp>;

  updateOne(
    request: UpdateOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UpdateOneRsp) => void
  ): grpcWeb.ClientReadableStream<UpdateOneRsp>;

  replaceOne(
    request: ReplaceOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ReplaceOneRsp) => void
  ): grpcWeb.ClientReadableStream<ReplaceOneRsp>;

  delete(
    request: DeleteRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteRsp>;

  deleteOne(
    request: DeleteOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteOneRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteOneRsp>;

  runAdminCmd(
    request: RunAdminCmdRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RunAdminCmdRsp) => void
  ): grpcWeb.ClientReadableStream<RunAdminCmdRsp>;

}

export class PersistenceServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StopResponse>;

  createDatabase(
    request: CreateDatabaseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateDatabaseRsp>;

  connect(
    request: ConnectRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ConnectRsp>;

  disconnect(
    request: DisconnectRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DisconnectRsp>;

  deleteDatabase(
    request: DeleteDatabaseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteDatabaseRsp>;

  createCollection(
    request: CreateCollectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateCollectionRsp>;

  deleteCollection(
    request: DeleteCollectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteCollectionRsp>;

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

  count(
    request: CountRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CountRsp>;

  insertOne(
    request: InsertOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<InsertOneRsp>;

  find(
    request: FindRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<FindResp>;

  findOne(
    request: FindOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<FindOneResp>;

  aggregate(
    request: AggregateRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<AggregateResp>;

  update(
    request: UpdateRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<UpdateRsp>;

  updateOne(
    request: UpdateOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<UpdateOneRsp>;

  replaceOne(
    request: ReplaceOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ReplaceOneRsp>;

  delete(
    request: DeleteRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteRsp>;

  deleteOne(
    request: DeleteOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteOneRsp>;

  runAdminCmd(
    request: RunAdminCmdRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<RunAdminCmdRsp>;

}

