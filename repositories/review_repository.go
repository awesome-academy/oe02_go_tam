package repositories

import (
	"errors"
	"gorm.io/gorm"
	"oe02_go_tam/models"
)

type ReviewRepository interface {
	GetByTourID(tourID uint) ([]models.Review, error)
	Create(review *models.Review) error
	GetByIDAndUserID(id, userID uint) (*models.Review, error)
	Update(review *models.Review) error
	Delete(review *models.Review) error
	Exists(reviewID uint) (bool, error)
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

func (r *reviewRepositoryImpl) Create(review *models.Review) error {
	err := r.db.Create(review).Error
	if err != nil {
		return err
	}

	return r.db.Preload("User").First(review, review.ID).Error
}

func (r *reviewRepositoryImpl) GetByIDAndUserID(id, userID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("id = ? AND user_id = ?", id, userID).Preload("User").First(&review).Error
	return &review, err
}

func (r *reviewRepositoryImpl) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepositoryImpl) Delete(review *models.Review) error {
	return r.db.Delete(review).Error
}

func (r *reviewRepositoryImpl) Exists(reviewID uint) (bool, error) {
	var review models.Review
	err := r.db.Select("id").Where("id = ?", reviewID).Limit(1).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
