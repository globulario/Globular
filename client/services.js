////////////////////////////////////////////////////////////////////////////
// Echo service
////////////////////////////////////////////////////////////////////////////
window.Echo = require('./echo/echopb/echo_pb.js');
window.Echo = Object.assign(window.Echo, require('./echo/echopb/echo_grpc_web_pb.js'));

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
// Storage service ( kv/cache )
////////////////////////////////////////////////////////////////////////////
window.Storage = require('./storage/storagepb/storage_pb.js');
window.Storage = Object.assign(window.Storage, require('./storage/storagepb/storage_grpc_web_pb.js'));

////////////////////////////////////////////////////////////////////////////
// File service
////////////////////////////////////////////////////////////////////////////
window.File = require('./file/filepb/file_pb.js');
window.File = Object.assign(window.File, require('./file/filepb/file_grpc_web_pb.js'));

////////////////////////////////////////////////////////////////////////////
// Server singleton object that give access to services.
////////////////////////////////////////////////////////////////////////////

// Global variables. Those variable are intercept by the actual server and are 
// change automaticaly with the network name and port.

/**
 * The singleton to access all services.
 */
class Globular {
    constructor(config) {
        console.log("init the services...")
        this.config = config
        if(this.config == undefined){
            this.config = globularConfig
        }

        if(this.config == undefined){
            console.log("no configuration found!")
        }

        // Now I will set serives...
        if (this.config.Services.sql_server != null) {
            this.echoService = new Echo.EchoServiceClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.echo_server.Proxy);
            this.echoServicePromise = new Echo.EchoServicePromiseClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.echo_server.Proxy);
            console.log("echo service is init.")
        }

        if (this.config.Services.sql_server != null) {
            this.sqlService = new Sql.SqlServiceClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.sql_server.Proxy);
            this.sqlServicePromise = new Sql.SqlServicePromiseClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.sql_server.Proxy);
            console.log("sql service is init.")
        }

        if (this.config.Services.ldap_server != null) {
            this.ldapService = new Ldap.LdapServiceClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.ldap_server.Proxy);
            this.ldapServicePromise = new Ldap.LdapServicePromiseClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.ldap_server.Proxy);
            console.log("ldap service is init.")
        }

        if (this.config.Services.smtp_server != null) {
            this.smtpService = new Smtp.SmtpServiceClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.smtp_server.Proxy);
            this.smtpServicePromise = new Smtp.SmtpServicePromiseClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.smtp_server.Proxy);
            console.log("smtp service is init.")
        }

        if (this.config.Services.spc_server != null) {
            this.spcService = new Spc.SpcServiceClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.spc_server.Proxy);
            this.spcServicePromise = new Spc.SpcServicePromiseClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.spc_server.Proxy);
            console.log("spc service is init.")
        }

        if (this.config.Services.persistence_server != null) {
            this.persistenceService = new Persistence.PersistenceServiceClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.persistence_server.Proxy);
            this.persistenceServicePromise = new Persistence.PersistenceServicePromiseClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.persistence_server.Proxy);
            console.log("persistence service is init.")
        }

        if (this.config.Storage.storage_server != null) {
            this.storageService = new Storage.StorageServiceClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.storage_server.Proxy);
            this.storageServicePromise = new Storage.StorageServicePromiseClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.storage_server.Proxy);
            console.log("storage service is init.")
        }

        if (this.config.Services.file_server != null) {
            this.fileService = new File.FileServiceClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.file_server.Proxy);
            this.fileServicePromise = new File.FileServicePromiseClient(this.config.Protocol + '://' + this.config.IP + ":" + this.config.Services.file_server.Proxy);
            console.log("file service is init.")
        }

        console.log("services are all initialyse!")
    }

}

// export the class Globular.
window.Globular = Globular;