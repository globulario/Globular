/**
 * That file contain a list of function that can be use as API instead of
 * using the globular object itself.
 */
import { LogInfo, LogType, Ressource, Peer } from "./ressource/ressource_pb";
import { ServiceDescriptor } from "./services/services_pb";
import { TagType } from "./plc/plcpb/plc_pb";
import { Globular, EventHub } from './services';
import { IConfig, IServiceConfig } from './services';
/**
 * Return the globular configuration file. The return config object
 * can contain sensible information so it must be called with appropriate
 * permission level.
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
export declare function readFullConfig(globular: Globular, application: string, domain: string, callback: (config: IConfig) => void, errorCallback: (err: any) => void): void;
export declare function saveConfig(globular: Globular, application: string, domain: string, config: IConfig, callback: (config: IConfig) => void, errorCallback: (err: any) => void): void;
/**
 * Retreive the list of ressource owner.
 * @param path
 * @param callback
 * @param errorCallback
 */
export declare function getRessourceOwners(globular: Globular, application: string, domain: string, path: string, callback: (infos: Array<any>) => void, errorCallback: (err: any) => void): void;
/**
 * The ressource owner to be set.
 * @param path The path of the ressource
 * @param owner The owner of the ressource
 * @param callback The success callback
 * @param errorCallback The error callback
 */
