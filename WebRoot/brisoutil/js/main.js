/**
 * That function create and return the sql connection if it not already exist.
 */
function initConnections(initCallback) {
    // Bris d'outil
    function initBrisOutilConnection(callback, initCallback) {
        console.log("Init BrisOutil sql connection")
        var rqst = new Sql.CreateConnectionRqst();
        var c = new Sql.Connection();
        c.setId("bris_outil")
        c.setName("BrisOutil")
        c.setUser("dbprog")
        c.setPassword("dbprog")
        c.setPort(1433)
        c.setDriver("odbc")
        c.setHost("mon-sql-v01")
        c.setCharset("utf8")
        rqst.setConnection(c)

        globular.sqlService.createConnection(rqst, {}, function () {
            return function (err, rsp) {
                if (err != null) {
                    console.log(err)
                } else {
                    console.log("Done Init BrisOutil sql connection")
                    callback(initCallback) // call the callback function.
                }
            }
        }(callback, initCallback));
    }
    // ProprieteMecanique
    function initProprieteMecaniqueConnection(callback, initCallback) {
        console.log("Init ProprieteMecanique sql connection")
        var rqst = new Sql.CreateConnectionRqst();
        var c = new Sql.Connection();
        c.setId("propriete_mecanique")
        c.setName("ProprieteMecanique")
        c.setUser("QISDBUserRO")
        c.setPassword("Grr78CHNN99a")
        c.setPort(1433)
        c.setDriver("odbc")
        c.setHost("mon-sql-v01")
        c.setCharset("utf8")
        rqst.setConnection(c)

        globular.sqlService.createConnection(rqst, {}, function () {
            return function (err, rsp) {
                if (err != null) {
                    console.log(err)
                } else {
                    console.log("Done Init ProprieteMecanique sql connection")
                    callback(initCallback) // call the callback function.
                }
            }
        }(callback, initCallback));
    }
    // Shopvue
    function initShopvueConnection(callback, initCallback) {
        console.log("Init Shopvue sql connection")
        var rqst = new Sql.CreateConnectionRqst();
        var c = new Sql.Connection();
        c.setId("shopvue")
        c.setName("ShopVue")
        c.setUser("ShopvueReport")
        c.setPassword("Dowty123")
        c.setPort(1433)
        c.setDriver("odbc")
        c.setHost("mon-sql-v02\\ShopVue")
        c.setCharset("utf8")
        rqst.setConnection(c)

        globular.sqlService.createConnection(rqst, {}, function () {
            return function (err, rsp) {
                if (err != null) {
                    console.log(err)
                } else {
                    console.log("Done Init ShopVue sql connection")
                    callback(initCallback) // call the callback function.
                }
            }
        }(callback, initCallback));
    }
    // GestionInformation
    function initGestionInformationConnection(callback, initCallback) {
        console.log("init GestionInformation sql connection")
        var rqst = new Sql.CreateConnectionRqst();
        var c = new Sql.Connection();
        c.setId("gestion_information")
        c.setName("GestionInformation")
        c.setUser("dbprog")
        c.setPassword("dbprog")
        c.setPort(1433)
        c.setDriver("odbc")
        c.setHost("mon-sql-v01")
        c.setCharset("utf8")
        rqst.setConnection(c)

        globular.sqlService.createConnection(rqst, {}, function () {
            return function (err, rsp) {
                if (err != null) {
                    console.log(err)
                } else {
                    console.log("Done Init GestionInformation sql connection")
                    callback(initCallback) // call the callback function.
                }
            }
        }(callback, initCallback));
    }
    // INVInterne
    function initINVInterneConnection(callback, initCallback) {
        console.log("init INVInterne sql connection")
        var rqst = new Sql.CreateConnectionRqst();
        var c = new Sql.Connection();
        c.setId("catalog_md")
        c.setName("INVInterne")
        c.setUser("dbprog")
        c.setPassword("dbprog")
        c.setPort(1433)
        c.setDriver("odbc")
        c.setHost("mon-sql-v01")
        c.setCharset("utf8")
        rqst.setConnection(c)

        globular.sqlService.createConnection(rqst, {}, function () {
            return function (err, rsp) {
                if (err != null) {
                    console.log(err)
                } else {
                    console.log("Done Init INVInterne sql connection")
                    callback(initCallback) // call the callback function.
                }
            }
        }(callback, initCallback));
    }
    // Eng_Dwg_Features
    function initEng_Dwg_FeaturesConnection(callback, initCallback) {
        console.log("init Eng_Dwg_Features sql connection")
        var rqst = new Sql.CreateConnectionRqst();
        var c = new Sql.Connection();
        c.setId("Eng_Dwg_Features")
        c.setName("Eng_Dwg_Features")
        c.setUser("EngDWG_User")
        c.setPassword("eNGdwg_rw%5")
        c.setPort(1433)
        c.setDriver("odbc")
        c.setHost("mon-sql-v01")
        c.setCharset("utf8")
        rqst.setConnection(c)

        globular.sqlService.createConnection(rqst, {}, function () {
            return function (err, rsp) {
                if (err != null) {
                    console.log(err)
                } else {
                    console.log("Done Init INVInterne sql connection")
                    callback(initCallback) // call the callback function.
                }
            }
        }(callback, initCallback));
    }

    // Init the ldap connection.
    function initLdapConnection(callback, initCallback) {
        console.log("init ldap connection")
        var rqst = new Ldap.CreateConnectionRqst();
        var c = new Ldap.Connection();
        c.setId("safran_ldap")
        c.setUser("mrmfct037@UD6.UF6")
        c.setPassword("Dowty123")
        c.setPort(389)
        c.setHost("mon-dc-p01.UD6.UF6")
        rqst.setConnection(c)

        globular.ldapService.createConnection(rqst, {}, function () {
            return function (err, rsp) {
                if (err != null) {
                    console.log(err)
                } else {
                    console.log("Done create ldap connection")
                    callback(initCallback) // call the callback function.
                }
            }
        }(callback, initCallback));
    }

    // Call initialisation one after one...
    initBrisOutilConnection(function (initCallback) {
        initProprieteMecaniqueConnection(function (initCallback) {
            initINVInterneConnection(function (initCallback) {
                initGestionInformationConnection(function (initCallback) {
                    initShopvueConnection(function (initCallback) {
                        initEng_Dwg_FeaturesConnection(function (initCallback) {
                            initLdapConnection(function (initCallback) {
                                initCallback()
                            }, initCallback)
                        }, initCallback)
                    }, initCallback)
                }, initCallback)
            }, initCallback)
        }, initCallback)
    }, initCallback)
}

