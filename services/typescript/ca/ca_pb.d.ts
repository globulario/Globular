import * as jspb from 'google-protobuf'



export class SignCertificateRequest extends jspb.Message {
  getCsr(): string;
  setCsr(value: string): SignCertificateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SignCertificateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SignCertificateRequest): SignCertificateRequest.AsObject;
  static serializeBinaryToWriter(message: SignCertificateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SignCertificateRequest;
  static deserializeBinaryFromReader(message: SignCertificateRequest, reader: jspb.BinaryReader): SignCertificateRequest;
}

export namespace SignCertificateRequest {
  export type AsObject = {
    csr: string,
  }
}

export class SignCertificateResponse extends jspb.Message {
  getCrt(): string;
  setCrt(value: string): SignCertificateResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SignCertificateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SignCertificateResponse): SignCertificateResponse.AsObject;
  static serializeBinaryToWriter(message: SignCertificateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SignCertificateResponse;
  static deserializeBinaryFromReader(message: SignCertificateResponse, reader: jspb.BinaryReader): SignCertificateResponse;
}

export namespace SignCertificateResponse {
  export type AsObject = {
    crt: string,
  }
}

export class GetCaCertificateRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCaCertificateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCaCertificateRequest): GetCaCertificateRequest.AsObject;
  static serializeBinaryToWriter(message: GetCaCertificateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCaCertificateRequest;
  static deserializeBinaryFromReader(message: GetCaCertificateRequest, reader: jspb.BinaryReader): GetCaCertificateRequest;
}

export namespace GetCaCertificateRequest {
  export type AsObject = {
  }
}

export class GetCaCertificateResponse extends jspb.Message {
  getCa(): string;
  setCa(value: string): GetCaCertificateResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCaCertificateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCaCertificateResponse): GetCaCertificateResponse.AsObject;
  static serializeBinaryToWriter(message: GetCaCertificateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCaCertificateResponse;
  static deserializeBinaryFromReader(message: GetCaCertificateResponse, reader: jspb.BinaryReader): GetCaCertificateResponse;
}

export namespace GetCaCertificateResponse {
  export type AsObject = {
    ca: string,
  }
}

