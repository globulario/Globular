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
 *   !proto.ressource.RemoveApplicationRqst,
 *   !proto.ressource.RemoveApplicationRsp>}
 */
const methodInfo_RessourceService_RemoveApplication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemoveApplicationRsp,
  /** @param {!proto.ressource.RemoveApplicationRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveApplicationRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RemoveApplicationRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RemoveApplicationRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RemoveApplicationRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.removeApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RemoveApplication',
      request,
      metadata || {},
      methodInfo_RessourceService_RemoveApplication,
      callback);
};


/**
 * @param {!proto.ressource.RemoveApplicationRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RemoveApplicationRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.removeApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RemoveApplication',
      request,
      metadata || {},
      methodInfo_RessourceService_RemoveApplication);
};


module.exports = proto.ressource;

