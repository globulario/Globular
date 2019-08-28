/**
 * @fileoverview gRPC-Web generated client stub for ldap
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.ldap = require('./ldap_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.ldap.LdapServiceClient =
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
proto.ldap.LdapServicePromiseClient =
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
 *   !proto.ldap.CreateConnectionRqst,
 *   !proto.ldap.CreateConnectionRsp>}
 */
const methodInfo_LdapService_CreateConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ldap.CreateConnectionRsp,
  /** @param {!proto.ldap.CreateConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ldap.CreateConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.ldap.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ldap.CreateConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ldap.CreateConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ldap.LdapServiceClient.prototype.createConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ldap.LdapService/CreateConnection',
      request,
      metadata || {},
      methodInfo_LdapService_CreateConnection,
      callback);
};


/**
 * @param {!proto.ldap.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ldap.CreateConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.ldap.LdapServicePromiseClient.prototype.createConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ldap.LdapService/CreateConnection',
      request,
      metadata || {},
      methodInfo_LdapService_CreateConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ldap.DeleteConnectionRqst,
 *   !proto.ldap.DeleteConnectionRsp>}
 */
const methodInfo_LdapService_DeleteConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ldap.DeleteConnectionRsp,
  /** @param {!proto.ldap.DeleteConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ldap.DeleteConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.ldap.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ldap.DeleteConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ldap.DeleteConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ldap.LdapServiceClient.prototype.deleteConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ldap.LdapService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_LdapService_DeleteConnection,
      callback);
};


/**
 * @param {!proto.ldap.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ldap.DeleteConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.ldap.LdapServicePromiseClient.prototype.deleteConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ldap.LdapService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_LdapService_DeleteConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ldap.CloseRqst,
 *   !proto.ldap.CloseRsp>}
 */
const methodInfo_LdapService_Close = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ldap.CloseRsp,
  /** @param {!proto.ldap.CloseRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ldap.CloseRsp.deserializeBinary
);


/**
 * @param {!proto.ldap.CloseRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ldap.CloseRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ldap.CloseRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ldap.LdapServiceClient.prototype.close =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ldap.LdapService/Close',
      request,
      metadata || {},
      methodInfo_LdapService_Close,
      callback);
};


/**
 * @param {!proto.ldap.CloseRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ldap.CloseRsp>}
 *     A native promise that resolves to the response
 */
proto.ldap.LdapServicePromiseClient.prototype.close =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ldap.LdapService/Close',
      request,
      metadata || {},
      methodInfo_LdapService_Close);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ldap.SearchRqst,
 *   !proto.ldap.SearchResp>}
 */
const methodInfo_LdapService_Search = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ldap.SearchResp,
  /** @param {!proto.ldap.SearchRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ldap.SearchResp.deserializeBinary
);


/**
 * @param {!proto.ldap.SearchRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ldap.SearchResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ldap.SearchResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ldap.LdapServiceClient.prototype.search =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ldap.LdapService/Search',
      request,
      metadata || {},
      methodInfo_LdapService_Search,
      callback);
};


/**
 * @param {!proto.ldap.SearchRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ldap.SearchResp>}
 *     A native promise that resolves to the response
 */
proto.ldap.LdapServicePromiseClient.prototype.search =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ldap.LdapService/Search',
      request,
      metadata || {},
      methodInfo_LdapService_Search);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ldap.AuthenticateRqst,
 *   !proto.ldap.AuthenticateRsp>}
 */
const methodInfo_LdapService_Authenticate = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ldap.AuthenticateRsp,
  /** @param {!proto.ldap.AuthenticateRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ldap.AuthenticateRsp.deserializeBinary
);


/**
 * @param {!proto.ldap.AuthenticateRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ldap.AuthenticateRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ldap.AuthenticateRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ldap.LdapServiceClient.prototype.authenticate =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ldap.LdapService/Authenticate',
      request,
      metadata || {},
      methodInfo_LdapService_Authenticate,
      callback);
};


/**
 * @param {!proto.ldap.AuthenticateRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ldap.AuthenticateRsp>}
 *     A native promise that resolves to the response
 */
proto.ldap.LdapServicePromiseClient.prototype.authenticate =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ldap.LdapService/Authenticate',
      request,
      metadata || {},
      methodInfo_LdapService_Authenticate);
};


module.exports = proto.ldap;

