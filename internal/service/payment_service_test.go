package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

// payServiceMockPaymentRepo mocks domain.PaymentRepository
type payServiceMockPaymentRepo struct {
	CreateFunc         func(ctx context.Context, p *domain.Payment) error
	GetByIDFunc        func(ctx context.Context, id uint) (*domain.Payment, error)
	GetByBookingIDFunc func(ctx context.Context, bookingID uint) ([]domain.Payment, error)
	UpdateStatusFunc   func(ctx context.Context, paymentID uint, status domain.PaymentStatus, verifierID *uint, reason string) error
}

func (m *payServiceMockPaymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, p)
	}
	return nil
}

func (m *payServiceMockPaymentRepo) GetByID(ctx context.Context, id uint) (*domain.Payment, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *payServiceMockPaymentRepo) GetByBookingID(ctx context.Context, bookingID uint) ([]domain.Payment, error) {
	if m.GetByBookingIDFunc != nil {
		return m.GetByBookingIDFunc(ctx, bookingID)
	}
	return nil, nil
}

func (m *payServiceMockPaymentRepo) UpdateStatus(ctx context.Context, paymentID uint, status domain.PaymentStatus, verifierID *uint, reason string) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, paymentID, status, verifierID, reason)
	}
	return nil
}

// payServiceMockBookingRepo mocks domain.BookingRepository
type payServiceMockBookingRepo struct {
	CreateQuoteFunc                    func(ctx context.Context, quote *domain.Quote) error
	GetQuoteByIDFunc                   func(ctx context.Context, id uint) (*domain.Quote, error)
	GetQuotesByClientIDFunc            func(ctx context.Context, clientID uint) ([]domain.Quote, error)
	UpdateQuoteStatusFunc              func(ctx context.Context, id uint, status domain.QuoteStatus) error
	CreateBookingFunc                  func(ctx context.Context, booking *domain.Booking) error
	GetBookingByIDFunc                 func(ctx context.Context, id uint) (*domain.Booking, error)
	GetBookingsByClientIDFunc          func(ctx context.Context, clientID uint) ([]domain.Booking, error)
	GetBookingsByPropertyIDFunc        func(ctx context.Context, propertyID uint) ([]domain.Booking, error)
	GetBookingsByOwnerIDFunc           func(ctx context.Context, ownerID uint) ([]domain.Booking, error)
	GetReservedDatesByPropertyIDFunc   func(ctx context.Context, propertyID uint) ([]domain.Booking, error)
	UpdateBookingStatusFunc            func(ctx context.Context, id uint, status domain.BookingStatus, reason string) error
	CheckAvailabilityFunc              func(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error)
	CreatePricingRuleFunc              func(ctx context.Context, rule *domain.PricingRule) error
	GetPricingRulesByPropertyIDFunc    func(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error)
	GetAllPricingRulesByPropertyIDFunc func(ctx context.Context, propertyID uint) ([]domain.PricingRule, error)
	DeletePricingRuleFunc              func(ctx context.Context, id uint) error
}

func (m *payServiceMockBookingRepo) CreateQuote(ctx context.Context, quote *domain.Quote) error {
	if m.CreateQuoteFunc != nil {
		return m.CreateQuoteFunc(ctx, quote)
	}
	return nil
}

