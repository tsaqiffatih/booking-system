package hotel_service

type HotelService interface {
	GetHotelsWithAvailableRooms(location, checkInDate, checkOutDate string) ([]Hotel, error)
	GetHotels(location string) ([]Hotel, error)
	CheckRoomAvailability(hotelID, roomID, checkInDate, checkOutDate string) (bool, error)
}

type hotelService struct {
	repo HotelRepository
}

func NewHotelService(repo HotelRepository) HotelService {
	return &hotelService{repo: repo}
}

func (s *hotelService) GetHotelsWithAvailableRooms(location, checkInDate, checkOutDate string) ([]Hotel, error) {
	hotels, err := s.repo.GetHotels(location)
	if err != nil {
		return nil, err
	}

	for i, hotel := range hotels {
		availableRooms, err := s.repo.GetAvailableRooms(hotel.HotelID)
		if err != nil {
			return nil, err
		}
		hotels[i].Rooms = availableRooms
	}

	return hotels, nil
}

func (s *hotelService) GetHotels(location string) ([]Hotel, error) {
	return s.repo.GetHotels(location)
}

func (s *hotelService) CheckRoomAvailability(hotelID, roomID, checkInDate, checkOutDate string) (bool, error) {
	return s.repo.CheckRoomAvailability(hotelID, roomID, checkInDate, checkOutDate)
}
