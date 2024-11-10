#!/bin/bash

# Exit on any error
set -e

# Get the project root directory and create build directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$PROJECT_ROOT/build"

# Build the application
go build -o "$PROJECT_ROOT/build/matrix" "$PROJECT_ROOT/cmd/main.go"

echo "Build complete! Binary is build/matrix"
