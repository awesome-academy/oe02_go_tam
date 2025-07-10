package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
	"oe02_go_tam/utils"
)

type AuthService interface {
	Register(name, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
	GoogleLogin(name, email, googleID string) (string, *models.User, error)
}

type authService struct {
	repo repositories.AuthRepository
}

func NewAuthService(r repositories.AuthRepository) AuthService {
	return &authService{r}
}

func (s *authService) Register(name, email, password string) (*models.User, error) {
	if existing, _ := s.repo.FindUserByEmail(email); existing != nil {
		return nil, constant.EmailAlreadyExists
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
		if errors.Is(err, constant.ErrUserNotFound) {
			return "", nil, constant.ErrUserNotFound
		}
		return "", nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, constant.LoginFailed
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) GoogleLogin(name, email, googleID string) (string, *models.User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil && !errors.Is(err, constant.ErrUserNotFound) {
		return "", nil, err
	}

	if user != nil {
		if user.GoogleID == "" {
			user.GoogleID = googleID
			if err := s.repo.UpdateUser(user); err != nil {
				return "", nil, err
			}
		} else if user.GoogleID != googleID {
			return "", nil, constant.ErrGoogleIDMismatch
		}
	} else {
		user = &models.User{
			Name:     name,
			Email:    email,
			GoogleID: googleID,
			Role:     "user",
		}
		if err := s.repo.CreateUser(user); err != nil {
			return "", nil, err
		}
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
