import * as jspb from "google-protobuf"

export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getHost(): string;
  setHost(value: string): void;

  getStore(): StoreType;
  setStore(value: StoreType): void;

  getUser(): string;
  setUser(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  getPort(): number;
  setPort(value: number): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Connection.AsObject;
  static toObject(includeInstance: boolean, msg: Connection): Connection.AsObject;
  static serializeBinaryToWriter(message: Connection, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Connection;
  static deserializeBinaryFromReader(message: Connection, reader: jspb.BinaryReader): Connection;
}

export namespace Connection {
  export type AsObject = {
    id: string,
    name: string,
    host: string,
    store: StoreType,
    user: string,
    password: string,
    port: number,
    timeout: number,
    options: string,
  }
}

export class CreateConnectionRqst extends jspb.Message {
  getConnection(): Connection | undefined;
  setConnection(value?: Connection): void;
  hasConnection(): boolean;
  clearConnection(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateConnectionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: CreateConnectionRqst): CreateConnectionRqst.AsObject;
  static serializeBinaryToWriter(message: CreateConnectionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateConnectionRqst;
  static deserializeBinaryFromReader(message: CreateConnectionRqst, reader: jspb.BinaryReader): CreateConnectionRqst;
}

export namespace CreateConnectionRqst {
  export type AsObject = {
    connection?: Connection.AsObject,
  }
}

export class CreateConnectionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateConnectionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: CreateConnectionRsp): CreateConnectionRsp.AsObject;
  static serializeBinaryToWriter(message: CreateConnectionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateConnectionRsp;
  static deserializeBinaryFromReader(message: CreateConnectionRsp, reader: jspb.BinaryReader): CreateConnectionRsp;
}

export namespace CreateConnectionRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteConnectionRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteConnectionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteConnectionRqst): DeleteConnectionRqst.AsObject;
  static serializeBinaryToWriter(message: DeleteConnectionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteConnectionRqst;
  static deserializeBinaryFromReader(message: DeleteConnectionRqst, reader: jspb.BinaryReader): DeleteConnectionRqst;
}

export namespace DeleteConnectionRqst {
  export type AsObject = {
    id: string,
  }
}

export class DeleteConnectionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteConnectionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteConnectionRsp): DeleteConnectionRsp.AsObject;
  static serializeBinaryToWriter(message: DeleteConnectionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteConnectionRsp;
  static deserializeBinaryFromReader(message: DeleteConnectionRsp, reader: jspb.BinaryReader): DeleteConnectionRsp;
}

export namespace DeleteConnectionRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class PingConnectionRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PingConnectionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: PingConnectionRqst): PingConnectionRqst.AsObject;
  static serializeBinaryToWriter(message: PingConnectionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PingConnectionRqst;
  static deserializeBinaryFromReader(message: PingConnectionRqst, reader: jspb.BinaryReader): PingConnectionRqst;
}

export namespace PingConnectionRqst {
  export type AsObject = {
    id: string,
  }
}

export class PingConnectionRsp extends jspb.Message {
  getResult(): string;
  setResult(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PingConnectionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: PingConnectionRsp): PingConnectionRsp.AsObject;
  static serializeBinaryToWriter(message: PingConnectionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PingConnectionRsp;
  static deserializeBinaryFromReader(message: PingConnectionRsp, reader: jspb.BinaryReader): PingConnectionRsp;
}

export namespace PingConnectionRsp {
  export type AsObject = {
    result: string,
  }
}

export class InsertManyRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getJsonstr(): string;
  setJsonstr(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsertManyRqst.AsObject;
  static toObject(includeInstance: boolean, msg: InsertManyRqst): InsertManyRqst.AsObject;
  static serializeBinaryToWriter(message: InsertManyRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsertManyRqst;
  static deserializeBinaryFromReader(message: InsertManyRqst, reader: jspb.BinaryReader): InsertManyRqst;
}

export namespace InsertManyRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    jsonstr: string,
    options: string,
  }
}

