package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authValue := headers.Get("Authorization")
	if authValue == "" {
		return  "", errors.New("no auth value")
	}

	fields := strings.Fields(authValue)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return "", errors.New("incorrect token format")
	}

	return fields[1], nil
}