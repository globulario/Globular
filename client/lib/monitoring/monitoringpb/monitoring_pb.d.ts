import * as jspb from "google-protobuf"

export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getHost(): string;
  setHost(value: string): void;

  getStore(): StoreType;
  setStore(value: StoreType): void;

  getPort(): number;
  setPort(value: number): void;

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
    host: string,
    store: StoreType,
    port: number,
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

export class AlertsRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AlertsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AlertsRequest): AlertsRequest.AsObject;
  static serializeBinaryToWriter(message: AlertsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AlertsRequest;
  static deserializeBinaryFromReader(message: AlertsRequest, reader: jspb.BinaryReader): AlertsRequest;
}

export namespace AlertsRequest {
  export type AsObject = {
    connectionid: string,
  }
}

export class AlertsResponse extends jspb.Message {
  getResults(): string;
  setResults(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AlertsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AlertsResponse): AlertsResponse.AsObject;
  static serializeBinaryToWriter(message: AlertsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AlertsResponse;
  static deserializeBinaryFromReader(message: AlertsResponse, reader: jspb.BinaryReader): AlertsResponse;
}

export namespace AlertsResponse {
  export type AsObject = {
    results: string,
  }
}

export class AlertManagersRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AlertManagersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AlertManagersRequest): AlertManagersRequest.AsObject;
  static serializeBinaryToWriter(message: AlertManagersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AlertManagersRequest;
  static deserializeBinaryFromReader(message: AlertManagersRequest, reader: jspb.BinaryReader): AlertManagersRequest;
}

export namespace AlertManagersRequest {
  export type AsObject = {
    connectionid: string,
  }
}

export class AlertManagersResponse extends jspb.Message {
  getResults(): string;
  setResults(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AlertManagersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AlertManagersResponse): AlertManagersResponse.AsObject;
  static serializeBinaryToWriter(message: AlertManagersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AlertManagersResponse;
  static deserializeBinaryFromReader(message: AlertManagersResponse, reader: jspb.BinaryReader): AlertManagersResponse;
}

export namespace AlertManagersResponse {
  export type AsObject = {
    results: string,
  }
}

export class CleanTombstonesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CleanTombstonesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CleanTombstonesRequest): CleanTombstonesRequest.AsObject;
  static serializeBinaryToWriter(message: CleanTombstonesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CleanTombstonesRequest;
  static deserializeBinaryFromReader(message: CleanTombstonesRequest, reader: jspb.BinaryReader): CleanTombstonesRequest;
}

export namespace CleanTombstonesRequest {
  export type AsObject = {
    connectionid: string,
  }
}

export class CleanTombstonesResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CleanTombstonesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CleanTombstonesResponse): CleanTombstonesResponse.AsObject;
  static serializeBinaryToWriter(message: CleanTombstonesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CleanTombstonesResponse;
  static deserializeBinaryFromReader(message: CleanTombstonesResponse, reader: jspb.BinaryReader): CleanTombstonesResponse;
}

export namespace CleanTombstonesResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class ConfigRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigRequest): ConfigRequest.AsObject;
  static serializeBinaryToWriter(message: ConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigRequest;
  static deserializeBinaryFromReader(message: ConfigRequest, reader: jspb.BinaryReader): ConfigRequest;
}

export namespace ConfigRequest {
  export type AsObject = {
    connectionid: string,
  }
}

export class ConfigResponse extends jspb.Message {
  getResults(): string;
  setResults(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigResponse): ConfigResponse.AsObject;
  static serializeBinaryToWriter(message: ConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigResponse;
  static deserializeBinaryFromReader(message: ConfigResponse, reader: jspb.BinaryReader): ConfigResponse;
}

export namespace ConfigResponse {
  export type AsObject = {
    results: string,
  }
}

export class DeleteSeriesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getMatchesList(): Array<string>;
  setMatchesList(value: Array<string>): void;
  clearMatchesList(): void;
  addMatches(value: string, index?: number): void;

  getStarttime(): number;
  setStarttime(value: number): void;

  getEndtime(): number;
  setEndtime(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSeriesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSeriesRequest): DeleteSeriesRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteSeriesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSeriesRequest;
  static deserializeBinaryFromReader(message: DeleteSeriesRequest, reader: jspb.BinaryReader): DeleteSeriesRequest;
}

export namespace DeleteSeriesRequest {
  export type AsObject = {
    connectionid: string,
    matchesList: Array<string>,
    starttime: number,
    endtime: number,
  }
}

export class DeleteSeriesResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSeriesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSeriesResponse): DeleteSeriesResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteSeriesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSeriesResponse;
  static deserializeBinaryFromReader(message: DeleteSeriesResponse, reader: jspb.BinaryReader): DeleteSeriesResponse;
}

