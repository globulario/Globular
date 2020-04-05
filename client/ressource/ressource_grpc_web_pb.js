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

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RegisterPeerRqst,
 *   !proto.ressource.RegisterPeerRsp>}
 */
const methodDescriptor_RessourceService_RegisterPeer = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RegisterPeer',
  grpc.web.MethodType.UNARY,
  proto.ressource.RegisterPeerRqst,
  proto.ressource.RegisterPeerRsp,
  /**
   * @param {!proto.ressource.RegisterPeerRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RegisterPeerRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RegisterPeerRqst,
 *   !proto.ressource.RegisterPeerRsp>}
 */
const methodInfo_RessourceService_RegisterPeer = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RegisterPeerRsp,
  /**
   * @param {!proto.ressource.RegisterPeerRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RegisterPeerRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RegisterPeerRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RegisterPeerRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RegisterPeerRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.registerPeer =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RegisterPeer',
      request,
      metadata || {},
      methodDescriptor_RessourceService_RegisterPeer,
      callback);
};


/**
 * @param {!proto.ressource.RegisterPeerRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RegisterPeerRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.registerPeer =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RegisterPeer',
      request,
      metadata || {},
      methodDescriptor_RessourceService_RegisterPeer);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.GetPeersRqst,
 *   !proto.ressource.GetPeersRsp>}
 */
const methodDescriptor_RessourceService_GetPeers = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/GetPeers',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.ressource.GetPeersRqst,
  proto.ressource.GetPeersRsp,
  /**
   * @param {!proto.ressource.GetPeersRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetPeersRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetPeersRqst,
 *   !proto.ressource.GetPeersRsp>}
 */
const methodInfo_RessourceService_GetPeers = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetPeersRsp,
  /**
   * @param {!proto.ressource.GetPeersRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetPeersRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.GetPeersRqst} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetPeersRsp>}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.getPeers =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/ressource.RessourceService/GetPeers',
      request,
      metadata || {},
      methodDescriptor_RessourceService_GetPeers);
};


/**
 * @param {!proto.ressource.GetPeersRqst} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetPeersRsp>}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServicePromiseClient.prototype.getPeers =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/ressource.RessourceService/GetPeers',
      request,
      metadata || {},
      methodDescriptor_RessourceService_GetPeers);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeletePeerRqst,
 *   !proto.ressource.DeletePeerRsp>}
 */
const methodDescriptor_RessourceService_DeletePeer = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeletePeer',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeletePeerRqst,
  proto.ressource.DeletePeerRsp,
  /**
   * @param {!proto.ressource.DeletePeerRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeletePeerRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeletePeerRqst,
 *   !proto.ressource.DeletePeerRsp>}
 */
const methodInfo_RessourceService_DeletePeer = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeletePeerRsp,
  /**
   * @param {!proto.ressource.DeletePeerRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeletePeerRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.DeletePeerRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.DeletePeerRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.DeletePeerRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.deletePeer =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/DeletePeer',
      request,
      metadata || {},
      methodDescriptor_RessourceService_DeletePeer,
      callback);
};


/**
 * @param {!proto.ressource.DeletePeerRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.DeletePeerRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.deletePeer =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/DeletePeer',
      request,
      metadata || {},
      methodDescriptor_RessourceService_DeletePeer);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.AddPeerActionRqst,
 *   !proto.ressource.AddPeerActionRsp>}
 */
const methodDescriptor_RessourceService_AddPeerAction = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/AddPeerAction',
  grpc.web.MethodType.UNARY,
  proto.ressource.AddPeerActionRqst,
  proto.ressource.AddPeerActionRsp,
  /**
   * @param {!proto.ressource.AddPeerActionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AddPeerActionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.AddPeerActionRqst,
 *   !proto.ressource.AddPeerActionRsp>}
 */
const methodInfo_RessourceService_AddPeerAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.AddPeerActionRsp,
  /**
   * @param {!proto.ressource.AddPeerActionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AddPeerActionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.AddPeerActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.AddPeerActionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.AddPeerActionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.addPeerAction =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/AddPeerAction',
      request,
      metadata || {},
      methodDescriptor_RessourceService_AddPeerAction,
      callback);
};


/**
 * @param {!proto.ressource.AddPeerActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.AddPeerActionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.addPeerAction =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/AddPeerAction',
      request,
      metadata || {},
      methodDescriptor_RessourceService_AddPeerAction);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RemovePeerActionRqst,
 *   !proto.ressource.RemovePeerActionRsp>}
 */
const methodDescriptor_RessourceService_RemovePeerAction = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RemovePeerAction',
  grpc.web.MethodType.UNARY,
  proto.ressource.RemovePeerActionRqst,
  proto.ressource.RemovePeerActionRsp,
  /**
   * @param {!proto.ressource.RemovePeerActionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemovePeerActionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RemovePeerActionRqst,
 *   !proto.ressource.RemovePeerActionRsp>}
 */
const methodInfo_RessourceService_RemovePeerAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemovePeerActionRsp,
  /**
   * @param {!proto.ressource.RemovePeerActionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemovePeerActionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RemovePeerActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RemovePeerActionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RemovePeerActionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.removePeerAction =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RemovePeerAction',
      request,
      metadata || {},
      methodDescriptor_RessourceService_RemovePeerAction,
      callback);
};


/**
 * @param {!proto.ressource.RemovePeerActionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RemovePeerActionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.removePeerAction =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RemovePeerAction',
      request,
      metadata || {},
      methodDescriptor_RessourceService_RemovePeerAction);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RegisterAccountRqst,
 *   !proto.ressource.RegisterAccountRsp>}
 */
