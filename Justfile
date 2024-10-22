@_default:
    just -l

tests: worker_test deno_test

# Launches the woker in debug mode
worker:
    #!/usr/bin/env bash
    cd {{ justfile_dir() / "worker" }}
    go fmt > /dev/null 2>&1
    go run . -dev

worker_test:
    cd {{ justfile_dir() / "worker" }} && \
    go test

deno:
    cd {{ justfile_dir() / "backend" }} && \
    deno run format && \
    deno run dev

deno_test:
    cd {{ justfile_dir() / "backend" }} && \
    deno test --allow-net --allow-env --allow-read --fail-fast --filter Worker

# This is a special test script that runs to debug only
[doc]
deno_ws:
    cd {{ justfile_dir() / "backend" }} && \
    deno test --allow-net --allow-env --allow-read --fail-fast --filter socket.io

spam_ws:
    #!/usr/bin/env bash
    cd {{ justfile_dir() }}
    for i in $(seq 1 50); do
        # The $i is not very relevant since it's used in parallel
        ((just deno_ws > /dev/null 2>&1) || echo "Test #$i failed") &
    done

vue:
    cd {{ justfile_dir() / "frontend" }} && \
    pnpm run dev



# Launches everything with docker compose
run-all:
    #!/usr/bin/env bash
    password=$(cat .env | grep "REDIS_PASSWORD" | grep -oE '".*"' | sed s/\"//g)
    sed "s/^requirepass.*/requirepass $password/m" -i redis.conf
    # docker compose up --build --remove-orphans --force-recreate
    docker compose up --build --remove-orphans
    sed "s/^requirepass.*/requirepass HIDDEN/m" -i redis.conf

# Uses .env REDIS_PASSWORD as a default user password
[doc("Launches only Redis with docker compose for dev only")]
redis:
    #!/usr/bin/env bash
    password=$(cat .env | grep "REDIS_PASSWORD" | grep -oE '".*"' | sed s/\"//g)
    sed "s/^requirepass.*/requirepass $password/m" -i redis.conf
    docker compose -f docker-compose-dev.yml build
    docker compose -f docker-compose-dev.yml up -d
    sed "s/^requirepass.*/requirepass HIDDEN/m" -i redis.conf

# Alias of docker compose down
@stop:
    docker compose down --remove-orphans
