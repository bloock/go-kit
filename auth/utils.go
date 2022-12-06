package auth

import (
	"fmt"
	http2 "github.com/bloock/go-kit/errors/http"
	"net/http"
)

func GetClientIDFromToken(token string) (string, error) {
	jwtToken := GetBearerToken(token)
	var claims JWTClaims
	err := DecodeJWTUnverified(jwtToken, &claims)
	if err != nil {
		return "", err
	}
	userId := claims.ClientID
	if userId == "" {
		appError := http2.NewHttpAppError(http.StatusUnauthorized, fmt.Errorf("couldn't get client ID from authentication token").Error())
		return "", appError
	}

	return userId, nil
}
