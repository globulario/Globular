import * as jspb from 'google-protobuf'



export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): Connection;

  getName(): string;
  setName(value: string): Connection;

  getHost(): string;
  setHost(value: string): Connection;

  getStore(): StoreType;
  setStore(value: StoreType): Connection;

  getUser(): string;
  setUser(value: string): Connection;

  getPassword(): string;
  setPassword(value: string): Connection;

  getPort(): number;
  setPort(value: number): Connection;

  getTimeout(): number;
  setTimeout(value: number): Connection;

  getOptions(): string;
  setOptions(value: string): Connection;

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
  setConnection(value?: Connection): CreateConnectionRqst;
  hasConnection(): boolean;
  clearConnection(): CreateConnectionRqst;

  getSave(): boolean;
  setSave(value: boolean): CreateConnectionRqst;

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
    save: boolean,
  }
}

export class CreateConnectionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): CreateConnectionRsp;

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
  setId(value: string): DeleteConnectionRqst;

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
  setResult(value: boolean): DeleteConnectionRsp;

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
  setId(value: string): PingConnectionRqst;

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
  setResult(value: string): PingConnectionRsp;

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
  setId(value: string): InsertManyRqst;

  getDatabase(): string;
  setDatabase(value: string): InsertManyRqst;

  getCollection(): string;
  setCollection(value: string): InsertManyRqst;

  getJsonstr(): string;
  setJsonstr(value: string): InsertManyRqst;

  getOptions(): string;
  setOptions(value: string): InsertManyRqst;

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
  setIds(value: string): InsertManyRsp;

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
  setId(value: string): InsertOneRqst;

  getDatabase(): string;
  setDatabase(value: string): InsertOneRqst;

  getCollection(): string;
  setCollection(value: string): InsertOneRqst;

  getJsonstr(): string;
  setJsonstr(value: string): InsertOneRqst;

  getOptions(): string;
  setOptions(value: string): InsertOneRqst;

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
  setId(value: string): InsertOneRsp;

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
  setId(value: string): FindRqst;

  getDatabase(): string;
  setDatabase(value: string): FindRqst;

  getCollection(): string;
  setCollection(value: string): FindRqst;

  getQuery(): string;
  setQuery(value: string): FindRqst;

  getOptions(): string;
  setOptions(value: string): FindRqst;

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
  setJsonstr(value: string): FindResp;

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
  setId(value: string): FindOneRqst;

  getDatabase(): string;
  setDatabase(value: string): FindOneRqst;

  getCollection(): string;
  setCollection(value: string): FindOneRqst;

  getQuery(): string;
  setQuery(value: string): FindOneRqst;

  getOptions(): string;
  setOptions(value: string): FindOneRqst;

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
  setJsonstr(value: string): FindOneResp;

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
  setId(value: string): AggregateRqst;

  getDatabase(): string;
  setDatabase(value: string): AggregateRqst;

  getCollection(): string;
  setCollection(value: string): AggregateRqst;

  getPipeline(): string;
  setPipeline(value: string): AggregateRqst;

  getOptions(): string;
  setOptions(value: string): AggregateRqst;

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
  setJsonstr(value: string): AggregateResp;

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
  setId(value: string): UpdateRqst;

  getDatabase(): string;
  setDatabase(value: string): UpdateRqst;

  getCollection(): string;
  setCollection(value: string): UpdateRqst;

  getQuery(): string;
  setQuery(value: string): UpdateRqst;

  getValue(): string;
  setValue(value: string): UpdateRqst;

  getOptions(): string;
  setOptions(value: string): UpdateRqst;

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
  setResult(value: boolean): UpdateRsp;

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
  setId(value: string): UpdateOneRqst;

  getDatabase(): string;
  setDatabase(value: string): UpdateOneRqst;

  getCollection(): string;
  setCollection(value: string): UpdateOneRqst;

  getQuery(): string;
  setQuery(value: string): UpdateOneRqst;

  getValue(): string;
  setValue(value: string): UpdateOneRqst;

  getOptions(): string;
  setOptions(value: string): UpdateOneRqst;

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
  setResult(value: boolean): UpdateOneRsp;

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
  setId(value: string): ReplaceOneRqst;

  getDatabase(): string;
  setDatabase(value: string): ReplaceOneRqst;

  getCollection(): string;
  setCollection(value: string): ReplaceOneRqst;

  getQuery(): string;
  setQuery(value: string): ReplaceOneRqst;

  getValue(): string;
  setValue(value: string): ReplaceOneRqst;

  getOptions(): string;
  setOptions(value: string): ReplaceOneRqst;

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
  setResult(value: boolean): ReplaceOneRsp;

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
  setId(value: string): DeleteRqst;

  getDatabase(): string;
  setDatabase(value: string): DeleteRqst;

  getCollection(): string;
  setCollection(value: string): DeleteRqst;

  getQuery(): string;
  setQuery(value: string): DeleteRqst;

  getOptions(): string;
  setOptions(value: string): DeleteRqst;

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
  setResult(value: boolean): DeleteRsp;

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
  setId(value: string): DeleteOneRqst;

  getDatabase(): string;
  setDatabase(value: string): DeleteOneRqst;

  getCollection(): string;
  setCollection(value: string): DeleteOneRqst;

  getQuery(): string;
  setQuery(value: string): DeleteOneRqst;

  getOptions(): string;
  setOptions(value: string): DeleteOneRqst;

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
  setResult(value: boolean): DeleteOneRsp;

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
  setId(value: string): CreateDatabaseRqst;

  getDatabase(): string;
  setDatabase(value: string): CreateDatabaseRqst;

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
  setResult(value: boolean): CreateDatabaseRsp;

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
  setId(value: string): DeleteDatabaseRqst;

  getDatabase(): string;
  setDatabase(value: string): DeleteDatabaseRqst;

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
  setResult(value: boolean): DeleteDatabaseRsp;

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
  setId(value: string): CreateCollectionRqst;

  getDatabase(): string;
  setDatabase(value: string): CreateCollectionRqst;

  getCollection(): string;
  setCollection(value: string): CreateCollectionRqst;

  getOptionsstr(): string;
  setOptionsstr(value: string): CreateCollectionRqst;

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
  setResult(value: boolean): CreateCollectionRsp;

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
  setId(value: string): DeleteCollectionRqst;

  getDatabase(): string;
  setDatabase(value: string): DeleteCollectionRqst;

  getCollection(): string;
  setCollection(value: string): DeleteCollectionRqst;

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
  setResult(value: boolean): DeleteCollectionRsp;

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
  setId(value: string): CountRqst;

  getDatabase(): string;
  setDatabase(value: string): CountRqst;

  getCollection(): string;
  setCollection(value: string): CountRqst;

  getQuery(): string;
  setQuery(value: string): CountRqst;

  getOptions(): string;
  setOptions(value: string): CountRqst;

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
  setResult(value: number): CountRsp;

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

