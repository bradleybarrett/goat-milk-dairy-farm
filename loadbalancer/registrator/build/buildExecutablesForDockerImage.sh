#!/bin/sh

# build executable for registrator
GOOS=linux GOARCH=amd64 go build -o executables/registrator ../register.go ../service.go