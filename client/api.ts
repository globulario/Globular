/**
 * That file contain a list of function that can be use as API instead of
 * using the globular object itself.
 */

import {
    GetConfigRequest,
    SaveConfigRequest,
    InstallServiceRequest,
    InstallServiceResponse,
    StopServiceRequest,
    StartServiceRequest,
    SaveConfigResponse,
    SetEmailRequest,
    SetPasswordRequest,

    UninstallServiceRequest,
    UninstallServiceResponse
} from "./admin/admin_pb";

import {
    QueryRangeRequest,
    QueryRequest
} from "./monitoring/monitoringpb/monitoring_pb";

import {
    RegisterAccountRqst,
    AuthenticateRqst,
    Account,
    GetAllActionsRqst,
    GetAllActionsRsp,
    AddRoleActionRqst,
    AddRoleActionRsp,
    RemoveRoleActionRqst,
    CreateRoleRqst,
    Role,
    CreateRoleRsp,
    DeleteRoleRqst,
    RefreshTokenRqst,
    RefreshTokenRsp,
    GetAllFilesInfoRqst,
    GetAllFilesInfoRsp,
    GetAllApplicationsInfoRqst,
    GetAllApplicationsInfoRsp,
    AddApplicationActionRqst,
    AddApplicationActionRsp,
    RemoveApplicationActionRqst,
    DeleteApplicationRqst,
    DeleteAccountRqst,
    AddAccountRoleRqst,
    RemoveAccountRoleRqst,
    GetPermissionsRqst,
    GetPermissionsRsp,
    DeletePermissionsRqst,
    DeletePermissionsRsp,
    SetPermissionRqst,
    RessourcePermission,
    SetPermissionRsp,
    SynchronizeLdapRqst,
    LdapSyncInfos,
    SynchronizeLdapRsp,
    UserSyncInfos,
    GroupSyncInfos,
    GetRessourceOwnersRqst,
    GetRessourceOwnersRsp,
    SetRessourceOwnerRqst,
    DeleteRessourceOwnerRqst,
    GetLogRqst,
    LogInfo,
    ClearAllLogRqst,
    LogType,
    SetActionPermissionRqst,
    Ressource,
    RemoveActionPermissionRqst,
    GetRessourcesRqst,
    RemoveRessourceRqst,
    DeleteLogRqst,
    GetPeersRqst,
    GetPeersRsp,
    Peer,
    AddPeerActionRqst,
    AddPeerActionRsp,
    RemovePeerActionRqst,
    DeletePeerRqst
} from "./ressource/ressource_pb";

import * as jwt from "jwt-decode";
import {
    InsertOneRqst,
    FindOneRqst,
    FindRqst,
    FindResp,
    FindOneResp,
    AggregateRqst,
    PingConnectionRqst,
    PingConnectionRsp,
    ReplaceOneRqst,
    ReplaceOneRsp
} from "./persistence/persistencepb/persistence_pb";

import {
    FindServicesDescriptorRequest,
    FindServicesDescriptorResponse,
    ServiceDescriptor,
    GetServiceDescriptorRequest,
    GetServiceDescriptorResponse,
    GetServicesDescriptorRequest,
    GetServicesDescriptorResponse,
    SetServiceDescriptorRequest
} from "./services/services_pb";

import {
    RenameRequest,
    RenameResponse,
    DeleteFileRequest,
    DeleteDirRequest,
    CreateArchiveRequest,
    CreateArchiveResponse,
    CreateDirRequest,
    ReadDirRequest,
} from "./file/filepb/file_pb";

import { TagType, ReadTagRqst } from "./plc/plcpb/plc_pb";
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
export function readFullConfig(
    globular: Globular,
    application: string,
    domain: string,
    callback: (config: IConfig) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new GetConfigRequest();
    if (globular.adminService !== undefined) {
        globular.adminService
            .getFullConfig(rqst, {
                token: <string>localStorage.getItem("user_token"),
                application: application, domain: domain
            })
            .then(rsp => {
                globular.config = JSON.parse(rsp.getResult());; // set the globular config with the full config.
                callback(globular.config);
            })
            .catch(err => {
                errorCallback(err);
            });
    }
}

// Save the configuration.
export function saveConfig(
    globular: Globular,
    application: string,
    domain: string,
    config: IConfig,
    callback: (config: IConfig) => void
    , errorCallback: (err: any) => void
) {
    let rqst = new SaveConfigRequest();
    rqst.setConfig(JSON.stringify(config));
    if (globular.adminService !== undefined) {
        globular.adminService
            .saveConfig(rqst, {
                token: <string>localStorage.getItem("user_token"),
                application: application, domain: domain
            })
            .then(rsp => {
                config = JSON.parse(rsp.getResult());
                callback(config);
            })
            .catch(err => {
                errorCallback(err);
            });
    }
}

///////////////////////////////////// Ressource & Permissions operations /////////////////////////////////

export function readAllActionPermission(
    globular: Globular,
    application: string,
    domain: string,
    callback: (results: any) => void,
    errorCallback: (err: any) => void) {
    let database = "local_ressource";
    let collection = "ActionPermission";

    let rqst = new FindRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setOptions("");
    rqst.setQuery("{}");

    // call persist data
    let stream = globular.persistenceService.find(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    let results = new Array();

    // Get the stream and set event on it...
    stream.on("data", rsp => {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });

    stream.on("status", status => {
        if (status.code == 0) {
            callback(results);
        } else {
            errorCallback({ "message": status.details })
        }
    });
}

export function getRessources(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    name: string,
    callback: (results: Ressource[]) => void,
    errorCallback: (err: any) => void) {

    let rqst = new GetRessourcesRqst
    rqst.setPath(path)
    rqst.setName(name)

    // call persist data
    let stream = globular.ressourceService.getRessources(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    });

    let results = new Array<Ressource>();

    // Get the stream and set event on it...
    stream.on("data", rsp => {
        results = results.concat(rsp.getRessourcesList())
    });

    stream.on("status", status => {
        if (status.code == 0) {
            callback(results);
        } else {
            errorCallback({ "message": status.details })
        }
    });
}

