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
 *   !proto.ressource.SynchronizeLdapRqst,
 *   !proto.ressource.SynchronizeLdapRsp>}
 */
const methodInfo_RessourceService_SynchronizeLdap = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.SynchronizeLdapRsp,
  /** @param {!proto.ressource.SynchronizeLdapRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SynchronizeLdapRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.SynchronizeLdapRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.SynchronizeLdapRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.SynchronizeLdapRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.synchronizeLdap =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/SynchronizeLdap',
      request,
      metadata || {},
      methodInfo_RessourceService_SynchronizeLdap,
      callback);
};


/**
 * @param {!proto.ressource.SynchronizeLdapRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.SynchronizeLdapRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.synchronizeLdap =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/SynchronizeLdap',
      request,
      metadata || {},
      methodInfo_RessourceService_SynchronizeLdap);
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


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.AddApplicationActionRqst,
 *   !proto.ressource.AddApplicationActionRsp>}
 */
const methodInfo_RessourceService_AddApplicationAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.AddApplicationActionRsp,
  /** @param {!proto.ressource.AddApplicationActionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AddApplicationActionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.AddApplicationActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.AddApplicationActionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.AddApplicationActionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.addApplicationAction =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/AddApplicationAction',
      request,
      metadata || {},
      methodInfo_RessourceService_AddApplicationAction,
      callback);
};


/**
 * @param {!proto.ressource.AddApplicationActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.AddApplicationActionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.addApplicationAction =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/AddApplicationAction',
      request,
      metadata || {},
      methodInfo_RessourceService_AddApplicationAction);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RemoveApplicationActionRqst,
 *   !proto.ressource.RemoveApplicationActionRsp>}
 */
const methodInfo_RessourceService_RemoveApplicationAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemoveApplicationActionRsp,
  /** @param {!proto.ressource.RemoveApplicationActionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveApplicationActionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RemoveApplicationActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RemoveApplicationActionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RemoveApplicationActionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.removeApplicationAction =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RemoveApplicationAction',
      request,
      metadata || {},
      methodInfo_RessourceService_RemoveApplicationAction,
      callback);
};


/**
 * @param {!proto.ressource.RemoveApplicationActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RemoveApplicationActionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.removeApplicationAction =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RemoveApplicationAction',
      request,
      metadata || {},
      methodInfo_RessourceService_RemoveApplicationAction);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetAllActionsRqst,
 *   !proto.ressource.GetAllActionsRsp>}
 */
const methodInfo_RessourceService_GetAllActions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetAllActionsRsp,
  /** @param {!proto.ressource.GetAllActionsRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetAllActionsRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.GetAllActionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.GetAllActionsRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetAllActionsRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.getAllActions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/GetAllActions',
      request,
      metadata || {},
      methodInfo_RessourceService_GetAllActions,
      callback);
};


/**
 * @param {!proto.ressource.GetAllActionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.GetAllActionsRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.getAllActions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/GetAllActions',
      request,
      metadata || {},
      methodInfo_RessourceService_GetAllActions);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetPermissionsRqst,
 *   !proto.ressource.GetPermissionsRsp>}
 */
const methodInfo_RessourceService_GetPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetPermissionsRsp,
  /** @param {!proto.ressource.GetPermissionsRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetPermissionsRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.GetPermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.GetPermissionsRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetPermissionsRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.getPermissions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/GetPermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_GetPermissions,
      callback);
};


/**
 * @param {!proto.ressource.GetPermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.GetPermissionsRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.getPermissions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/GetPermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_GetPermissions);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.SetPermissionRqst,
 *   !proto.ressource.SetPermissionRsp>}
 */
const methodInfo_RessourceService_SetPermission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.SetPermissionRsp,
  /** @param {!proto.ressource.SetPermissionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SetPermissionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.SetPermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.SetPermissionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.SetPermissionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.setPermission =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/SetPermission',
      request,
      metadata || {},
      methodInfo_RessourceService_SetPermission,
      callback);
};


/**
 * @param {!proto.ressource.SetPermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.SetPermissionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.setPermission =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/SetPermission',
      request,
      metadata || {},
      methodInfo_RessourceService_SetPermission);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeletePermissionsRqst,
 *   !proto.ressource.DeletePermissionsRsp>}
 */
