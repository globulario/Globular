# [Go to Globular website](https://globular.io)

## Why Globular?

While gRPC is helpful for service implementation, it doesn't address service management. This is where Globular comes in. Service management is the main purpose of Globular, as managing services manually becomes challenging as the number of applications and services increases.

To make a service manageable, it must have the following properties:

* Identifiable: Each service is identified by a unique service ID on a specific server and domain. The ID should remain constant over time to ensure proper functionality of applications using the service.

* Nameable: Multiple instances of the same service should be able to run simultaneously and share the same service name.

* Versionable: Service interfaces may change over time, and applications should be able to access specific versions of the service. Service functionality should not be frozen in time.

* Maintainable: Updating multiple service instances individually can be difficult and error-prone. Proper maintenance mechanisms should be in place.

* Available: If a service crashes, it should be automatically restarted. Losing service instances over time can lead to unstable applications.

* Reachable: Service instances require domain and port/proxy information to determine their network address.

* Trustable: HTTPS and TLS are integrated into gRPC, providing secure communication. Globular assists with certificate creation and management, allowing easy access control for resources and method execution.

* Scalable: Scalability should be considered from the beginning. Globular enables you to create a Globular Cluster, allowing you to run your own cloud. The architecture of Globular prioritizes scalability.

By using Globular, you can manage your microservices and make them available to your web applications.

Click [here](https://globular.io) to learn more about Globular!

General presentation of [Globular](https://medium.com/@dave.courtois60/here-comes-globular-5dee34eb52f8) used as personal cloud.

You can also install Globular using Docker images. Find the Docker image [here](https://hub.docker.com/r/globular/globular).

**Version 1.0 (beta) is now available.**
