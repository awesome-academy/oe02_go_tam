package responses

import "time"

type BookingResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	TourID        uint      `json:"tour_id"`
	Status        string    `json:"status"`
	NumberOfSeats int       `json:"number_of_seats"`
	TotalPrice    float64   `json:"total_price"`
	BookingDate   time.Time `json:"booking_date"`
}
