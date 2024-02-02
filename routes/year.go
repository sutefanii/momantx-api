package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sixfwa/fiber-api/database"
	"github.com/sixfwa/fiber-api/models"
)

type YearSerializer struct {
	ID   uint `json:"id" gorm:"primaryKey"`
	Year uint `json:"year"`
}

func CreateResponseYear(year models.Year) YearSerializer {
	return YearSerializer{ID: year.ID, Year: year.Year}
}

func CreateYear(c *fiber.Ctx) error {
	var year models.Year
	if err := c.BodyParser(&year); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.Database.Db.Create(&year)
	responseYear := CreateResponseYear(year)

	return c.Status(200).JSON(responseYear)
}

func GetYears(c *fiber.Ctx) error {
	years := []models.Year{}

	database.Database.Db.Find(&years)
	responseYears := []YearSerializer{}
	for _, year := range years {
		responseYear := CreateResponseYear(year)
		responseYears = append(responseYears, responseYear)
	}
	return c.Status(200).JSON(responseYears)
}

func findYear(id int, year *models.Year) error {
	database.Database.Db.Find(&year, "id = ?", id)
	if year.ID == 0 {
		return errors.New("year dose not exist")
	}
	return nil
}

func GetYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var year models.Year
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findYear(id, &year); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseYear := CreateResponseYear(year)
	return c.Status(200).JSON(responseYear)
}

func UpdateYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var year models.Year
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findYear(id, &year); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateYear struct {
		Year uint `json:"year"`
	}

	var updateData UpdateYear
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	year.Year = updateData.Year
	database.Database.Db.Save(&year)

	responseYear := CreateResponseYear(year)
	return c.Status(200).JSON(responseYear)
}

func DeleteYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var year models.Year
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findYear(id, &year); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&year).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).SendString("Successfully Deleted Year")
}
