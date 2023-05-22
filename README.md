
<img src="https://globular.io/images/globular_logo.svg" alt="Globular Logo" width="300">

## [Go to Globular website](https://globular.io)

Globular is an innovative software solution that seamlessly combines the functionalities of popular tools such as Skype, Dropbox, Medium, and Plex. Designed with developers in mind, Globular empowers the creation of web applications following microarchitecture principles. By leveraging Globular, developers can tap into a wide array of database options, including MongoDB, SQL Server (including SQLite, SQL Server, MariaDB, MySQL), LevelDB, Badger, Prometheus, and more, providing unparalleled flexibility and choice.

Simplifying the development of web applications, Globular offers essential features like robust role-based access control (RBAC) and efficient user management. It serves as a comprehensive solution that can potentially replace traditional cloud services like Firebase, streamlining the process of building and managing web applications.

## Why Globular?

While gRPC proves beneficial for service implementation, it often neglects the critical aspect of service management. This is where Globular excels. Service management lies at the core of Globular's purpose, particularly as the number of applications and services continues to grow, making manual management increasingly challenging.

To ensure optimal service manageability, Globular encompasses key properties, including:

- Identifiability: Each service is assigned a unique service ID on a specific server and domain, ensuring consistent functionality for applications relying on the service.
- Nameability: Multiple instances of the same service can seamlessly coexist and share a common service name.
- Versionability: Service interfaces may evolve over time, allowing applications to access specific versions of the service while avoiding stagnation of service functionality.
- Maintainability: Updating multiple service instances individually can be cumbersome and error-prone. Globular provides robust maintenance mechanisms to simplify the process.
- Availability: In the event of a service crash, Globular automatically restarts it, mitigating the risk of unstable applications resulting from service instance loss over time.
- Reachability: Service instances require domain and port/proxy information to determine their network address, ensuring seamless connectivity.
- Trustability: By integrating HTTPS and TLS into gRPC, Globular enables secure communication. It facilitates certificate creation and management, ensuring easy access control for resources and method execution.
- Scalability: With scalability at its core, Globular allows the creation of a Globular Cluster, empowering users to run their own cloud. The architectural design of Globular prioritizes scalability from the outset.

By adopting Globular, developers gain the ability to effectively manage microservices and make them readily available to their web applications.

To delve deeper into Globular's capabilities and explore its potential, visit the [official Globular website](https://globular.io). Additional resources are available to provide further insights:

- [General presentation of Globular](https://medium.com/@dave.courtois60/here-comes-globular-5dee34eb52f8) used as a personal cloud.
- [Installation and configuration guide for your server](https://medium.com/@dave.courtois60/in-this-article-i-will-guide-you-through-the-installation-and-configuration-of-your-personal-cloud-f8bdce33d33a).
- [Installing Globular using Docker](https://medium.com/@dave.courtois60/installing-globular-using-docker-fabd4f96b095).
- Docker image for Globular available on [Docker Hub](https://hub.docker.com/r/globular/globular).

Excitingly, **Version 1.0 (beta)** of Globular is now available, introducing a host of additional features and enhancements to empower developers in their journey.
