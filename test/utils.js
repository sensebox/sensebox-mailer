import { readFileSync } from "fs";
import { request } from "https";

const requestOptions = {
  hostname: "localhost",
  port: 3924,
  path: "/",
  method: "POST",
  key: readFileSync("./out/mailer_client.key"),
  cert: readFileSync("./out/mailer_client.crt"),
  ca: readFileSync("./out/openSenseMapCA.crt"),
  ecdhCurve: "auto"
};

function validateArgs (args) {
  if (!args) {
    return [false, 'sendMail requires nonempty arguments'];
  }

  const argKeys = Object.keys(args);
  for (const requiredKey of ['payload', 'callback']) {
    if (!argKeys.includes(requiredKey)) {
      return [false, `Could not find "${requiredKey}" in arguments`];
    }
  }

  return [true];
}

export function sendMail (args) {
  const [valid, msg] = validateArgs(args);
  if (!valid) {
    throw new Error(msg);
  }

  const req = request(requestOptions, function(res) {
    res.on("data", function(data) {
      process.stdout.write(data);
    });
    res.on("end", function () {
      args.callback();
      // checkMailhog();
    })
  });


  req.write(JSON.stringify(args.payload));

  req.end();

  req.on("error", function(e) {
    throw e;
  });
}
