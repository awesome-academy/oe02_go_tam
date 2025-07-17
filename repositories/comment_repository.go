package repositories

import (
	"errors"
	"gorm.io/gorm"
	"oe02_go_tam/models"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	Exists(commentID uint) (bool, error)
}

type commentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository { return &commentRepositoryImpl{db: db} }

func (r *commentRepositoryImpl) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepositoryImpl) Exists(commentID uint) (bool, error) {
	var comment models.Comment
	err := r.db.Select("id").Where("id = ?", commentID).Limit(1).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
