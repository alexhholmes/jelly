package util

import (
	"context"
	"log/slog"
	"net/http"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			slogWithReq := slog.Default().With(
				"method", r.Method,
				"path", r.URL.Path,
				"remote", r.RemoteAddr,
				"host", r.Host,
				"agent", r.UserAgent(),
				"tls", r.TLS,
				"proto", r.Proto,
				"user", func() string {
					if c, err := r.Cookie("user"); err == nil {
						return c.Value
					}
					return ""
				}(),
			)
			slogWithReq.Info("Request")

			ctx := context.WithValue(r.Context(), ContextLogger, slogWithReq)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

// Recovery recovers from panics, logs the error, and sends a 500 status code.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log, ok := r.Context().Value(ContextLogger).(*slog.Logger)
					if !ok {
						log = slog.Default()
					}

					log.Error("Recovered from panic", "error", err)
					http.Error(
						w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}
