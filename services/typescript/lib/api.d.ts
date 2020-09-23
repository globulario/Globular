/**
 * That file contain a list of function that can be use as API instead of
 * using the globular object itself.
 */
import { LogInfo, LogType, Ressource } from "./ressource/ressource_pb";
import { ServiceDescriptor } from "./services/services_pb";
import { TagType } from "./plc/plc_pb";
import { Globular, EventHub } from './services';
import { IConfig, IServiceConfig } from './services';
import { SearchResult } from "./search/search_pb";
/**
 * Return the globular configuration file. The return config object
 * can contain sensible information so it must be called with appropriate
 * level of permission.
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
export declare function readFullConfig(globular: Globular, callback: (config: IConfig) => void, errorCallback: (err: any) => void): void;
/**
 * Save a configuration
 * @param globular
 * @param application
 * @param domain
 * @param config The configuration to be save.
 * @param callback
 * @param errorCallback
 */
export declare function saveConfig(globular: Globular, config: IConfig, callback: (config: IConfig) => void, errorCallback: (err: any) => void): void;
/**
 * Return the list of all action permissions.
 * Action permission are apply on ressource managed by those action.
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
export declare function readAllActionPermissions(globular: Globular, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 * Return the list of ressources.
 * @param globular
 * @param application
 * @param domain
 * @param path
 * @param name
 * @param callback
 * @param errorCallback
 */
export declare function getRessources(globular: Globular, path: string, name: string, callback: (results: Ressource[]) => void, errorCallback: (err: any) => void): void;
/**
 * Set/create action permission.
 * @param globular
 * @param application
 * @param domain
 * @param action
 * @param permission
 * @param callback
 * @param errorCallback
 */
