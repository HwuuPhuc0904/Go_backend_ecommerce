package routers

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/middleware"
    "github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
    userController := controller.NewUserController()

    // Routes công khai không yêu cầu xác thực
    publicRoutes := router.Group("/auth")
    {
        publicRoutes.POST("/register", userController.RegisterUser)
        publicRoutes.POST("/login", userController.Login)
    }

    // Routes yêu cầu xác thực
    authenticatedRoutes := router.Group("/users")
    authenticatedRoutes.Use(middleware.AuthMiddleware())
    {
        // Quản lý thông tin người dùng
        authenticatedRoutes.GET("/profile", userController.GetProfile)
        authenticatedRoutes.PUT("/profile", userController.UpdateProfileUser)
        authenticatedRoutes.PUT("/change-password", userController.ChangePassword)

        // Quản lý danh sách người dùng (chỉ dành cho admin)
        adminRoutes := authenticatedRoutes.Group("admin")
        adminRoutes.Use(middleware.AdminMiddleware())
        {
            adminRoutes.GET("", userController.GetAllUser)
            adminRoutes.GET("/:id", userController.GetUserByID)
            adminRoutes.PUT("/:id", userController.UpdateUser)
            adminRoutes.DELETE("/:id", userController.DeleteUser)
        }
    }
}

