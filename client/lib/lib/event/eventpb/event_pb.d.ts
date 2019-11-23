import * as jspb from "google-protobuf"

export class Event extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Event.AsObject;
  static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
  static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Event;
  static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
  export type AsObject = {
    name: string,
    data: Uint8Array | string,
  }
}

export class SubscribeRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscribeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SubscribeRequest): SubscribeRequest.AsObject;
  static serializeBinaryToWriter(message: SubscribeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubscribeRequest;
  static deserializeBinaryFromReader(message: SubscribeRequest, reader: jspb.BinaryReader): SubscribeRequest;
}

export namespace SubscribeRequest {
  export type AsObject = {
    name: string,
  }
}

export class SubscribeResponse extends jspb.Message {
  getEvt(): Event | undefined;
  setEvt(value?: Event): void;
  hasEvt(): boolean;
  clearEvt(): void;
  hasEvt(): boolean;

  getUuid(): string;
  setUuid(value: string): void;
  hasUuid(): boolean;

  getResultCase(): SubscribeResponse.ResultCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscribeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SubscribeResponse): SubscribeResponse.AsObject;
  static serializeBinaryToWriter(message: SubscribeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubscribeResponse;
  static deserializeBinaryFromReader(message: SubscribeResponse, reader: jspb.BinaryReader): SubscribeResponse;
}

export namespace SubscribeResponse {
  export type AsObject = {
    evt?: Event.AsObject,
    uuid: string,
  }

  export enum ResultCase { 
    RESULT_NOT_SET = 0,
    EVT = 1,
    UUID = 2,
  }
}

export class UnSubscribeRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getUuid(): string;
  setUuid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnSubscribeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UnSubscribeRequest): UnSubscribeRequest.AsObject;
  static serializeBinaryToWriter(message: UnSubscribeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnSubscribeRequest;
  static deserializeBinaryFromReader(message: UnSubscribeRequest, reader: jspb.BinaryReader): UnSubscribeRequest;
}

export namespace UnSubscribeRequest {
  export type AsObject = {
    name: string,
    uuid: string,
  }
}

export class UnSubscribeResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnSubscribeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UnSubscribeResponse): UnSubscribeResponse.AsObject;
  static serializeBinaryToWriter(message: UnSubscribeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnSubscribeResponse;
  static deserializeBinaryFromReader(message: UnSubscribeResponse, reader: jspb.BinaryReader): UnSubscribeResponse;
}

export namespace UnSubscribeResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class PublishRequest extends jspb.Message {
  getEvt(): Event | undefined;
  setEvt(value?: Event): void;
  hasEvt(): boolean;
  clearEvt(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublishRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PublishRequest): PublishRequest.AsObject;
  static serializeBinaryToWriter(message: PublishRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublishRequest;
  static deserializeBinaryFromReader(message: PublishRequest, reader: jspb.BinaryReader): PublishRequest;
}

export namespace PublishRequest {
  export type AsObject = {
    evt?: Event.AsObject,
  }
}

export class PublishResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublishResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PublishResponse): PublishResponse.AsObject;
  static serializeBinaryToWriter(message: PublishResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublishResponse;
  static deserializeBinaryFromReader(message: PublishResponse, reader: jspb.BinaryReader): PublishResponse;
}

export namespace PublishResponse {
  export type AsObject = {
    result: boolean,
  }
}