const methodInfo_RessourceService_DeletePermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeletePermissionsRsp,
  /** @param {!proto.ressource.DeletePermissionsRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeletePermissionsRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeletePermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeletePermissionsRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeletePermissionsRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deletePermissions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeletePermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeletePermissions,
      callback);
};


/**
 * @param {!proto.ressource.DeletePermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeletePermissionsRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deletePermissions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeletePermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeletePermissions);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.SetRessourceOwnerRqst,
 *   !proto.ressource.SetRessourceOwnerRsp>}
 */
const methodInfo_RessourceService_SetRessourceOwner = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.SetRessourceOwnerRsp,
  /** @param {!proto.ressource.SetRessourceOwnerRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SetRessourceOwnerRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.SetRessourceOwnerRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.SetRessourceOwnerRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.SetRessourceOwnerRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.setRessourceOwner =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/SetRessourceOwner',
      request,
      metadata || {},
      methodInfo_RessourceService_SetRessourceOwner,
      callback);
};


/**
 * @param {!proto.ressource.SetRessourceOwnerRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.SetRessourceOwnerRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.setRessourceOwner =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/SetRessourceOwner',
      request,
      metadata || {},
      methodInfo_RessourceService_SetRessourceOwner);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetRessourceOwnersRqst,
 *   !proto.ressource.GetRessourceOwnersRsp>}
 */
const methodInfo_RessourceService_GetRessourceOwners = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetRessourceOwnersRsp,
  /** @param {!proto.ressource.GetRessourceOwnersRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetRessourceOwnersRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.GetRessourceOwnersRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.GetRessourceOwnersRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetRessourceOwnersRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.getRessourceOwners =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/GetRessourceOwners',
      request,
      metadata || {},
      methodInfo_RessourceService_GetRessourceOwners,
      callback);
};


/**
 * @param {!proto.ressource.GetRessourceOwnersRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.GetRessourceOwnersRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.getRessourceOwners =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/GetRessourceOwners',
      request,
      metadata || {},
      methodInfo_RessourceService_GetRessourceOwners);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteRessourceOwnerRqst,
 *   !proto.ressource.DeleteRessourceOwnerRsp>}
 */
const methodInfo_RessourceService_DeleteRessourceOwner = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteRessourceOwnerRsp,
  /** @param {!proto.ressource.DeleteRessourceOwnerRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteRessourceOwnerRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteRessourceOwnerRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteRessourceOwnerRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteRessourceOwnerRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteRessourceOwner =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteRessourceOwner',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteRessourceOwner,
      callback);
};


/**
 * @param {!proto.ressource.DeleteRessourceOwnerRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteRessourceOwnerRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteRessourceOwner =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteRessourceOwner',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteRessourceOwner);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteRessourceOwnersRqst,
 *   !proto.ressource.DeleteRessourceOwnersRsp>}
 */
const methodInfo_RessourceService_DeleteRessourceOwners = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteRessourceOwnersRsp,
  /** @param {!proto.ressource.DeleteRessourceOwnersRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteRessourceOwnersRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteRessourceOwnersRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteRessourceOwnersRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteRessourceOwnersRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteRessourceOwners =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteRessourceOwners',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteRessourceOwners,
      callback);
};


/**
 * @param {!proto.ressource.DeleteRessourceOwnersRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteRessourceOwnersRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteRessourceOwners =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteRessourceOwners',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteRessourceOwners);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetAllFilesInfoRqst,
 *   !proto.ressource.GetAllFilesInfoRsp>}
 */
const methodInfo_RessourceService_GetAllFilesInfo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetAllFilesInfoRsp,
  /** @param {!proto.ressource.GetAllFilesInfoRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetAllFilesInfoRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.GetAllFilesInfoRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.GetAllFilesInfoRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetAllFilesInfoRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.getAllFilesInfo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/GetAllFilesInfo',
      request,
      metadata || {},
      methodInfo_RessourceService_GetAllFilesInfo,
      callback);
};


/**
 * @param {!proto.ressource.GetAllFilesInfoRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.GetAllFilesInfoRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.getAllFilesInfo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/GetAllFilesInfo',
      request,
      metadata || {},
      methodInfo_RessourceService_GetAllFilesInfo);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidateUserFileAccessRqst,
 *   !proto.ressource.ValidateUserFileAccessRsp>}
 */
