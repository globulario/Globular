/**
 * @fileoverview gRPC-Web generated client stub for echo
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.echo = require('./search_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.echo.SearchServiceClient =
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
proto.echo.SearchServicePromiseClient =
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
 *   !proto.echo.GetVersionRequest,
 *   !proto.echo.GetVersionResponse>}
 */
const methodDescriptor_SearchService_GetVersion = new grpc.web.MethodDescriptor(
  '/echo.SearchService/GetVersion',
  grpc.web.MethodType.UNARY,
  proto.echo.GetVersionRequest,
  proto.echo.GetVersionResponse,
  /**
   * @param {!proto.echo.GetVersionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.GetVersionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.GetVersionRequest,
 *   !proto.echo.GetVersionResponse>}
 */
const methodInfo_SearchService_GetVersion = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.GetVersionResponse,
  /**
   * @param {!proto.echo.GetVersionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.GetVersionResponse.deserializeBinary
);


/**
 * @param {!proto.echo.GetVersionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.GetVersionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.GetVersionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.SearchServiceClient.prototype.getVersion =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.SearchService/GetVersion',
      request,
      metadata || {},
      methodDescriptor_SearchService_GetVersion,
      callback);
};


/**
 * @param {!proto.echo.GetVersionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.GetVersionResponse>}
 *     A native promise that resolves to the response
 */
proto.echo.SearchServicePromiseClient.prototype.getVersion =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.SearchService/GetVersion',
      request,
      metadata || {},
      methodDescriptor_SearchService_GetVersion);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.echo.IndexJsonObjectRequest,
 *   !proto.echo.IndexJsonObjectResponse>}
 */
const methodDescriptor_SearchService_IndexJsonObject = new grpc.web.MethodDescriptor(
  '/echo.SearchService/IndexJsonObject',
  grpc.web.MethodType.UNARY,
  proto.echo.IndexJsonObjectRequest,
  proto.echo.IndexJsonObjectResponse,
  /**
   * @param {!proto.echo.IndexJsonObjectRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.IndexJsonObjectResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.IndexJsonObjectRequest,
 *   !proto.echo.IndexJsonObjectResponse>}
 */
const methodInfo_SearchService_IndexJsonObject = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.IndexJsonObjectResponse,
  /**
   * @param {!proto.echo.IndexJsonObjectRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.IndexJsonObjectResponse.deserializeBinary
);


/**
 * @param {!proto.echo.IndexJsonObjectRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.IndexJsonObjectResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.IndexJsonObjectResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.SearchServiceClient.prototype.indexJsonObject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.SearchService/IndexJsonObject',
      request,
      metadata || {},
      methodDescriptor_SearchService_IndexJsonObject,
      callback);
};


/**
 * @param {!proto.echo.IndexJsonObjectRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.IndexJsonObjectResponse>}
 *     A native promise that resolves to the response
 */
proto.echo.SearchServicePromiseClient.prototype.indexJsonObject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.SearchService/IndexJsonObject',
      request,
      metadata || {},
      methodDescriptor_SearchService_IndexJsonObject);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.echo.IndexFileRequest,
 *   !proto.echo.IndexFileResponse>}
 */
const methodDescriptor_SearchService_IndexFile = new grpc.web.MethodDescriptor(
  '/echo.SearchService/IndexFile',
  grpc.web.MethodType.UNARY,
  proto.echo.IndexFileRequest,
  proto.echo.IndexFileResponse,
  /**
   * @param {!proto.echo.IndexFileRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.IndexFileResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.IndexFileRequest,
 *   !proto.echo.IndexFileResponse>}
 */
const methodInfo_SearchService_IndexFile = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.IndexFileResponse,
  /**
   * @param {!proto.echo.IndexFileRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.IndexFileResponse.deserializeBinary
);


/**
 * @param {!proto.echo.IndexFileRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.IndexFileResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.IndexFileResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.SearchServiceClient.prototype.indexFile =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.SearchService/IndexFile',
      request,
      metadata || {},
      methodDescriptor_SearchService_IndexFile,
      callback);
};


/**
 * @param {!proto.echo.IndexFileRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.IndexFileResponse>}
 *     A native promise that resolves to the response
 */
proto.echo.SearchServicePromiseClient.prototype.indexFile =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.SearchService/IndexFile',
      request,
      metadata || {},
      methodDescriptor_SearchService_IndexFile);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.echo.IndexDirRequest,
 *   !proto.echo.IndexDirResponse>}
 */
