#!/bin/bash

# Copy environment variables
echo "📝 Copying environment variables..."
cp .env-* graphql/
cp .env-* websocket/

# Build and deploy GraphQL containers
cd graphql/
echo "🔨 Building and deploying GraphQL containers (development)..."
docker-compose -f docker-compose-develop.yml up -d --build

echo "🔨 Building and deploying GraphQL containers (production)..."
docker-compose -f docker-compose-production.yml up -d --build
cd ..

# Build and deploy Websocket containers
cd websocket/
echo "🔨 Building and deploying Websocket containers (development)..."
docker-compose -f docker-compose-develop.yml up -d --build

echo "🔨 Building and deploying Websocket containers (production)..."
docker-compose -f docker-compose-production.yml up -d --build
cd ..

echo "✅ Deployment completed successfully!"