// Workorder serial number part number...
var productNumbers = []
var serialNumbers = []
var serialProduct = {}

// init all product part and serial at once.
function initSerialNumbers(initCallback) {
    var query = "SELECT [WorkOrder] "
    query += ",[SerialNumber]"
    query += ",[PartNumber]"
    query += ",[AssyNumber] "
    query += "FROM [Eng_Dwg_Features].[dbo].[PartInfo] "
    query += "WHERE PartNumber IS NOT NULL"

    var q = new Sql.Query()
    q.setQuery(query)
    q.setConnectionid("Eng_Dwg_Features")
    q.setParameters(JSON.stringify([]))

    var rqst = new Sql.QueryContextRqst()
    rqst.setQuery(q)

    console.log("Init serial number list")
    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.sqlService.queryContext(rqst, metadata);

    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if (response.hasHeader()) {
            var header = response.getHeader()
        } else if (response.hasRows()) {
            results = JSON.parse(response.getRows())
            for (var i = 0; i < results.length; i++) {
                var val = results[i]
                if (val[0].endsWith("M") || val[0].endsWith("C")) {
                    var workorder = val[0].replace("6M0", "")
                    var serial = val[1]
                    serialNumbers.push(workorder)
                    serialNumbers.push(serial)

                    // The product will not be the assembly number.
                    if(workorder.endsWith("C")){
                        serialProduct[workorder] = val[3]
                        serialProduct[serial] = val[3]
                    }else if(workorder.endsWith("M")){
                        serialProduct[workorder] = val[2]
                        serialProduct[serial] = val[2]
                    }else{
                        serialProduct[workorder] = val[2]
                        serialProduct[serial] = val[2]
                    }
                    
                    // Also append it to the product number.
                    if(productNumbers.indexOf(val[2]) == -1){
                        productNumbers.push(val[2])
                    }

                    // append the 
                    if(productNumbers.indexOf(val[3]) == -1){
                        productNumbers.push(val[3])
                    }
                }
            }
        }
    });

    // Get the results here if the statut is ok.
    stream.on('status', function (status) {
        if (status.code == 0) {
            console.log("Done init serial numbers")
            initCallback()
        }
    });

    stream.on('end', function () {
        return function (end) {
            // stream end signal
        }
    }(initCallback));
}

