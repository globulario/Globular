import * as jspb from 'google-protobuf'



export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): Connection;

  getDomain(): string;
  setDomain(value: string): Connection;

  getServiceid(): string;
  setServiceid(value: string): Connection;

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
    domain: string,
    serviceid: string,
  }
}

export class Tag extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): Tag;

  getServiceid(): string;
  setServiceid(value: string): Tag;

  getConnectionid(): string;
  setConnectionid(value: string): Tag;

  getName(): string;
  setName(value: string): Tag;

  getTypename(): string;
  setTypename(value: string): Tag;

  getOffset(): number;
  setOffset(value: number): Tag;

  getLength(): number;
  setLength(value: number): Tag;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Tag.AsObject;
  static toObject(includeInstance: boolean, msg: Tag): Tag.AsObject;
  static serializeBinaryToWriter(message: Tag, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Tag;
  static deserializeBinaryFromReader(message: Tag, reader: jspb.BinaryReader): Tag;
}

export namespace Tag {
  export type AsObject = {
    domain: string,
    serviceid: string,
    connectionid: string,
    name: string,
    typename: string,
    offset: number,
    length: number,
  }
}

export class Link extends jspb.Message {
  getId(): string;
  setId(value: string): Link;

  getFrequency(): number;
  setFrequency(value: number): Link;

  getSource(): Tag | undefined;
  setSource(value?: Tag): Link;
  hasSource(): boolean;
  clearSource(): Link;

  getTarget(): Tag | undefined;
  setTarget(value?: Tag): Link;
  hasTarget(): boolean;
  clearTarget(): Link;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Link.AsObject;
  static toObject(includeInstance: boolean, msg: Link): Link.AsObject;
  static serializeBinaryToWriter(message: Link, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Link;
  static deserializeBinaryFromReader(message: Link, reader: jspb.BinaryReader): Link;
}

export namespace Link {
  export type AsObject = {
    id: string,
    frequency: number,
    source?: Tag.AsObject,
    target?: Tag.AsObject,
  }
}

export class LinkRqst extends jspb.Message {
  getLnk(): Link | undefined;
  setLnk(value?: Link): LinkRqst;
  hasLnk(): boolean;
  clearLnk(): LinkRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LinkRqst.AsObject;
  static toObject(includeInstance: boolean, msg: LinkRqst): LinkRqst.AsObject;
  static serializeBinaryToWriter(message: LinkRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LinkRqst;
  static deserializeBinaryFromReader(message: LinkRqst, reader: jspb.BinaryReader): LinkRqst;
}

export namespace LinkRqst {
  export type AsObject = {
    lnk?: Link.AsObject,
  }
}

export class LinkRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): LinkRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LinkRsp.AsObject;
  static toObject(includeInstance: boolean, msg: LinkRsp): LinkRsp.AsObject;
  static serializeBinaryToWriter(message: LinkRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LinkRsp;
  static deserializeBinaryFromReader(message: LinkRsp, reader: jspb.BinaryReader): LinkRsp;
}

export namespace LinkRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class UnLinkRqst extends jspb.Message {
  getId(): string;
  setId(value: string): UnLinkRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnLinkRqst.AsObject;
  static toObject(includeInstance: boolean, msg: UnLinkRqst): UnLinkRqst.AsObject;
  static serializeBinaryToWriter(message: UnLinkRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnLinkRqst;
  static deserializeBinaryFromReader(message: UnLinkRqst, reader: jspb.BinaryReader): UnLinkRqst;
}

export namespace UnLinkRqst {
  export type AsObject = {
    id: string,
  }
}

export class UnLinkRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): UnLinkRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnLinkRsp.AsObject;
  static toObject(includeInstance: boolean, msg: UnLinkRsp): UnLinkRsp.AsObject;
  static serializeBinaryToWriter(message: UnLinkRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnLinkRsp;
  static deserializeBinaryFromReader(message: UnLinkRsp, reader: jspb.BinaryReader): UnLinkRsp;
}

export namespace UnLinkRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class SuspendRqst extends jspb.Message {
  getId(): string;
  setId(value: string): SuspendRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SuspendRqst.AsObject;
  static toObject(includeInstance: boolean, msg: SuspendRqst): SuspendRqst.AsObject;
  static serializeBinaryToWriter(message: SuspendRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SuspendRqst;
  static deserializeBinaryFromReader(message: SuspendRqst, reader: jspb.BinaryReader): SuspendRqst;
}

export namespace SuspendRqst {
  export type AsObject = {
    id: string,
  }
}

export class SuspendRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): SuspendRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SuspendRsp.AsObject;
  static toObject(includeInstance: boolean, msg: SuspendRsp): SuspendRsp.AsObject;
  static serializeBinaryToWriter(message: SuspendRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SuspendRsp;
  static deserializeBinaryFromReader(message: SuspendRsp, reader: jspb.BinaryReader): SuspendRsp;
}

export namespace SuspendRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class ResumeRqst extends jspb.Message {
  getId(): string;
  setId(value: string): ResumeRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResumeRqst.AsObject;
  static toObject(includeInstance: boolean, msg: ResumeRqst): ResumeRqst.AsObject;
  static serializeBinaryToWriter(message: ResumeRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResumeRqst;
  static deserializeBinaryFromReader(message: ResumeRqst, reader: jspb.BinaryReader): ResumeRqst;
}

export namespace ResumeRqst {
  export type AsObject = {
    id: string,
  }
}

export class ResumeRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): ResumeRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResumeRsp.AsObject;
  static toObject(includeInstance: boolean, msg: ResumeRsp): ResumeRsp.AsObject;
  static serializeBinaryToWriter(message: ResumeRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResumeRsp;
  static deserializeBinaryFromReader(message: ResumeRsp, reader: jspb.BinaryReader): ResumeRsp;
}

export namespace ResumeRsp {
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

