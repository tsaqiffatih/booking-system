package user_service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	FullName    string    `gorm:"type:varchar(50);not null" json:"full_name"`
	Password    string    `gorm:"type:varchar(50);not null" json:"password"`
	Email       string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	PhoneNumber string    `gorm:"type:varchar(100);not null" json:"phone_number"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UserID = uuid.New()
	return
}
