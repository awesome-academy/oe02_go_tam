package models

import (
	"gorm.io/gorm"
	"time"
)

type Booking struct {
	gorm.Model
	UserID        uint      `json:"user_id"`
	TourID        uint      `json:"tour_id"`
	Status        string    `gorm:"type:varchar(20);check:status IN ('pending','approved','cancelled','completed')" json:"status"`
	NumberOfSeats int       `json:"number_of_seats"`
	TotalPrice    float64   `gorm:"type:decimal(10,2)" json:"total_price"`
	BookingDate   time.Time `json:"booking_date"`

	User User `gorm:"foreignKey:UserID" json:"user"`
	Tour Tour `gorm:"foreignKey:TourID" json:"tour"`
}
