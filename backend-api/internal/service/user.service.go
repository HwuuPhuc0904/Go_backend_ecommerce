package service

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/repo"
	model "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"time"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/pkg/utils"
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
	extisting, err := us.UserRepo.GetUserByID(user_update.ID)
	if err != nil {
		return err
	}
	extisting.Name = user_update.Name
	extisting.Email = user_update.Email
	extisting.UpdatedAt = time.Now()
	
	return us.UserRepo.UpdateUser(extisting)
}

//Delete User
func (us *UserService) DeleteUser(id uint) error {
	return us.UserRepo.DeleteUser(id)
}

func (us *UserService) AuthenticateUser(email, password string) (*model.User,string,error) {
	user, err := us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil,"",err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil,"", errors.New("invalid password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
    if err != nil {
        global.Logger.Error("Failed to generate JWT", zap.Error(err))
        return nil, "", err
    }

    return user, token, nil
}

func (us *UserService) ChangePassword(userid uint, currentPassword, newPassword string) error {
	user, err := us.UserRepo.GetUserByID(userid)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return errors.New("invalid password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		global.Logger.Error("Failed to hash password", zap.Error(err))
		return errors.New("internal server error")
	}
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return us.UserRepo.UpdateUser(user)
}


// get all user for pagination
func (us *UserService) GetAllUsers(page, limit int) ([]model.User,int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit >100 {
		limit = 10
	}
	return us.UserRepo.GetAllUsers(page, limit)
}

