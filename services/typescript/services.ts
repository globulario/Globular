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
import { GetThumbnailsRequest } from './file/file_pb';

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
  Description: string;
  Keywords: Array<string>;
  Discoveries: Array<string>;
  Repositories: Array<string>;
  Proto: String;
  Path: string;
  CertAuthorityTrust: string;
  CertFile: string;
  KeyFile: string;
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

      stream.on('status', (status: any) => {
        if (status.code === 0) {
          /** Nothing to do here. */
        }
      });

      stream.on('end', () => {
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
    if (this.subscribers[name] === undefined) {
      return
    }
    if (this.subscribers[name][uuid] === undefined) {
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
      if (this.subscribers !== undefined) {
        if (this.subscribers[name] !== undefined) {
          if (this.subscribers[name][uuid] !== undefined) {
            this.subscribers[name][uuid](data);
          }
        }
      }
    }
  }
}

// Get the configuration from url
function getFileConfig(url: string, callback: (obj: any) => void, errorcallback: (err: any) => void) {
  var xmlhttp = new XMLHttpRequest();

  xmlhttp.onreadystatechange = function () {
    if (this.readyState == 4 && this.status == 201) {
      var obj = JSON.parse(this.responseText);
      callback(obj);
    }else if (this.readyState == 4){
      errorcallback("fail to get the configuration file at url " + url + " status " + this.status)
    }
  };

  xmlhttp.open("GET", url, true);
  xmlhttp.send();
}

/**
 * Globular regroup all serivces in one object that can be use by
 * application to get access to sql, ldap, persistence... service.
 */
export class Globular {


  /** The configuation. */
  constructor(url: string, callback: () => void, errorcallback: (err: any) => void) {
    // Keep the config...
    getFileConfig(url, (config:any)=>{
      this.config = config;
      callback();
    }, errorcallback)
  }

  private _config: IConfig;
  public get config(): IConfig {
    return this._config;
  }
  public set config(value: IConfig) {
    this._config = value;
  }

  /** The admin service to access to other configurations. */
  private _adminService: AdminServicePromiseClient
  public get adminService(): AdminServicePromiseClient | undefined {
    // refresh the config.
    if (this._adminService == null) {
      this._adminService = new AdminServicePromiseClient(
        this.config.Protocol + '://' + this.config.Domain + ':' + this.config.AdminProxy,
        null,
        null,
      )
    }
    return this._adminService;
  }

  private _loadBalancingService: LoadBalancingServicePromiseClient
  public get loadBalancingService(): LoadBalancingServicePromiseClient | undefined {
    // refresh the config.
    if (this._loadBalancingService == null) {
      this._loadBalancingService = new LoadBalancingServicePromiseClient(
        this.config.Protocol + '://' + this.config.Domain + ':' + this.config.LoadBalancingServiceProxy,
        null,
        null,
      );
    }
    return this._loadBalancingService;
  }

  private _ressourceService: RessourceServicePromiseClient
  public get ressourceService(): RessourceServicePromiseClient | undefined {
    // refresh the config.
    if (this._ressourceService == null) {
      this._ressourceService = new RessourceServicePromiseClient(
        this.config.Protocol + '://' + this.config.Domain + ':' + this.config.RessourceProxy,
        null,
        null,
      );
    }
    return this._ressourceService;
  }

  private _servicesDicovery: ServiceDiscoveryPromiseClient
  public get servicesDicovery(): ServiceDiscoveryPromiseClient | undefined {
    // refresh the config.
    if (this._servicesDicovery == null) {
      this._servicesDicovery = new ServiceDiscoveryPromiseClient(
        this.config.Protocol + '://' + this.config.Domain + ':' + this.config.ServicesDiscoveryProxy,
        null,
        null,
      );
    }
    return this._servicesDicovery;
  }

  private _servicesRepository: ServiceRepositoryPromiseClient
  public get servicesRepository(): ServiceRepositoryPromiseClient | undefined {
    // refresh the config.
    if (this._servicesRepository == null) {
      this._servicesRepository = new ServiceRepositoryPromiseClient(
        this.config.Protocol + '://' + this.config.Domain + ':' + this.config.ServicesRepositoryProxy,
        null,
        null,
      );
    }

    return this._servicesRepository;
  }

  private _certificateAuthority: CertificateAuthorityPromiseClient;
  public get certificateAuthority(): CertificateAuthorityPromiseClient | undefined {
    // refresh the config.
    if (this._certificateAuthority == null) {
      this._certificateAuthority = new CertificateAuthorityPromiseClient(
        this.config.Protocol + '://' + this.config.Domain + ':' + this.config.CertificateAuthorityProxy,
        null,
        null,
      );
    }
    return this._certificateAuthority;
  }

