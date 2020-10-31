#ifndef ECHOSERVICEIMPL_H
#define ECHOSERVICEIMPL_H

#include "echo.grpc.pb.h"
#include "echo.pb.h"
#include "globularserver.h"

#include <grpc++/grpc++.h>
using grpc::ServerContext;
using grpc::Status;
using echo::EchoResponse;
using echo::EchoRequest;

class EchoServiceImpl final: public echo::EchoService::Service, public Globular::GlobularService
{
public:
    EchoServiceImpl(std::string id = "echo",
                    std::string name = "echo.EchoService",
                    std::string domain = "localhost",
                    std::string publisher_id = "localhost",
                    bool allow_all_origins = false,
                    std::string allowed_origins = "",
                    std::string version = "0.0.1",
                    bool tls = false,
                    unsigned int defaultPort = 10023, unsigned int defaultProxy = 10024);

    Status Echo(ServerContext* /*context*/, const EchoRequest* /*request*/, EchoResponse* /*response*/) override;

};

#endif // ECHOSERVICEIMPL_H
