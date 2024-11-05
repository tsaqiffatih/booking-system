package hotel_service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HotelRepository interface {
	GetHotels(location string) ([]Hotel, error)
	GetAvailableRooms(hotelID uuid.UUID) ([]Room, error)
	CreateHotel(hotel *Hotel) error
	UpdateHotel(hotel *Hotel) error
	DeleteHotel(hotelID string) error
	CheckRoomAvailability(hotelID, roomID, checkInDate, checkOutDate string) (bool, error)
}

type hotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{db: db}
}

// GetHotels implements HotelRepository.
func (r *hotelRepository) GetHotels(location string) ([]Hotel, error) {
	var hotels []Hotel
	err := r.db.Where("location = ?", location).Find(&hotels).Error
	return hotels, err
}

func (r *hotelRepository) GetAvailableRooms(hotelID uuid.UUID) ([]Room, error) {
	var availableRooms []Room
	err := r.db.Where("hotel_id = ? AND available = ? ", hotelID, true).Find(&availableRooms).Error
	return availableRooms, err
}

// CreateHotel implements HotelRepository.
func (r *hotelRepository) CreateHotel(hotel *Hotel) error {
	return r.db.Create(hotel).Error
}

// UpdateHotel implements HotelRepository.
func (r *hotelRepository) UpdateHotel(hotel *Hotel) error {
	return r.db.Save(hotel).Error
}

// DeleteHotel implements HotelRepository.
func (r *hotelRepository) DeleteHotel(hotelID string) error {
	return r.db.Delete(&Hotel{}, "hotel_id = ?", hotelID).Error
}

// CheckRoomAvailability implements HotelRepository.
func (r *hotelRepository) CheckRoomAvailability(hotelID string, roomID string, checkInDate string, checkOutDate string) (bool, error) {
	return true, nil
}
