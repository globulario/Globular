import * as jspb from "google-protobuf"

export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getHost(): string;
  setHost(value: string): void;

  getCharset(): string;
  setCharset(value: string): void;

  getDriver(): string;
  setDriver(value: string): void;

  getUser(): string;
  setUser(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  getPort(): number;
  setPort(value: number): void;

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
    charset: string,
    driver: string,
    user: string,
    password: string,
    port: number,
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

export class Query extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getParameters(): string;
  setParameters(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Query.AsObject;
  static toObject(includeInstance: boolean, msg: Query): Query.AsObject;
  static serializeBinaryToWriter(message: Query, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Query;
  static deserializeBinaryFromReader(message: Query, reader: jspb.BinaryReader): Query;
}

export namespace Query {
  export type AsObject = {
    connectionid: string,
    query: string,
    parameters: string,
  }
}

export class QueryContextRqst extends jspb.Message {
  getQuery(): Query | undefined;
  setQuery(value?: Query): void;
  hasQuery(): boolean;
  clearQuery(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryContextRqst.AsObject;
  static toObject(includeInstance: boolean, msg: QueryContextRqst): QueryContextRqst.AsObject;
  static serializeBinaryToWriter(message: QueryContextRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryContextRqst;
  static deserializeBinaryFromReader(message: QueryContextRqst, reader: jspb.BinaryReader): QueryContextRqst;
}

export namespace QueryContextRqst {
  export type AsObject = {
    query?: Query.AsObject,
  }
}

export class QueryContextRsp extends jspb.Message {
  getHeader(): string;
  setHeader(value: string): void;
  hasHeader(): boolean;

  getRows(): string;
  setRows(value: string): void;
  hasRows(): boolean;

  getResultCase(): QueryContextRsp.ResultCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryContextRsp.AsObject;
  static toObject(includeInstance: boolean, msg: QueryContextRsp): QueryContextRsp.AsObject;
  static serializeBinaryToWriter(message: QueryContextRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryContextRsp;
  static deserializeBinaryFromReader(message: QueryContextRsp, reader: jspb.BinaryReader): QueryContextRsp;
}

export namespace QueryContextRsp {
  export type AsObject = {
    header: string,
    rows: string,
  }

  export enum ResultCase { 
    RESULT_NOT_SET = 0,
    HEADER = 1,
    ROWS = 2,
  }
}

export class ExecContextRqst extends jspb.Message {
  getQuery(): Query | undefined;
  setQuery(value?: Query): void;
  hasQuery(): boolean;
  clearQuery(): void;

  getTx(): boolean;
  setTx(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecContextRqst.AsObject;
  static toObject(includeInstance: boolean, msg: ExecContextRqst): ExecContextRqst.AsObject;
  static serializeBinaryToWriter(message: ExecContextRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecContextRqst;
  static deserializeBinaryFromReader(message: ExecContextRqst, reader: jspb.BinaryReader): ExecContextRqst;
}

export namespace ExecContextRqst {
  export type AsObject = {
    query?: Query.AsObject,
    tx: boolean,
  }
}

export class ExecContextRsp extends jspb.Message {
  getAffectedrows(): number;
  setAffectedrows(value: number): void;

  getLastid(): number;
  setLastid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecContextRsp.AsObject;
  static toObject(includeInstance: boolean, msg: ExecContextRsp): ExecContextRsp.AsObject;
  static serializeBinaryToWriter(message: ExecContextRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecContextRsp;
  static deserializeBinaryFromReader(message: ExecContextRsp, reader: jspb.BinaryReader): ExecContextRsp;
}

export namespace ExecContextRsp {
  export type AsObject = {
    affectedrows: number,
    lastid: number,
  }
}

