#include "globularserver.h"
#include <string>
#include <vector>
#include <fstream>
#include <iostream>
#include <sstream>
#include <map>
#include <cstddef>
#include <bitset>         // std::bitset
#include <math.h>       /* ceil */
#include "../json.hpp"
#include <fstream>
#include <filesystem>


using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;

std::string replaceAll(std::string str, const std::string& from, const std::string& to) {
    size_t start_pos = 0;
    while((start_pos = str.find(from, start_pos)) != std::string::npos) {
        str.replace(start_pos, from.length(), to);
        start_pos += to.length(); // Handles case where 'to' is a substring of 'from'
    }
    return str;
}


#ifdef WIN32
#include <windows.h>
std::string getexepath()
{
    char result[MAX_PATH];
    return std::string(result, GetModuleFileNameA(NULL, result, MAX_PATH));
}

void sleep(unsigned milliseconds)
{
    Sleep(milliseconds);
}

#else
#include <limits.h>
#include <unistd.h>
#include <linux/limits.h>

std::string getexepath()
{
    char result[PATH_MAX];
    ssize_t count = readlink("/proc/self/exe", result, PATH_MAX);
    return std::string(result, (count > 0) ? count : 0);
}

#endif // WIN32

Globular::GlobularService::GlobularService(std::string id,
                                           std::string name,
                                           std::string domain,
                                           std::string publisher_id,
                                           bool allow_all_origins,
                                           std::string allowed_origins,
                                           std::string version,
                                           bool tls,
                                           unsigned int defaultPort,
                                           unsigned int defaultProxy,
                                           bool keep_alive,
                                           bool keep_up_to_date
                                           ):
    id(id),
    name(name),
    domain(domain),
    publisher_id(publisher_id),
    allow_all_origins(allow_all_origins),
    port(defaultPort),
    proxy(defaultProxy),
    allowed_origins(allowed_origins),
    version(version),
    tls(tls),
    keep_alive(keep_alive),
    keep_up_to_date(keep_up_to_date)
{

    // first of all I will try to open the configuration from the file.
    std::ifstream inFile;
    std::string execPath = getexepath();
#ifdef WIN32
    std::size_t lastIndex = execPath.find_last_of("/\\");
    std::string configPath = execPath.substr(0, lastIndex) + "\\config.json";
#else
    std::size_t lastIndex = execPath.find_last_of("/");
    std::string configPath = execPath.substr(0, lastIndex) + "/config.json";
#endif
    inFile.open(configPath); //open the input file

    if (inFile.good()) {

        std::stringstream strStream;
        strStream << inFile.rdbuf(); //read the file
        std::string jsonStr = strStream.str(); //str holds the content of the file
        std::cout << "Load configuration file " << configPath << std::endl;
        // Parse the json file.
        auto j = nlohmann::json::parse(jsonStr);

        // Now I will initialyse the value from the configuration file.
        this->publisher_id = j["PublisherId"];
        this->version = j["Version"];
        this->keep_up_to_date = j["KeepUpToDate"];
        this->allow_all_origins = j["AllowAllOrigins"];
        this->cert_authority_trust = j["CertAuthorityTrust"];
        this->keep_alive = j["KeepAlive"];
        this->cert_file = j["CertFile"];
        this->domain = j["Domain"];
        this->key_file = j["KeyFile"];
        this->name = j["Name"];
        this->port = j["Port"];
        this->proxy = j["Proxy"];
        this->path = j["Path"];
        this->proto = j["Proto"];
        this->tls = j["TLS"];
        this->protocol = j["Protocol"];


        // can be a list of string
        this->allowed_origins = j["AllowedOrigins"];
    }else {
        // Set the application path.
        this->path = replaceAll(getexepath(), "\\", "/");
        this->save();
    }

    // Set the root path and configuration port from the local server configuration...

    auto path = std::filesystem::temp_directory_path().string() + "GLOBULAR_ROOT";
    std::ifstream input(path);
    std::stringstream sstr;

    while(input >> sstr.rdbuf());
    this->root = replaceAll(sstr.str().substr(0, sstr.str().find_last_of(":")), "\\", "/");

}

void
read ( const std::string& filename, std::string& data )
{
        std::ifstream file ( filename.c_str (), std::ios::in );

    if ( file.is_open () )
    {
        std::stringstream ss;
        ss << file.rdbuf ();

        file.close ();

        data = ss.str ();
    }

    return;
}



void Globular::GlobularService::save() {
    nlohmann::json j;
    j["PublisherId"] = this->publisher_id;
    j["Version"] = this->version;
    j["KeepUpToDate"] = this->keep_up_to_date;
    j["KeepAlive"] = this->keep_alive;
    j["AllowAllOrigins"] = this->allow_all_origins;
    j["AllowedOrigins"] = this->allowed_origins; // empty string
    j["CertAuthorityTrust"] = this->cert_authority_trust;
    j["CertFile"] = this->cert_file;
    j["Domain"] = this->domain;
    j["KeyFile"] = this->key_file;
    j["Name"] = this->name;
    j["Port"] = this->port;
    j["Id"] = this->id;
    j["Protocol"] = "grpc";
    j["Proto"] = this->proto;
    j["Proxy"] = this->proxy;
    j["TLS"] = this->tls;
    j["Path"] = this->path;

#ifdef WIN32
    std::size_t lastIndex = this->path.find_last_of("/\\");
    std::string configPath = this->path.substr(0, lastIndex) + "\\config.json";
#else
    std::size_t lastIndex = this->path.find_last_of("/");
    std::string configPath = this->path.substr(0, lastIndex) + "/config.json";
#endif

    std::ofstream file;
    file.open(configPath);
    file << j.dump();
    file.close();
}

// use it for shutdown only...
extern  std::unique_ptr<grpc::Server> server;

void Globular::GlobularService::stop(){

    server->Shutdown();
}


void Globular::GlobularService::run(Service* s) {
    ServerBuilder builder;
    std::stringstream ss;
    ss <<  "0.0.0.0" << ":" << this->port;

    if(this->tls){
        std::string key;
        std::string cert;
        std::string ca;

        read ( this->cert_file, cert );
        read ( this->key_file , key );
        read ( this->cert_authority_trust, ca );

        grpc::SslServerCredentialsOptions::PemKeyCertPair keycert =
        {
            key,
            cert
        };

        grpc::SslServerCredentialsOptions sslOps;
        sslOps.pem_root_certs = ca;
        sslOps.pem_key_cert_pairs.push_back ( keycert );

        builder.AddListeningPort(ss.str(), grpc::SslServerCredentials( sslOps ));

    }else{
        // Listen on the given address without any authentication mechanism.
        builder.AddListeningPort(ss.str(), grpc::InsecureServerCredentials());
    }


    // Register "service" as the instance through which we'll communicate with
    // clients. In this case it corresponds to an *synchronous* service.
    builder.RegisterService(s);

    // Set the interceptor creator.
    std::vector<std::unique_ptr<grpc::experimental::ServerInterceptorFactoryInterface>> creators;
    creators.push_back(std::unique_ptr<Globular::ServerInterceptorFactory>(new Globular::ServerInterceptorFactory()));
    builder.experimental().SetInterceptorCreators(std::move(creators));

    // Finally assemble the server.
    // std::unique_ptr<Server> server(builder.BuildAndStart());
    server = builder.BuildAndStart();

    std::cout << "Server listening on " << ss.str() << std::endl;

    // Wait for the server to shutdown. Note that some other thread must be
    // responsible for shutting down the server for this call to ever return.
    server->Wait();
}
