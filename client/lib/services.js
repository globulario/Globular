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
        // Iinitialisation of services.
        if (this.config.Services['catalog_server'] != undefined) {
            this.catalogService = new catalog_grpc_web_pb_1.CatalogServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['catalog_server'].Domain +
                ':' +
                this.config.Services['catalog_server'].Proxy, null, null);
        }
        if (this.config.Services['echo_server'] != null) {
            this.echoService = new echo_grpc_web_pb_1.EchoServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['echo_server'].Domain +
                ':' +
                this.config.Services['echo_server'].Proxy, null, null);
        }
        if (this.config.Services['event_server'] != null) {
            this.eventService = new event_grpc_web_pb_1.EventServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['event_server'].Domain +
                ':' +
                this.config.Services['event_server'].Proxy, null, null);
        }
        if (this.config.Services['file_server'] != null) {
            this.fileService = new file_grpc_web_pb_1.FileServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['file_server'].Domain +
                ':' +
                this.config.Services['file_server'].Proxy, null, null);
        }
        if (this.config.Services['ldap_server'] != null) {
            this.ldapService = new ldap_grpc_web_pb_1.LdapServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['ldap_server'].Domain +
                ':' +
                this.config.Services['ldap_server'].Proxy, null, null);
        }
        if (this.config.Services['persistence_server'] != null) {
            this.persistenceService = new persistence_grpc_web_pb_1.PersistenceServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['persistence_server'].Domain +
                ':' +
                this.config.Services['persistence_server'].Proxy, null, null);
        }
        if (this.config.Services['smtp_server'] != null) {
            this.smtpService = new smtp_grpc_web_pb_1.SmtpServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['smtp_server'].Domain +
                ':' +
                this.config.Services['smtp_server'].Proxy, null, null);
        }
        if (this.config.Services['sql_server'] != null) {
            this.sqlService = new sql_grpc_web_pb_1.SqlServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['sql_server'].Domain +
                ':' +
                this.config.Services['sql_server'].Proxy, null, null);
        }
        if (this.config.Services['storage_server'] != null) {
            this.storageService = new storage_grpc_web_pb_1.StorageServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['storage_server'].Domain +
                ':' +
                this.config.Services['storage_server'].Proxy, null, null);
        }
        if (this.config.Services['monitoring_server'] != null) {
            this.monitoringService = new monitoring_grpc_web_pb_1.MonitoringServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['monitoring_server'].Domain +
                ':' +
                this.config.Services['monitoring_server'].Proxy, null, null);
        }
        if (this.config.Services['spc_server'] != null) {
            this.spcService = new spc_grpc_web_pb_1.SpcServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['spc_server'].Domain +
                ':' +
                this.config.Services['spc_server'].Proxy, null, null);
        }
        // non open source services.
        if (this.config.Services['plc_server_ab'] != null) {
            this.plcService_ab = new plc_grpc_web_pb_1.PlcServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['plc_server_ab'].Domain +
                ':' +
                this.config.Services['plc_server_ab'].Proxy, null, null);
        }
        if (this.config.Services['plc_server_siemens'] != null) {
            this.plcService_siemens = new plc_grpc_web_pb_1.PlcServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['plc_server_siemens'].Domain +
                ':' +
                this.config.Services['plc_server_siemens'].Proxy, null, null);
        }
        if (this.config.Services['plc_link_server'] != null) {
            this.plcLinkService = new plc_link_grpc_web_pb_1.PlcLinkServicePromiseClient(this.config.Protocol +
                '://' +
                this.config.Services['plc_link_server'].Domain +
                ':' +
                this.config.Services['plc_link_server'].Proxy, null, null);
        }
    }
    return Globular;
}());
exports.Globular = Globular;
