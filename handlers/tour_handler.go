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
	Reviews []responses.ReviewResponse `json:"reviews"`
}

// ListTours godoc
// @Summary List all tours with optional filters and pagination
// @Description Retrieve all tours with optional filters: title, location, start_after, end_before, min_price, max_price
// @Tags Tour
// @Produce json
// @Param title query string false "Filter by title"
// @Param location query string false "Filter by location"
// @Param start_after query string false "Filter tours starting after this date (YYYY-MM-DD)"
// @Param end_before query string false "Filter tours ending before this date (YYYY-MM-DD)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param page query int false "Page number (default is 1)"
// @Param size query int false "Page size (default is 10)"
// @Success 200 {array} handlers.TourListResponse
// @Failure 500 {object} map[string]string "Failed to fetch tours"
// @Router /api/tours [get]
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.T("tour.fetch_failed"), "data": models.Tour{}})
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

	c.JSON(http.StatusOK, gin.H{"message": constant.T("tour.list_success"), "data": resp})
}

// GetTourDetail godoc
// @Summary Get detailed information about a tour
// @Description Retrieve full details of a tour including bookings and reviews
// @Tags Tour
// @Produce json
// @Param id path int true "Tour ID"
// @Success 200 {object} handlers.TourDetailResponse
// @Failure 400 {object} map[string]string "Invalid tour ID"
// @Failure 404 {object} map[string]string "Tour not found"
// @Router /api/tours/{id} [get]
func (h *TourHandler) GetTourDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.T("tour.id_invalid"), "data": []interface{}{}})
		return
	}

	tour, err := h.service.GetTourDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": constant.T("tour.not_found"), "data": []interface{}{}})
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

	for _, r := range tour.Reviews {
		resp.Reviews = append(resp.Reviews, utils.MapReviewToResponse(r))
	}

	c.JSON(http.StatusOK, gin.H{"message": constant.T("tour.get_success"), "data": resp})
}
