package util

import (
	"errors"
	"net/http"
	"strings"
)

func GetTokenString(r *http.Request) (string, error) {
	// Check cookie
	cookie, err := r.Cookie("Authorization")
	if err == nil && cookie != nil {
		return strings.TrimPrefix(cookie.Value, "Bearer "), nil
	}

	// Check header
	token := r.Header.Get("authorization")
	if token != "" {
		return strings.TrimPrefix(token, "Bearer "), nil
	}

	// Check query string
	token = r.URL.Query().Get("Authorization")
	if token != "" {
		return strings.TrimPrefix(token, "Bearer "), nil
	}

	token = r.URL.Query().Get("authorization")
	if token != "" {
		return strings.TrimPrefix(token, "Bearer "), nil
	}

	// Return error
	return "", errors.New("no token found")
}
