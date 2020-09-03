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
 *   !proto.catalog.SaveUnitOfMeasureRequest,
 *   !proto.catalog.SaveUnitOfMeasureResponse>}
 */
const methodInfo_CatalogService_SaveUnitOfMeasure = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveUnitOfMeasureResponse,
  /** @param {!proto.catalog.SaveUnitOfMeasureRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveUnitOfMeasureResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveUnitOfMeasureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveUnitOfMeasureResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveUnitOfMeasureResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveUnitOfMeasure =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveUnitOfMeasure',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveUnitOfMeasure,
      callback);
};


/**
 * @param {!proto.catalog.SaveUnitOfMeasureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveUnitOfMeasureResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveUnitOfMeasure =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveUnitOfMeasure',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveUnitOfMeasure);
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


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SaveInventoryRequest,
 *   !proto.catalog.SaveInventoryResponse>}
 */
const methodInfo_CatalogService_SaveInventory = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveInventoryResponse,
  /** @param {!proto.catalog.SaveInventoryRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveInventoryResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveInventoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveInventoryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveInventoryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveInventory =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveInventory',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveInventory,
      callback);
};


/**
 * @param {!proto.catalog.SaveInventoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveInventoryResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveInventory =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveInventory',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveInventory);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SaveManufacturerRequest,
 *   !proto.catalog.SaveManufacturerResponse>}
 */
const methodInfo_CatalogService_SaveManufacturer = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveManufacturerResponse,
  /** @param {!proto.catalog.SaveManufacturerRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveManufacturerResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveManufacturerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveManufacturerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveManufacturer =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveManufacturer,
      callback);
};


/**
 * @param {!proto.catalog.SaveManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveManufacturerResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveManufacturer =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveManufacturer);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SaveSupplierRequest,
 *   !proto.catalog.SaveSupplierResponse>}
 */
const methodInfo_CatalogService_SaveSupplier = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveSupplierResponse,
  /** @param {!proto.catalog.SaveSupplierRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveSupplierResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveSupplierResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveSupplierResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveSupplier =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveSupplier,
      callback);
};


/**
 * @param {!proto.catalog.SaveSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveSupplierResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveSupplier =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveSupplier);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SaveLocalisationRequest,
 *   !proto.catalog.SaveLocalisationResponse>}
 */
const methodInfo_CatalogService_SaveLocalisation = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveLocalisationResponse,
  /** @param {!proto.catalog.SaveLocalisationRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveLocalisationResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveLocalisationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveLocalisationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveLocalisationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveLocalisation =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveLocalisation',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveLocalisation,
      callback);
};


/**
 * @param {!proto.catalog.SaveLocalisationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveLocalisationResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveLocalisation =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveLocalisation',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveLocalisation);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SavePackageRequest,
 *   !proto.catalog.SavePackageResponse>}
 */
const methodInfo_CatalogService_SavePackage = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SavePackageResponse,
  /** @param {!proto.catalog.SavePackageRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SavePackageResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SavePackageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SavePackageResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SavePackageResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.savePackage =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SavePackage',
      request,
      metadata || {},
      methodInfo_CatalogService_SavePackage,
      callback);
};


/**
 * @param {!proto.catalog.SavePackageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SavePackageResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.savePackage =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SavePackage',
      request,
      metadata || {},
      methodInfo_CatalogService_SavePackage);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SavePackageSupplierRequest,
 *   !proto.catalog.SavePackageSupplierResponse>}
 */
const methodInfo_CatalogService_SavePackageSupplier = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SavePackageSupplierResponse,
  /** @param {!proto.catalog.SavePackageSupplierRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SavePackageSupplierResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SavePackageSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SavePackageSupplierResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SavePackageSupplierResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.savePackageSupplier =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SavePackageSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_SavePackageSupplier,
      callback);
};


/**
 * @param {!proto.catalog.SavePackageSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SavePackageSupplierResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.savePackageSupplier =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SavePackageSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_SavePackageSupplier);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SaveItemManufacturerRequest,
 *   !proto.catalog.SaveItemManufacturerResponse>}
 */
