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
import { PlcServicePromiseClient } from './plc/plcpb/plc_grpc_web_pb';
import { AdminServicePromiseClient } from './admin/admin_grpc_web_pb';
import { RessourceServicePromiseClient } from './ressource/ressource_grpc_web_pb';

/**
 * The service configuration information.
 */
export interface IServiceConfig {
  Name: string;
  Address: string;
  Domain: string;
  Port: Number;
  Proxy: Number;
  TLS: Boolean;
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
  AdminProxy: Number;
  RessourceProxy: Number;
  Protocol: string;
  IP: string;

  // The map of service object.
  Services: IServices;
}

/**
 * Globular regroup all serivces in one object that can be use by
 * application to get access to sql, ldap, persistence... service.
 */
export class Globular {
  config: IConfig | undefined;
  // The admin service.
  adminService: AdminServicePromiseClient | undefined;
  ressourceService: RessourceServicePromiseClient | undefined;

  // list of services.
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

  // Non open source services.
  plcService_ab: PlcServicePromiseClient | undefined;
  plcService_siemens: PlcServicePromiseClient | undefined;
  plcLinkService: PlcLinkServicePromiseClient | undefined;

  /** The */
  constructor(config: IConfig) {
    // Keep the config...
    this.config = config;

    /** The admin service to access to other configurations. */
    this.adminService = new AdminServicePromiseClient(
      this.config.Protocol + '://' + this.config.Domain + ':' + this.config.AdminProxy,
      null,
      null,
    );

    this.ressourceService = new RessourceServicePromiseClient(
      this.config.Protocol + '://' + this.config.Domain + ':' + this.config.RessourceProxy,
      null,
      null,
    );

    // Iinitialisation of services.
    if (this.config.Services['catalog_server'] != undefined) {
      let protocol = 'http';
      if (this.config.Services['catalog_server'].TLS == true) {
        protocol = 'https';
      }
      this.catalogService = new CatalogServicePromiseClient(
        protocol +
          '://' +
          this.config.Services['catalog_server'].Domain +
          ':' +
          this.config.Services['catalog_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['echo_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['echo_server'].TLS == true) {
        protocol = 'https';
      }
      this.echoService = new EchoServicePromiseClient(
        protocol + '://' + this.config.Services['echo_server'].Domain + ':' + this.config.Services['echo_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['event_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['event_server'].TLS == true) {
        protocol = 'https';
      }
      this.eventService = new EventServicePromiseClient(
        protocol +
          '://' +
          this.config.Services['event_server'].Domain +
          ':' +
          this.config.Services['event_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['file_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['file_server'].TLS == true) {
        protocol = 'https';
      }
      this.fileService = new FileServicePromiseClient(
        protocol + '://' + this.config.Services['file_server'].Domain + ':' + this.config.Services['file_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['ldap_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['ldap_server'].TLS == true) {
        protocol = 'https';
      }
      this.ldapService = new LdapServicePromiseClient(
        protocol + '://' + this.config.Services['ldap_server'].Domain + ':' + this.config.Services['ldap_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['persistence_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['persistence_server'].TLS == true) {
        protocol = 'https';
      }
      this.persistenceService = new PersistenceServicePromiseClient(
        protocol +
          '://' +
          this.config.Services['persistence_server'].Domain +
          ':' +
          this.config.Services['persistence_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['smtp_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['smtp_server'].TLS == true) {
        protocol = 'https';
      }
      this.smtpService = new SmtpServicePromiseClient(
        protocol + '://' + this.config.Services['smtp_server'].Domain + ':' + this.config.Services['smtp_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['sql_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['sql_server'].TLS == true) {
        protocol = 'https';
      }
      this.sqlService = new SqlServicePromiseClient(
        protocol + '://' + this.config.Services['sql_server'].Domain + ':' + this.config.Services['sql_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['storage_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['storage_server'].TLS == true) {
        protocol = 'https';
      }
      this.storageService = new StorageServicePromiseClient(
        protocol +
          '://' +
          this.config.Services['storage_server'].Domain +
          ':' +
          this.config.Services['storage_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['monitoring_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['monitoring_server'].TLS == true) {
        protocol = 'https';
      }
      this.monitoringService = new MonitoringServicePromiseClient(
        protocol +
          '://' +
          this.config.Services['monitoring_server'].Domain +
          ':' +
          this.config.Services['monitoring_server'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['spc_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['spc_server'].TLS == true) {
        protocol = 'https';
      }
      this.spcService = new SpcServicePromiseClient(
        protocol + '://' + this.config.Services['spc_server'].Domain + ':' + this.config.Services['spc_server'].Proxy,
        null,
        null,
      );
    }

    // non open source services.
    if (this.config.Services['plc_server_ab'] != null) {
      let protocol = 'http';
      if (this.config.Services['plc_server_ab'].TLS == true) {
        protocol = 'https';
      }
      this.plcService_ab = new PlcServicePromiseClient(
        protocol +
          '://' +
          this.config.Services['plc_server_ab'].Domain +
          ':' +
          this.config.Services['plc_server_ab'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['plc_server_siemens'] != null) {
      let protocol = 'http';
      if (this.config.Services['plc_server_siemens'].TLS == true) {
        protocol = 'https';
      }
      this.plcService_siemens = new PlcServicePromiseClient(
        protocol +
          '://' +
          this.config.Services['plc_server_siemens'].Domain +
          ':' +
          this.config.Services['plc_server_siemens'].Proxy,
        null,
        null,
      );
    }
    if (this.config.Services['plc_link_server'] != null) {
      let protocol = 'http';
      if (this.config.Services['plc_link_server'].TLS == true) {
        protocol = 'https';
      }
      this.plcLinkService = new PlcLinkServicePromiseClient(
        protocol +
          '://' +
          this.config.Services['plc_link_server'].Domain +
          ':' +
          this.config.Services['plc_link_server'].Proxy,
        null,
        null,
      );
    }
  }
}
