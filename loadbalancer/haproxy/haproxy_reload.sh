#!/bin/sh

[ -n "${DEBUG}" ] && set -x

echo "$0: Running reload script"

previousConfig=/etc/haproxy/haproxy.cfg.previous
newConfig=/etc/haproxy/haproxy.cfg

# first start
if [ ! -f ${previousConfig} ]; then
  echo "$0: First configuration"
  cp ${newConfig} ${previousConfig}
  haproxy -D -f ${newConfig}
  exit 0
fi

# trigger a reload if the configuration has changed
if ! $(cmp -s ${newConfig} ${previousConfig}); then
  if [ -S /run/haproxy.sock ]; then
    echo "show servers state" | socat /run/haproxy.sock - > /haproxy.serverstates
  fi
  prevConfigPid=$(cat /run/haproxy.pid)
  echo "$0: Configuration has changed."
  echo "$0: Start a process with the new config and... "
  echo "$0: Send a soft kill signal to the process running the previous config (pid=${prevConfigPid})."
  cp ${newConfig} ${previousConfig}
  haproxy -D -f ${newConfig} -sf ${prevConfigPid}
fi