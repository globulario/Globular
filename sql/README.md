# Sql Service
## Access sql directly from the browser

The *sql service* is the simplest way to get access to *sql* direcly from the browser. All you have to do it's to configure connections on the sever side and use it in clients.

### Configure a new connection
[Config](https://github.com/davecourtois/Globular/blob/master/sql/sql_server/config.json) file reside on server side and contain connections configurations. There is an example of configuration,
```
    "sql_server": {
      "AllowAllOrigins": true,
      "AllowedOrigins": "",
      "Connections": {
        "employees_db": {
          "Charset": "utf8",
          "Driver": "mysql",
          "Host": "localhost",
          "Id": "employees_db",
          "Name": "employees",
          "Password": "password",
          "Port": 3306,
          "User": "test"
        }
      },
      
```
Connection configuration contain fields,
* **Id** The id must unique
* **Name** The name of the database
* **Charset** The encoding of the database
* **Driver** The database driver string can be one of,
    * **mysql** MySql
    * **mssql** Microsoft Sql Server
    * **postgres** Postgres Sql driver
    * **odbc** Generic ODBC driver, work with (*Oracle*, *Microsoft Sql Server*)
* **Host** Can be domain host name or *IPV4* address of the sql server
* **User** The sql user (not the windows user in case of mssql)
* **Password** The user password
* **Port** The server port default are,
    * **MySql** 3306
    * **Microsoft Sql Server** 1433
    * **Postgre Sql** 5432
    
Because service [*config.json*](https://github.com/davecourtois/Globular/blob/master/sql/sql_server/config.json) file live on the server side, ther is no way for the client to access information from it. That's keep your connection information safe. All the client need to know it's the connection Id at time of connection.

### about odbc
if you plan to use odb on linux server you must install [unixODBC](http://www.unixodbc.org/).
to do so,

1. download the lastest version of [unixODBC](ftp://ftp.unixodbc.org/pub/unixODBC/unixODBC-2.3.7.tar.gz)
2. uncompress it 
    ```
    tar -xvzf unixODBC-2.3.7.tar.gz
    ```
3. now build it
    ```
    cd unixODBC-2.3.7
    ./configure
    sudo make all install clean
    ```
you are now ready to use odbc on linux.

## Setup connection with Globular service
1. First of all you must have a globular service configure, to do so
    in your index.html project file
    ```html
    <body onload="main();">
      <!-- Services files... -->
      <script src = "http://127.0.0.1:10000/config.json"></script>
      <script src = "http://127.0.0.1:10000/js/services.js"></script>
      ...
    ```
    in your javascript file initialyse a new globular object.
    
    ```javascript
    // The service configuration 
    globularConfig.IP = "127.0.0.1" // remove it when the site is publish.

    // The global service object.
    var globular = new Globular()
    ```
## Ping
Now you can test if your connection is correctly configure with help of ping.
```javascript
function testPing(){
    var rqst = new Sql.PingConnectionRqst()
    rqst.setId("employees_db")

    globular.sqlServicePromise.ping(rqst)
    .then((rsp) => {
        console.log(rsp.getResult())
    })
    .catch((error) => {
        console.log(error)
    })
}
```
if your connection is correctly configure you must receive answer *pong*.

## QueryContext
The *QueryContext* must be use for sql **select**.

### from javascript
Here I will show you how you can execute a query direclty from the browser. To do so we will use *QueryContext* function.
```javascript
function testSelectQuery() {
    var rqst = new Sql.QueryContextRqst()
    var q = new Sql.Query()
    q.setQuery("SELECT first_name, last_name FROM employees.employees WHERE gender=?")
    q.setConnectionid("employees_db")
    q.setParameters(JSON.stringify(["F"]))

    rqst.setQuery(q)
    var metadata = {}; // custom header value as needed...
    var stream = globular.sqlService.queryContext(rqst, metadata);
    
    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if(response.hasHeader()){
            var header = response.getHeader()
            console.log(JSON.parse(header));
        }else if(response.hasRows()){
            var rows = response.getRows()
            console.log(JSON.parse(rows));
        }
    });

    stream.on('status', function (status) {
        console.log(status.code);
        console.log(status.details);
        console.log(status.metadata);
    });

    stream.on('end', function (end) {
        // stream end signal
    });
 }
```
There is no magic here just plain gRpc-Web call. The service made use of a stream to tranfert data back to the client, in that way you can use it before the end of transfert and make your application more responsive. The result contain tow types of informations,
* **header** Contain a json string with column informations (name, data type, number format...) Usefull to create the *gui*.
    ```json
    0:
        name: "first_name"
        typeInfo:
            DatabaseTypeName: "VARCHAR"
            IsNull: false
            IsNullable: true
            Name: "VARCHAR"
    ```
* **data** Contain a json string *"[[v0,v1,v2], [v0,v1,v2], [v0,v1,v2]]"* that represent an array of rows.

### from web api
Here is the same query from the web api,
```http
http://127.0.0.1:10000/api/sql_service/QueryContext?p0=employees_db&p1=SELECT%20first_name,%20last_name%20FROM%20employees.employees%20WHERE%20gender=%3F%20AND%20first_name=%3F%20AND%20last_name=%3F&p2=F,Ebbe,Denis
```
Here is the return result.
```json
{  
    "data":[  
      [  
         "Ebbe",
         "Denis"
      ]
    ],
    "header":[  
      {  
         "name":"first_name",
         "typeInfo":{  
            "DatabaseTypeName":"VARCHAR",
            "IsNull":false,
            "IsNullable":true,
            "Name":"VARCHAR"
         }
      },
      {  
         "name":"last_name",
         "typeInfo":{  
            "DatabaseTypeName":"VARCHAR",
            "IsNull":false,
            "IsNullable":true,
            "Name":"VARCHAR"
         }
      }
    ]
}
```
## ExecContext
The *ExecContext* must be use for sql **INSERT**, **UPDATE**, **DELETE**, **CREATE TABLE**, **DROP**.

### from javascript
Here I will show you how to insert and delete data direclty from the browser... other operation are use the same way...
```javascript
function testInsertQuery(){
    var rqst = new Sql.ExecContextRqst()
    var q = new Sql.Query()
    q.setQuery("INSERT INTO employees.employees (emp_no, first_name, last_name, gender, hire_date, birth_date) VALUE(?,?,?,?,?,?)")
    q.setConnectionid("employees_db")
    q.setParameters(JSON.stringify([200000, 'Dave', 'Courtois', 'M', '2007-07-01', '1976-01-29']))

    rqst.setQuery(q)
    rqst.setTx(false)

    globular.sqlServicePromise.execContext(rqst)
    .then((rsp) => {
        var affectedRow = rsp.getAffectedrows()
        var lastId = rsp.getLastid()
        console.log("affected rows: ", affectedRow, " last id: ", lastId)
    })
    .catch((error) => {
        console.log(error)
    })
}

function testDeleteQuery(){
    var rqst = new Sql.ExecContextRqst()
    var q = new Sql.Query()
    q.setQuery("DELETE FROM employees.employees WHERE emp_no=?")
    q.setConnectionid("employees_db")
    q.setParameters(JSON.stringify([200000]))

    rqst.setQuery(q)
    rqst.setTx(false)

    globular.sqlServicePromise.execContext(rqst)
    .then((rsp) => {
        var affectedRow = rsp.getAffectedrows()
        var lastId = rsp.getLastid()
        console.log("affected rows: ", affectedRow, " last id: ", lastId)
    })
    .catch((error) => {
        console.log(error)
    })
}
```
Again all you see it's almost plain gRpc-Web call. The console display the number of affected rows and the last insert id in case of auto incremented key.

### from web api

Here is the http query,

* for **INSERT** 
    ```http
    http://127.0.0.1:10000/api/sql_service/ExecContext?p0=employees_db&p1=INSERT%20INTO%20employees.employees%20(emp_no,%20first_name,%20last_name,%20gender,%20hire_date,%20birth_date)%20VALUE(%3F,%3F,%3F,%3F,%3F,%3F)&p2=200000,Dave,Courtois,M,%202007-07-01,1976-01-29&p3=false
    ````
* for **DELETE** 
    ```http
    http://127.0.0.1:10000/api/sql_service/ExecContext?p0=employees_db&p1=DELETE%20FROM%20employees.employees%20WHERE%20emp_no=%3F&p2=200000&p3=false
    ```

Note that the *p3* paremeter tell if the query must use transaction.

Those tow query return the number of affected rows and the last id as a json object,

```json
{
    "affectRows":1,
    "lastId":0
}
```
