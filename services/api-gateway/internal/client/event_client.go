package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/maxluvds/booking-platform/api/proto"
)

type EventClient struct {
	client pb.EventServiceClient
	conn   *grpc.ClientConn
}

func NewEventClient(addr string, timeout time.Duration) (*EventClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(timeout),
	)
	if err != nil {
		return nil, err
	}

	return &EventClient{
		client: pb.NewEventServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *EventClient) Close() error {
	return c.conn.Close()
}

func (c *EventClient) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.EventResponse, error) {
	return c.client.CreateEvent(ctx, req)
}

func (c *EventClient) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.EventResponse, error) {
	return c.client.GetEvent(ctx, req)
}

func (c *EventClient) ListEvents(ctx context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	return c.client.ListEvents(ctx, req)
}

func (c *EventClient) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	return c.client.UpdateEvent(ctx, req)
}

func (c *EventClient) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	return c.client.DeleteEvent(ctx, req)
}

func (c *EventClient) GetAvailableSeats(ctx context.Context, req *pb.GetAvailableSeatsRequest) (*pb.GetAvailableSeatsResponse, error) {
	return c.client.GetAvailableSeats(ctx, req)
}
