using System;
using System.IO;
using System.Text.Json;
using Grpc.Core;
using Grpc.Core.Interceptors;
using System.Threading.Tasks;

namespace Globular
{

    /// <summary>
    /// That class contain the basic service class. Globular service are 
    /// plain gRPC service with required attributes to make it manageable.
    /// </summary>
    public class GlobularService
    {
        public string Name { get; set; }
        public int Port { get; set; }
        public int Proxy { get; set; }
        public string Protocol { get; set; }
        public bool AllowAllOrigins { get; set; }
        public string AllowedOrigins { get; set; }
        public string Domain { get; set; }
        public string CertAuthorityTrust { get; set; }
        public string CertFile { get; set; }
        public string KeyFile { get; set; }
        public bool TLS { get; set; }
        public string Version { get; set; }
        public string PublisherId { get; set; }
        public bool KeepUpToDate { get; set; }
        public bool KeepAlive { get; set; }

        private RessourceClient ressourceClient;
        public ServerUnaryInterceptor interceptor;

        /// <summary>
        /// The default constructor.
        /// </summary>
        public GlobularService(string address = "localhost")
        {
            // set default values.
            this.Domain = address;
            this.Protocol = "grpc";
            this.Version = "0.0.1";
            this.PublisherId = "localhost";

            // Create the interceptor.
            this.interceptor = new Globular.ServerUnaryInterceptor(this);
        }

        private RessourceClient getRessourceClient()
        {
            if (this.ressourceClient == null)
            {
                // there must be a globular server runing in order to validate ressources.
                ressourceClient = new RessourceClient(this.Domain, "Ressource");
            }
            return this.ressourceClient;
        }

        private string getPath()
        {
            string path = Directory.GetCurrentDirectory();
            Console.Write("----> path: " + path);
            
            return path;
        }

        /// <summary>
        /// Initialyse from json object from a file.
        /// </summary>
        public object init(object server)
        {
            var configPath = this.getPath() + Path.DirectorySeparatorChar + "config.json";
            // Here I will read the file that contain the object.
            if (File.Exists(configPath))
            {
                var jsonStr = File.ReadAllText(configPath);
                var s = JsonSerializer.Deserialize(jsonStr, server.GetType());
                return s;
            }
            this.save(server);
            return server;
        }

        /// <summary>
        /// Serialyse the object into json and save it in config.json file.
        /// </summary>
        public void save(object server)
        {
            var configPath = getPath() + Path.DirectorySeparatorChar + "config.json";
            string jsonStr;
            jsonStr = JsonSerializer.Serialize(server);
            File.WriteAllText(configPath, jsonStr);
        }

        /// <summary>
        /// Set a ressource on the globular ressource manager.
        /// </summary>
        /// <param name="path"></param>
        public void setRessource(string path)
        {
            this.getRessourceClient().SetRessource(path);
        }

        /// <summary>
        /// Set a batch of ressource on the globular ressource manager.
        /// </summary>
        /// <param name="paths"></param>
        public void setRessources(string[] paths)
        {
            this.getRessourceClient().SetRessouces(paths);
        }

        public bool validateUserAccess(string token, string method)
        {
            return this.getRessourceClient().ValidateApplicationAccess(token, method);
        }

        public bool validateApplicationAccess(string application, string method)
        {
            return this.getRessourceClient().ValidateApplicationAccess(application, method);
        }

        public void logInfo(string application, string token, string method, string message, int logType){
            this.getRessourceClient().Log(application, token, method, message, logType);
        }
    }

    public class ServerUnaryInterceptor : Interceptor
    {

        private GlobularService service;

        public ServerUnaryInterceptor(GlobularService srv)
        {
            this.service = srv;
        }

        public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(TRequest request, ServerCallContext context, UnaryServerMethod<TRequest, TResponse> continuation)
        {
            // Do method validations here.
            Metadata metadatas = context.RequestHeaders;
            string application = "";
            string token = "";
            string method = context.Method;
            bool hasAccess = false;

            // Get the metadata from the header.
            for (var i = 0; i < metadatas.Count; i++)
            {
                var metadata = metadatas[i];
                if (metadata.Key == "application")
                {
                    application = metadata.Value;
                    if(!hasAccess){
                        hasAccess = this.service.validateApplicationAccess(application, method);
                    }
                }
                else if (metadata.Key == "token")
                {
                    token = metadata.Value;
                    if(!hasAccess){
                        hasAccess = this.service.validateUserAccess(token, method);
                    }
                }
            }

            // Here I will validate the user for action.
            if(!hasAccess){
                // I will log the error in the log.
                this.service.logInfo(application, token, method, "Permission denied ", 1);
                // here I the user and the application has no access to the method 
                // I will throw an exception.
                throw new RpcException(new Status(StatusCode.PermissionDenied, "Permission denied"), metadatas);
            }

            // this.service.
            var response = await base.UnaryServerHandler(request, context, continuation);

            return response;
        }
    }
}