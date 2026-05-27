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

	"github.com/go-chi/chi/v5"
)

type mockUserServiceForHandler struct {
	CreateUserFunc  func(ctx context.Context, user *domain.User) error
	GetUserByIDFunc func(ctx context.Context, id uint) (*domain.User, error)
	GetAllUsersFunc func(ctx context.Context) ([]domain.User, error)
	UpdateUserFunc  func(ctx context.Context, user *domain.User) error
	DeleteUserFunc  func(ctx context.Context, id uint) error
	LoginFunc       func(ctx context.Context, email, password string) (string, *domain.User, error)
}

func (m *mockUserServiceForHandler) CreateUser(ctx context.Context, user *domain.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return nil
}

func (m *mockUserServiceForHandler) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserServiceForHandler) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	if m.GetAllUsersFunc != nil {
		return m.GetAllUsersFunc(ctx)
	}
	return nil, nil
}

func (m *mockUserServiceForHandler) UpdateUser(ctx context.Context, user *domain.User) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, user)
	}
	return nil
}

func (m *mockUserServiceForHandler) DeleteUser(ctx context.Context, id uint) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, id)
	}
	return nil
}

func (m *mockUserServiceForHandler) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, email, password)
	}
	return "", nil, nil
}

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        string
		mockSetup      func(m *mockUserServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:    "Happy Path with phone and role",
			reqBody: `{"name":"John Doe","email":"john@example.com","password":"securepassword","phone":"123456789","role":"CLIENT"}`,
			mockSetup: func(m *mockUserServiceForHandler) {
				m.CreateUserFunc = func(ctx context.Context, user *domain.User) error {
					user.ID = 1
					user.CreatedAt = time.Unix(1600000000, 0)
					user.IsActive = true
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.UserResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.ID != 1 || resp.Name != "John Doe" || resp.Email != "john@example.com" || resp.Role != "CLIENT" || resp.Phone == nil || *resp.Phone != "123456789" {
					t.Errorf("unexpected response content: %+v", resp)
				}
			},
		},
		{
			name:    "Invalid JSON",
			reqBody: `{"name":`,
			mockSetup: func(m *mockUserServiceForHandler) {
				// Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid error, got: %s", body)
				}
			},
		},
		{
			name:    "Service Error",
			reqBody: `{"name":"John Doe","email":"john@example.com","password":"securepassword"}`,
			mockSetup: func(m *mockUserServiceForHandler) {
				m.CreateUserFunc = func(ctx context.Context, user *domain.User) error {
					return errors.New("email already exists")
				}
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("email already exists")) {
					t.Errorf("expected email already exists error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserServiceForHandler{}
			tt.mockSetup(m)
			h := NewUserHandler(m)

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tt.reqBody))
			rec := httptest.NewRecorder()

			h.CreateUser(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        string
		mockSetup      func(m *mockUserServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:    "Happy Path",
			reqBody: `{"email":"john@example.com","password":"securepassword"}`,
			mockSetup: func(m *mockUserServiceForHandler) {
				m.LoginFunc = func(ctx context.Context, email, password string) (string, *domain.User, error) {
					return "mocked-jwt-token", &domain.User{
						ID:        1,
						Name:      "John Doe",
						Email:     "john@example.com",
						Role:      domain.RoleClient,
						IsActive:  true,
						CreatedAt: time.Unix(1600000000, 0),
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.LoginResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.Token != "mocked-jwt-token" || resp.User.ID != 1 {
					t.Errorf("unexpected response content: %+v", resp)
				}
			},
		},
		{
			name:    "Invalid JSON",
			reqBody: `{"email":`,
			mockSetup: func(m *mockUserServiceForHandler) {
				// Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid error, got: %s", body)
				}
			},
		},
		{
			name:    "Invalid Credentials",
			reqBody: `{"email":"john@example.com","password":"wrongpassword"}`,
			mockSetup: func(m *mockUserServiceForHandler) {
				m.LoginFunc = func(ctx context.Context, email, password string) (string, *domain.User, error) {
					return "", nil, errors.New("credenciales inválidas")
				}
			},
			expectedStatus: http.StatusUnauthorized,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("credenciales inválidas")) {
					t.Errorf("expected credentials error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserServiceForHandler{}
			tt.mockSetup(m)
			h := NewUserHandler(m)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.reqBody))
			rec := httptest.NewRecorder()

			h.Login(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userIDParam    string
		mockSetup      func(m *mockUserServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:        "Happy Path",
			userIDParam: "123",
			mockSetup: func(m *mockUserServiceForHandler) {
				m.GetUserByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					if id != 123 {
						return nil, errors.New("not matching ID")
					}
					return &domain.User{
						ID:        123,
						Name:      "John Doe",
						Email:     "john@example.com",
						Role:      domain.RoleClient,
						IsActive:  true,
						CreatedAt: time.Unix(1600000000, 0),
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.UserResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.ID != 123 || resp.Name != "John Doe" {
					t.Errorf("unexpected response content: %+v", resp)
				}
			},
		},
		{
			name:        "Invalid User ID",
			userIDParam: "abc",
			mockSetup: func(m *mockUserServiceForHandler) {
				// Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de usuario inválido")) {
					t.Errorf("expected user id invalid error, got: %s", body)
				}
			},
		},
		{
			name:        "User Not Found",
			userIDParam: "456",
			mockSetup: func(m *mockUserServiceForHandler) {
				m.GetUserByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					return nil, errors.New("not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("usuario no encontrado")) {
					t.Errorf("expected not found error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserServiceForHandler{}
			tt.mockSetup(m)
			h := NewUserHandler(m)

			r := chi.NewRouter()
			r.Get("/users/{id}", h.GetUser)

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userIDParam, nil)
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

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(m *mockUserServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name: "Happy Path",
			mockSetup: func(m *mockUserServiceForHandler) {
				m.GetAllUsersFunc = func(ctx context.Context) ([]domain.User, error) {
					return []domain.User{
						{
							ID:        1,
							Name:      "John Doe",
							Email:     "john@example.com",
							Role:      domain.RoleClient,
							IsActive:  true,
							CreatedAt: time.Unix(1600000000, 0),
						},
						{
							ID:        2,
							Name:      "Jane Doe",
							Email:     "jane@example.com",
							Role:      domain.RoleOwner,
							IsActive:  true,
							CreatedAt: time.Unix(1600000000, 0),
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp []dto.UserResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if len(resp) != 2 || resp[0].ID != 1 || resp[1].ID != 2 {
					t.Errorf("unexpected response content: %+v", resp)
				}
			},
		},
		{
			name: "Service Error",
			mockSetup: func(m *mockUserServiceForHandler) {
				m.GetAllUsersFunc = func(ctx context.Context) ([]domain.User, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("db error")) {
					t.Errorf("expected db error message, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserServiceForHandler{}
			tt.mockSetup(m)
			h := NewUserHandler(m)

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rec := httptest.NewRecorder()

			h.GetUsers(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec.Body.String())
			}
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userIDParam    string
		reqBody        string
		mockSetup      func(m *mockUserServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:        "Happy Path",
			userIDParam: "123",
			reqBody:     `{"name":"John Updated","phone":"987654321","role":"OWNER"}`,
			mockSetup: func(m *mockUserServiceForHandler) {
				m.UpdateUserFunc = func(ctx context.Context, user *domain.User) error {
					if user.ID != 123 || user.Name != "John Updated" || *user.Phone != "987654321" || user.Role != "OWNER" {
						return errors.New("incorrect parameters sent to service")
					}
					user.Email = "john@example.com"
					user.CreatedAt = time.Unix(1600000000, 0)
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var resp dto.UserResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.ID != 123 || resp.Name != "John Updated" || resp.Role != "OWNER" {
					t.Errorf("unexpected response content: %+v", resp)
				}
			},
		},
		{
			name:        "Invalid User ID",
			userIDParam: "abc",
			reqBody:     `{"name":"John Updated"}`,
			mockSetup: func(m *mockUserServiceForHandler) {
				// Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de usuario inválido")) {
					t.Errorf("expected user id invalid error, got: %s", body)
				}
			},
		},
		{
			name:        "Invalid JSON Payload",
			userIDParam: "123",
			reqBody:     `{"name":`,
			mockSetup: func(m *mockUserServiceForHandler) {
				// Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("payload inválido")) {
					t.Errorf("expected payload invalid error, got: %s", body)
				}
			},
		},
		{
			name:        "Service Error",
			userIDParam: "123",
			reqBody:     `{"name":"John Updated"}`,
			mockSetup: func(m *mockUserServiceForHandler) {
				m.UpdateUserFunc = func(ctx context.Context, user *domain.User) error {
					return errors.New("update database failed")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("update database failed")) {
					t.Errorf("expected update database failed error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserServiceForHandler{}
			tt.mockSetup(m)
			h := NewUserHandler(m)

			r := chi.NewRouter()
			r.Put("/users/{id}", h.UpdateUser)

			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userIDParam, bytes.NewBufferString(tt.reqBody))
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

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userIDParam    string
		mockSetup      func(m *mockUserServiceForHandler)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:        "Happy Path",
			userIDParam: "123",
			mockSetup: func(m *mockUserServiceForHandler) {
				m.DeleteUserFunc = func(ctx context.Context, id uint) error {
					if id != 123 {
						return errors.New("incorrect ID sent to service")
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
			name:        "Invalid User ID",
			userIDParam: "abc",
			mockSetup: func(m *mockUserServiceForHandler) {
				// Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("ID de usuario inválido")) {
					t.Errorf("expected user id invalid error, got: %s", body)
				}
			},
		},
		{
			name:        "Service Error",
			userIDParam: "123",
			mockSetup: func(m *mockUserServiceForHandler) {
				m.DeleteUserFunc = func(ctx context.Context, id uint) error {
					return errors.New("delete failed")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				if !bytes.Contains([]byte(body), []byte("delete failed")) {
					t.Errorf("expected delete failed error, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserServiceForHandler{}
			tt.mockSetup(m)
			h := NewUserHandler(m)

			r := chi.NewRouter()
			r.Delete("/users/{id}", h.DeleteUser)

			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userIDParam, nil)
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
