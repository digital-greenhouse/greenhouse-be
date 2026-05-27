package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

// Mock BookingRepository
type mockBookingRepository struct {
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

func (m *mockBookingRepository) CreateQuote(ctx context.Context, quote *domain.Quote) error {
	if m.CreateQuoteFunc != nil {
		return m.CreateQuoteFunc(ctx, quote)
	}
	return nil
}

func (m *mockBookingRepository) GetQuoteByID(ctx context.Context, id uint) (*domain.Quote, error) {
	if m.GetQuoteByIDFunc != nil {
		return m.GetQuoteByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockBookingRepository) GetQuotesByClientID(ctx context.Context, clientID uint) ([]domain.Quote, error) {
	if m.GetQuotesByClientIDFunc != nil {
		return m.GetQuotesByClientIDFunc(ctx, clientID)
	}
	return nil, nil
}

func (m *mockBookingRepository) UpdateQuoteStatus(ctx context.Context, id uint, status domain.QuoteStatus) error {
	if m.UpdateQuoteStatusFunc != nil {
		return m.UpdateQuoteStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *mockBookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	if m.CreateBookingFunc != nil {
		return m.CreateBookingFunc(ctx, booking)
	}
	return nil
}

func (m *mockBookingRepository) GetBookingByID(ctx context.Context, id uint) (*domain.Booking, error) {
	if m.GetBookingByIDFunc != nil {
		return m.GetBookingByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockBookingRepository) GetBookingsByClientID(ctx context.Context, clientID uint) ([]domain.Booking, error) {
	if m.GetBookingsByClientIDFunc != nil {
		return m.GetBookingsByClientIDFunc(ctx, clientID)
	}
	return nil, nil
}

func (m *mockBookingRepository) GetBookingsByPropertyID(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
	if m.GetBookingsByPropertyIDFunc != nil {
		return m.GetBookingsByPropertyIDFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *mockBookingRepository) GetBookingsByOwnerID(ctx context.Context, ownerID uint) ([]domain.Booking, error) {
	if m.GetBookingsByOwnerIDFunc != nil {
		return m.GetBookingsByOwnerIDFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *mockBookingRepository) GetReservedDatesByPropertyID(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
	if m.GetReservedDatesByPropertyIDFunc != nil {
		return m.GetReservedDatesByPropertyIDFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *mockBookingRepository) UpdateBookingStatus(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
	if m.UpdateBookingStatusFunc != nil {
		return m.UpdateBookingStatusFunc(ctx, id, status, reason)
	}
	return nil
}

func (m *mockBookingRepository) CheckAvailability(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
	if m.CheckAvailabilityFunc != nil {
		return m.CheckAvailabilityFunc(ctx, propertyID, checkIn, checkOut)
	}
	return false, nil
}

func (m *mockBookingRepository) CreatePricingRule(ctx context.Context, rule *domain.PricingRule) error {
	if m.CreatePricingRuleFunc != nil {
		return m.CreatePricingRuleFunc(ctx, rule)
	}
	return nil
}

func (m *mockBookingRepository) GetPricingRulesByPropertyID(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error) {
	if m.GetPricingRulesByPropertyIDFunc != nil {
		return m.GetPricingRulesByPropertyIDFunc(ctx, propertyID, start, end)
	}
	return nil, nil
}

func (m *mockBookingRepository) GetAllPricingRulesByPropertyID(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
	if m.GetAllPricingRulesByPropertyIDFunc != nil {
		return m.GetAllPricingRulesByPropertyIDFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *mockBookingRepository) DeletePricingRule(ctx context.Context, id uint) error {
	if m.DeletePricingRuleFunc != nil {
		return m.DeletePricingRuleFunc(ctx, id)
	}
	return nil
}

// Mock PropertyRepository
type mockPropertyRepository struct {
	CreateFunc       func(ctx context.Context, property *domain.Property) error
	GetAllFunc       func(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error)
	GetByIDFunc      func(ctx context.Context, id uint) (*domain.Property, error)
	GetByOwnerIDFunc func(ctx context.Context, ownerID uint) ([]domain.Property, error)
	UpdateFunc       func(ctx context.Context, property *domain.Property) error
	DeleteFunc       func(ctx context.Context, id uint) error
	AddImageFunc     func(ctx context.Context, image *domain.PropertyImage) error
	GetImageByIDFunc func(ctx context.Context, id uint) (*domain.PropertyImage, error)
	UpdateImageFunc  func(ctx context.Context, image *domain.PropertyImage) error
	DeleteImageFunc  func(ctx context.Context, id uint) error
}

func (m *mockPropertyRepository) Create(ctx context.Context, property *domain.Property) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, property)
	}
	return nil
}

func (m *mockPropertyRepository) GetAll(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, filter)
	}
	return nil, nil
}

func (m *mockPropertyRepository) GetByID(ctx context.Context, id uint) (*domain.Property, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockPropertyRepository) GetByOwnerID(ctx context.Context, ownerID uint) ([]domain.Property, error) {
	if m.GetByOwnerIDFunc != nil {
		return m.GetByOwnerIDFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *mockPropertyRepository) Update(ctx context.Context, property *domain.Property) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, property)
	}
	return nil
}

func (m *mockPropertyRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *mockPropertyRepository) AddImage(ctx context.Context, image *domain.PropertyImage) error {
	if m.AddImageFunc != nil {
		return m.AddImageFunc(ctx, image)
	}
	return nil
}

func (m *mockPropertyRepository) GetImageByID(ctx context.Context, id uint) (*domain.PropertyImage, error) {
	if m.GetImageByIDFunc != nil {
		return m.GetImageByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockPropertyRepository) UpdateImage(ctx context.Context, image *domain.PropertyImage) error {
	if m.UpdateImageFunc != nil {
		return m.UpdateImageFunc(ctx, image)
	}
	return nil
}

func (m *mockPropertyRepository) DeleteImage(ctx context.Context, id uint) error {
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(ctx, id)
	}
	return nil
}

func TestCalculateQuote(t *testing.T) {
	now := time.Now()
	checkIn := now.Add(24 * time.Hour)
	checkOut := now.Add(72 * time.Hour) // 2 nights

	clientID := uint(10)

	tests := []struct {
		name          string
		propertyID    uint
		clientID      *uint
		checkIn       time.Time
		checkOut      time.Time
		guests        int
		mockPropRepo  mockPropertyRepository
		mockBookRepo  mockBookingRepository
		wantErr       bool
		expectedErr   string
		verifyQuote   func(t *testing.T, q *domain.Quote)
	}{
		{
			name:       "Error: guests <= 0",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     0,
			wantErr:    true,
			expectedErr: "el número de huéspedes debe ser al menos 1",
		},
		{
			name:       "Error: checkOut before checkIn",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkOut,
			checkOut:   checkIn,
			guests:     2,
			wantErr:    true,
			expectedErr: "la fecha de salida debe ser posterior a la de entrada",
		},
		{
			name:       "Error: checkOut equal checkIn",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkIn,
			guests:     2,
			wantErr:    true,
			expectedErr: "la fecha de salida debe ser posterior a la de entrada",
		},
		{
			name:       "Error: property repo GetByID returns error",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     2,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr:     true,
			expectedErr: "propiedad no encontrada",
		},
		{
			name:       "Error: guests exceeds property MaxCapacity",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     5,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:          1,
						MaxCapacity: 4,
					}, nil
				},
			},
			wantErr:     true,
			expectedErr: "la cantidad de huéspedes excede la capacidad máxima",
		},
		{
			name:       "Error: repo CheckAvailability returns error",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     2,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:          1,
						MaxCapacity: 4,
					}, nil
				},
			},
			mockBookRepo: mockBookingRepository{
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
					return false, errors.New("check availability db error")
				},
			},
			wantErr:     true,
			expectedErr: "check availability db error",
		},
		{
			name:       "Error: property unavailable",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     2,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:          1,
						MaxCapacity: 4,
					}, nil
				},
			},
			mockBookRepo: mockBookingRepository{
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
					return false, nil
				},
			},
			wantErr:     true,
			expectedErr: "la propiedad no está disponible para las fechas seleccionadas",
		},
		{
			name:       "Error: repo GetPricingRulesByPropertyID returns error",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     2,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:          1,
						MaxCapacity: 4,
					}, nil
				},
			},
			mockBookRepo: mockBookingRepository{
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
					return true, nil
				},
				GetPricingRulesByPropertyIDFunc: func(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error) {
					return nil, errors.New("pricing rules db error")
				},
			},
			wantErr:     true,
			expectedErr: "pricing rules db error",
		},
		{
			name:       "Error: repo CreateQuote returns error",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     2,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:                1,
						MaxCapacity:       4,
						BasePricePerNight: 100.0,
					}, nil
				},
			},
			mockBookRepo: mockBookingRepository{
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
					return true, nil
				},
				GetPricingRulesByPropertyIDFunc: func(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error) {
					return nil, nil
				},
				CreateQuoteFunc: func(ctx context.Context, quote *domain.Quote) error {
					return errors.New("create quote db error")
				},
			},
			wantErr:     true,
			expectedErr: "create quote db error",
		},
		{
			name:       "Success with pricing rule applied",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     2,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:                1,
						MaxCapacity:       4,
						BasePricePerNight: 150.0,
					}, nil
				},
			},
			mockBookRepo: mockBookingRepository{
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
					return true, nil
				},
				GetPricingRulesByPropertyIDFunc: func(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error) {
					return []domain.PricingRule{
						{
							PriceModifier: 1.2, // 20% surcharge
						},
					}, nil
				},
				CreateQuoteFunc: func(ctx context.Context, quote *domain.Quote) error {
					quote.ID = 456
					return nil
				},
			},
			wantErr: false,
			verifyQuote: func(t *testing.T, q *domain.Quote) {
				if q.ID != 456 {
					t.Errorf("expected quote ID 456, got %d", q.ID)
				}
				// 2 nights * 2 guests * 150.0 * 1.2 modifier = 720.0
				expectedTotal := 720.0
				if q.CalculatedTotal != expectedTotal {
					t.Errorf("expected quote total %.2f, got %.2f", expectedTotal, q.CalculatedTotal)
				}
				if q.AppliedModifier != 1.2 {
					t.Errorf("expected quote modifier 1.2, got %.2f", q.AppliedModifier)
				}
				if q.NightsCount != 2 {
					t.Errorf("expected nights count 2, got %d", q.NightsCount)
				}
				if q.Status != domain.QuoteActive {
					t.Errorf("expected status ACTIVE, got %v", q.Status)
				}
				if q.ExpiresAt == nil {
					t.Error("expected ExpiresAt to not be nil")
				}
			},
		},
		{
			name:       "Success with no pricing rules (default modifier 1.0)",
			propertyID: 1,
			clientID:   &clientID,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     2,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:                1,
						MaxCapacity:       4,
						BasePricePerNight: 150.0,
					}, nil
				},
			},
			mockBookRepo: mockBookingRepository{
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
					return true, nil
				},
				GetPricingRulesByPropertyIDFunc: func(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error) {
					return nil, nil
				},
				CreateQuoteFunc: func(ctx context.Context, quote *domain.Quote) error {
					quote.ID = 789
					return nil
				},
			},
			wantErr: false,
			verifyQuote: func(t *testing.T, q *domain.Quote) {
				if q.ID != 789 {
					t.Errorf("expected quote ID 789, got %d", q.ID)
				}
				// 2 nights * 2 guests * 150.0 * 1.0 modifier = 600.0
				expectedTotal := 600.0
				if q.CalculatedTotal != expectedTotal {
					t.Errorf("expected quote total %.2f, got %.2f", expectedTotal, q.CalculatedTotal)
				}
				if q.AppliedModifier != 1.0 {
					t.Errorf("expected quote modifier 1.0, got %.2f", q.AppliedModifier)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBookingService(&tt.mockBookRepo, &tt.mockPropRepo)
			q, err := s.CalculateQuote(context.Background(), tt.propertyID, tt.clientID, tt.checkIn, tt.checkOut, tt.guests)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got error: %v", tt.wantErr, err)
			}
			if tt.wantErr {
				if err.Error() != tt.expectedErr {
					t.Errorf("expected error message %q, got %q", tt.expectedErr, err.Error())
				}
				return
			}
			if tt.verifyQuote != nil {
				tt.verifyQuote(t, q)
			}
		})
	}
}

