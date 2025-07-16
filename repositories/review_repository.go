package repositories

import (
	"gorm.io/gorm"
	"oe02_go_tam/models"
)

type ReviewRepository interface {
	GetByTourID(tourID uint) ([]models.Review, error)
}

type reviewRepositoryImpl struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepositoryImpl{db}
}

func (r *reviewRepositoryImpl) GetByTourID(tourID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("User").Preload("Comments").Preload("Likes").
		Where("tour_id = ?", tourID).Find(&reviews).Error
	return reviews, err
}
