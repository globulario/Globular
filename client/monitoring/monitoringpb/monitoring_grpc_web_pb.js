/**
 * @fileoverview gRPC-Web generated client stub for monitoring
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.monitoring = require('./monitoring_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.monitoring.MonitoringServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.monitoring.MonitoringServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.CreateConnectionRqst,
 *   !proto.monitoring.CreateConnectionRsp>}
 */
const methodInfo_MonitoringService_CreateConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.CreateConnectionRsp,
  /** @param {!proto.monitoring.CreateConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.CreateConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.monitoring.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.CreateConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.CreateConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.createConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/CreateConnection',
      request,
      metadata || {},
      methodInfo_MonitoringService_CreateConnection,
      callback);
};


/**
 * @param {!proto.monitoring.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.CreateConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.createConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/CreateConnection',
      request,
      metadata || {},
      methodInfo_MonitoringService_CreateConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.DeleteConnectionRqst,
 *   !proto.monitoring.DeleteConnectionRsp>}
 */
const methodInfo_MonitoringService_DeleteConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.DeleteConnectionRsp,
  /** @param {!proto.monitoring.DeleteConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.DeleteConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.monitoring.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.DeleteConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.DeleteConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.deleteConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_MonitoringService_DeleteConnection,
      callback);
};


/**
 * @param {!proto.monitoring.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.DeleteConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.deleteConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_MonitoringService_DeleteConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.AlertsRequest,
 *   !proto.monitoring.AlertsResponse>}
 */
const methodInfo_MonitoringService_Alerts = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.AlertsResponse,
  /** @param {!proto.monitoring.AlertsRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.AlertsResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.AlertsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.AlertsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.AlertsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.alerts =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/Alerts',
      request,
      metadata || {},
      methodInfo_MonitoringService_Alerts,
      callback);
};


/**
 * @param {!proto.monitoring.AlertsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.AlertsResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.alerts =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/Alerts',
      request,
      metadata || {},
      methodInfo_MonitoringService_Alerts);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.AlertManagersRequest,
 *   !proto.monitoring.AlertManagersResponse>}
 */
const methodInfo_MonitoringService_AlertManagers = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.AlertManagersResponse,
  /** @param {!proto.monitoring.AlertManagersRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.AlertManagersResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.AlertManagersRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.AlertManagersResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.AlertManagersResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.alertManagers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/AlertManagers',
      request,
      metadata || {},
      methodInfo_MonitoringService_AlertManagers,
      callback);
};


/**
 * @param {!proto.monitoring.AlertManagersRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.AlertManagersResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.alertManagers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/AlertManagers',
      request,
      metadata || {},
      methodInfo_MonitoringService_AlertManagers);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.CleanTombstonesRequest,
 *   !proto.monitoring.CleanTombstonesResponse>}
 */
const methodInfo_MonitoringService_CleanTombstones = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.CleanTombstonesResponse,
  /** @param {!proto.monitoring.CleanTombstonesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.CleanTombstonesResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.CleanTombstonesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.CleanTombstonesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.CleanTombstonesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.cleanTombstones =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/CleanTombstones',
      request,
      metadata || {},
      methodInfo_MonitoringService_CleanTombstones,
      callback);
};


/**
 * @param {!proto.monitoring.CleanTombstonesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.CleanTombstonesResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.cleanTombstones =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/CleanTombstones',
      request,
      metadata || {},
      methodInfo_MonitoringService_CleanTombstones);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.ConfigRequest,
 *   !proto.monitoring.ConfigResponse>}
 */
const methodInfo_MonitoringService_Config = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.ConfigResponse,
  /** @param {!proto.monitoring.ConfigRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.ConfigResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.ConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.ConfigResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.ConfigResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.config =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/Config',
      request,
      metadata || {},
      methodInfo_MonitoringService_Config,
      callback);
};


/**
 * @param {!proto.monitoring.ConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.ConfigResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.config =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/Config',
      request,
      metadata || {},
      methodInfo_MonitoringService_Config);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.DeleteSeriesRequest,
 *   !proto.monitoring.DeleteSeriesResponse>}
 */
const methodInfo_MonitoringService_DeleteSeries = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.DeleteSeriesResponse,
  /** @param {!proto.monitoring.DeleteSeriesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.DeleteSeriesResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.DeleteSeriesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.DeleteSeriesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.DeleteSeriesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.deleteSeries =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/DeleteSeries',
      request,
      metadata || {},
      methodInfo_MonitoringService_DeleteSeries,
      callback);
};


/**
 * @param {!proto.monitoring.DeleteSeriesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.DeleteSeriesResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.deleteSeries =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/DeleteSeries',
      request,
      metadata || {},
      methodInfo_MonitoringService_DeleteSeries);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.FlagsRequest,
 *   !proto.monitoring.FlagsResponse>}
 */
const methodInfo_MonitoringService_Flags = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.FlagsResponse,
  /** @param {!proto.monitoring.FlagsRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.FlagsResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.FlagsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.FlagsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.FlagsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.flags =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/Flags',
      request,
      metadata || {},
      methodInfo_MonitoringService_Flags,
      callback);
};


/**
 * @param {!proto.monitoring.FlagsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.FlagsResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.flags =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/Flags',
      request,
      metadata || {},
      methodInfo_MonitoringService_Flags);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.LabelNamesRequest,
 *   !proto.monitoring.LabelNamesResponse>}
 */
