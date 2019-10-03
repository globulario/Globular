// The service configuration 
// globularConfig.IP = "127.0.0.1" // comment it when the site is publish.

// The global service object.
var globular = new Globular()

/**
 * The main entry point.
 */
function main() {
    
    testEcho("Hello globular!");

    // Sql test.
    //  testCreateSqlConnection();
    // testPing()
    // testSelectQuery()
    // testDeleteQuery()
    // testInsertQuery()

    // testGetFileInfo()
    // testFilePane()

    // testCreatePersistenceConnection()
    // testPersistencePing()
    // testPersistenceFind()

    // display table.
    //displayTable()

    //testEvent()
}

/////////////////////////////////////////////////////////
// Echo test.
////////////////////////////////////////////////////////
function testEcho(str) {

    // Create a new request.
    var request = new Echo.EchoRequest();
    request.setMessage(str);

    globular.echoService.echo(request, {}, function (err, response) {
        // ...
        console.log(response.getMessage())
    });

    // Now I will test with promise
    globular.echoServicePromise.echo(request)
        .then((resp) => {
            console.log(resp.getMessage())
        })
        .catch((error) => {
            console.log(error)
        })
}

/////////////////////////////////////////////////////////
// event test.
////////////////////////////////////////////////////////
function testEvent(){
    // The first step is to subscribe to an event channel.
    var rqst = new EventBus.SubscribeRequest()
    rqst.setName("my topic")

    var stream = globular.eventService.subscribe(rqst, {});

    // Get the stream and set event on it...
    stream.on('data', function (rsp) {
        if (rsp.hasUuid()) {
            console.log(rsp.getUuid())
        } else if (rsp.hasEvt()) {
            var evt = rsp.getEvt()
            var data = new TextDecoder("utf-8").decode(evt.getData());
            console.log(data)
        }
    });

    stream.on('status', function (status) {
        if (status.code == 0) {
            console.log("---> end of subscription")
        }
        console.log(status)
    });

    stream.on('end', function (end) {
        // stream end signal
        console.log("---> end of subscription") 
    });
}

/////////////////////////////////////////////////////////
// file test.
////////////////////////////////////////////////////////
function testGetFileInfo() {
    var request = new File.GetFileInfoRequest();
    request.setPath("/home/dave/Pictures/unnamed.png")
    request.setThumnailheight(256)
    request.setThumnailwidth(256)

    globular.fileServicePromise.getFileInfo(request)
        .then((resp) => {
            var data = JSON.parse(resp.getData())
            console.log(data)
        })
        .catch((error) => {
            console.log(error)
        })
}

/////////////////////////////////////////////////////////
// file pane test
/////////////////////////////////////////////////////////
var uint8array
// Merge tow array together.
function mergeTypedArraysUnsafe(a, b) {
    var c = new a.constructor(a.length + b.length);
    c.set(a);
    c.set(b, a.length);
    return c;
}

