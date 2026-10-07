package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/maxluvds/booking-platform/api/proto"
)

type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(addr string, timeout time.Duration) (*UserClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(timeout),
	)
	if err != nil {
		return nil, err
	}

	return &UserClient{
		client: pb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return c.client.Register(ctx, req)
}

func (c *UserClient) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	return c.client.Login(ctx, req)
}

func (c *UserClient) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	return c.client.GetUser(ctx, req)
}

func (c *UserClient) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	return c.client.ValidateToken(ctx, req)
}
