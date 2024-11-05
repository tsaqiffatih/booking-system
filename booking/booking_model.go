package booking_service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	BookingID    uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"booking_id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	HotelID      uuid.UUID `gorm:"type:uuid;not null" json:"hotel_id"`
	RoomID       uuid.UUID `gorm:"type:uuid;not null" json:"room_id"`
	CheckInDate  time.Time `gorm:"type:date;not null" json:"check_in_date"`
	CheckOutDate time.Time `gorm:"type:date;not null" json:"check_out_date"`
	Status       string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	if b.BookingID == uuid.Nil {
		b.BookingID = uuid.New()
	}
	return
}
