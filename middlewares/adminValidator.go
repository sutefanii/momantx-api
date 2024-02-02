package middlewares

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var COOKIES_FIELDS = []string{"auth"}

func CheackAdminIsValid(c *fiber.Ctx) error {
	for i := 0; i < len(COOKIES_FIELDS); i++ {
		cookie := c.Cookies(COOKIES_FIELDS[i])

		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error dotenv file \n" + err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		if cookie != os.Getenv("ADMIN_TOKEN") {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
	}
	return c.Next()
}
