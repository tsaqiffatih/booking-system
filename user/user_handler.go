package user_service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tsaqiffatih/booking-system/user/userPb"
)

type userHandler struct {
	userPb.UnimplementedUserServiceServer
	service UserService
}

func NewUserHandler(service UserService) userPb.UserServiceServer {
	return &userHandler{service: service}
}

// CreateUser gRPC handler
func (h *userHandler) CreateUser(ctx context.Context, req *userPb.CreateUserRequest) (*userPb.CreateUserResponse, error) {
	user := &User{
		FullName:    req.GetFullName(),
		Email:       req.GetEmail(),
		Password:    req.GetPassword(),
		PhoneNumber: req.GetPhoneNumber(),
	}

	if err := h.service.RegisterUser(ctx, user); err != nil {
		return nil, err
	}

	return &userPb.CreateUserResponse{
		Id:      user.UserID.String(),
		Message: "User created successfully",
	}, nil
}

// GetUser gRPC handler
func (h *userHandler) GetUser(ctx context.Context, req *userPb.GetUserRequest) (*userPb.GetUserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &userPb.GetUserResponse{
		User: &userPb.User{
			Id:          user.UserID.String(),
			Name:        user.FullName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
	}, nil
}

// LoginUser gRPC handler
func (h *userHandler) LoginUser(ctx context.Context, req *userPb.LoginUserRequest) (*userPb.LoginUserResponse, error) {
	token, err := h.service.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userPb.LoginUserResponse{
		Token:   *token,
		Message: "Login successful",
	}, nil
}
