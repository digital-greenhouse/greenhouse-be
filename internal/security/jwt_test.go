package security

import (
	"os"
	"testing"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

func TestGetSecretKey(t *testing.T) {
	origSecret, exists := os.LookupEnv("JWT_SECRET")
	defer func() {
		if exists {
			os.Setenv("JWT_SECRET", origSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	// Test default
	os.Unsetenv("JWT_SECRET")
	key := getSecretKey()
	if string(key) != "default_secret_key_for_development_only" {
		t.Errorf("expected default secret key, got %s", string(key))
	}

	// Test custom
	os.Setenv("JWT_SECRET", "custom_secret_key_for_testing")
	key = getSecretKey()
	if string(key) != "custom_secret_key_for_testing" {
		t.Errorf("expected custom_secret_key_for_testing, got %s", string(key))
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	user := &domain.User{
		ID:    42,
		Name:  "Alice",
		Email: "alice@example.com",
		Role:  domain.RoleClient,
	}

	// Generate token
	tokenStr, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token string")
	}

	// Validate token
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	// Verify claims
	subVal, ok := claims["sub"].(float64)
	if !ok {
		t.Errorf("expected sub claim to be float64, got %T (%v)", claims["sub"], claims["sub"])
	} else if uint(subVal) != user.ID {
		t.Errorf("expected sub claim to be %d, got %v", user.ID, subVal)
	}

	nameVal, ok := claims["name"].(string)
	if !ok || nameVal != user.Name {
		t.Errorf("expected name claim to be %q, got %q", user.Name, nameVal)
	}

	emailVal, ok := claims["email"].(string)
	if !ok || emailVal != user.Email {
		t.Errorf("expected email claim to be %q, got %q", user.Email, emailVal)
	}

	roleVal, ok := claims["role"].(string)
	if !ok || roleVal != string(user.Role) {
		t.Errorf("expected role claim to be %q, got %q", string(user.Role), roleVal)
	}
}

func TestValidateTokenErrors(t *testing.T) {
	// 1. Expired Token
	t.Run("expired token", func(t *testing.T) {
		expiredClaims := jwt.MapClaims{
			"sub":   uint(42),
			"name":  "Alice",
			"email": "alice@example.com",
			"role":  string(domain.RoleClient),
			"exp":   time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
		}
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		expiredTokenStr, err := expiredToken.SignedString(getSecretKey())
		if err != nil {
			t.Fatalf("failed to sign expired token: %v", err)
		}

		_, err = ValidateToken(expiredTokenStr)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})

	// 2. Invalid Token String / Malformed
	t.Run("malformed token string", func(t *testing.T) {
		_, err := ValidateToken("invalid.token.string")
		if err == nil {
			t.Error("expected error for malformed token string, got nil")
		}
	})

	// 3. Token Signed with Invalid Method
	t.Run("invalid signing method (none)", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":   uint(42),
			"name":  "Alice",
			"email": "alice@example.com",
			"role":  string(domain.RoleClient),
			"exp":   time.Now().Add(time.Hour).Unix(),
		}
		noneToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		noneTokenStr, err := noneToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		if err != nil {
			t.Fatalf("failed to sign token with None method: %v", err)
		}

		_, err = ValidateToken(noneTokenStr)
		if err == nil {
			t.Error("expected error for token with invalid signing method, got nil")
		}
	})
}
