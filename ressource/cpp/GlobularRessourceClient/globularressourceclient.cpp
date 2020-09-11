#include "globularressourceclient.h"
#include <iostream>

Globular::RessourceClient::RessourceClient(std::string name, std::string domain, unsigned int configurationPort):
    Globular::Client(name,domain, configurationPort),
    stub_(ressource::RessourceService::NewStub(channel))
{

}

std::string Globular::RessourceClient::authenticate(std::string user, std::string password){
   ressource::AuthenticateRqst rqst;
   rqst.set_name(user);
   rqst.set_password(password);

   ressource::AuthenticateRsp rsp;
   auto ctx = &this->getClientContext();

   Status status = this->stub_->Authenticate(ctx, rqst, &rsp);

   // return the token.
   if(status.ok()){
    return rsp.token();
   }else{
       std::cout << "fail to autenticate user " << user;
       return "";
   }
}

bool Globular::RessourceClient::validateUserAccess(std::string token, std::string method){
    ressource::ValidateUserAccessRqst rqst;
    rqst.set_token(token);
    rqst.set_method(method);
    auto ctx = &this->getClientContext();
    ressource::ValidateUserAccessRsp rsp;

    Status status = this->stub_->ValidateUserAccess(ctx, rqst, &rsp);
    if(status.ok()){
        return rsp.result();
    }else{
        return false;
    }
}

bool Globular::RessourceClient::validateApplicationAccess(std::string name, std::string method){
    ressource::ValidateApplicationAccessRqst rqst;
    rqst.set_name(name);
    rqst.set_method(method);
    auto ctx = &this->getClientContext();
    ressource::ValidateApplicationAccessRsp rsp;
    Status status = this->stub_->ValidateApplicationAccess(ctx, rqst, &rsp);
    if(status.ok()){
        return rsp.result();
    }else{
        return false;
    }
}

bool  Globular::RessourceClient::validateApplicationRessourceAccess(std::string application, std::string path, std::string method, int32_t permission){
    ressource::ValidateApplicationRessourceAccessRqst rqst;
    rqst.set_name(application);
    rqst.set_method(method);
    rqst.set_path(path);
    rqst.set_permission(permission);

    auto ctx = &this->getClientContext();
    ressource::ValidateApplicationRessourceAccessRsp rsp;
    Status status = this->stub_->ValidateApplicationRessourceAccess(ctx, rqst, &rsp);
    if(status.ok()){
        return rsp.result();
    }else{
        return false;
    }
}

bool  Globular::RessourceClient::validateUserRessourceAccess(std::string token, std::string path, std::string method, int32_t permission){
    ressource::ValidateUserRessourceAccessRqst rqst;
    rqst.set_token(token);
    rqst.set_method(method);
    rqst.set_path(path);
    rqst.set_permission(permission);

    auto ctx = &this->getClientContext();
    ressource::ValidateUserRessourceAccessRsp rsp;
    Status status = this->stub_->ValidateUserRessourceAccess(ctx, rqst, &rsp);
    if(status.ok()){
        return rsp.result();
    }else{
        return false;
    }
}

void  Globular::RessourceClient::SetRessource(std::string path, std::string name, int modified, int size){

}

void  Globular::RessourceClient::removeRessouce(std::string path, std::string name){

}

int  Globular::RessourceClient::getActionPermission(std::string method){

}

void  Globular::RessourceClient::Log(std::string application, std::string token, std::string method, std::string message, int type){

}
