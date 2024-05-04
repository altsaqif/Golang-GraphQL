package models

import "time"

type Product struct {
	ID        int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string  `gorm:"not null" json:"name"`
	Stock     float64 `gorm:"not null" json:"stock"`
	Price     float64 `gorm:"not null" json:"price"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
