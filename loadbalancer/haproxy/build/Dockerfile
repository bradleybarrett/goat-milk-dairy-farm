FROM haproxy:2.1-alpine
LABEL Description="Runs consul template and haproxy based on the generated configuration" Version="0.1"

ENV CONSUL_TEMPLATE_VERSION=0.22.0

# Update wget to get support for SSL
RUN apk --update add wget

# Download consul-template
RUN wget --no-check-certificate https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_VERSION}/consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.tgz \
  && tar xfz consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.tgz \
  && mv consul-template /usr/bin/consul-template \
  && rm consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.tgz

RUN touch /run/haproxy.sock \
    && chmod 777 /run/haproxy.sock

RUN touch /run/haproxy.pid \
    && chmod 777 /run/haproxy.pid 

COPY ./haproxy_reload.sh /run/haproxy_reload.sh
RUN chmod +x /run/haproxy_reload.sh

COPY ./haproxy.hcl /tmp/haproxy.hcl
COPY ./haproxy.ctmpl /tmp/haproxy.ctmpl

ENV PATH /go/bin:$PATH
COPY ./build/executables/compute-weights /go/bin/compute-weights
RUN chmod +x /go/bin/compute-weights

# Environment variables that can be provided to the running container.
# These variables are used in the template for the haproxy frontend and backend definitions.
# Required:
# - APP_NAME (service name in consul of the app to be load balanced )
# Optional:
# - HAPROXY_APP_TARGET_PORT (defaults to 8202)
# - HAPROXY_STATS_TARGET_PORT (defaults to 8404)

ENTRYPOINT ["consul-template"]