func TestCreateBookingFromQuote(t *testing.T) {
	now := time.Now()
	checkIn := now.Add(24 * time.Hour)
	checkOut := now.Add(72 * time.Hour)
	activeExpiresAt := now.Add(48 * time.Hour)
	expiredExpiresAt := now.Add(-1 * time.Hour)

	tests := []struct {
		name         string
		quoteID      uint
		clientID     uint
		requests     string
		mockBookRepo mockBookingRepository
		wantErr      bool
		expectedErr  string
		verifyBook   func(t *testing.T, b *domain.Booking)
	}{
		{
			name:     "Error: quote not found",
			quoteID:  1,
			clientID: 10,
			requests: "extra pillows",
			mockBookRepo: mockBookingRepository{
				GetQuoteByIDFunc: func(ctx context.Context, id uint) (*domain.Quote, error) {
					return nil, errors.New("not found")
				},
			},
			wantErr:     true,
			expectedErr: "cotización no encontrada",
		},
		{
			name:     "Error: quote status is not ACTIVE",
			quoteID:  1,
			clientID: 10,
			requests: "extra pillows",
			mockBookRepo: mockBookingRepository{
				GetQuoteByIDFunc: func(ctx context.Context, id uint) (*domain.Quote, error) {
					return &domain.Quote{
						ID:     1,
						Status: domain.QuoteConverted,
					}, nil
				},
			},
			wantErr:     true,
			expectedErr: "la cotización ya no está activa o ya fue procesada",
		},
		{
			name:     "Error: quote has expired",
			quoteID:  1,
			clientID: 10,
			requests: "extra pillows",
			mockBookRepo: mockBookingRepository{
				GetQuoteByIDFunc: func(ctx context.Context, id uint) (*domain.Quote, error) {
					return &domain.Quote{
						ID:        1,
						Status:    domain.QuoteActive,
						ExpiresAt: &expiredExpiresAt,
					}, nil
				},
				UpdateQuoteStatusFunc: func(ctx context.Context, id uint, status domain.QuoteStatus) error {
					if id != 1 || status != domain.QuoteExpired {
						t.Errorf("unexpected status update: id %d, status %s", id, status)
					}
					return nil
				},
			},
			wantErr:     true,
			expectedErr: "la cotización ha expirado",
		},
		{
			name:     "Error: CheckAvailability returns error",
			quoteID:  1,
			clientID: 10,
			requests: "extra pillows",
			mockBookRepo: mockBookingRepository{
				GetQuoteByIDFunc: func(ctx context.Context, id uint) (*domain.Quote, error) {
					return &domain.Quote{
						ID:          1,
						PropertyID:  2,
						Status:      domain.QuoteActive,
						ExpiresAt:   &activeExpiresAt,
						CheckInDate: checkIn,
						CheckOutDate: checkOut,
					}, nil
				},
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, start, end time.Time) (bool, error) {
					return false, errors.New("check availability error")
				},
			},
			wantErr:     true,
			expectedErr: "check availability error",
		},
		{
			name:     "Error: property unavailable (someone booked in between)",
			quoteID:  1,
			clientID: 10,
			requests: "extra pillows",
			mockBookRepo: mockBookingRepository{
				GetQuoteByIDFunc: func(ctx context.Context, id uint) (*domain.Quote, error) {
					return &domain.Quote{
						ID:          1,
						PropertyID:  2,
						Status:      domain.QuoteActive,
						ExpiresAt:   &activeExpiresAt,
						CheckInDate: checkIn,
						CheckOutDate: checkOut,
					}, nil
				},
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, start, end time.Time) (bool, error) {
					return false, nil
				},
			},
			wantErr:     true,
			expectedErr: "la propiedad ya no está disponible para estas fechas",
		},
		{
			name:     "Error: CreateBooking returns error",
			quoteID:  1,
			clientID: 10,
			requests: "extra pillows",
			mockBookRepo: mockBookingRepository{
				GetQuoteByIDFunc: func(ctx context.Context, id uint) (*domain.Quote, error) {
					return &domain.Quote{
						ID:          1,
						PropertyID:  2,
						Status:      domain.QuoteActive,
						ExpiresAt:   &activeExpiresAt,
						CheckInDate: checkIn,
						CheckOutDate: checkOut,
					}, nil
				},
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, start, end time.Time) (bool, error) {
					return true, nil
				},
				CreateBookingFunc: func(ctx context.Context, booking *domain.Booking) error {
					return errors.New("create booking error")
				},
			},
			wantErr:     true,
			expectedErr: "create booking error",
		},
		{
			name:     "Success creating booking",
			quoteID:  1,
			clientID: 10,
			requests: "extra pillows",
			mockBookRepo: mockBookingRepository{
				GetQuoteByIDFunc: func(ctx context.Context, id uint) (*domain.Quote, error) {
					return &domain.Quote{
						ID:              1,
						PropertyID:      2,
						Status:          domain.QuoteActive,
						ExpiresAt:       &activeExpiresAt,
						CheckInDate:     checkIn,
						CheckOutDate:    checkOut,
						GuestCount:      3,
						NightsCount:     2,
						CalculatedTotal: 300.0,
					}, nil
				},
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, start, end time.Time) (bool, error) {
					return true, nil
				},
				CreateBookingFunc: func(ctx context.Context, booking *domain.Booking) error {
					booking.ID = 99
					return nil
				},
				UpdateQuoteStatusFunc: func(ctx context.Context, id uint, status domain.QuoteStatus) error {
					if id != 1 || status != domain.QuoteConverted {
						t.Errorf("unexpected status update: id %d, status %s", id, status)
					}
					return nil
				},
			},
			wantErr: false,
			verifyBook: func(t *testing.T, b *domain.Booking) {
				if b.ID != 99 {
					t.Errorf("expected booking ID 99, got %d", b.ID)
				}
				if b.PropertyID != 2 {
					t.Errorf("expected PropertyID 2, got %d", b.PropertyID)
				}
				if b.ClientID != 10 {
					t.Errorf("expected ClientID 10, got %d", b.ClientID)
				}
				if *b.QuoteID != 1 {
					t.Errorf("expected QuoteID 1, got %d", *b.QuoteID)
				}
				if b.Status != domain.BookingPending {
					t.Errorf("expected status Pending, got %s", b.Status)
				}
				if b.SpecialRequests != "extra pillows" {
					t.Errorf("expected special requests 'extra pillows', got %q", b.SpecialRequests)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBookingService(&tt.mockBookRepo, &mockPropertyRepository{})
			b, err := s.CreateBookingFromQuote(context.Background(), tt.quoteID, tt.clientID, tt.requests)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got error: %v", tt.wantErr, err)
			}
			if tt.wantErr {
				if err.Error() != tt.expectedErr {
					t.Errorf("expected error message %q, got %q", tt.expectedErr, err.Error())
				}
				return
			}
			if tt.verifyBook != nil {
				tt.verifyBook(t, b)
			}
		})
	}
}

