package context

import (
	"context"
)

var UserIDKey = "X-User-ID"
var RequestIDKey = "X-Request-ID"
var ClientIPKey = "X-Client-IP"
var AuthTokenKey = "X-Auth-JWT"

func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}

func GetClientIP(ctx context.Context) string {
	clientIP, ok := ctx.Value(ClientIPKey).(string)
	if !ok {
		return ""
	}
	return clientIP
}
