package models

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Question struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	TestRefer     int            `json:"test_id"`
	Test          Test           `gorm:"foreignKey:TestRefer"`
	CorrectAnswer string         `json:"correct_answer"`
	Answers       datatypes.JSON `json:"answers"`
	QuestionTitle string         `json:"title_question"`
	ItemRefer     int            `json:"item_id"`
	Item          Item           `gorm:"foreignKey:ItemRefer"`
}

func (q *Question) BeforeSave(tx *gorm.DB) (err error) {
	var (
		AnswersArr []string
	)
	err = json.Unmarshal(q.Answers, &AnswersArr)
	if err != nil {
		return err
	}
	return nil
}
