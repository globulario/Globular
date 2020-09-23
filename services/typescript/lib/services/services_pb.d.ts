import * as jspb from "google-protobuf"

export class ServiceDescriptor extends jspb.Message {
  getId(): string;
  setId(value: string): ServiceDescriptor;

  getPublisherid(): string;
  setPublisherid(value: string): ServiceDescriptor;

  getVersion(): string;
  setVersion(value: string): ServiceDescriptor;

  getDescription(): string;
  setDescription(value: string): ServiceDescriptor;

  getRepositoriesList(): Array<string>;
  setRepositoriesList(value: Array<string>): ServiceDescriptor;
  clearRepositoriesList(): ServiceDescriptor;
  addRepositories(value: string, index?: number): ServiceDescriptor;

  getDiscoveriesList(): Array<string>;
  setDiscoveriesList(value: Array<string>): ServiceDescriptor;
  clearDiscoveriesList(): ServiceDescriptor;
  addDiscoveries(value: string, index?: number): ServiceDescriptor;

  getKeywordsList(): Array<string>;
  setKeywordsList(value: Array<string>): ServiceDescriptor;
  clearKeywordsList(): ServiceDescriptor;
  addKeywords(value: string, index?: number): ServiceDescriptor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServiceDescriptor.AsObject;
  static toObject(includeInstance: boolean, msg: ServiceDescriptor): ServiceDescriptor.AsObject;
  static serializeBinaryToWriter(message: ServiceDescriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServiceDescriptor;
  static deserializeBinaryFromReader(message: ServiceDescriptor, reader: jspb.BinaryReader): ServiceDescriptor;
}

export namespace ServiceDescriptor {
  export type AsObject = {
    id: string,
    publisherid: string,
    version: string,
    description: string,
    repositoriesList: Array<string>,
    discoveriesList: Array<string>,
    keywordsList: Array<string>,
  }
}

export class ServiceBundle extends jspb.Message {
  getDescriptor(): ServiceDescriptor | undefined;
  setDescriptor(value?: ServiceDescriptor): ServiceBundle;
  hasDescriptor(): boolean;
  clearDescriptor(): ServiceBundle;

  getBuildnumber(): string;
  setBuildnumber(value: string): ServiceBundle;

  getPlaform(): Platform;
  setPlaform(value: Platform): ServiceBundle;

  getBinairies(): Uint8Array | string;
  getBinairies_asU8(): Uint8Array;
  getBinairies_asB64(): string;
  setBinairies(value: Uint8Array | string): ServiceBundle;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServiceBundle.AsObject;
  static toObject(includeInstance: boolean, msg: ServiceBundle): ServiceBundle.AsObject;
  static serializeBinaryToWriter(message: ServiceBundle, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServiceBundle;
  static deserializeBinaryFromReader(message: ServiceBundle, reader: jspb.BinaryReader): ServiceBundle;
}

export namespace ServiceBundle {
  export type AsObject = {
    descriptor?: ServiceDescriptor.AsObject,
    buildnumber: string,
    plaform: Platform,
    binairies: Uint8Array | string,
  }
}

export class PublishServiceDescriptorRequest extends jspb.Message {
  getDescriptor(): ServiceDescriptor | undefined;
  setDescriptor(value?: ServiceDescriptor): PublishServiceDescriptorRequest;
  hasDescriptor(): boolean;
  clearDescriptor(): PublishServiceDescriptorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublishServiceDescriptorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PublishServiceDescriptorRequest): PublishServiceDescriptorRequest.AsObject;
  static serializeBinaryToWriter(message: PublishServiceDescriptorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublishServiceDescriptorRequest;
  static deserializeBinaryFromReader(message: PublishServiceDescriptorRequest, reader: jspb.BinaryReader): PublishServiceDescriptorRequest;
}

export namespace PublishServiceDescriptorRequest {
  export type AsObject = {
    descriptor?: ServiceDescriptor.AsObject,
  }
}

export class PublishServiceDescriptorResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): PublishServiceDescriptorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublishServiceDescriptorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PublishServiceDescriptorResponse): PublishServiceDescriptorResponse.AsObject;
  static serializeBinaryToWriter(message: PublishServiceDescriptorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublishServiceDescriptorResponse;
  static deserializeBinaryFromReader(message: PublishServiceDescriptorResponse, reader: jspb.BinaryReader): PublishServiceDescriptorResponse;
}

