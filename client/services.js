
////////////////////////////////////////////////////////////////////////////
// Echo service
////////////////////////////////////////////////////////////////////////////

/**
 * Echo service is a simple test service.
 */
const { EchoRequest, EchoResponse } = require('./echo/echopb/echo_pb.js');
const { EchoServiceClient } = require('./echo/echopb/echo_grpc_web_pb.js');

// export the symbol
window.EchoServiceClient = EchoServiceClient;
window.EchoRequest = EchoRequest;
window.EchoResponse = EchoResponse;

////////////////////////////////////////////////////////////////////////////
// Sql service
////////////////////////////////////////////////////////////////////////////
window.Sql = require('./sql/sqlpb/sql_pb.js');
window.Sql = Object.assign(window.Sql, require('./sql/sqlpb/sql_grpc_web_pb.js'));

////////////////////////////////////////////////////////////////////////////
// Ldap service
////////////////////////////////////////////////////////////////////////////
window.Ldap = require('./ldap/ldappb/ldap_pb.js');
window.Ldap = Object.assign(window.Ldap, require('./ldap/ldappb/ldap_grpc_web_pb.js'));

////////////////////////////////////////////////////////////////////////////
// Smtp service
////////////////////////////////////////////////////////////////////////////
window.Smtp = require('./smtp/smtppb/smtp_pb.js');
window.Smtp = Object.assign(window.Smtp, require('./smtp/smtppb/smtp_grpc_web_pb.js'));

////////////////////////////////////////////////////////////////////////////
// Spc service ( statistical process control )
////////////////////////////////////////////////////////////////////////////
window.Spc = require('./spc/spcpb/spc_pb.js');
window.Spc = Object.assign(window.Spc, require('./spc/spcpb/spc_grpc_web_pb.js'));

////////////////////////////////////////////////////////////////////////////
// Persistence service ( rest )
////////////////////////////////////////////////////////////////////////////
window.Persistence = require('./persistence/persistencepb/persistence_pb.js');
window.Persistence = Object.assign(window.Persistence, require('./persistence/persistencepb/persistence_grpc_web_pb.js'));

////////////////////////////////////////////////////////////////////////////
// Server singleton object that give access to services.
////////////////////////////////////////////////////////////////////////////

/**
 * The singleton to access all services.
 */
class Globular {
    constructor() {
        this.config = null;
        console.log("init the services...")

        // So here I will get the configuration from the active server.
        var xmlhttp = new XMLHttpRequest();
        xmlhttp.onreadystatechange = function (globular) {
            return function () {
                if (this.readyState == 4 && this.status == 200) {
                    globular.config = JSON.parse(this.responseText);

                    // Now I will set serives...
                    globular.echoService = new EchoServiceClient('http://localhost:' + globular.config.Services.echo_server.Proxy);
                    console.log("echo service is init.")

                    if (globular.config.Services.sql_server != null) {
                        globular.sqlService = new Sql.SqlServiceClient('http://localhost:' + globular.config.Services.sql_server.Proxy);
                        console.log("sql service is init.")
                    }

                    if (globular.config.Services.ldap_server != null) {
                        globular.ldapService = new Ldap.LdapServiceClient('http://localhost:' + globular.config.Services.ldap_server.Proxy);
                        console.log("ldap service is init.")
                    }

                    if (globular.config.Services.smtp_server != null) {
                        globular.smtpService = new Smtp.SmtpServiceClient('http://localhost:' + globular.config.Services.smtp_server.Proxy);
                        console.log("smtp service is init.")
                    }

                    if (globular.config.Services.spc_server != null) {
                        globular.spcService = new Spc.SpcServiceClient('http://localhost:' + globular.config.Services.spc_server.Proxy);
                        console.log("spc service is init.")
                    }

                    if (globular.config.Services.persistence_server != null) {
                        globular.persistenceService = new Persistence.PersistenceServiceClient('http://localhost:' + globular.config.Services.persistence_server.Proxy);
                        console.log("persistence service is init.")
                    }

                    window.globular = globular

                    if (window.globularReady != null) {
                        window.globularReady()
                    }
                    console.log("init service done!")
                }
            }
        }(this);

        xmlhttp.open("GET", "/config.json", true);
        xmlhttp.send();
    }
}

// Create service connection and 
new Globular()