package models

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
	RestaurantName string    `gorm:"column:restaurant_name;type:text;not null"`
	Cuisine        *string   `gorm:"type:text"`
}
