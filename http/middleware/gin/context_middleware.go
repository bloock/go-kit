package gin

import (
	"errors"
	"github.com/bloock/go-kit/context"
	"github.com/bloock/go-kit/domain"
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

func ContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP, _ := findClientIP(ctx)
		if clientIP != "" {
			ctx.Set(context.ClientIPKey, clientIP)
		}
		ctx.Set(context.RequestIDKey, domain.GenUUID())
		ctx.Next()
	}
}

func findClientIP(ctx *gin.Context) (string, error) {
	ips := ctx.Request.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}
