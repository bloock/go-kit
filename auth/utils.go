package auth

import (
	"fmt"
	"net/http"

	"github.com/bloock/go-kit/errors"
)

func GetClientIDFromToken(token string) (string, error) {

	var claims JWTClaims
	err := DecodeJWTUnverified(token, &claims)
	if err != nil {
		return "", err
	}
	userId := claims.ClientID
	if userId == "" {
		appError := errors.NewHttpAppError(http.StatusUnauthorized, fmt.Errorf("couldn't get client ID from authentication token").Error())
		return "", appError
	}

	return userId, nil
}
