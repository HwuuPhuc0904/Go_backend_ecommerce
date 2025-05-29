package repo

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
    "go.uber.org/zap"
	"fmt"
	"time"

)

type OrderRepo struct{}

func NewOrderRepo() *OrderRepo {
    return &OrderRepo{}
}

// CreateOrder tạo đơn hàng mới và các mục trong đơn hàng
func (or *OrderRepo) CreateOrder(order *models.Order, orderItems []models.OrderItem, shippingMethod, paymentMethod string) (*models.Order, error) {
    // Bắt đầu transaction
    tx := global.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
	
	shipping := &models.Shipping{
        OrderID: order.ID,
        Status: "pending", 
        EstimatedDelivery: nil,
        ShippingCost: order.ShippingAmount,
        ShippingMethod: shippingMethod,
    }
	payment := &models.Payment{
        OrderID: order.ID,
        Amount: order.TotalAmount,
        Status: "pending",
        PaymentMethod: paymentMethod,
        PaymentDate: time.Now(),         
    }
    
    if err := tx.Create(payment).Error; err != nil {
        tx.Rollback()
        global.Logger.Error("Lỗi khi tạo bản ghi payment", zap.Error(err))
        return nil, fmt.Errorf("lỗi khi tạo bản ghi payment: %w", err)
    }

	if err := tx.Create(shipping).Error; err != nil {
        tx.Rollback()
        global.Logger.Error("Lỗi khi tạo bản ghi shipping", zap.Error(err))
        return nil, fmt.Errorf("lỗi khi tạo bản ghi shipping: %w", err)
    }

	order.ShippingID = shipping.ID
 	order.PaymentID = payment.ID

    // Tạo đơn hàng
    if err := tx.Create(order).Error; err != nil {
        tx.Rollback()
        global.Logger.Error("Lỗi khi tạo đơn hàng", zap.Error(err))
        return nil, err
    }

    // Thêm các mục vào đơn hàng
    for i := range orderItems {
        orderItems[i].OrderID = order.ID
        if err := tx.Create(&orderItems[i]).Error; err != nil {
            tx.Rollback()
            global.Logger.Error("Lỗi khi tạo mục đơn hàng", zap.Error(err))
            return nil, err
        }
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        global.Logger.Error("Lỗi khi commit transaction", zap.Error(err))
        return nil, err
    }

    // Lấy đơn hàng đầy đủ với các mục
    var createdOrder models.Order
    if err := global.DB.Preload("OrderItems").First(&createdOrder, order.ID).Error; err != nil {
        global.Logger.Error("Lỗi khi lấy đơn hàng đã tạo", zap.Error(err))
        return nil, err
    }

    return &createdOrder, nil
}

// GetOrdersByUserID lấy danh sách đơn hàng của người dùng với phân trang
func (or *OrderRepo) GetOrdersByUserID(userID uint, page, pageSize int, status string) ([]models.Order, int64, error) {
    var orders []models.Order
    var total int64
    offset := (page - 1) * pageSize

    // Xây dựng query
    query := global.DB.Model(&models.Order{}).Where("user_id = ?", userID)
    
    // Áp dụng bộ lọc trạng thái nếu có
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Đếm tổng số bản ghi
    if err := query.Count(&total).Error; err != nil {
        global.Logger.Error("Lỗi khi đếm đơn hàng", zap.Error(err))
        return nil, 0, err
    }
    
    // Lấy dữ liệu với phân trang
    if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").
        Preload("OrderItems").Find(&orders).Error; err != nil {
        global.Logger.Error("Lỗi khi lấy danh sách đơn hàng", zap.Error(err))
        return nil, 0, err
    }

    return orders, total, nil
}

