#!/bin/bash
openssl req \
    -x509 -newkey rsa:4096 -sha256 -nodes \
    -subj '/CN=localhost' \
    -addext 'subjectAltName=DNS:localhost' \
    -keyout selfsigned.key \
    -out selfsigned.pem \
    -days 3650
