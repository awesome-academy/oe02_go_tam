package handlers

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
	"strconv"
)

type AdminRevenueHandler struct {
	service services.AdminRevenueService
}

func NewAdminRevenueHandler(s services.AdminRevenueService) *AdminRevenueHandler {
	return &AdminRevenueHandler{s}
}

func (h *AdminRevenueHandler) ListRevenue(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0")) // 1-12
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))   // >= 2000

	bookings, total, err := h.service.ListRevenue(search, page, limit, month, year)
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{"error": "Failed to load revenue", "Title": constant.T("error.title")})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.HTML(http.StatusOK, constant.TemplateRevenue, gin.H{
		"Bookings": bookings,
		"Search":   search,
		"Pagination": gin.H{
			"Page":       page,
			"Limit":      limit,
			"TotalPages": totalPages,
		},
		"Month":       month,
		"Year":        year,
		"CurrentPage": constant.T("admin.revenue.current_page"),
		"Title":       constant.T("admin.revenue.title"),
	})
}
