package migrations

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    models "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
    "fmt"
    "go.uber.org/zap"
    "time"
)

func SeedRolesAndPermissions() error {
    global.Logger.Info("Seeding roles and permissions...")
    
    // Kiểm tra xem đã có roles chưa
    var roleCount int64
    global.DB.Model(&models.Role{}).Count(&roleCount)
    
    if roleCount == 0 {
        roles := []models.Role{
            {
                Name:        "admin",
                Description: "Administrator with full access",
                CreatedAt:   time.Now(),
                UpdatedAt:   time.Now(),
            },
            {
                Name:        "seller",
                Description: "Can manage own products and view own orders",
                CreatedAt:   time.Now(),
                UpdatedAt:   time.Now(),
            },
            {
                Name:        "customer",
                Description: "Regular user who can place orders",
                CreatedAt:   time.Now(),
                UpdatedAt:   time.Now(),
            },
        }
        
        for _, role := range roles {
            if err := global.DB.Create(&role).Error; err != nil {
                global.Logger.Error("Failed to seed roles", zap.Error(err))
                return fmt.Errorf("failed to seed roles: %w", err)
            }
        }
        
        // Tạo Permissions
        resourceActions := map[string][]string{
            "product": {"create", "read", "update", "delete"},
            "order":   {"create", "read", "update", "delete"},
            "user":    {"create", "read", "update", "delete"},
            "payment": {"create", "read", "update"},
            "review":  {"create", "read", "update", "delete"},
        }
        
        var permissions []models.Permission
        
        for resource, actions := range resourceActions {
            for _, action := range actions {
                perm := models.Permission{
                    Name:        fmt.Sprintf("%s:%s", resource, action),
                    Description: fmt.Sprintf("Can %s %s", action, resource),
                    Resource:    resource,
                    Action:      action,
                    CreatedAt:   time.Now(),
                    UpdatedAt:   time.Now(),
                }
                permissions = append(permissions, perm)
                
                if err := global.DB.Create(&perm).Error; err != nil {
                    global.Logger.Error("Failed to seed permissions", zap.Error(err))
                    return fmt.Errorf("failed to seed permissions: %w", err)
                }
            }
        }
        
        // Gán permissions cho roles
        
        // 1. Lấy role và permissions từ DB
        var adminRole models.Role
        var sellerRole models.Role
        var customerRole models.Role
        
        global.DB.Where("name = ?", "admin").First(&adminRole)
        global.DB.Where("name = ?", "seller").First(&sellerRole)
        global.DB.Where("name = ?", "customer").First(&customerRole)
        
        var allPermissions []models.Permission
        global.DB.Find(&allPermissions)
        
        var sellerPermissions []models.Permission
        global.DB.Where("resource = ? OR resource = ? OR (resource = ? AND action = ?)", 
            "product", "review", "order", "read").Find(&sellerPermissions)
        
        var customerPermissions []models.Permission
        global.DB.Where("(resource = ? AND action = ?) OR (resource = ? AND action IN ?) OR (resource = ? AND action IN ?)",
            "product", "read", 
            "order", []string{"create", "read"},
            "review", []string{"create", "read", "update"}).Find(&customerPermissions)
        
        // 2. Gán permissions cho admin (tất cả quyền)
        if err := global.DB.Model(&adminRole).Association("Permissions").Append(&allPermissions); err != nil {
            global.Logger.Error("Failed to assign permissions to admin", zap.Error(err))
            return fmt.Errorf("failed to assign permissions to admin: %w", err)
        }
        
        // 3. Gán permissions cho seller
        if err := global.DB.Model(&sellerRole).Association("Permissions").Append(&sellerPermissions); err != nil {
            global.Logger.Error("Failed to assign permissions to seller", zap.Error(err))
            return fmt.Errorf("failed to assign permissions to seller: %w", err)
        }
        
        // 4. Gán permissions cho customer
        if err := global.DB.Model(&customerRole).Association("Permissions").Append(&customerPermissions); err != nil {
            global.Logger.Error("Failed to assign permissions to customer", zap.Error(err))
            return fmt.Errorf("failed to assign permissions to customer: %w", err)
        }
        
        global.Logger.Info("Roles and permissions seeded successfully")
    } else {
        global.Logger.Info("Skipping roles and permissions seed - data already exists")
    }
    
    return nil
}

// SeedData khởi tạo dữ liệu mẫu
func SeedData() error {
    global.Logger.Info("Seeding test data...")
    
    // Gọi hàm seed roles và permissions
    if err := SeedRolesAndPermissions(); err != nil {
        return err
    }
    
    // Kiểm tra xem đã có dữ liệu chưa
    var count int64
    global.DB.Model(&models.User{}).Count(&count)
    
    if count == 0 {
        // Lấy roles từ database
        var adminRole models.Role
        var customerRole models.Role
        
        global.DB.Where("name = ?", "admin").First(&adminRole)
        global.DB.Where("name = ?", "customer").First(&customerRole)
        
        // Tạo users mẫu
        adminUser := models.User{
            Email:         "admin@example.com",
            Password:      "$2a$10$CwTycUXWue0Thq9StjUM0uQxTtwSdrsT/3e7S5ye5Z/3oPxgZF7Fq", // password: password
            Name:          "Admin User",
            Role:          "admin",
            IsActive:      true,
            EmailVerified: true,
			LastLogin:     time.Now(),
            CreatedAt:     time.Now(),
            UpdatedAt:     time.Now(),
        }
        
        customerUser := models.User{
            Email:         "customer@example.com",
            Password:      "$2a$10$CwTycUXWue0Thq9StjUM0uQxTtwSdrsT/3e7S5ye5Z/3oPxgZF7Fq", // password: password
            Name:          "Customer User",
            Role:          "customer",
            IsActive:      true,
            EmailVerified: false,
			LastLogin:     time.Now(),
            CreatedAt:     time.Now(),
            UpdatedAt:     time.Now(),
        }
        
        if err := global.DB.Create(&adminUser).Error; err != nil {
            global.Logger.Error("Failed to seed admin user", zap.Error(err))
            return fmt.Errorf("failed to seed admin user: %w", err)
        }
        
        if err := global.DB.Create(&customerUser).Error; err != nil {
            global.Logger.Error("Failed to seed customer user", zap.Error(err))
            return fmt.Errorf("failed to seed customer user: %w", err)
        }
        
        // Lấy các quyền từ database cho admin user
        var adminPermissions []models.Permission
        global.DB.Where("resource = ? AND action = ?", "user", "read").Find(&adminPermissions)
        
        // Gán quyền trực tiếp cho admin user
        if err := global.DB.Model(&adminUser).Association("Permissions").Append(&adminPermissions); err != nil {
            global.Logger.Error("Failed to assign permissions to admin user", zap.Error(err))
            return fmt.Errorf("failed to assign permissions to admin user: %w", err)
        }
        
        // Tạo địa chỉ mẫu
        address := models.Address{
            UserID:        adminUser.ID,
            RecipientName: "Admin User",
            StreetAddress: "123 Admin St",
            City:          "Ho Chi Minh City",
            State:         "HCM",
            PostalCode:    "70000",
            Country:       "Vietnam",
            Phone:         "0901234567",
        }
        
        if err := global.DB.Create(&address).Error; err != nil {
            global.Logger.Error("Failed to seed address", zap.Error(err))
            return fmt.Errorf("failed to seed address: %w", err)
        }
        
        global.Logger.Info("Sample data seeded successfully")
    } else {
        global.Logger.Info("Skipping user seed - data already exists")
    }
    
    return nil
}