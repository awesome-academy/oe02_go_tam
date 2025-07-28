package services

import (
	"fmt"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type AdminUsersService interface {
	GetUserList() ([]models.User, error)
	GetUsers(search string, page, limit int) ([]models.User, int64, error)
	GetUserByID(id uint) (*models.User, error)
	ToggleBanUser(id uint) error
	DeleteUser(id uint) error
}

type adminUsersServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewAdminUsersService(userRepo repositories.UserRepository) AdminUsersService {
	return &adminUsersServiceImpl{userRepo}
}

func (s *adminUsersServiceImpl) GetUserList() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *adminUsersServiceImpl) GetUsers(search string, page, limit int) ([]models.User, int64, error) {
	return s.userRepo.FindAllWithSearch(search, page, limit)
}

func (s *adminUsersServiceImpl) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *adminUsersServiceImpl) ToggleBanUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	user.Banned = !user.Banned
	return s.userRepo.Update(user)
}

func (s *adminUsersServiceImpl) DeleteUser(id uint) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("Failed to delete user: %w", err)
	}
	return nil
}
