package grpc

import (
	"context"
	"time"

	"calendar/api/pb/calendar" // geerated code
	"calendar/internal/logger"
	"calendar/internal/storage" // storage.Storage interface

	"google.golang.org/protobuf/types/known/timestamppb"
)

// calendar.CalendarServiceServer interface implemented for gRPC Server.
type Server struct {
	calendar.UnimplementedCalendarServiceServer
	storage storage.Storage
	logger  *logger.Logger
}

func New(storage storage.Storage, logger *logger.Logger) *Server {
	return &Server{
		storage: storage,
		logger:  logger,
	}
}

func convertFromPBEvent(from *calendar.Event) *storage.Event {
	return &storage.Event{
		ID:          from.Id,
		Title:       from.Title,
		Description: from.Description,
		StartTime:   from.StartTime.AsTime(),
		EndTime:     from.EndTime.AsTime(),
		UserID:      from.UserId,
	}
}

func convertToPBEvent(from *storage.Event) *calendar.Event {
	return &calendar.Event{
		Id:          from.ID,
		Title:       from.Title,
		Description: from.Description,
		StartTime:   timestamppb.New(from.StartTime),
		EndTime:     timestamppb.New(from.EndTime),
		UserId:      from.UserID,
	}
}

func (s *Server) LogCalendarEvent(method string, stage string, event *calendar.Event) {
	s.logger.Infof("gRPC %s/%s Id: %s Title: %q Description: %q StartTime:%s EndTime:%s UserId: %s",
		method, stage,
		event.Id, event.Title, event.Description,
		event.StartTime.AsTime().Format(time.RFC3339),
		event.EndTime.AsTime().Format(time.RFC3339),
		event.UserId)
}

func (s *Server) LogError(method string, err error) {
	s.logger.Infof("gRPC %s Error: %s", method, err.Error())
}

func (s *Server) LogDeleteRequest(req *calendar.DeleteEventRequest) {
	s.logger.Infof("gRPC Delete/Request Id: %s", req.Id)
}

func (s *Server) LogDeleteResponse(req *calendar.DeleteEventRequest) {
	s.logger.Infof("gRPC Delete/Response Id: %s", req.Id)
}

func (s *Server) CreateEvent(ctx context.Context,
	req *calendar.CreateEventRequest) (*calendar.CreateEventResponse, error) {
	s.LogCalendarEvent("Create", "Request", req.Event)

	event := convertFromPBEvent(req.Event)

	err := s.storage.CreateEvent(ctx, event)
	if err != nil {
		s.LogError("Create", err)
		return &calendar.CreateEventResponse{
			Event:        nil,
			ErrorMessage: err.Error(),
		}, err
	}

	resp := &calendar.CreateEventResponse{
		Event: convertToPBEvent(event),
	}

	s.LogCalendarEvent("Create", "Response", resp.Event)
	return resp, nil
}

func (s *Server) UpdateEvent(ctx context.Context,
	req *calendar.UpdateEventRequest) (*calendar.UpdateEventResponse, error) {
	s.LogCalendarEvent("Update", "Request", req.Event)

	event := convertFromPBEvent(req.Event)

	err := s.storage.UpdateEvent(ctx, event)
	if err != nil {
		s.LogError("Update", err)
		return &calendar.UpdateEventResponse{
			Event:        nil,
			ErrorMessage: err.Error(),
		}, err
	}

	resp := &calendar.UpdateEventResponse{
		Event: convertToPBEvent(event),
	}

	s.LogCalendarEvent("Update", "Response", resp.Event)
	return resp, nil
}

func (s *Server) DeleteEvent(ctx context.Context,
	req *calendar.DeleteEventRequest) (*calendar.DeleteEventResponse, error) {
	s.LogDeleteRequest(req)

	err := s.storage.DeleteEvent(ctx, req.Id)
	if err != nil {
		s.LogError("Delete", err)
		return &calendar.DeleteEventResponse{
			ErrorMessage: err.Error(),
		}, err
	}

	resp := &calendar.DeleteEventResponse{}

	s.LogDeleteResponse(req)
	return resp, nil
}
