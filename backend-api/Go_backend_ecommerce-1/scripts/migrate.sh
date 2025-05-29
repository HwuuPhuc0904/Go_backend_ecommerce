#!/bin/bash

# Navigate to the backend-api directory
cd backend-api

# Run database migrations
go run migrations/migration.go

# Optionally, you can add more migration commands here if needed
# For example, if you have seed data to run:
# go run migrations/seeddata.go

echo "Database migrations completed."