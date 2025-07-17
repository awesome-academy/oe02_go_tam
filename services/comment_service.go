package services

import (
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type CommentService interface {
	CreateComment(comment *models.Comment) error
}

type commentServiceImpl struct {
	repo       repositories.CommentRepository
	reviewRepo repositories.ReviewRepository
}

func NewCommentService(repo repositories.CommentRepository, rr repositories.ReviewRepository) CommentService {
	return &commentServiceImpl{repo, rr}
}

func (s *commentServiceImpl) CreateComment(comment *models.Comment) error {
	exists, _ := s.reviewRepo.Exists(comment.ReviewID)
	if !exists {
		return constant.ErrReviewNotFound
	}

	if comment.ParentID != nil {
		exists, _ := s.repo.Exists(*comment.ParentID)
		if !exists {
			return constant.ErrParentCommentNotFound
		}
	}

	return s.repo.Create(comment)
}
