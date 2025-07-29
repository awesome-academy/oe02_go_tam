package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
	"oe02_go_tam/utils"
	"strconv"
)

type VnpayHandler struct {
	service services.VnpayService
}

func NewVnpayHandler(service services.VnpayService) *VnpayHandler {
	return &VnpayHandler{service}
}

// CreatePaymentUrl godoc
// @Summary Create VNPAY payment URL
// @Description Generates a VNPAY payment URL for a given booking
// @Tags Payment
// @Produce json
// @Param booking_id query int true "Booking ID"
// @Success 200 {object} map[string]string "payment_url"
// @Failure 400 {object} map[string]string "Invalid booking ID or payment generation failed"
// @Security ApiKeyAuth
// @Router /api/payment/create-url [get]
func (h *VnpayHandler) CreatePaymentUrl(c *gin.Context) {
	bookingID, err := strconv.Atoi(c.Query("booking_id"))
	if err != nil || bookingID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("vnpay.invalid_booking_id")})
		return
	}
	userID := c.GetUint("user_id")
	clientIP := utils.GetClientIP(c.Request)

	url, err := h.service.CreatePaymentUrlFromBooking(userID, uint(bookingID), clientIP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_url": url})
}

// VnpayReturn godoc
// @Summary VNPAY return callback
// @Description Callback endpoint for VNPAY to update transaction status
// @Tags Payment
// @Produce json
// @Success 302 {string} string "Redirect to frontend after successful payment"
// @Failure 400 {object} map[string]string "Invalid signature or failed payment"
// @Failure 404 {object} map[string]string "Transaction or booking not found"
// @Failure 500 {object} map[string]string "Failed to update booking"
// @Router /api/payment/vnpay-return [get]
func (h *VnpayHandler) VnpayReturn(c *gin.Context) {
	params := c.Request.URL.Query()

	if !utils.VerifyVnpSignature(params, h.service.GetHashSecret()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("vnpay.invalid_signature")})
		return
	}

	if params.Get("vnp_ResponseCode") != "00" {
		txnRef := params.Get("vnp_TxnRef")
		tx, err := h.service.GetTransactionByTxnRef(txnRef)
		if err == nil {
			tx.Status = constant.BookingStatusFailed
			_ = h.service.UpdateTransaction(tx)
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("vnpay.payment_failed")})
		return
	}

	txnRef := params.Get("vnp_TxnRef")
	tx, err := h.service.GetTransactionByTxnRef(txnRef)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constant.T("vnpay.transaction_not_found")})
		return
	}

	booking, err := h.service.FindPendingBookingByID(uint(tx.BookingID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constant.T("vnpay.booking_not_found")})
		return
	}

	booking.Status = constant.BookingStatusCompleted
	if err := h.service.UpdateBooking(booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constant.T("vnpay.update_booking_failed")})
		return
	}

	tx.Status = "success"
	_ = h.service.UpdateTransaction(tx)

	c.Redirect(http.StatusFound, h.service.GetReturnSuccessURL())
}
