package internalhttp

import (
	"calendar/internal/app"
	"calendar/internal/logger"
	"calendar/internal/storage"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Server struct {
	addr       string
	httpServer *http.Server
	log        *logger.Logger
	storage    storage.Storage
}

func NewServer(cfg Config, app *app.App, storage storage.Storage) *Server {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return &Server{
		addr:    addr,
		log:     app.Logger,
		storage: storage,
	}
}

func (s *Server) Start(ctx context.Context) error {
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

	// Hello-world endpoint.
	mux.HandleFunc("/", s.handleHello)
	mux.HandleFunc("/hello", s.handleHello)

	// new API handlers.
	mux.HandleFunc("POST /events", s.handleCreateEvent)
	mux.HandleFunc("GET /events/{id}", s.handleGetEvent)
	mux.HandleFunc("PUT /events/{id}", s.handleUpdateEvent)
	mux.HandleFunc("DELETE /events/{id}", s.handleDeleteEvent)
	mux.HandleFunc("GET /events", s.handleListEvents)

	// Health check для мониторинга.
	mux.HandleFunc("GET /health", s.handleHealth)

	// Add logging before answer and after answer.
	handler := LoggingMiddleware(s.log)(mux)

	s.httpServer = &http.Server{
		Addr:         s.addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (s *Server) GetHandler() http.Handler {
	return s.httpServer.Handler
}

// handleHello обработчик для hello-world.
func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) {
	// Простой ответ.
	response := "Hello, World!\n"

	// to distinguish betwee /hello and just /.
	if r.URL.Path == "/hello" {
		response = "Hello from Calendar Service!\n"
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// handleCreateEvent — POST /events ...
func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	// 1. Читаем JSON из тела.
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Infof("HTTP Create Error: invalid json %s", err.Error())
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	s.log.Infof("HTTP Create/Request: title=%q, user_id=%s", req.Title, req.UserID)

	// 2. Преобразуем в storage.Event.
	event := &storage.Event{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		UserID:      req.UserID,
	}

	// 3. Вызываем ТОТ ЖЕ storage, что и gRPC!.
	err := s.storage.CreateEvent(r.Context(), event)
	if err != nil {
		// Возвращаем JSON с ошибкой.
		s.log.Infof("HTTP Create Error: event creating failed %s", err.Error())
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 4. Преобразуем в response DTO.
	resp := EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		UserID:      event.UserID,
	}

	s.log.Infof("HTTP Create/Response: title=%q, user_id=%s", req.Title, req.UserID)

	// 4. Возвращаем JSON с созданным событием.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// handleUpdateEvent — PUT /events/{id} ...
func (s *Server) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	// 1. Читаем JSON из тела.
	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Infof("HTTP Update Error: invalid json %s", err.Error())
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	s.log.Infof("HTTP Update/Request: title=%q, user_id=%s", req.Title, req.UserID)

	// 2. Преобразуем в storage.Event.
	event := &storage.Event{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		UserID:      req.UserID,
	}

	// 3. Вызываем ТОТ ЖЕ storage, что и gRPC!.
	err := s.storage.UpdateEvent(r.Context(), event)
	if err != nil {
		// Возвращаем JSON с ошибкой.
		s.log.Infof("HTTP Update Error: event creating failed %s", err.Error())
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 4. Преобразуем в response DTO.
	resp := EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		UserID:      event.UserID,
	}

	s.log.Infof("HTTP Update/Response: title=%q, user_id=%s", req.Title, req.UserID)

	// 4. Возвращаем JSON с созданным событием.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleGetEvent — GET /events/{id} ...
func (s *Server) handleGetEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: имплементация
}

// handleDeleteEvent — DELETE /events/{id} ...
func (s *Server) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: имплементация
}

// handleListEvents — GET /events ...
func (s *Server) handleListEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: имплементация
}

// handleHealth — проверка, что сервер жив.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
