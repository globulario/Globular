/**
 * @fileoverview gRPC-Web generated client stub for admin
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.admin = require('./admin_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.admin.AdminServiceClient =
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
proto.admin.AdminServicePromiseClient =
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
 *   !proto.admin.GetConfigRequest,
 *   !proto.admin.GetConfigResponse>}
 */
const methodInfo_AdminService_GetConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.GetConfigResponse,
  /** @param {!proto.admin.GetConfigRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.GetConfigResponse.deserializeBinary
);


/**
 * @param {!proto.admin.GetConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.GetConfigResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.GetConfigResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.getConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/GetConfig',
      request,
      metadata || {},
      methodInfo_AdminService_GetConfig,
      callback);
};


/**
 * @param {!proto.admin.GetConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.GetConfigResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.getConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/GetConfig',
      request,
      metadata || {},
      methodInfo_AdminService_GetConfig);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.GetConfigRequest,
 *   !proto.admin.GetConfigResponse>}
 */
const methodInfo_AdminService_GetFullConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.GetConfigResponse,
  /** @param {!proto.admin.GetConfigRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.GetConfigResponse.deserializeBinary
);


/**
 * @param {!proto.admin.GetConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.GetConfigResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.GetConfigResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.getFullConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/GetFullConfig',
      request,
      metadata || {},
      methodInfo_AdminService_GetFullConfig,
      callback);
};


/**
 * @param {!proto.admin.GetConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.GetConfigResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.getFullConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/GetFullConfig',
      request,
      metadata || {},
      methodInfo_AdminService_GetFullConfig);
};


module.exports = proto.admin;

