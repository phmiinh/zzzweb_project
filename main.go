package main

import (
	"Zzz_project/controllers"
	"Zzz_project/models"
	"Zzz_project/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"log"
)

func main() {
	models.ConnectDatabase()

	app := fiber.New(fiber.Config{
		Views: html.New("./views/shop", ".html"),
	})

	app.Static("/", "./public")
	controllers.Store = session.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
