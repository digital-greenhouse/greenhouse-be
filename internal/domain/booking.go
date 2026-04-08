package domain

import (
	"context"
	"time"
)

type QuoteStatus string
type BookingStatus string

const (
	QuoteActive    QuoteStatus = "ACTIVE"
	QuoteConverted QuoteStatus = "CONVERTED"
	QuoteExpired   QuoteStatus = "EXPIRED"
	QuoteAbandoned QuoteStatus = "ABANDONED"

	BookingPending   BookingStatus = "PENDING_PAYMENT"
	BookingConfirmed BookingStatus = "CONFIRMED"
	BookingCancelled BookingStatus = "CANCELLED"
	BookingCompleted BookingStatus = "COMPLETED"
)

type PricingRule struct {
	ID            uint
	PropertyID    uint
	Name          string
	StartDate     time.Time
	EndDate       time.Time
	PriceModifier float64
	Description   string
	IsActive      bool
	CreatedAt     time.Time
}

type Quote struct {
	ID                uint
	PropertyID        uint
	ClientID          *uint // Puede ser NULL (invitado)
	CheckInDate       time.Time
	CheckOutDate      time.Time
	GuestCount        int
	CalculatedTotal   float64
	NightsCount       int
	AppliedModifier   float64
	Status            QuoteStatus
	AbandonmentReason string
	ExpiresAt         *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Booking struct {
	ID                 uint
	PropertyID         uint
	ClientID           uint
	QuoteID            *uint
	CheckInDate        time.Time
	CheckOutDate       time.Time
	GuestCount         int
	NightsCount        int
	TotalPrice         float64
	Status             BookingStatus
	CancellationReason string
	SpecialRequests    string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type BookingRepository interface {
	// Quotes
	CreateQuote(ctx context.Context, quote *Quote) error
	GetQuoteByID(ctx context.Context, id uint) (*Quote, error)
	GetQuotesByClientID(ctx context.Context, clientID uint) ([]Quote, error)
	UpdateQuoteStatus(ctx context.Context, id uint, status QuoteStatus) error

	// Bookings
	CreateBooking(ctx context.Context, booking *Booking) error
	GetBookingByID(ctx context.Context, id uint) (*Booking, error)
	GetBookingsByClientID(ctx context.Context, clientID uint) ([]Booking, error)
	GetBookingsByPropertyID(ctx context.Context, propertyID uint) ([]Booking, error)
	UpdateBookingStatus(ctx context.Context, id uint, status BookingStatus, reason string) error
	CheckAvailability(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error)

	// Pricing Rules
	GetPricingRulesByPropertyID(ctx context.Context, propertyID uint, start, end time.Time) ([]PricingRule, error)
}

type BookingService interface {
	// Proceso de Cotización
	CalculateQuote(ctx context.Context, propertyID uint, clientID *uint, checkIn, checkOut time.Time, guests int) (*Quote, error)
	
	// Proceso de Reserva
	CreateBookingFromQuote(ctx context.Context, quoteID uint, clientID uint, requests string) (*Booking, error)
	CreateDirectBooking(ctx context.Context, propertyID uint, clientID uint, checkIn, checkOut time.Time, guests int, requests string) (*Booking, error)
	
	// Gestión
	CancelBooking(ctx context.Context, bookingID uint, reason string) error
	GetClientHistory(ctx context.Context, clientID uint) ([]Booking, error)
}
