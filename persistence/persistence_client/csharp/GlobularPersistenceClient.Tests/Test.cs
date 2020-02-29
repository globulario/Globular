using Xunit;
using System;

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

        // Test find all...
        [Fact]
        public void TestFind(){
            string jsonStr = this.client.Find("mongo_db_test_connection", "local_ressource", "Accounts", "{}", "");
            Assert.True(jsonStr.Length > 0);
        }

        [Fact]
        public void TestFindOne(){
            string jsonStr = this.client.FindOne("mongo_db_test_connection", "local_ressource", "Accounts", "{\"_id\":\"dave\"}", "");
            Assert.True(jsonStr.Length > 0);
        }


    }
}