export namespace PublishServiceDescriptorResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class UploadBundleRequest extends jspb.Message {
  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): UploadBundleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadBundleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UploadBundleRequest): UploadBundleRequest.AsObject;
  static serializeBinaryToWriter(message: UploadBundleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadBundleRequest;
  static deserializeBinaryFromReader(message: UploadBundleRequest, reader: jspb.BinaryReader): UploadBundleRequest;
}

export namespace UploadBundleRequest {
  export type AsObject = {
    data: Uint8Array | string,
  }
}

export class UploadBundleResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): UploadBundleResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadBundleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UploadBundleResponse): UploadBundleResponse.AsObject;
  static serializeBinaryToWriter(message: UploadBundleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadBundleResponse;
  static deserializeBinaryFromReader(message: UploadBundleResponse, reader: jspb.BinaryReader): UploadBundleResponse;
}

export namespace UploadBundleResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DownloadBundleRequest extends jspb.Message {
  getDescriptor(): ServiceDescriptor | undefined;
  setDescriptor(value?: ServiceDescriptor): DownloadBundleRequest;
  hasDescriptor(): boolean;
  clearDescriptor(): DownloadBundleRequest;

  getPlaform(): Platform;
  setPlaform(value: Platform): DownloadBundleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DownloadBundleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DownloadBundleRequest): DownloadBundleRequest.AsObject;
  static serializeBinaryToWriter(message: DownloadBundleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DownloadBundleRequest;
  static deserializeBinaryFromReader(message: DownloadBundleRequest, reader: jspb.BinaryReader): DownloadBundleRequest;
}

export namespace DownloadBundleRequest {
  export type AsObject = {
    descriptor?: ServiceDescriptor.AsObject,
    plaform: Platform,
  }
}

export class DownloadBundleResponse extends jspb.Message {
  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): DownloadBundleResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DownloadBundleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DownloadBundleResponse): DownloadBundleResponse.AsObject;
  static serializeBinaryToWriter(message: DownloadBundleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DownloadBundleResponse;
  static deserializeBinaryFromReader(message: DownloadBundleResponse, reader: jspb.BinaryReader): DownloadBundleResponse;
}

export namespace DownloadBundleResponse {
  export type AsObject = {
    data: Uint8Array | string,
  }
}

export class GetServiceDescriptorRequest extends jspb.Message {
  getServiceid(): string;
  setServiceid(value: string): GetServiceDescriptorRequest;

  getPublisherid(): string;
  setPublisherid(value: string): GetServiceDescriptorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetServiceDescriptorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetServiceDescriptorRequest): GetServiceDescriptorRequest.AsObject;
  static serializeBinaryToWriter(message: GetServiceDescriptorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetServiceDescriptorRequest;
  static deserializeBinaryFromReader(message: GetServiceDescriptorRequest, reader: jspb.BinaryReader): GetServiceDescriptorRequest;
}

export namespace GetServiceDescriptorRequest {
  export type AsObject = {
    serviceid: string,
    publisherid: string,
  }
}

export class GetServiceDescriptorResponse extends jspb.Message {
  getResultsList(): Array<ServiceDescriptor>;
  setResultsList(value: Array<ServiceDescriptor>): GetServiceDescriptorResponse;
  clearResultsList(): GetServiceDescriptorResponse;
  addResults(value?: ServiceDescriptor, index?: number): ServiceDescriptor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetServiceDescriptorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetServiceDescriptorResponse): GetServiceDescriptorResponse.AsObject;
  static serializeBinaryToWriter(message: GetServiceDescriptorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetServiceDescriptorResponse;
  static deserializeBinaryFromReader(message: GetServiceDescriptorResponse, reader: jspb.BinaryReader): GetServiceDescriptorResponse;
}

export namespace GetServiceDescriptorResponse {
  export type AsObject = {
    resultsList: Array<ServiceDescriptor.AsObject>,
  }
}

export class GetServicesDescriptorRequest extends jspb.Message {
  getQuery(): string;
  setQuery(value: string): GetServicesDescriptorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetServicesDescriptorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetServicesDescriptorRequest): GetServicesDescriptorRequest.AsObject;
  static serializeBinaryToWriter(message: GetServicesDescriptorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetServicesDescriptorRequest;
  static deserializeBinaryFromReader(message: GetServicesDescriptorRequest, reader: jspb.BinaryReader): GetServicesDescriptorRequest;
}

