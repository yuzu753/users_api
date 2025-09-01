#!/bin/bash

# Build Lambda deployment package only

set -e

echo "Building Lambda deployment package..."

# Clean up previous builds
rm -f bootstrap lambda.zip

# Build for Linux
echo "Building for Linux/AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap cmd/lambda/main.go

# Create deployment package
echo "Creating deployment package..."
zip lambda.zip bootstrap

echo "Lambda package created: lambda.zip"
ls -lh lambda.zip