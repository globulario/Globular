import * as jspb from 'google-protobuf'



export class Event extends jspb.Message {
  getName(): string;
  setName(value: string): Event;

  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): Event;

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

export class QuitRequest extends jspb.Message {
  getUuid(): string;
  setUuid(value: string): QuitRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QuitRequest.AsObject;
  static toObject(includeInstance: boolean, msg: QuitRequest): QuitRequest.AsObject;
  static serializeBinaryToWriter(message: QuitRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QuitRequest;
  static deserializeBinaryFromReader(message: QuitRequest, reader: jspb.BinaryReader): QuitRequest;
}

export namespace QuitRequest {
  export type AsObject = {
    uuid: string,
  }
}

export class QuitResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): QuitResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QuitResponse.AsObject;
  static toObject(includeInstance: boolean, msg: QuitResponse): QuitResponse.AsObject;
  static serializeBinaryToWriter(message: QuitResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QuitResponse;
  static deserializeBinaryFromReader(message: QuitResponse, reader: jspb.BinaryReader): QuitResponse;
}

export namespace QuitResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class OnEventRequest extends jspb.Message {
  getUuid(): string;
  setUuid(value: string): OnEventRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OnEventRequest.AsObject;
  static toObject(includeInstance: boolean, msg: OnEventRequest): OnEventRequest.AsObject;
  static serializeBinaryToWriter(message: OnEventRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OnEventRequest;
  static deserializeBinaryFromReader(message: OnEventRequest, reader: jspb.BinaryReader): OnEventRequest;
}

export namespace OnEventRequest {
  export type AsObject = {
    uuid: string,
  }
}

export class OnEventResponse extends jspb.Message {
  getEvt(): Event | undefined;
  setEvt(value?: Event): OnEventResponse;
  hasEvt(): boolean;
  clearEvt(): OnEventResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OnEventResponse.AsObject;
  static toObject(includeInstance: boolean, msg: OnEventResponse): OnEventResponse.AsObject;
  static serializeBinaryToWriter(message: OnEventResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OnEventResponse;
  static deserializeBinaryFromReader(message: OnEventResponse, reader: jspb.BinaryReader): OnEventResponse;
}

export namespace OnEventResponse {
  export type AsObject = {
    evt?: Event.AsObject,
  }
}

export class SubscribeRequest extends jspb.Message {
  getName(): string;
  setName(value: string): SubscribeRequest;

  getUuid(): string;
  setUuid(value: string): SubscribeRequest;

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
    uuid: string,
  }
}

export class SubscribeResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): SubscribeResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscribeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SubscribeResponse): SubscribeResponse.AsObject;
  static serializeBinaryToWriter(message: SubscribeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubscribeResponse;
  static deserializeBinaryFromReader(message: SubscribeResponse, reader: jspb.BinaryReader): SubscribeResponse;
}

export namespace SubscribeResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class UnSubscribeRequest extends jspb.Message {
  getName(): string;
  setName(value: string): UnSubscribeRequest;

  getUuid(): string;
  setUuid(value: string): UnSubscribeRequest;

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
  setResult(value: boolean): UnSubscribeResponse;

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
  setEvt(value?: Event): PublishRequest;
  hasEvt(): boolean;
  clearEvt(): PublishRequest;

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
  setResult(value: boolean): PublishResponse;

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

