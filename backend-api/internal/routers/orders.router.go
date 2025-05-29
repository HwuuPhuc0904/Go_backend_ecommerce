package routers

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/middleware"
    "github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(router *gin.RouterGroup) {
    orderController := controller.NewOrderController()
    
    // Nhóm các route liên quan đến đơn hàng (yêu cầu xác thực)
    ordersRoutes := router.Group("/orders")
    ordersRoutes.Use(middleware.AuthMiddleware())
    {
        // Tạo đơn hàng mới
        ordersRoutes.POST("", orderController.CreateOrder)
        
        // Lấy danh sách đơn hàng của người dùng hiện tại
        ordersRoutes.GET("", orderController.GetUserOrders)
        
        // Lấy chi tiết một đơn hàng
        ordersRoutes.GET("/:id", orderController.GetOrderByID)
        
        // Hủy đơn hàng
        ordersRoutes.PUT("/:id/cancel", orderController.CancelOrder)
        
        // Thanh toán đơn hàng
        ordersRoutes.POST("/:id/payment", orderController.ProcessPayment)
        
        // Xem thông tin vận chuyển
        ordersRoutes.GET("/:id/shipping", orderController.GetShippingInfo)
    }
}