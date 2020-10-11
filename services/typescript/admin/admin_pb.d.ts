import * as jspb from 'google-protobuf'



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
  setResult(value: string): GetConfigResponse;

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
  setConfig(value: string): SaveConfigRequest;

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
  setResult(value: string): SaveConfigResponse;

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
  setServiceId(value: string): StopServiceRequest;

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
  setResult(value: boolean): StopServiceResponse;

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
  setServiceId(value: string): StartServiceRequest;

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
  setServicePid(value: number): StartServiceResponse;

  getProxyPid(): number;
  setProxyPid(value: number): StartServiceResponse;

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
  setServiceId(value: string): RegisterExternalApplicationRequest;

  getPath(): string;
  setPath(value: string): RegisterExternalApplicationRequest;

  getArgsList(): Array<string>;
  setArgsList(value: Array<string>): RegisterExternalApplicationRequest;
  clearArgsList(): RegisterExternalApplicationRequest;
  addArgs(value: string, index?: number): RegisterExternalApplicationRequest;

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
  setServicePid(value: number): RegisterExternalApplicationResponse;

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
  setOldpassword(value: string): SetRootPasswordRequest;

  getNewpassword(): string;
  setNewpassword(value: string): SetRootPasswordRequest;

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
  setToken(value: string): SetRootPasswordResponse;

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

export class SetRootEmailRequest extends jspb.Message {
  getOldemail(): string;
  setOldemail(value: string): SetRootEmailRequest;

  getNewemail(): string;
  setNewemail(value: string): SetRootEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRootEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetRootEmailRequest): SetRootEmailRequest.AsObject;
  static serializeBinaryToWriter(message: SetRootEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRootEmailRequest;
  static deserializeBinaryFromReader(message: SetRootEmailRequest, reader: jspb.BinaryReader): SetRootEmailRequest;
}

export namespace SetRootEmailRequest {
  export type AsObject = {
    oldemail: string,
    newemail: string,
  }
}

export class SetRootEmailResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): SetRootEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRootEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetRootEmailResponse): SetRootEmailResponse.AsObject;
  static serializeBinaryToWriter(message: SetRootEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRootEmailResponse;
  static deserializeBinaryFromReader(message: SetRootEmailResponse, reader: jspb.BinaryReader): SetRootEmailResponse;
}

export namespace SetRootEmailResponse {
  export type AsObject = {
    token: string,
  }
}

export class SetPasswordRequest extends jspb.Message {
  getAccountid(): string;
  setAccountid(value: string): SetPasswordRequest;

  getOldpassword(): string;
  setOldpassword(value: string): SetPasswordRequest;

  getNewpassword(): string;
  setNewpassword(value: string): SetPasswordRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetPasswordRequest): SetPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: SetPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPasswordRequest;
  static deserializeBinaryFromReader(message: SetPasswordRequest, reader: jspb.BinaryReader): SetPasswordRequest;
}

export namespace SetPasswordRequest {
  export type AsObject = {
    accountid: string,
    oldpassword: string,
    newpassword: string,
  }
}

export class SetPasswordResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): SetPasswordResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetPasswordResponse): SetPasswordResponse.AsObject;
  static serializeBinaryToWriter(message: SetPasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPasswordResponse;
  static deserializeBinaryFromReader(message: SetPasswordResponse, reader: jspb.BinaryReader): SetPasswordResponse;
}

export namespace SetPasswordResponse {
  export type AsObject = {
    token: string,
  }
}

export class SetEmailRequest extends jspb.Message {
  getAccountid(): string;
  setAccountid(value: string): SetEmailRequest;

  getOldemail(): string;
  setOldemail(value: string): SetEmailRequest;

  getNewemail(): string;
  setNewemail(value: string): SetEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetEmailRequest): SetEmailRequest.AsObject;
  static serializeBinaryToWriter(message: SetEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetEmailRequest;
  static deserializeBinaryFromReader(message: SetEmailRequest, reader: jspb.BinaryReader): SetEmailRequest;
}

export namespace SetEmailRequest {
  export type AsObject = {
    accountid: string,
    oldemail: string,
    newemail: string,
  }
}

export class SetEmailResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): SetEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetEmailResponse): SetEmailResponse.AsObject;
  static serializeBinaryToWriter(message: SetEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetEmailResponse;
  static deserializeBinaryFromReader(message: SetEmailResponse, reader: jspb.BinaryReader): SetEmailResponse;
}

export namespace SetEmailResponse {
  export type AsObject = {
    token: string,
  }
}

export class HasRunningProcessRequest extends jspb.Message {
  getName(): string;
  setName(value: string): HasRunningProcessRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HasRunningProcessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: HasRunningProcessRequest): HasRunningProcessRequest.AsObject;
  static serializeBinaryToWriter(message: HasRunningProcessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HasRunningProcessRequest;
  static deserializeBinaryFromReader(message: HasRunningProcessRequest, reader: jspb.BinaryReader): HasRunningProcessRequest;
}