export class RunAdminCmdRqst extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): RunAdminCmdRqst;

  getUser(): string;
  setUser(value: string): RunAdminCmdRqst;

  getPassword(): string;
  setPassword(value: string): RunAdminCmdRqst;

  getScript(): string;
  setScript(value: string): RunAdminCmdRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunAdminCmdRqst.AsObject;
  static toObject(includeInstance: boolean, msg: RunAdminCmdRqst): RunAdminCmdRqst.AsObject;
  static serializeBinaryToWriter(message: RunAdminCmdRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunAdminCmdRqst;
  static deserializeBinaryFromReader(message: RunAdminCmdRqst, reader: jspb.BinaryReader): RunAdminCmdRqst;
}

export namespace RunAdminCmdRqst {
  export type AsObject = {
    connectionid: string,
    user: string,
    password: string,
    script: string,
  }
}

export class RunAdminCmdRsp extends jspb.Message {
  getResult(): string;
  setResult(value: string): RunAdminCmdRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunAdminCmdRsp.AsObject;
  static toObject(includeInstance: boolean, msg: RunAdminCmdRsp): RunAdminCmdRsp.AsObject;
  static serializeBinaryToWriter(message: RunAdminCmdRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunAdminCmdRsp;
  static deserializeBinaryFromReader(message: RunAdminCmdRsp, reader: jspb.BinaryReader): RunAdminCmdRsp;
}

export namespace RunAdminCmdRsp {
  export type AsObject = {
    result: string,
  }
}

export class ConnectRqst extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): ConnectRqst;

  getPassword(): string;
  setPassword(value: string): ConnectRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectRqst.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectRqst): ConnectRqst.AsObject;
  static serializeBinaryToWriter(message: ConnectRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectRqst;
  static deserializeBinaryFromReader(message: ConnectRqst, reader: jspb.BinaryReader): ConnectRqst;
}

export namespace ConnectRqst {
  export type AsObject = {
    connectionid: string,
    password: string,
  }
}

export class ConnectRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): ConnectRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectRsp.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectRsp): ConnectRsp.AsObject;
  static serializeBinaryToWriter(message: ConnectRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectRsp;
  static deserializeBinaryFromReader(message: ConnectRsp, reader: jspb.BinaryReader): ConnectRsp;
}

export namespace ConnectRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class DisconnectRqst extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): DisconnectRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisconnectRqst.AsObject;
  static toObject(includeInstance: boolean, msg: DisconnectRqst): DisconnectRqst.AsObject;
  static serializeBinaryToWriter(message: DisconnectRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisconnectRqst;
  static deserializeBinaryFromReader(message: DisconnectRqst, reader: jspb.BinaryReader): DisconnectRqst;
}

export namespace DisconnectRqst {
  export type AsObject = {
    connectionid: string,
  }
}

export class DisconnectRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): DisconnectRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisconnectRsp.AsObject;
  static toObject(includeInstance: boolean, msg: DisconnectRsp): DisconnectRsp.AsObject;
  static serializeBinaryToWriter(message: DisconnectRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisconnectRsp;
  static deserializeBinaryFromReader(message: DisconnectRsp, reader: jspb.BinaryReader): DisconnectRsp;
}

export namespace DisconnectRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class StopRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StopRequest): StopRequest.AsObject;
  static serializeBinaryToWriter(message: StopRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopRequest;
  static deserializeBinaryFromReader(message: StopRequest, reader: jspb.BinaryReader): StopRequest;
}

export namespace StopRequest {
  export type AsObject = {
  }
}

export class StopResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StopResponse): StopResponse.AsObject;
  static serializeBinaryToWriter(message: StopResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopResponse;
  static deserializeBinaryFromReader(message: StopResponse, reader: jspb.BinaryReader): StopResponse;
}

export namespace StopResponse {
  export type AsObject = {
  }
}

export enum StoreType { 
  MONGO = 0,
}
