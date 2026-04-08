package service

import (
	"context"
	"errors"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

type paymentService struct {
	repo        domain.PaymentRepository
	bookingRepo domain.BookingRepository
}

func NewPaymentService(repo domain.PaymentRepository, bookingRepo domain.BookingRepository) domain.PaymentService {
	return &paymentService{
		repo:        repo,
		bookingRepo: bookingRepo,
	}
}

func (s *paymentService) ProcessPaymentProof(ctx context.Context, bookingID uint, amount float64, method domain.PaymentMethod, proofData string, mimeType string) (*domain.Payment, error) {
	// 1. Verificar que la reserva existe
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, errors.New("reserva no encontrada")
	}

	// 2. Validar estado de la reserva
	if booking.Status != domain.BookingPending {
		return nil, errors.New("la reserva no está pendiente de pago")
	}

	// 3. Crear registro de pago
	payment := &domain.Payment{
		BookingID:     bookingID,
		Amount:        amount,
		PaymentMethod: method,
		ProofData:     proofData,
		ProofMimeType: mimeType,
		Status:        domain.PaymentPending,
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *paymentService) VerifyPayment(ctx context.Context, paymentID uint, verifierID uint, status domain.PaymentStatus, reason string) error {
	// 1. Obtener el pago
	payment, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		return errors.New("pago no encontrado")
	}

	if payment.Status != domain.PaymentPending {
		return errors.New("el pago ya ha sido procesado")
	}

	// 2. Actualizar estado del pago
	if err := s.repo.UpdateStatus(ctx, paymentID, status, &verifierID, reason); err != nil {
		return err
	}

	// 3. Si el pago es verificado, actualizar la reserva a CONFIRMED
	if status == domain.PaymentVerified {
		if err := s.bookingRepo.UpdateBookingStatus(ctx, payment.BookingID, domain.BookingConfirmed, ""); err != nil {
			return err
		}
	}

	return nil
}
