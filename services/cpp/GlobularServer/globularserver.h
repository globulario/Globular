#ifndef GLOBULARSERVER_H
#define GLOBULARSERVER_H



#include <string>
#include <sstream>

// gRpc stuff...
#include <grpc++/grpc++.h>
#include <grpcpp/support/server_interceptor.h>
#include "globularressourceclient.h"
using grpc::Service;
using grpc::Server;

namespace Globular {
/**
 * @brief That class contain the base for globular service.
 * It take care to get the basic attribute to make the service manageable.
 */
class GlobularService
{
protected:

    // The id of the service instance, must be unique on the globular server.
    std::string id;

    // The name of the service, must be the same from proto file.
    std::string name;

    // The path of the executable
    std::string path;

    // The path of the .proto path.
    std::string proto;

    // The grpc port
    unsigned int port;

    // The reverse proxy.
    unsigned int proxy;

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

    // The configuration port.
    int configurationPort;

    // The root path
    std::string root;

    std::unique_ptr<Server> server;

public:

    // The default constructor.
    GlobularService(std::string id,
                    std::string name,
                    std::string domain = "localhost",
                    std::string publisher_id = "localhost",
                    bool allow_all_origins = false,
                    std::string allowed_origins = "",
                    std::string version = "0.0.1",
                    bool tls = true,
                    unsigned int defaultPort = 10023,
                    unsigned int defaultProxy = 10024,
                    bool keep_alive = false,
                    bool keep_up_to_date = false
            );

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


    /**
     * @brief run Start listen and processing request.
     */
    void run(Service*);

    /**
     * Stop the server.
     */
    void stop();
};


/**
 * Intercept method call and validate application/user permission to execute method or ressource access.
 * @brief The ServerUnaryInterceptor class
 */
class ServerInterceptor : public grpc::experimental::Interceptor
{

public:
    ServerInterceptor( grpc::experimental::ServerRpcInfo* info) {
        info_ = info;

        // Check the method name and compare to the type
        const char* method = info->method();
        grpc::experimental::ServerRpcInfo::Type type = info->type();
        ressourceClient = 0;

    }

    /**
      * @brief Intercept Intercept method and validate access.
      * @param methods The intercepted method
      */
    void Intercept(grpc::experimental::InterceptorBatchMethods* methods) override {

        std::string method = this->info_->method();

        if (methods->QueryInterceptionHookPoint(
                    grpc::experimental::InterceptionHookPoints::POST_RECV_INITIAL_METADATA)) {
            auto* map = methods->GetRecvInitialMetadata();
            if(map->size()==0){
                methods->Proceed();
                return;
            }

            domain = getMetadata("domain", map);
            application = getMetadata("application", map);
            token = getMetadata("token", map);
        }

        if (methods->QueryInterceptionHookPoint(
                    grpc::experimental::InterceptionHookPoints::PRE_SEND_STATUS)) {

            auto hasAccess = false;
            if(domain.empty()){
                grpc::Status error(grpc::StatusCode::PERMISSION_DENIED, "Permission denied to execute " + method + " no domain was given!");
                methods->ModifySendStatus(error);
                methods->Proceed();
                return;
            }

            // Here I will create the ressource client.
            if(ressourceClient == 0){
                std::cout << "create ressouce client for domain: " << domain  << std::endl;
                auto index = domain.find(":");
                auto port = 80;
                if(index != 0){
                    port = atoi(domain.substr(index+1).c_str());
                    domain = domain.substr(0, index);
                    std::cout << "port" << domain.substr(index) << std::endl;
                    std::cout << "index "<< index << std::endl;
                    std::cout <<"domain " << domain << std::endl;
                    std::cout << "port int " << port << std::endl;
                }
                ressourceClient = new Globular::RessourceClient("ressource.RessourceService", domain, port);
            }

            if(!application.empty()){
                hasAccess = ressourceClient->validateApplicationAccess(application, method);
            }

            if(!hasAccess){
                std::cout << method << token << std::endl;
                hasAccess = ressourceClient->validateUserAccess(token, method);
            }

            if(!hasAccess){
                grpc::Status error(grpc::StatusCode::PERMISSION_DENIED, "Permission denied to execute " + method + "!");
                methods->ModifySendStatus(error);
                methods->Proceed();
                return;
            }else{

                // Now if the action has ressource access permission defines...
                auto permissions = ressourceClient->getActionPermission(method);
                for(auto it = permissions.cbegin(); it != permissions.cend(); ++it){

                    auto permission = *(it);
                    std::string path = "";

                    // permission.permission()
                    auto hasRessourcePermission = ServerInterceptor::ressourceClient->validateUserRessourceAccess(token, path, method, permission.permission());
                    if (!hasRessourcePermission)
                    {
                        ServerInterceptor::ressourceClient->validateApplicationRessourceAccess(application, path, method, permission.permission());
                    }
                    if (!hasRessourcePermission)
                    {
                        grpc::Status error(grpc::StatusCode::PERMISSION_DENIED, "Permission denied access denied on ressource " + path + "!");
                        methods->ModifySendStatus(error);
                        methods->Proceed();
                        return;
                    }

                }
            }
        }

        // run the method.
        methods->Proceed();
    }

private:
    std::string getMetadata(std::string key, std::multimap<grpc::string_ref, grpc::string_ref>* map){
        // Here I will retreive the metadata....
        bool found = false;
        for (const auto& pair : *map) {
            found = pair.first.find(key) == 0;

            if (found){
                return  std::string(pair.second.data());
            }
        }

        return "";
    }

    grpc::experimental::ServerRpcInfo* info_;
    std::string application;
    std::string token;
    std::string path;
    std::string domain;

    // The ressource client.
    RessourceClient* ressourceClient;

}; // namespace Globular.

/**
 * @brief The ServerInterceptorFactory class
 * Factory class to create Server interceptor.
 */
class ServerInterceptorFactory
        : public grpc::experimental::ServerInterceptorFactoryInterface {
public:
    virtual grpc::experimental::Interceptor* CreateServerInterceptor(
            grpc::experimental::ServerRpcInfo* info) override {
        return new ServerInterceptor(info);
    }
};

}

#endif // GLOBULARSERVER_H
