/**
 * @fileoverview gRPC-Web generated client stub for ressource
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.ressource = require('./ressource_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.ressource.RessourceServiceClient =
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
proto.ressource.RessourceServicePromiseClient =
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
 *   !proto.ressource.RegisterAccountRqst,
 *   !proto.ressource.RegisterAccountRsp>}
 */
const methodInfo_RessourceService_RegisterAccount = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RegisterAccountRsp,
  /** @param {!proto.ressource.RegisterAccountRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RegisterAccountRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RegisterAccountRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RegisterAccountRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RegisterAccountRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.registerAccount =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RegisterAccount',
      request,
      metadata || {},
      methodInfo_RessourceService_RegisterAccount,
      callback);
};


/**
 * @param {!proto.ressource.RegisterAccountRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RegisterAccountRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.registerAccount =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RegisterAccount',
      request,
      metadata || {},
      methodInfo_RessourceService_RegisterAccount);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteAccountRqst,
 *   !proto.ressource.DeleteAccountRsp>}
 */
const methodInfo_RessourceService_DeleteAccount = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteAccountRsp,
  /** @param {!proto.ressource.DeleteAccountRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteAccountRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteAccountRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteAccountRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteAccountRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteAccount =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteAccount',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteAccount,
      callback);
};


/**
 * @param {!proto.ressource.DeleteAccountRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteAccountRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteAccount =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteAccount',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteAccount);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.AuthenticateRqst,
 *   !proto.ressource.AuthenticateRsp>}
 */
const methodInfo_RessourceService_Authenticate = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.AuthenticateRsp,
  /** @param {!proto.ressource.AuthenticateRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AuthenticateRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.AuthenticateRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.AuthenticateRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.AuthenticateRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.authenticate =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/Authenticate',
      request,
      metadata || {},
      methodInfo_RessourceService_Authenticate,
      callback);
};


/**
 * @param {!proto.ressource.AuthenticateRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.AuthenticateRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.authenticate =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/Authenticate',
      request,
      metadata || {},
      methodInfo_RessourceService_Authenticate);
};


module.exports = proto.ressource;

