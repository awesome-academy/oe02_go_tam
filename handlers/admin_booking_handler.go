package handlers

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
	"strconv"
)

type AdminBookingHandler struct {
	service services.AdminBookingService
}

func NewAdminBookingHandler(s services.AdminBookingService) *AdminBookingHandler {
	return &AdminBookingHandler{s}
}

func (h *AdminBookingHandler) ListBookings(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	bookings, total, err := h.service.GetAllBookings(search, page, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{"error": constant.T("admin.booking.error.load"), "Title": constant.T("error.title")})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.HTML(http.StatusOK, constant.TemplateBooking, gin.H{
		"Bookings": bookings,
		"Search":   search,
		"Pagination": gin.H{
			"Page":       page,
			"Limit":      limit,
			"TotalPages": totalPages,
		},
		"CurrentPage": constant.T("admin.bookings.current_page"),
		"Title":       constant.T("admin.bookings.title"),
	})
}

func (h *AdminBookingHandler) ShowBookingDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateError, gin.H{"error": constant.T("admin.booking.error.invalid_id"), "Title": constant.T("error.title")})
		return
	}

	booking, err := h.service.GetBookingByID(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{"error": err.Error(), "Title": constant.T("error.title")})
		return
	}
	if booking == nil {
		c.HTML(http.StatusNotFound, constant.TemplateError, gin.H{"error": constant.T("admin.booking.error.not_found"), "Title": constant.T("error.title")})
		return
	}

	c.HTML(http.StatusOK, constant.TemplateBookingDetail, gin.H{
		"Booking": booking,
	})
}

func (h *AdminBookingHandler) DeleteBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateError, gin.H{"error": constant.T("admin.booking.error.invalid_id"), "Title": constant.T("error.title")})
		return
	}

	err = h.service.DeleteBooking(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{"error": constant.T("admin.booking.error.delete"), "Title": constant.T("error.title")})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/bookings")
}
