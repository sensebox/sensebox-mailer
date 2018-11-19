var fs = require("fs");
var https = require("https");

var options = {
  hostname: "localhost",
  port: 3924,
  path: "/",
  method: "POST",
  key: fs.readFileSync("./out/mailer_client.key"),
  cert: fs.readFileSync("./out/mailer_client.crt"),
  ca: fs.readFileSync("./out/openSenseMapCA.crt"),
  ecdhCurve: "auto"
};

var req = https.request(options, function(res) {
  res.on("data", function(data) {
    process.stdout.write(data);
  });
});

var payload = [
  {
    template: "newBoxHackAir",
    lang: "de_DE",
    recipient: {
      address: "address@mail.com",
      name: "Firstname Lastname"
    },
    payload: {
      origin: "webseite",
      user: {
        name: "Gera",
        firstname: "Gerald",
        lastname: "P",
        apikey: "123"
      },
      box: {
        name: "YOUR_SENSEBOX_NAME",
        _id: "YOUR_SENSEBOX_ID"
      }
    }
  }
];

req.write(JSON.stringify(payload));

req.end();

req.on("error", function(e) {
  console.error(e);
});