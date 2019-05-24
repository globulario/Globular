# Globular
Micro-Services Web application framework written in Go and JavaScript with help of gRpc-web. Each service are defined with help of protobufer. At this time available services are:

* echo (service example)
* [Sql](https://github.com/davecourtois/globulehub/tree/master/sql)
* [Ldap](https://github.com/davecourtois/globulehub/tree/master/ldap)
* [Smtp](https://github.com/davecourtois/globulehub/tree/master/smtp)
* [Persistence (MongoDB...)](https://github.com/davecourtois/globulehub/tree/master/persistence)
* [Storage (server side html5 storage)](https://github.com/davecourtois/globulehub/tree/master/storage)
* [File (give access to server file system)](https://github.com/davecourtois/globulehub/tree/master/file)

###  *share-as-little-as-possible*
Microservices, aka Microservice Architecture, is an architectural style that structures an application as a collection of small autonomous services, modeled around a business domain. Each service run inside their own process, as a result, each services are isolated from each other, so their impact are minimal. It's possbile to extend the level of functionality of the whole system by simply create a new service.
### What is *Globular*
If the serices are the brick's, then *Globular* is the mortar. All services define by Globular are plain gRpc service. I create those set of service as basic *web-application* brick's. But essensitly any gRpc service can by added to Globular whitout problem. The role of globular is to manage all those independent services and present them to web-application as a Whole. 

Functionality offer by Globular are:

* Give a global entry point to services via a web server (http/https)
* Starting/stopping service
* Keep services alive as needed
* Managing service configuration (name, port number, ...)
* Keep track of logging informations of the whole system
* Give access to services via regular web-api via http-query

