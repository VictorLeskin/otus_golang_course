package internalhttp

import (
	"net/http"
	"time"

	"calendar/internal/logger"
)

func LoggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log the start
			start := time.Now()
			logger.Debugf("Request started method: %s path: %s ip: %s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
			)

			// call handlers
			next.ServeHTTP(w, r)

			duration := time.Since(start)

			// Log the end
			logger.Infof("Request completed method: %s path: %s ip: %s latency:%d",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				duration,
			)
		})
	}
}