export namespace GetServicesDescriptorRequest {
  export type AsObject = {
    query: string,
  }
}

export class GetServicesDescriptorResponse extends jspb.Message {
  getResultsList(): Array<ServiceDescriptor>;
  setResultsList(value: Array<ServiceDescriptor>): GetServicesDescriptorResponse;
  clearResultsList(): GetServicesDescriptorResponse;
  addResults(value?: ServiceDescriptor, index?: number): ServiceDescriptor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetServicesDescriptorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetServicesDescriptorResponse): GetServicesDescriptorResponse.AsObject;
  static serializeBinaryToWriter(message: GetServicesDescriptorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetServicesDescriptorResponse;
  static deserializeBinaryFromReader(message: GetServicesDescriptorResponse, reader: jspb.BinaryReader): GetServicesDescriptorResponse;
}

export namespace GetServicesDescriptorResponse {
  export type AsObject = {
    resultsList: Array<ServiceDescriptor.AsObject>,
  }
}

export class SetServiceDescriptorRequest extends jspb.Message {
  getDescriptor(): ServiceDescriptor | undefined;
  setDescriptor(value?: ServiceDescriptor): SetServiceDescriptorRequest;
  hasDescriptor(): boolean;
  clearDescriptor(): SetServiceDescriptorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetServiceDescriptorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetServiceDescriptorRequest): SetServiceDescriptorRequest.AsObject;
  static serializeBinaryToWriter(message: SetServiceDescriptorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetServiceDescriptorRequest;
  static deserializeBinaryFromReader(message: SetServiceDescriptorRequest, reader: jspb.BinaryReader): SetServiceDescriptorRequest;
}

export namespace SetServiceDescriptorRequest {
  export type AsObject = {
    descriptor?: ServiceDescriptor.AsObject,
  }
}

export class SetServiceDescriptorResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): SetServiceDescriptorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetServiceDescriptorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetServiceDescriptorResponse): SetServiceDescriptorResponse.AsObject;
  static serializeBinaryToWriter(message: SetServiceDescriptorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetServiceDescriptorResponse;
  static deserializeBinaryFromReader(message: SetServiceDescriptorResponse, reader: jspb.BinaryReader): SetServiceDescriptorResponse;
}

export namespace SetServiceDescriptorResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class FindServicesDescriptorRequest extends jspb.Message {
  getKeywordsList(): Array<string>;
  setKeywordsList(value: Array<string>): FindServicesDescriptorRequest;
  clearKeywordsList(): FindServicesDescriptorRequest;
  addKeywords(value: string, index?: number): FindServicesDescriptorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindServicesDescriptorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FindServicesDescriptorRequest): FindServicesDescriptorRequest.AsObject;
  static serializeBinaryToWriter(message: FindServicesDescriptorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindServicesDescriptorRequest;
  static deserializeBinaryFromReader(message: FindServicesDescriptorRequest, reader: jspb.BinaryReader): FindServicesDescriptorRequest;
}

export namespace FindServicesDescriptorRequest {
  export type AsObject = {
    keywordsList: Array<string>,
  }
}

export class FindServicesDescriptorResponse extends jspb.Message {
  getResultsList(): Array<ServiceDescriptor>;
  setResultsList(value: Array<ServiceDescriptor>): FindServicesDescriptorResponse;
  clearResultsList(): FindServicesDescriptorResponse;
  addResults(value?: ServiceDescriptor, index?: number): ServiceDescriptor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindServicesDescriptorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FindServicesDescriptorResponse): FindServicesDescriptorResponse.AsObject;
  static serializeBinaryToWriter(message: FindServicesDescriptorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindServicesDescriptorResponse;
  static deserializeBinaryFromReader(message: FindServicesDescriptorResponse, reader: jspb.BinaryReader): FindServicesDescriptorResponse;
}

export namespace FindServicesDescriptorResponse {
  export type AsObject = {
    resultsList: Array<ServiceDescriptor.AsObject>,
  }
}

export enum Platform { 
  LINUX32 = 0,
  LINUX64 = 1,
  WIN32 = 2,
  WIN64 = 3,
}
