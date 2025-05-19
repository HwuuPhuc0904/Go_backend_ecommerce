package service

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/repo"
	"time"
	"errors"

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





 
