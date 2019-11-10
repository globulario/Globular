import * as grpcWeb from 'grpc-web';

import {
  AlertManagersRequest,
  AlertManagersResponse,
  AlertsRequest,
  AlertsResponse,
  CleanTombstonesRequest,
  CleanTombstonesResponse,
  ConfigRequest,
  ConfigResponse,
  CreateConnectionRqst,
  CreateConnectionRsp,
  DeleteConnectionRqst,
  DeleteConnectionRsp,
  DeleteSeriesRequest,
  DeleteSeriesResponse,
  FlagsRequest,
  FlagsResponse,
  LabelNamesRequest,
  LabelNamesResponse,
  LabelValuesRequest,
  LabelValuesResponse,
  QueryRangeRequest,
  QueryRangeResponse,
  QueryRequest,
  QueryResponse,
  RulesRequest,
  RulesResponse,
  SeriesRequest,
  SeriesResponse,
  SnapshotRequest,
  SnapshotResponse,
  TargetsMetadataRequest,
  TargetsMetadataResponse,
  TargetsRequest,
  TargetsResponse} from './monitoring_pb';

export class MonitoringServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  createConnection(
    request: CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteConnectionRsp>;

  alerts(
    request: AlertsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AlertsResponse) => void
  ): grpcWeb.ClientReadableStream<AlertsResponse>;

  alertManagers(
    request: AlertManagersRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AlertManagersResponse) => void
  ): grpcWeb.ClientReadableStream<AlertManagersResponse>;

  cleanTombstones(
    request: CleanTombstonesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CleanTombstonesResponse) => void
  ): grpcWeb.ClientReadableStream<CleanTombstonesResponse>;

  config(
    request: ConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ConfigResponse) => void
  ): grpcWeb.ClientReadableStream<ConfigResponse>;

  deleteSeries(
    request: DeleteSeriesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteSeriesResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteSeriesResponse>;

  flags(
    request: FlagsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: FlagsResponse) => void
  ): grpcWeb.ClientReadableStream<FlagsResponse>;

  labelNames(
    request: LabelNamesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: LabelNamesResponse) => void
  ): grpcWeb.ClientReadableStream<LabelNamesResponse>;

  labelValues(
    request: LabelValuesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: LabelValuesResponse) => void
  ): grpcWeb.ClientReadableStream<LabelValuesResponse>;

  query(
    request: QueryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: QueryResponse) => void
  ): grpcWeb.ClientReadableStream<QueryResponse>;

  queryRange(
    request: QueryRangeRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<QueryRangeResponse>;

  series(
    request: SeriesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SeriesResponse) => void
  ): grpcWeb.ClientReadableStream<SeriesResponse>;

  snapshot(
    request: SnapshotRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SnapshotResponse) => void
  ): grpcWeb.ClientReadableStream<SnapshotResponse>;

  rules(
    request: RulesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RulesResponse) => void
  ): grpcWeb.ClientReadableStream<RulesResponse>;

  targets(
    request: TargetsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: TargetsResponse) => void
  ): grpcWeb.ClientReadableStream<TargetsResponse>;

  targetsMetadata(
    request: TargetsMetadataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: TargetsMetadataResponse) => void
  ): grpcWeb.ClientReadableStream<TargetsMetadataResponse>;

}

export class MonitoringServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  createConnection(
    request: CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteConnectionRsp>;

  alerts(
    request: AlertsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<AlertsResponse>;

  alertManagers(
    request: AlertManagersRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<AlertManagersResponse>;

  cleanTombstones(
    request: CleanTombstonesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<CleanTombstonesResponse>;

  config(
    request: ConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ConfigResponse>;

  deleteSeries(
    request: DeleteSeriesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteSeriesResponse>;

  flags(
    request: FlagsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<FlagsResponse>;

  labelNames(
    request: LabelNamesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<LabelNamesResponse>;

  labelValues(
    request: LabelValuesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<LabelValuesResponse>;

  query(
    request: QueryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<QueryResponse>;

  queryRange(
    request: QueryRangeRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<QueryRangeResponse>;

  series(
    request: SeriesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SeriesResponse>;

  snapshot(
    request: SnapshotRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SnapshotResponse>;

  rules(
    request: RulesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<RulesResponse>;

  targets(
    request: TargetsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<TargetsResponse>;

  targetsMetadata(
    request: TargetsMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<TargetsMetadataResponse>;

}

