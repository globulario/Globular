import * as jspb from "google-protobuf"

export class GetConfigRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetConfigRequest): GetConfigRequest.AsObject;
  static serializeBinaryToWriter(message: GetConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConfigRequest;
  static deserializeBinaryFromReader(message: GetConfigRequest, reader: jspb.BinaryReader): GetConfigRequest;
}

export namespace GetConfigRequest {
  export type AsObject = {
  }
}

export class GetConfigResponse extends jspb.Message {
  getResult(): string;
  setResult(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetConfigResponse): GetConfigResponse.AsObject;
  static serializeBinaryToWriter(message: GetConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConfigResponse;
  static deserializeBinaryFromReader(message: GetConfigResponse, reader: jspb.BinaryReader): GetConfigResponse;
}

export namespace GetConfigResponse {
  export type AsObject = {
    result: string,
  }
}

export class SaveConfigRequest extends jspb.Message {
  getConfig(): string;
  setConfig(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveConfigRequest): SaveConfigRequest.AsObject;
  static serializeBinaryToWriter(message: SaveConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveConfigRequest;
  static deserializeBinaryFromReader(message: SaveConfigRequest, reader: jspb.BinaryReader): SaveConfigRequest;
}

export namespace SaveConfigRequest {
  export type AsObject = {
    config: string,
  }
}

export class SaveConfigResponse extends jspb.Message {
  getResult(): string;
  setResult(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveConfigResponse): SaveConfigResponse.AsObject;
  static serializeBinaryToWriter(message: SaveConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveConfigResponse;
  static deserializeBinaryFromReader(message: SaveConfigResponse, reader: jspb.BinaryReader): SaveConfigResponse;
}

export namespace SaveConfigResponse {
  export type AsObject = {
    result: string,
  }
}

export class StopServiceRequest extends jspb.Message {
  getServiceId(): string;
  setServiceId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StopServiceRequest): StopServiceRequest.AsObject;
  static serializeBinaryToWriter(message: StopServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopServiceRequest;
  static deserializeBinaryFromReader(message: StopServiceRequest, reader: jspb.BinaryReader): StopServiceRequest;
}

export namespace StopServiceRequest {
  export type AsObject = {
    serviceId: string,
  }
}

export class StopServiceResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StopServiceResponse): StopServiceResponse.AsObject;
  static serializeBinaryToWriter(message: StopServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopServiceResponse;
  static deserializeBinaryFromReader(message: StopServiceResponse, reader: jspb.BinaryReader): StopServiceResponse;
}

export namespace StopServiceResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class StartServiceRequest extends jspb.Message {
  getServiceId(): string;
  setServiceId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartServiceRequest): StartServiceRequest.AsObject;
  static serializeBinaryToWriter(message: StartServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartServiceRequest;
  static deserializeBinaryFromReader(message: StartServiceRequest, reader: jspb.BinaryReader): StartServiceRequest;
}

export namespace StartServiceRequest {
  export type AsObject = {
    serviceId: string,
  }
}

export class StartServiceResponse extends jspb.Message {
  getServicePid(): number;
  setServicePid(value: number): void;

  getProxyPid(): number;
  setProxyPid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StartServiceResponse): StartServiceResponse.AsObject;
  static serializeBinaryToWriter(message: StartServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartServiceResponse;
  static deserializeBinaryFromReader(message: StartServiceResponse, reader: jspb.BinaryReader): StartServiceResponse;
}

export namespace StartServiceResponse {
  export type AsObject = {
    servicePid: number,
    proxyPid: number,
  }
}

export class RegisterExternalApplicationRequest extends jspb.Message {
  getServiceId(): string;
  setServiceId(value: string): void;

  getPath(): string;
  setPath(value: string): void;

  getArgsList(): Array<string>;
  setArgsList(value: Array<string>): void;
  clearArgsList(): void;
  addArgs(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterExternalApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterExternalApplicationRequest): RegisterExternalApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterExternalApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterExternalApplicationRequest;
  static deserializeBinaryFromReader(message: RegisterExternalApplicationRequest, reader: jspb.BinaryReader): RegisterExternalApplicationRequest;
}

export namespace RegisterExternalApplicationRequest {
  export type AsObject = {
    serviceId: string,
    path: string,
    argsList: Array<string>,
  }
}

export class RegisterExternalApplicationResponse extends jspb.Message {
  getServicePid(): number;
  setServicePid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterExternalApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterExternalApplicationResponse): RegisterExternalApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterExternalApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterExternalApplicationResponse;
  static deserializeBinaryFromReader(message: RegisterExternalApplicationResponse, reader: jspb.BinaryReader): RegisterExternalApplicationResponse;
}

export namespace RegisterExternalApplicationResponse {
  export type AsObject = {
    servicePid: number,
  }
}

export class SetRootPasswordRqst extends jspb.Message {
  getOldpassword(): string;
  setOldpassword(value: string): void;

  getNewpassword(): string;
  setNewpassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRootPasswordRqst.AsObject;
  static toObject(includeInstance: boolean, msg: SetRootPasswordRqst): SetRootPasswordRqst.AsObject;
  static serializeBinaryToWriter(message: SetRootPasswordRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRootPasswordRqst;
  static deserializeBinaryFromReader(message: SetRootPasswordRqst, reader: jspb.BinaryReader): SetRootPasswordRqst;
}

export namespace SetRootPasswordRqst {
  export type AsObject = {
    oldpassword: string,
    newpassword: string,
  }
}

export class SetRootPasswordRsp extends jspb.Message {
  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRootPasswordRsp.AsObject;
  static toObject(includeInstance: boolean, msg: SetRootPasswordRsp): SetRootPasswordRsp.AsObject;
  static serializeBinaryToWriter(message: SetRootPasswordRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRootPasswordRsp;
  static deserializeBinaryFromReader(message: SetRootPasswordRsp, reader: jspb.BinaryReader): SetRootPasswordRsp;
}

export namespace SetRootPasswordRsp {
  export type AsObject = {
    token: string,
  }
}

