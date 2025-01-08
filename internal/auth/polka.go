package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetPolkaApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("empty authorization header")
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("empty authorization header")
	}

	key := authHeader[len(prefix):]
	return key, nil
}
