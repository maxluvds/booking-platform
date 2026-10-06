package delivery

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxluvds/booking-platform/api/proto"
	"github.com/maxluvds/booking-platform/services/event-service/internal/domain"
	"github.com/maxluvds/booking-platform/services/event-service/internal/usecase"
)

type EventGRPCHandler struct {
	pb.UnimplementedEventServiceServer
	usecase *usecase.EventUseCase
}

func NewEventGRPCHandler(usecase *usecase.EventUseCase) *EventGRPCHandler {
	return &EventGRPCHandler{
		usecase: usecase,
	}
}

func (h *EventGRPCHandler) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.EventResponse, error) {
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid date format: %v", err)
	}

	createReq := &domain.CreateEventRequest{
		Title:       req.Title,
		Description: req.Description,
		City:        req.City,
		Address:     req.Address,
		Category:    req.Category,
		Date:        date,
		TotalSeats:  int(req.TotalSeats),
		Price:       req.Price,
	}

	event, err := h.usecase.Create(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	return &pb.EventResponse{
		Event: convertToProto(event),
	}, nil
}

func (h *EventGRPCHandler) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.EventResponse, error) {
	event, err := h.usecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "event not found: %v", err)
	}

	return &pb.EventResponse{
		Event: convertToProto(event),
	}, nil
}

func (h *EventGRPCHandler) ListEvents(ctx context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	listReq := &domain.ListEventsRequest{
		City:     req.City,
		Category: req.Category,
		Status:   req.Status,
		Limit:    int(req.Limit),
		Offset:   int(req.Offset),
	}

	events, total, err := h.usecase.List(ctx, listReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events: %v", err)
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, e := range events {
		protoEvents[i] = convertToProto(e)
	}

	return &pb.ListEventsResponse{
		Events: protoEvents,
		Total:  int32(total),
	}, nil
}

func (h *EventGRPCHandler) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	updateReq := &domain.UpdateEventRequest{}

	if req.Title != nil {
		updateReq.Title = req.Title
	}
	if req.Description != nil {
		updateReq.Description = req.Description
	}
	if req.City != nil {
		updateReq.City = req.City
	}
	if req.Address != nil {
		updateReq.Address = req.Address
	}
	if req.Category != nil {
		updateReq.Category = req.Category
	}
	if req.Date != nil {
		date, err := time.Parse(time.RFC3339, *req.Date)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid date format: %v", err)
		}
		updateReq.Date = &date
	}
	if req.TotalSeats != nil {
		seats := int(*req.TotalSeats)
		updateReq.TotalSeats = &seats
	}
	if req.Price != nil {
		updateReq.Price = req.Price
	}
	if req.Status != nil {
		updateReq.Status = req.Status
	}

	event, err := h.usecase.Update(ctx, req.Id, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}

	return &pb.EventResponse{
		Event: convertToProto(event),
	}, nil
}

func (h *EventGRPCHandler) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	err := h.usecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}

	return &pb.DeleteEventResponse{
		Success: true,
	}, nil
}

func (h *EventGRPCHandler) GetAvailableSeats(ctx context.Context, req *pb.GetAvailableSeatsRequest) (*pb.GetAvailableSeatsResponse, error) {
	event, err := h.usecase.GetByID(ctx, req.EventId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "event not found: %v", err)
	}

	return &pb.GetAvailableSeatsResponse{
		EventId:        event.ID,
		AvailableSeats: int32(event.AvailableSeats),
	}, nil
}

func convertToProto(e *domain.Event) *pb.Event {
	return &pb.Event{
		Id:             e.ID,
		Title:          e.Title,
		Description:    e.Description,
		City:           e.City,
		Address:        e.Address,
		Category:       e.Category,
		Date:           e.Date.Format(time.RFC3339),
		TotalSeats:     int32(e.TotalSeats),
		AvailableSeats: int32(e.AvailableSeats),
		Price:          e.Price,
		Status:         string(e.Status),
		CreatedAt:      e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      e.UpdatedAt.Format(time.RFC3339),
	}
}
