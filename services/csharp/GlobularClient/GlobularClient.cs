using System.Net.Http;
using System.Text.Json;
using System.IO;
using System.Threading.Tasks;
using Grpc.Core;
using System.Collections.Generic;
using System.Diagnostics;
using System.Text;

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
        public int ConfigurationPort { get; set; }
        public int PortHttp { get; set; }
        public int PortHttps { get; set; }
        public int AdminPort { get; set; }
        public string AdminEmail { get; set; }
        public int RessourcePort { get; set; }
        public int ServicesDiscoveryPort { get; set; }
        public int ServicesRepositoryPort { get; set; }
        public int CertificateAuthorityPort { get; set; }
        public int LoadBalancingServicePort { get; set; }
        public int SessionTimeout { get; set; }
        public int CertExpirationDelay { get; set; }
        public int IdleTimeout { get; set; }

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
        public string Id { get; set; }
        public string Path { get; set; }
        public string Proto { get; set; }
        public int Port { get; set; }
        public bool TLS { get; set; }
    }

    public class Client
    {
    	private string id;
        private string name;
        private string domain;
        private int port;
        private bool hasTls;
        private string caFile;
        private string keyFile;
        private string certFile;
        
        protected Channel channel;

        // Get Domain return the client domain.
        public string GetDomain()
        {
            return this.domain;
        }
        
        public string GetId()
        {
            return this.id;
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

        // At firt the port contain the http(s) domain of the globular server.
        // The configuration will be get from that domain and the port will
        // be set back to the correct domain.
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

        // Get the CA public certificate.
        private string getCaCertificate(string domain, int ConfigurationPort)
        {
            // Get the configuration from the globular server.
            var client = new HttpClient();
            string rqst = "http://" + domain + ":" + ConfigurationPort + "/get_ca_certificate";
            var task = Task.Run(() => client.GetAsync(rqst));
            task.Wait();
            var rsp = task.Result;
            if (rsp.IsSuccessStatusCode == false)
            {
                throw new System.InvalidOperationException("Fail to get client configuration " + rqst);
            }

            return rsp.Content.ReadAsStringAsync().Result;
        }

        public static string Base64Encode(string plainText)
        {
            var plainTextBytes = System.Text.Encoding.UTF8.GetBytes(plainText);
            return System.Convert.ToBase64String(plainTextBytes);
        }

        // ask globular CA to sign the cerificate.
        private string signCaCertificate(string domain, int ConfigurationPort, string csr)
        {
            var client = new HttpClient();

            string csr_str = Base64Encode(csr);
            string rqst = "http://" + domain + ":" + ConfigurationPort + "/sign_ca_certificate?csr=" + csr_str;
            var task = Task.Run(() => client.GetAsync(rqst));
            task.Wait();
            var rsp = task.Result;
            if (rsp.IsSuccessStatusCode == false)
            {
                throw new System.InvalidOperationException("Fail to get sign ca certificate!");
            }

            return rsp.Content.ReadAsStringAsync().Result;
        }

        /**
         * I will made use of openssl as external command to be able to generate key and
         * certificate the same way in every language.
         */
        private void generateClientPrivateKey(string path, string pwd)
        {
            if (File.Exists(path + "/client.key"))
            {
                return;
            }

            Process process_0 = new Process();
            process_0.StartInfo.FileName = "openssl.exe";

            // Set args
            process_0.StartInfo.ArgumentList.Add("genrsa");
            process_0.StartInfo.ArgumentList.Add("-passout");
            process_0.StartInfo.ArgumentList.Add("pass:" + pwd);
            process_0.StartInfo.ArgumentList.Add("-des3");
            process_0.StartInfo.ArgumentList.Add("-out");
            process_0.StartInfo.ArgumentList.Add(path + "/client.pass.key");
            process_0.StartInfo.ArgumentList.Add("4096");

            // set options
            process_0.StartInfo.UseShellExecute = false;
            process_0.StartInfo.RedirectStandardOutput = true;
            process_0.StartInfo.RedirectStandardError = true;

            process_0.Start();
            process_0.WaitForExit();

            Process process_1 = new Process();
            process_1.StartInfo.FileName = "openssl.exe";

            // Set args
            process_1.StartInfo.ArgumentList.Add("rsa");
            process_1.StartInfo.ArgumentList.Add("-passin");
            process_1.StartInfo.ArgumentList.Add("pass:" + pwd);
            process_1.StartInfo.ArgumentList.Add("-in");
            process_1.StartInfo.ArgumentList.Add(path + "/client.pass.key");
            process_1.StartInfo.ArgumentList.Add("-out");
            process_1.StartInfo.ArgumentList.Add(path + "/client.key");

            // set options
            process_1.StartInfo.UseShellExecute = false;
            process_1.StartInfo.RedirectStandardOutput = true;
            process_1.StartInfo.RedirectStandardError = true;
            process_1.Start();
            process_1.WaitForExit();

            // remove the intermediary file.
            File.Delete(path + "/client.pass.key");
        }

        /**
         * Generate a client signing request for a given domain.
         */
        private void generateClientCertificateSigningRequest(string path, string domain)
        {
            if (File.Exists(path + "/client.csr"))
            {
                return;
            }

            Process process_0 = new Process();
            process_0.StartInfo.FileName = "openssl.exe";

            // Set args
            process_0.StartInfo.ArgumentList.Add("req");
            process_0.StartInfo.ArgumentList.Add("-new");
            process_0.StartInfo.ArgumentList.Add("-key");
            process_0.StartInfo.ArgumentList.Add(path + "/client.key");
            process_0.StartInfo.ArgumentList.Add("-out");
            process_0.StartInfo.ArgumentList.Add(path + "/client.csr");
            process_0.StartInfo.ArgumentList.Add("-subj");
            process_0.StartInfo.ArgumentList.Add("/CN=" + domain);

            // set options
            process_0.StartInfo.UseShellExecute = false;
            process_0.StartInfo.RedirectStandardOutput = true;
            process_0.StartInfo.RedirectStandardError = true;

            process_0.Start();
            process_0.WaitForExit();
        }

        private void keyToPem(string name, string path, string pwd)
        {
            if (File.Exists(path + "/" + name + ".pem"))
            {
                return;
            }

            Process process_0 = new Process();
            process_0.StartInfo.FileName = "openssl.exe";

            // Set args
            process_0.StartInfo.ArgumentList.Add("pkcs8");
            process_0.StartInfo.ArgumentList.Add("-topk8");
            process_0.StartInfo.ArgumentList.Add("-nocrypt");
            process_0.StartInfo.ArgumentList.Add("-passin");
            process_0.StartInfo.ArgumentList.Add("pass:" + pwd);
            process_0.StartInfo.ArgumentList.Add("-in");
            process_0.StartInfo.ArgumentList.Add(path + "/" + name + ".key");
            process_0.StartInfo.ArgumentList.Add("-out");
            process_0.StartInfo.ArgumentList.Add(path + "/" + name + ".pem");

            // set options
            process_0.StartInfo.UseShellExecute = false;
            process_0.StartInfo.RedirectStandardOutput = true;
            process_0.StartInfo.RedirectStandardError = true;

            process_0.Start();
            process_0.WaitForExit();
        }

        private static bool VerifyPeer(VerifyPeerContext context)
        {
            File.WriteAllText("c:/temp/toto.txt", "VerifiPeer!");
            return true;
        }

        private void init(string id , string domain="localhost",  int configurationPort=10000)
        {
            // Get the configuration from the globular server.
            var client = new HttpClient();
            string rqst = "http://" + domain + ":" + configurationPort +"/config";
            var task = Task.Run(() => client.GetAsync(rqst));
            task.Wait();
            var rsp = task.Result;
            if (rsp.IsSuccessStatusCode == false)
            {
                throw new System.InvalidOperationException("Fail to get client configuration " + rqst);
            }

            // Here I will parse the JSON object and initialyse values from it...
            var serverConfig = JsonSerializer.Deserialize<ServerConfig>(rsp.Content.ReadAsStringAsync().Result);
            ServiceConfig config = null;
            if (!serverConfig.Services.ContainsKey(id))
            {
                foreach(var s in serverConfig.Services.Values)
                {
                    if(s.Name == id){
                        config = s;
                        break;
                    }
                }
                if(config == null){
                    throw new System.InvalidOperationException("No serivce found with id " + id + "!");
                }
            }else{
                config = serverConfig.Services[id];
            }

            // get the service config.
            this.port = config.Port;
            this.hasTls = config.TLS;
            this.domain = config.Domain;
            this.id = config.Id;
            this.name = config.Name;

            // Here I will create grpc connection with the service...
            if (!this.HasTLS())
            {
                // Non secure connection.
                this.channel = new Channel(this.domain, this.port, ChannelCredentials.Insecure);
            }
            else
            {
                // if the client is not local I will generate TLS certificates.
                if (File.Exists(Path.GetTempPath() + "/" + this.domain + "_token"))
                {
                    // The ca certificate.
                    this.caFile = config.CertAuthorityTrust;

                    // get the client certificate and key here.
                    this.certFile = config.CertFile.Replace("server", "client");
                    this.keyFile = config.KeyFile.Replace("server", "client");
                }
                else
                {
                    // I will need to create certificate and make it sign by the CA.
                    var path = Path.GetTempPath() + "/config/tls/" + this.domain;

                    if (!Directory.Exists(path))
                    {
                        Directory.CreateDirectory(path);
                    }

                    // Now I will create the certificates.
                    var ca_crt = getCaCertificate(this.domain,  configurationPort);
                    File.WriteAllText(path + "/ca.crt", ca_crt);

                    var pwd = "1111"; // Set in the configuration...

                    // Now I will generate the certificate for the client...
                    // Step 1: Generate client private key.
                    this.generateClientPrivateKey(path, pwd);

                    // Step 2: Generate the client signing request.
                    this.generateClientCertificateSigningRequest(path, this.domain);

                    // Step 3: Generate client signed certificate.
                    var client_csr = File.ReadAllText(path + "/client.csr");
                    var client_crt = this.signCaCertificate(this.domain, configurationPort, client_csr);
                    File.WriteAllText(path + "/client.crt", client_crt);

                    // Step 4: Convert client.key to pem file.
                    this.keyToPem("client", path, pwd);

                    // Set path in the config.
                    this.keyFile = path + "/client.key";
                    this.caFile = path + "/ca.crt";
                    this.certFile = path + "/client.crt";
                }

                var cacert = File.ReadAllText(this.caFile);
                var clientcert = File.ReadAllText(this.certFile);
                var clientkey = File.ReadAllText(this.keyFile);
                var ssl = new SslCredentials(cacert, new KeyCertificatePair(clientcert, clientkey), VerifyPeer);

                //File.WriteAllText("c:/temp/toto.txt", clientkey);
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
                var path_ = Path.GetTempPath() + "/" + this.domain + "_token";
                if (File.Exists(path_))
                {
                    token = File.ReadAllText(path_);
                    metadata.Add("token", token);
                }
            }
            else
            {
                metadata.Add("token", token);
            }

            // set the local domain.
            if (domain.Length == 0)
            {
                metadata.Add("domain", this.domain);
            }
            else
            {
                metadata.Add("domain", domain);
            }

            if (application.Length > 0)
            {
                metadata.Add("application", application);
            }

            // The path of ressource if there one.
            if (path.Length > 0)
            {
                metadata.Add("path", path);
            }

            return metadata;
        }

        public Client(string id, string domain, int configurationPort)
        {
            this.id = id;
            this.domain = domain;
        
            // Now I will get the client configuration.
            this.init(id, domain, configurationPort);
        }
    }
}
