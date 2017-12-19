# sensebox-mailer

This project uses `dep` and can be built using `go build`.

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
| `HONEYBADGER_APIKEY` | api key for honeybadger error reporting | y |

