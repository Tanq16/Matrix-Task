#!/bin/bash

# Exit on any error
set -e

# Get the project root directory and create build directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$PROJECT_ROOT/build"

# Build the application
echo "Building application..."
go build -o "$PROJECT_ROOT/build/matrix" "$PROJECT_ROOT/cmd/main.go"

# Copy static assets and templates
echo "Copying assets..."
cp -r "$PROJECT_ROOT/static" "$PROJECT_ROOT/build/"
cp -r "$PROJECT_ROOT/internal/templates" "$PROJECT_ROOT/build/"

echo "Build complete! Binary is build/matrix"
