package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sixfwa/fiber-api/database"
	"github.com/sixfwa/fiber-api/middlewares"
	"github.com/sixfwa/fiber-api/routes"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to my awesome API")
}

func setupRoutes(app *fiber.App) {
	// welcome
	app.Get("/api", welcome)
	// Year endpoints
	app.Get("/api/years/:id", routes.GetYear)
	app.Get("/api/years", routes.GetYears)

	app.Post("/api/years", middlewares.CheackAdminIsValid, routes.CreateYear)
	app.Put("/api/years/:id", middlewares.CheackAdminIsValid, routes.UpdateYear)
	app.Delete("/api/years/:id", middlewares.CheackAdminIsValid, routes.DeleteYear)
	// Item endpoints
	app.Get("/api/items", routes.GetItems)
	app.Get("/api/items/:id", routes.GetItem)

	app.Post("/api/items", middlewares.CheackAdminIsValid, routes.CreateItem)
	app.Put("/api/items/:id", middlewares.CheackAdminIsValid, routes.UpdateItem)
	app.Delete("/api/items/:id", middlewares.CheackAdminIsValid, routes.DeleteItem)
	// Test endpoints
	app.Get("/api/tests/:id", routes.GetTest)
	app.Get("/api/tests", routes.GetTests)

	app.Post("/api/tests", middlewares.CheackAdminIsValid, routes.CreateTest)
	app.Put("/api/tests/:id", middlewares.CheackAdminIsValid, routes.UpdateTest)
	app.Delete("/api/tests/:id", middlewares.CheackAdminIsValid, routes.DeleteTest)
	// Questions endpoints
	app.Get("/api/questions/:id", routes.GetQuestion)
	app.Get("/api/questions", routes.GetQuestions)

	app.Post("/api/questions", middlewares.CheackAdminIsValid, routes.CreateQuestion)
	app.Put("/api/questions/:id", middlewares.CheackAdminIsValid, routes.UpdateQuestion)
	app.Delete("/api/questions/:id", middlewares.CheackAdminIsValid, routes.DeleteQuestion)
}

func main() {
	database.ConnectDb()
	app := fiber.New()

	setupRoutes(app)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error dotenv file \n" + err.Error())
		return
	}

	log.Fatal(app.Listen(os.Getenv("PORT")))
}
