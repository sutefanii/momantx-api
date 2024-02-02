package routes

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sixfwa/fiber-api/database"
	"github.com/sixfwa/fiber-api/models"
	"gorm.io/datatypes"
)

type ItemSerializer struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Year       YearSerializer `gorm:"foreignKey:YaerRefer"`
	Date       string         `json:"date"`
	Name       string         `json:"name"`
	Text       string         `json:"text" gorm:"text"`
	SourceLink string         `json:"source_link"`
	ImageReal  datatypes.JSON `json:"imageReal"`
	ImageAi    datatypes.JSON `json:"imageAi"`
	Slug       string         `json:"slug"`
	CreatedAt  time.Time
}

func CreateResponseItem(item models.Item, year YearSerializer) ItemSerializer {
	return ItemSerializer{
		ID:         item.ID,
		CreatedAt:  item.CreatedAt,
		Year:       year,
		Date:       item.Date,
		Name:       item.Name,
		Text:       item.Text,
		SourceLink: item.SourceLink,
		ImageReal:  item.ImageReal,
		ImageAi:    item.ImageAi,
		Slug:       item.Slug,
	}
}

func CreateItem(c *fiber.Ctx) error {
	var item models.Item

	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var year models.Year
	if err := findYear(item.YaerRefer, &year); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err := database.Database.Db.Create(&item)
	if err != nil {
		return c.Status(500).SendString("Не корректный формат вводимых данных")
	}

	responseYear := CreateResponseYear(year)
	responseItem := CreateResponseItem(item, responseYear)

	return c.Status(200).JSON(responseItem)
}

func GetItems(c *fiber.Ctx) error {
	items := []models.Item{}

	database.Database.Db.Find(&items)
	responseItems := []ItemSerializer{}
	for _, item := range items {
		var year models.Year
		database.Database.Db.Find(&year, "id = ?", item.YaerRefer)

		responseItem := CreateResponseItem(item, CreateResponseYear(year))
		responseItems = append(responseItems, responseItem)
	}
	return c.Status(200).JSON(responseItems)
}

func findItem(id int, item *models.Item) error {
	database.Database.Db.Find(&item, "id = ?", id)
	if item.ID == 0 {
		return errors.New("order does not exist")
	}

	return nil
}

func GetItem(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var item models.Item

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findItem(id, &item); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var year models.Year
	database.Database.Db.First(&year, item.YaerRefer)
	responseYear := CreateResponseYear(year)

	responseItem := CreateResponseItem(item, responseYear)
	return c.Status(200).JSON(responseItem)
}

func UpdateItem(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var item models.Item
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findItem(id, &item); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateItem struct {
		YaerRefer  int            `json:"year_id"`
		Date       string         `json:"date"`
		Name       string         `json:"name"`
		Text       string         `json:"text" gorm:"text"`
		SourceLink string         `json:"source_link"`
		ImageReal  datatypes.JSON `json:"imageReal"`
		ImageAi    datatypes.JSON `json:"imageAi"`
	}

	var updateData UpdateItem
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	item.YaerRefer = updateData.YaerRefer
	item.Date = updateData.Date
	item.Name = updateData.Name
	item.Text = updateData.Text
	item.SourceLink = updateData.SourceLink
	item.ImageReal = updateData.ImageReal
	item.ImageAi = updateData.ImageAi

	database.Database.Db.Save(&item)

	var year models.Year
	database.Database.Db.Find(&year, "id = ?", updateData.YaerRefer)

	responseYear := CreateResponseYear(year)
	responseItem := CreateResponseItem(item, responseYear)

	return c.Status(200).JSON(responseItem)
}

func DeleteItem(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var item models.Item
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findItem(id, &item); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&item).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).SendString("Successfully Deleted Item")
}