export declare function setActionPermission(globular: Globular, action: string, permission: number, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 * Delete action permission.
 * @param globular
 * @param application
 * @param domain
 * @param action
 * @param callback
 * @param errorCallback
 */
export declare function removeActionPermission(globular: Globular, action: string, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 * Delete a ressource
 * @param globular
 * @param application
 * @param domain
 * @param path
 * @param name
 * @param callback
 * @param errorCallback
 */
export declare function removeRessource(globular: Globular, path: string, name: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Retreive the list of ressource owner.
 * @param path
 * @param callback
 * @param errorCallback
 */
export declare function getRessourceOwners(globular: Globular, path: string, callback: (infos: any[]) => void, errorCallback: (err: any) => void): void;
/**
 * The ressource owner to be set.
 * @param path The path of the ressource
 * @param owner The owner of the ressource
 * @param callback The success callback
 * @param errorCallback The error callback
 */
export declare function setRessourceOwners(globular: Globular, path: string, owner: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Delete a given ressource owner
 * @param path The path of the ressource.
 * @param owner The owner to be remove
 * @param callback The sucess callback
 * @param errorCallback The error callback
 */
export declare function deleteRessourceOwners(globular: Globular, path: string, owner: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Retreive the permission for a given ressource.
 * @param path
 * @param callback
 * @param errorCallback
 */
export declare function getRessourcePermissions(globular: Globular, path: string, callback: (infos: any[]) => void, errorCallback: (err: any) => void): void;
/**
 * The permission can be assigned to
 * a User, a Role or an Application.
 */
export declare enum OwnerType {
    User = 1,
    Role = 2,
    Application = 3
}
/**
 * Create a file permission.
 * @param path The path on the server from the root.
 * @param owner The owner of the permission
 * @param ownerType The owner type
 * @param number The (unix) permission number.
 * @param callback The success callback
 * @param errorCallback The error callback
 */
export declare function setRessourcePermission(globular: Globular, path: string, owner: string, ownerType: OwnerType, permissionNumber: number, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Delete a file permission for a give user.
 * @param path The path of the file on the server.
 * @param owner The owner of the file
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
export declare function deleteRessourcePermissions(globular: Globular, path: string, owner: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Return server files operations.
 * @param globular
 * @param application
 * @param domain
 * @param callbak
 * @param errorCallback
 */
export declare function getAllFilesInfo(globular: Globular, callbak: (filesInfo: any) => void, errorCallback: (err: any) => void): void;
/**
 * Rename a file or a directorie with given name.
 * @param path The path inside webroot
 * @param newName The new file name
 * @param oldName  The old file name
 * @param callback  The success callback.
 * @param errorCallback The error callback.
 */
export declare function renameFile(globular: Globular, path: string, newName: string, oldName: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Delete a file with a given path.
 * @param path The path of the file to be deleted.
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
export declare function deleteFile(globular: Globular, path: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Remove a given directory and all element it contain.
 * @param path The path of the directory to be deleted.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export declare function deleteDir(globular: Globular, path: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Create a dir archive.
 * @param path
 * @param name
 * @param callback
 * @param errorCallback
 */
export declare function createArchive(globular: Globular, path: string, name: string, callback: (path: string) => void, errorCallback: (err: any) => void): void;
/**
 * Download a file from the server.
 * @param urlToSend
 */
export declare function downloadFileHttp(urlToSend: string, fileName: string, callback: () => void): void;
/**
 * Download a directory as archive file. (.tar.gz)
 * @param path The path of the directory to dowload.
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
export declare function downloadDir(globular: Globular, path: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Read the content of a dir from a given path.
 * @param path The parent path of the dir to be read.
 * @param callback  Return the path of the dir with more information.
 * @param errorCallback Return a error if the file those not contain the value.
 */
export declare function readDir(globular: Globular, path: string, callback: (dir: any) => void, errorCallback: (err: any) => void): void;
/**
 * Create a new directory inside existing one.
 * @param path The path of the directory
 * @param callback The callback
 * @param errorCallback The error callback
 */
export declare function createDir(globular: Globular, path: string, callback: (dirName: string) => void, errorCallback: (err: any) => void): void;
/**
 * Run a query over a time series database.
 * @param globular
 * @param application
 * @param domain
 * @param connectionId
 * @param query
 * @param ts
 * @param callback
 * @param errorCallback
 */
export declare function queryTs(globular: Globular, connectionId: string, query: string, ts: number, callback: (value: any) => void, errorCallback: (error: any) => void): void;
/**
 * Run query over a time series
 * @param globular
 * @param application
 * @param domain
 * @param connectionId
 * @param query
 * @param startTime
 * @param endTime
 * @param step
 * @param callback
 * @param errorCallback
 */
export declare function queryTsRange(globular: Globular, connectionId: string, query: string, startTime: number, endTime: number, step: number, callback: (values: any) => void, errorCallback: (err: any) => void): void;
/**
 * Return the list of all account on the server, guest and admin are new account...
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
export declare function GetAllAccountsInfo(globular: Globular, callback: (accounts: any[]) => void, errorCallback: (err: any) => void): void;
/**
 * Register a new account.
 * @param userName The name of the account
 * @param email The email
 * @param password The password
 * @param confirmPassword
 * @param callback
 * @param errorCallback
 */
export declare function registerAccount(globular: Globular, userName: string, email: string, password: string, confirmPassword: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Remove an account from the server.
 * @param name  The _id of the account.
 * @param callback The callback when the action succed
 * @param errorCallback The error callback.
 */
export declare function DeleteAccount(globular: Globular, id: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Remove a role from an account.
 * @param accountId The account id
 * @param roleId The role name (id)
 * @param callback The success callback
 * @param errorCallback The error callback
 */
export declare function RemoveRoleFromAccount(globular: Globular, accountId: string, roleId: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Append a role to an account.
 * @param accountId The account id
 * @param roleId The role name (id)
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export declare function AppendRoleToAccount(globular: Globular, accountId: string, roleId: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Update the account email
 * @param accountId The account id
 * @param old_email the old email
 * @param new_email the new email
 * @param callback  the callback when success
 * @param errorCallback the error callback in case of error
 */
export declare function updateAccountEmail(globular: Globular, accountId: string, oldEmail: string, newEmail: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * The update account password
 * @param accountId The account id
 * @param old_password The old password
 * @param new_password The new password
 * @param confirm_password The new password confirmation
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export declare function updateAccountPassword(globular: Globular, accountId: string, oldPassword: string, newPassword: string, confirmPassword: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Authenticate the user and get the token
 * @param globular
 * @param eventHub
 * @param application
 * @param domain
 * @param userName
 * @param password
 * @param callback
 * @param errorCallback
 */
export declare function authenticate(globular: Globular, eventHub: EventHub, userName: string, password: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Function to be use to refresh token.
 * @param globular
 * @param eventHub
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
export declare function refreshToken(globular: Globular, eventHub: EventHub, callback: (token: any) => void, errorCallback: (err: any) => void): void;
/**
 * Save user data into the user_data collection.
 * @param globular
 * @param application
 * @param domain
 * @param data
 * @param callback
 * @param errorCallback
 */
export declare function appendUserData(globular: Globular, data: any, callback: (id: string) => void, errorCallback: (err: any) => void): void;
/**
 * Read user data one result at time.
 * @param globular
 * @param application
 * @param domain
 * @param query
 * @param callback
 * @param errorCallback
 */
export declare function readOneUserData(globular: Globular, query: string, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 * Read all user data.
 * @param globular
 * @param application
 * @param domain
 * @param query
 * @param callback
 * @param errorCallback
 */
export declare function readUserData(globular: Globular, query: string, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 * Retreive all available actions on the server.
 * @param callback That function is call in case of success.
 * @param errorCallback That function is call in case error.
 */
export declare function getAllActions(globular: Globular, callback: (ations: string[]) => void, errorCallback: (err: any) => void): void;
/**
 * Retreive the list of all available roles on the server.
 * @param callback That function is call in case of success.
 * @param errorCallback That function is call in case error.
 */
export declare function getAllRoles(globular: Globular, callback: (roles: any[]) => void, errorCallback: (err: any) => void): void;
/**
 * Append Action to a given role.
 * @param action The action name.
 * @param role The role.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export declare function AppendActionToRole(globular: Globular, role: string, action: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Remove the action from a given role.
 * @param action The action id
 * @param role The role id
 * @param callback success callback
 * @param errorCallback error callback
 */
export declare function RemoveActionFromRole(globular: Globular, role: string, action: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Create a new Role
 * @param globular
 * @param application
 * @param domain
 * @param id
 * @param callback
 * @param errorCallback
 */
export declare function CreateRole(globular: Globular, id: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Delete a given role
 * @param globular
 * @param application
 * @param domain
 * @param id
 * @param callback
 * @param errorCallback
 */
export declare function DeleteRole(globular: Globular, id: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Return the list of all application
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
export declare function GetAllApplicationsInfo(globular: Globular, callback: (infos: any) => void, errorCallback: (err: any) => void): void;
/**
 * Append action to application.
 * @param globular
 * @param application
 * @param domain
 * @param applicationId
 * @param action
 * @param callback
 * @param errorCallback
 */
export declare function AppendActionToApplication(globular: Globular, applicationId: string, action: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Remove action from application.
 * @param globular
 * @param application
 * @param domain
 * @param action
 * @param callback
 * @param errorCallback
 */
export declare function RemoveActionFromApplication(globular: Globular, action: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Delete application
 * @param globular
 * @param application
 * @param domain
 * @param applicationId
 * @param callback
 * @param errorCallback
 */
export declare function DeleteApplication(globular: Globular, applicationId: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Save application
 * @param globular
 * @param eventHub
 * @param applicationId
 * @param domain
 * @param application
 * @param callback
 * @param errorCallback
 */
export declare function SaveApplication(globular: Globular, eventHub: EventHub, _application: any, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Return a list of service descriptor related to a service host by a server.
 * @param globular
 * @param application
 * @param domain
 * @param serviceId
 * @param publisherId
 * @param callback
 * @param errorCallback
 */
export declare function GetServiceDescriptor(globular: Globular, serviceId: string, publisherId: string, callback: (descriptors: ServiceDescriptor[]) => void, errorCallback: (err: any) => void): void;
/**
 * Get the list of all service descriptor hosted on a server.
 * @param globular The globular object instance
 * @param application The application name who called the function.
 * @param domain The domain where the application reside.
 * @param callback
 * @param errorCallback
 */
export declare function GetServicesDescriptor(globular: Globular, callback: (descriptors: ServiceDescriptor[]) => void, errorCallback: (err: any) => void): void;
/**
 * Create or update a service descriptor.
 * @param globular
 * @param application
 * @param domain
 * @param descriptor
 * @param callback
 * @param errorCallback
 */
export declare function SetServicesDescriptor(globular: Globular, descriptor: ServiceDescriptor, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Find services by keywords.
 * @param query
 * @param callback
 */
export declare function findServices(globular: Globular, keywords: string[], callback: (results: ServiceDescriptor[]) => void, errorCallback: (err: any) => void): void;
/**
 * Install a service
 * @param globular
 * @param application
 * @param domain
 * @param discoveryId
 * @param serviceId
 * @param publisherId
 * @param version
 * @param callback
 * @param errorCallback
 */
export declare function installService(globular: Globular, discoveryId: string, serviceId: string, publisherId: string, version: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Stop a service.
 */
export declare function stopService(globular: Globular, serviceId: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Start a service
 * @param serviceId The id of the service to start.
 * @param callback  The callback on success.
 */
export declare function startService(globular: Globular, serviceId: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Here I will save the service configuration.
 * @param service The configuration to save.
 */
export declare function saveService(globular: Globular, service: IServiceConfig, callback: (config: any) => void, errorCallback: (err: any) => void): void;
/**
 * Uninstall a service from the server.
 * @param globular
 * @param application
 * @param domain
 * @param service
 * @param deletePermissions
 * @param callback
 * @param errorCallback
 */
export declare function uninstallService(globular: Globular, service: IServiceConfig, deletePermissions: boolean, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Return the list of service bundles.
 * @param callback
 */
export declare function GetServiceBundles(globular: Globular, publisherId: string, serviceId: string, version: string, callback: (bundles: any[]) => void, errorCallback: (err: any) => void): void;
/**
 * Get the object pointed by a reference.
 * @param globular
 * @param application
 * @param domain
 * @param ref
 * @param callback
 * @param errorCallback
 */
export declare function getReferencedValue(globular: Globular, ref: any, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 * Read all errors data for server log.
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
export declare function readErrors(globular: Globular, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 *  Read all logs
 * @param globular
 * @param application
 * @param domain
 * @param query
 * @param callback
 * @param errorCallback
 */
export declare function readLogs(globular: Globular, query: string, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 * Clear all log of a given type.
 * @param globular
 * @param application
 * @param domain
 * @param logType
 * @param callback
 * @param errorCallback
 */
export declare function clearAllLog(globular: Globular, logType: LogType, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Delete log entry.
 * @param globular
 * @param application
 * @param domain
 * @param log
 * @param callback
 * @param errorCallback
 */
export declare function deleteLogEntry(globular: Globular, log: LogInfo, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Return the logged method and their count.
 * @param pipeline
 * @param callback
 * @param errorCallback
 */
export declare function getNumbeOfLogsByMethod(globular: Globular, callback: (resuts: any[]) => void, errorCallback: (err: any) => void): void;
export declare enum PLC_TYPE {
    ALEN_BRADLEY = 1,
    SIEMENS = 2,
    MODBUS = 3
}
/**
 * Read a plc tag.
 * @param globular
 * @param application
 * @param domain
 * @param plcType
 * @param connectionId
 * @param name
 * @param type
 * @param offset
 */
export declare function readPlcTag(globular: Globular, plcType: PLC_TYPE, connectionId: string, name: string, type: TagType, offset: number): Promise<any>;
/**
 * Synchronize LDAP and Globular/MongoDB user and roles.
 * @param info The synchronisations informations.
 * @param callback success callback.
 */
export declare function syncLdapInfos(globular: Globular, info: any, timeout: number, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Ping globular sql service.
 * @param globular
 * @param application
 * @param domain
 * @param connectionId
 * @param callback
 * @param errorCallback
 */
export declare function pingSql(globular: Globular, connectionId: string, callback: (pong: string) => {}, errorCallback: (err: any) => void): void;
/**
 * Search Documents from given database(s) and return results. The search engine in use is
 * xapian, so query must follow xapian query rules.
 * @param globular The server object.
 * @param paths The list of database paths.
 * @param query The query to execute.
 * @param language The language of the database.
 * @param fields The list of field to query, can be empty if all fields must be search or fields are specified in the query.
 * @param offset The offset of resultset
 * @param pageSize The number of result to return
 * @param snippetLength The length of the snippet result.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export declare function searchDocuments(globular: Globular, paths: string[], query: string, language: string, fields: string[], offset: number, pageSize: number, snippetLength: number, callback: (results: SearchResult[]) => void, errorCallback: (err: any) => void): void;
