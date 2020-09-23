"use strict";
/**
 * That file contain a list of function that can be use as API instead of
 * using the globular object itself.
 */
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.searchDocuments = exports.pingSql = exports.syncLdapInfos = exports.readPlcTag = exports.PLC_TYPE = exports.getNumbeOfLogsByMethod = exports.deleteLogEntry = exports.clearAllLog = exports.readLogs = exports.readErrors = exports.getReferencedValue = exports.GetServiceBundles = exports.uninstallService = exports.saveService = exports.startService = exports.stopService = exports.installService = exports.findServices = exports.SetServicesDescriptor = exports.GetServicesDescriptor = exports.GetServiceDescriptor = exports.SaveApplication = exports.DeleteApplication = exports.RemoveActionFromApplication = exports.AppendActionToApplication = exports.GetAllApplicationsInfo = exports.DeleteRole = exports.CreateRole = exports.RemoveActionFromRole = exports.AppendActionToRole = exports.getAllRoles = exports.getAllActions = exports.readUserData = exports.readOneUserData = exports.appendUserData = exports.refreshToken = exports.authenticate = exports.updateAccountPassword = exports.updateAccountEmail = exports.AppendRoleToAccount = exports.RemoveRoleFromAccount = exports.DeleteAccount = exports.registerAccount = exports.GetAllAccountsInfo = exports.queryTsRange = exports.queryTs = exports.createDir = exports.readDir = exports.downloadDir = exports.downloadFileHttp = exports.createArchive = exports.deleteDir = exports.deleteFile = exports.renameFile = exports.getAllFilesInfo = exports.deleteRessourcePermissions = exports.setRessourcePermission = exports.OwnerType = exports.getRessourcePermissions = exports.deleteRessourceOwners = exports.setRessourceOwners = exports.getRessourceOwners = exports.removeRessource = exports.removeActionPermission = exports.setActionPermission = exports.getRessources = exports.readAllActionPermissions = exports.saveConfig = exports.readFullConfig = void 0;
var admin_pb_1 = require("./admin/admin_pb");
var monitoring_pb_1 = require("./monitoring/monitoring_pb");
var ressource_pb_1 = require("./ressource/ressource_pb");
var jwt = require("jwt-decode");
var persistence_pb_1 = require("./persistence/persistence_pb");
var services_pb_1 = require("./services/services_pb");
var file_pb_1 = require("./file/file_pb");
var plc_pb_1 = require("./plc/plc_pb");
var search_pb_1 = require("./search/search_pb");
// Here I will get the authentication information.
var domain = window.location.hostname;
var application = window.location.pathname.split('/')[1];
var token = localStorage.getItem("user_token");
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
function readFullConfig(globular, callback, errorCallback) {
    var rqst = new admin_pb_1.GetConfigRequest();
    if (globular.adminService !== undefined) {
        globular.adminService
            .getFullConfig(rqst, {
            "token": token,
            "application": application,
            "domain": domain
        })
            .then(function (rsp) {
            globular.config = JSON.parse(rsp.getResult());
            ; // set the globular config with the full config.
            callback(globular.config);
        })
            .catch(function (err) {
            errorCallback(err);
        });
    }
}
exports.readFullConfig = readFullConfig;
/**
 * Save a configuration
 * @param globular
 * @param application
 * @param domain
 * @param config The configuration to be save.
 * @param callback
 * @param errorCallback
 */
function saveConfig(globular, config, callback, errorCallback) {
    var rqst = new admin_pb_1.SaveConfigRequest();
    rqst.setConfig(JSON.stringify(config));
    if (globular.adminService !== undefined) {
        globular.adminService
            .saveConfig(rqst, {
            "token": token,
            "application": application,
            "domain": domain
        })
            .then(function (rsp) {
            config = JSON.parse(rsp.getResult());
            callback(config);
        })
            .catch(function (err) {
            errorCallback(err);
        });
    }
}
exports.saveConfig = saveConfig;
///////////////////////////////////// Ressource & Permissions operations /////////////////////////////////
/**
 * Return the list of all action permissions.
 * Action permission are apply on ressource managed by those action.
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
function readAllActionPermissions(globular, callback, errorCallback) {
    var database = "local_ressource";
    var collection = "ActionPermission";
    var rqst = new persistence_pb_1.FindRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setOptions("");
    rqst.setQuery("{}");
    // call persist data
    var stream = globular.persistenceService.find(rqst, {
        "token": token,
        "application": application, "domain": domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            callback(results);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.readAllActionPermissions = readAllActionPermissions;
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
function getRessources(globular, path, name, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetRessourcesRqst;
    rqst.setPath(path);
    rqst.setName(name);
    // call persist data
    var stream = globular.ressourceService.getRessources(rqst, {
        "token": token,
        "application": application, "domain": domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(rsp.getRessourcesList());
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            callback(results);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.getRessources = getRessources;
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
function setActionPermission(globular, action, permission, callback, errorCallback) {
    var rqst = new ressource_pb_1.SetActionPermissionRqst;
    rqst.setAction(action);
    rqst.setPermission(permission);
    // Call set action permission.
    globular.ressourceService.setActionPermission(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.setActionPermission = setActionPermission;
/**
 * Delete action permission.
 * @param globular
 * @param application
 * @param domain
 * @param action
 * @param callback
 * @param errorCallback
 */
