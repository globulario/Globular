import * as jspb from "google-protobuf"

export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDomain(): string;
  setDomain(value: string): void;

  getAddress(): string;
  setAddress(value: string): void;

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
    address: string,
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

export class Tag extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): void;

  getAddress(): string;
  setAddress(value: string): void;

  getConnectionid(): string;
  setConnectionid(value: string): void;

  getName(): string;
  setName(value: string): void;

  getLabel(): string;
  setLabel(value: string): void;

  getTypename(): string;
  setTypename(value: string): void;

  getOffset(): number;
  setOffset(value: number): void;

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
    address: string,
    connectionid: string,
    name: string,
    label: string,
    typename: string,
    offset: number,
  }
}

export class Link extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFrequency(): number;
  setFrequency(value: number): void;

  getSource(): Tag | undefined;
  setSource(value?: Tag): void;
  hasSource(): boolean;
  clearSource(): void;

  getTarget(): Tag | undefined;
  setTarget(value?: Tag): void;
  hasTarget(): boolean;
  clearTarget(): void;

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
  setLnk(value?: Link): void;
  hasLnk(): boolean;
  clearLnk(): void;

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
  setResult(value: boolean): void;

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
  setId(value: string): void;

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
  setResult(value: boolean): void;

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
  setId(value: string): void;

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
  setResult(value: boolean): void;

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
  setId(value: string): void;

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
  setResult(value: boolean): void;

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

