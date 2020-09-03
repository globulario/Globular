"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// Here is the list of services from the backend.
var event_grpc_web_pb_1 = require("./event/eventpb/event_grpc_web_pb");
var echo_grpc_web_pb_1 = require("./echo/echopb/echo_grpc_web_pb");
var catalog_grpc_web_pb_1 = require("./catalog/catalogpb/catalog_grpc_web_pb");
var file_grpc_web_pb_1 = require("./file/filepb/file_grpc_web_pb");
var ldap_grpc_web_pb_1 = require("./ldap/ldappb/ldap_grpc_web_pb");
var persistence_grpc_web_pb_1 = require("./persistence/persistencepb/persistence_grpc_web_pb");
var plc_link_grpc_web_pb_1 = require("./plc_link/plc_link_pb/plc_link_grpc_web_pb");
var smtp_grpc_web_pb_1 = require("./smtp/smtppb/smtp_grpc_web_pb");
var spc_grpc_web_pb_1 = require("./spc/spcpb/spc_grpc_web_pb");
var sql_grpc_web_pb_1 = require("./sql/sqlpb/sql_grpc_web_pb");
var storage_grpc_web_pb_1 = require("./storage/storagepb/storage_grpc_web_pb");
var monitoring_grpc_web_pb_1 = require("./monitoring/monitoringpb/monitoring_grpc_web_pb");
var admin_grpc_web_pb_1 = require("./admin/admin_grpc_web_pb");
var ressource_grpc_web_pb_1 = require("./ressource/ressource_grpc_web_pb");
var services_grpc_web_pb_1 = require("./services/services_grpc_web_pb");
var ca_grpc_web_pb_1 = require("./ca/ca_grpc_web_pb");
var event_pb_1 = require("./event/eventpb/event_pb");
var search_grpc_web_pb_1 = require("./search/searchpb/search_grpc_web_pb");
var lb_grpc_web_pb_1 = require("./lb/lb_grpc_web_pb");
/**
 * Create a "version 4" RFC-4122 UUID (Universal Unique Identifier) string.
 * @returns {string} A string containing the UUID.
 */
function randomUUID() {
    var s = new Array();
    var itoh = '0123456789abcdef'; // Make array of random hex digits. The UUID only has 32 digits in it, but we
    // allocate an extra items to make room for the '-'s we'll be inserting.
    for (var i = 0; i < 36; i++)
        s[i] = Math.floor(Math.random() * 0x10); // Conform to RFC-4122, section 4.4
    s[14] = 4; // Set 4 high bits of time_high field to version
    s[19] = s[19] & 0x3 | 0x8; // Specify 2 high bits of clock sequence
    // Convert to hex chars
    for (var i = 0; i < 36; i++)
        s[i] = itoh[s[i]]; // Insert '-'s
    s[8] = s[13] = s[18] = s[23] = '-';
    return s.join('');
}
/**
 * That local and distant event hub.
 */