const methodDescriptor_RessourceService_RegisterAccount = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RegisterAccount',
  grpc.web.MethodType.UNARY,
  proto.ressource.RegisterAccountRqst,
  proto.ressource.RegisterAccountRsp,
  /**
   * @param {!proto.ressource.RegisterAccountRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RegisterAccountRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RegisterAccountRqst,
 *   !proto.ressource.RegisterAccountRsp>}
 */
const methodInfo_RessourceService_RegisterAccount = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RegisterAccountRsp,
  /**
   * @param {!proto.ressource.RegisterAccountRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_RegisterAccount,
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
      methodDescriptor_RessourceService_RegisterAccount);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteAccountRqst,
 *   !proto.ressource.DeleteAccountRsp>}
 */
const methodDescriptor_RessourceService_DeleteAccount = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteAccount',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteAccountRqst,
  proto.ressource.DeleteAccountRsp,
  /**
   * @param {!proto.ressource.DeleteAccountRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteAccountRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteAccountRqst,
 *   !proto.ressource.DeleteAccountRsp>}
 */
const methodInfo_RessourceService_DeleteAccount = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteAccountRsp,
  /**
   * @param {!proto.ressource.DeleteAccountRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteAccount,
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
      methodDescriptor_RessourceService_DeleteAccount);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.AuthenticateRqst,
 *   !proto.ressource.AuthenticateRsp>}
 */
const methodDescriptor_RessourceService_Authenticate = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/Authenticate',
  grpc.web.MethodType.UNARY,
  proto.ressource.AuthenticateRqst,
  proto.ressource.AuthenticateRsp,
  /**
   * @param {!proto.ressource.AuthenticateRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AuthenticateRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.AuthenticateRqst,
 *   !proto.ressource.AuthenticateRsp>}
 */
const methodInfo_RessourceService_Authenticate = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.AuthenticateRsp,
  /**
   * @param {!proto.ressource.AuthenticateRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_Authenticate,
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
      methodDescriptor_RessourceService_Authenticate);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.SynchronizeLdapRqst,
 *   !proto.ressource.SynchronizeLdapRsp>}
 */
const methodDescriptor_RessourceService_SynchronizeLdap = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/SynchronizeLdap',
  grpc.web.MethodType.UNARY,
  proto.ressource.SynchronizeLdapRqst,
  proto.ressource.SynchronizeLdapRsp,
  /**
   * @param {!proto.ressource.SynchronizeLdapRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SynchronizeLdapRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.SynchronizeLdapRqst,
 *   !proto.ressource.SynchronizeLdapRsp>}
 */
const methodInfo_RessourceService_SynchronizeLdap = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.SynchronizeLdapRsp,
  /**
   * @param {!proto.ressource.SynchronizeLdapRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_SynchronizeLdap,
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
      methodDescriptor_RessourceService_SynchronizeLdap);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RefreshTokenRqst,
 *   !proto.ressource.RefreshTokenRsp>}
 */
const methodDescriptor_RessourceService_RefreshToken = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RefreshToken',
  grpc.web.MethodType.UNARY,
  proto.ressource.RefreshTokenRqst,
  proto.ressource.RefreshTokenRsp,
  /**
   * @param {!proto.ressource.RefreshTokenRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RefreshTokenRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RefreshTokenRqst,
 *   !proto.ressource.RefreshTokenRsp>}
 */
const methodInfo_RessourceService_RefreshToken = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RefreshTokenRsp,
  /**
   * @param {!proto.ressource.RefreshTokenRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_RefreshToken,
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
      methodDescriptor_RessourceService_RefreshToken);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.AddAccountRoleRqst,
 *   !proto.ressource.AddAccountRoleRsp>}
 */
const methodDescriptor_RessourceService_AddAccountRole = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/AddAccountRole',
  grpc.web.MethodType.UNARY,
  proto.ressource.AddAccountRoleRqst,
  proto.ressource.AddAccountRoleRsp,
  /**
   * @param {!proto.ressource.AddAccountRoleRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AddAccountRoleRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.AddAccountRoleRqst,
 *   !proto.ressource.AddAccountRoleRsp>}
 */
const methodInfo_RessourceService_AddAccountRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.AddAccountRoleRsp,
  /**
   * @param {!proto.ressource.AddAccountRoleRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_AddAccountRole,
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
      methodDescriptor_RessourceService_AddAccountRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RemoveAccountRoleRqst,
 *   !proto.ressource.RemoveAccountRoleRsp>}
 */
const methodDescriptor_RessourceService_RemoveAccountRole = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RemoveAccountRole',
  grpc.web.MethodType.UNARY,
  proto.ressource.RemoveAccountRoleRqst,
  proto.ressource.RemoveAccountRoleRsp,
  /**
   * @param {!proto.ressource.RemoveAccountRoleRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveAccountRoleRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RemoveAccountRoleRqst,
 *   !proto.ressource.RemoveAccountRoleRsp>}
 */
const methodInfo_RessourceService_RemoveAccountRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemoveAccountRoleRsp,
  /**
   * @param {!proto.ressource.RemoveAccountRoleRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_RemoveAccountRole,
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
      methodDescriptor_RessourceService_RemoveAccountRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.CreateRoleRqst,
 *   !proto.ressource.CreateRoleRsp>}
 */
const methodDescriptor_RessourceService_CreateRole = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/CreateRole',
  grpc.web.MethodType.UNARY,
  proto.ressource.CreateRoleRqst,
  proto.ressource.CreateRoleRsp,
  /**
   * @param {!proto.ressource.CreateRoleRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.CreateRoleRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.CreateRoleRqst,
 *   !proto.ressource.CreateRoleRsp>}
 */
const methodInfo_RessourceService_CreateRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.CreateRoleRsp,
  /**
   * @param {!proto.ressource.CreateRoleRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_CreateRole,
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
      methodDescriptor_RessourceService_CreateRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteRoleRqst,
 *   !proto.ressource.DeleteRoleRsp>}
 */