export function setActionPermission(
    globular: Globular,
    application: string,
    domain: string,
    action: string,
    permission: number,
    callback: (results: any) => void,
    errorCallback: (err: any) => void) {
    let rqst = new SetActionPermissionRqst
    rqst.setAction(action)
    rqst.setPermission(permission)

    // Call set action permission.
    globular.ressourceService.setActionPermission(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch((err: any) => {
            errorCallback(err)
        })
}

export function removeActionPermission(
    globular: Globular,
    application: string,
    domain: string,
    action: string,
    callback: (results: any) => void,
    errorCallback: (err: any) => void) {
    let rqst = new RemoveActionPermissionRqst
    rqst.setAction(action)
    // Call set action permission.
    globular.ressourceService.removeActionPermission(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch((err: any) => {
            errorCallback(err)
        })
}

export function removeRessource(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    name: string,
    callback: () => void,
    errorCallback: (err: any) => void) {
    let rqst = new RemoveRessourceRqst
    let ressource = new Ressource
    ressource.setPath(path)
    ressource.setName(name)
    rqst.setRessource(ressource)
    globular.ressourceService.removeRessource(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch((err: any) => {
            errorCallback(err)
        })

}

/**
 * Retreive the list of ressource owner.
 * @param path 
 * @param callback 
 * @param errorCallback 
 */
export function getRessourceOwners(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    callback: (infos: Array<any>) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new GetRessourceOwnersRqst
    path = path.replace("/webroot", "");
    rqst.setPath(path);

    globular.ressourceService.getRessourceOwners(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then((rsp: GetRessourceOwnersRsp) => {
        callback(rsp.getOwnersList())
    }).catch((err: any) => {
        errorCallback(err);
    });

}

/**
 * The ressource owner to be set.
 * @param path The path of the ressource
 * @param owner The owner of the ressource
 * @param callback The success callback
 * @param errorCallback The error callback
 */
export function setRessourceOwners(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    owner: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {

    let rqst = new SetRessourceOwnerRqst
    path = path.replace("/webroot", ""); // remove the /webroot part.
    rqst.setPath(path);
    rqst.setOwner(owner);

    globular.ressourceService.setRessourceOwner(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(() => {
        callback()
    }).catch((err: any) => {
        errorCallback(err);
    });

}

/**
 * Delete a given ressource owner
 * @param path The path of the ressource.
 * @param owner The owner to be remove
 * @param callback The sucess callback
 * @param errorCallback The error callback
 */
export function deleteRessourceOwners(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    owner: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {

    let rqst = new DeleteRessourceOwnerRqst
    path = path.replace("/webroot", ""); // remove the /webroot part.
    rqst.setPath(path);
    rqst.setOwner(owner);

    globular.ressourceService.deleteRessourceOwner(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(() => {
        callback()
    }).catch((err: any) => {
        errorCallback(err);
    });

}

/**
 * Retreive the permission for a given file.
 * @param path 
 * @param callback 
 * @param errorCallback 
 */
export function getRessourcePermissions(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    callback: (infos: Array<any>) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new GetPermissionsRqst
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);

    globular.ressourceService
        .getPermissions(rqst, {
            application: application, domain: domain
        })
        .then((rsp: GetPermissionsRsp) => {
            let permissions = JSON.parse(rsp.getPermissions())
            callback(permissions);

        })
        .catch(error => {
            if (errorCallback != undefined) {
                errorCallback(error);
            }
        });
}

/**
 * The permission can be assigned to 
 * a User, a Role or an Application.
 */
export enum OwnerType {
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
export function setRessourcePermission(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    owner: string,
    ownerType: OwnerType,
    number: number,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new SetPermissionRqst
    path = path.replace("/webroot", ""); // remove the /webroot part.

    if (path.length == 0) {
        path = "/";
    }

    let permission = new RessourcePermission
    permission.setPath(path)
    permission.setNumber(number)
    if (ownerType == OwnerType.User) {
        permission.setUser(owner)
    } else if (ownerType == OwnerType.Role) {
        permission.setRole(owner)
    } else if (ownerType == OwnerType.Application) {
        permission.setApplication(owner)
    }

    rqst.setPermission(permission)

    globular.ressourceService
        .setPermission(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: SetPermissionRsp) => {
            callback();
        })
        .catch(error => {
            if (errorCallback != undefined) {
                errorCallback(error);
            }
        });
}

/**
 * Delete a file permission for a give user.
 * @param path The path of the file on the server.
 * @param owner The owner of the file
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
export function deleteRessourcePermissions(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    owner: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {

    let rqst = new DeletePermissionsRqst
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }

    rqst.setPath(path);
    rqst.setOwner(owner);

    globular.ressourceService
        .deletePermissions(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: DeletePermissionsRsp) => {
            callback();
        })
        .catch(error => {
            if (errorCallback != undefined) {
                errorCallback(error);
            }
        });
}

///////////////////////////////////// File operations /////////////////////////////////

/**
 * Return server files operations.
 * @param globular 
 * @param application 
 * @param domain 
 * @param callbak 
 * @param errorCallback 
 */
export function getAllFilesInfo(
    globular: Globular,
    application: string,
    domain: string,
    callbak: (filesInfo: any) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new GetAllFilesInfoRqst();
    globular.ressourceService
        .getAllFilesInfo(rqst, { application: application, domain: domain })
        .then((rsp: GetAllFilesInfoRsp) => {
            let filesInfo = JSON.parse(rsp.getResult());
            callbak(filesInfo);
        })
        .catch((err: any) => {
            errorCallback(err);
        });
}

/**
 * Rename a file or a directorie with given name.
 * @param path The path inside webroot
 * @param newName The new file name
 * @param oldName  The old file name
 * @param callback  The success callback.
 * @param errorCallback The error callback.
 */
export function renameFile(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    newName: string,
    oldName: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new RenameRequest();
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);
    rqst.setOldName(oldName);
    rqst.setNewName(newName);

    globular.fileService
        .rename(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: RenameResponse) => {
            callback();
        })
        .catch(error => {
            if (errorCallback != undefined) {
                errorCallback(error);
            }
        });
}

/**
 * Delete a file with a given path.
 * @param path The path of the file to be deleted.
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
export function deleteFile(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new DeleteFileRequest();
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);

    globular.fileService
        .deleteFile(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: RenameResponse) => {
            callback();
        })
        .catch(error => {
            if (errorCallback != undefined) {
                errorCallback(error);
            }
        });
}

/**
 * 
 * @param path The path of the directory to be deleted.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export function deleteDir(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new DeleteDirRequest();
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);
    globular.fileService
        .deleteDir(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: RenameResponse) => {
            callback();
        })
        .catch(error => {
            if (errorCallback != undefined) {
                errorCallback(error);
            }
        });
}

/**
 * Create a dir archive.
 * @param path 
 * @param name 
 * @param callback 
 * @param errorCallback 
 */
export function createArchive(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    name: string,
    callback: (path: string) => void,
    errorCallback: (err: any) => void) {
    let rqst = new CreateArchiveRequest;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path)
    rqst.setName(name)

    globular.fileService.createAchive(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(
        (rsp: CreateArchiveResponse) => {
            callback(rsp.getResult())
        }
    ).catch(error => {
        if (errorCallback != undefined) {
            errorCallback(error);
        }
    });
}

/**
 * 
 * @param urlToSend 
 */
export function downloadFileHttp(
    globular: Globular,
    application: string,
    domain: string,
    urlToSend: string,
    fileName: string,
    callback: () => void) {
    var req = new XMLHttpRequest();
    req.open("GET", urlToSend, true);

    // Set the token to manage downlaod access.
    req.setRequestHeader("token", <string>localStorage.getItem("user_token"))
    req.setRequestHeader("application", application)
    req.setRequestHeader("domain", domain)

    req.responseType = "blob";
    req.onload = function (event) {
        var blob = req.response;
        var link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link.download = fileName;
        link.click();
        callback();
    };

    req.send();
}

/**
 * Download a directory as archive file. (.tar.gz)
 * @param path The path of the directory to dowload.
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
export function downloadDir(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    callback: () => void,
    errorCallback: (err: any) => void) {
    let name = path.split("/")[path.split("/").length - 1]
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }

    // Create an archive-> download it-> delete it...
    createArchive(globular, application, domain, path, name, (path: string) => {
        let name = path.split("/")[path.split("/").length - 1]
        // display the archive path...
        downloadFileHttp(globular, application, domain, window.location.origin + path, name, () => {
            // Here the file was downloaded I will now delete it.
            setTimeout(() => {
                // wait a little and remove the archive from the server.
                let rqst = new DeleteFileRequest
                rqst.setPath(path)
                globular.fileService.deleteFile(rqst, {
                    token: <string>localStorage.getItem("user_token"),
                    application: application, domain: domain
                }).then(callback)
                    .catch(errorCallback)
            }, 5000); // wait 5 second, arbritary...
        });

    }, errorCallback)
}

// Merge tow array together.
function mergeTypedArraysUnsafe(a: any, b: any) {
    var c = new a.constructor(a.length + b.length);
    c.set(a);
    c.set(b, a.length);
    return c;
}

/**
 * Read the content of a dir from a given path.
 * @param path The parent path of the dir to be read.
 * @param callback  Return the path of the dir with more information.
 * @param errorCallback Return a error if the file those not contain the value.
 */
export function readDir(
    globular: Globular,
    application: string,
    domain: string,
    path: string, callback: (dir: any) => void, errorCallback: (err: any) => void) {
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }

    let rqst = new ReadDirRequest
    rqst.setPath(path)
    rqst.setRecursive(true)
    rqst.setThumnailheight(256)
    rqst.setThumnailwidth(256)

    var uint8array = new Uint8Array(0);

    var stream = globular.fileService.readDir(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    });

    stream.on("data", rsp => {
        uint8array = mergeTypedArraysUnsafe(uint8array, rsp.getData())
    });

    stream.on("status", function (status) {
        if (status.code == 0) {
            var jsonStr = new TextDecoder("utf-8").decode(uint8array);
            var content = JSON.parse(jsonStr)
            callback(content)
        } else {
            // error here...
            errorCallback({ "message": status.details })
        }
    });

}

/**
 * 
 * @param files 
 */
function fileExist(fileName: string, files: Array<any>): boolean {
    if (files != null) {
        for (var i = 0; i < files.length; i++) {
            if (files[i].Name == fileName) {
                return true;
            }
        }
    }
    return false
}

/**
 * Create a new directory inside existing one.
 * @param path The path of the directory
 * @param callback The callback
 * @param errorCallback The error callback
 */
export function createDir(
    globular: Globular,
    application: string,
    domain: string,
    path: string,
    callback: (dirName: string) => void,
    errorCallback: (err: any) => void) {
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    // first of all I will read the directory content...
    readDir(globular, application, domain, path, (dir: any) => {
        let newDirName = "New Folder"
        for (var i = 0; i < 1024; i++) {
            if (!fileExist(newDirName, dir.Files)) {
                break
            }
            newDirName = "New Folder (" + i + ")"
        }

        // Set the request.
        let rqst = new CreateDirRequest
        rqst.setPath(path)
        rqst.setName(newDirName)

        // Create a directory at the given path.
        globular.fileService.createDir(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        }).then(() => {
            // The new directory was created.
            callback(newDirName)
        })
            .catch(
                (err: any) => {
                    errorCallback(err)
                }
            )

    }, errorCallback)
}

///////////////////////////////////// Time series Query //////////////////////////////////////

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
export function queryTs(
    globular: Globular,
    application: string,
    domain: string,
    connectionId: string,
    query: string,
    ts: number,
    callback: (value: any) => void,
    errorCallback: (error: any) => void
) {
    // Create a new request.
    var request = new QueryRequest();
    request.setConnectionid(connectionId);
    request.setQuery(query);
    request.setTs(ts);

    // Now I will test with promise
    globular.monitoringService
        .query(request, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then(resp => {
            if (callback != undefined) {
                callback(JSON.parse(resp.getValue()));
            }
        })
        .catch(error => {
            if (errorCallback != undefined) {
                errorCallback(error);
            }
        });
}

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
export function queryTsRange(
    globular: Globular,
    application: string,
    domain: string,
    connectionId: string,
    query: string,
    startTime: number,
    endTime: number,
    step: number,
    callback: (values: any) => void,
    errorCallback: (err: any) => void
) {
    // Create a new request.
    var request = new QueryRangeRequest();
    request.setConnectionid(connectionId);
    request.setQuery(query);
    request.setStarttime(startTime);
    request.setEndtime(endTime);
    request.setStep(step);

    let buffer = { value: "", warning: "" };

    var stream = globular.monitoringService.queryRange(request, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    stream.on("data", rsp => {
        buffer.value += rsp.getValue();
        buffer.warning = rsp.getWarnings();
    });

    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(JSON.parse(buffer.value));
        } else {
            errorCallback({ "message": status.details })
        }
    });

    stream.on("end", () => {
        // stream end signal
    });
}

///////////////////////////////////// Account management action //////////////////////////////////////

/**
 * Register a new user.
 * @param userName The name of the account
 * @param email The email
 * @param password The password
 * @param confirmPassword
 * @param callback
 * @param errorCallback
 */
export function registerAccount(
    globular: Globular,
    application: string,
    domain: string,
    userName: string,
    email: string,
    password: string,
    confirmPassword: string,
    callback: (value: any) => void,
    errorCallback: (err: any) => void
) {
    var request = new RegisterAccountRqst();
    var account = new Account();
    account.setName(userName);
    account.setEmail(email);
    request.setAccount(account);
    request.setPassword(password);
    request.setConfirmPassword(confirmPassword);

    // Create the user account.
    globular.ressourceService
        .registerAccount(request, { application: application, domain: domain })
        .then(rsp => {
            callback(rsp.getResult());
        })
        .catch(err => {
            errorCallback(err);
        });
}

/**
 * Remove an account from the server.
 * @param name  The _id of the account.
 * @param callback The callback when the action succed
 * @param errorCallback The error callback.
 */
export function DeleteAccount(
    globular: Globular,
    application: string,
    domain: string,
    id: string,
    callback: (value: any) => void,
    errorCallback: (err: any) => void) {
    let rqst = new DeleteAccountRqst
    rqst.setId(id)

    // Remove the account from the database.
    globular.ressourceService
        .deleteAccount(rqst, { token: <string>localStorage.getItem("user_token"), application: application, domain: domain })
        .then(rsp => {
            callback(rsp.getResult());
        })
        .catch(err => {
            errorCallback(err);
        });
}

/**
 * Remove a role from an account.
 * @param accountId The account id
 * @param roleId The role name (id)
 * @param callback The success callback
 * @param errorCallback The error callback
 */
export function RemoveRoleFromAccount(
    globular: Globular,
    application: string,
    domain: string,
    accountId: string,
    roleId: string,
    callback: (value: any) => void,
    errorCallback: (err: any) => void
) {

    let rqst = new RemoveAccountRoleRqst
    rqst.setAccountid(accountId)
    rqst.setRoleid(roleId)

    globular.ressourceService
        .removeAccountRole(rqst, { token: <string>localStorage.getItem("user_token"), application: application, domain: domain })
        .then(rsp => {
            callback(rsp.getResult());
        })
        .catch(err => {
            errorCallback(err);
        });

}

/**
 * Append a role to an account.
 * @param accountId The account id
 * @param roleId The role name (id)
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export function AppendRoleToAccount(
    globular: Globular,
    application: string,
    domain: string,
    accountId: string,
    roleId: string,
    callback: (value: any) => void,
    errorCallback: (err: any) => void
) {

    let rqst = new AddAccountRoleRqst
    rqst.setAccountid(accountId)
    rqst.setRoleid(roleId)

    globular.ressourceService
        .addAccountRole(rqst, { token: <string>localStorage.getItem("user_token"), application: application, domain: domain })
        .then(rsp => {
            callback(rsp.getResult());
        })
        .catch(err => {
            errorCallback(err);
        });

}

/**
 * Update the account email
 * @param accountId The account id
 * @param old_email the old email
 * @param new_email the new email
 * @param callback  the callback when success
 * @param errorCallback the error callback in case of error
 */
export function updateAccountEmail(
    globular: Globular,
    application: string,
    domain: string,
    accountId: string,
    old_email: string,
    new_email: string,
    callback: () => void,
    errorCallback: (err: any) => void) {
    let rqst = new SetEmailRequest
    rqst.setAccountid(accountId)
    rqst.setOldemail(old_email)
    rqst.setNewemail(new_email)

    globular.adminService.setEmail(rqst,
        {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        }).then(rsp => {
            callback()
        })
        .catch(err => {
            console.log("fail to save config ", err);
        });
}

/**
 * The update account password
 * @param accountId The account id
 * @param old_password The old password
 * @param new_password The new password
 * @param confirm_password The new password confirmation
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export function updateAccountPassword(
    globular: Globular,
    application: string,
    domain: string,
    accountId: string,
    old_password: string,
    new_password: string,
    confirm_password: string,
    callback: () => void,
    errorCallback: (err: any) => void) {
    let rqst = new SetPasswordRequest
    rqst.setAccountid(accountId)
    rqst.setOldpassword(old_password)
    rqst.setNewpassword(new_password)

    if (confirm_password != new_password) {
        errorCallback("password not match!")
        return
    }

    globular.adminService.setPassword(rqst,
        {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        }).then(rsp => {
            callback()
        })
        .catch(error => {
            if (errorCallback != undefined) {
                errorCallback(error);
            }
        });
}

/**
 * Authenticate the user and get the token
 * @param userName The account name or email
 * @param password  The user password
 * @param callback
 * @param errorCallback
 */
export function authenticate(
    globular: Globular,
    eventHub: EventHub,
    application: string,
    domain: string,
    userName: string,
    password: string,
    callback: (value: any) => void,
    errorCallback: (err: any) => void
) {
    var rqst = new AuthenticateRqst();
    rqst.setName(userName);
    rqst.setPassword(password);

    // Create the user account.
    globular.ressourceService
        .authenticate(rqst, { application: application, domain: domain })
        .then(rsp => {
            // Here I will set the token in the localstorage.
            let token = rsp.getToken();
            let decoded = jwt(token);

            // here I will save the user token and user_name in the local storage.
            localStorage.setItem("user_token", token);
            localStorage.setItem("user_name", (<any>decoded).username);

            readFullConfig(globular, application, domain, (config: any) => {
                // Publish local login event.
                eventHub.publish("onlogin", config, true); // return the full config...
                callback(decoded);
            }, (err: any) => {
                errorCallback(err)
            })
        })
        .catch(err => {
            console.log(err)
            errorCallback(err);
        });
}

/**
 * Function to be use to refresh token or full configuration.
 * @param callback On success callback
 * @param errorCallback On error callback
 */
export function refreshToken(
    globular: Globular,
    eventHub: EventHub,
    application: string,
    domain: string,
    callback: (token: any) => void,
    errorCallback: (err: any) => void) {
    let rqst = new RefreshTokenRqst();
    rqst.setToken(<string>localStorage.getItem("user_token"));

    globular.ressourceService
        .refreshToken(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: RefreshTokenRsp) => {
            // Here I will set the token in the localstorage.
            let token = rsp.getToken();
            let decoded = jwt(token);

            // here I will save the user token and user_name in the local storage.
            localStorage.setItem("user_token", token);
            localStorage.setItem("user_name", (<any>decoded).username);

            readFullConfig(globular, application, domain, (config: any) => {

                // Publish local login event.
                eventHub.publish("onlogin", config, true); // return the full config...

                callback(decoded);
            }, (err: any) => {
                errorCallback(err)
            })

        })
        .catch((err: any) => {
            onerror(err);
        });
}

