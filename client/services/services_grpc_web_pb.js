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
 *   !proto.services.FindServicesDescriptorRequest,
 *   !proto.services.FindServicesDescriptorResponse>}
 */
const methodInfo_ServiceDiscovery_FindServices = new grpc.web.AbstractClientBase.MethodInfo(
  proto.services.FindServicesDescriptorResponse,
  /** @param {!proto.services.FindServicesDescriptorRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.services.FindServicesDescriptorResponse.deserializeBinary
);


/**
 * @param {!proto.services.FindServicesDescriptorRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.services.FindServicesDescriptorResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.services.FindServicesDescriptorResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.services.ServiceDiscoveryClient.prototype.findServices =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/services.ServiceDiscovery/FindServices',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_FindServices,
      callback);
};


/**
 * @param {!proto.services.FindServicesDescriptorRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.services.FindServicesDescriptorResponse>}
 *     A native promise that resolves to the response
 */
proto.services.ServiceDiscoveryPromiseClient.prototype.findServices =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/services.ServiceDiscovery/FindServices',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_FindServices);
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
 *   !proto.services.PublishServiceDescriptorRequest,
 *   !proto.services.PublishServiceDescriptorResponse>}
 */
const methodInfo_ServiceDiscovery_publishServiceDescriptor = new grpc.web.AbstractClientBase.MethodInfo(
  proto.services.PublishServiceDescriptorResponse,
  /** @param {!proto.services.PublishServiceDescriptorRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.services.PublishServiceDescriptorResponse.deserializeBinary
);


/**
 * @param {!proto.services.PublishServiceDescriptorRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.services.PublishServiceDescriptorResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.services.PublishServiceDescriptorResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.services.ServiceDiscoveryClient.prototype.publishServiceDescriptor =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/services.ServiceDiscovery/publishServiceDescriptor',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_publishServiceDescriptor,
      callback);
};


/**
 * @param {!proto.services.PublishServiceDescriptorRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.services.PublishServiceDescriptorResponse>}
 *     A native promise that resolves to the response
 */
proto.services.ServiceDiscoveryPromiseClient.prototype.publishServiceDescriptor =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/services.ServiceDiscovery/publishServiceDescriptor',
      request,
      metadata || {},
      methodInfo_ServiceDiscovery_publishServiceDescriptor);
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

