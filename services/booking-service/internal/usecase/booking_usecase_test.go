package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/maxluvds/booking-platform/services/booking-service/internal/domain"
)

func TestBookingUseCase_CreateBooking_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	req := &domain.CreateBookingRequest{
		UserID:  1,
		EventID: 10,
		Seats:   2,
	}

	expectedBooking := &domain.Booking{
		ID:         1,
		UserID:     1,
		EventID:    10,
		Seats:      2,
		TotalPrice: 200.0,
		Status:     domain.StatusConfirmed,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockBookingRepo.On("Create",
		mock.Anything,
		mock.AnythingOfType("*domain.Booking"),
		mock.AnythingOfType("*domain.OutboxEvent"),
	).Return(expectedBooking, nil)

	result, err := uc.CreateBooking(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, 2, result.Seats)
	assert.Equal(t, 200.0, result.TotalPrice)

	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_CreateBooking_InvalidUserID(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	req := &domain.CreateBookingRequest{
		UserID:  0,
		EventID: 10,
		Seats:   2,
	}

	result, err := uc.CreateBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user id")
}

func TestBookingUseCase_CreateBooking_InvalidEventID(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	req := &domain.CreateBookingRequest{
		UserID:  1,
		EventID: 0,
		Seats:   2,
	}

	result, err := uc.CreateBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid event id")
}

func TestBookingUseCase_CreateBooking_TooManySeats(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	req := &domain.CreateBookingRequest{
		UserID:  1,
		EventID: 10,
		Seats:   100,
	}

	result, err := uc.CreateBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "seats must be between 1 and 10")
}

func TestBookingUseCase_CreateBooking_ZeroSeats(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	req := &domain.CreateBookingRequest{
		UserID:  1,
		EventID: 10,
		Seats:   0,
	}

	result, err := uc.CreateBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "seats must be between 1 and 10")
}

func TestBookingUseCase_CreateBooking_RepoError(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	req := &domain.CreateBookingRequest{
		UserID:  1,
		EventID: 10,
		Seats:   2,
	}

	mockBookingRepo.On("Create",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, errors.New("database error"))

	result, err := uc.CreateBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_GetBooking_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	expected := &domain.Booking{
		ID:     1,
		UserID: 1,
		Seats:  2,
		Status: domain.StatusConfirmed,
	}

	mockBookingRepo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

	result, err := uc.GetBooking(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)

	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_GetBooking_NotFound(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	mockBookingRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)

	result, err := uc.GetBooking(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "booking not found")

	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_GetBooking_InvalidID(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	result, err := uc.GetBooking(context.Background(), 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid booking id")
}

func TestBookingUseCase_ListUserBookings_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	bookings := []*domain.Booking{
		{ID: 1, UserID: 1, Seats: 2},
		{ID: 2, UserID: 1, Seats: 1},
	}

	mockBookingRepo.On("ListByUser", mock.Anything, int64(1)).Return(bookings, nil)

	result, err := uc.ListUserBookings(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, result, 2)

	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_ListUserBookings_InvalidUserID(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	result, err := uc.ListUserBookings(context.Background(), 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user id")
}

func TestBookingUseCase_CancelBooking_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	existing := &domain.Booking{
		ID:     1,
		UserID: 1,
		Status: domain.StatusConfirmed,
	}

	mockBookingRepo.On("GetByID", mock.Anything, int64(1)).Return(existing, nil)
	mockBookingRepo.On("UpdateStatus", mock.Anything, int64(1), domain.StatusCancelled).Return(nil)

	req := &domain.CancelBookingRequest{
		BookingID: 1,
		UserID:    1,
	}

	err := uc.CancelBooking(context.Background(), req)

	require.NoError(t, err)

	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_CancelBooking_NotOwner(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	existing := &domain.Booking{
		ID:     1,
		UserID: 1,
		Status: domain.StatusConfirmed,
	}

	mockBookingRepo.On("GetByID", mock.Anything, int64(1)).Return(existing, nil)

	req := &domain.CancelBookingRequest{
		BookingID: 1,
		UserID:    2, // другой пользователь
	}

	err := uc.CancelBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "you can only cancel your own bookings")

	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_CancelBooking_AlreadyCancelled(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	existing := &domain.Booking{
		ID:     1,
		UserID: 1,
		Status: domain.StatusCancelled,
	}

	mockBookingRepo.On("GetByID", mock.Anything, int64(1)).Return(existing, nil)

	req := &domain.CancelBookingRequest{
		BookingID: 1,
		UserID:    1,
	}

	err := uc.CancelBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "booking is already cancelled")

	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_CancelBooking_NotFound(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	mockBookingRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)

	req := &domain.CancelBookingRequest{
		BookingID: 999,
		UserID:    1,
	}

	err := uc.CancelBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "booking not found")

	mockBookingRepo.AssertExpectations(t)
}

func TestBookingUseCase_CancelBooking_InvalidID(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockOutboxRepo := new(MockOutboxRepository)
	uc := NewBookingUseCase(mockBookingRepo, mockOutboxRepo)

	req := &domain.CancelBookingRequest{
		BookingID: 0,
		UserID:    1,
	}

	err := uc.CancelBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid booking id")
}
