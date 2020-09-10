#include "globularressourceclient.h"
#include <iostream>

Globular::RessourceClient::RessourceClient(std::string name, std::string domain, unsigned int configurationPort):
    Globular::Client(name,domain, configurationPort)
{

}
