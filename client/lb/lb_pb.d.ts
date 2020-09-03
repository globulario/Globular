import * as jspb from "google-protobuf"

export class ServerInfo extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDomain(): string;
  setDomain(value: string): void;

  getPort(): number;
  setPort(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerInfo.AsObject;
  static toObject(includeInstance: boolean, msg: ServerInfo): ServerInfo.AsObject;
  static serializeBinaryToWriter(message: ServerInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerInfo;
  static deserializeBinaryFromReader(message: ServerInfo, reader: jspb.BinaryReader): ServerInfo;
}

export namespace ServerInfo {
  export type AsObject = {
    id: string,
    name: string,
    domain: string,
    port: number,
  }
}

export class LoadInfo extends jspb.Message {
  getServerinfo(): ServerInfo | undefined;
  setServerinfo(value?: ServerInfo): void;
  hasServerinfo(): boolean;
  clearServerinfo(): void;

  getLoad1(): number;
  setLoad1(value: number): void;

  getLoad5(): number;
  setLoad5(value: number): void;

  getLoad15(): number;
  setLoad15(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoadInfo.AsObject;
  static toObject(includeInstance: boolean, msg: LoadInfo): LoadInfo.AsObject;
  static serializeBinaryToWriter(message: LoadInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoadInfo;
  static deserializeBinaryFromReader(message: LoadInfo, reader: jspb.BinaryReader): LoadInfo;
}

export namespace LoadInfo {
  export type AsObject = {
    serverinfo?: ServerInfo.AsObject,
    load1: number,
    load5: number,
    load15: number,
  }
}

export class GetCanditatesRequest extends jspb.Message {
  getServicename(): string;
  setServicename(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCanditatesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCanditatesRequest): GetCanditatesRequest.AsObject;
  static serializeBinaryToWriter(message: GetCanditatesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCanditatesRequest;
  static deserializeBinaryFromReader(message: GetCanditatesRequest, reader: jspb.BinaryReader): GetCanditatesRequest;
}

export namespace GetCanditatesRequest {
  export type AsObject = {
    servicename: string,
  }
}

export class GetCanditatesResponse extends jspb.Message {
  getServersList(): Array<ServerInfo>;
  setServersList(value: Array<ServerInfo>): void;
  clearServersList(): void;
  addServers(value?: ServerInfo, index?: number): ServerInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCanditatesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCanditatesResponse): GetCanditatesResponse.AsObject;
  static serializeBinaryToWriter(message: GetCanditatesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCanditatesResponse;
  static deserializeBinaryFromReader(message: GetCanditatesResponse, reader: jspb.BinaryReader): GetCanditatesResponse;
}

export namespace GetCanditatesResponse {
  export type AsObject = {
    serversList: Array<ServerInfo.AsObject>,
  }
}

export class ReportLoadInfoRequest extends jspb.Message {
  getInfo(): LoadInfo | undefined;
  setInfo(value?: LoadInfo): void;
  hasInfo(): boolean;
  clearInfo(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReportLoadInfoRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReportLoadInfoRequest): ReportLoadInfoRequest.AsObject;
  static serializeBinaryToWriter(message: ReportLoadInfoRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReportLoadInfoRequest;
  static deserializeBinaryFromReader(message: ReportLoadInfoRequest, reader: jspb.BinaryReader): ReportLoadInfoRequest;
}

export namespace ReportLoadInfoRequest {
  export type AsObject = {
    info?: LoadInfo.AsObject,
  }
}

export class ReportLoadInfoResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReportLoadInfoResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReportLoadInfoResponse): ReportLoadInfoResponse.AsObject;
  static serializeBinaryToWriter(message: ReportLoadInfoResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReportLoadInfoResponse;
  static deserializeBinaryFromReader(message: ReportLoadInfoResponse, reader: jspb.BinaryReader): ReportLoadInfoResponse;
}

export namespace ReportLoadInfoResponse {
  export type AsObject = {
  }
}