const methodInfo_MonitoringService_LabelNames = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.LabelNamesResponse,
  /** @param {!proto.monitoring.LabelNamesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.LabelNamesResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.LabelNamesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.LabelNamesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.LabelNamesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.labelNames =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/LabelNames',
      request,
      metadata || {},
      methodInfo_MonitoringService_LabelNames,
      callback);
};


/**
 * @param {!proto.monitoring.LabelNamesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.LabelNamesResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.labelNames =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/LabelNames',
      request,
      metadata || {},
      methodInfo_MonitoringService_LabelNames);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.LabelValuesRequest,
 *   !proto.monitoring.LabelValuesResponse>}
 */
const methodInfo_MonitoringService_LabelValues = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.LabelValuesResponse,
  /** @param {!proto.monitoring.LabelValuesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.LabelValuesResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.LabelValuesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.LabelValuesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.LabelValuesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.labelValues =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/LabelValues',
      request,
      metadata || {},
      methodInfo_MonitoringService_LabelValues,
      callback);
};


/**
 * @param {!proto.monitoring.LabelValuesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.LabelValuesResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.labelValues =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/LabelValues',
      request,
      metadata || {},
      methodInfo_MonitoringService_LabelValues);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.QueryRequest,
 *   !proto.monitoring.QueryResponse>}
 */
const methodInfo_MonitoringService_Query = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.QueryResponse,
  /** @param {!proto.monitoring.QueryRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.QueryResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.QueryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.QueryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.QueryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.query =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/Query',
      request,
      metadata || {},
      methodInfo_MonitoringService_Query,
      callback);
};


/**
 * @param {!proto.monitoring.QueryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.QueryResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.query =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/Query',
      request,
      metadata || {},
      methodInfo_MonitoringService_Query);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.QueryRangeRequest,
 *   !proto.monitoring.QueryRangeResponse>}
 */
const methodInfo_MonitoringService_QueryRange = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.QueryRangeResponse,
  /** @param {!proto.monitoring.QueryRangeRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.QueryRangeResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.QueryRangeRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.QueryRangeResponse>}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.queryRange =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/monitoring.MonitoringService/QueryRange',
      request,
      metadata || {},
      methodInfo_MonitoringService_QueryRange);
};


/**
 * @param {!proto.monitoring.QueryRangeRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.QueryRangeResponse>}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.queryRange =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/monitoring.MonitoringService/QueryRange',
      request,
      metadata || {},
      methodInfo_MonitoringService_QueryRange);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.SeriesRequest,
 *   !proto.monitoring.SeriesResponse>}
 */
const methodInfo_MonitoringService_Series = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.SeriesResponse,
  /** @param {!proto.monitoring.SeriesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.SeriesResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.SeriesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.SeriesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.SeriesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.series =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/Series',
      request,
      metadata || {},
      methodInfo_MonitoringService_Series,
      callback);
};


/**
 * @param {!proto.monitoring.SeriesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.SeriesResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.series =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/Series',
      request,
      metadata || {},
      methodInfo_MonitoringService_Series);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.SnapshotRequest,
 *   !proto.monitoring.SnapshotResponse>}
 */
const methodInfo_MonitoringService_Snapshot = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.SnapshotResponse,
  /** @param {!proto.monitoring.SnapshotRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.SnapshotResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.SnapshotRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.SnapshotResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.SnapshotResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.snapshot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/Snapshot',
      request,
      metadata || {},
      methodInfo_MonitoringService_Snapshot,
      callback);
};


/**
 * @param {!proto.monitoring.SnapshotRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.SnapshotResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.snapshot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/Snapshot',
      request,
      metadata || {},
      methodInfo_MonitoringService_Snapshot);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.RulesRequest,
 *   !proto.monitoring.RulesResponse>}
 */
const methodInfo_MonitoringService_Rules = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.RulesResponse,
  /** @param {!proto.monitoring.RulesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.RulesResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.RulesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.RulesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.RulesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.rules =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/Rules',
      request,
      metadata || {},
      methodInfo_MonitoringService_Rules,
      callback);
};


/**
 * @param {!proto.monitoring.RulesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.RulesResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.rules =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/Rules',
      request,
      metadata || {},
      methodInfo_MonitoringService_Rules);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.TargetsRequest,
 *   !proto.monitoring.TargetsResponse>}
 */
const methodInfo_MonitoringService_Targets = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.TargetsResponse,
  /** @param {!proto.monitoring.TargetsRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.TargetsResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.TargetsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.TargetsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.TargetsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.targets =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/Targets',
      request,
      metadata || {},
      methodInfo_MonitoringService_Targets,
      callback);
};


/**
 * @param {!proto.monitoring.TargetsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.TargetsResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.targets =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/Targets',
      request,
      metadata || {},
      methodInfo_MonitoringService_Targets);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.monitoring.TargetsMetadataRequest,
 *   !proto.monitoring.TargetsMetadataResponse>}
 */
const methodInfo_MonitoringService_TargetsMetadata = new grpc.web.AbstractClientBase.MethodInfo(
  proto.monitoring.TargetsMetadataResponse,
  /** @param {!proto.monitoring.TargetsMetadataRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.monitoring.TargetsMetadataResponse.deserializeBinary
);


/**
 * @param {!proto.monitoring.TargetsMetadataRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.monitoring.TargetsMetadataResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.monitoring.TargetsMetadataResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.monitoring.MonitoringServiceClient.prototype.targetsMetadata =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/monitoring.MonitoringService/TargetsMetadata',
      request,
      metadata || {},
      methodInfo_MonitoringService_TargetsMetadata,
      callback);
};


/**
 * @param {!proto.monitoring.TargetsMetadataRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.monitoring.TargetsMetadataResponse>}
 *     A native promise that resolves to the response
 */
proto.monitoring.MonitoringServicePromiseClient.prototype.targetsMetadata =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/monitoring.MonitoringService/TargetsMetadata',
      request,
      metadata || {},
      methodInfo_MonitoringService_TargetsMetadata);
};


module.exports = proto.monitoring;

