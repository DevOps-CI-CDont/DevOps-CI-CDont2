// Listen on a specific host via the HOST environment variable
var host = process.env.HOST || "localhost";
// Listen on a specific port via the PORT environment variable
var port = 3001;
var cors_proxy = require("cors-anywhere");
cors_proxy
  .createServer({
    originWhitelist: [], // Allow all origins
    requireHeader: ["origin", "x-requested-with"],
    redirectSameOrigin: true,
    setHeaders: function(res, req) {
      res.setHeader("Access-Control-Allow-Credentials", "true");
    }
  })
  .listen(port, host, function () {
    console.log("Running CORS Anywhere on " + host + ":" + port);
  });
