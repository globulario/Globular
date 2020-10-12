// Here is the list of services from the backend.
import { EventServicePromiseClient } from './event/event_grpc_web_pb';
import { EchoServicePromiseClient } from './echo/echo_grpc_web_pb';
import { CatalogServicePromiseClient } from './catalog/catalog_grpc_web_pb';
import { FileServicePromiseClient } from './file/file_grpc_web_pb';
import { LdapServicePromiseClient } from './ldap/ldap_grpc_web_pb';
import { PersistenceServicePromiseClient } from './persistence/persistence_grpc_web_pb';
import { PlcLinkServicePromiseClient } from './plc_link/plc_link_grpc_web_pb';
import { MailServicePromiseClient } from './mail/mail_grpc_web_pb';
import { SpcServicePromiseClient } from './spc/spc_grpc_web_pb';
import { SqlServicePromiseClient } from './sql/sql_grpc_web_pb';
import { StorageServicePromiseClient } from './storage/storage_grpc_web_pb';
import { MonitoringServicePromiseClient } from './monitoring/monitoring_grpc_web_pb';
import { PlcServicePromiseClient } from './plc/plc_grpc_web_pb';
import { AdminServicePromiseClient } from './admin/admin_grpc_web_pb';
import { RessourceServicePromiseClient } from './ressource/ressource_grpc_web_pb';
import { ServiceDiscoveryPromiseClient, ServiceRepositoryPromiseClient } from './services/services_grpc_web_pb';
import { CertificateAuthorityPromiseClient } from './ca/ca_grpc_web_pb';
import { SubscribeRequest, UnSubscribeRequest, PublishRequest, Event, OnEventRequest, SubscribeResponse } from './event/event_pb';
import { SearchServiceClient, SearchServicePromiseClient } from './search/search_grpc_web_pb';
import { LoadBalancingServiceClient, LoadBalancingServicePromiseClient } from './lb/lb_grpc_web_pb';

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

  // The map of service object.
  Services: IServices;
}

/**
 * Create a "version 4" RFC-4122 UUID (Universal Unique Identifier) string.
 * @returns {string} A string containing the UUID.
 */
function randomUUID(): string {
  const s = new Array();
  const itoh = '0123456789abcdef'; // Make array of random hex digits. The UUID only has 32 digits in it, but we

  // allocate an extra items to make room for the '-'s we'll be inserting.
  for (let i = 0; i < 36; i++) s[i] = Math.floor(Math.random() * 0x10); // Conform to RFC-4122, section 4.4
  s[14] = 4; // Set 4 high bits of time_high field to version
  s[19] = s[19] & 0x3 | 0x8; // Specify 2 high bits of clock sequence
  // Convert to hex chars
  for (let i = 0; i < 36; i++) s[i] = itoh[s[i]]; // Insert '-'s
  s[8] = s[13] = s[18] = s[23] = '-';
  return s.join('');
}

/**
 * That local and distant event hub.
 */
export class EventHub {
  private service: any;
  private subscribers: any;
  private subscriptions: any;
  private uuid: string;

  /**
   * @param {*} service If undefined only local event will be allow.
   */
  constructor(service: any) {
    // The network event bus.
    this.service = service
    // Subscriber function map.
    this.subscribers = {}
    // Subscription name/uuid's maps
    this.subscriptions = {}
    // This is the client uuid
    this.uuid = randomUUID();

    // Open the connection with the server.
    if (this.service !== undefined) {
      // The first step is to subscribe to an event channel.
      const rqst = new OnEventRequest()
      rqst.setUuid(this.uuid)

      const stream = this.service.onEvent(rqst, {});

      // Get the stream and set event on it...
      stream.on('data', (rsp: any) => {
        const evt = rsp.getEvt()
        const data = new TextDecoder("utf-8").decode(evt.getData());
        // dispatch the event localy.
        this.dispatch(evt.getName(), data)
      });

      stream.on('status', (status: any)=> {
        if (status.code === 0) {
          /** Nothing to do here. */
        }
      });

      stream.on('end',  () => {
        // stream end signal
        /** Nothing to do here. */
      });
    }

  }

