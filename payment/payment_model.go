package payment_service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "PENDING"
	StatusCompleted PaymentStatus = "COMPLETED"
	StatusFailed    PaymentStatus = "FAILED"
	StatusRefunded  PaymentStatus = "REFUNDED"
	StatusCancelled PaymentStatus = "CANCELLED"
)

type Payment struct {
	PaymentID uuid.UUID     `gorm:"type:uuid;primaryKey" json:"payment_id"`
	BookingID uuid.UUID     `gorm:"type:uuid;not null" json:"booking_id"`
	UserID    uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	Amount    float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status    PaymentStatus `gorm:"type:varchar(20);not null" json:"status"`
	ExpiresAt time.Time     `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.PaymentID == uuid.Nil {
		p.PaymentID = uuid.New()
	}
	return
}
