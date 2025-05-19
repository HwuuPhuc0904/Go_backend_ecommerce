package repo

import(
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	model "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"gorm.io/gorm"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

type ProductRepo struct {
	db * gorm.DB
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{
		db: global.DB,
	}
}

func(pr *ProductRepo) GetProductByCategory(category string, limit int , offset int, shortBy, orderBy string) ([]model.Product, int64, error) {
    var products []model.Product 
	var total int64

    searchPattern := "%" + category + "%"
        
	query := pr.db.Model(&model.Product{}).Where("categories LIKE ?", searchPattern)
	
	if err := query.Count(&total).Error; err != nil {
        global.Logger.Error("Error counting products", zap.Error(err))
        return nil, 0, err
    }

	if shortBy == "" {
		shortBy = "id"
	}
	if orderBy != "asc" && orderBy != "desc" {
		orderBy = "asc"
	}

	order := shortBy + " " + orderBy

	if err := pr.db.Where("categories LIKE ?", searchPattern).
        Order(order).
        Limit(limit).
        Offset(offset).
        Find(&products).Error; err != nil {
        global.Logger.Error("Error fetching products", zap.Error(err))
        return nil, 0, err
    }

    return products, total, nil
}

func (pr * ProductRepo) GetAllProducts(limit int, offset int, sortBy string, orderBy string) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := pr.db.Model(&model.Product{})
	
	if err := query.Count(&total).Error; err != nil {
        global.Logger.Error("Error counting products", zap.Error(err))
        return nil, 0, err
    }

	if err := query.
        Order(fmt.Sprintf("%s %s", sortBy, orderBy)).
        Limit(limit).
        Offset(offset).
        Find(&products).Error; err != nil {
        global.Logger.Error("Error fetching products", zap.Error(err))
        return nil, 0, err
    }
	return products, total, nil
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