  /**
   * @param {*} name The name of the event to subcribe to. 
   * @param {*} onsubscribe That function return the uuid of the subscriber.
   * @param {*} onevent That function is call when the event is use.
   */
  subscribe(name: string, onsubscribe: (uuid: string) => any, onevent: (data: any) => any, local: boolean) {
    // Register the local subscriber.
    const uuid = randomUUID()
    if (!local) {
      const rqst = new SubscribeRequest
      rqst.setName(name)
      rqst.setUuid(this.uuid)
      this.service.subscribe(rqst).then((rsp: SubscribeResponse) => {
        if (this.subscribers[name] === undefined) {
          this.subscribers[name] = {}
        }
        this.subscribers[name][uuid] = onevent
        onsubscribe(uuid)
      })
    } else {
      // create a uuid and call onsubscribe callback.
      if (this.subscribers[name] === undefined) {
        this.subscribers[name] = {}
      }
      this.subscribers[name][uuid] = onevent
      onsubscribe(uuid)
    }
  }

  /**
   * 
   * @param {*} name 
   * @param {*} uuid 
   */
  unSubscribe(name: string, uuid: string) {
    if(this.subscribers[name]=== undefined){
      return
    }
    if(this.subscribers[name][uuid]=== undefined){
      return
    }
    // Remove the local subscriber.
    delete this.subscribers[name][uuid]
    if (Object.keys(this.subscribers[name]).length === 0) {
      delete this.subscribers[name]
      // disconnect from the distant server.
      if (this.service !== undefined) {
        const rqst = new UnSubscribeRequest();
        rqst.setName(name);
        rqst.setUuid(this.subscriptions[name])

        // remove the subcription uuid.
        delete this.subscriptions[name]

        // Now I will test with promise
        this.service.unSubscribe(rqst)
          .then((resp: any) => {
            /** Nothing to do here */
          })
          .catch((error: any) => {
            console.log(error)
          })
      }
    }
  }

  /**
   * Publish an event on the bus, or locally in case of local event.
   * @param {*} name The  name of the event to publish
   * @param {*} data The data associated with the event
   * @param {*} local If the event is not local the data must be seraliaze before sent.
   */
  publish(name: string, data: any, local: boolean) {
    if (local === true) {
      this.dispatch(name, data)
    } else {
      // Create a new request.
      const rqst = new PublishRequest();
      const evt = new Event();
      evt.setName(name)

      const enc = new TextEncoder(); // always utf-8
      // encode the string to a array of byte
      evt.setData(enc.encode(data))
      rqst.setEvt(evt);

      // Now I will test with promise
      this.service.publish(rqst)
        .then((resp: any) => {
          /** Nothing to do here. */
        })
        .catch((error: any) => {
          console.log(error)
        })
    }
  }

  /** Dispatch the event localy */
  dispatch(name: string, data: any) {
      for (const uuid in this.subscribers[name]) {
          // call the event callback function.
          if(this.subscribers !== undefined){
              if(this.subscribers[name] !== undefined){
                  if(this.subscribers[name][uuid]!== undefined){
                      this.subscribers[name][uuid](data);
                  }
              }
          }
      }
  }
}

/**
 * Globular regroup all serivces in one object that can be use by
 * application to get access to sql, ldap, persistence... service.
 */
export class Globular {
  config: IConfig | undefined;
  // The admin service.
  adminService: AdminServicePromiseClient | undefined;
  loadBalancingService: LoadBalancingServicePromiseClient | undefined;
  ressourceService: RessourceServicePromiseClient | undefined;
  servicesDicovery: ServiceDiscoveryPromiseClient | undefined;
  servicesRepository: ServiceRepositoryPromiseClient | undefined;
  certificateAuthority: CertificateAuthorityPromiseClient | undefined;

  // list of services.
  catalogService: CatalogServicePromiseClient | undefined;
  echoService: EchoServicePromiseClient | undefined;
  eventService: EventServicePromiseClient | undefined;
  fileService: FileServicePromiseClient | undefined;
  ldapService: LdapServicePromiseClient | undefined;
  persistenceService: PersistenceServicePromiseClient | undefined;
  mailService: MailServicePromiseClient | undefined;
  sqlService: SqlServicePromiseClient | undefined;
  storageService: StorageServicePromiseClient | undefined;
  monitoringService: MonitoringServicePromiseClient | undefined;
  spcService: SpcServicePromiseClient | undefined;
  searchService: SearchServicePromiseClient | undefined;

