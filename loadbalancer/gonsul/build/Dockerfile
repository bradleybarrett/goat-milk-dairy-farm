FROM alpine
LABEL Description="Runs a gonsul instance to persist values from git to consul" Version="0.1"

# Install required packages
RUN apk add --no-cache su-exec

# Create a non-root user to run the gonsul process
RUN addgroup -S gonsul && adduser --home /home/gonsul -S gonsul -G gonsul

# Copy executables and modify permissions as needed.
# Gonsul executable retrieved from: https://github.com/miniclip/gonsul/releases
WORKDIR /usr/local/bin
COPY ./gonsul ./gonsul
COPY ./build/docker-entrypoint.sh ./docker-entrypoint.sh
RUN chmod 755 ./docker-entrypoint.sh 

# Gonsul clones the config repo using git (optionally over an ssh connection)
# Gonsul can use ssh keys bind-mounted to the container at start-up.
# An .ssh directory on the host should be bind-mounted to /tmp/.ssh
# The entrypoint script copies the read-only .ssh directory to the user home directory (/home/gonsul/.ssh) and modifies permissions as needed
# Note: gonsul is accessable at http://host:port/v1/run when run with strategy=HOOK
WORKDIR /home/gonsul
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["gonsul"]