#include "echoserviceimpl.h"

EchoServiceImpl::EchoServiceImpl(std::string id,
                                 std::string name,
                                 std::string domain,
                                 std::string publisher_id,
                                 bool allow_all_origins,
                                 std::string allowed_origins,
                                 std::string version,
                                 bool tls,
                                 unsigned int defaultPort, unsigned int defaultProxy):
    Globular::GlobularService(id, name, domain, publisher_id, allow_all_origins, allowed_origins, version, tls, defaultPort, defaultProxy )
{
    // Set the proto path if is not already set.
    if(this->proto.length() == 0){
        EchoRequest request;
        this->proto = this->root + "/" + request.GetDescriptor()->file()->name();
        this->save();
    }
}

Status EchoServiceImpl::Echo(ServerContext* context, const EchoRequest* request, EchoResponse* response) {
    std::string prefix("Echo ");
    response->set_message(prefix + request->message());
    return Status::OK;
}
