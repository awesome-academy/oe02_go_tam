package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oe02_go_tam/constant"
	"oe02_go_tam/models"
	"oe02_go_tam/responses"
	"oe02_go_tam/services"
	"oe02_go_tam/utils"
	"strconv"
	"time"
)

type TourHandler struct {
	service services.TourService
}

func NewTourHandler(s services.TourService) *TourHandler {
	return &TourHandler{s}
}

//type TourListResponse struct {
//	Title       string    `json:"title"`
//	Description string    `json:"description"`
//	Location    string    `json:"location"`
//	StartDate   time.Time `json:"start_date"`
//	EndDate     time.Time `json:"end_date"`
//	Price       float64   `json:"price"`
//	Seats       int       `json:"seats"`
//}
//
//type TourDetailResponse struct {
//	ID          uint      `json:"id"`
//	Title       string    `json:"title"`
//	Description string    `json:"description"`
//	Location    string    `json:"location"`
//	StartDate   time.Time `json:"start_date"`
//	EndDate     time.Time `json:"end_date"`
//	Price       float64   `json:"price"`
//	Seats       int       `json:"seats"`
//	CreatedBy   uint      `json:"created_by"`
//
//	Creator struct {
//		ID    uint   `json:"id"`
//		Name  string `json:"name"`
//		Email string `json:"email"`
//	} `json:"creator"`
//
//	Bookings []models.Booking `json:"bookings"`
//
//	Reviews []responses.ReviewResponse `json:"reviews"`
//}

type TourBase struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Price       float64   `json:"price"`
	Seats       int       `json:"seats"`
}

type TourListResponse = TourBase

type TourDetailResponse struct {
	ID uint `json:"id"`
	TourBase
	CreatedBy uint `json:"created_by"`
	Creator   struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"creator"`
	Bookings []models.Booking           `json:"bookings"`
	Reviews  []responses.ReviewResponse `json:"reviews"`
}

func (h *TourHandler) ListTours(c *gin.Context) {
	filters := map[string]string{
		"title":       c.Query("title"),
		"location":    c.Query("location"),
		"start_after": c.Query("start_after"),
		"end_before":  c.Query("end_before"),
		"min_price":   c.Query("min_price"),
		"max_price":   c.Query("max_price"),
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	tours, err := h.service.ListToursWithFilters(filters, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constant.T("tour.fetch_failed")})
		return
	}

	var resp []TourListResponse
	for _, t := range tours {
		resp = append(resp, TourListResponse{
			Title:       t.Title,
			Description: t.Description,
			Location:    t.Location,
			StartDate:   t.StartDate,
			EndDate:     t.EndDate,
			Price:       t.Price,
			Seats:       t.Seats,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TourHandler) GetTourDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.T("tour.id_invalid")})
		return
	}

	tour, err := h.service.GetTourDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constant.T("tour.not_found")})
		return
	}

	var resp TourDetailResponse
	resp.ID = tour.ID
	resp.Title = tour.Title
	resp.Description = tour.Description
	resp.Location = tour.Location
	resp.StartDate = tour.StartDate
	resp.EndDate = tour.EndDate
	resp.Price = tour.Price
	resp.Seats = tour.Seats
	resp.CreatedBy = tour.CreatedBy

	resp.Creator.ID = tour.Creator.ID
	resp.Creator.Name = tour.Creator.Name
	resp.Creator.Email = tour.Creator.Email

	resp.Bookings = tour.Bookings

	for _, r := range tour.Reviews {
		resp.Reviews = append(resp.Reviews, utils.MapReviewToResponse(r))
	}

	c.JSON(http.StatusOK, resp)
}