const methodInfo_RessourceService_ValidateUserFileAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidateUserFileAccessRsp,
  /** @param {!proto.ressource.ValidateUserFileAccessRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateUserFileAccessRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ValidateUserFileAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ValidateUserFileAccessRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ValidateUserFileAccessRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.validateUserFileAccess =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ValidateUserFileAccess',
      request,
      metadata || {},
      methodInfo_RessourceService_ValidateUserFileAccess,
      callback);
};


/**
 * @param {!proto.ressource.ValidateUserFileAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ValidateUserFileAccessRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.validateUserFileAccess =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ValidateUserFileAccess',
      request,
      metadata || {},
      methodInfo_RessourceService_ValidateUserFileAccess);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidateApplicationFileAccessRqst,
 *   !proto.ressource.ValidateApplicationFileAccessRsp>}
 */
const methodInfo_RessourceService_ValidateApplicationFileAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidateApplicationFileAccessRsp,
  /** @param {!proto.ressource.ValidateApplicationFileAccessRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateApplicationFileAccessRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ValidateApplicationFileAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ValidateApplicationFileAccessRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ValidateApplicationFileAccessRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.validateApplicationFileAccess =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ValidateApplicationFileAccess',
      request,
      metadata || {},
      methodInfo_RessourceService_ValidateApplicationFileAccess,
      callback);
};


/**
 * @param {!proto.ressource.ValidateApplicationFileAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ValidateApplicationFileAccessRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.validateApplicationFileAccess =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ValidateApplicationFileAccess',
      request,
      metadata || {},
      methodInfo_RessourceService_ValidateApplicationFileAccess);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidateUserAccessRqst,
 *   !proto.ressource.ValidateUserAccessRsp>}
 */
const methodInfo_RessourceService_ValidateUserAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidateUserAccessRsp,
  /** @param {!proto.ressource.ValidateUserAccessRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateUserAccessRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ValidateUserAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ValidateUserAccessRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ValidateUserAccessRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.validateUserAccess =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ValidateUserAccess',
      request,
      metadata || {},
      methodInfo_RessourceService_ValidateUserAccess,
      callback);
};


/**
 * @param {!proto.ressource.ValidateUserAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ValidateUserAccessRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.validateUserAccess =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ValidateUserAccess',
      request,
      metadata || {},
      methodInfo_RessourceService_ValidateUserAccess);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidateApplicationAccessRqst,
 *   !proto.ressource.ValidateApplicationAccessRsp>}
 */
const methodInfo_RessourceService_ValidateApplicationAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidateApplicationAccessRsp,
  /** @param {!proto.ressource.ValidateApplicationAccessRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateApplicationAccessRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ValidateApplicationAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ValidateApplicationAccessRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ValidateApplicationAccessRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.validateApplicationAccess =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ValidateApplicationAccess',
      request,
      metadata || {},
      methodInfo_RessourceService_ValidateApplicationAccess,
      callback);
};


/**
 * @param {!proto.ressource.ValidateApplicationAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ValidateApplicationAccessRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.validateApplicationAccess =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ValidateApplicationAccess',
      request,
      metadata || {},
      methodInfo_RessourceService_ValidateApplicationAccess);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.CreateDirPermissionsRqst,
 *   !proto.ressource.CreateDirPermissionsRsp>}
 */
const methodInfo_RessourceService_CreateDirPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.CreateDirPermissionsRsp,
  /** @param {!proto.ressource.CreateDirPermissionsRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.CreateDirPermissionsRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.CreateDirPermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.CreateDirPermissionsRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.CreateDirPermissionsRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.createDirPermissions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/CreateDirPermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_CreateDirPermissions,
      callback);
};


/**
 * @param {!proto.ressource.CreateDirPermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.CreateDirPermissionsRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.createDirPermissions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/CreateDirPermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_CreateDirPermissions);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RenameFilePermissionRqst,
 *   !proto.ressource.RenameFilePermissionRsp>}
 */
