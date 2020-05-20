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

RUN mkdir -p /var/haproxy && chmod 0 /var/haproxy
WORKDIR /var/haproxy

# Files used by consul-template to configure and run haproxy
COPY ./haproxy_reload.sh ./tmp/haproxy_reload.sh
COPY ./haproxy.hcl ./tmp/haproxy.hcl
COPY ./haproxy.ctmpl ./tmp/haproxy.ctmpl

# Files used by haproxy to store runtime info
RUN mkdir ./run && chmod 0 /run
RUN touch ./run/haproxy.sock
RUN touch ./run/haproxy.pid

# Create a user for consul-template, haproxy will be started in a separate process by this user/group
RUN addgroup -S consul-template && adduser -S consul-template -G consul-template

# Give consul-template (and ultimately haproxy) appropriate permission to the resource directory
# traverse directory, traverse its sub folders, rw- /run files, r-- /tmp files, r-x /tmp.sh files
RUN chown consul-template:consul-template -R . \
  && chmod 700 . \
  && chmod 700 $(find . -type d) \
  && chmod 600 $(find ./run/ -type f) \
  && chmod 400 $(find ./tmp/ -type f) \
  && chmod 500 $(find ./tmp/ -type f -name "*.sh")

USER consul-template

# Environment variables that can be provided to the running container.
# These variables are used in the template for the haproxy frontend and backend definitions.
# Required:
# - APP_NAME (service name in consul of the app to be load balanced )
# Optional:
# - HAPROXY_APP_TARGET_PORT (defaults to 8202)
# - HAPROXY_STATS_TARGET_PORT (defaults to 8404)

ENTRYPOINT ["consul-template"]