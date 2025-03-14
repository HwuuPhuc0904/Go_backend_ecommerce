package service

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/repo"
)

type UserService struct {
	userRepo  *repo.UserRepo // repo là tên của trường
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repo.NewUserRepo(),
	}
}

func (us *UserService) GetUserByID() string {
	return us.userRepo.GetInforUser()
}