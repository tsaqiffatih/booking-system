package payment_service

import (
	"context"
	"fmt"

	"github.com/tsaqiffatih/booking-system/payment/paymentPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type PaymentClient struct {
	conn           *grpc.ClientConn
	servicePayment paymentPb.PaymentServiceClient
}

func NewPaymentClient(url string) (*PaymentClient, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	// conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := paymentPb.NewPaymentServiceClient(conn)
	return &PaymentClient{conn: conn, servicePayment: client}, nil
}

func (pc *PaymentClient) Close() {
	pc.conn.Close()
}

func (pc *PaymentClient) ProcessPayment(ctx context.Context, bookingID, userID, paymentType string, amount float64) (*paymentPb.PaymentResponse, error) {
	req := &paymentPb.PaymentRequest{
		BookingId:   bookingID,
		UserId:      userID,
		PaymentType: paymentType,
		Amount:      float32(amount),
	}

	resp, err := pc.servicePayment.ProcessPayment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	return resp, nil
}

func (pc *PaymentClient) CheckPaymentStatus(ctx context.Context, paymentID string) (*paymentPb.PaymentResponse, error) {
	req := &paymentPb.PaymentStatusRequest{
		PaymentId: paymentID,
	}

	resp, err := pc.servicePayment.CheckPaymentStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check payment status: %w", err)
	}

	return resp, nil
}

func (pc *PaymentClient) CancelPayment(ctx context.Context, paymentID string) (*paymentPb.PaymentResponse, error) {
	req := &paymentPb.PaymentStatusRequest{
		PaymentId: paymentID,
	}

	resp, err := pc.servicePayment.CancelPayment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment: %w", err)
	}

	return resp, nil
}
