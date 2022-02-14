package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type planClaims struct {
	ID       string                 `json:"id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type userClaims struct {
	Name      string `json:"name,omitempty"`
	Surname   string `json:"surname,omitempty"`
	Email     string `json:"email,omitempty"`
	Activated bool   `json:"activated,omitempty"`
	Verified bool `json:"verified,omitempty"`
	Deleted   bool   `json:"deleted,omitempty"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	ClientID string              `json:"client_id,omitempty"`
	Plan     planClaims          `json:"plan,omitempty"`
	User     userClaims          `json:"user,omitempty"`
	Scopes   map[string][]string `json:"scopes,omitempty"`
}

func (c JWTClaims) Valid() error {
	return nil
}

func NewJWTClaim(expiresAt, issuedAt, notBefore time.Time, clientID string, planID string, planMetadata map[string]interface{}, userName string, userSurname string, userEmail string, userActivated, userDeleted bool, scopes map[string][]string, userVerified bool) JWTClaims {
	return JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(notBefore),
		},
		ClientID: clientID,
		Plan: planClaims{
			ID:       planID,
			Metadata: planMetadata,
		},
		User: userClaims{
			Name:      userName,
			Surname:   userSurname,
			Email:     userEmail,
			Activated: userActivated,
			Verified: userVerified,
			Deleted:   userDeleted,
		},
		Scopes: scopes,
	}
}

func NewJWT(claims JWTClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidJWT(tokenString string, secret string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		log.Println(err.Error())
		return false
	}

	return token.Valid
}

func DecodeJWT(tokenString string, secret string, claims *JWTClaims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if _, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return nil
	} else {
		return err
	}
}