func (m *payServiceMockBookingRepo) GetQuoteByID(ctx context.Context, id uint) (*domain.Quote, error) {
	if m.GetQuoteByIDFunc != nil {
		return m.GetQuoteByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *payServiceMockBookingRepo) GetQuotesByClientID(ctx context.Context, clientID uint) ([]domain.Quote, error) {
	if m.GetQuotesByClientIDFunc != nil {
		return m.GetQuotesByClientIDFunc(ctx, clientID)
	}
	return nil, nil
}

func (m *payServiceMockBookingRepo) UpdateQuoteStatus(ctx context.Context, id uint, status domain.QuoteStatus) error {
	if m.UpdateQuoteStatusFunc != nil {
		return m.UpdateQuoteStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *payServiceMockBookingRepo) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	if m.CreateBookingFunc != nil {
		return m.CreateBookingFunc(ctx, booking)
	}
	return nil
}

func (m *payServiceMockBookingRepo) GetBookingByID(ctx context.Context, id uint) (*domain.Booking, error) {
	if m.GetBookingByIDFunc != nil {
		return m.GetBookingByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *payServiceMockBookingRepo) GetBookingsByClientID(ctx context.Context, clientID uint) ([]domain.Booking, error) {
	if m.GetBookingsByClientIDFunc != nil {
		return m.GetBookingsByClientIDFunc(ctx, clientID)
	}
	return nil, nil
}

func (m *payServiceMockBookingRepo) GetBookingsByPropertyID(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
	if m.GetBookingsByPropertyIDFunc != nil {
		return m.GetBookingsByPropertyIDFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *payServiceMockBookingRepo) GetBookingsByOwnerID(ctx context.Context, ownerID uint) ([]domain.Booking, error) {
	if m.GetBookingsByOwnerIDFunc != nil {
		return m.GetBookingsByOwnerIDFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *payServiceMockBookingRepo) GetReservedDatesByPropertyID(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
	if m.GetReservedDatesByPropertyIDFunc != nil {
		return m.GetReservedDatesByPropertyIDFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *payServiceMockBookingRepo) UpdateBookingStatus(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
	if m.UpdateBookingStatusFunc != nil {
		return m.UpdateBookingStatusFunc(ctx, id, status, reason)
	}
	return nil
}

func (m *payServiceMockBookingRepo) CheckAvailability(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
	if m.CheckAvailabilityFunc != nil {
		return m.CheckAvailabilityFunc(ctx, propertyID, checkIn, checkOut)
	}
	return false, nil
}

func (m *payServiceMockBookingRepo) CreatePricingRule(ctx context.Context, rule *domain.PricingRule) error {
	if m.CreatePricingRuleFunc != nil {
		return m.CreatePricingRuleFunc(ctx, rule)
	}
	return nil
}

func (m *payServiceMockBookingRepo) GetPricingRulesByPropertyID(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error) {
	if m.GetPricingRulesByPropertyIDFunc != nil {
		return m.GetPricingRulesByPropertyIDFunc(ctx, propertyID, start, end)
	}
	return nil, nil
}

func (m *payServiceMockBookingRepo) GetAllPricingRulesByPropertyID(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
	if m.GetAllPricingRulesByPropertyIDFunc != nil {
		return m.GetAllPricingRulesByPropertyIDFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *payServiceMockBookingRepo) DeletePricingRule(ctx context.Context, id uint) error {
	if m.DeletePricingRuleFunc != nil {
		return m.DeletePricingRuleFunc(ctx, id)
	}
	return nil
}

// payServiceMockPropertyRepo mocks domain.PropertyRepository
type payServiceMockPropertyRepo struct {
	CreateFunc       func(ctx context.Context, property *domain.Property) error
	GetAllFunc       func(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error)
	GetByIDFunc      func(ctx context.Context, id uint) (*domain.Property, error)
	GetByOwnerIDFunc func(ctx context.Context, ownerID uint) ([]domain.Property, error)
	UpdateFunc       func(ctx context.Context, property *domain.Property) error
	DeleteFunc       func(ctx context.Context, id uint) error

	AddImageFunc      func(ctx context.Context, image *domain.PropertyImage) error
	GetImageByIDFunc  func(ctx context.Context, id uint) (*domain.PropertyImage, error)
	UpdateImageFunc   func(ctx context.Context, image *domain.PropertyImage) error
	DeleteImageFunc   func(ctx context.Context, id uint) error
}

func (m *payServiceMockPropertyRepo) Create(ctx context.Context, property *domain.Property) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, property)
	}
	return nil
}

func (m *payServiceMockPropertyRepo) GetAll(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, filter)
	}
	return nil, nil
}

func (m *payServiceMockPropertyRepo) GetByID(ctx context.Context, id uint) (*domain.Property, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *payServiceMockPropertyRepo) GetByOwnerID(ctx context.Context, ownerID uint) ([]domain.Property, error) {
	if m.GetByOwnerIDFunc != nil {
		return m.GetByOwnerIDFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *payServiceMockPropertyRepo) Update(ctx context.Context, property *domain.Property) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, property)
	}
	return nil
}

func (m *payServiceMockPropertyRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *payServiceMockPropertyRepo) AddImage(ctx context.Context, image *domain.PropertyImage) error {
	if m.AddImageFunc != nil {
		return m.AddImageFunc(ctx, image)
	}
	return nil
}

func (m *payServiceMockPropertyRepo) GetImageByID(ctx context.Context, id uint) (*domain.PropertyImage, error) {
	if m.GetImageByIDFunc != nil {
		return m.GetImageByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *payServiceMockPropertyRepo) UpdateImage(ctx context.Context, image *domain.PropertyImage) error {
	if m.UpdateImageFunc != nil {
		return m.UpdateImageFunc(ctx, image)
	}
	return nil
}

func (m *payServiceMockPropertyRepo) DeleteImage(ctx context.Context, id uint) error {
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(ctx, id)
	}
	return nil
}

