package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/responses"
	"oe02_go_tam/services"
	"oe02_go_tam/utils"
	"strconv"
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
