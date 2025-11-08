package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID                      uint           `gorm:"primaryKey" json:"id"`
	MerchantID              string         `gorm:"type:varchar(50);not null" json:"merchant_id"`
	Amount                  float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	Currency                string         `gorm:"type:varchar(3);default:'IDR'" json:"currency"`
	TrxID                   string         `gorm:"type:varchar(100)" json:"trx_id"`
	PartnerReferenceNumber  string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"partner_reference_number"`
	ReferenceNumber         string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"reference_number"`
	Status                  string         `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	TransactionDate         time.Time      `gorm:"not null" json:"transaction_date"`
	PaidDate                *time.Time     `json:"paid_date,omitempty"`
	QRContent               string         `gorm:"type:text" json:"qr_content,omitempty"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	StatusPending = "PENDING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)