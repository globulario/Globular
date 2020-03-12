using System;

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

        /// <summary>
        /// The default constructor.
        /// </summary>
        public GlobularService(string address="localhost")
        {
            // set default values.
            this.Domain = "localhost";
            this.Protocol = "grpc";
            this.Version = "0.0.1";
            this.PublisherId = "localhost";

            // there must be a globular server runing in order to validate ressources.
            ressourceClient = new RessourceClient(address, "Ressource");
        }

        private string getPath()
        {
            var rootDir = System.IO.Path.GetDirectoryName(System.Reflection.Assembly.GetExecutingAssembly().CodeBase);
            return rootDir;
        }

        /// <summary>
        /// Initialyse from json object from a file.
        /// </summary>
        public void init()
        {
            Console.Write("init service from: ", getPath() + System.IO.Path.PathSeparator + "config.json");
        }

        /// <summary>
        /// Serialyse the object into json and save it in config.json file.
        /// </summary>
        public void save(){
            Console.Write("save service to: ", getPath() + System.IO.Path.PathSeparator + "config.json");
        }

        /// <summary>
        /// Set a ressource on the globular ressource manager.
        /// </summary>
        /// <param name="path"></param>
        public void setRessource(string path){
            this.ressourceClient.SetRessource(path);
        }

        /// <summary>
        /// Set a batch of ressource on the globular ressource manager.
        /// </summary>
        /// <param name="paths"></param>
        public void setRessources(string[] paths){
            this.ressourceClient.SetRessouces(paths);
        }

    }


}
