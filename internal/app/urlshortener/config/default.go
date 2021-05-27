package config

// nolint:lll
const defaultConfig = `
logger:
  access:
    enabled: false
    path: "/var/log/access.log"
    format: "${remote_ip} - - [${time_rfc3339}] \"${method} ${uri} HTTP/1.1\" ${status} \
    ${bytes_out} ${bytes_in} ${latency} \"${referer}\" \"${user_agent}\"\n"
    max-size: 1024
    max-backups: 7
    max-age: 7
  app:
    level: info
    path: "/var/log/app.log"
    max-size: 1024
    max-backups: 7
    max-age: 7
    stdout: true
    
server:
  address: :8080
  read-timeout: 20s
  write-timeout: 20s
  graceful-timeout: 5s

database:
  driver: postgres
  master-conn-string: postgresql://urldb:secret@127.0.0.1:5432/urldb?sslmode=disable&connect_timeout=30
  slave-conn-string: postgresql://urldb:secret@127.0.0.1:5432/urldb?sslmode=disable&connect_timeout=30

monitoring:
  prometheus:
    enabled: true
    address: ":9001"
`
