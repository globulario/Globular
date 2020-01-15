/**
 * @fileoverview gRPC-Web generated client stub for admin
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.admin = require('./admin_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.admin.AdminServiceClient =
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
proto.admin.AdminServicePromiseClient =
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
 *   !proto.admin.SetRootPasswordRequest,
 *   !proto.admin.SetRootPasswordResponse>}
 */
const methodInfo_AdminService_SetRootPassword = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.SetRootPasswordResponse,
  /** @param {!proto.admin.SetRootPasswordRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.SetRootPasswordResponse.deserializeBinary
);


/**
 * @param {!proto.admin.SetRootPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.SetRootPasswordResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.SetRootPasswordResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.setRootPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/SetRootPassword',
      request,
      metadata || {},
      methodInfo_AdminService_SetRootPassword,
      callback);
};


/**
 * @param {!proto.admin.SetRootPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.SetRootPasswordResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.setRootPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/SetRootPassword',
      request,
      metadata || {},
      methodInfo_AdminService_SetRootPassword);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.SetRootEmailRequest,
 *   !proto.admin.SetRootEmailResponse>}
 */
const methodInfo_AdminService_SetRootEmail = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.SetRootEmailResponse,
  /** @param {!proto.admin.SetRootEmailRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.SetRootEmailResponse.deserializeBinary
);


/**
 * @param {!proto.admin.SetRootEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.SetRootEmailResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.SetRootEmailResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.setRootEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/SetRootEmail',
      request,
      metadata || {},
      methodInfo_AdminService_SetRootEmail,
      callback);
};


/**
 * @param {!proto.admin.SetRootEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.SetRootEmailResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.setRootEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/SetRootEmail',
      request,
      metadata || {},
      methodInfo_AdminService_SetRootEmail);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.GetConfigRequest,
 *   !proto.admin.GetConfigResponse>}
 */
const methodInfo_AdminService_GetConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.GetConfigResponse,
  /** @param {!proto.admin.GetConfigRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.GetConfigResponse.deserializeBinary
);


/**
 * @param {!proto.admin.GetConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.GetConfigResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.GetConfigResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.getConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/GetConfig',
      request,
      metadata || {},
      methodInfo_AdminService_GetConfig,
      callback);
};


/**
 * @param {!proto.admin.GetConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.GetConfigResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.getConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/GetConfig',
      request,
      metadata || {},
      methodInfo_AdminService_GetConfig);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.GetConfigRequest,
 *   !proto.admin.GetConfigResponse>}
 */
const methodInfo_AdminService_GetFullConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.GetConfigResponse,
  /** @param {!proto.admin.GetConfigRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.GetConfigResponse.deserializeBinary
);


/**
 * @param {!proto.admin.GetConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.GetConfigResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.GetConfigResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.getFullConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/GetFullConfig',
      request,
      metadata || {},
      methodInfo_AdminService_GetFullConfig,
      callback);
};


/**
 * @param {!proto.admin.GetConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.GetConfigResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.getFullConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/GetFullConfig',
      request,
      metadata || {},
      methodInfo_AdminService_GetFullConfig);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.SaveConfigRequest,
 *   !proto.admin.SaveConfigResponse>}
 */
const methodInfo_AdminService_SaveConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.SaveConfigResponse,
  /** @param {!proto.admin.SaveConfigRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.SaveConfigResponse.deserializeBinary
);


/**
 * @param {!proto.admin.SaveConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.SaveConfigResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.SaveConfigResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.saveConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/SaveConfig',
      request,
      metadata || {},
      methodInfo_AdminService_SaveConfig,
      callback);
};


/**
 * @param {!proto.admin.SaveConfigRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.SaveConfigResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.saveConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/SaveConfig',
      request,
      metadata || {},
      methodInfo_AdminService_SaveConfig);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.StopServiceRequest,
 *   !proto.admin.StopServiceResponse>}
 */
