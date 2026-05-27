package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
	"digital-greenhouse/greenhouse-be/internal/security"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKeyForTest() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_key_for_development_only"
	}
	return []byte(secret)
}

func generateCustomToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKeyForTest())
}

func TestAuthMiddleware(t *testing.T) {
	// Generate a valid token
	validUser := &domain.User{
		ID:    42,
		Name:  "Alice",
		Email: "alice@example.com",
		Role:  domain.RoleClient,
	}
	validToken, err := security.GenerateToken(validUser)
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	// Generate an expired token
	expiredClaims := jwt.MapClaims{
		"sub":   float64(42),
		"name":  "Alice",
		"email": "alice@example.com",
		"role":  string(domain.RoleClient),
		"exp":   time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken, err := generateCustomToken(expiredClaims)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	// Generate a token with missing "sub" claim
	missingSubClaims := jwt.MapClaims{
		"name":  "Alice",
		"email": "alice@example.com",
		"role":  string(domain.RoleClient),
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	missingSubToken, err := generateCustomToken(missingSubClaims)
	if err != nil {
		t.Fatalf("failed to generate missing sub token: %v", err)
	}

	// Generate a token with invalid type for "sub" claim
	invalidSubClaims := jwt.MapClaims{
		"sub":   "not-a-number",
		"name":  "Alice",
		"email": "alice@example.com",
		"role":  string(domain.RoleClient),
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	invalidSubToken, err := generateCustomToken(invalidSubClaims)
	if err != nil {
		t.Fatalf("failed to generate invalid sub token: %v", err)
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectHandler  bool
		expectedUserID uint
		expectedRole   string
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectHandler:  true,
			expectedUserID: 42,
			expectedRole:   string(domain.RoleClient),
		},
		{
			name:           "Missing Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectHandler:  false,
		},
		{
			name:           "Invalid prefix (no Bearer)",
			authHeader:     "Token " + validToken,
			expectedStatus: http.StatusUnauthorized,
			expectHandler:  false,
		},
		{
			name:           "Invalid prefix (too many parts)",
			authHeader:     "Bearer token extra",
			expectedStatus: http.StatusUnauthorized,
			expectHandler:  false,
		},
		{
			name:           "Invalid prefix (too few parts)",
			authHeader:     "Bearer",
			expectedStatus: http.StatusUnauthorized,
			expectHandler:  false,
		},
		{
			name:           "Malformed token",
			authHeader:     "Bearer malformed-token-string",
			expectedStatus: http.StatusUnauthorized,
			expectHandler:  false,
		},
		{
			name:           "Expired token",
			authHeader:     "Bearer " + expiredToken,
			expectedStatus: http.StatusUnauthorized,
			expectHandler:  false,
		},
		{
			name:           "Missing userID (sub)",
			authHeader:     "Bearer " + missingSubToken,
			expectedStatus: http.StatusUnauthorized,
			expectHandler:  false,
		},
		{
			name:           "Invalid userID type",
			authHeader:     "Bearer " + invalidSubToken,
			expectedStatus: http.StatusUnauthorized,
			expectHandler:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				userID := GetUserID(r.Context())
				role := GetUserRole(r.Context())
				if userID != tt.expectedUserID {
					t.Errorf("expected userID %d, got %d", tt.expectedUserID, userID)
				}
				if role != tt.expectedRole {
					t.Errorf("expected role %q, got %q", tt.expectedRole, role)
				}
				w.WriteHeader(http.StatusOK)
			})

			middlewareToTest := AuthMiddleware(nextHandler)

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			middlewareToTest.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if handlerCalled != tt.expectHandler {
				t.Errorf("expected handlerCalled to be %t, got %t", tt.expectHandler, handlerCalled)
			}
		})
	}
}

func TestOptionalAuthMiddleware(t *testing.T) {
	// Generate a valid token
	validUser := &domain.User{
		ID:    42,
		Name:  "Alice",
		Email: "alice@example.com",
		Role:  domain.RoleClient,
	}
	validToken, err := security.GenerateToken(validUser)
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	// Generate an expired token
	expiredClaims := jwt.MapClaims{
		"sub":   float64(42),
		"name":  "Alice",
		"email": "alice@example.com",
		"role":  string(domain.RoleClient),
		"exp":   time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken, err := generateCustomToken(expiredClaims)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	// Generate a token with missing "sub" claim
	missingSubClaims := jwt.MapClaims{
		"name":  "Alice",
		"email": "alice@example.com",
		"role":  string(domain.RoleClient),
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	missingSubToken, err := generateCustomToken(missingSubClaims)
	if err != nil {
		t.Fatalf("failed to generate missing sub token: %v", err)
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedUserID uint
		expectedRole   string
	}{
		{
			name:           "Valid token (Success)",
			authHeader:     "Bearer " + validToken,
			expectedUserID: 42,
			expectedRole:   string(domain.RoleClient),
		},
		{
			name:           "Missing Authorization header (Pass-through)",
			authHeader:     "",
			expectedUserID: 0,
			expectedRole:   "",
		},
		{
			name:           "Invalid prefix (Pass-through)",
			authHeader:     "Token " + validToken,
			expectedUserID: 0,
			expectedRole:   "",
		},
		{
			name:           "Invalid prefix - parts count (Pass-through)",
			authHeader:     "Bearer",
			expectedUserID: 0,
			expectedRole:   "",
		},
		{
			name:           "Invalid/expired token (Pass-through)",
			authHeader:     "Bearer " + expiredToken,
			expectedUserID: 0,
			expectedRole:   "",
		},
		{
			name:           "Missing sub claim (Pass-through)",
			authHeader:     "Bearer " + missingSubToken,
			expectedUserID: 0,
			expectedRole:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				userID := GetUserID(r.Context())
				role := GetUserRole(r.Context())
				if userID != tt.expectedUserID {
					t.Errorf("expected userID %d, got %d", tt.expectedUserID, userID)
				}
				if role != tt.expectedRole {
					t.Errorf("expected role %q, got %q", tt.expectedRole, role)
				}
				w.WriteHeader(http.StatusOK)
			})

			middlewareToTest := OptionalAuthMiddleware(nextHandler)

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			middlewareToTest.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
			}

			if !handlerCalled {
				t.Error("expected handler to be called")
			}
		})
	}
}

func TestGetUserIDAndRoleFromContext(t *testing.T) {
	// 1. Missing values
	ctxEmpty := context.Background()
	if uid := GetUserID(ctxEmpty); uid != 0 {
		t.Errorf("expected 0, got %d", uid)
	}
	if role := GetUserRole(ctxEmpty); role != "" {
		t.Errorf("expected empty string, got %q", role)
	}

	// 2. Invalid types
	ctxInvalid := context.WithValue(context.Background(), UserIDKey, "invalid-type")
	ctxInvalid = context.WithValue(ctxInvalid, UserRoleKey, 123)

	if uid := GetUserID(ctxInvalid); uid != 0 {
		t.Errorf("expected 0 for invalid type, got %d", uid)
	}
	if role := GetUserRole(ctxInvalid); role != "" {
		t.Errorf("expected empty string for invalid type, got %q", role)
	}
}
