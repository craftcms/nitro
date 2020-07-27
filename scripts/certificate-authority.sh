#!/bin/bash
COMMON_NAME=nitro-dev

# files
CA_KEY_PATH="/etc/nitrod/certs/ca.key"
CA_TRUST_PATH="/etc/nitrod/certs/ca.crt"
SERVER_PRIVATE_KEY_PATH="/etc/nitrod/certs/nitro.key"
SERVER_CSR_PATH="/etc/nitrod/certs/nitro.csr"

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

# generate a private key for the nitro server
if [ ! -f ${SERVER_PRIVATE_KEY_PATH} ]; then
  echo "INFO: creating the server private key..."
  openssl genrsa \
    -passout pass:1111 \
    -des3 \
    -out ${SERVER_PRIVATE_KEY_PATH} 4096
fi

# create the certificate signing request for the nitro server
if [ ! -f ${SERVER_CSR_PATH} ]; then
  echo "INFO: creating the server certificate signing request..."
  openssl req \
    -passin pass:1111 \
    -new \
    -key ${SERVER_PRIVATE_KEY_PATH} \
    -out ${SERVER_CSR_PATH} \
    -subj "/CN=${COMMON_NAME}"
fi
