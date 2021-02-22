#!/bin/bash

set -e

function cleanup() {
  echo 'cleanup!'
  sudo docker stop mailer
  sudo docker rm mailer
}
trap cleanup EXIT

sudo docker run \
  --name=mailer \
  --network=host \
  -e "SENSEBOX_MAILER_SERVER_CERT=$(cat out/mailer_server.crt)" \
  -e "SENSEBOX_MAILER_SERVER_KEY=$(cat out/mailer_server.key)" \
  -e "SENSEBOX_MAILER_CA_CERT=$(cat out/openSenseMapCA.crt)" \
  -e "SENSEBOX_MAILER_SMTP_SERVER=localhost" \
  -e "SENSEBOX_MAILER_SMTP_PORT=1025" \
  -e "SENSEBOX_MAILER_SMTP_USER=smtpuser" \
  -e "SENSEBOX_MAILER_SMTP_PASSWORD=smtppassword" \
  -e "SENSEBOX_MAILER_FROM_DOMAIN=sensebox.de" \
  -e "SENSEBOX_MAILER_FROM_NAME_PREFIX=senseBox" \
  mailer
