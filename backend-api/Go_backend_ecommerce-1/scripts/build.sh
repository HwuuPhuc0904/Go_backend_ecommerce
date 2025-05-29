#!/bin/bash

# Build the Docker images for the backend API and other services

# Navigate to the backend-api directory
cd backend-api

# Build the backend API Docker image
docker build -t backend-api .

# Navigate to the deployments/docker/postgres directory and build the PostgreSQL image
cd ../deployments/docker/postgres
docker build -t postgres-db .

# Navigate to the deployments/docker/nginx directory and build the Nginx image
cd ../nginx
docker build -t nginx-server .

# Return to the root directory
cd ../../..

# Optionally, you can run docker-compose to start all services
# docker-compose up -d

echo "Docker images built successfully."