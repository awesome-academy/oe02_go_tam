package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
)

type LikeHandler struct {
	service services.LikeService
}

type LikeRequest struct {
	ReviewID uint `json:"review_id" binding:"required"`
}

func NewLikeHandler(service services.LikeService) *LikeHandler {
	return &LikeHandler{service}
}

// LikeReview godoc
// @Summary Like a review
// @Description Like a review by review ID. User can only like once.
// @Tags Like
// @Accept json
// @Produce json
// @Param body body LikeRequest true "Like Review Payload"
// @Success 200 {object} map[string]string "Like success message"
// @Failure 400 {object} map[string]string "Bad request or already liked"
// @Failure 401 {object} map[string]string "Unauthorized (user not logged in)"
// @Failure 404 {object} map[string]string "Review not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/likes [post]
func (h *LikeHandler) LikeReview(c *gin.Context) {
	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	err := h.service.LikeReview(userID, req.ReviewID)
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"message": constant.T("like.review.success")})
	case constant.ErrReviewNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": constant.T("like.review.not_found")})
	case constant.ErrAlreadyLiked:
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("like.review.already_liked")})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": constant.T("like.review.failed")})
	}
}