const methodInfo_RessourceService_RenameFilePermission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RenameFilePermissionRsp,
  /** @param {!proto.ressource.RenameFilePermissionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RenameFilePermissionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RenameFilePermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RenameFilePermissionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RenameFilePermissionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.renameFilePermission =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RenameFilePermission',
      request,
      metadata || {},
      methodInfo_RessourceService_RenameFilePermission,
      callback);
};


/**
 * @param {!proto.ressource.RenameFilePermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RenameFilePermissionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.renameFilePermission =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RenameFilePermission',
      request,
      metadata || {},
      methodInfo_RessourceService_RenameFilePermission);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteDirPermissionsRqst,
 *   !proto.ressource.DeleteDirPermissionsRsp>}
 */
const methodInfo_RessourceService_DeleteDirPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteDirPermissionsRsp,
  /** @param {!proto.ressource.DeleteDirPermissionsRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteDirPermissionsRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteDirPermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteDirPermissionsRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteDirPermissionsRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteDirPermissions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteDirPermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteDirPermissions,
      callback);
};


/**
 * @param {!proto.ressource.DeleteDirPermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteDirPermissionsRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteDirPermissions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteDirPermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteDirPermissions);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteFilePermissionsRqst,
 *   !proto.ressource.DeleteFilePermissionsRsp>}
 */
const methodInfo_RessourceService_DeleteFilePermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteFilePermissionsRsp,
  /** @param {!proto.ressource.DeleteFilePermissionsRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteFilePermissionsRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteFilePermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteFilePermissionsRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteFilePermissionsRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteFilePermissions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteFilePermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteFilePermissions,
      callback);
};


/**
 * @param {!proto.ressource.DeleteFilePermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteFilePermissionsRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteFilePermissions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteFilePermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteFilePermissions);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteAccountPermissionsRqst,
 *   !proto.ressource.DeleteAccountPermissionsRsp>}
 */
const methodInfo_RessourceService_DeleteAccountPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteAccountPermissionsRsp,
  /** @param {!proto.ressource.DeleteAccountPermissionsRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteAccountPermissionsRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteAccountPermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteAccountPermissionsRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteAccountPermissionsRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteAccountPermissions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteAccountPermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteAccountPermissions,
      callback);
};


/**
 * @param {!proto.ressource.DeleteAccountPermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteAccountPermissionsRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteAccountPermissions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteAccountPermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteAccountPermissions);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteRolePermissionsRqst,
 *   !proto.ressource.DeleteRolePermissionsRsp>}
 */
const methodInfo_RessourceService_DeleteRolePermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteRolePermissionsRsp,
  /** @param {!proto.ressource.DeleteRolePermissionsRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteRolePermissionsRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteRolePermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteRolePermissionsRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteRolePermissionsRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteRolePermissions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteRolePermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteRolePermissions,
      callback);
};


/**
 * @param {!proto.ressource.DeleteRolePermissionsRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteRolePermissionsRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteRolePermissions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteRolePermissions',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteRolePermissions);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetAllApplicationsInfoRqst,
 *   !proto.ressource.GetAllApplicationsInfoRsp>}
 */
const methodInfo_RessourceService_GetAllApplicationsInfo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetAllApplicationsInfoRsp,
  /** @param {!proto.ressource.GetAllApplicationsInfoRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetAllApplicationsInfoRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.GetAllApplicationsInfoRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.GetAllApplicationsInfoRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetAllApplicationsInfoRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.getAllApplicationsInfo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/GetAllApplicationsInfo',
      request,
      metadata || {},
      methodInfo_RessourceService_GetAllApplicationsInfo,
      callback);
};


/**
 * @param {!proto.ressource.GetAllApplicationsInfoRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.GetAllApplicationsInfoRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.getAllApplicationsInfo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/GetAllApplicationsInfo',
      request,
      metadata || {},
      methodInfo_RessourceService_GetAllApplicationsInfo);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteApplicationRqst,
 *   !proto.ressource.DeleteApplicationRsp>}
 */
