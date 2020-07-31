#!/bin/bash
COMMON_NAME=nitro-dev

# make the directories if they do not exist
if [ ! -d "/etc/nitrod/certs" ]; then
    mkdir -p /etc/nitrod/certs
fi

# files
CA_KEY_PATH="/etc/nitrod/certs/ca.key"
CA_TRUST_PATH="/etc/nitrod/certs/ca.crt"
SERVER_PRIVATE_KEY_PATH="/etc/nitrod/certs/nitrod.key"
SERVER_CSR_PATH="/etc/nitrod/certs/nitrod.csr"

# generate the certificate authority if it does not exist
if [ ! -f ${CA_KEY_PATH} ]; then
  echo "INFO: creating the certificate authority..."
  openssl genrsa \
    -passout pass:1111 \
    -des3 \
    -out ${CA_KEY_PATH} 4096
fi

# generate the trust certificate if we don't have one
if [ ! -f ${CA_TRUST_PATH} ]; then
  echo "INFO: creating the trust certificate..."
  openssl req \
    -passin pass:1111 \
    -new -x509 \
    -days 3650 \
    -key ${CA_KEY_PATH} \
    -out ${CA_TRUST_PATH} \
    -subj "/CN=${COMMON_NAME}"
fi

# generate a private key for the nitrod api
if [ ! -f ${SERVER_PRIVATE_KEY_PATH} ]; then
  echo "INFO: creating the nitrod private key..."
  openssl genrsa \
    -passout pass:1111 \
    -des3 \
    -out ${SERVER_PRIVATE_KEY_PATH} 4096
fi

# create the certificate signing request for the nitrod api
if [ ! -f ${SERVER_CSR_PATH} ]; then
  echo "INFO: creating the server certificate signing request for nitrod..."
  openssl req \
    -passin pass:1111 \
    -new \
    -key ${SERVER_PRIVATE_KEY_PATH} \
    -out ${SERVER_CSR_PATH} \
    -subj "/CN=${COMMON_NAME}"
fi