export namespace DeleteSeriesResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class FlagsRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FlagsRequest): FlagsRequest.AsObject;
  static serializeBinaryToWriter(message: FlagsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FlagsRequest;
  static deserializeBinaryFromReader(message: FlagsRequest, reader: jspb.BinaryReader): FlagsRequest;
}

export namespace FlagsRequest {
  export type AsObject = {
    connectionid: string,
  }
}

export class FlagsResponse extends jspb.Message {
  getResults(): string;
  setResults(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FlagsResponse): FlagsResponse.AsObject;
  static serializeBinaryToWriter(message: FlagsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FlagsResponse;
  static deserializeBinaryFromReader(message: FlagsResponse, reader: jspb.BinaryReader): FlagsResponse;
}

export namespace FlagsResponse {
  export type AsObject = {
    results: string,
  }
}

export class LabelNamesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LabelNamesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LabelNamesRequest): LabelNamesRequest.AsObject;
  static serializeBinaryToWriter(message: LabelNamesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LabelNamesRequest;
  static deserializeBinaryFromReader(message: LabelNamesRequest, reader: jspb.BinaryReader): LabelNamesRequest;
}

export namespace LabelNamesRequest {
  export type AsObject = {
    connectionid: string,
  }
}

export class LabelNamesResponse extends jspb.Message {
  getLabelsList(): Array<string>;
  setLabelsList(value: Array<string>): void;
  clearLabelsList(): void;
  addLabels(value: string, index?: number): void;

  getWarnings(): string;
  setWarnings(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LabelNamesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LabelNamesResponse): LabelNamesResponse.AsObject;
  static serializeBinaryToWriter(message: LabelNamesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LabelNamesResponse;
  static deserializeBinaryFromReader(message: LabelNamesResponse, reader: jspb.BinaryReader): LabelNamesResponse;
}

export namespace LabelNamesResponse {
  export type AsObject = {
    labelsList: Array<string>,
    warnings: string,
  }
}

export class LabelValuesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getLabel(): string;
  setLabel(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LabelValuesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LabelValuesRequest): LabelValuesRequest.AsObject;
  static serializeBinaryToWriter(message: LabelValuesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LabelValuesRequest;
  static deserializeBinaryFromReader(message: LabelValuesRequest, reader: jspb.BinaryReader): LabelValuesRequest;
}

export namespace LabelValuesRequest {
  export type AsObject = {
    connectionid: string,
    label: string,
  }
}

export class LabelValuesResponse extends jspb.Message {
  getLabelvalues(): string;
  setLabelvalues(value: string): void;

  getWarnings(): string;
  setWarnings(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LabelValuesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LabelValuesResponse): LabelValuesResponse.AsObject;
  static serializeBinaryToWriter(message: LabelValuesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LabelValuesResponse;
  static deserializeBinaryFromReader(message: LabelValuesResponse, reader: jspb.BinaryReader): LabelValuesResponse;
}

export namespace LabelValuesResponse {
  export type AsObject = {
    labelvalues: string,
    warnings: string,
  }
}

export class QueryRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getTs(): number;
  setTs(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: QueryRequest): QueryRequest.AsObject;
  static serializeBinaryToWriter(message: QueryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryRequest;
  static deserializeBinaryFromReader(message: QueryRequest, reader: jspb.BinaryReader): QueryRequest;
}

export namespace QueryRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    ts: number,
  }
}

export class QueryResponse extends jspb.Message {
  getValue(): string;
  setValue(value: string): void;

  getWarnings(): string;
  setWarnings(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: QueryResponse): QueryResponse.AsObject;
  static serializeBinaryToWriter(message: QueryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryResponse;
  static deserializeBinaryFromReader(message: QueryResponse, reader: jspb.BinaryReader): QueryResponse;
}

export namespace QueryResponse {
  export type AsObject = {
    value: string,
    warnings: string,
  }
}

export class QueryRangeRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getStarttime(): number;
  setStarttime(value: number): void;

  getEndtime(): number;
  setEndtime(value: number): void;

  getStep(): number;
  setStep(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryRangeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: QueryRangeRequest): QueryRangeRequest.AsObject;
  static serializeBinaryToWriter(message: QueryRangeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryRangeRequest;
  static deserializeBinaryFromReader(message: QueryRangeRequest, reader: jspb.BinaryReader): QueryRangeRequest;
}

export namespace QueryRangeRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    starttime: number,
    endtime: number,
    step: number,
  }
}

export class QueryRangeResponse extends jspb.Message {
  getValue(): string;
  setValue(value: string): void;

  getWarnings(): string;
  setWarnings(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryRangeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: QueryRangeResponse): QueryRangeResponse.AsObject;
  static serializeBinaryToWriter(message: QueryRangeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryRangeResponse;
  static deserializeBinaryFromReader(message: QueryRangeResponse, reader: jspb.BinaryReader): QueryRangeResponse;
}

export namespace QueryRangeResponse {
  export type AsObject = {
    value: string,
    warnings: string,
  }
}

export class SeriesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getMatchesList(): Array<string>;
  setMatchesList(value: Array<string>): void;
  clearMatchesList(): void;
  addMatches(value: string, index?: number): void;

