package middleware

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

// CustomClaimsExample contains custom data we want from the token.
type CustomClaims struct {
	OrgId       string   `json:"org_id"`
	Permissions []string `json:"permissions"`
}

func (c *CustomClaims) Validate(ctx context.Context) error {
	return nil
}

type auth0Middleware struct {
	middleware *jwtmiddleware.JWTMiddleware
}

func NewAuth0Middleware(issuer string, audience string) *auth0Middleware {
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return &auth0Middleware{
		middleware,
	}
}

func (m *auth0Middleware) AuthorizeHandler(permissions []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r

			claims, ok := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			if !ok {
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					map[string]string{"message": "Failed to get validated JWT claims."},
				)
				return
			}

			customClaims, ok := claims.CustomClaims.(*CustomClaims)
			if !ok {
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					map[string]string{"message": "Failed to cast custom JWT claims to specific type."},
				)
				return
			}

			if !tokenHasPermission(customClaims.Permissions, permissions) {
				ctx.AbortWithStatusJSON(
					http.StatusForbidden,
					map[string]string{
						"message": "You don't have the required permissions to access this resource",
						"error":   "forbidden",
					},
				)
				return
			}

			ctx.Next()
		}

		m.middleware.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"message": "JWT is invalid."},
			)
		}
	}
}

func tokenHasPermission(userPermissions []string, expectedPermissions []string) bool {
	// If no permissions are required, allow access
	if len(expectedPermissions) == 0 {
		return true
	}

	// If user has no permissions but some are required, deny access
	if len(userPermissions) == 0 {
		return false
	}

	permissionMap := make(map[string]bool)
	for _, perm := range userPermissions {
		permissionMap[perm] = true
	}

	for _, expectedPerm := range expectedPermissions {
		if permissionMap[expectedPerm] {
			return true
		}
	}

	return false
}
