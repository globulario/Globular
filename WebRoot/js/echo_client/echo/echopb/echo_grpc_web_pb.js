/**
 * @fileoverview gRPC-Web generated client stub for echo
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.echo = require('./echo_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.echo.EchoServiceClient =
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
proto.echo.EchoServicePromiseClient =
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
 *   !proto.echo.EchoRqst,
 *   !proto.echo.EchoRsp>}
 */
const methodInfo_EchoService_Echo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.EchoRsp,
  /** @param {!proto.echo.EchoRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.EchoRsp.deserializeBinary
);


/**
 * @param {!proto.echo.EchoRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.EchoRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.EchoRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.EchoServiceClient.prototype.echo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.EchoService/Echo',
      request,
      metadata || {},
      methodInfo_EchoService_Echo,
      callback);
};


/**
 * @param {!proto.echo.EchoRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.EchoRsp>}
 *     A native promise that resolves to the response
 */
proto.echo.EchoServicePromiseClient.prototype.echo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.EchoService/Echo',
      request,
      metadata || {},
      methodInfo_EchoService_Echo);
};


module.exports = proto.echo;

