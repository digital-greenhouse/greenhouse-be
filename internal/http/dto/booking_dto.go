package dto

import (
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

type CreateQuoteRequest struct {
	PropertyID   uint      `json:"property_id"`
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
	GuestCount   int       `json:"guest_count"`
}

type QuoteResponse struct {
	ID              uint               `json:"id"`
	PropertyID      uint               `json:"property_id"`
	ClientID        *uint              `json:"client_id,omitempty"`
	CheckInDate     time.Time          `json:"check_in_date"`
	CheckOutDate    time.Time          `json:"check_out_date"`
	GuestCount      int                `json:"guest_count"`
	CalculatedTotal float64            `json:"calculated_total"`
	NightsCount     int                `json:"nights_count"`
	AppliedModifier float64            `json:"applied_modifier"`
	Status          domain.QuoteStatus `json:"status"`
	ExpiresAt       *time.Time         `json:"expires_at,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
}

type CreateBookingFromQuoteRequest struct {
	QuoteID         uint   `json:"quote_id"`
	SpecialRequests string `json:"special_requests"`
}

type BookingResponse struct {
	ID                 uint                 `json:"id"`
	PropertyID         uint                 `json:"property_id"`
	ClientID           uint                 `json:"client_id"`
	QuoteID            *uint                `json:"quote_id,omitempty"`
	CheckInDate        time.Time            `json:"check_in_date"`
	CheckOutDate       time.Time            `json:"check_out_date"`
	GuestCount         int                  `json:"guest_count"`
	NightsCount        int                  `json:"nights_count"`
	TotalPrice         float64              `json:"total_price"`
	Status             domain.BookingStatus `json:"status"`
	CancellationReason string               `json:"cancellation_reason,omitempty"`
	SpecialRequests    string               `json:"special_requests,omitempty"`
	CreatedAt          time.Time            `json:"created_at"`

	// Datos Dashboard
	ClientName   string `json:"client_name,omitempty"`
	ClientPhone  string `json:"client_phone,omitempty"`
	PropertyName string `json:"property_name,omitempty"`
}

type ReservedDateResponse struct {
	CheckInDate  string `json:"check_in"`
	CheckOutDate string `json:"check_out"`
}

func ToQuoteResponse(q *domain.Quote) QuoteResponse {
	return QuoteResponse{
		ID:              q.ID,
		PropertyID:      q.PropertyID,
		ClientID:        q.ClientID,
		CheckInDate:     q.CheckInDate,
		CheckOutDate:    q.CheckOutDate,
		GuestCount:      q.GuestCount,
		CalculatedTotal: q.CalculatedTotal,
		NightsCount:     q.NightsCount,
		AppliedModifier: q.AppliedModifier,
		Status:          q.Status,
		ExpiresAt:       q.ExpiresAt,
		CreatedAt:       q.CreatedAt,
	}
}

func ToBookingResponse(b *domain.Booking) BookingResponse {
	return BookingResponse{
		ID:                 b.ID,
		PropertyID:         b.PropertyID,
		ClientID:           b.ClientID,
		QuoteID:            b.QuoteID,
		CheckInDate:        b.CheckInDate,
		CheckOutDate:       b.CheckOutDate,
		GuestCount:         b.GuestCount,
		NightsCount:        b.NightsCount,
		TotalPrice:         b.TotalPrice,
		Status:             b.Status,
		CancellationReason: b.CancellationReason,
		SpecialRequests:    b.SpecialRequests,
		CreatedAt:          b.CreatedAt,
		ClientName:         b.ClientName,
		ClientPhone:        b.ClientPhone,
		PropertyName:       b.PropertyName,
	}
}
