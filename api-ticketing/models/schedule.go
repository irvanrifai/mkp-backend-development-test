package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MovieID        uint           `gorm:"not null" json:"movie_id"`
	StudioID       uint           `gorm:"not null" json:"studio_id"`
	ShowTime       time.Time      `json:"show_time"`
	PricePerTicket float64        `json:"price_per_ticket"`
	Status         string         `gorm:"not null" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
