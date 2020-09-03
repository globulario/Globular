/**
 * @fileoverview gRPC-Web generated client stub for lb
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.lb = require('./lb_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.lb.LoadBalancingServiceClient =
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

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.lb.LoadBalancingServicePromiseClient =
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

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.lb.GetCanditatesRequest,
 *   !proto.lb.GetCanditatesResponse>}
 */
const methodDescriptor_LoadBalancingService_GetCanditates = new grpc.web.MethodDescriptor(
  '/lb.LoadBalancingService/GetCanditates',
  grpc.web.MethodType.UNARY,
  proto.lb.GetCanditatesRequest,
  proto.lb.GetCanditatesResponse,
  /**
   * @param {!proto.lb.GetCanditatesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.lb.GetCanditatesResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.lb.GetCanditatesRequest,
 *   !proto.lb.GetCanditatesResponse>}
 */
const methodInfo_LoadBalancingService_GetCanditates = new grpc.web.AbstractClientBase.MethodInfo(
  proto.lb.GetCanditatesResponse,
  /**
   * @param {!proto.lb.GetCanditatesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.lb.GetCanditatesResponse.deserializeBinary
);


/**
 * @param {!proto.lb.GetCanditatesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.lb.GetCanditatesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.lb.GetCanditatesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.lb.LoadBalancingServiceClient.prototype.getCanditates =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/lb.LoadBalancingService/GetCanditates',
      request,
      metadata || {},
      methodDescriptor_LoadBalancingService_GetCanditates,
      callback);
};


/**
 * @param {!proto.lb.GetCanditatesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.lb.GetCanditatesResponse>}
 *     A native promise that resolves to the response
 */
proto.lb.LoadBalancingServicePromiseClient.prototype.getCanditates =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/lb.LoadBalancingService/GetCanditates',
      request,
      metadata || {},
      methodDescriptor_LoadBalancingService_GetCanditates);
};


module.exports = proto.lb;

