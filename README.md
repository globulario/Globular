# [Globular](https://www.globular.io)

## Why Globular?
If gRPC is of great help with service implementation, it does nothing about service management. In fact, service management is out of the scope of gRPC, but it’s the main purpose of globular. Managing service manually can be easy if you have only couple of service, and the number of application is limited. Over time number of application and services has tendency to increase, and when that happened, you can became a victim of your success.

What propertie must be define on a service to make it manageable? To be manageable a service must be,

* Identifiable: The service Id identified a running service on a given server with a given domain. The Id must be unique on the server, and must not change over time, or application that using it will stop working correctly.

* Nameable: Multiple instances of the same service must be able to run at the same time, in redundancy, instance must share the same service name.

* Version able: Service interface can change overtime, application must be able to get access to specific service version. With version service functionality are not freeze in time.

* Updateable: When many service instances are running it can be difficult and error prone to update them one by one.

* Available: If a service crash for any reason, it must be restart. Over time loosing service instances can result in unstable applications.

* Reachable: Here tree properties are requires, the domain and the port/proxy pair. Those properties are used to get the instance network address.

* Trustable: The publisher Id defines the identity of service creator. Globular can be used to authenticate and validate service publisher.

* Securable: The TLS variable defines if the service must use a secure network or not.

By using Globular you will be able to manage your microservices and make it avaible to your web-application.

Click [here](https://www.globular.io) and learn more about globular!

** The vesion 1.0 is available. The website is not 100% finish but installation and quickstart are ready to help you to make your first step. A complete tutorial he's on the way. All documentation must be written before the end of feburary.