/**
 * Save user data into the user_data collection.
 */
export function appendUserData(
    globular: Globular,
    application: string,
    domain: string,
    data: any,
    callback: (id: string) => void) {
    let userName = <string>localStorage.getItem("user_name");
    let database = userName + "_db";
    let collection = "user_data";

    let rqst = new InsertOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setJsonstr(JSON.stringify(data));
    rqst.setOptions("");

    // call persist data
    globular.persistenceService
        .insertOne(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: any) => {
            callback(rsp.getId());
        })
        .catch((err: any) => {
            console.log(err);
        });
}

/**
 * Read user data one result at time.
 */
export function readOneUserData(
    globular: Globular,
    application: string,
    domain: string,
    query: string,
    callback: (results: any) => void
) {
    let userName = <string>localStorage.getItem("user_name");
    let database = userName + "_db";
    let collection = "user_data";

    let rqst = new FindOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setQuery(query);
    rqst.setOptions("");

    // call persist data
    globular.persistenceService
        .findOne(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: any) => {
            callback(JSON.parse(rsp.getJsonstr()));
        })
        .catch((err: any) => {
            console.log(err);
        });
}

/**
 * Read all user data.
 */
export function readUserData(
    globular: Globular,
    application: string,
    domain: string,
    query: string,
    callback: (results: any) => void,
    errorCallback: (err: any) => void) {
    let userName = <string>localStorage.getItem("user_name");
    let database = userName + "_db";
    let collection = "user_data";

    let rqst = new FindOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setQuery(query);
    rqst.setOptions("");

    // call persist data
    let stream = globular.persistenceService.find(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    let results = new Array();

    // Get the stream and set event on it...
    stream.on("data", rsp => {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });

    stream.on("status", status => {
        if (status.code == 0) {
            callback(results);
        } else {
            errorCallback({ "message": status.details })
        }
    });
}



///////////////////////////////////// Role action //////////////////////////////////////

/**
 * Return the list of all account on the server, guest and admin are new account...
 * @param callback 
 */
export function GetAllAccountsInfo(
    globular: Globular,
    application: string,
    domain: string,
    callback: (accounts: Array<any>) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new FindRqst();
    rqst.setCollection("Accounts");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setQuery("{}"); // means all values.

    var stream = globular.persistenceService.find(rqst, {
        application: application, domain: domain
    });

    let accounts = new Array<any>()

    stream.on("data", (rsp: FindResp) => {
        accounts = accounts.concat(JSON.parse(rsp.getJsonstr()))
    });

    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(accounts);
        } else {
            errorCallback({ "message": status.details })
        }
    });
}

