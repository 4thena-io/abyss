#!/bin/sh

set -ex

echo "Building the app"

mkdir -p build

# --- Build linux/amd64 binary ---
echo "Building linux/amd64 binary"
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build \
    -tags 'sqlite sqlite_unlock_notify' \
    -ldflags '-s -w' \
    -o ./build/app-amd64.out \
    ./cmd/api

# --- Build linux/arm64 binary (Cross-Compilation) ---
echo "Building linux/arm64 binary"
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
    CC=aarch64-linux-gnu-gcc \
    go build \
    -tags 'sqlite sqlite_unlock_notify' \
    -ldflags '-s -w' \
    -o ./build/app-arm64.out \
    ./cmd/api

echo "Build Complete"
ls -lh ./build/
