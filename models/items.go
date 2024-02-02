package models

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"encoding/json"
	"fmt"
	"strings"
)

type Item struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	YaerRefer  int            `json:"year_id"`
	Year       Year           `gorm:"foreignKey:YaerRefer"`
	Date       string         `json:"date"`
	Name       string         `json:"name"`
	Text       string         `json:"text" gorm:"text"`
	SourceLink string         `json:"source_link"`
	ImageReal  datatypes.JSON `json:"imageReal"`
	ImageAi    datatypes.JSON `json:"imageAi"`
	Slug       string         `json:"slug"`
	CreatedAt  time.Time
}

func (item *Item) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.New()
	item.Slug = uuid.String()
	return
}

func (item *Item) BeforeSave(tx *gorm.DB) (err error) {
	var (
		arrReal []string
		arrAi   []string
	)

	err = json.Unmarshal(item.ImageReal, &arrReal)
	if err != nil {
		return err
	}

	err = json.Unmarshal(item.ImageAi, &arrAi)
	if err != nil {
		return err
	}

	for _, item := range append(arrReal, arrAi...) {
		if !strings.HasPrefix(item, "http://") && !strings.HasPrefix(item, "https://") {
			return fmt.Errorf("элемент '%s' не начинается с http:// или https://", item)
		}
	}

	return nil
}