/**
 * Retreive all available actions on the server.
 * @param callback That function is call in case of success.
 * @param errorCallback That function is call in case error.
 */
export function getAllActions(
    globular: Globular,
    application: string,
    domain: string,
    callback: (ations: Array<string>) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new GetAllActionsRqst();
    globular.ressourceService
        .getAllActions(rqst, { application: application, domain: domain })
        .then((rsp: GetAllActionsRsp) => {
            callback(rsp.getActionsList());
        })
        .catch((err: any) => {
            errorCallback(err);
        });
}

/**
 * Retreive the list of all available roles on the server.
 * @param callback That function is call in case of success.
 * @param errorCallback That function is call in case error.
 */
export function getAllRoles(
    globular: Globular,
    application: string,
    domain: string,
    callback: (roles: Array<any>) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new FindRqst();
    rqst.setCollection("Roles");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setQuery("{}"); // means all values.

    var stream = globular.persistenceService.find(rqst, {
        application: application, domain: domain
    });

    var roles = new Array<any>();

    stream.on("data", (rsp: FindResp) => {
        roles = roles.concat(JSON.parse(rsp.getJsonstr()));
    });

    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(roles);
        } else {
            errorCallback({ "message": status.details })
        }
    });
}

/**
 * Append Action to a given role.
 * @param action The action name.
 * @param role The role.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
export function AppendActionToRole(
    globular: Globular,
    application: string,
    domain: string,
    role: string,
    action: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new AddRoleActionRqst();
    rqst.setRoleid(role);
    rqst.setAction(action);

    globular.ressourceService
        .addRoleAction(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: AddRoleActionRsp) => {
            callback();
        })
        .catch((err: any) => {
            errorCallback(err);
        });
}

/**
 * Remove the action from a given role.
 * @param action The action id
 * @param role The role id
 * @param callback success callback
 * @param errorCallback error callback
 */
