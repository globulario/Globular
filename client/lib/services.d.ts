import { EventServicePromiseClient } from "./event/eventpb/event_grpc_web_pb";
import { EchoServicePromiseClient } from "./echo/echopb/echo_grpc_web_pb";
import { CatalogServicePromiseClient } from "./catalog/catalogpb/catalog_grpc_web_pb";
import { FileServicePromiseClient } from "./file/filepb/file_grpc_web_pb";
import { LdapServicePromiseClient } from "./ldap/ldappb/ldap_grpc_web_pb";
import { PersistenceServicePromiseClient } from "./persistence/persistencepb/persistence_grpc_web_pb";
import { PlcLinkServicePromiseClient } from "./plc_link/plc_link_pb/plc_link_grpc_web_pb";
import { SmtpServicePromiseClient } from "./smtp/smtppb/smtp_grpc_web_pb";
import { SpcServicePromiseClient } from "./spc/spcpb/spc_grpc_web_pb";
import { SqlServicePromiseClient } from "./sql/sqlpb/sql_grpc_web_pb";
import { StorageServicePromiseClient } from "./storage/storagepb/storage_grpc_web_pb";
import { MonitoringServicePromiseClient } from "./monitoring/monitoringpb/monitoring_grpc_web_pb";
import { PlcServicePromiseClient } from "./plc/plcpb/plc_grpc_web_pb";
/**
 * The service configuration information.
 */
interface IServiceConfig {
    Address: string;
    Domain: string;
    Port: Number;
    Proxy: Number;
}
/**
 * Define a map of services.
 */
interface IServices {
    [index: string]: IServiceConfig;
}
/**
 * The application server informations.
 */
interface IConfig {
    Name: string;
    Domain: string;
    PortHttp: Number;
    PortHttps: Number;
    Protocol: string;
    IP: string;
    Services: IServices;
}
/**
 * Globular regroup all serivces in one object that can be use by
 * application to get access to sql, ldap, persistence... service.
 */
export declare class Globular {
    config: IConfig | undefined;
    catalogService: CatalogServicePromiseClient | undefined;
    echoService: EchoServicePromiseClient | undefined;
    eventService: EventServicePromiseClient | undefined;
    fileService: FileServicePromiseClient | undefined;
    ldapService: LdapServicePromiseClient | undefined;
    persistenceService: PersistenceServicePromiseClient | undefined;
    smtpService: SmtpServicePromiseClient | undefined;
    sqlService: SqlServicePromiseClient | undefined;
    storageService: StorageServicePromiseClient | undefined;
    monitoringService: MonitoringServicePromiseClient | undefined;
    spcService: SpcServicePromiseClient | undefined;
    plcService_ab: PlcServicePromiseClient | undefined;
    plcService_simens: PlcServicePromiseClient | undefined;
    plcLinkService: PlcLinkServicePromiseClient | undefined;
    /** The */
    constructor(config: IConfig);
}
export {};
