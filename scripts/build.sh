#!/bin/sh

echo "Building the app with xgo"

mkdir -p build

if ! command -v xgo &> /dev/null; then
    go install src.techknowlogick.com/xgo@latest
fi

# Build for both architectures with CGO support
xgo -out app \
    -targets 'linux/amd64,linux/arm64' \
    -dest ./build \
    -ldflags '-s -w' \
    ./cmd/api

# Rename outputs to match your naming convention
mv ./build/app-linux-amd64 ./build/app-amd64.out
mv ./build/app-linux-arm64 ./build/app-arm64.out

echo "Build Complete"
