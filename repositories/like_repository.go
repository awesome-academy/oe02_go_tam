package repositories

import (
	"errors"
	"gorm.io/gorm"
	"oe02_go_tam/models"
)

type LikeRepository interface {
	CreateIfNotExists(like *models.Like) error
	HasLiked(userID, reviewID uint) (bool, error)
}

type likeRepositoryImpl struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepositoryImpl{db}
}

func (r *likeRepositoryImpl) CreateIfNotExists(like *models.Like) error {
	var existing models.Like
	err := r.db.Where("user_id = ? AND review_id = ?", like.UserID, like.ReviewID).First(&existing).Error
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(like).Error
	}
	return err
}

func (r *likeRepositoryImpl) HasLiked(userID, reviewID uint) (bool, error) {
	var like models.Like
	err := r.db.Select("id").Where("user_id = ? AND review_id = ?", userID, reviewID).
		Limit(1).First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
