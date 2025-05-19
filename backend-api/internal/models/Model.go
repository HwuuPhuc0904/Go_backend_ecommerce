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
    Birthday      *time.Time  
    Language      string    `gorm:"size:20"` // 'en', 'vn', etc.
    Country       string    `gorm:"size:50"`
    Gender        string    `gorm:"size:10"` // male , female, gay
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
    ID                uint    `gorm:"primaryKey"`
    SKU               string  `gorm:"size:100"`
    Name              string  `gorm:"type:text"` // product_name
    ShortDescription  string  `gorm:"type:text"`
    Description       string  `gorm:"type:text"`
    Feature           string  `gorm:"type:text"`
    OriginalPrice     float64
    Price             float64
    Discount          float64
    DiscountPercent   float64
    RatingAverage     float64
    ReviewCount       int
    OrderCount        int
    InventoryStatus   string  `gorm:"size:50"`
    IsVisible         bool
    IsAvailable       bool
    Stock             int     // Corresponds to stock_item_qty
    MaxSaleQty        int     // stock_item_max_sale_qty
    BrandID           uint
    Brand             string  `gorm:"size:255"` // brand_name
    MainImage         string  `gorm:"type:text"`
    ImagesURL         string  `gorm:"type:text"` // all_images joined with '|'
    ImageCount        int
    ProductDimension  string  `gorm:"type:text"`
    ProductWeight     string  `gorm:"type:text"`
    SpecificationsFull string  `gorm:"type:text"` // JSON string of specifications
    UserID            uint
    ParentID          string  `gorm:"size:100"`
    Categories        string  `gorm:"type:text"`
}

// Order Model
type Order struct {
    ID                 uint          `gorm:"primaryKey" json:"id"`
    OrderNumber        string        `gorm:"size:50;uniqueIndex" json:"order_number"`
    UserID             uint          `json:"user_id"`
    User               User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
    OrderItems         []OrderItem   `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
    OrderDate          time.Time     `json:"order_date"`
    Status             string        `gorm:"size:30" json:"status"` // 'pending', 'processing', 'shipped', 'delivered', 'canceled'
    SubtotalAmount     float64       `json:"subtotal_amount"` // Sum of all items before discounts
    DiscountAmount     float64       `json:"discount_amount"` // Total discount applied
    ShippingAmount     float64       `json:"shipping_amount"` // Shipping fee
    TaxAmount          float64       `json:"tax_amount"`      // Tax applied
    TotalAmount        float64       `json:"total_amount"`    // Final amount charged
    PaymentID          uint          `json:"payment_id,omitempty"`
    Payment            Payment       `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
    ShippingID         uint          `json:"shipping_id,omitempty"`
    Shipping           Shipping      `gorm:"foreignKey:ShippingID" json:"shipping,omitempty"`
    BillingAddressID   uint          `json:"billing_address_id"`
    BillingAddress     Address       `gorm:"foreignKey:BillingAddressID" json:"billing_address,omitempty"`
    ShippingAddressID  uint          `json:"shipping_address_id"`
    ShippingAddress    Address       `gorm:"foreignKey:ShippingAddressID" json:"shipping_address,omitempty"`
    CouponCode         string        `gorm:"size:50" json:"coupon_code,omitempty"`
    Notes              string        `gorm:"type:text" json:"notes,omitempty"`
    CreatedAt          time.Time     `json:"created_at"`
    UpdatedAt          time.Time     `json:"updated_at"`
    DeletedAt          *time.Time    `gorm:"index" json:"deleted_at,omitempty"`
}

// Order Item Model
type OrderItem struct {
    ID               uint          `gorm:"primaryKey" json:"id"`
    OrderID          uint          `json:"order_id"`
    Order            Order         `gorm:"foreignKey:OrderID" json:"-"` // Avoid circular reference in JSON
    ProductID        uint          `json:"product_id"`
    Product          Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
    SKU              string        `gorm:"size:100" json:"sku"`         // Stored at time of order
    ProductName      string        `gorm:"size:255" json:"product_name"` // Stored at time of order
    ProductImage     string        `gorm:"type:text" json:"product_image"` // Stored at time of order
    Quantity         int           `json:"quantity"`
    UnitPrice        float64       `json:"unit_price"` // Price per unit at time of order
    Subtotal         float64       `json:"subtotal"`   // unit_price * quantity
    Discount         float64       `json:"discount"`   // Discount applied to this item
    FinalPrice       float64       `json:"final_price"` // after discount
    Attributes       string        `gorm:"type:text" json:"attributes,omitempty"` // Color, size, etc. as JSON
    CreatedAt        time.Time     `json:"created_at"`
    UpdatedAt        time.Time     `json:"updated_at"`
}

