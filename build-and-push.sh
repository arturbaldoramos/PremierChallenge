#!/bin/bash

# Build and Push Script for Premiere Challenge
# Usage: ./build-and-push.sh [dockerhub-username]

set -e

DOCKERHUB_USERNAME=${1:-"blamet"}
VERSION=${2:-"latest"}

echo "🏗️  Building Premiere Challenge Docker Images..."
echo "DockerHub Username: $DOCKERHUB_USERNAME"
echo "Version: $VERSION"

# Build Backend
echo "📦 Building Backend..."
cd Backend
docker build -t "$DOCKERHUB_USERNAME/premiere-challenge-backend:$VERSION" .
docker tag "$DOCKERHUB_USERNAME/premiere-challenge-backend:$VERSION" "$DOCKERHUB_USERNAME/premiere-challenge-backend:latest"
cd ..

# Build Frontend
echo "🌐 Building Frontend..."
cd Frontend
docker build -t "$DOCKERHUB_USERNAME/premiere-challenge-frontend:$VERSION" .
docker tag "$DOCKERHUB_USERNAME/premiere-challenge-frontend:$VERSION" "$DOCKERHUB_USERNAME/premiere-challenge-frontend:latest"
cd ..

echo "✅ Build completed successfully!"

# Ask for push confirmation
read -p "🚀 Push images to DockerHub? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🔑 Logging into DockerHub..."
    docker login

    echo "⬆️  Pushing Backend..."
    docker push "$DOCKERHUB_USERNAME/premiere-challenge-backend:$VERSION"
    docker push "$DOCKERHUB_USERNAME/premiere-challenge-backend:latest"

    echo "⬆️  Pushing Frontend..."
    docker push "$DOCKERHUB_USERNAME/premiere-challenge-frontend:$VERSION"
    docker push "$DOCKERHUB_USERNAME/premiere-challenge-frontend:latest"

    echo "🎉 Images pushed successfully!"
    echo ""
    echo "📋 Images available:"
    echo "   Backend:  $DOCKERHUB_USERNAME/premiere-challenge-backend:$VERSION"
    echo "   Frontend: $DOCKERHUB_USERNAME/premiere-challenge-frontend:$VERSION"

    echo ""
    echo "🔄 To update your docker-compose, the images are:"
    echo "   blamet/premiere-challenge-backend:latest"
    echo "   blamet/premiere-challenge-frontend:latest"
else
    echo "⏭️  Skipping push to DockerHub"
fi

echo ""
echo "🐳 To test locally:"
echo "   docker-compose up -d"