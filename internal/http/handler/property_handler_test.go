package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
	"digital-greenhouse/greenhouse-be/internal/http/dto"
	"digital-greenhouse/greenhouse-be/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

type mockPropertyServiceForHandler struct {
	CreatePropertyFunc              func(ctx context.Context, property *domain.Property) error
	ListPropertiesFunc              func(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error)
	GetPropertyByIDFunc             func(ctx context.Context, id uint) (*domain.Property, error)
	GetPropertiesByOwnerFunc         func(ctx context.Context, ownerID uint) ([]domain.Property, error)
	UpdatePropertyFunc              func(ctx context.Context, property *domain.Property) error
	DeletePropertyFunc              func(ctx context.Context, id uint) error
	AddPropertyImageFunc            func(ctx context.Context, image *domain.PropertyImage) error
	UpdatePropertyImageFunc         func(ctx context.Context, image *domain.PropertyImage) error
	DeletePropertyImageFunc         func(ctx context.Context, imageID uint) error
	CreatePricingRuleFunc           func(ctx context.Context, rule *domain.PricingRule) error
	ListPricingRulesByPropertyFunc  func(ctx context.Context, propertyID uint) ([]domain.PricingRule, error)
	DeletePricingRuleFunc           func(ctx context.Context, id uint) error
	AutoGenerateHighSeasonRulesFunc func(ctx context.Context, propertyID uint) error
}

func (m *mockPropertyServiceForHandler) CreateProperty(ctx context.Context, property *domain.Property) error {
	if m.CreatePropertyFunc != nil {
		return m.CreatePropertyFunc(ctx, property)
	}
	return nil
}

func (m *mockPropertyServiceForHandler) ListProperties(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error) {
	if m.ListPropertiesFunc != nil {
		return m.ListPropertiesFunc(ctx, filter)
	}
	return nil, nil
}

