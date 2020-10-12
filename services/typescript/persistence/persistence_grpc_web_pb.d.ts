import * as grpcWeb from 'grpc-web';

import * as persistence_pb from './persistence_pb';


export class PersistenceServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: persistence_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.StopResponse>;

  createDatabase(
    request: persistence_pb.CreateDatabaseRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.CreateDatabaseRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.CreateDatabaseRsp>;

  connect(
    request: persistence_pb.ConnectRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.ConnectRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.ConnectRsp>;

  disconnect(
    request: persistence_pb.DisconnectRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.DisconnectRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.DisconnectRsp>;

  deleteDatabase(
    request: persistence_pb.DeleteDatabaseRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.DeleteDatabaseRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.DeleteDatabaseRsp>;

  createCollection(
    request: persistence_pb.CreateCollectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.CreateCollectionRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.CreateCollectionRsp>;

  deleteCollection(
    request: persistence_pb.DeleteCollectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.DeleteCollectionRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.DeleteCollectionRsp>;

  createConnection(
    request: persistence_pb.CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.CreateConnectionRsp>;

  deleteConnection(
    request: persistence_pb.DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.DeleteConnectionRsp>;

  ping(
    request: persistence_pb.PingConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.PingConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.PingConnectionRsp>;

  count(
    request: persistence_pb.CountRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.CountRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.CountRsp>;

  insertOne(
    request: persistence_pb.InsertOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.InsertOneRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.InsertOneRsp>;

  find(
    request: persistence_pb.FindRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<persistence_pb.FindResp>;

  findOne(
    request: persistence_pb.FindOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.FindOneResp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.FindOneResp>;

  aggregate(
    request: persistence_pb.AggregateRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<persistence_pb.AggregateResp>;

  update(
    request: persistence_pb.UpdateRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.UpdateRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.UpdateRsp>;

  updateOne(
    request: persistence_pb.UpdateOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.UpdateOneRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.UpdateOneRsp>;

  replaceOne(
    request: persistence_pb.ReplaceOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.ReplaceOneRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.ReplaceOneRsp>;

  delete(
    request: persistence_pb.DeleteRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.DeleteRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.DeleteRsp>;

  deleteOne(
    request: persistence_pb.DeleteOneRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.DeleteOneRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.DeleteOneRsp>;

  runAdminCmd(
    request: persistence_pb.RunAdminCmdRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: persistence_pb.RunAdminCmdRsp) => void
  ): grpcWeb.ClientReadableStream<persistence_pb.RunAdminCmdRsp>;

}

export class PersistenceServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: persistence_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.StopResponse>;

  createDatabase(
    request: persistence_pb.CreateDatabaseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.CreateDatabaseRsp>;

  connect(
    request: persistence_pb.ConnectRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.ConnectRsp>;

  disconnect(
    request: persistence_pb.DisconnectRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.DisconnectRsp>;

  deleteDatabase(
    request: persistence_pb.DeleteDatabaseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.DeleteDatabaseRsp>;

  createCollection(
    request: persistence_pb.CreateCollectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.CreateCollectionRsp>;

  deleteCollection(
    request: persistence_pb.DeleteCollectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.DeleteCollectionRsp>;

  createConnection(
    request: persistence_pb.CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.CreateConnectionRsp>;

  deleteConnection(
    request: persistence_pb.DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.DeleteConnectionRsp>;

  ping(
    request: persistence_pb.PingConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.PingConnectionRsp>;

  count(
    request: persistence_pb.CountRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.CountRsp>;

  insertOne(
    request: persistence_pb.InsertOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.InsertOneRsp>;

  find(
    request: persistence_pb.FindRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<persistence_pb.FindResp>;

  findOne(
    request: persistence_pb.FindOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.FindOneResp>;

  aggregate(
    request: persistence_pb.AggregateRqst,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<persistence_pb.AggregateResp>;

  update(
    request: persistence_pb.UpdateRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.UpdateRsp>;

  updateOne(
    request: persistence_pb.UpdateOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.UpdateOneRsp>;

  replaceOne(
    request: persistence_pb.ReplaceOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.ReplaceOneRsp>;

  delete(
    request: persistence_pb.DeleteRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.DeleteRsp>;

  deleteOne(
    request: persistence_pb.DeleteOneRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.DeleteOneRsp>;

  runAdminCmd(
    request: persistence_pb.RunAdminCmdRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<persistence_pb.RunAdminCmdRsp>;

}

