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
import { SearchServicePromiseClient } from './search/searchpb/search_grpc_web_pb';
import { LoadBalancingServicePromiseClient } from './lb/lb_grpc_web_pb';
/**
 * The service configuration information.
 */
export interface IServiceConfig {
    Id: string;
    Name: string;
    State: string;
    Domain: string;
    Port: number;
    Proxy: number;
    TLS: boolean;
    KeepUpToDate: boolean;
    KeepAlive: boolean;
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
    PortHttp: number;
    PortHttps: number;
    AdminPort: number;
    AdminProxy: number;
    AdminEmail: string;
    RessourcePort: number;
    RessourceProxy: number;
    ServicesDiscoveryPort: number;
    ServicesDiscoveryProxy: number;
    ServicesRepositoryPort: number;
    ServicesRepositoryProxy: number;
    CertificateAuthorityPort: number;
    CertificateAuthorityProxy: number;
    LoadBalancingServiceProxy: number;
    SessionTimeout: number;
    Protocol: string;
    Discoveries: string[];
    DNS: string[];
    CertExpirationDelay: number;
    CertStableURL: string;
    CertURL: string;
    IdleTimeout: number;
    Services: IServices;
}
/**
 * That local and distant event hub.
 */
export declare class EventHub {
    private service;
    private subscribers;
    private subscriptions;
    private uuid;
    /**
     * @param {*} service If undefined only local event will be allow.
     */
    constructor(service: any);
    /**
     * @param {*} name The name of the event to subcribe to.
     * @param {*} onsubscribe That function return the uuid of the subscriber.
     * @param {*} onevent That function is call when the event is use.
     */
    subscribe(name: string, onsubscribe: (uuid: string) => any, onevent: (data: any) => any, local: boolean): void;
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
    loadBalancingService: LoadBalancingServicePromiseClient | undefined;
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
    searchService: SearchServicePromiseClient | undefined;
    plcService_ab: PlcServicePromiseClient | undefined;
    plcService_siemens: PlcServicePromiseClient | undefined;
    plcLinkService: PlcLinkServicePromiseClient | undefined;
    /** The configuation. */
    constructor(config: IConfig);
    getFirstConfigByName(name: string): IServiceConfig;
}
