package models

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	BookingID       uint       `gorm:"uniqueIndex" json:"booking_id"`
	Amount          float64    `gorm:"type:decimal(10,2)" json:"amount"`
	Method          string     `gorm:"type:varchar(20);check:method IN ('vnpay','momo','paypal')" json:"method"`
	Status          string     `gorm:"type:varchar(20);check:status IN ('pending','paid','failed')" json:"status"`
	VnpTxnRef       string     `gorm:"size:255" json:"vnp_txn_ref"`
	BankCode        string     `gorm:"size:50" json:"bank_code"`
	PayDate         *time.Time `json:"pay_date"`
	TransactionTime *time.Time `json:"transaction_time"`

	Booking Booking `gorm:"foreignKey:BookingID;references:ID;constraint:OnDelete:CASCADE" json:"booking"`
}
