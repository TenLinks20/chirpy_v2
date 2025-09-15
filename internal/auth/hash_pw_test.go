package auth

import (
	"testing"
)

func TestHashPasswordAndCheck(t *testing.T) {
	password := "supersecret123"

	// Hash the password
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	// Check that the hash is not empty and looks valid
	if hash == "" {
		t.Fatal("expected a non-empty hash")
	}

	// bcrypt hashes should always start with $2a, $2b, or $2y
	if len(hash) < 4 || hash[0:2] != "$2" {
		t.Errorf("unexpected hash format: %s", hash)
	}

	// Correct password should pass
	if err := CheckHashPassword(hash, password); err != nil {
		t.Errorf("expected password to match, got error: %v", err)
	}

	// Wrong password should fail
	if err := CheckHashPassword(hash, "wrongpassword"); err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}