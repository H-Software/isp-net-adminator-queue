# isp-net-adminator-queue
task manager for [ISP-Adminator](https://github.com/H-Software/isp-net-adminator)

# development

## env-vars for adminator
```
export REDIS_ADDR=127.0.0.1:16379

export MYSQL_SERVER=127.0.0.1

export MYSQL_USER=root
export MYSQL_PASSWD=isp-net-passwd

export POSTGRES_SERVER=127.0.0.1
export POSTGRES_USER=adminator
export POSTGRES_PASSWD=isp-net-passwd
export POSTGRES_DB=adminator.new
```

## env-vars for adminator - devcontainers
```
export REDIS_ADDR=host.docker.internal:16379

export MYSQL_SERVER=host.docker.internal

export MYSQL_USER=root
export MYSQL_PASSWD=isp-net-passwd

export POSTGRES_SERVER=host.docker.internal
export POSTGRES_USER=adminator
export POSTGRES_PASSWD=isp-net-passwd
export POSTGRES_DB=adminator.new
```

# links
- https://github.com/czhujer/h-platform-automation-cc-server/blob/master/main.go
## boilerplates
- https://github.com/ardanlabs/service
- https://github.com/allaboutapps/go-starter
## asynq
- https://github.com/hibiken/asynq/wiki/Getting-Started
- https://github.com/hibiken/asynq/wiki/Task-aggregation
## otel
- https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/examples/prometheus/main.go
- https://github.com/open-telemetry/opentelemetry-demo/blob/main/src/product-catalog/
  
# Author
Patrik Majer

# Licence
MIT
