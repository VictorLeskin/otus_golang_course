package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"calendar/internal/app"
	"calendar/internal/logger"
	"calendar/internal/storage"
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

func ConvertEventResponse(event *storage.Event) EventResponse {
	return EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		UserID:      event.UserID,
	}
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
	resp := ConvertEventResponse(event)

	s.log.Infof("HTTP Create/Response: title=%q, user_id=%s", req.Title, req.UserID)

	// 4. Возвращаем JSON с созданным событием.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) EventIDFromURL(urlPath string) (string, error) {
	_ = s
	// 1. Достаём ID из URL: /events/123 ....
	pathParts := strings.Split(urlPath, "/")
	if len(pathParts) != 3 {
		return "", fmt.Errorf("invalid URL")
	}
	id := pathParts[2] //
	if id == "" {
		return "", fmt.Errorf("invalid URL")
	}
	return id, nil
}

// handleUpdateEvent — PUT /events/{id} ...
func (s *Server) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	// 1. Достаём ID из URL: /events/123 ....
	id, err := s.EventIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	// 2. Читаем JSON из тела.
	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Infof("HTTP Update Error: invalid json %s", err.Error())
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	s.log.Infof("HTTP Update/Request: title=%q, user_id=%s", req.Title, req.UserID)

	// 3. Преобразуем в storage.Event.
	event := &storage.Event{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		UserID:      req.UserID,
	}

	// 4. Вызываем ТОТ ЖЕ storage, что и gRPC!.
	err = s.storage.UpdateEvent(r.Context(), event)
	if err != nil {
		// Возвращаем JSON с ошибкой.
		s.log.Infof("HTTP Update Error: event updating failed %s", err.Error())
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 5. Преобразуем в response DTO.
	resp := ConvertEventResponse(event)

	s.log.Infof("HTTP Update/Response: title=%q, user_id=%s", req.Title, req.UserID)

	// 6. Возвращаем JSON с созданным событием.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleGetEvent — GET /events/{id} ...
func (s *Server) handleGetEvent(w http.ResponseWriter, r *http.Request) {
	// 1. Достаём ID из URL: /events/123 ....
	id, err := s.EventIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	s.log.Infof("HTTP Get/Request: id=%s", id)

	// 2. Вызываем ТОТ ЖЕ storage, что и gRPC!.
	event, err := s.storage.GetEvent(r.Context(), id)
	if err != nil {
		// Возвращаем JSON с ошибкой.
		s.log.Infof("HTTP Get Error: event getting failed %s", err.Error())
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 3. Преобразуем в response DTO.
	resp := ConvertEventResponse(event)

	s.log.Infof("HTTP Get/Response: title=%q, user_id=%s", resp.Title, resp.UserID)

	// 4. Возвращаем JSON с созданным событием.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleDeleteEvent — DELETE /events/{id} ...
func (s *Server) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	// 1. Достаём ID из URL: /events/123 ....
	id, err := s.EventIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	s.log.Infof("HTTP Delete/Request: id=%s", id)

	// 2. Вызываем ТОТ ЖЕ storage, что и gRPC!.
	err = s.storage.DeleteEvent(r.Context(), id)
	if err != nil {
		// Возвращаем JSON с ошибкой.
		s.log.Infof("HTTP Delete Error: event deleting failed %s", err.Error())
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	s.log.Infof("HTTP Delete/Response: id=%s", id)

	// 4. Возвращаем JSON с созданным событием.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"deleted id": id,
	})
}

// ListEventsParams for parsed parameters of a request.
type ListEventsParams struct {
	UserID    string
	From      time.Time
	To        time.Time
	HasFromTo bool // interval parameters are valid.
}

// ParseListEventsParams разбирает query-параметры запроса.
func ParseListEventsParameters(r *http.Request) (ListEventsParams, error) {
	query := r.URL.Query()
	params := ListEventsParams{
		UserID: query.Get("user_id"),
	}
	if params.UserID != "" {
		return params, nil
	}

	fromStr := query.Get("from")
	toStr := query.Get("to")

	// should be presented both.
	if fromStr == "" {
		if toStr == "" {
			return params, fmt.Errorf("missing input parameters")
		}
		return params, fmt.Errorf("missing 'from' parameter")
	}
	if toStr == "" {
		return params, fmt.Errorf("missing 'to' parameter")
	}

	// paring from and to.
	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return params, fmt.Errorf("invalid 'from' format: %w (use RFC3339)", err)
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return params, fmt.Errorf("invalid 'to' format: %w (use RFC3339)", err)
	}

	// chek interval validity.
	if from.After(to) {
		return params, fmt.Errorf("'from' must be before or equal to 'to'")
	}

	params.From = from
	params.To = to
	params.HasFromTo = true

	return params, nil
}

func (s *Server) handleListEventsInInterval(from, to time.Time, w http.ResponseWriter, r *http.Request) {
	s.log.Infof("HTTP ListEventsInInterval/Request: from=%s to=%s", from.String(), to.String())

	// 2. Вызываем ТОТ ЖЕ storage, что и gRPC!.
	events, err := s.storage.ListEventsInInterval(r.Context(), from, to)
	if err != nil {
		// Возвращаем JSON с ошибкой.
		s.log.Infof("HTTP ListEventsInInterval Error: list events getting failed %s", err.Error())
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 3. Преобразуем в response DTO.
	var resp ListEventsResponse
	for _, event := range events {
		resp.Events = append(resp.Events, ConvertEventResponse(event))
	}

	s.log.Infof("HTTP ListEventsInInterval/Response: from=%s to=%s", from.String(), to.String())

	// 4. Возвращаем JSON с созданным событием.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleListEvents — GET /events ...
func (s *Server) handleListEvents(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем user_id из query-параметра
	params, err := ParseListEventsParameters(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing input parameters %s", err.Error()), http.StatusBadRequest)
		return
	}
	if params.HasFromTo {
		s.handleListEventsInInterval(params.From, params.To, w, r)
		return
	}

	// 1. Получаем user_id из query-параметра
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	s.log.Infof("HTTP ListEvents/Request: id=%s", userID)

	// 2. Вызываем ТОТ ЖЕ storage, что и gRPC!.
	events, err := s.storage.ListEvents(r.Context(), userID)
	if err != nil {
		// Возвращаем JSON с ошибкой.
		s.log.Infof("HTTP ListEvents Error: list events getting failed %s", err.Error())
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 3. Преобразуем в response DTO.
	var resp ListEventsResponse
	for _, event := range events {
		resp.Events = append(resp.Events, ConvertEventResponse(event))
	}

	s.log.Infof("HTTP ListEvents/Response: user_id=%s", userID)

	// 4. Возвращаем JSON с созданным событием.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleHealth — проверка, что сервер жив.
func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
