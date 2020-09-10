
using System;
using Globular;
using grpc = global::Grpc.Core;
using System.Threading.Tasks;
using System.Text.Json;

// The first thing to do is derived the service base class with GlobularService class.
namespace Echo
{
    public static partial class EchoService
    {
        public abstract partial class EchoServiceBase : GlobularService
        {

        }

    }

    public class EchoServiceImpl : EchoService.EchoServiceBase
    {
        public string Value { get; set; }

        public EchoServiceImpl(string id, string domain, uint port, uint proxy)
        {
            // Here I will set the default values.
            this.Port = port; // The default port value
            this.Proxy = proxy; // The reverse proxy port
            this.Name = "echo.EchoService"; // The service name
            this.Version = "0.0.1";
            this.PublisherId = "localhost"; // must be the publisher id here...
            this.Domain = domain;
            this.Protocol = "grpc";
            this.Version = "0.0.1";
            this.Value = "echo value!";
        }

        // Overide method of the service to implement in C#
        public override Task<global::Echo.EchoResponse> Echo(global::Echo.EchoRequest request, grpc::ServerCallContext context)
        {
            Echo.EchoResponse rsp = new EchoResponse();
            rsp.Message = "echo " + request.Message;
            return Task.FromResult(rsp);
        }

        // Here I will set the default config values...
        public EchoServiceImpl init()
        {
            // call save on init
            return (EchoServiceImpl)base.init(this);
        }
    }
}
