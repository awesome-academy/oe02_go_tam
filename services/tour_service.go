package services

import (
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type TourService interface {
	ListToursWithFilters(filters map[string]string, page, size int) ([]models.Tour, error)
	GetTourDetail(id uint) (*models.Tour, error)
}

type tourServiceImpl struct {
	repo repositories.TourRepository
}

func NewTourService(r repositories.TourRepository) TourService {
	return &tourServiceImpl{r}
}

func (s *tourServiceImpl) ListToursWithFilters(filters map[string]string, page, size int) ([]models.Tour, error) {
	return s.repo.GetAllWithFilters(filters, page, size)
}

func (s *tourServiceImpl) GetTourDetail(id uint) (*models.Tour, error) {
	return s.repo.GetByID(id)
}
