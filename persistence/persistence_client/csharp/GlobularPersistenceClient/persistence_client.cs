using Grpc.Core;
using System;
using System.Threading.Tasks;
using System.Collections;
using System.Text.Json;

namespace Globular
{
    public class PersistenceClient : Client
    {
        private Persistence.PersistenceService.PersistenceServiceClient client;

        /// <summary>
        /// gRPC client for persistence service.
        /// </summary>
        /// <param name="address">Can be a domain or a IP address ex: localhos or 127.0.0.1</param>
        /// <param name="name">The name of the service on the server. ex: persistence_server</param>
        /// <returns>Return the instance of the client with it connection ready to be use.</returns>
        public PersistenceClient(string address, string name) : base(address, name)
        {
            // Here I will create grpc connection with the service...
            this.client = new Persistence.PersistenceService.PersistenceServiceClient(this.channel);
        }

        /// <summary>
        /// Create a new persistence connection
        /// </summary>
        /// <param name="connection">The connection information</param>
        /// <param name="save">If true the connection will be save in the configuation file.</param>
        public void CreateConnection(Persistence.Connection connection, bool save)
        {
            // Here I will create the new connection.
            Persistence.CreateConnectionRqst rqst = new Persistence.CreateConnectionRqst();
            rqst.Connection = connection;
            rqst.Save = save;

            // Create a new connection
            this.client.CreateConnection(rqst, this.GetClientContext());
        }

        /// <summary>
        /// Delete a connection with a given id.
        /// </summary>
        /// <param name="connectionId">The connection to delete</param>
        public void DeleteConnection(string connectionId)
        {
            var rqst = new Persistence.DeleteConnectionRqst();
            rqst.Id = connectionId;
            this.client.DeleteConnection(rqst, this.GetClientContext());
        }

        /// <summary>
        /// Open a connection with the datastore.
        /// </summary>
        /// <param name="connectionId"></param>
        public void Connect(string connectionId)
        {
            var rqst = new Persistence.ConnectRqst();
            rqst.ConnectionId = connectionId;
            this.client.Connect(rqst, this.GetClientContext());
        }

        /// <summary>
        /// Disconnect from the  server.
        /// </summary>
        /// <param name="connectionId">The connection id</param>
        public void Disconnect(string connectionId)
        {
            var rqst = new Persistence.DisconnectRqst();
            rqst.ConnectionId = connectionId;
            this.client.Disconnect(rqst, this.GetClientContext());
        }

        /// <summary>
        /// Ping a given persistence service.
        /// </summary>
        /// <param name="connectionId">The connection id, not it name.</param>
        /// <returns>Must return 'pong'</returns>
        public string Ping(string connectionId)
        {
            // Here I will create the new connection.
            Persistence.PingConnectionRqst rqst = new Persistence.PingConnectionRqst();
            rqst.Id = connectionId;

            // Create a new connection
            var rsp = this.client.Ping(rqst, this.GetClientContext());
            return rsp.Result;
        }

        ///////////////////////////////////// Quering /////////////////////////////////////

        /// <summary>
        /// Find one object from the database.
        /// </summary>
        /// <param name="connectionId">The connection id</param>
        /// <param name="database">The database name</param>
        /// <param name="collection">The collection name</param>
        /// <param name="query">The filter</param>
        /// <param name="options">a list of option, must be a json array</param>
        /// <returns></returns>
        public string FindOne(string connectionId, string database, string collection, string query, string options)
        {
            var rqst = new Persistence.FindOneRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.Query = query;
            rqst.Options = options;

            var rsp = this.client.FindOne(rqst, this.GetClientContext());
            return rsp.JsonStr;
        }

        /// <summary>
        /// Find multiple values from the data store.
        /// </summary>
        /// <param name="connectionId">The connection Id to be used</param>
        /// <param name="database">The database name</param>
        /// <param name="collection">The collection name</param>
        /// <param name="query">The query</param>
        /// <param name="options">a list of option, must be a json array</param>
        /// <returns></returns>
        public string Find(string connectionId, string database, string collection, string query, string options)
        {
            var rqst = new Persistence.FindRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.Query = query;
            rqst.Options = options;

            var call = this.client.Find(rqst, this.GetClientContext());

            // Make the function synchrone...
            string jsonStr = "[";
            bool hasNext = true;

            // read until no more values found...
            while (hasNext)
            {
                var task = Task.Run(() => call.ResponseStream.MoveNext());
                task.Wait(); // wait for the next value...
                hasNext = task.Result;
                if (hasNext)
                {
                    string str = call.ResponseStream.Current.JsonStr;
                    if (jsonStr.Length > 1)
                    {
                        jsonStr += ",";
                    }
                    jsonStr += str.Substring(1, str.Length - 1);
                }
            }

            return jsonStr + "]";
        }

