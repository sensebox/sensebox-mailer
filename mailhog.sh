#!/bin/bash

set -e

docker run -d --name mailhog -p "1025:1025" -p "8025:8025" mailhog/mailhog

branch=$(git rev-parse --abbrev-ref HEAD)
ts=$(TZ=UTC git log --date=local --pretty=format:"%ct" -n 1)
hash=$(TZ=UTC git log --date=local --pretty=format:"%h" -n 1)

$GOBIN/statik -src=templates -f
go build -o sensebox-mailer -ldflags "-X main.branch=$branch -X main.ts=$ts -X main.hash=$hash" cmd/sensebox-mailer/*.go

export SENSEBOX_MAILER_SERVER_CERT=$(cat out/mailer_server.crt)
export SENSEBOX_MAILER_SERVER_KEY=$(cat out/mailer_server.key)
export SENSEBOX_MAILER_CA_CERT=$(cat out/openSenseMapCA.crt)
export SENSEBOX_MAILER_SMTP_SERVER=localhost
export SENSEBOX_MAILER_SMTP_PORT=1025
export SENSEBOX_MAILER_SMTP_USER=smtpuser
export SENSEBOX_MAILER_SMTP_PASSWORD=smtppassword
export SENSEBOX_MAILER_FROM_DOMAIN=sensebox.de
export SENSEBOX_MAILER_FROM_NAME_PREFIX=senseBox

function cleanup() {
  echo 'cleanup!'
  docker stop mailhog
  docker rm mailhog
}
trap cleanup EXIT

./sensebox-mailer