func TestCreateDirectBooking(t *testing.T) {
	now := time.Now()
	checkIn := now.Add(24 * time.Hour)
	checkOut := now.Add(72 * time.Hour)

	tests := []struct {
		name         string
		propertyID   uint
		clientID     uint
		checkIn      time.Time
		checkOut     time.Time
		guests       int
		requests     string
		mockPropRepo mockPropertyRepository
		mockBookRepo mockBookingRepository
		wantErr      bool
		expectedErr  string
	}{
		{
			name:       "Error: CalculateQuote returns error (invalid guest count)",
			propertyID: 1,
			clientID:   10,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     -1,
			wantErr:    true,
			expectedErr: "el número de huéspedes debe ser al menos 1",
		},
		{
			name:       "Success Direct Booking",
			propertyID: 1,
			clientID:   10,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     2,
			requests:   "direct special request",
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{
						ID:                1,
						MaxCapacity:       4,
						BasePricePerNight: 100.0,
					}, nil
				},
			},
			mockBookRepo: mockBookingRepository{
				CheckAvailabilityFunc: func(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
					return true, nil
				},
				GetPricingRulesByPropertyIDFunc: func(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error) {
					return nil, nil
				},
				CreateQuoteFunc: func(ctx context.Context, quote *domain.Quote) error {
					quote.ID = 123
					return nil
				},
				GetQuoteByIDFunc: func(ctx context.Context, id uint) (*domain.Quote, error) {
					activeExpiresAt := now.Add(48 * time.Hour)
					return &domain.Quote{
						ID:              123,
						PropertyID:      1,
						Status:          domain.QuoteActive,
						ExpiresAt:       &activeExpiresAt,
						CheckInDate:     checkIn,
						CheckOutDate:    checkOut,
						GuestCount:      2,
						NightsCount:     2,
						CalculatedTotal: 400.0,
					}, nil
				},
				CreateBookingFunc: func(ctx context.Context, booking *domain.Booking) error {
					booking.ID = 456
					return nil
				},
				UpdateQuoteStatusFunc: func(ctx context.Context, id uint, status domain.QuoteStatus) error {
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBookingService(&tt.mockBookRepo, &tt.mockPropRepo)
			b, err := s.CreateDirectBooking(context.Background(), tt.propertyID, tt.clientID, tt.checkIn, tt.checkOut, tt.guests, tt.requests)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got error: %v", tt.wantErr, err)
			}
			if tt.wantErr {
				if err.Error() != tt.expectedErr {
					t.Errorf("expected error message %q, got %q", tt.expectedErr, err.Error())
				}
				return
			}
			if b.ID != 456 {
				t.Errorf("expected booking ID 456, got %d", b.ID)
			}
			if b.SpecialRequests != "direct special request" {
				t.Errorf("expected special requests 'direct special request', got %q", b.SpecialRequests)
			}
		})
	}
}

