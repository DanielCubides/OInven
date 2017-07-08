#!/usr/bin/env sh

mvn -f pom.xml package

docker-compose build
docker-compose up