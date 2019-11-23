/**
 * @fileoverview gRPC-Web generated client stub for services
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.services = require('./services_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.services.ServiceDiscoveryClient =
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
proto.services.ServiceDiscoveryPromiseClient =
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
 *   !proto.services.GetServiceDescriptorRequest,
 *   !proto.services.GetServiceDescriptorResponse>}
 */
const methodInfo_ServiceDiscovery_GetServiceDescriptor = new grpc.web.AbstractClientBase.MethodInfo(
  proto.services.GetServiceDescriptorResponse,
  /** @param {!proto.services.GetServiceDescriptorRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.services.GetServiceDescriptorResponse.deserializeBinary
);


/**
 * @param {!proto.services.GetServiceDescriptorRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.services.GetServiceDescriptorResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.services.GetServiceDescriptorResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.services.ServiceDiscoveryClient.prototype.getServiceDescriptor =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/services.ServiceDiscovery/GetServiceDescriptor',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_GetServiceDescriptor,
      callback);
};


/**
 * @param {!proto.services.GetServiceDescriptorRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.services.GetServiceDescriptorResponse>}
 *     A native promise that resolves to the response
 */
proto.services.ServiceDiscoveryPromiseClient.prototype.getServiceDescriptor =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/services.ServiceDiscovery/GetServiceDescriptor',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_GetServiceDescriptor);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.services.GetServicesDescriptorRequest,
 *   !proto.services.GetServicesDescriptorResponse>}
 */
const methodInfo_ServiceDiscovery_GetServicesDescriptor = new grpc.web.AbstractClientBase.MethodInfo(
  proto.services.GetServicesDescriptorResponse,
  /** @param {!proto.services.GetServicesDescriptorRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.services.GetServicesDescriptorResponse.deserializeBinary
);


/**
 * @param {!proto.services.GetServicesDescriptorRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.services.GetServicesDescriptorResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.services.GetServicesDescriptorResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.services.ServiceDiscoveryClient.prototype.getServicesDescriptor =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/services.ServiceDiscovery/GetServicesDescriptor',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_GetServicesDescriptor,
      callback);
};


/**
 * @param {!proto.services.GetServicesDescriptorRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.services.GetServicesDescriptorResponse>}
 *     A native promise that resolves to the response
 */
proto.services.ServiceDiscoveryPromiseClient.prototype.getServicesDescriptor =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/services.ServiceDiscovery/GetServicesDescriptor',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_GetServicesDescriptor);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.services.PublishServiceRequest,
 *   !proto.services.PublishServiceResponse>}
 */
const methodInfo_ServiceDiscovery_publishService = new grpc.web.AbstractClientBase.MethodInfo(
  proto.services.PublishServiceResponse,
  /** @param {!proto.services.PublishServiceRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.services.PublishServiceResponse.deserializeBinary
);


/**
 * @param {!proto.services.PublishServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.services.PublishServiceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.services.PublishServiceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.services.ServiceDiscoveryClient.prototype.publishService =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/services.ServiceDiscovery/publishService',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_publishService,
      callback);
};


/**
 * @param {!proto.services.PublishServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.services.PublishServiceResponse>}
 *     A native promise that resolves to the response
 */
proto.services.ServiceDiscoveryPromiseClient.prototype.publishService =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/services.ServiceDiscovery/publishService',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_publishService);
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.services.ServiceRepositoryClient =
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
proto.services.ServiceRepositoryPromiseClient =
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
 *   !proto.services.DownloadBundleRequest,
 *   !proto.services.DownloadBundleResponse>}
 */
const methodInfo_ServiceRepository_downloadBundle = new grpc.web.AbstractClientBase.MethodInfo(
  proto.services.DownloadBundleResponse,
  /** @param {!proto.services.DownloadBundleRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.services.DownloadBundleResponse.deserializeBinary
);


/**
 * @param {!proto.services.DownloadBundleRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.services.DownloadBundleResponse>}
 *     The XHR Node Readable Stream
 */
proto.services.ServiceRepositoryClient.prototype.downloadBundle =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/services.ServiceRepository/downloadBundle',
      request,
      metadata || {},
      methodInfo_ServiceRepository_downloadBundle);
};


/**
 * @param {!proto.services.DownloadBundleRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.services.DownloadBundleResponse>}
 *     The XHR Node Readable Stream
 */
proto.services.ServiceRepositoryPromiseClient.prototype.downloadBundle =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/services.ServiceRepository/downloadBundle',
      request,
      metadata || {},
      methodInfo_ServiceRepository_downloadBundle);
};


module.exports = proto.services;

