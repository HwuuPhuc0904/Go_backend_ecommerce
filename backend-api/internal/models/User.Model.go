package model

import(
	"gorm.io/gorm"
	"time"
)

type User struct {
		ID             uint           `gorm:"primaryKey" json:"id"`
		Email          string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
		Password       string         `gorm:"size:100;not null" json:"-"` 
		FirstName      string         `gorm:"size:50" json:"first_name"`
		LastName       string         `gorm:"size:50" json:"last_name"`
		Phone          string         `gorm:"size:20" json:"phone"`
		Role           string         `gorm:"size:20;default:'user'" json:"role"`
		EmailVerified  bool           `gorm:"default:false" json:"email_verified"`
		LastLogin      *time.Time     `json:"last_login"`
		ProfilePicture string         `json:"profile_picture"`
		CreatedAt      time.Time      `json:"created_at"`
		UpdatedAt      time.Time      `json:"updated_at"`
		DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}


func (User) TableName() string {
	return "users"
}

// UserAddress đại diện cho địa chỉ của người dùng
type UserAddress struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    UserID    uint           `gorm:"not null" json:"user_id"`
    User      User           `gorm:"foreignKey:UserID" json:"-"`
    Name      string         `gorm:"size:100" json:"name"`
    Phone     string         `gorm:"size:20" json:"phone"`
    City      string         `gorm:"size:100" json:"city"`
    State     string         `gorm:"size:100" json:"state"`
    Country   string         `gorm:"size:100" json:"country"`
    IsDefault bool           `gorm:"default:false" json:"is_default"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserAddress) TableName() string {
	return "user_addresses"
}