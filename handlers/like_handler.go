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
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": []interface{}{}})
		return
	}

	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.T("auth.unauthorized"), "data": []interface{}{}})
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.T("auth.user_id.invalid"), "data": []interface{}{}})
		return
	}

	err := h.service.LikeReview(userID, req.ReviewID)
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"message": constant.T("like.review.success"), "data": []interface{}{}})
	case constant.ErrReviewNotFound:
		c.JSON(http.StatusNotFound, gin.H{"message": constant.T("like.review.not_found"), "data": []interface{}{}})
	case constant.ErrAlreadyLiked:
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.T("like.review.already_liked"), "data": []interface{}{}})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.T("like.review.failed"), "data": []interface{}{}})
	}
}
