package hotel_service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomType string

const (
	Standard          RoomType = "Standard"
	Superior          RoomType = "Superior"
	Deluxe            RoomType = "Deluxe"
	Executive         RoomType = "Executive"
	Suite             RoomType = "Suite"
	JuniorSuite       RoomType = "Junior Suite"
	FamilyRoom        RoomType = "Family Room"
	ConnectingRoom    RoomType = "Connecting Room"
	Cabana            RoomType = "Cabana"
	Penthouse         RoomType = "Penthouse"
	Luxury            RoomType = "Luxury"
	PresidentialSuite RoomType = "Presidential Suite"
	Accessible        RoomType = "Accessible"
)

type Hotel struct {
	HotelID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"hotel_id"`
	HotelsName string    `gorm:"type:varchar(100);not null" json:"hotel_name"`
	Location   string    `gorm:"type:varchar(100);not null" json:"location"`
	Rooms      []Room    `gorm:"foreignKey:HotelID;references:HotelID"`
}

func (h *Hotel) BeforeCreate(tx *gorm.DB) (err error) {
	h.HotelID = uuid.New()
	return
}

type Room struct {
	RoomID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"room_id"`
	HotelID   uuid.UUID `gorm:"type:uuid;not null" json:"hotel_id"`
	Type      RoomType  `gorm:"type:varchar(50);not null" json:"type"`
	Price     float64   `gorm:"type:decimal;not null" json:"price"`
	Available bool      `gorm:"not null" json:"available"`

	Hotel Hotel `gorm:"foreignKey:HotelID;references:HotelID"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) (err error) {
	r.RoomID = uuid.New()
	return
}