// The map of employes.
var employes = {}
var employeByName = {}
var employeById = {}
var employeNames = []

function initEmployes(initCallback) {
    var query = "SELECT [Employee_CrewID],[LastName],[FirstName],[MiddleInitial] FROM [ShopVue].[dbo].[fBadge] where FirstName <> '' order by FirstName, LastName asc"

    var q = new Sql.Query()
    q.setQuery(query)
    q.setConnectionid("shopvue")
    q.setParameters(JSON.stringify([]))

    var rqst = new Sql.QueryContextRqst()
    rqst.setQuery(q)

    console.log("Init Employes from ShopVue")
    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.sqlService.queryContext(rqst, metadata);

    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if (response.hasHeader()) {
            var header = response.getHeader()
        } else if (response.hasRows()) {
            results = JSON.parse(response.getRows())
            for (var i = 0; i < results.length; i++) {
                var result = results[i][0]
                var employe = { "id": results[i][0], "firstName": results[i][2], "lastName": results[i][1], "middle": results[i][3] }
                employes[employe.id] = employe
                employeName = employe.firstName
                if (employe.middle !== null) {
                    employeName += " " + employe.middle
                }

                employeName += " " + employe.lastName
                if (employeNames.indexOf(employeName) == -1) {
                    employeNames.push(employeName)
                }

                employeById[employe.id.toUpperCase()] = employe
                employeByName[employeName] = employe
            }
        }
    });

    // Get the results here if the statut is ok.
    stream.on('status', function (status) {
        if (status.code == 0) {
            console.log("Done init Employees")
            initCallback()
        }
    });

    stream.on('end', function () {
        return function (end) {
            // stream end signal
        }
    }(initCallback));
}

// The tools list...
var toolNumbers = []
function initToolNumbers(initCallback) {
    var query = "SELECT DISTINCT EMPLACEMENT FROM OUTIL WHERE [EMPLACEMENT] IS NOT NULL ORDER BY EMPLACEMENT ASC"
    var q = new Sql.Query()
    q.setQuery(query)
    q.setConnectionid("gestion_information")
    q.setParameters(JSON.stringify([]))

    var rqst = new Sql.QueryContextRqst()
    rqst.setQuery(q)

    console.log("Init tool number list")
    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.sqlService.queryContext(rqst, metadata);

    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if (response.hasHeader()) {
            var header = response.getHeader()
        } else if (response.hasRows()) {
            results = JSON.parse(response.getRows())
            var re = /[A-Z]{2}-[A-Z]{3}[0-9]{3}-[0-9]{4}/
            for (var i = 0; i < results.length; i++) {
                var result = results[i][0].trim().toUpperCase()
                toolNumber = result.match(re)
                if (toolNumber != null) {
                    if (toolNumbers.indexOf(toolNumber[0]) == -1) {
                        toolNumbers.push(toolNumber[0])
                    }
                }
            }
        }
    });

    // Get the results here if the statut is ok.
    stream.on('status', function (status) {
        if (status.code == 0) {
            console.log("Done init tool number list")
            initCallback()
        }
    });

    stream.on('end', function () {
        return function (end) {
            // stream end signal
        }
    }(initCallback));
}