var EventHub = /** @class */ (function () {
    /**
     * @param {*} service If undefined only local event will be allow.
     */
    function EventHub(service) {
        var _this = this;
        // The network event bus.
        this.service = service;
        // Subscriber function map.
        this.subscribers = {};
        // Subscription name/uuid's maps
        this.subscriptions = {};
        // This is the client uuid
        this.uuid = randomUUID();
        // Open the connection with the server.
        if (this.service !== undefined) {
            // The first step is to subscribe to an event channel.
            var rqst = new event_pb_1.OnEventRequest();
            rqst.setUuid(this.uuid);
            var stream = this.service.onEvent(rqst, {});
            // Get the stream and set event on it...
            stream.on('data', function (rsp) {
                var evt = rsp.getEvt();
                var data = new TextDecoder("utf-8").decode(evt.getData());
                // dispatch the event localy.
                _this.dispatch(evt.getName(), data);
            });
            stream.on('status', function (status) {
                if (status.code === 0) {
                    /** Nothing to do here. */
                }
            });
            stream.on('end', function () {
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
    EventHub.prototype.subscribe = function (name, onsubscribe, onevent, local) {
        var _this = this;
        // Register the local subscriber.
        var uuid = randomUUID();
        if (!local) {
            var rqst = new event_pb_1.SubscribeRequest;
            rqst.setName(name);
            rqst.setUuid(this.uuid);
            this.service.subscribe(rqst).then(function (rsp) {
                if (_this.subscribers[name] === undefined) {
                    _this.subscribers[name] = {};
                }
                _this.subscribers[name][uuid] = onevent;
                onsubscribe(uuid);
            });
        }
        else {
            // create a uuid and call onsubscribe callback.
            if (this.subscribers[name] === undefined) {
                this.subscribers[name] = {};
            }
            this.subscribers[name][uuid] = onevent;
            onsubscribe(uuid);
        }
    };
    /**
     *
     * @param {*} name
     * @param {*} uuid
     */
    EventHub.prototype.unSubscribe = function (name, uuid) {
        if (this.subscribers[name] === undefined) {
            return;
        }
        if (this.subscribers[name][uuid] === undefined) {
            return;
        }
        // Remove the local subscriber.
        delete this.subscribers[name][uuid];
        if (Object.keys(this.subscribers[name]).length === 0) {
            delete this.subscribers[name];
            // disconnect from the distant server.
            if (this.service !== undefined) {
                var rqst = new event_pb_1.UnSubscribeRequest();
                rqst.setName(name);
                rqst.setUuid(this.subscriptions[name]);
                // remove the subcription uuid.
                delete this.subscriptions[name];
                // Now I will test with promise
                this.service.unSubscribe(rqst)
                    .then(function (resp) {
                    /** Nothing to do here */
                })
                    .catch(function (error) {
                    console.log(error);
                });
            }
        }
    };
    /**
     * Publish an event on the bus, or locally in case of local event.
     * @param {*} name The  name of the event to publish
     * @param {*} data The data associated with the event
     * @param {*} local If the event is not local the data must be seraliaze before sent.
     */
    EventHub.prototype.publish = function (name, data, local) {
        if (local === true) {
            this.dispatch(name, data);
        }
        else {
            // Create a new request.
            var rqst = new event_pb_1.PublishRequest();
            var evt = new event_pb_1.Event();
            evt.setName(name);
            var enc = new TextEncoder(); // always utf-8
            // encode the string to a array of byte
            evt.setData(enc.encode(data));
            rqst.setEvt(evt);
            // Now I will test with promise
            this.service.publish(rqst)
                .then(function (resp) {
                /** Nothing to do here. */
            })
                .catch(function (error) {
                console.log(error);
            });
        }
    };
    /** Dispatch the event localy */
    EventHub.prototype.dispatch = function (name, data) {
        for (var uuid in this.subscribers[name]) {
            // call the event callback function.
            if (this.subscribers !== undefined) {
                if (this.subscribers[name] !== undefined) {
                    if (this.subscribers[name][uuid] !== undefined) {
                        this.subscribers[name][uuid](data);
                    }
                }
            }
        }
    };
    return EventHub;
}());
exports.EventHub = EventHub;
/**
 * Globular regroup all serivces in one object that can be use by
 * application to get access to sql, ldap, persistence... service.
 */
var Globular = /** @class */ (function () {
    /** The configuation. */
    function Globular(config) {
        // Keep the config...
        this.config = config;
        /** The admin service to access to other configurations. */
        this.adminService = new admin_grpc_web_pb_1.AdminServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.AdminProxy, null, null);
        /** That service is use to control acces to ressource like method access and account. */
        this.ressourceService = new ressource_grpc_web_pb_1.RessourceServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.RessourceProxy, null, null);
        this.loadBalancingService = new lb_grpc_web_pb_1.LoadBalancingServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.LoadBalancingServiceProxy, null, null);
        /** That service help to find and install or publish new service on the backend. */
        this.servicesDicovery = new services_grpc_web_pb_1.ServiceDiscoveryPromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.ServicesDiscoveryProxy, null, null);
        /** Functionality to use service repository server. */
        this.servicesRepository = new services_grpc_web_pb_1.ServiceRepositoryPromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.ServicesRepositoryProxy, null, null);
        /** Certificate authority function needed for TLS client. */
        this.certificateAuthority = new ca_grpc_web_pb_1.CertificateAuthorityPromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.CertificateAuthorityProxy, null, null);
        // Iinitialisation of services.
        // The catalog server
        var catalog_server = this.getFirstConfigByName('catalog.CatalogService');
        if (catalog_server != undefined) {
            var protocol = 'http';
            if (catalog_server.TLS == true) {
                protocol = 'https';
            }
            this.catalogService = new catalog_grpc_web_pb_1.CatalogServicePromiseClient(protocol +
                '://' +
                catalog_server.Domain +
                ':' +
                catalog_server.Proxy, null, null);
        }
        // The echo server
        var echo_server = this.getFirstConfigByName('echo.EchoService');
        if (echo_server != null) {
            var protocol = 'http';
            if (echo_server.TLS == true) {
                protocol = 'https';
            }
            this.echoService = new echo_grpc_web_pb_1.EchoServicePromiseClient(protocol + '://' + echo_server.Domain + ':' + echo_server.Proxy, null, null);
        }
        // The search service
        var search_server = this.getFirstConfigByName('search.SearchService');
        if (search_server != null) {
            var protocol = 'http';
            if (search_server.TLS == true) {
                protocol = 'https';
            }
            this.searchService = new search_grpc_web_pb_1.SearchServicePromiseClient(protocol + '://' + search_server.Domain + ':' + search_server.Proxy, null, null);
        }
        // The event server.
        var event_server = this.getFirstConfigByName('event.EventService');
        if (event_server != null) {
            var protocol = 'http';
            if (event_server.TLS == true) {
                protocol = 'https';
            }
            this.eventService = new event_grpc_web_pb_1.EventServicePromiseClient(protocol +
                '://' +
                event_server.Domain +
                ':' +
                event_server.Proxy, null, null);
        }
        // The file server.
        var file_server = this.getFirstConfigByName('file.FileService');
        if (file_server != null) {
            var protocol = 'http';
            if (file_server.TLS == true) {
                protocol = 'https';
            }
            this.fileService = new file_grpc_web_pb_1.FileServicePromiseClient(protocol + '://' + file_server.Domain + ':' + file_server.Proxy, null, null);
        }
        // The ldap server
        var ldap_server = this.getFirstConfigByName('ldap.LdapService');
        if (ldap_server != null) {
            var protocol = 'http';
            if (ldap_server.TLS == true) {
                protocol = 'https';
            }
            this.ldapService = new ldap_grpc_web_pb_1.LdapServicePromiseClient(protocol + '://' + ldap_server.Domain + ':' + ldap_server.Proxy, null, null);
        }
        // The persistence server.
        var persistence_server = this.getFirstConfigByName('persistence.PersistenceService');
        if (persistence_server != null) {
            var protocol = 'http';
            if (persistence_server.TLS == true) {
                protocol = 'https';
            }
            this.persistenceService = new persistence_grpc_web_pb_1.PersistenceServicePromiseClient(protocol +
                '://' +
                persistence_server.Domain +
                ':' +
                persistence_server.Proxy, null, null);
        }
        // The smtp server
        var smtp_server = this.getFirstConfigByName('smtp.SmtpService');
        if (smtp_server != null) {
            var protocol = 'http';
            if (smtp_server.TLS == true) {
                protocol = 'https';
            }
            this.smtpService = new smtp_grpc_web_pb_1.SmtpServicePromiseClient(protocol + '://' + smtp_server.Domain + ':' + smtp_server.Proxy, null, null);
        }
        // The sql service.
        var sql_server = this.getFirstConfigByName('sql.SqlService');
        if (sql_server != null) {
            var protocol = 'http';
            if (sql_server.TLS == true) {
                protocol = 'https';
            }
            this.sqlService = new sql_grpc_web_pb_1.SqlServicePromiseClient(protocol + '://' + sql_server.Domain + ':' + sql_server.Proxy, null, null);
        }
        // The storage service.
        var storage_server = this.getFirstConfigByName('storage.StorageService');
        if (storage_server != null) {
            var protocol = 'http';
            if (storage_server.TLS == true) {
                protocol = 'https';
            }
            this.storageService = new storage_grpc_web_pb_1.StorageServicePromiseClient(protocol +
                '://' +
                storage_server.Domain +
                ':' +
                storage_server.Proxy, null, null);
        }
        // The monitoring service.
        var monitoring_server = this.getFirstConfigByName('monitoring.MonitoringService');
        if (monitoring_server != null) {
            var protocol = 'http';
            if (monitoring_server.TLS == true) {
                protocol = 'https';
            }
            this.monitoringService = new monitoring_grpc_web_pb_1.MonitoringServicePromiseClient(protocol +
                '://' +
                monitoring_server.Domain +
                ':' +
                monitoring_server.Proxy, null, null);
        }
        // The spc server.
        var spc_server = this.getFirstConfigByName('spc.SpcService');
        if (spc_server != null) {
            var protocol = 'http';
            if (spc_server.TLS == true) {
                protocol = 'https';
            }
            this.spcService = new spc_grpc_web_pb_1.SpcServicePromiseClient(protocol + '://' + spc_server.Domain + ':' + spc_server.Proxy, null, null);
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
        var plc_link_server = this.getFirstConfigByName('plc_link.PlcLinkService');
        if (plc_link_server != null) {
            var protocol = 'http';
            if (plc_link_server.TLS == true) {
                protocol = 'https';
            }
            this.plcLinkService = new plc_link_grpc_web_pb_1.PlcLinkServicePromiseClient(protocol +
                '://' +
                plc_link_server.Domain +
                ':' +
                plc_link_server.Proxy, null, null);
        }
    }
    // Return the first configuration that match the given name.
    // The load balancer will be in charge to select the correct service instance from the list
    // The first instance is the entry point of the services.
    Globular.prototype.getFirstConfigByName = function (name) {
        for (var id in this.config.Services) {
            var service = this.config.Services[id];
            if (service.Name == name) {
                return service;
            }
        }
        return null;
    };
    return Globular;
}());
exports.Globular = Globular;
