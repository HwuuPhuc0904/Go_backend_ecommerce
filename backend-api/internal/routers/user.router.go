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

