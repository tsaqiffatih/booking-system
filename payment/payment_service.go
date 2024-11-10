package payment_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, bookingID, userID string, paymentType string, amount float64) (*PaymentDetails, error)
	CheckPaymentStatus(ctx context.Context, paymentID string) (*Payment, error)
	CancelPayment(ctx context.Context, paymentID string) error
	UpdatePaymentStatus(ctx context.Context, paymentID string, status PaymentStatus) error
}

type paymentService struct {
	repo      PaymentRepository
	midClient coreapi.Client
}

type PaymentDetails struct {
	PaymentID                uuid.UUID
	Status                   PaymentStatus
	RedirectURL              string
	QRCodeURL                string
	BankTransferInstructions *BankTransferInstructions
	ExpiresAt                time.Time
}

type BankTransferInstructions struct {
	Bank           string
	VirtualAccount string
}

const (
	PaymentTypeGopay               = "gopay"
	PaymentTypeBankTransferBCA     = "bank_transfer_bca"
	PaymentTypeBankTransferBNI     = "bank_transfer_bni"
	PaymentTypeBankTransferPermata = "bank_transfer_permata"
	PaymentTypeBankTransferMandiri = "bank_transfer_mandiri"
	PaymentTypeShopeepay           = "shopeepay"
	PaymentTypeQris                = "qris"
)

func NewPaymentService(repo PaymentRepository, serverKey string, isProduction bool) PaymentService {
	midtrans.ServerKey = serverKey
	if isProduction {
		midtrans.Environment = midtrans.Production
	} else {
		midtrans.Environment = midtrans.Sandbox
	}

	client := coreapi.Client{}
	client.New(midtrans.ServerKey, midtrans.Environment)

	return &paymentService{
		repo:      repo,
		midClient: client,
	}
}

// ProcessPayment implements PaymentService.
func (s *paymentService) ProcessPayment(ctx context.Context, bookingID string, userID string, paymentType string, amount float64) (*PaymentDetails, error) {
	paymentID := uuid.New()
	payment := &Payment{
		PaymentID: paymentID,
		BookingID: uuid.MustParse(bookingID),
		UserID:    uuid.MustParse(userID),
		Amount:    amount,
		Status:    StatusPending,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// create request pembayaran
	req, err := s.createChargeRequest(paymentType, paymentID, amount)
	if err != nil {
		return nil, err
	}

	// mengirim request ke midtrans
	resp, midtransError := s.midClient.ChargeTransaction(req)
	if midtransError != nil {
		return nil, fmt.Errorf("failed to charge payment: %w", err)
	}

	// detail pembayaran berdasarkan respons
	paymentDetails := s.populatePaymentDetails(paymentType, resp)
	paymentDetails.PaymentID = paymentID

	// Simpan ke database
	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return paymentDetails, nil
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

// fucntion helper buat request ChargeReq berdasarkan tipe pembayaran
func (s *paymentService) createChargeRequest(paymentType string, paymentID uuid.UUID, amount float64) (*coreapi.ChargeReq, error) {
	switch paymentType {
	case PaymentTypeGopay:
		return &coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeGopay,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  paymentID.String(),
				GrossAmt: int64(amount),
			},
			Gopay: &coreapi.GopayDetails{
				EnableCallback: true,
				CallbackUrl:    "https://your-callback-url.com",
			},
		}, nil
	case PaymentTypeBankTransferBCA:
		return &coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeBankTransfer,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  paymentID.String(),
				GrossAmt: int64(amount),
			},
			BankTransfer: &coreapi.BankTransferDetails{
				Bank: midtrans.BankBca,
			},
		}, nil
	case PaymentTypeBankTransferBNI:
		return &coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeBankTransfer,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  paymentID.String(),
				GrossAmt: int64(amount),
			},
			BankTransfer: &coreapi.BankTransferDetails{
				Bank: midtrans.BankBni,
			},
		}, nil
	case PaymentTypeBankTransferPermata:
		return &coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeBankTransfer,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  paymentID.String(),
				GrossAmt: int64(amount),
			},
			BankTransfer: &coreapi.BankTransferDetails{
				Bank: midtrans.BankPermata,
			},
		}, nil
	case PaymentTypeBankTransferMandiri:
		return &coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeBankTransfer,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  paymentID.String(),
				GrossAmt: int64(amount),
			},
			BankTransfer: &coreapi.BankTransferDetails{
				Bank: midtrans.BankMandiri,
			},
		}, nil
	case PaymentTypeShopeepay:
		return &coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeShopeepay,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  paymentID.String(),
				GrossAmt: int64(amount),
			},
		}, nil
	case PaymentTypeQris:
		return &coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeQris,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  paymentID.String(),
				GrossAmt: int64(amount),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported payment type: %s", paymentType)
	}
}

// function helper mengisi detail pembayaran berdasarkan respons dari Midtrans
func (s *paymentService) populatePaymentDetails(paymentType string, resp *coreapi.ChargeResponse) *PaymentDetails {
	paymentDetails := &PaymentDetails{
		Status: StatusPending,
	}

	switch paymentType {
	case PaymentTypeGopay, PaymentTypeShopeepay:
		if len(resp.Actions) > 0 {
			paymentDetails.RedirectURL = resp.Actions[0].URL
		}
	case PaymentTypeBankTransferBCA, PaymentTypeBankTransferBNI, PaymentTypeBankTransferPermata, PaymentTypeBankTransferMandiri:
		if len(resp.VaNumbers) > 0 {
			paymentDetails.BankTransferInstructions = &BankTransferInstructions{
				Bank:           string(resp.VaNumbers[0].Bank),
				VirtualAccount: resp.VaNumbers[0].VANumber,
			}
		}
	case PaymentTypeQris:
		if len(resp.Actions) > 0 {
			paymentDetails.QRCodeURL = resp.Actions[0].URL
		}
	}
	return paymentDetails
}