const methodInfo_RessourceService_DeleteApplication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteApplicationRsp,
  /** @param {!proto.ressource.DeleteApplicationRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteApplicationRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteApplicationRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteApplicationRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteApplicationRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteApplication',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteApplication,
      callback);
};


/**
 * @param {!proto.ressource.DeleteApplicationRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteApplicationRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteApplication',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteApplication);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.LogRqst,
 *   !proto.ressource.LogRsp>}
 */
const methodInfo_RessourceService_Log = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.LogRsp,
  /** @param {!proto.ressource.LogRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.LogRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.LogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.LogRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.LogRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.log =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/Log',
      request,
      metadata || {},
      methodInfo_RessourceService_Log,
      callback);
};


/**
 * @param {!proto.ressource.LogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.LogRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.log =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/Log',
      request,
      metadata || {},
      methodInfo_RessourceService_Log);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.SetLogRqst,
 *   !proto.ressource.SetLogRsp>}
 */
const methodInfo_RessourceService_SetLog = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.SetLogRsp,
  /** @param {!proto.ressource.SetLogRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SetLogRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.SetLogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.SetLogRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.SetLogRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.setLog =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/SetLog',
      request,
      metadata || {},
      methodInfo_RessourceService_SetLog,
      callback);
};


/**
 * @param {!proto.ressource.SetLogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.SetLogRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.setLog =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/SetLog',
      request,
      metadata || {},
      methodInfo_RessourceService_SetLog);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ResetLogRqst,
 *   !proto.ressource.ResetLogRsp>}
 */
const methodInfo_RessourceService_ResetLog = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ResetLogRsp,
  /** @param {!proto.ressource.ResetLogRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ResetLogRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ResetLogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ResetLogRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ResetLogRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.resetLog =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ResetLog',
      request,
      metadata || {},
      methodInfo_RessourceService_ResetLog,
      callback);
};


/**
 * @param {!proto.ressource.ResetLogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ResetLogRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.resetLog =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ResetLog',
      request,
      metadata || {},
      methodInfo_RessourceService_ResetLog);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetLogRqst,
 *   !proto.ressource.GetLogRsp>}
 */
const methodInfo_RessourceService_GetLog = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetLogRsp,
  /** @param {!proto.ressource.GetLogRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetLogRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.GetLogRqst} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetLogRsp>}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.getLog =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/ressource.RessourceService/GetLog',
      request,
      metadata || {},
      methodInfo_RessourceService_GetLog);
};


/**
 * @param {!proto.ressource.GetLogRqst} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetLogRsp>}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServicePromiseClient.prototype.getLog =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/ressource.RessourceService/GetLog',
      request,
      metadata || {},
      methodInfo_RessourceService_GetLog);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteLogRqst,
 *   !proto.ressource.DeleteLogRsp>}
 */
const methodInfo_RessourceService_DeleteLog = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteLogRsp,
  /** @param {!proto.ressource.DeleteLogRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteLogRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeleteLogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeleteLogRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeleteLogRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deleteLog =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeleteLog',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteLog,
      callback);
};


/**
 * @param {!proto.ressource.DeleteLogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeleteLogRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deleteLog =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeleteLog',
      request,
      metadata || {},
      methodInfo_RessourceService_DeleteLog);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ClearAllLogRqst,
 *   !proto.ressource.ClearAllLogRsp>}
 */
const methodInfo_RessourceService_ClearAllLog = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ClearAllLogRsp,
  /** @param {!proto.ressource.ClearAllLogRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ClearAllLogRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ClearAllLogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ClearAllLogRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ClearAllLogRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.clearAllLog =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ClearAllLog',
      request,
      metadata || {},
      methodInfo_RessourceService_ClearAllLog,
      callback);
};


/**
 * @param {!proto.ressource.ClearAllLogRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ClearAllLogRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.clearAllLog =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ClearAllLog',
      request,
      metadata || {},
      methodInfo_RessourceService_ClearAllLog);
};


module.exports = proto.ressource;

