#include "globularclient.h"
#include "HTTPRequest.hpp"
#include <sstream>
#include <fstream>
#include <iostream>
#include <cstdio>
#include <memory>
#include <stdexcept>
#include <string>
#include <array>
#include <filesystem>

//  https://github.com/nlohmann/json
#include "json.hpp"
#include "Base64.h"
#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>


/**
 * @brief writeAllText Write a text file into a given path.
 * @param path The path where to write the file.
 * @param text The text to save.
 */
void writeAllText(std::string path, std::string text){
    std::ofstream f;
    f.open(path);
    f << text;
    f.close();
}

/**
 * @brief readAllText Read all text from a given file
 * @param path The path of the file to be read.
 * @return
 */
std::string readAllText(std::string path){
    std::ifstream t(path);
    std::string str;

    t.seekg(0, std::ios::end);
    str.reserve(t.tellg());
    t.seekg(0, std::ios::beg);

    str.assign((std::istreambuf_iterator<char>(t)),
               std::istreambuf_iterator<char>());

    return str;
}

std::string getTempDir(){
    // std::filesystem::temp_directory_path()
    return std::filesystem::temp_directory_path().string();
}

bool exists(std::string path){
    //
    std::ifstream infile(path);
    return infile.good();
}

void createDir(std::string path){
    std::filesystem::create_directory(path);
}

std::string getPathName(const std::string& s) {

    char sep = '/';

#ifdef _WIN32
    sep = '\\';
#endif

    size_t i = s.rfind(sep, s.length());
    if (i != std::string::npos) {
        return(s.substr(0, i));
    }

    return("");
}
Globular::Client::Client(std::string name, std::string domain, unsigned int configurationPort)
{
    this->config = new Globular::ServiceConfig();
    this->config->Name = name;
    this->config->Domain = domain;

    // init internal values.
    this->init(configurationPort);
    std::cout << "Domain config port " << configurationPort << std::endl;
    std::stringstream ss;
    ss << this->config->Domain << ":" << this->config->Port;

    // Now I will create the grpc channel.
    if(this->getTls()){
        grpc::SslCredentialsOptions opts =
        {
            readAllText( this->config->CertAuthorityTrust),
            readAllText( this->config->KeyFile),
            readAllText( this->config->CertFile)
        };
        std::cout << "try to open channel with address " << ss.str() << std::endl;

        this->channel = grpc::CreateChannel(ss.str(), grpc::SslCredentials ( opts ) );

    }else{
        // Create insecure connection to the service.
        this->channel = grpc::CreateChannel(ss.str(), grpc::InsecureChannelCredentials());
        if(this->channel){
            std::cout << "client channel is now initialysed!" << std::endl;
            std::cout << this->channel->GetServiceConfigJSON()	<< std::endl;
        }
    }

}

void Globular::Client::init(unsigned int configurationPort){

    // Initialyse client stuff here.
    this->initServiceConfig(configurationPort);

}

