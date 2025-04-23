package controllers

import (
	"Zzz_project/models"
	"github.com/gofiber/fiber/v2"
)

// Trang dashboard chính của Admin
func AdminDashboard(c *fiber.Ctx) error {
	// Kiểm tra session để lấy thông tin user
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	role, _ := sess.Get("role").(string)
	if role != "admin" {
		return c.Status(403).SendString("Bạn không có quyền truy cập!")
	}

	// Truy vấn số lượng sản phẩm và đơn hàng
	var productCount int64
	var orderCount int64

	models.DB.Model(&models.Product{}).Count(&productCount)
	models.DB.Model(&models.Order{}).Count(&orderCount)

	// Render giao diện admin dashboard
	return c.Render("admin_dashboard", fiber.Map{
		"Title":        "Admin Dashboard",
		"ProductCount": productCount,
		"OrderCount":   orderCount,
		"Username":     sess.Get("username"),
	}, "layouts/admin")
}

func AdminUsersHandler(c *fiber.Ctx) error {

	var users []models.User
	if err := models.DB.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Lỗi truy vấn database")
	}

	return c.Render("admin/admin_users", fiber.Map{
		"Title": "Users management",
		"Users": users,
	}, "layouts/admin")
}
