import * as jspb from "google-protobuf"

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

export class CarbonCopy extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): CarbonCopy;

  getName(): string;
  setName(value: string): CarbonCopy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CarbonCopy.AsObject;
  static toObject(includeInstance: boolean, msg: CarbonCopy): CarbonCopy.AsObject;
  static serializeBinaryToWriter(message: CarbonCopy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CarbonCopy;
  static deserializeBinaryFromReader(message: CarbonCopy, reader: jspb.BinaryReader): CarbonCopy;
}

export namespace CarbonCopy {
  export type AsObject = {
    address: string,
    name: string,
  }
}

export class Attachement extends jspb.Message {
  getFilename(): string;
  setFilename(value: string): Attachement;

  getFiledata(): Uint8Array | string;
  getFiledata_asU8(): Uint8Array;
  getFiledata_asB64(): string;
  setFiledata(value: Uint8Array | string): Attachement;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Attachement.AsObject;
  static toObject(includeInstance: boolean, msg: Attachement): Attachement.AsObject;
  static serializeBinaryToWriter(message: Attachement, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Attachement;
  static deserializeBinaryFromReader(message: Attachement, reader: jspb.BinaryReader): Attachement;
}

export namespace Attachement {
  export type AsObject = {
    filename: string,
    filedata: Uint8Array | string,
  }
}

export class Email extends jspb.Message {
  getFrom(): string;
  setFrom(value: string): Email;

  getToList(): Array<string>;
  setToList(value: Array<string>): Email;
  clearToList(): Email;
  addTo(value: string, index?: number): Email;

  getCcList(): Array<CarbonCopy>;
  setCcList(value: Array<CarbonCopy>): Email;
  clearCcList(): Email;
  addCc(value?: CarbonCopy, index?: number): CarbonCopy;

  getSubject(): string;
  setSubject(value: string): Email;

  getBody(): string;
  setBody(value: string): Email;

  getBodytype(): BodyType;
  setBodytype(value: BodyType): Email;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Email.AsObject;
  static toObject(includeInstance: boolean, msg: Email): Email.AsObject;
  static serializeBinaryToWriter(message: Email, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Email;
  static deserializeBinaryFromReader(message: Email, reader: jspb.BinaryReader): Email;
}

export namespace Email {
  export type AsObject = {
    from: string,
    toList: Array<string>,
    ccList: Array<CarbonCopy.AsObject>,
    subject: string,
    body: string,
    bodytype: BodyType,
  }
}

export class SendEmailRqst extends jspb.Message {
  getId(): string;
  setId(value: string): SendEmailRqst;

  getEmail(): Email | undefined;
  setEmail(value?: Email): SendEmailRqst;
  hasEmail(): boolean;
  clearEmail(): SendEmailRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendEmailRqst.AsObject;
  static toObject(includeInstance: boolean, msg: SendEmailRqst): SendEmailRqst.AsObject;
  static serializeBinaryToWriter(message: SendEmailRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendEmailRqst;
  static deserializeBinaryFromReader(message: SendEmailRqst, reader: jspb.BinaryReader): SendEmailRqst;
}

export namespace SendEmailRqst {
  export type AsObject = {
    id: string,
    email?: Email.AsObject,
  }
}

export class SendEmailRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): SendEmailRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendEmailRsp.AsObject;
  static toObject(includeInstance: boolean, msg: SendEmailRsp): SendEmailRsp.AsObject;
  static serializeBinaryToWriter(message: SendEmailRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendEmailRsp;
  static deserializeBinaryFromReader(message: SendEmailRsp, reader: jspb.BinaryReader): SendEmailRsp;
}

export namespace SendEmailRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class SendEmailWithAttachementsRqst extends jspb.Message {
  getId(): string;
  setId(value: string): SendEmailWithAttachementsRqst;

  getEmail(): Email | undefined;
  setEmail(value?: Email): SendEmailWithAttachementsRqst;
  hasEmail(): boolean;
  clearEmail(): SendEmailWithAttachementsRqst;

  getAttachements(): Attachement | undefined;
  setAttachements(value?: Attachement): SendEmailWithAttachementsRqst;
  hasAttachements(): boolean;
  clearAttachements(): SendEmailWithAttachementsRqst;

  getDataCase(): SendEmailWithAttachementsRqst.DataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendEmailWithAttachementsRqst.AsObject;
  static toObject(includeInstance: boolean, msg: SendEmailWithAttachementsRqst): SendEmailWithAttachementsRqst.AsObject;
  static serializeBinaryToWriter(message: SendEmailWithAttachementsRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendEmailWithAttachementsRqst;
  static deserializeBinaryFromReader(message: SendEmailWithAttachementsRqst, reader: jspb.BinaryReader): SendEmailWithAttachementsRqst;
}

export namespace SendEmailWithAttachementsRqst {
  export type AsObject = {
    id: string,
    email?: Email.AsObject,
    attachements?: Attachement.AsObject,
  }

  export enum DataCase { 
    DATA_NOT_SET = 0,
    EMAIL = 2,
    ATTACHEMENTS = 3,
  }
}

export class SendEmailWithAttachementsRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): SendEmailWithAttachementsRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendEmailWithAttachementsRsp.AsObject;
  static toObject(includeInstance: boolean, msg: SendEmailWithAttachementsRsp): SendEmailWithAttachementsRsp.AsObject;
  static serializeBinaryToWriter(message: SendEmailWithAttachementsRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendEmailWithAttachementsRsp;
  static deserializeBinaryFromReader(message: SendEmailWithAttachementsRsp, reader: jspb.BinaryReader): SendEmailWithAttachementsRsp;
}

export namespace SendEmailWithAttachementsRsp {
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

export enum BodyType { 
  TEXT = 0,
  HTML = 1,
}
