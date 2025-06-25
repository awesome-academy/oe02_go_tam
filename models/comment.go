package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	ReviewID uint   `json:"review_id"`
	ParentID *uint  `json:"parent_id,omitempty"`
	Content  string `json:"content"`

	User    User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Review  Review    `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE" json:"review"`
	Replies []Comment `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"replies"`
}
