package booking_service

import (
	"context"
	"fmt"

	"github.com/tsaqiffatih/booking-system/booking/bookingPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type BookingClient struct {
	conn    *grpc.ClientConn
	service bookingPb.BookingServiceClient
}

func NewBookingClient(url string) (*BookingClient, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	// conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := bookingPb.NewBookingServiceClient(conn)
	return &BookingClient{conn: conn, service: client}, nil
}

func (c *BookingClient) Close() {
	c.conn.Close()
}

func (c *BookingClient) CreateBooking(ctx context.Context, userID, hotelID, roomID, checkInDate, checkOutDate string) (*bookingPb.BookingResponse, error) {
	req := &bookingPb.CreateBookingRequest{
		UserId:       userID,
		HotelId:      hotelID,
		RoomId:       roomID,
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}
	resp, err := c.service.CreateBooking(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}
	return resp, nil
}

func (c *BookingClient) GetBooking(ctx context.Context, bookingID string) (*bookingPb.BookingResponse, error) {
	req := &bookingPb.GetBookingRequest{
		BookingId: bookingID,
	}
	resp, err := c.service.GetBooking(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	return resp, nil
}

func (c *BookingClient) UpdateBookingStatus(ctx context.Context, bookingID, status string) (*bookingPb.UpdateBookingStatusResponse, error) {
	req := &bookingPb.UpdateBookingStatusRequest{
		BookingId: bookingID,
		Status:    status,
	}
	resp, err := c.service.UpdateBookingStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update booking status: %w", err)
	}
	return resp, nil
}
