package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
	"digital-greenhouse/greenhouse-be/internal/http/dto"
	"digital-greenhouse/greenhouse-be/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

type mockBookingServiceForHandler struct {
	CalculateQuoteFunc         func(ctx context.Context, propertyID uint, clientID *uint, checkIn, checkOut time.Time, guests int) (*domain.Quote, error)
	CreateBookingFromQuoteFunc func(ctx context.Context, quoteID uint, clientID uint, requests string) (*domain.Booking, error)
	CreateDirectBookingFunc    func(ctx context.Context, propertyID uint, clientID uint, checkIn, checkOut time.Time, guests int, requests string) (*domain.Booking, error)
	CancelBookingFunc          func(ctx context.Context, bookingID uint, reason string) error
	GetClientHistoryFunc       func(ctx context.Context, clientID uint) ([]domain.Booking, error)
	GetReservedDatesFunc       func(ctx context.Context, propertyID uint) ([]domain.Booking, error)
	GetOwnerBookingsFunc       func(ctx context.Context, ownerID uint) ([]domain.Booking, error)
}

func (m *mockBookingServiceForHandler) CalculateQuote(ctx context.Context, propertyID uint, clientID *uint, checkIn, checkOut time.Time, guests int) (*domain.Quote, error) {
	if m.CalculateQuoteFunc != nil {
		return m.CalculateQuoteFunc(ctx, propertyID, clientID, checkIn, checkOut, guests)
	}
	return nil, nil
}

func (m *mockBookingServiceForHandler) CreateBookingFromQuote(ctx context.Context, quoteID uint, clientID uint, requests string) (*domain.Booking, error) {
	if m.CreateBookingFromQuoteFunc != nil {
		return m.CreateBookingFromQuoteFunc(ctx, quoteID, clientID, requests)
	}
	return nil, nil
}

func (m *mockBookingServiceForHandler) CreateDirectBooking(ctx context.Context, propertyID uint, clientID uint, checkIn, checkOut time.Time, guests int, requests string) (*domain.Booking, error) {
	if m.CreateDirectBookingFunc != nil {
		return m.CreateDirectBookingFunc(ctx, propertyID, clientID, checkIn, checkOut, guests, requests)
	}
	return nil, nil
}

func (m *mockBookingServiceForHandler) CancelBooking(ctx context.Context, bookingID uint, reason string) error {
	if m.CancelBookingFunc != nil {
		return m.CancelBookingFunc(ctx, bookingID, reason)
	}
	return nil
}

func (m *mockBookingServiceForHandler) GetClientHistory(ctx context.Context, clientID uint) ([]domain.Booking, error) {
	if m.GetClientHistoryFunc != nil {
		return m.GetClientHistoryFunc(ctx, clientID)
	}
	return nil, nil
}

