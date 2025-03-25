package models

import (
    "time"
)

// User Model
type User struct {
    ID            uint      `gorm:"primaryKey"`
    Name          string    `gorm:"size:100"`
    Email         string    `gorm:"size:100;unique"`
    Password      string    `gorm:"size:100"`
    Role          string    `gorm:"size:20"` // 'admin', 'customer', 'seller'
    IsActive      bool      `gorm:"default:true"`
    EmailVerified bool      `gorm:"default:false"`
    LastLogin     time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
    DeletedAt     *time.Time `gorm:"index"`
    Permissions   []Permission `gorm:"many2many:user_permissions;"`
}
// Role Model
type Role struct {
    ID          uint   `gorm:"primaryKey"`
    Name        string `gorm:"size:50;unique"`
    Description string `gorm:"size:255"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Permissions []Permission `gorm:"many2many:role_permissions;"`
}

// Permission Model - Định nghĩa các quyền chi tiết
type Permission struct {
    ID          uint   `gorm:"primaryKey"`
    Name        string `gorm:"size:100;unique"` 
    Description string `gorm:"size:255"`
    Resource    string `gorm:"size:50"` // 'product', 'order', 'user', etc.
    Action      string `gorm:"size:50"` // 'create', 'read', 'update', 'delete'
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Bảng trung gian để map User và Permission
type UserPermission struct {
    UserID       uint `gorm:"primaryKey"`
    PermissionID uint `gorm:"primaryKey"`
}

// Bảng trung gian để map Role và Permission
type RolePermission struct {
    RoleID       uint `gorm:"primaryKey"`
    PermissionID uint `gorm:"primaryKey"`
}

// Product Model
type Product struct {
    ID               uint   `gorm:"primaryKey"`
    Name             string `gorm:"type:text"`
    Brand            string `gorm:"size:255"`
    ProductDimension string `gorm:"type:text"`
	ProductWeight    string `gorm:"type:text"`
    Description      string `gorm:"type:text"`
    Feature          string `gorm:"type:text"`
    Price            float64
    Stock            int
    ImagesURL        string `gorm:"type:text"`
    UserID           uint
    ParentID         string `gorm:"size:100"`
    IsAvailable      bool
	Categories 	 string `gorm:"type:text"`
}

// Order Model
type Order struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      uint
    OrderDate   time.Time
    Status      string  `gorm:"size:20"`
    TotalAmount float64
}

// Order Item Model
type OrderItem struct {
    ID        uint    `gorm:"primaryKey"`
    OrderID   uint
    ProductID uint
    Quantity  int
    Price     float64
}

// Payment Model
type Payment struct {
    ID            uint      `gorm:"primaryKey"`
    OrderID       uint
    PaymentDate   time.Time
    Amount        float64
    PaymentMethod string `gorm:"size:50"`
    Status        string `gorm:"size:20"`
}

// Shipping Model
type Shipping struct {
    ID             uint      `gorm:"primaryKey"`
    OrderID        uint
    AddressID      uint
    ShippingMethod string `gorm:"size:50"`
    TrackingNumber string `gorm:"size:100"`
    ShippedDate    time.Time
    DeliveredDate  time.Time
}
// Review Model
type Review struct {
    ID        uint      `gorm:"primaryKey"`
    ProductID uint
    UserID    uint
    Rating    int
    Comment   string    `gorm:"type:text"`
    ReviewDate time.Time
}

// Cart Model
type Cart struct {
    ID        uint `gorm:"primaryKey"`
    UserID    uint
    ProductID uint
    Quantity  int
}

// Address Model
type Address struct {
    ID            uint   `gorm:"primaryKey"`
    UserID        uint
    RecipientName string `gorm:"size:100"`
    StreetAddress string `gorm:"size:255"`
    City          string `gorm:"size:100"`
    State         string `gorm:"size:100"`
    PostalCode    string `gorm:"size:20"`
    Country       string `gorm:"size:50"`
    Phone         string `gorm:"size:20"`
}
