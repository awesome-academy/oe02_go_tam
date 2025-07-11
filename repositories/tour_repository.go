package repositories

import (
	"gorm.io/gorm"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"strconv"
	"time"
)

type TourRepository interface {
	GetByID(id uint) (*models.Tour, error)
	GetAllWithFilters(filters map[string]string, page, size int) ([]models.Tour, error)
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
