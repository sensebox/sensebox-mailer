#!/bin/bash

# Generate a CA
# TODO check for ca cert files and skip
openssl req -new -x509 -days 3650 -subj "/emailAddress=ca@sensebox/C=DE/ST=NRW/L=Muenster/O=senseBox/OU=senseBoxDevOps/CN=mailer_ca" -keyout ca_key.pem -out ca_cert.pem

# Generate a Key for the server
openssl ecparam -genkey -name secp521r1 -out server_key.pem

# Generate a signing request for the server certificate
openssl req -new -subj "/emailAddress=server@sensebox/C=DE/ST=NRW/L=Muenster/O=senseBox/OU=senseBoxDevOps/CN=localhost/subjectAltName=DNS.1=mailer,DNS.2=sensebox-mailer" -key server_key.pem -out server_csr.pem

# Sign it
openssl x509 -req -days 3650 -in server_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -CAcreateserial -out server_cert.pem

# Generate a Key for the client
openssl ecparam -genkey -name secp521r1 -out client_key.pem

# Generate a signing request for the client certificate
openssl req -new -subj "/emailAddress=client@sensebox/C=DE/ST=NRW/L=Muenster/O=senseBox/OU=senseBoxDevOps/CN=localhost/subjectAltName=DNS.1=api,DNS.2=osem-api" -key client_key.pem -out client_csr.pem

# Sign it
openssl x509 -req -days 3650 -in client_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -CAcreateserial -out client_cert.pem

# generate docker-compose.yml environment sections..
cat << EOF > client_env.yml
      OSEM_mailer_cert: |-
$(cat client_cert.pem | while read line; do
  echo "        $line"
done)
      OSEM_mailer_key: |-
$(cat client_key.pem | while read line; do
  echo "        $line"
done)
      OSEM_mailer_ca: |-
$(cat ca_cert.pem | while read line; do
  echo "        $line"
done)
EOF

cat << EOF > server_env.yml
      SENSEBOX_MAILER_SERVER_CERT: |-
$(cat server_cert.pem | while read line; do
  echo "        $line"
done)
      SENSEBOX_MAILER_SERVER_KEY: |-
$(cat server_key.pem | while read line; do
  echo "        $line"
done)
      SENSEBOX_MAILER_CA_CERT: |-
$(cat ca_cert.pem | while read line; do
  echo "        $line"
done)
EOF