void Globular::Client::initServiceConfig(unsigned int configurationPort){

    std::stringstream ss;
    ss << "http://" << this->config->Domain << ":" << configurationPort << "/config";
    std::cout << "get configuration from " << ss.str() << std::endl;;
    http::Request request(ss.str());
    const http::Response response = request.send("GET");
    ss.flush();
    ss << std::string(response.body.begin(), response.body.end()) << '\n'; // print the result
    std::string jsonStr = ss.str();
    size_t index = jsonStr.find_first_of("{");
    jsonStr = jsonStr.substr(index, jsonStr.length() - index);
    auto j = nlohmann::json::parse(jsonStr);

    // Now I will initialyse the value from the configuration file.
    auto services = j["Services"];

    //
    for (auto it = services.begin(); it != services.end(); ++it)
    {
        if(it.key() == this->config->Name || (*it)["Id"].get<std::string>() == this->config->Name){

            this->config->Id = (*it)["Id"].get<std::string>();
            this->config->Name = (*it)["Name"].get<std::string>();
            this->config->Domain = (*it)["Domain"].get<std::string>();
            this->config->Description = (*it)["Domain"].get<std::string>();
            this->config->Port = (*it)["Port"].get<unsigned int>();
            this->config->Proxy = (*it)["Proxy"].get<unsigned int>();
            this->config->TLS= (*it)["TLS"].get<bool>();
            this->config->Keywords = (*it)["Keywords"].get<std::vector<std::string>>();
            this->config->Discoveries = (*it)["Discoveries"].get<std::vector<std::string>>();
            this->config->Repositories = (*it)["Repositories"].get<std::vector<std::string>>();

            if(this->config->TLS){
                // The service is secure.
                std::stringstream ss;
                ss << getTempDir() <<   "/" << this->config->Domain + "_token";
                auto path = ss.str();
                std::cout << "read: " << path << std::endl;
                if(exists(path)){
                    // TODO make correction here the CertFile and KeyFile are the one of the server not the client.
                    this->config->CertAuthorityTrust =  (*it)["CertAuthorityTrust"].get<std::string>();
                    std::string path = getPathName(this->config->CertAuthorityTrust);

                    this->config->CertFile =  path + "/client.crt";
                    this->config->KeyFile =   path + "/client.pem";;

                    std::cout << this->config->KeyFile << std::endl;
                    std::cout << this->config->CertFile << std::endl;
                    std::cout << this->config->CertAuthorityTrust << std::endl;
                }else{
                    ss.flush();
                    ss << getTempDir() <<   "/config/tls/" << this->config->Domain;
                    auto path = ss.str();
                    // Here I will create a directory named
                    if(!exists(path)){
                        createDir(path);
                    }

                    // I will create a certificate request and make it sign by the server.
                    auto ca_crt = this->getCaCertificate(this->config->Domain, configurationPort);
                    writeAllText(path + "/ca.crt", ca_crt);

                    // The password must be store in the client configuration...
                    auto pwd = "1111";

                    // Now I will generate the certificate for the client...
                    // Step 1: Generate client private key.
                    this->generateClientPrivateKey(path, pwd);

                    // Step 2: Generate the client signing request.
                    this->generateClientCertificateSigningRequest(path, this->config->Domain);


                    // Step 3: Generate client signed certificate.
                    auto client_csr = readAllText(path + "/client.csr");
                    auto client_crt = this->signCaCertificate(this->config->Domain, configurationPort, client_csr);
                    writeAllText(path + "/client.crt", client_crt);

                    // Step 4: Convert client.key to pem file.
                    this->keyToPem("client", path, pwd);

                    // Set path in the config.
                    this->config->KeyFile= path + "/client.key";
                    this->config->CertAuthorityTrust  = path + "/ca.crt";
                    this->config->CertFile  = path + "/client.crt";
                }
            }
            break; // exit the loop
        }
    }
}

std::string Globular::Client::getName() const
{
    return this->config->Name;
}

std::string Globular::Client::getDomain() const
{
    return this->config->Domain;
}

unsigned int Globular::Client::getPort() const
{
    return this->config->Port;
}

bool Globular::Client::getTls() const
{
    return this->config->TLS;
}

std::string Globular::Client::getCaFile() const
{
    return this->config->CertAuthorityTrust;
}

void Globular::Client::setCaFile(const std::string &value)
{
    this->config->CertAuthorityTrust = value;
}

std::string Globular::Client::getKeyFile() const
{
    return this->config->KeyFile;
}

void Globular::Client::setKeyFile(const std::string &value)
{
    this->config->KeyFile = value;
}

std::string Globular::Client::getCertFile() const
{
    return this->config->CertFile;
}

void Globular::Client::setCertFile(const std::string &value)
{
    this->config->CertFile = value;
}

void Globular::Client::setTls(bool value)
{
    this->config->TLS = value;
}

void Globular::Client::close()
{
}

