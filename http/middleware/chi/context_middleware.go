package chi

import (
	goctx "context"
	"errors"
	"github.com/bloock/go-kit/context"
	"github.com/bloock/go-kit/domain"
	"net"
	"net/http"
	"strings"
)

func ContextMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientIP, _ := findClientIP(r)
		if clientIP != "" {
			ctx = goctx.WithValue(ctx, context.ClientIPKey, clientIP)
		}

		ctx = goctx.WithValue(ctx, context.RequestIDKey, domain.GenUUID())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func findClientIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
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