export namespace HasRunningProcessRequest {
  export type AsObject = {
    name: string,
  }
}

export class HasRunningProcessResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): HasRunningProcessResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HasRunningProcessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: HasRunningProcessResponse): HasRunningProcessResponse.AsObject;
  static serializeBinaryToWriter(message: HasRunningProcessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HasRunningProcessResponse;
  static deserializeBinaryFromReader(message: HasRunningProcessResponse, reader: jspb.BinaryReader): HasRunningProcessResponse;
}

export namespace HasRunningProcessResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class PublishServiceRequest extends jspb.Message {
  getServiceid(): string;
  setServiceid(value: string): PublishServiceRequest;

  getPublisherid(): string;
  setPublisherid(value: string): PublishServiceRequest;

  getPath(): string;
  setPath(value: string): PublishServiceRequest;

  getDicorveryid(): string;
  setDicorveryid(value: string): PublishServiceRequest;

  getRepositoryid(): string;
  setRepositoryid(value: string): PublishServiceRequest;

  getDescription(): string;
  setDescription(value: string): PublishServiceRequest;

  getKeywordsList(): Array<string>;
  setKeywordsList(value: Array<string>): PublishServiceRequest;
  clearKeywordsList(): PublishServiceRequest;
  addKeywords(value: string, index?: number): PublishServiceRequest;

  getVersion(): string;
  setVersion(value: string): PublishServiceRequest;

  getPlatform(): Platform;
  setPlatform(value: Platform): PublishServiceRequest;

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
    publisherid: string,
    path: string,
    dicorveryid: string,
    repositoryid: string,
    description: string,
    keywordsList: Array<string>,
    version: string,
    platform: Platform,
  }
}

export class UploadServicePackageRequest extends jspb.Message {
  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): UploadServicePackageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadServicePackageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UploadServicePackageRequest): UploadServicePackageRequest.AsObject;
  static serializeBinaryToWriter(message: UploadServicePackageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadServicePackageRequest;
  static deserializeBinaryFromReader(message: UploadServicePackageRequest, reader: jspb.BinaryReader): UploadServicePackageRequest;
}

export namespace UploadServicePackageRequest {
  export type AsObject = {
    data: Uint8Array | string,
  }
}

export class UploadServicePackageResponse extends jspb.Message {
  getPath(): string;
  setPath(value: string): UploadServicePackageResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadServicePackageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UploadServicePackageResponse): UploadServicePackageResponse.AsObject;
  static serializeBinaryToWriter(message: UploadServicePackageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadServicePackageResponse;
  static deserializeBinaryFromReader(message: UploadServicePackageResponse, reader: jspb.BinaryReader): UploadServicePackageResponse;
}

export namespace UploadServicePackageResponse {
  export type AsObject = {
    path: string,
  }
}

export class PublishServiceResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): PublishServiceResponse;

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
  setDicorveryid(value: string): InstallServiceRequest;

  getServiceid(): string;
  setServiceid(value: string): InstallServiceRequest;

  getPublisherid(): string;
  setPublisherid(value: string): InstallServiceRequest;

  getVersion(): string;
  setVersion(value: string): InstallServiceRequest;

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
  setResult(value: boolean): InstallServiceResponse;

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
  setServiceid(value: string): UninstallServiceRequest;

  getPublisherid(): string;
  setPublisherid(value: string): UninstallServiceRequest;

  getVersion(): string;
  setVersion(value: string): UninstallServiceRequest;

  getDeletepermissions(): boolean;
  setDeletepermissions(value: boolean): UninstallServiceRequest;

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
    deletepermissions: boolean,
  }
}

export class UninstallServiceResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): UninstallServiceResponse;

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

export class DeployApplicationRequest extends jspb.Message {
  getName(): string;
  setName(value: string): DeployApplicationRequest;

  getDomain(): string;
  setDomain(value: string): DeployApplicationRequest;

  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): DeployApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeployApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeployApplicationRequest): DeployApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: DeployApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeployApplicationRequest;
  static deserializeBinaryFromReader(message: DeployApplicationRequest, reader: jspb.BinaryReader): DeployApplicationRequest;
}

export namespace DeployApplicationRequest {
  export type AsObject = {
    name: string,
    domain: string,
    data: Uint8Array | string,
  }
}

export class DeployApplicationResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): DeployApplicationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeployApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeployApplicationResponse): DeployApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: DeployApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeployApplicationResponse;
  static deserializeBinaryFromReader(message: DeployApplicationResponse, reader: jspb.BinaryReader): DeployApplicationResponse;
}

export namespace DeployApplicationResponse {
  export type AsObject = {
    result: boolean,
  }
}

export enum Platform { 
  LINUX32 = 0,
  LINUX64 = 1,
  WIN32 = 2,
  WIN64 = 3,
}
