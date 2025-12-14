#!/bin/bash

# Release script for template-cli
# Usage: ./scripts/release.sh <version>
# Example: ./scripts/release.sh 1.0.0

set -e

if [ -z "$1" ]; then
  echo "Error: Version required"
  echo "Usage: $0 <version>"
  echo "Example: $0 1.0.0"
  exit 1
fi

VERSION=$1
TAG="v${VERSION}"

echo "Building release ${VERSION}..."

# Create dist directory
mkdir -p dist

# Build for multiple platforms
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o dist/template-cli-linux-amd64

echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -o dist/template-cli-linux-arm64

echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o dist/template-cli-darwin-amd64

echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o dist/template-cli-darwin-arm64

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o dist/template-cli-windows-amd64.exe

echo "Building for Windows (arm64)..."
GOOS=windows GOARCH=arm64 go build -o dist/template-cli-windows-arm64.exe

# Create checksums
echo "Creating checksums..."
cd dist
sha256sum * > checksums.txt
cd ..

echo ""
echo "✓ Build complete! Binaries are in the dist/ directory"
echo ""
echo "To create a GitHub release:"
echo "  1. git tag ${TAG}"
echo "  2. git push origin ${TAG}"
echo ""
echo "The GitHub Actions workflow will automatically create a release with these binaries."

