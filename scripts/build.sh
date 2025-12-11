#!/bin/sh

set -e  # Exit on error

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo "${RED}[ERROR]${NC} $1"
}

log_section() {
    echo ""
    echo "${BLUE}========================================${NC}"
    echo "${BLUE}  $1${NC}"
    echo "${BLUE}========================================${NC}"
}

# Start build process
log_section "Starting Abyss Build Process"

# Build the frontend
log_section "Building Frontend"

if [ ! -d "frontend" ]; then
    log_error "Frontend directory not found!"
    exit 1
fi

log_info "Changing to frontend directory..."
cd frontend

log_info "Installing bun dependencies..."
bun install

log_info "Running bun build..."
if bun run build; then
    log_success "Frontend build completed successfully"
else
    log_error "Frontend build failed"
    exit 1
fi

log_info "Listing frontend build output..."
ls -lah dist/

log_info "Returning to root directory..."
cd ..

# Install xgo
log_section "Installing xgo"

log_info "Installing xgo from src.techknowlogick.com/xgo@latest..."
if go install src.techknowlogick.com/xgo@latest; then
    log_success "xgo installed successfully"
else
    log_error "Failed to install xgo"
    exit 1
fi

# Verify xgo is in PATH
log_info "Verifying xgo installation..."
if command -v xgo >/dev/null 2>&1; then
    XGO_PATH=$(which xgo)
    log_success "xgo found at: $XGO_PATH"
else
    log_error "xgo not found in PATH"
    exit 1
fi

# Build binaries
log_section "Building Go Binaries with xgo"

log_info "Target platforms: linux/amd64, linux/arm64"
log_info "Output directory: build/"
log_info "Binary name: abyss"
log_info "Package: ./cmd/api"

log_info "Starting cross-compilation..."
if xgo \
  --targets=linux/amd64,linux/arm64 \
  --dest=build \
  --out=abyss \
  ./cmd/api; then
    log_success "Cross-compilation completed successfully"
else
    log_error "Cross-compilation failed"
    exit 1
fi

# Verify build outputs
log_section "Verifying Build Artifacts"

log_info "Listing build directory contents..."
ls -lah build/

if [ -f "build/abyss-linux-amd64" ]; then
    SIZE=$(du -h build/abyss-linux-amd64 | cut -f1)
    log_success "abyss-linux-amd64 created (Size: $SIZE)"
else
    log_error "abyss-linux-amd64 not found!"
    exit 1
fi

if [ -f "build/abyss-linux-arm64" ]; then
    SIZE=$(du -h build/abyss-linux-arm64 | cut -f1)
    log_success "abyss-linux-arm64 created (Size: $SIZE)"
else
    log_error "abyss-linux-arm64 not found!"
    exit 1
fi

# Final summary
log_section "Build Summary"
log_success "All builds completed successfully!"
log_info "Build artifacts:"
for file in build/*; do
    if [ -f "$file" ]; then
        SIZE=$(du -h "$file" | cut -f1)
        log_info "  - $(basename $file) ($SIZE)"
    fi
done

log_section "Build Process Complete"
