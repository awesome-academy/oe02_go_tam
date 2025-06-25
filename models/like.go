package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID   uint `gorm:"uniqueIndex:idx_user_review" json:"user_id"`
	ReviewID uint `gorm:"uniqueIndex:idx_user_review" json:"review_id"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Review Review `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE" json:"review"`
}