// Payment Model
type Payment struct {
    ID                uint          `gorm:"primaryKey" json:"id"`
    OrderID           uint          `json:"order_id"`
    TransactionID     string        `gorm:"size:100" json:"transaction_id"` // ID from payment gateway
    PaymentMethod     string        `gorm:"size:50" json:"payment_method"`  // 'credit_card', 'paypal', etc.
    PaymentProvider   string        `gorm:"size:50" json:"payment_provider,omitempty"` // 'stripe', 'paypal', etc.
    Amount            float64       `json:"amount"`
    Currency          string        `gorm:"size:3;default:USD" json:"currency"` // ISO currency code
    Status            string        `gorm:"size:30" json:"status"` // 'pending', 'completed', 'failed', 'refunded'
    PaymentDate       time.Time     `json:"payment_date"`
    LastFourDigits    string        `gorm:"size:4" json:"last_four_digits,omitempty"` // For cards
    CardType          string        `gorm:"size:20" json:"card_type,omitempty"` // 'visa', 'mastercard', etc.
    BillingEmail      string        `gorm:"size:100" json:"billing_email,omitempty"`
    PaymentDetails    string        `gorm:"type:text" json:"payment_details,omitempty"` // Additional JSON details
    RefundAmount      float64       `json:"refund_amount,omitempty"`
    RefundDate        *time.Time    `json:"refund_date,omitempty"`
    RefundReason      string        `gorm:"type:text" json:"refund_reason,omitempty"`
    CreatedAt         time.Time     `json:"created_at"`
    UpdatedAt         time.Time     `json:"updated_at"`
}

// Shipping Model
type Shipping struct {
    ID                uint          `gorm:"primaryKey" json:"id"`
    OrderID           uint          `json:"order_id"`
    ShippingMethod    string        `gorm:"size:50" json:"shipping_method"` // 'standard', 'express', etc.
    Carrier           string        `gorm:"size:50" json:"carrier"`         // 'UPS', 'FedEx', etc.
    TrackingNumber    string        `gorm:"size:100" json:"tracking_number,omitempty"`
    TrackingURL       string        `gorm:"type:text" json:"tracking_url,omitempty"`
    ShippingCost      float64       `json:"shipping_cost"`
    EstimatedDelivery *time.Time    `json:"estimated_delivery,omitempty"`
    ShippedDate       *time.Time    `json:"shipped_date,omitempty"`
    DeliveredDate     *time.Time    `json:"delivered_date,omitempty"`
    Status            string        `gorm:"size:30" json:"status"` // 'pending', 'shipped', 'delivered'
    RequireSignature  bool          `gorm:"default:false" json:"require_signature"`
    ShippingWeight    float64       `json:"shipping_weight,omitempty"` // in kg
    ShippingNotes     string        `gorm:"type:text" json:"shipping_notes,omitempty"`
    CreatedAt         time.Time     `json:"created_at"`
    UpdatedAt         time.Time     `json:"updated_at"`
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
    ID                uint          `gorm:"primaryKey" json:"id"`
    UserID            uint          `json:"user_id"`
    User              User          `gorm:"foreignKey:UserID" json:"-"` 
    AddressType       string        `gorm:"size:20" json:"address_type"` // 'billing', 'shipping', 'both'
    IsDefault         bool          `gorm:"default:false" json:"is_default"`
    RecipientName     string        `gorm:"size:100" json:"recipient_name"`
    CompanyName       string        `gorm:"size:100" json:"company_name,omitempty"`
    StreetAddress1    string        `gorm:"size:255" json:"street_address_1"`
    StreetAddress2    string        `gorm:"size:255" json:"street_address_2,omitempty"`
    City              string        `gorm:"size:100" json:"city"`
    State             string        `gorm:"size:100" json:"state"`
    PostalCode        string        `gorm:"size:20" json:"postal_code"`
    Country           string        `gorm:"size:50" json:"country"`
    Phone             string        `gorm:"size:20" json:"phone"`
    AlternatePhone    string        `gorm:"size:20" json:"alternate_phone,omitempty"`
    Email             string        `gorm:"size:100" json:"email,omitempty"`
    DeliveryNotes     string        `gorm:"type:text" json:"delivery_notes,omitempty"`
    Latitude          float64       `json:"latitude,omitempty"`  // For map integration
    Longitude         float64       `json:"longitude,omitempty"` // For map integration
    CreatedAt         time.Time     `json:"created_at"`
    UpdatedAt         time.Time     `json:"updated_at"`
}