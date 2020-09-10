#ifndef ECHOSERVICEIMPL_H
#define ECHOSERVICEIMPL_H

#include "echo/echopb/echo.grpc.pb.h"
#include "echo/echopb/echo.pb.h"
#include "globularserver.h"

class EchoServiceImpl final: public echo::EchoService::Service, Globular::GlobularService
{
public:
    EchoServiceImpl(std::string id = "",
                    std::string name = "",
                    std::string domain = "localhost",
                    std::string publisher_id = "localhost",
                    bool allow_all_origins = false,
                    std::string allowed_origins = "",
                    bool tls = false,
                    unsigned int defaultPort = 10023, unsigned int defaultProxy = 10024);
};

#endif // ECHOSERVICEIMPL_H
