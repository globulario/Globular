# Persistence
### Nothing in this world can take the place of persistence. 
With emergence of [NoSQL](https://en.wikipedia.org/wiki/NoSQL), new kind of databases became widely used by developper to store application data. The persistence service give access to those data-stores. The interface propose by that service is well suited for document store. At the moment only [MongoDB](https://www.mongodb.com/) is implement by the service, but [ArangoDB](https://www.arangodb.com/) and [CouchDB](http://couchdb.apache.org/) will be available soon. The *Storage* service is also available to store application data, it was created to interface *key-value* stores.

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
    * **CouchDB** 2
* **Host** Can be domain host name or *IPV4* address of the sql server
* **User** The datastore user (can be empty)
* **Password** The datastore password (can be empty)
* **Port** The server port default are,
    * **MongoDB** 3306
    * **ArangoDB** 8529
    * **CouchDB** 5984
* **Timeout** The number of time before give up on connection
* **Options** A json string containing connection options.

### Find
The find operation is use to retreive multiple value from the data store.
```javascript
var testPersistenceResults = []
function testPersistenceFind(){
    var rqst = new Persistence.FindRqst()
    rqst.setId("mongo_db_test_connection")
    rqst.setDatabase("TestMongoDB")
    rqst.setCollection("Employees")
    rqst.setQuery( '{"first_name": "Anneke"}')
    rqst.setFieldsList(["_id", "birth_date"]) // here I will get only the _id and the birth date.

    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.persistenceService.find(rqst, metadata);
    
    // Get the stream and set event on it...
    stream.on('data', function (rsp) {
        testPersistenceResults = testPersistenceResults.concat(JSON.parse(rsp.getJsonstr()))
    });

    stream.on('status', function (status) {
        if(status.code == 0){
            console.log(testPersistenceResults)
        }
    });

    stream.on('end', function (end) {
        // stream end signal
    });
}
```

The same query with the web api,
```http
http://127.0.0.1:10000/api/persistence_service/Find?p0=mongo_db_test_connection&p1=TestMongoDB&p2=Employees&p3={%22first_name%22:%20%22Anneke%22}&p4=_id,birth_date&p5=
```
What we receive as a result is an array of array of tow values. The *_id* and the *birth_date* as specified in the *fields* parameter.
```
["5cd841f5c46c04131d092657","1953-04-20"],["5cd841f5c46c04131d09286c","1955-02-06"]...]
```
