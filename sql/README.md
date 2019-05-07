# Sql Service
## Access sql directly from the browser

The *sql service* is the simplest way to get access to sql server(s) direcly from the browser. In order to get access to a *sql* server you
must create a connection. To do so you can append the connection information in the config.json file by hand, see the sample config file.

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
Here is the definitions of each field defined in the connection,
* **Id** The id must unique it reprensente
