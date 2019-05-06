
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

const { Connection, CreateConnectionRqst,  CreateConnectionRsp, DeleteConnectionRqst, 
    DeleteConnectionRsp, PingConnectionRqst, PingConnectionRsp, Query,  QueryContextRqst,
    QueryContextRsp, ExecContextRqst, ExecContextRsp } = require('./sql/sqlpb/sql_pb.js');

const { SqlServiceClient } = require('./sql/sqlpb/sql_grpc_web_pb.js');

// export the symbol.
window.SqlServiceClient = SqlServiceClient;
window.SqlConnection = Connection;
window.CreateConnectionRqst = CreateConnectionRqst;
window.CreateConnectionRsp = CreateConnectionRsp;
window.DeleteConnectionRqst = DeleteConnectionRqst;
window.DeleteConnectionRsp = DeleteConnectionRsp;
window.PingConnectionRqst = PingConnectionRqst;
window.PingConnectionRsp = PingConnectionRsp;
window.Query = Query;
window.QueryContextRqst = QueryContextRqst;
window.DeleteConnectionRsp = DeleteConnectionRsp;
window.QueryContextRsp = QueryContextRsp;
window.ExecContextRqst = ExecContextRqst;
window.ExecContextRsp = ExecContextRsp;

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

                    globular.sqlService = new SqlServiceClient('http://localhost:' + globular.config.Services.sql_server.Proxy, null, {"format":"binary"});
                    console.log("sql service is init.")
                    
                    window.globular = globular

                    if (window.globularReady != null) {
                        window.globularReady()
                    }
                    console.log("init service done!")
                }
            }
        }(this);

        xmlhttp.open("GET", "config.json", true);
        xmlhttp.send();
    }
}

// Create service connection and 
new Globular()