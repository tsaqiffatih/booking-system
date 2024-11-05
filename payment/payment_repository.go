package payment_service

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *Payment) error
	GetPaymentByID(ctx context.Context, paymentID string) (*Payment, error)
	UpdatePayment(ctx context.Context, payment *Payment) error
	CancelExpiredPayments(ctx context.Context) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymenRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// CreatePayment implements PaymentRepository.
func (r *paymentRepository) CreatePayment(ctx context.Context, payment *Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetPaymentByID implements PaymentRepository.
func (r *paymentRepository) GetPaymentByID(ctx context.Context, paymentID string) (*Payment, error) {
	var payment Payment
	err := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).First(&payment).Error
	return &payment, err
}

// UpdatePayment implements PaymentRepository.
func (r *paymentRepository) UpdatePayment(ctx context.Context, payment *Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// CancelExpiredPayments implements PaymentRepository.
func (r *paymentRepository) CancelExpiredPayments(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(Payment{}).
		Where("status = ? AND expires_at < ?", StatusPending, now).
		Update("status", StatusCancelled).Error
}
