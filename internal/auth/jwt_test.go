package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// helper to make a test token
func makeTestToken(t *testing.T, userID uuid.UUID, secret string, expiresIn time.Duration) string {
	t.Helper()

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "test-suite",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

func TestValidateJWT_Success(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()
	token := makeTestToken(t, userID, secret, time.Minute)

	uid, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if uid != userID {
		t.Errorf("expected %v, got %v", userID, uid)
	}
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()
	token := makeTestToken(t, userID, secret, time.Minute)

	_, err := ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Fatal("expected error for invalid secret, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()
	token := makeTestToken(t, userID, secret, -time.Minute) // already expired

	_, err := ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}
