package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sixfwa/fiber-api/database"
	"github.com/sixfwa/fiber-api/models"
)

type TestSerializer struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title_test"`
}

func CreateResponseTest(test models.Test) TestSerializer {
	return TestSerializer{ID: test.ID, Title: test.Title}
}

func CreateTest(c *fiber.Ctx) error {
	var test models.Test
	if err := c.BodyParser(&test); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.Database.Db.Create(&test)
	responseTest := CreateResponseTest(test)

	return c.Status(200).JSON(responseTest)
}

func GetTests(c *fiber.Ctx) error {
	tests := []models.Test{}

	database.Database.Db.Find(&tests)
	responseTests := []TestSerializer{}
	for _, test := range tests {
		responseTest := CreateResponseTest(test)
		responseTests = append(responseTests, responseTest)
	}
	return c.Status(200).JSON(responseTests)
}

func findTest(id int, test *models.Test) error {
	database.Database.Db.Find(&test, "id = ?", id)
	if test.ID == 0 {
		return errors.New("test dose not exist")
	}
	return nil
}

func GetTest(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var test models.Test
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findTest(id, &test); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseTest := CreateResponseTest(test)
	return c.Status(200).JSON(responseTest)
}

func UpdateTest(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var test models.Test
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findTest(id, &test); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateTest struct {
		Title string `json:"title_test"`
	}

	var updateData UpdateTest
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	test.Title = updateData.Title
	database.Database.Db.Save(&test)

	responseTest := CreateResponseTest(test)
	return c.Status(200).JSON(responseTest)
}

func DeleteTest(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var test models.Test
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findTest(id, &test); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&test).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).SendString("Successfully Deleted Test")
}
