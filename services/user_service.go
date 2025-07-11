package services

import (
	"fmt"
	"net/url"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type UserService interface {
	GetProfile(userID uint) (*models.User, error)
	UpdateProfile(userID uint, name string, avatar string) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetProfile(userID uint) (*models.User, error) {
	return s.repo.FindByID(userID)
}

func (s *userService) UpdateProfile(userID uint, name string, avatarURL string) (*models.User, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name cannot be empty", constant.ErrValidation)
	}
	if avatarURL != "" {
		if _, err := url.ParseRequestURI(avatarURL); err != nil {
			return nil, fmt.Errorf("%w: invalid avatar URL format", constant.ErrValidation)
		}
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.AvatarURL = avatarURL

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
