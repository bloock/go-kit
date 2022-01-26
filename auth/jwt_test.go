package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestJwt(t *testing.T) {

	t.Run("Given a correct claim should not return err", func(t *testing.T) {
		jc := JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{},
			UserID:           "5fe1ff5d-dd31-496a-9dfa-c95cc7847df8",
			Plan: PlanClaims{
				ID: "930b5f10-f457-4ad8-9977-56c38f2f1aa9",
				Metadata: PlanMetadataClaims{
					Scope:                  []string{"live", "test"},
					MaxSubscriptionRecords: 2000,
					MaxApiKeys:             10,
				},
			},
			User: UserClaims{
				Name:      "Joe",
				Surname:   "Doe",
				Email:     "joe@doe.com",
				Activated: true,
				Deleted:   false,
			},
			Scopes: map[string][]string{
				"foo.bar": {"create"},
			},
		}

		token, err := NewJWT(jc, "secret")
		assert.NoError(t, err)
		assert.Equal(t, token, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNWZlMWZmNWQtZGQzMS00OTZhLTlkZmEtYzk1Y2M3ODQ3ZGY4IiwicGxhbiI6eyJpZCI6IjkzMGI1ZjEwLWY0NTctNGFkOC05OTc3LTU2YzM4ZjJmMWFhOSIsIm1ldGFkYXRhIjp7InNjb3BlIjpbImxpdmUiLCJ0ZXN0Il0sIm1heF9zdWJzY3JpcHRpb25fcmVjb3JkcyI6MjAwMCwibWF4X2FwaV9rZXlzIjoxMH19LCJ1c2VyIjp7Im5hbWUiOiJKb2UiLCJzdXJuYW1lIjoiRG9lIiwiZW1haWwiOiJqb2VAZG9lLmNvbSIsImFjdGl2YXRlZCI6dHJ1ZX0sInNjb3BlcyI6eyJmb28uYmFyIjpbImNyZWF0ZSJdfX0.kP9nC4kuROYoTF6igzTzOJhfN7btg2v9q2Conkka9hM")
	})

	t.Run("Given a valid claim should parse jwt claims successfully", func(t *testing.T) {
		jc := JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{},
			UserID:           "5fe1ff5d-dd31-496a-9dfa-c95cc7847df8",
			Plan: PlanClaims{
				ID: "930b5f10-f457-4ad8-9977-56c38f2f1aa9",
				Metadata: PlanMetadataClaims{
					Scope:                  []string{"live", "test"},
					MaxSubscriptionRecords: 2000,
					MaxApiKeys:             10,
				},
			},
			User: UserClaims{
				Name:      "Joe",
				Surname:   "Doe",
				Email:     "joe@doe.com",
				Activated: true,
				Deleted:   false,
			},
			Scopes: map[string][]string{
				"foo.bar": {"create"},
			},
		}

		var claims JWTClaims
		err := DecodeJWT(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNWZlMWZmNWQtZGQzMS00OTZhLTlkZmEtYzk1Y2M3ODQ3ZGY4IiwicGxhbiI6eyJpZCI6IjkzMGI1ZjEwLWY0NTctNGFkOC05OTc3LTU2YzM4ZjJmMWFhOSIsIm1ldGFkYXRhIjp7InNjb3BlIjpbImxpdmUiLCJ0ZXN0Il0sIm1heF9zdWJzY3JpcHRpb25fcmVjb3JkcyI6MjAwMCwibWF4X2FwaV9rZXlzIjoxMH19LCJ1c2VyIjp7Im5hbWUiOiJKb2UiLCJzdXJuYW1lIjoiRG9lIiwiZW1haWwiOiJqb2VAZG9lLmNvbSIsImFjdGl2YXRlZCI6dHJ1ZX0sInNjb3BlcyI6eyJmb28uYmFyIjpbImNyZWF0ZSJdfX0.kP9nC4kuROYoTF6igzTzOJhfN7btg2v9q2Conkka9hM",
			"secret",
			&claims,
		)
		assert.NoError(t, err)
		assert.Equal(t, claims, jc)
	})

	t.Run("Given an invalid secret should return error", func(t *testing.T) {
		jc := JWTClaims{}

		token, err := NewJWT(jc, "secret")
		assert.NoError(t, err)

		var claims JWTClaims

		err = DecodeJWT(
			token,
			"invalid_secret",
			&claims,
		)
		assert.Error(t, err)
	})

}