const methodDescriptor_RessourceService_DeleteRole = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteRole',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteRoleRqst,
  proto.ressource.DeleteRoleRsp,
  /**
   * @param {!proto.ressource.DeleteRoleRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteRoleRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteRoleRqst,
 *   !proto.ressource.DeleteRoleRsp>}
 */
const methodInfo_RessourceService_DeleteRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteRoleRsp,
  /**
   * @param {!proto.ressource.DeleteRoleRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteRole,
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
      methodDescriptor_RessourceService_DeleteRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.AddRoleActionRqst,
 *   !proto.ressource.AddRoleActionRsp>}
 */
const methodDescriptor_RessourceService_AddRoleAction = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/AddRoleAction',
  grpc.web.MethodType.UNARY,
  proto.ressource.AddRoleActionRqst,
  proto.ressource.AddRoleActionRsp,
  /**
   * @param {!proto.ressource.AddRoleActionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AddRoleActionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.AddRoleActionRqst,
 *   !proto.ressource.AddRoleActionRsp>}
 */
const methodInfo_RessourceService_AddRoleAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.AddRoleActionRsp,
  /**
   * @param {!proto.ressource.AddRoleActionRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_AddRoleAction,
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
      methodDescriptor_RessourceService_AddRoleAction);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RemoveRoleActionRqst,
 *   !proto.ressource.RemoveRoleActionRsp>}
 */
const methodDescriptor_RessourceService_RemoveRoleAction = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RemoveRoleAction',
  grpc.web.MethodType.UNARY,
  proto.ressource.RemoveRoleActionRqst,
  proto.ressource.RemoveRoleActionRsp,
  /**
   * @param {!proto.ressource.RemoveRoleActionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveRoleActionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RemoveRoleActionRqst,
 *   !proto.ressource.RemoveRoleActionRsp>}
 */
const methodInfo_RessourceService_RemoveRoleAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemoveRoleActionRsp,
  /**
   * @param {!proto.ressource.RemoveRoleActionRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_RemoveRoleAction,
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
      methodDescriptor_RessourceService_RemoveRoleAction);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.AddApplicationActionRqst,
 *   !proto.ressource.AddApplicationActionRsp>}
 */
const methodDescriptor_RessourceService_AddApplicationAction = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/AddApplicationAction',
  grpc.web.MethodType.UNARY,
  proto.ressource.AddApplicationActionRqst,
  proto.ressource.AddApplicationActionRsp,
  /**
   * @param {!proto.ressource.AddApplicationActionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.AddApplicationActionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.AddApplicationActionRqst,
 *   !proto.ressource.AddApplicationActionRsp>}
 */
const methodInfo_RessourceService_AddApplicationAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.AddApplicationActionRsp,
  /**
   * @param {!proto.ressource.AddApplicationActionRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_AddApplicationAction,
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
      methodDescriptor_RessourceService_AddApplicationAction);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RemoveApplicationActionRqst,
 *   !proto.ressource.RemoveApplicationActionRsp>}
 */
const methodDescriptor_RessourceService_RemoveApplicationAction = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RemoveApplicationAction',
  grpc.web.MethodType.UNARY,
  proto.ressource.RemoveApplicationActionRqst,
  proto.ressource.RemoveApplicationActionRsp,
  /**
   * @param {!proto.ressource.RemoveApplicationActionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveApplicationActionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RemoveApplicationActionRqst,
 *   !proto.ressource.RemoveApplicationActionRsp>}
 */
const methodInfo_RessourceService_RemoveApplicationAction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemoveApplicationActionRsp,
  /**
   * @param {!proto.ressource.RemoveApplicationActionRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_RemoveApplicationAction,
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
      methodDescriptor_RessourceService_RemoveApplicationAction);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.GetAllActionsRqst,
 *   !proto.ressource.GetAllActionsRsp>}
 */
const methodDescriptor_RessourceService_GetAllActions = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/GetAllActions',
  grpc.web.MethodType.UNARY,
  proto.ressource.GetAllActionsRqst,
  proto.ressource.GetAllActionsRsp,
  /**
   * @param {!proto.ressource.GetAllActionsRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetAllActionsRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetAllActionsRqst,
 *   !proto.ressource.GetAllActionsRsp>}
 */
const methodInfo_RessourceService_GetAllActions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetAllActionsRsp,
  /**
   * @param {!proto.ressource.GetAllActionsRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_GetAllActions,
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
      methodDescriptor_RessourceService_GetAllActions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.GetPermissionsRqst,
 *   !proto.ressource.GetPermissionsRsp>}
 */
const methodDescriptor_RessourceService_GetPermissions = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/GetPermissions',
  grpc.web.MethodType.UNARY,
  proto.ressource.GetPermissionsRqst,
  proto.ressource.GetPermissionsRsp,
  /**
   * @param {!proto.ressource.GetPermissionsRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetPermissionsRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetPermissionsRqst,
 *   !proto.ressource.GetPermissionsRsp>}
 */
const methodInfo_RessourceService_GetPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetPermissionsRsp,
  /**
   * @param {!proto.ressource.GetPermissionsRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_GetPermissions,
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
      methodDescriptor_RessourceService_GetPermissions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.SetPermissionRqst,
 *   !proto.ressource.SetPermissionRsp>}
 */
const methodDescriptor_RessourceService_SetPermission = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/SetPermission',
  grpc.web.MethodType.UNARY,
  proto.ressource.SetPermissionRqst,
  proto.ressource.SetPermissionRsp,
  /**
   * @param {!proto.ressource.SetPermissionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SetPermissionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.SetPermissionRqst,
 *   !proto.ressource.SetPermissionRsp>}
 */
const methodInfo_RessourceService_SetPermission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.SetPermissionRsp,
  /**
   * @param {!proto.ressource.SetPermissionRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_SetPermission,
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
      methodDescriptor_RessourceService_SetPermission);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeletePermissionsRqst,
 *   !proto.ressource.DeletePermissionsRsp>}
 */
