# Sql Service
## Access sql directly from the browser

The *sql service* is the simplest way to get access to *sql* direcly from the browser.

### Set a new connection
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
* **USER** The sql user (not the windows user in case of mssql)
* **Password** The user password
* **Port** The server port default are,
    * **MySql** 3306
    * **Microsoft Sql Server** 1433
    * **Postgre Sql** 5432
