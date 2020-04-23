using System.Net.Http;
using System.Text.Json;
using System.IO;
using System.Threading.Tasks;
using Grpc.Core;
using System.Collections.Generic;

namespace Globular
{
    /** Globular server config. **/
    public class ServerConfig
    {
        public string Domain { get; set; }
        public string Name { get; set; }
        public string Protocol { get; set; }
        public string CertStableURL { get; set; }
        public string CertURL { get; set; }
        public uint PortHttp { get; set; }
        public uint PortHttps { get; set; }
        public uint AdminPort { get; set; }
        public uint AdminProxy { get; set; }
        public string AdminEmail { get; set; }
        public uint RessourcePort { get; set; }
        public uint RessourceProxy { get; set; }
        public uint ServicesDiscoveryPort { get; set; }
        public uint ServicesDiscoveryProxy { get; set; }
        public uint ServicesRepositoryPort { get; set; }
        public uint ServicesRepositoryProxy { get; set; }
        public uint CertificateAuthorityPort { get; set; }
        public uint CertificateAuthorityProxy { get; set; }
        public uint SessionTimeout { get; set; }
        public uint CertExpirationDelay { get; set; }
        public uint IdleTimeout { get; set; }

        public string[] Discoveries { get; set; }
        public string[] DNS { get; set; }

        // The map of service object.
        public Dictionary<string, ServiceConfig> Services { get; set; }

    }

    /**
     * Used by JSON serialysation.
     */
    public class ServiceConfig
    {
        public string CertAuthorityTrust { get; set; }
        public string CertFile { get; set; }
        public string KeyFile { get; set; }
        public string Domain { get; set; }
        public string Name { get; set; }
        public int Port { get; set; }
        public bool TLS { get; set; }
    }

    public class Client
    {
        private string name;
        private string address;
        private string domain;
        private int port;
        private bool hasTls;
        private string caFile;
        private string keyFile;
        private string certFile;
        protected Channel channel;

        // Return the ipv4 address
        public string GetAddress()
        {
            return this.address;
        }

        // Get Domain return the client domain.
        public string GetDomain()
        {
            return this.domain;
        }

        public string GetName()
        {
            return this.name;
        }

        public int GetPort()
        {
            return this.port;
        }

        // Close the client.
        public void Close()
        {
            // close the connection channel.
            this.channel.ShutdownAsync().Wait();
        }

        // At firt the port contain the http(s) address of the globular server.
        // The configuration will be get from that address and the port will
        // be set back to the correct address.
        public void SetPort(int port)
        {
            this.port = port;
        }

        // Set the name of the client
        public void SetName(string name)
        {
            this.name = name;
        }

        // Set the domain of the client
        public void SetDomain(string domain)
        {
            this.domain = domain;
        }

        ////////////////// TLS ///////////////////

        /// <summary>
        /// Test if the server is secure with TLS.
        /// </summary>
        /// <returns>True if it secure.</returns>
        public bool HasTLS()
        {
            return this.hasTls;
        }

        // Get the TLS certificate file path
        public string GetCertFile()
        {
            return this.certFile;
        }

        // Get the TLS key file path
        public string GetKeyFile()
        {
            return this.keyFile;
        }

        // Get the TLS key file path
        public string GetCaFile()
        {
            return this.certFile;
        }

        // Set the client is a secure client.
        public void SetTLS(bool hasTls)
        {
            this.hasTls = hasTls;
        }

        // Set TLS certificate file path
        public void SetCertFile(string certFile)
        {
            this.certFile = certFile;
        }

        // Set TLS key file path
        public void SetKeyFile(string keyFile)
        {
            this.keyFile = keyFile;
        }

        // Set TLS authority trust certificate file path
        public void SetCaFile(string caFile)
        {
            this.caFile = caFile;
        }
        private void init(string address, string name)
        {
            // Get the configuration from the globular server.
            var client = new HttpClient();
            string rqst = "http://" + address + ":10000/config";
            var task = Task.Run(() => client.GetAsync(rqst));
            task.Wait();
            var rsp = task.Result;
            if (rsp.IsSuccessStatusCode == false)
            {
                throw new System.InvalidOperationException("Fail to get client configuration " + rqst);
            }

            // Here I will parse the JSON object and initialyse values from it...
            var serverConfig = JsonSerializer.Deserialize<ServerConfig>(rsp.Content.ReadAsStringAsync().Result);
            if (!serverConfig.Services.ContainsKey(name))
            {
                throw new System.InvalidOperationException("No serivce found with name " + name + "!");
            }

            // get the service config.
            var config = serverConfig.Services[name];
            this.port = config.Port;
            this.hasTls = config.TLS;
            this.domain = config.Domain;
            this.caFile = config.CertAuthorityTrust;
            this.certFile = config.CertFile;
            this.keyFile = config.KeyFile;

            // Here I will create grpc connection with the service...
            if (!this.HasTLS())
            {
                // Non secure connection.
                this.channel = new Channel(this.GetDomain() + ":" + this.GetPort(), ChannelCredentials.Insecure);
            }
            else
            {
                // TODO test if the service is local. 


                // Secure connection.
                var cacert = File.ReadAllText(this.caFile);
                var clientcert = File.ReadAllText(this.certFile);
                var clientkey = File.ReadAllText(this.keyFile);
                var ssl = new SslCredentials(cacert, new KeyCertificatePair(clientcert, clientkey));

                this.channel = new Channel(this.domain, this.port, ssl);
            }
        }

        protected Metadata GetClientContext(string token = "", string application = "", string domain = "", string path = "")
        {
            // Set the token in the metadata.
            var metadata = new Metadata();

            // Here I will get the token from the file.
            if (token.Length == 0)
            {
                var path_ = Path.GetTempPath() + Path.PathSeparator + this.domain + "_token";
                if (File.Exists(path_))
                {
                    token = File.ReadAllText(path_);
                    metadata.Add("token", token);
                }
            }else{
                 metadata.Add("token", token);
            }

            // set the local domain.
            if(domain.Length == 0){
                metadata.Add("domain", this.domain);
            }else{
                 metadata.Add("domain", domain);
            }

            if(application.Length>0){
                 metadata.Add("application", application);
            }

            // The path of ressource if there one.
            if(path.Length>0){
                 metadata.Add("path", path);
            }

            return metadata;
        }

        public Client(string address, string name)
        {
            this.name = name;
            this.address = address;

            // Here I will initialyse the configuration.
            if (address.IndexOf(':') > 0)
            {
                address = address.Split(':')[0];
            }

            // Now I will get the client configuration.
            this.init(address, name);
        }
    }
}
