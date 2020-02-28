using Xunit;
using Globular;
using Persistence;

namespace Globular {
    public class PersistenceClient_Test{
        private PersistenceClient client = new PersistenceClient("localhost", "persistence_server");

        // Test create connection and also ping the connection to see if it exist and ready...
        [Fact]
        public void TestCreateConnection(){
            Persistence.Connection connection = new Persistence.Connection();
            connection.Id = "mongo_db_test_connection";
            connection.Name = "TestMongoDB";
            connection.Host = "localhost";
            connection.Port = 27017;
            connection.Store = Persistence.StoreType.Mongo;
            connection.User = "sa";
            connection.Password = "adminadmin";
            connection.Timeout = 3000;
            connection.Options = "";

            this.client.CreateConnection(connection, true);
            
            var pong = this.client.Ping("mongo_db_test_connection");
            Assert.Equal("pong", pong);
        }


    }
}