package payment_service

import (
	"context"
	"time"

	"github.com/tsaqiffatih/booking-system/payment/paymentPb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type paymentHandler struct {
	paymentPb.UnimplementedPaymentServiceServer
	service PaymentService
}

func NewPaymentHandler(service PaymentService) paymentPb.PaymentServiceServer {
	return &paymentHandler{service: service}
}

func (s *paymentHandler) ProcessPayment(ctx context.Context, req *paymentPb.PaymentRequest) (*paymentPb.PaymentResponse, error) {
	paymentDetails, err := s.service.ProcessPayment(ctx, req.BookingId, req.UserId, req.PaymentType, float64(req.Amount))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	response := &paymentPb.PaymentResponse{
		PaymentId: paymentDetails.PaymentID.String(),
		BookingId: req.BookingId,
		UserId:    req.UserId,
		Amount:    req.Amount,
		Status:    paymentPb.PaymentStatus_PENDING,
		Details: &paymentPb.PaymentDetails{
			RedirectUrl:    paymentDetails.RedirectURL,
			VirtualAccount: paymentDetails.BankTransferInstructions.VirtualAccount,
			QrCodeUrl:      paymentDetails.QRCodeURL,
		},
		ExpiresAt: paymentDetails.ExpiresAt.Format(time.RFC3339),
	}

	return response, nil
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
