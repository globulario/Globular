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
import { PlcServicePromiseClient } from './plc/plcpb/plc_grpc_web_pb';
import { AdminServicePromiseClient } from './admin/admin_grpc_web_pb';
import { RessourceServicePromiseClient } from './ressource/ressource_grpc_web_pb';
import { ServiceDiscoveryPromiseClient, ServiceRepositoryPromiseClient } from './services/services_grpc_web_pb';
import { CertificateAuthorityPromiseClient } from './ca/ca_grpc_web_pb';
/**
 * The service configuration information.
 */
export interface IServiceConfig {
    Name: string;
    Domain: string;
    Port: Number;
    Proxy: Number;
    TLS: Boolean;
    KeepUpToDate: Boolean;
    PublisherId: string;
    Version: string;
}
/**
 * Define a map of services.
 */
export interface IServices {
    [key: string]: IServiceConfig;
}
/**
 * The application server informations.
 */
export interface IConfig {
    Name: string;
    Domain: string;
    PortHttp: Number;
    PortHttps: Number;
    AdminPort: Number;
    AdminProxy: Number;
    AdminEmail: Number;
    RessourcePort: Number;
    RessourceProxy: Number;
    ServicesDiscoveryPort: Number;
    ServicesDiscoveryProxy: Number;
    ServicesRepositoryPort: Number;
    ServicesRepositoryProxy: Number;
    CertificateAuthorityPort: Number;
    CertificateAuthorityProxy: Number;
    SessionTimeout: Number;
    Protocol: string;
    Discoveries: Array<string>;
    DNS: Array<string>;
    CertExpirationDelay: Number;
    CertStableURL: string;
    CertURL: string;
    IdleTimeout: number;
    Services: IServices;
}
/**
 * That local and distant event hub.
 */
export declare class EventHub {
    readonly service: any;
    readonly subscribers: any;
    readonly subscriptions: any;
    /**
     * @param {*} service If undefined only local event will be allow.
     */
    constructor(service: any);
    /**
     * @param {*} name The name of the event to subcribe to.
     * @param {*} onsubscribe That function return the uuid of the subscriber.
     * @param {*} onevent That function is call when the event is use.
     */
    subscribe(name: string, onsubscribe: (uuid: string) => any, onevent: (data: any) => any): void;
    /**
     *
     * @param {*} name
     * @param {*} uuid
     */
    unSubscribe(name: string, uuid: string): void;
    /**
     * Publish an event on the bus, or locally in case of local event.
     * @param {*} name The  name of the event to publish
     * @param {*} data The data associated with the event
     * @param {*} local If the event is not local the data must be seraliaze before sent.
     */
    publish(name: string, data: any, local: boolean): void;
    /** Dispatch the event localy */
    dispatch(name: string, data: any): void;
}
/**
 * Globular regroup all serivces in one object that can be use by
 * application to get access to sql, ldap, persistence... service.
 */
export declare class Globular {
    config: IConfig | undefined;
    adminService: AdminServicePromiseClient | undefined;
    ressourceService: RessourceServicePromiseClient | undefined;
    servicesDicovery: ServiceDiscoveryPromiseClient | undefined;
    servicesRepository: ServiceRepositoryPromiseClient | undefined;
    certificateAuthority: CertificateAuthorityPromiseClient | undefined;
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
    plcService_siemens: PlcServicePromiseClient | undefined;
    plcLinkService: PlcLinkServicePromiseClient | undefined;
    /** The */
    constructor(config: IConfig);
}
