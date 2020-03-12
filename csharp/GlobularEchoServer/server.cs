
using Globular;
using grpc = global::Grpc.Core;
using System.Threading.Tasks;

// The first thing to do is derived the service base class with GlobularService class.
// 
namespace Echo {
  public static partial class EchoService
  {
     public abstract partial class EchoServiceBase : GlobularService {
         
     }
  }

  public class EchoServiceImpl : EchoService.EchoServiceBase {
      public EchoServiceImpl(){
          
      }
      public override Task<global::Echo.EchoResponse> Echo(global::Echo.EchoRequest request, grpc::ServerCallContext context)
      {
        Echo.EchoResponse rsp = new EchoResponse();
        rsp.Message = "echo " + request.Message;
        return Task.FromResult(rsp);
      }

  }
}

