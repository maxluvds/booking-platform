package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/maxluvds/booking-platform/services/booking-service/internal/domain"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking, outboxEvent *domain.OutboxEvent) (*domain.Booking, error)
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	ListByUser(ctx context.Context, userID int64) ([]*domain.Booking, error)
	UpdateStatus(ctx context.Context, id int64, status domain.BookingStatus) error
}

type OutboxRepository interface {
	GetPending(ctx context.Context, limit int) ([]*domain.OutboxEvent, error)
	MarkAsSent(ctx context.Context, id int64) error
	MarkAsFailed(ctx context.Context, id int64, errMsg string) error
}

type BookingUseCase struct {
	bookingRepo BookingRepository
	outboxRepo  OutboxRepository
}

func NewBookingUseCase(
	bookingRepo BookingRepository,
	outboxRepo OutboxRepository,
) *BookingUseCase {
	return &BookingUseCase{
		bookingRepo: bookingRepo,
		outboxRepo:  outboxRepo,
	}
}

func (uc *BookingUseCase) CreateBooking(ctx context.Context, req *domain.CreateBookingRequest) (*domain.Booking, error) {
	if req.UserID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	if req.EventID <= 0 {
		return nil, fmt.Errorf("invalid event id")
	}
	if req.Seats <= 0 || req.Seats > 10 {
		return nil, fmt.Errorf("seats must be between 1 and 10")
	}

	pricePerSeat := 100.0
	totalPrice := pricePerSeat * float64(req.Seats)

	booking := &domain.Booking{
		UserID:     req.UserID,
		EventID:    req.EventID,
		Seats:      req.Seats,
		TotalPrice: totalPrice,
		Status:     domain.StatusConfirmed,
	}

	payload, err := json.Marshal(map[string]interface{}{
		"user_id":     req.UserID,
		"event_id":    req.EventID,
		"seats":       req.Seats,
		"total_price": totalPrice,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	outboxEvent := &domain.OutboxEvent{
		EventType: "booking.created",
		Payload:   payload,
		Status:    "pending",
	}

	createdBooking, err := uc.bookingRepo.Create(ctx, booking, outboxEvent)
	if err != nil {
		return nil, err
	}

	return createdBooking, nil
}

func (uc *BookingUseCase) GetBooking(ctx context.Context, id int64) (*domain.Booking, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid booking id")
	}

	booking, err := uc.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}

	return booking, nil
}

func (uc *BookingUseCase) ListUserBookings(ctx context.Context, userID int64) ([]*domain.Booking, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	return uc.bookingRepo.ListByUser(ctx, userID)
}

func (uc *BookingUseCase) CancelBooking(ctx context.Context, req *domain.CancelBookingRequest) error {
	if req.BookingID <= 0 {
		return fmt.Errorf("invalid booking id")
	}

	booking, err := uc.bookingRepo.GetByID(ctx, req.BookingID)
	if err != nil {
		return err
	}
	if booking == nil {
		return fmt.Errorf("booking not found")
	}

	if booking.UserID != req.UserID {
		return fmt.Errorf("you can only cancel your own bookings")
	}

	if booking.Status == domain.StatusCancelled {
		return fmt.Errorf("booking is already cancelled")
	}

	return uc.bookingRepo.UpdateStatus(ctx, req.BookingID, domain.StatusCancelled)
}
