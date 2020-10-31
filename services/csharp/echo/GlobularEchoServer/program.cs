using System.Collections.Generic;
using System.Security.Principal;
using System.Text.RegularExpressions;
using System;
using Grpc.Core;
using Grpc.Core.Interceptors;
using System.Threading;
using System.Threading.Tasks;
using System.IO;

namespace Echo
{
    public class Prorgam
    {
        private static readonly AutoResetEvent _closing = new AutoResetEvent(false);
        private static Server server;

        public static void Main(string[] args)
        {
            Task.Factory.StartNew(() =>
            {
                // Create a new echo server instance.
                var echoServer = new EchoServiceImpl();
                // init values from the configuration file.
                echoServer = echoServer.init();
                if (echoServer.TLS == true)
                {
                    // Read ssl certificate and initialyse credential with it.
                    var cacert = File.ReadAllText(echoServer.CertAuthorityTrust);
                    var servercert = File.ReadAllText(echoServer.CertFile);
                    var serverkey = File.ReadAllText(echoServer.KeyFile);
                    var keypair = new KeyCertificatePair(servercert, serverkey);
                    // secure connection parameters.
                    var ssl = new SslServerCredentials(new List<KeyCertificatePair>() { keypair }, cacert, false);
                    // create the server.
                    server = new Server
                    {
                        Services = { EchoService.BindService(echoServer).Intercept(echoServer.interceptor) },
                        Ports = { new ServerPort(echoServer.Domain, echoServer.Port, ssl) }
                    };
                }
                else
                {
                    // non secure server.
                    server = new Server
                    {
                        Services = { EchoService.BindService(echoServer).Intercept(echoServer.interceptor) },
                        Ports = { new ServerPort(echoServer.Domain, echoServer.Port, ServerCredentials.Insecure) }
                    };
                }

                Console.WriteLine("Echo server listening on port " + echoServer.Port);

                // GRPC server.
                server.Start();
            });
            Console.CancelKeyPress += new ConsoleCancelEventHandler(OnExit);
            _closing.WaitOne();
        }

        protected static void OnExit(object sender, ConsoleCancelEventArgs args)
        {
            server.ShutdownAsync().Wait();
            _closing.Set();
        }
    }
}