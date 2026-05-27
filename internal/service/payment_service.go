package service

import (
	"context"
	"errors"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

type paymentService struct {
	repo         domain.PaymentRepository
	bookingRepo  domain.BookingRepository
	propertyRepo domain.PropertyRepository
}

func NewPaymentService(repo domain.PaymentRepository, bookingRepo domain.BookingRepository, propertyRepo domain.PropertyRepository) domain.PaymentService {
	return &paymentService{
		repo:         repo,
		bookingRepo:  bookingRepo,
		propertyRepo: propertyRepo,
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

	// 2.1 Validar que el monto sea al menos el 50% del total
	if amount < (booking.TotalPrice * 0.5) {
		return nil, errors.New("el pago debe ser de al menos el 50% del total para confirmar la reserva")
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

	// 3. Actualizar la reserva según el estado del pago
	if status == domain.PaymentVerified {
		if err := s.bookingRepo.UpdateBookingStatus(ctx, payment.BookingID, domain.BookingConfirmed, ""); err != nil {
			return err
		}
	} else if status == domain.PaymentRejected {
		// Al rechazar el pago, cancelamos la reserva para liberar las fechas
		reasonStr := "Pago rechazado: " + reason
		if err := s.bookingRepo.UpdateBookingStatus(ctx, payment.BookingID, domain.BookingCancelled, reasonStr); err != nil {
			return err
		}
	}

	return nil
}

func (s *paymentService) GetPaymentProof(ctx context.Context, paymentID uint, requesterID uint) (*domain.Payment, error) {
	payment, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, errors.New("pago no encontrado")
	}

	booking, err := s.bookingRepo.GetBookingByID(ctx, payment.BookingID)
	if err != nil {
		return nil, errors.New("reserva asociada no encontrada")
	}

	// 1. Validar si el solicitante es el cliente
	if booking.ClientID == requesterID {
		return payment, nil
	}

	// 2. Validar si el solicitante es el dueño de la propiedad
	property, err := s.propertyRepo.GetByID(ctx, booking.PropertyID)
	if err != nil {
		return nil, errors.New("propiedad asociada no encontrada")
	}

	if property.OwnerID == requesterID {
		return payment, nil
	}

	return nil, errors.New("no tienes permiso para ver este comprobante")
}
