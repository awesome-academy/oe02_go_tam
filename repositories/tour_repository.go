package repositories

import (
	"errors"
	"gorm.io/gorm"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"strconv"
	"time"
)

type TourRepository interface {
	GetByID(id uint) (*models.Tour, error)
	GetAllWithFilters(filters map[string]string, page, size int) ([]models.Tour, error)
	FindAllWithSearch(search string, page, limit int) ([]models.Tour, int64, error)
	FindByID(id uint) (*models.Tour, error)
	Create(tour *models.Tour) error
	Update(tour *models.Tour) error
	Delete(id uint) error
}

type tourRepositoryImpl struct {
	db *gorm.DB
}

func NewTourRepository(db *gorm.DB) TourRepository {
	return &tourRepositoryImpl{db}
}

func (r *tourRepositoryImpl) GetAllWithFilters(filters map[string]string, page, size int) ([]models.Tour, error) {
	var tours []models.Tour

	query := r.db.Preload("Creator")

	if title, ok := filters[constant.FilterTitle]; ok && title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if location, ok := filters[constant.FilterLocation]; ok && location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}
	if start, ok := filters[constant.FilterStartAfter]; ok && start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			query = query.Where("start_date >= ?", t)
		}
	}
	if end, ok := filters[constant.FilterEndBefore]; ok && end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			query = query.Where("end_date <= ?", t)
		}
	}
	if minPrice, ok := filters[constant.FilterMinPrice]; ok && minPrice != "" {
		if p, err := strconv.ParseFloat(minPrice, 64); err == nil {
			query = query.Where("price >= ?", p)
		}
	}
	if maxPrice, ok := filters[constant.FilterMaxPrice]; ok && maxPrice != "" {
		if p, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query = query.Where("price <= ?", p)
		}
	}

	offset := (page - 1) * size
	query = query.Offset(offset).Limit(size)

	err := query.Find(&tours).Error
	return tours, err
}

func (r *tourRepositoryImpl) GetByID(id uint) (*models.Tour, error) {
	var tour models.Tour

	err := r.db.Preload("Creator").
		Preload("Bookings").
		Preload("Reviews").
		Preload("Reviews.User").
		First(&tour, id).Error

	if err != nil {
		return nil, err
	}

	for i := range tour.Reviews {
		r.db.Model(&tour.Reviews[i]).Association("Comments").Find(&tour.Reviews[i].Comments)
		r.db.Model(&tour.Reviews[i]).Association("Likes").Find(&tour.Reviews[i].Likes)
	}

	return &tour, nil
}

func (r *tourRepositoryImpl) FindAllWithSearch(search string, page, limit int) ([]models.Tour, int64, error) {
	var tours []models.Tour
	var total int64

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	query := r.db.Model(&models.Tour{})
	if search != "" {
		query = query.Where("title LIKE ? OR location LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tours).Error
	return tours, total, err
}

func (r *tourRepositoryImpl) FindByID(id uint) (*models.Tour, error) {
	var tour models.Tour
	err := r.db.First(&tour, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tour, nil
}

func (r *tourRepositoryImpl) Create(tour *models.Tour) error {
	return r.db.Create(tour).Error
}

func (r *tourRepositoryImpl) Update(tour *models.Tour) error {
	return r.db.Save(tour).Error
}

func (r *tourRepositoryImpl) Delete(id uint) error {
	var booking models.Booking
	if err := r.db.Select("id").Where("tour_id = ?", id).Limit(1).First(&booking).Error; err == nil {
		return errors.New("cannot delete tour: has related bookings")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var review models.Review
	if err := r.db.Select("id").Where("tour_id = ?", id).Limit(1).First(&review).Error; err == nil {
		return errors.New("cannot delete tour: has related reviews")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	
	var comment models.Comment
	if err := r.db.Raw(`
		SELECT comments.id
		FROM comments
		JOIN reviews ON comments.review_id = reviews.id
		WHERE reviews.tour_id = ?
		LIMIT 1`, id).Scan(&comment).Error; err == nil && comment.ID != 0 {
		return errors.New("cannot delete tour: has related comments via reviews")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.Delete(&models.Tour{}, id).Error
}