const methodDescriptor_RessourceService_DeletePermissions = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeletePermissions',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeletePermissionsRqst,
  proto.ressource.DeletePermissionsRsp,
  /**
   * @param {!proto.ressource.DeletePermissionsRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeletePermissionsRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeletePermissionsRqst,
 *   !proto.ressource.DeletePermissionsRsp>}
 */
const methodInfo_RessourceService_DeletePermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeletePermissionsRsp,
  /**
   * @param {!proto.ressource.DeletePermissionsRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeletePermissions,
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
      methodDescriptor_RessourceService_DeletePermissions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.SetRessourceOwnerRqst,
 *   !proto.ressource.SetRessourceOwnerRsp>}
 */
const methodDescriptor_RessourceService_SetRessourceOwner = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/SetRessourceOwner',
  grpc.web.MethodType.UNARY,
  proto.ressource.SetRessourceOwnerRqst,
  proto.ressource.SetRessourceOwnerRsp,
  /**
   * @param {!proto.ressource.SetRessourceOwnerRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SetRessourceOwnerRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.SetRessourceOwnerRqst,
 *   !proto.ressource.SetRessourceOwnerRsp>}
 */
const methodInfo_RessourceService_SetRessourceOwner = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.SetRessourceOwnerRsp,
  /**
   * @param {!proto.ressource.SetRessourceOwnerRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_SetRessourceOwner,
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
      methodDescriptor_RessourceService_SetRessourceOwner);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.GetRessourceOwnersRqst,
 *   !proto.ressource.GetRessourceOwnersRsp>}
 */
const methodDescriptor_RessourceService_GetRessourceOwners = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/GetRessourceOwners',
  grpc.web.MethodType.UNARY,
  proto.ressource.GetRessourceOwnersRqst,
  proto.ressource.GetRessourceOwnersRsp,
  /**
   * @param {!proto.ressource.GetRessourceOwnersRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetRessourceOwnersRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetRessourceOwnersRqst,
 *   !proto.ressource.GetRessourceOwnersRsp>}
 */
const methodInfo_RessourceService_GetRessourceOwners = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetRessourceOwnersRsp,
  /**
   * @param {!proto.ressource.GetRessourceOwnersRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_GetRessourceOwners,
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
      methodDescriptor_RessourceService_GetRessourceOwners);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteRessourceOwnerRqst,
 *   !proto.ressource.DeleteRessourceOwnerRsp>}
 */
const methodDescriptor_RessourceService_DeleteRessourceOwner = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteRessourceOwner',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteRessourceOwnerRqst,
  proto.ressource.DeleteRessourceOwnerRsp,
  /**
   * @param {!proto.ressource.DeleteRessourceOwnerRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteRessourceOwnerRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteRessourceOwnerRqst,
 *   !proto.ressource.DeleteRessourceOwnerRsp>}
 */
const methodInfo_RessourceService_DeleteRessourceOwner = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteRessourceOwnerRsp,
  /**
   * @param {!proto.ressource.DeleteRessourceOwnerRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteRessourceOwner,
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
      methodDescriptor_RessourceService_DeleteRessourceOwner);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteRessourceOwnersRqst,
 *   !proto.ressource.DeleteRessourceOwnersRsp>}
 */
const methodDescriptor_RessourceService_DeleteRessourceOwners = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteRessourceOwners',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteRessourceOwnersRqst,
  proto.ressource.DeleteRessourceOwnersRsp,
  /**
   * @param {!proto.ressource.DeleteRessourceOwnersRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteRessourceOwnersRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteRessourceOwnersRqst,
 *   !proto.ressource.DeleteRessourceOwnersRsp>}
 */
const methodInfo_RessourceService_DeleteRessourceOwners = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteRessourceOwnersRsp,
  /**
   * @param {!proto.ressource.DeleteRessourceOwnersRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteRessourceOwners,
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
      methodDescriptor_RessourceService_DeleteRessourceOwners);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.GetAllFilesInfoRqst,
 *   !proto.ressource.GetAllFilesInfoRsp>}
 */
const methodDescriptor_RessourceService_GetAllFilesInfo = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/GetAllFilesInfo',
  grpc.web.MethodType.UNARY,
  proto.ressource.GetAllFilesInfoRqst,
  proto.ressource.GetAllFilesInfoRsp,
  /**
   * @param {!proto.ressource.GetAllFilesInfoRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetAllFilesInfoRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetAllFilesInfoRqst,
 *   !proto.ressource.GetAllFilesInfoRsp>}
 */
const methodInfo_RessourceService_GetAllFilesInfo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetAllFilesInfoRsp,
  /**
   * @param {!proto.ressource.GetAllFilesInfoRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_GetAllFilesInfo,
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
      methodDescriptor_RessourceService_GetAllFilesInfo);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.ValidateTokenRqst,
 *   !proto.ressource.ValidateTokenRsp>}
 */
const methodDescriptor_RessourceService_ValidateToken = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/ValidateToken',
  grpc.web.MethodType.UNARY,
  proto.ressource.ValidateTokenRqst,
  proto.ressource.ValidateTokenRsp,
  /**
   * @param {!proto.ressource.ValidateTokenRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateTokenRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidateTokenRqst,
 *   !proto.ressource.ValidateTokenRsp>}
 */
const methodInfo_RessourceService_ValidateToken = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidateTokenRsp,
  /**
   * @param {!proto.ressource.ValidateTokenRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateTokenRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ValidateTokenRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ValidateTokenRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ValidateTokenRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.validateToken =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ValidateToken',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidateToken,
      callback);
};


/**
 * @param {!proto.ressource.ValidateTokenRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ValidateTokenRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.validateToken =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ValidateToken',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidateToken);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.ValidateUserRessourceAccessRqst,
 *   !proto.ressource.ValidateUserRessourceAccessRsp>}
 */
