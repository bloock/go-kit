package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJwt(t *testing.T) {

	issued, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	expire, _ := time.Parse(time.RFC3339, "2099-11-12T11:45:26.371Z")
	notBefore, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")

	t.Run("Given a correct claim should not return err", func(t *testing.T) {
		jc := NewJWTClaim(
			expire,
			issued,
			notBefore,
			"5fe1ff5d-dd31-496a-9dfa-c95cc7847df8",
			"930b5f10-f457-4ad8-9977-56c38f2f1aa9",
			map[string]interface{}{
				"scope":                    "live,test",
				"max_subscription_records": 2000,
				"max_api_keys":             10,
			},
			"Joe", "Doe", "joe@doe.com", true, false,
			map[string][]string{
				"foo.bar": {"create"},
			},
		)

		token, err := NewJWT(jc, "secret")
		assert.NoError(t, err)
		assert.Equal(t, token, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwbGFuIjp7ImlkIjoiOTMwYjVmMTAtZjQ1Ny00YWQ4LTk5NzctNTZjMzhmMmYxYWE5IiwibWV0YWRhdGEiOnsibWF4X2FwaV9rZXlzIjoxMCwibWF4X3N1YnNjcmlwdGlvbl9yZWNvcmRzIjoyMDAwLCJzY29wZSI6ImxpdmUsdGVzdCJ9fSwidXNlciI6eyJuYW1lIjoiSm9lIiwic3VybmFtZSI6IkRvZSIsImVtYWlsIjoiam9lQGRvZS5jb20iLCJhY3RpdmF0ZWQiOnRydWV9LCJzY29wZXMiOnsiZm9vLmJhciI6WyJjcmVhdGUiXX19.gp1GL7IigMGZ17iyJAReK-AxNNhPXTm7A7cY4rZfrmY")
	})

	t.Run("Given a valid claim should parse jwt claims successfully", func(t *testing.T) {
		jc := NewJWTClaim(
			expire,
			issued,
			notBefore,
			"5fe1ff5d-dd31-496a-9dfa-c95cc7847df8",
			"930b5f10-f457-4ad8-9977-56c38f2f1aa9",
			map[string]interface{}{
				"scope":                    "live,test",
				"max_subscription_records": 2000,
				"max_api_keys":             10,
			},
			"Joe", "Doe", "joe@doe.com", true, false,
			map[string][]string{
				"foo.bar": {"create"},
			},
		)

		var claims JWTClaims
		err := DecodeJWT(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwbGFuIjp7ImlkIjoiOTMwYjVmMTAtZjQ1Ny00YWQ4LTk5NzctNTZjMzhmMmYxYWE5IiwibWV0YWRhdGEiOnsibWF4X2FwaV9rZXlzIjoxMCwibWF4X3N1YnNjcmlwdGlvbl9yZWNvcmRzIjoyMDAwLCJzY29wZSI6ImxpdmUsdGVzdCJ9fSwidXNlciI6eyJuYW1lIjoiSm9lIiwic3VybmFtZSI6IkRvZSIsImVtYWlsIjoiam9lQGRvZS5jb20iLCJhY3RpdmF0ZWQiOnRydWV9LCJzY29wZXMiOnsiZm9vLmJhciI6WyJjcmVhdGUiXX19.gp1GL7IigMGZ17iyJAReK-AxNNhPXTm7A7cY4rZfrmY",
			"secret",
			&claims,
		)
		assert.NoError(t, err)
		assert.NotEqual(t, claims, jc)
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

	t.Run("Given a valid token should return true", func(t *testing.T) {
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwbGFuIjp7ImlkIjoiOTMwYjVmMTAtZjQ1Ny00YWQ4LTk5NzctNTZjMzhmMmYxYWE5IiwibWV0YWRhdGEiOnsibWF4X2FwaV9rZXlzIjoxMCwibWF4X3N1YnNjcmlwdGlvbl9yZWNvcmRzIjoyMDAwLCJzY29wZSI6ImxpdmUsdGVzdCJ9fSwidXNlciI6eyJuYW1lIjoiSm9lIiwic3VybmFtZSI6IkRvZSIsImVtYWlsIjoiam9lQGRvZS5jb20iLCJhY3RpdmF0ZWQiOnRydWV9LCJzY29wZXMiOnsiZm9vLmJhciI6WyJjcmVhdGUiXX19.gp1GL7IigMGZ17iyJAReK-AxNNhPXTm7A7cY4rZfrmY"

		assert.True(t, ValidJWT(token, "secret"))
	})

	t.Run("Given an ivalid token should return false", func(t *testing.T) {
		token := "cyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwbGFuIjp7ImlkIjoiOTMwYjVmMTAtZjQ1Ny00YWQ4LTk5NzctNTZjMzhmMmYxYWE5IiwibWV0YWRhdGEiOnsibWF4X2FwaV9rZXlzIjoxMCwibWF4X3N1YnNjcmlwdGlvbl9yZWNvcmRzIjoyMDAwLCJzY29wZSI6ImxpdmUsdGVzdCJ9fSwidXNlciI6eyJuYW1lIjoiSm9lIiwic3VybmFtZSI6IkRvZSIsImVtYWlsIjoiam9lQGRvZS5jb20iLCJhY3RpdmF0ZWQiOnRydWV9LCJzY29wZXMiOnsiZm9vLmJhciI6WyJjcmVhdGUiXX19.gp1GL7IigMGZ17iyJAReK-AxNNhPXTm7A7cY4rZfrmY"

		assert.False(t, ValidJWT(token, "secret"))
	})

}
