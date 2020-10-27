#ifndef GLOBULARRESSOURCECLIENT_H
#define GLOBULARRESSOURCECLIENT_H
#include "globularclient.h"
#include <thread>

#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/client_context.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>

#include "ressource.pb.h"
#include "ressource.grpc.pb.h"

// GRPC stuff.
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

public:

    // The constructor.
    RessourceClient(std::string name, std::string domain="localhost", unsigned int configurationPort=8080);

    // Now the ressource client functionnalites.

    /**
     * @brief authenticate Authenticate a user on the services.
     * @param user The user id
     * @param password The user password
     * @return the token (valid for a given delay)
     */
    std::string authenticate(std::string user, std::string password);

    /**
     * @brief validateUserAccess Validate if a user can access or not a given method.
     * @param token The token receive from authenticate.
     * @param method The gRpc method path. ex. /module/methodName/
     * @return
     */
    bool validateUserAccess(std::string token, std::string method);


    /**
     * @brief validateApplicationAccess Validate if application can access or not a given method.
     * @param name The name the application (must be unique on the server).
     * @param method The gRpc method path. ex. /module/methodName/
     * @return
     */
    bool validateApplicationAccess(std::string name, std::string method);

    /**
     * @brief validateUserRessourceAccess Validate if user can access a given ressource on the server.
     * @param token The token received from authentication
     * @param path The path of the ressource (must be unique on the server)
     * @param method The gRpc method path. ex. /module/methodName/
     * @param permission The permission number, see chmod number... (ReadWriteDelete)
     * @return
     */
    bool validateUserRessourceAccess(std::string token, std::string path, std::string method, int32_t permission);

    /**
     * @brief validateApplicationRessourceAccess
     * @param application The name of the application to be validate.
     * @param path The path of the ressource (must be unique on the server)
     * @param method The gRpc method path. ex. /module/methodName/
     * @param permission The permission number, see chmod number... (ReadWriteDelete)
     * @return
     */
    bool validateApplicationRessourceAccess(std::string application, std::string path, std::string method, int32_t permission);


    /**
     * @brief SetRessource
     * @param path The path of the ressource (must be unique on the server)
     * @param name The name of the ressource
     * @param modified The modified date.
     * @param size The size of the ressource.
     */
    void SetRessource(std::string path, std::string name, int modified, int size);

    /**
     * @brief removeRessouce
     * @param path The path where the ressource is located.
     * @param name The name of the ressource.
     */
    void removeRessouce(std::string path, std::string name);


    /**
     * @brief getActionPermission
     * @param method The gRpc method path. ex. /module/methodName/
     * @return
     */
    std::vector<::ressource::ActionParameterRessourcePermission> getActionPermission(std::string method);

    /**
     * @brief Log
     * @param application The application name
     * @param method The gRpc method path. ex. /module/methodName/
     * @param message The message to log.
     * @param type can be 0 for INFO_MESSAGE and 1 for ERROR_MESSAGE.
     */
    void Log(std::string application, std::string method, std::string message, int type = 0);
};

}

#endif // GLOBULARRESSOURCECLIENT_H
