package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
	"strconv"
)

type BookingHandler struct {
	service services.BookingService
}

func NewBookingHandler(service services.BookingService) *BookingHandler {
	return &BookingHandler{service}
}

type BookTourRequest struct {
	TourID        uint `json:"tour_id" binding:"required"`
	NumberOfSeats int  `json:"number_of_seats" binding:"required"`
}

// BookTour godoc
// @Summary Book a tour
// @Description Book a tour with a given number of seats
// @Tags Booking
// @Accept json
// @Produce json
// @Param body body BookTourRequest true "Booking request payload"
// @Success 201 {object} map[string]interface{} "Booking created successfully"
// @Failure 400 {object} map[string]string "Bad request or not enough seats or already booked"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Tour not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/bookings [post]
func (h *BookingHandler) BookTour(c *gin.Context) {
	var req BookTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": []interface{}{}})
		return
	}

	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.T("auth.unauthorized"), "data": []interface{}{}})
		return
	}

	userID, ok := val.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.T("auth.user_id.invalid"), "data": []interface{}{}})
		return
	}

	booking, err := h.service.BookTour(userID, req.TourID, req.NumberOfSeats)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, gin.H{
			"message": constant.T("booking.create.success"),
			"data": gin.H{
				"booking_id":      booking.ID,
				"tour_id":         booking.TourID,
				"number_of_seats": booking.NumberOfSeats,
				"total_price":     booking.TotalPrice,
				"status":          booking.Status,
				"booking_date":    booking.BookingDate,
			},
		})
	case constant.ErrTourNotFound:
		c.JSON(http.StatusNotFound, gin.H{"message": constant.T("booking.error.tour_not_found"), "data": []interface{}{}})
	case constant.ErrNotEnoughSeats:
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.T("booking.error.not_enough_seats"), "data": []interface{}{}})
	case constant.ErrAlreadyBooked:
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.T("booking.error.already_booked"), "data": []interface{}{}})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.T("booking.error.failed"), "data": []interface{}{}})
	}
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a booking by ID for the authenticated user
// @Tags Booking
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string "Booking cancelled successfully"
// @Failure 400 {object} map[string]string "Invalid booking ID or already cancelled"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Booking not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/bookings/{id}/cancel [post]
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.T("booking.error.invalid_id"), "data": []interface{}{}})
		return
	}

	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.T("auth.unauthorized"), "data": []interface{}{}})
		return
	}

	userID, ok := val.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.T("auth.user_id.invalid"), "data": []interface{}{}})
		return
	}

	err = h.service.CancelBooking(userID, uint(id))
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"message": constant.T("booking.cancel.success"), "data": []interface{}{}})
	case constant.ErrBookingNotFound:
		c.JSON(http.StatusNotFound, gin.H{"message": constant.T("booking.error.not_found"), "data": []interface{}{}})
	case constant.ErrAlreadyCancelled:
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.T("booking.error.already_cancelled"), "data": []interface{}{}})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.T("booking.error.failed"), "data": []interface{}{}})
	}
}
