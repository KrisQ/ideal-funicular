package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	length, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	if length != 32 {
		return "", errors.New("failed to randomize a 32 bytes")
	}
	encodedStr := hex.EncodeToString(b)
	return encodedStr, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return parsedID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("empty authorization header")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("empty authorization header")
	}

	tokenString := authHeader[len(bearerPrefix):]
	return tokenString, nil
}