void Globular::Client::getClientContext(ClientContext& context, std::string token, std::string application, std::string domain, std::string path){
    if(domain.empty()){
        context.AddMetadata("domain", this->config->Domain);
        domain = this->config->Domain;
    }else{
        context.AddMetadata("domain", domain);
    }

    if(token.empty()){
        std::stringstream ss;
        ss << getTempDir() << "/" << domain << "_token";
        if(exists(ss.str())){
            token = readAllText(ss.str());
            context.AddMetadata("token", token);
        }
    }else{
        context.AddMetadata("token", token);
    }

    if(!application.empty()){
        context.AddMetadata("application", application);
    }

    if(!path.empty()){
        context.AddMetadata("path", path);
    }
}

std::string exec(const char* cmd) {
    std::array<char, 128> buffer;
    std::string result;
    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd, "r"), pclose);
    if (!pipe) {
        throw std::runtime_error("popen() failed!");
    }
    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }
    return result;
}

std::string Globular::Client::getCaCertificate(std::string domain, unsigned int configurationPort){
    std::stringstream ss;
    ss << "http://" << domain << ":" << configurationPort << "/get_ca_certificate";
    http::Request request(ss.str());
    const http::Response response = request.send("GET");
    ss.flush();
    ss << std::string(response.body.begin(), response.body.end()) << '\n'; // print the result
    return ss.str();
}

std::string Globular::Client::signCaCertificate(std::string domain, unsigned int configurationPort, std::string csr){
    std::stringstream ss;
    std::string csr_str = macaron::Base64::Encode(csr);
    ss << "http://" << domain << ":" << configurationPort << "/sign_ca_certificate?=" + csr_str;
    http::Request request(ss.str());
    const http::Response response = request.send("GET");
    ss.flush();
    ss << std::string(response.body.begin(), response.body.end()) << '\n'; // print the result
    return ss.str();
}

// TODO fix to new certificate with multiple domains.
void Globular::Client::generateClientPrivateKey(std::string path, std::string pwd){
    std::stringstream ss;
    ss << path <<   "/client.key";
    auto path_ = ss.str();
    if(exists(path)){
        return;
    }

    ss.flush();

    ss << "openssl.exe";
    ss<< " genrsa";
    ss << " -passout";
    ss << " pass:"  << pwd;
    ss << " -des3";
    ss << " -out ";
    ss << " " << path << "/client.pass.key";
    ss << " 4096";
    std::system(ss.str().c_str());

    ss.flush();
    ss << "openssl.exe";
    ss<< " genrsa";
    ss << " -passin";
    ss << " pass:"  << pwd;
    ss << " -in";
    ss << " " << path << "/client.pass.key";
    ss << " -out ";
    ss << " " << path << "/client.key";
    std::system(ss.str().c_str());

    ss.flush();
    ss << path << "/client.pass.key";
    remove(ss.str().c_str());

}

void Globular::Client::generateClientCertificateSigningRequest(std::string path, std::string domain){
    std::stringstream ss;
    ss << path <<   "/client.csr";
    auto path_ = ss.str();
    if(exists(path)){
        return;
    }

    ss.flush();

    ss << "openssl.exe";
    ss<< " req";
    ss << " -new";
    ss << " -key";
    ss << " " << path << "/client.key";
    ss << " -out ";
    ss << " " << path << "/client.csr";
    ss << " -subj";
    ss << " /CN=" << domain;

    // run the command...
    std::system(ss.str().c_str());
}

void Globular::Client::keyToPem(std::string name, std::string path, std::string pwd){
    std::stringstream ss;
    ss << path <<   "/" << name + ".pem";
    auto path_ = ss.str();
    if(exists(path)){
        return;
    }

    ss.flush();

    ss << "openssl.exe";
    ss<< " pkcs8";
    ss << " -topk8";
    ss << " -nocrypt";
    ss << " -passin";
    ss << " pass:"  << pwd;
    ss << " -in ";
    ss << " " << path << "/" << name << ".key";
    ss << " -out ";
    ss << " " << path << "/" << name << ".pem";

    // run the command...
    std::system(ss.str().c_str());
}