        public string Aggregate(string connectionId, string database, string collection, string pipeline, string options)
        {
            var rqst = new Persistence.AggregateRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.Pipeline = pipeline;
            rqst.Options = options;

            var call = this.client.Aggregate(rqst, this.GetClientContext());

            // Make the function synchrone...
            string jsonStr = "[";
            bool hasNext = true;

            // read until no more values found...
            while (hasNext)
            {
                var task = Task.Run(() => call.ResponseStream.MoveNext());
                task.Wait(); // wait for the next value...
                hasNext = task.Result;
                if (hasNext)
                {
                    string str = call.ResponseStream.Current.JsonStr;
                    if (jsonStr.Length > 1)
                    {
                        jsonStr += ",";
                    }
                    jsonStr += str.Substring(1, str.Length - 1);
                }
            }

            return jsonStr + "]";
        }

        /// <summary>
        /// Count the number of document that match a given query
        /// </summary>
        /// <param name="connectionId">The connection id</param>
        /// <param name="database">The datase</param>
        /// <param name="collection">The collection</param>
        /// <param name="query">The query</param>
        /// <param name="options">A list of options in form of json string</param>
        /// <returns></returns>
        public long Count(string connectionId, string database, string collection, string query, string options)
        {
            var rqst = new Persistence.CountRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.Query = query;
            rqst.Options = options;

            var rsp = this.client.Count(rqst, this.GetClientContext());
            return rsp.Result;
        }

        /// <summary>
        /// Insert one document in the database and return the newly create document id.
        /// </summary>
        /// <param name="connectionId">The connection id</param>
        /// <param name="database">The database name</param>
        /// <param name="collection">The collection</param>
        /// <param name="jsonStr">The oject stringnify value</param>
        /// <param name="options">The options</param>
        /// <returns></returns>
        public string InsertOne(string connectionId, string database, string collection, string jsonStr, string options)
        {
            var rqst = new Persistence.InsertOneRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.JsonStr = jsonStr;
            rqst.Options = options;

            var rsp = this.client.InsertOne(rqst, this.GetClientContext());
            return rsp.Id;
        }

        public void ReplaceOne(string connectionId, string database, string collection, string query, string value, string options)
        {
            var rqst = new Persistence.ReplaceOneRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.Query = query;
            rqst.Value = value;
            rqst.Options = options;

            this.client.ReplaceOne(rqst, this.GetClientContext());
        }

        public void UpdateOne(string connectionId, string database, string collection, string query, string value, string options)
        {
            var rqst = new Persistence.UpdateOneRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.Query = query;
            rqst.Value = value;
            rqst.Options = options;

            this.client.UpdateOne(rqst, this.GetClientContext());
        }


        public void Update(string connectionId, string database, string collection, string query, string value, string options)
        {
            var rqst = new Persistence.UpdateRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.Query = query;
            rqst.Value = value;
            rqst.Options = options;

            this.client.Update(rqst, this.GetClientContext());
        }

        public void DeleteOne(string connectionId, string database, string collection, string query, string options)
        {
            var rqst = new Persistence.DeleteOneRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.Query = query;
            rqst.Options = options;

            this.client.DeleteOne(rqst, this.GetClientContext());
        }

        public void Delete(string connectionId, string database, string collection, string query, string options)
        {
            var rqst = new Persistence.DeleteRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;
            rqst.Query = query;
            rqst.Options = options;

            this.client.Delete(rqst, this.GetClientContext());
        }

        public string InsertMany(string connectionId, string database, string collection, ArrayList objects, string jsonStr, string options)
        {

            // Open a stream with the server.
            var call = this.client.InsertMany(this.GetClientContext());

            // Here i will iterate over the list of object contain in the collection and persist 500 object at time.
            var chunkSize = 500;
            for (var i = 0; i < objects.Count; i += chunkSize)
            {
                var rqst = new Persistence.InsertManyRqst();
                rqst.Id = connectionId;
                rqst.Database = database;
                rqst.Collection = collection;

                if (i + chunkSize < objects.Count)
                {
                    rqst.JsonStr = "[" + JsonSerializer.Serialize(objects.GetRange(i, chunkSize)) + "]";
                }
                else
                {
                    rqst.JsonStr = "[" + JsonSerializer.Serialize(objects.GetRange(i, objects.Count - i)) + "]";
                }

                var task = Task.Run(() => call.RequestStream.WriteAsync(rqst));
                task.Wait(); // wait until the message was sent...
            }

            Task.Run(() => call.RequestStream.CompleteAsync()).Wait();

            var rsp = Task.Run(() => call.ResponseAsync);
            rsp.Wait();

            return rsp.Result.Ids;
        }

        public void DeleteCollection(string connectionId, string database, string collection)
        {
            var rqst = new Persistence.DeleteCollectionRqst();
            rqst.Id = connectionId;
            rqst.Database = database;
            rqst.Collection = collection;

            this.client.DeleteCollection(rqst, this.GetClientContext());
        }

        public void DeleteDatabase(string connectionId, string database)
        {
            var rqst = new Persistence.DeleteDatabaseRqst();
            rqst.Id = connectionId;
            rqst.Database = database;

            this.client.DeleteDatabase(rqst, this.GetClientContext());
        }

        public void RunAdminCmd(string connectionId, string user, string pwd, string script)
        {
            var rqst = new Persistence.RunAdminCmdRqst();
            rqst.ConnectionId = connectionId;
            rqst.Script = script;
            rqst.User = user;
            rqst.Password = pwd;

            this.client.RunAdminCmd(rqst, this.GetClientContext());
        }
        
    }

}
