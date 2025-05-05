package adminInterface

import (
	"Zzz_project/controllers/loginSignup"
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
)

func AdminDashboard(c *fiber.Ctx) error {
	sess, err := loginSignup.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	role, _ := sess.Get("role").(string)
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).SendString("Bạn không có quyền truy cập!")
	}

	var userCount int64
	var productCount int64
	var categoryCount int64
	var orderCount int64
	var bannerCount int64

	models.DB.Model(&models.User{}).Where("role = ?", "customer").Count(&userCount)
	models.DB.Model(&models.Product{}).Count(&productCount)
	models.DB.Raw("SELECT COUNT(DISTINCT category) FROM products").Scan(&categoryCount)
	models.DB.Model(&models.Order{}).Count(&orderCount)
	models.DB.Model(&models.Banner{}).Count(&bannerCount) // ✨ Thêm count Banner đây

	return c.Render("admin/dashboard", fiber.Map{
		"Title":         "Admin Dashboard",
		"Username":      sess.Get("username"),
		"UserCount":     userCount,
		"ProductCount":  productCount,
		"CategoryCount": categoryCount,
		"OrderCount":    orderCount,
		"BannerCount":   bannerCount,
	}, "layouts/admin")
}
