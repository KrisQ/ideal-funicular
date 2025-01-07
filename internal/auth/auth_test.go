package auth_test

import (
	"chirpy/internal/auth"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeRefreshToken(t *testing.T) {
	str, err := auth.MakeRefreshToken()
	if err != nil {
		t.Errorf("exected a string got: %v", err)
	}
	// verifies the full round-trip 👀 encoding and decoding 🤪
	decoded, err := hex.DecodeString(str)
	if err != nil {
		t.Fatalf("failed to decode token from hex: %v", err)
	}

	if len(decoded) != 32 {
		t.Errorf("expected 32 bytes, got %d bytes", len(decoded))
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	password := "mySecret123"
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hashedPassword == "" {
		t.Error("expected a non-empty hashed password")
	}

	err = auth.CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("expected no error when comparing correct password, got %v", err)
	}

	err = auth.CheckPasswordHash("wrongPassword", hashedPassword)
	if err == nil {
		t.Errorf("expected an error when comparing wrong password, got none")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {

		userID, err := uuid.NewUUID()
		secret := []byte("mySecret123")
		if err != nil {
			t.Fatal(err)
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
		jwt, _ := auth.MakeJWT(userID, secret, 1*time.Millisecond)
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
