#!/bin/sh

# Build the docker image for the load balancer
# cd haproxy/build && ./buildExecutablesForDockerImage.sh && cd ../..
docker-compose -f haproxy/build/docker-compose-build.yml build


# Build the docker image for the registrator
#cd registrator/build && ./buildExecutablesForDockerImage.sh && cd ../..
docker-compose -f registrator/build/docker-compose-build.yml build

# Build the docker image for the gonsul instance
docker-compose -f gonsul/build/docker-compose-build.yml build