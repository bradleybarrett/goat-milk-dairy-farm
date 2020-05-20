template {
  source = "/var/haproxy/tmp/haproxy.ctmpl"
  destination = "/var/haproxy/etc/haproxy/haproxy.cfg"
  command = "./tmp/haproxy_reload.sh"
}