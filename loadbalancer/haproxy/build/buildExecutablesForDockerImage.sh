#!/bin/sh

# build executable for compute-weights.go
GOOS=linux GOARCH=amd64 go build -o executables/compute-weights ../compute-weights.go