  // Non open source services.
  plcService_ab: PlcServicePromiseClient | undefined;
  plcService_siemens: PlcServicePromiseClient | undefined;
  plcLinkService: PlcLinkServicePromiseClient | undefined;

  /** The configuation. */
  constructor(config: IConfig) {
    // Keep the config...
    this.config = config;

    /** The admin service to access to other configurations. */
    this.adminService = new AdminServicePromiseClient(
      this.config.Protocol + '://' + this.config.Domain + ':' + this.config.AdminProxy,
      null,
      null,
    );

    /** That service is use to control acces to ressource like method access and account. */
    this.ressourceService = new RessourceServicePromiseClient(
      this.config.Protocol + '://' + this.config.Domain + ':' + this.config.RessourceProxy,
      null,
      null,
    );

    this.loadBalancingService = new LoadBalancingServicePromiseClient(
      this.config.Protocol + '://' + this.config.Domain + ':' + this.config.LoadBalancingServiceProxy,
      null,
      null,
    );

    /** That service help to find and install or publish new service on the backend. */
    this.servicesDicovery = new ServiceDiscoveryPromiseClient(
      this.config.Protocol + '://' + this.config.Domain + ':' + this.config.ServicesDiscoveryProxy,
      null,
      null,
    );

    /** Functionality to use service repository server. */
    this.servicesRepository = new ServiceRepositoryPromiseClient(
      this.config.Protocol + '://' + this.config.Domain + ':' + this.config.ServicesRepositoryProxy,
      null,
      null,
    );

    /** Certificate authority function needed for TLS client. */
    this.certificateAuthority = new CertificateAuthorityPromiseClient(
      this.config.Protocol + '://' + this.config.Domain + ':' + this.config.CertificateAuthorityProxy,
      null,
      null,
    );

    // Iinitialisation of services.

    // The catalog server
    let catalog_server = this.getFirstConfigByName('catalog.CatalogService')
    if (catalog_server != undefined) {
      let protocol = 'http';
      if (catalog_server.TLS == true) {
        protocol = 'https';
      }
      this.catalogService = new CatalogServicePromiseClient(
        protocol +
        '://' +
        catalog_server.Domain +
        ':' +
        catalog_server.Proxy,
        null,
        null,
      );
    }

    // The echo server
    let echo_server = this.getFirstConfigByName('echo.EchoService')
    if (echo_server != null) {
      let protocol = 'http';
      if (echo_server.TLS == true) {
        protocol = 'https';
      }
      this.echoService = new EchoServicePromiseClient(
        protocol + '://' + echo_server.Domain + ':' + echo_server.Proxy,
        null,
        null,
      );
    }

    // The search service
    let search_server = this.getFirstConfigByName('search.SearchService')
    if (search_server != null) {
      let protocol = 'http';
      if (search_server.TLS == true) {
        protocol = 'https';
      }
      this.searchService = new SearchServicePromiseClient(
        protocol + '://' + search_server.Domain + ':' + search_server.Proxy,
        null,
        null,
      );
    }

    // The event server.
    let event_server = this.getFirstConfigByName('event.EventService')
    if (event_server != null) {
      let protocol = 'http';
      if (event_server.TLS == true) {
        protocol = 'https';
      }
      this.eventService = new EventServicePromiseClient(
        protocol +
        '://' +
        event_server.Domain +
        ':' +
        event_server.Proxy,
        null,
        null,
      );
    }

    // The file server.
    let file_server = this.getFirstConfigByName('file.FileService')
    if (file_server != null) {
      let protocol = 'http';
      if (file_server.TLS == true) {
        protocol = 'https';
      }
      this.fileService = new FileServicePromiseClient(
        protocol + '://' + file_server.Domain + ':' + file_server.Proxy,
        null,
        null,
      );
    }

    // The ldap server
    let ldap_server = this.getFirstConfigByName('ldap.LdapService')
    if (ldap_server != null) {
      let protocol = 'http';
      if (ldap_server.TLS == true) {
        protocol = 'https';
      }
      this.ldapService = new LdapServicePromiseClient(
        protocol + '://' + ldap_server.Domain + ':' + ldap_server.Proxy,
        null,
        null,
      );
    }

    // The persistence server.
    let persistence_server = this.getFirstConfigByName('persistence.PersistenceService')
    if (persistence_server != null) {
      let protocol = 'http';
      if (persistence_server.TLS == true) {
        protocol = 'https';
      }
      this.persistenceService = new PersistenceServicePromiseClient(
        protocol +
        '://' +
        persistence_server.Domain +
        ':' +
        persistence_server.Proxy,
        null,
        null,
      );
    }

    // The mail server
    let mail_server = this.getFirstConfigByName('mail.MailService')

    if (mail_server != null) {
      let protocol = 'http';
      if (mail_server.TLS == true) {
        protocol = 'https';
      }
      this.mailService = new MailServicePromiseClient(
        protocol + '://' + mail_server.Domain + ':' + mail_server.Proxy,
        null,
        null,
      );
    }

    // The sql service.
    let sql_server = this.getFirstConfigByName('sql.SqlService')
    if (sql_server != null) {
      let protocol = 'http';
      if (sql_server.TLS == true) {
        protocol = 'https';
      }
      this.sqlService = new SqlServicePromiseClient(
        protocol + '://' + sql_server.Domain + ':' + sql_server.Proxy,
        null,
        null,
      );
    }

    // The storage service.
    let storage_server = this.getFirstConfigByName('storage.StorageService')
    if (storage_server != null) {
      let protocol = 'http';
      if (storage_server.TLS == true) {
        protocol = 'https';
      }
      this.storageService = new StorageServicePromiseClient(
        protocol +
        '://' +
        storage_server.Domain +
        ':' +
        storage_server.Proxy,
        null,
        null,
      );
    }

    // The monitoring service.
    let monitoring_server = this.getFirstConfigByName('monitoring.MonitoringService')

    if (monitoring_server != null) {
      let protocol = 'http';
      if (monitoring_server.TLS == true) {
        protocol = 'https';
      }
      this.monitoringService = new MonitoringServicePromiseClient(
        protocol +
        '://' +
        monitoring_server.Domain +
        ':' +
        monitoring_server.Proxy,
        null,
        null,
      );
    }

    // The spc server.
    let spc_server = this.getFirstConfigByName('spc.SpcService')
    if (spc_server != null) {
      let protocol = 'http';
      if (spc_server.TLS == true) {
        protocol = 'https';
      }
      this.spcService = new SpcServicePromiseClient(
        protocol + '://' + spc_server.Domain + ':' + spc_server.Proxy,
        null,
        null,
      );
    }

    // TODO Here I got tow implementation of the same service.
    // So I will use it Path instead of name...
    /*
    // The plc_server_ab
    let plc_server_ab = this.getFirstConfigByName('plc.PlcService')
    if (plc_server_ab != null) {
      let protocol = 'http';
      if (plc_server_ab.TLS == true) {
        protocol = 'https';
      }
      this.plcService_ab = new PlcServicePromiseClient(
        protocol +
        '://' +
        plc_server_ab.Domain +
        ':' +
        plc_server_ab.Proxy,
        null,
        null,
      );
    }

    // PLC simmens
    let plc_server_siemens = this.getFirstConfigByName('plc.PlcService')
    if (plc_server_siemens != null) {
      let protocol = 'http';
      if (plc_server_siemens.TLS == true) {
        protocol = 'https';
      }
      this.plcService_siemens = new PlcServicePromiseClient(
        protocol +
        '://' +
        plc_server_siemens.Domain +
        ':' +
        plc_server_siemens.Proxy,
        null,
        null,
      );
    }
    */

    // PLC Link server.
    let plc_link_server = this.getFirstConfigByName('plc_link.PlcLinkService')
    if (plc_link_server != null) {
      let protocol = 'http';
      if (plc_link_server.TLS == true) {
        protocol = 'https';
      }
      this.plcLinkService = new PlcLinkServicePromiseClient(
        protocol +
        '://' +
        plc_link_server.Domain +
        ':' +
        plc_link_server.Proxy,
        null,
        null,
      );
    }
  }

  // Return the first configuration that match the given name.
  // The load balancer will be in charge to select the correct service instance from the list
  // The first instance is the entry point of the services.
  getFirstConfigByName(name:string): IServiceConfig {
    for(const id in this.config.Services){
      const service = this.config.Services[id]
      if(service.Name == name){
        return service;
      }
    }
    return null;
  }
}
