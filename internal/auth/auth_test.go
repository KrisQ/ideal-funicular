package auth_test

import (
	"chirpy/internal/auth"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestMakeAndValidateJWT(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {

		userID, err := uuid.NewUUID()
		secret := []byte("mySecret123") // Convert to byte slice
		if err != nil {
			t.Fatal(err) // Use Fatal instead of Error to stop test if UUID fails
		}
		jwt, err := auth.MakeJWT(userID, string(secret), 30*time.Second)
		if err != nil {
			t.Fatalf("expected to create a JWT: %v", err)
		}

		idFromJwt, err := auth.ValidateJWT(jwt, string(secret))
		if err != nil {
			t.Fatalf("expected to validate jwt: %v", err)
		}

		if userID != idFromJwt {
			t.Errorf("expected id from jwt to match user id: %v", err)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		userID, _ := uuid.NewUUID()
		secret := "mySecret123"
		// Create a token that expires very quickly
		jwt, _ := auth.MakeJWT(userID, secret, 1*time.Millisecond)
		// Wait for token to expire
		time.Sleep(2 * time.Millisecond)

		_, err := auth.ValidateJWT(jwt, secret)
		if err == nil {
			t.Error("expected error for expired token")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		userID, _ := uuid.NewUUID()
		jwt, _ := auth.MakeJWT(userID, "correctSecret", time.Hour)

		_, err := auth.ValidateJWT(jwt, "wrongSecret")
		if err == nil {
			t.Error("expected error for invalid secret")
		}
	})
}
