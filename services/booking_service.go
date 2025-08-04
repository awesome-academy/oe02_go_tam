package services

import (
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
	"time"
)

type BookingService interface {
	BookTour(userID, tourID uint, numberOfSeats int, startTime, endTime time.Time) (*models.Booking, error)
	CancelBooking(userID, bookingID uint) error
}

type bookingServiceImpl struct {
	bookingRepo repositories.BookingRepository
	tourRepo    repositories.TourRepository
}

func NewBookingService(br repositories.BookingRepository, tr repositories.TourRepository) BookingService {
	return &bookingServiceImpl{br, tr}
}

func (s *bookingServiceImpl) BookTour(userID, tourID uint, seats int, startTime, endTime time.Time) (*models.Booking, error) {
	tour, err := s.tourRepo.GetByID(tourID)
	if err != nil {
		return nil, constant.ErrTourNotFound
	}

	if tour.Seats < seats {
		return nil, constant.ErrNotEnoughSeats
	}

	// Check if user already booked
	if _, err := s.bookingRepo.FindByUserAndTour(userID, tourID); err == nil {
		return nil, constant.ErrAlreadyBooked
	}

	totalBooked, err := s.bookingRepo.TotalBookedSeats(tourID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	availableSeats := tour.Seats - totalBooked
	if availableSeats < seats {
		return nil, constant.ErrNotEnoughSeats
	}

	total := float64(seats) * tour.Price
	booking := &models.Booking{
		UserID:        userID,
		TourID:        tourID,
		Status:        constant.BookingStatusPending,
		NumberOfSeats: seats,
		TotalPrice:    total,
		BookingDate:   time.Now(),
		StartTime:     startTime,
		EndTime:       endTime,
	}

	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *bookingServiceImpl) CancelBooking(userID, bookingID uint) error {
	booking, err := s.bookingRepo.GetByIDAndUser(bookingID, userID)
	if err != nil {
		return constant.ErrBookingNotFound
	}

	if booking.Status == constant.BookingStatusCancelled {
		return constant.ErrAlreadyCancelled
	}

	booking.Status = constant.BookingStatusCancelled
	return s.bookingRepo.Update(booking)
}
