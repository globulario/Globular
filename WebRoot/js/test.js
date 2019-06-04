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

    // testCreatePersistenceConnection()
    // testPersistencePing()
    // testPersistenceFind()
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
// file test.
////////////////////////////////////////////////////////
function testGetFileInfo(){
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
// Persistence test. (MongoDB backend)
////////////////////////////////////////////////////////

function testCreatePersistenceConnection(){
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

function testPersistencePing(){
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
function testPersistenceFind(){
    var rqst = new Persistence.FindRqst()
    rqst.setId("mongo_db_test_connection")
    rqst.setDatabase("TestMongoDB")
    rqst.setCollection("Employees")
    rqst.setQuery( '{"first_name": "Anneke"}' /*"{}"*/)
    rqst.setFieldsList(["_id", "birth_date"]) // here I will get only the _id and the birth date.

    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.persistenceService.find(rqst, metadata);
    
    // Get the stream and set event on it...
    stream.on('data', function (rsp) {
        testPersistenceResults = testPersistenceResults.concat(JSON.parse(rsp.getJsonstr()))
    });

    stream.on('status', function (status) {
        if(status.code == 0){
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
    q.setQuery("SELECT first_name, last_name FROM employees.employees WHERE gender=?")
    q.setConnectionid("employees_db")
    q.setParameters(JSON.stringify(["F"]))

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
    });
}

function testSqlPing(){
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

function testInsertQuery(){
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

function testDeleteQuery(){
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
