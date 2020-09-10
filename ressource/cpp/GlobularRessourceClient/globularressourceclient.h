#ifndef GLOBULARRESSOURCECLIENT_H
#define GLOBULARRESSOURCECLIENT_H
#include "globularclient.h"
#include <thread>

#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/client_context.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>

#include "ressource/ressource.pb.h"
#include "ressource/ressource.grpc.pb.h"

// GRPC stuff.
using grpc::Channel;
using grpc::Channel;
using grpc::ClientContext;
using grpc::ClientReader;
using grpc::ClientReaderWriter;
using grpc::ClientWriter;
using grpc::Status;

namespace Globular {

class RessourceClient : Client
{
    // the underlying grpc ressource client.
    std::unique_ptr<ressource::RessourceService::Stub> stub_;
    std::shared_ptr<Channel> channel;

public:
    RessourceClient(std::string name, std::string domain, unsigned int configurationPort=10000);
};

}

#endif // GLOBULARRESSOURCECLIENT_H
