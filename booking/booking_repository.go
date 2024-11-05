package booking_service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *Booking) error
	GetBookingByID(id uuid.UUID) (*Booking, error)
	GetBookingsByUserID(userID uuid.UUID) ([]Booking, error)
	UpdateBookingStatus(id uuid.UUID, status string) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

func (r *bookingRepository) CreateBooking(booking *Booking) error {
	return r.db.Create(booking).Error
}
func (r *bookingRepository) GetBookingByID(id uuid.UUID) (*Booking, error) {
	var booking Booking
	err := r.db.First(&booking, "booking_id = ?", id).Error
	return &booking, err
}
func (r *bookingRepository) GetBookingsByUserID(userID uuid.UUID) ([]Booking, error) {
	var bookings []Booking
	err := r.db.Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}
func (r *bookingRepository) UpdateBookingStatus(id uuid.UUID, status string) error {
	return r.db.Model(Booking{}).Where("booking_id = ?", id).Update("status", status).Error
}
