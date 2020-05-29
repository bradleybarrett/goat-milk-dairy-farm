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

# Get the cluster leader from the consul server running at the provided address. 
# ex. getClusterLeader http://localhost:8500
getClusterLeader()
{
    local timeout=1 # seconds
    echo $(curl -s --connect-timeout ${timeout} "${1}/v1/status/leader" | tr -d '"')
}

# Wait for consul to start. 
# ex. waitForConsulServer http://localhost:8500
waitForConsulServer()
{
    local waitPeriod=2 # seconds
    local timeElapsed=0 # seconds
    
    # Wait for the consul server to start and elect a cluster leader.
    while [ -z $(getClusterLeader $1) ];
    do
        echo "Waiting for consul server at ${1}, time elapsed: ${timeElapsed}s"
        timeElapsed=$((timeElapsed + $waitPeriod))
        sleep $waitPeriod;
    done

    echo "Consul server is UP at ${1}"
}

# Wait for consul to start.
waitForConsulServer ${consulAddr}

# Run gonsul with the provided commands.
gonsul ${gonsulArgs}