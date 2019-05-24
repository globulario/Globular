# Globular
Micro-Services Web application framework written in Go and JavaScript with help of gRpc-web. Each service are defined with help of protobufer. At this time available services are:

* [Echo (example service )](https://github.com/davecourtois/globulehub/tree/master#Echo)
* [Sql](https://github.com/davecourtois/globulehub/tree/master/sql)
* [Ldap](https://github.com/davecourtois/globulehub/tree/master/ldap)
* [Smtp](https://github.com/davecourtois/globulehub/tree/master/smtp)
* [Persistence (MongoDB...)](https://github.com/davecourtois/globulehub/tree/master/persistence)
* [Storage (server side html5 storage)](https://github.com/davecourtois/globulehub/tree/master/storage)
* [File (give access to server file system)](https://github.com/davecourtois/globulehub/tree/master/file)

###  *share-as-little-as-possible*
Microservices, aka Microservice Architecture, is an architectural style that structures an application as a collection of small autonomous services, modeled around a business domain. Each service run inside their own process, as a result, each services are isolated from each other, so their impact are minimal. It's possbile to extend the level of functionality of the whole system by simply create a new service.
### What is *Globular*
If the serices are the brick's, then *Globular* is the mortar. All services define by Globular are plain gRpc service. I create those set of service as basic *web-application* brick's. But essensitly any gRpc services can by added to Globular whitout any problems. The role of globular is to manage all those independent services and present them to web-application as a Whole. 

Funtionalities offer by Globular are:

* Give a global entry point to services via a web server (http/https)
* Starting/stopping services
* Keep services alive as needed
* Managing service configuration (name, port number, ...)
* Keep track of logging informations of the whole system
* Give access to services via regular web-api via http-query
* Keep services configuration details hidden from client side

### Echo
Here I will show you how you can create your own personnal service in Globular and use it in your web application. You are welcome to share it here with the rest of pepole as you want, in fact it will be nice to have a micro-services repository ready to use by web-applications.
#### Define your service
The first thing to do is to create a directory named *echo*pb (pb stand for Protocol Buffer). In that directory you will define your service interface. The file [*echo*.proto](https://github.com/davecourtois/Globular/blob/master/echo/echopb/echo.proto) contain the grpc service definition.
```proto
  package echo;

  option go_package="echopb";

  message EchoRequest {
    string message = 1;
  }

  message EchoResponse {
    string message = 1;
    int32 message_count = 2;
  }

  service EchoService {
    // One request followed by one response
    // The server returns the client message as-is.
    rpc Echo(EchoRequest) returns (EchoResponse);
  }
```
#### Create the server
The next step is to create the server directory, that directory will contain three sub-directories:
* *echo*_server That directory contain the gRpc service side code replace *echo* by your actual service name. You can start from the [*echo*_server.go](https://github.com/davecourtois/Globular/blob/master/echo/echo_server/echo_server.go) as starting point. If you use *echo*_server.go your server will create a [*config.json*](https://github.com/davecourtois/Globular/blob/master/echo/echo_server/config.json) file for you the first time it start. That must contain nessary configuration informations use by your service. 
  ``` JSON
   {
    "Name": "echo_server",
    "Port": 10001,
    "Proxy": 10002,
    "AllowAllOrigins": true,
    "AllowedOrigins": "",
    "Protocol": "grpc"
  }
  ```
  * The *Name* of your service (must be unique on your Globular server)
  * The *Port* number (That is the gRpc port)
  * The *Proxy* port number (Use by the browser to access the gRpc service)
  * *AllowAllOrigins* You can give list of address that can access your service, that will block all other origin. By default all address are allow.
  * *Protocol* must be *grpc*, but other *rpc* protocol can be added in the futur.
  If your gRpc server is written in a different language than *Go* you can put your code here. At the end you must have a *config.json* file and an executable file named *echo*_server.*exe* (as example).
  
