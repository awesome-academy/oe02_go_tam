package handlers

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
	"strconv"
)

type AdminReviewHandler struct {
	service services.AdminReviewService
}

func NewAdminReviewHandler(s services.AdminReviewService) *AdminReviewHandler {
	return &AdminReviewHandler{s}
}

func (h *AdminReviewHandler) ListReviews(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.service.GetAllReviews(search, page, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{"error": "Failed to load reviews", "Title": constant.T("error.title")})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.HTML(http.StatusOK, constant.TemplateReview, gin.H{
		"Reviews": reviews,
		"Search":  search,
		"Pagination": gin.H{
			"Page":       page,
			"Limit":      limit,
			"TotalPages": totalPages,
		},
		"CurrentPage": constant.T("admin.reviews.current_page"),
		"Title":       constant.T("admin.reviews.title"),
	})
}

func (h *AdminReviewHandler) ShowReviewDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateError, gin.H{"error": "Invalid review ID", "Title": constant.T("error.title")})
		return
	}

	review, err := h.service.GetReviewByID(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{"error": err.Error(), "Title": constant.T("error.title")})
		return
	}
	if review == nil {
		c.HTML(http.StatusNotFound, constant.TemplateError, gin.H{"error": "Review not found", "Title": constant.T("error.title")})
		return
	}

	c.HTML(http.StatusOK, constant.TemplateReviewDetail, gin.H{
		"Review": review,
		"Title":  constant.T("admin.reviews.detail_title"),
	})
}

func (h *AdminReviewHandler) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateError, gin.H{"error": "Invalid review ID", "Title": constant.T("error.title")})
		return
	}

	err = h.service.DeleteReview(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{"error": "Failed to delete review", "Title": constant.T("error.title")})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/reviews")
}
