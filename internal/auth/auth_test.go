package auth_test

import (
	"chirpy/internal/auth"
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	// Given
	password := "mySecret123"

	// When
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Then
	if hashedPassword == "" {
		t.Error("expected a non-empty hashed password")
	}

	// Now test CheckPasswordHash with a valid password
	err = auth.CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("expected no error when comparing correct password, got %v", err)
	}

	// Test CheckPasswordHash with an invalid password
	err = auth.CheckPasswordHash("wrongPassword", hashedPassword)
	if err == nil {
		t.Errorf("expected an error when comparing wrong password, got none")
	}
}
