FROM alpine
LABEL Description="Runs a gonsul instance to persist values from git to consul" Version="0.1"

# Copy the gonsul executable and edit the file permissions to be executable
COPY ./build/executables/gonsul /usr/bin/gonsul
RUN chmod +x /usr/bin/gonsul

# serves gonsul accessable at http://localhost:8000/v1/run when run with strategy=HOOK
# gonsul clones the config repo using git and ssh
# need to bind mount ssh keys from host to container and provide the ssh key file as a command line argument to gonsul

ENTRYPOINT ["./usr/bin/gonsul"]