frontend := "frontend"

watch:
    #!/usr/bin/env sh
    trap 'kill 0' EXIT
    just watch-backend &
    just watch-frontend &
    wait

watch-backend:
    air -c .air.toml

watch-frontend:
    cd {{frontend}} && bun run dev

build: build-frontend build-backend

build-backend:
    go build -o build/abyss .

build-frontend:
    cd {{frontend}} && bun run build

mocks:
    mockery
    go mod tidy
