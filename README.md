sensebox-mailer (WIP)
====================

This project can be built using [gb](getgb.io).

In order to run it, you have to follow these steps:

### 1. Generate server and client certificates

Read `generateServerAndClientCertificates.sh`. The file generates a self signed CA certificate file and two certificate and key files signed by the CA. You need the CA certificate and the server certificate and key on the server and the CA certificate and the client certificate and key on the client.

### 2. Build it

Preferred way to build the mailer is to build a Docker image from it. Just run `docker build -t sensebox/sensebox-mailer .` in this directory and you should have a working version of the mailer.

### 3. Run it

The server can only be configured through environment variables. The easiest way is to use [docker-compose](https://github.com/docker/compose) to achieve this. Just use the supplied `docker-compose.yml`. All config keys are prefixed with `SENSEBOX_MAILER_`. Consider this when consulting the table below.

You should configure the following variables:

| key | comment | optional |
|-----|---------|---------------------------------------------------------------------------|
| `CA_CERT` | the certificate of your CA. Server and client should be signed by this CA | y |
| `SERVER_CERT` | the server certificate | y |
| `SERVER_KEY` | the key of the server certificate | y |
| `SMTP_SERVER` | the smtp server address | y |
| `SMTP_PORT` | the smtp server port | y |
| `SMTP_USER` | the smtp server user | y |
| `SMTP_PASSWORD` | the smtp server password | y |
| `FROM_DOMAIN` | the domain you are sending from | y |
| `HONEYBADGER_APIKEY` | api key for honeybadger error reporting |  |