func TestCancelBooking(t *testing.T) {
	tests := []struct {
		name         string
		bookingID    uint
		reason       string
		mockBookRepo mockBookingRepository
		wantErr      bool
		expectedErr  string
	}{
		{
			name:      "Error: booking not found",
			bookingID: 1,
			reason:    "change of plans",
			mockBookRepo: mockBookingRepository{
				GetBookingByIDFunc: func(ctx context.Context, id uint) (*domain.Booking, error) {
					return nil, errors.New("not found")
				},
			},
			wantErr:     true,
			expectedErr: "reserva no encontrada",
		},
		{
			name:      "Error: booking is already completed",
			bookingID: 1,
			reason:    "change of plans",
			mockBookRepo: mockBookingRepository{
				GetBookingByIDFunc: func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:     1,
						Status: domain.BookingCompleted,
					}, nil
				},
			},
			wantErr:     true,
			expectedErr: "no se puede cancelar una reserva completada o ya cancelada",
		},
		{
			name:      "Error: booking is already cancelled",
			bookingID: 1,
			reason:    "change of plans",
			mockBookRepo: mockBookingRepository{
				GetBookingByIDFunc: func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:     1,
						Status: domain.BookingCancelled,
					}, nil
				},
			},
			wantErr:     true,
			expectedErr: "no se puede cancelar una reserva completada o ya cancelada",
		},
		{
			name:      "Error: UpdateBookingStatus returns error",
			bookingID: 1,
			reason:    "change of plans",
			mockBookRepo: mockBookingRepository{
				GetBookingByIDFunc: func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:     1,
						Status: domain.BookingPending,
					}, nil
				},
				UpdateBookingStatusFunc: func(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
					return errors.New("db error")
				},
			},
			wantErr:     true,
			expectedErr: "db error",
		},
		{
			name:      "Success cancelling booking",
			bookingID: 1,
			reason:    "change of plans",
			mockBookRepo: mockBookingRepository{
				GetBookingByIDFunc: func(ctx context.Context, id uint) (*domain.Booking, error) {
					return &domain.Booking{
						ID:     1,
						Status: domain.BookingPending,
					}, nil
				},
				UpdateBookingStatusFunc: func(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
					if id != 1 || status != domain.BookingCancelled || reason != "change of plans" {
						t.Errorf("unexpected parameters: id %d, status %s, reason %q", id, status, reason)
					}
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBookingService(&tt.mockBookRepo, &mockPropertyRepository{})
			err := s.CancelBooking(context.Background(), tt.bookingID, tt.reason)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got error: %v", tt.wantErr, err)
			}
			if tt.wantErr && err.Error() != tt.expectedErr {
				t.Errorf("expected error message %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestGetClientHistory(t *testing.T) {
	tests := []struct {
		name         string
		clientID     uint
		mockBookRepo mockBookingRepository
		wantErr      bool
		expectedLen  int
	}{
		{
			name:     "Error getting client history",
			clientID: 10,
			mockBookRepo: mockBookingRepository{
				GetBookingsByClientIDFunc: func(ctx context.Context, clientID uint) ([]domain.Booking, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
		{
			name:     "Success getting client history",
			clientID: 10,
			mockBookRepo: mockBookingRepository{
				GetBookingsByClientIDFunc: func(ctx context.Context, clientID uint) ([]domain.Booking, error) {
					return []domain.Booking{
						{ID: 1, ClientID: 10},
						{ID: 2, ClientID: 10},
					}, nil
				},
			},
			wantErr:     false,
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBookingService(&tt.mockBookRepo, &mockPropertyRepository{})
			history, err := s.GetClientHistory(context.Background(), tt.clientID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got error: %v", tt.wantErr, err)
			}
			if !tt.wantErr && len(history) != tt.expectedLen {
				t.Errorf("expected history length %d, got %d", tt.expectedLen, len(history))
			}
		})
	}
}

func TestGetReservedDates(t *testing.T) {
	tests := []struct {
		name         string
		propertyID   uint
		mockPropRepo mockPropertyRepository
		mockBookRepo mockBookingRepository
		wantErr      bool
		expectedErr  string
		expectedLen  int
	}{
		{
			name:       "Error: property not found",
			propertyID: 1,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return nil, errors.New("not found")
				},
			},
			wantErr:     true,
			expectedErr: "propiedad no encontrada",
		},
		{
			name:       "Error: GetReservedDatesByPropertyID returns error",
			propertyID: 1,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{ID: 1}, nil
				},
			},
			mockBookRepo: mockBookingRepository{
				GetReservedDatesByPropertyIDFunc: func(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr:     true,
			expectedErr: "db error",
		},
		{
			name:       "Success getting reserved dates",
			propertyID: 1,
			mockPropRepo: mockPropertyRepository{
				GetByIDFunc: func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{ID: 1}, nil
				},
			},
			mockBookRepo: mockBookingRepository{
				GetReservedDatesByPropertyIDFunc: func(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
					return []domain.Booking{
						{ID: 1, PropertyID: 1},
					}, nil
				},
			},
			wantErr:     false,
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBookingService(&tt.mockBookRepo, &tt.mockPropRepo)
			dates, err := s.GetReservedDates(context.Background(), tt.propertyID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got error: %v", tt.wantErr, err)
			}
			if tt.wantErr {
				if err.Error() != tt.expectedErr {
					t.Errorf("expected error message %q, got %q", tt.expectedErr, err.Error())
				}
				return
			}
			if len(dates) != tt.expectedLen {
				t.Errorf("expected reserved bookings length %d, got %d", tt.expectedLen, len(dates))
			}
		})
	}
}

func TestGetOwnerBookings(t *testing.T) {
	tests := []struct {
		name         string
		ownerID      uint
		mockBookRepo mockBookingRepository
		wantErr      bool
		expectedLen  int
	}{
		{
			name:    "Error getting owner bookings",
			ownerID: 5,
			mockBookRepo: mockBookingRepository{
				GetBookingsByOwnerIDFunc: func(ctx context.Context, ownerID uint) ([]domain.Booking, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
		{
			name:    "Success getting owner bookings",
			ownerID: 5,
			mockBookRepo: mockBookingRepository{
				GetBookingsByOwnerIDFunc: func(ctx context.Context, ownerID uint) ([]domain.Booking, error) {
					return []domain.Booking{
						{ID: 10, PropertyID: 1},
					}, nil
				},
			},
			wantErr:     false,
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBookingService(&tt.mockBookRepo, &mockPropertyRepository{})
			bookings, err := s.GetOwnerBookings(context.Background(), tt.ownerID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got error: %v", tt.wantErr, err)
			}
			if !tt.wantErr && len(bookings) != tt.expectedLen {
				t.Errorf("expected bookings length %d, got %d", tt.expectedLen, len(bookings))
			}
		})
	}
}
