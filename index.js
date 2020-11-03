var fs = require("fs");
var https = require("https");
var http = require("http");

function checkMailhog () {
  console.log("")
  setTimeout(function () {
    var apiReq = http.request({
      hostname: "localhost",
      port: 8025,
      path: "/api/v2/messages",
      method: "GET"
    }, function (res) {
      var responseString = "";
      res.on("data", function(data) {
        // process.stdout.write(data);
        responseString += data;
      });
      res.on("end", function () {
        var j = JSON.parse(responseString);
        var mail = j.items[0];

        var body = mail.Content.Body;

        console.log(body)
      })
    });

    apiReq.end();

    apiReq.on("error", function(e) {
      console.error(e);
    });
  }, 5000);
}

function sendMail () {
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
    res.on("end", function () {
      checkMailhog();
    })
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
        origin: "https://testing.opensensemap.org",
        user: {
          name: "Gerald",
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
}

sendMail();
