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

	metadata := make(map[string]interface{})
	metadata["name"] = "basic"
	productClaims := NewProductClaims("product_id", metadata)

	t.Run("Given a correct claim should not return err", func(t *testing.T) {
		jc := NewJWTClaim(
			expire,
			issued,
			notBefore,
			"5fe1ff5d-dd31-496a-9dfa-c95cc7847df8",
			[]ProductClaims{productClaims},
			"Joe", "Doe", "joe@doe.com",
			map[string][]string{
				"foo.bar": {"create"},
			},
			"free_plan",
		)

		token, err := NewJWT(jc, "secret")
		assert.NoError(t, err)
		assert.Equal(t, token, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwcm9kdWN0IjpbeyJpZCI6InByb2R1Y3RfaWQiLCJtZXRhZGF0YSI6eyJuYW1lIjoiYmFzaWMifX1dLCJ1c2VyIjp7Im5hbWUiOiJKb2UiLCJzdXJuYW1lIjoiRG9lIiwiZW1haWwiOiJqb2VAZG9lLmNvbSIsInVzZXJfcm9sZSI6ImZyZWVfcGxhbiJ9LCJzY29wZXMiOnsiZm9vLmJhciI6WyJjcmVhdGUiXX19.HSCVfLRjAv49LGTHQeozs8n7LJamaqeKTd7AXjuwonw")
	})

	t.Run("Given a valid claim should parse jwt claims successfully", func(t *testing.T) {
		jc := NewJWTClaim(
			expire,
			issued,
			notBefore,
			"5fe1ff5d-dd31-496a-9dfa-c95cc7847df8",
			[]ProductClaims{productClaims},
			"Joe", "Doe", "joe@doe.com",
			map[string][]string{
				"foo.bar": {"create"},
			},
			"disabled",
		)

		var claims JWTClaims
		err := DecodeJWT(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwbGFuIjp7ImlkIjoiOTMwYjVmMTAtZjQ1Ny00YWQ4LTk5NzctNTZjMzhmMmYxYWE5IiwibWV0YWRhdGEiOnsibWF4X2FwaV9rZXlzIjoxMCwibWF4X3N1YnNjcmlwdGlvbl9yZWNvcmRzIjoyMDAwLCJzY29wZSI6ImxpdmUsdGVzdCJ9fSwidXNlciI6eyJuYW1lIjoiSm9lIiwic3VybmFtZSI6IkRvZSIsImVtYWlsIjoiam9lQGRvZS5jb20iLCJhY3RpdmF0ZWQiOnRydWUsInZlcmlmaWVkIjp0cnVlfSwic2NvcGVzIjp7ImZvby5iYXIiOlsiY3JlYXRlIl19fQ.BYr_ac0VB4-5CM4n9rZR9idTUTyn-TnRq4uK9Gh6WfM",
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

	t.Run("Given an invalid token should return false", func(t *testing.T) {
		token := "cyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwbGFuIjp7ImlkIjoiOTMwYjVmMTAtZjQ1Ny00YWQ4LTk5NzctNTZjMzhmMmYxYWE5IiwibWV0YWRhdGEiOnsibWF4X2FwaV9rZXlzIjoxMCwibWF4X3N1YnNjcmlwdGlvbl9yZWNvcmRzIjoyMDAwLCJzY29wZSI6ImxpdmUsdGVzdCJ9fSwidXNlciI6eyJuYW1lIjoiSm9lIiwic3VybmFtZSI6IkRvZSIsImVtYWlsIjoiam9lQGRvZS5jb20iLCJhY3RpdmF0ZWQiOnRydWV9LCJzY29wZXMiOnsiZm9vLmJhciI6WyJjcmVhdGUiXX19.gp1GL7IigMGZ17iyJAReK-AxNNhPXTm7A7cY4rZfrmY"

		assert.False(t, ValidJWT(token, "secret"))
	})

	t.Run("Given a token with bearer prefix should return its token string", func(t *testing.T) {
		token := "Bearer cyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwbGFuIjp7ImlkIjoiOTMwYjVmMTAtZjQ1Ny00YWQ4LTk5NzctNTZjMzhmMmYxYWE5IiwibWV0YWRhdGEiOnsibWF4X2FwaV9rZXlzIjoxMCwibWF4X3N1YnNjcmlwdGlvbl9yZWNvcmRzIjoyMDAwLCJzY29wZSI6ImxpdmUsdGVzdCJ9fSwidXNlciI6eyJuYW1lIjoiSm9lIiwic3VybmFtZSI6IkRvZSIsImVtYWlsIjoiam9lQGRvZS5jb20iLCJhY3RpdmF0ZWQiOnRydWV9LCJzY29wZXMiOnsiZm9vLmJhciI6WyJjcmVhdGUiXX19.gp1GL7IigMGZ17iyJAReK-AxNNhPXTm7A7cY4rZfrmY"
		expectedTokenString := "cyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwbGFuIjp7ImlkIjoiOTMwYjVmMTAtZjQ1Ny00YWQ4LTk5NzctNTZjMzhmMmYxYWE5IiwibWV0YWRhdGEiOnsibWF4X2FwaV9rZXlzIjoxMCwibWF4X3N1YnNjcmlwdGlvbl9yZWNvcmRzIjoyMDAwLCJzY29wZSI6ImxpdmUsdGVzdCJ9fSwidXNlciI6eyJuYW1lIjoiSm9lIiwic3VybmFtZSI6IkRvZSIsImVtYWlsIjoiam9lQGRvZS5jb20iLCJhY3RpdmF0ZWQiOnRydWV9LCJzY29wZXMiOnsiZm9vLmJhciI6WyJjcmVhdGUiXX19.gp1GL7IigMGZ17iyJAReK-AxNNhPXTm7A7cY4rZfrmY"

		tokenString := GetBearerToken(token)

		assert.Equal(t, expectedTokenString, tokenString)
	})

	t.Run("Given a token with basic prefix should return an empty token string", func(t *testing.T) {
		token := "Basic cyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTgxNjcxMjYsIm5iZiI6MTQxNTc5MjcyNiwiaWF0IjoxNDE1NzkyNzI2LCJjbGllbnRfaWQiOiI1ZmUxZmY1ZC1kZDMxLTQ5NmEtOWRmYS1jOTVjYzc4NDdkZjgiLCJwbGFuIjp7ImlkIjoiOTMwYjVmMTAtZjQ1Ny00YWQ4LTk5NzctNTZjMzhmMmYxYWE5IiwibWV0YWRhdGEiOnsibWF4X2FwaV9rZXlzIjoxMCwibWF4X3N1YnNjcmlwdGlvbl9yZWNvcmRzIjoyMDAwLCJzY29wZSI6ImxpdmUsdGVzdCJ9fSwidXNlciI6eyJuYW1lIjoiSm9lIiwic3VybmFtZSI6IkRvZSIsImVtYWlsIjoiam9lQGRvZS5jb20iLCJhY3RpdmF0ZWQiOnRydWV9LCJzY29wZXMiOnsiZm9vLmJhciI6WyJjcmVhdGUiXX19.gp1GL7IigMGZ17iyJAReK-AxNNhPXTm7A7cY4rZfrmY"

		tokenString := GetBearerToken(token)

		assert.Equal(t, "", tokenString)
	})
}
