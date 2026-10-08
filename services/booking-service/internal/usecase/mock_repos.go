package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/maxluvds/booking-platform/services/booking-service/internal/domain"
)

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) Create(
	ctx context.Context,
	booking *domain.Booking,
	outboxEvent *domain.OutboxEvent,
) (*domain.Booking, error) {
	args := m.Called(ctx, booking, outboxEvent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) ListByUser(ctx context.Context, userID int64) ([]*domain.Booking, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) UpdateStatus(ctx context.Context, id int64, status domain.BookingStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

type MockOutboxRepository struct {
	mock.Mock
}

func (m *MockOutboxRepository) GetPending(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.OutboxEvent), args.Error(1)
}

func (m *MockOutboxRepository) MarkAsSent(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOutboxRepository) MarkAsFailed(ctx context.Context, id int64, errMsg string) error {
	args := m.Called(ctx, id, errMsg)
	return args.Error(0)
}
