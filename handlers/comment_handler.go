package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/services"
	"oe02_go_tam/utils"
	"strings"
)

type CommentHandler struct {
	service services.CommentService
}

type CreateCommentRequest struct {
	ReviewID uint   `json:"review_id" binding:"required"`
	ParentID *uint  `json:"parent_id"`
	Content  string `json:"content" binding:"required"`
}

func NewCommentHandler(service services.CommentService) *CommentHandler {
	return &CommentHandler{service}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content must not be empty or whitespace"})
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

	comment := &models.Comment{
		UserID:   userID,
		ReviewID: req.ReviewID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	err := h.service.CreateComment(comment)
	switch err {
	case nil:
		response := utils.MapCommentToResponse(*comment)
		c.JSON(http.StatusCreated, gin.H{"data": response})
	case constant.ErrReviewNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": constant.T("comment.create.review_not_found")})
	case constant.ErrParentCommentNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": constant.T("comment.create.parent_not_found")})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": constant.T("comment.create.failed")})
	}
}
