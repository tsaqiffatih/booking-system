package booking_service

import "github.com/google/uuid"

type BookingService interface {
	CreateBooking(booking *Booking) error
	GetBookingByID(id uuid.UUID) (*Booking, error)
	GetBookingsByUserID(userID uuid.UUID) ([]Booking, error)
	UpdateBookingStatus(id uuid.UUID, status string) error
}

type bookingService struct {
	repo BookingRepository
}

func NewBookingService(repo BookingRepository) BookingService {
	return &bookingService{repo}
}

func (s *bookingService) CreateBooking(booking *Booking) error {
	booking.Status = "pending"
	return s.repo.CreateBooking(booking)
}

func (s *bookingService) GetBookingByID(id uuid.UUID) (*Booking, error) {
	return s.repo.GetBookingByID(id)
}

func (s *bookingService) GetBookingsByUserID(userID uuid.UUID) ([]Booking, error) {
	return s.repo.GetBookingsByUserID(userID)
}

func (s *bookingService) UpdateBookingStatus(id uuid.UUID, status string) error {
	return s.repo.UpdateBookingStatus(id, status)
}