const methodDescriptor_RessourceService_ValidateUserRessourceAccess = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/ValidateUserRessourceAccess',
  grpc.web.MethodType.UNARY,
  proto.ressource.ValidateUserRessourceAccessRqst,
  proto.ressource.ValidateUserRessourceAccessRsp,
  /**
   * @param {!proto.ressource.ValidateUserRessourceAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateUserRessourceAccessRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidateUserRessourceAccessRqst,
 *   !proto.ressource.ValidateUserRessourceAccessRsp>}
 */
const methodInfo_RessourceService_ValidateUserRessourceAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidateUserRessourceAccessRsp,
  /**
   * @param {!proto.ressource.ValidateUserRessourceAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateUserRessourceAccessRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ValidateUserRessourceAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ValidateUserRessourceAccessRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ValidateUserRessourceAccessRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.validateUserRessourceAccess =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ValidateUserRessourceAccess',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidateUserRessourceAccess,
      callback);
};


/**
 * @param {!proto.ressource.ValidateUserRessourceAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ValidateUserRessourceAccessRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.validateUserRessourceAccess =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ValidateUserRessourceAccess',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidateUserRessourceAccess);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.ValidateApplicationRessourceAccessRqst,
 *   !proto.ressource.ValidateApplicationRessourceAccessRsp>}
 */
const methodDescriptor_RessourceService_ValidateApplicationRessourceAccess = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/ValidateApplicationRessourceAccess',
  grpc.web.MethodType.UNARY,
  proto.ressource.ValidateApplicationRessourceAccessRqst,
  proto.ressource.ValidateApplicationRessourceAccessRsp,
  /**
   * @param {!proto.ressource.ValidateApplicationRessourceAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateApplicationRessourceAccessRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidateApplicationRessourceAccessRqst,
 *   !proto.ressource.ValidateApplicationRessourceAccessRsp>}
 */
const methodInfo_RessourceService_ValidateApplicationRessourceAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidateApplicationRessourceAccessRsp,
  /**
   * @param {!proto.ressource.ValidateApplicationRessourceAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateApplicationRessourceAccessRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ValidateApplicationRessourceAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ValidateApplicationRessourceAccessRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ValidateApplicationRessourceAccessRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.validateApplicationRessourceAccess =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ValidateApplicationRessourceAccess',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidateApplicationRessourceAccess,
      callback);
};


/**
 * @param {!proto.ressource.ValidateApplicationRessourceAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ValidateApplicationRessourceAccessRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.validateApplicationRessourceAccess =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ValidateApplicationRessourceAccess',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidateApplicationRessourceAccess);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.ValidatePeerRessourceAccessRqst,
 *   !proto.ressource.ValidatePeerRessourceAccessRsp>}
 */
const methodDescriptor_RessourceService_ValidatePeerRessourceAccess = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/ValidatePeerRessourceAccess',
  grpc.web.MethodType.UNARY,
  proto.ressource.ValidatePeerRessourceAccessRqst,
  proto.ressource.ValidatePeerRessourceAccessRsp,
  /**
   * @param {!proto.ressource.ValidatePeerRessourceAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidatePeerRessourceAccessRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidatePeerRessourceAccessRqst,
 *   !proto.ressource.ValidatePeerRessourceAccessRsp>}
 */
const methodInfo_RessourceService_ValidatePeerRessourceAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidatePeerRessourceAccessRsp,
  /**
   * @param {!proto.ressource.ValidatePeerRessourceAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidatePeerRessourceAccessRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ValidatePeerRessourceAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ValidatePeerRessourceAccessRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ValidatePeerRessourceAccessRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.validatePeerRessourceAccess =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ValidatePeerRessourceAccess',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidatePeerRessourceAccess,
      callback);
};


/**
 * @param {!proto.ressource.ValidatePeerRessourceAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ValidatePeerRessourceAccessRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.validatePeerRessourceAccess =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ValidatePeerRessourceAccess',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidatePeerRessourceAccess);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.ValidateUserAccessRqst,
 *   !proto.ressource.ValidateUserAccessRsp>}
 */
const methodDescriptor_RessourceService_ValidateUserAccess = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/ValidateUserAccess',
  grpc.web.MethodType.UNARY,
  proto.ressource.ValidateUserAccessRqst,
  proto.ressource.ValidateUserAccessRsp,
  /**
   * @param {!proto.ressource.ValidateUserAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateUserAccessRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidateUserAccessRqst,
 *   !proto.ressource.ValidateUserAccessRsp>}
 */
const methodInfo_RessourceService_ValidateUserAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidateUserAccessRsp,
  /**
   * @param {!proto.ressource.ValidateUserAccessRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_ValidateUserAccess,
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
      methodDescriptor_RessourceService_ValidateUserAccess);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.ValidateApplicationAccessRqst,
 *   !proto.ressource.ValidateApplicationAccessRsp>}
 */
const methodDescriptor_RessourceService_ValidateApplicationAccess = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/ValidateApplicationAccess',
  grpc.web.MethodType.UNARY,
  proto.ressource.ValidateApplicationAccessRqst,
  proto.ressource.ValidateApplicationAccessRsp,
  /**
   * @param {!proto.ressource.ValidateApplicationAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidateApplicationAccessRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidateApplicationAccessRqst,
 *   !proto.ressource.ValidateApplicationAccessRsp>}
 */
const methodInfo_RessourceService_ValidateApplicationAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidateApplicationAccessRsp,
  /**
   * @param {!proto.ressource.ValidateApplicationAccessRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_ValidateApplicationAccess,
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
      methodDescriptor_RessourceService_ValidateApplicationAccess);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.ValidatePeerAccessRqst,
 *   !proto.ressource.ValidatePeerAccessRsp>}
 */
