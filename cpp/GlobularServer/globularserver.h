#ifndef GLOBULARSERVER_H
#define GLOBULARSERVER_H

#include <string>
#include <sstream>
#include "globularressourceclient.h"

namespace Globular {

/**
 * @brief That class contain the base for globular service.
 * It take care to get the basic attribute to make the service manageable.
 */
class GlobularService
{
    // The id of the service instance, must be unique on the globular server.
    std::string id;

    // The name of the service, must be the same from proto file.
    std::string name;

    // The path of the executable
    std::string path;

    // The path of the .proto path.
    std::string proto;

    // The grpc port
    unsigned port;

    // The reverse proxy.
    unsigned proxy;

    // GRPC
    std::string protocol;

    // The domain
    std::string domain;

    // The publisher id
    std::string publisher_id;

    // The service version
    std::string version;

    // Keep service up to date if new version came out.
    bool keep_up_to_date;

    // Restart the service if it stop.
    bool keep_alive;

    // allow all origin
    bool allow_all_origins;

    // comma separated list of origin.
    std::string allowed_origins;

    // true if connection use TLS
    bool tls;

    // Now certificates.

    // The CA certificate
    std::string cert_authority_trust;

    // The private key file
    std::string key_file;

    // The server certificate.
    std::string cert_file;

    // The ressource client.
    Globular::RessourceClient *ressourceClient;

public:

    // The default constructor.
    GlobularService(std::string id,
                    std::string name,
                    std::string domain = "localhost",
                    std::string publisher_id = "localhost",
                    bool allow_all_origins = false,
                    std::string allowed_origins = "",
                    bool tls = true,
                    unsigned int defaultPort = 10023, unsigned int defaultProxy = 10024);

    // Getter/Setter
    const std::string& getName() {
        return this->name;
    }

    const std::string getAddress() {
        std::stringstream ss;
        ss << this->domain << ":" << this->port;
        return  ss.str();
    }

    unsigned int getDefaultPort() {
        return this->port;
    }

    unsigned int getDefaulProxy() {
        return this->proxy;
    }

    /**
     * @brief save Save the service configuration.
     */
    void save();
};

}

#endif // GLOBULARSERVER_H
