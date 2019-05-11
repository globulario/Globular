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
 *   !proto.persistence.CreateConnectionRqst,
 *   !proto.persistence.CreateConnectionRsp>}
 */
const methodInfo_PersistenceService_CreateConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.persistence.CreateConnectionRsp,
  /** @param {!proto.persistence.CreateConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.persistence.CreateConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.persistence.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.persistence.CreateConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.persistence.CreateConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.persistence.PersistenceServiceClient.prototype.createConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/persistence.PersistenceService/CreateConnection',
      request,
      metadata || {},
      methodInfo_PersistenceService_CreateConnection,
      callback);
};


/**
 * @param {!proto.persistence.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.persistence.CreateConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.persistence.PersistenceServicePromiseClient.prototype.createConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/persistence.PersistenceService/CreateConnection',
      request,
      metadata || {},
      methodInfo_PersistenceService_CreateConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.persistence.DeleteConnectionRqst,
 *   !proto.persistence.DeleteConnectionRsp>}
 */
const methodInfo_PersistenceService_DeleteConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.persistence.DeleteConnectionRsp,
  /** @param {!proto.persistence.DeleteConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.persistence.DeleteConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.persistence.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.persistence.DeleteConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.persistence.DeleteConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.persistence.PersistenceServiceClient.prototype.deleteConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/persistence.PersistenceService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_PersistenceService_DeleteConnection,
      callback);
};


/**
 * @param {!proto.persistence.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.persistence.DeleteConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.persistence.PersistenceServicePromiseClient.prototype.deleteConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/persistence.PersistenceService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_PersistenceService_DeleteConnection);
};


module.exports = proto.persistence;