const methodDescriptor_RessourceService_ValidatePeerAccess = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/ValidatePeerAccess',
  grpc.web.MethodType.UNARY,
  proto.ressource.ValidatePeerAccessRqst,
  proto.ressource.ValidatePeerAccessRsp,
  /**
   * @param {!proto.ressource.ValidatePeerAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidatePeerAccessRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ValidatePeerAccessRqst,
 *   !proto.ressource.ValidatePeerAccessRsp>}
 */
const methodInfo_RessourceService_ValidatePeerAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ValidatePeerAccessRsp,
  /**
   * @param {!proto.ressource.ValidatePeerAccessRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ValidatePeerAccessRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.ValidatePeerAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.ValidatePeerAccessRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.ValidatePeerAccessRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.validatePeerAccess =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/ValidatePeerAccess',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidatePeerAccess,
      callback);
};


/**
 * @param {!proto.ressource.ValidatePeerAccessRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.ValidatePeerAccessRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.validatePeerAccess =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/ValidatePeerAccess',
      request,
      metadata || {},
      methodDescriptor_RessourceService_ValidatePeerAccess);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.CreateDirPermissionsRqst,
 *   !proto.ressource.CreateDirPermissionsRsp>}
 */
const methodDescriptor_RessourceService_CreateDirPermissions = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/CreateDirPermissions',
  grpc.web.MethodType.UNARY,
  proto.ressource.CreateDirPermissionsRqst,
  proto.ressource.CreateDirPermissionsRsp,
  /**
   * @param {!proto.ressource.CreateDirPermissionsRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.CreateDirPermissionsRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.CreateDirPermissionsRqst,
 *   !proto.ressource.CreateDirPermissionsRsp>}
 */
const methodInfo_RessourceService_CreateDirPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.CreateDirPermissionsRsp,
  /**
   * @param {!proto.ressource.CreateDirPermissionsRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_CreateDirPermissions,
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
      methodDescriptor_RessourceService_CreateDirPermissions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RenameFilePermissionRqst,
 *   !proto.ressource.RenameFilePermissionRsp>}
 */
const methodDescriptor_RessourceService_RenameFilePermission = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RenameFilePermission',
  grpc.web.MethodType.UNARY,
  proto.ressource.RenameFilePermissionRqst,
  proto.ressource.RenameFilePermissionRsp,
  /**
   * @param {!proto.ressource.RenameFilePermissionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RenameFilePermissionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RenameFilePermissionRqst,
 *   !proto.ressource.RenameFilePermissionRsp>}
 */
const methodInfo_RessourceService_RenameFilePermission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RenameFilePermissionRsp,
  /**
   * @param {!proto.ressource.RenameFilePermissionRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_RenameFilePermission,
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
      methodDescriptor_RessourceService_RenameFilePermission);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteDirPermissionsRqst,
 *   !proto.ressource.DeleteDirPermissionsRsp>}
 */
const methodDescriptor_RessourceService_DeleteDirPermissions = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteDirPermissions',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteDirPermissionsRqst,
  proto.ressource.DeleteDirPermissionsRsp,
  /**
   * @param {!proto.ressource.DeleteDirPermissionsRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteDirPermissionsRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteDirPermissionsRqst,
 *   !proto.ressource.DeleteDirPermissionsRsp>}
 */
const methodInfo_RessourceService_DeleteDirPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteDirPermissionsRsp,
  /**
   * @param {!proto.ressource.DeleteDirPermissionsRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteDirPermissions,
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
      methodDescriptor_RessourceService_DeleteDirPermissions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteFilePermissionsRqst,
 *   !proto.ressource.DeleteFilePermissionsRsp>}
 */
const methodDescriptor_RessourceService_DeleteFilePermissions = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteFilePermissions',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteFilePermissionsRqst,
  proto.ressource.DeleteFilePermissionsRsp,
  /**
   * @param {!proto.ressource.DeleteFilePermissionsRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteFilePermissionsRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteFilePermissionsRqst,
 *   !proto.ressource.DeleteFilePermissionsRsp>}
 */
const methodInfo_RessourceService_DeleteFilePermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteFilePermissionsRsp,
  /**
   * @param {!proto.ressource.DeleteFilePermissionsRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteFilePermissions,
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
      methodDescriptor_RessourceService_DeleteFilePermissions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteAccountPermissionsRqst,
 *   !proto.ressource.DeleteAccountPermissionsRsp>}
 */
const methodDescriptor_RessourceService_DeleteAccountPermissions = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteAccountPermissions',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteAccountPermissionsRqst,
  proto.ressource.DeleteAccountPermissionsRsp,
  /**
   * @param {!proto.ressource.DeleteAccountPermissionsRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteAccountPermissionsRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteAccountPermissionsRqst,
 *   !proto.ressource.DeleteAccountPermissionsRsp>}
 */
const methodInfo_RessourceService_DeleteAccountPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteAccountPermissionsRsp,
  /**
   * @param {!proto.ressource.DeleteAccountPermissionsRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteAccountPermissions,
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
      methodDescriptor_RessourceService_DeleteAccountPermissions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteRolePermissionsRqst,
 *   !proto.ressource.DeleteRolePermissionsRsp>}
 */
const methodDescriptor_RessourceService_DeleteRolePermissions = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteRolePermissions',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteRolePermissionsRqst,
  proto.ressource.DeleteRolePermissionsRsp,
  /**
   * @param {!proto.ressource.DeleteRolePermissionsRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteRolePermissionsRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteRolePermissionsRqst,
 *   !proto.ressource.DeleteRolePermissionsRsp>}
 */
const methodInfo_RessourceService_DeleteRolePermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteRolePermissionsRsp,
  /**
   * @param {!proto.ressource.DeleteRolePermissionsRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteRolePermissions,
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
      methodDescriptor_RessourceService_DeleteRolePermissions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.GetAllApplicationsInfoRqst,
 *   !proto.ressource.GetAllApplicationsInfoRsp>}
 */
