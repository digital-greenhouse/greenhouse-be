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

func TestPropertyRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		property      *domain.Property
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			property: &domain.Property{
				OwnerID:           100,
				Name:              "Beach Villa",
				Description:       "Lovely beach villa",
				Address:           "123 Ocean Drive",
				BasePricePerNight: 150.0,
				MaxCapacity:       6,
				Status:            domain.PropertyActive,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `properties`").
					WithArgs(
						uint(100),
						"Beach Villa",
						"Lovely beach villa",
						"123 Ocean Drive",
						150.0,
						6,
						"ACTIVE",
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			property: &domain.Property{
				OwnerID: 100,
				Name:    "Error Villa",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `properties`").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("db error"),
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

			repo := NewPropertyRepository(gormDB)
			err = repo.Create(context.Background(), tt.property)

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
				if tt.property.ID != 1 {
					t.Errorf("expected auto-increment ID to be injected as 1, got %d", tt.property.ID)
				}
				if tt.property.CreatedAt.IsZero() || tt.property.UpdatedAt.IsZero() {
					t.Errorf("expected timestamps to be updated")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPropertyRepository_GetAll(t *testing.T) {
	now := time.Now()
	checkIn := now.Add(24 * time.Hour)
	checkOut := now.Add(48 * time.Hour)

	tests := []struct {
		name          string
		filter        domain.PropertyFilter
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedCount int
	}{
		{
			name:   "success empty filter",
			filter: domain.PropertyFilter{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// 1. SELECT query on properties
				propRows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "address", "base_price_per_night", "max_capacity", "status", "created_at", "updated_at"}).
					AddRow(uint(1), uint(100), "Beach Villa", "Lovely beach villa", "123 Ocean Drive", 150.0, 6, "ACTIVE", now, now)
				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE status = \\?").
					WithArgs("ACTIVE").
					WillReturnRows(propRows)

				// 2. SELECT query on property_images
				imgRows := sqlmock.NewRows([]string{"id", "property_id", "image_data", "mime_type", "alt_text", "is_cover", "sort_order", "created_at"}).
					AddRow(uint(10), uint(1), "base64data", "image/jpeg", "Cover image", true, 0, now)
				mock.ExpectQuery("^SELECT \\* FROM `property_images` WHERE .*property_id.* = \\?").
					WithArgs(uint(1)).
					WillReturnRows(imgRows)
			},
			wantErr:       false,
			expectedCount: 1,
		},
		{
			name: "success with search/location/price/capacity filters",
			filter: domain.PropertyFilter{
				Search:     "beach",
				Location:   "florida",
				MinPrice:   100.0,
				MaxPrice:   200.0,
				GuestCount: 4,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// 1. SELECT query on properties with filters
				propRows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "address", "base_price_per_night", "max_capacity", "status", "created_at", "updated_at"}).
					AddRow(uint(1), uint(100), "Beach Villa", "Lovely beach villa", "123 Ocean Drive", 150.0, 6, "ACTIVE", now, now)

				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE status = \\? AND .*name LIKE.*description LIKE.* AND .*address LIKE.* AND .*base_price_per_night >=.* AND .*base_price_per_night <=.* AND .*max_capacity >=.*").
					WithArgs("ACTIVE", "%beach%", "%beach%", "%florida%", 100.0, 200.0, 4).
					WillReturnRows(propRows)

				// 2. SELECT query on property_images (preload)
				imgRows := sqlmock.NewRows([]string{"id", "property_id", "image_data", "mime_type", "alt_text", "is_cover", "sort_order", "created_at"}).
					AddRow(uint(10), uint(1), "base64data", "image/jpeg", "Cover image", true, 0, now)
				mock.ExpectQuery("^SELECT \\* FROM `property_images` WHERE .*property_id.* = \\?").
					WithArgs(uint(1)).
					WillReturnRows(imgRows)
			},
			wantErr:       false,
			expectedCount: 1,
		},
		{
			name: "success with dates filter",
			filter: domain.PropertyFilter{
				CheckInDate:  &checkIn,
				CheckOutDate: &checkOut,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// 1. SELECT query on properties with bookings subquery
				propRows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "address", "base_price_per_night", "max_capacity", "status", "created_at", "updated_at"}).
					AddRow(uint(1), uint(100), "Beach Villa", "Lovely beach villa", "123 Ocean Drive", 150.0, 6, "ACTIVE", now, now)

				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE status = \\? AND id NOT IN \\(SELECT .* FROM `bookings` WHERE status IN \\(\\?, \\?\\) AND \\(check_in_date < \\? AND check_out_date > \\?\\)\\)").
					WithArgs("ACTIVE", "CONFIRMED", "PENDING_PAYMENT", checkOut, checkIn).
					WillReturnRows(propRows)

				// 2. SELECT query on property_images (preload)
				imgRows := sqlmock.NewRows([]string{"id", "property_id", "image_data", "mime_type", "alt_text", "is_cover", "sort_order", "created_at"}).
					AddRow(uint(10), uint(1), "base64data", "image/jpeg", "Cover image", true, 0, now)
				mock.ExpectQuery("^SELECT \\* FROM `property_images` WHERE .*property_id.* = \\?").
					WithArgs(uint(1)).
					WillReturnRows(imgRows)
			},
			wantErr:       false,
			expectedCount: 1,
		},
		{
			name:   "database error on properties query",
			filter: domain.PropertyFilter{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE status = \\?").
					WithArgs("ACTIVE").
					WillReturnError(errors.New("db query error"))
			},
			wantErr:       true,
			expectedError: errors.New("db query error"),
		},
		{
			name:   "database error on preload query",
			filter: domain.PropertyFilter{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				propRows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "address", "base_price_per_night", "max_capacity", "status", "created_at", "updated_at"}).
					AddRow(uint(1), uint(100), "Beach Villa", "Lovely beach villa", "123 Ocean Drive", 150.0, 6, "ACTIVE", now, now)
				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE status = \\?").
					WithArgs("ACTIVE").
					WillReturnRows(propRows)

				mock.ExpectQuery("^SELECT \\* FROM `property_images` WHERE .*property_id.* = \\?").
					WithArgs(uint(1)).
					WillReturnError(errors.New("preload error"))
			},
			wantErr:       true,
			expectedError: errors.New("preload error"),
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

			repo := NewPropertyRepository(gormDB)
			res, err := repo.GetAll(context.Background(), tt.filter)

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
				if len(res) != tt.expectedCount {
					t.Errorf("expected %d properties, got %d", tt.expectedCount, len(res))
				}
				if len(res) > 0 && len(res[0].Images) != 1 {
					t.Errorf("expected 1 preloaded image, got %d", len(res[0].Images))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPropertyRepository_GetByID(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		id            uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedProp  *domain.Property
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				propRows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "address", "base_price_per_night", "max_capacity", "status", "created_at", "updated_at"}).
					AddRow(uint(1), uint(100), "Beach Villa", "Lovely beach villa", "123 Ocean Drive", 150.0, 6, "ACTIVE", now, now)
				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE `properties`.`id` = \\?").
					WithArgs(uint(1), 1).
					WillReturnRows(propRows)

				imgRows := sqlmock.NewRows([]string{"id", "property_id", "image_data", "mime_type", "alt_text", "is_cover", "sort_order", "created_at"}).
					AddRow(uint(10), uint(1), "base64data", "image/jpeg", "Cover image", true, 0, now)
				mock.ExpectQuery("^SELECT \\* FROM `property_images` WHERE .*property_id.* = \\?").
					WithArgs(uint(1)).
					WillReturnRows(imgRows)
			},
			wantErr: false,
			expectedProp: &domain.Property{
				ID:                1,
				OwnerID:           100,
				Name:              "Beach Villa",
				Description:       "Lovely beach villa",
				Address:           "123 Ocean Drive",
				BasePricePerNight: 150.0,
				MaxCapacity:       6,
				Status:            domain.PropertyActive,
			},
		},
		{
			name: "record not found",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE `properties`.`id` = \\?").
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
				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE `properties`.`id` = \\?").
					WithArgs(uint(3), 1).
					WillReturnError(errors.New("db error"))
			},
			wantErr:       true,
			expectedError: errors.New("db error"),
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

			repo := NewPropertyRepository(gormDB)
			res, err := repo.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() && err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if res == nil {
					t.Errorf("expected non-nil property")
				} else {
					if res.ID != tt.expectedProp.ID || res.Name != tt.expectedProp.Name || len(res.Images) != 1 {
						t.Errorf("returned property mismatch")
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPropertyRepository_GetByOwnerID(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		ownerID       uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedCount int
	}{
		{
			name:    "success",
			ownerID: 100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				propRows := sqlmock.NewRows([]string{"id", "owner_id", "name", "description", "address", "base_price_per_night", "max_capacity", "status", "created_at", "updated_at"}).
					AddRow(uint(1), uint(100), "Beach Villa", "Lovely beach villa", "123 Ocean Drive", 150.0, 6, "ACTIVE", now, now)
				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE owner_id = \\?").
					WithArgs(uint(100)).
					WillReturnRows(propRows)

				imgRows := sqlmock.NewRows([]string{"id", "property_id", "image_data", "mime_type", "alt_text", "is_cover", "sort_order", "created_at"}).
					AddRow(uint(10), uint(1), "base64data", "image/jpeg", "Cover image", true, 0, now)
				mock.ExpectQuery("^SELECT \\* FROM `property_images` WHERE .*property_id.* = \\?").
					WithArgs(uint(1)).
					WillReturnRows(imgRows)
			},
			wantErr:       false,
			expectedCount: 1,
		},
		{
			name:    "database error",
			ownerID: 101,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `properties` WHERE owner_id = \\?").
					WithArgs(uint(101)).
					WillReturnError(errors.New("db error"))
			},
			wantErr:       true,
			expectedError: errors.New("db error"),
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

			repo := NewPropertyRepository(gormDB)
			res, err := repo.GetByOwnerID(context.Background(), tt.ownerID)

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
				if len(res) != tt.expectedCount {
					t.Errorf("expected %d properties, got %d", tt.expectedCount, len(res))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPropertyRepository_Update(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		property      *domain.Property
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			property: &domain.Property{
				ID:                1,
				OwnerID:           100,
				Name:              "Updated Villa",
				Description:       "Lovely beach villa",
				Address:           "123 Ocean Drive",
				BasePricePerNight: 160.0,
				MaxCapacity:       7,
				Status:            domain.PropertyActive,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `properties` SET").
					WithArgs(
						uint(100),
						"Updated Villa",
						"Lovely beach villa",
						"123 Ocean Drive",
						160.0,
						7,
						"ACTIVE",
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
						uint(1),          // WHERE ID
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			property: &domain.Property{
				ID: 1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `properties` SET").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("db error"),
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

			repo := NewPropertyRepository(gormDB)
			err = repo.Update(context.Background(), tt.property)

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

func TestPropertyRepository_Delete(t *testing.T) {
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
				mock.ExpectExec("^DELETE FROM `properties` WHERE `properties`.`id` = \\?").
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
				mock.ExpectExec("^DELETE FROM `properties` WHERE `properties`.`id` = \\?").
					WithArgs(uint(2)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("db error"),
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

			repo := NewPropertyRepository(gormDB)
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

func TestPropertyRepository_AddImage(t *testing.T) {
	tests := []struct {
		name          string
		image         *domain.PropertyImage
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			image: &domain.PropertyImage{
				PropertyID: 1,
				ImageData:  "base64data",
				MimeType:   "image/jpeg",
				AltText:    "Cover image",
				IsCover:    true,
				SortOrder:  0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `property_images`").
					WithArgs(
						uint(1),
						"base64data",
						"image/jpeg",
						"Cover image",
						true,
						0,
						sqlmock.AnyArg(), // CreatedAt
					).
					WillReturnResult(sqlmock.NewResult(10, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			image: &domain.PropertyImage{
				PropertyID: 1,
				ImageData:  "base64data",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `property_images`").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("db error"),
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

			repo := NewPropertyRepository(gormDB)
			err = repo.AddImage(context.Background(), tt.image)

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
				if tt.image.ID != 10 {
					t.Errorf("expected auto-increment ID to be injected as 10, got %d", tt.image.ID)
				}
				if tt.image.CreatedAt.IsZero() {
					t.Errorf("expected CreatedAt timestamp to be updated")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPropertyRepository_GetImageByID(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		id            uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
		expectedImg   *domain.PropertyImage
	}{
		{
			name: "success",
			id:   10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				imgRows := sqlmock.NewRows([]string{"id", "property_id", "image_data", "mime_type", "alt_text", "is_cover", "sort_order", "created_at"}).
					AddRow(uint(10), uint(1), "base64data", "image/jpeg", "Cover image", true, 0, now)
				mock.ExpectQuery("^SELECT \\* FROM `property_images` WHERE `property_images`.`id` = \\?").
					WithArgs(uint(10), 1).
					WillReturnRows(imgRows)
			},
			wantErr: false,
			expectedImg: &domain.PropertyImage{
				ID:         10,
				PropertyID: 1,
				ImageData:  "base64data",
				MimeType:   "image/jpeg",
				AltText:    "Cover image",
				IsCover:    true,
				SortOrder:  0,
			},
		},
		{
			name: "record not found",
			id:   11,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `property_images` WHERE `property_images`.`id` = \\?").
					WithArgs(uint(11), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr:       true,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "database error",
			id:   12,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `property_images` WHERE `property_images`.`id` = \\?").
					WithArgs(uint(12), 1).
					WillReturnError(errors.New("db error"))
			},
			wantErr:       true,
			expectedError: errors.New("db error"),
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

			repo := NewPropertyRepository(gormDB)
			res, err := repo.GetImageByID(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.expectedError != nil && err.Error() != tt.expectedError.Error() && err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if res == nil {
					t.Errorf("expected non-nil image")
				} else {
					if res.ID != tt.expectedImg.ID || res.ImageData != tt.expectedImg.ImageData {
						t.Errorf("returned image mismatch")
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPropertyRepository_UpdateImage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		image         *domain.PropertyImage
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			image: &domain.PropertyImage{
				ID:         10,
				PropertyID: 1,
				ImageData:  "newBase64",
				MimeType:   "image/png",
				AltText:    "New Alt",
				IsCover:    false,
				SortOrder:  1,
				CreatedAt:  now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `property_images` SET").
					WithArgs(
						uint(1),
						"newBase64",
						"image/png",
						"New Alt",
						false,
						1,
						sqlmock.AnyArg(), // CreatedAt
						uint(10),         // WHERE ID
					).
					WillReturnResult(sqlmock.NewResult(10, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			image: &domain.PropertyImage{
				ID: 10,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `property_images` SET").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("db error"),
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

			repo := NewPropertyRepository(gormDB)
			err = repo.UpdateImage(context.Background(), tt.image)

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

func TestPropertyRepository_DeleteImage(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockSetup     func(mock sqlmock.Sqlmock)
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			id:   10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `property_images` WHERE `property_images`.`id` = \\?").
					WithArgs(uint(10)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			id:   11,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `property_images` WHERE `property_images`.`id` = \\?").
					WithArgs(uint(11)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: errors.New("db error"),
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

			repo := NewPropertyRepository(gormDB)
			err = repo.DeleteImage(context.Background(), tt.id)

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
