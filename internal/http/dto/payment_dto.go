package dto

import (
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

type UploadPaymentRequest struct {
	BookingID     uint                 `json:"booking_id"`
	Amount        float64              `json:"amount"`
	PaymentMethod domain.PaymentMethod `json:"payment_method"`
	ProofData     string               `json:"proof_data"` // Base64
	ProofMimeType string               `json:"proof_mime_type"`
}

type VerifyPaymentRequest struct {
	Status          domain.PaymentStatus `json:"status"`
	RejectionReason string               `json:"rejection_reason"`
}

type PaymentResponse struct {
	ID              uint                 `json:"id"`
	BookingID       uint                 `json:"booking_id"`
	Amount          float64              `json:"amount"`
	PaymentMethod   domain.PaymentMethod `json:"payment_method"`
	Status          domain.PaymentStatus `json:"status"`
	RejectionReason string               `json:"rejection_reason,omitempty"`
	PaymentDate     time.Time            `json:"payment_date"`
}

func ToPaymentResponse(p *domain.Payment) PaymentResponse {
	return PaymentResponse{
		ID:              p.ID,
		BookingID:       p.BookingID,
		Amount:          p.Amount,
		PaymentMethod:   p.PaymentMethod,
		Status:          p.Status,
		RejectionReason: p.RejectionReason,
		PaymentDate:     p.PaymentDate,
	}
}
