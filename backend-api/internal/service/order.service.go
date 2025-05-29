package service

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/repo"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	"time"
	"errors"
	"math/rand"
	"fmt"
	"go.uber.org/zap"
	"strconv"

)
type OrderService struct {
	orderRepo *repo.OrderRepo
	productRepo *repo.ProductRepo
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo: repo.NewOrderRepo(),
		productRepo: repo.NewProductRepo(),
	}
}

func generateOrderNumber() string {
    now := time.Now()
    dateStr := now.Format("20060102")
    rand.Seed(time.Now().UnixNano())
    randomNum := rand.Intn(1000000)
    return fmt.Sprintf("ORD-%s-%06d", dateStr, randomNum)
}

func (os *OrderService) CreateOrder(order *models.Order, item []models.OrderItem,shippingMethod, paymentMethod string) (*models.Order, error) {
	order.OrderNumber = generateOrderNumber()


	// set default values
	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()


	var subtotalAmount float64
	var processedItems []models.OrderItem

	for _, item := range item {
		product, err := os.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("not enough stock for product ID %d", item.ProductID)
		}

		itemTotal := product.Price * float64(item.Quantity)
		processedItem := models.OrderItem{
            ProductID:    product.ID,
            SKU:          product.SKU,
            ProductName:  product.Name,
            ProductImage: product.MainImage,
            Quantity:     item.Quantity,
            UnitPrice:    product.Price,
            Subtotal:     itemTotal,
            Discount:     0, // Có thể tính toán giảm giá nếu cần
            FinalPrice:   itemTotal,
            Attributes:   "", // Có thể thêm thông tin thuộc tính nếu cần
        }

		processedItems = append(processedItems, processedItem)
		subtotalAmount += itemTotal

	}

	var discountAmount float64
	if order.CouponCode != "" {
		// Giả sử bạn có một hàm để tính toán giảm giá từ mã giảm giá
		//discountAmount = os.calculateDiscount(order.CouponCode, subtotalAmount)
	}
	
	var shippingAmount float64 = 300000 // add shipping logic here

	var taxAmount = subtotalAmount * 0.1 // VAT 10%

	order.SubtotalAmount = subtotalAmount
    order.DiscountAmount = discountAmount
    order.ShippingAmount = shippingAmount
    order.TaxAmount = taxAmount
    order.TotalAmount = subtotalAmount - discountAmount + shippingAmount + taxAmount


	createdOrder, err := os.orderRepo.CreateOrder(order, processedItems, shippingMethod, paymentMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range processedItems {
        err = os.productRepo.UpdateProductStock(item.ProductID, -item.Quantity)
        if err != nil {
            global.Logger.Error("Lỗi khi cập nhật tồn kho sản phẩm", zap.Error(err))
        }
    }
	return createdOrder, nil	
}

// Get orders by user ID with pagination and status filter
func (os *OrderService) GetOrdersByUserID(userID uint, page, pageSize int, status string) ([]models.Order, int64, error) {
    return os.orderRepo.GetOrdersByUserID(userID, page, pageSize, status)
}

// Get order information by order ID
// and check if the user has permission to access it
func (os *OrderService) GetOrderByID(orderID, userID uint) (*models.Order, error) {
    order, err := os.orderRepo.GetOrderByID(orderID)
    if err != nil {
        return nil, err
    }
    
    if order.UserID != userID {
        return nil, errors.New("permission denied")
    }
    
    return order, nil
}

func (os *OrderService) CancelOrder(orderID, userID uint, cancelReason string) error {
    order, err := os.orderRepo.GetOrderByID(orderID)
    if err != nil {
        return err
    }
    
    // Kiểm tra xem đơn hàng có thuộc về người dùng không
    if order.UserID != userID {
        return errors.New("permission denied")
    }
    
    // Chỉ cho phép hủy đơn hàng ở trạng thái "pending" hoặc "processing"
    if order.Status != "pending" && order.Status != "processing" {
        return errors.New("order cannot be cancelled")
    }
    
    // Cập nhật trạng thái đơn hàng
    err = os.orderRepo.UpdateOrderStatus(orderID, "cancelled", cancelReason)
    if err != nil {
        return err
    }
    
    // Hoàn trả số lượng sản phẩm vào kho
    orderItems, err := os.orderRepo.GetOrderItems(orderID)
    if err != nil {
        return err
    }
    
    for _, item := range orderItems {
        err = os.productRepo.UpdateProductStock(item.ProductID, item.Quantity)
        if err != nil {
            global.Logger.Error("Lỗi khi cập nhật lại tồn kho sản phẩm", zap.Error(err))
            // Trong thực tế nên xử lý riêng việc cập nhật tồn kho thất bại
        }
    }
    
    return nil
}


func (os *OrderService) ProcessPayment(orderID, userID uint, payment *models.Payment) (*models.Payment, error) {
    order, err := os.orderRepo.GetOrderByID(orderID)
    if err != nil {
        return nil, err
    }
    
    // Kiểm tra xem đơn hàng có thuộc về người dùng không
    if order.UserID != userID {
        return nil, errors.New("permission denied")
    }
    
    // Kiểm tra xem đơn hàng có thể thanh toán không
    if order.Status != "pending" && order.Status != "processing" {
        return nil, errors.New("invalid order status for payment")
    }
    
    // Kiểm tra số tiền thanh toán
    if payment.Amount != order.TotalAmount {
        return nil, errors.New("payment amount does not match order total")
    }
    
    // Trong thực tế, ở đây sẽ gọi đến payment gateway để xử lý thanh toán
    // Mô phỏng xử lý thanh toán thành công
    payment.Status = "completed"
    payment.TransactionID = "TXN-" + strconv.FormatInt(time.Now().Unix(), 10)
    
    // Lưu thông tin thanh toán
    createdPayment, err := os.orderRepo.CreatePayment(payment)
    if err != nil {
        return nil, err
    }
    
    // Cập nhật đơn hàng với ID thanh toán
    err = os.orderRepo.UpdateOrderPayment(orderID, createdPayment.ID)
    if err != nil {
        return nil, err
    }
    
    // Cập nhật trạng thái đơn hàng sang "processing"
    err = os.orderRepo.UpdateOrderStatus(orderID, "processing", "")
    if err != nil {
        return nil, err
    }
    
    return createdPayment, nil
}

// GetShippingInfo lấy thông tin vận chuyển của đơn hàng
func (os *OrderService) GetShippingInfo(orderID, userID uint) (*models.Shipping, error) {
    // Kiểm tra xem đơn hàng có thuộc về người dùng không
    order, err := os.orderRepo.GetOrderByID(orderID)
    if err != nil {
        return nil, err
    }
    
    if order.UserID != userID {
        return nil, errors.New("permission denied")
    }
    
    // Lấy thông tin vận chuyển
    shipping, err := os.orderRepo.GetShippingByOrderID(orderID)
    if err != nil {
        return nil, err
    }
    
    return shipping, nil
}
 