func TestProcessPaymentProof(t *testing.T) {
	tests := []struct {
		name          string
		bookingID     uint
		amount        float64
		method        domain.PaymentMethod
		proofData     string
		mimeType      string
		mockBooking   func(mock *payServiceMockBookingRepo)
		mockPayment   func(mock *payServiceMockPaymentRepo)
		expectedError string
		checkPayment  func(t *testing.T, p *domain.Payment)
	}{
		{
			name:      "booking not found",
			bookingID: 1,
			amount:    100.0,
			method:    domain.PaymentMethodTransfer,
			proofData: "base64data",
			mimeType:  "image/png",
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return nil, errors.New("db error")
				}
			},
			mockPayment:   func(mock *payServiceMockPaymentRepo) {},
			expectedError: "reserva no encontrada",
		},
		{
			name:      "status is not PENDING_PAYMENT",
			bookingID: 1,
			amount:    100.0,
			method:    domain.PaymentMethodTransfer,
			proofData: "base64data",
			mimeType:  "image/png",
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:         1,
						Status:     domain.BookingConfirmed, // Not pending_payment
						TotalPrice: 200.0,
					}, nil
				}
			},
			mockPayment:   func(mock *payServiceMockPaymentRepo) {},
			expectedError: "la reserva no está pendiente de pago",
		},
		{
			name:      "amount is less than 50 percent of total price",
			bookingID: 1,
			amount:    99.9, // 50% is 100.0
			method:    domain.PaymentMethodTransfer,
			proofData: "base64data",
			mimeType:  "image/png",
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:         1,
						Status:     domain.BookingPending,
						TotalPrice: 200.0,
					}, nil
				}
			},
			mockPayment:   func(mock *payServiceMockPaymentRepo) {},
			expectedError: "el pago debe ser de al menos el 50% del total",
		},
		{
			name:      "repository error on Create",
			bookingID: 1,
			amount:    100.0,
			method:    domain.PaymentMethodTransfer,
			proofData: "base64data",
			mimeType:  "image/png",
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:         1,
						Status:     domain.BookingPending,
						TotalPrice: 200.0,
					}, nil
				}
			},
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.CreateFunc = func(ctx context.Context, p *domain.Payment) error {
					return errors.New("failed to insert")
				}
			},
			expectedError: "failed to insert",
		},
		{
			name:      "success exactly 50 percent",
			bookingID: 1,
			amount:    100.0,
			method:    domain.PaymentMethodTransfer,
			proofData: "base64data",
			mimeType:  "image/png",
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:         1,
						Status:     domain.BookingPending,
						TotalPrice: 200.0,
					}, nil
				}
			},
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.CreateFunc = func(ctx context.Context, p *domain.Payment) error {
					p.ID = 42
					return nil
				}
			},
			expectedError: "",
			checkPayment: func(t *testing.T, p *domain.Payment) {
				if p == nil {
					t.Fatal("expected payment not to be nil")
				}
				if p.ID != 42 {
					t.Errorf("expected payment ID 42, got %d", p.ID)
				}
				if p.BookingID != 1 {
					t.Errorf("expected BookingID 1, got %d", p.BookingID)
				}
				if p.Amount != 100.0 {
					t.Errorf("expected Amount 100.0, got %f", p.Amount)
				}
				if p.Status != domain.PaymentPending {
					t.Errorf("expected Status %s, got %s", domain.PaymentPending, p.Status)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mBook := &payServiceMockBookingRepo{}
			mPay := &payServiceMockPaymentRepo{}
			mProp := &payServiceMockPropertyRepo{}
			tt.mockBooking(mBook)
			tt.mockPayment(mPay)

			service := NewPaymentService(mPay, mBook, mProp)
			res, err := service.ProcessPaymentProof(context.Background(), tt.bookingID, tt.amount, tt.method, tt.proofData, tt.mimeType)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedError)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error %q to contain %q", err.Error(), tt.expectedError)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if tt.checkPayment != nil {
					tt.checkPayment(t, res)
				}
			}
		})
	}
}