// The plaquette
var plaquetteNumbers = []

function initPlaquetteNumbers(initCallback) {
    var query = "SELECT [PART_NBR] FROM [PIECEMASTER] where [TYPE] ='outil' ORDER BY [PART_NBR] ASC"
    var q = new Sql.Query()
    q.setQuery(query)
    q.setConnectionid("catalog_md")
    q.setParameters(JSON.stringify([]))

    var rqst = new Sql.QueryContextRqst()
    rqst.setQuery(q)

    console.log("Init palquette numebers")
    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.sqlService.queryContext(rqst, metadata);

    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if (response.hasHeader()) {
            var header = response.getHeader()
        } else if (response.hasRows()) {
            results = JSON.parse(response.getRows())
            for (var i = 0; i < results.length; i++) {
                var result = results[i][0]
                plaquetteNumbers.push(result.trim())
            }
        }
    });

    // Get the results here if the statut is ok.
    stream.on('status', function (status) {
        if (status.code == 0) {
            console.log("Done init plaquette numbers")
            initCallback()
        }
    });

    stream.on('end', function () {
        return function (end) {
            // stream end signal
        }
    }(initCallback));
}

var programNumbers = []
function initProgramNumbers(initCallback) {

    var query = "SELECT DISTINCT NO_PROGRAMME FROM PROGRAMME order by NO_PROGRAMME asc"

    var q = new Sql.Query()
    q.setQuery(query)
    q.setConnectionid("gestion_information")
    q.setParameters(JSON.stringify([]))

    var rqst = new Sql.QueryContextRqst()
    rqst.setQuery(q)

    console.log("Init Programs numbers")

    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.sqlService.queryContext(rqst, metadata);

    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if (response.hasHeader()) {
            var header = response.getHeader()
        } else if (response.hasRows()) {
            results = JSON.parse(response.getRows())
            for (var i = 0; i < results.length; i++) {
                var result = results[i][0]
                if(programNumbers.indexOf(result.trim()) == -1){
                    programNumbers.push(result.trim())
                }
            }
        }
    });

    // Get the results here if the statut is ok.
    stream.on('status', function (status) {
        if (status.code == 0) {
            console.log("Done init Product numbers")
            initCallback()
        }
    });

    stream.on('end', function () {
        return function (end) {
            // stream end signal
        }
    }(initCallback));
}

var machines = []
var machinesByName = {}

// initialyse machine information here.
function initMachines(initCallback) {
    var s = new Ldap.Search()
    s.setId("safran_ldap")
    s.setBasedn("OU=Shared,OU=Users,OU=MON,OU=CA,DC=UD6,DC=UF6")
    s.setFilter("(&(givenName=Machine*)(objectClass=user))")
    s.setAttributesList(["sAMAccountName", "givenName", "mail"])

    var rqst = new Ldap.SearchRqst()
    rqst.setSearch(s)
    var stream = globular.ldapService.search(rqst)

    stream.on('data', function (response) {
            results = JSON.parse(response.getResult())
            for(var i=0; i < results.length; i++){
                var machine = { "name": results[i][1].split(" ")[1], "email": results[i][2] }
                if(machine.name != undefined){
                    machines.push(machine.name)
                    machinesByName[machine.name] = machine
                }
            }
    });

    // Get the results here if the statut is ok.
    stream.on('status', function (status) {
        if (status.code == 0) {
            console.log("Done init Product numbers")
            initCallback()
        }
    });

    stream.on('end', function () {
        return function (end) {
            // stream end signal
        }
    }(initCallback));
}

var mainPage;

function main() {
    console.log("bienvenue dans bris d'outil!")
    initConnections(function () {
        console.log("connection created successfully")
        // Now I will init the list of product number.
        initEmployes(function () {
            initToolNumbers(function () {
                initSerialNumbers(function () {
                    initPlaquetteNumbers(function () {
                        initMachines(function(){
                            initProgramNumbers(function(){
                                console.log("Data intialisation done!")
                                mainPage = new MainPage()
                            })
                        })
                    })
                })
            })
        })
    })
}


