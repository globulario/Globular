using System;
using Grpc.Core;

namespace Echo
{
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Hello World!");
            try
            {    
                Server server = new Server
                {
                    Services = { AccountService.BindService(new AccountsImpl()) },
                    Ports = { new ServerPort("localhost", Port, ServerCredentials.Insecure) }
                };
                server.Start();
                Console.WriteLine("Accounts server listening on port " + Port);
                Console.WriteLine("Press any key to stop the server...");
                Console.ReadKey();
                server.ShutdownAsync().Wait();
            }
            catch(Exception ex)
            {
                Console.WriteLine($"Exception encountered: {ex}");
            }
        }
    }
}
