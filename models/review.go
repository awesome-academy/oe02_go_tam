package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	TourID  uint   `json:"tour_id"`
	Rating  int    `json:"rating"`
	Content string `json:"content"`

	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Tour     Tour      `gorm:"foreignKey:TourID;constraint:OnDelete:CASCADE" json:"tour"`
	Comments []Comment `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE" json:"comments"`
	Likes    []Like    `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE" json:"likes"`
}
