using System;
using Grpc.Core;

namespace Globular
{
    public class RessourceClient : Client
    {
        private Ressource.RessourceService.RessourceServiceClient client;

        /// <summary>
        /// The ressource client is use by the interceptor to validate user access.
        /// </summary>
        /// <param name="id"></param> The name or the id of the services.
        /// <param name="domain"></param> The domain of the services
        /// <param name="configurationPort"></param> The domain of the services
        /// <returns></returns>
        public RessourceClient( string id, string domain, int configurationPort) : base(id, domain, configurationPort)
        {
            // Here I will create grpc connection with the service...
            this.client = new Ressource.RessourceService.RessourceServiceClient(this.channel);
        }

        public string Authenticate(string user, string password){
            Ressource.AuthenticateRqst rqst = new Ressource.AuthenticateRqst();
            rqst.Name = user;
            rqst.Password = password;
            var rsp = this.client.Authenticate(rqst, this.GetClientContext());
            return rsp.Token;
        }

        /// <summary>
        /// Validate if the user can access a given method.
        /// </summary>
        /// <param name="token">The user token</param>
        /// <param name="method">The method </param>
        /// <returns></returns>
        public bool ValidateUserAccess(string token, string method)
        {
            Ressource.ValidateUserAccessRqst rqst = new Ressource.ValidateUserAccessRqst();
            rqst.Token = token;
            rqst.Method = method;
            var rsp = this.client.ValidateUserAccess(rqst, this.GetClientContext());
            return rsp.Result;
        }

        /// <summary>
        /// Validate if an application have access a given method.
        /// </summary>
        /// <param name="token"></param>
        /// <param name="method"></param>
        /// <returns></returns>
        public bool ValidateApplicationAccess(string name, string method)
        {
            Ressource.ValidateApplicationAccessRqst rqst = new Ressource.ValidateApplicationAccessRqst();
            rqst.Name = name;
            rqst.Method = method;
            var rsp = this.client.ValidateApplicationAccess(rqst, this.GetClientContext());
            return rsp.Result;
        }

        /// <summary>
        /// Validate if the user can access a given method.
        /// </summary>
        /// <param name="token">The user token</param>
        /// <param name="method">The method </param>
        /// <returns></returns>
        public bool ValidateUserRessourceAccess(string token, string path, string method, int permission)
        {
            Ressource.ValidateUserRessourceAccessRqst rqst = new Ressource.ValidateUserRessourceAccessRqst();
            rqst.Token = token;
            rqst.Method = method;
            rqst.Path = path; // the path of the ressource... 
            rqst.Permission = permission;

            var rsp = this.client.ValidateUserRessourceAccess(rqst, this.GetClientContext());
            return rsp.Result;
        }

        /// <summary>
        /// Validate if an application have access a given method.
        /// </summary>
        /// <param name="name"></param>
        /// <param name="method"></param>
        /// <param name="permission"></param>
        /// <returns></returns>
        public bool ValidateApplicationRessourceAccess (string name, string path, string method, int permission)
        {
            Ressource.ValidateApplicationRessourceAccessRqst rqst = new Ressource.ValidateApplicationRessourceAccessRqst();
            rqst.Name = name;
            rqst.Method = method;
            rqst.Path = path;
            rqst.Permission = permission;

            var rsp = this.client.ValidateApplicationRessourceAccess(rqst, this.GetClientContext());
            return rsp.Result;
        }

        /// <summary>
        /// Set a ressource path.
        /// </summary>
        /// <param name="path">The path of the ressource in form /toto/titi/tata</param>
        public void SetRessource(string path, string name, int modified, int size){
            Ressource.SetRessourceRqst rqst = new Ressource.SetRessourceRqst();
            Ressource.Ressource ressource = new Ressource.Ressource();
            ressource.Path = path;
            ressource.Name = name;
            ressource.Modified = modified;
            ressource.Size = size;
            rqst.Ressource = ressource;
            this.client.SetRessource(rqst);
        }

        /// <summary>
        /// Remove a ressource from globular. It also remove asscociated permissions.
        /// </summary>
        /// <param name="path"></param>
        public void RemoveRessouce(string path, string name){
            Ressource.RemoveRessourceRqst rqst = new Ressource.RemoveRessourceRqst();
            Ressource.Ressource ressource = new Ressource.Ressource();
            ressource.Path = path;
            ressource.Name = name;
            rqst.Ressource = ressource;
            this.client.RemoveRessource(rqst);
        }

        /// <summary>
        /// Get the ressource Action permission for a given ressource.
        /// </summary>
        /// <param name="path">The ressource path</param>
        /// <param name="action">The gRPC action</param>
        /// <returns></returns>
        public Int32 GetActionPermission(string action) {
            Ressource.GetActionPermissionRqst rqst = new Ressource.GetActionPermissionRqst();
            rqst.Action = action;
            var rsp = this.client.GetActionPermission(rqst);
            return rsp.Permission;
        }

        /// <summary>
        /// That method id use to log information/error
        /// </summary>
        /// <param name="application">The name of the application (given in the context)</param>
        /// <param name="token">Ths user token (logged end user)</param>
        /// <param name="method">The method called</param>
        /// <param name="message">The message info</param>
        /// <param name="type">Information or Error</param>
        public void Log(string application, string token, string method, string message, int type = 0)
        {
            var rqst = new Ressource.LogRqst();
            var info = new Ressource.LogInfo();
            info.Application = application;
            info.UserId = token; // can be a token or the user id...
            info.Method = method;
            if(type == 0){
                 info.Type = Ressource.LogType.InfoMessage;
            }else{
                info.Type = Ressource.LogType.ErrorMessage;
            }
            info.Message = message;
            rqst.Info = info;

            // Set the log.
            this.client.Log(rqst, this.GetClientContext());
        }
    }
}
