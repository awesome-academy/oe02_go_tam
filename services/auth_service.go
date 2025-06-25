package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
	"oe02_go_tam/utils"
)

type AuthService interface {
	Register(name, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

type authService struct {
	repo repositories.AuthRepository
}

func NewAuthService(r repositories.AuthRepository) AuthService {
	return &authService{r}
}

func (s *authService) Register(name, email, password string) (*models.User, error) {
	if existing, _ := s.repo.FindUserByEmail(email); existing != nil {
		return nil, errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
