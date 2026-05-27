package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestBookingRepository_CreateQuote(t *testing.T) {
	clientID := uint(5)
	now := time.Now()
	expires := now.Add(24 * time.Hour)

	tests := []struct {
		name          string
		quote         *domain.Quote
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			quote: &domain.Quote{
				PropertyID:      10,
				ClientID:        &clientID,
				CheckInDate:     now,
				CheckOutDate:    now.Add(48 * time.Hour),
				GuestCount:      2,
				CalculatedTotal: 250.50,
				NightsCount:     2,
				AppliedModifier: 1.0,
				Status:          domain.QuoteActive,
				ExpiresAt:       &expires,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `quotes`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(100, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			quote: &domain.Quote{
				PropertyID: 10,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `quotes`").
					WillReturnError(errors.New("insert quote failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("insert quote failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			err = repo.CreateQuote(context.Background(), tt.quote)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.quote.ID != 100 {
					t.Errorf("expected quote ID to be updated to 100, got %d", tt.quote.ID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_GetQuoteByID(t *testing.T) {
	clientID := uint(5)
	now := time.Now()
	expires := now.Add(24 * time.Hour)

	tests := []struct {
		name          string
		id            uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedQuote *domain.Quote
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "property_id", "client_id", "check_in_date", "check_out_date", "guest_count", "calculated_total", "nights_count", "applied_modifier", "status", "abandonment_reason", "expires_at", "created_at", "updated_at"}).
					AddRow(uint(1), uint(10), &clientID, now, now.Add(48*time.Hour), 2, 250.50, 2, 1.0, "ACTIVE", "", &expires, now, now)
				mock.ExpectQuery("^SELECT \\* FROM `quotes` WHERE `quotes`.`id` = \\?").
					WithArgs(uint(1), 1).
					WillReturnRows(rows)
			},
			wantErr: false,
			expectedQuote: &domain.Quote{
				ID:              1,
				PropertyID:      10,
				ClientID:        &clientID,
				GuestCount:      2,
				CalculatedTotal: 250.50,
				NightsCount:     2,
				AppliedModifier: 1.0,
				Status:          domain.QuoteActive,
			},
		},
		{
			name: "record not found",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `quotes` WHERE `quotes`.`id` = \\?").
					WithArgs(uint(2), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr:       true,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			quote, err := repo.GetQuoteByID(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err != tt.expectedError && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if quote == nil {
					t.Errorf("expected quote not to be nil")
				} else if quote.ID != tt.expectedQuote.ID || quote.PropertyID != tt.expectedQuote.PropertyID || *quote.ClientID != *tt.expectedQuote.ClientID {
					t.Errorf("returned quote %+v does not match expected %+v", quote, tt.expectedQuote)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_GetQuotesByClientID(t *testing.T) {
	clientID := uint(5)
	now := time.Now()
	expires := now.Add(24 * time.Hour)

	tests := []struct {
		name          string
		clientID      uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedLen   int
	}{
		{
			name:     "success",
			clientID: 5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "property_id", "client_id", "check_in_date", "check_out_date", "guest_count", "calculated_total", "nights_count", "applied_modifier", "status", "abandonment_reason", "expires_at", "created_at", "updated_at"}).
					AddRow(uint(1), uint(10), &clientID, now, now.Add(48*time.Hour), 2, 250.50, 2, 1.0, "ACTIVE", "", &expires, now, now).
					AddRow(uint(2), uint(11), &clientID, now, now.Add(48*time.Hour), 1, 150.00, 2, 1.0, "CONVERTED", "", &expires, now, now)
				mock.ExpectQuery("^SELECT \\* FROM `quotes` WHERE client_id = \\?").
					WithArgs(uint(5)).
					WillReturnRows(rows)
			},
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name:     "database error",
			clientID: 6,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `quotes` WHERE client_id = \\?").
					WithArgs(uint(6)).
					WillReturnError(errors.New("db query failed"))
			},
			wantErr:       true,
			expectedError: errors.New("db query failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			quotes, err := repo.GetQuotesByClientID(context.Background(), tt.clientID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(quotes) != tt.expectedLen {
					t.Errorf("expected %d quotes, got %d", tt.expectedLen, len(quotes))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_UpdateQuoteStatus(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		status        domain.QuoteStatus
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name:   "success",
			id:     1,
			status: domain.QuoteConverted,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `quotes` SET `status`=\\?,`updated_at`=\\? WHERE id = \\?").
					WithArgs("CONVERTED", sqlmock.AnyArg(), uint(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:   "database error",
			id:     2,
			status: domain.QuoteAbandoned,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `quotes` SET `status`=\\?,`updated_at`=\\? WHERE id = \\?").
					WillReturnError(errors.New("update status failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("update status failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			err = repo.UpdateQuoteStatus(context.Background(), tt.id, tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_CreateBooking(t *testing.T) {
	quoteID := uint(100)
	now := time.Now()

	tests := []struct {
		name          string
		booking       *domain.Booking
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			booking: &domain.Booking{
				PropertyID:      10,
				ClientID:        5,
				QuoteID:         &quoteID,
				CheckInDate:     now,
				CheckOutDate:    now.Add(48 * time.Hour),
				GuestCount:      2,
				NightsCount:     2,
				TotalPrice:      250.50,
				Status:          domain.BookingPending,
				SpecialRequests: "Non-smoking room please",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `bookings`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(200, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			booking: &domain.Booking{
				PropertyID: 10,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `bookings`").
					WillReturnError(errors.New("insert booking failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("insert booking failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			err = repo.CreateBooking(context.Background(), tt.booking)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.booking.ID != 200 {
					t.Errorf("expected booking ID to be updated to 200, got %d", tt.booking.ID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_GetBookingByID(t *testing.T) {
	quoteID := uint(100)
	now := time.Now()
	paymentID := uint(999)
	paymentStatus := domain.PaymentVerified

	tests := []struct {
		name            string
		id              uint
		mockSetup       func(mock sqlmock.Sqlmock)
		wantErr         bool
		expectedError   error
		expectedBooking *domain.Booking
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "property_id", "client_id", "quote_id", "check_in_date", "check_out_date",
					"guest_count", "nights_count", "total_price", "status", "cancellation_reason",
					"special_requests", "created_at", "updated_at", "payment_id", "payment_status",
				}).AddRow(
					uint(1), uint(10), uint(5), &quoteID, now, now.Add(48*time.Hour),
					2, 2, 250.50, "PENDING_PAYMENT", "",
					"requests", now, now, &paymentID, &paymentStatus,
				)

				mock.ExpectQuery("^SELECT bookings.* FROM `bookings` LEFT JOIN payments p ON p.booking_id = bookings.id AND p.id = \\(SELECT MAX\\(id\\) FROM payments WHERE booking_id = bookings.id\\) WHERE bookings.id = \\? .* LIMIT \\?").
					WithArgs(uint(1), 1).
					WillReturnRows(rows)
			},
			wantErr: false,
			expectedBooking: &domain.Booking{
				ID:            1,
				PropertyID:    10,
				ClientID:      5,
				QuoteID:       &quoteID,
				NightsCount:   2,
				TotalPrice:    250.50,
				Status:        domain.BookingPending,
				PaymentID:     &paymentID,
				PaymentStatus: &paymentStatus,
			},
		},
		{
			name: "record not found",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT bookings.* FROM `bookings`").
					WithArgs(uint(2), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr:       true,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			booking, err := repo.GetBookingByID(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err != tt.expectedError && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if booking == nil {
					t.Errorf("expected booking not to be nil")
				} else if booking.ID != tt.expectedBooking.ID || booking.PropertyID != tt.expectedBooking.PropertyID || *booking.PaymentStatus != *tt.expectedBooking.PaymentStatus {
					t.Errorf("returned booking %+v does not match expected %+v", booking, tt.expectedBooking)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_GetBookingsByClientID(t *testing.T) {
	clientID := uint(5)
	quoteID := uint(100)
	now := time.Now()
	paymentID := uint(999)
	paymentStatus := domain.PaymentVerified

	tests := []struct {
		name          string
		clientID      uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedLen   int
	}{
		{
			name:     "success",
			clientID: 5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "property_id", "client_id", "quote_id", "check_in_date", "check_out_date",
					"guest_count", "nights_count", "total_price", "status", "cancellation_reason",
					"special_requests", "created_at", "updated_at", "payment_id", "payment_status",
				}).AddRow(
					uint(1), uint(10), clientID, &quoteID, now, now.Add(48*time.Hour),
					2, 2, 250.50, "PENDING_PAYMENT", "",
					"requests", now, now, &paymentID, &paymentStatus,
				)

				mock.ExpectQuery("^SELECT bookings.* FROM `bookings` LEFT JOIN payments p ON p.booking_id = bookings.id AND p.id = \\(SELECT MAX\\(id\\) FROM payments WHERE booking_id = bookings.id\\) WHERE bookings.client_id = \\? ORDER BY bookings.created_at DESC").
					WithArgs(clientID).
					WillReturnRows(rows)
			},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:     "database error",
			clientID: 6,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT bookings.* FROM `bookings`").
					WithArgs(uint(6)).
					WillReturnError(errors.New("select failed"))
			},
			wantErr:       true,
			expectedError: errors.New("select failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			bookings, err := repo.GetBookingsByClientID(context.Background(), tt.clientID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(bookings) != tt.expectedLen {
					t.Errorf("expected %d bookings, got %d", tt.expectedLen, len(bookings))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_GetBookingsByPropertyID(t *testing.T) {
	propertyID := uint(10)
	quoteID := uint(100)
	now := time.Now()

	tests := []struct {
		name          string
		propertyID    uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedLen   int
	}{
		{
			name:       "success",
			propertyID: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "property_id", "client_id", "quote_id", "check_in_date", "check_out_date",
					"guest_count", "nights_count", "total_price", "status", "cancellation_reason",
					"special_requests", "created_at", "updated_at",
				}).AddRow(
					uint(1), propertyID, uint(5), &quoteID, now, now.Add(48*time.Hour),
					2, 2, 250.50, "PENDING_PAYMENT", "",
					"requests", now, now,
				)

				mock.ExpectQuery("^SELECT \\* FROM `bookings` WHERE property_id = \\?").
					WithArgs(propertyID).
					WillReturnRows(rows)
			},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:       "database error",
			propertyID: 11,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `bookings` WHERE property_id = \\?").
					WithArgs(uint(11)).
					WillReturnError(errors.New("query failed"))
			},
			wantErr:       true,
			expectedError: errors.New("query failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			bookings, err := repo.GetBookingsByPropertyID(context.Background(), tt.propertyID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(bookings) != tt.expectedLen {
					t.Errorf("expected %d bookings, got %d", tt.expectedLen, len(bookings))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_GetBookingsByOwnerID(t *testing.T) {
	ownerID := uint(2)
	quoteID := uint(100)
	now := time.Now()
	paymentID := uint(999)
	paymentStatus := domain.PaymentVerified

	tests := []struct {
		name          string
		ownerID       uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedLen   int
	}{
		{
			name:    "success",
			ownerID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "property_id", "client_id", "quote_id", "check_in_date", "check_out_date",
					"guest_count", "nights_count", "total_price", "status", "cancellation_reason",
					"special_requests", "created_at", "updated_at", "client_name", "client_phone",
					"property_name", "payment_id", "payment_status",
				}).AddRow(
					uint(1), uint(10), uint(5), &quoteID, now, now.Add(48*time.Hour),
					2, 2, 250.50, "PENDING_PAYMENT", "",
					"requests", now, now, "John Client", "555-1234",
					"Nice Villa", &paymentID, &paymentStatus,
				)

				mock.ExpectQuery("^SELECT bookings.*, users.name as client_name, users.phone as client_phone, properties.name as property_name, p.id as payment_id, p.status as payment_status FROM `bookings` JOIN properties ON bookings.property_id = properties.id JOIN users ON bookings.client_id = users.id LEFT JOIN payments p ON p.booking_id = bookings.id AND p.id = \\(SELECT MAX\\(id\\) FROM payments WHERE booking_id = bookings.id\\) WHERE properties.owner_id = \\? ORDER BY bookings.created_at DESC").
					WithArgs(ownerID).
					WillReturnRows(rows)
			},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:    "database error",
			ownerID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT bookings.* FROM `bookings`").
					WithArgs(uint(3)).
					WillReturnError(errors.New("scan failed"))
			},
			wantErr:       true,
			expectedError: errors.New("scan failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			bookings, err := repo.GetBookingsByOwnerID(context.Background(), tt.ownerID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(bookings) != tt.expectedLen {
					t.Errorf("expected %d bookings, got %d", tt.expectedLen, len(bookings))
				}
				if len(bookings) > 0 {
					if bookings[0].ClientName != "John Client" || bookings[0].PropertyName != "Nice Villa" {
						t.Errorf("scanned join fields incorrectly: %+v", bookings[0])
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_GetReservedDatesByPropertyID(t *testing.T) {
	propertyID := uint(10)
	now := time.Now()

	tests := []struct {
		name          string
		propertyID    uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedLen   int
	}{
		{
			name:       "success",
			propertyID: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"check_in_date", "check_out_date"}).
					AddRow(now, now.Add(48*time.Hour))

				mock.ExpectQuery("^SELECT check_in_date, check_out_date FROM `bookings` WHERE property_id = \\? AND status IN \\(\\?, \\?\\) AND check_out_date >= \\? ORDER BY check_in_date ASC").
					WithArgs(propertyID, "CONFIRMED", "PENDING_PAYMENT", sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:       "database error",
			propertyID: 11,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT check_in_date, check_out_date FROM `bookings`").
					WithArgs(uint(11), "CONFIRMED", "PENDING_PAYMENT", sqlmock.AnyArg()).
					WillReturnError(errors.New("reserved dates failed"))
			},
			wantErr:       true,
			expectedError: errors.New("reserved dates failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			bookings, err := repo.GetReservedDatesByPropertyID(context.Background(), tt.propertyID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(bookings) != tt.expectedLen {
					t.Errorf("expected %d bookings, got %d", tt.expectedLen, len(bookings))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_UpdateBookingStatus(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		status        domain.BookingStatus
		reason        string
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name:   "success without reason",
			id:     1,
			status: domain.BookingConfirmed,
			reason: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `bookings` SET `status`=\\?,`updated_at`=\\? WHERE id = \\?").
					WithArgs("CONFIRMED", sqlmock.AnyArg(), uint(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:   "success with reason",
			id:     1,
			status: domain.BookingCancelled,
			reason: "User cancelled",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `bookings` SET .* WHERE id = \\?").
					WithArgs("User cancelled", "CANCELLED", sqlmock.AnyArg(), uint(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:   "database error",
			id:     2,
			status: domain.BookingConfirmed,
			reason: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `bookings` SET `status`=\\?,`updated_at`=\\? WHERE id = \\?").
					WillReturnError(errors.New("update booking status failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("update booking status failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			err = repo.UpdateBookingStatus(context.Background(), tt.id, tt.status, tt.reason)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_CheckAvailability(t *testing.T) {
	propertyID := uint(10)
	now := time.Now()

	tests := []struct {
		name          string
		propertyID    uint
		checkIn       time.Time
		checkOut      time.Time
		mockSetup     func(mock sqlmock.Sqlmock)
		wantAvail     bool
		wantErr       bool
		expectedError error
	}{
		{
			name:       "available (count = 0)",
			propertyID: propertyID,
			checkIn:    now,
			checkOut:   now.Add(48 * time.Hour),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(int64(0))
				mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `bookings` WHERE property_id = \\? AND status IN \\(\\?, \\?\\) AND check_in_date < \\? AND check_out_date > \\?").
					WithArgs(propertyID, "CONFIRMED", "PENDING_PAYMENT", now.Add(48*time.Hour), now).
					WillReturnRows(rows)
			},
			wantAvail: true,
			wantErr:   false,
		},
		{
			name:       "not available (count > 0)",
			propertyID: propertyID,
			checkIn:    now,
			checkOut:   now.Add(48 * time.Hour),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(int64(1))
				mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `bookings` WHERE property_id = \\? AND status IN \\(\\?, \\?\\) AND check_in_date < \\? AND check_out_date > \\?").
					WithArgs(propertyID, "CONFIRMED", "PENDING_PAYMENT", now.Add(48*time.Hour), now).
					WillReturnRows(rows)
			},
			wantAvail: false,
			wantErr:   false,
		},
		{
			name:       "database error",
			propertyID: propertyID,
			checkIn:    now,
			checkOut:   now.Add(48 * time.Hour),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `bookings` WHERE property_id = \\? AND status IN \\(\\?, \\?\\) AND check_in_date < \\? AND check_out_date > \\?").
					WillReturnError(errors.New("count query failed"))
			},
			wantAvail:     false,
			wantErr:       true,
			expectedError: errors.New("count query failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			avail, err := repo.CheckAvailability(context.Background(), tt.propertyID, tt.checkIn, tt.checkOut)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if avail != tt.wantAvail {
					t.Errorf("expected availability %t, got %t", tt.wantAvail, avail)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_CreatePricingRule(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		rule          *domain.PricingRule
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			rule: &domain.PricingRule{
				PropertyID:    10,
				Name:          "Holiday Season",
				StartDate:     now,
				EndDate:       now.Add(240 * time.Hour),
				PriceModifier: 1.25,
				Description:   "Increase price",
				IsActive:      true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `pricing_rules`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(50, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			rule: &domain.PricingRule{
				PropertyID: 10,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `pricing_rules`").
					WillReturnError(errors.New("insert pricing rule failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("insert pricing rule failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			err = repo.CreatePricingRule(context.Background(), tt.rule)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.rule.ID != 50 {
					t.Errorf("expected pricing rule ID to be updated to 50, got %d", tt.rule.ID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_GetPricingRulesByPropertyID(t *testing.T) {
	propertyID := uint(10)
	now := time.Now()

	tests := []struct {
		name          string
		propertyID    uint
		start         time.Time
		end           time.Time
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedLen   int
	}{
		{
			name:       "success",
			propertyID: 10,
			start:      now,
			end:        now.Add(48 * time.Hour),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "property_id", "name", "start_date", "end_date", "price_modifier", "description", "is_active", "created_at"}).
					AddRow(uint(1), propertyID, "Rule 1", now, now.Add(48*time.Hour), 1.2, "Desc", true, now)

				mock.ExpectQuery("^SELECT \\* FROM `pricing_rules` WHERE property_id = \\? AND is_active = TRUE AND \\(\\(start_date <= \\? AND end_date >= \\?\\) OR \\(start_date <= \\? AND end_date >= \\?\\)\\)").
					WithArgs(propertyID, now.Add(48*time.Hour), now, now.Add(48*time.Hour), now).
					WillReturnRows(rows)
			},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:       "database error",
			propertyID: 11,
			start:      now,
			end:        now.Add(48 * time.Hour),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `pricing_rules` WHERE property_id = \\?").
					WithArgs(uint(11), now.Add(48*time.Hour), now, now.Add(48*time.Hour), now).
					WillReturnError(errors.New("query pricing rules failed"))
			},
			wantErr:       true,
			expectedError: errors.New("query pricing rules failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			rules, err := repo.GetPricingRulesByPropertyID(context.Background(), tt.propertyID, tt.start, tt.end)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(rules) != tt.expectedLen {
					t.Errorf("expected %d rules, got %d", tt.expectedLen, len(rules))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_GetAllPricingRulesByPropertyID(t *testing.T) {
	propertyID := uint(10)
	now := time.Now()

	tests := []struct {
		name          string
		propertyID    uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedLen   int
	}{
		{
			name:       "success",
			propertyID: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "property_id", "name", "start_date", "end_date", "price_modifier", "description", "is_active", "created_at"}).
					AddRow(uint(1), propertyID, "Rule 1", now, now.Add(48*time.Hour), 1.2, "Desc", true, now)

				mock.ExpectQuery("^SELECT \\* FROM `pricing_rules` WHERE property_id = \\? ORDER BY start_date ASC").
					WithArgs(propertyID).
					WillReturnRows(rows)
			},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:       "database error",
			propertyID: 11,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `pricing_rules` WHERE property_id = \\? ORDER BY start_date ASC").
					WithArgs(uint(11)).
					WillReturnError(errors.New("all rules query failed"))
			},
			wantErr:       true,
			expectedError: errors.New("all rules query failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			rules, err := repo.GetAllPricingRulesByPropertyID(context.Background(), tt.propertyID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(rules) != tt.expectedLen {
					t.Errorf("expected %d rules, got %d", tt.expectedLen, len(rules))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestBookingRepository_DeletePricingRule(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `pricing_rules` WHERE `pricing_rules`.`id` = \\?").
					WithArgs(uint(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `pricing_rules` WHERE `pricing_rules`.`id` = \\?").
					WillReturnError(errors.New("delete rule failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("delete rule failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.mockSetup(mock)

			repo := NewBookingRepository(gormDB)
			err = repo.DeletePricingRule(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}
