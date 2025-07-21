package repositories

import (
	"errors"
	"gorm.io/gorm"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
)

type BookingRepository interface {
	Create(booking *models.Booking) error
	FindByUserAndTour(userID, tourID uint) (*models.Booking, error)
	Update(booking *models.Booking) error
	GetByIDAndUser(id, userID uint) (*models.Booking, error)
	TotalBookedSeats(tourID uint) (int, error)
	GetByID(bookingID uint) (*models.Booking, error)
}

type bookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepositoryImpl{db}
}

func (r *bookingRepositoryImpl) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepositoryImpl) FindByUserAndTour(userID, tourID uint) (*models.Booking, error) {
	var b models.Booking
	err := r.db.Where("user_id = ? AND tour_id = ? AND status != ?", userID, tourID, constant.BookingStatusCancelled).First(&b).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &b, err
}

func (r *bookingRepositoryImpl) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepositoryImpl) GetByIDAndUser(id, userID uint) (*models.Booking, error) {
	var b models.Booking
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&b).Error
	return &b, err
}

func (r *bookingRepositoryImpl) TotalBookedSeats(tourID uint) (int, error) {
	var total int64
	err := r.db.Model(&models.Booking{}).
		Where("tour_id = ? AND status != ?", tourID, constant.BookingStatusCancelled).
		Select("SUM(number_of_seats)").Scan(&total).Error
	return int(total), err
}

func (r *bookingRepositoryImpl) GetByID(bookingID uint) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.Where("id = ?", bookingID).First(&booking).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}