function removeActionPermission(globular, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveActionPermissionRqst;
    rqst.setAction(action);
    // Call set action permission.
    globular.ressourceService.removeActionPermission(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.removeActionPermission = removeActionPermission;
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
function removeRessource(globular, path, name, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveRessourceRqst;
    var ressource = new ressource_pb_1.Ressource;
    ressource.setPath(path);
    ressource.setName(name);
    rqst.setRessource(ressource);
    globular.ressourceService.removeRessource(rqst, {
        "token": token,
        "application": application,
        "domain": domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.removeRessource = removeRessource;
/**
 * Retreive the list of ressource owner.
 * @param path
 * @param callback
 * @param errorCallback
 */
function getRessourceOwners(globular, path, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetRessourceOwnersRqst;
    path = path.replace("/webroot", "");
    rqst.setPath(path);
    globular.ressourceService.getRessourceOwners(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(function (rsp) {
        callback(rsp.getOwnersList());
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.getRessourceOwners = getRessourceOwners;
/**
 * The ressource owner to be set.
 * @param path The path of the ressource
 * @param owner The owner of the ressource
 * @param callback The success callback
 * @param errorCallback The error callback
 */
function setRessourceOwners(globular, path, owner, callback, errorCallback) {
    var rqst = new ressource_pb_1.SetRessourceOwnerRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    rqst.setPath(path);
    rqst.setOwner(owner);
    globular.ressourceService.setRessourceOwner(rqst, {
        "token": token,
        "application": application,
        "domain": domain
    }).then(function () {
        callback();
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.setRessourceOwners = setRessourceOwners;
/**
 * Delete a given ressource owner
 * @param path The path of the ressource.
 * @param owner The owner to be remove
 * @param callback The sucess callback
 * @param errorCallback The error callback
 */
function deleteRessourceOwners(globular, path, owner, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteRessourceOwnerRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    rqst.setPath(path);
    rqst.setOwner(owner);
    globular.ressourceService.deleteRessourceOwner(rqst, {
        "token": token,
        "application": application,
        "domain": domain
    }).then(function () {
        callback();
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.deleteRessourceOwners = deleteRessourceOwners;
/**
 * Retreive the permission for a given ressource.
 * @param path
 * @param callback
 * @param errorCallback
 */
function getRessourcePermissions(globular, path, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetPermissionsRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length === 0) {
        path = "/";
    }
    rqst.setPath(path);
    globular.ressourceService
        .getPermissions(rqst, {
        "application": application,
        "domain": domain
    })
        .then(function (rsp) {
        var permissions = JSON.parse(rsp.getPermissions());
        callback(permissions);
    })
        .catch(function (error) {
        if (errorCallback !== undefined) {
            errorCallback(error);
        }
    });
}
exports.getRessourcePermissions = getRessourcePermissions;
/**
 * The permission can be assigned to
 * a User, a Role or an Application.
 */
var OwnerType;
(function (OwnerType) {
    OwnerType[OwnerType["User"] = 1] = "User";
    OwnerType[OwnerType["Role"] = 2] = "Role";
    OwnerType[OwnerType["Application"] = 3] = "Application";
})(OwnerType = exports.OwnerType || (exports.OwnerType = {}));
/**
 * Create a file permission.
 * @param path The path on the server from the root.
 * @param owner The owner of the permission
 * @param ownerType The owner type
 * @param number The (unix) permission number.
 * @param callback The success callback
 * @param errorCallback The error callback
 */
function setRessourcePermission(globular, path, owner, ownerType, permissionNumber, callback, errorCallback) {
    var rqst = new ressource_pb_1.SetPermissionRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length === 0) {
        path = "/";
    }
    var permission = new ressource_pb_1.RessourcePermission;
    permission.setPath(path);
    permission.setNumber(permissionNumber);
    if (ownerType === OwnerType.User) {
        permission.setUser(owner);
    }
    else if (ownerType === OwnerType.Role) {
        permission.setRole(owner);
    }
    else if (ownerType === OwnerType.Application) {
        permission.setApplication(owner);
    }
    rqst.setPermission(permission);
    globular.ressourceService
        .setPermission(rqst, {
        "token": token,
        "application": application,
        "domain": domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback !== undefined) {
            errorCallback(error);
        }
    });
}
exports.setRessourcePermission = setRessourcePermission;
/**
 * Delete a file permission for a give user.
 * @param path The path of the file on the server.
 * @param owner The owner of the file
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
function deleteRessourcePermissions(globular, path, owner, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeletePermissionsRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length === 0) {
        path = "/";
    }
    rqst.setPath(path);
    rqst.setOwner(owner);
    globular.ressourceService
        .deletePermissions(rqst, {
        "token": token,
        "application": application,
        "domain": domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback !== undefined) {
            errorCallback(error);
        }
    });
}
exports.deleteRessourcePermissions = deleteRessourcePermissions;
///////////////////////////////////// File operations /////////////////////////////////
/**
 * Return server files operations.
 * @param globular
 * @param application
 * @param domain
 * @param callbak
 * @param errorCallback
 */
function getAllFilesInfo(globular, callbak, errorCallback) {
    var rqst = new ressource_pb_1.GetAllFilesInfoRqst();
    globular.ressourceService
        .getAllFilesInfo(rqst, { "application": application, "domain": domain })
        .then(function (rsp) {
        var filesInfo = JSON.parse(rsp.getResult());
        callbak(filesInfo);
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.getAllFilesInfo = getAllFilesInfo;
/**
 * Rename a file or a directorie with given name.
 * @param path The path inside webroot
 * @param newName The new file name
 * @param oldName  The old file name
 * @param callback  The success callback.
 * @param errorCallback The error callback.
 */
function renameFile(globular, path, newName, oldName, callback, errorCallback) {
    var rqst = new file_pb_1.RenameRequest();
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length === 0) {
        path = "/";
    }
    rqst.setPath(path);
    rqst.setOldName(oldName);
    rqst.setNewName(newName);
    globular.fileService
        .rename(rqst, {
        "token": token,
        "application": application,
        "domain": domain,
        "path": path + "/" + oldName
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback !== undefined) {
            errorCallback(error);
        }
    });
}
exports.renameFile = renameFile;
/**
 * Delete a file with a given path.
 * @param path The path of the file to be deleted.
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
function deleteFile(globular, path, callback, errorCallback) {
    var rqst = new file_pb_1.DeleteFileRequest();
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length === 0) {
        path = "/";
    }
    rqst.setPath(path);
    globular.fileService
        .deleteFile(rqst, {
        "token": token,
        "application": application,
        "domain": domain,
        "path": path
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback !== undefined) {
            errorCallback(error);
        }
    });
}
exports.deleteFile = deleteFile;
/**
 * Remove a given directory and all element it contain.
 * @param path The path of the directory to be deleted.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
function deleteDir(globular, path, callback, errorCallback) {
    var rqst = new file_pb_1.DeleteDirRequest();
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length === 0) {
        path = "/";
    }
    rqst.setPath(path);
    globular.fileService
        .deleteDir(rqst, {
        "token": token,
        "application": application,
        "domain": domain,
        "path": path
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback !== undefined) {
            errorCallback(error);
        }
    });
}
exports.deleteDir = deleteDir;
/**
 * Create a dir archive.
 * @param path
 * @param name
 * @param callback
 * @param errorCallback
 */
function createArchive(globular, path, name, callback, errorCallback) {
    var rqst = new file_pb_1.CreateArchiveRequest;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length === 0) {
        path = "/";
    }
    rqst.setPath(path);
    rqst.setName(name);
    globular.fileService.createAchive(rqst, {
        "token": token,
        "application": application,
        "domain": domain,
        "path": path
    }).then(function (rsp) {
        callback(rsp.getResult());
    }).catch(function (error) {
        if (errorCallback !== undefined) {
            errorCallback(error);
        }
    });
}
exports.createArchive = createArchive;
/**
 * Download a file from the server.
 * @param urlToSend
 */
function downloadFileHttp(urlToSend, fileName, callback) {
    var req = new XMLHttpRequest();
    req.open("GET", urlToSend, true);
    // Set the token to manage downlaod access.
    req.setRequestHeader("token", token);
    req.setRequestHeader("application", application);
    req.setRequestHeader("domain", domain);
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
exports.downloadFileHttp = downloadFileHttp;
/**
 * Download a directory as archive file. (.tar.gz)
 * @param path The path of the directory to dowload.
 * @param callback The success callback.
 * @param errorCallback The error callback.
 */
function downloadDir(globular, path, callback, errorCallback) {
    var name = path.split("/")[path.split("/").length - 1];
    path = path.replace("/webroot", ""); // remove the /webroot part.
    // Create an archive-> download it-> delete it...
    createArchive(globular, path, name, function (_path) {
        // display the archive path...
        downloadFileHttp(window.location.origin + _path, name, function () {
            // Here the file was downloaded I will now delete it.
            setTimeout(function () {
                // wait a little and remove the archive from the server.
                var rqst = new file_pb_1.DeleteFileRequest;
                rqst.setPath(path + "/" + name);
                globular.fileService.deleteFile(rqst, {
                    "token": token,
                    "application": application,
                    "domain": domain,
                    "path": path
                }).then(callback)
                    .catch(errorCallback);
            }, 5000); // wait 5 second, arbritary...
        });
    }, errorCallback);
}
exports.downloadDir = downloadDir;
// Merge tow array together.
function mergeTypedArraysUnsafe(a, b) {
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
function readDir(globular, path, callback, errorCallback) {
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length === 0) {
        path = "/";
    }
    var rqst = new file_pb_1.ReadDirRequest;
    rqst.setPath(path);
    rqst.setRecursive(true);
    rqst.setThumnailheight(256);
    rqst.setThumnailwidth(256);
    var uint8array = new Uint8Array(0);
    var stream = globular.fileService.readDir(rqst, {
        "token": token,
        "application": application,
        "domain": domain,
        "path": path
    });
    stream.on("data", function (rsp) {
        uint8array = mergeTypedArraysUnsafe(uint8array, rsp.getData());
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            var content = JSON.parse(new TextDecoder("utf-8").decode(uint8array));
            callback(content);
        }
        else {
            // error here...
            errorCallback({ "message": status.details });
        }
    });
}
exports.readDir = readDir;
/**
 * Test if a file is contain in a list of files.
 * @param files
 */
function fileExist(fileName, files) {
    if (files != null) {
        for (var _i = 0, files_1 = files; _i < files_1.length; _i++) {
            var file = files_1[_i];
            if (file.Name === fileName) {
                return true;
            }
        }
    }
    return false;
}
/**
 * Create a new directory inside existing one.
 * @param path The path of the directory
 * @param callback The callback
 * @param errorCallback The error callback
 */
function createDir(globular, path, callback, errorCallback) {
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length === 0) {
        path = "/";
    }
    // first of all I will read the directory content...
    readDir(globular, path, function (dir) {
        var newDirName = "New Folder";
        for (var i = 0; i < 1024; i++) {
            if (!fileExist(newDirName, dir.Files)) {
                break;
            }
            newDirName = "New Folder (" + i + ")";
        }
        // Set the request.
        var rqst = new file_pb_1.CreateDirRequest;
        rqst.setPath(path);
        rqst.setName(newDirName);
        // Create a directory at the given path.
        globular.fileService.createDir(rqst, {
            "token": token,
            "application": application,
            "domain": domain,
            "path": path
        }).then(function () {
            // The new directory was created.
            callback(newDirName);
        })
            .catch(function (err) {
            errorCallback(err);
        });
    }, errorCallback);
}
exports.createDir = createDir;
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
function queryTs(globular, connectionId, query, ts, callback, errorCallback) {
    // Create a new request.
    var rqst = new monitoring_pb_1.QueryRequest();
    rqst.setConnectionid(connectionId);
    rqst.setQuery(query);
    rqst.setTs(ts);
    // Now I will test with promise
    globular.monitoringService
        .query(rqst, {
        "token": token,
        "application": application,
        "domain": domain
    })
        .then(function (resp) {
        if (callback !== undefined) {
            callback(JSON.parse(resp.getValue()));
        }
    })
        .catch(function (error) {
        if (errorCallback !== undefined) {
            errorCallback(error);
        }
    });
}
exports.queryTs = queryTs;
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
function queryTsRange(globular, connectionId, query, startTime, endTime, step, callback, errorCallback) {
    // Create a new request.
    var rqst = new monitoring_pb_1.QueryRangeRequest();
    rqst.setConnectionid(connectionId);
    rqst.setQuery(query);
    rqst.setStarttime(startTime);
    rqst.setEndtime(endTime);
    rqst.setStep(step);
    var buffer = { value: "", warning: "" };
    var stream = globular.monitoringService.queryRange(rqst, {
        "token": token,
        "application": application,
        "domain": domain
    });
    stream.on("data", function (rsp) {
        buffer.value += rsp.getValue();
        buffer.warning = rsp.getWarnings();
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            callback(JSON.parse(buffer.value));
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
    stream.on("end", function () {
        // stream end signal
    });
}
exports.queryTsRange = queryTsRange;
///////////////////////////////////// Account management action //////////////////////////////////////
/**
 * Return the list of all account on the server, guest and admin are new account...
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
function GetAllAccountsInfo(globular, callback, errorCallback) {
    var rqst = new persistence_pb_1.FindRqst();
    rqst.setCollection("Accounts");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setQuery("{}"); // means all values.
    var stream = globular.persistenceService.find(rqst, {
        "application": application,
        "domain": domain
    });
    var accounts = [];
    stream.on("data", function (rsp) {
        accounts = accounts.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            callback(accounts);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.GetAllAccountsInfo = GetAllAccountsInfo;
/**
 * Register a new account.
 * @param userName The name of the account
 * @param email The email
 * @param password The password
 * @param confirmPassword
 * @param callback
 * @param errorCallback
 */
function registerAccount(globular, userName, email, password, confirmPassword, callback, errorCallback) {
    var request = new ressource_pb_1.RegisterAccountRqst();
    var account = new ressource_pb_1.Account();
    account.setName(userName);
    account.setEmail(email);
    request.setAccount(account);
    request.setPassword(password);
    request.setConfirmPassword(confirmPassword);
    // Create the user account.
    globular.ressourceService
        .registerAccount(request, {
        "application": application,
        "domain": domain
    })
        .then(function (rsp) {
        callback(rsp.getResult());
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.registerAccount = registerAccount;
/**
 * Remove an account from the server.
 * @param name  The _id of the account.
 * @param callback The callback when the action succed
 * @param errorCallback The error callback.
 */
function DeleteAccount(globular, id, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteAccountRqst;
    rqst.setId(id);
    // Remove the account from the database.
    globular.ressourceService
        .deleteAccount(rqst, { "token": token, "application": application, "domain": domain })
        .then(function (rsp) {
        callback(rsp.getResult());
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.DeleteAccount = DeleteAccount;
/**
 * Remove a role from an account.
 * @param accountId The account id
 * @param roleId The role name (id)
 * @param callback The success callback
 * @param errorCallback The error callback
 */
function RemoveRoleFromAccount(globular, accountId, roleId, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveAccountRoleRqst;
    rqst.setAccountid(accountId);
    rqst.setRoleid(roleId);
    globular.ressourceService
        .removeAccountRole(rqst, { "token": token, "application": application, "domain": domain })
        .then(function (rsp) {
        callback(rsp.getResult());
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.RemoveRoleFromAccount = RemoveRoleFromAccount;
/**
 * Append a role to an account.
 * @param accountId The account id
 * @param roleId The role name (id)
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
function AppendRoleToAccount(globular, accountId, roleId, callback, errorCallback) {
    var rqst = new ressource_pb_1.AddAccountRoleRqst;
    rqst.setAccountid(accountId);
    rqst.setRoleid(roleId);
    globular.ressourceService
        .addAccountRole(rqst, { "token": token, "application": application, "domain": domain })
        .then(function (rsp) {
        callback(rsp.getResult());
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.AppendRoleToAccount = AppendRoleToAccount;
/**
 * Update the account email
 * @param accountId The account id
 * @param old_email the old email
 * @param new_email the new email
 * @param callback  the callback when success
 * @param errorCallback the error callback in case of error
 */
function updateAccountEmail(globular, accountId, oldEmail, newEmail, callback, errorCallback) {
    var rqst = new admin_pb_1.SetEmailRequest;
    rqst.setAccountid(accountId);
    rqst.setOldemail(oldEmail);
    rqst.setNewemail(newEmail);
    globular.adminService.setEmail(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.updateAccountEmail = updateAccountEmail;
/**
 * The update account password
 * @param accountId The account id
 * @param old_password The old password
 * @param new_password The new password
 * @param confirm_password The new password confirmation
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
function updateAccountPassword(globular, accountId, oldPassword, newPassword, confirmPassword, callback, errorCallback) {
    var rqst = new admin_pb_1.SetPasswordRequest;
    rqst.setAccountid(accountId);
    rqst.setOldpassword(oldPassword);
    rqst.setNewpassword(newPassword);
    if (confirmPassword !== newPassword) {
        errorCallback("password not match!");
        return;
    }
    globular.adminService.setPassword(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback !== undefined) {
            errorCallback(error);
        }
    });
}
exports.updateAccountPassword = updateAccountPassword;
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
function authenticate(globular, eventHub, userName, password, callback, errorCallback) {
    var rqst = new ressource_pb_1.AuthenticateRqst();
    rqst.setName(userName);
    rqst.setPassword(password);
    // Create the user account.
    globular.ressourceService
        .authenticate(rqst, { "application": application, "domain": domain })
        .then(function (rsp) {
        // Here I will set the token in the localstorage.
        token = rsp.getToken();
        var decoded = jwt(token);
        // here I will save the user token and user_name in the local storage.
        localStorage.setItem("user_token", token);
        localStorage.setItem("user_name", decoded.username);
        readFullConfig(globular, function (config) {
            // Publish local login event.
            eventHub.publish("onlogin", config, true); // return the full config...
            callback(decoded);
        }, function (err) {
            errorCallback(err);
        });
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.authenticate = authenticate;
/**
 * Function to be use to refresh token.
 * @param globular
 * @param eventHub
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
function refreshToken(globular, eventHub, callback, errorCallback) {
    var rqst = new ressource_pb_1.RefreshTokenRqst();
    rqst.setToken(localStorage.getItem("user_token"));
    globular.ressourceService
        .refreshToken(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        // Here I will set the token in the localstorage.
        token = rsp.getToken();
        var decoded = jwt(token);
        // here I will save the user token and user_name in the local storage.
        localStorage.setItem("user_token", token);
        localStorage.setItem("user_name", decoded.username);
        // Publish local login event.
        eventHub.publish("onlogin", globular.config, true); // return the full config...
        callback(decoded);
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.refreshToken = refreshToken;
/**
 * Save user data into the user_data collection.
 * @param globular
 * @param application
 * @param domain
 * @param data
 * @param callback
 * @param errorCallback
 */
function appendUserData(globular, data, callback, errorCallback) {
    var userName = localStorage.getItem("user_name");
    var database = userName + "_db";
    var collection = "user_data";
    var rqst = new persistence_pb_1.InsertOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setJsonstr(JSON.stringify(data));
    rqst.setOptions("");
    // call persist data
    globular.persistenceService
        .insertOne(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        callback(rsp.getId());
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.appendUserData = appendUserData;
/**
 * Read user data one result at time.
 * @param globular
 * @param application
 * @param domain
 * @param query
 * @param callback
 * @param errorCallback
 */
function readOneUserData(globular, query, callback, errorCallback) {
    var userName = localStorage.getItem("user_name");
    var database = userName + "_db";
    var collection = "user_data";
    var rqst = new persistence_pb_1.FindOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setQuery(query);
    rqst.setOptions("");
    // call persist data
    globular.persistenceService
        .findOne(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        callback(JSON.parse(rsp.getJsonstr()));
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.readOneUserData = readOneUserData;
/**
 * Read all user data.
 * @param globular
 * @param application
 * @param domain
 * @param query
 * @param callback
 * @param errorCallback
 */
function readUserData(globular, query, callback, errorCallback) {
    var userName = localStorage.getItem("user_name");
    var database = userName + "_db";
    var collection = "user_data";
    var rqst = new persistence_pb_1.FindOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setQuery(query);
    rqst.setOptions("");
    // call persist data
    var stream = globular.persistenceService.find(rqst, {
        "token": token,
        "application": application, "domain": domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            callback(results);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.readUserData = readUserData;
///////////////////////////////////// Role action //////////////////////////////////////
/**
 * Retreive all available actions on the server.
 * @param callback That function is call in case of success.
 * @param errorCallback That function is call in case error.
 */
function getAllActions(globular, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetAllActionsRqst();
    globular.ressourceService
        .getAllActions(rqst, { "application": application, "domain": domain })
        .then(function (rsp) {
        callback(rsp.getActionsList());
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.getAllActions = getAllActions;
/**
 * Retreive the list of all available roles on the server.
 * @param callback That function is call in case of success.
 * @param errorCallback That function is call in case error.
 */
function getAllRoles(globular, callback, errorCallback) {
    var rqst = new persistence_pb_1.FindRqst();
    rqst.setCollection("Roles");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setQuery("{}"); // means all values.
    var stream = globular.persistenceService.find(rqst, {
        "application": application, "domain": domain
    });
    var roles = [];
    stream.on("data", function (rsp) {
        roles = roles.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            callback(roles);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.getAllRoles = getAllRoles;
/**
 * Append Action to a given role.
 * @param action The action name.
 * @param role The role.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
function AppendActionToRole(globular, role, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.AddRoleActionRqst();
    rqst.setRoleid(role);
    rqst.setAction(action);
    globular.ressourceService
        .addRoleAction(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.AppendActionToRole = AppendActionToRole;
/**
 * Remove the action from a given role.
 * @param action The action id
 * @param role The role id
 * @param callback success callback
 * @param errorCallback error callback
 */
function RemoveActionFromRole(globular, role, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveRoleActionRqst();
    rqst.setRoleid(role);
    rqst.setAction(action);
    globular.ressourceService
        .removeRoleAction(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.RemoveActionFromRole = RemoveActionFromRole;
/**
 * Create a new Role
 * @param globular
 * @param application
 * @param domain
 * @param id
 * @param callback
 * @param errorCallback
 */
function CreateRole(globular, id, callback, errorCallback) {
    var rqst = new ressource_pb_1.CreateRoleRqst();
    var role = new ressource_pb_1.Role();
    role.setId(id);
    role.setName(id);
    rqst.setRole(role);
    globular.ressourceService
        .createRole(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.CreateRole = CreateRole;
/**
 * Delete a given role
 * @param globular
 * @param application
 * @param domain
 * @param id
 * @param callback
 * @param errorCallback
 */
function DeleteRole(globular, id, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteRoleRqst();
    rqst.setRoleid(id);
    globular.ressourceService
        .deleteRole(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.DeleteRole = DeleteRole;
///////////////////////////////////// Application operations /////////////////////////////////
/**
 * Return the list of all application
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
function GetAllApplicationsInfo(globular, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetAllApplicationsInfoRqst();
    globular.ressourceService
        .getAllApplicationsInfo(rqst)
        .then(function (rsp) {
        var infos = JSON.parse(rsp.getResult());
        callback(infos);
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.GetAllApplicationsInfo = GetAllApplicationsInfo;
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
function AppendActionToApplication(globular, applicationId, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.AddApplicationActionRqst;
    rqst.setApplicationid(applicationId);
    rqst.setAction(action);
    globular.ressourceService.addApplicationAction(rqst, { "token": token, "application": application, "domain": domain })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.AppendActionToApplication = AppendActionToApplication;
/**
 * Remove action from application.
 * @param globular
 * @param application
 * @param domain
 * @param action
 * @param callback
 * @param errorCallback
 */
function RemoveActionFromApplication(globular, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveApplicationActionRqst;
    rqst.setApplicationid(application);
    rqst.setAction(action);
    globular.ressourceService.removeApplicationAction(rqst, { "token": token, "application": application, "domain": domain })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.RemoveActionFromApplication = RemoveActionFromApplication;
/**
 * Delete application
 * @param globular
 * @param application
 * @param domain
 * @param applicationId
 * @param callback
 * @param errorCallback
 */
function DeleteApplication(globular, applicationId, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteApplicationRqst;
    rqst.setApplicationid(applicationId);
    globular.ressourceService.deleteApplication(rqst, { "token": token, "application": application, "domain": domain })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.DeleteApplication = DeleteApplication;
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
function SaveApplication(globular, eventHub, _application, callback, errorCallback) {
    var rqst = new persistence_pb_1.ReplaceOneRqst;
    rqst.setCollection("Applications");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setValue(JSON.stringify(_application));
    rqst.setQuery("{\"_id\":\"" + _application._id + "\"}"); // means all values.
    globular.persistenceService.replaceOne(rqst, { "token": token, "application": application, "domain": domain })
        .then(function (rsp) {
        eventHub.publish("update_application_info_event", JSON.stringify(application), false);
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.SaveApplication = SaveApplication;
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
function GetServiceDescriptor(globular, serviceId, publisherId, callback, errorCallback) {
    var rqst = new services_pb_1.GetServiceDescriptorRequest;
    rqst.setServiceid(serviceId);
    rqst.setPublisherid(publisherId);
    globular.servicesDicovery.getServiceDescriptor(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(function (rsp) {
        callback(rsp.getResultsList());
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.GetServiceDescriptor = GetServiceDescriptor;
/**
 * Get the list of all service descriptor hosted on a server.
 * @param globular The globular object instance
 * @param application The application name who called the function.
 * @param domain The domain where the application reside.
 * @param callback
 * @param errorCallback
 */
function GetServicesDescriptor(globular, callback, errorCallback) {
    var rqst = new services_pb_1.GetServicesDescriptorRequest;
    var stream = globular.servicesDicovery.getServicesDescriptor(rqst, {
        "application": application, "domain": domain
    });
    var descriptors = new Array();
    stream.on("data", function (rsp) {
        descriptors = descriptors.concat(rsp.getResultsList());
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            callback(descriptors);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.GetServicesDescriptor = GetServicesDescriptor;
/**
 * Create or update a service descriptor.
 * @param globular
 * @param application
 * @param domain
 * @param descriptor
 * @param callback
 * @param errorCallback
 */
function SetServicesDescriptor(globular, descriptor, callback, errorCallback) {
    var rqst = new services_pb_1.SetServiceDescriptorRequest;
    rqst.setDescriptor(descriptor);
    globular.servicesDicovery.setServiceDescriptor(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.SetServicesDescriptor = SetServicesDescriptor;
/**
 * Find services by keywords.
 * @param query
 * @param callback
 */
function findServices(globular, keywords, callback, errorCallback) {
    var rqst = new services_pb_1.FindServicesDescriptorRequest();
    rqst.setKeywordsList(keywords);
    // Find services by keywords.
    globular.servicesDicovery
        .findServices(rqst, { "application": application, "domain": domain })
        .then(function (rsp) {
        callback(rsp.getResultsList());
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.findServices = findServices;
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
function installService(globular, discoveryId, serviceId, publisherId, version, callback, errorCallback) {
    var rqst = new admin_pb_1.InstallServiceRequest();
    rqst.setPublisherid(publisherId);
    rqst.setDicorveryid(discoveryId);
    rqst.setServiceid(serviceId);
    rqst.setVersion(version);
    // Install the service.
    globular.adminService
        .installService(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        readFullConfig(globular, callback, errorCallback);
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.installService = installService;
/**
 * Stop a service.
 */
function stopService(globular, serviceId, callback, errorCallback) {
    var rqst = new admin_pb_1.StopServiceRequest();
    rqst.setServiceId(serviceId);
    globular.adminService
        .stopService(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function () {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.stopService = stopService;
/**
 * Start a service
 * @param serviceId The id of the service to start.
 * @param callback  The callback on success.
 */
function startService(globular, serviceId, callback, errorCallback) {
    var rqst = new admin_pb_1.StartServiceRequest();
    rqst.setServiceId(serviceId);
    globular.adminService
        .startService(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function () {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.startService = startService;
/**
 * Here I will save the service configuration.
 * @param service The configuration to save.
 */
function saveService(globular, service, callback, errorCallback) {
    var rqst = new admin_pb_1.SaveConfigRequest();
    rqst.setConfig(JSON.stringify(service));
    globular.adminService
        .saveConfig(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        // The service with updated values...
        callback(JSON.parse(rsp.getResult()));
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.saveService = saveService;
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
function uninstallService(globular, service, deletePermissions, callback, errorCallback) {
    var rqst = new admin_pb_1.UninstallServiceRequest;
    rqst.setServiceid(service.Id);
    rqst.setPublisherid(service.PublisherId);
    rqst.setVersion(service.Version);
    rqst.setDeletepermissions(deletePermissions);
    globular.adminService
        .uninstallService(rqst, {
        "token": token,
        "application": application, "domain": domain
    })
        .then(function (rsp) {
        delete globular.config.Services[service.Id];
        // The service with updated values...
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.uninstallService = uninstallService;
/**
 * Return the list of service bundles.
 * @param callback
 */
function GetServiceBundles(globular, publisherId, serviceId, version, callback, errorCallback) {
    var rqst = new persistence_pb_1.FindRqst();
    rqst.setCollection("ServiceBundle");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setQuery("{}"); // means all values.
    var stream = globular.persistenceService.find(rqst, {
        "application": application, "domain": domain
    });
    var bundles = [];
    stream.on("data", function (rsp) {
        bundles = bundles.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            // filter localy.
            callback(bundles.filter(function (bundle) { return String(bundle._id).startsWith(publisherId + '%' + serviceId + '%' + version); }));
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.GetServiceBundles = GetServiceBundles;
/**
 * Get the object pointed by a reference.
 * @param globular
 * @param application
 * @param domain
 * @param ref
 * @param callback
 * @param errorCallback
 */
function getReferencedValue(globular, ref, callback, errorCallback) {
    var database = ref.$db;
    var collection = ref.$ref;
    var rqst = new persistence_pb_1.FindOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setQuery("{\"_id\":\"" + ref.$id + "\"}");
    rqst.setOptions("");
    globular.persistenceService.findOne(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(function (rsp) {
        callback(JSON.parse(rsp.getJsonstr()));
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.getReferencedValue = getReferencedValue;
///////////////////////////// Logging Operations ////////////////////////////////////////
/**
 * Read all errors data for server log.
 * @param globular
 * @param application
 * @param domain
 * @param callback
 * @param errorCallback
 */
function readErrors(globular, callback, errorCallback) {
    var database = "local_ressource";
    var collection = "Logs";
    var rqst = new persistence_pb_1.FindOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setOptions("");
    rqst.setQuery("{}");
    // call persist data
    var stream = globular.persistenceService.find(rqst, {
        "token": token,
        "application": application, "domain": domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            callback(results);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.readErrors = readErrors;
/**
 *  Read all logs
 * @param globular
 * @param application
 * @param domain
 * @param query
 * @param callback
 * @param errorCallback
 */
function readLogs(globular, query, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetLogRqst();
    rqst.setQuery(query);
    // call persist data
    var stream = globular.ressourceService.getLog(rqst, {
        "token": token,
        "application": application, "domain": domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(rsp.getInfoList());
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            results = results.sort(function (t1, t2) {
                var name1 = t1.getDate();
                var name2 = t2.getDate();
                if (name1 < name2) {
                    return 1;
                }
                if (name1 > name2) {
                    return -1;
                }
                return 0;
            });
            callback(results);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.readLogs = readLogs;
/**
 * Clear all log of a given type.
 * @param globular
 * @param application
 * @param domain
 * @param logType
 * @param callback
 * @param errorCallback
 */
function clearAllLog(globular, logType, callback, errorCallback) {
    var rqst = new ressource_pb_1.ClearAllLogRqst;
    rqst.setType(logType);
    globular.ressourceService.clearAllLog(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.clearAllLog = clearAllLog;
/**
 * Delete log entry.
 * @param globular
 * @param application
 * @param domain
 * @param log
 * @param callback
 * @param errorCallback
 */
function deleteLogEntry(globular, log, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteLogRqst;
    rqst.setLog(log);
    globular.ressourceService.deleteLog(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.deleteLogEntry = deleteLogEntry;
/**
 * Return the logged method and their count.
 * @param pipeline
 * @param callback
 * @param errorCallback
 */
function getNumbeOfLogsByMethod(globular, callback, errorCallback) {
    var database = "local_ressource";
    var collection = "Logs";
    var rqst = new persistence_pb_1.AggregateRqst;
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setOptions("");
    var pipeline = "[{\"$group\":{\"_id\":{\"method\":\"$method\"}, \"count\":{\"$sum\":1}}}]";
    rqst.setPipeline(pipeline);
    // call persist data
    var stream = globular.persistenceService.aggregate(rqst, {
        "token": token,
        "application": application, "domain": domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code === 0) {
            callback(results);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.getNumbeOfLogsByMethod = getNumbeOfLogsByMethod;
//////////////////////////// PLC Operations ///////////////////////////////////
var PLC_TYPE;
(function (PLC_TYPE) {
    PLC_TYPE[PLC_TYPE["ALEN_BRADLEY"] = 1] = "ALEN_BRADLEY";
    PLC_TYPE[PLC_TYPE["SIEMENS"] = 2] = "SIEMENS";
    PLC_TYPE[PLC_TYPE["MODBUS"] = 3] = "MODBUS";
})(PLC_TYPE = exports.PLC_TYPE || (exports.PLC_TYPE = {}));
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
function readPlcTag(globular, plcType, connectionId, name, type, offset) {
    return __awaiter(this, void 0, void 0, function () {
        var rqst, result, rsp, rsp, err_1;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    rqst = new plc_pb_1.ReadTagRqst();
                    rqst.setName(name);
                    rqst.setType(type);
                    rqst.setOffset(offset);
                    rqst.setConnectionId(connectionId);
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 11, , 12]);
                    if (!(plcType === PLC_TYPE.ALEN_BRADLEY)) return [3 /*break*/, 5];
                    if (!(globular.plcService_ab !== undefined)) return [3 /*break*/, 3];
                    return [4 /*yield*/, globular.plcService_ab.readTag(rqst, {
                            "token": token,
                            "application": application, "domain": domain
                        })];
                case 2:
                    rsp = _a.sent();
                    result = rsp.getValues();
                    return [3 /*break*/, 4];
                case 3: return [2 /*return*/, "No Alen Bradlay PLC server configured!"];
                case 4: return [3 /*break*/, 10];
                case 5:
                    if (!(plcType === PLC_TYPE.SIEMENS)) return [3 /*break*/, 9];
                    if (!(globular.plcService_siemens !== undefined)) return [3 /*break*/, 7];
                    return [4 /*yield*/, globular.plcService_siemens.readTag(rqst, {
                            "token": token,
                            "application": application, "domain": domain
                        })];
                case 6:
                    rsp = _a.sent();
                    result = rsp.getValues();
                    return [3 /*break*/, 8];
                case 7: return [2 /*return*/, "No Siemens PLC server configured!"];
                case 8: return [3 /*break*/, 10];
                case 9: return [2 /*return*/, "No PLC server configured!"];
                case 10: return [3 /*break*/, 12];
                case 11:
                    err_1 = _a.sent();
                    return [2 /*return*/, err_1];
                case 12:
                    // Here I got the value in a string I will convert it into it type.
                    if (type === plc_pb_1.TagType.BOOL) {
                        return [2 /*return*/, result === "true" ? true : false];
                    }
                    else if (type === plc_pb_1.TagType.REAL) {
                        return [2 /*return*/, parseFloat(result)];
                    }
                    else { // Must be cinsidere a integer.
                        return [2 /*return*/, parseInt(result, 10)];
                    }
                    return [2 /*return*/];
            }
        });
    });
}
exports.readPlcTag = readPlcTag;
///////////////////////////////////// LDAP operations /////////////////////////////////
/**
 * Synchronize LDAP and Globular/MongoDB user and roles.
 * @param info The synchronisations informations.
 * @param callback success callback.
 */
function syncLdapInfos(globular, info, timeout, callback, errorCallback) {
    var rqst = new ressource_pb_1.SynchronizeLdapRqst;
    var syncInfos = new ressource_pb_1.LdapSyncInfos;
    syncInfos.setConnectionid(info.connectionId);
    syncInfos.setLdapseriveid(info.ldapSeriveId);
    syncInfos.setRefresh(info.refresh);
    var userSyncInfos = new ressource_pb_1.UserSyncInfos;
    userSyncInfos.setBase(info.userSyncInfos.base);
    userSyncInfos.setQuery(info.userSyncInfos.query);
    userSyncInfos.setId(info.userSyncInfos.id);
    userSyncInfos.setEmail(info.userSyncInfos.email);
    syncInfos.setUsersyncinfos(userSyncInfos);
    var groupSyncInfos = new ressource_pb_1.GroupSyncInfos;
    groupSyncInfos.setBase(info.groupSyncInfos.base);
    groupSyncInfos.setQuery(info.groupSyncInfos.query);
    groupSyncInfos.setId(info.groupSyncInfos.id);
    syncInfos.setGroupsyncinfos(groupSyncInfos);
    rqst.setSyncinfo(syncInfos);
    // Try to synchronyze the ldap service.
    globular.ressourceService.synchronizeLdap(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(function (rsp) {
        callback();
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.syncLdapInfos = syncLdapInfos;
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
function pingSql(globular, connectionId, callback, errorCallback) {
    var rqst = new persistence_pb_1.PingConnectionRqst;
    rqst.setId(connectionId);
    globular.sqlService.ping(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(function (rsp) {
        callback(rsp.getResult());
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.pingSql = pingSql;
///////////////////////////// Search Operations ////////////////////////////////////////
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
function searchDocuments(globular, paths, query, language, fields, offset, pageSize, snippetLength, callback, errorCallback) {
    var rqst = new search_pb_1.SearchDocumentsRequest;
    rqst.setPathsList(paths);
    rqst.setQuery(query);
    rqst.setLanguage(language);
    rqst.setFieldsList(fields);
    rqst.setOffset(offset);
    rqst.setPagesize(pageSize);
    rqst.setSnippetlength(snippetLength);
    globular.searchService.searchDocuments(rqst, {
        "token": token,
        "application": application, "domain": domain
    }).then(function (rsp) {
        callback(rsp.getResultsList());
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.searchDocuments = searchDocuments;
