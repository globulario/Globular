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

## Select Query

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
