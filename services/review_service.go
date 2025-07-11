package services

import (
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type ReviewService interface {
	GetReviews(tourID uint) ([]models.Review, error)
}

type reviewServiceImpl struct {
	repo repositories.ReviewRepository
}

func NewReviewService(r repositories.ReviewRepository) ReviewService {
	return &reviewServiceImpl{r}
}

func (s *reviewServiceImpl) GetReviews(tourID uint) ([]models.Review, error) {
	return s.repo.GetByTourID(tourID)
}
