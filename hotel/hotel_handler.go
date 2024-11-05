package hotel_service

import (
	"context"

	"github.com/tsaqiffatih/booking-system/hotel/hotelPb"
)

type HotelHandler struct {
	hotelPb.UnimplementedHotelServiceServer
	service HotelService
}

func NewHotelHandler(service HotelService) *HotelHandler {
	return &HotelHandler{service: service}
}

func (h *HotelHandler) ListHotels(ctx context.Context, req *hotelPb.ListHotelsRequest) (*hotelPb.ListHotelsResponse, error) {
	hotels, err := h.service.GetHotelsWithAvailableRooms(req.Location, req.CheckInDate, req.CheckOutDate)
	if err != nil {
		return nil, err
	}

	var hotelList []*hotelPb.Hotel
	for _, hotel := range hotels {
		rooms := make([]*hotelPb.Room, len(hotel.Rooms))
		for i, room := range hotel.Rooms {
			rooms[i] = &hotelPb.Room{
				RoomId:    room.RoomID.String(),
				Type:      string(room.Type),
				Price:     float32(room.Price),
				Available: room.Available,
			}
		}

		hotelList = append(hotelList, &hotelPb.Hotel{
			HotelId:  hotel.HotelID.String(),
			Name:     hotel.HotelsName,
			Location: hotel.Location,
			Rooms:    rooms,
		})
	}

	return &hotelPb.ListHotelsResponse{Hotels: hotelList}, nil
}
