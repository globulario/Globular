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
var admin_pb_1 = require("./admin/admin_pb");
var monitoring_pb_1 = require("./monitoring/monitoringpb/monitoring_pb");
var ressource_pb_1 = require("./ressource/ressource_pb");
var jwt = require("jwt-decode");
var persistence_pb_1 = require("./persistence/persistencepb/persistence_pb");
var services_pb_1 = require("./services/services_pb");
var file_pb_1 = require("./file/filepb/file_pb");
var plc_pb_1 = require("./plc/plcpb/plc_pb");
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
function readFullConfig(globular, application, domain, callback, errorCallback) {
    var rqst = new admin_pb_1.GetConfigRequest();
    if (globular.adminService !== undefined) {
        globular.adminService
            .getFullConfig(rqst, {
            token: localStorage.getItem("user_token"),
            application: application, domain: domain
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
// Save the configuration.
function saveConfig(globular, application, domain, config, callback, errorCallback) {
    var rqst = new admin_pb_1.SaveConfigRequest();
    rqst.setConfig(JSON.stringify(config));
    if (globular.adminService !== undefined) {
        globular.adminService
            .saveConfig(rqst, {
            token: localStorage.getItem("user_token"),
            application: application, domain: domain
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
///////////////////////////////////// Ressources actions /////////////////////////////////////
/**
 * Retreive the list of ressource owner.
 * @param path
 * @param callback
 * @param errorCallback
 */
function getRessourceOwners(globular, application, domain, path, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetRessourceOwnersRqst;
    path = path.replace("/webroot", "");
    rqst.setPath(path);
    globular.ressourceService.getRessourceOwners(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
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
function setRessourceOwners(globular, application, domain, path, owner, callback, errorCallback) {
    var rqst = new ressource_pb_1.SetRessourceOwnerRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    rqst.setPath(path);
    rqst.setOwner(owner);
    globular.ressourceService.setRessourceOwner(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
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
function deleteRessourceOwners(globular, application, domain, path, owner, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteRessourceOwnerRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    rqst.setPath(path);
    rqst.setOwner(owner);
    globular.ressourceService.deleteRessourceOwner(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(function () {
        callback();
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.deleteRessourceOwners = deleteRessourceOwners;
/**
 * Retreive the permission for a given file.
 * @param path
 * @param callback
 * @param errorCallback
 */
function getRessourcePermissions(globular, application, domain, path, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetPermissionsRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);
    globular.ressourceService
        .getPermissions(rqst, {
        application: application, domain: domain
    })
        .then(function (rsp) {
        var permissions = JSON.parse(rsp.getPermissions());
        callback(permissions);
    })
        .catch(function (error) {
        if (errorCallback != undefined) {
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
function setRessourcePermission(globular, application, domain, path, owner, ownerType, number, callback, errorCallback) {
    var rqst = new ressource_pb_1.SetPermissionRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    var permission = new ressource_pb_1.RessourcePermission;
    permission.setPath(path);
    permission.setNumber(number);
    if (ownerType == OwnerType.User) {
        permission.setUser(owner);
    }
    else if (ownerType == OwnerType.Role) {
        permission.setRole(owner);
    }
    else if (ownerType == OwnerType.Application) {
        permission.setApplication(owner);
    }
    rqst.setPermission(permission);
    globular.ressourceService
        .setPermission(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback != undefined) {
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
function deleteRessourcePermissions(globular, application, domain, path, owner, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeletePermissionsRqst;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);
    rqst.setOwner(owner);
    globular.ressourceService
        .deletePermissions(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback != undefined) {
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
function getAllFilesInfo(globular, application, domain, callbak, errorCallback) {
    var rqst = new ressource_pb_1.GetAllFilesInfoRqst();
    globular.ressourceService
        .getAllFilesInfo(rqst, { application: application, domain: domain })
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
function renameFile(globular, application, domain, path, newName, oldName, callback, errorCallback) {
    var rqst = new file_pb_1.RenameRequest();
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);
    rqst.setOldName(oldName);
    rqst.setNewName(newName);
    globular.fileService
        .rename(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback != undefined) {
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
function deleteFile(globular, application, domain, path, callback, errorCallback) {
    var rqst = new file_pb_1.DeleteFileRequest();
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);
    globular.fileService
        .deleteFile(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback != undefined) {
            errorCallback(error);
        }
    });
}
exports.deleteFile = deleteFile;
/**
 *
 * @param path The path of the directory to be deleted.
 * @param callback The success callback
 * @param errorCallback The error callback.
 */
function deleteDir(globular, application, domain, path, callback, errorCallback) {
    var rqst = new file_pb_1.DeleteDirRequest();
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);
    globular.fileService
        .deleteDir(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback != undefined) {
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
function createArchive(globular, application, domain, path, name, callback, errorCallback) {
    var rqst = new file_pb_1.CreateArchiveRequest;
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    rqst.setPath(path);
    rqst.setName(name);
    globular.fileService.createAchive(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(function (rsp) {
        callback(rsp.getResult());
    }).catch(function (error) {
        if (errorCallback != undefined) {
            errorCallback(error);
        }
    });
}
exports.createArchive = createArchive;
/**
 *
 * @param urlToSend
 */
function downloadFileHttp(globular, application, domain, urlToSend, fileName, callback) {
    var req = new XMLHttpRequest();
    req.open("GET", urlToSend, true);
    // Set the token to manage downlaod access.
    req.setRequestHeader("token", localStorage.getItem("user_token"));
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
function downloadDir(globular, application, domain, path, callback, errorCallback) {
    var name = path.split("/")[path.split("/").length - 1];
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    // Create an archive-> download it-> delete it...
    createArchive(globular, application, domain, path, name, function (path) {
        var name = path.split("/")[path.split("/").length - 1];
        // display the archive path...
        downloadFileHttp(globular, application, domain, window.location.origin + path, name, function () {
            // Here the file was downloaded I will now delete it.
            setTimeout(function () {
                // wait a little and remove the archive from the server.
                var rqst = new file_pb_1.DeleteFileRequest;
                rqst.setPath(path);
                globular.fileService.deleteFile(rqst, {
                    token: localStorage.getItem("user_token"),
                    application: application, domain: domain
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
function readDir(globular, application, domain, path, callback, errorCallback) {
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    var rqst = new file_pb_1.ReadDirRequest;
    rqst.setPath(path);
    rqst.setRecursive(true);
    rqst.setThumnailheight(256);
    rqst.setThumnailwidth(256);
    var uint8array = new Uint8Array(0);
    var stream = globular.fileService.readDir(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    stream.on("data", function (rsp) {
        uint8array = mergeTypedArraysUnsafe(uint8array, rsp.getData());
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
            var jsonStr = new TextDecoder("utf-8").decode(uint8array);
            var content = JSON.parse(jsonStr);
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
 *
 * @param files
 */
function fileExist(fileName, files) {
    if (files != null) {
        for (var i = 0; i < files.length; i++) {
            if (files[i].Name == fileName) {
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
function createDir(globular, application, domain, path, callback, errorCallback) {
    path = path.replace("/webroot", ""); // remove the /webroot part.
    if (path.length == 0) {
        path = "/";
    }
    // first of all I will read the directory content...
    readDir(globular, application, domain, path, function (dir) {
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
            token: localStorage.getItem("user_token"),
            application: application, domain: domain
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
function queryTs(globular, application, domain, connectionId, query, ts, callback, errorCallback) {
    // Create a new request.
    var request = new monitoring_pb_1.QueryRequest();
    request.setConnectionid(connectionId);
    request.setQuery(query);
    request.setTs(ts);
    // Now I will test with promise
    globular.monitoringService
        .query(request, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (resp) {
        if (callback != undefined) {
            callback(JSON.parse(resp.getValue()));
        }
    })
        .catch(function (error) {
        if (errorCallback != undefined) {
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
function queryTsRange(globular, application, domain, connectionId, query, startTime, endTime, step, callback, errorCallback) {
    // Create a new request.
    var request = new monitoring_pb_1.QueryRangeRequest();
    request.setConnectionid(connectionId);
    request.setQuery(query);
    request.setStarttime(startTime);
    request.setEndtime(endTime);
    request.setStep(step);
    var buffer = { value: "", warning: "" };
    var stream = globular.monitoringService.queryRange(request, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    stream.on("data", function (rsp) {
        buffer.value += rsp.getValue();
        buffer.warning = rsp.getWarnings();
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
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
 * Register a new user.
 * @param userName The name of the account
 * @param email The email
 * @param password The password
 * @param confirmPassword
 * @param callback
 * @param errorCallback
 */
function registerAccount(globular, application, domain, userName, email, password, confirmPassword, callback, errorCallback) {
    var request = new ressource_pb_1.RegisterAccountRqst();
    var account = new ressource_pb_1.Account();
    account.setName(userName);
    account.setEmail(email);
    request.setAccount(account);
    request.setPassword(password);
    request.setConfirmPassword(confirmPassword);
    // Create the user account.
    globular.ressourceService
        .registerAccount(request, { application: application, domain: domain })
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
function DeleteAccount(globular, application, domain, id, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteAccountRqst;
    rqst.setId(id);
    // Remove the account from the database.
    globular.ressourceService
        .deleteAccount(rqst, { token: localStorage.getItem("user_token"), application: application, domain: domain })
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
function RemoveRoleFromAccount(globular, application, domain, accountId, roleId, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveAccountRoleRqst;
    rqst.setAccountid(accountId);
    rqst.setRoleid(roleId);
    globular.ressourceService
        .removeAccountRole(rqst, { token: localStorage.getItem("user_token"), application: application, domain: domain })
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
function AppendRoleToAccount(globular, application, domain, accountId, roleId, callback, errorCallback) {
    var rqst = new ressource_pb_1.AddAccountRoleRqst;
    rqst.setAccountid(accountId);
    rqst.setRoleid(roleId);
    globular.ressourceService
        .addAccountRole(rqst, { token: localStorage.getItem("user_token"), application: application, domain: domain })
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
function updateAccountEmail(globular, application, domain, accountId, old_email, new_email, callback, errorCallback) {
    var rqst = new admin_pb_1.SetEmailRequest;
    rqst.setAccountid(accountId);
    rqst.setOldemail(old_email);
    rqst.setNewemail(new_email);
    globular.adminService.setEmail(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        console.log("fail to save config ", err);
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
function updateAccountPassword(globular, application, domain, accountId, old_password, new_password, confirm_password, callback, errorCallback) {
    var rqst = new admin_pb_1.SetPasswordRequest;
    rqst.setAccountid(accountId);
    rqst.setOldpassword(old_password);
    rqst.setNewpassword(new_password);
    if (confirm_password != new_password) {
        errorCallback("password not match!");
        return;
    }
    globular.adminService.setPassword(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(function (rsp) {
        callback();
    })
        .catch(function (error) {
        if (errorCallback != undefined) {
            errorCallback(error);
        }
    });
}
exports.updateAccountPassword = updateAccountPassword;
/**
 * Authenticate the user and get the token
 * @param userName The account name or email
 * @param password  The user password
 * @param callback
 * @param errorCallback
 */
function authenticate(globular, eventHub, application, domain, userName, password, callback, errorCallback) {
    var rqst = new ressource_pb_1.AuthenticateRqst();
    rqst.setName(userName);
    rqst.setPassword(password);
    // Create the user account.
    globular.ressourceService
        .authenticate(rqst, { application: application, domain: domain })
        .then(function (rsp) {
        // Here I will set the token in the localstorage.
        var token = rsp.getToken();
        var decoded = jwt(token);
        // here I will save the user token and user_name in the local storage.
        localStorage.setItem("user_token", token);
        localStorage.setItem("user_name", decoded.username);
        readFullConfig(globular, application, domain, function (config) {
            // Publish local login event.
            eventHub.publish("onlogin", config, true); // return the full config...
            callback(decoded);
        }, function (err) {
            errorCallback(err);
        });
    })
        .catch(function (err) {
        console.log(err);
        errorCallback(err);
    });
}
exports.authenticate = authenticate;
/**
 * Function to be use to refresh token or full configuration.
 * @param callback On success callback
 * @param errorCallback On error callback
 */
function refreshToken(globular, eventHub, application, domain, callback, errorCallback) {
    var rqst = new ressource_pb_1.RefreshTokenRqst();
    rqst.setToken(localStorage.getItem("user_token"));
    globular.ressourceService
        .refreshToken(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        // Here I will set the token in the localstorage.
        var token = rsp.getToken();
        var decoded = jwt(token);
        // here I will save the user token and user_name in the local storage.
        localStorage.setItem("user_token", token);
        localStorage.setItem("user_name", decoded.username);
        readFullConfig(globular, application, domain, function (config) {
            // Publish local login event.
            eventHub.publish("onlogin", config, true); // return the full config...
            callback(decoded);
        }, function (err) {
            errorCallback(err);
        });
    })
        .catch(function (err) {
        onerror(err);
    });
}
exports.refreshToken = refreshToken;
/**
 * Save user data into the user_data collection.
 */
function appendUserData(globular, application, domain, data, callback) {
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
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        callback(rsp.getId());
    })
        .catch(function (err) {
        console.log(err);
    });
}
exports.appendUserData = appendUserData;
/**
 * Read user data one result at time.
 */
function readOneUserData(globular, application, domain, query, callback) {
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
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        callback(JSON.parse(rsp.getJsonstr()));
    })
        .catch(function (err) {
        console.log(err);
    });
}
exports.readOneUserData = readOneUserData;
/**
 * Read all user data.
 */
function readUserData(globular, application, domain, query, callback, errorCallback) {
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
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
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
 * Return the list of all account on the server, guest and admin are new account...
 * @param callback
 */
function GetAllAccountsInfo(globular, application, domain, callback, errorCallback) {
    var rqst = new persistence_pb_1.FindRqst();
    rqst.setCollection("Accounts");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setQuery("{}"); // means all values.
    var stream = globular.persistenceService.find(rqst, {
        application: application, domain: domain
    });
    var accounts = new Array();
    stream.on("data", function (rsp) {
        accounts = accounts.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(accounts);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.GetAllAccountsInfo = GetAllAccountsInfo;
/**
 * Retreive all available actions on the server.
 * @param callback That function is call in case of success.
 * @param errorCallback That function is call in case error.
 */
function getAllActions(globular, application, domain, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetAllActionsRqst();
    globular.ressourceService
        .getAllActions(rqst, { application: application, domain: domain })
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
function getAllRoles(globular, application, domain, callback, errorCallback) {
    var rqst = new persistence_pb_1.FindRqst();
    rqst.setCollection("Roles");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setQuery("{}"); // means all values.
    var stream = globular.persistenceService.find(rqst, {
        application: application, domain: domain
    });
    var roles = new Array();
    stream.on("data", function (rsp) {
        roles = roles.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
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
function AppendActionToRole(globular, application, domain, role, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.AddRoleActionRqst();
    rqst.setRoleid(role);
    rqst.setAction(action);
    globular.ressourceService
        .addRoleAction(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
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
function RemoveActionFromRole(globular, application, domain, role, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveRoleActionRqst();
    rqst.setRoleid(role);
    rqst.setAction(action);
    globular.ressourceService
        .removeRoleAction(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.RemoveActionFromRole = RemoveActionFromRole;
function CreateRole(globular, application, domain, id, callback, errorCallback) {
    var rqst = new ressource_pb_1.CreateRoleRqst();
    var role = new ressource_pb_1.Role();
    role.setId(id);
    role.setName(id);
    rqst.setRole(role);
    globular.ressourceService
        .createRole(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.CreateRole = CreateRole;
function DeleteRole(globular, application, domain, id, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteRoleRqst();
    rqst.setRoleid(id);
    globular.ressourceService
        .deleteRole(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
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
function GetAllApplicationsInfo(globular, application, domain, callback, errorCallback) {
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
function AppendActionToApplication(globular, application, domain, applicationId, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.AddApplicationActionRqst;
    rqst.setApplicationid(applicationId);
    rqst.setAction(action);
    globular.ressourceService.addApplicationAction(rqst, { token: localStorage.getItem("user_token"), application: application, domain: domain })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        console.log(err);
        errorCallback(err);
    });
}
exports.AppendActionToApplication = AppendActionToApplication;
function RemoveActionFromApplication(globular, application, domain, applicationId, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveApplicationActionRqst;
    rqst.setApplicationid(applicationId);
    rqst.setAction(action);
    globular.ressourceService.removeApplicationAction(rqst, { token: localStorage.getItem("user_token"), application: application, domain: domain })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        console.log(err);
        errorCallback(err);
    });
}
exports.RemoveActionFromApplication = RemoveActionFromApplication;
function DeleteApplication(globular, application, domain, applicationId, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteApplicationRqst;
    rqst.setApplicationid(applicationId);
    globular.ressourceService.deleteApplication(rqst, { token: localStorage.getItem("user_token"), application: application, domain: domain })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        console.log(err);
        errorCallback(err);
    });
}
exports.DeleteApplication = DeleteApplication;
function SaveApplication(globular, eventHub, application_, domain, application, callback, errorCallback) {
    var rqst = new persistence_pb_1.ReplaceOneRqst;
    rqst.setCollection("Applications");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setValue(JSON.stringify(application));
    rqst.setQuery("{\"_id\":\"" + application._id + "\"}"); // means all values.
    globular.persistenceService.replaceOne(rqst, { token: localStorage.getItem("user_token"), application: application_, domain: domain })
        .then(function (rsp) {
        eventHub.publish("update_application_info_event", JSON.stringify(application), false);
        callback();
    })
        .catch(function (err) {
        console.log(err);
        errorCallback(err);
    });
}
exports.SaveApplication = SaveApplication;
///////////////////////////////////// Peers operations /////////////////////////////////
function GetAllPeersInfo(globular, application, domain, query, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetPeersRqst();
    rqst.setQuery(query);
    var peers = new Array();
    var stream = globular.ressourceService.getPeers(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        peers = peers.concat(rsp.getPeersList());
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(peers);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.GetAllPeersInfo = GetAllPeersInfo;
function AppendActionToPeer(globular, application, domain, id, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.AddPeerActionRqst;
    rqst.setPeerid(id);
    rqst.setAction(action);
    globular.ressourceService.addPeerAction(rqst, { token: localStorage.getItem("user_token"), application: application, domain: domain })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        console.log(err);
        errorCallback(err);
    });
}
exports.AppendActionToPeer = AppendActionToPeer;
function RemoveActionFromPeer(globular, application, domain, id, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemovePeerActionRqst;
    rqst.setPeerid(id);
    rqst.setAction(action);
    globular.ressourceService.removePeerAction(rqst, { token: localStorage.getItem("user_token"), application: application, domain: domain })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        console.log(err);
        errorCallback(err);
    });
}
exports.RemoveActionFromPeer = RemoveActionFromPeer;
function DeletePeer(globular, application, domain, peer, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeletePeerRqst;
    rqst.setPeer(peer);
    globular.ressourceService.deletePeer(rqst, { token: localStorage.getItem("user_token"), application: application, domain: domain })
        .then(function (rsp) {
        callback();
    })
        .catch(function (err) {
        console.log(err);
        errorCallback(err);
    });
}
exports.DeletePeer = DeletePeer;
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
function GetServiceDescriptor(globular, application, domain, serviceId, publisherId, callback, errorCallback) {
    var rqst = new services_pb_1.GetServiceDescriptorRequest;
    rqst.setServiceid(serviceId);
    rqst.setPublisherid(publisherId);
    globular.servicesDicovery.getServiceDescriptor(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
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
function GetServicesDescriptor(globular, application, domain, callback, errorCallback) {
    var rqst = new services_pb_1.GetServicesDescriptorRequest;
    var stream = globular.servicesDicovery.getServicesDescriptor(rqst, {
        application: application, domain: domain
    });
    var descriptors = new Array();
    stream.on("data", function (rsp) {
        descriptors = descriptors.concat(rsp.getResultsList());
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(descriptors);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.GetServicesDescriptor = GetServicesDescriptor;
function SetServicesDescriptor(globular, application, domain, descriptor, callback, errorCallback) {
    var rqst = new services_pb_1.SetServiceDescriptorRequest;
    rqst.setDescriptor(descriptor);
    globular.servicesDicovery.setServiceDescriptor(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
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
function findServices(globular, application, domain, keywords, callback, errorCallback) {
    var rqst = new services_pb_1.FindServicesDescriptorRequest();
    rqst.setKeywordsList(keywords);
    // Find services by keywords.
    globular.servicesDicovery
        .findServices(rqst, { application: application, domain: domain })
        .then(function (rsp) {
        var results = rsp.getResultsList();
        callback(results);
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.findServices = findServices;
function installService(globular, application, domain, discoveryId, serviceId, publisherId, version, callback, errorCallback) {
    var rqst = new admin_pb_1.InstallServiceRequest();
    rqst.setPublisherid(publisherId);
    rqst.setDicorveryid(discoveryId);
    rqst.setServiceid(serviceId);
    rqst.setVersion(version);
    // Install the service.
    globular.adminService
        .installService(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        readFullConfig(globular, application, domain, callback, errorCallback);
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.installService = installService;
/**
 * Stop a service.
 */
function stopService(globular, application, domain, serviceId, callback, errorCallback) {
    var rqst = new admin_pb_1.StopServiceRequest();
    rqst.setServiceId(serviceId);
    globular.adminService
        .stopService(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
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
function startService(globular, application, domain, serviceId, callback, errorCallback) {
    var rqst = new admin_pb_1.StartServiceRequest();
    rqst.setServiceId(serviceId);
    globular.adminService
        .startService(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
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
function saveService(globular, application, domain, service, callback, errorCallback) {
    var rqst = new admin_pb_1.SaveConfigRequest();
    rqst.setConfig(JSON.stringify(service));
    globular.adminService
        .saveConfig(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    })
        .then(function (rsp) {
        // The service with updated values...
        var service = JSON.parse(rsp.getResult());
        callback(service);
    })
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.saveService = saveService;
function uninstallService(globular, application, domain, service, deletePermissions, callback, errorCallback) {
    var rqst = new admin_pb_1.UninstallServiceRequest;
    rqst.setServiceid(service.Id);
    rqst.setPublisherid(service.PublisherId);
    rqst.setVersion(service.Version);
    rqst.setDeletepermissions(deletePermissions);
    globular.adminService
        .uninstallService(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
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
function GetServiceBundles(globular, application, domain, publisherId, serviceId, version, callback, errorCallback) {
    var rqst = new persistence_pb_1.FindRqst();
    rqst.setCollection("ServiceBundle");
    rqst.setDatabase("local_ressource");
    rqst.setId("local_ressource");
    rqst.setQuery("{}"); // means all values.
    var stream = globular.persistenceService.find(rqst, {
        application: application, domain: domain
    });
    var bundles = new Array();
    stream.on("data", function (rsp) {
        bundles = bundles.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
            // filter localy.
            callback(bundles.filter(function (bundle) { return String(bundle._id).startsWith(publisherId + '%' + serviceId + '%' + version); }));
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.GetServiceBundles = GetServiceBundles;
// Get the object pointed by a reference.
function getReferencedValue(globular, application, domain, ref, callback, errorCallback) {
    var database = ref.$db;
    var collection = ref.$ref;
    var rqst = new persistence_pb_1.FindOneRqst();
    rqst.setId(database);
    rqst.setDatabase(database);
    rqst.setCollection(collection);
    rqst.setQuery("{\"_id\":\"" + ref.$id + "\"}");
    rqst.setOptions("");
    globular.persistenceService.findOne(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(function (rsp) {
        callback(JSON.parse(rsp.getJsonstr()));
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.getReferencedValue = getReferencedValue;
/**
 * Read all errors data.
 * @param callback
 */
function readErrors(globular, application, domain, callback, errorCallback) {
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
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(results);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.readErrors = readErrors;
///////////////////////////////////// Ressource & Permissions operations /////////////////////////////////
function readAllActionPermission(globular, application, domain, callback, errorCallback) {
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
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(results);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.readAllActionPermission = readAllActionPermission;
function getRessources(globular, application, domain, path, name, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetRessourcesRqst;
    rqst.setPath(path);
    rqst.setName(name);
    // call persist data
    var stream = globular.ressourceService.getRessources(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(rsp.getRessourcesList());
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
            callback(results);
        }
        else {
            errorCallback({ "message": status.details });
        }
    });
}
exports.getRessources = getRessources;
function setActionPermission(globular, application, domain, action, permission, callback, errorCallback) {
    var rqst = new ressource_pb_1.SetActionPermissionRqst;
    rqst.setAction(action);
    rqst.setPermission(permission);
    // Call set action permission.
    globular.ressourceService.setActionPermission(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.setActionPermission = setActionPermission;
function removeActionPermission(globular, application, domain, action, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveActionPermissionRqst;
    rqst.setAction(action);
    // Call set action permission.
    globular.ressourceService.removeActionPermission(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.removeActionPermission = removeActionPermission;
function removeRessource(globular, application, domain, path, name, callback, errorCallback) {
    var rqst = new ressource_pb_1.RemoveRessourceRqst;
    var ressource = new ressource_pb_1.Ressource;
    ressource.setPath(path);
    ressource.setName(name);
    rqst.setRessource(ressource);
    globular.ressourceService.removeRessource(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.removeRessource = removeRessource;
///////////////////////////// Logging Operations ////////////////////////////////////////
/**
 * Read all logs
 * @param callback The success callback.
 */
function readLogs(globular, application, domain, query, callback, errorCallback) {
    var rqst = new ressource_pb_1.GetLogRqst();
    rqst.setQuery(query);
    // call persist data
    var stream = globular.ressourceService.getLog(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(rsp.getInfoList());
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
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
            console.log(status.details);
            errorCallback({ "message": status.details });
        }
    });
}
exports.readLogs = readLogs;
function clearAllLog(globular, application, domain, logType, callback, errorCallback) {
    var rqst = new ressource_pb_1.ClearAllLogRqst;
    rqst.setType(logType);
    globular.ressourceService.clearAllLog(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.clearAllLog = clearAllLog;
function deleteLog(globular, application, domain, log, callback, errorCallback) {
    var rqst = new ressource_pb_1.DeleteLogRqst;
    rqst.setLog(log);
    globular.ressourceService.deleteLog(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(callback)
        .catch(function (err) {
        errorCallback(err);
    });
}
exports.deleteLog = deleteLog;
/**
 * Return the logged method and their count.
 * @param pipeline
 * @param callback
 * @param errorCallback
 */
function getNumbeOfLogsByMethod(globular, application, domain, callback, errorCallback) {
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
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    });
    var results = new Array();
    // Get the stream and set event on it...
    stream.on("data", function (rsp) {
        results = results.concat(JSON.parse(rsp.getJsonstr()));
    });
    stream.on("status", function (status) {
        if (status.code == 0) {
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
* Read a plc tag from the defined backend.
* @param plcType  The plc type can be Alen Bradley or Simens, modbus is on the planned.
* @param connectionId  The connection id defined for that plc.
* @param name The name of the tag to read.
* @param type The type name of the plc.
* @param offset The offset in the memory.
*/
function readPlcTag(globular, application, domain, plcType, connectionId, name, type, offset) {
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
                    if (!(plcType == PLC_TYPE.ALEN_BRADLEY)) return [3 /*break*/, 5];
                    if (!(globular.plcService_ab != undefined)) return [3 /*break*/, 3];
                    return [4 /*yield*/, globular.plcService_ab.readTag(rqst, {
                            token: localStorage.getItem("user_token"),
                            application: application, domain: domain
                        })];
                case 2:
                    rsp = _a.sent();
                    result = rsp.getValues();
                    return [3 /*break*/, 4];
                case 3: return [2 /*return*/, "No Alen Bradlay PLC server configured!"];
                case 4: return [3 /*break*/, 10];
                case 5:
                    if (!(plcType == PLC_TYPE.SIEMENS)) return [3 /*break*/, 9];
                    if (!(globular.plcService_siemens != undefined)) return [3 /*break*/, 7];
                    return [4 /*yield*/, globular.plcService_siemens.readTag(rqst, {
                            token: localStorage.getItem("user_token"),
                            application: application, domain: domain
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
                    if (type == plc_pb_1.TagType.BOOL) {
                        return [2 /*return*/, result == "true" ? true : false];
                    }
                    else if (type == plc_pb_1.TagType.REAL) {
                        return [2 /*return*/, parseFloat(result)];
                    }
                    else { // Must be cinsidere a integer.
                        return [2 /*return*/, parseInt(result)];
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
function syncLdapInfos(globular, application, domain, info, timeout, callback, errorCallback) {
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
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(function (rsp) {
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
function pingSql(globular, application, domain, connectionId, callback, errorCallback) {
    var rqst = new persistence_pb_1.PingConnectionRqst;
    rqst.setId(connectionId);
    globular.sqlService.ping(rqst, {
        token: localStorage.getItem("user_token"),
        application: application, domain: domain
    }).then(function (rsp) {
        callback(rsp.getResult());
    }).catch(function (err) {
        errorCallback(err);
    });
}
exports.pingSql = pingSql;
