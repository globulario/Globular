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
/**
 * Globular regroup all serivces in one object that can be use by
 * application to get access to sql, ldap, persistence... service.
 */
var Globular = /** @class */ (function () {
    /** The */
    function Globular(config) {
        // Keep the config...
        this.config = config;
        /** The admin service to access to other configurations. */
        this.adminService = new admin_grpc_web_pb_1.AdminServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.AdminProxy, null, null);
        this.ressourceService = new ressource_grpc_web_pb_1.RessourceServicePromiseClient(this.config.Protocol + '://' + this.config.Domain + ':' + this.config.RessourceProxy, null, null);
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
