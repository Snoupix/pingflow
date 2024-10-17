worker:
    #!/usr/bin/env bash
    cd {{justfile_dir() / "worker"}}
    go fmt
    go run . -dev

@run-all:
    docker compose build
    docker compose up -d

@redis:
    docker compose -f docker-compose-dev.yml build
    docker compose -f docker-compose-dev.yml up -d

@shutdown:
    docker compose down
