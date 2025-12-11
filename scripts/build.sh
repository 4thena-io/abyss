#!/bin/sh

# Build the frontend
cd frontend
bun run build
cd ..

# Install xgo if not already installed
go install src.techknowlogick.com/xgo@latest

# Cross-compile with xgo
xgo \
  --targets=linux/amd64,linux/arm64 \
  --dest=build \
  --out=abyss \
  ./cmd/api
