package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"calendar/internal/app"
	"calendar/internal/logger"
)

type Config struct {
	Host string
	Port string
}

type Server struct {
	addr       string
	httpServer *http.Server
	log        *logger.Logger
}

func NewServer(cfg Config, app *app.App) *Server {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return &Server{
		addr: addr,
		log:  app.Logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	s.log.Infof("Starting HTTP server on %s", s.addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Stopping HTTP server...")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) RegisterHandlers() {
	mux := http.NewServeMux()

	// Hello-world endpoint
	mux.HandleFunc("/", s.handleHello)
	mux.HandleFunc("/hello", s.handleHello)

	// Add logging before answer and after answer
	handler := LoggingMiddleware(s.log)(mux)

	s.httpServer = &http.Server{
		Addr:         s.addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// handleHello обработчик для hello-world
func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) {
	// Простой ответ
	response := "Hello, World! 🌍\n"

	// Добавляем немного информации для наглядности
	if r.URL.Path == "/hello" {
		response = "Hello from Calendar Service!\n"
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

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
