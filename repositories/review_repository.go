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
	FindAllWithUserAndTour(search string, page, limit int) ([]models.Review, int64, error)
	FindByIDWithRelations(id uint) (*models.Review, error)
	DeleteByID(id uint) error
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

func (r *reviewRepositoryImpl) FindAllWithUserAndTour(search string, page, limit int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.Model(&models.Review{}).
		Preload("User").
		Preload("Tour").
		Preload("Comments").
		Preload("Likes").
		Joins("JOIN users ON users.id = reviews.user_id").
		Joins("JOIN tours ON tours.id = reviews.tour_id")

	if search != "" {
		query = query.Where("users.name LIKE ? OR tours.title LIKE ? OR reviews.content LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("reviews.created_at DESC").Offset(offset).Limit(limit).Find(&reviews).Error
	return reviews, total, err
}

func (r *reviewRepositoryImpl) FindByIDWithRelations(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("User").
		Preload("Tour").
		Preload("Comments.User").
		Preload("Comments.Replies").
		Preload("Likes").
		First(&review, id).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &review, err
}

func (r *reviewRepositoryImpl) DeleteByID(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}