function testFilePane() {
    // Here I will create the file panel...
    var filePane = document.createElement("file-pane-element")
    filePane.height = 180;
    
    // This is where the file will be save.
    filePane.path = "/test/filePane";

    // Now the pane event...
    filePane.ondelete = function (path) {
        console.log("delete the file", path)
        var request = new File.DeleteFileRequest();
        request.setPath(path)

        globular.fileServicePromise.deleteFile(request)
            .then((resp) => {
                mainPage.displayMessage("<div>Le fichier " + path + "</span>est supprimer! </div>", 3000)
            })
            .catch((error) => {
                console.log(error)
            })
    }

    // Here I will display the file when the user click on it.
    filePane.onopen = function (fileInfo) {
        // I will get the file from the server.
        var req = new XMLHttpRequest();
        var urlToSend = "http://" + globular.config.Domain + ":" + globular.config.PortHttp + fileInfo.Path
        req.open("GET", urlToSend, true);
        req.responseType = "blob";
        req.onload = function (fileInfo) {
            return function (event) {
                var blob = req.response;
                var fileName = req.getResponseHeader("fileName") //if you have the fileName header available
                var link = document.createElement('a');
                if (blob.size > 100) {
                    link.href = window.URL.createObjectURL(blob);
                } else {
                    link.href = fileInfo.Thumbnail
                }
                link.download = fileName;
                link.click();

            }
        }(fileInfo);
        req.send();
    }

    // The upload function.
    filePane.uploadHandler = function (fileInfo) {

        // First of all the value must be upload on the server uploads directory...
        // and move after it.
        // Upload the file here.
        var formData = new FormData()
        formData.append("multiplefiles", fileInfo.Local, fileInfo.Name)
        
        // The path is the parent path here...
        formData.append("path", fileInfo.Path.replace("/" + fileInfo.Name, ""))
        
        // Use the post function to upload the file to the server.
        var xhr = new XMLHttpRequest()
        xhr.open('POST', '/uploads', true)

        // In case of error or success...
        xhr.onload = function (e) {
            if (xhr.readyState === 4) {

            }
        }

        // now the progress event...
        xhr.upload.onprogress = function (e) {

        }

        xhr.send(formData);
    }

    // set the new file event.
    filePane.onnewfile = function (fileInfo) {
        console.log(fileInfo)
    }

    var div = document.createElement("div")
    div.style.width = "350px"
    div.style.height = "100px"
    document.body.appendChild(div)
    var btn = document.createElement("paper-button")
    btn.innerHTML = "save"
    btn.style.marginTop = "10px"
    div.appendChild(filePane)
    div.appendChild(btn)

    // Save the files.
    btn.onclick = function (filePane) {
        return function () {
            filePane.saveAll()
        }
    }(filePane)

    // Here I will set the file pane.
    var rqst = new File.ReadDirRequest()
    rqst.setPath(filePane.path)
    rqst.setRecursive(false)
    rqst.setThumnailwidth(256)
    rqst.setThumnailheight(256)

    var stream = globular.fileService.readDir(rqst, {});
    uint8array = new Uint8Array();

    // Get the stream and set event on it...
    stream.on('data', function (resp) {
        uint8array = mergeTypedArraysUnsafe(uint8array, resp.getData())
    });

    stream.on('status', function(filePane){
    return function (status) {
        if (status.code == 0) {
            var jsonStr = new TextDecoder("utf-8").decode(uint8array);
            filePane.setDirInfo(jsonStr)
        }
    }}(filePane));

    stream.on('end', function (end) {
        // stream end signal
        console.log("---> end: ", end)
    });
}

/////////////////////////////////////////////////////////
// Persistence test. (MongoDB backend)
////////////////////////////////////////////////////////

function testCreatePersistenceConnection() {
    var rqst = new Persistence.CreateConnectionRqst();
    var c = new Persistence.Connection();
    c.setId("mongo_db_test_connection")
    c.setName("TestMongoDB") // already exist...
    c.setUser("")
    c.setPassword("")
    c.setPort(27017)
    c.setStore(0) // mongodb
    c.setHost("localhost")
    c.setTimeout(10)

    rqst.setConnection(c)

    globular.persistenceService.createConnection(rqst, {}, function (err, rsp) {
        // ...
        console.log(rsp.getResult())
    });
}

function testPersistencePing() {
    var rqst = new Persistence.PingConnectionRqst()
    rqst.setId("mongo_db_test_connection")

    globular.persistenceServicePromise.ping(rqst)
        .then((rsp) => {
            console.log(rsp.getResult())
        })
        .catch((error) => {
            console.log(error)
        })
}

// Test Find existing values...
var testPersistenceResults = []
function testPersistenceFind() {
    var rqst = new Persistence.FindRqst()
    rqst.setId("mongo_db_test_connection")
    rqst.setDatabase("TestMongoDB")
    rqst.setCollection("Employees")
    rqst.setQuery('{"first_name": "Anneke"}' /*"{}"*/)
    rqst.setFieldsList(["_id", "birth_date"]) // here I will get only the _id and the birth date.

    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.persistenceService.find(rqst, metadata);

    // Get the stream and set event on it...
    stream.on('data', function (rsp) {
        testPersistenceResults = testPersistenceResults.concat(JSON.parse(rsp.getJsonstr()))
    });

    stream.on('status', function (status) {
        if (status.code == 0) {
            console.log(testPersistenceResults)
        }
    });

    stream.on('end', function (end) {
        // stream end signal
    });
}

