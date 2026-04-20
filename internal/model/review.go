package model

import (
	"time"

	"github.com/google/uuid"
)

// Review is a placeholder model for future review feature migrations.
type Review struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;column:id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	DeviceID  string    `gorm:"type:varchar(64);not null;column:device_id"`
	Rating    int       `gorm:"not null;column:rating"`
	Comment   string    `gorm:"type:text;not null;column:comment"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (Review) TableName() string {
	return "reviews"
}