const methodDescriptor_SearchService_IndexDir = new grpc.web.MethodDescriptor(
  '/echo.SearchService/IndexDir',
  grpc.web.MethodType.UNARY,
  proto.echo.IndexDirRequest,
  proto.echo.IndexDirResponse,
  /**
   * @param {!proto.echo.IndexDirRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.IndexDirResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.IndexDirRequest,
 *   !proto.echo.IndexDirResponse>}
 */
const methodInfo_SearchService_IndexDir = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.IndexDirResponse,
  /**
   * @param {!proto.echo.IndexDirRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.IndexDirResponse.deserializeBinary
);


/**
 * @param {!proto.echo.IndexDirRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.IndexDirResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.IndexDirResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.SearchServiceClient.prototype.indexDir =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.SearchService/IndexDir',
      request,
      metadata || {},
      methodDescriptor_SearchService_IndexDir,
      callback);
};


/**
 * @param {!proto.echo.IndexDirRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.IndexDirResponse>}
 *     A native promise that resolves to the response
 */
proto.echo.SearchServicePromiseClient.prototype.indexDir =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.SearchService/IndexDir',
      request,
      metadata || {},
      methodDescriptor_SearchService_IndexDir);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.echo.CountRequest,
 *   !proto.echo.CountResponse>}
 */
const methodDescriptor_SearchService_Count = new grpc.web.MethodDescriptor(
  '/echo.SearchService/Count',
  grpc.web.MethodType.UNARY,
  proto.echo.CountRequest,
  proto.echo.CountResponse,
  /**
   * @param {!proto.echo.CountRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.CountResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.CountRequest,
 *   !proto.echo.CountResponse>}
 */
const methodInfo_SearchService_Count = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.CountResponse,
  /**
   * @param {!proto.echo.CountRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.CountResponse.deserializeBinary
);


/**
 * @param {!proto.echo.CountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.CountResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.CountResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.SearchServiceClient.prototype.count =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.SearchService/Count',
      request,
      metadata || {},
      methodDescriptor_SearchService_Count,
      callback);
};


/**
 * @param {!proto.echo.CountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.CountResponse>}
 *     A native promise that resolves to the response
 */
proto.echo.SearchServicePromiseClient.prototype.count =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.SearchService/Count',
      request,
      metadata || {},
      methodDescriptor_SearchService_Count);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.echo.DeleteDocumentRequest,
 *   !proto.echo.DeleteDocumentResponse>}
 */
const methodDescriptor_SearchService_DeleteDocument = new grpc.web.MethodDescriptor(
  '/echo.SearchService/DeleteDocument',
  grpc.web.MethodType.UNARY,
  proto.echo.DeleteDocumentRequest,
  proto.echo.DeleteDocumentResponse,
  /**
   * @param {!proto.echo.DeleteDocumentRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.DeleteDocumentResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.DeleteDocumentRequest,
 *   !proto.echo.DeleteDocumentResponse>}
 */
const methodInfo_SearchService_DeleteDocument = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.DeleteDocumentResponse,
  /**
   * @param {!proto.echo.DeleteDocumentRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.DeleteDocumentResponse.deserializeBinary
);


/**
 * @param {!proto.echo.DeleteDocumentRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.DeleteDocumentResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.DeleteDocumentResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.SearchServiceClient.prototype.deleteDocument =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.SearchService/DeleteDocument',
      request,
      metadata || {},
      methodDescriptor_SearchService_DeleteDocument,
      callback);
};


/**
 * @param {!proto.echo.DeleteDocumentRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.DeleteDocumentResponse>}
 *     A native promise that resolves to the response
 */
proto.echo.SearchServicePromiseClient.prototype.deleteDocument =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.SearchService/DeleteDocument',
      request,
      metadata || {},
      methodDescriptor_SearchService_DeleteDocument);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.echo.SearchDocumentsRequest,
 *   !proto.echo.SearchDocumentsResponse>}
 */
const methodDescriptor_SearchService_SearchDocuments = new grpc.web.MethodDescriptor(
  '/echo.SearchService/SearchDocuments',
  grpc.web.MethodType.UNARY,
  proto.echo.SearchDocumentsRequest,
  proto.echo.SearchDocumentsResponse,
  /**
   * @param {!proto.echo.SearchDocumentsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.SearchDocumentsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.echo.SearchDocumentsRequest,
 *   !proto.echo.SearchDocumentsResponse>}
 */
const methodInfo_SearchService_SearchDocuments = new grpc.web.AbstractClientBase.MethodInfo(
  proto.echo.SearchDocumentsResponse,
  /**
   * @param {!proto.echo.SearchDocumentsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.echo.SearchDocumentsResponse.deserializeBinary
);


/**
 * @param {!proto.echo.SearchDocumentsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.echo.SearchDocumentsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.echo.SearchDocumentsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.echo.SearchServiceClient.prototype.searchDocuments =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/echo.SearchService/SearchDocuments',
      request,
      metadata || {},
      methodDescriptor_SearchService_SearchDocuments,
      callback);
};


/**
 * @param {!proto.echo.SearchDocumentsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.echo.SearchDocumentsResponse>}
 *     A native promise that resolves to the response
 */
proto.echo.SearchServicePromiseClient.prototype.searchDocuments =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/echo.SearchService/SearchDocuments',
      request,
      metadata || {},
      methodDescriptor_SearchService_SearchDocuments);
};


module.exports = proto.echo;

