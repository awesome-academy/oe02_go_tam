package services

import (
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type AdminTourService interface {
	GetTours(search string, page, limit int) ([]models.Tour, int64, error)
	GetTourByID(id uint) (*models.Tour, error)
	CreateTour(tour *models.Tour) error
	UpdateTour(tour *models.Tour) error
	DeleteTour(id uint) error
}

type adminTourService struct {
	tourRepo repositories.TourRepository
}

func NewAdminTourService(r repositories.TourRepository) AdminTourService {
	return &adminTourService{tourRepo: r}
}

func (s *adminTourService) GetTours(search string, page, limit int) ([]models.Tour, int64, error) {
	return s.tourRepo.FindAllWithSearch(search, page, limit)
}

func (s *adminTourService) GetTourByID(id uint) (*models.Tour, error) {
	tour, err := s.tourRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if tour == nil {
		return nil, constant.ErrTourNotFound
	}

	return tour, nil
}

func (s *adminTourService) CreateTour(tour *models.Tour) error {
	return s.tourRepo.Create(tour)
}

func (s *adminTourService) UpdateTour(tour *models.Tour) error {
	return s.tourRepo.Update(tour)
}

func (s *adminTourService) DeleteTour(id uint) error {
	return s.tourRepo.Delete(id)
}
