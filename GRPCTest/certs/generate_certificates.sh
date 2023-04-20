#!/bin/bash

rm -f ca.crt ca.key server.crt server.csr server.key client*.key client*.csr client*.crt

# Certificate authority
openssl genrsa -out ca.key 2048

openssl req -new -x509 -days 3650 -key ca.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out ca.crt

# Server certificate
openssl req -newkey rsa:2048 -nodes -keyout server.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=localhost" -out server.csr

openssl x509 -req -extfile <(printf "subjectAltName=DNS:localhost") -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt


# Client certificates
for i in {1..2}; do
    openssl genrsa -out client$i.key 2048
    openssl req -new -key client$i.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=client$i" -out client$i.csr
    openssl x509 -req -days 3650 -in client$i.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client$i.crt
done
