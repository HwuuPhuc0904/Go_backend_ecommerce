package migrations

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    models "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
    "fmt"
    "go.uber.org/zap"
)

// MigrateDB tạo các bảng dựa trên model đã định nghĩa
func MigrateDB() error {
    global.Logger.Info("Running database migrations...")
    
    // Auto Migrate tạo bảng dựa trên struct model
    err := global.DB.AutoMigrate(
        &models.User{}, &models.Product{},
		&models.Order{}, &models.OrderItem{}, &models.Payment{}, &models.Shipping{},
		&models.Review{}, &models.Cart{}, &models.Address{},
    )
    if err != nil {
        global.Logger.Error("Failed to migrate database", zap.Error(err))
        return fmt.Errorf("failed to migrate database: %w", err)
    }
    
    global.Logger.Info("Database migration completed successfully")
    return nil
}

