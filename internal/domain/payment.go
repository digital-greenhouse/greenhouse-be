package domain

import (
	"context"
	"time"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentPending  PaymentStatus = "PENDING_VERIFICATION"
	PaymentVerified PaymentStatus = "VERIFIED"
	PaymentRejected PaymentStatus = "REJECTED"

	PaymentMethodTransfer PaymentMethod = "TRANSFERENCIA"
	PaymentMethodCash     PaymentMethod = "EFECTIVO"
)

type Payment struct {
	ID              uint
	BookingID       uint
	Amount          float64
	PaymentMethod   PaymentMethod
	ProofData       string // Base64
	ProofMimeType   string // image/jpeg, etc.
	Status          PaymentStatus
	RejectionReason string
	VerifiedBy      *uint
	PaymentDate     time.Time
	VerifiedAt      *time.Time
}

type PaymentRepository interface {
	Create(ctx context.Context, p *Payment) error
	GetByID(ctx context.Context, id uint) (*Payment, error)
	GetByBookingID(ctx context.Context, bookingID uint) ([]Payment, error)
	UpdateStatus(ctx context.Context, paymentID uint, status PaymentStatus, verifierID *uint, reason string) error
}

type PaymentService interface {
	ProcessPaymentProof(ctx context.Context, bookingID uint, amount float64, method PaymentMethod, proofData string, mimeType string) (*Payment, error)
	VerifyPayment(ctx context.Context, paymentID uint, verifierID uint, status PaymentStatus, reason string) error
	GetPaymentProof(ctx context.Context, paymentID uint, requesterID uint) (*Payment, error)
}
