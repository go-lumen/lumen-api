#!/bin/bash
openssl genrsa -out key.rsa 2048
openssl base64 -in key.rsa -out key64.rsa
docker run --name demo-plb-mongo -d mongo
docker build -t api-demo-plb-img .
docker run --name demo-plb-api -p 127.0.0.1:4000:4000 --link demo-plb-mongo:mongo -d api-demo-plb-img
