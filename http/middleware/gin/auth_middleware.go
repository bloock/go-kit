package gin

import (
	"errors"
	"fmt"
	"github.com/bloock/go-kit/auth"
	bloockCtx "github.com/bloock/go-kit/context"
	httpError "github.com/bloock/go-kit/errors"
	bloockHttp "github.com/bloock/go-kit/http"
	"github.com/bloock/go-kit/observability"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"time"
)

type AuthMiddleware struct {
	httpClient bloockHttp.HttpRequest
	authHost   string
	logger     observability.Logger
}

func NewAuthMiddlewareEntity(httpClient bloockHttp.HttpRequest, authHost string, l observability.Logger) AuthMiddleware {
	l.UpdateLogger(l.With().Caller().Str("component", "auth-middleware").Logger())

	return AuthMiddleware{
		httpClient: httpClient,
		authHost:   authHost,
		logger:     l,
	}
}

type CredentialAuthResponse struct {
	JWT string `json:"jwt"`
}

func (a AuthMiddleware) Authorize(ability auth.Ability) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		credAuthResp, err := a.GetCredentialsAuthenticate(ctx)
		if err != nil {
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		var claims auth.JWTClaims
		err = auth.DecodeJWTUnverified(credAuthResp.JWT, &claims)
		if err != nil {
			_ = ctx.Error(httpError.ErrUnauthorized(errors.New("invalid token content provided")))
			ctx.Abort()
			return
		}
		ctx.Set(bloockCtx.UserIDKey, claims.ClientID)
		ctx.Set(bloockCtx.AuthTokenKey, credAuthResp.JWT)

		if allowedActions := claims.Scopes[ability.Resource()]; allowedActions != nil {
			for _, a := range allowedActions {
				if a == ability.Action() {
					ctx.Next()
					return
				}
			}
		}

		_ = ctx.Error(httpError.ErrForbidden(errors.New("action forbidden")))
		ctx.Abort()
	}
}

func (a AuthMiddleware) GetCredentialsAuthenticate(ctx *gin.Context) (CredentialAuthResponse, error) {
	requestID, ok := ctx.Get(bloockCtx.RequestIDKey)
	if !ok {
		err := httpError.ErrUnexpected(errors.New("request id not found"))
		a.logger.Info(ctx).Err(err).Msg("")
		return CredentialAuthResponse{}, err
	}

	url := fmt.Sprintf("%s/v1/authenticate", a.authHost)

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ctx.Request.Header.Set(bloockCtx.RequestIDKey, requestID.(string))

	var response CredentialAuthResponse
	err := a.httpClient.GetWithHeaders(ctxTimeout, url, &response, ctx.Request.Header)
	if err != nil {
		a.logger.Info(ctx).Err(err).Msg("")
		return CredentialAuthResponse{}, err
	}

	return response, nil
}
