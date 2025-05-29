#!/bin/bash

# Start the development environment for the Go backend e-commerce project

# Build the Docker images
docker-compose -f ../docker-compose.dev.yml build

# Start the services
docker-compose -f ../docker-compose.dev.yml up
