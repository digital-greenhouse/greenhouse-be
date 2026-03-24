package service

import (
	"context"
	"errors"
	"math"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

type bookingService struct {
	repo         domain.BookingRepository
	propertyRepo domain.PropertyRepository
}

func NewBookingService(repo domain.BookingRepository, propertyRepo domain.PropertyRepository) domain.BookingService {
	return &bookingService{
		repo:         repo,
		propertyRepo: propertyRepo,
	}
}

func (s *bookingService) CalculateQuote(ctx context.Context, propertyID uint, clientID *uint, checkIn, checkOut time.Time, guests int) (*domain.Quote, error) {
	// 1. Validaciones básicas
	if checkOut.Before(checkIn) || checkOut.Equal(checkIn) {
		return nil, errors.New("la fecha de salida debe ser posterior a la de entrada")
	}

	property, err := s.propertyRepo.GetByID(ctx, propertyID)
	if err != nil {
		return nil, errors.New("propiedad no encontrada")
	}

	if guests > property.MaxCapacity {
		return nil, errors.New("la cantidad de huéspedes excede la capacidad máxima")
	}

	// 2. Calcular noches
	nights := int(math.Ceil(checkOut.Sub(checkIn).Hours() / 24))
	
	// 3. Obtener reglas de precio
	rules, err := s.repo.GetPricingRulesByPropertyID(ctx, propertyID, checkIn, checkOut)
	if err != nil {
		return nil, err
	}

	modifier := 1.0
	if len(rules) > 0 {
		// Por simplicidad, tomamos el modificador más alto (o el primero activo)
		modifier = rules[0].PriceModifier
	}

	// 4. Calcular total
	total := float64(nights) * property.BasePricePerNight * modifier

	// 5. Crear cotización
	expiresAt := time.Now().Add(48 * time.Hour) // 48 horas de validez
	quote := &domain.Quote{
		PropertyID:      propertyID,
		ClientID:        clientID,
		CheckInDate:     checkIn,
		CheckOutDate:    checkOut,
		GuestCount:      guests,
		CalculatedTotal: total,
		NightsCount:     nights,
		AppliedModifier: modifier,
		Status:          domain.QuoteActive,
		ExpiresAt:       &expiresAt,
	}

	if err := s.repo.CreateQuote(ctx, quote); err != nil {
		return nil, err
	}

	return quote, nil
}

func (s *bookingService) CreateBookingFromQuote(ctx context.Context, quoteID uint, clientID uint, requests string) (*domain.Booking, error) {
	quote, err := s.repo.GetQuoteByID(ctx, quoteID)
	if err != nil {
		return nil, errors.New("cotización no encontrada")
	}

	if quote.Status != domain.QuoteActive {
		return nil, errors.New("la cotización ya no está activa o ya fue procesada")
	}

	if quote.ExpiresAt != nil && time.Now().After(*quote.ExpiresAt) {
		s.repo.UpdateQuoteStatus(ctx, quoteID, domain.QuoteExpired)
		return nil, errors.New("la cotización ha expirado")
	}

	// Crear reserva
	booking := &domain.Booking{
		PropertyID:      quote.PropertyID,
		ClientID:        clientID,
		QuoteID:         &quote.ID,
		CheckInDate:     quote.CheckInDate,
		CheckOutDate:    quote.CheckOutDate,
		GuestCount:      quote.GuestCount,
		NightsCount:     quote.NightsCount,
		TotalPrice:      quote.CalculatedTotal,
		Status:          domain.BookingPending,
		SpecialRequests: requests,
	}

	if err := s.repo.CreateBooking(ctx, booking); err != nil {
		return nil, err
	}

	// Actualizar estado de cotización
	s.repo.UpdateQuoteStatus(ctx, quoteID, domain.QuoteConverted)

	return booking, nil
}

func (s *bookingService) CreateDirectBooking(ctx context.Context, propertyID uint, clientID uint, checkIn, checkOut time.Time, guests int, requests string) (*domain.Booking, error) {
	// Primero generamos una cotización interna para validar precios
	quote, err := s.CalculateQuote(ctx, propertyID, &clientID, checkIn, checkOut, guests)
	if err != nil {
		return nil, err
	}

	return s.CreateBookingFromQuote(ctx, quote.ID, clientID, requests)
}

func (s *bookingService) CancelBooking(ctx context.Context, bookingID uint, reason string) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return errors.New("reserva no encontrada")
	}

	if booking.Status == domain.BookingCompleted || booking.Status == domain.BookingCancelled {
		return errors.New("no se puede cancelar una reserva completada o ya cancelada")
	}

	return s.repo.UpdateBookingStatus(ctx, bookingID, domain.BookingCancelled, reason)
}

func (s *bookingService) GetClientHistory(ctx context.Context, clientID uint) ([]domain.Booking, error) {
	return s.repo.GetBookingsByClientID(ctx, clientID)
}
