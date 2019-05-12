
/**
 * The main entry point.
 */
function main() {
    testEcho("Hello globular!");

    // Sql test.
    //  testCreateSqlConnection();

    testSelectQuery()
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
// Sql test.
////////////////////////////////////////////////////////

// Test with MySQL

// Test create a new sql connection.
function testCreateSqlConnection() {
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
}

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

/*
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
function testSelectQuery() {
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
}
*/