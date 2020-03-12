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
        /// <param name="address"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public RessourceClient(string address, string name) : base(address, name)
        {
            // Here I will create grpc connection with the service...
            this.client = new Ressource.RessourceService.RessourceServiceClient(this.channel);
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
        public bool ValidateApplicationRessourceAccess(string token, string path, string method)
        {
            Ressource.ValidateUserRessourceAccessRqst rqst = new Ressource.ValidateUserRessourceAccessRqst();
            rqst.Token = token;
            rqst.Method = method;
            rqst.Path = path;

            var rsp = this.client.ValidateUserRessourceAccess(rqst, this.GetClientContext());
            return rsp.Result;
        }

        /// <summary>
        /// Validate if an application have access a given method.
        /// </summary>
        /// <param name="token"></param>
        /// <param name="method"></param>
        /// <returns></returns>
        public bool ValidateUserRessourceAccess(string name, string path, string method)
        {
            Ressource.ValidateApplicationRessourceAccessRqst rqst = new Ressource.ValidateApplicationRessourceAccessRqst();
            rqst.Name = name;
            rqst.Method = method;
            rqst.Path = path;

            var rsp = this.client.ValidateApplicationRessourceAccess(rqst, this.GetClientContext());
            return rsp.Result;
        }

        /// <summary>
        /// Set a ressource path.
        /// </summary>
        /// <param name="path">The path of the ressource in form /toto/titi/tata</param>
        public void SetRessource(string path){
            Ressource.SetRessourceRqst rqst = new Ressource.SetRessourceRqst();
            rqst.Ressource = path;
            this.client.setRessource(rqst);
        }

        public void SetRessouces(string[] paths){
 
            // append paths to the field.
            var call = this.client.setRessources();

            // call.RequestStream.
            Ressource.SetRessourcesRqst rqst = new Ressource.SetRessourcesRqst();;
            for(var i=0; i < paths.Length; i++){
                rqst.Ressources.Add(paths[i]);
                if(i % 1000 == 0 && i > 0){
                    call.RequestStream.WriteAsync(rqst);
                    // set a new request...
                    rqst = new Ressource.SetRessourcesRqst();
                }
            }

            if(rqst.Ressources.Count > 0){
                call.RequestStream.WriteAsync(rqst);
            }

            // Close the stream.
            call.RequestStream.CompleteAsync();
        }

        /// <summary>
        /// Remove a ressource from globular. It also remove asscociated permissions.
        /// </summary>
        /// <param name="path"></param>
        public void RemoveRessouce(string path){
            Ressource.RemoveRessourceRqst rqst = new Ressource.RemoveRessourceRqst();
            rqst.Ressource = path;
            this.client.removeRessource(rqst);
        }

        /// <summary>
        /// That method id use to log information/error
        /// </summary>
        /// <param name="application">The name of the application (given in the context)</param>
        /// <param name="userId">Ths user id (logged end user)</param>
        /// <param name="method">The method called</param>
        /// <param name="message">The message info</param>
        /// <param name="type">Information or Error</param>
        public void Log(string application, string userId, string method, string message, Ressource.LogType type)
        {
            var rqst = new Ressource.LogRqst();
            var info = new Ressource.LogInfo();
            info.Application = application;
            info.UserId = userId;
            info.Method = method;
            info.Type = type;
            info.Message = message;
            rqst.Info = info;

            // Set the log.
            this.client.Log(rqst, this.GetClientContext());
        }
    }
}