func TestVerifyPayment(t *testing.T) {
	tests := []struct {
		name          string
		paymentID     uint
		verifierID    uint
		status        domain.PaymentStatus
		reason        string
		mockPayment   func(mock *payServiceMockPaymentRepo)
		mockBooking   func(mock *payServiceMockBookingRepo)
		expectedError string
	}{
		{
			name:       "payment not found",
			paymentID:  1,
			verifierID: 9,
			status:     domain.PaymentVerified,
			reason:     "",
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return nil, errors.New("db error")
				}
			},
			mockBooking:   func(mock *payServiceMockBookingRepo) {},
			expectedError: "pago no encontrado",
		},
		{
			name:       "payment status is not PENDING_VERIFICATION",
			paymentID:  1,
			verifierID: 9,
			status:     domain.PaymentVerified,
			reason:     "",
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{
						ID:     1,
						Status: domain.PaymentVerified, // already verified
					}, nil
				}
			},
			mockBooking:   func(mock *payServiceMockBookingRepo) {},
			expectedError: "el pago ya ha sido procesado",
		},
		{
			name:       "repository status update failure",
			paymentID:  1,
			verifierID: 9,
			status:     domain.PaymentVerified,
			reason:     "",
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{
						ID:        1,
						BookingID: 10,
						Status:    domain.PaymentPending, // matches PENDING_VERIFICATION
					}, nil
				}
				mock.UpdateStatusFunc = func(ctx context.Context, paymentID uint, status domain.PaymentStatus, verifierID *uint, reason string) error {
					return errors.New("failed to update status")
				}
			},
			mockBooking:   func(mock *payServiceMockBookingRepo) {},
			expectedError: "failed to update status",
		},
		{
			name:       "booking status update failure - verified",
			paymentID:  1,
			verifierID: 9,
			status:     domain.PaymentVerified,
			reason:     "",
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{
						ID:        1,
						BookingID: 10,
						Status:    domain.PaymentPending,
					}, nil
				}
				mock.UpdateStatusFunc = func(ctx context.Context, paymentID uint, status domain.PaymentStatus, verifierID *uint, reason string) error {
					return nil
				}
			},
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.UpdateBookingStatusFunc = func(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
					return errors.New("failed booking status update")
				}
			},
			expectedError: "failed booking status update",
		},
		{
			name:       "booking status update failure - rejected",
			paymentID:  1,
			verifierID: 9,
			status:     domain.PaymentRejected,
			reason:     "bad photo",
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{
						ID:        1,
						BookingID: 10,
						Status:    domain.PaymentPending,
					}, nil
				}
				mock.UpdateStatusFunc = func(ctx context.Context, paymentID uint, status domain.PaymentStatus, verifierID *uint, reason string) error {
					return nil
				}
			},
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.UpdateBookingStatusFunc = func(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
					return errors.New("failed booking status cancel")
				}
			},
			expectedError: "failed booking status cancel",
		},
		{
			name:       "success path for verified payment (booking confirmed)",
			paymentID:  1,
			verifierID: 9,
			status:     domain.PaymentVerified,
			reason:     "",
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{
						ID:        1,
						BookingID: 10,
						Status:    domain.PaymentPending,
					}, nil
				}
				mock.UpdateStatusFunc = func(ctx context.Context, paymentID uint, status domain.PaymentStatus, verifierID *uint, reason string) error {
					if paymentID != 1 || status != domain.PaymentVerified || *verifierID != 9 || reason != "" {
						return errors.New("incorrect status update params")
					}
					return nil
				}
			},
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.UpdateBookingStatusFunc = func(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
					if id != 10 || status != domain.BookingConfirmed || reason != "" {
						return errors.New("incorrect booking update params")
					}
					return nil
				}
			},
			expectedError: "",
		},
		{
			name:       "success path for rejected payment (booking cancelled)",
			paymentID:  1,
			verifierID: 9,
			status:     domain.PaymentRejected,
			reason:     "invalid screenshot",
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{
						ID:        1,
						BookingID: 10,
						Status:    domain.PaymentPending,
					}, nil
				}
				mock.UpdateStatusFunc = func(ctx context.Context, paymentID uint, status domain.PaymentStatus, verifierID *uint, reason string) error {
					if paymentID != 1 || status != domain.PaymentRejected || *verifierID != 9 || reason != "invalid screenshot" {
						return errors.New("incorrect status update params")
					}
					return nil
				}
			},
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.UpdateBookingStatusFunc = func(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
					expectedReason := "Pago rechazado: invalid screenshot"
					if id != 10 || status != domain.BookingCancelled || reason != expectedReason {
						return errors.New("incorrect booking update params")
					}
					return nil
				}
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mBook := &payServiceMockBookingRepo{}
			mPay := &payServiceMockPaymentRepo{}
			mProp := &payServiceMockPropertyRepo{}
			tt.mockBooking(mBook)
			tt.mockPayment(mPay)

			service := NewPaymentService(mPay, mBook, mProp)
			err := service.VerifyPayment(context.Background(), tt.paymentID, tt.verifierID, tt.status, tt.reason)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedError)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error %q to contain %q", err.Error(), tt.expectedError)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestGetPaymentProof(t *testing.T) {
	tests := []struct {
		name          string
		paymentID     uint
		requesterID   uint
		mockPayment   func(mock *payServiceMockPaymentRepo)
		mockBooking   func(mock *payServiceMockBookingRepo)
		mockProperty  func(mock *payServiceMockPropertyRepo)
		expectedError string
		checkPayment  func(t *testing.T, p *domain.Payment)
	}{
		{
			name:        "payment not found",
			paymentID:   1,
			requesterID: 100,
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return nil, errors.New("db error")
				}
			},
			mockBooking:   func(mock *payServiceMockBookingRepo) {},
			mockProperty:  func(mock *payServiceMockPropertyRepo) {},
			expectedError: "pago no encontrado",
		},
		{
			name:        "booking not found",
			paymentID:   1,
			requesterID: 100,
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{ID: 1, BookingID: 10}, nil
				}
			},
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return nil, errors.New("db error")
				}
			},
			mockProperty:  func(mock *payServiceMockPropertyRepo) {},
			expectedError: "reserva asociada no encontrada",
		},
		{
			name:        "property not found",
			paymentID:   1,
			requesterID: 100, // not client (which is 200)
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{ID: 1, BookingID: 10}, nil
				}
			},
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:         10,
						ClientID:   200,
						PropertyID: 300,
					}, nil
				}
			},
			mockProperty: func(mock *payServiceMockPropertyRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Property, error) {
					return nil, errors.New("db error")
				}
			},
			expectedError: "propiedad asociada no encontrada",
		},
		{
			name:        "permission denied (requester is not client or owner)",
			paymentID:   1,
			requesterID: 400, // Client is 200, Owner is 300, Requester is 400
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{ID: 1, BookingID: 10}, nil
				}
			},
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:         10,
						ClientID:   200,
						PropertyID: 50,
					}, nil
				}
			},
			mockProperty: func(mock *payServiceMockPropertyRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:      50,
						OwnerID: 300,
					}, nil
				}
			},
			expectedError: "no tienes permiso para ver este comprobante",
		},
		{
			name:        "success path for client",
			paymentID:   1,
			requesterID: 200, // Client is 200
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{ID: 1, BookingID: 10, ProofData: "proven"}, nil
				}
			},
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:         10,
						ClientID:   200,
						PropertyID: 50,
					}, nil
				}
			},
			mockProperty: func(mock *payServiceMockPropertyRepo) {}, // shouldn't be called because requester is client
			expectedError: "",
			checkPayment: func(t *testing.T, p *domain.Payment) {
				if p == nil {
					t.Fatal("expected payment not to be nil")
				}
				if p.ProofData != "proven" {
					t.Errorf("expected proof data 'proven', got %q", p.ProofData)
				}
			},
		},
		{
			name:        "success path for owner",
			paymentID:   1,
			requesterID: 300, // Owner is 300, Client is 200
			mockPayment: func(mock *payServiceMockPaymentRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Payment, error) {
					return &domain.Payment{ID: 1, BookingID: 10, ProofData: "proven"}, nil
				}
			},
			mockBooking: func(mock *payServiceMockBookingRepo) {
				mock.GetBookingByIDFunc = func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:         10,
						ClientID:   200,
						PropertyID: 50,
					}, nil
				}
			},
			mockProperty: func(mock *payServiceMockPropertyRepo) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:      50,
						OwnerID: 300,
					}, nil
				}
			},
			expectedError: "",
			checkPayment: func(t *testing.T, p *domain.Payment) {
				if p == nil {
					t.Fatal("expected payment not to be nil")
				}
				if p.ProofData != "proven" {
					t.Errorf("expected proof data 'proven', got %q", p.ProofData)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mBook := &payServiceMockBookingRepo{}
			mPay := &payServiceMockPaymentRepo{}
			mProp := &payServiceMockPropertyRepo{}
			tt.mockBooking(mBook)
			tt.mockPayment(mPay)
			tt.mockProperty(mProp)

			service := NewPaymentService(mPay, mBook, mProp)
			res, err := service.GetPaymentProof(context.Background(), tt.paymentID, tt.requesterID)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedError)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error %q to contain %q", err.Error(), tt.expectedError)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if tt.checkPayment != nil {
					tt.checkPayment(t, res)
				}
			}
		})
	}
}
