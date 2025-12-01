#!/bin/sh

set -ex  # Exit on error, print commands

echo "Building the app"

mkdir -p build

for arch in amd64 arm64; do
    echo "Building linux/$arch binary"
    CGO_ENABLED=1 GOOS=linux GOARCH=$arch \
        go build \
        -tags 'sqlite sqlite_unlock_notify' \
        -ldflags '-s -w' \
        -o ./build/app-$arch.out \
        ./cmd/api
    
    if [ $? -ne 0 ]; then
        echo "Error: Build for $arch failed."
        exit 1
    fi
done

echo "Build Complete"
ls -lh ./build/
