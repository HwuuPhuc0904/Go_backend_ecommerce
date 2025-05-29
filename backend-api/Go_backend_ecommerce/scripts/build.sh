#!/bin/bash

# Build the Docker images for the backend API, Nginx, and PostgreSQL

# Navigate to the backend-api directory and build the Docker image
cd backend-api
docker build -t backend-api .

# Navigate to the nginx directory and build the Docker image
cd ../nginx
docker build -t nginx-server .

# Navigate to the postgres directory and build the Docker image
cd ../postgres
docker build -t postgres-db .

# Return to the root directory
cd ..

# Optionally, you can run docker-compose to start all services
# Uncomment the following line if you want to start the services after building
# docker-compose up -d

echo "Docker images built successfully."