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
        public  bool ValidateUserAccess( string token, string method ){
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
        public  bool ValidateApplicationAccess( string name, string method ){
          Ressource.ValidateApplicationAccessRqst rqst = new Ressource.ValidateApplicationAccessRqst();
            rqst.Name = name;
            rqst.Method = method;
            var rsp = this.client.ValidateApplicationAccess(rqst, this.GetClientContext());
            return rsp.Result;
        }
    }
}