// GetOrderByID lấy chi tiết đơn hàng theo ID
func (or *OrderRepo) GetOrderByID(orderID uint) (*models.Order, error) {
    var order models.Order
    err := global.DB.Preload("OrderItems.Product").
        Preload("Payment").
        Preload("Shipping").
        Preload("BillingAddress").
        Preload("ShippingAddress").
        First(&order, orderID).Error
    if err != nil {
        global.Logger.Error("Lỗi khi lấy đơn hàng theo ID", zap.Uint("orderID", orderID), zap.Error(err))
        return nil, err
    }
    return &order, nil
}

// GetUserByID lấy thông tin người dùng theo ID
func (or *OrderRepo) GetUserByID(userID uint) (*models.User, error) {
    var user models.User
    if err := global.DB.First(&user, userID).Error; err != nil {
        global.Logger.Error("Lỗi khi lấy thông tin người dùng", zap.Uint("userID", userID), zap.Error(err))
        return nil, err
    }
    return &user, nil
}

// UpdateOrderStatus cập nhật trạng thái đơn hàng
func (or *OrderRepo) UpdateOrderStatus(orderID uint, status, cancelReason string) error {
    updateFields := map[string]interface{}{
        "status": status,
    }
    
    // Cập nhật cả ghi chú nếu có
    if cancelReason != "" {
        updateFields["cancel_reason"] = cancelReason
    }

    err := global.DB.Model(&models.Order{}).Where("id = ?", orderID).
        Updates(updateFields).Error
    if err != nil {
        global.Logger.Error("Lỗi khi cập nhật trạng thái đơn hàng", 
            zap.Uint("orderID", orderID), 
            zap.String("status", status), 
            zap.String("cancelReason", cancelReason),
            zap.Error(err))
    }
     global.Logger.Info("Update order status successfully", 
            zap.Uint("orderID", orderID), 
            zap.String("status", status), 
            zap.String("cancelReason", cancelReason),
            )

    return nil
}

// GetOrderItems lấy danh sách các mục trong đơn hàng
func (or *OrderRepo) GetOrderItems(orderID uint) ([]models.OrderItem, error) {
    var orderItems []models.OrderItem
    err := global.DB.Where("order_id = ?", orderID).Find(&orderItems).Error
    if err != nil {
        global.Logger.Error("Lỗi khi lấy các mục trong đơn hàng", zap.Uint("orderID", orderID), zap.Error(err))
        return nil, err
    }
    return orderItems, nil
}

// CreatePayment tạo mới thanh toán
func (or *OrderRepo) CreatePayment(payment *models.Payment) (*models.Payment, error) {
    err := global.DB.Create(payment).Error
    if err != nil {
        global.Logger.Error("Lỗi khi tạo thanh toán", zap.Error(err))
        return nil, err
    }
    return payment, nil
}

// UpdateOrderPayment cập nhật thông tin thanh toán cho đơn hàng
func (or *OrderRepo) UpdateOrderPayment(orderID, paymentID uint) error {
    err := global.DB.Model(&models.Order{}).Where("id = ?", orderID).
        Update("payment_id", paymentID).Error
    if err != nil {
        global.Logger.Error("Lỗi khi cập nhật ID thanh toán cho đơn hàng", 
            zap.Uint("orderID", orderID), 
            zap.Uint("paymentID", paymentID), 
            zap.Error(err))
    }
    return err
}

// GetShippingByOrderID lấy thông tin vận chuyển của đơn hàng
func (or *OrderRepo) GetShippingByOrderID(orderID uint) (*models.Shipping, error) {
    var order models.Order
    err := global.DB.Select("shipping_id").First(&order, orderID).Error
    if err != nil {
        return nil, err
    }
    
    // Nếu đơn hàng chưa có thông tin vận chuyển
    if order.ShippingID == 0 {
        return nil, nil
    }
    
    var shipping models.Shipping
    err = global.DB.First(&shipping, order.ShippingID).Error
    if err != nil {
        global.Logger.Error("Lỗi khi lấy thông tin vận chuyển", 
            zap.Uint("orderID", orderID), 
            zap.Uint("shippingID", order.ShippingID), 
            zap.Error(err))
        return nil, err
    }
    
    return &shipping, nil
}