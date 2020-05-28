#!/bin/sh
set -e

gonsulArgs="$@"
consulAddr=localhost:8500

# Get the address of the consul server from the gonsul --consul-url= argument
while [ "$1" != "" ]; do
    # Need to escape first '-' character in grep pattern so that '--' is not interpeted as an arg to grep.
    if echo $1 | grep -q "\--consul-url="; then
        consulAddr="$(echo $1 | cut -d "=" -f 2)"
        break
    fi
    shift
done

# Ping consul server on the provided address: ex. pingConsulServer http://localhost:8500
pingConsulServer()
{
    local timeout=1 # seconds
    echo $(curl -s -o /dev/null -w %{http_code} --connect-timeout ${timeout} "${1}/v1/status/leader")
}

# Poll mock server running on the provided port: ex. pollConsulServer http://localhost:8500
pollConsulServer()
{
    local waitPeriod=2 # seconds
    local timeElapsed=0 # seconds

    while [ $(pingConsulServer $1) != "200" ];
    do
    echo "Waiting for consul server at ${1}, time elapsed: ${timeElapsed}s"
    timeElapsed=$((timeElapsed + $waitPeriod))
    sleep $waitPeriod;
    done

    echo "Consul server is UP at ${1}"
}

# Wait for consul to start
pollConsulServer ${consulAddr}

# Run gonsul with the provided commands.
gonsul ${gonsulArgs}