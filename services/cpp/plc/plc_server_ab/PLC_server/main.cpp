
#include "PlcServiceImpl.h"
#include "cxxopts.hpp" // argument parser.

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;

using namespace std;

int main(int argc, char** argv) {

	cxxopts::Options options("plc service", "A gRpc service to communicate with PLC.");
	
	auto result = options.parse(argc, argv);

	PlcServiceImpl service;
	std::stringstream ss;
	ss << "0.0.0.0" << ":" << service.getDefaultPort();

	ServerBuilder builder;
	// Listen on the given address without any authentication mechanism.
	builder.AddListeningPort(ss.str(), grpc::InsecureServerCredentials());

	// Register "service" as the instance through which we'll communicate with
	// clients. In this case it corresponds to an *synchronous* service.
	builder.RegisterService(&service);
	// Finally assemble the server.
	std::unique_ptr<Server> server(builder.BuildAndStart());
	std::cout << "Server listening on " << ss.str() << std::endl;

	// Wait for the server to shutdown. Note that some other thread must be
	// responsible for shutting down the server for this call to ever return.
	server->Wait();

	return 0;
}
