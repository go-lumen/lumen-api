#!/bin/bash
docker stop lumen-api
docker rm lumen-api
docker volume prune
docker rmi api-things-img:latest

docker build -t api-things-img .
docker run --name lumen-api -p 127.0.0.1:4000:4000 --link things-mongo:mongo -d api-things-img