const methodDescriptor_RessourceService_GetAllApplicationsInfo = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/GetAllApplicationsInfo',
  grpc.web.MethodType.UNARY,
  proto.ressource.GetAllApplicationsInfoRqst,
  proto.ressource.GetAllApplicationsInfoRsp,
  /**
   * @param {!proto.ressource.GetAllApplicationsInfoRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetAllApplicationsInfoRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetAllApplicationsInfoRqst,
 *   !proto.ressource.GetAllApplicationsInfoRsp>}
 */
const methodInfo_RessourceService_GetAllApplicationsInfo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetAllApplicationsInfoRsp,
  /**
   * @param {!proto.ressource.GetAllApplicationsInfoRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_GetAllApplicationsInfo,
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
      methodDescriptor_RessourceService_GetAllApplicationsInfo);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteApplicationRqst,
 *   !proto.ressource.DeleteApplicationRsp>}
 */
const methodDescriptor_RessourceService_DeleteApplication = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteApplication',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteApplicationRqst,
  proto.ressource.DeleteApplicationRsp,
  /**
   * @param {!proto.ressource.DeleteApplicationRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteApplicationRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteApplicationRqst,
 *   !proto.ressource.DeleteApplicationRsp>}
 */
const methodInfo_RessourceService_DeleteApplication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteApplicationRsp,
  /**
   * @param {!proto.ressource.DeleteApplicationRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteApplication,
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
      methodDescriptor_RessourceService_DeleteApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.LogRqst,
 *   !proto.ressource.LogRsp>}
 */
const methodDescriptor_RessourceService_Log = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/Log',
  grpc.web.MethodType.UNARY,
  proto.ressource.LogRqst,
  proto.ressource.LogRsp,
  /**
   * @param {!proto.ressource.LogRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.LogRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.LogRqst,
 *   !proto.ressource.LogRsp>}
 */
const methodInfo_RessourceService_Log = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.LogRsp,
  /**
   * @param {!proto.ressource.LogRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_Log,
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
      methodDescriptor_RessourceService_Log);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.GetLogRqst,
 *   !proto.ressource.GetLogRsp>}
 */
const methodDescriptor_RessourceService_GetLog = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/GetLog',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.ressource.GetLogRqst,
  proto.ressource.GetLogRsp,
  /**
   * @param {!proto.ressource.GetLogRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetLogRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetLogRqst,
 *   !proto.ressource.GetLogRsp>}
 */
const methodInfo_RessourceService_GetLog = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetLogRsp,
  /**
   * @param {!proto.ressource.GetLogRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_GetLog);
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
      methodDescriptor_RessourceService_GetLog);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.DeleteLogRqst,
 *   !proto.ressource.DeleteLogRsp>}
 */
const methodDescriptor_RessourceService_DeleteLog = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/DeleteLog',
  grpc.web.MethodType.UNARY,
  proto.ressource.DeleteLogRqst,
  proto.ressource.DeleteLogRsp,
  /**
   * @param {!proto.ressource.DeleteLogRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.DeleteLogRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.DeleteLogRqst,
 *   !proto.ressource.DeleteLogRsp>}
 */
const methodInfo_RessourceService_DeleteLog = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.DeleteLogRsp,
  /**
   * @param {!proto.ressource.DeleteLogRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_DeleteLog,
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
      methodDescriptor_RessourceService_DeleteLog);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.ClearAllLogRqst,
 *   !proto.ressource.ClearAllLogRsp>}
 */
const methodDescriptor_RessourceService_ClearAllLog = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/ClearAllLog',
  grpc.web.MethodType.UNARY,
  proto.ressource.ClearAllLogRqst,
  proto.ressource.ClearAllLogRsp,
  /**
   * @param {!proto.ressource.ClearAllLogRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.ClearAllLogRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.ClearAllLogRqst,
 *   !proto.ressource.ClearAllLogRsp>}
 */
const methodInfo_RessourceService_ClearAllLog = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.ClearAllLogRsp,
  /**
   * @param {!proto.ressource.ClearAllLogRqst} request
   * @return {!Uint8Array}
   */
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
      methodDescriptor_RessourceService_ClearAllLog,
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
      methodDescriptor_RessourceService_ClearAllLog);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.GetRessourcesRqst,
 *   !proto.ressource.GetRessourcesRsp>}
 */
const methodDescriptor_RessourceService_GetRessources = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/GetRessources',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.ressource.GetRessourcesRqst,
  proto.ressource.GetRessourcesRsp,
  /**
   * @param {!proto.ressource.GetRessourcesRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetRessourcesRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetRessourcesRqst,
 *   !proto.ressource.GetRessourcesRsp>}
 */
const methodInfo_RessourceService_GetRessources = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetRessourcesRsp,
  /**
   * @param {!proto.ressource.GetRessourcesRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetRessourcesRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.GetRessourcesRqst} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetRessourcesRsp>}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.getRessources =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/ressource.RessourceService/GetRessources',
      request,
      metadata || {},
      methodDescriptor_RessourceService_GetRessources);
};


/**
 * @param {!proto.ressource.GetRessourcesRqst} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetRessourcesRsp>}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServicePromiseClient.prototype.getRessources =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/ressource.RessourceService/GetRessources',
      request,
      metadata || {},
      methodDescriptor_RessourceService_GetRessources);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.SetRessourceRqst,
 *   !proto.ressource.SetRessourceRsp>}
 */
const methodDescriptor_RessourceService_SetRessource = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/SetRessource',
  grpc.web.MethodType.UNARY,
  proto.ressource.SetRessourceRqst,
  proto.ressource.SetRessourceRsp,
  /**
   * @param {!proto.ressource.SetRessourceRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SetRessourceRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.SetRessourceRqst,
 *   !proto.ressource.SetRessourceRsp>}
 */
