package userInterface

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
)

func HomePage(c *fiber.Ctx) error {
	var banners []models.Banner
	if err := models.DB.Order("created_at desc").Find(&banners).Error; err != nil {
		return c.Status(500).SendString("Failed to fetch banners")
	}

	return c.Render("index/home", fiber.Map{
		"Title":   "Home",
		"Banners": banners,
	}, "layouts/main")
}
