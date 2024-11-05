package booking_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tsaqiffatih/booking-system/booking/bookingPb"
)

type BookingHandler struct {
	service BookingService
	bookingPb.UnimplementedBookingServiceServer
}

func NewBookingHandler(service BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) CreateBooking(ctx context.Context, req *bookingPb.CreateBookingRequest) (*bookingPb.BookingResponse, error) {
	checkInDate, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		return nil, fmt.Errorf("invalid check-in date format: %v", err)
	}

	checkOutDate, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		return nil, fmt.Errorf("invalid check-out date format: %v", err)
	}

	booking := Booking{
		UserID:       uuid.MustParse(req.UserId),
		HotelID:      uuid.MustParse(req.HotelId),
		RoomID:       uuid.MustParse(req.RoomId),
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
	}

	err = h.service.CreateBooking(&booking)
	if err != nil {
		return nil, err
	}

	return &bookingPb.BookingResponse{
		BookingId:    booking.BookingID.String(),
		UserId:       booking.UserID.String(),
		HotelId:      booking.HotelID.String(),
		RoomId:       booking.RoomID.String(),
		CheckInDate:  req.CheckInDate,
		CheckOutDate: req.CheckOutDate,
		Status:       booking.Status,
	}, nil
}

func (h *BookingHandler) GetBooking(ctx context.Context, req *bookingPb.GetBookingRequest) (*bookingPb.BookingResponse, error) {
	id, err := uuid.Parse(req.BookingId)
	if err != nil {
		return nil, err
	}

	booking, err := h.service.GetBookingByID(id)
	if err != nil {
		return nil, err
	}

	return &bookingPb.BookingResponse{
		BookingId:    booking.BookingID.String(),
		UserId:       booking.UserID.String(),
		HotelId:      booking.HotelID.String(),
		RoomId:       booking.RoomID.String(),
		CheckInDate:  booking.CheckInDate.String(),
		CheckOutDate: booking.CheckOutDate.String(),
		Status:       booking.Status,
	}, nil
}

func (h *BookingHandler) GetBookings(ctx context.Context, req *bookingPb.GetBookingsRequest) (*bookingPb.GetBookingsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	bookings, err := h.service.GetBookingsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var bookingResponses []*bookingPb.BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, &bookingPb.BookingResponse{
			BookingId:    booking.BookingID.String(),
			UserId:       booking.UserID.String(),
			HotelId:      booking.HotelID.String(),
			RoomId:       booking.RoomID.String(),
			CheckInDate:  booking.CheckInDate.String(),
			CheckOutDate: booking.CheckOutDate.String(),
			Status:       booking.Status,
		})
	}

	return &bookingPb.GetBookingsResponse{Bookings: bookingResponses}, nil
}

func (h *BookingHandler) UpdateBookingStatus(ctx context.Context, req *bookingPb.UpdateBookingStatusRequest) (*bookingPb.UpdateBookingStatusResponse, error) {
	id, err := uuid.Parse(req.BookingId)
	if err != nil {
		return nil, err
	}

	err = h.service.UpdateBookingStatus(id, req.Status)
	if err != nil {
		return nil, err
	}

	return &bookingPb.UpdateBookingStatusResponse{
		BookingId: req.BookingId,
		Status:    req.Status,
	}, nil
}
