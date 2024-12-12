package chi

import (
	context2 "github.com/bloock/go-kit/context"
	"github.com/bloock/go-kit/observability"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"time"
)

func Logger(logger observability.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()

			l := logger.With().Timestamp().Logger()

			clientIp := ""
			cip, ok := ctx.Value(context2.ClientIPKey).(string)
			if ok {
				clientIp = cip
			}

			defer func() {
				l = l.With().
					Int("status", ww.Status()).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("ip", clientIp).
					Dur("latency", time.Since(start)).
					Str("user-agent", r.UserAgent()).
					Logger()

				switch {
				case ww.Status() >= http.StatusBadRequest && ww.Status() < http.StatusInternalServerError:
					{
						l.WithLevel(zerolog.WarnLevel).
							Msg("")
					}
				case ww.Status() >= http.StatusInternalServerError:
					{
						l.WithLevel(zerolog.ErrorLevel).
							Msg("")
					}
				default:
					l.WithLevel(zerolog.InfoLevel).
						Msg("")
				}
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
