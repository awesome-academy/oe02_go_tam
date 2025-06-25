package models

import "gorm.io/gorm"

type BankAccount struct {
	gorm.Model
	UserID        uint   `gorm:"uniqueIndex:idx_user_account"`
	BankName      string `gorm:"size:100;uniqueIndex:idx_user_account"`
	AccountNumber string `gorm:"size:50;uniqueIndex:idx_user_account"`
	OwnerName     string `gorm:"size:100"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}
