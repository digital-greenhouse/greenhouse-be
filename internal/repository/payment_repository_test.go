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

func TestPaymentRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		payment       *domain.Payment
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			payment: &domain.Payment{
				BookingID:     10,
				Amount:        150.75,
				PaymentMethod: domain.PaymentMethodTransfer,
				ProofData:     "base64_encoded_proof",
				ProofMimeType: "image/png",
				Status:        domain.PaymentPending,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `payments`").
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
					).
					WillReturnResult(sqlmock.NewResult(5, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			payment: &domain.Payment{
				BookingID: 10,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `payments`").
					WillReturnError(errors.New("insert payment failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("insert payment failed"),
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

			repo := NewPaymentRepository(gormDB)
			err = repo.Create(context.Background(), tt.payment)

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
				if tt.payment.ID != 5 {
					t.Errorf("expected payment ID to be updated to 5, got %d", tt.payment.ID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPaymentRepository_GetByID(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name            string
		id              uint
		mockSetup       func(mock sqlmock.Sqlmock)
		wantErr         bool
		expectedError   error
		expectedPayment *domain.Payment
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "booking_id", "amount", "payment_method", "proof_data", "proof_mime_type", "status", "rejection_reason", "verified_by", "payment_date", "verified_at"}).
					AddRow(uint(1), uint(10), 150.75, "TRANSFERENCIA", "base64_encoded_proof", "image/png", "PENDING_VERIFICATION", "", nil, now, nil)
				mock.ExpectQuery("^SELECT \\* FROM `payments` WHERE `payments`.`id` = \\?").
					WithArgs(uint(1), 1).
					WillReturnRows(rows)
			},
			wantErr: false,
			expectedPayment: &domain.Payment{
				ID:            1,
				BookingID:     10,
				Amount:        150.75,
				PaymentMethod: domain.PaymentMethodTransfer,
				ProofData:     "base64_encoded_proof",
				ProofMimeType: "image/png",
				Status:        domain.PaymentPending,
			},
		},
		{
			name: "record not found",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `payments` WHERE `payments`.`id` = \\?").
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

			repo := NewPaymentRepository(gormDB)
			payment, err := repo.GetByID(context.Background(), tt.id)

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
				if payment == nil {
					t.Errorf("expected payment not to be nil")
				} else if payment.ID != tt.expectedPayment.ID || payment.BookingID != tt.expectedPayment.BookingID || payment.Amount != tt.expectedPayment.Amount {
					t.Errorf("returned payment %+v does not match expected %+v", payment, tt.expectedPayment)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPaymentRepository_GetByBookingID(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		bookingID     uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedLen   int
	}{
		{
			name:      "success",
			bookingID: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "booking_id", "amount", "payment_method", "proof_data", "proof_mime_type", "status", "rejection_reason", "verified_by", "payment_date", "verified_at"}).
					AddRow(uint(1), uint(10), 150.75, "TRANSFERENCIA", "base64_encoded_proof", "image/png", "PENDING_VERIFICATION", "", nil, now, nil).
					AddRow(uint(2), uint(10), 50.00, "EFECTIVO", "base64_encoded_proof_2", "image/jpeg", "VERIFIED", "", nil, now, nil)
				mock.ExpectQuery("^SELECT \\* FROM `payments` WHERE booking_id = \\?").
					WithArgs(uint(10)).
					WillReturnRows(rows)
			},
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name:      "database error",
			bookingID: 11,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `payments` WHERE booking_id = \\?").
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

			repo := NewPaymentRepository(gormDB)
			payments, err := repo.GetByBookingID(context.Background(), tt.bookingID)

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
				if len(payments) != tt.expectedLen {
					t.Errorf("expected %d payments, got %d", tt.expectedLen, len(payments))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPaymentRepository_UpdateStatus(t *testing.T) {
	verifierID := uint(3)

	tests := []struct {
		name          string
		id            uint
		status        domain.PaymentStatus
		verifierID    *uint
		reason        string
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name:       "success with verifier",
			id:         1,
			status:     domain.PaymentVerified,
			verifierID: &verifierID,
			reason:     "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `payments` SET").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:       "success with rejection reason",
			id:         2,
			status:     domain.PaymentRejected,
			verifierID: nil,
			reason:     "Invalid bank receipt",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `payments` SET").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:       "database error",
			id:         3,
			status:     domain.PaymentVerified,
			verifierID: nil,
			reason:     "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `payments` SET").
					WillReturnError(errors.New("update payment status failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("update payment status failed"),
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

			repo := NewPaymentRepository(gormDB)
			err = repo.UpdateStatus(context.Background(), tt.id, tt.status, tt.verifierID, tt.reason)

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
