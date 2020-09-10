#include "globularserver.h"
#include <string>
#include <fstream>
#include <iostream>
#include <sstream>
#include <map>
#include <cstddef>
#include <bitset>         // std::bitset
#include <math.h>       /* ceil */
#include "json.hpp"
#include <fstream>

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
                                           bool tls,
                                           unsigned int defaultPort,
                                           unsigned int defaultProxy):
    id(id),
    name(name),
    domain(domain),
    publisher_id(publisher_id),
    allow_all_origins(allow_all_origins),
    port(defaultPort),
    proxy(defaultProxy),
    allowed_origins(allowed_origins),
    tls(tls)
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
    }

    // Now I will create the new ressource client.
    this->ressourceClient = new RessourceClient("ressource.RessourceService", this->domain);
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

    std::string execPath = getexepath();
    j["Path"] = execPath;

#ifdef WIN32
    std::size_t lastIndex = execPath.find_last_of("/\\");
    std::string configPath = execPath.substr(0, lastIndex) + "\\config.json";
#else
    std::size_t lastIndex = execPath.find_last_of("/");
    std::string configPath = execPath.substr(0, lastIndex) + "/config.json";
#endif
    std::ofstream file;
    file.open(configPath);
    file << j.dump();
    file.close();
}
