package repositories

import (
	"oe02_go_tam/database"
	"oe02_go_tam/models"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
}

type authRepositoryImpl struct {
}

func NewAuthRepository() AuthRepository {
	return &authRepositoryImpl{}
}

func (r *authRepositoryImpl) CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *authRepositoryImpl) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
