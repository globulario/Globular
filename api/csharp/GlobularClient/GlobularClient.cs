using System.Net.Http;
using System.Text.Json;
using System.IO;
using System.Threading.Tasks;
using Grpc.Core;

namespace Globular
{
    /**
     * Used by JSON serialysation.
     */
    public class Config
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
            this.channel.ShutdownAsync();
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

        //if the client is secure.
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
            string rqst = "http://localhost:10000/client_config?address=" + address + "&name=" + name;

            var task = Task.Run(() => client.GetAsync(rqst));
            task.Wait();
            var rsp = task.Result;
            if (rsp.IsSuccessStatusCode == false)
            {
                throw new System.InvalidOperationException("Fail to get client configuration " + rqst);
            }

            // Here I will parse the JSON object and initialyse values from it...
            var config = JsonSerializer.Deserialize<Config>(rsp.Content.ReadAsStringAsync().Result);
            this.port = 10035;//config.Port;
            this.hasTls = config.TLS;
            this.domain = config.Domain;
            this.caFile = config.CertAuthorityTrust;
            this.certFile = config.CertFile;
            this.keyFile = config.KeyFile;

            // Now I will open the channel with the server.

            // Here I will create grpc connection with the service...
            if (!this.HasTLS())
            {
                // Non secure connection.
                this.channel = new Channel(this.GetDomain() + ":" + this.GetPort(), ChannelCredentials.Insecure);
            }
            else
            {
                // Secure connection.
                // TODO create a secure connection here.
            }
        }

        protected Metadata GetClientContext()
        {
            // Here I will get the token from the file.
            var path = Path.GetTempPath() + Path.PathSeparator + this.domain + "_token";
            string token = "";
            // Set the token in the metadata.
            var metadata = new Metadata();

            if(File.Exists(path)){
                token = File.ReadAllText(path);
                metadata.Add("token", token);
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