export function RemoveActionFromRole(
    globular: Globular,
    application: string,
    domain: string,
    role: string,
    action: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new RemoveRoleActionRqst();
    rqst.setRoleid(role);
    rqst.setAction(action);

    globular.ressourceService
        .removeRoleAction(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: AddRoleActionRsp) => {
            callback();
        })
        .catch((err: any) => {
            errorCallback(err);
        });
}

export function CreateRole(
    globular: Globular,
    application: string,
    domain: string,
    id: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new CreateRoleRqst();
    let role = new Role();
    role.setId(id);
    role.setName(id);
    rqst.setRole(role);

    globular.ressourceService
        .createRole(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: CreateRoleRsp) => {
            callback();
        })
        .catch((err: any) => {
            errorCallback(err);
        });
}

export function DeleteRole(
    globular: Globular,
    application: string,
    domain: string,
    id: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new DeleteRoleRqst();
    rqst.setRoleid(id);

    globular.ressourceService
        .deleteRole(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: CreateRoleRsp) => {
            callback();
        })
        .catch((err: any) => {
            errorCallback(err);
        });
}

///////////////////////////////////// Application operations /////////////////////////////////

export function GetAllApplicationsInfo(
    globular: Globular,
    application: string,
    domain: string,
    callback: (infos: any) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new GetAllApplicationsInfoRqst();
    globular.ressourceService
        .getAllApplicationsInfo(rqst)
        .then((rsp: GetAllApplicationsInfoRsp) => {
            let infos = JSON.parse(rsp.getResult());
            callback(infos);
        })
        .catch((err: any) => {
            errorCallback(err);
        });
}

