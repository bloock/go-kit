package auth

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v4"
)

type PlanMetadataClaims struct {
	Scope                  []string `json:"scope,omitempty"`
	MaxSubscriptionRecords int64    `json:"max_subscription_records,omitempty"`
	MaxApiKeys             int64    `json:"max_api_keys,omitempty"`
}

type PlanClaims struct {
	ID       string             `json:"id,omitempty"`
	Metadata PlanMetadataClaims `json:"metadata,omitempty"`
}

type UserClaims struct {
	Name      string `json:"name,omitempty"`
	Surname   string `json:"surname,omitempty"`
	Email     string `json:"email,omitempty"`
	Activated bool   `json:"activated,omitempty"`
	Deleted   bool   `json:"deleted,omitempty"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string              `json:"user_id,omitempty"`
	Plan   PlanClaims          `json:"plan,omitempty"`
	User   UserClaims          `json:"user,omitempty"`
	Scopes map[string][]string `json:"scopes,omitempty"`
}

func (c JWTClaims) Valid() error {
	return nil
}

func NewJWT(claims JWTClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func DecodeJWT(tokenString string, secret string, claims *JWTClaims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	log.Printf("%+v", token.Claims)

	if _, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return nil
	} else {
		return err
	}
}
