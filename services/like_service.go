package services

import (
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type LikeService interface {
	LikeReview(userID, reviewID uint) error
}

type likeServiceImpl struct {
	repo       repositories.LikeRepository
	reviewRepo repositories.ReviewRepository
}

func NewLikeService(repo repositories.LikeRepository, reviewRepo repositories.ReviewRepository) LikeService {
	return &likeServiceImpl{repo, reviewRepo}
}

func (s *likeServiceImpl) LikeReview(userID, reviewID uint) error {
	exists, _ := s.reviewRepo.Exists(reviewID)
	if !exists {
		return constant.ErrReviewNotFound
	}

	liked, err := s.repo.HasLiked(userID, reviewID)
	if err != nil {
		return err
	}
	if liked {
		return constant.ErrAlreadyLiked
	}

	like := &models.Like{
		UserID:   userID,
		ReviewID: reviewID,
	}
	return s.repo.CreateIfNotExists(like)
}