export function AppendActionToApplication(
    globular: Globular,
    application: string,
    domain: string,
    applicationId: string,
    action: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new AddApplicationActionRqst;
    rqst.setApplicationid(applicationId)
    rqst.setAction(action)
    globular.ressourceService.addApplicationAction(rqst, { token: <string>localStorage.getItem("user_token"), application: application, domain: domain })
        .then((rsp: AddApplicationActionRsp) => {
            callback()
        })
        .catch((err: any) => {
            console.log(err)
            errorCallback(err);
        });

}

export function RemoveActionFromApplication(
    globular: Globular,
    application: string,
    domain: string,
    applicationId: string,
    action: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new RemoveApplicationActionRqst;
    rqst.setApplicationid(applicationId)
    rqst.setAction(action)
    globular.ressourceService.removeApplicationAction(rqst, { token: <string>localStorage.getItem("user_token"), application: application, domain: domain })
        .then((rsp: AddApplicationActionRsp) => {
            callback()
        })
        .catch((err: any) => {
            console.log(err)
            errorCallback(err);
        });

}

export function DeleteApplication(
    globular: Globular,
    application: string,
    domain: string,
    applicationId: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new DeleteApplicationRqst;
    rqst.setApplicationid(applicationId)
    globular.ressourceService.deleteApplication(rqst, { token: <string>localStorage.getItem("user_token"), application: application, domain: domain })
        .then((rsp: AddApplicationActionRsp) => {
            callback()
        })
        .catch((err: any) => {
            console.log(err)
            errorCallback(err);
        });

}

export function SaveApplication(
    globular: Globular,
    eventHub: EventHub,
    application_: string,
    domain: string,
    application: any,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new ReplaceOneRqst;
    rqst.setCollection("Applications");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setValue(JSON.stringify(application))
    rqst.setQuery(`{"_id":"${application._id}"}`); // means all values.

    globular.persistenceService.replaceOne(rqst, { token: <string>localStorage.getItem("user_token"), application: application_, domain: domain })
        .then((rsp: ReplaceOneRsp) => {
            eventHub.publish("update_application_info_event", JSON.stringify(application), false);
            callback()
        })
        .catch((err: any) => {
            console.log(err)
            errorCallback(err);
        });

}

///////////////////////////////////// Peers operations /////////////////////////////////

export function GetAllPeersInfo(
    globular: Globular,
    application: string,
    domain: string,
    query: string,
    callback: (peers: Peer[]) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new GetPeersRqst();
    rqst.setQuery(query)

    let peers = new Array<Peer>();

    let stream = globular.ressourceService.getPeers(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    });

    // Get the stream and set event on it...
    stream.on("data", (rsp: GetPeersRsp) => {
        peers = peers.concat(rsp.getPeersList());
    });

    stream.on("status", status => {
        if (status.code == 0) {
            callback(peers);
        } else {
            errorCallback({ "message": status.details })
        }
    });
}

export function AppendActionToPeer(
    globular: Globular,
    application: string,
    domain: string,
    id: string,
    action: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new AddPeerActionRqst;
    rqst.setPeerid(id)
    rqst.setAction(action)
    globular.ressourceService.addPeerAction(rqst, { token: <string>localStorage.getItem("user_token"), application: application, domain: domain })
        .then((rsp: AddPeerActionRsp) => {
            callback()
        })
        .catch((err: any) => {
            console.log(err)
            errorCallback(err);
        });

}

export function RemoveActionFromPeer(
    globular: Globular,
    application: string,
    domain: string,
    id: string,
    action: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new RemovePeerActionRqst;
    rqst.setPeerid(id)
    rqst.setAction(action)
    globular.ressourceService.removePeerAction(rqst, { token: <string>localStorage.getItem("user_token"), application: application, domain: domain })
        .then((rsp: AddApplicationActionRsp) => {
            callback()
        })
        .catch((err: any) => {
            console.log(err)
            errorCallback(err);
        });

}

export function DeletePeer(
    globular: Globular,
    application: string,
    domain: string,
    peer: Peer,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new DeletePeerRqst;
    rqst.setPeer(peer)
    globular.ressourceService.deletePeer(rqst, { token: <string>localStorage.getItem("user_token"), application: application, domain: domain })
        .then((rsp: AddApplicationActionRsp) => {
            callback()
        })
        .catch((err: any) => {
            console.log(err)
            errorCallback(err);
        });

}

///////////////////////////////////// Services operations /////////////////////////////////

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
export function GetServiceDescriptor(
    globular: Globular,
    application: string,
    domain: string,
    serviceId: string,
    publisherId: string,
    callback: (descriptors: Array<ServiceDescriptor>) => void, errorCallback: (err: any) => void) {
    let rqst = new GetServiceDescriptorRequest
    rqst.setServiceid(serviceId);
    rqst.setPublisherid(publisherId);

    globular.servicesDicovery.getServiceDescriptor(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then((rsp: GetServiceDescriptorResponse) => {
        callback(rsp.getResultsList())
    }).catch(
        (err: any) => {
            errorCallback(err);
        }
    );
}

/**
 * Get the list of all service descriptor hosted on a server.
 * @param globular The globular object instance
 * @param application The application name who called the function.
 * @param domain The domain where the application reside.
 * @param callback 
 * @param errorCallback 
 */
export function GetServicesDescriptor(
    globular: Globular,
    application: string,
    domain: string,
    callback: (descriptors: Array<ServiceDescriptor>) => void, errorCallback: (err: any) => void) {
    let rqst = new GetServicesDescriptorRequest

    var stream = globular.servicesDicovery.getServicesDescriptor(rqst, {
        application: application, domain: domain
    });

    let descriptors = new Array<ServiceDescriptor>()

    stream.on("data", (rsp: GetServicesDescriptorResponse) => {
        descriptors = descriptors.concat(rsp.getResultsList())
    });

    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(descriptors);
        } else {
            errorCallback({ "message": status.details })
        }
    });
}

