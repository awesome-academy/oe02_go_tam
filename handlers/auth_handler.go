package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
)

type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(s services.AuthService) *AuthHandler {
	return &AuthHandler{s}
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Register payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("auth.password.mismatch")})
		return
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// Login godoc
// @Summary Login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constant.T(err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// Logout godoc
// @Summary Logout the user (frontend discards token)
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Frontend should discard token
	c.JSON(http.StatusOK, gin.H{
		"message": constant.T("auth.logout.success"),
	})
}

// GoogleCallback godoc
// @Summary Google OAuth2 login callback
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	googleUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.service.GoogleLogin(googleUser.Name, googleUser.Email, googleUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
