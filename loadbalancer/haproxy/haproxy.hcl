template {
  source = "/tmp/haproxy.ctmpl"
  destination = "/etc/haproxy/haproxy.cfg"
  command = "./run/haproxy_reload.sh"
}