package services

import (
	"errors"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type AdminBookingService interface {
	GetAllBookings(search string, page, limit int) ([]models.Booking, int64, error)
	GetBookingByID(id uint) (*models.Booking, error)
	DeleteBooking(id uint) error
	CancelBooking(id uint) error
}

type adminBookingServiceImpl struct {
	repo repositories.BookingRepository
}

func NewAdminBookingService(repo repositories.BookingRepository) AdminBookingService {
	return &adminBookingServiceImpl{repo}
}

func (s *adminBookingServiceImpl) GetAllBookings(search string, page, limit int) ([]models.Booking, int64, error) {
	return s.repo.FindAllWithUserAndTour(search, page, limit)
}

func (s *adminBookingServiceImpl) GetBookingByID(id uint) (*models.Booking, error) {
	booking, err := s.repo.FindByIDWithUserAndTour(id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}
	return booking, nil
}

func (s *adminBookingServiceImpl) DeleteBooking(id uint) error {
	return s.repo.Delete(id)
}

func (s *adminBookingServiceImpl) CancelBooking(id uint) error {
	return s.repo.CancelBooking(id)
}
