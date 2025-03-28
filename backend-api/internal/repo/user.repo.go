package repo

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	model "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepo struct {
	db * gorm.DB
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		db: global.DB,
	}
}

func (ur * UserRepo) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := ur.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")		
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur * UserRepo) GetUserByEmail(email string) (*model.User, error){
	var user model.User
	result := ur.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// Create User
func (ur * UserRepo) CreateUser(user *model.User) error {
	result := ur.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//Update User
func (ur * UserRepo) UpdateUser(user_update *model.User) error {
	result := ur.db.Save(user_update)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//Delete User
func (ur * UserRepo) DeleteUser(id uint) error {
	return ur.db.Delete(&model.User{}, id).Error
}

// get all user cho phân trang
func (ur *UserRepo) GetAllUsers(page, limit int) ([]model.User, int64, error) {
    var users []model.User
    var total int64

	// location to start get data from database
    offset := (page - 1) * limit

    // Đếm tổng số lượng người dùng
    if err := ur.db.Model(&model.User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Lấy danh sách người dùng với phân trang
    if err := ur.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
}