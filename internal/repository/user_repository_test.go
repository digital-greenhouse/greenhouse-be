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

func TestUserRepository_Create(t *testing.T) {
	phoneStr := "+12345678"

	tests := []struct {
		name          string
		user          *domain.User
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success with phone and default active",
			user: &domain.User{
				Name:         "John Doe",
				Email:        "john@example.com",
				PasswordHash: "hashed_password",
				Role:         domain.RoleClient,
				Phone:        &phoneStr,
				IsActive:     true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `users`").
					WithArgs(
						"John Doe",
						"john@example.com",
						"hashed_password",
						"CLIENT",
						&phoneStr,
						true,
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
					).
					WillReturnResult(sqlmock.NewResult(10, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "success with nil phone",
			user: &domain.User{
				Name:         "Jane Doe",
				Email:        "jane@example.com",
				PasswordHash: "hashed_password",
				Role:         domain.RoleSuperAdmin,
				Phone:        nil,
				IsActive:     false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `users`").
					WithArgs(
						"Jane Doe",
						"jane@example.com",
						"hashed_password",
						"SUPERADMIN",
						nil,
						true,
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
					).
					WillReturnResult(sqlmock.NewResult(11, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			user: &domain.User{
				Name:  "Error User",
				Email: "error@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `users`").
					WillReturnError(errors.New("insert failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("insert failed"),
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

			repo := NewUserRepository(gormDB)
			err = repo.Create(context.Background(), tt.user)

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
				if tt.user.ID == 0 {
					t.Errorf("expected user ID to be updated, got 0")
				}
				if tt.user.CreatedAt.IsZero() || tt.user.UpdatedAt.IsZero() {
					t.Errorf("expected timestamps to be updated, got zero")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	phoneStr := "+12345678"
	now := time.Now()

	tests := []struct {
		name          string
		id            uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedUser  *domain.User
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role", "phone", "is_active", "created_at", "updated_at"}).
					AddRow(uint(1), "John Doe", "john@example.com", "hashed_password", "CLIENT", &phoneStr, true, now, now)
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\?").
					WithArgs(uint(1), 1).
					WillReturnRows(rows)
			},
			wantErr: false,
			expectedUser: &domain.User{
				ID:           1,
				Name:         "John Doe",
				Email:        "john@example.com",
				PasswordHash: "hashed_password",
				Role:         domain.RoleClient,
				Phone:        &phoneStr,
				IsActive:     true,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		{
			name: "record not found",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\?").
					WithArgs(uint(2), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr:       true,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "database error",
			id:   3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\?").
					WithArgs(uint(3), 1).
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

			repo := NewUserRepository(gormDB)
			user, err := repo.GetByID(context.Background(), tt.id)

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
				if user == nil {
					t.Errorf("expected user not to be nil")
				} else {
					if user.ID != tt.expectedUser.ID || user.Name != tt.expectedUser.Name || user.Email != tt.expectedUser.Email || user.Role != tt.expectedUser.Role {
						t.Errorf("returned user %+v does not match expected %+v", user, tt.expectedUser)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	phoneStr := "+12345678"
	now := time.Now()

	tests := []struct {
		name          string
		email         string
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedUser  *domain.User
	}{
		{
			name:  "success",
			email: "john@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role", "phone", "is_active", "created_at", "updated_at"}).
					AddRow(uint(1), "John Doe", "john@example.com", "hashed_password", "CLIENT", &phoneStr, true, now, now)
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE email = \\?").
					WithArgs("john@example.com", 1).
					WillReturnRows(rows)
			},
			wantErr: false,
			expectedUser: &domain.User{
				ID:           1,
				Name:         "John Doe",
				Email:        "john@example.com",
				PasswordHash: "hashed_password",
				Role:         domain.RoleClient,
				Phone:        &phoneStr,
				IsActive:     true,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		{
			name:  "record not found",
			email: "notfound@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE email = \\?").
					WithArgs("notfound@example.com", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr:       true,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "database error",
			email: "error@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE email = \\?").
					WithArgs("error@example.com", 1).
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

			repo := NewUserRepository(gormDB)
			user, err := repo.GetByEmail(context.Background(), tt.email)

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
				if user == nil {
					t.Errorf("expected user not to be nil")
				} else {
					if user.ID != tt.expectedUser.ID || user.Name != tt.expectedUser.Name || user.Email != tt.expectedUser.Email {
						t.Errorf("returned user %+v does not match expected %+v", user, tt.expectedUser)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepository_GetAll(t *testing.T) {
	phoneStr := "+12345678"
	now := time.Now()

	tests := []struct {
		name          string
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedLen   int
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role", "phone", "is_active", "created_at", "updated_at"}).
					AddRow(uint(1), "John Doe", "john@example.com", "hashed_password", "CLIENT", &phoneStr, true, now, now).
					AddRow(uint(2), "Jane Doe", "jane@example.com", "hashed_password_2", "OWNER", nil, false, now, now)
				mock.ExpectQuery("^SELECT \\* FROM `users`").
					WillReturnRows(rows)
			},
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users`").
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

			repo := NewUserRepository(gormDB)
			users, err := repo.GetAll(context.Background())

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
				if len(users) != tt.expectedLen {
					t.Errorf("expected %d users, got %d", tt.expectedLen, len(users))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	phoneStr := "+12345678"
	now := time.Now()

	tests := []struct {
		name          string
		user          *domain.User
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			user: &domain.User{
				ID:           5,
				Name:         "John Updated",
				Email:        "john@example.com",
				PasswordHash: "new_hash",
				Role:         domain.RoleClient,
				Phone:        &phoneStr,
				IsActive:     true,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `users` SET").
					WithArgs(
						"John Updated",
						"john@example.com",
						"new_hash",
						"CLIENT",
						&phoneStr,
						true,
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
						5,                // ID
					).
					WillReturnResult(sqlmock.NewResult(5, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			user: &domain.User{
				ID:           6,
				Name:         "Error User",
				Email:        "error@example.com",
				PasswordHash: "hash",
				Role:         domain.RoleClient,
				Phone:        nil,
				IsActive:     true,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `users` SET").
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("update failed"),
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

			repo := NewUserRepository(gormDB)
			err = repo.Update(context.Background(), tt.user)

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

func TestUserRepository_Delete(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			id:   5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `users` WHERE `users`.`id` = \\?").
					WithArgs(uint(5)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			id:   6,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `users` WHERE `users`.`id` = \\?").
					WithArgs(uint(6)).
					WillReturnError(errors.New("delete failed"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("delete failed"),
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

			repo := NewUserRepository(gormDB)
			err = repo.Delete(context.Background(), tt.id)

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