/////////////////////////////////////////////////////////
// Sql test.
////////////////////////////////////////////////////////

// Test with MySQL

// Test create a new sql connection.
/*function testCreateSqlConnection() {
    var rqst = new Sql.CreateConnectionRqst();
    var c = new Sql.Connection();
    c.setId("employees_db")
    c.setName("employees")
    c.setUser("test")
    c.setPassword("password")
    c.setPort(3306)
    c.setDriver("mysql")
    c.setHost("localhost")
    c.setCharset("utf8")

    rqst.setConnection(c)

    globular.sqlService.createConnection(rqst, {}, function (err, rsp) {
        // ...
        console.log(rsp.getResult())
    });
}*/

// Test a select query.
function testSelectQuery() {
    var rqst = new Sql.QueryContextRqst()
    var q = new Sql.Query()
    q.setQuery("SELECT * FROM employees.employees WHERE gender=?")
    q.setConnectionid("employees_db")
    q.setParameters(JSON.stringify(["F"]))

    rqst.setQuery(q)
    // var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.sqlService.queryContext(rqst, {});

    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if (response.hasHeader()) {
            var header = response.getHeader()
            console.log(JSON.parse(header));
        } else if (response.hasRows()) {
            var rows = response.getRows()
            console.log(JSON.parse(rows));
        }
    });

    stream.on('status', function (status) {
        console.log(status.code);
        console.log(status.details);
        console.log(status.metadata);
    });

    stream.on('end', function (end) {
        // stream end signal
    });
}

function testSqlPing() {
    var rqst = new Sql.PingConnectionRqst()
    rqst.setId("employees_db")

    globular.sqlServicePromise.ping(rqst)
        .then((rsp) => {
            console.log(rsp.getResult())
        })
        .catch((error) => {
            console.log(error)
        })
}

function testInsertQuery() {
    var rqst = new Sql.ExecContextRqst()
    var q = new Sql.Query()
    q.setQuery("INSERT INTO employees.employees (emp_no, first_name, last_name, gender, hire_date, birth_date) VALUE(?,?,?,?,?,?)")
    q.setConnectionid("employees_db")
    q.setParameters(JSON.stringify([200000, 'Dave', 'Courtois', 'M', '2007-07-01', '1976-01-29']))

    rqst.setQuery(q)
    rqst.setTx(false)

    globular.sqlServicePromise.execContext(rqst)
        .then((rsp) => {
            var affectedRow = rsp.getAffectedrows()
            var lastId = rsp.getLastid()
            console.log("affected rows: ", affectedRow, " last id: ", lastId)
        })
        .catch((error) => {
            console.log(error)
        })
}

function testDeleteQuery() {
    var rqst = new Sql.ExecContextRqst()
    var q = new Sql.Query()
    q.setQuery("DELETE FROM employees.employees WHERE emp_no=?")
    q.setConnectionid("employees_db")
    q.setParameters(JSON.stringify([200000]))

    rqst.setQuery(q)
    rqst.setTx(false)

    globular.sqlServicePromise.execContext(rqst)
        .then((rsp) => {
            var affectedRow = rsp.getAffectedrows()
            var lastId = rsp.getLastid()
            console.log("affected rows: ", affectedRow, " last id: ", lastId)
        })
        .catch((error) => {
            console.log(error)
        })
}

// Test with Sql Server and odbc connector.
function testCreateSqlConnection() {
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

    globular.sqlService.createConnection(rqst, {}, function (err, rsp) {
        // ...
        console.log(rsp.getResult())
    });
}

var table_example_data = [];
var table_example_header = null;
/**
 * That function is use in table example to get a large array of values.
 * @param {*} input 
 */