const methodInfo_CatalogService_SaveItemManufacturer = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveItemManufacturerResponse,
  /** @param {!proto.catalog.SaveItemManufacturerRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveItemManufacturerResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveItemManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveItemManufacturerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveItemManufacturerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveItemManufacturer =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveItemManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveItemManufacturer,
      callback);
};


/**
 * @param {!proto.catalog.SaveItemManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveItemManufacturerResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveItemManufacturer =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveItemManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveItemManufacturer);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.SaveCategoryRequest,
 *   !proto.catalog.SaveCategoryResponse>}
 */
const methodInfo_CatalogService_SaveCategory = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.SaveCategoryResponse,
  /** @param {!proto.catalog.SaveCategoryRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.SaveCategoryResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.SaveCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.SaveCategoryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.SaveCategoryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.saveCategory =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/SaveCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveCategory,
      callback);
};


/**
 * @param {!proto.catalog.SaveCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.SaveCategoryResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.saveCategory =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/SaveCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_SaveCategory);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.AppendItemDefinitionCategoryRequest,
 *   !proto.catalog.AppendItemDefinitionCategoryResponse>}
 */
const methodInfo_CatalogService_AppendItemDefinitionCategory = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.AppendItemDefinitionCategoryResponse,
  /** @param {!proto.catalog.AppendItemDefinitionCategoryRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.AppendItemDefinitionCategoryResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.AppendItemDefinitionCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.AppendItemDefinitionCategoryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.AppendItemDefinitionCategoryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.appendItemDefinitionCategory =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/AppendItemDefinitionCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_AppendItemDefinitionCategory,
      callback);
};


/**
 * @param {!proto.catalog.AppendItemDefinitionCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.AppendItemDefinitionCategoryResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.appendItemDefinitionCategory =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/AppendItemDefinitionCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_AppendItemDefinitionCategory);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.RemoveItemDefinitionCategoryRequest,
 *   !proto.catalog.RemoveItemDefinitionCategoryResponse>}
 */
const methodInfo_CatalogService_RemoveItemDefinitionCategory = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.RemoveItemDefinitionCategoryResponse,
  /** @param {!proto.catalog.RemoveItemDefinitionCategoryRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.RemoveItemDefinitionCategoryResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.RemoveItemDefinitionCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.RemoveItemDefinitionCategoryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.RemoveItemDefinitionCategoryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.removeItemDefinitionCategory =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/RemoveItemDefinitionCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_RemoveItemDefinitionCategory,
      callback);
};


/**
 * @param {!proto.catalog.RemoveItemDefinitionCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.RemoveItemDefinitionCategoryResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.removeItemDefinitionCategory =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/RemoveItemDefinitionCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_RemoveItemDefinitionCategory);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetSupplierRequest,
 *   !proto.catalog.GetSupplierResponse>}
 */
const methodInfo_CatalogService_getSupplier = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetSupplierResponse,
  /** @param {!proto.catalog.GetSupplierRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetSupplierResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetSupplierResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetSupplierResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getSupplier =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_getSupplier,
      callback);
};


/**
 * @param {!proto.catalog.GetSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetSupplierResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getSupplier =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_getSupplier);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetSuppliersRequest,
 *   !proto.catalog.GetSuppliersResponse>}
 */
const methodInfo_CatalogService_getSuppliers = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetSuppliersResponse,
  /** @param {!proto.catalog.GetSuppliersRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetSuppliersResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetSuppliersRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetSuppliersResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetSuppliersResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getSuppliers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getSuppliers',
      request,
      metadata || {},
      methodInfo_CatalogService_getSuppliers,
      callback);
};


/**
 * @param {!proto.catalog.GetSuppliersRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetSuppliersResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getSuppliers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getSuppliers',
      request,
      metadata || {},
      methodInfo_CatalogService_getSuppliers);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetManufacturerRequest,
 *   !proto.catalog.GetManufacturerResponse>}
 */
const methodInfo_CatalogService_getManufacturer = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetManufacturerResponse,
  /** @param {!proto.catalog.GetManufacturerRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetManufacturerResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetManufacturerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetManufacturerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getManufacturer =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_getManufacturer,
      callback);
};


/**
 * @param {!proto.catalog.GetManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetManufacturerResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getManufacturer =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_getManufacturer);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetManufacturersRequest,
 *   !proto.catalog.GetManufacturersResponse>}
 */