export class InsertManyRsp extends jspb.Message {
  getIds(): string;
  setIds(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsertManyRsp.AsObject;
  static toObject(includeInstance: boolean, msg: InsertManyRsp): InsertManyRsp.AsObject;
  static serializeBinaryToWriter(message: InsertManyRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsertManyRsp;
  static deserializeBinaryFromReader(message: InsertManyRsp, reader: jspb.BinaryReader): InsertManyRsp;
}

export namespace InsertManyRsp {
  export type AsObject = {
    ids: string,
  }
}

export class InsertOneRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getJsonstr(): string;
  setJsonstr(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsertOneRqst.AsObject;
  static toObject(includeInstance: boolean, msg: InsertOneRqst): InsertOneRqst.AsObject;
  static serializeBinaryToWriter(message: InsertOneRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsertOneRqst;
  static deserializeBinaryFromReader(message: InsertOneRqst, reader: jspb.BinaryReader): InsertOneRqst;
}

export namespace InsertOneRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    jsonstr: string,
    options: string,
  }
}

export class InsertOneRsp extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsertOneRsp.AsObject;
  static toObject(includeInstance: boolean, msg: InsertOneRsp): InsertOneRsp.AsObject;
  static serializeBinaryToWriter(message: InsertOneRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsertOneRsp;
  static deserializeBinaryFromReader(message: InsertOneRsp, reader: jspb.BinaryReader): InsertOneRsp;
}

export namespace InsertOneRsp {
  export type AsObject = {
    id: string,
  }
}

export class FindRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindRqst.AsObject;
  static toObject(includeInstance: boolean, msg: FindRqst): FindRqst.AsObject;
  static serializeBinaryToWriter(message: FindRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindRqst;
  static deserializeBinaryFromReader(message: FindRqst, reader: jspb.BinaryReader): FindRqst;
}

export namespace FindRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    query: string,
    options: string,
  }
}

export class FindResp extends jspb.Message {
  getJsonstr(): string;
  setJsonstr(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindResp.AsObject;
  static toObject(includeInstance: boolean, msg: FindResp): FindResp.AsObject;
  static serializeBinaryToWriter(message: FindResp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindResp;
  static deserializeBinaryFromReader(message: FindResp, reader: jspb.BinaryReader): FindResp;
}

export namespace FindResp {
  export type AsObject = {
    jsonstr: string,
  }
}

export class FindOneRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindOneRqst.AsObject;
  static toObject(includeInstance: boolean, msg: FindOneRqst): FindOneRqst.AsObject;
  static serializeBinaryToWriter(message: FindOneRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindOneRqst;
  static deserializeBinaryFromReader(message: FindOneRqst, reader: jspb.BinaryReader): FindOneRqst;
}

export namespace FindOneRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    query: string,
    options: string,
  }
}

export class FindOneResp extends jspb.Message {
  getJsonstr(): string;
  setJsonstr(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindOneResp.AsObject;
  static toObject(includeInstance: boolean, msg: FindOneResp): FindOneResp.AsObject;
  static serializeBinaryToWriter(message: FindOneResp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindOneResp;
  static deserializeBinaryFromReader(message: FindOneResp, reader: jspb.BinaryReader): FindOneResp;
}

export namespace FindOneResp {
  export type AsObject = {
    jsonstr: string,
  }
}

export class AggregateRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getPipeline(): string;
  setPipeline(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AggregateRqst.AsObject;
  static toObject(includeInstance: boolean, msg: AggregateRqst): AggregateRqst.AsObject;
  static serializeBinaryToWriter(message: AggregateRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AggregateRqst;
  static deserializeBinaryFromReader(message: AggregateRqst, reader: jspb.BinaryReader): AggregateRqst;
}

export namespace AggregateRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    pipeline: string,
    options: string,
  }
}

export class AggregateResp extends jspb.Message {
  getJsonstr(): string;
  setJsonstr(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AggregateResp.AsObject;
  static toObject(includeInstance: boolean, msg: AggregateResp): AggregateResp.AsObject;
  static serializeBinaryToWriter(message: AggregateResp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AggregateResp;
  static deserializeBinaryFromReader(message: AggregateResp, reader: jspb.BinaryReader): AggregateResp;
}

export namespace AggregateResp {
  export type AsObject = {
    jsonstr: string,
  }
}

export class UpdateRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRqst.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRqst): UpdateRqst.AsObject;
  static serializeBinaryToWriter(message: UpdateRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRqst;
  static deserializeBinaryFromReader(message: UpdateRqst, reader: jspb.BinaryReader): UpdateRqst;
}

export namespace UpdateRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    query: string,
    value: string,
    options: string,
  }
}

