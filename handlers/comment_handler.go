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

// CreateComment godoc
// @Summary Create a new comment for a review
// @Description Create a comment or reply (nested) on a review. ParentID is optional for nested comments.
// @Tags Comment
// @Accept json
// @Produce json
// @Param body body CreateCommentRequest true "Create Comment Payload"
// @Success 201 {object} map[string]interface{} "Returns the created comment data"
// @Failure 400 {object} map[string]string "Bad request, invalid input or empty content"
// @Failure 401 {object} map[string]string "Unauthorized (user not logged in)"
// @Failure 404 {object} map[string]string "Review or Parent comment not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": []interface{}{}})
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.T("comment.content_empty"), "data": []interface{}{}})
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
		c.JSON(http.StatusCreated, gin.H{"message": constant.T("comment.create.success"), "data": response})
	case constant.ErrReviewNotFound:
		c.JSON(http.StatusNotFound, gin.H{"message": constant.T("comment.create.review_not_found"), "data": []interface{}{}})
	case constant.ErrParentCommentNotFound:
		c.JSON(http.StatusNotFound, gin.H{"message": constant.T("comment.create.parent_not_found"), "data": []interface{}{}})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.T("comment.create.failed"), "data": []interface{}{}})
	}
}
