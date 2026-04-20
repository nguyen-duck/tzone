package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;column:id"`
	Email        string    `gorm:"uniqueIndex;not null;column:email"`
	PasswordHash *string   `gorm:"column:password_hash"`
	GoogleSub    *string   `gorm:"uniqueIndex;column:google_sub"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at"`
}

func (User) TableName() string {
	return "users"
}
