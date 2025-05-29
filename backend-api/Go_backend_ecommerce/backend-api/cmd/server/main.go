package main

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/initialize"
    "log"
)

func main() {
    // Initialize global configurations
    if err := initialize.LoadConfig(); err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }

    // Initialize database connection
    if err := initialize.ConnectDatabase(); err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    // Run database migrations
    if err := initialize.MigrateDB(); err != nil {
        log.Fatalf("Error migrating database: %v", err)
    }

    // Set up routers and start the server
    router := initialize.InitRouters()
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}