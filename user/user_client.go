package user_service

import (
	"context"
	"fmt"

	"github.com/tsaqiffatih/booking-system/user/userPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type UserClient struct {
	conn    *grpc.ClientConn
	service userPb.UserServiceClient
}

func NewUserClient(url string) (*UserClient, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	// conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := userPb.NewUserServiceClient(conn)
	return &UserClient{conn: conn, service: client}, nil
}

func (c *UserClient) Close() {
	c.conn.Close()
}

func (c *UserClient) RegisterUser(ctx context.Context, fullName, password, email, phoneNumber string) (*userPb.CreateUserResponse, error) {
	req := &userPb.CreateUserRequest{
		FullName:    fullName,
		Email:       email,
		PhoneNumber: phoneNumber,
		Password:    password,
	}
	resp, err := c.service.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	return &userPb.CreateUserResponse{
		Id:      resp.Id,
		Message: resp.Message,
	}, nil
}

func (c *UserClient) GetUser(ctx context.Context, userID string) (*userPb.User, error) {
	req := &userPb.GetUserRequest{
		Id: userID,
	}
	resp, err := c.service.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &userPb.User{
		Id:          resp.User.Id,
		Name:        resp.User.Name,
		Email:       resp.User.Email,
		PhoneNumber: resp.User.PhoneNumber,
	}, nil
}

func (c *UserClient) LoginUser(ctx context.Context, email, password string) (*userPb.LoginUserResponse, error) {
	req := &userPb.LoginUserRequest{
		Email:    email,
		Password: password,
	}
	resp, err := c.service.LoginUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &userPb.LoginUserResponse{
		Token:   resp.Token,
		Message: resp.Message,
	}, nil
}
