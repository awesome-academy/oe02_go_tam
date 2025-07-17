package services

import (
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
	"strings"
)

type ReviewService interface {
	GetReviews(tourID uint) ([]models.Review, error)
	CreateReview(userID uint, tourID uint, rating int, content string) (*models.Review, error)
	GetOwnReview(id, userID uint) (*models.Review, error)
	UpdateReview(id, userID uint, rating int, content string) (*models.Review, error)
	DeleteReview(id, userID uint) error
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

func (s *reviewServiceImpl) CreateReview(userID uint, tourID uint, rating int, content string) (*models.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, constant.ErrInvalidRating
	}
	if strings.TrimSpace(content) == "" {
		return nil, constant.ErrEmptyContent
	}

	review := &models.Review{
		UserID:  userID,
		TourID:  tourID,
		Rating:  rating,
		Content: content,
	}
	if err := s.repo.Create(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *reviewServiceImpl) GetOwnReview(id, userID uint) (*models.Review, error) {
	return s.repo.GetByIDAndUserID(id, userID)
}

func (s *reviewServiceImpl) UpdateReview(id, userID uint, rating int, content string) (*models.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, constant.ErrInvalidRating
	}
	if strings.TrimSpace(content) == "" {
		return nil, constant.ErrEmptyContent
	}

	review, err := s.repo.GetByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	review.Rating = rating
	review.Content = content

	if err := s.repo.Update(review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewServiceImpl) DeleteReview(id, userID uint) error {
	review, err := s.repo.GetByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	return s.repo.Delete(review)
}
