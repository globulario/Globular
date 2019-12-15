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
import { ServiceDiscoveryPromiseClient, ServiceRepositoryPromiseClient } from './services/services_grpc_web_pb';
import { CertificateAuthorityPromiseClient } from './ca/ca_grpc_web_pb';
import { SubscribeRequest, UnSubscribeRequest, PublishRequest, Event } from './event/eventpb/event_pb';

/**
 * The service configuration information.
 */
export interface IServiceConfig {
  Id: string;
  Name: string;
  State: string;
  Domain: string;
  Port: Number;
  Proxy: Number;
  TLS: Boolean;
  KeepUpToDate: Boolean;
  KeepAlive: Boolean;
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
  CertURL:string;
  IdleTimeout: number;
  
  // The map of service object.
  Services: IServices;
}

/**
 * Create a "version 4" RFC-4122 UUID (Universal Unique Identifier) string.
 * @returns {string} A string containing the UUID.
 */
function randomUUID(): string {
    var s = new Array();
    var itoh = '0123456789abcdef'; // Make array of random hex digits. The UUID only has 32 digits in it, but we

    // allocate an extra items to make room for the '-'s we'll be inserting.
    for (var i = 0; i < 36; i++) s[i] = Math.floor(Math.random() * 0x10); // Conform to RFC-4122, section 4.4
    s[14] = 4; // Set 4 high bits of time_high field to version
    s[19] = s[19] & 0x3 | 0x8; // Specify 2 high bits of clock sequence
    // Convert to hex chars
    for (var i = 0; i < 36; i++) s[i] = itoh[s[i]]; // Insert '-'s
    s[8] = s[13] = s[18] = s[23] = '-';
    return s.join('');
}

/**
 * That local and distant event hub.
 */
export class EventHub {
    readonly service: any;
    readonly subscribers: any;
    readonly subscriptions: any;

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
    }

    /**
     * @param {*} name The name of the event to subcribe to. 
     * @param {*} onsubscribe That function return the uuid of the subscriber.
     * @param {*} onevent That function is call when the event is use.
     */
    subscribe(name: string, onsubscribe: (uuid: string) => any, onevent: (data: any) => any) {
        // Register the local subscriber.
        var uuid = randomUUID()
        if (this.subscribers[name] == undefined) {
            this.subscribers[name] = {}
            if (this.service != undefined) {
                // The first step is to subscribe to an event channel.
                var rqst = new SubscribeRequest()
                rqst.setName(name)
                if (this.service != null) {
                    var stream = this.service.subscribe(rqst, {});
                    // Get the stream and set event on it...
                    stream.on('data', function (hub, name) {
                        return function (rsp: any) {
                            if (rsp.hasUuid()) {
                                hub.subscriptions[name] = rsp.getUuid()
                            } else if (rsp.hasEvt()) {
                                var evt = rsp.getEvt()
                                var data = new TextDecoder("utf-8").decode(evt.getData());
                                // dispatch the event localy.
                                hub.dispatch(name, data)
                            }
                        }
                    }(this, name));

                    stream.on('status', function (status: any) {
                        if (status.code == 0) {
                            /** Nothing to do here. */
                        }
                    });

                    stream.on('end', function () {
                        // stream end signal
                        /** Nothing to do here. */
                    });
                }
            }
        }

        // Set the event callback function.
        this.subscribers[name][uuid] = onevent

        // call on subscribe call back.
        onsubscribe(uuid)
    }

    /**
     * 
     * @param {*} name 
     * @param {*} uuid 
     */
    unSubscribe(name: string, uuid: string) {
        // Remove the local subscriber.
        delete this.subscribers[name][uuid]
        if (Object.keys(this.subscribers[name]).length == 0) {
            delete this.subscribers[name]
            // disconnect from the distant server.
            if (this.service != undefined) {
                var request = new UnSubscribeRequest();
                request.setName(name);
                request.setUuid(this.subscriptions[name])

                // remove the subcription uuid.
                delete this.subscriptions[name]

                // Now I will test with promise
                this.service.unSubscribe(request)
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
        if (this.service == undefined || local) {
            this.dispatch(name, data)
        } else {
            // Create a new request.
            var request = new PublishRequest();
            var evt = new Event();
            evt.setName(name)

            var enc = new TextEncoder(); // always utf-8
            // encode the string to a array of byte
            evt.setData(enc.encode(data))
            request.setEvt(evt);

            // Now I will test with promise
            this.service.publish(request)
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
        for (var uuid in this.subscribers[name]) {
            // call the event callback function.
            this.subscribers[name][uuid](data)
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
  smtpService: SmtpServicePromiseClient | undefined;
  sqlService: SqlServicePromiseClient | undefined;
  storageService: StorageServicePromiseClient | undefined;
  monitoringService: MonitoringServicePromiseClient | undefined;
  spcService: SpcServicePromiseClient | undefined;

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
