package repo

import(
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	model "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"gorm.io/gorm"
	"errors"
)

type ProductRepo struct {
	db * gorm.DB
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{
		db: global.DB,
	}
}

func (pr * ProductRepo) GetProductByID(id uint) (*model.Product, error) {
	var product model.Product
	result := pr.db.First(&product, id)
	if(result.Error != nil) {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, result.Error
	}
	return &product , nil 
}

func (pr * ProductRepo) CreateProduct(product *model.Product) error {
	return pr.db.Create(product).Error	
}

func (pr * ProductRepo) UpdateProduct(product *model.Product) error {
	return pr.db.Save(product).Error
}

func (pr * ProductRepo) DeleteProductByID(id uint) error {
	return pr.db.Delete(&model.Product{}, id).Error
}

func (pr * ProductRepo) UpdateStockProduct(id uint, quantity int) error {
	return pr.db.Transaction(func(tx *gorm.DB) error {
		var product model.Product
		result := tx.First(&product, id)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("product not found")
			}
			return result.Error
		}
		product.Stock = quantity
		result = tx.Save(&product)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func (pr * ProductRepo) SearchProduct(keyword string) ([]model.Product, error) {
	var products []model.Product
	result := pr.db.Where("name LIKE ?", "%" + keyword + "%").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