func (m *mockPropertyServiceForHandler) GetPropertyByID(ctx context.Context, id uint) (*domain.Property, error) {
	if m.GetPropertyByIDFunc != nil {
		return m.GetPropertyByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockPropertyServiceForHandler) GetPropertiesByOwner(ctx context.Context, ownerID uint) ([]domain.Property, error) {
	if m.GetPropertiesByOwnerFunc != nil {
		return m.GetPropertiesByOwnerFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *mockPropertyServiceForHandler) UpdateProperty(ctx context.Context, property *domain.Property) error {
	if m.UpdatePropertyFunc != nil {
		return m.UpdatePropertyFunc(ctx, property)
	}
	return nil
}

func (m *mockPropertyServiceForHandler) DeleteProperty(ctx context.Context, id uint) error {
	if m.DeletePropertyFunc != nil {
		return m.DeletePropertyFunc(ctx, id)
	}
	return nil
}

func (m *mockPropertyServiceForHandler) AddPropertyImage(ctx context.Context, image *domain.PropertyImage) error {
	if m.AddPropertyImageFunc != nil {
		return m.AddPropertyImageFunc(ctx, image)
	}
	return nil
}

func (m *mockPropertyServiceForHandler) UpdatePropertyImage(ctx context.Context, image *domain.PropertyImage) error {
	if m.UpdatePropertyImageFunc != nil {
		return m.UpdatePropertyImageFunc(ctx, image)
	}
	return nil
}

func (m *mockPropertyServiceForHandler) DeletePropertyImage(ctx context.Context, imageID uint) error {
	if m.DeletePropertyImageFunc != nil {
		return m.DeletePropertyImageFunc(ctx, imageID)
	}
	return nil
}

func (m *mockPropertyServiceForHandler) CreatePricingRule(ctx context.Context, rule *domain.PricingRule) error {
	if m.CreatePricingRuleFunc != nil {
		return m.CreatePricingRuleFunc(ctx, rule)
	}
	return nil
}

func (m *mockPropertyServiceForHandler) ListPricingRulesByProperty(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
	if m.ListPricingRulesByPropertyFunc != nil {
		return m.ListPricingRulesByPropertyFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *mockPropertyServiceForHandler) DeletePricingRule(ctx context.Context, id uint) error {
	if m.DeletePricingRuleFunc != nil {
		return m.DeletePricingRuleFunc(ctx, id)
	}
	return nil
}

func (m *mockPropertyServiceForHandler) AutoGenerateHighSeasonRules(ctx context.Context, propertyID uint) error {
	if m.AutoGenerateHighSeasonRulesFunc != nil {
		return m.AutoGenerateHighSeasonRulesFunc(ctx, propertyID)
	}
	return nil
}

func TestPropertyHandler_CreateProperty(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		reqBody        string
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:   "Happy Path with Images",
			userID: 10,
			reqBody: `{
				"name": "Beach Villa",
				"description": "Lovely villa",
				"address": "123 Coast Rd",
				"base_price_per_night": 150.0,
				"max_capacity": 6,
				"images": [{"image_data":"base64data","mime_type":"image/png","alt_text":"Cover","is_cover":true,"sort_order":1}]
			}`,
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.CreatePropertyFunc = func(ctx context.Context, property *domain.Property) error {
					if property.OwnerID != 10 || property.Name != "Beach Villa" || len(property.Images) != 1 {
						return errors.New("invalid arguments to service")
					}
					property.ID = 100
					property.Status = domain.PropertyActive
					property.CreatedAt = time.Unix(1600000000, 0)
					property.UpdatedAt = time.Unix(1600000000, 0)
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.PropertyResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if resp.ID != 100 || resp.OwnerID != 10 || len(resp.Images) != 1 || resp.Images[0].ImageData != "base64data" {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:           "Unauthorized",
			userID:         0,
			reqBody:        `{}`,
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusUnauthorized,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("se requiere autenticación")) {
					t.Errorf("expected unauthorized, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid JSON",
			userID:         10,
			reqBody:        `{"name":`,
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid, got: %s", body)
				}
			},
		},
		{
			name:   "Service Error",
			userID: 10,
			reqBody: `{"name":"Villa"}`,
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.CreatePropertyFunc = func(ctx context.Context, property *domain.Property) error {
					return errors.New("db error")
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
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			req := httptest.NewRequest(http.MethodPost, "/properties", bytes.NewBufferString(tt.reqBody))
			if tt.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}
			rec := httptest.NewRecorder()

			h.CreateProperty(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestPropertyHandler_ListProperties(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    url.Values
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name: "Happy Path - With Filters",
			queryParams: url.Values{
				"search":    {"beach"},
				"location":  {"miami"},
				"min_price": {"100.0"},
				"max_price": {"500.0"},
				"guests":    {"4"},
				"check_in":  {"2026-06-01"},
				"check_out": {"2026-06-10"},
			},
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.ListPropertiesFunc = func(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error) {
					if filter.Search != "beach" || filter.Location != "miami" || filter.MinPrice != 100.0 || filter.MaxPrice != 500.0 || filter.GuestCount != 4 {
						return nil, errors.New("filter mismatch")
					}
					if filter.CheckInDate == nil || filter.CheckInDate.Format("2006-01-02") != "2026-06-01" {
						return nil, errors.New("check in date mismatch")
					}
					if filter.CheckOutDate == nil || filter.CheckOutDate.Format("2006-01-02") != "2026-06-10" {
						return nil, errors.New("check out date mismatch")
					}
					return []domain.Property{
						{ID: 1, Name: "Villa 1", CreatedAt: time.Unix(1600000000, 0), UpdatedAt: time.Unix(1600000000, 0)},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp []dto.PropertyResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(resp) != 1 || resp[0].ID != 1 {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name: "Service Error",
			queryParams: url.Values{},
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.ListPropertiesFunc = func(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error) {
					return nil, errors.New("list failed")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("list failed")) {
					t.Errorf("expected list failed error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			uri := "/properties"
			if len(tt.queryParams) > 0 {
				uri += "?" + tt.queryParams.Encode()
			}
			req := httptest.NewRequest(http.MethodGet, uri, nil)
			rec := httptest.NewRecorder()

			h.ListProperties(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestPropertyHandler_GetPropertyByID(t *testing.T) {
	tests := []struct {
		name           string
		propertyID     string
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:       "Happy Path",
			propertyID: "100",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.GetPropertyByIDFunc = func(ctx context.Context, id uint) (*domain.Property, error) {
					if id != 100 {
						return nil, errors.New("wrong ID")
					}
					return &domain.Property{
						ID:        100,
						Name:      "Unique Villa",
						CreatedAt: time.Unix(1600000000, 0),
						UpdatedAt: time.Unix(1600000000, 0),
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.PropertyResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if resp.ID != 100 || resp.Name != "Unique Villa" {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:       "Property Not Found",
			propertyID: "100",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.GetPropertyByIDFunc = func(ctx context.Context, id uint) (*domain.Property, error) {
					return nil, nil
				}
			},
			expectedStatus: http.StatusNotFound,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("propiedad no encontrada")) {
					t.Errorf("expected not found error, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid Property ID",
			propertyID:     "abc",
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de propiedad inválido")) {
					t.Errorf("expected invalid id error, got: %s", body)
				}
			},
		},
		{
			name:       "Service Error",
			propertyID: "100",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.GetPropertyByIDFunc = func(ctx context.Context, id uint) (*domain.Property, error) {
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
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			r := chi.NewRouter()
			r.Get("/properties/{id}", h.GetPropertyByID)

			req := httptest.NewRequest(http.MethodGet, "/properties/"+tt.propertyID, nil)
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

func TestPropertyHandler_GetPropertiesByOwner(t *testing.T) {
	tests := []struct {
		name           string
		ownerID        string
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:    "Happy Path",
			ownerID: "10",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.GetPropertiesByOwnerFunc = func(ctx context.Context, ownerID uint) ([]domain.Property, error) {
					if ownerID != 10 {
						return nil, errors.New("wrong ID")
					}
					return []domain.Property{
						{ID: 1, Name: "Villa 1", CreatedAt: time.Unix(1600000000, 0), UpdatedAt: time.Unix(1600000000, 0)},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp []dto.PropertyResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(resp) != 1 || resp[0].ID != 1 {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:           "Invalid Owner ID",
			ownerID:        "abc",
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de dueño inválido")) {
					t.Errorf("expected invalid owner id, got: %s", body)
				}
			},
		},
		{
			name:    "Service Error",
			ownerID: "10",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.GetPropertiesByOwnerFunc = func(ctx context.Context, ownerID uint) ([]domain.Property, error) {
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
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			r := chi.NewRouter()
			r.Get("/owners/{id}/properties", h.GetPropertiesByOwner)

			req := httptest.NewRequest(http.MethodGet, "/owners/"+tt.ownerID+"/properties", nil)
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

func TestPropertyHandler_AddImage(t *testing.T) {
	tests := []struct {
		name           string
		propertyID     string
		reqBody        string
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:       "Happy Path",
			propertyID: "100",
			reqBody:    `{"image_data":"base64","mime_type":"image/png","alt_text":"Alt","is_cover":true,"sort_order":1}`,
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.AddPropertyImageFunc = func(ctx context.Context, image *domain.PropertyImage) error {
					if image.PropertyID != 100 || image.ImageData != "base64" {
						return errors.New("invalid arguments to service")
					}
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.PropertyImageDTO
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if resp.ImageData != "base64" || resp.MimeType != "image/png" {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:           "Invalid Property ID",
			propertyID:     "abc",
			reqBody:        `{}`,
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de propiedad inválido")) {
					t.Errorf("expected invalid property id error, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid JSON",
			propertyID:     "100",
			reqBody:        `{"image_data":`,
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid error, got: %s", body)
				}
			},
		},
		{
			name:       "Service Error",
			propertyID: "100",
			reqBody:    `{"image_data":"base64"}`,
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.AddPropertyImageFunc = func(ctx context.Context, image *domain.PropertyImage) error {
					return errors.New("db error")
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
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			r := chi.NewRouter()
			r.Post("/properties/{id}/images", h.AddImage)

			req := httptest.NewRequest(http.MethodPost, "/properties/"+tt.propertyID+"/images", bytes.NewBufferString(tt.reqBody))
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

func TestPropertyHandler_DeleteImage(t *testing.T) {
	tests := []struct {
		name           string
		imageID        string
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:    "Happy Path",
			imageID: "50",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.DeletePropertyImageFunc = func(ctx context.Context, id uint) error {
					if id != 50 {
						return errors.New("wrong ID")
					}
					return nil
				}
			},
			expectedStatus: http.StatusNoContent,
			verifyResponse: func(t *testing.T, body string) {
				if body != "" {
					t.Errorf("expected empty body, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid Image ID",
			imageID:        "abc",
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de imagen inválido")) {
					t.Errorf("expected invalid image id error, got: %s", body)
				}
			},
		},
		{
			name:    "Service Error",
			imageID: "50",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.DeletePropertyImageFunc = func(ctx context.Context, id uint) error {
					return errors.New("db error")
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
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			r := chi.NewRouter()
			r.Delete("/properties/images/{imageID}", h.DeleteImage)

			req := httptest.NewRequest(http.MethodDelete, "/properties/images/"+tt.imageID, nil)
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

func TestPropertyHandler_CreatePricingRule(t *testing.T) {
	tests := []struct {
		name           string
		propertyID     string
		reqBody        string
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:       "Happy Path",
			propertyID: "10",
			reqBody:    `{"name":"High Season","start_date":"2026-12-01","end_date":"2026-12-31","price_modifier":1.5,"description":"Winter holidays"}`,
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.CreatePricingRuleFunc = func(ctx context.Context, rule *domain.PricingRule) error {
					if rule.PropertyID != 10 || rule.Name != "High Season" || rule.PriceModifier != 1.5 {
						return errors.New("invalid arguments to service")
					}
					rule.ID = 101
					rule.IsActive = true
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.PricingRuleDTO
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if resp.ID != 101 || resp.StartDate != "2026-12-01" || resp.EndDate != "2026-12-31" {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:           "Invalid Property ID",
			propertyID:     "abc",
			reqBody:        `{}`,
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de propiedad inválido")) {
					t.Errorf("expected invalid property id error, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid JSON",
			propertyID:     "10",
			reqBody:        `{"name":`,
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid error, got: %s", body)
				}
			},
		},
		{
			name:       "Invalid Start Date Format",
			propertyID: "10",
			reqBody:    `{"name":"Rule","start_date":"2026/12/01","end_date":"2026-12-31","price_modifier":1.5}`,
			mockSetup:  func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("fecha de inicio inválida")) {
					t.Errorf("expected invalid start date format error, got: %s", body)
				}
			},
		},
		{
			name:       "Invalid End Date Format",
			propertyID: "10",
			reqBody:    `{"name":"Rule","start_date":"2026-12-01","end_date":"2026/12/31","price_modifier":1.5}`,
			mockSetup:  func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("fecha de fin inválida")) {
					t.Errorf("expected invalid end date format error, got: %s", body)
				}
			},
		},
		{
			name:       "Service Error",
			propertyID: "10",
			reqBody:    `{"name":"Rule","start_date":"2026-12-01","end_date":"2026-12-31","price_modifier":1.5}`,
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.CreatePricingRuleFunc = func(ctx context.Context, rule *domain.PricingRule) error {
					return errors.New("db error")
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
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			r := chi.NewRouter()
			r.Post("/properties/{id}/pricing-rules", h.CreatePricingRule)

			req := httptest.NewRequest(http.MethodPost, "/properties/"+tt.propertyID+"/pricing-rules", bytes.NewBufferString(tt.reqBody))
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

func TestPropertyHandler_GetPricingRules(t *testing.T) {
	tests := []struct {
		name           string
		propertyID     string
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:       "Happy Path",
			propertyID: "10",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.ListPricingRulesByPropertyFunc = func(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
					if propertyID != 10 {
						return nil, errors.New("wrong ID")
					}
					return []domain.PricingRule{
						{
							ID:            1,
							PropertyID:    10,
							Name:          "Rule 1",
							StartDate:     time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC),
							EndDate:       time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
							PriceModifier: 1.2,
							IsActive:      true,
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp []dto.PricingRuleDTO
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(resp) != 1 || resp[0].ID != 1 || resp[0].StartDate != "2026-12-01" {
					t.Errorf("unexpected response: %+v", resp)
				}
			},
		},
		{
			name:           "Invalid Property ID",
			propertyID:     "abc",
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de propiedad inválido")) {
					t.Errorf("expected invalid property id error, got: %s", body)
				}
			},
		},
		{
			name:       "Service Error",
			propertyID: "10",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.ListPricingRulesByPropertyFunc = func(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
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
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			r := chi.NewRouter()
			r.Get("/properties/{id}/pricing-rules", h.GetPricingRules)

			req := httptest.NewRequest(http.MethodGet, "/properties/"+tt.propertyID+"/pricing-rules", nil)
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

func TestPropertyHandler_DeletePricingRule(t *testing.T) {
	tests := []struct {
		name           string
		ruleID         string
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:   "Happy Path",
			ruleID: "50",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.DeletePricingRuleFunc = func(ctx context.Context, id uint) error {
					if id != 50 {
						return errors.New("wrong ID")
					}
					return nil
				}
			},
			expectedStatus: http.StatusNoContent,
			verifyResponse: func(t *testing.T, body string) {
				if body != "" {
					t.Errorf("expected empty body, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid Rule ID",
			ruleID:         "abc",
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de regla inválido")) {
					t.Errorf("expected invalid rule id error, got: %s", body)
				}
			},
		},
		{
			name:   "Service Error",
			ruleID: "50",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.DeletePricingRuleFunc = func(ctx context.Context, id uint) error {
					return errors.New("db error")
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
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			r := chi.NewRouter()
			r.Delete("/properties/pricing-rules/{ruleId}", h.DeletePricingRule)

			req := httptest.NewRequest(http.MethodDelete, "/properties/pricing-rules/"+tt.ruleID, nil)
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

func TestPropertyHandler_AutoGeneratePricingRules(t *testing.T) {
	tests := []struct {
		name           string
		propertyID     string
		mockSetup      func(m *mockPropertyServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:       "Happy Path",
			propertyID: "10",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.AutoGenerateHighSeasonRulesFunc = func(ctx context.Context, id uint) error {
					if id != 10 {
						return errors.New("wrong ID")
					}
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("reglas de temporada alta generadas exitosamente")) {
					t.Errorf("expected success message, got: %s", body)
				}
			},
		},
		{
			name:           "Invalid Property ID",
			propertyID:     "abc",
			mockSetup:      func(m *mockPropertyServiceForHandler) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de propiedad inválido")) {
					t.Errorf("expected invalid property id error, got: %s", body)
				}
			},
		},
		{
			name:       "Service Error",
			propertyID: "10",
			mockSetup: func(m *mockPropertyServiceForHandler) {
				m.AutoGenerateHighSeasonRulesFunc = func(ctx context.Context, id uint) error {
					return errors.New("db error")
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
			m := &mockPropertyServiceForHandler{}
			tt.mockSetup(m)
			h := NewPropertyHandler(m)

			r := chi.NewRouter()
			r.Post("/properties/{id}/pricing-rules/auto", h.AutoGeneratePricingRules)

			req := httptest.NewRequest(http.MethodPost, "/properties/"+tt.propertyID+"/pricing-rules/auto", nil)
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
