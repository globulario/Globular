using Grpc.Core;

namespace Globular
{
    public class PersistenceClient : Client
    {
        private Persistence.PersistenceService.PersistenceServiceClient client;
 
        public PersistenceClient(string address, string name) : base(address, name)
        {
            // Here I will create grpc connection with the service...
            this.client = new Persistence.PersistenceService.PersistenceServiceClient(this.channel);
        }

        public void CreateConnection(Persistence.Connection connection, bool save){
            // Here I will create the new connection.
            Persistence.CreateConnectionRqst rqst = new Persistence.CreateConnectionRqst();
            rqst.Connection = connection;
            rqst.Save = save;

            // Create a new connection
            this.client.CreateConnection(rqst, this.GetClientContext());
        }

        public string Ping(string connectionId){
            // Here I will create the new connection.
            Persistence.PingConnectionRqst rqst = new Persistence.PingConnectionRqst();
            rqst.Id = connectionId;

            // Create a new connection
            var rsp = this.client.Ping(rqst, this.GetClientContext());
            return rsp.Result;
        }
    }
}
