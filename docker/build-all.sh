#!/bin/bash

VERSION=${1:-v0.0.1}

cd "$(dirname "$0")/.."

echo "Building Docker images with version: $VERSION"

echo "Building frontend image..."
docker build -f docker/Dockerfile.frontend -t online-editer-frontend:$VERSION .

echo "Building backend image..."
docker build -f docker/Dockerfile.backend -t online-editer-backend:$VERSION .

echo "Building sandbox image..."
docker build -f docker/Dockerfile.sandbox -t online-editer-sandbox:$VERSION .

echo "All images built successfully!"
echo ""
echo "Images:"
echo "  - online-editer-frontend:$VERSION"
echo "  - online-editer-backend:$VERSION"
echo "  - online-editer-sandbox:$VERSION"
