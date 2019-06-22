/**
 * @fileoverview gRPC-Web generated client stub for event
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.event = require('./event_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.event.EventServiceClient =
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
proto.event.EventServicePromiseClient =
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
 *   !proto.event.SubscribeRequest,
 *   !proto.event.SubscribeResponse>}
 */
const methodInfo_EventService_Subscribe = new grpc.web.AbstractClientBase.MethodInfo(
  proto.event.SubscribeResponse,
  /** @param {!proto.event.SubscribeRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.event.SubscribeResponse.deserializeBinary
);


/**
 * @param {!proto.event.SubscribeRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.event.SubscribeResponse>}
 *     The XHR Node Readable Stream
 */
proto.event.EventServiceClient.prototype.subscribe =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/event.EventService/Subscribe',
      request,
      metadata || {},
      methodInfo_EventService_Subscribe);
};


/**
 * @param {!proto.event.SubscribeRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.event.SubscribeResponse>}
 *     The XHR Node Readable Stream
 */
proto.event.EventServicePromiseClient.prototype.subscribe =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/event.EventService/Subscribe',
      request,
      metadata || {},
      methodInfo_EventService_Subscribe);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.event.UnSubscribeRequest,
 *   !proto.event.UnSubscribeResponse>}
 */
const methodInfo_EventService_UnSubscribe = new grpc.web.AbstractClientBase.MethodInfo(
  proto.event.UnSubscribeResponse,
  /** @param {!proto.event.UnSubscribeRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.event.UnSubscribeResponse.deserializeBinary
);


/**
 * @param {!proto.event.UnSubscribeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.event.UnSubscribeResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.event.UnSubscribeResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.event.EventServiceClient.prototype.unSubscribe =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/event.EventService/UnSubscribe',
      request,
      metadata || {},
      methodInfo_EventService_UnSubscribe,
      callback);
};


/**
 * @param {!proto.event.UnSubscribeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.event.UnSubscribeResponse>}
 *     A native promise that resolves to the response
 */
proto.event.EventServicePromiseClient.prototype.unSubscribe =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/event.EventService/UnSubscribe',
      request,
      metadata || {},
      methodInfo_EventService_UnSubscribe);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.event.PublishRequest,
 *   !proto.event.PublishResponse>}
 */
const methodInfo_EventService_Publish = new grpc.web.AbstractClientBase.MethodInfo(
  proto.event.PublishResponse,
  /** @param {!proto.event.PublishRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.event.PublishResponse.deserializeBinary
);


/**
 * @param {!proto.event.PublishRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.event.PublishResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.event.PublishResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.event.EventServiceClient.prototype.publish =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/event.EventService/Publish',
      request,
      metadata || {},
      methodInfo_EventService_Publish,
      callback);
};


/**
 * @param {!proto.event.PublishRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.event.PublishResponse>}
 *     A native promise that resolves to the response
 */
proto.event.EventServicePromiseClient.prototype.publish =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/event.EventService/Publish',
      request,
      metadata || {},
      methodInfo_EventService_Publish);
};


module.exports = proto.event;