const methodInfo_CatalogService_getManufacturers = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetManufacturersResponse,
  /** @param {!proto.catalog.GetManufacturersRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetManufacturersResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetManufacturersRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetManufacturersResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetManufacturersResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getManufacturers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getManufacturers',
      request,
      metadata || {},
      methodInfo_CatalogService_getManufacturers,
      callback);
};


/**
 * @param {!proto.catalog.GetManufacturersRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetManufacturersResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getManufacturers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getManufacturers',
      request,
      metadata || {},
      methodInfo_CatalogService_getManufacturers);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetSupplierPackagesRequest,
 *   !proto.catalog.GetSupplierPackagesResponse>}
 */
const methodInfo_CatalogService_getSupplierPackages = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetSupplierPackagesResponse,
  /** @param {!proto.catalog.GetSupplierPackagesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetSupplierPackagesResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetSupplierPackagesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetSupplierPackagesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetSupplierPackagesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getSupplierPackages =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getSupplierPackages',
      request,
      metadata || {},
      methodInfo_CatalogService_getSupplierPackages,
      callback);
};


/**
 * @param {!proto.catalog.GetSupplierPackagesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetSupplierPackagesResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getSupplierPackages =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getSupplierPackages',
      request,
      metadata || {},
      methodInfo_CatalogService_getSupplierPackages);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetPackageRequest,
 *   !proto.catalog.GetPackageResponse>}
 */
const methodInfo_CatalogService_getPackage = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetPackageResponse,
  /** @param {!proto.catalog.GetPackageRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetPackageResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetPackageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetPackageResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetPackageResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getPackage =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getPackage',
      request,
      metadata || {},
      methodInfo_CatalogService_getPackage,
      callback);
};


/**
 * @param {!proto.catalog.GetPackageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetPackageResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getPackage =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getPackage',
      request,
      metadata || {},
      methodInfo_CatalogService_getPackage);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetPackagesRequest,
 *   !proto.catalog.GetPackagesResponse>}
 */
const methodInfo_CatalogService_getPackages = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetPackagesResponse,
  /** @param {!proto.catalog.GetPackagesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetPackagesResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetPackagesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetPackagesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetPackagesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getPackages =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getPackages',
      request,
      metadata || {},
      methodInfo_CatalogService_getPackages,
      callback);
};


/**
 * @param {!proto.catalog.GetPackagesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetPackagesResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getPackages =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getPackages',
      request,
      metadata || {},
      methodInfo_CatalogService_getPackages);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetUnitOfMeasureRequest,
 *   !proto.catalog.GetUnitOfMeasureResponse>}
 */
const methodInfo_CatalogService_getUnitOfMeasure = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetUnitOfMeasureResponse,
  /** @param {!proto.catalog.GetUnitOfMeasureRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetUnitOfMeasureResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetUnitOfMeasureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetUnitOfMeasureResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetUnitOfMeasureResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getUnitOfMeasure =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getUnitOfMeasure',
      request,
      metadata || {},
      methodInfo_CatalogService_getUnitOfMeasure,
      callback);
};


/**
 * @param {!proto.catalog.GetUnitOfMeasureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetUnitOfMeasureResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getUnitOfMeasure =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getUnitOfMeasure',
      request,
      metadata || {},
      methodInfo_CatalogService_getUnitOfMeasure);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetUnitOfMeasuresRequest,
 *   !proto.catalog.GetUnitOfMeasuresResponse>}
 */
const methodInfo_CatalogService_getUnitOfMeasures = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetUnitOfMeasuresResponse,
  /** @param {!proto.catalog.GetUnitOfMeasuresRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetUnitOfMeasuresResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetUnitOfMeasuresRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetUnitOfMeasuresResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetUnitOfMeasuresResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getUnitOfMeasures =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getUnitOfMeasures',
      request,
      metadata || {},
      methodInfo_CatalogService_getUnitOfMeasures,
      callback);
};


/**
 * @param {!proto.catalog.GetUnitOfMeasuresRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetUnitOfMeasuresResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getUnitOfMeasures =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getUnitOfMeasures',
      request,
      metadata || {},
      methodInfo_CatalogService_getUnitOfMeasures);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetItemDefinitionRequest,
 *   !proto.catalog.GetItemDefinitionResponse>}
 */
