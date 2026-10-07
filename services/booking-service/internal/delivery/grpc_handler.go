package delivery

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxluvds/booking-platform/api/proto"
	"github.com/maxluvds/booking-platform/services/booking-service/internal/domain"
	"github.com/maxluvds/booking-platform/services/booking-service/internal/usecase"
)

type BookingGRPCHandler struct {
	pb.UnimplementedBookingServiceServer
	usecase *usecase.BookingUseCase
}

func NewBookingGRPCHandler(usecase *usecase.BookingUseCase) *BookingGRPCHandler {
	return &BookingGRPCHandler{
		usecase: usecase,
	}
}

func (h *BookingGRPCHandler) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error) {
	createReq := &domain.CreateBookingRequest{
		UserID:  req.UserId,
		EventID: req.EventId,
		Seats:   int(req.Seats),
	}

	booking, err := h.usecase.CreateBooking(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create booking: %v", err)
	}

	return &pb.BookingResponse{
		Booking: convertToProto(booking),
	}, nil
}

func (h *BookingGRPCHandler) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingResponse, error) {
	booking, err := h.usecase.GetBooking(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "booking not found: %v", err)
	}

	return &pb.BookingResponse{
		Booking: convertToProto(booking),
	}, nil
}

func (h *BookingGRPCHandler) ListUserBookings(ctx context.Context, req *pb.ListUserBookingsRequest) (*pb.ListBookingsResponse, error) {
	bookings, err := h.usecase.ListUserBookings(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list bookings: %v", err)
	}

	protoBookings := make([]*pb.Booking, len(bookings))
	for i, b := range bookings {
		protoBookings[i] = convertToProto(b)
	}

	return &pb.ListBookingsResponse{
		Bookings: protoBookings,
		Total:    int32(len(bookings)),
	}, nil
}

func (h *BookingGRPCHandler) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	cancelReq := &domain.CancelBookingRequest{
		BookingID: req.BookingId,
		UserID:    req.UserId,
	}

	err := h.usecase.CancelBooking(ctx, cancelReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cancel booking: %v", err)
	}

	return &pb.CancelBookingResponse{
		Success: true,
	}, nil
}

func convertToProto(b *domain.Booking) *pb.Booking {
	return &pb.Booking{
		Id:         b.ID,
		UserId:     b.UserID,
		EventId:    b.EventID,
		Seats:      int32(b.Seats),
		TotalPrice: b.TotalPrice,
		Status:     string(b.Status),
		CreatedAt:  b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  b.UpdatedAt.Format(time.RFC3339),
	}
}
