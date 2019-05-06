
/**
 * The main entry point.
 */
function main() {
    testEcho("Hello globular!");

    // Sql test.
    testCreateSqlConnection();

    testSelectQuery()
}

/////////////////////////////////////////////////////////
// Echo test.
////////////////////////////////////////////////////////
function testEcho(str) {

    // Create a new request.
    var request = new EchoRequest();
    request.setMessage(str);

    globular.echoService.echo(request, {}, function (err, response) {
        // ...
        console.log(response.getMessage())
    });
}

/////////////////////////////////////////////////////////
// Sql test.
////////////////////////////////////////////////////////

// Test create a new sql connection.
function testCreateSqlConnection() {
    var rqst = new CreateConnectionRqst();
    var c = new SqlConnection();
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
    var rqst = new QueryContextRqst()
    var q = new Query()
    q.setQuery("SELECT first_name, last_name FROM employees.employees WHERE gender=?")
    q.setConnectionid("employees_db")
    q.setParameters(JSON.stringify(["F"]))

    rqst.setQuery(q)
    var metadata = { 'custom-header-1': 'value1' };
    var stream = globular.sqlService.queryContext(rqst, metadata);
    
    // Get the stream and set event on it...
    stream.on('data', function (response) {
        if(response.hasRows()){
            var rows = response.getRows()
            console.log(JSON.parse(rows));
        }else if(response.hasHeader()){
            var header = response.getHeader()
            console.log(JSON.parse(header));
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