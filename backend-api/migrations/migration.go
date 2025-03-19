package migrations

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    model "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
    "fmt"
    "go.uber.org/zap"
)

// MigrateDB tạo các bảng dựa trên model đã định nghĩa
func MigrateDB() error {
    global.Logger.Info("Running database migrations...")
    
    // Auto Migrate tạo bảng dựa trên struct model
    err := global.DB.AutoMigrate(
        &model.User{},
        &model.UserAddress{},
        // Thêm các model khác ở đây khi cần
    )
    
    if err != nil {
        global.Logger.Error("Failed to migrate database", zap.Error(err))
        return fmt.Errorf("failed to migrate database: %w", err)
    }
    
    global.Logger.Info("Database migration completed successfully")
    return nil
}

// SeedData thêm dữ liệu mẫu
func SeedData() error {
    global.Logger.Info("Seeding test data...")
    
    // Kiểm tra xem đã có dữ liệu chưa
    var count int64
    global.DB.Model(&model.User{}).Count(&count)
    
    if count == 0 {
        // Tạo users mẫu
        users := []model.User{
            {
                Email:         "admin@example.com",
                Password:      "$2a$10$CwTycUXWue0Thq9StjUM0uQxTtwSdrsT/3e7S5ye5Z/3oPxgZF7Fq", // password: password
                FirstName:     "Admin",
                LastName:      "User",
                Phone:         "0901234567",
                Role:          "admin",
                EmailVerified: true,
            },
            {
                Email:         "user@example.com",
                Password:      "$2a$10$CwTycUXWue0Thq9StjUM0uQxTtwSdrsT/3e7S5ye5Z/3oPxgZF7Fq", // password: password
                FirstName:     "Regular",
                LastName:      "User",
                Phone:         "0909876543",
                Role:          "user",
                EmailVerified: false,
            },
        }
        
        for _, user := range users {
            if err := global.DB.Create(&user).Error; err != nil {
                global.Logger.Error("Failed to seed users", zap.Error(err))
                return fmt.Errorf("failed to seed users: %w", err)
            }
        }
        
        // Tạo addresses mẫu
        addresses := []model.UserAddress{
            {
                UserID:    1,
                Name:      "Home Address",
                Phone:     "0901234567",
                City:      "Ho Chi Minh City",
                State:     "HCM",
                Country:   "Vietnam",
                IsDefault: true,
            },
            {
                UserID:    1,
                Name:      "Work Address",
                Phone:     "0901234567",
                City:      "Hanoi",
                State:     "HN",
                Country:   "Vietnam",
                IsDefault: false,
            },
        }
        
        for _, address := range addresses {
            if err := global.DB.Create(&address).Error; err != nil {
                global.Logger.Error("Failed to seed addresses", zap.Error(err))
                return fmt.Errorf("failed to seed addresses: %w", err)
            }
        }
        
        global.Logger.Info("Sample data seeded successfully")
    } else {
        global.Logger.Info("Skipping user seed - data already exists")
    }
    
    return nil
}