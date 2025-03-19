package service

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/repo"
	model "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"

)

type UserService struct {
	UserRepo * repo.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		UserRepo: repo.NewUserRepo(),
	}
}

func (us * UserService) GetUserByID(id uint) (*model.User, error) {
	global.Logger.Info("Getting user by Id", zap.Uint("id", id))
	return us.UserRepo.GetUserByID(id)
}

func(us * UserService) GetUserByEmail(email string) (*model.User, error) {
	global.Logger.Info("Getting user by email", zap.String("email", email))
	return us.UserRepo.GetUserByEmail(email)	
}

func (s *UserService) CreateUser(user *model.User) error {
    // Hash password trước khi lưu
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    
    return s.UserRepo.CreateUser(user)
}

func (s *UserService) UpdateProfile(id uint, firstName, lastName, phone string) (*model.User, error) {
    user, err := s.UserRepo.GetUserByID(id)
    if err != nil {
        return nil, err
    }
    
    user.FirstName = firstName
    user.LastName = lastName
    user.Phone = phone
    user.UpdatedAt = time.Now()
    
    if err := s.UserRepo.UpdateUser(user); err != nil {
        return nil, err
    }
    
    return user, nil
}