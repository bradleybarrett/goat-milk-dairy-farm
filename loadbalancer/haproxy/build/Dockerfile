# Build compute-weights executable with GO in throwaway container
FROM golang:alpine as gobuilder

RUN mkdir /build
WORKDIR /build 
COPY ./compute-weights.go .
RUN go build -o compute-weights compute-weights.go

# Build HA-Proxy image with consul-template and compute-weights executable
FROM haproxy:2.1-alpine
LABEL Description="Runs consul template and haproxy based on the generated configuration" Version="0.1"

# Update wget to get support for SSL
RUN apk --update add wget

# Download consul-template
ENV CONSUL_TEMPLATE_VERSION=0.22.0
RUN wget --no-check-certificate https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_VERSION}/consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.tgz \
  && tar xfz consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.tgz \
  && mv consul-template /usr/bin/consul-template \
  && rm consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.tgz

# Add compute-weight executable from build stage
COPY --from=gobuilder /build/compute-weights /go/bin/compute-weights
# Add the /go/bin directory to the path so consul-template can find the executable
ENV PATH /go/bin:$PATH

COPY ./haproxy_reload.sh /run/haproxy_reload.sh
RUN chmod +x /run/haproxy_reload.sh

RUN touch /run/haproxy.sock
RUN touch /run/haproxy.pid
COPY ./haproxy.hcl /tmp/haproxy.hcl
COPY ./haproxy.ctmpl /tmp/haproxy.ctmpl

# Environment variables that can be provided to the running container.
# These variables are used in the template for the haproxy frontend and backend definitions.
# Required:
# - APP_NAME (service name in consul of the app to be load balanced )
# Optional:
# - HAPROXY_APP_TARGET_PORT (defaults to 8202)
# - HAPROXY_STATS_TARGET_PORT (defaults to 8404)

ENTRYPOINT ["consul-template"]