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
