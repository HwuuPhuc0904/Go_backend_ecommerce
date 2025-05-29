package controller

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	"go.uber.org/zap"
	"time"
	"strconv"
)


type OrderController struct {
	orderService *service.OrderService
}

func NewOrderController() *OrderController {
	return &OrderController{
		orderService: service.NewOrderService(),
	}
}

func (oc *OrderController) CreateOrder(c * gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}
	
	var orderRequest struct {
		OrderItems []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required,min=1"`
		} `json:"order_items" binding:"required,dive"`
		ShippingAddressID  uint   `json:"shipping_address_id" binding:"required"`
		BillingAddressID   uint   `json:"billing_address_id" binding:"required"`
		ShippingMethod     string `json:"shipping_method" binding:"required"`
		PaymentMethod      string `json:"payment_method" binding:"required"` 
		RequireSignature   bool   `json:"require_signature"`
		ShippingNotes      string `json:"shipping_notes"`
		CouponCode         string `json:"coupon_code"`
		Notes              string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		global.Logger.Error("Failed to bind JSON products", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error should by json ": err.Error()})
		return
	}
	now := time.Now()
	order := &models.Order{
		UserID:            userID.(uint),
		OrderDate:         now,
		Status:            "pending",
		ShippingAddressID: orderRequest.ShippingAddressID,
		BillingAddressID:  orderRequest.BillingAddressID,
		CouponCode:        orderRequest.CouponCode,
		Notes:             orderRequest.Notes,
	}

	var orderItems []models.OrderItem
	for _, item := range orderRequest.OrderItems {
		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	createdOrder, err := oc.orderService.CreateOrder(order, orderItems, orderRequest.ShippingMethod, orderRequest.PaymentMethod)
	if err != nil {
		global.Logger.Error("Failed to create order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": createdOrder,
		"message": "Order created successfully",
	})

}

func (oc *OrderController) GetUserOrders(c * gin.Context){
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	orders, total, err := oc.orderService.GetOrdersByUserID(userID.(uint), page, pageSize, status)
	if err != nil {
		global.Logger.Error("Failed to get user orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"pagination": gin.H{
			"total": total,
			"page": page,
			"page_size": pageSize,
		},
	})
}

func (oc *OrderController) GetOrderByID(c * gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		global.Logger.Error("User ID not found")
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		global.Logger.Error("Invalid order ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := oc.orderService.GetOrderByID(uint(orderID), userID.(uint))
	if err != nil {
		global.Logger.Error("Failed to get order by ID", zap.Error(err))
	
		if order == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order by ID"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (oc *OrderController) CancelOrder(c *gin.Context) {
	userID, exists := c.Get("userID") 
	if !exists {
		global.Logger.Error("User ID not found")
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}

	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		global.Logger.Error("Invalid order ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var cancelRequest struct {
		cancelReason string `json:"cancel_reason"`
	}

	if err := c.ShouldBindJSON(&cancelRequest); err != nil {
		global.Logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
		return
	}

	err = oc.orderService.CancelOrder(uint(orderID), userID.(uint), cancelRequest.cancelReason)
	if err != nil {
		global.Logger.Error("Failed to update order status", zap.Error(err), zap.Int64("orderID", orderID))		
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		if err.Error() == "order cannot be canceled" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be canceled in this status"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancled order "})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Canceled order successfully"})
}


func (oc *OrderController) ProcessPayment(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Người dùng chưa được xác thực"})
        return
    }

    orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID đơn hàng không hợp lệ"})
        return
    }

    var paymentRequest struct {
        PaymentMethod   string  `json:"payment_method" binding:"required"`
        PaymentProvider string  `json:"payment_provider"`
        Amount          float64 `json:"amount" binding:"required"`
        // Các thông tin thẻ hoặc thông tin thanh toán khác
        CardNumber      string  `json:"card_number"`
        CardHolderName  string  `json:"card_holder_name"`
        ExpiryMonth     string  `json:"expiry_month"`
        ExpiryYear      string  `json:"expiry_year"`
        CVV             string  `json:"cvv"`
    }

    if err := c.ShouldBindJSON(&paymentRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu thanh toán không hợp lệ"})
        return
    }

    // Tạo thông tin thanh toán
    payment := &models.Payment{
        OrderID:         uint(orderID),
        PaymentMethod:   paymentRequest.PaymentMethod,
        PaymentProvider: paymentRequest.PaymentProvider,
        Amount:          paymentRequest.Amount,
        Status:          "pending",
        PaymentDate:     time.Now(),
        // Lưu chi tiết thanh toán dưới dạng JSON
        PaymentDetails:  "",  // Trong thực tế, cần mã hóa thông tin thẻ thành JSON
    }

    result, err := oc.orderService.ProcessPayment(uint(orderID), userID.(uint), payment)
    if err != nil {
        global.Logger.Error("Failed to process payment", 
            zap.Uint64("orderID", orderID), 
            zap.Error(err))
            
        if err.Error() == "record not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Not found order"})
            return
        } else if err.Error() == "invalid order status for payment" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Order not in a valid status for payment"})
            return
        }
        
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xử lý thanh toán: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Payment processed successfully",
        "payment": result,
    })
}

func (oc *OrderController) GetShippingInfo(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
        return
    }

    shippingInfo, err := oc.orderService.GetShippingInfo(uint(orderID), userID.(uint))
    if err != nil {
        global.Logger.Error("Failed to get shipping info", 
            zap.Uint64("orderID", orderID), 
            zap.Error(err))
            
        if err.Error() == "record not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Not found order or shipping info"})
            return
        }
        
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Not found order or shipping info"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "shipping_info": shippingInfo,
    })
}


	

