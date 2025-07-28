package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/services"
	"strconv"
)

type AdminUsersHandler struct {
	service services.AdminUsersService
}

func NewAdminUsersHandler(service services.AdminUsersService) *AdminUsersHandler {
	return &AdminUsersHandler{service}
}

func (h *AdminUsersHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	search := c.Query("search")

	users, total, err := h.service.GetUsers(search, page, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{
			"error": err.Error(),
			"Title": constant.T("error.title"),
		})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.HTML(http.StatusOK, constant.TemplateUserDashboard, gin.H{
		"Users":       users,
		"Pagination":  gin.H{"Page": page, "Limit": limit, "TotalPages": totalPages},
		"Search":      search,
		"Title":       constant.T("admin.users.title"),
		"CurrentPage": constant.T("admin.users.current_page"),
	})
}

func (h *AdminUsersHandler) ViewUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, constant.TemplateError, gin.H{
			"error": constant.T("user.not_found"),
			"Title": constant.T("error.title"),
		})
		return
	}
	c.HTML(http.StatusOK, constant.TemplateUserDetail, gin.H{
		"User":  user,
		"Title": constant.T("admin.users.detail"),
	})
}

func (h *AdminUsersHandler) ToggleBanUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.service.ToggleBanUser(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{
			"error": err.Error(),
			"Title": constant.T("error.title"),
		})
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}

func (h *AdminUsersHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.HTML(http.StatusBadRequest, constant.TemplateError, gin.H{
			"error": constant.T("user.invalid_id"),
			"Title": constant.T("error.title"),
		})
		return
	}

	err = h.service.DeleteUser(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{
			"error": err.Error(),
			"Title": constant.T("error.title"),
		})
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}
