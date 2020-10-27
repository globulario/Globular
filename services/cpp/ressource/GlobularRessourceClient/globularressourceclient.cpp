#include "globularressourceclient.h"
#include <iostream>

Globular::RessourceClient::RessourceClient(std::string name, std::string domain, unsigned int configurationPort):
    Globular::Client(name,domain, configurationPort),
    stub_(ressource::RessourceService::NewStub(this->channel))
{
    std::cout << "init the ressource client!" << std::endl;
}

std::string Globular::RessourceClient::authenticate(std::string user, std::string password){
   ressource::AuthenticateRqst rqst;
   rqst.set_name(user);
   rqst.set_password(password);

   ressource::AuthenticateRsp rsp;
   grpc::ClientContext ctx;
   this->getClientContext(ctx);
   Status status = this->stub_->Authenticate(&ctx, rqst, &rsp);

   // return the token.
   if(status.ok()){
    return rsp.token();
   }else{
       std::cout << "fail to autenticate user " << user  << std::endl;
       return "";
   }
}

bool Globular::RessourceClient::validateUserAccess(std::string token, std::string method){
    ressource::ValidateUserAccessRqst rqst;
    rqst.set_token(token);
    rqst.set_method(method);
    grpc::ClientContext ctx;
    this->getClientContext(ctx);
    ressource::ValidateUserAccessRsp rsp;

    Status status = this->stub_->ValidateUserAccess(&ctx, rqst, &rsp);
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

    grpc::ClientContext ctx;
    this->getClientContext(ctx);
    ressource::ValidateApplicationAccessRsp rsp;
    Status status = this->stub_->ValidateApplicationAccess(&ctx, rqst, &rsp);
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

    grpc::ClientContext ctx;
    this->getClientContext(ctx);
    ressource::ValidateApplicationRessourceAccessRsp rsp;
    Status status = this->stub_->ValidateApplicationRessourceAccess(&ctx, rqst, &rsp);
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

    grpc::ClientContext ctx;
    this->getClientContext(ctx);
    ressource::ValidateUserRessourceAccessRsp rsp;
    Status status = this->stub_->ValidateUserRessourceAccess(&ctx, rqst, &rsp);
    if(status.ok()){
        return rsp.result();
    }else{
        return false;
    }
}

void  Globular::RessourceClient::SetRessource(std::string path, std::string name, int modified, int size){
    ressource::SetRessourceRqst rqst;
    ressource::Ressource* r = rqst.mutable_ressource();
    r->set_path(path);
    r->set_name(name);
    r->set_modified(modified);
    r->set_size(size);

    grpc::ClientContext ctx;
    this->getClientContext(ctx);
    ressource::SetRessourceRsp rsp;
    Status status = this->stub_->SetRessource(&ctx, rqst, &rsp);
    if(!status.ok()){
        std::cout << "Fail to set ressource " << name << std::endl;
    }
}

void  Globular::RessourceClient::removeRessouce(std::string path, std::string name){
    ressource::RemoveRessourceRqst rqst;
    ressource::Ressource* r = rqst.mutable_ressource();
    r->set_path(path);
    r->set_name(name);

    grpc::ClientContext ctx;
    this->getClientContext(ctx);
    ressource::RemoveRessourceRsp rsp;
    Status status = this->stub_->RemoveRessource(&ctx, rqst, &rsp);
    if(!status.ok()){
        std::cout << "Fail to remove ressource " << name << std::endl;
    }
}

std::vector<::ressource::ActionParameterRessourcePermission>  Globular::RessourceClient::getActionPermission(std::string method){
    ressource::GetActionPermissionRqst rqst;
    rqst.set_action(method);

    grpc::ClientContext ctx;
    this->getClientContext(ctx);
    ressource::GetActionPermissionRsp rsp;
    Status status = this->stub_->GetActionPermission(&ctx, rqst, &rsp);
    std::vector<::ressource::ActionParameterRessourcePermission>  results;
    if(status.ok()){
        for(auto i=0; i < rsp.actionparameterressourcepermissions().size(); i++){
            results.push_back(rsp.actionparameterressourcepermissions()[i]);
        }
    }
    return results;
}

void  Globular::RessourceClient::Log(std::string application, std::string method, std::string message, int type){

    ressource::LogRqst rqst;
    ressource::LogInfo* info = rqst.mutable_info();
    info->set_type(ressource::LogType(type));
    info->set_message(message);
    info->set_application(application);
    info->set_method(method);
    info->set_date(std::time(0));
    grpc::ClientContext ctx;
    this->getClientContext(ctx);
    ressource::LogRsp rsp;
    Status status = this->stub_->Log(&ctx, rqst, &rsp);
    if(status.ok()){
        return;
    }else{
        std::cout << "Fail to log information " << application << ":" << method << std::endl;
    }
}
