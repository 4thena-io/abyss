#!/bin/sh

set -ex

echo "Building Abyss"

mkdir -p build

for arch in amd64 arm64; do
    echo "Building linux/$arch binary"
    CGO_ENABLED=0 GOOS=linux GOARCH=$arch \
        go build \
        -ldflags="-s -w -X 'main.Version=${VERSION}'" \
        -o ./build/app-$arch \
        ./cmd/api
done

echo "Build Complete"
