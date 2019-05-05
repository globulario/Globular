/**
 * @fileoverview gRPC-Web generated client stub for persistence
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.persistence = require('./persistence_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.persistence.PersistenceServiceClient =
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
proto.persistence.PersistenceServicePromiseClient =
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
 *   !proto.persistence.PersistEntityRqst,
 *   !proto.persistence.PersistEntityRsp>}
 */
const methodInfo_PersistenceService_PersistEntity = new grpc.web.AbstractClientBase.MethodInfo(
  proto.persistence.PersistEntityRsp,
  /** @param {!proto.persistence.PersistEntityRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.persistence.PersistEntityRsp.deserializeBinary
);


/**
 * @param {!proto.persistence.PersistEntityRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.persistence.PersistEntityRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.persistence.PersistEntityRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.persistence.PersistenceServiceClient.prototype.persistEntity =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/persistence.PersistenceService/PersistEntity',
      request,
      metadata || {},
      methodInfo_PersistenceService_PersistEntity,
      callback);
};


/**
 * @param {!proto.persistence.PersistEntityRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.persistence.PersistEntityRsp>}
 *     A native promise that resolves to the response
 */
proto.persistence.PersistenceServicePromiseClient.prototype.persistEntity =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/persistence.PersistenceService/PersistEntity',
      request,
      metadata || {},
      methodInfo_PersistenceService_PersistEntity);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.persistence.GetEntityByUuidRqst,
 *   !proto.persistence.GetEntityByUuidRsp>}
 */
const methodInfo_PersistenceService_GetEntityByUuid = new grpc.web.AbstractClientBase.MethodInfo(
  proto.persistence.GetEntityByUuidRsp,
  /** @param {!proto.persistence.GetEntityByUuidRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.persistence.GetEntityByUuidRsp.deserializeBinary
);


/**
 * @param {!proto.persistence.GetEntityByUuidRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.persistence.GetEntityByUuidRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.persistence.GetEntityByUuidRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.persistence.PersistenceServiceClient.prototype.getEntityByUuid =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/persistence.PersistenceService/GetEntityByUuid',
      request,
      metadata || {},
      methodInfo_PersistenceService_GetEntityByUuid,
      callback);
};


/**
 * @param {!proto.persistence.GetEntityByUuidRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.persistence.GetEntityByUuidRsp>}
 *     A native promise that resolves to the response
 */
proto.persistence.PersistenceServicePromiseClient.prototype.getEntityByUuid =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/persistence.PersistenceService/GetEntityByUuid',
      request,
      metadata || {},
      methodInfo_PersistenceService_GetEntityByUuid);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.persistence.GetEntitiesByTypenameRqst,
 *   !proto.persistence.GetEntitiesByTypenameRsp>}
 */
const methodInfo_PersistenceService_GetEntitiesByTypename = new grpc.web.AbstractClientBase.MethodInfo(
  proto.persistence.GetEntitiesByTypenameRsp,
  /** @param {!proto.persistence.GetEntitiesByTypenameRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.persistence.GetEntitiesByTypenameRsp.deserializeBinary
);


/**
 * @param {!proto.persistence.GetEntitiesByTypenameRqst} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.persistence.GetEntitiesByTypenameRsp>}
 *     The XHR Node Readable Stream
 */
proto.persistence.PersistenceServiceClient.prototype.getEntitiesByTypename =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/persistence.PersistenceService/GetEntitiesByTypename',
      request,
      metadata || {},
      methodInfo_PersistenceService_GetEntitiesByTypename);
};


/**
 * @param {!proto.persistence.GetEntitiesByTypenameRqst} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.persistence.GetEntitiesByTypenameRsp>}
 *     The XHR Node Readable Stream
 */
proto.persistence.PersistenceServicePromiseClient.prototype.getEntitiesByTypename =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/persistence.PersistenceService/GetEntitiesByTypename',
      request,
      metadata || {},
      methodInfo_PersistenceService_GetEntitiesByTypename);
};


module.exports = proto.persistence;

