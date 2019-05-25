# Persistence
### Nothing in this world can take the place of persistence. 
With emergence of [NoSQL](https://en.wikipedia.org/wiki/NoSQL), new kind of databases became widely used by developper to store application data. The persistence service give access to those data-stores. The interface propose by that service is well suited for document store. At the moment only [MongoDB](https://www.mongodb.com/) is implement by the service, but [ArangoDB](https://www.arangodb.com/) will be available soon, depending of interest of developper. The *Storage* service is also available to store application data, it was created to interface *key-value* stores.

### Configure a new connection
[Config](https://github.com/davecourtois/Globular/blob/master/persistence/persistence_server/config.json) file reside on server side and contain connections configurations. There is an example of configuration,
  ```
  {
    "Name": "persistence_server",
    "Port": 10005,
    "Proxy": 10006,
    "Protocol": "grpc",
    "AllowAllOrigins": true,
    "AllowedOrigins": "",
    "Connections": {
      "mongo_db_test_connection": {
        "Id": "mongo_db_test_connection",
        "Name": "TestMongoDB",
        "Host": "localhost",
        "Store": 0,
        "User": "",
        "Password": "",
        "Port": 27017,
        "Timeout": 0,
        "Options": ""
      }
    }
  }
  ```
Connection configuration contain fields,
* **Id** The id must unique
* **Name** The name of the database
* **Store** The database driver string can be one of,
    * **MongoDB** 0
    * **ArangoDB** 1 
* **Host** Can be domain host name or *IPV4* address of the sql server
* **User** The datastore user (can be empty)
* **Password** The datastore password (can be empty)
* **Port** The server port default are,
    * **MongoDB** 3306
    * **ArangoDB** 8529
* **Timeout** The number of time before give up on connection
* **Options** A json string containing connection options.

