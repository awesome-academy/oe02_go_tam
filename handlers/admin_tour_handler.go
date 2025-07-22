package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/services"
	"strconv"
	"time"
)

type AdminTourHandler struct {
	service services.AdminTourService
}

func NewAdminTourHandler(s services.AdminTourService) *AdminTourHandler {
	return &AdminTourHandler{s}
}

type TourForm struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	Location    string `form:"location"`
	StartDate   string `form:"start_date" binding:"required"` // "2006-01-02"
	EndDate     string `form:"end_date" binding:"required"`
	Price       string `form:"price" binding:"required"` // parse to float64
	Seats       string `form:"seats" binding:"required"` // parse to int
}

func (h *AdminTourHandler) ListTours(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	search := c.Query("search")

	tours, total, err := h.service.GetTours(search, page, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{"error": err.Error(), "Title": constant.T("error.title")})
		return
	}
	totalPages := (int(total) + limit - 1) / limit

	c.HTML(http.StatusOK, constant.TemplateAdminTour, gin.H{
		"Tours":       tours,
		"Pagination":  gin.H{"Page": page, "Limit": limit, "TotalPages": totalPages},
		"Search":      search,
		"Title":       constant.T("admin.tours.title"),
		"CurrentPage": constant.T("admin.tours.current_page"),
	})
}

func (h *AdminTourHandler) ViewTour(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tour, err := h.service.GetTourByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, constant.TemplateError, gin.H{"error": constant.T("tour.not_found"), "Title": constant.T("error.title")})
		return
	}
	c.HTML(http.StatusOK, constant.TemplateTourDetail, gin.H{"Tour": tour, "Title": constant.T("admin.tours.title_detail")})
}

func (h *AdminTourHandler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, constant.TemplateTourCreate, nil)
}

func (h *AdminTourHandler) HandleCreate(c *gin.Context) {
	var form TourForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourCreate, gin.H{"error": err.Error(), "Title": constant.T("admin.tours.title_create")})
		return
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, form.StartDate)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourCreate, gin.H{"error": constant.T("tour.invalid_start_date"), "Title": constant.T("admin.tours.title_create")})
		return
	}
	endDate, err := time.Parse(layout, form.EndDate)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourCreate, gin.H{"error": constant.T("tour.invalid_end_date"), "Title": constant.T("admin.tours.title_create")})
		return
	}

	price, err := strconv.ParseFloat(form.Price, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourCreate, gin.H{"error": constant.T("tour.invalid_price"), "Title": constant.T("admin.tours.title_create")})
		return
	}
	seats, err := strconv.Atoi(form.Seats)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourCreate, gin.H{"error": constant.T("tour.invalid_seats"), "Title": constant.T("admin.tours.title_create")})
		return
	}
	createdBy, exists := c.Get("user_id")
	if !exists {
		c.HTML(http.StatusUnauthorized, constant.TemplateTourCreate, gin.H{"error": constant.T("auth.unauthorized"), "Title": constant.T("admin.tours.title_create")})
		return
	}

	tour := &models.Tour{
		Title:       form.Title,
		Description: form.Description,
		Location:    form.Location,
		StartDate:   startDate,
		EndDate:     endDate,
		Price:       price,
		Seats:       seats,
		CreatedBy:   createdBy.(uint),
	}

	if err := h.service.CreateTour(tour); err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateTourCreate, gin.H{
			"error": err.Error(),
			"Title": constant.T("admin.tours.title_create"),
		})
		return
	}
	c.Redirect(http.StatusFound, "/admin/tours")
}

func (h *AdminTourHandler) ShowEditForm(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tour, err := h.service.GetTourByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, constant.TemplateError, gin.H{
			"error": constant.T("tour.not_found"),
			"Title": constant.T("error.title"),
		})
		return
	}
	c.HTML(http.StatusOK, constant.TemplateTourEdit, gin.H{
		"Tour":  tour,
		"Title": constant.T("admin.tours.title_update"),
	})
}

func (h *AdminTourHandler) HandleUpdate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tour, err := h.service.GetTourByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, constant.TemplateError, gin.H{
			"error": constant.T("tour.not_found"),
			"Title": constant.T("error.title"),
		})
		return
	}

	var form TourForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourEdit, gin.H{
			"Tour":  tour,
			"error": err.Error(),
			"Title": constant.T("admin.tours.title_update"),
		})
		return
	}

	// Parse values
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, form.StartDate)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourEdit, gin.H{"Tour": tour, "error": constant.T("tour.invalid_start_date")})
		return
	}
	endDate, err := time.Parse(layout, form.EndDate)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourEdit, gin.H{
			"Tour":  tour,
			"error": constant.T("tour.invalid_end_date"),
			"Title": constant.T("admin.tours.title_update"),
		})
		return
	}
	price, err := strconv.ParseFloat(form.Price, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourEdit, gin.H{
			"Tour":  tour,
			"error": constant.T("tour.invalid_price"),
			"Title": constant.T("admin.tours.title_update"),
		})
		return
	}
	seats, err := strconv.Atoi(form.Seats)
	if err != nil {
		c.HTML(http.StatusBadRequest, constant.TemplateTourEdit, gin.H{
			"Tour":  tour,
			"error": constant.T("tour.invalid_seats"),
			"Title": constant.T("admin.tours.title_update"),
		})
		return
	}

	tour.Title = form.Title
	tour.Description = form.Description
	tour.Location = form.Location
	tour.StartDate = startDate
	tour.EndDate = endDate
	tour.Price = price
	tour.Seats = seats

	if err := h.service.UpdateTour(tour); err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateTourEdit, gin.H{
			"Tour":  tour,
			"error": err.Error(),
			"Title": constant.T("admin.tours.title_update"),
		})
		return
	}
	c.Redirect(http.StatusFound, "/admin/tours")
}

func (h *AdminTourHandler) DeleteTour(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.service.DeleteTour(uint(id))
	if err != nil {
		c.HTML(http.StatusInternalServerError, constant.TemplateError, gin.H{
			"error": err.Error(),
			"Title": constant.T("error.title"),
		})
		return
	}
	c.Redirect(http.StatusFound, "/admin/tours")
}
