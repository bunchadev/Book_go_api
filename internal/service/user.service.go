package service

import (
	"Book_market_api/internal/models"
	"Book_market_api/internal/repo"
	"time"
)

type UserService struct {
	userRepo *repo.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repo.NewUserRepo(),
	}
}

func (us *UserService) GetUserPagination(page int, limit int, order string, field string, username string, email string) ([]models.UserPagination, error) {
	return us.userRepo.GetUserPagination(page, limit, order, field, username, email)
}

func (us *UserService) CreateUser(user *models.UserCreate) error {
	return us.userRepo.CreateUser(user)
}

func (us *UserService) CheckUserName(username string) bool {
	return us.userRepo.CheckUserName(username)
}

func (us *UserService) UpdateUser(user *models.UpdateUser) error {
	return us.userRepo.UpdateUser(user)
}

func (us *UserService) DeleteUser(id string) error {
	return us.userRepo.DeleteUser(id)
}

func (us *UserService) LoginUser(user *models.LoginUser) (*models.TokenResponse, error) {
	return us.userRepo.LoginUser(user)
}

func (us *UserService) GetNewToken(id string, hour time.Duration) (*models.TokenResponse, error) {
	return us.userRepo.GetNewToken(id, hour)
}

func (us *UserService) LoginSocialMedia(user *models.SocialMedia) (*models.TokenResponse, error) {
	return us.userRepo.LoginSocialMedia(user)
}