export function SetServicesDescriptor(
    globular: Globular,
    application: string,
    domain: string,
    descriptor: ServiceDescriptor, callback: () => void, errorCallback: (err: any) => void) {

    let rqst = new SetServiceDescriptorRequest
    rqst.setDescriptor(descriptor);

    globular.servicesDicovery.setServiceDescriptor(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch(
            (err: any) => {
                errorCallback(err);
            }
        );
}

/**
 * Find services by keywords.
 * @param query
 * @param callback
 */
export function findServices(
    globular: Globular,
    application: string,
    domain: string,
    keywords: Array<string>,
    callback: (results: Array<ServiceDescriptor>) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new FindServicesDescriptorRequest();
    rqst.setKeywordsList(keywords);

    // Find services by keywords.
    globular.servicesDicovery
        .findServices(rqst, { application: application, domain: domain })
        .then((rsp: FindServicesDescriptorResponse) => {
            let results = rsp.getResultsList()
            callback(results);
        })
        .catch(
            (err: any) => {
                errorCallback(err);
            }
        );
}

export function installService(
    globular: Globular,
    application: string,
    domain: string,
    discoveryId: string,
    serviceId: string,
    publisherId: string,
    version: string,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new InstallServiceRequest();
    rqst.setPublisherid(publisherId);
    rqst.setDicorveryid(discoveryId);
    rqst.setServiceid(serviceId);
    rqst.setVersion(version);

    // Install the service.
    globular.adminService
        .installService(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: InstallServiceResponse) => {
            readFullConfig(globular, application, domain, callback, errorCallback)
        }).catch(
            (err: any) => {
                errorCallback(err);
            }
        );
}

/**
 * Stop a service.
 */
