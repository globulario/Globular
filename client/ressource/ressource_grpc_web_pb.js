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


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RefreshTokenRqst,
 *   !proto.ressource.RefreshTokenRsp>}
 */
const methodInfo_RessourceService_RefreshToken = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RefreshTokenRsp,
  /** @param {!proto.ressource.RefreshTokenRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RefreshTokenRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RefreshTokenRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RefreshTokenRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RefreshTokenRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.refreshToken =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RefreshToken',
      request,
      metadata || {},
      methodInfo_RessourceService_RefreshToken,
      callback);
};


/**
 * @param {!proto.ressource.RefreshTokenRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RefreshTokenRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.refreshToken =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RefreshToken',
      request,
      metadata || {},
      methodInfo_RessourceService_RefreshToken);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.AddAccountRoleRqst,
 *   !proto.ressource.AddAccountRoleRsp>}
 */
const methodInfo_RessourceService_AddAccountRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.AddAccountRoleRsp,
  /** @param {!proto.ressource.AddAccountRoleRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AddAccountRoleRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.AddAccountRoleRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.AddAccountRoleRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.AddAccountRoleRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.addAccountRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/AddAccountRole',
      request,
      metadata || {},
      methodInfo_RessourceService_AddAccountRole,
      callback);
};


/**
 * @param {!proto.ressource.AddAccountRoleRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.AddAccountRoleRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.addAccountRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/AddAccountRole',
      request,
      metadata || {},
      methodInfo_RessourceService_AddAccountRole);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RemoveAccountRoleRqst,
 *   !proto.ressource.RemoveAccountRoleRsp>}
 */
const methodInfo_RessourceService_RemoveAccountRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemoveAccountRoleRsp,
  /** @param {!proto.ressource.RemoveAccountRoleRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveAccountRoleRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RemoveAccountRoleRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RemoveAccountRoleRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RemoveAccountRoleRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.removeAccountRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RemoveAccountRole',
      request,
      metadata || {},
      methodInfo_RessourceService_RemoveAccountRole,
      callback);
};


/**
 * @param {!proto.ressource.RemoveAccountRoleRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RemoveAccountRoleRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.removeAccountRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RemoveAccountRole',
      request,
      metadata || {},
      methodInfo_RessourceService_RemoveAccountRole);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.CreateRoleRqst,
 *   !proto.ressource.CreateRoleRsp>}
 */
const methodInfo_RessourceService_CreateRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.CreateRoleRsp,
  /** @param {!proto.ressource.CreateRoleRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.CreateRoleRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.CreateRoleRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.CreateRoleRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.CreateRoleRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.createRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/CreateRole',
      request,
      metadata || {},
      methodInfo_RessourceService_CreateRole,
      callback);
};


/**
 * @param {!proto.ressource.CreateRoleRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.CreateRoleRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.createRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/CreateRole',
      request,
      metadata || {},
      methodInfo_RessourceService_CreateRole);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteRoleRqst,
 *   !proto.ressource.DeleteRoleRsp>}
 */
const methodInfo_RessourceService_DeleteRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteRoleRsp,
  /** @param {!proto.ressource.DeleteRoleRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteRoleRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteRoleRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteRoleRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteRoleRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteRole',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteRole,
      callback);
};


/**
 * @param {!proto.ressource.DeleteRoleRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteRoleRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteRole',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteRole);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.AddRoleActionRqst,
 *   !proto.ressource.AddRoleActionRsp>}
 */
const methodInfo_RessourceService_AddRoleAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.AddRoleActionRsp,
  /** @param {!proto.ressource.AddRoleActionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AddRoleActionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.AddRoleActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.AddRoleActionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.AddRoleActionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.addRoleAction =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/AddRoleAction',
      request,
      metadata || {},
      methodInfo_RessourceService_AddRoleAction,
      callback);
};


/**
 * @param {!proto.ressource.AddRoleActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.AddRoleActionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.addRoleAction =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/AddRoleAction',
      request,
      metadata || {},
      methodInfo_RessourceService_AddRoleAction);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RemoveRoleActionRqst,
 *   !proto.ressource.RemoveRoleActionRsp>}
 */
const methodInfo_RessourceService_RemoveRoleAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemoveRoleActionRsp,
  /** @param {!proto.ressource.RemoveRoleActionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveRoleActionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RemoveRoleActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RemoveRoleActionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RemoveRoleActionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.removeRoleAction =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RemoveRoleAction',
      request,
      metadata || {},
      methodInfo_RessourceService_RemoveRoleAction,
      callback);
};


/**
 * @param {!proto.ressource.RemoveRoleActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RemoveRoleActionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.removeRoleAction =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RemoveRoleAction',
      request,
      metadata || {},
      methodInfo_RessourceService_RemoveRoleAction);
};


module.exports = proto.ressource;

