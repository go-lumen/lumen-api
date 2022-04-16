#!/bin/bash
docker stop things-api
docker rm things-api
docker volume prune
docker rmi api-things-img:latest

docker build -t api-things-img .
docker run --name things-api -p 127.0.0.1:4000:4000 --link things-mongo:mongo -d api-things-img