const methodInfo_CatalogService_getItemDefinition = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetItemDefinitionResponse,
  /** @param {!proto.catalog.GetItemDefinitionRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetItemDefinitionResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetItemDefinitionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetItemDefinitionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetItemDefinitionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getItemDefinition =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getItemDefinition',
      request,
      metadata || {},
      methodInfo_CatalogService_getItemDefinition,
      callback);
};


/**
 * @param {!proto.catalog.GetItemDefinitionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetItemDefinitionResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getItemDefinition =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getItemDefinition',
      request,
      metadata || {},
      methodInfo_CatalogService_getItemDefinition);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetItemDefinitionsRequest,
 *   !proto.catalog.GetItemDefinitionsResponse>}
 */
const methodInfo_CatalogService_getItemDefinitions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetItemDefinitionsResponse,
  /** @param {!proto.catalog.GetItemDefinitionsRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetItemDefinitionsResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetItemDefinitionsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetItemDefinitionsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetItemDefinitionsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getItemDefinitions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getItemDefinitions',
      request,
      metadata || {},
      methodInfo_CatalogService_getItemDefinitions,
      callback);
};


/**
 * @param {!proto.catalog.GetItemDefinitionsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetItemDefinitionsResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getItemDefinitions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getItemDefinitions',
      request,
      metadata || {},
      methodInfo_CatalogService_getItemDefinitions);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetItemInstanceRequest,
 *   !proto.catalog.GetItemInstanceResponse>}
 */
const methodInfo_CatalogService_getItemInstance = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetItemInstanceResponse,
  /** @param {!proto.catalog.GetItemInstanceRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetItemInstanceResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetItemInstanceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetItemInstanceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetItemInstanceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getItemInstance =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getItemInstance',
      request,
      metadata || {},
      methodInfo_CatalogService_getItemInstance,
      callback);
};


/**
 * @param {!proto.catalog.GetItemInstanceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetItemInstanceResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getItemInstance =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getItemInstance',
      request,
      metadata || {},
      methodInfo_CatalogService_getItemInstance);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetItemInstancesRequest,
 *   !proto.catalog.GetItemInstancesResponse>}
 */
const methodInfo_CatalogService_getItemInstances = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetItemInstancesResponse,
  /** @param {!proto.catalog.GetItemInstancesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetItemInstancesResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetItemInstancesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetItemInstancesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetItemInstancesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getItemInstances =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getItemInstances',
      request,
      metadata || {},
      methodInfo_CatalogService_getItemInstances,
      callback);
};


/**
 * @param {!proto.catalog.GetItemInstancesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetItemInstancesResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getItemInstances =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getItemInstances',
      request,
      metadata || {},
      methodInfo_CatalogService_getItemInstances);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetLocalisationRequest,
 *   !proto.catalog.GetLocalisationResponse>}
 */
const methodInfo_CatalogService_getLocalisation = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetLocalisationResponse,
  /** @param {!proto.catalog.GetLocalisationRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetLocalisationResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetLocalisationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetLocalisationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetLocalisationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getLocalisation =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getLocalisation',
      request,
      metadata || {},
      methodInfo_CatalogService_getLocalisation,
      callback);
};


/**
 * @param {!proto.catalog.GetLocalisationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetLocalisationResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getLocalisation =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getLocalisation',
      request,
      metadata || {},
      methodInfo_CatalogService_getLocalisation);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetLocalisationsRequest,
 *   !proto.catalog.GetLocalisationsResponse>}
 */
const methodInfo_CatalogService_getLocalisations = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetLocalisationsResponse,
  /** @param {!proto.catalog.GetLocalisationsRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetLocalisationsResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetLocalisationsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetLocalisationsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetLocalisationsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getLocalisations =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getLocalisations',
      request,
      metadata || {},
      methodInfo_CatalogService_getLocalisations,
      callback);
};


/**
 * @param {!proto.catalog.GetLocalisationsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetLocalisationsResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getLocalisations =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getLocalisations',
      request,
      metadata || {},
      methodInfo_CatalogService_getLocalisations);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetCategoryRequest,
 *   !proto.catalog.GetCategoryResponse>}
 */
const methodInfo_CatalogService_getCategory = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetCategoryResponse,
  /** @param {!proto.catalog.GetCategoryRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetCategoryResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetCategoryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetCategoryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getCategory =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_getCategory,
      callback);
};


/**
 * @param {!proto.catalog.GetCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetCategoryResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getCategory =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_getCategory);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetCategoriesRequest,
 *   !proto.catalog.GetCategoriesResponse>}
 */
