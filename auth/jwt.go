package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ProductClaims struct {
	ID       string                 `json:"id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func NewProductClaims(id string, metadata map[string]interface{}) ProductClaims {
	return ProductClaims{
		ID:       id,
		Metadata: metadata,
	}
}

type userClaims struct {
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Email    string `json:"email,omitempty"`
	UserRole string `json:"user_role"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	ClientID string              `json:"client_id,omitempty"`
	Products []ProductClaims     `json:"product,omitempty"`
	User     userClaims          `json:"user,omitempty"`
	Scopes   map[string][]string `json:"scopes,omitempty"`
}

func (c JWTClaims) Valid() error {
	return nil
}

func NewJWTClaim(expiresAt, issuedAt, notBefore time.Time, clientID string, products []ProductClaims, userName string, userSurname string, userEmail string, scopes map[string][]string, userRole string) JWTClaims {
	return JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(notBefore),
		},
		ClientID: clientID,
		Products: products,
		User: userClaims{
			Name:     userName,
			Surname:  userSurname,
			Email:    userEmail,
			UserRole: userRole,
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
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(*JWTClaims); ok {
		if token.Valid {
			return nil
		} else {
			return fmt.Errorf("invalid token provided")
		}
	} else {
		return fmt.Errorf("couldn't cast jwt claims")
	}
}

func DecodeJWTUnverified(tokenString string, claims *JWTClaims) error {
	tokenClaims, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return err
	}

	if _, ok := tokenClaims.Claims.(*JWTClaims); ok {
		return nil
	} else {
		return fmt.Errorf("couldn't cast jwt claims")
	}
}

func GetBearerToken(token string) string {
	splitToken := strings.Split(token, BEARER_PREFIX)

	if len(splitToken) != 2 {
		return ""
	}
	return strings.TrimSpace(splitToken[1])
}
