package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bloock/go-kit/auth"
	bloockCtx "github.com/bloock/go-kit/context"
	httpError "github.com/bloock/go-kit/errors"
	"github.com/bloock/go-kit/observability"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"time"
)

type AuthMiddleware interface {
	GetCredentialsAuthenticate(ctx *gin.Context) (CredentialAuthResponse, error)
}

type AuthMiddlewareEntity struct {
	httpClient http.Client
	authHost   string
	logger     observability.Logger
}

func NewAuthMiddlewareEntity(authHost string, l observability.Logger) AuthMiddlewareEntity {
	l.UpdateLogger(l.With().Caller().Str("component", "auth-middleware").Logger())

	return AuthMiddlewareEntity{
		httpClient: http.Client{},
		authHost:   authHost,
		logger:     l,
	}
}

type CredentialAuthResponse struct {
	JWT string `json:"jwt"`
}

type CredentialAuthErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (a AuthMiddlewareEntity) Authorize(ability auth.Ability) gin.HandlerFunc {
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

func (a AuthMiddlewareEntity) GetCredentialsAuthenticate(ctx *gin.Context) (CredentialAuthResponse, error) {
	requestID, ok := ctx.Get(bloockCtx.RequestIDKey)
	if !ok {
		err := httpError.ErrUnexpected(errors.New("request id not found"))
		a.logger.Info(ctx).Err(err).Msg("")
		return CredentialAuthResponse{}, err
	}

	url := fmt.Sprintf("%s/credentials/v1/authenticate", a.authHost)

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxTimeout, http.MethodGet, url, nil)
	if err != nil {
		err = httpError.ErrUnexpected(err)
		a.logger.Info(ctx).Err(err).Msg("")
		return CredentialAuthResponse{}, err
	}
	req.Header = ctx.Request.Header
	req.Header.Set(bloockCtx.RequestIDKey, requestID.(string))

	resp, err := a.httpClient.Do(req)
	if err != nil {
		err = httpError.ErrUnexpected(err)
		a.logger.Info(ctx).Err(err).Msg("")
		return CredentialAuthResponse{}, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		err = httpError.ErrUnexpected(err)
		a.logger.Info(ctx).Err(err).Msg("")
		return CredentialAuthResponse{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorMessage string
		var respError CredentialAuthErrorResponse
		err = json.Unmarshal(respByte, &respError)
		if err != nil {
			errorMessage = string(respByte)
		}
		errorMessage = respError.Message
		err = httpError.NewHttpAppError(resp.StatusCode, errorMessage)
		a.logger.Info(ctx).Err(err).Msg("")
		return CredentialAuthResponse{}, err
	}

	var response CredentialAuthResponse
	err = json.Unmarshal(respByte, &response)
	if err != nil {
		err = httpError.ErrUnexpected(err)
		a.logger.Info(ctx).Err(err).Msg("")
		return CredentialAuthResponse{}, err
	}

	return response, nil
}
