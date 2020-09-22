import * as jspb from "google-protobuf"

export class CreateAnalyseRqst extends jspb.Message {
  getData(): string;
  setData(value: string): CreateAnalyseRqst;

  getTolzon(): number;
  setTolzon(value: number): CreateAnalyseRqst;

  getLotol(): number;
  setLotol(value: number): CreateAnalyseRqst;

  getUptol(): number;
  setUptol(value: number): CreateAnalyseRqst;

  getToltype(): string;
  setToltype(value: string): CreateAnalyseRqst;

  getIspopulation(): boolean;
  setIspopulation(value: boolean): CreateAnalyseRqst;

  getTests(): string;
  setTests(value: string): CreateAnalyseRqst;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAnalyseRqst.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAnalyseRqst): CreateAnalyseRqst.AsObject;
  static serializeBinaryToWriter(message: CreateAnalyseRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAnalyseRqst;
  static deserializeBinaryFromReader(message: CreateAnalyseRqst, reader: jspb.BinaryReader): CreateAnalyseRqst;
}

export namespace CreateAnalyseRqst {
  export type AsObject = {
    data: string,
    tolzon: number,
    lotol: number,
    uptol: number,
    toltype: string,
    ispopulation: boolean,
    tests: string,
  }
}

export class CreateAnalyseRsp extends jspb.Message {
  getResult(): string;
  setResult(value: string): CreateAnalyseRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAnalyseRsp.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAnalyseRsp): CreateAnalyseRsp.AsObject;
  static serializeBinaryToWriter(message: CreateAnalyseRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAnalyseRsp;
  static deserializeBinaryFromReader(message: CreateAnalyseRsp, reader: jspb.BinaryReader): CreateAnalyseRsp;
}

export namespace CreateAnalyseRsp {
  export type AsObject = {
    result: string,
  }
}

