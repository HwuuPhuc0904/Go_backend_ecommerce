#!/bin/bash

# This script deploys the Go backend e-commerce application.

# Set environment variables
export $(cat .env | xargs)

# Build the Docker images
echo "Building Docker images..."
docker-compose -f docker-compose.yml build

# Run database migrations
echo "Running database migrations..."
./scripts/migrate.sh

# Start the application
echo "Starting the application..."
docker-compose -f docker-compose.yml up -d

echo "Deployment completed successfully."