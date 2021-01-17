# sensebox-mailer

This project is the mailer used by the [openSenseMap-API](https://github.com/sensebox/openSenseMap-API) and other services in the openSenseMap stack. It is written in Golang and thus can be compiled into a single binary.

## Development

- Clone the project into your GOPATH.
- Install certstrap (`go get -u github.com/square/certstrap`)
- Install statik (`go get -u github.com/rakyll/statik`)

### Compilation

    go build -o sensebox-mailer cmd/sensebox-mailer/*.go

### Adding new mail templates

1. Create new html files in `templates`. Ideally named `yourNewTemplate_en.html` or `yourNewTemplate_de.html`. Please always create both files for german (`_de.html`) and english (`_en.html`).
1. Edit `templates/templates.json` and add your template to the array.
1. Run `statik -src=templates -f`
1. Compile the project `go build -o sensebox-mailer cmd/sensebox-mailer/*.go`

### Running

The mailer relies on some environment variables for configuration. You can find these at the bottom of the README. Before running, you should generate certificates (`./genCerts.sh`).

A good mailserver for development and testing is [mailhog](https://github.com/mailhog/MailHog). You should start an instance of it

    docker run -d --name mailhog -p "1025:1025" -p "8025:8025" mailhog/mailhog

A good starting point for a bash script for development is `mailhog.sh`.

Running this script:

- starts `mailhog` docker container
- runs statik and compiles all templates
- compiles the sensebox-mailer go binary
- exports all necessary env variables
- runs the sensebox-mailer go binary

Before running this script its important to run [`genCerts.sh`](#1-generate-server-and-client-certificates).

You run `node index.js` to test your templates. Visit `localhost:8025` and check the **Inbox** for new mails.

Change the value of `template` inside `index.js` to test different `templates`.


## HTTP interface

Upon running the compiled binary, the mailer will expose a HTTP interface running on port 3924. Clients can send mails by sending POST requests to `/` with JSON payload. Clients should authenticate using TLS client certificates.

An example payload should look like this:

    [
      {
        "template": "registration",       // required. The template you want to render
        "lang": "en",                     // required. the language to use
        "recipient": {                    // required. sould have the keys address and name
          "address": "email@address.com",
          "name": "Philip J. Fry"
        },
        "payload": {
          "foo": "bar",
          "baz": {
            "boing": "boom"
          }
        },
        "attachment": {               // optional. should contain keys filename and contents
          "filename": "senseBox.ino", // filename of the attachment
          "contents": "<file contents in base64>" // file contents encoded in base64
        }
      },
      ...
    ]

The root of the JSON payload should always be an array containing single requests. Required keys for the single requests are `template`, `lang`, `recipient` and `payload`. The key `attachment` is optional.

## Production use

In order to run it, you have to follow these steps:

### 1. Generate server and client certificates

- Install certstrap (`go get -u github.com/square/certstrap`)
- Run `./genCerts.sh` to generate a self signed CA along with certificates for the mailer and a client. The certificates are used for client/server TLS.

### 2. Build it

Preferred way to build the mailer is to build a Docker image from it. Just run `docker build -t sensebox/sensebox-mailer .` in this directory and you should have a working version of the mailer.

You can also build it using `docker-compose build`

### 3. Run it

The server can only be configured through environment variables. The easiest way is to use [docker-compose](https://github.com/docker/compose) to achieve this. Just use the supplied `docker-compose.yml`. All config keys are prefixed with `SENSEBOX_MAILER_`. Consider this when consulting the table below.

Another option is to create a `docker-compose.override.yml` to override the values in the original `docker-compose.yml` with values of the generated `server_env.yml` and run it with `docker-compose up`

You should configure the following variables:

| key | comment | optional |
|-----|---------|---------------------------------------------------------------------------|
| `CA_CERT` | the certificate of your CA. Server and client should be signed by this CA |  |
| `SERVER_CERT` | the server certificate |  |
| `SERVER_KEY` | the key of the server certificate |  |
| `SMTP_SERVER` | the smtp server address |  |
| `SMTP_PORT` | the smtp server port |  |
| `SMTP_USER` | the smtp server user |  |
| `SMTP_PASSWORD` | the smtp server password |  |
| `FROM_DOMAIN` | the domain you are sending from |  |
