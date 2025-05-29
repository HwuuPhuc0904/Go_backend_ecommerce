#!/bin/bash

# Navigate to the backend-api directory
cd backend-api

# Build the Docker images
docker-compose build

# Run the Docker containers
docker-compose up -d

# Run database migrations
docker-compose exec backend-api go run migrations/migration.go

# Optionally, you can run seed data if needed
# docker-compose exec backend-api go run migrations/seeddata.go

# Display the logs of the backend service
docker-compose logs -f backend-api