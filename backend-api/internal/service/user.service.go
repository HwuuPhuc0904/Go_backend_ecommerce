package service

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/repo"
	model "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"errors"
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

func  (us *UserService) GetUserByID(id uint) (*model.User, error) {
	global.Logger.Info("GetUserByID : ", zap.Uint("id", id))
	return us.UserRepo.GetUserByID(id)
}

func (us *UserService) GetUserByEmail(email string) (*model.User, error) {
	global.Logger.Info("GetUserByEmail : ", zap.String("email", email))
	return us.UserRepo.GetUserByEmail(email)
}

func (us *UserService) CreateUser(user *model.User) error {
	//checking for existing user
	existing, err := us.UserRepo.GetUserByEmail(user.Email)
	if err == nil && existing != nil {
		return errors.New("user already exists")
	} else if err != nil && err.Error() != "user not found" {
		return err
	}

	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		global.Logger.Error("Failed to hash password", zap.Error(err))
		return errors.New("internal server error")
	}
	user.Password = string(hashedPassword)

	// setting default 
	user.Role = "user"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return us.UserRepo.CreateUser(user)
	
}

//Update User Information

func (us *UserService) UpdateUser(user_update *model.User) error {
	
}