const methodInfo_AdminService_StopService = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.StopServiceResponse,
  /** @param {!proto.admin.StopServiceRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.StopServiceResponse.deserializeBinary
);


/**
 * @param {!proto.admin.StopServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.StopServiceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.StopServiceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.stopService =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/StopService',
      request,
      metadata || {},
      methodInfo_AdminService_StopService,
      callback);
};


/**
 * @param {!proto.admin.StopServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.StopServiceResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.stopService =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/StopService',
      request,
      metadata || {},
      methodInfo_AdminService_StopService);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.StartServiceRequest,
 *   !proto.admin.StartServiceResponse>}
 */
const methodInfo_AdminService_StartService = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.StartServiceResponse,
  /** @param {!proto.admin.StartServiceRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.StartServiceResponse.deserializeBinary
);


/**
 * @param {!proto.admin.StartServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.StartServiceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.StartServiceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.startService =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/StartService',
      request,
      metadata || {},
      methodInfo_AdminService_StartService,
      callback);
};


/**
 * @param {!proto.admin.StartServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.StartServiceResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.startService =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/StartService',
      request,
      metadata || {},
      methodInfo_AdminService_StartService);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.PublishServiceRequest,
 *   !proto.admin.PublishServiceResponse>}
 */
const methodInfo_AdminService_PublishService = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.PublishServiceResponse,
  /** @param {!proto.admin.PublishServiceRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.PublishServiceResponse.deserializeBinary
);


/**
 * @param {!proto.admin.PublishServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.PublishServiceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.PublishServiceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.publishService =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/PublishService',
      request,
      metadata || {},
      methodInfo_AdminService_PublishService,
      callback);
};


/**
 * @param {!proto.admin.PublishServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.PublishServiceResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.publishService =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/PublishService',
      request,
      metadata || {},
      methodInfo_AdminService_PublishService);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.InstallServiceRequest,
 *   !proto.admin.InstallServiceResponse>}
 */
const methodInfo_AdminService_InstallService = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.InstallServiceResponse,
  /** @param {!proto.admin.InstallServiceRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.InstallServiceResponse.deserializeBinary
);


/**
 * @param {!proto.admin.InstallServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.InstallServiceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.InstallServiceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.installService =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/InstallService',
      request,
      metadata || {},
      methodInfo_AdminService_InstallService,
      callback);
};


/**
 * @param {!proto.admin.InstallServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.InstallServiceResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.installService =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/InstallService',
      request,
      metadata || {},
      methodInfo_AdminService_InstallService);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.UninstallServiceRequest,
 *   !proto.admin.UninstallServiceResponse>}
 */
const methodInfo_AdminService_UninstallService = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.UninstallServiceResponse,
  /** @param {!proto.admin.UninstallServiceRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.UninstallServiceResponse.deserializeBinary
);


/**
 * @param {!proto.admin.UninstallServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.UninstallServiceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.UninstallServiceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.uninstallService =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/UninstallService',
      request,
      metadata || {},
      methodInfo_AdminService_UninstallService,
      callback);
};


/**
 * @param {!proto.admin.UninstallServiceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.UninstallServiceResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.uninstallService =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/UninstallService',
      request,
      metadata || {},
      methodInfo_AdminService_UninstallService);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.admin.RegisterExternalApplicationRequest,
 *   !proto.admin.RegisterExternalApplicationResponse>}
 */
const methodInfo_AdminService_RegisterExternalApplication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.admin.RegisterExternalApplicationResponse,
  /** @param {!proto.admin.RegisterExternalApplicationRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.admin.RegisterExternalApplicationResponse.deserializeBinary
);


/**
 * @param {!proto.admin.RegisterExternalApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.admin.RegisterExternalApplicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.admin.RegisterExternalApplicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.admin.AdminServiceClient.prototype.registerExternalApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/admin.AdminService/RegisterExternalApplication',
      request,
      metadata || {},
      methodInfo_AdminService_RegisterExternalApplication,
      callback);
};


/**
 * @param {!proto.admin.RegisterExternalApplicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.admin.RegisterExternalApplicationResponse>}
 *     A native promise that resolves to the response
 */
proto.admin.AdminServicePromiseClient.prototype.registerExternalApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/admin.AdminService/RegisterExternalApplication',
      request,
      metadata || {},
      methodInfo_AdminService_RegisterExternalApplication);
};


module.exports = proto.admin;

