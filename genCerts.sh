
#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

CA_NAME=openSenseMapCA
SERVICE=mailer

echo "Generate root CA"
certstrap init --passphrase "" --expires "10 years" --common-name "${CA_NAME}"

echo "Create certificate requests for the server and client"
certstrap request-cert --passphrase "" --key-bits "4096" --common-name "${SERVICE}_server" --domain "${SERVICE},localhost"
certstrap request-cert --passphrase "" --key-bits "4096" --common-name "${SERVICE}_client" --domain "${SERVICE},localhost"

echo "Sign the certificate requests"
certstrap sign --passphrase "" --expires "10 years" --CA "$CA_NAME" "${SERVICE}_server"
certstrap sign --passphrase "" --expires "10 years" --CA "$CA_NAME" "${SERVICE}_client"

# generate docker-compose.yml environment sections..
cat << EOF > client_env.yml
      OSEM_mailer_cert: |-
$(cat out/${SERVICE}_client.crt | while read line; do
  echo "        $line"
done)
      OSEM_mailer_key: |-
$(cat out/${SERVICE}_client.key | while read line; do
  echo "        $line"
done)
      OSEM_mailer_ca: |-
$(cat out/${CA_NAME}.crt | while read line; do
  echo "        $line"
done)
EOF

cat << EOF > server_env.yml
      SENSEBOX_MAILER_SERVER_CERT: |-
$(cat out/${SERVICE}_server.crt | while read line; do
  echo "        $line"
done)
      SENSEBOX_MAILER_SERVER_KEY: |-
$(cat out/${SERVICE}_server.key | while read line; do
  echo "        $line"
done)
      SENSEBOX_MAILER_CA_CERT: |-
$(cat out/${CA_NAME}.crt | while read line; do
  echo "        $line"
done)
EOF
