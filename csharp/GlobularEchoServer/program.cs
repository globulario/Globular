using System;
using Grpc.Core;

namespace Echo
{
    class Program
    {
        static void Main(string[] args)
        {
            try
            {    
                // Create a new echo server instance.
                var echoServer = new EchoServiceImpl();
                echoServer.init();
                
                Server server = new Server
                {
                    Services = { EchoService.BindService(echoServer) },
                    Ports = { new ServerPort("localhost", echoServer.Port, ServerCredentials.Insecure) }
                };

                server.Start();

                Console.WriteLine("Echo server listening on port " + echoServer.Port);
                Console.WriteLine("Press any key to stop the server...");
                Console.ReadKey();

                server.ShutdownAsync().Wait(); // Close the server.
            }
            catch(Exception ex)
            {
                Console.WriteLine($"Exception encountered: {ex}");
            }
        }
    }
}
