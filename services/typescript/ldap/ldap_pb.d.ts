import * as jspb from 'google-protobuf'



export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): Connection;

  getHost(): string;
  setHost(value: string): Connection;

  getUser(): string;
  setUser(value: string): Connection;

  getPassword(): string;
  setPassword(value: string): Connection;

  getPort(): number;
  setPort(value: number): Connection;

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
    host: string,
    user: string,
    password: string,
    port: number,
  }
}

export class CreateConnectionRqst extends jspb.Message {
  getConnection(): Connection | undefined;
  setConnection(value?: Connection): CreateConnectionRqst;
  hasConnection(): boolean;
  clearConnection(): CreateConnectionRqst;

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

export class CloseRqst extends jspb.Message {
  getId(): string;
  setId(value: string): CloseRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseRqst.AsObject;
  static toObject(includeInstance: boolean, msg: CloseRqst): CloseRqst.AsObject;
  static serializeBinaryToWriter(message: CloseRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseRqst;
  static deserializeBinaryFromReader(message: CloseRqst, reader: jspb.BinaryReader): CloseRqst;
}

export namespace CloseRqst {
  export type AsObject = {
    id: string,
  }
}

export class CloseRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): CloseRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseRsp.AsObject;
  static toObject(includeInstance: boolean, msg: CloseRsp): CloseRsp.AsObject;
  static serializeBinaryToWriter(message: CloseRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseRsp;
  static deserializeBinaryFromReader(message: CloseRsp, reader: jspb.BinaryReader): CloseRsp;
}

export namespace CloseRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class Search extends jspb.Message {
  getId(): string;
  setId(value: string): Search;

  getBasedn(): string;
  setBasedn(value: string): Search;

  getFilter(): string;
  setFilter(value: string): Search;

  getAttributesList(): Array<string>;
  setAttributesList(value: Array<string>): Search;
  clearAttributesList(): Search;
  addAttributes(value: string, index?: number): Search;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Search.AsObject;
  static toObject(includeInstance: boolean, msg: Search): Search.AsObject;
  static serializeBinaryToWriter(message: Search, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Search;
  static deserializeBinaryFromReader(message: Search, reader: jspb.BinaryReader): Search;
}

export namespace Search {
  export type AsObject = {
    id: string,
    basedn: string,
    filter: string,
    attributesList: Array<string>,
  }
}

export class SearchRqst extends jspb.Message {
  getSearch(): Search | undefined;
  setSearch(value?: Search): SearchRqst;
  hasSearch(): boolean;
  clearSearch(): SearchRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchRqst.AsObject;
  static toObject(includeInstance: boolean, msg: SearchRqst): SearchRqst.AsObject;
  static serializeBinaryToWriter(message: SearchRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchRqst;
  static deserializeBinaryFromReader(message: SearchRqst, reader: jspb.BinaryReader): SearchRqst;
}

export namespace SearchRqst {
  export type AsObject = {
    search?: Search.AsObject,
  }
}

export class SearchResp extends jspb.Message {
  getResult(): string;
  setResult(value: string): SearchResp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchResp.AsObject;
  static toObject(includeInstance: boolean, msg: SearchResp): SearchResp.AsObject;
  static serializeBinaryToWriter(message: SearchResp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchResp;
  static deserializeBinaryFromReader(message: SearchResp, reader: jspb.BinaryReader): SearchResp;
}

export namespace SearchResp {
  export type AsObject = {
    result: string,
  }
}

export class AuthenticateRqst extends jspb.Message {
  getId(): string;
  setId(value: string): AuthenticateRqst;

  getLogin(): string;
  setLogin(value: string): AuthenticateRqst;

  getPwd(): string;
  setPwd(value: string): AuthenticateRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthenticateRqst.AsObject;
  static toObject(includeInstance: boolean, msg: AuthenticateRqst): AuthenticateRqst.AsObject;
  static serializeBinaryToWriter(message: AuthenticateRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthenticateRqst;
  static deserializeBinaryFromReader(message: AuthenticateRqst, reader: jspb.BinaryReader): AuthenticateRqst;
}

export namespace AuthenticateRqst {
  export type AsObject = {
    id: string,
    login: string,
    pwd: string,
  }
}

export class AuthenticateRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): AuthenticateRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthenticateRsp.AsObject;
  static toObject(includeInstance: boolean, msg: AuthenticateRsp): AuthenticateRsp.AsObject;
  static serializeBinaryToWriter(message: AuthenticateRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthenticateRsp;
  static deserializeBinaryFromReader(message: AuthenticateRsp, reader: jspb.BinaryReader): AuthenticateRsp;
}

export namespace AuthenticateRsp {
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

