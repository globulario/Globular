/**
 * @fileoverview gRPC-Web generated client stub for echo
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.echo = require('./plc_link_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.echo.PlcLinkServiceClient =
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
proto.echo.PlcLinkServicePromiseClient =
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
 *   !proto.echo.CreateConnectionRqst,
 *   !proto.echo.CreateConnectionRsp>}
 */
const methodInfo_PlcLinkService_CreateConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.CreateConnectionRsp,
  /** @param {!proto.echo.CreateConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.CreateConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.echo.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.CreateConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.CreateConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.PlcLinkServiceClient.prototype.createConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.PlcLinkService/CreateConnection',
      request,
      metadata || {},
      methodInfo_PlcLinkService_CreateConnection,
      callback);
};


/**
 * @param {!proto.echo.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.CreateConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.echo.PlcLinkServicePromiseClient.prototype.createConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.PlcLinkService/CreateConnection',
      request,
      metadata || {},
      methodInfo_PlcLinkService_CreateConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.DeleteConnectionRqst,
 *   !proto.echo.DeleteConnectionRsp>}
 */
const methodInfo_PlcLinkService_DeleteConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.DeleteConnectionRsp,
  /** @param {!proto.echo.DeleteConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.DeleteConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.echo.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.DeleteConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.DeleteConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.PlcLinkServiceClient.prototype.deleteConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.PlcLinkService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_PlcLinkService_DeleteConnection,
      callback);
};


/**
 * @param {!proto.echo.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.DeleteConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.echo.PlcLinkServicePromiseClient.prototype.deleteConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.PlcLinkService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_PlcLinkService_DeleteConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.LinkRqst,
 *   !proto.echo.LinkRsp>}
 */
const methodInfo_PlcLinkService_Link = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.LinkRsp,
  /** @param {!proto.echo.LinkRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.LinkRsp.deserializeBinary
);


/**
 * @param {!proto.echo.LinkRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.LinkRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.LinkRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.PlcLinkServiceClient.prototype.link =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.PlcLinkService/Link',
      request,
      metadata || {},
      methodInfo_PlcLinkService_Link,
      callback);
};


/**
 * @param {!proto.echo.LinkRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.LinkRsp>}
 *     A native promise that resolves to the response
 */
proto.echo.PlcLinkServicePromiseClient.prototype.link =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.PlcLinkService/Link',
      request,
      metadata || {},
      methodInfo_PlcLinkService_Link);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.UnLinkRqst,
 *   !proto.echo.UnLinkRsp>}
 */
const methodInfo_PlcLinkService_UnLink = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.UnLinkRsp,
  /** @param {!proto.echo.UnLinkRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.UnLinkRsp.deserializeBinary
);


/**
 * @param {!proto.echo.UnLinkRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.UnLinkRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.UnLinkRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.PlcLinkServiceClient.prototype.unLink =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.PlcLinkService/UnLink',
      request,
      metadata || {},
      methodInfo_PlcLinkService_UnLink,
      callback);
};


/**
 * @param {!proto.echo.UnLinkRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.UnLinkRsp>}
 *     A native promise that resolves to the response
 */
proto.echo.PlcLinkServicePromiseClient.prototype.unLink =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.PlcLinkService/UnLink',
      request,
      metadata || {},
      methodInfo_PlcLinkService_UnLink);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.SuspendRqst,
 *   !proto.echo.SuspendRsp>}
 */
const methodInfo_PlcLinkService_Suspend = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.SuspendRsp,
  /** @param {!proto.echo.SuspendRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.SuspendRsp.deserializeBinary
);


/**
 * @param {!proto.echo.SuspendRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.SuspendRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.SuspendRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.PlcLinkServiceClient.prototype.suspend =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.PlcLinkService/Suspend',
      request,
      metadata || {},
      methodInfo_PlcLinkService_Suspend,
      callback);
};


/**
 * @param {!proto.echo.SuspendRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.SuspendRsp>}
 *     A native promise that resolves to the response
 */
proto.echo.PlcLinkServicePromiseClient.prototype.suspend =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.PlcLinkService/Suspend',
      request,
      metadata || {},
      methodInfo_PlcLinkService_Suspend);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.ResumeRqst,
 *   !proto.echo.ResumeRsp>}
 */
const methodInfo_PlcLinkService_Resume = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.ResumeRsp,
  /** @param {!proto.echo.ResumeRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.ResumeRsp.deserializeBinary
);


/**
 * @param {!proto.echo.ResumeRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.ResumeRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.ResumeRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.PlcLinkServiceClient.prototype.resume =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.PlcLinkService/Resume',
      request,
      metadata || {},
      methodInfo_PlcLinkService_Resume,
      callback);
};


/**
 * @param {!proto.echo.ResumeRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.ResumeRsp>}
 *     A native promise that resolves to the response
 */
proto.echo.PlcLinkServicePromiseClient.prototype.resume =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.PlcLinkService/Resume',
      request,
      metadata || {},
      methodInfo_PlcLinkService_Resume);
};


module.exports = proto.echo;

