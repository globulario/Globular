#include "echoserviceimpl.h"

EchoServiceImpl::EchoServiceImpl(std::string id,
                                 std::string name,
                                 std::string domain,
                                 std::string publisher_id,
                                 bool allow_all_origins,
                                 std::string allowed_origins,
                                 bool tls,
                                 unsigned int defaultPort, unsigned int defaultProxy):
    Globular::GlobularService(id, name, domain, publisher_id, allow_all_origins, allowed_origins, tls, defaultPort, defaultProxy )
{

}
