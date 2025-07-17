package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/responses"
	"oe02_go_tam/services"
	"oe02_go_tam/utils"
	"strconv"
	"strings"
)

type ReviewHandler struct {
	service services.ReviewService
}

func NewReviewHandler(s services.ReviewService) *ReviewHandler {
	return &ReviewHandler{s}
}

func (h *ReviewHandler) GetReviews(c *gin.Context) {
	tourIDStr := c.Param("id")
	tourID, err := strconv.Atoi(tourIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("review.tour_id_invalid")})
		return
	}

	reviews, err := h.service.GetReviews(uint(tourID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constant.T("review.fetch_failed")})
		return
	}

	var response []responses.ReviewResponse
	for _, r := range reviews {
		response = append(response, utils.MapReviewToResponse(r))
	}

	c.JSON(http.StatusOK, response)
}

type ReviewRequest struct {
	TourID  uint   `json:"tour_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("user_id")

	review, err := h.service.CreateReview(userID, req.TourID, req.Rating, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create review failed"})
		return
	}

	reviewResponse := utils.MapReviewToResponse(*review)

	c.JSON(http.StatusCreated, reviewResponse)
}

func (h *ReviewHandler) GetOwnReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("review.id_invalid")})
		return
	}
	userID := c.GetUint("user_id")

	review, err := h.service.GetOwnReview(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constant.T("review.not_found")})
		return
	}

	reviewResponse := utils.MapReviewToResponse(*review)

	c.JSON(http.StatusOK, reviewResponse)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("review.id_invalid")})
		return
	}

	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("review.invalid_input")})
		return
	}

	userID := c.GetUint("user_id")

	review, err := h.service.UpdateReview(uint(id), userID, req.Rating, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constant.T("review.update_failed")})
		return
	}

	reviewResponse := utils.MapReviewToResponse(*review)
	c.JSON(http.StatusOK, reviewResponse)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("review.id_invalid")})
		return
	}

	userID := c.GetUint("user_id")

	if err := h.service.DeleteReview(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constant.T("review.delete_failed")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constant.T("review.delete_success")})
}
