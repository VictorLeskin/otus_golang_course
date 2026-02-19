package grpc

import (
	"context"
	"log/slog"

	"calendar/api/pb/calendar"  // geerated code
	"calendar/internal/storage" // storage.Storage interface

	"google.golang.org/protobuf/types/known/timestamppb"
)

// calendar.CalendarServiceServer interface implemented for gRPC Server.
type Server struct {
	calendar.UnimplementedCalendarServiceServer
	storage storage.Storage
	logger  *slog.Logger
}

func New(storage storage.Storage, logger *slog.Logger) *Server {
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

func (s *Server) CreateEvent(ctx context.Context,
	req *calendar.CreateEventRequest) (*calendar.CreateEventResponse, error) {
	event := convertFromPBEvent(req.Event)

	err := s.storage.CreateEvent(ctx, event)
	if err != nil {
		return &calendar.CreateEventResponse{
			Event:        nil,
			ErrorMessage: err.Error(),
		}, err
	}

	return &calendar.CreateEventResponse{
		Event: convertToPBEvent(event),
	}, nil
}
