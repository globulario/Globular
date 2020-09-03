import * as jspb from "google-protobuf"

export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getType(): StoreType;
  setType(value: StoreType): void;

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
    type: StoreType,
  }
}

export class OpenRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpenRqst.AsObject;
  static toObject(includeInstance: boolean, msg: OpenRqst): OpenRqst.AsObject;
  static serializeBinaryToWriter(message: OpenRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpenRqst;
  static deserializeBinaryFromReader(message: OpenRqst, reader: jspb.BinaryReader): OpenRqst;
}

export namespace OpenRqst {
  export type AsObject = {
    id: string,
    options: string,
  }
}

export class OpenRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpenRsp.AsObject;
  static toObject(includeInstance: boolean, msg: OpenRsp): OpenRsp.AsObject;
  static serializeBinaryToWriter(message: OpenRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpenRsp;
  static deserializeBinaryFromReader(message: OpenRsp, reader: jspb.BinaryReader): OpenRsp;
}

export namespace OpenRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class CloseRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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
  setResult(value: boolean): void;

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

export class SetItemRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getKey(): string;
  setKey(value: string): void;

  getValue(): Uint8Array | string;
  getValue_asU8(): Uint8Array;
  getValue_asB64(): string;
  setValue(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetItemRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetItemRequest): SetItemRequest.AsObject;
  static serializeBinaryToWriter(message: SetItemRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetItemRequest;
  static deserializeBinaryFromReader(message: SetItemRequest, reader: jspb.BinaryReader): SetItemRequest;
}

export namespace SetItemRequest {
  export type AsObject = {
    id: string,
    key: string,
    value: Uint8Array | string,
  }
}

export class SetItemResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetItemResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetItemResponse): SetItemResponse.AsObject;
  static serializeBinaryToWriter(message: SetItemResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetItemResponse;
  static deserializeBinaryFromReader(message: SetItemResponse, reader: jspb.BinaryReader): SetItemResponse;
}

export namespace SetItemResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class GetItemRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getKey(): string;
  setKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemRequest): GetItemRequest.AsObject;
  static serializeBinaryToWriter(message: GetItemRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemRequest;
  static deserializeBinaryFromReader(message: GetItemRequest, reader: jspb.BinaryReader): GetItemRequest;
}

export namespace GetItemRequest {
  export type AsObject = {
    id: string,
    key: string,
  }
}

export class GetItemResponse extends jspb.Message {
  getResult(): Uint8Array | string;
  getResult_asU8(): Uint8Array;
  getResult_asB64(): string;
  setResult(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemResponse): GetItemResponse.AsObject;
  static serializeBinaryToWriter(message: GetItemResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemResponse;
  static deserializeBinaryFromReader(message: GetItemResponse, reader: jspb.BinaryReader): GetItemResponse;
}

export namespace GetItemResponse {
  export type AsObject = {
    result: Uint8Array | string,
  }
}

export class RemoveItemRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getKey(): string;
  setKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveItemRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveItemRequest): RemoveItemRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveItemRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveItemRequest;
  static deserializeBinaryFromReader(message: RemoveItemRequest, reader: jspb.BinaryReader): RemoveItemRequest;
}

export namespace RemoveItemRequest {
  export type AsObject = {
    id: string,
    key: string,
  }
}

export class RemoveItemResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveItemResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveItemResponse): RemoveItemResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveItemResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveItemResponse;
  static deserializeBinaryFromReader(message: RemoveItemResponse, reader: jspb.BinaryReader): RemoveItemResponse;
}

export namespace RemoveItemResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class ClearRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClearRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ClearRequest): ClearRequest.AsObject;
  static serializeBinaryToWriter(message: ClearRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClearRequest;
  static deserializeBinaryFromReader(message: ClearRequest, reader: jspb.BinaryReader): ClearRequest;
}

export namespace ClearRequest {
  export type AsObject = {
    id: string,
  }
}

export class ClearResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClearResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ClearResponse): ClearResponse.AsObject;
  static serializeBinaryToWriter(message: ClearResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClearResponse;
  static deserializeBinaryFromReader(message: ClearResponse, reader: jspb.BinaryReader): ClearResponse;
}

export namespace ClearResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DropRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DropRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DropRequest): DropRequest.AsObject;
  static serializeBinaryToWriter(message: DropRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DropRequest;
  static deserializeBinaryFromReader(message: DropRequest, reader: jspb.BinaryReader): DropRequest;
}

export namespace DropRequest {
  export type AsObject = {
    id: string,
  }
}

export class DropResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DropResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DropResponse): DropResponse.AsObject;
  static serializeBinaryToWriter(message: DropResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DropResponse;
  static deserializeBinaryFromReader(message: DropResponse, reader: jspb.BinaryReader): DropResponse;
}

export namespace DropResponse {
  export type AsObject = {
    result: boolean,
  }
}

export enum StoreType { 
  LEVEL_DB = 0,
  BIG_CACHE = 1,
}