function getEmployees(callback) {
    var rqst = new Sql.QueryContextRqst()
    var q = new Sql.Query()
    q.setQuery("SELECT * FROM employees.employees WHERE gender=?")
    q.setConnectionid("employees_db")
    q.setParameters(JSON.stringify(["F"]))
    rqst.setQuery(q)
    var metadata = {}// { 'content-type': 'application/grpc-web-text' };
    var stream = globular.sqlService.queryContext(rqst, metadata);
    table_example_data = [] // reset it content.

    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if (response.hasHeader()) {
            table_example_header = JSON.parse(response.getHeader())
        } else if (response.hasRows()) {
            table_example_data = table_example_data.concat(JSON.parse(response.getRows()))
            for (var i = 0; i < table_example_header.length; i++) {
                for(var j=0; j < table_example_data.length; j++){
                    if (table_example_header[i].typeInfo.Name == "DATE") {
                        // convert to date.
                        table_example_data[j][i] = new Date(table_example_data[j][i])
                    }
                }
            }
        }
    });

    stream.on('status', function (status) {
        if (status.code == 0) {
            callback(table_example_header, table_example_data)
        }
    });

    stream.on('end', function (end) {
        // stream end signal
    });
}

function displayTable() {
    var div = document.createElement("div")
    div.id = "table_example_div"
    div.style.position = "relative"
    document.body.appendChild(div)

    /** Get employes get values from sql tables... */
    getEmployees(function (infos, employees) {
        /* Let's the hack begin */

        /* Step 1 create the table */
        var table = document.createElement("table-element")

        /* Step 2 create the header */
        var header = document.createElement("table-header-element")
        table.appendChild(header)

        /* Set various style here... */
        table.rowheight = 29
        table.style.width = "880px"
        table.style.maxHeight = "500px";

        /* set the table data */
        table.data = employees

        header.fixed = true;
        for (var i = 0; i < infos.length; i++) {
            /** Here I will create the header cell */
            var headerCell = document.createElement("table-header-cell-element")

            /** Optionaly set a table sorter and filter */
            headerCell.innerHTML = "<table-sorter-element></table-sorter-element><div>" + infos[i].name + "</div> <table-filter-element></table-filter-element>"
            if (i == 0) {
                /** Here if you want to overide onrender to modify how cell in that row will be display... */
                headerCell.onrender = function (div, value) {
                    div.innerHTML = value
                    div.style.textDecoration = "underline"
                    div.onclick = function (value) {
                        return function () {
                            alert("you select " + value + "!")
                        }
                    }(value)

                    div.onmouseover = function () {
                        this.style.cursor = "pointer"
                    }

                    div.onmouseout = function () {
                        this.style.cursor = "default"
                    }
                }
            }

            /** Here if the sql type is a date I will convert-it to a JavaScript date and set it display function. */
            if (infos[i].typeInfo.Name == "DATE") {
                /** Here is the cell is a date */
                headerCell.onrender = function (div, value) {
                    if (value != undefined) {
                        /* Convert string value to date. */
                        div.innerHTML = value.toLocaleDateString()
                    }
                }
            }

            if (i == 4) {
                /** set the size of a column */
                headerCell.style.minWidth = "60px"
            }

            header.appendChild(headerCell)
        }

        /** append the table in a div. **/
        var parentDiv = document.getElementById("table_example_div")
        parentDiv.appendChild(table)
    })
}

// Test a select query.
/*function testSelectQuery() {
    var rqst = new Sql.QueryContextRqst()
    var q = new Sql.Query()
    q.setQuery("SELECT * FROM [BrisOutil].[dbo].[Bris] WHERE product_id LIKE ?")
    q.setConnectionid("bris_outil")
    q.setParameters(JSON.stringify(["50-%"]))

    rqst.setQuery(q)
    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.sqlService.queryContext(rqst, metadata);

    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if(response.hasHeader()){
            var header = response.getHeader()
            console.log(JSON.parse(header));
        }else if(response.hasRows()){
            var rows = response.getRows()
            console.log(JSON.parse(rows));
        }
    });

    stream.on('status', function (status) {
        console.log(status.code);
        console.log(status.details);
        console.log(status.metadata);
    });

    stream.on('end', function (end) {
        // stream end signal
        console.log("---> end: ", end)
    });
}*/
