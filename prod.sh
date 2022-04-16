#!/bin/bash
openssl genrsa -out key.rsa 2048
openssl base64 -in key.rsa -out key64.rsa
cat key64.rsa | tr -d '\n'
docker run --name things-mongo -d mongo
docker build -t api-things-img .
docker run --name things-api -p 127.0.0.1:4000:4000 --link things-mongo:mongo -d api-things-img