using System.Collections.Generic;
using System.Security.Principal;
using System.Text.RegularExpressions;
using System;
using Grpc.Core;
using Grpc.Core.Interceptors;
using System.Threading;
using System.Threading.Tasks;


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

                server = new Server
                {
                    Services = { EchoService.BindService(echoServer).Intercept(echoServer.interceptor) },
                    Ports = { new ServerPort("localhost", echoServer.Port, ServerCredentials.Insecure) }
                };

                Console.WriteLine("Echo server listening on port " + echoServer.Port);

                // GRPC server.
                server.Start();

                while (true)
                {
                    Thread.Sleep(1000);
                }
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