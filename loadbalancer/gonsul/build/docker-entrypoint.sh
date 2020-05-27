#!/bin/sh
set -e

user=gonsul
home=/home/${user}

# Copy .ssh files from bind-mount to the user's home directory. 
# Permissions and ownership of the copy can be modified without propagating changes to the host.
cp -R /tmp/.ssh ${home}/.ssh

# Give the user ownership of the .ssh files in their home directory.
chown -R ${user}:${user} ${home}/.ssh

# If the first arg looks like a flag, assume we want to run gonsul server.
if [ "${1:0:1}" = '-' ]; then
    # Prefix the input command args with the gonsul command.
    set -- gonsul "$@"
fi

# Run gonsul in the container as a non-root user.
exec su-exec ${user}:${user} "$@"