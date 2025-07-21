package services

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"oe02_go_tam/config"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/repositories"
	"oe02_go_tam/utils"
	"strconv"
	"time"
)

type VnpayService interface {
	CreatePaymentUrlFromBooking(userID, bookingID uint, clientIP string) (string, error)
	GetHashSecret() string
	FindPendingBooking(userID, tourID uint) (*models.Booking, error)
	UpdateBooking(booking *models.Booking) error
	GetReturnSuccessURL() string
	FindPendingBookingByID(bookingID uint) (*models.Booking, error)
	GetTransactionByTxnRef(txnRef string) (*models.PaymentTransaction, error)
	UpdateTransaction(tx *models.PaymentTransaction) error
}

type vnpayServiceImpl struct {
	bookingRepo     repositories.BookingRepository
	tourRepo        repositories.TourRepository
	transactionRepo repositories.TransactionRepository
	cfg             config.VnpayConfig
}

func NewVnpayService(br repositories.BookingRepository, tr repositories.TourRepository, trr repositories.TransactionRepository, cfg config.VnpayConfig) VnpayService {
	return &vnpayServiceImpl{br, tr, trr, cfg}
}

func (s *vnpayServiceImpl) CreatePaymentUrlFromBooking(userID, bookingID uint, clientIP string) (string, error) {
	booking, err := s.bookingRepo.GetByIDAndUser(bookingID, userID)
	if err != nil {
		return "", constant.ErrBookingNotFound
	}

	if booking.Status != constant.BookingStatusPending {
		return "", constant.ErrBookingAlreadyProcessed
	}

	tour, err := s.tourRepo.GetByID(booking.TourID)
	if err != nil {
		return "", constant.ErrTourNotFound
	}

	txnRef := strconv.FormatInt(time.Now().UnixNano(), 10)
	amount := int(math.Round(float64(booking.NumberOfSeats) * tour.Price * 100))
	err = s.transactionRepo.Create(&models.PaymentTransaction{
		TxnRef:    txnRef,
		BookingID: booking.ID,
		Status:    "pending",
	})
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return "", err
	}

	params := url.Values{}
	params.Add("vnp_Version", "2.1.0")
	params.Add("vnp_Command", "pay")
	params.Add("vnp_TmnCode", s.cfg.TmnCode)
	params.Add("vnp_Amount", strconv.Itoa(amount))
	params.Add("vnp_CurrCode", "VND")
	params.Add("vnp_TxnRef", txnRef)
	params.Add("vnp_OrderInfo", fmt.Sprintf("Booking tour %s", tour.Title))
	params.Add("vnp_OrderType", "other")
	params.Add("vnp_Locale", "vn")
	params.Add("vnp_ReturnUrl", s.cfg.ReturnURL)
	params.Add("vnp_IpAddr", clientIP)
	params.Add("vnp_CreateDate", time.Now().Format("20060102150405"))
	params.Add("vnp_ExpireDate", time.Now().Add(15*time.Minute).Format("20060102150405"))

	signedUrl := utils.BuildVnpUrl(params, s.cfg.HashSecret, s.cfg.PayURL, s.cfg.HashType)

	return signedUrl, nil
}

func (s *vnpayServiceImpl) GetHashSecret() string {
	return s.cfg.HashSecret
}

func (s *vnpayServiceImpl) FindPendingBooking(userID, tourID uint) (*models.Booking, error) {
	booking, err := s.bookingRepo.FindByUserAndTour(userID, tourID)
	if err != nil {
		return nil, err
	}
	if booking.Status != constant.BookingStatusPending {
		return nil, constant.ErrBookingAlreadyProcessed
	}
	return booking, nil
}

func (s *vnpayServiceImpl) UpdateBooking(booking *models.Booking) error {
	return s.bookingRepo.Update(booking)
}

func (s *vnpayServiceImpl) GetReturnSuccessURL() string {
	return s.cfg.ReturnURL
}

func (s *vnpayServiceImpl) FindPendingBookingByID(bookingID uint) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}
	if booking.Status != constant.BookingStatusPending {
		return nil, constant.ErrBookingAlreadyProcessed
	}
	return booking, nil
}

func (s *vnpayServiceImpl) GetTransactionByTxnRef(txnRef string) (*models.PaymentTransaction, error) {
	return s.transactionRepo.FindByTxnRef(txnRef)
}

func (s *vnpayServiceImpl) UpdateTransaction(tx *models.PaymentTransaction) error {
	return s.transactionRepo.Update(tx)
}
