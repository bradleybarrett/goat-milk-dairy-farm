#!/bin/bash
serverPort=8091

port=$serverPort
version=0
host=localhost
consulHost=localhost

while [ "$1" != "" ]; do
    case $1 in
        -n | --name )
            shift
            name=$1
            ;;
        -p | --port )
            shift
            port=$1
            ;;
        -v | --version )
            shift  
            version=$1
            ;;
        -h | --host )
            shift  
            host=$1
            ;;
        -ch | --consul-host )
            shift  
            consulHost=$1
            ;;
        * )
            exit 1
    esac
    shift
done

docker run -d \
    -p $port:$serverPort \
    bbarrett/farmer:latest \
    --host.ip=$host \
    --host.port=$port \
    --consul.host=$consulHost \
    --version=$version