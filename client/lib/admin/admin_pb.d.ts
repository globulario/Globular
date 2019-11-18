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

export class SetRootPasswordRequest extends jspb.Message {
  getOldpassword(): string;
  setOldpassword(value: string): void;

  getNewpassword(): string;
  setNewpassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRootPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetRootPasswordRequest): SetRootPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: SetRootPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRootPasswordRequest;
  static deserializeBinaryFromReader(message: SetRootPasswordRequest, reader: jspb.BinaryReader): SetRootPasswordRequest;
}

export namespace SetRootPasswordRequest {
  export type AsObject = {
    oldpassword: string,
    newpassword: string,
  }
}

export class SetRootPasswordResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRootPasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetRootPasswordResponse): SetRootPasswordResponse.AsObject;
  static serializeBinaryToWriter(message: SetRootPasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRootPasswordResponse;
  static deserializeBinaryFromReader(message: SetRootPasswordResponse, reader: jspb.BinaryReader): SetRootPasswordResponse;
}

export namespace SetRootPasswordResponse {
  export type AsObject = {
    token: string,
  }
}

export class PublishServiceRequest extends jspb.Message {
  getServiceid(): string;
  setServiceid(value: string): void;

  getDicorveryid(): string;
  setDicorveryid(value: string): void;

  getRepositoryid(): string;
  setRepositoryid(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getKeywordsList(): Array<string>;
  setKeywordsList(value: Array<string>): void;
  clearKeywordsList(): void;
  addKeywords(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublishServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PublishServiceRequest): PublishServiceRequest.AsObject;
  static serializeBinaryToWriter(message: PublishServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublishServiceRequest;
  static deserializeBinaryFromReader(message: PublishServiceRequest, reader: jspb.BinaryReader): PublishServiceRequest;
}

export namespace PublishServiceRequest {
  export type AsObject = {
    serviceid: string,
    dicorveryid: string,
    repositoryid: string,
    description: string,
    keywordsList: Array<string>,
  }
}

export class PublishServiceResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublishServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PublishServiceResponse): PublishServiceResponse.AsObject;
  static serializeBinaryToWriter(message: PublishServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublishServiceResponse;
  static deserializeBinaryFromReader(message: PublishServiceResponse, reader: jspb.BinaryReader): PublishServiceResponse;
}

export namespace PublishServiceResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class InstallServiceRequest extends jspb.Message {
  getDicorveryid(): string;
  setDicorveryid(value: string): void;

  getServiceid(): string;
  setServiceid(value: string): void;

  getPublisherid(): string;
  setPublisherid(value: string): void;

  getVersion(): string;
  setVersion(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InstallServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: InstallServiceRequest): InstallServiceRequest.AsObject;
  static serializeBinaryToWriter(message: InstallServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InstallServiceRequest;
  static deserializeBinaryFromReader(message: InstallServiceRequest, reader: jspb.BinaryReader): InstallServiceRequest;
}

export namespace InstallServiceRequest {
  export type AsObject = {
    dicorveryid: string,
    serviceid: string,
    publisherid: string,
    version: string,
  }
}

export class InstallServiceResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InstallServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: InstallServiceResponse): InstallServiceResponse.AsObject;
  static serializeBinaryToWriter(message: InstallServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InstallServiceResponse;
  static deserializeBinaryFromReader(message: InstallServiceResponse, reader: jspb.BinaryReader): InstallServiceResponse;
}

export namespace InstallServiceResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class UninstallServiceRequest extends jspb.Message {
  getServiceid(): string;
  setServiceid(value: string): void;

  getPublisherid(): string;
  setPublisherid(value: string): void;

  getVersion(): string;
  setVersion(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UninstallServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UninstallServiceRequest): UninstallServiceRequest.AsObject;
  static serializeBinaryToWriter(message: UninstallServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UninstallServiceRequest;
  static deserializeBinaryFromReader(message: UninstallServiceRequest, reader: jspb.BinaryReader): UninstallServiceRequest;
}

export namespace UninstallServiceRequest {
  export type AsObject = {
    serviceid: string,
    publisherid: string,
    version: string,
  }
}

export class UninstallServiceResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UninstallServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UninstallServiceResponse): UninstallServiceResponse.AsObject;
  static serializeBinaryToWriter(message: UninstallServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UninstallServiceResponse;
  static deserializeBinaryFromReader(message: UninstallServiceResponse, reader: jspb.BinaryReader): UninstallServiceResponse;
}

export namespace UninstallServiceResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeployedApplicationRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeployedApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeployedApplicationRequest): DeployedApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: DeployedApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeployedApplicationRequest;
  static deserializeBinaryFromReader(message: DeployedApplicationRequest, reader: jspb.BinaryReader): DeployedApplicationRequest;
}

export namespace DeployedApplicationRequest {
  export type AsObject = {
    name: string,
    data: Uint8Array | string,
  }
}

export class DeployedApplicationResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeployedApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeployedApplicationResponse): DeployedApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: DeployedApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeployedApplicationResponse;
  static deserializeBinaryFromReader(message: DeployedApplicationResponse, reader: jspb.BinaryReader): DeployedApplicationResponse;
}

export namespace DeployedApplicationResponse {
  export type AsObject = {
    result: boolean,
  }
}

