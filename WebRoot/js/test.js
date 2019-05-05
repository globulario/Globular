
/**
 * The main entry point.
 */
function main() {
    testEcho("Hello globular!");
}

function testEcho(str) {

    // Create a new request.
    var request = new EchoRequest();
    request.setMessage(str);

    globular.echoService.echo(request, {}, function (err, response) {
        // ...
        console.log(response.getMessage())
    });
}