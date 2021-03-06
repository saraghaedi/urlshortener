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

redis:
  master-address: 127.0.0.1:6379
  slave-address: 127.0.0.1:6379
  options:
    sentinel: false
    master-name: mymaster
    password: ""
    pool-size: 0
    min-idle-conns: 20
    dial-timeout: 5s
    read-timeout: 3s
    write-timeout: 3s
    pool-timeout: 4s
    idle-timeout: 5m
    max-retries: 5
    min-retry-backoff: 1s
    max-retry-backoff: 3s

nats:
  addresses:
    - nats://127.0.0.1:4222

monitoring:
  prometheus:
    enabled: true
    address: ":9001"
`
