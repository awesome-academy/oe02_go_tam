package repositories

import (
	"errors"
	"gorm.io/gorm"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}

func (r *authRepositoryImpl) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepositoryImpl) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, constant.ErrUserNotFound
	}
	return &user, err
}

func (r *authRepositoryImpl) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}
