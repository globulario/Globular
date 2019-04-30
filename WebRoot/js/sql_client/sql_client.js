/**
 * Sql gRPC client 
 */
const {EchoRequest,ServerStreamingEchoRequest} = require('./echo_pb.js');
const {EchoServiceClient} = require('./echo_grpc_web_pb.js');
const {EchoApp} = require('../echoapp.js');

const grpc = {};
grpc.web = require('grpc-web');

// 
var echoService = new EchoServiceClient('http://'+window.location.hostname+':8080', null, null);