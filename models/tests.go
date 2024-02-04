package models

import "time"

type Test struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Title     string `json:"title_test"`
	ImagePath string
	CreatedAt time.Time
}
