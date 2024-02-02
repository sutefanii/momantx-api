package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sixfwa/fiber-api/database"
	"github.com/sixfwa/fiber-api/models"
	"gorm.io/datatypes"
)

type QuestionSerializer struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Test          TestSerializer `gorm:"foreignKey:TestRefer"`
	CorrectAnswer string         `json:"correct_answer"`
	Answers       datatypes.JSON `json:"answers"`
	QuestionTitle string         `json:"title_question"`
	Item          ItemSerializer `gorm:"foreignKey:ItemRefer"`
}

func CreateResponseQuestion(question models.Question, test TestSerializer, item ItemSerializer) QuestionSerializer {
	return QuestionSerializer{
		ID:            question.ID,
		Test:          test,
		CorrectAnswer: question.CorrectAnswer,
		Answers:       question.Answers,
		QuestionTitle: question.QuestionTitle,
		Item:          item,
	}
}

func CreateQuestion(c *fiber.Ctx) error {
	var question models.Question

	if err := c.BodyParser(&question); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var test models.Test
	if err := findTest(question.TestRefer, &test); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var item models.Item
	if err := findItem(question.ItemRefer, &item); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var year models.Year
	if err := findYear(item.YaerRefer, &year); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&question)

	responseTest := CreateResponseTest(test)
	responseYear := CreateResponseYear(year)
	responseItem := CreateResponseItem(item, responseYear)

	responseQuestion := CreateResponseQuestion(question, responseTest, responseItem)
	return c.Status(200).JSON(responseQuestion)
}

func GetQuestions(c *fiber.Ctx) error {
	questions := []models.Question{}

	database.Database.Db.Find(&questions)
	responseQuestions := []QuestionSerializer{}
	for _, question := range questions {
		var year models.Year
		database.Database.Db.Find(&year, "id = ?", question.Item.YaerRefer)
		responseItem := CreateResponseItem(question.Item, CreateResponseYear(year))
		responseTest := CreateResponseTest(question.Test)

		responseQuestion := CreateResponseQuestion(question, responseTest, responseItem)

		responseQuestions = append(responseQuestions, responseQuestion)
	}
	return c.Status(200).JSON(responseQuestions)
}

func FindQuestion(id int, question *models.Question) error {
	database.Database.Db.Find(&question, "id = ?", id)
	if question.ID == 0 {
		return errors.New("order does not exist")
	}

	return nil
}

func GetQuestion(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var question models.Question

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := FindQuestion(id, &question); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var year models.Year
	database.Database.Db.Find(&year, "id = ?", question.Item.YaerRefer)
	responseItem := CreateResponseItem(question.Item, CreateResponseYear(year))
	responseTest := CreateResponseTest(question.Test)

	responseQuestion := CreateResponseQuestion(question, responseTest, responseItem)
	return c.Status(200).JSON(responseQuestion)
}

func UpdateQuestion(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var question models.Question
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := FindQuestion(id, &question); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateQuestionData struct {
		Test          models.Test    `gorm:"foreignKey:TestRefer"`
		CorrectAnswer string         `json:"correct_answer"`
		Answers       datatypes.JSON `json:"answers"`
		QuestionTitle string         `json:"title_question"`
		Item          models.Item    `gorm:"foreignKey:ItemRefer"`
	}

	var updateData UpdateQuestionData
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	question.Test = updateData.Test
	question.CorrectAnswer = updateData.CorrectAnswer
	question.Answers = updateData.Answers
	question.QuestionTitle = updateData.QuestionTitle
	question.Item = updateData.Item

	database.Database.Db.Save(&question)

	var year models.Year
	database.Database.Db.Find(&year, "id = ?", question.Item.YaerRefer)

	responseYear := CreateResponseYear(year)

	responseTest := CreateResponseTest(updateData.Test)
	responseItem := CreateResponseItem(updateData.Item, responseYear)

	responseQuestion := CreateResponseQuestion(question, responseTest, responseItem)
	return c.Status(200).JSON(responseQuestion)
}

func DeleteQuestion(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var question models.Question
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := FindQuestion(id, &question); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&question).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).SendString("Successfully Deleted Question")
}
