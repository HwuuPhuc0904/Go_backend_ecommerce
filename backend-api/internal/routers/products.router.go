package routers

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller"
    "github.com/gin-gonic/gin"
)

// RegisterProductRoutes đăng ký các routes liên quan đến sản phẩm
func RegisterProductRoutes(router *gin.RouterGroup) {
    productController := controller.NewProductController()

    // Routes công khai
    publicRoutes := router.Group("/products")
    {
        publicRoutes.GET("/:id", productController.GetProductByID)
		publicRoutes.POST("", productController.CreateProduct)
    }
}