export class UpdateRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRsp.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRsp): UpdateRsp.AsObject;
  static serializeBinaryToWriter(message: UpdateRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRsp;
  static deserializeBinaryFromReader(message: UpdateRsp, reader: jspb.BinaryReader): UpdateRsp;
}

export namespace UpdateRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class UpdateOneRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOneRqst.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOneRqst): UpdateOneRqst.AsObject;
  static serializeBinaryToWriter(message: UpdateOneRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOneRqst;
  static deserializeBinaryFromReader(message: UpdateOneRqst, reader: jspb.BinaryReader): UpdateOneRqst;
}

export namespace UpdateOneRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    query: string,
    value: string,
    options: string,
  }
}

export class UpdateOneRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOneRsp.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOneRsp): UpdateOneRsp.AsObject;
  static serializeBinaryToWriter(message: UpdateOneRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOneRsp;
  static deserializeBinaryFromReader(message: UpdateOneRsp, reader: jspb.BinaryReader): UpdateOneRsp;
}

export namespace UpdateOneRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class ReplaceOneRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReplaceOneRqst.AsObject;
  static toObject(includeInstance: boolean, msg: ReplaceOneRqst): ReplaceOneRqst.AsObject;
  static serializeBinaryToWriter(message: ReplaceOneRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReplaceOneRqst;
  static deserializeBinaryFromReader(message: ReplaceOneRqst, reader: jspb.BinaryReader): ReplaceOneRqst;
}

export namespace ReplaceOneRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    query: string,
    value: string,
    options: string,
  }
}

export class ReplaceOneRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReplaceOneRsp.AsObject;
  static toObject(includeInstance: boolean, msg: ReplaceOneRsp): ReplaceOneRsp.AsObject;
  static serializeBinaryToWriter(message: ReplaceOneRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReplaceOneRsp;
  static deserializeBinaryFromReader(message: ReplaceOneRsp, reader: jspb.BinaryReader): ReplaceOneRsp;
}

export namespace ReplaceOneRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRqst.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRqst): DeleteRqst.AsObject;
  static serializeBinaryToWriter(message: DeleteRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRqst;
  static deserializeBinaryFromReader(message: DeleteRqst, reader: jspb.BinaryReader): DeleteRqst;
}

export namespace DeleteRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    query: string,
    options: string,
  }
}

export class DeleteRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRsp.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRsp): DeleteRsp.AsObject;
  static serializeBinaryToWriter(message: DeleteRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRsp;
  static deserializeBinaryFromReader(message: DeleteRsp, reader: jspb.BinaryReader): DeleteRsp;
}

export namespace DeleteRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteOneRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteOneRqst.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteOneRqst): DeleteOneRqst.AsObject;
  static serializeBinaryToWriter(message: DeleteOneRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteOneRqst;
  static deserializeBinaryFromReader(message: DeleteOneRqst, reader: jspb.BinaryReader): DeleteOneRqst;
}

export namespace DeleteOneRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    query: string,
    options: string,
  }
}

export class DeleteOneRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteOneRsp.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteOneRsp): DeleteOneRsp.AsObject;
  static serializeBinaryToWriter(message: DeleteOneRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteOneRsp;
  static deserializeBinaryFromReader(message: DeleteOneRsp, reader: jspb.BinaryReader): DeleteOneRsp;
}

export namespace DeleteOneRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class CreateDatabaseRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateDatabaseRqst.AsObject;
  static toObject(includeInstance: boolean, msg: CreateDatabaseRqst): CreateDatabaseRqst.AsObject;
  static serializeBinaryToWriter(message: CreateDatabaseRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateDatabaseRqst;
  static deserializeBinaryFromReader(message: CreateDatabaseRqst, reader: jspb.BinaryReader): CreateDatabaseRqst;
}

