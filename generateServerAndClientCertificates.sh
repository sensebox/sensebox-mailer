#!/bin/bash

# Generate a CA
# TODO check for ca cert files and skip
openssl req -new -x509 -days 3650 -subj "/emailAddress=ca@sensebox/C=DE/ST=NRW/L=Muenster/O=senseBox/OU=senseBoxDevOps/CN=mailer_ca" -keyout ca_key.pem -out ca_cert.pem

# Generate a Key for the server
openssl ecparam -genkey -name secp521r1 -out server_key.pem

# Generate a signing request for the server certificate
openssl req -new -subj "/emailAddress=server@sensebox/C=DE/ST=NRW/L=Muenster/O=senseBox/OU=senseBoxDevOps/CN=localhost" -key server_key.pem -out server_csr.pem

# Sign it
openssl x509 -req -days 3650 -in server_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -CAcreateserial -out server_cert.pem

# Generate a Key for the client
openssl ecparam -genkey -name secp521r1 -out client_key.pem

# Generate a signing request for the client certificate
openssl req -new -subj "/emailAddress=client@sensebox/C=DE/ST=NRW/L=Muenster/O=senseBox/OU=senseBoxDevOps/CN=localhost" -key client_key.pem -out client_csr.pem

# Sign it
openssl x509 -req -days 3650 -in client_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -CAcreateserial -out client_cert.pem

