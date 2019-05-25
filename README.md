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
Microservices, aka Microservice Architecture, is an architectural style that structures an application as a collection of small autonomous services, modeled around a business domain. Each service run inside it own process, as a result, each service is isolated from all other, so it impact is minimal. It's possbile to extend the level of functionality of the whole system by simply create a new service.
### What is *Globular*
If the services are the brick's, then *Globular* is the mortar. All services define by Globular are plain gRpc service. I create those set of service as basic *web-application* brick's. But essensitly any gRpc services can by added to Globular whitout any problems. The role of globular is to manage all those independent services and present them to web-application as a Whole. 

Funtionalities offer by Globular are,
* Give a global entry point to services via a web server (http/https)
* Starting/stopping services
* Keep services alive as needed
* Managing service configuration (name, port number, ...)
* Keep track of logging informations of the whole system
* Give access to services via regular web-api via http-query
* Keep services configuration details hidden from client side

## Install from binairy
Here's a link to get *Globular* binairy for win64 and linux
[globular.1.0.tar.gz](https://www.dropbox.com/home/git?preview=globular.1.0.tar.gz)
```
tar xvzf globular.1.0.tar.gz -C /path/to/somedirectory
cd /path/to/somedirectory
./Globular(.exe)
```
Globular is up and running at [http://127.0.0.1:10000/](http://127.0.0.1:10000/)

## How to create your own service with Globular
### Echo
Here I will show you how you can create your own personnal service in Globular and use it in your web application. You are welcome to share it here with the rest of pepole as you want, in fact it will be nice to have a micro-services repository ready to use by web-applications.
#### Define your service
The first thing to do is to create a directory named [*echo*pb](https://github.com/davecourtois/Globular/blob/master/echo/echopb) (pb stand for Protocol Buffer). In that directory you will define your service interface. The file [*echo*.proto](https://github.com/davecourtois/Globular/blob/master/echo/echopb/echo.proto) contain the grpc service definition.
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
    // --> hello ---> hello.
    // ...Is There Anybody Out There?
    rpc Echo(EchoRequest) returns (EchoResponse);
  }
```
Now you must generate the gRpc code for the server (*Go*),
```bash
protoc echo/echopb/echo.proto --go_out=plugins=grpc:.
```
And also the client (*JavaScript*)
```bash
protoc echo/echopb/echo.proto --js_out=import_style=commonjs:client
protoc echo/echopb/echo.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
```
You append your command inside the [*generateCode.sh*](https://github.com/davecourtois/Globular/blob/master/generateCode.sh) in case you want to share your service with the rest of the planet.

#### Create the server
The next step is to create the server directory, that directory will contain three sub-directories:
* [*echo*_server](https://github.com/davecourtois/Globular/tree/master/echo/echo_server) That directory contain the gRpc service side code replace *echo* by your actual service name. You can start from the [*echo*_server.go](https://github.com/davecourtois/Globular/blob/master/echo/echo_server/echo_server.go) as starting point. If you use *echo*_server.go your server will create a [*config.json*](https://github.com/davecourtois/Globular/blob/master/echo/echo_server/config.json) file for you the first time it start. That must contain nessary configuration informations use by your service. 
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
  
  #### Test-it!
  Unit test are the best way to make sure your service is working correctly. When your create your test at the same time you almost create the client code that will be use in the next step, so it's not a waste of time! Create a new directory name [*echo*_test](https://github.com/davecourtois/Globular/tree/master/echo/echo_test). Create a file named *echo*_test.go and do your homework!-)
 Now you can start your server in the command shell and call go test...
 ```
 ./echo/echo_server/echo_server &
 cd echo/echo_test
 go test
 ```
 #### Create the web-api
 That step is optional, but if you plan to give access to your service in the old fashion way this is how it's done in Globular. In the file [*clients.go*](https://github.com/davecourtois/Globular/blob/master/clients.go) you must wrote the client side of your gRpc service. (see your test code...)
``` go
   type Echo_Client struct {
    cc *grpc.ClientConn
    c  echopb.EchoServiceClient
  }

  // Create a connection to the service.
  func NewEcho_Client(addresse string) *Echo_Client {
    client := new(Echo_Client)
    client.cc = getClientConnection(addresse)
    client.c = echopb.NewEchoServiceClient(client.cc)
    return client
  }

  // must be close when no more needed.
  func (self *Echo_Client) Close() {
    self.cc.Close()
  }

  func (self *Echo_Client) Echo(msg interface{}) (string, error) {
    rqst := &echopb.EchoRequest{
      Message: Utility.ToString(msg),
    }

    rsp, err := self.c.Echo(context.Background(), rqst)
    if err != nil {
      return "", err
    }

    return rsp.Message, nil
  }
```
Finaly in [*globular.go*](https://github.com/davecourtois/Globular/blob/master/globular.go) register your *NewEcho_Client* function,
```go
  /**
   * Init the service client.
   */
  func (self *Globule) initClients() {
    // Register service constructor function here.
    // The name of the contructor must follow the same pattern.
    Utility.RegisterFunction("NewEcho_Client", NewEcho_Client)
    Utility.RegisterFunction("NewSql_Client", NewSql_Client)
    Utility.RegisterFunction("NewFile_Client", NewFile_Client)
    Utility.RegisterFunction("NewPersistence_Client", NewPersistence_Client)
    Utility.RegisterFunction("NewSmtp_Client", NewSmtp_Client)
    Utility.RegisterFunction("NewLdap_Client", NewLdap_Client)

    // The echo service
    for k, _ := range self.services {
      name := strings.Split(k, "_")[0]
      self.initClient(name)
    }
  }
```

**Important** The syntax of your function name must be like New**Echo**_Client where *Echo* is the name of your service with the first letter capitalized.

Your service will be reachable at the address,
```
http://127.0.0.1:10000/api/echo_service/Echo?p0=Hello
```
Where the echo_service is the name of your gRpc service and Echo is the name of your function. Each parameter must be named
*p*0, *p*1... *p*n.

#### Access your service form the browser with JavaScript
We previously generate the JS code it's now time to append it into Globular. In file [*services.js*](https://github.com/davecourtois/Globular/blob/master/client/services.js) wrote,
```javascript
////////////////////////////////////////////////////////////////////////////
// Echo service
////////////////////////////////////////////////////////////////////////////
window.Echo = require('./echo/echopb/echo_pb.js');
window.Echo = Object.assign(window.Echo, require('./echo/echopb/echo_grpc_web_pb.js'));

```
And in the Globular constructor append line,
```javascript
// Now I will set serives...
if (this.config.Services.echo_server != null) {
  this.echoService = new Echo.EchoServiceClient(this.config.Protocol + ': //' + this.config.IP + ":" + this.config.Services.echo_server.Proxy);
  this.echoServicePromise = new Echo.EchoServicePromiseClient(this.config.Protocol + ': //' + this.config.IP + ":" + this.config.Services.echo_server.Proxy);
  console.log("echo service is init.")
}
```
The next step is to compile the *services.js* file with [*webpack*](https://webpack.js.org/), from the client directory run,
```
npx webpack
```
that will ouput a new file named [services.js](https://github.com/davecourtois/Globular/blob/master/client/dist/services.js) in [*dist*](https://github.com/davecourtois/Globular/tree/master/client/dist) directory

That file must replace the existing file named [*services.js*](https://github.com/davecourtois/Globular/blob/master/WebRoot/js/services.js) in [*WebRoot*](https://github.com/davecourtois/Globular/tree/master/WebRoot/js) directory.

### Start the server
Now all it need is to compile Globular and start it... from the [*root directory*](https://github.com/davecourtois/Globular) of the project
```bash
go build
```
Now you must have an executable file named *Globular(.exe)* in your directory. The server configuation will be created the first time you will start your server. The configuration will vary depending of services found, (directory containing a [*config.json*](https://github.com/davecourtois/Globular/blob/master/WebRoot/config.json) file and executable). If you need to change services configuration change the service [*config.json*](https://github.com/davecourtois/Globular/blob/master/echo/echo_server/config.json)

```bash
./Globular
```
Congratulation you wrote your first gRpc service in Globular!

### Access the service from the browser via JavaScript.
Now there is the steps to access service within the browser,

1. Import those file into your [*index.html*](https://github.com/davecourtois/Globular/blob/master/WebRoot/index.html) file.
```html
<body onload="main();">
  <!-- Services files... -->
  <script src = "http://127.0.0.1:10000/config.json"></script>
  <script src = "http://127.0.0.1:10000/js/services.js"></script>
  ...
</body>
```
2. At the top of your js file here [*test.js*](https://github.com/davecourtois/Globular/blob/master/WebRoot/js/test.js) append,
```javascript
// The service configuration 
globularConfig.IP = "127.0.0.1" // remove it when the site is publish.

// The global service object.
var globular = new Globular()
```
Note that globularConfig is a global variable and it contain the default service connection. The IP address is the external IP address of your server, so here I change it to the local address (*127.0.0.1*) because it's a test...

your service is now ready to use!

#### Use your service...
Here is a simple exemple how to use your Globular service.
``` javascript
function testEcho(str) {
    // Create a new request.
    var request = new Echo.EchoRequest();
    request.setMessage(str);

    globular.echoService.echo(request, {}, function (err, response) {
        // ...
        console.log(response.getMessage())
    });

    // Now I will test with promise
    globular.echoServicePromise.echo(request)
        .then((resp) => {
            console.log(resp.getMessage())
        })
        .catch((error) => {
            console.log(error)
        })
}
```
Has you can see, the *globular* object contain a reference to all services that you have define one your application server. With gRpc-Web distant object are used just like local one. As an application developper I found it more natural to work with object than http call.

#### Connect more than one service at time.
You are not limited to one server connection, to connect your application to another server all you have to do is to create a new *Globular* object and specify a different configuration as parameter,
```javascript
var myOtherServer = new Globular({
  "Name": "MyOtherSeverName",
  "Port": "10001",
  "Protocol": "http",
  "IP": "127.0.0.1",
  "Services": {
    "echo_server": {
      "Port": 10001,
      "Proxy": 10002
    },
    "file_server": {
      "Port": 10011,
      "Proxy": 10012
    }
  }
})
```
## Conclusion
Imagine a constellation of micro-services availables to create any type of application with different language accessible in the browser. Imagine a globular cluster...
