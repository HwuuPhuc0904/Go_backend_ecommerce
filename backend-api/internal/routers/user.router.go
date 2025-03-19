package routers

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller"
    "github.com/gin-gonic/gin"
)

// RegisterUserRoutes đăng ký các routes liên quan đến người dùng
func RegisterUserRoutes(router *gin.RouterGroup) {
    userController := controller.NewUserController()

    // Nhóm routes cho user
    userRoutes := router.Group("/users")
    {
        userRoutes.GET("/:id", userController.GetUserByID)
    }
}

// SetupRouter cài đặt toàn bộ router chính
func SetupRouter() *gin.Engine {
    r := gin.Default() // Đã bao gồm Logger và Recovery middleware
    
    // API version group
    v1 := r.Group("/api/v1")
    
    // Đăng ký các routes
    RegisterUserRoutes(v1)
    
    // Thêm route ping để kiểm tra server hoạt động
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    
    return r
}