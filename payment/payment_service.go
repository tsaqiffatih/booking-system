package payment_service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, bookingID, userID string, amount float64) (*Payment, error)
	CheckPaymentStatus(ctx context.Context, paymentID string) (*Payment, error)
	CancelPayment(ctx context.Context, paymentID string) error
	UpdatePaymentStatus(ctx context.Context, paymentID string, status PaymentStatus) error
}

type paymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

// ProcessPayment implements PaymentService.
func (s *paymentService) ProcessPayment(ctx context.Context, bookingID string, userID string, amount float64) (*Payment, error) {
	payment := &Payment{
		PaymentID: uuid.New(),
		BookingID: uuid.MustParse(bookingID),
		UserID:    uuid.MustParse(userID),
		Amount:    amount,
		Status:    StatusPending,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err := s.repo.CreatePayment(ctx, payment)
	return payment, err
}

// CheckPaymentStatus implements PaymentService.
func (s *paymentService) CheckPaymentStatus(ctx context.Context, paymentID string) (*Payment, error) {
	return s.repo.GetPaymentByID(ctx, paymentID)
}

// CancelPayment implements PaymentService.
func (s *paymentService) CancelPayment(ctx context.Context, paymentID string) error {
	payment, err := s.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return err
	}
	payment.Status = StatusCancelled
	return s.repo.UpdatePayment(ctx, payment)
}

// UpdatePaymentStatus implements PaymentService.
func (s *paymentService) UpdatePaymentStatus(ctx context.Context, paymentID string, status PaymentStatus) error {
	payment, err := s.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return err
	}
	payment.Status = status
	return s.repo.UpdatePayment(ctx, payment)
}
