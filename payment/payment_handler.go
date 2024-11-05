package payment_service

import (
	"context"

	"github.com/tsaqiffatih/booking-system/payment/paymentPb"
)

type paymentHandler struct {
	paymentPb.UnimplementedPaymentServiceServer
	service PaymentService
}

func NewPaymentHandler(service PaymentService) paymentPb.PaymentServiceServer {
	return &paymentHandler{service: service}
}

func (s *paymentHandler) ProcessPayment(ctx context.Context, req *paymentPb.PaymentRequest) (*paymentPb.PaymentResponse, error) {
	payment, err := s.service.ProcessPayment(ctx, req.BookingId, req.UserId, float64(req.Amount))
	if err != nil {
		return nil, err
	}
	return &paymentPb.PaymentResponse{
		PaymentId: payment.PaymentID.String(),
		BookingId: payment.BookingID.String(),
		UserId:    payment.UserID.String(),
		Amount:    float32(payment.Amount),
		Status:    paymentPb.PaymentStatus(paymentPb.PaymentStatus_value[string(payment.Status)]),
		CreatedAt: payment.CreatedAt.String(),
		UpdatedAt: payment.UpdatedAt.String(),
		ExpiresAt: payment.ExpiresAt.String(),
	}, nil
}

func (s *paymentHandler) CheckPaymentStatus(ctx context.Context, req *paymentPb.PaymentStatusRequest) (*paymentPb.PaymentResponse, error) {
	payment, err := s.service.CheckPaymentStatus(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}
	return &paymentPb.PaymentResponse{
		PaymentId: payment.PaymentID.String(),
		BookingId: payment.BookingID.String(),
		UserId:    payment.UserID.String(),
		Amount:    float32(payment.Amount),
		Status:    paymentPb.PaymentStatus(paymentPb.PaymentStatus_value[string(payment.Status)]),
		CreatedAt: payment.CreatedAt.String(),
		UpdatedAt: payment.UpdatedAt.String(),
		ExpiresAt: payment.ExpiresAt.String(),
	}, nil
}

func (s *paymentHandler) CancelPayment(ctx context.Context, req *paymentPb.PaymentStatusRequest) (*paymentPb.PaymentResponse, error) {
	err := s.service.CancelPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}
	return &paymentPb.PaymentResponse{
		PaymentId: req.PaymentId,
		Status:    paymentPb.PaymentStatus_CANCELLED,
	}, nil
}

func (s *paymentHandler) HandlePaymentCallback(ctx context.Context, req *paymentPb.PaymentCallbackRequest) (*paymentPb.PaymentResponse, error) {
	err := s.service.UpdatePaymentStatus(ctx, req.PaymentId, PaymentStatus(req.Status.String()))
	if err != nil {
		return nil, err
	}
	return &paymentPb.PaymentResponse{
		PaymentId: req.PaymentId,
		Status:    req.Status,
	}, nil
}
