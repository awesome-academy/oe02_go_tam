package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string `gorm:"size:100" json:"name"`
	Email      string `gorm:"size:100;uniqueIndex" json:"email"`
	Password   string `gorm:"size:250" json:"password"`
	AvatarURL  string `gorm:"size:250" json:"avatar_url"`
	Role       string `gorm:"type:enum('admin','user')" json:"role"`
	Banned     bool   `gorm:"default:false" json:"banned"`
	FacebookID string `gorm:"size:255" json:"facebook_id"`
	GoogleID   string `gorm:"size:255" json:"google_id"`

	Tours        []Tour        `gorm:"foreignKey:CreatedBy" json:"tours"`
	Bookings     []Booking     `gorm:"foreignKey:UserID" json:"bookings"`
	Reviews      []Review      `gorm:"foreignKey:UserID" json:"reviews"`
	Comments     []Comment     `gorm:"foreignKey:UserID" json:"comments"`
	Likes        []Like        `gorm:"foreignKey:UserID" json:"likes"`
	BankAccounts []BankAccount `gorm:"foreignKey:UserID" json:"bank_accounts"`
}
