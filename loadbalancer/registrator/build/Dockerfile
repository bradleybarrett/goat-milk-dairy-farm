FROM alpine:3.11.2
LABEL Description="Runs a registrator which registers an haproxy instance with consul" Version="0.1"

ENV PATH /go/bin:$PATH
COPY ./build/executables/registrator /go/bin/registrator
RUN chmod +x /go/bin/registrator

ENTRYPOINT ["./go/bin/registrator"]