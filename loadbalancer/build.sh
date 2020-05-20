#!/bin/sh

# Build the docker image for the load balancer
docker-compose -f haproxy/build/docker-compose-build.yml build

# Build the docker image for the registrator
docker-compose -f registrator/build/docker-compose-build.yml build

# Build the docker image for the gonsul instance
docker-compose -f gonsul/build/docker-compose-build.yml build