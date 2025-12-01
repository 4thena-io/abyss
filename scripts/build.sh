#!/bin/sh

echo "Building the app with xgo"

mkdir -p build

xgo -out app \
    -targets 'linux/amd64,linux/arm64' \
    -dest ./build \
		-tags 'sqlite sqlite_unlock_notify' \
    -ldflags '-s -w' \
    ./cmd/api

echo "Build Complete"
