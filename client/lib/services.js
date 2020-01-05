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
var plc_grpc_web_pb_1 = require("./plc/plcpb/plc_grpc_web_pb");
var admin_grpc_web_pb_1 = require("./admin/admin_grpc_web_pb");
var ressource_grpc_web_pb_1 = require("./ressource/ressource_grpc_web_pb");
var services_grpc_web_pb_1 = require("./services/services_grpc_web_pb");
var ca_grpc_web_pb_1 = require("./ca/ca_grpc_web_pb");
var event_pb_1 = require("./event/eventpb/event_pb");
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
        if (this.service != undefined) {
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
                if (_this.subscribers[name] == undefined) {
                    _this.subscribers[name] = {};
                }
                _this.subscribers[name][uuid] = onevent;
                onsubscribe(uuid);
            });
        }
        else {
            // create a uuid and call onsubscribe callback.
            if (this.subscribers[name] == undefined) {
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
        // Remove the local subscriber.
        delete this.subscribers[name][uuid];
        if (Object.keys(this.subscribers[name]).length == 0) {
            delete this.subscribers[name];
            // disconnect from the distant server.
            if (this.service != undefined) {
                var request = new event_pb_1.UnSubscribeRequest();
                request.setName(name);
                request.setUuid(this.subscriptions[name]);
                // remove the subcription uuid.
                delete this.subscriptions[name];
                // Now I will test with promise
                this.service.unSubscribe(request)
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
        if (local == true) {
            this.dispatch(name, data);
        }
        else {
            // Create a new request.
            var request = new event_pb_1.PublishRequest();
            var evt = new event_pb_1.Event();
            evt.setName(name);
            var enc = new TextEncoder(); // always utf-8
            // encode the string to a array of byte
            evt.setData(enc.encode(data));
            request.setEvt(evt);
            // Now I will test with promise
            this.service.publish(request)
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
            this.subscribers[name][uuid](data);
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
        /** That service help to find and install or publish new service on the backend. */
        this.servicesDicovery = new services_grpc_web_pb_1.ServiceDiscoveryPromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.ServicesDiscoveryProxy, null, null);
        /** Functionality to use service repository server. */
        this.servicesRepository = new services_grpc_web_pb_1.ServiceRepositoryPromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.ServicesRepositoryProxy, null, null);
        /** Certificate authority function needed for TLS client. */
        this.certificateAuthority = new ca_grpc_web_pb_1.CertificateAuthorityPromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.CertificateAuthorityProxy, null, null);
        // Iinitialisation of services.
        if (this.config.Services['catalog_server'] != undefined) {
            var protocol = 'http';
            if (this.config.Services['catalog_server'].TLS == true) {
                protocol = 'https';
            }
            this.catalogService = new catalog_grpc_web_pb_1.CatalogServicePromiseClient(protocol +
                '://' +
                this.config.Services['catalog_server'].Domain +
                ':' +
                this.config.Services['catalog_server'].Proxy, null, null);
        }
        if (this.config.Services['echo_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['echo_server'].TLS == true) {
                protocol = 'https';
            }
            this.echoService = new echo_grpc_web_pb_1.EchoServicePromiseClient(protocol + '://' + this.config.Services['echo_server'].Domain + ':' + this.config.Services['echo_server'].Proxy, null, null);
        }
        if (this.config.Services['event_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['event_server'].TLS == true) {
                protocol = 'https';
            }
            this.eventService = new event_grpc_web_pb_1.EventServicePromiseClient(protocol +
                '://' +
                this.config.Services['event_server'].Domain +
                ':' +
                this.config.Services['event_server'].Proxy, null, null);
        }
        if (this.config.Services['file_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['file_server'].TLS == true) {
                protocol = 'https';
            }
            this.fileService = new file_grpc_web_pb_1.FileServicePromiseClient(protocol + '://' + this.config.Services['file_server'].Domain + ':' + this.config.Services['file_server'].Proxy, null, null);
        }
        if (this.config.Services['ldap_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['ldap_server'].TLS == true) {
                protocol = 'https';
            }
            this.ldapService = new ldap_grpc_web_pb_1.LdapServicePromiseClient(protocol + '://' + this.config.Services['ldap_server'].Domain + ':' + this.config.Services['ldap_server'].Proxy, null, null);
        }
        if (this.config.Services['persistence_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['persistence_server'].TLS == true) {
                protocol = 'https';
            }
            this.persistenceService = new persistence_grpc_web_pb_1.PersistenceServicePromiseClient(protocol +
                '://' +
                this.config.Services['persistence_server'].Domain +
                ':' +
                this.config.Services['persistence_server'].Proxy, null, null);
        }
        if (this.config.Services['smtp_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['smtp_server'].TLS == true) {
                protocol = 'https';
            }
            this.smtpService = new smtp_grpc_web_pb_1.SmtpServicePromiseClient(protocol + '://' + this.config.Services['smtp_server'].Domain + ':' + this.config.Services['smtp_server'].Proxy, null, null);
        }
        if (this.config.Services['sql_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['sql_server'].TLS == true) {
                protocol = 'https';
            }
            this.sqlService = new sql_grpc_web_pb_1.SqlServicePromiseClient(protocol + '://' + this.config.Services['sql_server'].Domain + ':' + this.config.Services['sql_server'].Proxy, null, null);
        }
        if (this.config.Services['storage_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['storage_server'].TLS == true) {
                protocol = 'https';
            }
            this.storageService = new storage_grpc_web_pb_1.StorageServicePromiseClient(protocol +
                '://' +
                this.config.Services['storage_server'].Domain +
                ':' +
                this.config.Services['storage_server'].Proxy, null, null);
        }
        if (this.config.Services['monitoring_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['monitoring_server'].TLS == true) {
                protocol = 'https';
            }
            this.monitoringService = new monitoring_grpc_web_pb_1.MonitoringServicePromiseClient(protocol +
                '://' +
                this.config.Services['monitoring_server'].Domain +
                ':' +
                this.config.Services['monitoring_server'].Proxy, null, null);
        }
        if (this.config.Services['spc_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['spc_server'].TLS == true) {
                protocol = 'https';
            }
            this.spcService = new spc_grpc_web_pb_1.SpcServicePromiseClient(protocol + '://' + this.config.Services['spc_server'].Domain + ':' + this.config.Services['spc_server'].Proxy, null, null);
        }
        // non open source services.
        if (this.config.Services['plc_server_ab'] != null) {
            var protocol = 'http';
            if (this.config.Services['plc_server_ab'].TLS == true) {
                protocol = 'https';
            }
            this.plcService_ab = new plc_grpc_web_pb_1.PlcServicePromiseClient(protocol +
                '://' +
                this.config.Services['plc_server_ab'].Domain +
                ':' +
                this.config.Services['plc_server_ab'].Proxy, null, null);
        }
        if (this.config.Services['plc_server_siemens'] != null) {
            var protocol = 'http';
            if (this.config.Services['plc_server_siemens'].TLS == true) {
                protocol = 'https';
            }
            this.plcService_siemens = new plc_grpc_web_pb_1.PlcServicePromiseClient(protocol +
                '://' +
                this.config.Services['plc_server_siemens'].Domain +
                ':' +
                this.config.Services['plc_server_siemens'].Proxy, null, null);
        }
        if (this.config.Services['plc_link_server'] != null) {
            var protocol = 'http';
            if (this.config.Services['plc_link_server'].TLS == true) {
                protocol = 'https';
            }
            this.plcLinkService = new plc_link_grpc_web_pb_1.PlcLinkServicePromiseClient(protocol +
                '://' +
                this.config.Services['plc_link_server'].Domain +
                ':' +
                this.config.Services['plc_link_server'].Proxy, null, null);
        }
    }
    return Globular;
}());
exports.Globular = Globular;