export namespace CreateDatabaseRqst {
  export type AsObject = {
    id: string,
    database: string,
  }
}

export class CreateDatabaseRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateDatabaseRsp.AsObject;
  static toObject(includeInstance: boolean, msg: CreateDatabaseRsp): CreateDatabaseRsp.AsObject;
  static serializeBinaryToWriter(message: CreateDatabaseRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateDatabaseRsp;
  static deserializeBinaryFromReader(message: CreateDatabaseRsp, reader: jspb.BinaryReader): CreateDatabaseRsp;
}

export namespace CreateDatabaseRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteDatabaseRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteDatabaseRqst.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteDatabaseRqst): DeleteDatabaseRqst.AsObject;
  static serializeBinaryToWriter(message: DeleteDatabaseRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteDatabaseRqst;
  static deserializeBinaryFromReader(message: DeleteDatabaseRqst, reader: jspb.BinaryReader): DeleteDatabaseRqst;
}

export namespace DeleteDatabaseRqst {
  export type AsObject = {
    id: string,
    database: string,
  }
}

export class DeleteDatabaseRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteDatabaseRsp.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteDatabaseRsp): DeleteDatabaseRsp.AsObject;
  static serializeBinaryToWriter(message: DeleteDatabaseRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteDatabaseRsp;
  static deserializeBinaryFromReader(message: DeleteDatabaseRsp, reader: jspb.BinaryReader): DeleteDatabaseRsp;
}

export namespace DeleteDatabaseRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class CreateCollectionRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getOptionsstr(): string;
  setOptionsstr(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateCollectionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: CreateCollectionRqst): CreateCollectionRqst.AsObject;
  static serializeBinaryToWriter(message: CreateCollectionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateCollectionRqst;
  static deserializeBinaryFromReader(message: CreateCollectionRqst, reader: jspb.BinaryReader): CreateCollectionRqst;
}

export namespace CreateCollectionRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    optionsstr: string,
  }
}

export class CreateCollectionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateCollectionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: CreateCollectionRsp): CreateCollectionRsp.AsObject;
  static serializeBinaryToWriter(message: CreateCollectionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateCollectionRsp;
  static deserializeBinaryFromReader(message: CreateCollectionRsp, reader: jspb.BinaryReader): CreateCollectionRsp;
}

export namespace CreateCollectionRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteCollectionRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteCollectionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteCollectionRqst): DeleteCollectionRqst.AsObject;
  static serializeBinaryToWriter(message: DeleteCollectionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteCollectionRqst;
  static deserializeBinaryFromReader(message: DeleteCollectionRqst, reader: jspb.BinaryReader): DeleteCollectionRqst;
}

export namespace DeleteCollectionRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
  }
}

export class DeleteCollectionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteCollectionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteCollectionRsp): DeleteCollectionRsp.AsObject;
  static serializeBinaryToWriter(message: DeleteCollectionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteCollectionRsp;
  static deserializeBinaryFromReader(message: DeleteCollectionRsp, reader: jspb.BinaryReader): DeleteCollectionRsp;
}

export namespace DeleteCollectionRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class CountRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDatabase(): string;
  setDatabase(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CountRqst.AsObject;
  static toObject(includeInstance: boolean, msg: CountRqst): CountRqst.AsObject;
  static serializeBinaryToWriter(message: CountRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CountRqst;
  static deserializeBinaryFromReader(message: CountRqst, reader: jspb.BinaryReader): CountRqst;
}

export namespace CountRqst {
  export type AsObject = {
    id: string,
    database: string,
    collection: string,
    query: string,
    options: string,
  }
}

export class CountRsp extends jspb.Message {
  getResult(): number;
  setResult(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CountRsp.AsObject;
  static toObject(includeInstance: boolean, msg: CountRsp): CountRsp.AsObject;
  static serializeBinaryToWriter(message: CountRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CountRsp;
  static deserializeBinaryFromReader(message: CountRsp, reader: jspb.BinaryReader): CountRsp;
}

export namespace CountRsp {
  export type AsObject = {
    result: number,
  }
}

export enum StoreType { 
  MONGO = 0,
}
