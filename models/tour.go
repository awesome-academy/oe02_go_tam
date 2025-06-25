package models

import (
	"gorm.io/gorm"
	"time"
)

type Tour struct {
	gorm.Model
	Title       string    `gorm:"size:200" json:"title"`
	Description string    `json:"description"`
	Location    string    `gorm:"size:200" json:"location"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Price       float64   `gorm:"type:decimal(10,2)" json:"price"`
	Seats       int       `json:"seats"`
	CreatedBy   uint      `json:"created_by"`

	Creator  User      `gorm:"foreignKey:CreatedBy" json:"creator"`
	Bookings []Booking `gorm:"foreignKey:TourID" json:"bookings"`
	Reviews  []Review  `gorm:"foreignKey:TourID" json:"reviews"`
}
