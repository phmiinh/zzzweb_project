package adminInterface

import (
	"Zzz_project/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func ListBanners(c *fiber.Ctx) error {
	var banners []models.Banner
	if err := models.DB.Order("created_at desc").Find(&banners).Error; err != nil {
		return c.Status(500).SendString("Failed to fetch banners")
	}

	return c.Render("admin/banners", fiber.Map{
		"Banners": banners,
	}, "layouts/admin")
}

func CreateBanner(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).SendString("Image upload failed")
	}

	filePath := fmt.Sprintf("./public/uploads/banners/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).SendString("Failed to save image")
	}

	banner := models.Banner{
		ImageURL: strings.TrimPrefix(filePath, "."),
		Title:    c.FormValue("title"),
		Subtitle: c.FormValue("subtitle"),
		Link:     c.FormValue("link"),
	}

	if err := models.DB.Create(&banner).Error; err != nil {
		return c.Status(500).SendString("Failed to create banner")
	}

	return c.Redirect("/admin/banners")
}

func UpdateBanner(c *fiber.Ctx) error {
	id := c.Params("id")
	var banner models.Banner

	if err := models.DB.First(&banner, id).Error; err != nil {
		return c.Status(404).SendString("Banner not found")
	}

	// Nếu có upload ảnh mới
	file, err := c.FormFile("image")
	if err == nil && file.Size > 0 {
		filePath := fmt.Sprintf("./public/uploads/banners/%s", file.Filename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(500).SendString("Failed to save new image")
		}
		banner.ImageURL = strings.TrimPrefix(filePath, ".")
	}

	banner.Title = c.FormValue("title")
	banner.Subtitle = c.FormValue("subtitle")
	banner.Link = c.FormValue("link")

	if err := models.DB.Save(&banner).Error; err != nil {
		return c.Status(500).SendString("Failed to update banner")
	}

	return c.Redirect("/admin/banners")
}

func DeleteBanner(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := models.DB.Delete(&models.Banner{}, id).Error; err != nil {
		return c.Status(500).SendString("Failed to delete banner")
	}

	return c.SendStatus(200)
}