func (m *mockBookingServiceForHandler) GetReservedDates(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
	if m.GetReservedDatesFunc != nil {
		return m.GetReservedDatesFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *mockBookingServiceForHandler) GetOwnerBookings(ctx context.Context, ownerID uint) ([]domain.Booking, error) {
	if m.GetOwnerBookingsFunc != nil {
		return m.GetOwnerBookingsFunc(ctx, ownerID)
	}
	return nil, nil
}

func TestBookingHandler_CreateQuote(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		reqBody        string
		mockSetup      func(m *mockBookingServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:   "Happy Path - Guest User",
			userID: 0,
			reqBody: func() string {
				b, _ := json.Marshal(dto.CreateQuoteRequest{
					PropertyID:   1,
					CheckInDate:  time.Now().Add(24 * time.Hour),
					CheckOutDate: time.Now().Add(48 * time.Hour),
					GuestCount:   2,
				})
				return string(b)
			}(),
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.CalculateQuoteFunc = func(ctx context.Context, propertyID uint, clientID *uint, checkIn, checkOut time.Time, guests int) (*domain.Quote, error) {
					if clientID != nil {
						return nil, errors.New("expected clientID to be nil")
					}
					return &domain.Quote{
						ID:              10,
						PropertyID:      propertyID,
						ClientID:        clientID,
						CheckInDate:     checkIn,
						CheckOutDate:    checkOut,
						GuestCount:      guests,
						CalculatedTotal: 250.0,
						NightsCount:     1,
						AppliedModifier: 1.0,
						Status:          domain.QuoteActive,
						CreatedAt:       time.Now(),
					}, nil
				}
			},
			expectedStatus: http.StatusCreated,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.QuoteResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if resp.ID != 10 || resp.CalculatedTotal != 250.0 || resp.ClientID != nil {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:   "Happy Path - Authenticated User",
			userID: 99,
			reqBody: func() string {
				b, _ := json.Marshal(dto.CreateQuoteRequest{
					PropertyID:   1,
					CheckInDate:  time.Now().Add(24 * time.Hour),
					CheckOutDate: time.Now().Add(48 * time.Hour),
					GuestCount:   2,
				})
				return string(b)
			}(),
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.CalculateQuoteFunc = func(ctx context.Context, propertyID uint, clientID *uint, checkIn, checkOut time.Time, guests int) (*domain.Quote, error) {
					if clientID == nil || *clientID != 99 {
						return nil, errors.New("expected clientID to be 99")
					}
					return &domain.Quote{
						ID:              11,
						PropertyID:      propertyID,
						ClientID:        clientID,
						CheckInDate:     checkIn,
						CheckOutDate:    checkOut,
						GuestCount:      guests,
						CalculatedTotal: 250.0,
						NightsCount:     1,
						AppliedModifier: 1.0,
						Status:          domain.QuoteActive,
						CreatedAt:       time.Now(),
					}, nil
				}
			},
			expectedStatus: http.StatusCreated,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.QuoteResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if resp.ID != 11 || resp.ClientID == nil || *resp.ClientID != 99 {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:           "Invalid JSON",
			userID:         0,
			reqBody:        `{"property_id":`,
			mockSetup:      func(m *mockBookingServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid, got: %s", body)
				}
			},
		},
		{
			name:   "Service Error",
			userID: 0,
			reqBody: func() string {
				b, _ := json.Marshal(dto.CreateQuoteRequest{
					PropertyID:   1,
					CheckInDate:  time.Now().Add(24 * time.Hour),
					CheckOutDate: time.Now().Add(48 * time.Hour),
					GuestCount:   2,
				})
				return string(b)
			}(),
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.CalculateQuoteFunc = func(ctx context.Context, propertyID uint, clientID *uint, checkIn, checkOut time.Time, guests int) (*domain.Quote, error) {
					return nil, errors.New("property not available")
				}
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("property not available")) {
					t.Errorf("expected property not available, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookingServiceForHandler{}
			tt.mockSetup(m)
			h := NewBookingHandler(m)

			req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBufferString(tt.reqBody))
			if tt.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}
			rec := httptest.NewRecorder()

			h.CreateQuote(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestBookingHandler_CreateBooking(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		reqBody        string
		mockSetup      func(m *mockBookingServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:   "Happy Path",
			userID: 42,
			reqBody: func() string {
				b, _ := json.Marshal(dto.CreateBookingFromQuoteRequest{
					QuoteID:         10,
					SpecialRequests: "Late check-in",
				})
				return string(b)
			}(),
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.CreateBookingFromQuoteFunc = func(ctx context.Context, quoteID uint, clientID uint, requests string) (*domain.Booking, error) {
					if quoteID != 10 || clientID != 42 || requests != "Late check-in" {
						return nil, errors.New("unexpected params")
					}
					return &domain.Booking{
						ID:              5,
						PropertyID:      1,
						ClientID:        42,
						TotalPrice:      300.0,
						SpecialRequests: requests,
						Status:          domain.BookingPending,
						CreatedAt:       time.Now(),
					}, nil
				}
			},
			expectedStatus: http.StatusCreated,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.BookingResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if resp.ID != 5 || resp.ClientID != 42 || resp.TotalPrice != 300.0 || resp.SpecialRequests != "Late check-in" {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:           "Unauthorized - No User ID",
			userID:         0,
			reqBody:        `{"quote_id":10}`,
			mockSetup:      func(m *mockBookingServiceForHandler) {},
			expectedStatus: http.StatusUnauthorized,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("se requiere autenticación")) {
					t.Errorf("expected unauthorized, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid JSON",
			userID:         42,
			reqBody:        `{"quote_id":`,
			mockSetup:      func(m *mockBookingServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid, got: %s", body)
				}
			},
		},
		{
			name:   "Service Error",
			userID: 42,
			reqBody: func() string {
				b, _ := json.Marshal(dto.CreateBookingFromQuoteRequest{
					QuoteID:         10,
					SpecialRequests: "Late check-in",
				})
				return string(b)
			}(),
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.CreateBookingFromQuoteFunc = func(ctx context.Context, quoteID uint, clientID uint, requests string) (*domain.Booking, error) {
					return nil, errors.New("quote expired")
				}
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("quote expired")) {
					t.Errorf("expected quote expired error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookingServiceForHandler{}
			tt.mockSetup(m)
			h := NewBookingHandler(m)

			req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBufferString(tt.reqBody))
			if tt.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}
			rec := httptest.NewRecorder()

			h.CreateBooking(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestBookingHandler_GetMyHistory(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		mockSetup      func(m *mockBookingServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:   "Happy Path",
			userID: 42,
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.GetClientHistoryFunc = func(ctx context.Context, clientID uint) ([]domain.Booking, error) {
					if clientID != 42 {
						return nil, errors.New("incorrect client ID")
					}
					return []domain.Booking{
						{ID: 1, ClientID: 42, TotalPrice: 100.0, CreatedAt: time.Now()},
						{ID: 2, ClientID: 42, TotalPrice: 200.0, CreatedAt: time.Now()},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp []dto.BookingResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(resp) != 2 || resp[0].ID != 1 || resp[1].ID != 2 {
					t.Errorf("unexpected bookings response: %+v", resp)
				}
			},
		},
		{
			name:           "Unauthorized",
			userID:         0,
			mockSetup:      func(m *mockBookingServiceForHandler) {},
			expectedStatus: http.StatusUnauthorized,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("se requiere autenticación")) {
					t.Errorf("expected unauthorized, got: %s", body)
				}
			},
		},
		{
			name:   "Service Error",
			userID: 42,
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.GetClientHistoryFunc = func(ctx context.Context, clientID uint) ([]domain.Booking, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("db error")) {
					t.Errorf("expected db error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookingServiceForHandler{}
			tt.mockSetup(m)
			h := NewBookingHandler(m)

			req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)
			if tt.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}
			rec := httptest.NewRecorder()

			h.GetMyHistory(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestBookingHandler_GetBookingsByUserID(t *testing.T) {
	tests := []struct {
		name           string
		userIDParam    string
		mockSetup      func(m *mockBookingServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:        "Happy Path",
			userIDParam: "42",
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.GetClientHistoryFunc = func(ctx context.Context, clientID uint) ([]domain.Booking, error) {
					if clientID != 42 {
						return nil, errors.New("incorrect client ID")
					}
					return []domain.Booking{
						{ID: 10, ClientID: 42, TotalPrice: 150.0, CreatedAt: time.Now()},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp []dto.BookingResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(resp) != 1 || resp[0].ID != 10 {
					t.Errorf("unexpected bookings response: %+v", resp)
				}
			},
		},
		{
			name:           "Invalid User ID Param",
			userIDParam:    "abc",
			mockSetup:      func(m *mockBookingServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de usuario inválido")) {
					t.Errorf("expected invalid user id error, got: %s", body)
				}
			},
		},
		{
			name:        "Service Error",
			userIDParam: "42",
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.GetClientHistoryFunc = func(ctx context.Context, clientID uint) ([]domain.Booking, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("db error")) {
					t.Errorf("expected db error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookingServiceForHandler{}
			tt.mockSetup(m)
			h := NewBookingHandler(m)

			r := chi.NewRouter()
			r.Get("/users/{userId}/bookings", h.GetBookingsByUserID)

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userIDParam+"/bookings", nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestBookingHandler_CancelBooking(t *testing.T) {
	tests := []struct {
		name           string
		bookingIDParam string
		reqBody        string
		mockSetup      func(m *mockBookingServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:           "Happy Path",
			bookingIDParam: "5",
			reqBody:        `{"reason":"Change of plans"}`,
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.CancelBookingFunc = func(ctx context.Context, bookingID uint, reason string) error {
					if bookingID != 5 || reason != "Change of plans" {
						return errors.New("unexpected params")
					}
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("reserva cancelada exitosamente")) {
					t.Errorf("expected success message, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid Booking ID",
			bookingIDParam: "abc",
			reqBody:        `{}`,
			mockSetup:      func(m *mockBookingServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de reserva inválido")) {
					t.Errorf("expected invalid booking id error, got: %s", body)
				}
			},
		},
		{
			name:           "Service Error",
			bookingIDParam: "5",
			reqBody:        `{"reason":"Change of plans"}`,
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.CancelBookingFunc = func(ctx context.Context, bookingID uint, reason string) error {
					return errors.New("cannot cancel within 24h")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("cannot cancel within 24h")) {
					t.Errorf("expected cancel error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookingServiceForHandler{}
			tt.mockSetup(m)
			h := NewBookingHandler(m)

			r := chi.NewRouter()
			r.Post("/bookings/{id}/cancel", h.CancelBooking)

			req := httptest.NewRequest(http.MethodPost, "/bookings/"+tt.bookingIDParam+"/cancel", bytes.NewBufferString(tt.reqBody))
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestBookingHandler_GetReservedDates(t *testing.T) {
	tests := []struct {
		name           string
		propertyID     string
		mockSetup      func(m *mockBookingServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:       "Happy Path",
			propertyID: "100",
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.GetReservedDatesFunc = func(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
					if propertyID != 100 {
						return nil, errors.New("incorrect property ID")
					}
					t1, _ := time.Parse("2006-01-02", "2026-06-01")
					t2, _ := time.Parse("2006-01-02", "2026-06-05")
					return []domain.Booking{
						{
							CheckInDate:  t1,
							CheckOutDate: t2,
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp []dto.ReservedDateResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(resp) != 1 || resp[0].CheckInDate != "2026-06-01" || resp[0].CheckOutDate != "2026-06-05" {
					t.Errorf("unexpected reserved dates response: %+v", resp)
				}
			},
		},
		{
			name:           "Invalid Property ID",
			propertyID:     "abc",
			mockSetup:      func(m *mockBookingServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de propiedad inválido")) {
					t.Errorf("expected invalid property id error, got: %s", body)
				}
			},
		},
		{
			name:       "Service Error",
			propertyID: "100",
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.GetReservedDatesFunc = func(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("db error")) {
					t.Errorf("expected db error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookingServiceForHandler{}
			tt.mockSetup(m)
			h := NewBookingHandler(m)

			r := chi.NewRouter()
			r.Get("/properties/{propertyId}/reserved-dates", h.GetReservedDates)

			req := httptest.NewRequest(http.MethodGet, "/properties/"+tt.propertyID+"/reserved-dates", nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestBookingHandler_GetOwnerBookings(t *testing.T) {
	tests := []struct {
		name           string
		ownerID        uint
		mockSetup      func(m *mockBookingServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:    "Happy Path",
			ownerID: 10,
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.GetOwnerBookingsFunc = func(ctx context.Context, ownerID uint) ([]domain.Booking, error) {
					if ownerID != 10 {
						return nil, errors.New("incorrect owner ID")
					}
					return []domain.Booking{
						{ID: 1, ClientID: 42, TotalPrice: 100.0, CreatedAt: time.Now()},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp []dto.BookingResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(resp) != 1 || resp[0].ID != 1 {
					t.Errorf("unexpected bookings response: %+v", resp)
				}
			},
		},
		{
			name:           "Unauthorized",
			ownerID:        0,
			mockSetup:      func(m *mockBookingServiceForHandler) {},
			expectedStatus: http.StatusUnauthorized,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("se requiere autenticación")) {
					t.Errorf("expected unauthorized, got: %s", body)
				}
			},
		},
		{
			name:    "Service Error",
			ownerID: 10,
			mockSetup: func(m *mockBookingServiceForHandler) {
				m.GetOwnerBookingsFunc = func(ctx context.Context, ownerID uint) ([]domain.Booking, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("db error")) {
					t.Errorf("expected db error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookingServiceForHandler{}
			tt.mockSetup(m)
			h := NewBookingHandler(m)

			req := httptest.NewRequest(http.MethodGet, "/bookings/owner", nil)
			if tt.ownerID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.ownerID)
				req = req.WithContext(ctx)
			}
			rec := httptest.NewRecorder()

			h.GetOwnerBookings(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}
