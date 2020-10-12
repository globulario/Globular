import * as grpcWeb from 'grpc-web';

import * as monitoring_pb from './monitoring_pb';


export class MonitoringServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: monitoring_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.StopResponse>;

  createConnection(
    request: monitoring_pb.CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.CreateConnectionRsp>;

  deleteConnection(
    request: monitoring_pb.DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.DeleteConnectionRsp>;

  alerts(
    request: monitoring_pb.AlertsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.AlertsResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.AlertsResponse>;

  alertManagers(
    request: monitoring_pb.AlertManagersRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.AlertManagersResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.AlertManagersResponse>;

  cleanTombstones(
    request: monitoring_pb.CleanTombstonesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.CleanTombstonesResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.CleanTombstonesResponse>;

  config(
    request: monitoring_pb.ConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.ConfigResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.ConfigResponse>;

  deleteSeries(
    request: monitoring_pb.DeleteSeriesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.DeleteSeriesResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.DeleteSeriesResponse>;

  flags(
    request: monitoring_pb.FlagsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.FlagsResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.FlagsResponse>;

  labelNames(
    request: monitoring_pb.LabelNamesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.LabelNamesResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.LabelNamesResponse>;

  labelValues(
    request: monitoring_pb.LabelValuesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.LabelValuesResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.LabelValuesResponse>;

  query(
    request: monitoring_pb.QueryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.QueryResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.QueryResponse>;

  queryRange(
    request: monitoring_pb.QueryRangeRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<monitoring_pb.QueryRangeResponse>;

  series(
    request: monitoring_pb.SeriesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.SeriesResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.SeriesResponse>;

  snapshot(
    request: monitoring_pb.SnapshotRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.SnapshotResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.SnapshotResponse>;

  rules(
    request: monitoring_pb.RulesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.RulesResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.RulesResponse>;

  targets(
    request: monitoring_pb.TargetsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.TargetsResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.TargetsResponse>;

  targetsMetadata(
    request: monitoring_pb.TargetsMetadataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: monitoring_pb.TargetsMetadataResponse) => void
  ): grpcWeb.ClientReadableStream<monitoring_pb.TargetsMetadataResponse>;

}

export class MonitoringServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: monitoring_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.StopResponse>;

  createConnection(
    request: monitoring_pb.CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.CreateConnectionRsp>;

  deleteConnection(
    request: monitoring_pb.DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.DeleteConnectionRsp>;

  alerts(
    request: monitoring_pb.AlertsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.AlertsResponse>;

  alertManagers(
    request: monitoring_pb.AlertManagersRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.AlertManagersResponse>;

  cleanTombstones(
    request: monitoring_pb.CleanTombstonesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.CleanTombstonesResponse>;

  config(
    request: monitoring_pb.ConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.ConfigResponse>;

  deleteSeries(
    request: monitoring_pb.DeleteSeriesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.DeleteSeriesResponse>;

  flags(
    request: monitoring_pb.FlagsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.FlagsResponse>;

  labelNames(
    request: monitoring_pb.LabelNamesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.LabelNamesResponse>;

  labelValues(
    request: monitoring_pb.LabelValuesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.LabelValuesResponse>;

  query(
    request: monitoring_pb.QueryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.QueryResponse>;

  queryRange(
    request: monitoring_pb.QueryRangeRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<monitoring_pb.QueryRangeResponse>;

  series(
    request: monitoring_pb.SeriesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.SeriesResponse>;

  snapshot(
    request: monitoring_pb.SnapshotRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.SnapshotResponse>;

  rules(
    request: monitoring_pb.RulesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.RulesResponse>;

  targets(
    request: monitoring_pb.TargetsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.TargetsResponse>;

  targetsMetadata(
    request: monitoring_pb.TargetsMetadataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<monitoring_pb.TargetsMetadataResponse>;

}

