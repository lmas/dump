#!/bin/bash
echo "Generating an SSL private key to sign your certificate..."
openssl genrsa -des3 -out myssl.key 2048

echo "Generating a Certificate Signing Request..."
openssl req -new -key myssl.key -out myssl.csr

echo "Removing passphrase from key (for nginx)..."
cp myssl.key myssl.key.org
openssl rsa -in myssl.key.org -out myssl.key
rm myssl.key.org

echo "Generating certificate..."
openssl x509 -req -days 365 -in myssl.csr -signkey myssl.key -out myssl.crt

echo "Setting stricter permissions on key..."
chmod 400 myssl.key

echo "Done! You can now check your website's certs at: https://www.ssllabs.com/ssltest/"