const methodInfo_CatalogService_getCategories = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetCategoriesResponse,
  /** @param {!proto.catalog.GetCategoriesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetCategoriesResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetCategoriesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetCategoriesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetCategoriesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getCategories =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getCategories',
      request,
      metadata || {},
      methodInfo_CatalogService_getCategories,
      callback);
};


/**
 * @param {!proto.catalog.GetCategoriesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetCategoriesResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getCategories =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getCategories',
      request,
      metadata || {},
      methodInfo_CatalogService_getCategories);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.GetInventoriesRequest,
 *   !proto.catalog.GetInventoriesResponse>}
 */
const methodInfo_CatalogService_getInventories = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.GetInventoriesResponse,
  /** @param {!proto.catalog.GetInventoriesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.GetInventoriesResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.GetInventoriesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.GetInventoriesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.GetInventoriesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.getInventories =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/getInventories',
      request,
      metadata || {},
      methodInfo_CatalogService_getInventories,
      callback);
};


/**
 * @param {!proto.catalog.GetInventoriesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.GetInventoriesResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.getInventories =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/getInventories',
      request,
      metadata || {},
      methodInfo_CatalogService_getInventories);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeleteInventoryRequest,
 *   !proto.catalog.DeleteInventoryResponse>}
 */
const methodInfo_CatalogService_deleteInventory = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeleteInventoryResponse,
  /** @param {!proto.catalog.DeleteInventoryRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeleteInventoryResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeleteInventoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeleteInventoryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeleteInventoryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deleteInventory =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deleteInventory',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteInventory,
      callback);
};


/**
 * @param {!proto.catalog.DeleteInventoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeleteInventoryResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deleteInventory =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deleteInventory',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteInventory);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeletePackageRequest,
 *   !proto.catalog.DeletePackageResponse>}
 */
const methodInfo_CatalogService_deletePackage = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeletePackageResponse,
  /** @param {!proto.catalog.DeletePackageRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeletePackageResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeletePackageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeletePackageResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeletePackageResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deletePackage =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deletePackage',
      request,
      metadata || {},
      methodInfo_CatalogService_deletePackage,
      callback);
};


/**
 * @param {!proto.catalog.DeletePackageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeletePackageResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deletePackage =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deletePackage',
      request,
      metadata || {},
      methodInfo_CatalogService_deletePackage);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeletePackageSupplierRequest,
 *   !proto.catalog.DeletePackageSupplierResponse>}
 */
const methodInfo_CatalogService_deletePackageSupplier = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeletePackageSupplierResponse,
  /** @param {!proto.catalog.DeletePackageSupplierRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeletePackageSupplierResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeletePackageSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeletePackageSupplierResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeletePackageSupplierResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deletePackageSupplier =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deletePackageSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_deletePackageSupplier,
      callback);
};


/**
 * @param {!proto.catalog.DeletePackageSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeletePackageSupplierResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deletePackageSupplier =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deletePackageSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_deletePackageSupplier);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeleteSupplierRequest,
 *   !proto.catalog.DeleteSupplierResponse>}
 */
const methodInfo_CatalogService_deleteSupplier = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeleteSupplierResponse,
  /** @param {!proto.catalog.DeleteSupplierRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeleteSupplierResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeleteSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeleteSupplierResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeleteSupplierResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deleteSupplier =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deleteSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteSupplier,
      callback);
};


/**
 * @param {!proto.catalog.DeleteSupplierRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeleteSupplierResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deleteSupplier =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deleteSupplier',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteSupplier);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeletePropertyDefinitionRequest,
 *   !proto.catalog.DeletePropertyDefinitionResponse>}
 */
const methodInfo_CatalogService_deletePropertyDefinition = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeletePropertyDefinitionResponse,
  /** @param {!proto.catalog.DeletePropertyDefinitionRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeletePropertyDefinitionResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeletePropertyDefinitionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeletePropertyDefinitionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeletePropertyDefinitionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deletePropertyDefinition =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deletePropertyDefinition',
      request,
      metadata || {},
      methodInfo_CatalogService_deletePropertyDefinition,
      callback);
};


/**
 * @param {!proto.catalog.DeletePropertyDefinitionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeletePropertyDefinitionResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deletePropertyDefinition =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deletePropertyDefinition',
      request,
      metadata || {},
      methodInfo_CatalogService_deletePropertyDefinition);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeleteUnitOfMeasureRequest,
 *   !proto.catalog.DeleteUnitOfMeasureResponse>}
 */
