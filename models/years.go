package models

import "time"

type Year struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	Year      uint `json:"year"`
	CreatedAt time.Time
}
