@_default:
    just -l

# Launches the woker in debug mode
worker:
    #!/usr/bin/env bash
    cd {{justfile_dir() / "worker"}}
    go fmt > /dev/null 2>&1
    go run . -dev

# Launches everything with docker compose
@run-all:
    docker compose build
    docker compose up -d

# Uses .env REDIS_PASSWORD as a default user password
[doc("Launches only Redis with docker compose")]
redis:
    #!/usr/bin/env bash
    password=$(cat .env | grep "REDIS_PASSWORD" | grep -oE '".*"' | sed s/\"//g)
    sed "s/^requirepass.*/requirepass $password/m" -i redis.conf
    docker compose -f docker-compose-dev.yml build
    docker compose -f docker-compose-dev.yml up -d
    sed "s/^requirepass.*/requirepass HIDDEN/m" -i redis.conf

# Alias of docker compose down
@stop:
    docker compose down
