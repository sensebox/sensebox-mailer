import { request as _request } from "http";
import { sendMail } from "./utils.js";

function checkMailhog () {
  console.log("")
  setTimeout(function () {
    var apiReq = _request({
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

sendMail({ payload, callback: checkMailhog });
