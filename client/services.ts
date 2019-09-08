
// Here is the list of services from the backend.
import { EventServicePromiseClient } from './event/eventpb/event_grpc_web_pb';
import { EchoServicePromiseClient } from './echo/echopb/echo_grpc_web_pb';
import { CatalogServicePromiseClient } from './catalog/catalogpb/catalog_grpc_web_pb';
import { FileServicePromiseClient } from './file/filepb/file_grpc_web_pb';
import { LdapServicePromiseClient } from './ldap/ldappb/ldap_grpc_web_pb';
import { PersistenceServicePromiseClient } from './persistence/persistencepb/persistence_grpc_web_pb';
import { PlcLinkServicePromiseClient } from './plc_link/plc_link_pb/plc_link_grpc_web_pb';
import { SmtpServicePromiseClient } from './smtp/smtppb/smtp_grpc_web_pb';
import { SpcServicePromiseClient } from './spc/spcpb/spc_grpc_web_pb';
import { SqlServicePromiseClient } from './sql/sqlpb/sql_grpc_web_pb';
import { StorageServicePromiseClient } from './storage/storagepb/storage_grpc_web_pb';
import { MonitoringServicePromiseClient } from './monitoring/monitoringpb/monitoring_grpc_web_pb';

/**
 * The service configuration information.
 */
interface IServiceConfig {
    Address: string
    Domain: string
    Port: Number
    Proxy: Number
}

/**
 * Define a map of services.
 */
interface IServices{
    [index: string]: IServiceConfig;
}

/**
 * The application server informations.
 */
interface IConfig {
    Name: string
    Domain: string
    PortHttp: Number
    PortHttps: Number
    Protocol: string
    IP: string

    // The map of service object.
    Services: IServices
}

/**
 * Globular regroup all serivces in one object that can be use by 
 * application to get access to sql, ldap, persistence... service.
 */
export class Globular {
    config : IConfig

    // list of services.
    catalogService: CatalogServicePromiseClient
    echoService: EchoServicePromiseClient
    eventService: EventServicePromiseClient
    fileService: FileServicePromiseClient
    ldapService: LdapServicePromiseClient
    persistenceService: PersistenceServicePromiseClient
    smtpService: SmtpServicePromiseClient
    sqlService: SqlServicePromiseClient
    storageService: StorageServicePromiseClient
    monitoringService: MonitoringServicePromiseClient
    spcService: SpcServicePromiseClient
    
    // Non open source services.
    plcService_ab: PlcLinkServicePromiseClient
    plcService_simens: PlcLinkServicePromiseClient
    plcLinkService: PlcLinkServicePromiseClient

    /** The */
    constructor(config: IConfig){
        // Iinitialisation of services.
        if(this.config.Services["catalog_server"]!=null){
            this.catalogService = new CatalogServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["catalog_server"].Proxy, null, null);
        }
        if(this.config.Services["echo_server"]!=null){
            this.echoService = new EchoServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["echo_server"].Proxy, null, null);
        }
        if(this.config.Services["event_server"]!=null){
            this.eventService = new EventServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["event_server"].Proxy, null, null);
        }
        if(this.config.Services["file_server"]!=null){
            this.fileService = new FileServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["file_server"].Proxy, null, null);
        }
        if(this.config.Services["ldap_server"]!=null){
            this.ldapService = new LdapServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["ldap_server"].Proxy, null, null);
        }
        if(this.config.Services["persistence_server"]!=null){
            this.persistenceService = new PersistenceServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["persistence_server"].Proxy, null, null);
        }
        if(this.config.Services["smtp_server"]!=null){
            this.smtpService = new SmtpServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["smtp_server"].Proxy, null, null);
        }
        if(this.config.Services["sql_server"]!=null){
            this.sqlService = new SqlServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["sql_server"].Proxy, null, null);
        }
        if(this.config.Services["storage_server"]!=null){
            this.storageService = new StorageServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["storage_server"].Proxy, null, null);
        }
        if(this.config.Services["monitoring_server"]!=null){
            this.monitoringService = new MonitoringServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["monitoring_server"].Proxy, null, null);
        }
        if(this.config.Services["spc_server"]!=null){
            this.spcService = new SpcServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["spc_server"].Proxy, null, null);
        }

        // non open source services.
        if(this.config.Services["plc_server_ab"]!=null){
            this.plcService_ab = new PlcLinkServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["plc_server_ab"].Proxy, null, null);
        }
        if(this.config.Services["plc_server_simens"]!=null){
            this.plcService_simens = new PlcLinkServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["plc_server_simens"].Proxy, null, null);
        }
        if(this.config.Services["plc_link_server"]!=null){
            this.plcLinkService = new PlcLinkServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ":" + this.config.Services["plc_link_server"].Proxy, null, null);
        }
    }
}