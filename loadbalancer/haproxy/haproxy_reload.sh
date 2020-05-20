#!/bin/sh

[ -n "${DEBUG}" ] && set -x

echo "$0: Running reload script"
#TODO: put this in an env variable, created in the dockerfile
haproxyHome=/var/haproxy

previousConfig=${haproxyHome}/etc/haproxy/haproxy.cfg.previous
newConfig=${haproxyHome}/etc/haproxy/haproxy.cfg

pidFile=${haproxyHome}/run/haproxy.pid
sockFile=${haproxyHome}/run/haproxy.sock

# first start
if [ ! -f ${previousConfig} ]; then
  echo "$0: First configuration"
  cp ${newConfig} ${previousConfig}
  haproxy -D -f ${newConfig}
  exit 0
fi

# trigger a reload if the configuration has changed
if ! $(cmp -s ${newConfig} ${previousConfig}); then
  if [ -S ${sockFile} ]; then
    echo "show servers state" | socat ${sockFile} - > ${haproxyHome}/haproxy.serverstates
  fi
  prevConfigPid=$(cat ${pidFile})
  echo "$0: Configuration has changed."
  echo "$0: Start a process with the new config and... "
  echo "$0: Send a soft kill signal to the process running the previous config (pid=${prevConfigPid})."
  cp ${newConfig} ${previousConfig}
  haproxy -D -f ${newConfig} -sf ${prevConfigPid}
fi