export function stopService(
    globular: Globular,
    application: string,
    domain: string,
    serviceId: string,
    callback: () => void,
    errorCallback: (err: any) => void) {
    let rqst = new StopServiceRequest();
    rqst.setServiceId(serviceId);
    globular.adminService
        .stopService(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then(() => {
            callback();
        })
        .catch((err: any) => {
            errorCallback(err);
        });
}

/**
 * Start a service
 * @param serviceId The id of the service to start.
 * @param callback  The callback on success.
 */
export function startService(
    globular: Globular,
    application: string,
    domain: string,
    serviceId: string,
    callback: () => void,
    errorCallback: (err: any) => void) {
    let rqst = new StartServiceRequest();
    rqst.setServiceId(serviceId);
    globular.adminService
        .startService(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then(() => {
            callback();
        })
        .catch((err: any) => {
            errorCallback(err)
        });
}

/**
 * Here I will save the service configuration.
 * @param service The configuration to save.
 */
export function saveService(
    globular: Globular,
    application: string,
    domain: string,
    service: IServiceConfig,
    callback: (config: any) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new SaveConfigRequest();

    rqst.setConfig(JSON.stringify(service));
    globular.adminService
        .saveConfig(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: SaveConfigResponse) => {
            // The service with updated values...
            let service = JSON.parse(rsp.getResult());
            callback(service);
        })
        .catch((err: any) => {
            errorCallback(err)
        });
}

export function uninstallService(
    globular: Globular,
    application: string,
    domain: string,
    service: IServiceConfig,
    deletePermissions: boolean,
    callback: () => void,
    errorCallback: (err: any) => void
) {
    let rqst = new UninstallServiceRequest
    rqst.setServiceid(service.Id)
    rqst.setPublisherid(service.PublisherId)
    rqst.setVersion(service.Version)
    rqst.setDeletepermissions(deletePermissions)

    globular.adminService
        .uninstallService(rqst, {
            token: <string>localStorage.getItem("user_token"),
            application: application, domain: domain
        })
        .then((rsp: UninstallServiceResponse) => {
            delete globular.config.Services[service.Id]
            // The service with updated values...
            callback();
        })
        .catch((err: any) => {
            errorCallback(err)
        });
}

/**
 * Return the list of service bundles.
 * @param callback 
 */
export function GetServiceBundles(
    globular: Globular,
    application: string,
    domain: string,
    publisherId: string,
    serviceId: string,
    version: string, callback: (
        bundles: Array<any>) => void,
    errorCallback: (err: any) => void
) {
    let rqst = new FindRqst();
    rqst.setCollection("ServiceBundle");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setQuery(`{}`); // means all values.

    var stream = globular.persistenceService.find(rqst, {
        application: application, domain: domain
    });

    let bundles = new Array<any>()

    stream.on("data", (rsp: FindResp) => {
        bundles = bundles.concat(JSON.parse(rsp.getJsonstr()))
    });

    stream.on("status", function (status) {
        if (status.code == 0) {
            // filter localy.
            callback(bundles.filter(bundle => String(bundle._id).startsWith(publisherId + '%' + serviceId + '%' + version)));
        } else {
            errorCallback({ "message": status.details })
        }
    });
}

// Get the object pointed by a reference.
export function getReferencedValue(
    globular: Globular,
    application: string,
    domain: string,
    ref: any,
    callback: (results: any) => void,
    errorCallback: (err: any) => void) {

    let database = ref.$db;
    let collection = ref.$ref;

    let rqst = new FindOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setQuery(`{"_id":"${ref.$id}"}`);
    rqst.setOptions("");

    globular.persistenceService.findOne(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then((rsp: FindOneResp) => {
        callback(JSON.parse(rsp.getJsonstr()))
    }).catch((err: any) => {
        errorCallback(err)
    });
}


/**
 * Read all errors data.
 * @param callback 
 */
export function readErrors(
    globular: Globular,
    application: string,
    domain: string,
    callback: (results: any) => void,
    errorCallback: (err: any) => void) {
    let database = "local_ressource";
    let collection = "Logs";

    let rqst = new FindOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setOptions("");
    rqst.setQuery("{}");

    // call persist data
    let stream = globular.persistenceService.find(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    let results = new Array();

    // Get the stream and set event on it...
    stream.on("data", rsp => {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });

    stream.on("status", status => {
        if (status.code == 0) {
            callback(results);
        } else {
            errorCallback({ "message": status.details })
        }
    });
}



///////////////////////////// Logging Operations ////////////////////////////////////////

/**
 * Read all logs
 * @param callback The success callback.
 */
export function readLogs(
    globular: Globular,
    application: string,
    domain: string,
    query: string,
    callback: (results: any) => void,
    errorCallback: (err: any) => void) {

    let rqst = new GetLogRqst();
    rqst.setQuery(query);

    // call persist data
    let stream = globular.ressourceService.getLog(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    });

    let results = new Array<LogInfo>();

    // Get the stream and set event on it...
    stream.on("data", rsp => {
        results = results.concat(rsp.getInfoList());
    });

    stream.on("status", status => {
        if (status.code == 0) {
            results = results.sort((t1, t2) => {
                const name1 = t1.getDate();
                const name2 = t2.getDate();
                if (name1 < name2) { return 1; }
                if (name1 > name2) { return -1; }
                return 0;
            });

            callback(results);
        } else {
            console.log(status.details)
            errorCallback({ "message": status.details })
        }
    });
}

export function clearAllLog(
    globular: Globular,
    application: string,
    domain: string,
    logType: LogType,
    callback: () => void,
    errorCallback: (err: any) => void) {
    let rqst = new ClearAllLogRqst
    rqst.setType(logType)
    globular.ressourceService.clearAllLog(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch((err: any) => {
            errorCallback(err)
        })
}

export function deleteLog(
    globular: Globular,
    application: string,
    domain: string,
    log: LogInfo,
    callback: () => void,
    errorCallback: (err: any) => void) {
    let rqst = new DeleteLogRqst
    rqst.setLog(log)
    globular.ressourceService.deleteLog(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch((err: any) => {
            errorCallback(err)
        })
}

/**
 * Return the logged method and their count.
 * @param pipeline
 * @param callback 
 * @param errorCallback 
 */
export function getNumbeOfLogsByMethod(
    globular: Globular,
    application: string,
    domain: string,
    callback: (resuts: Array<any>) => void,
    errorCallback: (err: any) => void) {

    let database = "local_ressource";
    let collection = "Logs";
    let rqst = new AggregateRqst
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setOptions("");

    let pipeline = `[{"$group":{"_id":{"method":"$method"}, "count":{"$sum":1}}}]`;

    rqst.setPipeline(pipeline);

    // call persist data
    let stream = globular.persistenceService.aggregate(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    let results = new Array();

    // Get the stream and set event on it...
    stream.on("data", rsp => {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });

    stream.on("status", status => {
        if (status.code == 0) {
            callback(results);
        } else {
            errorCallback({ "message": status.details })
        }
    });

}

//////////////////////////// PLC Operations ///////////////////////////////////
export enum PLC_TYPE {
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
export async function readPlcTag(
    globular: Globular,
    application: string,
    domain: string,
    plcType: PLC_TYPE,
    connectionId: string,
    name: string,
    type: TagType, offset: number) {
    let rqst = new ReadTagRqst();
    rqst.setName(name)
    rqst.setType(type)
    rqst.setOffset(offset)
    rqst.setConnectionId(connectionId)
    let result: string

    // Try to get the value from the server.
    try {
        if (plcType == PLC_TYPE.ALEN_BRADLEY) {
            if (globular.plcService_ab != undefined) {
                let rsp = await globular.plcService_ab.readTag(rqst, {
                    token: <string>localStorage.getItem("user_token"),
                    application: application, domain: domain
                });
                result = rsp.getValues()
            } else {
                return "No Alen Bradlay PLC server configured!"
            }
        } else if (plcType == PLC_TYPE.SIEMENS) {
            if (globular.plcService_siemens != undefined) {
                let rsp = await globular.plcService_siemens.readTag(rqst, {
                    token: <string>localStorage.getItem("user_token"),
                    application: application, domain: domain
                });
                result = rsp.getValues()
            } else {
                return "No Siemens PLC server configured!"
            }
        } else {
            return "No PLC server configured!"
        }
    } catch (err) {
        return err
    }

    // Here I got the value in a string I will convert it into it type.
    if (type == TagType.BOOL) {
        return result == "true" ? true : false
    } else if (type == TagType.REAL) {
        return parseFloat(result)
    } else { // Must be cinsidere a integer.
        return parseInt(result)
    }
}

///////////////////////////////////// LDAP operations /////////////////////////////////

/**
 * Synchronize LDAP and Globular/MongoDB user and roles.
 * @param info The synchronisations informations.
 * @param callback success callback.
 */
export function syncLdapInfos(
    globular: Globular,
    application:
        string, domain: string,
    info: any, timeout: number, callback: () => void, errorCallback: (err: any) => void) {
    let rqst = new SynchronizeLdapRqst
    let syncInfos = new LdapSyncInfos
    syncInfos.setConnectionid(info.connectionId)
    syncInfos.setLdapseriveid(info.ldapSeriveId)
    syncInfos.setRefresh(info.refresh)

    let userSyncInfos = new UserSyncInfos
    userSyncInfos.setBase(info.userSyncInfos.base)
    userSyncInfos.setQuery(info.userSyncInfos.query)
    userSyncInfos.setId(info.userSyncInfos.id)
    userSyncInfos.setEmail(info.userSyncInfos.email)
    syncInfos.setUsersyncinfos(userSyncInfos)

    let groupSyncInfos = new GroupSyncInfos
    groupSyncInfos.setBase(info.groupSyncInfos.base)
    groupSyncInfos.setQuery(info.groupSyncInfos.query)
    groupSyncInfos.setId(info.groupSyncInfos.id)
    syncInfos.setGroupsyncinfos(groupSyncInfos)

    rqst.setSyncinfo(syncInfos)

    // Try to synchronyze the ldap service.
    globular.ressourceService.synchronizeLdap(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then((rsp: SynchronizeLdapRsp) => {

    }).catch((err: any) => {
        errorCallback(err);
    })
}

///////////////////////////////////// SQL operations /////////////////////////////////

/**
 * Ping globular sql service.
 * @param globular 
 * @param application 
 * @param domain 
 * @param connectionId 
 * @param callback 
 * @param errorCallback 
 */
export function pingSql(
    globular: Globular,
    application: string,
    domain: string,
    connectionId: string, callback: (pong: string) => {}, errorCallback: (err: any) => void) {
    let rqst = new PingConnectionRqst
    rqst.setId(connectionId)

    globular.sqlService.ping(rqst, {
        token: <string>localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then((rsp: PingConnectionRsp) => {
        callback(rsp.getResult());
    }).catch((err: any) => {
        errorCallback(err);
    })
}