  // list of services.
  private _catalogService: CatalogServicePromiseClient
  public get catalogService(): CatalogServicePromiseClient | undefined {
    if (this._catalogService != null) {
      return this._catalogService
    }
    let config = this.getFirstConfigByName('catalog.CatalogService')
    if (config != undefined) {
      if (this._catalogService == null) {
        this._catalogService = new CatalogServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._catalogService;
    }
    return undefined;
  }

  private _echoService: EchoServicePromiseClient
  public get echoService(): EchoServicePromiseClient | undefined {
    if (this._echoService != null) {
      return this._echoService
    }
    let config = this.getFirstConfigByName('echo.EchoService')
    if (config != undefined) {
      if (this._echoService == null) {
        this._echoService = new EchoServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._echoService;
    }
    return undefined;
  }

  private _eventService: EventServicePromiseClient
  public get eventService(): EventServicePromiseClient | undefined {
    if (this._eventService != null) {
      return this._eventService
    }
    let config = this.getFirstConfigByName('event.EventService')
    if (config != undefined) {
      if (this._eventService == null) {
        this._eventService = new EventServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._eventService;
    }
    return undefined;
  }

  private _fileService: FileServicePromiseClient
  public get fileService(): FileServicePromiseClient | undefined {
    if (this._fileService != null) {
      return this._fileService
    }
    let config = this.getFirstConfigByName('file.FileService')
    if (config != undefined) {
      if (this._fileService == null) {
        this._fileService = new FileServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._fileService;
    }
    return undefined;
  }

  private _ldapService: LdapServicePromiseClient
  public get ldapService(): LdapServicePromiseClient | undefined {
    if (this._ldapService != null) {
      return this._ldapService
    }
    let config = this.getFirstConfigByName('ldap.LdapService')
    if (config != undefined) {
      if (this._ldapService == null) {
        this._ldapService = new LdapServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._ldapService;
    }
    return undefined;
  }

  private _persistenceService: PersistenceServicePromiseClient
  public get persistenceService(): PersistenceServicePromiseClient | undefined {
    if (this._persistenceService != null) {
      return this._persistenceService
    }
    let config = this.getFirstConfigByName('persistence.PersistenceService')
    if (config != undefined) {
      if (this._persistenceService == null) {
        this._persistenceService = new PersistenceServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._persistenceService;
    }
    return undefined;
  }

  private _mailService: MailServicePromiseClient
  public get mailService(): MailServicePromiseClient | undefined {
    if (this._mailService != null) {
      return this._mailService
    }
    let config = this.getFirstConfigByName('mail.MailService')
    if (config != undefined) {
      if (this._mailService == null) {
        this._mailService = new MailServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._mailService;
    }
    return undefined;
  }

  private _sqlService: SqlServicePromiseClient
  public get sqlService(): SqlServicePromiseClient | undefined {
    if (this._sqlService != null) {
      return this._sqlService
    }
    let config = this.getFirstConfigByName('sql.SqlService')
    if (config != undefined) {
      if (this._sqlService == null) {
        this._sqlService = new SqlServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._sqlService;
    }
    return undefined;
  }

  private _storageService: StorageServicePromiseClient
  public get storageService(): StorageServicePromiseClient | undefined {
    if (this._storageService != null) {
      return this._storageService
    }
    let config = this.getFirstConfigByName('storage.StorageService')
    if (config != undefined) {
      if (this._storageService == null) {
        this._storageService = new StorageServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._storageService;
    }
    return undefined;
  }

  private _monitoringService: MonitoringServicePromiseClient
  public get monitoringService(): MonitoringServicePromiseClient | undefined {
    if (this._monitoringService != null) {
      return this._monitoringService
    }
    let config = this.getFirstConfigByName('monitoring.MonitoringService')
    if (config != undefined) {
      if (this._monitoringService == null) {
        this._monitoringService = new MonitoringServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._monitoringService;
    }
    return undefined;
  }

  private _spcService: SpcServicePromiseClient
  public get spcService(): SpcServicePromiseClient | undefined {
    if (this._spcService != null) {
      return this._spcService
    }
    let config = this.getFirstConfigByName('spc.SpcService')
    if (config != undefined) {
      if (this._spcService == null) {
        this._spcService = new SpcServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._spcService;
    }
    return undefined;
  }

  private _searchService: SearchServicePromiseClient
  public get searchService(): SearchServicePromiseClient | undefined {
    if (this._searchService != null) {
      return this.searchService
    }
    let config = this.getFirstConfigByName('search.SearchService')
    if (config != undefined) {
      if (this._searchService == null) {
        this._searchService = new SearchServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._searchService;
    }
    return undefined;
  }

  // Non open source services.
  private _plcService_ab: PlcServicePromiseClient
  public get plcService_ab(): PlcServicePromiseClient | undefined {
    let config = this.getFirstConfigByName('plc.PlcService')
    if (config != undefined) {
      if (this._plcService_ab == null) {
        this._plcService_ab = new PlcServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._plcService_ab;
    }
    return undefined;
  }

  private _plcService_siemens: PlcServicePromiseClient
  public get plcService_siemens(): PlcServicePromiseClient | undefined {
    let config = this.getFirstConfigByName('plc.PlcService')
    if (config != undefined) {
      if (this._plcService_siemens == null) {
        this._plcService_siemens = new PlcServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._plcService_siemens;
    }
    return undefined;
  }

  private _plcLinkService: PlcLinkServicePromiseClient
  public get plcLinkService(): PlcLinkServicePromiseClient | undefined {
    let config = this.getFirstConfigByName('plc_link.PlcLinkService')
    if (config != undefined) {
      if (this._plcLinkService == null) {
        this._plcLinkService = new PlcLinkServicePromiseClient(
          this.config.Protocol +
          '://' +
          config.Domain +
          ':' +
          config.Proxy,
          null,
          null,
        );
      }
      return this._plcLinkService;
    }
    return undefined;
  }

  // Return the first configuration that match the given name.
  // The load balancer will be in charge to select the correct service instance from the list
  // The first instance is the entry point of the services.
  getFirstConfigByName(name: string): IServiceConfig {
    for (const id in this.config.Services) {
      const service = this.config.Services[id]
      if (service.Name == name) {
        return service;
      }
    }
    return null;
  }
}