export declare function setRessourceOwners(globular: Globular, application: string, domain: string, path: string, owner: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Delete a given ressource owner
 * @param path The path of the ressource.
 * @param owner The owner to be remove
 * @param callback The sucess callback
 * @param errorCallback The error callback
 */
export declare function deleteRessourceOwners(globular: Globular, application: string, domain: string, path: string, owner: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Retreive the permission for a given file.
 * @param path
 * @param callback
 * @param errorCallback
 */
export declare function getRessourcePermissions(globular: Globular, application: string, domain: string, path: string, callback: (infos: Array<any>) => void, errorCallback: (err: any) => void): void;
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
export declare function setRessourcePermission(globular: Globular, application: string, domain: string, path: string, owner: string, ownerType: OwnerType, number: number, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Delete a file permission for a give user.
 * @param path The path of the file on the server.
 * @param owner The owner of the file
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
export declare function deleteRessourcePermissions(globular: Globular, application: string, domain: string, path: string, owner: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Return server files operations.
 * @param globular
 * @param application
 * @param domain
 * @param callbak
 * @param errorCallback
 */
export declare function getAllFilesInfo(globular: Globular, application: string, domain: string, callbak: (filesInfo: any) => void, errorCallback: (err: any) => void): void;
/**
 * Rename a file or a directorie with given name.
 * @param path The path inside webroot
 * @param newName The new file name
 * @param oldName  The old file name
 * @param callback  The success callback.
 * @param errorCallback The error callback.
 */
export declare function renameFile(globular: Globular, application: string, domain: string, path: string, newName: string, oldName: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Delete a file with a given path.
 * @param path The path of the file to be deleted.
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
export declare function deleteFile(globular: Globular, application: string, domain: string, path: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 *
 * @param path The path of the directory to be deleted.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export declare function deleteDir(globular: Globular, application: string, domain: string, path: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Create a dir archive.
 * @param path
 * @param name
 * @param callback
 * @param errorCallback
 */
export declare function createArchive(globular: Globular, application: string, domain: string, path: string, name: string, callback: (path: string) => void, errorCallback: (err: any) => void): void;
/**
 *
 * @param urlToSend
 */
export declare function downloadFileHttp(globular: Globular, application: string, domain: string, urlToSend: string, fileName: string, callback: () => void): void;
/**
 * Download a directory as archive file. (.tar.gz)
 * @param path The path of the directory to dowload.
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
export declare function downloadDir(globular: Globular, application: string, domain: string, path: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Read the content of a dir from a given path.
 * @param path The parent path of the dir to be read.
 * @param callback  Return the path of the dir with more information.
 * @param errorCallback Return a error if the file those not contain the value.
 */
export declare function readDir(globular: Globular, application: string, domain: string, path: string, callback: (dir: any) => void, errorCallback: (err: any) => void): void;
/**
 * Create a new directory inside existing one.
 * @param path The path of the directory
 * @param callback The callback
 * @param errorCallback The error callback
 */
export declare function createDir(globular: Globular, application: string, domain: string, path: string, callback: (dirName: string) => void, errorCallback: (err: any) => void): void;
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
export declare function queryTs(globular: Globular, application: string, domain: string, connectionId: string, query: string, ts: number, callback: (value: any) => void, errorCallback: (error: any) => void): void;
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
export declare function queryTsRange(globular: Globular, application: string, domain: string, connectionId: string, query: string, startTime: number, endTime: number, step: number, callback: (values: any) => void, errorCallback: (err: any) => void): void;
/**
 * Register a new user.
 * @param userName The name of the account
 * @param email The email
 * @param password The password
 * @param confirmPassword
 * @param callback
 * @param errorCallback
 */
export declare function registerAccount(globular: Globular, application: string, domain: string, userName: string, email: string, password: string, confirmPassword: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Remove an account from the server.
 * @param name  The _id of the account.
 * @param callback The callback when the action succed
 * @param errorCallback The error callback.
 */
export declare function DeleteAccount(globular: Globular, application: string, domain: string, id: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Remove a role from an account.
 * @param accountId The account id
 * @param roleId The role name (id)
 * @param callback The success callback
 * @param errorCallback The error callback
 */
export declare function RemoveRoleFromAccount(globular: Globular, application: string, domain: string, accountId: string, roleId: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Append a role to an account.
 * @param accountId The account id
 * @param roleId The role name (id)
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export declare function AppendRoleToAccount(globular: Globular, application: string, domain: string, accountId: string, roleId: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Update the account email
 * @param accountId The account id
 * @param old_email the old email
 * @param new_email the new email
 * @param callback  the callback when success
 * @param errorCallback the error callback in case of error
 */
export declare function updateAccountEmail(globular: Globular, application: string, domain: string, accountId: string, old_email: string, new_email: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * The update account password
 * @param accountId The account id
 * @param old_password The old password
 * @param new_password The new password
 * @param confirm_password The new password confirmation
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export declare function updateAccountPassword(globular: Globular, application: string, domain: string, accountId: string, old_password: string, new_password: string, confirm_password: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Authenticate the user and get the token
 * @param userName The account name or email
 * @param password  The user password
 * @param callback
 * @param errorCallback
 */
export declare function authenticate(globular: Globular, eventHub: EventHub, application: string, domain: string, userName: string, password: string, callback: (value: any) => void, errorCallback: (err: any) => void): void;
/**
 * Function to be use to refresh token or full configuration.
 * @param callback On success callback
 * @param errorCallback On error callback
 */
export declare function refreshToken(globular: Globular, eventHub: EventHub, application: string, domain: string, callback: (token: any) => void, errorCallback: (err: any) => void): void;
/**
 * Save user data into the user_data collection.
 */
export declare function appendUserData(globular: Globular, application: string, domain: string, data: any, callback: (id: string) => void): void;
/**
 * Read user data one result at time.
 */
export declare function readOneUserData(globular: Globular, application: string, domain: string, query: string, callback: (results: any) => void): void;
/**
 * Read all user data.
 */
export declare function readUserData(globular: Globular, application: string, domain: string, query: string, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 * Return the list of all account on the server, guest and admin are new account...
 * @param callback
 */
export declare function GetAllAccountsInfo(globular: Globular, application: string, domain: string, callback: (accounts: Array<any>) => void, errorCallback: (err: any) => void): void;
/**
 * Retreive all available actions on the server.
 * @param callback That function is call in case of success.
 * @param errorCallback That function is call in case error.
 */
export declare function getAllActions(globular: Globular, application: string, domain: string, callback: (ations: Array<string>) => void, errorCallback: (err: any) => void): void;
/**
 * Retreive the list of all available roles on the server.
 * @param callback That function is call in case of success.
 * @param errorCallback That function is call in case error.
 */
export declare function getAllRoles(globular: Globular, application: string, domain: string, callback: (roles: Array<any>) => void, errorCallback: (err: any) => void): void;
/**
 * Append Action to a given role.
 * @param action The action name.
 * @param role The role.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export declare function AppendActionToRole(globular: Globular, application: string, domain: string, role: string, action: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Remove the action from a given role.
 * @param action The action id
 * @param role The role id
 * @param callback success callback
 * @param errorCallback error callback
 */
export declare function RemoveActionFromRole(globular: Globular, application: string, domain: string, role: string, action: string, callback: () => void, errorCallback: (err: any) => void): void;
export declare function CreateRole(globular: Globular, application: string, domain: string, id: string, callback: () => void, errorCallback: (err: any) => void): void;
export declare function DeleteRole(globular: Globular, application: string, domain: string, id: string, callback: () => void, errorCallback: (err: any) => void): void;
export declare function GetAllApplicationsInfo(globular: Globular, application: string, domain: string, callback: (infos: any) => void, errorCallback: (err: any) => void): void;
export declare function AppendActionToApplication(globular: Globular, application: string, domain: string, applicationId: string, action: string, callback: () => void, errorCallback: (err: any) => void): void;
export declare function RemoveActionFromApplication(globular: Globular, application: string, domain: string, applicationId: string, action: string, callback: () => void, errorCallback: (err: any) => void): void;
export declare function DeleteApplication(globular: Globular, application: string, domain: string, applicationId: string, callback: () => void, errorCallback: (err: any) => void): void;
export declare function SaveApplication(globular: Globular, eventHub: EventHub, application_: string, domain: string, application: any, callback: () => void, errorCallback: (err: any) => void): void;
export declare function GetAllPeersInfo(globular: Globular, application: string, domain: string, query: string, callback: (peers: Peer[]) => void, errorCallback: (err: any) => void): void;
export declare function AppendActionToPeer(globular: Globular, application: string, domain: string, id: string, action: string, callback: () => void, errorCallback: (err: any) => void): void;
export declare function RemoveActionFromPeer(globular: Globular, application: string, domain: string, id: string, action: string, callback: () => void, errorCallback: (err: any) => void): void;
export declare function DeletePeer(globular: Globular, application: string, domain: string, peer: Peer, callback: () => void, errorCallback: (err: any) => void): void;
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
export declare function GetServiceDescriptor(globular: Globular, application: string, domain: string, serviceId: string, publisherId: string, callback: (descriptors: Array<ServiceDescriptor>) => void, errorCallback: (err: any) => void): void;
/**
 * Get the list of all service descriptor hosted on a server.
 * @param globular The globular object instance
 * @param application The application name who called the function.
 * @param domain The domain where the application reside.
 * @param callback
 * @param errorCallback
 */
export declare function GetServicesDescriptor(globular: Globular, application: string, domain: string, callback: (descriptors: Array<ServiceDescriptor>) => void, errorCallback: (err: any) => void): void;
export declare function SetServicesDescriptor(globular: Globular, application: string, domain: string, descriptor: ServiceDescriptor, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Find services by keywords.
 * @param query
 * @param callback
 */
export declare function findServices(globular: Globular, application: string, domain: string, keywords: Array<string>, callback: (results: Array<ServiceDescriptor>) => void, errorCallback: (err: any) => void): void;
export declare function installService(globular: Globular, application: string, domain: string, discoveryId: string, serviceId: string, publisherId: string, version: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Stop a service.
 */
export declare function stopService(globular: Globular, application: string, domain: string, serviceId: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Start a service
 * @param serviceId The id of the service to start.
 * @param callback  The callback on success.
 */
export declare function startService(globular: Globular, application: string, domain: string, serviceId: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Here I will save the service configuration.
 * @param service The configuration to save.
 */
export declare function saveService(globular: Globular, application: string, domain: string, service: IServiceConfig, callback: (config: any) => void, errorCallback: (err: any) => void): void;
export declare function uninstallService(globular: Globular, application: string, domain: string, service: IServiceConfig, deletePermissions: boolean, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Return the list of service bundles.
 * @param callback
 */
export declare function GetServiceBundles(globular: Globular, application: string, domain: string, publisherId: string, serviceId: string, version: string, callback: (bundles: Array<any>) => void, errorCallback: (err: any) => void): void;
export declare function getReferencedValue(globular: Globular, application: string, domain: string, ref: any, callback: (results: any) => void, errorCallback: (err: any) => void): void;
/**
 * Read all errors data.
 * @param callback
 */
export declare function readErrors(globular: Globular, application: string, domain: string, callback: (results: any) => void, errorCallback: (err: any) => void): void;
export declare function readAllActionPermission(globular: Globular, application: string, domain: string, callback: (results: any) => void, errorCallback: (err: any) => void): void;
export declare function getRessources(globular: Globular, application: string, domain: string, path: string, name: string, callback: (results: Ressource[]) => void, errorCallback: (err: any) => void): void;
export declare function setActionPermission(globular: Globular, application: string, domain: string, action: string, permission: number, callback: (results: any) => void, errorCallback: (err: any) => void): void;
export declare function removeActionPermission(globular: Globular, application: string, domain: string, action: string, callback: (results: any) => void, errorCallback: (err: any) => void): void;
export declare function removeRessource(globular: Globular, application: string, domain: string, path: string, name: string, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Read all logs
 * @param callback The success callback.
 */
export declare function readLogs(globular: Globular, application: string, domain: string, query: string, callback: (results: any) => void, errorCallback: (err: any) => void): void;
export declare function clearAllLog(globular: Globular, application: string, domain: string, logType: LogType, callback: () => void, errorCallback: (err: any) => void): void;
export declare function deleteLog(globular: Globular, application: string, domain: string, log: LogInfo, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Return the logged method and their count.
 * @param pipeline
 * @param callback
 * @param errorCallback
 */
export declare function getNumbeOfLogsByMethod(globular: Globular, application: string, domain: string, callback: (resuts: Array<any>) => void, errorCallback: (err: any) => void): void;
export declare enum PLC_TYPE {
    ALEN_BRADLEY = 1,
    SIEMENS = 2,
    MODBUS = 3
}
/**
* Read a plc tag from the defined backend.
* @param plcType  The plc type can be Alen Bradley or Simens, modbus is on the planned.
* @param connectionId  The connection id defined for that plc.
* @param name The name of the tag to read.
* @param type The type name of the plc.
* @param offset The offset in the memory.
*/
export declare function readPlcTag(globular: Globular, application: string, domain: string, plcType: PLC_TYPE, connectionId: string, name: string, type: TagType, offset: number): Promise<any>;
/**
 * Synchronize LDAP and Globular/MongoDB user and roles.
 * @param info The synchronisations informations.
 * @param callback success callback.
 */
export declare function syncLdapInfos(globular: Globular, application: string, domain: string, info: any, timeout: number, callback: () => void, errorCallback: (err: any) => void): void;
/**
 * Ping globular sql service.
 * @param globular
 * @param application
 * @param domain
 * @param connectionId
 * @param callback
 * @param errorCallback
 */
export declare function pingSql(globular: Globular, application: string, domain: string, connectionId: string, callback: (pong: string) => {}, errorCallback: (err: any) => void): void;
