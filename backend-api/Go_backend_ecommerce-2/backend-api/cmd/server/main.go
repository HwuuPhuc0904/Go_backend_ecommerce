package main

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/routers"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/migrations"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize the database and run migrations
    if err := migrations.MigrateDB(); err != nil {
        global.Logger.Fatal("Failed to migrate database", zap.Error(err))
    }

    // Set up the router
    r := gin.Default()
    routers.SetupRouter(r)

    // Start the server
    if err := r.Run(":8080"); err != nil {
        global.Logger.Fatal("Failed to start server", zap.Error(err))
    }
}