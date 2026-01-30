#!/bin/bash
# Build script for Cadence - automatically injects version info from git tags
# Usage: ./scripts/build.sh [output-file]

set -e

OUTPUT="${1:-cadence}"

# Get version from git tags, fall back to commit hash
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")

# Get short commit hash
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Get build time in UTC
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

# Build ldflags
LDFLAGS="-X github.com/codemeapixel/cadence/internal/version.Version=$VERSION -X github.com/codemeapixel/cadence/internal/version.GitCommit=$COMMIT -X github.com/codemeapixel/cadence/internal/version.BuildTime=$BUILD_TIME"

echo "Building Cadence..."
echo "  Version: $VERSION"
echo "  Commit:  $COMMIT"
echo "  Time:    $BUILD_TIME"
echo ""

go build -ldflags="$LDFLAGS" -o "$OUTPUT" ./cmd/cadence

echo "Build complete: $OUTPUT"
