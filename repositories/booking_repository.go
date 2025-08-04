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
	FindAllWithUserAndTour(search string, page, limit int) ([]models.Booking, int64, error)
	FindByIDWithUserAndTour(id uint) (*models.Booking, error)
	Delete(id uint) error
	GetCompletedBookings(search string, page, limit, month, year int) ([]models.Booking, int64, error)
	CancelBooking(id uint) error
	GetMonthlyRevenue(year int) ([]MonthlyRevenue, error)
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
		Select("COALESCE(SUM(number_of_seats), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}

	return int(total), nil
}

func (r *bookingRepositoryImpl) GetByID(bookingID uint) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.Where("id = ?", bookingID).First(&booking).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepositoryImpl) FindAllWithUserAndTour(search string, page, limit int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).
		Preload("User").
		Preload("Tour").
		Joins("JOIN users ON users.id = bookings.user_id").
		Joins("JOIN tours ON tours.id = bookings.tour_id")

	if search != "" {
		query = query.Where("users.name LIKE ? OR tours.title LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("bookings.created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *bookingRepositoryImpl) FindByIDWithUserAndTour(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("User").Preload("Tour").First(&booking, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Booking{}, id).Error
}

func (r *bookingRepositoryImpl) GetCompletedBookings(search string, page, limit, month, year int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).
		Where("status = ?", "completed").
		Preload("User").
		Preload("Tour")

	if search != "" {
		query = query.Joins("JOIN users ON users.id = bookings.user_id").
			Joins("JOIN tours ON tours.id = bookings.tour_id").
			Where("users.name LIKE ? OR tours.title LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if month > 0 {
		query = query.Where("MONTH(booking_date) = ?", month)
	}

	if year > 0 {
		query = query.Where("YEAR(booking_date) = ?", year)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("booking_date DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *bookingRepositoryImpl) CancelBooking(id uint) error {
	return r.db.Model(&models.Booking{}).
		Where("id = ?", id).
		Update("status", constant.BookingStatusCancelled).Error
}

type MonthlyRevenue struct {
	Month int
	Total float64
}

func (r *bookingRepositoryImpl) GetMonthlyRevenue(year int) ([]MonthlyRevenue, error) {
	var result []MonthlyRevenue
	err := r.db.
		Model(&models.Booking{}).
		Select("MONTH(booking_date) as month, SUM(total_price) as total").
		Where("YEAR(booking_date) = ? AND status = ?", year, constant.BookingStatusCompleted).
		Group("MONTH(booking_date)").
		Order("month").
		Scan(&result).Error
	return result, err
}
