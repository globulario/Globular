/**
 * @fileoverview gRPC-Web generated client stub for catalog
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.catalog = require('./catalog_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.catalog.CatalogServiceClient =
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
proto.catalog.CatalogServicePromiseClient =
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
 *   !proto.catalog.CreateConnectionRqst,
 *   !proto.catalog.CreateConnectionRsp>}
 */
const methodInfo_CatalogService_CreateConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.CreateConnectionRsp,
  /** @param {!proto.catalog.CreateConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.CreateConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.catalog.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.CreateConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.CreateConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.createConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/CreateConnection',
      request,
      metadata || {},
      methodInfo_CatalogService_CreateConnection,
      callback);
};


/**
 * @param {!proto.catalog.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.CreateConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.createConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/CreateConnection',
      request,
      metadata || {},
      methodInfo_CatalogService_CreateConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeleteConnectionRqst,
 *   !proto.catalog.DeleteConnectionRsp>}
 */
const methodInfo_CatalogService_DeleteConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeleteConnectionRsp,
  /** @param {!proto.catalog.DeleteConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeleteConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.catalog.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeleteConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeleteConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deleteConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_CatalogService_DeleteConnection,
      callback);
};


/**
 * @param {!proto.catalog.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeleteConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deleteConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_CatalogService_DeleteConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SaveUnitOfMesureRequest,
 *   !proto.catalog.SaveUnitOfMesureResponse>}
 */
const methodInfo_CatalogService_SaveUnitOfMesure = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveUnitOfMesureResponse,
  /** @param {!proto.catalog.SaveUnitOfMesureRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveUnitOfMesureResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveUnitOfMesureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveUnitOfMesureResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveUnitOfMesureResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveUnitOfMesure =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveUnitOfMesure',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveUnitOfMesure,
      callback);
};


/**
 * @param {!proto.catalog.SaveUnitOfMesureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveUnitOfMesureResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveUnitOfMesure =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveUnitOfMesure',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveUnitOfMesure);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SavePropertyDefinitionRequest,
 *   !proto.catalog.SavePropertyDefinitionResponse>}
 */
const methodInfo_CatalogService_SavePropertyDefinition = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SavePropertyDefinitionResponse,
  /** @param {!proto.catalog.SavePropertyDefinitionRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SavePropertyDefinitionResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SavePropertyDefinitionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SavePropertyDefinitionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SavePropertyDefinitionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.savePropertyDefinition =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SavePropertyDefinition',
      request,
      metadata || {},
      methodInfo_CatalogService_SavePropertyDefinition,
      callback);
};


/**
 * @param {!proto.catalog.SavePropertyDefinitionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SavePropertyDefinitionResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.savePropertyDefinition =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SavePropertyDefinition',
      request,
      metadata || {},
      methodInfo_CatalogService_SavePropertyDefinition);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SaveItemDefinitionRequest,
 *   !proto.catalog.SaveItemDefinitionResponse>}
 */
const methodInfo_CatalogService_SaveItemDefinition = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveItemDefinitionResponse,
  /** @param {!proto.catalog.SaveItemDefinitionRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveItemDefinitionResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveItemDefinitionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveItemDefinitionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveItemDefinitionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveItemDefinition =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveItemDefinition',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveItemDefinition,
      callback);
};


/**
 * @param {!proto.catalog.SaveItemDefinitionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveItemDefinitionResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveItemDefinition =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveItemDefinition',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveItemDefinition);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SaveItemInstanceRequest,
 *   !proto.catalog.SaveItemInstanceResponse>}
 */
const methodInfo_CatalogService_SaveItemInstance = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveItemInstanceResponse,
  /** @param {!proto.catalog.SaveItemInstanceRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveItemInstanceResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveItemInstanceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveItemInstanceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveItemInstanceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveItemInstance =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveItemInstance',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveItemInstance,
      callback);
};


/**
 * @param {!proto.catalog.SaveItemInstanceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveItemInstanceResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveItemInstance =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveItemInstance',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveItemInstance);
};


module.exports = proto.catalog;

