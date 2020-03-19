
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

        public EchoServiceImpl()
        {
            // Here I will set the default values.
            this.Port = 10029; // The default port value
            this.Proxy = 10030; // The reverse proxy port
            this.Name = "echo_server"; // The service name
            this.Version = "0.0.1";
            this.PublisherId = "localhost";
            this.Domain = "localhost";
            this.Protocol = "grpc";
            this.Version = "0.0.1";
            this.Value = "echo value!";
        }
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

