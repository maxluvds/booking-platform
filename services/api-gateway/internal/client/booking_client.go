package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/maxluvds/booking-platform/api/proto"
)

type BookingClient struct {
	client pb.BookingServiceClient
	conn   *grpc.ClientConn
}

func NewBookingClient(addr string, timeout time.Duration) (*BookingClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(timeout),
	)
	if err != nil {
		return nil, err
	}

	return &BookingClient{
		client: pb.NewBookingServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *BookingClient) Close() error {
	return c.conn.Close()
}

func (c *BookingClient) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error) {
	return c.client.CreateBooking(ctx, req)
}

func (c *BookingClient) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingResponse, error) {
	return c.client.GetBooking(ctx, req)
}

func (c *BookingClient) ListUserBookings(ctx context.Context, req *pb.ListUserBookingsRequest) (*pb.ListBookingsResponse, error) {
	return c.client.ListUserBookings(ctx, req)
}

func (c *BookingClient) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	return c.client.CancelBooking(ctx, req)
}