  getStarttime(): number;
  setStarttime(value: number): void;

  getEndtime(): number;
  setEndtime(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SeriesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SeriesRequest): SeriesRequest.AsObject;
  static serializeBinaryToWriter(message: SeriesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SeriesRequest;
  static deserializeBinaryFromReader(message: SeriesRequest, reader: jspb.BinaryReader): SeriesRequest;
}

export namespace SeriesRequest {
  export type AsObject = {
    connectionid: string,
    matchesList: Array<string>,
    starttime: number,
    endtime: number,
  }
}

export class SeriesResponse extends jspb.Message {
  getLabelset(): string;
  setLabelset(value: string): void;

  getWarnings(): string;
  setWarnings(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SeriesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SeriesResponse): SeriesResponse.AsObject;
  static serializeBinaryToWriter(message: SeriesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SeriesResponse;
  static deserializeBinaryFromReader(message: SeriesResponse, reader: jspb.BinaryReader): SeriesResponse;
}

export namespace SeriesResponse {
  export type AsObject = {
    labelset: string,
    warnings: string,
  }
}

export class SnapshotRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getSkiphead(): boolean;
  setSkiphead(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SnapshotRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SnapshotRequest): SnapshotRequest.AsObject;
  static serializeBinaryToWriter(message: SnapshotRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SnapshotRequest;
  static deserializeBinaryFromReader(message: SnapshotRequest, reader: jspb.BinaryReader): SnapshotRequest;
}

export namespace SnapshotRequest {
  export type AsObject = {
    connectionid: string,
    skiphead: boolean,
  }
}

export class SnapshotResponse extends jspb.Message {
  getResult(): string;
  setResult(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SnapshotResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SnapshotResponse): SnapshotResponse.AsObject;
  static serializeBinaryToWriter(message: SnapshotResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SnapshotResponse;
  static deserializeBinaryFromReader(message: SnapshotResponse, reader: jspb.BinaryReader): SnapshotResponse;
}

export namespace SnapshotResponse {
  export type AsObject = {
    result: string,
  }
}

export class RulesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RulesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RulesRequest): RulesRequest.AsObject;
  static serializeBinaryToWriter(message: RulesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RulesRequest;
  static deserializeBinaryFromReader(message: RulesRequest, reader: jspb.BinaryReader): RulesRequest;
}

export namespace RulesRequest {
  export type AsObject = {
    connectionid: string,
  }
}

export class RulesResponse extends jspb.Message {
  getResult(): string;
  setResult(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RulesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RulesResponse): RulesResponse.AsObject;
  static serializeBinaryToWriter(message: RulesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RulesResponse;
  static deserializeBinaryFromReader(message: RulesResponse, reader: jspb.BinaryReader): RulesResponse;
}

export namespace RulesResponse {
  export type AsObject = {
    result: string,
  }
}

export class TargetsRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TargetsRequest): TargetsRequest.AsObject;
  static serializeBinaryToWriter(message: TargetsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TargetsRequest;
  static deserializeBinaryFromReader(message: TargetsRequest, reader: jspb.BinaryReader): TargetsRequest;
}

export namespace TargetsRequest {
  export type AsObject = {
    connectionid: string,
  }
}

export class TargetsResponse extends jspb.Message {
  getResult(): string;
  setResult(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TargetsResponse): TargetsResponse.AsObject;
  static serializeBinaryToWriter(message: TargetsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TargetsResponse;
  static deserializeBinaryFromReader(message: TargetsResponse, reader: jspb.BinaryReader): TargetsResponse;
}

export namespace TargetsResponse {
  export type AsObject = {
    result: string,
  }
}

export class TargetsMetadataRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getMatchtarget(): string;
  setMatchtarget(value: string): void;

  getMetric(): string;
  setMetric(value: string): void;

  getLimit(): string;
  setLimit(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetsMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TargetsMetadataRequest): TargetsMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: TargetsMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TargetsMetadataRequest;
  static deserializeBinaryFromReader(message: TargetsMetadataRequest, reader: jspb.BinaryReader): TargetsMetadataRequest;
}

export namespace TargetsMetadataRequest {
  export type AsObject = {
    connectionid: string,
    matchtarget: string,
    metric: string,
    limit: string,
  }
}

export class TargetsMetadataResponse extends jspb.Message {
  getResult(): string;
  setResult(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetsMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TargetsMetadataResponse): TargetsMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: TargetsMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TargetsMetadataResponse;
  static deserializeBinaryFromReader(message: TargetsMetadataResponse, reader: jspb.BinaryReader): TargetsMetadataResponse;
}

export namespace TargetsMetadataResponse {
  export type AsObject = {
    result: string,
  }
}

export enum StoreType { 
  PROMETHEUS = 0,
}
