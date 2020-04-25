/**
 * @fileoverview gRPC-Web generated client stub for smtp
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.smtp = require('./smtp_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smtp.SmtpServiceClient =
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
proto.smtp.SmtpServicePromiseClient =
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
 *   !proto.smtp.CreateConnectionRqst,
 *   !proto.smtp.CreateConnectionRsp>}
 */
const methodInfo_SmtpService_CreateConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smtp.CreateConnectionRsp,
  /** @param {!proto.smtp.CreateConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.smtp.CreateConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.smtp.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smtp.CreateConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smtp.CreateConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smtp.SmtpServiceClient.prototype.createConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smtp.SmtpService/CreateConnection',
      request,
      metadata || {},
      methodInfo_SmtpService_CreateConnection,
      callback);
};


/**
 * @param {!proto.smtp.CreateConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smtp.CreateConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.smtp.SmtpServicePromiseClient.prototype.createConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smtp.SmtpService/CreateConnection',
      request,
      metadata || {},
      methodInfo_SmtpService_CreateConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smtp.DeleteConnectionRqst,
 *   !proto.smtp.DeleteConnectionRsp>}
 */
const methodInfo_SmtpService_DeleteConnection = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smtp.DeleteConnectionRsp,
  /** @param {!proto.smtp.DeleteConnectionRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.smtp.DeleteConnectionRsp.deserializeBinary
);


/**
 * @param {!proto.smtp.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smtp.DeleteConnectionRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smtp.DeleteConnectionRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smtp.SmtpServiceClient.prototype.deleteConnection =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smtp.SmtpService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_SmtpService_DeleteConnection,
      callback);
};


/**
 * @param {!proto.smtp.DeleteConnectionRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smtp.DeleteConnectionRsp>}
 *     A native promise that resolves to the response
 */
proto.smtp.SmtpServicePromiseClient.prototype.deleteConnection =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smtp.SmtpService/DeleteConnection',
      request,
      metadata || {},
      methodInfo_SmtpService_DeleteConnection);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smtp.SendEmailRqst,
 *   !proto.smtp.SendEmailRsp>}
 */
const methodInfo_SmtpService_SendEmail = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smtp.SendEmailRsp,
  /** @param {!proto.smtp.SendEmailRqst} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.smtp.SendEmailRsp.deserializeBinary
);


/**
 * @param {!proto.smtp.SendEmailRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smtp.SendEmailRsp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smtp.SendEmailRsp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smtp.SmtpServiceClient.prototype.sendEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smtp.SmtpService/SendEmail',
      request,
      metadata || {},
      methodInfo_SmtpService_SendEmail,
      callback);
};


/**
 * @param {!proto.smtp.SendEmailRqst} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smtp.SendEmailRsp>}
 *     A native promise that resolves to the response
 */
proto.smtp.SmtpServicePromiseClient.prototype.sendEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smtp.SmtpService/SendEmail',
      request,
      metadata || {},
      methodInfo_SmtpService_SendEmail);
};


module.exports = proto.smtp;

