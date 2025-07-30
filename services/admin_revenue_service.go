package services

import (
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
)

type AdminRevenueService interface {
	ListRevenue(search string, page, limit, month, year int) ([]models.Booking, int64, error)
	GetMonthlyRevenue(year int) ([]repositories.MonthlyRevenue, error)
}

type adminRevenueServiceImpl struct {
	repo repositories.BookingRepository
}

func NewAdminRevenueService(repo repositories.BookingRepository) AdminRevenueService {
	return &adminRevenueServiceImpl{repo}
}

func (s *adminRevenueServiceImpl) ListRevenue(search string, page, limit, month, year int) ([]models.Booking, int64, error) {
	return s.repo.GetCompletedBookings(search, page, limit, month, year)
}

func (s *adminRevenueServiceImpl) GetMonthlyRevenue(year int) ([]repositories.MonthlyRevenue, error) {
	return s.repo.GetMonthlyRevenue(year)
}
