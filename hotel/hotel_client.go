package hotel_service

import (
	"context"
	"fmt"

	"github.com/tsaqiffatih/booking-system/hotel/hotelPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type HotelClient struct {
	conn    *grpc.ClientConn
	service hotelPb.HotelServiceClient
}

func NewHotelClient(url string) (*HotelClient, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}

	client := hotelPb.NewHotelServiceClient(conn)
	return &HotelClient{conn: conn, service: client}, nil
}

func (c *HotelClient) ListHotelsClient(ctx context.Context, location, checkInDate, checkOutDate string) (*hotelPb.ListHotelsResponse, error) {
	req := &hotelPb.ListHotelsRequest{
		Location:     location,
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}

	resp, err := c.service.ListHotels(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HotelClient) CheckRoomAvailability(ctx context.Context, hotelID, roomID, checkInDate, checkOutDate string) (*hotelPb.CheckRoomResponse, error) {
	req := &hotelPb.CheckRoomRequest{
		HotelId:      hotelID,
		RoomId:       roomID,
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}

	resp, err := c.service.CheckRoomAvailability(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check room availability: %w", err)
	}
	return resp, nil
}
