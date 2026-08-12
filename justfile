ui := "ui"

watch:
    #!/usr/bin/env sh
    trap 'kill 0' EXIT
    just watch-backend &
    just watch-frontend &
    wait

watch-backend:
    air -c .air.toml

watch-frontend:
    cd {{ui}} && bun run dev

build: build-frontend build-backend

build-backend:
    go build -tags embed_ui -o build/abyss .

build-frontend:
    cd {{ui}} && bun run build

mocks:
    mockery
    go mod tidy
