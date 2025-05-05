package main

import (
	"Zzz_project/controllers/loginSignup"
	"Zzz_project/models"
	"Zzz_project/routes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"log"
	"time"
)

func main() {
	models.ConnectDatabase()

	engine := html.New("./views/shop", ".html")

	engine.AddFunc("formatFloat", func(f float64) string {
		return fmt.Sprintf("%.2f", f)
	})
	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})
	engine.AddFunc("formatDate", func(t time.Time) string {
		return t.Format("2006-01-02")
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./public")
	loginSignup.Store = session.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