const methodInfo_RessourceService_SetRessource = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.SetRessourceRsp,
  /**
   * @param {!proto.ressource.SetRessourceRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SetRessourceRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.SetRessourceRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.SetRessourceRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.SetRessourceRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.setRessource =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/SetRessource',
      request,
      metadata || {},
      methodDescriptor_RessourceService_SetRessource,
      callback);
};


/**
 * @param {!proto.ressource.SetRessourceRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.SetRessourceRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.setRessource =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/SetRessource',
      request,
      metadata || {},
      methodDescriptor_RessourceService_SetRessource);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RemoveRessourceRqst,
 *   !proto.ressource.RemoveRessourceRsp>}
 */
const methodDescriptor_RessourceService_RemoveRessource = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RemoveRessource',
  grpc.web.MethodType.UNARY,
  proto.ressource.RemoveRessourceRqst,
  proto.ressource.RemoveRessourceRsp,
  /**
   * @param {!proto.ressource.RemoveRessourceRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveRessourceRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RemoveRessourceRqst,
 *   !proto.ressource.RemoveRessourceRsp>}
 */
const methodInfo_RessourceService_RemoveRessource = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemoveRessourceRsp,
  /**
   * @param {!proto.ressource.RemoveRessourceRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveRessourceRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RemoveRessourceRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RemoveRessourceRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RemoveRessourceRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.removeRessource =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RemoveRessource',
      request,
      metadata || {},
      methodDescriptor_RessourceService_RemoveRessource,
      callback);
};


/**
 * @param {!proto.ressource.RemoveRessourceRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RemoveRessourceRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.removeRessource =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RemoveRessource',
      request,
      metadata || {},
      methodDescriptor_RessourceService_RemoveRessource);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.SetActionPermissionRqst,
 *   !proto.ressource.SetActionPermissionRsp>}
 */
const methodDescriptor_RessourceService_SetActionPermission = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/SetActionPermission',
  grpc.web.MethodType.UNARY,
  proto.ressource.SetActionPermissionRqst,
  proto.ressource.SetActionPermissionRsp,
  /**
   * @param {!proto.ressource.SetActionPermissionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SetActionPermissionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.SetActionPermissionRqst,
 *   !proto.ressource.SetActionPermissionRsp>}
 */
const methodInfo_RessourceService_SetActionPermission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.SetActionPermissionRsp,
  /**
   * @param {!proto.ressource.SetActionPermissionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.SetActionPermissionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.SetActionPermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.SetActionPermissionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.SetActionPermissionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.setActionPermission =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/SetActionPermission',
      request,
      metadata || {},
      methodDescriptor_RessourceService_SetActionPermission,
      callback);
};


/**
 * @param {!proto.ressource.SetActionPermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.SetActionPermissionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.setActionPermission =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/SetActionPermission',
      request,
      metadata || {},
      methodDescriptor_RessourceService_SetActionPermission);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.RemoveActionPermissionRqst,
 *   !proto.ressource.RemoveActionPermissionRsp>}
 */
const methodDescriptor_RessourceService_RemoveActionPermission = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/RemoveActionPermission',
  grpc.web.MethodType.UNARY,
  proto.ressource.RemoveActionPermissionRqst,
  proto.ressource.RemoveActionPermissionRsp,
  /**
   * @param {!proto.ressource.RemoveActionPermissionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveActionPermissionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.RemoveActionPermissionRqst,
 *   !proto.ressource.RemoveActionPermissionRsp>}
 */
const methodInfo_RessourceService_RemoveActionPermission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.RemoveActionPermissionRsp,
  /**
   * @param {!proto.ressource.RemoveActionPermissionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.RemoveActionPermissionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.RemoveActionPermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.RemoveActionPermissionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.RemoveActionPermissionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.removeActionPermission =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/RemoveActionPermission',
      request,
      metadata || {},
      methodDescriptor_RessourceService_RemoveActionPermission,
      callback);
};


/**
 * @param {!proto.ressource.RemoveActionPermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.RemoveActionPermissionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.removeActionPermission =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/RemoveActionPermission',
      request,
      metadata || {},
      methodDescriptor_RessourceService_RemoveActionPermission);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ressource.GetActionPermissionRqst,
 *   !proto.ressource.GetActionPermissionRsp>}
 */
const methodDescriptor_RessourceService_GetActionPermission = new grpc.web.MethodDescriptor(
  '/ressource.RessourceService/GetActionPermission',
  grpc.web.MethodType.UNARY,
  proto.ressource.GetActionPermissionRqst,
  proto.ressource.GetActionPermissionRsp,
  /**
   * @param {!proto.ressource.GetActionPermissionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetActionPermissionRsp.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ressource.GetActionPermissionRqst,
 *   !proto.ressource.GetActionPermissionRsp>}
 */
const methodInfo_RessourceService_GetActionPermission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ressource.GetActionPermissionRsp,
  /**
   * @param {!proto.ressource.GetActionPermissionRqst} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ressource.GetActionPermissionRsp.deserializeBinary
);


/**
 * @param {!proto.ressource.GetActionPermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ressource.GetActionPermissionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ressource.GetActionPermissionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ressource.RessourceServiceClient.prototype.getActionPermission =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ressource.RessourceService/GetActionPermission',
      request,
      metadata || {},
      methodDescriptor_RessourceService_GetActionPermission,
      callback);
};


/**
 * @param {!proto.ressource.GetActionPermissionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ressource.GetActionPermissionRsp>}
 *     A native promise that resolves to the response
 */
proto.ressource.RessourceServicePromiseClient.prototype.getActionPermission =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ressource.RessourceService/GetActionPermission',
      request,
      metadata || {},
      methodDescriptor_RessourceService_GetActionPermission);
};


module.exports = proto.ressource;

