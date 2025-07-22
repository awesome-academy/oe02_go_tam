package services

import (
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type AdminReviewService interface {
	GetAllReviews(search string, page, limit int) ([]models.Review, int64, error)
	GetReviewByID(id uint) (*models.Review, error)
	DeleteReview(id uint) error
}

type adminReviewService struct {
	repo repositories.ReviewRepository
}

func NewAdminReviewService(repo repositories.ReviewRepository) AdminReviewService {
	return &adminReviewService{repo}
}

func (s *adminReviewService) GetAllReviews(search string, page, limit int) ([]models.Review, int64, error) {
	return s.repo.FindAllWithUserAndTour(search, page, limit)
}

func (s *adminReviewService) GetReviewByID(id uint) (*models.Review, error) {
	return s.repo.FindByIDWithRelations(id)
}

func (s *adminReviewService) DeleteReview(id uint) error {
	return s.repo.DeleteByID(id)
}
