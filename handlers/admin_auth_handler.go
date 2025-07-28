package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
)

type AdminAuthHandler struct {
	service services.AuthService // reuse user AuthService
}

func NewAdminAuthHandler(s services.AuthService) *AdminAuthHandler {
	return &AdminAuthHandler{s}
}

func (h *AdminAuthHandler) ShowLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, constant.TemplateAdminLogin, gin.H{})
}

func (h *AdminAuthHandler) HandleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	token, _, err := h.service.Login(email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, constant.TemplateAdminLogin, gin.H{
			"Error": constant.T("auth.login.invalid_email_password"),
			"Title": constant.T("auth.header.login_title"),
		})
		return
	}

	// Store in session or cookie (you can use secure cookies or Gin sessions)
	c.SetCookie("admin_token", token, 86400, "/", "", false, true)

	// Redirect to admin dashboard
	c.Redirect(http.StatusFound, "/admin/users")
}

func (h *AdminAuthHandler) HandleLogout(c *gin.Context) {
	c.SetCookie("admin_token", "", -1, "/", "", false, true)

	c.Redirect(http.StatusFound, "/admin/login")
}
