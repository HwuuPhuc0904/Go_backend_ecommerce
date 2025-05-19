package service

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    model "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/repo"
    "errors"
    "strings"
    "go.uber.org/zap"
)


type ProductService struct {
    productRepo *repo.ProductRepo
}

func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repo.NewProductRepo(),
	}
}



// Get Product by Category
func(ps * ProductService) GetProductsByCategory(category string, limit int, page int , shortBy string, orderBy string ) ([]model.Product,int64, error){
	global.Logger.Info("GetProductByCategory", zap.String("Category", category))
	

	offset := (page - 1) * limit

	// validate sortby and orderby
	allowedSortFileds := map[string]bool{
		"id": true,
		"name": true,
		"price": true,
		"created_at": true,	
	}

	if !allowedSortFileds[shortBy] {
		shortBy = "id"
	}
	
	allowedOrderBy := map[string]bool {
		"asc": true,
		"desc": true,
	}

	if !allowedOrderBy[orderBy] {
		orderBy = "asc"
	}

	return ps.productRepo.GetProductByCategory(category, limit, offset, shortBy, orderBy)

}


func (ps * ProductService) GetAllProducts(limit int, page int, sortBy string, orderBy string ) ([]model.Product,int64, error) {
	
	offset := (page - 1) * limit
	// validate sortby and orderby 
	allowedSortFileds := map[string]bool{
		"id" : true, 
		"name": true,
		"price": true,
		"created_at": true,
	}

	if !allowedSortFileds[sortBy] {
		sortBy = "id"
	}

	// validate orderby
	allowedOrderBy := map[string]bool{
		"asc": true,
		"desc": true,
	}
	if !allowedOrderBy[orderBy] {
		orderBy = "asc"
	}

	products, total, err := ps.productRepo.GetAllProducts(limit, offset, sortBy, orderBy)

	if err != nil {	
		global.Logger.Error("Error to get products", zap.Error(err))
		return nil, 0, err
	}

	return products, total, nil
}

func (ps *ProductService) GetProductByID(id uint) (*model.Product, error) {
	global.Logger.Info("GetProductByID", zap.Uint("id", id))
	return ps.productRepo.GetProductByID(id)
}

func (ps *ProductService) CreateProduct(product *model.Product, userID uint) error {
	product.UserID = userID

	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)

	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price < 0 {
		return errors.New("price must be greater than 0")
	}
	global.Logger.Info("CreateProduct", zap.String("Name", product.Name), zap.Float64("Price", product.Price))	
    
	return ps.productRepo.CreateProduct(product)
}

func (ps *ProductService) UpdateProduct(product *model.Product) error {
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)

	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price < 0 {
		return errors.New("price must be greater than 0")
	}
	global.Logger.Info("UpdateProduct", zap.String("Name", product.Name), zap.Float64("Price", product.Price))	
	
	return ps.productRepo.UpdateProduct(product)
}

func (ps * ProductService) DeleteProductByID(id uint) error {
	existingProduct, err := ps.GetProductByID(id)
	if err != nil {
		return err
	}
	
	if existingProduct == nil {
		return errors.New("product not found")
	}	
	global.Logger.Info("DeleteProductByID", zap.Uint("id", id))

	return ps.productRepo.DeleteProductByID(id)
}

func (ps *ProductService) SearchProducts(keyword string) ([]model.Product, error) {
    global.Logger.Info("Searching products", zap.String("keyword", keyword))
    return ps.productRepo.SearchProduct(keyword)
}

func (ps *ProductService) UpdateStockProduct(id uint, quantity int) error {
	global.Logger.Info("UpdateStockProduct", zap.Uint("id", id), zap.Int("quantity", quantity))
	return ps.productRepo.UpdateStockProduct(id, quantity)
}