const methodInfo_CatalogService_deleteUnitOfMeasure = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeleteUnitOfMeasureResponse,
  /** @param {!proto.catalog.DeleteUnitOfMeasureRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeleteUnitOfMeasureResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeleteUnitOfMeasureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeleteUnitOfMeasureResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeleteUnitOfMeasureResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deleteUnitOfMeasure =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deleteUnitOfMeasure',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteUnitOfMeasure,
      callback);
};


/**
 * @param {!proto.catalog.DeleteUnitOfMeasureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeleteUnitOfMeasureResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deleteUnitOfMeasure =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deleteUnitOfMeasure',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteUnitOfMeasure);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeleteItemInstanceRequest,
 *   !proto.catalog.DeleteItemInstanceResponse>}
 */
const methodInfo_CatalogService_deleteItemInstance = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeleteItemInstanceResponse,
  /** @param {!proto.catalog.DeleteItemInstanceRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeleteItemInstanceResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeleteItemInstanceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeleteItemInstanceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeleteItemInstanceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deleteItemInstance =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deleteItemInstance',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteItemInstance,
      callback);
};


/**
 * @param {!proto.catalog.DeleteItemInstanceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeleteItemInstanceResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deleteItemInstance =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deleteItemInstance',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteItemInstance);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeleteManufacturerRequest,
 *   !proto.catalog.DeleteManufacturerResponse>}
 */
const methodInfo_CatalogService_deleteManufacturer = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeleteManufacturerResponse,
  /** @param {!proto.catalog.DeleteManufacturerRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeleteManufacturerResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeleteManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeleteManufacturerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeleteManufacturerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deleteManufacturer =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deleteManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteManufacturer,
      callback);
};


/**
 * @param {!proto.catalog.DeleteManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeleteManufacturerResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deleteManufacturer =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deleteManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteManufacturer);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeleteItemManufacturerRequest,
 *   !proto.catalog.DeleteItemManufacturerResponse>}
 */
const methodInfo_CatalogService_deleteItemManufacturer = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeleteItemManufacturerResponse,
  /** @param {!proto.catalog.DeleteItemManufacturerRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeleteItemManufacturerResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeleteItemManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeleteItemManufacturerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeleteItemManufacturerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deleteItemManufacturer =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deleteItemManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteItemManufacturer,
      callback);
};


/**
 * @param {!proto.catalog.DeleteItemManufacturerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeleteItemManufacturerResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deleteItemManufacturer =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deleteItemManufacturer',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteItemManufacturer);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeleteCategoryRequest,
 *   !proto.catalog.DeleteCategoryResponse>}
 */
const methodInfo_CatalogService_deleteCategory = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeleteCategoryResponse,
  /** @param {!proto.catalog.DeleteCategoryRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeleteCategoryResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeleteCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeleteCategoryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeleteCategoryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deleteCategory =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deleteCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteCategory,
      callback);
};


/**
 * @param {!proto.catalog.DeleteCategoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeleteCategoryResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deleteCategory =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deleteCategory',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteCategory);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.catalog.DeleteLocalisationRequest,
 *   !proto.catalog.DeleteLocalisationResponse>}
 */
const methodInfo_CatalogService_deleteLocalisation = new grpc.web.AbstractClientBase.MethodInfo(
  proto.catalog.DeleteLocalisationResponse,
  /** @param {!proto.catalog.DeleteLocalisationRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.catalog.DeleteLocalisationResponse.deserializeBinary
);


/**
 * @param {!proto.catalog.DeleteLocalisationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.catalog.DeleteLocalisationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.catalog.DeleteLocalisationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.catalog.CatalogServiceClient.prototype.deleteLocalisation =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/catalog.CatalogService/deleteLocalisation',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteLocalisation,
      callback);
};


/**
 * @param {!proto.catalog.DeleteLocalisationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.catalog.DeleteLocalisationResponse>}
 *     A native promise that resolves to the response
 */
proto.catalog.CatalogServicePromiseClient.prototype.deleteLocalisation =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/catalog.CatalogService/deleteLocalisation',
      request,
      metadata || {},
      methodInfo_CatalogService_deleteLocalisation);
};


module.exports = proto.catalog;

