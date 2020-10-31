#include <iostream>
#include "../cxxopts.hpp"
#include "echoserviceimpl.h"

using namespace std;
//#pragma comment(lib,"ws2_32.lib")

int main(int argc, char** argv)
{
    cxxopts::Options options("c++ echo service", "A c++ gRpc service implementation");
    auto result = options.parse(argc, argv);

    // Instantiate a new server.
    EchoServiceImpl service("test", "echo.EchoService");

    // Start the service.
    service.run(&service);

    return 0;
}
