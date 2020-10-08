import * as jspb from "google-protobuf"

export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): Connection;

  getIp(): string;
  setIp(value: string): Connection;

  getProtocol(): ProtocolType;
  setProtocol(value: ProtocolType): Connection;

  getCpu(): CpuType;
  setCpu(value: CpuType): Connection;

  getPorttype(): PortType;
  setPorttype(value: PortType): Connection;

  getSlot(): number;
  setSlot(value: number): Connection;

  getRack(): number;
  setRack(value: number): Connection;

  getTimeout(): number;
  setTimeout(value: number): Connection;

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
    ip: string,
    protocol: ProtocolType,
    cpu: CpuType,
    porttype: PortType,
    slot: number,
    rack: number,
    timeout: number,
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

export class GetConnectionRqst extends jspb.Message {
  getId(): string;
  setId(value: string): GetConnectionRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConnectionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: GetConnectionRqst): GetConnectionRqst.AsObject;
  static serializeBinaryToWriter(message: GetConnectionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConnectionRqst;
  static deserializeBinaryFromReader(message: GetConnectionRqst, reader: jspb.BinaryReader): GetConnectionRqst;
}

export namespace GetConnectionRqst {
  export type AsObject = {
    id: string,
  }
}

export class GetConnectionRsp extends jspb.Message {
  getConnection(): Connection | undefined;
  setConnection(value?: Connection): GetConnectionRsp;
  hasConnection(): boolean;
  clearConnection(): GetConnectionRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConnectionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: GetConnectionRsp): GetConnectionRsp.AsObject;
  static serializeBinaryToWriter(message: GetConnectionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConnectionRsp;
  static deserializeBinaryFromReader(message: GetConnectionRsp, reader: jspb.BinaryReader): GetConnectionRsp;
}

export namespace GetConnectionRsp {
  export type AsObject = {
    connection?: Connection.AsObject,
  }
}

export class CloseConnectionRqst extends jspb.Message {
  getConnectionId(): string;
  setConnectionId(value: string): CloseConnectionRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseConnectionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: CloseConnectionRqst): CloseConnectionRqst.AsObject;
  static serializeBinaryToWriter(message: CloseConnectionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseConnectionRqst;
  static deserializeBinaryFromReader(message: CloseConnectionRqst, reader: jspb.BinaryReader): CloseConnectionRqst;
}

export namespace CloseConnectionRqst {
  export type AsObject = {
    connectionId: string,
  }
}

export class CloseConnectionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): CloseConnectionRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloseConnectionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: CloseConnectionRsp): CloseConnectionRsp.AsObject;
  static serializeBinaryToWriter(message: CloseConnectionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloseConnectionRsp;
  static deserializeBinaryFromReader(message: CloseConnectionRsp, reader: jspb.BinaryReader): CloseConnectionRsp;
}

export namespace CloseConnectionRsp {
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

export class WriteTagRqst extends jspb.Message {
  getConnectionId(): string;
  setConnectionId(value: string): WriteTagRqst;

  getName(): string;
  setName(value: string): WriteTagRqst;

  getType(): TagType;
  setType(value: TagType): WriteTagRqst;

  getValues(): string;
  setValues(value: string): WriteTagRqst;

  getOffset(): number;
  setOffset(value: number): WriteTagRqst;

  getLength(): number;
  setLength(value: number): WriteTagRqst;

  getUnsigned(): boolean;
  setUnsigned(value: boolean): WriteTagRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WriteTagRqst.AsObject;
  static toObject(includeInstance: boolean, msg: WriteTagRqst): WriteTagRqst.AsObject;
  static serializeBinaryToWriter(message: WriteTagRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WriteTagRqst;
  static deserializeBinaryFromReader(message: WriteTagRqst, reader: jspb.BinaryReader): WriteTagRqst;
}

export namespace WriteTagRqst {
  export type AsObject = {
    connectionId: string,
    name: string,
    type: TagType,
    values: string,
    offset: number,
    length: number,
    unsigned: boolean,
  }
}

export class WriteTagRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): WriteTagRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WriteTagRsp.AsObject;
  static toObject(includeInstance: boolean, msg: WriteTagRsp): WriteTagRsp.AsObject;
  static serializeBinaryToWriter(message: WriteTagRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WriteTagRsp;
  static deserializeBinaryFromReader(message: WriteTagRsp, reader: jspb.BinaryReader): WriteTagRsp;
}

export namespace WriteTagRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class ReadTagRqst extends jspb.Message {
  getConnectionId(): string;
  setConnectionId(value: string): ReadTagRqst;

  getName(): string;
  setName(value: string): ReadTagRqst;

  getType(): TagType;
  setType(value: TagType): ReadTagRqst;

  getOffset(): number;
  setOffset(value: number): ReadTagRqst;

  getLength(): number;
  setLength(value: number): ReadTagRqst;

  getUnsigned(): boolean;
  setUnsigned(value: boolean): ReadTagRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReadTagRqst.AsObject;
  static toObject(includeInstance: boolean, msg: ReadTagRqst): ReadTagRqst.AsObject;
  static serializeBinaryToWriter(message: ReadTagRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReadTagRqst;
  static deserializeBinaryFromReader(message: ReadTagRqst, reader: jspb.BinaryReader): ReadTagRqst;
}

export namespace ReadTagRqst {
  export type AsObject = {
    connectionId: string,
    name: string,
    type: TagType,
    offset: number,
    length: number,
    unsigned: boolean,
  }
}

export class ReadTagRsp extends jspb.Message {
  getValues(): string;
  setValues(value: string): ReadTagRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReadTagRsp.AsObject;
  static toObject(includeInstance: boolean, msg: ReadTagRsp): ReadTagRsp.AsObject;
  static serializeBinaryToWriter(message: ReadTagRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReadTagRsp;
  static deserializeBinaryFromReader(message: ReadTagRsp, reader: jspb.BinaryReader): ReadTagRsp;
}

export namespace ReadTagRsp {
  export type AsObject = {
    values: string,
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

export enum CpuType { 
  PLC = 0,
  PLC5 = 1,
  SLC = 2,
  SLC500 = 3,
  MICROLOGIX = 4,
  MLGX = 5,
  COMPACTLOGIX = 6,
  CLGX = 7,
  LGX = 8,
  CONTROLLOGIX = 9,
  CONTROLOGIX = 10,
  FLEXLOGIX = 11,
  FLGX = 12,
  SIMMENS = 14,
}
export enum ProtocolType { 
  AB_EIP = 0,
  AB_CIP = 1,
}
export enum PortType { 
  BACKPLANE = 0,
  NET_ETHERNET = 1,
  DH_PLUS_CHANNEL_A = 2,
  DH_PLUS_CHANNEL_B = 3,
  SERIAL = 4,
  TCP = 5,
}
export enum TagType { 
  BOOL = 0,
  SINT = 1,
  INT = 2,
  DINT = 3,
  REAL = 4,
  LREAL = 5,
  LINT = 6,
}
