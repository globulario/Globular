import * as GlobularWebClient from "globular-web-client";
import { GetConfigRequest } from "globular-web-client/lib/admin/admin_pb";
import { AuthenticateRqst, Account, RegisterAccountRqst, RegisterAccountRsp, RefreshTokenRqst, RefreshTokenRsp } from "globular-web-client/lib/ressource/ressource_pb";
import * as jwt from 'jsonwebtoken'

// Set the basic configuration without services details.
let config = {
    Name: "Globular",
    PortHttp: 80,
    PortHttps: 443,
    AdminPort: 10001,
    AdminProxy: 10002,
    RessourcePort: 10003,
    RessourceProxy: 10004,
    ServicesDiscoveryPort: 10005,
    ServicesDiscoveryProxy: 10006,
    ServicesRepositoryPort: 10007,
    ServicesRepositoryProxy: 10008,
    Protocol: "https",
    Domain: "www.omniscient.app",
    Services: {}, // empty for start.
    SessionTimeout: 5
};


// Create a new connection with the backend.
let globular = new GlobularWebClient.Globular(config);


function initServices(callback:()=>void) {
    // Here I will set the token...
    let rqst = new GetConfigRequest;
    if (globular.adminService !== undefined) {
        globular.adminService.getConfig(rqst).then((rsp:any) => {
            let config = JSON.parse(rsp.getResult())
            // init the services from the configuration retreived.
            globular = new GlobularWebClient.Globular(config);
            console.log(globular)
            callback()
        }).catch((err:any) => {
            console.log("fail to get config ", err)
        })
    }
}

initServices(()=>{
    test()
})

// let config = globular.adminService.GetConfig()
function readFullConfig(callback:(config:GlobularWebClient.IConfig)=>void) {
    let rqst = new GetConfigRequest;
    if (globular.adminService !== undefined) {
        globular.adminService.getFullConfig(rqst, {"token":localStorage.getItem("user_token")}).then((rsp) => {
           config = JSON.parse(rsp.getResult())
           // Keep the config session timeout in the local storage.
           localStorage.setItem("session_timeout", config.SessionTimeout.toString())
           callback(config)
        }).catch((err) => {
            console.log("fail to get config ", err)
        })
    }
}

function test(){
    Authentitcate("sa","adminadmin", (decoded:any)=>{
        console.log("----> success!", decoded)
        // here I will renew the token each time

        // I will try to read the full configuration.
        readFullConfig((config:GlobularWebClient.IConfig)=>{
            console.log("---> read full config success", config)
        })
    }, (err:any)=>{
        console.log("----> error!", err)
    })
}

function RegisterAccount(name:string, email:string, password:string, confirmPassword: string, callback: (token:string)=>void){
    let account = new Account
    account.setEmail(email)
    account.setName(name)
    let rqst = new RegisterAccountRqst
    rqst.setPassword(password)
    rqst.setConfirmPassword(confirmPassword)
    rqst.setAccount(account)
    globular.ressourceService.registerAccount(rqst).then((rsp: RegisterAccountRsp)=>{
        let token = rsp.getResult()
        if(token.length > 0){
            let decoded = jwt.decode(token);
            localStorage.setItem("user_token", token)
            localStorage.setItem("user_name", (<any>decoded).username)
        }
        callback(token)
    }).catch((err: any)=>{

    })
}

function Authentitcate(userName: string, password: string, successCallback:(userInfo: any)=>void, errorCallback:(err:string)=>void){
    let rqst = new AuthenticateRqst()
    rqst.setName(userName)
    rqst.setPassword(password)
    // Authenticate as admin...
    globular.ressourceService.authenticate(rqst).then((rsp)=>{
        let token = rsp.getToken()
        if(token.length > 0){
            // Here I will set the token in the localstorage.
            let decoded = jwt.decode(token);

            // here I will save the user token and user_name in the local storage.
            localStorage.setItem("user_token", token)
            localStorage.setItem("user_name", (<any>decoded).username)

            // Here I will refresh the token... 
            setInterval(()=>{
                let rqst = new RefreshTokenRqst
                // Get the last token from the local storage.
                let token = localStorage.getItem("user_token")
                rqst.setToken(token)

                globular.ressourceService.refreshToken(rqst).then((rsp:RefreshTokenRsp)=>{
                    let token = rsp.getToken()
                    let decoded = jwt.decode(token);
                    console.log("---> new token from refresh...")
                    // here I will save the user token and user_name in the local storage.
                    localStorage.setItem("user_token", token)
                    localStorage.setItem("user_name", (<any>decoded).username)
        
                })

            }, (15*60*60*1000) - 100)

            successCallback(decoded) 
        }

    }).catch((err)=>{
        errorCallback(err);
    })
}

