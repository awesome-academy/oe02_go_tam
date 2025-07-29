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

// GetReviews godoc
// @Summary Get reviews for a tour
// @Description Retrieve all reviews of a specific tour by tour ID
// @Tags Review
// @Produce json
// @Param id path int true "Tour ID"
// @Success 200 {array} responses.ReviewResponse
// @Failure 400 {object} map[string]string "Invalid tour ID"
// @Failure 500 {object} map[string]string "Failed to fetch reviews"
// @Router /api/tours/{id}/reviews [get]
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

// CreateReview godoc
// @Summary Create a new review
// @Description Create a review for a tour. Requires authentication.
// @Tags Review
// @Accept json
// @Produce json
// @Param body body ReviewRequest true "Review payload"
// @Success 201 {object} responses.ReviewResponse
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Create review failed"
// @Security ApiKeyAuth
// @Router /api/reviews [post]
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

// GetOwnReview godoc
// @Summary Get own review by ID
// @Description Get the review written by the authenticated user by review ID
// @Tags Review
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} responses.ReviewResponse
// @Failure 400 {object} map[string]string "Invalid review ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Review not found"
// @Security ApiKeyAuth
// @Router /api/reviews/{id} [get]
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

// UpdateReview godoc
// @Summary Update own review
// @Description Update a review by ID owned by the authenticated user
// @Tags Review
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param body body ReviewRequest true "Updated review payload"
// @Success 200 {object} responses.ReviewResponse
// @Failure 400 {object} map[string]string "Invalid input or review ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Update failed"
// @Security ApiKeyAuth
// @Router /api/reviews/{id} [put]
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

// DeleteReview godoc
// @Summary Delete own review
// @Description Delete a review by ID owned by the authenticated user
// @Tags Review
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} map[string]string "Delete success message"
// @Failure 400 {object} map[string]string "Invalid review ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Delete failed"
// @Security ApiKeyAuth
// @Router /api/reviews/{id} [delete]
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
