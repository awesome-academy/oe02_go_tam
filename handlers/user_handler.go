package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{s}
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GetProfile godoc
// @Summary Get current user's profile
// @Description Retrieve the authenticated user's profile details
// @Tags User
// @Produce json
// @Success 200 {object} handlers.UserResponse
// @Failure 400 {object} map[string]string "Invalid user ID format"
// @Failure 404 {object} map[string]string "User not found"
// @Security ApiKeyAuth
// @Router /api/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")

	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constant.T("user.not_found")})
		return
	}

	resp := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}

	c.JSON(http.StatusOK, resp)
}

type UpdateProfileRequest struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// UpdateProfile godoc
// @Summary Update current user's profile
// @Description Update the authenticated user's name and avatar
// @Tags User
// @Accept json
// @Produce json
// @Param body body handlers.UpdateProfileRequest true "Profile update payload"
// @Success 200 {object} handlers.UserResponse
// @Failure 400 {object} map[string]string "Invalid input or user ID format"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /api/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")

	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateProfile(userID, req.Name, req.AvatarURL)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else if errors.Is(err, constant.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		}
		return
	}

	resp := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}

	c.JSON(http.StatusOK, resp)
}
