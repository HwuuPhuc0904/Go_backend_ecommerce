package repo

type UserRepo struct {}

func  NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (ur * UserRepo) GetInforUser() string {
	return "infor user"
}