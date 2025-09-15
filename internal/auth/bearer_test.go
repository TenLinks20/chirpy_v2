package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken_Success(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer sometoken123")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token != "sometoken123" {
		t.Errorf("expected sometoken123, got %s", token)
	}
}

func TestGetBearerToken_MissingHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for missing Authorization header, got nil")
	}
}

func TestGetBearerToken_BadFormat(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Token sometoken123")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for wrong scheme, got nil")
	}
}

func TestGetBearerToken_NoToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for missing token value, got nil")
	}
}
